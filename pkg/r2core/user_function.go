package r2core

type UserFunction struct {
	Args     []string
	Body     *BlockStatement
	Env      *Environment
	IsMethod bool
	code     string
}

func (uf *UserFunction) NativeCall(currentEnv *Environment, args ...interface{}) interface{} {
	newEnv := currentEnv
	if newEnv == nil {
		newEnv = NewInnerEnv(uf.Env)
	} else {
		newEnv = currentEnv
	}

	// Detección de recursión infinita
	limiter := newEnv.GetLimiter()
	if limiter.Enabled {
		// Verificar límite de profundidad de recursión
		if limiter.CheckRecursionDepth() {
			panic(NewRecursionError("max_depth", len(limiter.CallStack)))
		}

		// Verificar timeout global
		if limiter.CheckTimeLimit() {
			panic(NewTimeoutError("function_timeout", newEnv.GetContext()))
		}

		// Verificar context cancelation
		if limiter.CheckContext() {
			panic(NewTimeoutError("function_context_canceled", newEnv.GetContext()))
		}

		// Entrar en función (incrementar stack)
		limiter.EnterFunction(uf.code)
		defer limiter.ExitFunction()
	}

	if uf.IsMethod {
		if uf.Env != nil {
			if selfVal, ok := uf.Env.Get("self"); ok {
				newEnv.Set("self", selfVal)
				newEnv.Set("this", selfVal)
			}
		} else {
			if selfVal, ok := currentEnv.Get("self"); ok {
				newEnv.Set("self", selfVal)
				newEnv.Set("this", selfVal)
				if s, ok := newEnv.Get("super"); ok {
					if smap, ok := s.(map[string]interface{}); ok {
						newEnv.Set("super", smap["super"])
					}
				}
			}
		}

	}
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

func (uf *UserFunction) Call(args ...interface{}) interface{} {
	tmp := uf.Env.CurrenFx
	uf.Env.CurrenFx = uf.code
	out := uf.NativeCall(nil, args...)
	uf.Env.CurrenFx = tmp
	return out
}

func (uf *UserFunction) SuperCall(env *Environment, args ...interface{}) interface{} {
	tmp := env.CurrenFx
	env.CurrenFx = uf.code
	out := uf.NativeCall(env, args...)
	env.CurrenFx = tmp
	return out
}

func (uf *UserFunction) CallStep(env *Environment, args ...interface{}) interface{} {
	tmp := env.CurrenFx
	env.CurrenFx = uf.code
	out := uf.NativeCall(env, args...)
	env.CurrenFx = tmp
	return out
}

type BuiltinFunction func(args ...interface{}) interface{}

type ObjectInstance struct {
	Env *Environment
}

func instantiateObject(env *Environment, blueprint map[string]interface{}, argVals []interface{}) *ObjectInstance {
	objEnv := NewInnerEnv(env)
	instance := &ObjectInstance{Env: objEnv}
	for k, v := range blueprint {
		switch vv := v.(type) {
		case *UserFunction:
			newFn := &UserFunction{
				Args:     vv.Args,
				Body:     vv.Body,
				Env:      objEnv,
				IsMethod: true,
			}
			objEnv.Set(k, newFn)
		default:
			objEnv.Set(k, vv)
		}
	}

	objEnv.Set("self", instance)
	objEnv.Set("this", instance)
	if constructor, ok := objEnv.Get("constructor"); ok {
		if constructorFn, isFn := constructor.(*UserFunction); isFn {
			constructorFn.Call(argVals...)
		}
	}

	return instance
}
