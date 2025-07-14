package r2core

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
	newEnv := NewInnerEnv(env)

	if fs.inFlag {
		return fs.evalForIn(newEnv)
	}
	return fs.evalStandardFor(newEnv)
}

func (fs *ForStatement) evalForIn(env *Environment) interface{} {
	var result interface{}
	raw, _ := env.Get(fs.inArray)
	env.Set("$c", raw)

	if arr, ok := raw.([]interface{}); ok {
		for i, v := range arr {
			env.Set(fs.inIndexName, float64(i))
			env.Set("$k", float64(i))
			env.Set("$v", v)
			val := fs.Body.Eval(env)
			if rv, ok := val.(ReturnValue); ok {
				return rv
			}
			result = val
		}
	} else if mapVal, ok := raw.(map[string]interface{}); ok {
		for k, v := range mapVal {
			env.Set(fs.inIndexName, k)
			env.Set("$k", k)
			env.Set("$v", v)
			val := fs.Body.Eval(env)
			if rv, ok := val.(ReturnValue); ok {
				return rv
			}
			result = val
		}
	} else {
		panic("Not an array or map for ‘for’")
	}
	return result
}

func (fs *ForStatement) evalStandardFor(env *Environment) interface{} {
	// Intentar optimización específica para loops simples
	if optimized := fs.trySimpleLoopOptimization(env); optimized != nil {
		return optimized
	}
	
	// Ejecución normal para loops complejos
	return fs.executeStandardLoop(env)
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
	var result interface{}
	if fs.Init != nil {
		fs.Init.Eval(env)
	}

	for {
		condVal := fs.Condition.Eval(env)
		if !toBool(condVal) {
			break
		}
		val := fs.Body.Eval(env)
		if rv, ok := val.(ReturnValue); ok {
			return rv
		}
		result = val
		if fs.Post != nil {
			fs.Post.Eval(env)
		}
	}
	return result
}
