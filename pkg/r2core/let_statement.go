package r2core

// LetStatement => let x = expr;
type LetStatement struct {
	Name  string
	Value Node
}

func (ls *LetStatement) Eval(env *Environment) interface{} {
	var val interface{}
	if ls.Value != nil {
		val = ls.Value.Eval(env)
	}
	env.Set(ls.Name, val)
	return nil
}
