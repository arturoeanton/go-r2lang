# Propuesta: Palabra Reservada "sync" para R2Lang

## Resumen Ejecutivo

Esta propuesta introduce la palabra reservada `sync` en R2Lang, inspirada en la palabra clave `synchronized` de Java, para proporcionar sincronización thread-safe en operaciones concurrentes utilizando el ecosistema de goroutines existente de R2Lang.

## Motivación

R2Lang ya cuenta con primitivas de concurrencia a través de la función `r2()` que crea goroutines. Sin embargo, carece de mecanismos nativos para sincronización segura entre threads cuando múltiples goroutines acceden a recursos compartidos. La adición de `sync` completaría el modelo de concurrencia del lenguaje.

## Casos de Uso Principales

### 1. Protección de Variables Compartidas
```r2
let contador = 0;
let mutex = createMutex();

func incrementar() {
    sync(mutex) {
        contador = contador + 1;
        print("Contador: " + contador);
    }
}

// Múltiples goroutines pueden llamar incrementar() de forma segura
r2(incrementar);
r2(incrementar);
r2(incrementar);
```

### 2. Sincronización de Métodos de Clase
```r2
class ContadorSeguro {
    constructor() {
        this.valor = 0;
        this.mutex = createMutex();
    }
    
    sync incrementar() {  // Método sincronizado
        this.valor = this.valor + 1;
        return this.valor;
    }
    
    sync decrementar() {
        this.valor = this.valor - 1;
        return this.valor;
    }
}
```

### 3. Sincronización de Bloques Críticos
```r2
let recursoCompartido = createSharedResource();
let accessMutex = createMutex();

func procesarDatos(datos) {
    // Procesamiento local sin sincronización
    let resultado = transformar(datos);
    
    // Sección crítica sincronizada
    sync(accessMutex) {
        recursoCompartido.agregar(resultado);
        recursoCompartido.actualizar();
        print("Procesado: " + resultado.id);
    }
    
    // Más procesamiento local
    notificarCompletado(resultado);
}
```

## Especificación Técnica

### Sintaxis Propuesta

#### 1. Bloque Sync con Mutex Explícito
```r2
sync(mutexObject) {
    // código sincronizado
}
```

#### 2. Método Sync (Mutex Implícito por Instancia)
```r2
class MiClase {
    sync metodoSincronizado() {
        // código sincronizado usando mutex interno de la instancia
    }
}
```

#### 3. Función Sync Global
```r2
sync func funcionGlobal() {
    // código sincronizado usando mutex global del módulo
}
```

### Implementación en el Parser

#### Modificaciones al Lexer (`pkg/r2core/lexer.go`)
```go
// Agregar nuevo token
case "sync":
    return Token{Type: SYNC, Literal: "sync"}
```

#### Modificaciones al Parser (`pkg/r2core/parse.go`)
```go
// Nuevo nodo AST para bloques sync
type SyncBlock struct {
    Mutex      Node
    Body       *BlockStatement
    IsMethod   bool
    IsFunction bool
}

func (p *Parser) parseSyncStatement() Node {
    // Implementación del parsing de sync
}
```

#### Nuevo archivo AST (`pkg/r2core/sync_statement.go`)
```go
package r2core

import "sync"

type SyncStatement struct {
    Token    Token
    Mutex    Node
    Body     *BlockStatement
    IsMethod bool
}

func (ss *SyncStatement) Eval(env *Environment) interface{} {
    var mutex *sync.Mutex
    
    if ss.Mutex != nil {
        // Mutex explícito
        mutexVal := ss.Mutex.Eval(env)
        if m, ok := mutexVal.(*sync.Mutex); ok {
            mutex = m
        } else {
            panic("sync: argument must be a mutex")
        }
    } else {
        // Mutex implícito (método o función sync)
        mutex = env.GetImplicitMutex()
    }
    
    mutex.Lock()
    defer mutex.Unlock()
    
    return ss.Body.Eval(env)
}
```

### Nuevas Funciones Built-in

#### Agregar a `pkg/r2libs/r2goroutine.go`
```go
"createMutex": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
    return &sync.Mutex{}
}),

"createRWMutex": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
    return &sync.RWMutex{}
}),

"createWaitGroup": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
    return &sync.WaitGroup{}
}),
```

### Modificaciones al Environment

#### Agregar soporte para mutex implícitos
```go
type Environment struct {
    // campos existentes...
    implicitMutex *sync.Mutex
    classMutexes  map[string]*sync.Mutex
}

func (env *Environment) GetImplicitMutex() *sync.Mutex {
    if env.implicitMutex == nil {
        env.implicitMutex = &sync.Mutex{}
    }
    return env.implicitMutex
}

func (env *Environment) GetClassMutex(className string) *sync.Mutex {
    if env.classMutexes == nil {
        env.classMutexes = make(map[string]*sync.Mutex)
    }
    if env.classMutexes[className] == nil {
        env.classMutexes[className] = &sync.Mutex{}
    }
    return env.classMutexes[className]
}
```

## Ejemplos de Uso Avanzados

### 1. Patrón Productor-Consumidor
```r2
class BufferSeguro {
    constructor(tamaño) {
        this.buffer = [];
        this.tamaño = tamaño;
        this.mutex = createMutex();
        this.condVar = createConditionVariable();
    }
    
    sync agregar(item) {
        while (this.buffer.length >= this.tamaño) {
            this.condVar.wait();
        }
        this.buffer.push(item);
        this.condVar.notifyAll();
    }
    
    sync obtener() {
        while (this.buffer.length == 0) {
            this.condVar.wait();
        }
        let item = this.buffer.shift();
        this.condVar.notifyAll();
        return item;
    }
}

let buffer = new BufferSeguro(10);

// Productor
r2(func() {
    for (let i = 0; i < 100; i++) {
        buffer.agregar("item_" + i);
        sleep(100);
    }
});

// Consumidor
r2(func() {
    for (let i = 0; i < 100; i++) {
        let item = buffer.obtener();
        print("Consumido: " + item);
        sleep(150);
    }
});
```

