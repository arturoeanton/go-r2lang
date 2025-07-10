# R2Lang Implementation Book - Guía Completa de Implementación

## Tabla de Contenidos

1. [Arquitectura General](#arquitectura-general)
2. [Lexer - Análisis Léxico](#lexer---análisis-léxico)
3. [Parser - Análisis Sintáctico](#parser---análisis-sintáctico)
4. [AST - Árbol de Sintaxis Abstracta](#ast---árbol-de-sintaxis-abstracta)
5. [Environment - Sistema de Entornos](#environment---sistema-de-entornos)
6. [Evaluación - Tree Walking Interpreter](#evaluación---tree-walking-interpreter)
7. [Sistema de Tipos](#sistema-de-tipos)
8. [Orientación a Objetos](#orientación-a-objetos)
9. [Concurrencia](#concurrencia)
10. [Sistema de Testing](#sistema-de-testing)
11. [Bibliotecas Nativas](#bibliotecas-nativas)
12. [Optimizaciones](#optimizaciones)

---

## Arquitectura General

### Visión de Alto Nivel

R2Lang implementa un intérprete **tree-walking** que procesa el código en las siguientes fases:

```
Código Fuente → Lexer → Tokens → Parser → AST → Evaluator → Resultado
```

**Ubicación**: `r2lang/r2lang.go`

### Estructura de Datos Principal

```go
type Environment struct {
    store    map[string]interface{}  // Variables y funciones
    outer    *Environment           // Environment padre (scoping)
    imported map[string]bool        // Cache de imports
    dir      string                 // Directorio actual
    CurrenFx string                // Función en ejecución
}
```

**Líneas**: `r2lang/r2lang.go:1429-1435`

La arquitectura se basa en environments anidados que permiten lexical scoping, donde cada nuevo scope (función, bloque, objeto) crea un nuevo environment que referencia al padre.

---

## Lexer - Análisis Léxico

### Propósito y Responsabilidades

El lexer convierte el texto fuente en una secuencia de tokens estructurados.

**Ubicación**: `r2lang/r2lang.go:72-321`

### Estructura del Lexer

```go
type Lexer struct {
    input        string  // Código fuente
    pos          int     // Posición actual
    col          int     // Columna actual
    line         int     // Línea actual
    length       int     // Longitud total
    currentToken Token   // Token actual
}

type Token struct {
    Type  string  // Tipo del token
    Value string  // Valor literal
    Line  int     // Línea donde aparece
    Pos   int     // Posición en el archivo
    Col   int     // Columna
}
```

### Tipos de Tokens Soportados

**Líneas**: `r2lang/r2lang.go:17-54`

```go
const (
    TOKEN_EOF     = "EOF"      // Fin de archivo
    TOKEN_NUMBER  = "NUMBER"   // 123, 45.67, -89.1
    TOKEN_STRING  = "STRING"   // "hola", 'mundo'
    TOKEN_IDENT   = "IDENT"    // identificadores
    TOKEN_ARROW   = "ARROW"    // =>
    TOKEN_SYMBOL  = "SYMBOL"   // (, ), {, }, [, ], ;, etc.
    
    // Palabras clave
    LET      = "let"
    FUNC     = "func"
    IF       = "if"
    WHILE    = "while"
    FOR      = "for"
    CLASS    = "class"
    EXTENDS  = "extends"
    // ... más keywords
)
```

### Algoritmo de Tokenización

**Método Principal**: `NextToken()` - `r2lang/r2lang.go:139-321`

1. **Skip Whitespace**: Salta espacios, tabs, newlines
2. **Comentarios**: Maneja `//` y `/* */`
3. **Números**: Reconoce enteros, decimales, negativos
4. **Strings**: Procesa comillas simples y dobles
5. **Operadores**: Detecta `==`, `<=`, `++`, `--`, etc.
6. **Identificadores**: Convierte keywords o nombres de variables

### Casos Especiales

#### Números con Signo
**Líneas**: `r2lang/r2lang.go:114-137`

```go
func (l *Lexer) parseNumberOrSign() Token {
    start := l.pos
    if l.input[l.pos] == '-' || l.input[l.pos] == '+' {
        l.nextch()
    }
    hasDigits := false
    // Procesar dígitos enteros
    for l.pos < l.length && isDigit(l.input[l.pos]) {
        hasDigits = true
        l.nextch()
    }
    // Procesar parte decimal
    if l.pos < l.length && l.input[l.pos] == '.' {
        l.nextch()
        for l.pos < l.length && isDigit(l.input[l.pos]) {
            hasDigits = true
            l.nextch()
        }
    }
    // Validación
    if !hasDigits {
        panic("Invalid number in " + l.input[start:l.pos])
    }
    val := l.input[start:l.pos]
    return Token{Type: TOKEN_NUMBER, Value: val, Line: l.line, Pos: l.pos, Col: l.col}
}
```

#### Detección Contextual de Signos
**Líneas**: `r2lang/r2lang.go:180-190`

El lexer determina si `-` o `+` son operadores binarios o parte de un número basándose en el contexto (qué carácter los precede).

---

## Parser - Análisis Sintáctico

### Método de Parsing

R2Lang implementa un **Recursive Descent Parser** con **precedencia de operadores**.

**Ubicación**: `r2lang/r2lang.go:1662-2331`

### Estructura del Parser

```go
type Parser struct {
    lexer   *Lexer  // Lexer para obtener tokens
    savTok  Token   // Token guardado
    prevTok Token   // Token anterior
    curTok  Token   // Token actual
    peekTok Token   // Siguiente token (lookahead)
    baseDir string  // Directorio base
}
```

### Gramática Implementada

#### Statements (Declaraciones)

**Líneas**: `r2lang/r2lang.go:1773-1820`

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

#### Expressions (Expresiones)

**Líneas**: `r2lang/r2lang.go:2128-2318`

```
expression := factor (binary_op factor)*

factor := NUMBER 
        | STRING
        | IDENTIFIER postfix*
        | "(" expression ")"
        | "[" expression_list "]"      // Array literal
        | "{" key_value_pairs "}"      // Map literal
        | "func" "(" args ")" block    // Anonymous function

postfix := "(" args ")"               // Function call
         | "." IDENTIFIER             // Property access
         | "[" expression "]"         // Index access
```

### Parsing de Construcciones Específicas

#### Clases con Herencia
**Líneas**: `r2lang/r2lang.go:2050-2090`

```go
func (p *Parser) parseObjectDeclaration() Node {
    p.nextToken() // "class" o "obj"
    objName := p.curTok.Value
    p.nextToken()
    
    parentName := ""
    if p.curTok.Value == EXTENDS {
        p.nextToken() // "extends"
        parentName = p.curTok.Value
        p.nextToken()
    }
    
    // Parsear miembros: { let x; func y() {...} }
    var members []Node
    for p.curTok.Value != "}" {
        if p.curTok.Value == LET || p.curTok.Value == VAR {
            members = append(members, p.parseLetStatement())
        } else if p.curTok.Value == FUNC {
            members = append(members, p.parseFunctionDeclaration())
        }
    }
    
    return &ObjectDeclaration{
        Name: objName, 
        Members: members, 
        ParentName: parentName
    }
}
```

#### For-In Loops
**Líneas**: `r2lang/r2lang.go:1991-2047`

```go
// Detecta: for (let item in collection)
if p.curTok.Value == LET {
    init = p.parseLetStatement()
    indexName := init.(*LetStatement).Name
    if p.curTok.Value == IN {
        p.nextToken()
        exp := p.parseExpression()
        body := p.parseBlockStatement()
        return &ForStatement{
            Init: exp, 
            Body: body, 
            inFlag: true, 
            inIndexName: indexName
        }
    }
}
```

#### TestCase Syntax
**Líneas**: `r2lang/r2lang.go:1711-1746`

```go
func (p *Parser) parseTestCase() Node {
    p.nextToken() // "TestCase"
    name := p.curTok.Value // String del nombre
    p.nextToken()
    
    var steps []TestStep
    for p.curTok.Value != "}" {
        var stepType string
        switch p.curTok.Type {
        case TOKEN_GIVEN, TOKEN_WHEN, TOKEN_THEN, TOKEN_AND:
            stepType = p.curTok.Value
            p.nextToken()
        }
        command := p.parseExpression()
        steps = append(steps, TestStep{Type: stepType, Command: command})
    }
    
    return &TestCase{Name: name, Steps: steps}
}
```

---

## AST - Árbol de Sintaxis Abstracta

### Interface Base

**Líneas**: `r2lang/r2lang.go:327-334`

```go
type Node interface {
    Eval(env *Environment) interface{}
}

type NodeTest interface {
    Eval(env *Environment) interface{}
    EvalStep(env *Environment) interface{}
}
```

Todos los nodos del AST implementan la interfaz `Node`, permitiendo evaluación polimórfica.

### Jerarquía de Nodos

#### Statements (Declaraciones)

**Program** - `r2lang/r2lang.go:433-449`
```go
type Program struct {
    Statements []Node
}

func (p *Program) Eval(env *Environment) interface{} {
    var result interface{}
    for _, stmt := range p.Statements {
        val := stmt.Eval(env)
        if rv, ok := val.(ReturnValue); ok {
            return rv.Value  // Early return
        }
        result = val
    }
    return result
}
```

**LetStatement** - `r2lang/r2lang.go:456-468`
```go
type LetStatement struct {
    Name  string
    Value Node
}

func (ls *LetStatement) Eval(env *Environment) interface{} {
    var val interface{}
    if ls.Value != nil {
        val = ls.Value.Eval(env)
    }
    env.Set(ls.Name, val)
    return nil
}
```

**FunctionDeclaration** - `r2lang/r2lang.go:756-772`
```go
type FunctionDeclaration struct {
    Name string
    Args []string
    Body *BlockStatement
}

func (fd *FunctionDeclaration) Eval(env *Environment) interface{} {
    fn := &UserFunction{
        Args:     fd.Args,
        Body:     fd.Body,
        Env:      env,      // Closure!
        IsMethod: false,
        code:     fd.Name,
    }
    env.Set(fd.Name, fn)
    return nil
}
```

#### Expressions (Expresiones)

**BinaryExpression** - `r2lang/r2lang.go:888-922`
```go
type BinaryExpression struct {
    Left  Node
    Op    string
    Right Node
}

func (be *BinaryExpression) Eval(env *Environment) interface{} {
    lv := be.Left.Eval(env)
    rv := be.Right.Eval(env)
    
    switch be.Op {
    case "+": return addValues(lv, rv)
    case "-": return subValues(lv, rv)
    case "*": return mulValues(lv, rv)
    case "/": return divValues(lv, rv)
    case "<": return toFloat(lv) < toFloat(rv)
    case ">": return toFloat(lv) > toFloat(rv)
    case "==": return equals(lv, rv)
    // ... más operadores
    default:
        panic("Unsupported binary operator: " + be.Op)
    }
}
```

**CallExpression** - `r2lang/r2lang.go:925-962`
```go
type CallExpression struct {
    Callee Node
    Args   []Node
}

func (ce *CallExpression) Eval(env *Environment) interface{} {
    calleeVal := ce.Callee.Eval(env)
    
    var argVals []interface{}
    for _, a := range ce.Args {
        argVals = append(argVals, a.Eval(env))
    }
    
    switch cv := calleeVal.(type) {
    case BuiltinFunction:
        return cv(argVals...)
    case *UserFunction:
        return cv.Call(argVals...)
    case map[string]interface{}:
        // Instanciar blueprint de clase
        return instantiateObject(env, cv, argVals)
    default:
        panic("Attempt to call non-function")
    }
}
```

---

## Environment - Sistema de Entornos

### Scoping Léxico

**Líneas**: `r2lang/r2lang.go:1429-1467`

El sistema de environments implementa **lexical scoping** mediante una cadena de environments padre-hijo.

```go
type Environment struct {
    store    map[string]interface{}  // Variables locales
    outer    *Environment           // Environment padre
    imported map[string]bool        // Cache de imports
    dir      string                 // Directorio de trabajo
    CurrenFx string                // Función actual (para debugging)
}

func NewInnerEnv(outer *Environment) *Environment {
    return &Environment{
        store:    make(map[string]interface{}),
        outer:    outer,                    // Referencia al padre
        imported: make(map[string]bool),
        dir:      outer.dir,
    }
}

func (e *Environment) Get(name string) (interface{}, bool) {
    val, ok := e.store[name]
    if ok {
        return val, true
    }
    if e.outer != nil {
        return e.outer.Get(name)  // Buscar en padre
    }
    return nil, false
}
```

### Closures

Las funciones capturan su environment de definición:

**Líneas**: `r2lang/r2lang.go:419-427`

```go
func (fl *FunctionLiteral) Eval(env *Environment) interface{} {
    fn := &UserFunction{
        Args:     fl.Args,
        Body:     fl.Body,
        Env:      env,  // ¡Captura el environment actual!
        IsMethod: false,
    }
    return fn
}
```

Cuando la función se ejecuta, crea un nuevo environment hijo del capturado:

**Líneas**: `r2lang/r2lang.go:1326-1371`

```go
func (uf *UserFunction) Call(args ...interface{}) interface{} {
    newEnv := NewInnerEnv(uf.Env)  // Nuevo scope basado en closure
    
    // Bind argumentos
    for i, param := range uf.Args {
        if i < len(args) {
            newEnv.Set(param, args[i])
        } else {
            newEnv.Set(param, nil)
        }
    }
    
    val := uf.Body.Eval(newEnv)
    if rv, ok := val.(ReturnValue); ok {
        return rv.Value
    }
    return val
}
```

---

## Evaluación - Tree Walking Interpreter

### Método de Evaluación

R2Lang usa **evaluación inmediata** mediante tree-walking. Cada nodo se evalúa recursivamente visitando sus hijos.

### Polimorfismo de Evaluación

El método `Eval(env *Environment) interface{}` es llamado polimórficamente:

```go
// En cualquier lugar del intérprete
result := node.Eval(environment)
```

### Manejo de Tipos Runtime

#### Conversiones Automáticas
**Líneas**: `r2lang/r2lang.go:1513-1581`

```go
func toFloat(val interface{}) float64 {
    switch v := val.(type) {
    case float64: return v
    case int:     return float64(v)
    case bool:    
        if v { return 1 }
        return 0
    case string:
        f, err := strconv.ParseFloat(v, 64)
        if err != nil {
            panic("Cannot convert string to number:" + v)
        }
        return f
    }
    panic("Cannot convert value to number")
}

func toBool(val interface{}) bool {
    if val == nil { return false }
    switch v := val.(type) {
    case bool:    return v
    case float64: return v != 0
    case int:     return v != 0
    case string:  return v != ""
    }
    return true  // Objetos y arrays son truthy
}
```

#### Operaciones Aritméticas Polimórficas
**Líneas**: `r2lang/r2lang.go:1583-1621`

```go
func addValues(a, b interface{}) interface{} {
    // Números
    if isNumeric(a) && isNumeric(b) {
        return toFloat(a) + toFloat(b)
    }
    
    // Concatenación de arrays
    if aa, ok := a.([]interface{}); ok {
        if bb, ok := b.([]interface{}); ok {
            return append(aa, bb...)  // [1,2] + [3,4] = [1,2,3,4]
        }
        return append(aa, b)          // [1,2] + 3 = [1,2,3]
    }
    
    // Concatenación de strings
    if sa, ok := a.(string); ok {
        return sa + fmt.Sprint(b)     // "hello" + 123 = "hello123"
    }
    if sb, ok := b.(string); ok {
        return fmt.Sprint(a) + sb     // 123 + "world" = "123world"
    }
    
    // Fallback a suma numérica
    return toFloat(a) + toFloat(b)
}
```

### Control de Flujo

#### Return Values
**Líneas**: `r2lang/r2lang.go:451-453`

```go
type ReturnValue struct {
    Value interface{}
}
```

Los returns se propagan mediante el tipo especial `ReturnValue`, que es detectado y unwrapped en cada nivel de la evaluación.

#### Try-Catch-Finally
**Líneas**: `r2lang/r2lang.go:609-623`

```go
func (ts *TryStatement) Eval(env *Environment) interface{} {
    defer func() {
        if r := recover(); r != nil {
            if ts.CatchBlock != nil {
                newEnv := NewInnerEnv(env)
                newEnv.Set(ts.ExceptionVar, r)  // Bind excepción
                ts.CatchBlock.Eval(newEnv)
            }
        }
        if ts.FinallyBlock != nil {
            ts.FinallyBlock.Eval(env)  // Siempre ejecuta
        }
    }()
    return ts.Body.Eval(env)
}
```

---

## Sistema de Tipos

### Tipos Nativos

R2Lang maneja los siguientes tipos en runtime:

1. **`float64`** - Números (enteros y decimales)
2. **`string`** - Cadenas de texto
3. **`bool`** - Booleanos (true/false)
4. **`[]interface{}`** - Arrays dinámicos
5. **`map[string]interface{}`** - Maps/Objetos
6. **`*UserFunction`** - Funciones definidas por usuario
7. **`BuiltinFunction`** - Funciones nativas
8. **`*ObjectInstance`** - Instancias de clases
9. **`nil`** - Valor nulo

### Type Assertions y Reflection

**Líneas**: `r2lang/r2std.go:14-20`

```go
env.Set("typeOf", BuiltinFunction(func(args ...interface{}) interface{} {
    if len(args) < 1 {
        return "nil"
    }
    val := args[0]
    return fmt.Sprintf("%T", val)  // Reflection de Go
}))
```

### Equality Comparison
**Líneas**: `r2lang/r2lang.go:1562-1581`

```go
func equals(a, b interface{}) bool {
    // Comparación numérica normalizada
    if isNumeric(a) && isNumeric(b) {
        return toFloat(a) == toFloat(b)
    }
    
    // Comparación por tipo
    switch aa := a.(type) {
    case string:
        if bb, ok := b.(string); ok {
            return aa == bb
        }
    case bool:
        if bb, ok := b.(bool); ok {
            return aa == bb
        }
    case nil:
        return b == nil
    }
    return false
}
```

---

## Orientación a Objetos

### Blueprints de Clases

Las clases en R2Lang son **blueprints** almacenados como `map[string]interface{}`.

**Líneas**: `r2lang/r2lang.go:775-824`

```go
func (od *ObjectDeclaration) Eval(env *Environment) interface{} {
    blueprint := make(map[string]interface{})
    
    // Herencia: copiar miembros del padre
    if od.ParentName != "" {
        if parent, ok := env.Get(od.ParentName); ok {
            blueprint["super"] = parent
        }
        raw, _ := env.Get(od.ParentName)
        if props, ok := raw.(map[string]interface{}); ok {
            for k, v := range props {
                if k == "ClassName" || k == "super" { continue }
                blueprint[k] = v  // Heredar propiedades
            }
        }
    }
    
    blueprint["ClassName"] = od.Name
    
    // Procesar miembros definidos
    for _, m := range od.Members {
        switch node := m.(type) {
        case *LetStatement:
            blueprint[node.Name] = nil  // Propiedad
        case *FunctionDeclaration:
            fn := &UserFunction{
                Args:     node.Args,
                Body:     node.Body,
                Env:      nil,     // Se asigna en instanciación
                IsMethod: true,
            }
            blueprint[node.Name] = fn
        }
    }
    
    env.Set(od.Name, blueprint)
    return nil
}
```

### Instanciación de Objetos

**Líneas**: `r2lang/r2lang.go:1396-1423`

```go
func instantiateObject(env *Environment, blueprint map[string]interface{}, argVals []interface{}) *ObjectInstance {
    objEnv := NewInnerEnv(env)
    instance := &ObjectInstance{Env: objEnv}
    
    // Copiar miembros del blueprint
    for k, v := range blueprint {
        switch vv := v.(type) {
        case *UserFunction:
            // Crear nueva función con environment de la instancia
            newFn := &UserFunction{
                Args:     vv.Args,
                Body:     vv.Body,
                Env:      objEnv,  // Bind al objeto
                IsMethod: true,
            }
            objEnv.Set(k, newFn)
        default:
            objEnv.Set(k, vv)
        }
    }
    
    // Self-reference
    objEnv.Set("self", instance)
    objEnv.Set("this", instance)
    
    // Llamar constructor si existe
    if constructor, ok := objEnv.Get("constructor"); ok {
        if constructorFn, isFn := constructor.(*UserFunction); isFn {
            constructorFn.Call(argVals...)
        }
    }
    
    return instance
}
```

### Method Binding y `this`

**Líneas**: `r2lang/r2lang.go:1326-1351`

```go
func (uf *UserFunction) NativeCall(currentEnv *Environment, args ...interface{}) interface{} {
    newEnv := currentEnv
    if newEnv == nil {
        newEnv = NewInnerEnv(uf.Env)
    }
    
    if uf.IsMethod {
        // Bind `this` y `self` para métodos
        if uf.Env != nil {
            if selfVal, ok := uf.Env.Get("self"); ok {
                newEnv.Set("self", selfVal)
                newEnv.Set("this", selfVal)
            }
        }
    }
    
    // Bind argumentos
    for i, param := range uf.Args {
        if i < len(args) {
            newEnv.Set(param, args[i])
        }
    }
    
    val := uf.Body.Eval(newEnv)
    return val
}
```

### Super Calls

**Líneas**: `r2lang/r2lang.go:931-940` y `1374-1380`

```go
// Detección de super calls en CallExpression
flagSuper := false
if o := ce.Callee; o != nil {
    if ae, ok := o.(*AccessExpression); ok {
        if id, ok := ae.Object.(*Identifier); ok {
            if id.Name == "super" {
                flagSuper = true
            }
        }
    }
}

// Ejecución especial para super
if flagSuper {
    return cv.SuperCall(env, argVals...)
}
```

---

## Concurrencia

### Goroutines R2

**Ubicación**: `r2lang/r2lib.go:17-38`

```go
env.Set("r2", BuiltinFunction(func(args ...interface{}) interface{} {
    if len(args) < 1 {
        panic("r2 need at least one function as argument")
    }
    
    fn, ok := args[0].(*UserFunction)
    if !ok {
        panic("r2 first argument must be a function")
    }
    
    wg.Add(1)  // WaitGroup global
    go func() {
        defer wg.Done()
        defer func() {
            if r := recover(); r != nil {
                fmt.Println("Error en goroutine:", r)
            }
        }()
        fn.Call(args[1:]...)  // Ejecutar función con argumentos restantes
    }()
    return nil
}))
```

### Sincronización Global

**Líneas**: `r2lang/r2lang.go:56-58` y `1471-1472`

```go
var (
    wg sync.WaitGroup  // WaitGroup global
)

func (e *Environment) Run(parser *Parser) (result interface{}) {
    defer wg.Wait()           // Esperar todas las goroutines
    wg = sync.WaitGroup{}     // Reset para nueva ejecución
    
    // ... evaluación del programa
}
```

### Primitivas de Concurrencia Adicionales

El sistema puede extenderse con:
- Channels (similar a Go)
- Mutexes
- Atomic operations
- Actor model

---

## Sistema de Testing

### Sintaxis TestCase

**Líneas**: `r2lang/r2lang.go:341-392`

```go
type TestCase struct {
    Name  string
    Steps []TestStep
}

type TestStep struct {
    Type    string // "Given", "When", "Then", "And"
    Command Node
}

func (tc *TestCase) Eval(env *Environment) interface{} {
    fmt.Printf("Executing Test Case: %s\n", tc.Name)
    var previousStepType string
    
    for _, step := range tc.Steps {
        stepType := step.Type
        if stepType == "And" {
            stepType = previousStepType  // And hereda el tipo previo
        } else {
            previousStepType = stepType
        }
        
        fmt.Printf("  %s: ", stepType)
        
        // Ejecutar comando del step
        if ce, ok := step.Command.(*CallExpression); ok {
            // Call directo
            result := ce.Eval(env)
            if result != nil {
                fmt.Println(result)
            }
        } else if fl, ok := step.Command.(*FunctionLiteral); ok {
            // Lambda function
            currentStep := fl.Eval(env).(*UserFunction)
            out := currentStep.CallStep(env)
            if out != nil {
                fmt.Println(out)
            }
        }
    }
    
    fmt.Println("Test Case Executed Successfully.")
    return nil
}
```

### Funciones de Testing

**Ubicación**: `r2lang/r2test.go`

```go
func RegisterTest(env *Environment) {
    env.Set("assertEqual", BuiltinFunction(func(args ...interface{}) interface{} {
        if len(args) < 2 {
            panic("assertEqual needs 2 arguments")
        }
        if !equals(args[0], args[1]) {
            panic(fmt.Sprintf("Assertion failed: %v != %v", args[0], args[1]))
        }
        return true
    }))
    
    env.Set("assertTrue", BuiltinFunction(func(args ...interface{}) interface{} {
        if len(args) < 1 {
            panic("assertTrue needs 1 argument")
        }
        if !toBool(args[0]) {
            panic(fmt.Sprintf("Assertion failed: %v is not true", args[0]))
        }
        return true
    }))
}
```

---

## Bibliotecas Nativas

### Arquitectura de Extensión

Cada biblioteca nativa se implementa como un archivo `r2*.go` con una función `Register*()`:

```go
func RegisterNombre(env *Environment) {
    env.Set("funcion1", BuiltinFunction(func(args ...interface{}) interface{} {
        // Implementación en Go
        return resultado
    }))
    
    env.Set("funcion2", BuiltinFunction(func(args ...interface{}) interface{} {
        // Otra función
        return resultado
    }))
}
```

### Ejemplo: HTTP Server
**Ubicación**: `r2lang/r2http.go`

```go
func RegisterHTTP(env *Environment) {
    env.Set("http", map[string]interface{}{
        "server": BuiltinFunction(func(args ...interface{}) interface{} {
            if len(args) < 2 {
                panic("http.server needs port and handler")
            }
            
            port := int(toFloat(args[0]))
            handler := args[1].(*UserFunction)
            
            http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
                // Crear objetos request/response para R2
                req := createR2Request(r)
                res := createR2Response(w)
                
                // Llamar handler de R2
                handler.Call(req, res)
            })
            
            fmt.Printf("Server running on port %d\n", port)
            log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
            return nil
        }),
    })
}
```

### Integración con el Runtime

**Líneas**: `r2lang/r2lang.go:2349-2364`

```go
func RunCode(filename string) {
    env := NewEnvironment()
    env.Set("true", true)
    env.Set("false", false)
    env.Set("nil", nil)
    
    // Registrar todas las bibliotecas
    RegisterLib(env)        // Core functions
    RegisterStd(env)        // Standard library
    RegisterIO(env)         // File I/O
    RegisterHTTPClient(env) // HTTP client
    RegisterString(env)     // String manipulation
    RegisterMath(env)       // Math functions
    RegisterRand(env)       // Random numbers
    RegisterTest(env)       // Testing utilities
    RegisterHTTP(env)       // HTTP server
    RegisterPrint(env)      // Pretty printing
    RegisterOS(env)         // OS interface
    RegisterHack(env)       // Debug utilities
    RegisterConcurrency(env)// Concurrency primitives
    RegisterCollections(env)// Array/Map operations
    
    parser := NewParser(code)
    env.Run(parser)
}
```

---

## Optimizaciones

### Optimizaciones Actuales

1. **Token Caching**: El parser mantiene `curTok` y `peekTok` para lookahead eficiente
2. **Environment Hierarchy**: Búsqueda de variables optimizada con early termination
3. **Type Assertions**: Uso de switch statements para type dispatch rápido

### Optimizaciones Futuras

1. **Bytecode Compilation**: Convertir AST a bytecode intermedio
2. **JIT Compilation**: Hot path optimization
3. **Constant Folding**: Evaluar expresiones constantes en parse time
4. **Dead Code Elimination**: Remover código inalcanzable
5. **Tail Call Optimization**: Optimizar recursión de cola

### Profiling Points

Para optimización futura, los puntos críticos son:

1. **Parser.parseExpression()** - Llamado frecuentemente
2. **Environment.Get()** - Variable lookup
3. **BinaryExpression.Eval()** - Operaciones aritméticas
4. **CallExpression.Eval()** - Function calls
5. **Array/Map operations** - Colecciones grandes

---

## Conclusión

La implementación de R2Lang demuestra un intérprete moderno y extensible que balances simplicidad con funcionalidad. La arquitectura modular permite fácil extensión mientras mantiene el core limpio y comprensible.

Las decisiones de diseño clave incluyen:

- **Tree-walking interpretation** para simplicidad
- **Lexical scoping** con environments anidados
- **Duck typing** con conversiones automáticas
- **Blueprint-based OOP** para flexibilidad
- **Native libraries** para extensibilidad

Esta implementación sirve como base sólida para futuras mejoras y optimizaciones descritas en el roadmap del proyecto.