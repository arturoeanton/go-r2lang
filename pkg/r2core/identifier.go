package r2core

type Identifier struct {
	Name string
}

func (id *Identifier) Eval(env *Environment) interface{} {
	val, ok := env.Get(id.Name)
	if !ok {
		panic("Undeclared variable: " + id.Name)
	}
	return val
}
