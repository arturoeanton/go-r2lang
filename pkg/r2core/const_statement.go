package r2core

// ConstStatement => const x = expr;
type ConstStatement struct {
	Name  string
	Value Node
}

func (cs *ConstStatement) Eval(env *Environment) interface{} {
	var val interface{}
	if cs.Value != nil {
		val = cs.Value.Eval(env)
	}
	env.SetConst(cs.Name, val)
	return nil
}

// MultipleConstStatement => const x = 1, y = 2, z = 3;
type MultipleConstStatement struct {
	Declarations []ConstDeclaration
}

type ConstDeclaration struct {
	Name  string
	Value Node
}

func (mcs *MultipleConstStatement) Eval(env *Environment) interface{} {
	for _, decl := range mcs.Declarations {
		var val interface{}
		if decl.Value != nil {
			val = decl.Value.Eval(env)
		}
		env.SetConst(decl.Name, val)
	}
	return nil
}
