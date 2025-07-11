package r2core

type TryStatement struct {
	Body         *BlockStatement
	CatchBlock   *BlockStatement
	FinallyBlock *BlockStatement
	ExceptionVar string
}

func (ts *TryStatement) Eval(env *Environment) interface{} {
	defer func() {
		if r := recover(); r != nil {
			if ts.CatchBlock != nil {
				newEnv := NewInnerEnv(env)
				newEnv.Set(ts.ExceptionVar, r)
				ts.CatchBlock.Eval(newEnv)
			}
		}
		if ts.FinallyBlock != nil {
			ts.FinallyBlock.Eval(env)
		}
	}()
	return ts.Body.Eval(env)
}
