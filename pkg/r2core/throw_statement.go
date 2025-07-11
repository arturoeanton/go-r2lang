package r2core

type ThrowStatement struct {
	Message string
}

func (ts *ThrowStatement) Eval(env *Environment) interface{} {
	panic(ts.Message)
	//return nil
}
