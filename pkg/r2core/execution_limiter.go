package r2core

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Error sentinels para identificación externa
var (
	ErrBudgetExceeded = errors.New("execution budget exceeded")
	ErrInfiniteLoop   = errors.New("infinite loop detected")
	ErrRecursionLimit = errors.New("recursion limit exceeded")
	ErrTimeout        = errors.New("execution timeout")
)

// ExecutionLimiter controla límites de ejecución para prevenir loops infinitos
type ExecutionLimiter struct {
	MaxIterations     int64
	CurrentIterations int64
	MaxRecursionDepth int
	CallStack         []string
	StartTime         time.Time
	MaxExecutionTime  time.Duration
	Enabled           bool
	Context           context.Context
	Cancel            context.CancelFunc
}

// LoopContext representa el contexto de un bucle específico
type LoopContext struct {
	Type          string    // "while", "for", "for-in"
	Iterations    int64     // Comienza en 0 cada bucle
	MaxIterations int64     // Límite específico del bucle
	StartTime     time.Time
	Location      string
}

// InfiniteLoopError representa un error de loop infinito con contexto
type InfiniteLoopError struct {
	Type        string            // "while", "for", "recursion", "timeout"
	Location    string            // Ubicación en el código
	Iterations  int64             // Número de iteraciones
	Duration    time.Duration     // Tiempo transcurrido
	Suggestion  string            // Sugerencia de solución
	Stats       map[string]interface{} // Estadísticas adicionales
	Sentinel    error             // Error sentinel para identificación
}

func (ile *InfiniteLoopError) Error() string {
	return fmt.Sprintf(
		"Loop infinito detectado (%s) en %s:\n"+
			"- Iteraciones: %d\n"+
			"- Tiempo transcurrido: %v\n"+
			"- Sugerencia: %s",
		ile.Type, ile.Location, ile.Iterations, ile.Duration, ile.Suggestion,
	)
}

func (ile *InfiniteLoopError) ShortMessage() string {
	return fmt.Sprintf("%s loop tras %d iteraciones", ile.Type, ile.Iterations)
}

func (ile *InfiniteLoopError) Is(target error) bool {
	return ile.Sentinel == target
}

// NewExecutionLimiter crea un nuevo limitador con valores por defecto
func NewExecutionLimiter() *ExecutionLimiter {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &ExecutionLimiter{
		MaxIterations:     1000000,  // 1M iteraciones por defecto
		MaxRecursionDepth: 1000,     // 1K niveles de recursión
		MaxExecutionTime:  30 * time.Second,
		Enabled:           true,
		StartTime:         time.Now(),
		Context:           ctx,
		Cancel:            cancel,
		CallStack:         make([]string, 0),
	}
}

// NewExecutionLimiterWithTimeout crea un limitador con timeout específico
func NewExecutionLimiterWithTimeout(timeout time.Duration) *ExecutionLimiter {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	
	return &ExecutionLimiter{
		MaxIterations:     1000000,
		MaxRecursionDepth: 1000,
		MaxExecutionTime:  timeout,
		Enabled:           true,
		StartTime:         time.Now(),
		Context:           ctx,
		Cancel:            cancel,
		CallStack:         make([]string, 0),
	}
}

// CheckIterationLimit verifica si se ha alcanzado el límite de iteraciones
func (el *ExecutionLimiter) CheckIterationLimit() bool {
	if !el.Enabled {
		return false
	}
	return el.CurrentIterations >= el.MaxIterations
}

// CheckRecursionDepth verifica si se ha alcanzado el límite de profundidad de recursión
func (el *ExecutionLimiter) CheckRecursionDepth() bool {
	if !el.Enabled {
		return false
	}
	return len(el.CallStack) >= el.MaxRecursionDepth
}

// CheckTimeLimit verifica si se ha alcanzado el límite de tiempo
func (el *ExecutionLimiter) CheckTimeLimit() bool {
	if !el.Enabled || el.MaxExecutionTime == 0 {
		return false
	}
	return time.Since(el.StartTime) >= el.MaxExecutionTime
}

// CheckContext verifica si el contexto ha sido cancelado
func (el *ExecutionLimiter) CheckContext() bool {
	if !el.Enabled || el.Context == nil {
		return false
	}
	
	select {
	case <-el.Context.Done():
		return true
	default:
		return false
	}
}

// IncrementIterations incrementa el contador de iteraciones
func (el *ExecutionLimiter) IncrementIterations() {
	if el.Enabled {
		el.CurrentIterations++
	}
}

// EnterFunction registra la entrada a una función (para recursión)
func (el *ExecutionLimiter) EnterFunction(functionName string) {
	if el.Enabled {
		el.CallStack = append(el.CallStack, functionName)
	}
}

// ExitFunction registra la salida de una función
func (el *ExecutionLimiter) ExitFunction() {
	if el.Enabled && len(el.CallStack) > 0 {
		el.CallStack = el.CallStack[:len(el.CallStack)-1]
	}
}

// Reset reinicia los contadores
func (el *ExecutionLimiter) Reset() {
	el.CurrentIterations = 0
	el.CallStack = el.CallStack[:0]
	el.StartTime = time.Now()
}

// SetLimits configura los límites
func (el *ExecutionLimiter) SetLimits(maxIter int64, maxDepth int, maxTime time.Duration) {
	el.MaxIterations = maxIter
	el.MaxRecursionDepth = maxDepth
	el.MaxExecutionTime = maxTime
}

// Enable habilita el limitador
func (el *ExecutionLimiter) Enable() {
	el.Enabled = true
}

// Disable deshabilita el limitador
func (el *ExecutionLimiter) Disable() {
	el.Enabled = false
}

// NewInfiniteLoopError crea un error de loop infinito
func NewInfiniteLoopError(loopType string, ctx *LoopContext) *InfiniteLoopError {
	suggestions := map[string]string{
		"while":     "Verifica que la condición del while pueda volverse false",
		"for":       "Asegúrate de que el incremento modifique la condición",
		"recursion": "Añade un caso base que termine la recursión",
		"timeout":   "Considera dividir el trabajo en partes más pequeñas",
	}
	
	var iterations int64
	var duration time.Duration
	var location string
	
	if ctx != nil {
		iterations = ctx.Iterations
		duration = time.Since(ctx.StartTime)
		location = ctx.Location
	}
	
	return &InfiniteLoopError{
		Type:       loopType,
		Location:   location,
		Iterations: iterations,
		Duration:   duration,
		Suggestion: suggestions[loopType],
		Sentinel:   ErrInfiniteLoop,
		Stats: map[string]interface{}{
			"loop_type": loopType,
		},
	}
}

// NewRecursionError crea un error de recursión infinita
func NewRecursionError(errorType string, context interface{}) *InfiniteLoopError {
	return &InfiniteLoopError{
		Type:     "recursion",
		Location: fmt.Sprintf("recursion depth: %v", context),
		Sentinel: ErrRecursionLimit,
		Suggestion: "Añade un caso base que termine la recursión",
		Stats: map[string]interface{}{
			"recursion_type": errorType,
			"context":        context,
		},
	}
}

// NewTimeoutError crea un error de timeout
func NewTimeoutError(timeoutType string, ctx context.Context) *InfiniteLoopError {
	return &InfiniteLoopError{
		Type:     "timeout",
		Location: "global execution",
		Sentinel: ErrTimeout,
		Suggestion: "Considera dividir el trabajo en partes más pequeñas",
		Stats: map[string]interface{}{
			"timeout_type": timeoutType,
			"context_err":  ctx.Err(),
		},
	}
}