# Propuesta: Sistema de Debugging para R2Lang

**Versión:** 1.0  
**Fecha:** 2025-07-15  
**Estado:** Propuesta  

## Resumen Ejecutivo

Esta propuesta presenta un sistema integral de debugging para R2Lang que incluye debugger interactivo, breakpoints, step debugging, inspección de variables, stack traces y integración con IDEs, proporcionando una experiencia de desarrollo moderna y productiva.

## Problema Actual

R2Lang actualmente carece de herramientas de debugging avanzadas:

- **No hay debugger interactivo**
- **Debugging limitado a print statements**
- **No hay breakpoints**
- **Stack traces básicos**
- **No hay inspección de variables en runtime**
- **No hay integración con IDEs**

```r2
// Debugging actual - primitivo
func problematicFunction(x) {
    print("Debug: x =", x);  // Print debugging
    
    if (x > 0) {
        print("Debug: entering if branch");
        return x * 2;
    }
    
    print("Debug: entering else branch");
    return -1;
}
```

## Solución Propuesta

### 1. Debugger Interactivo CLI

#### 1.1 Comandos Básicos
```bash
# Iniciar debugger
r2 debug script.r2
r2 debug --port 9229 script.r2    # Debug server mode

# Comandos dentro del debugger
(r2db) break main.r2:10           # Set breakpoint
(r2db) break functionName         # Break on function
(r2db) run                        # Start execution
(r2db) continue                   # Continue execution
(r2db) step                       # Step into
(r2db) next                       # Step over
(r2db) finish                     # Step out
(r2db) print variable             # Print variable
(r2db) info locals                # Show local variables
(r2db) info stack                 # Show call stack
(r2db) list                       # Show source code
(r2db) quit                       # Exit debugger
```

#### 1.2 Interfaz del Debugger
```
R2Lang Debugger v1.0
Loading script: examples/user_service.r2
Reading symbols...done.

Breakpoint 1 at user_service.r2:15
(r2db) run
Starting program: user_service.r2

Breakpoint 1, registerUser (name="John", email="john@example.com") at user_service.r2:15
15    func registerUser(name, email) {
(r2db) list
10        return users.length;
11    }
12    
13    // Register new user
14    func registerUser(name, email) {
15 ->     if (!isValidEmail(email)) {
16            throw new Error("Invalid email format");
17        }
18        
19        let newUser = {
20            id: generateId(),

(r2db) print name
$1 = "John"
(r2db) print email  
$2 = "john@example.com"
(r2db) step
```

### 2. Implementación del Debug Runtime

#### 2.1 Debug Context
```go
// pkg/r2core/debug_context.go
type DebugContext struct {
    Enabled          bool
    Breakpoints      map[string]*Breakpoint
    CurrentLocation  *Location
    CallStack        []*StackFrame
    Variables        map[string]interface{}
    StepMode         StepMode
    DebuggerClient   *DebuggerClient
    TraceMode        bool
    WatchExpressions []string
}

type Breakpoint struct {
    ID          int
    File        string
    Line        int
    Column      int
    Condition   string
    Enabled     bool
    HitCount    int
    HitCondition string
}

type StackFrame struct {
    Function    string
    File        string
    Line        int
    Column      int
    Variables   map[string]interface{}
    Arguments   []interface{}
}

type StepMode int

const (
    StepModeNone StepMode = iota
    StepModeInto
    StepModeOver
    StepModeOut
)

func (dc *DebugContext) SetBreakpoint(file string, line int, condition string) *Breakpoint {
    key := fmt.Sprintf("%s:%d", file, line)
    
    bp := &Breakpoint{
        ID:        len(dc.Breakpoints) + 1,
        File:      file,
        Line:      line,
        Condition: condition,
        Enabled:   true,
    }
    
    dc.Breakpoints[key] = bp
    return bp
}

func (dc *DebugContext) ShouldBreak(file string, line int) bool {
    if !dc.Enabled {
        return false
    }
    
    // Check step mode
    if dc.StepMode != StepModeNone {
        return dc.checkStepMode(file, line)
    }
    
    // Check breakpoints
    key := fmt.Sprintf("%s:%d", file, line)
    if bp, exists := dc.Breakpoints[key]; exists && bp.Enabled {
        return dc.evaluateBreakpointCondition(bp)
    }
    
    return false
}
```

