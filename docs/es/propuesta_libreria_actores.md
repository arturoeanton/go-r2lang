# Propuesta: Librería de Actores para R2Lang

**Versión:** 1.0  
**Fecha:** 2025-07-15  
**Estado:** Propuesta  

## Resumen Ejecutivo

Esta propuesta presenta un sistema de actores (Actor Model) para R2Lang que permite manejar concurrencia de manera más estructurada y escalable, proporcionando comunicación asíncrona mediante mensajes, aislamiento de estado y tolerancia a fallos.

## Problema Actual

R2Lang actualmente maneja concurrencia con:
- **r2() goroutines básicas:** Sin estructura ni comunicación formal
- **Estado compartido:** Posibles race conditions
- **Falta de supervisión:** No hay mecanismo de recuperación de fallos
- **Comunicación primitiva:** No hay sistema de mensajería estructurado

```r2
// Concurrencia actual - limitada
r2(func() {
    // Goroutine básica sin estructura
    print("Ejecutando en paralelo");
});
```

## Solución Propuesta: Sistema de Actores

### 1. Modelo de Actores Básico

#### 1.1 Definición de Actores
```r2
// actor_system.r2
import "actors" as actors;

// Definir un actor
class CounterActor extends actors.Actor {
    let count = 0;
    
    receive(message) {
        switch (message.type) {
            case "INCREMENT":
                this.count += message.value || 1;
                this.sender().tell({
                    type: "COUNTER_RESPONSE",
                    value: this.count
                });
                break;
                
            case "GET_COUNT":
                this.sender().tell({
                    type: "COUNT_VALUE",
                    value: this.count
                });
                break;
                
            case "RESET":
                this.count = 0;
                this.sender().tell({type: "RESET_OK"});
                break;
                
            default:
                this.unhandled(message);
        }
    }
}

// Crear sistema de actores
let system = actors.createSystem("MySystem");

// Crear y usar actores
let counter = system.actorOf(CounterActor, "counter");

counter.tell({type: "INCREMENT", value: 5});
counter.ask({type: "GET_COUNT"}).then(func(response) {
    print("Count:", response.value);
});
```

#### 1.2 Actor System Architecture
```r2
// Diferentes tipos de actores
class WorkerActor extends actors.Actor {
    receive(message) {
        switch (message.type) {
            case "WORK":
                let result = this.processWork(message.data);
                this.sender().tell({
                    type: "WORK_RESULT",
                    result: result
                });
                break;
        }
    }
    
    processWork(data) {
        // Simular trabajo pesado
        sleep(1000);
        return data * 2;
    }
}

class SupervisorActor extends actors.Actor {
    let workers = [];
    
    preStart() {
        // Crear workers al iniciar
        for (let i = 0; i < 5; i++) {
            let worker = this.context().actorOf(WorkerActor, `worker-${i}`);
            this.workers.push(worker);
        }
    }
    
    receive(message) {
        switch (message.type) {
            case "DISTRIBUTE_WORK":
                this.distributeWork(message.tasks);
                break;
                
            case "WORKER_FAILED":
                this.handleWorkerFailure(message.worker);
                break;
        }
    }
    
    distributeWork(tasks) {
        let workerIndex = 0;
        for (let task of tasks) {
            let worker = this.workers[workerIndex];
            worker.tell({type: "WORK", data: task});
            workerIndex = (workerIndex + 1) % this.workers.length;
        }
    }
}
```

### 2. Implementación Técnica del Actor System

