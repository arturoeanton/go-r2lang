package r2core

type BlockStatement struct {
	Statements []Node
}

func (bs *BlockStatement) Eval(env *Environment) interface{} {
	var result interface{}
	for _, stmt := range bs.Statements {
		val := stmt.Eval(env)
		if rv, ok := val.(ReturnValue); ok {
			return rv
		}
		if _, ok := val.(BreakValue); ok {
			return val
		}
		if _, ok := val.(ContinueValue); ok {
			return val
		}
		result = val
	}
	return result
}