#### 2.2 Integración con Parser y Evaluator
```go
// Instrumentar evaluación para debugging
func (ws *WhileStatement) Eval(env *Environment) interface{} {
    debugCtx := env.GetDebugContext()
    
    if debugCtx != nil && debugCtx.Enabled {
        location := &Location{
            File:   ws.GetFile(),
            Line:   ws.GetLine(),
            Column: ws.GetColumn(),
        }
        
        if debugCtx.ShouldBreak(location.File, location.Line) {
            debugCtx.HandleBreakpoint(location, env)
        }
    }
    
    // Evaluación normal
    for {
        condition := ws.Condition.Eval(env)
        if !isTruthy(condition) {
            break
        }
        
        // Debug trace para cada iteración
        if debugCtx != nil && debugCtx.TraceMode {
            debugCtx.TraceStatement("while loop iteration", ws.GetLocation())
        }
        
        ws.Body.Eval(env)
    }
    
    return nil
}

// Instrumentar llamadas a funciones
func (fc *FunctionCall) Eval(env *Environment) interface{} {
    debugCtx := env.GetDebugContext()
    
    if debugCtx != nil && debugCtx.Enabled {
        // Crear stack frame
        frame := &StackFrame{
            Function:  fc.FunctionName,
            File:      fc.GetFile(),
            Line:      fc.GetLine(),
            Column:    fc.GetColumn(),
            Arguments: fc.evaluateArguments(env),
        }
        
        debugCtx.PushStackFrame(frame)
        defer debugCtx.PopStackFrame()
        
        // Check breakpoint
        if debugCtx.ShouldBreak(frame.File, frame.Line) {
            debugCtx.HandleBreakpoint(&Location{
                File:   frame.File,
                Line:   frame.Line,
                Column: frame.Column,
            }, env)
        }
    }
    
    return fc.executeFunction(env)
}
```

### 3. Debug Server Protocol

#### 3.1 Debug Adapter Protocol (DAP)
```go
// pkg/r2core/debug_server.go
type DebugServer struct {
    Port           int
    Clients        map[string]*DebugClient
    Runtime        *Runtime
    SessionManager *SessionManager
}

type DebugClient struct {
    ID          string
    Connection  net.Conn
    Encoder     *json.Encoder
    Decoder     *json.Decoder
    Session     *DebugSession
}

type DebugSession struct {
    ID            string
    Runtime       *Runtime
    Breakpoints   map[string]*Breakpoint
    StoppedReason string
    StoppedThread int
}

func (ds *DebugServer) Start() error {
    listener, err := net.Listen("tcp", fmt.Sprintf(":%d", ds.Port))
    if err != nil {
        return err
    }
    
    fmt.Printf("Debug server listening on port %d\n", ds.Port)
    
    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        
        client := &DebugClient{
            ID:         generateClientID(),
            Connection: conn,
            Encoder:    json.NewEncoder(conn),
            Decoder:    json.NewDecoder(conn),
        }
        
        ds.Clients[client.ID] = client
        go ds.handleClient(client)
    }
}

func (ds *DebugServer) handleClient(client *DebugClient) {
    defer func() {
        client.Connection.Close()
        delete(ds.Clients, client.ID)
    }()
    
    for {
        var request DebugRequest
        if err := client.Decoder.Decode(&request); err != nil {
            return
        }
        
        response := ds.processRequest(client, &request)
        if err := client.Encoder.Encode(response); err != nil {
            return
        }
    }
}

// Implementar DAP messages
func (ds *DebugServer) processRequest(client *DebugClient, request *DebugRequest) *DebugResponse {
    switch request.Command {
    case "initialize":
        return ds.handleInitialize(client, request)
    case "launch":
        return ds.handleLaunch(client, request)
    case "setBreakpoints":
        return ds.handleSetBreakpoints(client, request)
    case "continue":
        return ds.handleContinue(client, request)
    case "next":
        return ds.handleNext(client, request)
    case "stepIn":
        return ds.handleStepIn(client, request)
    case "stepOut":
        return ds.handleStepOut(client, request)
    case "pause":
        return ds.handlePause(client, request)
    case "variables":
        return ds.handleVariables(client, request)
    case "evaluate":
        return ds.handleEvaluate(client, request)
    case "stackTrace":
        return ds.handleStackTrace(client, request)
    default:
        return &DebugResponse{
            Success: false,
            Message: fmt.Sprintf("Unknown command: %s", request.Command),
        }
    }
}
```