#### 2.1 Core Actor Infrastructure
```go
// pkg/r2libs/r2actors.go
type ActorSystem struct {
    Name       string
    Actors     map[string]*ActorRef
    Mailboxes  map[string]*Mailbox
    Scheduler  *Scheduler
    Supervisor *Supervisor
    Config     *ActorConfig
}

type ActorRef struct {
    Path     string
    System   *ActorSystem
    Mailbox  *Mailbox
    Actor    Actor
    Context  *ActorContext
    State    ActorState
}

type Actor interface {
    Receive(message *Message) error
    PreStart() error
    PostStop() error
    PreRestart(reason error) error
    PostRestart(reason error) error
}

type Message struct {
    Type     string
    Data     interface{}
    Sender   *ActorRef
    ID       string
    Created  time.Time
}

type Mailbox struct {
    Messages chan *Message
    Capacity int
    Actor    *ActorRef
}

func NewActorSystem(name string) *ActorSystem {
    return &ActorSystem{
        Name:      name,
        Actors:    make(map[string]*ActorRef),
        Mailboxes: make(map[string]*Mailbox),
        Scheduler: NewScheduler(),
        Supervisor: NewSupervisor(),
        Config:    DefaultActorConfig(),
    }
}

func (as *ActorSystem) ActorOf(actorType string, name string) *ActorRef {
    path := fmt.Sprintf("/%s/%s", as.Name, name)
    
    mailbox := &Mailbox{
        Messages: make(chan *Message, as.Config.MailboxCapacity),
        Capacity: as.Config.MailboxCapacity,
    }
    
    ref := &ActorRef{
        Path:    path,
        System:  as,
        Mailbox: mailbox,
        Context: NewActorContext(as, path),
        State:   ActorStateStarting,
    }
    
    // Create actor instance
    actor := as.createActorInstance(actorType)
    ref.Actor = actor
    mailbox.Actor = ref
    
    as.Actors[path] = ref
    as.Mailboxes[path] = mailbox
    
    // Start actor
    go ref.start()
    
    return ref
}
```

#### 2.2 Message Processing
```go
func (ar *ActorRef) start() {
    // Ejecutar PreStart
    if err := ar.Actor.PreStart(); err != nil {
        ar.System.Supervisor.HandleFailure(ar, err)
        return
    }
    
    ar.State = ActorStateRunning
    
    // Loop principal de procesamiento de mensajes
    for {
        select {
        case message := <-ar.Mailbox.Messages:
            ar.processMessage(message)
            
        case <-ar.Context.StopSignal:
            ar.stop()
            return
        }
    }
}

func (ar *ActorRef) processMessage(message *Message) {
    defer func() {
        if r := recover(); r != nil {
            err := fmt.Errorf("actor panic: %v", r)
            ar.System.Supervisor.HandleFailure(ar, err)
        }
    }()
    
    // Set current sender in context
    ar.Context.CurrentSender = message.Sender
    
    // Process message
    if err := ar.Actor.Receive(message); err != nil {
        ar.System.Supervisor.HandleFailure(ar, err)
    }
}

func (ar *ActorRef) Tell(message *Message) {
    if ar.State != ActorStateRunning {
        return // Actor not ready
    }
    
    select {
    case ar.Mailbox.Messages <- message:
        // Message sent successfully
    default:
        // Mailbox full, handle overflow
        ar.System.Supervisor.HandleMailboxOverflow(ar, message)
    }
}

func (ar *ActorRef) Ask(message *Message, timeout time.Duration) (*Message, error) {
    // Create response channel
    responseID := generateMessageID()
    responseChan := make(chan *Message, 1)
    
    // Register response handler
    ar.System.RegisterResponseHandler(responseID, responseChan)
    defer ar.System.UnregisterResponseHandler(responseID)
    
    // Send message with response ID
    message.ID = responseID
    ar.Tell(message)
    
    // Wait for response
    select {
    case response := <-responseChan:
        return response, nil
    case <-time.After(timeout):
        return nil, fmt.Errorf("ask timeout after %v", timeout)
    }
}
```

### 3. Supervisión y Tolerancia a Fallos

