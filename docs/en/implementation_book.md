# R2Lang Implementation Book - Complete Implementation Guide

## Table of Contents

1. [General Architecture](#general-architecture)
2. [Lexer - Lexical Analysis](#lexer---lexical-analysis)
3. [Parser - Syntactic Analysis](#parser---syntactic-analysis)
4. [AST - Abstract Syntax Tree](#ast---abstract-syntax-tree)
5. [Environment - Environment System](#environment---environment-system)
6. [Evaluation - Tree Walking Interpreter](#evaluation---tree-walking-interpreter)
7. [Type System](#type-system)
8. [Object-Oriented Programming](#object-oriented-programming)
9. [Concurrency](#concurrency)
10. [Testing System](#testing-system)
11. [Native Libraries](#native-libraries)
12. [Optimizations](#optimizations)

---

## General Architecture

### High-Level Overview

R2Lang implements a **tree-walking** interpreter that processes code in the following phases:

```
Source Code → Lexer → Tokens → Parser → AST → Evaluator → Result
```

**Location**: `r2lang/r2lang.go`

### Main Data Structure

```go
type Environment struct {
    store    map[string]interface{}  // Variables and functions
    outer    *Environment           // Parent environment (scoping)
    imported map[string]bool        // Import cache
    dir      string                 // Current directory
    CurrenFx string                // Function in execution
}
```

**Lines**: `r2lang/r2lang.go:1429-1435`

The architecture is based on nested environments that enable lexical scoping, where each new scope (function, block, object) creates a new environment that references the parent.

---

## Lexer - Lexical Analysis

### Purpose and Responsibilities

The lexer converts source text into a sequence of structured tokens.

**Location**: `r2lang/r2lang.go:72-321`

### Lexer Structure

```go
type Lexer struct {
    input        string  // Source code
    pos          int     // Current position
    col          int     // Current column
    line         int     // Current line
    length       int     // Total length
    currentToken Token   // Current token
}

type Token struct {
    Type  string  // Token type
    Value string  // Literal value
    Line  int     // Line where it appears
    Pos   int     // Position in file
    Col   int     // Column
}
```

### Supported Token Types

**Lines**: `r2lang/r2lang.go:17-54`

```go
const (
    TOKEN_EOF     = "EOF"      // End of file
    TOKEN_NUMBER  = "NUMBER"   // 123, 45.67, -89.1
    TOKEN_STRING  = "STRING"   // "hello", 'world'
    TOKEN_IDENT   = "IDENT"    // identifiers
    TOKEN_ARROW   = "ARROW"    // =>
    TOKEN_SYMBOL  = "SYMBOL"   // (, ), {, }, [, ], ;, etc.
    
    // Keywords
    LET      = "let"
    FUNC     = "func"
    IF       = "if"
    WHILE    = "while"
    FOR      = "for"
    CLASS    = "class"
    EXTENDS  = "extends"
    // ... more keywords
)
```

### Tokenization Algorithm

**Main Method**: `NextToken()` - `r2lang/r2lang.go:139-321`

1. **Skip Whitespace**: Skips spaces, tabs, newlines
2. **Comments**: Handles `//` and `/* */`
3. **Numbers**: Recognizes integers, decimals, negatives
4. **Strings**: Processes single and double quotes
5. **Operators**: Detects `==`, `<=`, `++`, `--`, etc.
6. **Identifiers**: Converts keywords or variable names

### Special Cases

#### Signed Numbers
**Lines**: `r2lang/r2lang.go:114-137`

```go
func (l *Lexer) parseNumberOrSign() Token {
    start := l.pos
    if l.input[l.pos] == '-' || l.input[l.pos] == '+' {
        l.nextch()
    }
    hasDigits := false
    // Process integer digits
    for l.pos < l.length && isDigit(l.input[l.pos]) {
        hasDigits = true
        l.nextch()
    }
    // Process decimal part
    if l.pos < l.length && l.input[l.pos] == '.' {
        l.nextch()
        for l.pos < l.length && isDigit(l.input[l.pos]) {
            hasDigits = true
            l.nextch()
        }
    }
    // Validation
    if !hasDigits {
        panic("Invalid number in " + l.input[start:l.pos])
    }
    val := l.input[start:l.pos]
    return Token{Type: TOKEN_NUMBER, Value: val, Line: l.line, Pos: l.pos, Col: l.col}
}
```

#### Contextual Sign Detection
**Lines**: `r2lang/r2lang.go:180-190`

The lexer determines if `-` or `+` are binary operators or part of a number based on context (what character precedes them).

---

## Parser - Syntactic Analysis

### Parsing Method

R2Lang implements a **Recursive Descent Parser** with **operator precedence**.

**Location**: `r2lang/r2lang.go:1662-2331`

### Parser Structure

```go
type Parser struct {
    lexer   *Lexer  // Lexer to get tokens
    savTok  Token   // Saved token
    prevTok Token   // Previous token
    curTok  Token   // Current token
    peekTok Token   // Next token (lookahead)
    baseDir string  // Base directory
}
```

### Implemented Grammar

#### Statements (Declarations)

**Lines**: `r2lang/r2lang.go:1773-1820`

```
statement := import_statement
           | testcase_statement  
           | try_statement
           | throw_statement
           | return_statement
           | let_statement
           | function_declaration
           | if_statement
           | while_statement
           | for_statement
           | object_declaration
           | assignment_or_expression
```

#### Expressions

**Lines**: `r2lang/r2lang.go:1950-2040`

```
expression := assignment_expression
            | ternary_expression
            | logical_or_expression
            | logical_and_expression
            | equality_expression
            | relational_expression
            | additive_expression
            | multiplicative_expression
            | unary_expression
            | postfix_expression
            | primary_expression
```

### Operator Precedence

**Lines**: `r2lang/r2lang.go:2041-2080`

R2Lang implements precedence climbing with the following hierarchy (highest to lowest):

1. **Primary**: `()`, literals, identifiers
2. **Postfix**: `[]`, `.`, `()`, `++`, `--`
3. **Unary**: `!`, `-`, `+`, `++`, `--`
4. **Multiplicative**: `*`, `/`, `%`
5. **Additive**: `+`, `-`
6. **Relational**: `<`, `<=`, `>`, `>=`
7. **Equality**: `==`, `!=`
8. **Logical AND**: `&&`
9. **Logical OR**: `||`
10. **Ternary**: `?:`
11. **Assignment**: `=`, `+=`, `-=`, etc.

### Key Parsing Functions

#### Let Statement
**Lines**: `r2lang/r2lang.go:1821-1839`

```go
func (p *Parser) parseLetStatement() *LetStatement {
    stmt := &LetStatement{Token: p.curTok}
    
    if !p.expectPeek(TOKEN_IDENT) {
        return nil
    }
    
    stmt.Name = &Identifier{Token: p.curTok, Value: p.curTok.Value}
    
    if !p.expectPeek("=") {
        return nil
    }
    
    p.nextToken()
    stmt.Value = p.parseExpression(LOWEST)
    
    // Skip optional semicolon
    if p.peekTokenIs(";") {
        p.nextToken()
    }
    
    return stmt
}
```

#### Function Declaration
**Lines**: `r2lang/r2lang.go:1840-1890`

```go
func (p *Parser) parseFunctionDeclaration() *FunctionDeclaration {
    stmt := &FunctionDeclaration{Token: p.curTok}
    
    if !p.expectPeek(TOKEN_IDENT) {
        return nil
    }
    
    stmt.Name = &Identifier{Token: p.curTok, Value: p.curTok.Value}
    
    if !p.expectPeek("(") {
        return nil
    }
    
    stmt.Parameters = p.parseFunctionParameters()
    
    if !p.expectPeek("{") {
        return nil
    }
    
    stmt.Body = p.parseBlockStatement()
    
    return stmt
}
```

#### Class Declaration
**Lines**: `r2lang/r2lang.go:2100-2180`

```go
func (p *Parser) parseClassStatement() *ClassStatement {
    stmt := &ClassStatement{Token: p.curTok}
    
    if !p.expectPeek(TOKEN_IDENT) {
        return nil
    }
    
    stmt.Name = &Identifier{Token: p.curTok, Value: p.curTok.Value}
    
    // Check for inheritance
    if p.peekTokenIs(EXTENDS) {
        p.nextToken()
        if !p.expectPeek(TOKEN_IDENT) {
            return nil
        }
        stmt.SuperClass = &Identifier{Token: p.curTok, Value: p.curTok.Value}
    }
    
    if !p.expectPeek("{") {
        return nil
    }
    
    stmt.Body = p.parseClassBody()
    
    return stmt
}
```

---

## AST - Abstract Syntax Tree

### Node Interface

**Lines**: `r2lang/r2lang.go:327-334`

```go
type Node interface {
    Eval(env *Environment) interface{}
    String() string
}

type Statement interface {
    Node
    statementNode()
}

type Expression interface {
    Node
    expressionNode()
}
```

All AST nodes implement the `Node` interface with two key methods:
- `Eval(env *Environment) interface{}`: Executes the node
- `String() string`: Returns string representation

### Main AST Nodes

#### Expressions

**Identifier** - `r2lang/r2lang.go:335-349`
```go
type Identifier struct {
    Token Token
    Value string
}

func (i *Identifier) Eval(env *Environment) interface{} {
    val, ok := env.Get(i.Value)
    if !ok {
        panic("identifier not found: " + i.Value)
    }
    return val
}
```

**IntegerLiteral** - `r2lang/r2lang.go:350-370`
```go
type IntegerLiteral struct {
    Token Token
    Value int64
}

func (il *IntegerLiteral) Eval(env *Environment) interface{} {
    return float64(il.Value)
}
```

**InfixExpression** - `r2lang/r2lang.go:450-520`
```go
type InfixExpression struct {
    Token    Token
    Left     Expression
    Operator string
    Right    Expression
}

func (ie *InfixExpression) Eval(env *Environment) interface{} {
    left := ie.Left.Eval(env)
    right := ie.Right.Eval(env)
    
    return evalInfixExpression(ie.Operator, left, right)
}
```

#### Statements

**LetStatement** - `r2lang/r2lang.go:400-420`
```go
type LetStatement struct {
    Token Token
    Name  *Identifier
    Value Expression
}

func (ls *LetStatement) Eval(env *Environment) interface{} {
    val := ls.Value.Eval(env)
    env.Set(ls.Name.Value, val)
    return val
}
```

**FunctionDeclaration** - `r2lang/r2lang.go:580-620`
```go
type FunctionDeclaration struct {
    Token      Token
    Name       *Identifier
    Parameters []*Identifier
    Body       *BlockStatement
}

func (fd *FunctionDeclaration) Eval(env *Environment) interface{} {
    fn := &Function{
        Parameters: fd.Parameters,
        Body:       fd.Body,
        Env:        env,
    }
    env.Set(fd.Name.Value, fn)
    return fn
}
```

**ClassStatement** - `r2lang/r2lang.go:700-780`
```go
type ClassStatement struct {
    Token      Token
    Name       *Identifier
    SuperClass *Identifier
    Body       *BlockStatement
}

func (cs *ClassStatement) Eval(env *Environment) interface{} {
    classEnv := NewEnclosedEnvironment(env)
    
    // Execute class body to collect methods
    cs.Body.Eval(classEnv)
    
    class := &Class{
        Name:       cs.Name.Value,
        Methods:    classEnv.store,
        SuperClass: nil,
        Env:        env,
    }
    
    // Handle inheritance
    if cs.SuperClass != nil {
        superClass, ok := env.Get(cs.SuperClass.Value)
        if !ok {
            panic("superclass not found: " + cs.SuperClass.Value)
        }
        class.SuperClass = superClass.(*Class)
    }
    
    env.Set(cs.Name.Value, class)
    return class
}
```

---

## Environment - Environment System

### Environment Structure

**Lines**: `r2lang/r2lang.go:1429-1507`

The environment system implements lexical scoping through a chain of nested environments.

```go
type Environment struct {
    store    map[string]interface{}
    outer    *Environment  // Parent environment
    imported map[string]bool
    dir      string
    CurrenFx string
}

func NewEnvironment() *Environment {
    return &Environment{
        store:    make(map[string]interface{}),
        outer:    nil,
        imported: make(map[string]bool),
        dir:      "",
        CurrenFx: "",
    }
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
    env := NewEnvironment()
    env.outer = outer
    env.dir = outer.dir
    return env
}
```

### Key Methods

#### Get - Variable Resolution
```go
func (e *Environment) Get(name string) (interface{}, bool) {
    value, ok := e.store[name]
    if !ok && e.outer != nil {
        value, ok = e.outer.Get(name)
    }
    return value, ok
}
```

Variable resolution follows the scope chain: first check current environment, then parent environments recursively.

#### Set - Variable Assignment
```go
func (e *Environment) Set(name string, val interface{}) interface{} {
    e.store[name] = val
    return val
}
```

Variables are always set in the current environment, creating new bindings or shadowing outer ones.

#### SetOuter - Update in Outer Scope
```go
func (e *Environment) SetOuter(name string, val interface{}) interface{} {
    _, ok := e.store[name]
    if ok {
        e.store[name] = val
        return val
    }
    if e.outer != nil {
        return e.outer.SetOuter(name, val)
    }
    e.store[name] = val
    return val
}
```

Used for assignment operations to existing variables in outer scopes.

---

## Evaluation - Tree Walking Interpreter

### Evaluation Strategy

R2Lang uses a **tree-walking interpreter** where each AST node implements an `Eval()` method that recursively evaluates child nodes.

### Built-in Type System

**Lines**: `r2lang/r2lang.go:55-71`

```go
type Function struct {
    Parameters []*Identifier
    Body       *BlockStatement
    Env        *Environment
}

type Class struct {
    Name       string
    Methods    map[string]interface{}
    SuperClass *Class
    Env        *Environment
}

type ObjectInstance struct {
    Class      *Class
    Properties map[string]interface{}
}

type BuiltinFunction func(args ...interface{}) interface{}
```

### Expression Evaluation

#### Infix Operations
**Lines**: `r2lang/r2lang.go:1200-1350`

```go
func evalInfixExpression(operator string, left, right interface{}) interface{} {
    switch {
    case leftType == "FLOAT" && rightType == "FLOAT":
        return evalFloatInfixExpression(operator, left, right)
    case leftType == "STRING" && rightType == "STRING":
        return evalStringInfixExpression(operator, left, right)
    case leftType == "BOOLEAN" && rightType == "BOOLEAN":
        return evalBooleanInfixExpression(operator, left, right)
    default:
        return evalMixedInfixExpression(operator, left, right)
    }
}
```

#### Function Calls
**Lines**: `r2lang/r2lang.go:950-1050`

```go
func (ce *CallExpression) Eval(env *Environment) interface{} {
    function := ce.Function.Eval(env)
    
    args := []interface{}{}
    for _, arg := range ce.Arguments {
        args = append(args, arg.Eval(env))
    }
    
    return applyFunction(function, args, env)
}

func applyFunction(fn interface{}, args []interface{}, env *Environment) interface{} {
    switch function := fn.(type) {
    case *Function:
        extendedEnv := extendFunctionEnv(function, args)
        evaluated := function.Body.Eval(extendedEnv)
        return unwrapReturnValue(evaluated)
    case BuiltinFunction:
        return function(args...)
    default:
        panic("not a function: " + fmt.Sprintf("%T", fn))
    }
}
```

### Statement Execution

#### Block Statements
```go
func (bs *BlockStatement) Eval(env *Environment) interface{} {
    var result interface{}
    
    for _, statement := range bs.Statements {
        result = statement.Eval(env)
        
        if returnValue, ok := result.(*ReturnValue); ok {
            return returnValue
        }
    }
    
    return result
}
```

#### Control Flow
```go
func (ifs *IfStatement) Eval(env *Environment) interface{} {
    condition := ifs.Condition.Eval(env)
    
    if isTruthy(condition) {
        return ifs.Consequence.Eval(env)
    } else if ifs.Alternative != nil {
        return ifs.Alternative.Eval(env)
    } else {
        return nil
    }
}
```

---

## Type System

### Runtime Type System

R2Lang implements a dynamic type system with runtime type checking.

#### Supported Types
- **float64**: All numbers (integers and decimals)
- **string**: Text values
- **bool**: Boolean values
- **Array**: Ordered collections
- **Map**: Key-value pairs (objects)
- **Function**: Callable functions
- **Class**: Class definitions
- **ObjectInstance**: Object instances

#### Type Checking
**Lines**: `r2lang/r2lang.go:1100-1150`

```go
func getType(obj interface{}) string {
    switch obj.(type) {
    case float64:
        return "FLOAT"
    case string:
        return "STRING"
    case bool:
        return "BOOLEAN"
    case []interface{}:
        return "ARRAY"
    case map[string]interface{}:
        return "MAP"
    case *Function:
        return "FUNCTION"
    case *Class:
        return "CLASS"
    case *ObjectInstance:
        return "OBJECT"
    case BuiltinFunction:
        return "BUILTIN"
    default:
        return "UNKNOWN"
    }
}
```

#### Type Conversions
**Lines**: `r2lang/r2lang.go:1151-1199`

```go
func convertToFloat(obj interface{}) (float64, bool) {
    switch v := obj.(type) {
    case float64:
        return v, true
    case string:
        if f, err := strconv.ParseFloat(v, 64); err == nil {
            return f, true
        }
    case bool:
        if v {
            return 1.0, true
        }
        return 0.0, true
    }
    return 0, false
}

func convertToString(obj interface{}) string {
    switch v := obj.(type) {
    case string:
        return v
    case float64:
        if v == float64(int(v)) {
            return strconv.Itoa(int(v))
        }
        return strconv.FormatFloat(v, 'f', -1, 64)
    case bool:
        return strconv.FormatBool(v)
    default:
        return fmt.Sprintf("%v", v)
    }
}
```

---

## Object-Oriented Programming

### Class System Implementation

R2Lang supports full object-oriented programming with classes, inheritance, and method overriding.

#### Class Definition
**Lines**: `r2lang/r2lang.go:700-780`

Classes are defined with the `class` keyword and can extend other classes:

```r2
class Animal {
    let name
    
    constructor(name) {
        this.name = name
    }
    
    speak() {
        print(this.name + " makes a sound")
    }
}

class Dog extends Animal {
    constructor(name) {
        super.constructor(name)
    }
    
    speak() {
        print(this.name + " barks")
    }
}
```

#### Constructor Handling
**Lines**: `r2lang/r2lang.go:1350-1400`

```go
func createObjectInstance(class *Class, args []interface{}, env *Environment) *ObjectInstance {
    instance := &ObjectInstance{
        Class:      class,
        Properties: make(map[string]interface{}),
    }
    
    // Set 'this' reference
    instanceEnv := NewEnclosedEnvironment(env)
    instanceEnv.Set("this", instance)
    
    // Call constructor if exists
    if constructor, ok := class.Methods["constructor"]; ok {
        if fn, ok := constructor.(*Function); ok {
            extendedEnv := extendFunctionEnv(fn, args)
            extendedEnv.Set("this", instance)
            fn.Body.Eval(extendedEnv)
        }
    }
    
    return instance
}
```

#### Method Resolution
**Lines**: `r2lang/r2lang.go:800-850`

```go
func (oi *ObjectInstance) GetMethod(name string) interface{} {
    // Check instance methods first
    if method, ok := oi.Class.Methods[name]; ok {
        return method
    }
    
    // Check superclass chain
    for super := oi.Class.SuperClass; super != nil; super = super.SuperClass {
        if method, ok := super.Methods[name]; ok {
            return method
        }
    }
    
    return nil
}
```

#### Property Access
```go
func (ae *AccessExpression) Eval(env *Environment) interface{} {
    object := ae.Object.Eval(env)
    
    switch obj := object.(type) {
    case *ObjectInstance:
        if ae.Property.Value == "constructor" {
            return obj.GetMethod("constructor")
        }
        
        // Check properties first
        if prop, ok := obj.Properties[ae.Property.Value]; ok {
            return prop
        }
        
        // Then check methods
        return obj.GetMethod(ae.Property.Value)
        
    case map[string]interface{}:
        return obj[ae.Property.Value]
        
    default:
        panic("invalid property access")
    }
}
```

---

## Concurrency

### Goroutine Implementation

R2Lang provides native concurrency through the `r2()` function, which creates goroutines.

**Lines**: `r2std.go:150-180`

```go
func r2Function(args ...interface{}) interface{} {
    if len(args) < 1 {
        panic("r2 requires at least one argument")
    }
    
    fn := args[0]
    fnArgs := args[1:]
    
    go func() {
        switch f := fn.(type) {
        case *Function:
            applyFunction(f, fnArgs, f.Env)
        case BuiltinFunction:
            f(fnArgs...)
        default:
            panic("r2: first argument must be a function")
        }
    }()
    
    return nil
}
```

### Usage Example

```r2
func worker(id) {
    for (let i = 0; i < 5; i++) {
        print("Worker " + id + " - iteration " + i)
        sleep(1)
    }
}

func main() {
    r2(worker, "A")
    r2(worker, "B")
    r2(worker, "C")
    sleep(6)  // Wait for workers to complete
}
```

### Thread Safety Considerations

Currently, R2Lang relies on Go's runtime for goroutine management and memory safety. However, shared state access between goroutines is not protected by mutexes, which means careful programming is required to avoid race conditions.

---

## Testing System

### BDD-Style Testing

R2Lang includes a built-in testing system with Behavior-Driven Development (BDD) syntax.

#### TestCase Structure
**Lines**: `r2lang/r2lang.go:2200-2280`

```go
type TestCase struct {
    Token       Token
    Name        string
    Given       *FunctionLiteral
    When        *FunctionLiteral
    Then        *FunctionLiteral
    And         *FunctionLiteral
}

func (tc *TestCase) Eval(env *Environment) interface{} {
    fmt.Printf("Running TestCase: %s\\n", tc.Name)
    
    // Execute Given
    if tc.Given != nil {
        givenResult := tc.Given.Eval(env)
        fmt.Printf("Given: %v\\n", givenResult)
    }
    
    // Execute When
    if tc.When != nil {
        whenResult := tc.When.Eval(env)
        fmt.Printf("When: %v\\n", whenResult)
    }
    
    // Execute Then
    if tc.Then != nil {
        thenResult := tc.Then.Eval(env)
        fmt.Printf("Then: %v\\n", thenResult)
    }
    
    // Execute And
    if tc.And != nil {
        andResult := tc.And.Eval(env)
        fmt.Printf("And: %v\\n", andResult)
    }
    
    return nil
}
```

#### Testing Functions
Built-in assertion functions are provided for testing:

```go
func assertEqualFunction(args ...interface{}) interface{} {
    if len(args) != 2 {
        panic("assertEqual requires exactly 2 arguments")
    }
    
    actual := args[0]
    expected := args[1]
    
    if !reflect.DeepEqual(actual, expected) {
        panic(fmt.Sprintf("Assertion failed: expected %v, got %v", expected, actual))
    }
    
    return true
}

func assertTrueFunction(args ...interface{}) interface{} {
    if len(args) != 1 {
        panic("assertTrue requires exactly 1 argument")
    }
    
    if !isTruthy(args[0]) {
        panic(fmt.Sprintf("Assertion failed: expected true, got %v", args[0]))
    }
    
    return true
}
```

---

## Native Libraries

### Library Registration System

Native libraries are implemented in Go and registered during environment initialization.

**Location**: `r2lang/r2*.go` files

### Standard Library (`r2std.go`)

```go
func RegisterStd(env *Environment) {
    env.Set("print", BuiltinFunction(printFunction))
    env.Set("len", BuiltinFunction(lenFunction))
    env.Set("typeOf", BuiltinFunction(typeOfFunction))
    env.Set("sleep", BuiltinFunction(sleepFunction))
    env.Set("assertEqual", BuiltinFunction(assertEqualFunction))
    env.Set("assertTrue", BuiltinFunction(assertTrueFunction))
    env.Set("r2", BuiltinFunction(r2Function))
}
```

### Math Library (`r2math.go`)

```go
func RegisterMath(env *Environment) {
    mathObj := map[string]interface{}{
        "sqrt": BuiltinFunction(sqrtFunction),
        "pow":  BuiltinFunction(powFunction),
        "sin":  BuiltinFunction(sinFunction),
        "cos":  BuiltinFunction(cosFunction),
        "tan":  BuiltinFunction(tanFunction),
        "abs":  BuiltinFunction(absFunction),
        "pi":   math.Pi,
        "e":    math.E,
    }
    env.Set("math", mathObj)
}
```

### HTTP Library (`r2http.go`)

```go
func RegisterHTTP(env *Environment) {
    httpObj := map[string]interface{}{
        "get":    BuiltinFunction(httpGetFunction),
        "post":   BuiltinFunction(httpPostFunction),
        "server": BuiltinFunction(httpServerFunction),
    }
    env.Set("http", httpObj)
    
    httpClientObj := map[string]interface{}{
        "get":  BuiltinFunction(httpClientGetFunction),
        "post": BuiltinFunction(httpClientPostFunction),
    }
    env.Set("httpClient", httpClientObj)
}
```

### I/O Library (`r2io.go`)

```go
func RegisterIO(env *Environment) {
    ioObj := map[string]interface{}{
        "readFile":  BuiltinFunction(readFileFunction),
        "writeFile": BuiltinFunction(writeFileFunction),
        "exists":    BuiltinFunction(fileExistsFunction),
    }
    env.Set("io", ioObj)
}
```

---

## Optimizations

### Current Optimizations

1. **Token Lookahead**: Parser uses one-token lookahead for efficient parsing
2. **Environment Caching**: Import statements cache loaded modules
3. **String Interning**: Common strings reuse memory
4. **Tail Call Recognition**: Basic return statement optimization

### Performance Characteristics

- **Lexing**: O(n) linear time complexity
- **Parsing**: O(n) for most constructs, O(n²) worst case for deeply nested expressions
- **Evaluation**: O(n) per AST node visit
- **Memory**: Proportional to source code size and runtime objects

### Future Optimization Opportunities

1. **Bytecode Compilation**: Convert AST to bytecode for faster execution
2. **JIT Compilation**: Compile hot code paths to native code
3. **Constant Folding**: Evaluate constant expressions at parse time
4. **Dead Code Elimination**: Remove unreachable code
5. **Type Inference**: Static type analysis for better optimization
6. **Memory Pooling**: Reuse object allocations
7. **Inline Caching**: Cache method lookups and property access

---

## Conclusion

R2Lang's implementation demonstrates a clean, extensible architecture that balances simplicity with power. The tree-walking interpreter provides excellent debuggability and extensibility, while the modular library system allows for easy addition of new functionality.

The codebase is well-structured for further development, with clear separation between lexical analysis, parsing, AST representation, and evaluation. This foundation supports the planned improvements outlined in the roadmap, including advanced type systems, performance optimizations, and ecosystem enhancements.