#### 3.2 Debug Events
```go
// Debug events para comunicación con IDE
type DebugEvent struct {
    Type string      `json:"type"`
    Body interface{} `json:"body"`
}

type StoppedEvent struct {
    Reason      string `json:"reason"`
    ThreadID    int    `json:"threadId"`
    Text        string `json:"text,omitempty"`
    AllThreads  bool   `json:"allThreadsStopped,omitempty"`
}

type BreakpointEvent struct {
    Reason     string      `json:"reason"`
    Breakpoint *Breakpoint `json:"breakpoint"`
}

func (ds *DebugServer) sendStoppedEvent(client *DebugClient, reason string, threadID int) {
    event := &DebugEvent{
        Type: "stopped",
        Body: &StoppedEvent{
            Reason:     reason,
            ThreadID:   threadID,
            AllThreads: true,
        },
    }
    
    client.Encoder.Encode(event)
}

func (ds *DebugServer) sendBreakpointEvent(client *DebugClient, reason string, bp *Breakpoint) {
    event := &DebugEvent{
        Type: "breakpoint",
        Body: &BreakpointEvent{
            Reason:     reason,
            Breakpoint: bp,
        },
    }
    
    client.Encoder.Encode(event)
}
```

### 4. Variable Inspection y Evaluation

#### 4.1 Variable Inspector
```go
// pkg/r2core/variable_inspector.go
type VariableInspector struct {
    Context *DebugContext
    Scopes  []*Scope
}

type Scope struct {
    Name      string
    Variables map[string]*Variable
    Expensive bool
}

type Variable struct {
    Name             string
    Value            interface{}
    Type             string
    VariablesRef     int
    NamedVariables   int
    IndexedVariables int
    Expensive        bool
}

func (vi *VariableInspector) GetScopes() []*Scope {
    var scopes []*Scope
    
    // Local scope
    if frame := vi.Context.GetCurrentFrame(); frame != nil {
        localScope := &Scope{
            Name:      "Local",
            Variables: vi.convertVariables(frame.Variables),
        }
        scopes = append(scopes, localScope)
    }
    
    // Global scope
    globalScope := &Scope{
        Name:      "Global",
        Variables: vi.convertVariables(vi.Context.GetGlobalVariables()),
    }
    scopes = append(scopes, globalScope)
    
    return scopes
}

func (vi *VariableInspector) GetVariable(varRef int) (*Variable, error) {
    // Resolver referencia de variable
    if varRef == 0 {
        return nil, fmt.Errorf("invalid variable reference")
    }
    
    variable := vi.Context.GetVariableByRef(varRef)
    if variable == nil {
        return nil, fmt.Errorf("variable not found")
    }
    
    return variable, nil
}

func (vi *VariableInspector) EvaluateExpression(expr string, frameID int) (*Variable, error) {
    // Evaluar expresión en contexto del frame
    frame := vi.Context.GetFrame(frameID)
    if frame == nil {
        return nil, fmt.Errorf("frame not found")
    }
    
    // Crear environment temporal con variables del frame
    tempEnv := NewEnvironment()
    for name, value := range frame.Variables {
        tempEnv.Set(name, value)
    }
    
    // Parsear y evaluar expresión
    parser := NewParser(expr)
    ast := parser.ParseExpression()
    
    result := ast.Eval(tempEnv)
    
    return &Variable{
        Name:  expr,
        Value: result,
        Type:  vi.getTypeString(result),
    }, nil
}
```

#### 4.2 Watch Expressions
```r2
// En el debugger, soporte para watch expressions
(r2db) watch user.name
Watch 1: user.name
(r2db) watch userCount > 10
Watch 2: userCount > 10
(r2db) info watches
Num  Expression    Value
1    user.name     "John Doe"
2    userCount > 10  false
```

### 5. CLI Debugger Interface