#### 3.1 Supervisor Strategy
```go
// pkg/r2libs/r2supervisor.go
type SupervisorStrategy int

const (
    RestartStrategy SupervisorStrategy = iota
    ResumeStrategy
    StopStrategy
    EscalateStrategy
)

type Supervisor struct {
    Strategy     SupervisorStrategy
    MaxRetries   int
    TimeWindow   time.Duration
    FailureStats map[string]*FailureStats
}

type FailureStats struct {
    Failures    int
    LastFailure time.Time
    FirstFailure time.Time
}

func (s *Supervisor) HandleFailure(actor *ActorRef, err error) {
    stats := s.getFailureStats(actor.Path)
    stats.Failures++
    stats.LastFailure = time.Now()
    
    if stats.FirstFailure.IsZero() {
        stats.FirstFailure = time.Now()
    }
    
    // Check if within time window
    if time.Since(stats.FirstFailure) > s.TimeWindow {
        // Reset failure count
        stats.Failures = 1
        stats.FirstFailure = time.Now()
    }
    
    // Decide action based on strategy
    switch s.Strategy {
    case RestartStrategy:
        if stats.Failures <= s.MaxRetries {
            s.restartActor(actor, err)
        } else {
            s.stopActor(actor, err)
        }
        
    case ResumeStrategy:
        s.resumeActor(actor, err)
        
    case StopStrategy:
        s.stopActor(actor, err)
        
    case EscalateStrategy:
        s.escalateFailure(actor, err)
    }
}

func (s *Supervisor) restartActor(actor *ActorRef, err error) {
    // PreRestart hook
    if restartErr := actor.Actor.PreRestart(err); restartErr != nil {
        s.stopActor(actor, restartErr)
        return
    }
    
    // Clear mailbox
    s.clearMailbox(actor)
    
    // Restart actor
    actor.State = ActorStateRestarting
    go func() {
        // PostRestart hook
        if restartErr := actor.Actor.PostRestart(err); restartErr != nil {
            s.stopActor(actor, restartErr)
            return
        }
        
        actor.State = ActorStateRunning
        actor.start()
    }()
}
```

#### 3.2 Fault Tolerance en R2Lang
```r2
// Configuración de supervisor
class ResilientActor extends actors.Actor {
    supervisorStrategy() {
        return {
            strategy: "restart",
            maxRetries: 3,
            timeWindow: "1m",
            decider: func(error) {
                switch (error.type) {
                    case "NetworkError":
                        return "restart";
                    case "ValidationError":
                        return "resume";
                    case "CriticalError":
                        return "stop";
                    default:
                        return "escalate";
                }
            }
        };
    }
    
    receive(message) {
        switch (message.type) {
            case "RISKY_OPERATION":
                this.performRiskyOperation(message.data);
                break;
        }
    }
    
    performRiskyOperation(data) {
        if (Math.random() < 0.3) {
            throw new Error("NetworkError", "Connection failed");
        }
        
        // Operación normal
        return data.process();
    }
    
    preRestart(reason) {
        print(`Actor restarting due to: ${reason}`);
        // Cleanup logic
    }
    
    postRestart(reason) {
        print(`Actor restarted successfully`);
        // Reinitialization logic
    }
}
```

### 4. Patrones de Comunicación

#### 4.1 Fire-and-Forget (Tell)
```r2
// Comunicación asíncrona sin respuesta
let logger = system.actorOf(LoggerActor, "logger");
logger.tell({
    type: "LOG",
    level: "INFO",
    message: "Operation completed"
});
```

#### 4.2 Request-Response (Ask)
```r2
// Comunicación con respuesta
let calculator = system.actorOf(CalculatorActor, "calc");

calculator.ask({
    type: "CALCULATE",
    operation: "add",
    values: [10, 20]
}, "5s").then(func(response) {
    print("Result:", response.result);
}).catch(func(error) {
    print("Error:", error.message);
});
```

