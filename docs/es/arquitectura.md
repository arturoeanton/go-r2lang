# Arquitectura de R2Lang

## Visión General del Sistema

R2Lang implementa un intérprete tree-walking con arquitectura modular que separa claramente las responsabilidades del análisis léxico, sintáctico, semántico y ejecución.

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Código R2     │───▶│     Lexer       │───▶│     Parser      │
│                 │    │                 │    │                 │
│ func main() {   │    │ TOKEN_FUNC      │    │   AST Nodes     │
│   print("hi")   │    │ TOKEN_IDENT     │    │                 │
│ }               │    │ TOKEN_SYMBOL    │    │ FunctionDecl    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
┌─────────────────┐    ┌─────────────────┐              │
│   Resultado     │◀───│   Evaluator     │◀─────────────┘
│                 │    │                 │
│ "hi"            │    │ Tree Walking    │
│                 │    │ Interpreter     │
└─────────────────┘    └─────────────────┘
```

## Componentes Principales

### 1. Lexer (Analizador Léxico)
**Archivo**: `r2lang/r2lang.go:72-321`

**Responsabilidades**:
- Tokenización del código fuente
- Reconocimiento de patrones léxicos
- Manejo de comentarios y whitespace
- Detección de errores léxicos

**Algoritmos Clave**:
```go
// Finite State Machine para números
parseNumberOrSign() → Detecta: -123.45, +67, 89.0

// Context-sensitive operator detection  
// Diferencia entre: 
// x-y    (resta)
// (-5)   (número negativo)
```

**Patterns Soportados**:
- Números: `123`, `45.67`, `-89.1`, `+42`
- Strings: `"texto"`, `'texto'`
- Identificadores: `variable`, `función`, `_privado`
- Operadores: `+`, `-`, `*`, `/`, `==`, `<=`, `++`, `--`
- Delimitadores: `(`, `)`, `{`, `}`, `[`, `]`, `;`, `,`

### 2. Parser (Analizador Sintáctico)
**Archivo**: `r2lang/r2lang.go:1662-2331`

**Método**: Recursive Descent Parser con precedencia de operadores
**Lookahead**: 1 token (LL(1))

**Gramática Simplificada**:
```bnf
program         := statement*
statement       := let_stmt | func_decl | class_decl | if_stmt | 
                   while_stmt | for_stmt | try_stmt | expr_stmt
expression      := factor (binary_op factor)*
factor          := number | string | identifier | "(" expression ")" |
                   array_literal | map_literal | function_literal
postfix         := "(" args ")" | "." identifier | "[" expression "]"
```

**Construcciones Especiales**:
- **For-in loops**: `for (let item in array)`
- **Class inheritance**: `class Child extends Parent`
- **TestCase syntax**: `TestCase "name" { Given ... When ... Then ... }`
- **Lambda functions**: `func(x, y) { return x + y }`

### 3. AST (Árbol de Sintaxis Abstracta)
**Archivo**: `r2lang/r2lang.go:327-1657`

**Diseño**: Visitor Pattern implícito con polimorfismo
**Interface Base**:
```go
type Node interface {
    Eval(env *Environment) interface{}
}
```

**Jerarquía de Nodos**:
```
Node
├── Statements
│   ├── Program
│   ├── LetStatement
│   ├── FunctionDeclaration
│   ├── ObjectDeclaration
│   ├── IfStatement
│   ├── WhileStatement
│   ├── ForStatement
│   ├── TryStatement
│   ├── ReturnStatement
│   └── ExprStatement
└── Expressions
    ├── Identifier
    ├── NumberLiteral
    ├── StringLiteral
    ├── BooleanLiteral
    ├── ArrayLiteral
    ├── MapLiteral
    ├── FunctionLiteral
    ├── BinaryExpression
    ├── CallExpression
    ├── AccessExpression
    └── IndexExpression
