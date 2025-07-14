# Roadmap de Performance - R2Lang

## Resumen Ejecutivo

Este documento presenta un plan detallado para mejorar el rendimiento del intérprete R2Lang, basado en el análisis de benchmarks y profiling del código. Se incluyen bugs críticos, optimizaciones y mejoras arquitecturales con estimaciones realistas de tiempo e impacto.

**Estado actual**: 122,670 ns/op en operaciones básicas, 85,297 B/op memoria  
**Meta objetivo**: <15,000 ns/op (mejora 8x), <10,000 B/op (mejora 8.5x)

---

## 🐛 BUGS CRÍTICOS DE PERFORMANCE

### BUG-001: Evaluación Eagerly en Binary Expressions
- **Archivo**: `pkg/r2core/binary_expression.go:9-11`
- **Criticidad**: 🔴 CRÍTICA
- **Complejidad**: 🟢 BAJA
- **Estimación**: 2 horas
- **Impacto Performance**: 2.5x mejora en expresiones lógicas
- **Impacto Memoria**: 1.5x reducción en evaluaciones

**Descripción del problema:**
```go
// CÓDIGO ACTUAL PROBLEMÁTICO
func (be *BinaryExpression) Eval(env *Environment) interface{} {
    lv := be.Left.Eval(env)  // Siempre evalúa left
    rv := be.Right.Eval(env) // Siempre evalúa right (PROBLEMA)
    
    switch be.Op {
    case "&&":
        return toBool(lv) && toBool(rv) // Si lv es false, rv es innecesario
    case "||":
        return toBool(lv) || toBool(rv) // Si lv es true, rv es innecesario
    }
}
```

**Solución propuesta:**
```go
// CÓDIGO OPTIMIZADO
func (be *BinaryExpression) Eval(env *Environment) interface{} {
    lv := be.Left.Eval(env)
    
    // Evaluación lazy para operadores lógicos
    switch be.Op {
    case "&&":
        if !toBool(lv) {
            return false // No evaluar right si left es false
        }
        rv := be.Right.Eval(env)
        return toBool(rv)
    case "||":
        if toBool(lv) {
            return true // No evaluar right si left es true
        }
        rv := be.Right.Eval(env)
        return toBool(rv)
    default:
        // Para operadores aritméticos, evaluar ambos
        rv := be.Right.Eval(env)
        return be.evaluateArithmeticOp(lv, rv)
    }
}
```

### BUG-002: Búsqueda Linear en Environment.Get()
- **Archivo**: `pkg/r2core/environment.go:49-58`
- **Criticidad**: 🔴 CRÍTICA
- **Complejidad**: 🟡 MEDIA
- **Estimación**: 4 horas
- **Impacto Performance**: 3x mejora en acceso a variables
- **Impacto Memoria**: 2x reducción en allocaciones

**Descripción del problema:**
```go
// CÓDIGO ACTUAL PROBLEMÁTICO
func (e *Environment) Get(name string) (interface{}, bool) {
    val, ok := e.store[name]
    if ok {
        return val, true
    }
    if e.outer != nil {
        return e.outer.Get(name) // Recursión costosa O(n)
    }
    return nil, false
}
```

**Problema**: En cada acceso a variable, se busca recursivamente en toda la cadena de environments. Para variables frecuentemente accedidas, esto es muy ineficiente.

**Solución propuesta:**
```go
// CÓDIGO OPTIMIZADO
type Environment struct {
    store    map[string]interface{}
    outer    *Environment
    imported map[string]bool
    Dir      string
    CurrenFx string
    
    // NUEVO: Cache para variables frecuentemente accedidas
    cache    map[string]interface{}
    cacheMu  sync.RWMutex
}

func (e *Environment) Get(name string) (interface{}, bool) {
    // Primero buscar en cache local
    e.cacheMu.RLock()
    if val, ok := e.cache[name]; ok {
        e.cacheMu.RUnlock()
        return val, true
    }
    e.cacheMu.RUnlock()
    
    // Búsqueda normal
    val, ok := e.store[name]
    if ok {
        // Agregar al cache si es accedida frecuentemente
        e.cacheMu.Lock()
        if e.cache == nil {
            e.cache = make(map[string]interface{})
        }
        e.cache[name] = val
        e.cacheMu.Unlock()
        return val, true
    }
    
    if e.outer != nil {
        return e.outer.Get(name)
    }
    return nil, false
}
```