#### 4.3 Publish-Subscribe
```r2
// Sistema de eventos
class EventBusActor extends actors.Actor {
    let subscribers = {};
    
    receive(message) {
        switch (message.type) {
            case "SUBSCRIBE":
                this.subscribe(message.topic, message.subscriber);
                break;
                
            case "PUBLISH":
                this.publish(message.topic, message.data);
                break;
                
            case "UNSUBSCRIBE":
                this.unsubscribe(message.topic, message.subscriber);
                break;
        }
    }
    
    subscribe(topic, subscriber) {
        if (!this.subscribers[topic]) {
            this.subscribers[topic] = [];
        }
        this.subscribers[topic].push(subscriber);
    }
    
    publish(topic, data) {
        let topicSubscribers = this.subscribers[topic] || [];
        for (let subscriber of topicSubscribers) {
            subscriber.tell({
                type: "EVENT",
                topic: topic,
                data: data
            });
        }
    }
}

// Uso del EventBus
let eventBus = system.actorOf(EventBusActor, "eventBus");
let subscriber = system.actorOf(NotificationActor, "notifications");

eventBus.tell({
    type: "SUBSCRIBE",
    topic: "user.registered",
    subscriber: subscriber
});

// Publicar evento
eventBus.tell({
    type: "PUBLISH",
    topic: "user.registered",
    data: {userId: 123, email: "user@example.com"}
});
```

### 5. Routing y Load Balancing

#### 5.1 Router Actors
```r2
// Router para balanceo de carga
class RoundRobinRouter extends actors.Actor {
    let routees = [];
    let currentIndex = 0;
    
    constructor(routeeClass, count) {
        super();
        
        // Crear routees
        for (let i = 0; i < count; i++) {
            let routee = this.context().actorOf(routeeClass, `routee-${i}`);
            this.routees.push(routee);
        }
    }
    
    receive(message) {
        if (message.type === "ROUTE") {
            let routee = this.getNextRoutee();
            routee.tell(message.data);
        }
    }
    
    getNextRoutee() {
        let routee = this.routees[this.currentIndex];
        this.currentIndex = (this.currentIndex + 1) % this.routees.length;
        return routee;
    }
}

// Uso del router
let router = system.actorOf(RoundRobinRouter, "worker-router", WorkerActor, 5);

// Distribuir trabajo
for (let i = 0; i < 100; i++) {
    router.tell({
        type: "ROUTE",
        data: {type: "WORK", task: i}
    });
}
```

#### 5.2 Diferentes Estrategias de Routing
```r2
// Estrategias de routing disponibles
class RouterStrategies {
    static roundRobin(routees, message) {
        // Implementación round-robin
    }
    
    static random(routees, message) {
        let index = Math.floor(Math.random() * routees.length);
        return routees[index];
    }
    
    static consistentHashing(routees, message) {
        // Hash basado en algún field del mensaje
        let hash = this.hashFunction(message.key);
        let index = hash % routees.length;
        return routees[index];
    }
    
    static smallestMailbox(routees, message) {
        // Seleccionar actor con menor carga
        let smallest = routees[0];
        for (let routee of routees) {
            if (routee.mailboxSize() < smallest.mailboxSize()) {
                smallest = routee;
            }
        }
        return smallest;
    }
}
```

### 6. Integración con R2Lang

#### 6.1 Built-in Actor Functions
```go
// pkg/r2libs/r2actors.go
func RegisterActors(env *r2core.Environment) {
    // Sistema de actores
    env.Set("actors", map[string]interface{}{
        "createSystem":    createSystem,
        "Actor":          BaseActor{},
        "Router":         BaseRouter{},
        "Supervisor":     BaseSupervisor{},
    })
    
    // Utilidades
    env.Set("actorSystem", nil) // Se seteará cuando se cree un sistema
}

func createSystem(args ...interface{}) interface{} {
    if len(args) == 0 {
        panic("createSystem requires a name")
    }
    
    name := args[0].(string)
    system := NewActorSystem(name)
    
    // Registrar sistema global
    env.Set("actorSystem", system)
    
    return system
}
```

#### 6.2 Actor Lifecycle Hooks
```r2
class MyActor extends actors.Actor {
    preStart() {
        // Inicialización
        print("Actor starting:", this.self().path());
        this.initializeResources();
    }
    
    postStop() {
        // Limpieza
        print("Actor stopping:", this.self().path());
        this.cleanupResources();
    }
    
    preRestart(reason) {
        print("Actor restarting due to:", reason);
        this.saveState();
    }
    
    postRestart(reason) {
        print("Actor restarted");
        this.restoreState();
    }
    
    receive(message) {
        // Lógica principal
    }
}
```