### 2. Singleton Thread-Safe
```r2
class DatabaseConnection {
    static instance = null;
    static mutex = createMutex();
    
    static sync getInstance() {
        if (DatabaseConnection.instance == null) {
            DatabaseConnection.instance = new DatabaseConnection();
        }
        return DatabaseConnection.instance;
    }
    
    constructor() {
        if (DatabaseConnection.instance != null) {
            throw "Use getInstance() to create DatabaseConnection";
        }
        this.connected = false;
    }
    
    sync connect() {
        if (!this.connected) {
            // Lógica de conexión
            this.connected = true;
        }
    }
}
```

### 3. Pool de Recursos Sincronizado
```r2
class ResourcePool {
    constructor(createResource, maxSize) {
        this.createResource = createResource;
        this.maxSize = maxSize;
        this.available = [];
        this.inUse = new Set();
        this.mutex = createMutex();
    }
    
    sync acquire() {
        while (this.available.length == 0 && this.inUse.size >= this.maxSize) {
            sleep(10); // Esperar recursos disponibles
        }
        
        let resource;
        if (this.available.length > 0) {
            resource = this.available.pop();
        } else {
            resource = this.createResource();
        }
        
        this.inUse.add(resource);
        return resource;
    }
    
    sync release(resource) {
        if (this.inUse.has(resource)) {
            this.inUse.delete(resource);
            this.available.push(resource);
        }
    }
}
```

## Compatibilidad y Migración

### Retrocompatibilidad
- El código R2Lang existente continuará funcionando sin modificaciones
- `sync` se introduce como nueva palabra reservada sin conflictos
- Las funciones de goroutines existentes (`r2()`) se mantienen inalteradas

### Estrategia de Migración
1. **Fase 1**: Implementar funciones mutex básicas
2. **Fase 2**: Agregar soporte para bloques `sync`
3. **Fase 3**: Implementar métodos `sync` en clases
4. **Fase 4**: Optimizaciones y herramientas de debugging

## Beneficios

### Para Desarrolladores
- **Sintaxis Familiar**: Inspirada en Java, fácil adopción
- **Type Safety**: Prevención de race conditions en tiempo de ejecución
- **Productividad**: Menos código boilerplate para sincronización
- **Debugging**: Mejor trazabilidad de problemas de concurrencia

### Para el Ecosistema R2Lang
- **Completitud**: Modelo de concurrencia robusto y completo
- **Performance**: Sincronización eficiente basada en Go's sync package
- **Escalabilidad**: Soporte nativo para aplicaciones multi-threaded
- **Competitividad**: Paridad con lenguajes modernos (Java, C#, Kotlin)

## Consideraciones de Implementación

### Performance
- Utilización directa de `sync.Mutex` de Go para máximo rendimiento
- Overhead mínimo en contextos no sincronizados
- Pool de mutex para optimizar allocación/deallocación

### Memory Management
- Garbage collection automático de mutex no utilizados
- Weak references para evitar memory leaks en objetos sincronizados
- Detección de deadlocks en modo debug

### Error Handling
```r2
try {
    sync(mutexInvalido) {
        // código
    }
} catch (SyncError e) {
    print("Error de sincronización: " + e.message);
}
```

## Testing Strategy

### Unit Tests
```r2
test "sync básico funciona correctamente" {
    given {
        let contador = 0;
        let mutex = createMutex();
        let incrementos = 1000;
    }
    
    when {
        for (let i = 0; i < incrementos; i++) {
            r2(func() {
                sync(mutex) {
                    contador = contador + 1;
                }
            });
        }
        waitForAllGoroutines();
    }
    
    then {
        assert(contador == incrementos, "Contador debe ser " + incrementos);
    }
}
```

### Integration Tests
- Pruebas de stress con múltiples goroutines
- Pruebas de deadlock detection
- Benchmarks de performance vs implementaciones manuales

## Roadmap de Desarrollo

### Milestone 1 (1-2 semanas)
- [ ] Implementar funciones mutex básicas
- [ ] Agregar token SYNC al lexer
- [ ] Tests unitarios para mutex

### Milestone 2 (2-3 semanas)
- [ ] Parser para bloques sync
- [ ] Implementar SyncStatement AST node
- [ ] Tests de integración básicos

### Milestone 3 (3-4 semanas)
- [ ] Soporte para métodos sync en clases
- [ ] Mutex implícitos por instancia
- [ ] Documentación y ejemplos

### Milestone 4 (2-3 semanas)
- [ ] Optimizaciones de performance
- [ ] Herramientas de debugging
- [ ] Tests de stress y benchmarks

## Conclusión

La introducción de la palabra reservada `sync` en R2Lang representaría un avance significativo en las capacidades de concurrencia del lenguaje. Esta propuesta:

1. **Mantiene la filosofía**: Sintaxis clara e intuitiva
2. **Aprovecha Go**: Utiliza las primitivas de sincronización de Go
3. **Es práctica**: Casos de uso reales y comunes
4. **Es escalable**: Desde scripts simples hasta aplicaciones enterprise

La implementación propuesta es técnicamente factible, mantiene la retrocompatibilidad y proporciona herramientas poderosas para el desarrollo de aplicaciones concurrentes robustas.

---

**Autor**: Propuesta para R2Lang  
**Fecha**: Julio 2025  
**Estado**: Borrador para Revisión  
**Versión**: 1.0