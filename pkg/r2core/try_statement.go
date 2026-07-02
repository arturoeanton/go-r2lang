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
					// Run the catch block under its own recover so that an
					// exception thrown from inside catch doesn't skip the
					// finally block below (it still propagates, just after
					// finally has run).
					func() {
						defer func() {
							if r2 := recover(); r2 != nil {
								unhandled = true
								caught = r2
							}
						}()
						newEnv := NewInnerEnv(env)
						newEnv.Set(ts.ExceptionVar, r)
						result = ts.CatchBlock.Eval(newEnv)
					}()
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
		finallyResult := ts.FinallyBlock.Eval(env)
		// A return/break/continue executed inside finally overrides any
		// pending result or exception from the try/catch block, matching
		// standard try/finally semantics (JS, Java, Python).
		switch finallyResult.(type) {
		case ReturnValue, BreakValue, ContinueValue:
			return finallyResult
		}
	}

	if unhandled {
		panic(caught)
	}

	return result
}
