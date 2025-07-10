# R2Lang Architecture

## System Overview

R2Lang implements a tree-walking interpreter with modular architecture that clearly separates responsibilities of lexical analysis, syntactic analysis, semantic analysis and execution.

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   R2 Code       │───▶│     Lexer       │───▶│     Parser      │
│                 │    │                 │    │                 │
│ func main() {   │    │ TOKEN_FUNC      │    │   AST Nodes     │
│   print("hi")   │    │ TOKEN_IDENT     │    │                 │
│ }               │    │ TOKEN_SYMBOL    │    │ FunctionDecl    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
┌─────────────────┐    ┌─────────────────┐              │
│   Result        │◀───│   Evaluator     │◀─────────────┘
│                 │    │                 │
│ "hi"            │    │ Tree Walking    │
│                 │    │ Interpreter     │
└─────────────────┘    └─────────────────┘
```

## Main Components

### 1. Lexer (Lexical Analyzer)
**File**: `r2lang/r2lang.go:72-321`

**Responsibilities**:
- Source code tokenization
- Lexical pattern recognition
- Comment and whitespace handling
- Lexical error detection

**Key Algorithms**:
```go
// Finite State Machine for numbers
parseNumberOrSign() → Detects: -123.45, +67, 89.0

// Context-sensitive operator detection  
// Differentiates between: 
// x-y    (subtraction)
// (-5)   (negative number)
```

**Supported Patterns**:
- Numbers: `123`, `45.67`, `-89.1`, `+42`
- Strings: `"text"`, `'text'`
- Identifiers: `variable`, `function`, `_private`
- Operators: `+`, `-`, `*`, `/`, `==`, `<=`, `++`, `--`
- Delimiters: `(`, `)`, `{`, `}`, `[`, `]`, `;`, `,`

### 2. Parser (Syntactic Analyzer)
**File**: `r2lang/r2lang.go:1662-2331`

**Method**: Recursive Descent Parser with operator precedence
**Lookahead**: 1 token (LL(1))

**Simplified Grammar**:
```bnf
program         := statement*
statement       := let_stmt | func_decl | class_decl | if_stmt | 
                   while_stmt | for_stmt | try_stmt | expr_stmt
expression      := factor (binary_op factor)*
factor          := number | string | identifier | "(" expression ")" |
                   array_literal | map_literal | function_literal