#### 5.1 Interactive Debugger
```go
// pkg/r2core/cli_debugger.go
type CLIDebugger struct {
    Runtime     *Runtime
    Context     *DebugContext
    Input       *bufio.Scanner
    Output      io.Writer
    Commands    map[string]DebugCommand
    History     []string
    CurrentFile string
    CurrentLine int
}

type DebugCommand struct {
    Name        string
    Aliases     []string
    Description string
    Handler     func(args []string) error
}

func NewCLIDebugger() *CLIDebugger {
    debugger := &CLIDebugger{
        Input:    bufio.NewScanner(os.Stdin),
        Output:   os.Stdout,
        Commands: make(map[string]DebugCommand),
    }
    
    debugger.registerCommands()
    return debugger
}

func (cd *CLIDebugger) registerCommands() {
    commands := []DebugCommand{
        {
            Name:        "break",
            Aliases:     []string{"b"},
            Description: "Set breakpoint",
            Handler:     cd.handleBreakpoint,
        },
        {
            Name:        "continue",
            Aliases:     []string{"c"},
            Description: "Continue execution",
            Handler:     cd.handleContinue,
        },
        {
            Name:        "step",
            Aliases:     []string{"s"},
            Description: "Step into",
            Handler:     cd.handleStep,
        },
        {
            Name:        "next",
            Aliases:     []string{"n"},
            Description: "Step over",
            Handler:     cd.handleNext,
        },
        {
            Name:        "print",
            Aliases:     []string{"p"},
            Description: "Print variable",
            Handler:     cd.handlePrint,
        },
        {
            Name:        "list",
            Aliases:     []string{"l"},
            Description: "List source code",
            Handler:     cd.handleList,
        },
        {
            Name:        "info",
            Aliases:     []string{"i"},
            Description: "Show information",
            Handler:     cd.handleInfo,
        },
    }
    
    for _, cmd := range commands {
        cd.Commands[cmd.Name] = cmd
        for _, alias := range cmd.Aliases {
            cd.Commands[alias] = cmd
        }
    }
}

func (cd *CLIDebugger) Run() {
    cd.showWelcome()
    
    for {
        cd.showPrompt()
        
        if !cd.Input.Scan() {
            break
        }
        
        line := strings.TrimSpace(cd.Input.Text())
        if line == "" {
            continue
        }
        
        cd.History = append(cd.History, line)
        
        if line == "quit" || line == "q" {
            break
        }
        
        cd.processCommand(line)
    }
}

func (cd *CLIDebugger) processCommand(line string) {
    parts := strings.Fields(line)
    if len(parts) == 0 {
        return
    }
    
    cmdName := parts[0]
    args := parts[1:]
    
    if cmd, exists := cd.Commands[cmdName]; exists {
        if err := cmd.Handler(args); err != nil {
            fmt.Fprintf(cd.Output, "Error: %v\n", err)
        }
    } else {
        fmt.Fprintf(cd.Output, "Unknown command: %s\n", cmdName)
        fmt.Fprintf(cd.Output, "Type 'help' for available commands\n")
    }
}
```