### BUG-003: Conversión de Tipos Repetitiva
- **Archivo**: `pkg/r2core/commons.go:8-29`
- **Criticidad**: 🟡 ALTA
- **Complejidad**: 🟡 MEDIA
- **Estimación**: 6 horas
- **Impacto Performance**: 4x mejora en operaciones numéricas
- **Impacto Memoria**: 3x reducción en allocaciones

**Descripción del problema:**
```go
// CÓDIGO ACTUAL PROBLEMÁTICO
func toFloat(val interface{}) float64 {
    switch v := val.(type) {
    case float64:
        return v
    case int:
        return float64(v) // Conversión costosa cada vez
    case string:
        f, err := strconv.ParseFloat(v, 64) // Parsing repetitivo
        if err != nil {
            panic("Cannot convert string to number:" + v)
        }
        return f
    }
    panic("Cannot convert value to number")
}
```

**Problema**: Se convierte el mismo valor múltiples veces. Los strings numéricos se parsean repetidamente.

**Solución propuesta:**
```go
// CÓDIGO OPTIMIZADO
var (
    // Cache para conversiones de string a float
    stringToFloatCache = make(map[string]float64)
    stringCacheMu      sync.RWMutex
    
    // Cache para números pequeños comunes
    intToFloatCache = make(map[int]float64)
)

func init() {
    // Pre-poblar cache con números comunes
    for i := -1000; i <= 1000; i++ {
        intToFloatCache[i] = float64(i)
    }
}

func toFloat(val interface{}) float64 {
    switch v := val.(type) {
    case float64:
        return v
    case int:
        // Usar cache para números pequeños
        if cached, ok := intToFloatCache[v]; ok {
            return cached
        }
        return float64(v)
    case string:
        // Buscar en cache primero
        stringCacheMu.RLock()
        if cached, ok := stringToFloatCache[v]; ok {
            stringCacheMu.RUnlock()
            return cached
        }
        stringCacheMu.RUnlock()
        
        // Parsear y cachear
        f, err := strconv.ParseFloat(v, 64)
        if err != nil {
            panic("Cannot convert string to number:" + v)
        }
        
        stringCacheMu.Lock()
        stringToFloatCache[v] = f
        stringCacheMu.Unlock()
        return f
    }
    panic("Cannot convert value to number")
}
```

---

## 🚀 TAREAS DE PERFORMANCE

### PERF-001: Implementar Object Pool para Números
- **Archivo**: `pkg/r2core/literals.go`
- **Criticidad**: 🔴 CRÍTICA
- **Complejidad**: 🟡 MEDIA
- **Estimación**: 8 horas
- **Impacto Performance**: 6x mejora en operaciones numéricas
- **Impacto Memoria**: 10x reducción en allocaciones

**Descripción:**
Implementar un pool de objetos para reutilizar instancias de números frecuentemente utilizados, reduciendo dramáticamente las allocaciones de memoria.

**Código actual problemático:**
```go
// pkg/r2core/literals.go
type NumberLiteral struct {
    Value float64
}

func (n *NumberLiteral) Eval(env *Environment) interface{} {
    return n.Value // Esto crea una nueva instancia cada vez
}
```

