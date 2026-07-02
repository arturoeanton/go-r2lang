package r2core

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
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
	storeMu     sync.RWMutex
	outer       *Environment
	imported    map[string]bool
	importStack []string   // Track import chain for cyclic detection
	callStack   *CallStack // Track R2Lang function call stack
	Dir         string
	CurrenFx    string
	currenFxMu  sync.Mutex // Guards CurrenFx: mutated by UserFunction.Call, which can run concurrently across goroutines sharing the same closure Environment
	CurrentFile string     // Track current file for error reporting

	// Cache optimizado para lookup frecuente
	lookupCache   map[string]interface{}
	lookupCacheMu sync.RWMutex
	cacheHits     int64
	cacheMisses   int64

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
	e.storeMu.RLock()
	defer e.storeMu.RUnlock()
	result := make(map[string]interface{})
	for name, variable := range e.store {
		result[name] = variable.Value
	}
	return result
}

func (e *Environment) Set(name string, value interface{}) {
	e.storeMu.Lock()
	// Check if variable already exists and is const
	if existing, exists := e.store[name]; exists && existing.IsConst {
		e.storeMu.Unlock()
		panic("cannot assign to const variable '" + name + "'")
	}
	e.store[name] = &Variable{Value: value, IsConst: false}
	e.storeMu.Unlock()
	// Limpiar cache cuando se modifica una variable
	e.lookupCacheMu.Lock()
	delete(e.lookupCache, name)
	e.lookupCacheMu.Unlock()
}

// SetConst creates an immutable variable
func (e *Environment) SetConst(name string, value interface{}) {
	e.storeMu.Lock()
	// Check if variable already exists
	if _, exists := e.store[name]; exists {
		e.storeMu.Unlock()
		panic("variable '" + name + "' already declared")
	}
	e.store[name] = &Variable{Value: value, IsConst: true}
	e.storeMu.Unlock()
	// Limpiar cache cuando se modifica una variable
	e.lookupCacheMu.Lock()
	delete(e.lookupCache, name)
	e.lookupCacheMu.Unlock()
}

// Update modifica una variable existente en el scope correcto
func (e *Environment) Update(name string, value interface{}) {
	// Buscar la variable en el scope actual
	e.storeMu.Lock()
	if existing, ok := e.store[name]; ok {
		if existing.IsConst {
			e.storeMu.Unlock()
			panic("cannot assign to const variable '" + name + "'")
		}
		e.store[name] = &Variable{Value: value, IsConst: false}
		e.storeMu.Unlock()
		// Limpiar cache cuando se modifica una variable
		e.lookupCacheMu.Lock()
		delete(e.lookupCache, name)
		e.lookupCacheMu.Unlock()
		return
	}
	e.storeMu.Unlock()

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
		e.lookupCacheMu.RUnlock()
		atomic.AddInt64(&e.cacheHits, 1)
		return val, true
	}
	e.lookupCacheMu.RUnlock()

	// Búsqueda en store local
	e.storeMu.RLock()
	variable, ok := e.store[name]
	e.storeMu.RUnlock()
	if ok {
		val := variable.Value
		// Cachear solo si el cache no está lleno (evitar memory leak)
		e.lookupCacheMu.Lock()
		if len(e.lookupCache) < 100 { // Limitar tamaño del cache
			e.lookupCache[name] = val
		}
		e.lookupCacheMu.Unlock()
		atomic.AddInt64(&e.cacheMisses, 1)
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

	e.SetCurrenFx(".")

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

// GetCurrenFx retorna de forma segura el nombre de la función actualmente en ejecución
func (e *Environment) GetCurrenFx() string {
	e.currenFxMu.Lock()
	defer e.currenFxMu.Unlock()
	return e.CurrenFx
}

// SetCurrenFx establece de forma segura el nombre de la función actualmente en ejecución
func (e *Environment) SetCurrenFx(name string) {
	e.currenFxMu.Lock()
	e.CurrenFx = name
	e.currenFxMu.Unlock()
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