#### 5.2 Command Handlers
```go
func (cd *CLIDebugger) handleBreakpoint(args []string) error {
    if len(args) == 0 {
        return fmt.Errorf("breakpoint requires location")
    }
    
    location := args[0]
    
    // Parse location (file:line or function name)
    if strings.Contains(location, ":") {
        parts := strings.Split(location, ":")
        if len(parts) != 2 {
            return fmt.Errorf("invalid breakpoint format")
        }
        
        file := parts[0]
        line, err := strconv.Atoi(parts[1])
        if err != nil {
            return fmt.Errorf("invalid line number")
        }
        
        bp := cd.Context.SetBreakpoint(file, line, "")
        fmt.Fprintf(cd.Output, "Breakpoint %d set at %s:%d\n", bp.ID, file, line)
        
    } else {
        // Function breakpoint
        bp := cd.Context.SetFunctionBreakpoint(location)
        fmt.Fprintf(cd.Output, "Breakpoint %d set at function %s\n", bp.ID, location)
    }
    
    return nil
}

func (cd *CLIDebugger) handlePrint(args []string) error {
    if len(args) == 0 {
        return fmt.Errorf("print requires variable name")
    }
    
    varName := args[0]
    
    variable, err := cd.Context.GetVariable(varName)
    if err != nil {
        return fmt.Errorf("variable not found: %s", varName)
    }
    
    value := cd.formatValue(variable.Value)
    fmt.Fprintf(cd.Output, "$%d = %s\n", cd.Context.GetNextPrintID(), value)
    
    return nil
}

func (cd *CLIDebugger) handleList(args []string) error {
    file := cd.CurrentFile
    startLine := cd.CurrentLine - 5
    endLine := cd.CurrentLine + 5
    
    if len(args) > 0 {
        if num, err := strconv.Atoi(args[0]); err == nil {
            startLine = num - 5
            endLine = num + 5
        }
    }
    
    lines, err := cd.getSourceLines(file, startLine, endLine)
    if err != nil {
        return err
    }
    
    for lineNum, line := range lines {
        marker := "  "
        if lineNum == cd.CurrentLine {
            marker = "->"
        }
        
        fmt.Fprintf(cd.Output, "%s %d\t%s\n", marker, lineNum, line)
    }
    
    return nil
}
```

### 6. IDE Integration

#### 6.1 VS Code Extension
```typescript
// vscode-r2lang-debug/src/extension.ts
import * as vscode from 'vscode';
import { R2DebugAdapterFactory } from './r2DebugAdapter';

export function activate(context: vscode.ExtensionContext) {
    const factory = new R2DebugAdapterFactory();
    
    context.subscriptions.push(
        vscode.debug.registerDebugAdapterDescriptorFactory('r2lang', factory)
    );
    
    context.subscriptions.push(
        vscode.debug.registerDebugConfigurationProvider('r2lang', {
            provideDebugConfigurations: () => {
                return [{
                    name: 'Launch R2Lang',
                    type: 'r2lang',
                    request: 'launch',
                    program: '${workspaceFolder}/main.r2',
                    stopOnEntry: false,
                    args: [],
                    cwd: '${workspaceFolder}',
                    env: {},
                    debugServer: 9229
                }];
            }
        })
    );
}

// vscode-r2lang-debug/src/r2DebugAdapter.ts
import { DebugProtocol } from 'vscode-debugprotocol';
import { DebugSession } from 'vscode-debugadapter';

export class R2DebugSession extends DebugSession {
    private runtime: R2Runtime;
    
    protected initializeRequest(response: DebugProtocol.InitializeResponse): void {
        response.body = {
            supportsConfigurationDoneRequest: true,
            supportsEvaluateForHovers: true,
            supportsStepBack: false,
            supportsDataBreakpoints: false,
            supportsCompletionsRequest: true,
            supportsBreakpointLocationsRequest: true,
            supportsSetVariable: true,
            supportsRestartRequest: true,
            supportsExceptionOptions: true,
            supportsValueFormattingOptions: true,
            supportsExceptionInfoRequest: true,
            supportTerminateDebuggee: true,
            supportsDelayedStackTraceLoading: true,
            supportsLoadedSourcesRequest: true,
            supportsLogPoints: true,
            supportsTerminateThreadsRequest: true,
            supportsSetExpression: true,
            supportsTerminateRequest: true,
        };
        
        this.sendResponse(response);
    }
    
    protected launchRequest(response: DebugProtocol.LaunchResponse, args: any): void {
        this.runtime = new R2Runtime(args.program);
        this.runtime.on('stopOnEntry', () => {
            this.sendEvent(new StoppedEvent('entry', 1));
        });
        
        this.runtime.on('stopOnBreakpoint', () => {
            this.sendEvent(new StoppedEvent('breakpoint', 1));
        });
        
        this.runtime.start();
        this.sendResponse(response);
    }
}
```

#### 6.2 Debug Configuration
```json
// .vscode/launch.json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug R2Lang",
            "type": "r2lang",
            "request": "launch",
            "program": "${workspaceFolder}/main.r2",
            "stopOnEntry": false,
            "args": [],
            "cwd": "${workspaceFolder}",
            "env": {},
            "debugServer": 9229,
            "trace": false
        },
        {
            "name": "Attach to R2Lang",
            "type": "r2lang",
            "request": "attach",
            "port": 9229,
            "host": "localhost"
        }
    ]
}
```