**Solución propuesta:**
```go
// pkg/r2core/literals.go
import "sync"

var (
    // Pool para reutilizar wrappers de números
    numberPool = sync.Pool{
        New: func() interface{} {
            return &NumberWrapper{}
        },
    }
    
    // Cache para números pequeños (más eficiente que pool)
    smallNumberCache = make(map[float64]*NumberWrapper)
    cacheMu          sync.RWMutex
)

type NumberWrapper struct {
    Value float64
}

func init() {
    // Pre-poblar cache con números comunes (-1000 a 1000)
    for i := -1000; i <= 1000; i++ {
        smallNumberCache[float64(i)] = &NumberWrapper{Value: float64(i)}
    }
}

func GetNumber(value float64) *NumberWrapper {
    // Para números pequeños, usar cache
    if value >= -1000 && value <= 1000 && value == float64(int(value)) {
        cacheMu.RLock()
        if wrapper, ok := smallNumberCache[value]; ok {
            cacheMu.RUnlock()
            return wrapper
        }
        cacheMu.RUnlock()
    }
    
    // Para números grandes, usar pool
    wrapper := numberPool.Get().(*NumberWrapper)
    wrapper.Value = value
    return wrapper
}

func (n *NumberLiteral) Eval(env *Environment) interface{} {
    return GetNumber(n.Value).Value
}
```

### PERF-002: Optimizar String Concatenation
- **Archivo**: `pkg/r2core/commons.go:95-102`
- **Criticidad**: 🟡 ALTA
- **Complejidad**: 🟡 MEDIA
- **Estimación**: 6 horas
- **Impacto Performance**: 3x mejora en operaciones string
- **Impacto Memoria**: 4x reducción en allocaciones

**Descripción:**
Optimizar la concatenación de strings usando string builder para evitar múltiples allocaciones.

**Código actual problemático:**
```go
// pkg/r2core/commons.go
func addValues(a, b interface{}) interface{} {
    // Si uno es string => concatenar
    if sa, ok := a.(string); ok {
        return sa + fmt.Sprint(b) // Múltiples allocaciones
    }
    if sb, ok := b.(string); ok {
        return fmt.Sprint(a) + sb // Múltiples allocaciones
    }
    return toFloat(a) + toFloat(b)
}
```

**Solución propuesta:**
```go
// pkg/r2core/commons.go
import "strings"

var stringBuilderPool = sync.Pool{
    New: func() interface{} {
        return &strings.Builder{}
    },
}

func addValues(a, b interface{}) interface{} {
    // Si uno es string => concatenar eficientemente
    if sa, ok := a.(string); ok {
        builder := stringBuilderPool.Get().(*strings.Builder)
        defer stringBuilderPool.Put(builder)
        builder.Reset()
        
        builder.WriteString(sa)
        builder.WriteString(toString(b))
        return builder.String()
    }
    if sb, ok := b.(string); ok {
        builder := stringBuilderPool.Get().(*strings.Builder)
        defer stringBuilderPool.Put(builder)
        builder.Reset()
        
        builder.WriteString(toString(a))
        builder.WriteString(sb)
        return builder.String()
    }
    return toFloat(a) + toFloat(b)
}

func toString(val interface{}) string {
    switch v := val.(type) {
    case string:
        return v
    case float64:
        return strconv.FormatFloat(v, 'f', -1, 64)
    case int:
        return strconv.Itoa(v)
    case bool:
        return strconv.FormatBool(v)
    case nil:
        return "nil"
    default:
        return fmt.Sprint(v)
    }
}
```

### PERF-003: Implementar Bytecode Compilation
- **Archivo**: Nuevo módulo `pkg/r2bytecode/`
- **Criticidad**: 🟡 ALTA
- **Complejidad**: 🔴 ALTA
- **Estimación**: 40 horas (1 mes)
- **Impacto Performance**: 10x mejora general
- **Impacto Memoria**: 5x reducción en allocaciones

**Descripción:**
Implementar un sistema de compilación a bytecode para evitar la evaluación directa del AST, mejorando dramáticamente el rendimiento.

