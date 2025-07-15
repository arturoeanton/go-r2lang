package r2core

import (
	"fmt"
	"os"
	"sync"
)

// ============================================================
// 6) ENVIRONMENT
// ============================================================

type Environment struct {
	store    map[string]interface{}
	outer    *Environment
	imported map[string]bool
	Dir      string
	CurrenFx string

	// Cache optimizado para lookup frecuente
	lookupCache   map[string]interface{}
	lookupCacheMu sync.RWMutex
	cacheHits     int
	cacheMisses   int
}

func NewEnvironment() *Environment {
	return &Environment{
		store:       make(map[string]interface{}),
		outer:       nil,
		imported:    make(map[string]bool),
		lookupCache: make(map[string]interface{}),
	}
}

func NewInnerEnv(outer *Environment) *Environment {
	return &Environment{
		store:       make(map[string]interface{}),
		outer:       outer,
		imported:    make(map[string]bool),
		Dir:         outer.Dir,
		lookupCache: make(map[string]interface{}),
	}
}

func (e *Environment) GetStore() map[string]interface{} {
	if e == nil {
		return nil
	}
	return e.store
}

func (e *Environment) Set(name string, value interface{}) {
	e.store[name] = value
	// Limpiar cache cuando se modifica una variable
	e.lookupCacheMu.Lock()
	delete(e.lookupCache, name)
	e.lookupCacheMu.Unlock()
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
	val, ok := e.store[name]
	if ok {
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
