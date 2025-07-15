# Propuesta: Detección de Loops Infinitos en R2Lang

**Versión:** 1.0  
**Fecha:** 2025-07-15  
**Estado:** Propuesta  

## Resumen Ejecutivo

Esta propuesta presenta un sistema integral para detectar y prevenir loops infinitos en R2Lang, proporcionando mecanismos de protección tanto en tiempo de ejecución como configurables por el usuario.

## Problema Actual

R2Lang actualmente no tiene protección contra loops infinitos, lo que puede causar:

- **Consumo excesivo de CPU:** Loops que nunca terminan
- **Bloqueo del intérprete:** Imposibilidad de interrumpir ejecución
- **Consumo de memoria:** Acumulación de stack frames en recursión infinita
- **Mala experiencia de usuario:** Programas que se "cuelgan"

### Ejemplos de Código Problemático

```r2
// Loop while infinito
while(true) {
    print("infinito");
}

// Loop for infinito
for(let i = 0; i >= 0; i++) {
    print(i);
}

// Recursión infinita
func infiniteRecursion() {
    infiniteRecursion();
}
```

## Solución Propuesta

### 1. Sistema de Contadores de Iteración

#### 1.1 Contador Global de Iteraciones
```go
type ExecutionLimiter struct {
    MaxIterations     int64
    CurrentIterations int64
    MaxRecursionDepth int
    CurrentDepth      int
    StartTime         time.Time
    MaxExecutionTime  time.Duration
}
```

#### 1.2 Integración en Estructuras de Control

**En WhileStatement:**
```go
func (ws *WhileStatement) Eval(env *Environment) interface{} {
    limiter := env.GetLimiter()
    
    for {
        if limiter.CheckIterationLimit() {
            panic("Loop infinito detectado: Máximo de iteraciones excedido")
        }
        
        condition := ws.Condition.Eval(env)
        if !isTruthy(condition) {
            break
        }
        
        limiter.IncrementIterations()
        ws.Body.Eval(env)
    }
    return nil
}
```

**En ForStatement:**
```go
func (fs *ForStatement) Eval(env *Environment) interface{} {
    limiter := env.GetLimiter()
    
    if fs.Init != nil {
        fs.Init.Eval(env)
    }
    
    for {
        if limiter.CheckIterationLimit() {
            panic("Loop infinito detectado en for: Máximo de iteraciones excedido")
        }
        
        if fs.Condition != nil {
            condition := fs.Condition.Eval(env)
            if !isTruthy(condition) {
                break
            }
        }
        
        limiter.IncrementIterations()
        fs.Body.Eval(env)
        
        if fs.Post != nil {
            fs.Post.Eval(env)
        }
    }
    return nil
}
```

### 2. Detección de Recursión Infinita

#### 2.1 Stack de Llamadas con Límite
```go
func (fc *FunctionCall) Eval(env *Environment) interface{} {
    limiter := env.GetLimiter()
    
    if limiter.CheckRecursionDepth() {
        panic("Recursión infinita detectada: Máxima profundidad excedida")
    }
    
    limiter.EnterFunction(fc.FunctionName)
    defer limiter.ExitFunction()
    
    // Lógica existente de llamada a función
    return result
}
```

#### 2.2 Detección de Patrones Recursivos
```go
type CallStack struct {
    Functions []string
    Counts    map[string]int
}

func (cs *CallStack) DetectInfinitePattern(funcName string) bool {
    cs.Counts[funcName]++
    
    // Si una función se llama más de N veces consecutivas
    if cs.Counts[funcName] > 1000 {
        return true
    }
    
    // Detectar patrones cíclicos simples (A->B->A->B...)
    if len(cs.Functions) >= 4 {
        pattern := cs.Functions[len(cs.Functions)-2:]
        prevPattern := cs.Functions[len(cs.Functions)-4 : len(cs.Functions)-2]
        return reflect.DeepEqual(pattern, prevPattern)
    }
    
    return false
}
```

### 3. Límites de Tiempo de Ejecución

#### 3.1 Timeout Global
```go
func (env *Environment) ExecuteWithTimeout(node Node, timeout time.Duration) interface{} {
    limiter := env.GetLimiter()
    limiter.StartTime = time.Now()
    limiter.MaxExecutionTime = timeout
    
    done := make(chan interface{}, 1)
    go func() {
        defer func() {
            if r := recover(); r != nil {
                done <- r
            }
        }()
        result := node.Eval(env)
        done <- result
    }()
    
    select {
    case result := <-done:
        return result
    case <-time.After(timeout):
        panic("Tiempo de ejecución excedido: Posible loop infinito")
    }
}
```

### 4. Configuración Flexible

#### 4.1 Configuración por Usuario
```r2
// Configurar límites en código R2Lang
setExecutionLimits({
    maxIterations: 1000000,
    maxRecursionDepth: 1000,
    maxExecutionTime: "30s"
});

// Configuración por scope
func riskyFunction() {
    withLimits({maxIterations: 100}, func() {
        // Código con límites más estrictos
    });
}
```

