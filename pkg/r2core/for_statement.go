package r2core

import (
	"time"
)

type ForStatement struct {
	Init      Node
	Condition Node
	Post      Node
	Body      *BlockStatement
	inFlag    bool
	inArray   string
	//inMap       string
	inIndexName string
	LoopID      string // Para identificación JIT
}

func (fs *ForStatement) Eval(env *Environment) interface{} {
	if fs.inFlag {
		return fs.evalForIn(env)
	}
	return fs.evalStandardFor(env)
}

func (fs *ForStatement) evalForIn(env *Environment) interface{} {
	limiter := env.GetLimiter()

	// Crear contexto de bucle
	loopCtx := &LoopContext{
		Type:          "for-in",
		Iterations:    0,
		MaxIterations: limiter.MaxIterations,
		StartTime:     time.Now(),
		Location:      "for-in statement", // TODO: agregar ubicación real del archivo
	}

	var result interface{}
	raw, _ := env.Get(fs.inArray)
	env.Set("$c", raw)

	if arr, ok := raw.(interfaceSlice); ok {
		for i, v := range arr {
			// Verificar límites antes de cada iteración
			if limiter.Enabled {
				// Verificar timeout global
				if limiter.CheckTimeLimit() {
					panic(NewTimeoutError("for_in_timeout", env.GetContext()))
				}

				// Verificar context cancelation
				if limiter.CheckContext() {
					panic(NewTimeoutError("for_in_context_canceled", env.GetContext()))
				}

				// Verificar límite de iteraciones del bucle
				if loopCtx.Iterations >= loopCtx.MaxIterations {
					panic(NewInfiniteLoopError("for-in", loopCtx))
				}
			}

			env.Set(fs.inIndexName, float64(i))
			env.Set("$k", float64(i))
			env.Set("$v", v)

			// Incrementar contador de iteraciones del bucle
			loopCtx.Iterations++

			val := fs.Body.Eval(env)
			if rv, ok := val.(ReturnValue); ok {
				return rv
			}
			if _, ok := val.(BreakValue); ok {
				break
			}
			if _, ok := val.(ContinueValue); ok {
				continue
			}
			result = val
		}
	} else if arr, ok := raw.([]interface{}); ok {
		for i, v := range arr {
			// Verificar límites antes de cada iteración
			if limiter.Enabled {
				// Verificar timeout global
				if limiter.CheckTimeLimit() {
					panic(NewTimeoutError("for_in_timeout", env.GetContext()))
				}

				// Verificar context cancelation
				if limiter.CheckContext() {
					panic(NewTimeoutError("for_in_context_canceled", env.GetContext()))
				}

				// Verificar límite de iteraciones del bucle
				if loopCtx.Iterations >= loopCtx.MaxIterations {
					panic(NewInfiniteLoopError("for-in", loopCtx))
				}
			}

			env.Set(fs.inIndexName, float64(i))
			env.Set("$k", float64(i))
			env.Set("$v", v)

			// Incrementar contador de iteraciones del bucle
			loopCtx.Iterations++

			val := fs.Body.Eval(env)
			if rv, ok := val.(ReturnValue); ok {
				return rv
			}
			if _, ok := val.(BreakValue); ok {
				break
			}
			if _, ok := val.(ContinueValue); ok {
				continue
			}
			result = val
		}
	} else if mapVal, ok := raw.(map[string]interface{}); ok {
		for k, v := range mapVal {
			// Verificar límites antes de cada iteración
			if limiter.Enabled {
				// Verificar timeout global
				if limiter.CheckTimeLimit() {
					panic(NewTimeoutError("for_in_timeout", env.GetContext()))
				}

				// Verificar context cancelation
				if limiter.CheckContext() {
					panic(NewTimeoutError("for_in_context_canceled", env.GetContext()))
				}

				// Verificar límite de iteraciones del bucle
				if loopCtx.Iterations >= loopCtx.MaxIterations {
					panic(NewInfiniteLoopError("for-in", loopCtx))
				}
			}

			env.Set(fs.inIndexName, k)
			env.Set("$k", k)
			env.Set("$v", v)

			// Incrementar contador de iteraciones del bucle
			loopCtx.Iterations++

			val := fs.Body.Eval(env)
			if rv, ok := val.(ReturnValue); ok {
				return rv
			}
			if _, ok := val.(BreakValue); ok {
				break
			}
			if _, ok := val.(ContinueValue); ok {
				continue
			}
			result = val
		}
	} else {
		panic("Not an array or map for 'for'")
	}
	return result
}

func (fs *ForStatement) evalStandardFor(env *Environment) interface{} {
	// Crear un nuevo scope para la inicialización del loop
	newEnv := NewInnerEnv(env)

	// Intentar optimización específica para loops simples
	if optimized := fs.trySimpleLoopOptimization(newEnv); optimized != nil {
		return optimized
	}

	// Ejecución normal para loops complejos
	return fs.executeStandardLoop(newEnv)
}

// trySimpleLoopOptimization intenta optimizar loops aritméticos simples
func (fs *ForStatement) trySimpleLoopOptimization(env *Environment) interface{} {
	// Solo optimizar loops muy específicos y seguros
	// Por ejemplo: for (var i = 0; i < N; i++) { suma += i; }

	// TODO: Implementar detección de patrones simples
	// Por ahora, no optimizar para mantener estabilidad
	return nil
}

func (fs *ForStatement) executeStandardLoop(env *Environment) interface{} {
	limiter := env.GetLimiter()

	// Crear contexto de bucle
	loopCtx := &LoopContext{
		Type:          "for",
		Iterations:    0,
		MaxIterations: limiter.MaxIterations,
		StartTime:     time.Now(),
		Location:      "for statement", // TODO: agregar ubicación real del archivo
	}

	var result interface{}
	if fs.Init != nil {
		fs.Init.Eval(env)
	}

	for {
		// Verificar límites antes de cada iteración
		if limiter.Enabled {
			// Verificar timeout global
			if limiter.CheckTimeLimit() {
				panic(NewTimeoutError("for_timeout", env.GetContext()))
			}

			// Verificar context cancelation
			if limiter.CheckContext() {
				panic(NewTimeoutError("for_context_canceled", env.GetContext()))
			}

			// Verificar límite de iteraciones del bucle
			if loopCtx.Iterations >= loopCtx.MaxIterations {
				panic(NewInfiniteLoopError("for", loopCtx))
			}
		}

		condVal := fs.Condition.Eval(env)
		if !toBool(condVal) {
			break
		}

		// Incrementar contador de iteraciones del bucle
		loopCtx.Iterations++

		val := fs.Body.Eval(env)
		if rv, ok := val.(ReturnValue); ok {
			return rv
		}
		if _, ok := val.(BreakValue); ok {
			break
		}
		if _, ok := val.(ContinueValue); ok {
			if fs.Post != nil {
				fs.Post.Eval(env)
			}
			continue
		}
		result = val
		if fs.Post != nil {
			fs.Post.Eval(env)
		}
	}
	return result
}
