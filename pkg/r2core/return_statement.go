package r2core

type ReturnStatement struct {
	Value Node
}

func (rs *ReturnStatement) Eval(env *Environment) interface{} {
	if rs.Value == nil {
		return ReturnValue{Value: nil}
	}
	val := rs.Value.Eval(env)
	return ReturnValue{Value: val}
}
