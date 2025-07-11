package r2core

type ExprStatement struct {
	Expr Node
}

func (es *ExprStatement) Eval(env *Environment) interface{} {
	return es.Expr.Eval(env)
}