### 7. Clustering y Distribución

#### 7.1 Remote Actors
```r2
// Configuración de cluster
let clusterConfig = {
    seedNodes: ["127.0.0.1:2551", "127.0.0.1:2552"],
    port: 2551,
    hostname: "127.0.0.1",
    actorSystem: "ClusterSystem"
};

let system = actors.createClusterSystem("MyCluster", clusterConfig);

// Crear actor remoto
let remoteActor = system.actorOf(WorkerActor, "worker", {
    deploy: {
        target: "remote://MyCluster@127.0.0.1:2552"
    }
});

// Uso normal - transparente si es local o remoto
remoteActor.tell({type: "WORK", data: "some task"});
```

#### 7.2 Cluster Management
```go
// pkg/r2libs/r2cluster.go
type ClusterSystem struct {
    *ActorSystem
    Nodes        map[string]*ClusterNode
    SeedNodes    []string
    LocalNode    *ClusterNode
    MemberStatus MemberStatus
}

type ClusterNode struct {
    Address    string
    Port       int
    Status     NodeStatus
    LastSeen   time.Time
    ActorRefs  map[string]*RemoteActorRef
}

func (cs *ClusterSystem) Join(seedNodes []string) error {
    // Conectar a seed nodes
    for _, seed := range seedNodes {
        if err := cs.connectToSeed(seed); err != nil {
            continue // Try next seed
        }
        return nil
    }
    return fmt.Errorf("failed to join cluster")
}

func (cs *ClusterSystem) DeployActor(actorType, name string, target string) *ActorRef {
    if cs.isLocalDeployment(target) {
        return cs.ActorSystem.ActorOf(actorType, name)
    }
    
    return cs.deployRemoteActor(actorType, name, target)
}
```

### 8. Monitoring y Metrics

#### 8.1 Actor Metrics
```r2
// Métricas automáticas
class MonitoredActor extends actors.Actor {
    receive(message) {
        // Métricas automáticas:
        // - Número de mensajes procesados
        // - Tiempo de procesamiento
        // - Errores
        // - Tamaño de mailbox
        
        this.context().metrics().recordMessage(message.type);
        
        let startTime = Date.now();
        try {
            this.processMessage(message);
            this.context().metrics().recordProcessingTime(
                message.type, 
                Date.now() - startTime
            );
        } catch (error) {
            this.context().metrics().recordError(message.type, error);
            throw error;
        }
    }
}

// Acceso a métricas
let metrics = system.metrics();
print("Messages processed:", metrics.messagesProcessed);
print("Average processing time:", metrics.averageProcessingTime);
print("Error rate:", metrics.errorRate);
```

#### 8.2 Health Checks
```r2
// Health check endpoints
class HealthCheckActor extends actors.Actor {
    receive(message) {
        switch (message.type) {
            case "HEALTH_CHECK":
                this.sender().tell({
                    type: "HEALTH_RESPONSE",
                    status: "healthy",
                    timestamp: Date.now(),
                    metrics: this.gatherMetrics()
                });
                break;
        }
    }
    
    gatherMetrics() {
        return {
            actorCount: this.context().system().getActorCount(),
            messageRate: this.context().system().getMessageRate(),
            errorRate: this.context().system().getErrorRate(),
            uptime: this.context().system().getUptime()
        };
    }
}
```

### 9. Ejemplos Prácticos

