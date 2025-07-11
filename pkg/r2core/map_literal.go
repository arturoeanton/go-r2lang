package r2core

type MapLiteral struct {
	Pairs map[string]Node
}

func (ml *MapLiteral) Eval(env *Environment) interface{} {
	m := make(map[string]interface{})
	for k, nd := range ml.Pairs {
		m[k] = nd.Eval(env)
	}
	return m
}