```

### 4. Environment (Sistema de Entornos)
**Archivo**: `r2lang/r2lang.go:1429-1467`

**Patrón**: Chain of Responsibility para scoping
**Estructura**:
```go
type Environment struct {
    store    map[string]interface{}  // Variables locales
    outer    *Environment           // Environment padre
    imported map[string]bool        // Cache de módulos
    dir      string                 // Directorio actual
    CurrenFx string                // Función en ejecución
}
```

**Scoping Rules**:
```
Global Environment
├── Function Environment (closure)
│   ├── Block Environment
│   └── Inner Block Environment
├── Class Environment
│   ├── Method Environment
│   └── Constructor Environment
└── Module Environment
```

### 5. Evaluator (Intérprete)
**Método**: Tree-walking with immediate evaluation
**Patrón**: Visitor Pattern + Interpreter Pattern

**Flow de Evaluación**:
```
1. Node.Eval(env) llamado polimórficamente
2. Cada nodo evalúa recursivamente sus hijos
3. Los valores se propagan hacia arriba
4. Control flow manejado via ReturnValue exceptions
```

## Patrones de Diseño Utilizados

### 1. Interpreter Pattern
Cada nodo del AST implementa su propia lógica de evaluación:
```go
func (be *BinaryExpression) Eval(env *Environment) interface{} {
    left := be.Left.Eval(env)    // Recursión
    right := be.Right.Eval(env)  // Recursión
    return applyOperator(be.Op, left, right)
}
```

### 2. Chain of Responsibility
Para variable lookup en environments anidados:
```go
func (e *Environment) Get(name string) (interface{}, bool) {
    if val, ok := e.store[name]; ok {
        return val, true
    }
    if e.outer != nil {
        return e.outer.Get(name)  // Delegar al padre
    }
    return nil, false
}
```

### 3. Builder Pattern
Para construcción incremental del AST:
```go
func (p *Parser) parseExpression() Node {
    left := p.parseFactor()
    for isBinaryOp(p.curTok.Value) {
        op := p.curTok.Value
        right := p.parseFactor()
        left = &BinaryExpression{Left: left, Op: op, Right: right}
    }
    return left
}
```

### 4. Strategy Pattern
Para diferentes tipos de funciones:
```go
switch cv := calleeVal.(type) {
case BuiltinFunction:
    return cv(argVals...)        // Función nativa
case *UserFunction:
    return cv.Call(argVals...)   // Función R2
case map[string]interface{}:
    return instantiateObject(env, cv, argVals)  // Constructor
}
```

## Subsistemas Especializados

### Sistema de Tipos Runtime
**Ubicación**: `r2lang/r2lang.go:1513-1621`

**Características**:
- Dynamic typing con duck typing
- Automatic type coercion
- Type reflection via `typeOf()`

**Type Hierarchy**:
```
interface{}
├── float64        (números)
├── string         (cadenas)
├── bool           (booleanos)
├── []interface{}  (arrays)
├── map[string]interface{}  (maps/objetos)
├── *UserFunction  (funciones usuario)
├── BuiltinFunction (funciones nativas)
├── *ObjectInstance (instancias de clase)
└── nil            (valor nulo)
```

### Sistema de Orientación a Objetos
**Patrón**: Prototype-based inheritance (como JavaScript)

**Blueprints de Clase**:
```go
// Una clase es un map que sirve como template
blueprint := map[string]interface{}{
    "ClassName": "MiClase",
    "propiedad": nil,
    "metodo": &UserFunction{...},
    "super": parentBlueprint,  // Para herencia
}
```

**Instanciación**:
```go
// Crear nuevo environment para la instancia
objEnv := NewInnerEnv(globalEnv)
instance := &ObjectInstance{Env: objEnv}

// Copiar y bindear métodos
for k, v := range blueprint {
    if method, ok := v.(*UserFunction); ok {
        method.Env = objEnv  // Bind al objeto
    }
    objEnv.Set(k, v)
}

// Self-reference
objEnv.Set("this", instance)
objEnv.Set("self", instance)
```

### Sistema de Concurrencia
**Modelo**: Goroutines de Go con WaitGroup global

**Arquitectura**:
```go
var wg sync.WaitGroup  // Global wait group

// r2(function) crea nueva goroutine
wg.Add(1)
go func() {
    defer wg.Done()
    userFunction.Call(args...)
}()

// Al final del programa: wg.Wait()
```

**Limitaciones Actuales**:
- No hay comunicación entre goroutines (channels)
- No hay primitivas de sincronización (mutexes)
- Error handling básico en goroutines

### Sistema de Módulos
**Patrón**: Import estático con cache

**Algoritmo de Import**:
```go
1. Resolver path relativo/absoluto
2. Verificar cache de imports (evitar ciclos)
3. Leer archivo y parsear
4. Crear nuevo environment para el módulo
5. Evaluar módulo en su environment
6. Exportar símbolos según alias
```

**Resolución de Paths**:
```go
// Relativo al archivo actual
import "./utils.r2"

// Relativo al directorio base
import "lib/math.r2"