### 7. Advanced Debugging Features

#### 7.1 Conditional Breakpoints
```bash
# Breakpoint condicional
(r2db) break user_service.r2:25 if user.age > 18
(r2db) break validateUser if email.contains("@admin")

# Hit count breakpoints
(r2db) break main.r2:10 hit 5      # Break on 5th hit
(r2db) break loop.r2:15 hit > 10   # Break after 10 hits
```

#### 7.2 Exception Breakpoints
```r2
// Configurar exception breakpoints
(r2db) catch all                   # Break on any exception
(r2db) catch NetworkError          # Break on specific exception
(r2db) catch uncaught             # Break on uncaught exceptions

class CustomError extends Error {
    constructor(message) {
        super(message);
        this.name = "CustomError";
    }
}

func riskyOperation() {
    throw new CustomError("Something went wrong");
}
```

#### 7.3 Memory Debugging
```r2
// Memory debugging commands
(r2db) info memory                 # Show memory usage
(r2db) info heap                   # Show heap statistics
(r2db) info gc                     # Show GC statistics
(r2db) watch memory > 100MB        # Break when memory exceeds limit

// Memory profiling integration
(r2db) profile memory start
(r2db) continue
(r2db) profile memory stop
(r2db) profile memory report
```

### 8. Testing Integration

#### 8.1 Test Debugging
```r2
// Debug tests
r2 debug test user_service_test.r2
r2 debug test --grep "should validate email"

// En el test
describe("UserService", func() {
    it("should validate email", func() {
        let user = {name: "John", email: "invalid-email"};
        
        // Debugger se detiene aquí si hay breakpoint
        let result = validateUser(user);
        
        assert.isFalse(result.valid);
    });
});
```

#### 8.2 Debug Test Runner
```bash
# Test runner con debugging
r2 test debug user_service_test.r2
r2 test debug --break-on-failure
r2 test debug --inspect-failures
```

### 9. Remote Debugging

#### 9.1 Remote Debug Server
```bash
# Iniciar servidor de debug remoto
r2 debug --remote --host 0.0.0.0 --port 9229 script.r2

# Conectar desde cliente remoto
r2 debug --attach --host remote-server --port 9229
```

#### 9.2 Debug en Containers
```dockerfile
# Dockerfile con debugging
FROM r2lang:latest

EXPOSE 9229

# Habilitar debugging
ENV R2_DEBUG=true
ENV R2_DEBUG_PORT=9229

CMD ["r2", "debug", "--remote", "--host", "0.0.0.0", "app.r2"]
```

## Plan de Implementación

### Fase 1: Core Debugging Infrastructure
- [ ] Debug context y runtime integration
- [ ] Básicos breakpoints y step debugging
- [ ] CLI debugger básico
- [ ] Stack traces y variable inspection

### Fase 2: Advanced Features
- [ ] Conditional breakpoints
- [ ] Exception breakpoints
- [ ] Watch expressions
- [ ] Memory debugging

### Fase 3: IDE Integration
- [ ] Debug Adapter Protocol implementation
- [ ] VS Code extension
- [ ] IntelliJ/JetBrains plugin
- [ ] Vim/Neovim integration

### Fase 4: Remote y Advanced
- [ ] Remote debugging
- [ ] Debug server mode
- [ ] Performance profiling integration
- [ ] Test debugging

### Fase 5: Ecosystem
- [ ] Debug visualization tools
- [ ] Debug log analysis
- [ ] CI/CD debugging integration
- [ ] Documentation completa

## Beneficios

1. **Productividad:** Debugging eficiente reduce tiempo de desarrollo
2. **Calidad:** Identificación temprana de bugs
3. **Aprendizaje:** Mejor comprensión del código
4. **IDE Integration:** Experiencia de desarrollo moderna
5. **Escalabilidad:** Debugging en aplicaciones complejas

## Conclusión

Este sistema de debugging proporcionará a R2Lang capacidades de debugging de nivel profesional, permitiendo a los desarrolladores identificar y resolver problemas de manera eficiente, con herramientas modernas y integración completa con el ecosistema de desarrollo.