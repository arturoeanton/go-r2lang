package r2core

type WhileStatement struct {
	Condition Node
	Body      *BlockStatement
}

func (ws *WhileStatement) Eval(env *Environment) interface{} {
	var result interface{}
	for {
		condVal := ws.Condition.Eval(env)
		if !toBool(condVal) {
			break
		}
		val := ws.Body.Eval(env)
		if rv, ok := val.(ReturnValue); ok {
			return rv
		}
		result = val
	}
	return result
}