postfix         := "(" args ")" | "." identifier | "[" expression "]"
```

**Special Constructs**:
- **For-in loops**: `for (let item in array)`
- **Class inheritance**: `class Child extends Parent`
- **TestCase syntax**: `TestCase "name" { Given ... When ... Then ... }`
- **Lambda functions**: `func(x, y) { return x + y }`

### 3. AST (Abstract Syntax Tree)
**File**: `r2lang/r2lang.go:327-1657`

**Design**: Implicit Visitor Pattern with polymorphism
**Base Interface**:
```go
type Node interface {
    Eval(env *Environment) interface{}
}
```

**Node Hierarchy**:
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

### 4. Environment (Environment System)
**File**: `r2lang/r2lang.go:1429-1467`

**Pattern**: Chain of Responsibility for scoping
**Structure**:
```go
type Environment struct {
    store    map[string]interface{}  // Local variables
    outer    *Environment           // Parent environment
    imported map[string]bool        // Module cache
    dir      string                 // Current directory
    CurrenFx string                // Function in execution
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

### 5. Evaluator (Interpreter)
**Method**: Tree-walking with immediate evaluation
**Pattern**: Visitor Pattern + Interpreter Pattern

**Evaluation Flow**:
```
1. Node.Eval(env) called polymorphically
2. Each node recursively evaluates its children
3. Values propagate upward
4. Control flow handled via ReturnValue exceptions
```

## Design Patterns Used

### 1. Interpreter Pattern
Each AST node implements its own evaluation logic:
```go
func (be *BinaryExpression) Eval(env *Environment) interface{} {
    left := be.Left.Eval(env)    // Recursion
    right := be.Right.Eval(env)  // Recursion
    return applyOperator(be.Op, left, right)
}
```

### 2. Chain of Responsibility
For variable lookup in nested environments:
```go
func (e *Environment) Get(name string) (interface{}, bool) {
    if val, ok := e.store[name]; ok {
        return val, true
    }
    if e.outer != nil {
        return e.outer.Get(name)  // Delegate to parent
    }
    return nil, false
}
```

### 3. Builder Pattern
For incremental AST construction:
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
For different function types:
```go
switch cv := calleeVal.(type) {
case BuiltinFunction:
    return cv(argVals...)        // Native function
case *UserFunction:
    return cv.Call(argVals...)   // R2 function
case map[string]interface{}:
    return instantiateObject(env, cv, argVals)  // Constructor
}
```

## Specialized Subsystems

### Runtime Type System
**Location**: `r2lang/r2lang.go:1513-1621`

**Features**:
- Dynamic typing with duck typing
- Automatic type coercion
- Type reflection via `typeOf()`

**Type Hierarchy**:
```
interface{}
├── float64        (numbers)
├── string         (strings)
├── bool           (booleans)
├── []interface{}  (arrays)
├── map[string]interface{}  (maps/objects)
├── *UserFunction  (user functions)
├── BuiltinFunction (native functions)
├── *ObjectInstance (class instances)
└── nil            (null value)
```

### Object-Oriented System
**Pattern**: Prototype-based inheritance (like JavaScript)

**Class Blueprints**:
```go
// A class is a map serving as template
blueprint := map[string]interface{}{
    "ClassName": "MyClass",
    "property": nil,
    "method": &UserFunction{...},
    "super": parentBlueprint,  // For inheritance
}
```

**Instantiation**:
```go
// Create new environment for instance
objEnv := NewInnerEnv(globalEnv)
instance := &ObjectInstance{Env: objEnv}

// Copy and bind methods
for k, v := range blueprint {
    if method, ok := v.(*UserFunction); ok {
        method.Env = objEnv  // Bind to object
    }
    objEnv.Set(k, v)
}

// Self-reference
objEnv.Set("this", instance)
objEnv.Set("self", instance)
```

### Concurrency System
**Model**: Go goroutines with global WaitGroup

**Architecture**:
```go
var wg sync.WaitGroup  // Global wait group

// r2(function) creates new goroutine
wg.Add(1)
go func() {
    defer wg.Done()
    userFunction.Call(args...)
}()

// At program end: wg.Wait()
```

**Current Limitations**:
- No communication between goroutines (channels)
- No synchronization primitives (mutexes)
- Basic error handling in goroutines

### Module System
**Pattern**: Static import with cache

**Import Algorithm**:
```go
1. Resolve relative/absolute path
2. Check import cache (avoid cycles)
3. Read file and parse
4. Create new environment for module
5. Evaluate module in its environment
6. Export symbols according to alias
```

**Path Resolution**:
```go
// Relative to current file
import "./utils.r2"

// Relative to base directory
import "lib/math.r2"

// With alias
import "./helpers.r2" as h
```

## Data Flow

### Execution Pipeline
```
1. main.go receives CLI arguments
2. Read .r2 file or start REPL
3. NewLexer(sourceCode)
4. NewParser(lexer)
5. parser.ParseProgram() → AST
6. NewEnvironment()
7. RegisterBuiltins(env)
8. ast.Eval(env) → Result
```

### Error Handling
**Strategy**: Panic/Recover with propagation

```go
// Lexer errors
panic("Unexpected character: " + ch)

// Parser errors  
panic("Expected ')' after expression")

// Runtime errors
panic("Undefined variable: " + name)

// Try-catch handles via recover()
defer func() {
    if r := recover(); r != nil {
        // Bind error to catch variable
        catchEnv.Set(errorVar, r)
        catchBlock.Eval(catchEnv)
    }
}()
```

### Memory Management
**Strategy**: Delegation to Go's GC

**Object Lifecycle**:
- Variables: Lifetime = environment scope
- Functions: Lifetime = closure + active references
- Objects: Lifetime = references in environments
- Arrays/Maps: Lifetime = active references

**Potential Memory Leaks**:
- Closures capturing large environments
- Reference cycles in objects
- Goroutines that don't terminate

## Extensibility

### Adding New Language Features

**1. New Tokens**:
```go
// In constants
const TOKEN_NEWFEATURE = "NEWFEATURE"

// In lexer
case "newkeyword":
    return Token{Type: TOKEN_NEWFEATURE, Value: literal}
```

**2. New AST Nodes**:
```go
type NewFeatureNode struct {
    Property string
    Child    Node
}

func (nfn *NewFeatureNode) Eval(env *Environment) interface{} {
    // Implement semantics
    return result
}
```

**3. Parsing Logic**:
```go
func (p *Parser) parseNewFeature() Node {
    // Parse specific syntax
    return &NewFeatureNode{...}
}
```

### Adding Native Libraries

**Template**:
```go
// r2lang/r2mynewlib.go
func RegisterMyNewLib(env *Environment) {
    env.Set("myFunction", BuiltinFunction(func(args ...interface{}) interface{} {
        // Implementation in Go
        return result
    }))
    
    env.Set("myObject", map[string]interface{}{
        "method1": BuiltinFunction(func(args ...interface{}) interface{} {
            // Method 1
            return result
        }),
        "property": "value",
    })
}

// In RunCode():
RegisterMyNewLib(env)
```

## Performance Considerations

### Current Bottlenecks

1. **Variable Lookup**: O(depth) in nested environments
2. **Function Calls**: New environment creation
3. **Type Checking**: Runtime type assertions
4. **AST Walking**: No hot path optimization
5. **String Operations**: Frequent concatenation

### Implementable Optimizations

1. **Variable Caching**: Global hash table for frequent variables
2. **Bytecode Compilation**: AST → Bytecode → Execution
3. **JIT Compilation**: Hot path optimization
4. **Constant Folding**: Evaluate constant expressions at parse time
5. **Tail Call Optimization**: For efficient recursion

## Comparison with Other Interpreters

### R2Lang Advantages
- **Simplicity**: Clear and understandable architecture
- **Extensibility**: Easy to add new libraries
- **Modernity**: Familiar syntax (JavaScript-like)
- **Integration**: Built-in testing and concurrency

### Limitations vs. Other Languages
- **Performance**: Slower than optimized interpreters (Python, Ruby)
- **Memory**: No optimized GC for the language
- **Ecosystem**: Limited standard library
- **Tooling**: No debugger, profiler, etc.

## Architectural Roadmap

### Short Term (1-3 months)
- Bytecode compilation layer
- Optimized variable lookup
- Improved error reporting

### Medium Term (3-6 months)  
- JIT compilation
- Specific garbage collector
- Expanded standard library

### Long Term (6-12 months)
- Language Server Protocol
- Multi-target compilation (WASM, native)
- Advanced optimization passes

## Conclusion

R2Lang's architecture prioritizes clarity and extensibility over raw performance, resulting in an interpreter that is easy to understand, modify, and extend. This decision facilitates rapid development of new features and learning the codebase, although it limits performance in intensive applications.

The modular design allows incremental optimizations without complete restructuring, positioning R2Lang for continuous evolution toward a robust production language.