// Con alias
import "./helpers.r2" as h
```

## Flujo de Datos

### Pipeline de Ejecución
```
1. main.go recibe argumentos CLI
2. Lee archivo .r2 o inicia REPL
3. NewLexer(sourceCode)
4. NewParser(lexer)
5. parser.ParseProgram() → AST
6. NewEnvironment()
7. RegisterBuiltins(env)
8. ast.Eval(env) → Resultado
```

### Manejo de Errores
**Estrategia**: Panic/Recover con propagación

```go
// Lexer errors
panic("Unexpected character: " + ch)

// Parser errors  
panic("Expected ')' after expression")

// Runtime errors
panic("Undefined variable: " + name)

// Try-catch maneja via recover()
defer func() {
    if r := recover(); r != nil {
        // Bind error to catch variable
        catchEnv.Set(errorVar, r)
        catchBlock.Eval(catchEnv)
    }
}()
```

### Memory Management
**Estrategia**: Delegación al GC de Go

**Lifecycle de Objetos**:
- Variables: Lifetime = scope del environment
- Functions: Lifetime = closure + referencias activas
- Objects: Lifetime = referencias en environments
- Arrays/Maps: Lifetime = referencias activas

**Memory Leaks Potenciales**:
- Closures que capturan environments grandes
- Ciclos de referencia en objetos
- Goroutines que no terminan

## Extensibilidad

### Añadir Nuevas Características del Lenguaje

**1. Nuevos Tokens**:
```go
// En constants
const TOKEN_NEWFEATURE = "NEWFEATURE"

// En lexer
case "newkeyword":
    return Token{Type: TOKEN_NEWFEATURE, Value: literal}
```

**2. Nuevos Nodos AST**:
```go
type NewFeatureNode struct {
    Property string
    Child    Node
}

func (nfn *NewFeatureNode) Eval(env *Environment) interface{} {
    // Implementar semántica
    return result
}
```

**3. Parsing Logic**:
```go
func (p *Parser) parseNewFeature() Node {
    // Parsear sintaxis específica
    return &NewFeatureNode{...}
}
```

### Añadir Bibliotecas Nativas

**Template**:
```go
// r2lang/r2mynewlib.go
func RegisterMyNewLib(env *Environment) {
    env.Set("myFunction", BuiltinFunction(func(args ...interface{}) interface{} {
        // Implementación en Go
        return result
    }))
    
    env.Set("myObject", map[string]interface{}{
        "method1": BuiltinFunction(func(args ...interface{}) interface{} {
            // Método 1
            return result
        }),
        "property": "value",
    })
}

// En RunCode():
RegisterMyNewLib(env)
```

## Performance Considerations

### Bottlenecks Actuales

1. **Variable Lookup**: O(depth) en environments anidados
2. **Function Calls**: Creación de nuevos environments
3. **Type Checking**: Runtime type assertions
4. **AST Walking**: No optimización de hot paths
5. **String Operations**: Concatenación frecuente

### Optimizaciones Implementables

1. **Variable Caching**: Hash table global para variables frecuentes
2. **Bytecode Compilation**: AST → Bytecode → Execution
3. **JIT Compilation**: Hot path optimization
4. **Constant Folding**: Evaluar expresiones constantes en parse time
5. **Tail Call Optimization**: Para recursión eficiente

## Comparación con Otros Intérpretes

### Ventajas de R2Lang
- **Simplicidad**: Arquitectura clara y comprensible
- **Extensibilidad**: Fácil agregar nuevas bibliotecas
- **Modernidad**: Sintaxis familiar (JavaScript-like)
- **Integración**: Built-in testing y concurrencia

### Limitaciones vs. Otros Lenguajes
- **Performance**: Más lento que intérpretes optimizados (Python, Ruby)
- **Memory**: Sin GC optimizado para el lenguaje
- **Ecosystem**: Biblioteca estándar limitada
- **Tooling**: Sin debugger, profiler, etc.

## Roadmap Arquitectural

### Corto Plazo (1-3 meses)
- Bytecode compilation layer
- Optimized variable lookup
- Error reporting mejorado

### Medio Plazo (3-6 meses)  
- JIT compilation
- Garbage collector específico
- Standard library expandida

### Largo Plazo (6-12 meses)
- Language Server Protocol
- Multi-target compilation (WASM, native)
- Advanced optimization passes

## Conclusión

La arquitectura de R2Lang prioriza claridad y extensibilidad sobre performance bruta, resultando en un intérprete que es fácil de entender, modificar y extender. Esta decisión facilita el desarrollo rápido de nuevas características y el aprendizaje del codebase, aunque limita el rendimiento en aplicaciones intensivas.

La modularidad del diseño permite optimizaciones incrementales sin reestructuración completa, posicionando a R2Lang para evolución continua hacia un lenguaje de producción robusto.