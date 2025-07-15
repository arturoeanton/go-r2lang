package r2core

// MultipleLetStatement representa múltiples declaraciones de variables: let a = 1, b = 2, c = 3;
type MultipleLetStatement struct {
	Declarations []LetDeclaration
}

// LetDeclaration representa una declaración individual en una declaración múltiple
type LetDeclaration struct {
	Name  string
	Value Node
}

// Eval evalúa todas las declaraciones de variables
func (mls *MultipleLetStatement) Eval(env *Environment) interface{} {
	var lastValue interface{}

	for _, decl := range mls.Declarations {
		if decl.Value != nil {
			value := decl.Value.Eval(env)
			env.Set(decl.Name, value)
			lastValue = value
		} else {
			env.Set(decl.Name, nil)
			lastValue = nil
		}
	}

	return lastValue
}