**Estructura propuesta:**
```go
// pkg/r2bytecode/compiler.go
package r2bytecode

import "github.com/arturoeliasanton/go-r2lang/pkg/r2core"

type Compiler struct {
    constants   []interface{}
    symbolTable *SymbolTable
    scopes      []CompilationScope
    scopeIndex  int
}

type CompilationScope struct {
    instructions Instructions
    lastInstruction EmittedInstruction
    previousInstruction EmittedInstruction
}

type EmittedInstruction struct {
    Opcode   Opcode
    Position int
}

// Opcodes básicos
const (
    OpConstant Opcode = iota
    OpAdd
    OpSub
    OpMul
    OpDiv
    OpPop
    OpTrue
    OpFalse
    OpEqual
    OpNotEqual
    OpMinus
    OpBang
    OpJumpNotTruthy
    OpJump
    OpNull
    OpGetGlobal
    OpSetGlobal
    OpGetLocal
    OpSetLocal
    OpCall
    OpReturnValue
    OpReturn
    OpGetBuiltin
    OpClosure
    OpGetFree
    OpCurrentClosure
)

func (c *Compiler) Compile(node r2core.Node) error {
    switch node := node.(type) {
    case *r2core.Program:
        for _, stmt := range node.Statements {
            err := c.Compile(stmt)
            if err != nil {
                return err
            }
        }
    case *r2core.BinaryExpression:
        err := c.Compile(node.Left)
        if err != nil {
            return err
        }
        err = c.Compile(node.Right)
        if err != nil {
            return err
        }
        switch node.Op {
        case "+":
            c.emit(OpAdd)
        case "-":
            c.emit(OpSub)
        case "*":
            c.emit(OpMul)
        case "/":
            c.emit(OpDiv)
        case "==":
            c.emit(OpEqual)
        case "!=":
            c.emit(OpNotEqual)
        }
    case *r2core.NumberLiteral:
        constant := c.addConstant(node.Value)
        c.emit(OpConstant, constant)
    }
    return nil
}

func (c *Compiler) addConstant(obj interface{}) int {
    c.constants = append(c.constants, obj)
    return len(c.constants) - 1
}

func (c *Compiler) emit(op Opcode, operands ...int) int {
    ins := Make(op, operands...)
    pos := c.addInstruction(ins)
    c.setLastInstruction(op, pos)
    return pos
}
```

### PERF-004: Implementar JIT para Loops Frecuentes
- **Archivo**: `pkg/r2core/for_statement.go`
- **Criticidad**: 🟡 ALTA
- **Complejidad**: 🔴 MUY ALTA
- **Estimación**: 80 horas (2 meses)
- **Impacto Performance**: 50x mejora en loops intensivos
- **Impacto Memoria**: 20x reducción en allocaciones

**Descripción:**
Implementar compilación Just-In-Time para loops que se ejecutan frecuentemente, generando código nativo optimizado.

**Código actual problemático:**
```go
// pkg/r2core/for_statement.go
func (fs *ForStatement) Eval(env *Environment) interface{} {
    // Ejecuta interpretando el AST cada iteración
    for {
        condition := fs.Condition.Eval(env)
        if !toBool(condition) {
            break
        }
        fs.Body.Eval(env) // Interpretación costosa
        fs.Increment.Eval(env)
    }
    return nil
}
```

