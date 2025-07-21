package r2core

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

// ============================================================
// 6) ENVIRONMENT
// ============================================================

// Variable represents a variable with its value and mutability status
type Variable struct {
	Value   interface{}
	IsConst bool
}

type Environment struct {
	store       map[string]*Variable
	outer       *Environment
	imported    map[string]bool
	importStack []string   // Track import chain for cyclic detection
	callStack   *CallStack // Track R2Lang function call stack
	Dir         string
	CurrenFx    string
	CurrentFile string // Track current file for error reporting

	// Cache optimizado para lookup frecuente
	lookupCache   map[string]interface{}
	lookupCacheMu sync.RWMutex
	cacheHits     int
	cacheMisses   int

	// Execution limiter para prevenir loops infinitos
	limiter *ExecutionLimiter
	context context.Context
}

func NewEnvironment() *Environment {
	return &Environment{
		store:       make(map[string]*Variable),
		outer:       nil,
		imported:    make(map[string]bool),
		importStack: make([]string, 0),
		callStack:   &CallStack{Frames: make([]StackFrame, 0)},
		lookupCache: make(map[string]interface{}),
		limiter:     NewExecutionLimiter(),
		context:     context.Background(),
	}
}

func NewInnerEnv(outer *Environment) *Environment {
	return &Environment{
		store:       make(map[string]*Variable),
		outer:       outer,
		imported:    outer.imported,    // Share imported map to prevent duplicates
		importStack: outer.importStack, // Share import stack for cyclic detection
		callStack:   outer.callStack,   // Share call stack for debugging
		Dir:         outer.Dir,
		CurrentFile: outer.CurrentFile,
		lookupCache: make(map[string]interface{}),
		limiter:     outer.limiter, // Compartir limiter con el outer environment
		context:     outer.context,
	}
}

func (e *Environment) GetStore() map[string]interface{} {
	if e == nil {
		return nil
	}
	// Convert Variable map to interface{} map for backward compatibility
	result := make(map[string]interface{})
	for name, variable := range e.store {
		result[name] = variable.Value
	}
	return result
}

func (e *Environment) Set(name string, value interface{}) {
	// Check if variable already exists and is const
	if existing, exists := e.store[name]; exists && existing.IsConst {
		panic("cannot assign to const variable '" + name + "'")
	}
	e.store[name] = &Variable{Value: value, IsConst: false}
	// Limpiar cache cuando se modifica una variable
	e.lookupCacheMu.Lock()
	delete(e.lookupCache, name)
	e.lookupCacheMu.Unlock()
}

// SetConst creates an immutable variable
func (e *Environment) SetConst(name string, value interface{}) {
	// Check if variable already exists
	if _, exists := e.store[name]; exists {
		panic("variable '" + name + "' already declared")
	}
	e.store[name] = &Variable{Value: value, IsConst: true}
	// Limpiar cache cuando se modifica una variable
	e.lookupCacheMu.Lock()
	delete(e.lookupCache, name)
	e.lookupCacheMu.Unlock()
}

// Update modifica una variable existente en el scope correcto
func (e *Environment) Update(name string, value interface{}) {
	// Buscar la variable en el scope actual
	if existing, ok := e.store[name]; ok {
		if existing.IsConst {
			panic("cannot assign to const variable '" + name + "'")
		}
		e.store[name] = &Variable{Value: value, IsConst: false}
		// Limpiar cache cuando se modifica una variable
		e.lookupCacheMu.Lock()
		delete(e.lookupCache, name)
		e.lookupCacheMu.Unlock()
		return
	}

	// Si no está en el scope actual, buscar en el outer scope
	if e.outer != nil {
		e.outer.Update(name, value)
		return
	}

	// Si no existe en ningún scope, crear en el scope actual
	e.Set(name, value)
}

