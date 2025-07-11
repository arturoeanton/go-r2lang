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
		result = val
	}
	return result
}