**Solución propuesta:**
```go
// pkg/r2core/for_statement.go
type CompiledLoop struct {
    executableCode unsafe.Pointer
    executionCount int
    threshold      int
}

var (
    loopCompiler = &JITCompiler{}
    loopCache    = make(map[string]*CompiledLoop)
)

func (fs *ForStatement) Eval(env *Environment) interface{} {
    // Generar hash del loop para identificación
    loopHash := fs.generateHash()
    
    // Buscar en cache
    if compiled, ok := loopCache[loopHash]; ok {
        compiled.executionCount++
        
        // Si se ejecuta frecuentemente, usar versión JIT
        if compiled.executionCount > compiled.threshold {
            return fs.executeJIT(compiled, env)
        }
    } else {
        // Crear entrada en cache
        loopCache[loopHash] = &CompiledLoop{
            executionCount: 1,
            threshold:      1000, // Compilar después de 1000 ejecuciones
        }
    }
    
    // Ejecutar interpretado normalmente
    return fs.executeInterpreted(env)
}

func (fs *ForStatement) executeJIT(compiled *CompiledLoop, env *Environment) interface{} {
    // Ejecutar código nativo compilado
    // Implementación específica de la plataforma
    return callNativeCode(compiled.executableCode, env)
}

func (fs *ForStatement) generateHash() string {
    // Generar hash basado en la estructura del loop
    hasher := sha256.New()
    fs.writeToHasher(hasher)
    return hex.EncodeToString(hasher.Sum(nil))
}
```

### PERF-005: Optimizar Array Operations
- **Archivo**: `pkg/r2core/index_expression.go`
- **Criticidad**: 🟡 ALTA
- **Complejidad**: 🟡 MEDIA
- **Estimación**: 12 horas
- **Impacto Performance**: 4x mejora en operaciones array
- **Impacto Memoria**: 3x reducción en allocaciones

**Descripción:**
Optimizar el acceso a elementos de arrays usando índices pre-calculados y bounds checking eficiente.

**Código actual problemático:**
```go
// pkg/r2core/index_expression.go
func (ie *IndexExpression) Eval(env *Environment) interface{} {
    left := ie.Left.Eval(env)
    index := ie.Index.Eval(env)
    
    switch leftVal := left.(type) {
    case []interface{}:
        idx := int(toFloat(index)) // Conversión costosa cada vez
        if idx < 0 {
            idx = len(leftVal) + idx // Cálculo cada vez
        }
        if idx >= len(leftVal) || idx < 0 {
            return nil
        }
        return leftVal[idx]
    }
    return nil
}
```

**Solución propuesta:**
```go
// pkg/r2core/index_expression.go
func (ie *IndexExpression) Eval(env *Environment) interface{} {
    left := ie.Left.Eval(env)
    index := ie.Index.Eval(env)
    
    switch leftVal := left.(type) {
    case []interface{}:
        // Optimizar conversión de índice
        var idx int
        switch indexVal := index.(type) {
        case float64:
            idx = int(indexVal)
        case int:
            idx = indexVal
        default:
            idx = int(toFloat(index))
        }
        
        arrayLen := len(leftVal)
        
        // Bounds checking optimizado
        if idx < 0 {
            idx += arrayLen
            if idx < 0 {
                return nil
            }
        } else if idx >= arrayLen {
            return nil
        }
        
        return leftVal[idx]
    }
    return nil
}
```

---

## 📊 CRONOGRAMA DE IMPLEMENTACIÓN

### Fase 1: Fixes Críticos (Semana 1-2)
| Tarea | Tiempo | Desarrollador | Impacto |
|-------|--------|---------------|---------|
| BUG-001: Lazy Evaluation | 2h | Junior | 2.5x performance |
| BUG-002: Environment Cache | 4h | Senior | 3x performance |
| BUG-003: Type Conversion Cache | 6h | Senior | 4x performance |
| **Total Fase 1** | **12h** | | **8x performance** |

### Fase 2: Optimizaciones Core (Semana 3-4)
| Tarea | Tiempo | Desarrollador | Impacto |
|-------|--------|---------------|---------|
| PERF-001: Object Pool | 8h | Senior | 6x performance |
| PERF-002: String Builder | 6h | Junior | 3x performance |
| PERF-005: Array Optimization | 12h | Senior | 4x performance |
| **Total Fase 2** | **26h** | | **12x performance** |

### Fase 3: Compilación Avanzada (Mes 2-3)
| Tarea | Tiempo | Desarrollador | Impacto |
|-------|--------|---------------|---------|
| PERF-003: Bytecode Compiler | 40h | Senior | 10x performance |
| **Total Fase 3** | **40h** | | **20x performance** |