func (e *Environment) Get(name string) (interface{}, bool) {
	// Fast path: buscar en cache optimizado
	e.lookupCacheMu.RLock()
	if val, ok := e.lookupCache[name]; ok {
		e.cacheHits++
		e.lookupCacheMu.RUnlock()
		return val, true
	}
	e.lookupCacheMu.RUnlock()

	// Búsqueda en store local
	variable, ok := e.store[name]
	if ok {
		val := variable.Value
		// Cachear solo si el cache no está lleno (evitar memory leak)
		e.lookupCacheMu.Lock()
		if len(e.lookupCache) < 100 { // Limitar tamaño del cache
			e.lookupCache[name] = val
		}
		e.cacheMisses++
		e.lookupCacheMu.Unlock()
		return val, true
	}

	if e.outer != nil {
		return e.outer.Get(name)
	}
	return nil, false
}

func (e *Environment) Run(parser *Parser) (result interface{}) {

	defer wg.Wait()
	wg = sync.WaitGroup{}

	ast := parser.ParseProgram()

	/*
		defer func() {
			if r := recover(); r != nil {
				_, err := fmt.Fprintln(os.Stderr, "Exception:", r)
				if err != nil {
					panic(err)
				}
				_, err = fmt.Fprintln(os.Stderr, "Current fx -> ", e.CurrenFx)
				if err != nil {
					panic(err)
				}
				os.Exit(1)
			}
		}()//*/

	e.CurrenFx = "."

	// Ejecutar
	result = ast.Eval(e)

	// Llamar a main() si está
	mainVal, ok := e.Get("main")
	if ok {
		mainFn, isFn := mainVal.(*UserFunction)
		if !isFn {
			fmt.Println("Error: ‘main’ is not a function.")
			os.Exit(1)
		}
		result = mainFn.Call()
	}
	return result
}

// GetLimiter retorna el ExecutionLimiter
func (e *Environment) GetLimiter() *ExecutionLimiter {
	if e.limiter == nil {
		e.limiter = NewExecutionLimiter()
	}
	return e.limiter
}

// SetLimiter establece un ExecutionLimiter personalizado
func (e *Environment) SetLimiter(limiter *ExecutionLimiter) {
	e.limiter = limiter
}

// GetContext retorna el contexto de ejecución
func (e *Environment) GetContext() context.Context {
	return e.context
}

// SetContext establece el contexto de ejecución
func (e *Environment) SetContext(ctx context.Context) {
	e.context = ctx
}

// SetLimits configura los límites de ejecución
func (e *Environment) SetLimits(maxIter int64, maxDepth int, maxTime time.Duration) {
	limiter := e.GetLimiter()
	limiter.SetLimits(maxIter, maxDepth, maxTime)
}

// ExecuteWithTimeout ejecuta código con un timeout específico
func (e *Environment) ExecuteWithTimeout(node Node, timeout time.Duration) interface{} {
	// Crear contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Crear limiter con timeout
	limiter := NewExecutionLimiterWithTimeout(timeout)

	// Crear environment temporal
	tempEnv := NewInnerEnv(e)
	tempEnv.SetLimiter(limiter)
	tempEnv.SetContext(ctx)

	// Canal para resultado
	done := make(chan interface{}, 1)
	errChan := make(chan error, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					errChan <- err
				} else {
					errChan <- fmt.Errorf("panic: %v", r)
				}
			}
		}()

		result := node.Eval(tempEnv)
		done <- result
	}()

	select {
	case result := <-done:
		return result
	case err := <-errChan:
		panic(err)
	case <-ctx.Done():
		panic(NewTimeoutError("execution_timeout", ctx))
	}
}

// IsImportCycle checks if adding the given file would create a circular import
func (e *Environment) IsImportCycle(filePath string) bool {
	for _, importedFile := range e.importStack {
		if importedFile == filePath {
			return true
		}
	}
	return false
}

// PushImport adds a file to the import stack
func (e *Environment) PushImport(filePath string) {
	e.importStack = append(e.importStack, filePath)
}

// PopImport removes the last file from the import stack
func (e *Environment) PopImport() {
	if len(e.importStack) > 0 {
		e.importStack = e.importStack[:len(e.importStack)-1]
	}
}

// GetImportChain returns the current import chain for error reporting
func (e *Environment) GetImportChain() []string {
	chain := make([]string, len(e.importStack))
	copy(chain, e.importStack)
	return chain
}
