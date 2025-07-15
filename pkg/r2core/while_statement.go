package r2core

import (
	"time"
)

type WhileStatement struct {
	Condition Node
	Body      *BlockStatement
}

func (ws *WhileStatement) Eval(env *Environment) interface{} {
	limiter := env.GetLimiter()

	// Crear contexto de bucle
	loopCtx := &LoopContext{
		Type:          "while",
		Iterations:    0,
		MaxIterations: limiter.MaxIterations,
		StartTime:     time.Now(),
		Location:      "while statement", // TODO: agregar ubicación real del archivo
	}

	var result interface{}
	for {
		// Verificar límites antes de cada iteración
		if limiter.Enabled {
			// Verificar timeout global
			if limiter.CheckTimeLimit() {
				panic(NewTimeoutError("while_timeout", env.GetContext()))
			}

			// Verificar context cancelation
			if limiter.CheckContext() {
				panic(NewTimeoutError("while_context_canceled", env.GetContext()))
			}

			// Verificar límite de iteraciones del bucle
			if loopCtx.Iterations >= loopCtx.MaxIterations {
				panic(NewInfiniteLoopError("while", loopCtx))
			}
		}

		condVal := ws.Condition.Eval(env)
		if !toBool(condVal) {
			break
		}

		// Incrementar contador de iteraciones del bucle
		loopCtx.Iterations++

		val := ws.Body.Eval(env)
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
	return result
}
