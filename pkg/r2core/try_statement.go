package r2core

type TryStatement struct {
	Body         *BlockStatement
	CatchBlock   *BlockStatement
	FinallyBlock *BlockStatement
	ExceptionVar string
}

func (ts *TryStatement) Eval(env *Environment) interface{} {
	var result interface{}
	var caught interface{}
	unhandled := false

	func() {
		defer func() {
			if r := recover(); r != nil {
				if ts.CatchBlock != nil {
					newEnv := NewInnerEnv(env)
					newEnv.Set(ts.ExceptionVar, r)
					result = ts.CatchBlock.Eval(newEnv)
				} else {
					// No catch block: remember the panic so it can be
					// re-raised after the finally block runs, instead
					// of being silently swallowed.
					unhandled = true
					caught = r
				}
			}
		}()
		result = ts.Body.Eval(env)
	}()

	if ts.FinallyBlock != nil {
		ts.FinallyBlock.Eval(env)
	}

	if unhandled {
		panic(caught)
	}

	return result
}