#### 9.1 Web Server con Actores
```r2
// Web server actor-based
class HttpHandlerActor extends actors.Actor {
    receive(message) {
        switch (message.type) {
            case "HTTP_REQUEST":
                this.handleRequest(message.request, message.response);
                break;
        }
    }
    
    handleRequest(req, res) {
        // Delegar a worker específico
        let worker = this.context().actorSelection("/system/api-workers/*");
        worker.tell({
            type: "PROCESS_REQUEST",
            request: req,
            response: res
        });
    }
}

class ApiWorkerActor extends actors.Actor {
    receive(message) {
        switch (message.type) {
            case "PROCESS_REQUEST":
                this.processApiRequest(message.request, message.response);
                break;
        }
    }
    
    processApiRequest(req, res) {
        // Procesar request
        let result = this.handleApiCall(req);
        res.json(result);
    }
}

// Setup del sistema
let system = actors.createSystem("WebServer");
let httpHandler = system.actorOf(HttpHandlerActor, "http-handler");
let workerRouter = system.actorOf(RoundRobinRouter, "api-workers", ApiWorkerActor, 10);

// Integrar con servidor HTTP
http.createServer(func(req, res) {
    httpHandler.tell({
        type: "HTTP_REQUEST",
        request: req,
        response: res
    });
}).listen(8080);
```

#### 9.2 Sistema de Chat
```r2
// Sistema de chat con actores
class ChatRoomActor extends actors.Actor {
    let participants = {};
    let messageHistory = [];
    
    receive(message) {
        switch (message.type) {
            case "JOIN":
                this.handleJoin(message.user, message.connection);
                break;
                
            case "LEAVE":
                this.handleLeave(message.user);
                break;
                
            case "MESSAGE":
                this.handleMessage(message.user, message.text);
                break;
        }
    }
    
    handleJoin(user, connection) {
        this.participants[user] = connection;
        this.broadcast({
            type: "USER_JOINED",
            user: user,
            timestamp: Date.now()
        });
        
        // Enviar historial al nuevo usuario
        connection.tell({
            type: "MESSAGE_HISTORY",
            messages: this.messageHistory
        });
    }
    
    handleMessage(user, text) {
        let message = {
            type: "CHAT_MESSAGE",
            user: user,
            text: text,
            timestamp: Date.now()
        };
        
        this.messageHistory.push(message);
        this.broadcast(message);
    }
    
    broadcast(message) {
        for (let user in this.participants) {
            this.participants[user].tell(message);
        }
    }
}

// Usar el sistema de chat
let chatSystem = actors.createSystem("ChatSystem");
let chatRoom = chatSystem.actorOf(ChatRoomActor, "general-chat");

// Simular conexión de usuario
chatRoom.tell({
    type: "JOIN",
    user: "alice",
    connection: userConnection
});
```

## Plan de Implementación

### Fase 1: Core Actor System
- [ ] Implementar ActorSystem básico
- [ ] Message passing (Tell/Ask)
- [ ] Actor lifecycle management
- [ ] Mailbox implementation

### Fase 2: Supervision y Fault Tolerance
- [ ] Supervisor strategies
- [ ] Actor restart mechanisms
- [ ] Failure escalation
- [ ] Health monitoring

### Fase 3: Routing y Load Balancing
- [ ] Router actors
- [ ] Diferentes estrategias de routing
- [ ] Dynamic routee management
- [ ] Load balancing metrics

### Fase 4: Advanced Features
- [ ] Publish-Subscribe pattern
- [ ] Actor selection/discovery
- [ ] Metrics y monitoring
- [ ] Performance optimizations

### Fase 5: Distribution y Clustering
- [ ] Remote actor deployment
- [ ] Cluster management
- [ ] Network communication
- [ ] Distributed supervision

## Beneficios

1. **Concurrencia Estructurada:** Modelo claro y predecible
2. **Tolerancia a Fallos:** Supervisión automática y recuperación
3. **Escalabilidad:** Distribución transparente
4. **Aislamiento:** Sin estado compartido, sin race conditions
5. **Maintainability:** Código más organizado y testeable

## Conclusión

El sistema de actores proporcionará a R2Lang un modelo de concurrencia moderno y escalable, permitiendo desarrollar aplicaciones distribuidas robustas con patrones de comunicación bien definidos y tolerancia a fallos integrada.