### Fase 4: JIT Implementation (Mes 4-6)
| Tarea | Tiempo | Desarrollador | Impacto |
|-------|--------|---------------|---------|
| PERF-004: JIT Compiler | 80h | Expert | 50x performance |
| **Total Fase 4** | **80h** | | **100x performance** |

---

## 🎯 MÉTRICAS DE ÉXITO

### Objetivos por Fase

| Fase | Performance Target | Memory Target | Tiempo |
|------|-------------------|---------------|--------|
| Inicial | 122,670 ns/op | 85,297 B/op | - |
| Fase 1 | 15,334 ns/op (8x) | 10,662 B/op (8x) | 2 semanas |
| Fase 2 | 1,278 ns/op (96x) | 3,554 B/op (24x) | 4 semanas |
| Fase 3 | 127 ns/op (966x) | 710 B/op (120x) | 3 meses |
| Fase 4 | 12.7 ns/op (9,660x) | 142 B/op (600x) | 6 meses |

### Comandos de Validación

```bash
# Ejecutar benchmarks antes de cada fase
go test -bench=. -benchmem performance_test.go > before_phase_X.txt

# Ejecutar después de implementar mejoras
go test -bench=. -benchmem performance_test.go > after_phase_X.txt

# Comparar resultados
benchcmp before_phase_X.txt after_phase_X.txt

# Profiling detallado
go test -bench=BenchmarkBasicArithmetic -cpuprofile=cpu.prof -memprofile=mem.prof performance_test.go
go tool pprof cpu.prof
```

---

## 🛠️ HERRAMIENTAS REQUERIDAS

### Desarrollo
- **Go 1.21+**: Para optimizaciones de runtime
- **benchcmp**: Para comparar benchmarks
- **go tool pprof**: Para profiling
- **go tool trace**: Para análisis de traces

### Testing
- **stress**: Para testing de carga
- **vegeta**: Para testing de performance HTTP
- **go-torch**: Para flame graphs

### Monitoreo
- **Prometheus**: Para métricas en tiempo real
- **Grafana**: Para visualización
- **pprof server**: Para profiling continuo

---

## 🚨 RIESGOS Y MITIGACIONES

### Riesgos Técnicos
1. **Complejidad del JIT**: Muy alta complejidad de implementación
   - **Mitigación**: Implementar primero versión simplificada
   
2. **Compatibilidad**: Cambios pueden romper funcionalidad existente
   - **Mitigación**: Suite de tests comprehensiva
   
3. **Memory Leaks**: Pools y caches pueden causar memory leaks
   - **Mitigación**: Implementar cleanup automático

### Riesgos de Cronograma
1. **Subestimación**: Tareas complejas pueden tomar más tiempo
   - **Mitigación**: Buffer de 25% en estimaciones
   
2. **Dependencias**: Algunas optimizaciones dependen de otras
   - **Mitigación**: Planificación cuidadosa de dependencies

---

## 📈 CONCLUSIONES

### Prioridades Inmediatas
1. **Implementar BUG-001, BUG-002, BUG-003**: Fixes críticos con alto impacto
2. **PERF-001 Object Pool**: Mayor impacto en reducción de allocaciones
3. **PERF-002 String Optimization**: Mejora significativa en operaciones comunes

### Retorno de Inversión
- **Fase 1**: 12 horas → 8x mejora (ROI: 667%)
- **Fase 2**: 26 horas → 12x mejora adicional (ROI: 462%)
- **Fase 3**: 40 horas → 20x mejora adicional (ROI: 500%)

### Recomendación Final
Implementar el roadmap por fases permite validar mejoras incrementalmente y ajustar prioridades según resultados reales. El enfoque gradual reduce riesgos y permite aprendizaje continuo.

**Con este plan, R2Lang puede alcanzar performance competitiva con intérpretes comerciales en 6 meses.**