#### 4.2 Variables de Entorno
```bash
export R2LANG_MAX_ITERATIONS=1000000
export R2LANG_MAX_RECURSION=1000
export R2LANG_MAX_TIME=30s
export R2LANG_INFINITE_DETECTION=true
```

### 5. Implementación Técnica

#### 5.1 Estructura del Limitador
```go
// pkg/r2core/execution_limiter.go
type ExecutionLimiter struct {
    MaxIterations     int64
    CurrentIterations int64
    MaxRecursionDepth int
    CallStack         []string
    StartTime         time.Time
    MaxExecutionTime  time.Duration
    Enabled           bool
}

func NewExecutionLimiter() *ExecutionLimiter {
    return &ExecutionLimiter{
        MaxIterations:     1000000,  // 1M iteraciones por defecto
        MaxRecursionDepth: 1000,     // 1K niveles de recursión
        MaxExecutionTime:  30 * time.Second,
        Enabled:           true,
    }
}

func (el *ExecutionLimiter) CheckIterationLimit() bool {
    if !el.Enabled {
        return false
    }
    return el.CurrentIterations >= el.MaxIterations
}

func (el *ExecutionLimiter) CheckRecursionDepth() bool {
    if !el.Enabled {
        return false
    }
    return len(el.CallStack) >= el.MaxRecursionDepth
}

func (el *ExecutionLimiter) CheckTimeLimit() bool {
    if !el.Enabled || el.MaxExecutionTime == 0 {
        return false
    }
    return time.Since(el.StartTime) >= el.MaxExecutionTime
}
```

#### 5.2 Integración en Environment
```go
// pkg/r2core/environment.go
type Environment struct {
    // ... campos existentes
    limiter *ExecutionLimiter
}

func (env *Environment) GetLimiter() *ExecutionLimiter {
    if env.limiter == nil {
        env.limiter = NewExecutionLimiter()
    }
    return env.limiter
}

func (env *Environment) SetLimits(maxIter int64, maxDepth int, maxTime time.Duration) {
    limiter := env.GetLimiter()
    limiter.MaxIterations = maxIter
    limiter.MaxRecursionDepth = maxDepth
    limiter.MaxExecutionTime = maxTime
}
```

### 6. Funciones Built-in para Control

```go
// pkg/r2libs/r2execution.go
func RegisterExecution(env *r2core.Environment) {
    env.Set("setExecutionLimits", r2core.BuiltinFunction(setExecutionLimits))
    env.Set("getExecutionStats", r2core.BuiltinFunction(getExecutionStats))
    env.Set("resetExecutionCounters", r2core.BuiltinFunction(resetCounters))
    env.Set("withLimits", r2core.BuiltinFunction(withLimits))
}

func setExecutionLimits(args ...interface{}) interface{} {
    if len(args) != 1 {
        panic("setExecutionLimits requiere un objeto de configuración")
    }
    
    config, ok := args[0].(map[string]interface{})
    if !ok {
        panic("Configuración debe ser un objeto")
    }
    
    // Aplicar configuración
    if maxIter, exists := config["maxIterations"]; exists {
        // Configurar límite de iteraciones
    }
    
    return nil
}
```

### 7. Mensajes de Error Informativos

```go
type InfiniteLoopError struct {
    Type        string
    Location    string
    Iterations  int64
    Duration    time.Duration
    Suggestion  string
}

func (ile *InfiniteLoopError) Error() string {
    return fmt.Sprintf(
        "Loop infinito detectado (%s) en %s:\n" +
        "- Iteraciones: %d\n" +
        "- Tiempo transcurrido: %v\n" +
        "- Sugerencia: %s",
        ile.Type, ile.Location, ile.Iterations, ile.Duration, ile.Suggestion,
    )
}
```

## Beneficios

1. **Protección Automática:** Previene loops infinitos sin intervención manual
2. **Configurabilidad:** Límites ajustables según necesidades
3. **Rendimiento:** Overhead mínimo en ejecución normal
4. **Debugging:** Información detallada sobre el problema
5. **Flexibilidad:** Diferentes estrategias de detección

## Consideraciones de Rendimiento

- **Overhead mínimo:** Solo incremento de contadores
- **Configuración inteligente:** Límites razonables por defecto
- **Deshabilitación:** Posibilidad de desactivar en producción
- **Granularidad:** Control por función o bloque específico

## Plan de Implementación

### Fase 1: Infraestructura Base
- [ ] Crear ExecutionLimiter
- [ ] Integrar en Environment
- [ ] Tests unitarios básicos

### Fase 2: Detección de Loops
- [ ] Implementar en WhileStatement
- [ ] Implementar en ForStatement
- [ ] Tests de detección

### Fase 3: Recursión Infinita
- [ ] Stack de llamadas
- [ ] Detección de patrones
- [ ] Tests de recursión

### Fase 4: Configuración Avanzada
- [ ] Built-in functions
- [ ] Variables de entorno
- [ ] Documentación completa

## Conclusión

Esta implementación proporcionará un sistema robusto de detección de loops infinitos que mejorará significativamente la estabilidad y usabilidad de R2Lang, manteniendo un overhead mínimo en el rendimiento normal.