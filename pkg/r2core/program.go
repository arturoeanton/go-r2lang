package r2core

type Program struct {
	Statements []Node
}

func (p *Program) Eval(env *Environment) interface{} {

	var result interface{}
	for _, stmt := range p.Statements {
		val := stmt.Eval(env)

		if rv, ok := val.(ReturnValue); ok {
			return rv.Value
		}
		result = val
	}
	return result
}
