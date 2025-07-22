package r2core

type Identifier struct {
	BaseNode
	Name string
}

func (id *Identifier) Eval(env *Environment) interface{} {
	// P6 Feature: Handle '_' as placeholder
	if id.Name == "_" {
		return &Placeholder{BaseNode: id.BaseNode}
	}

	val, ok := env.Get(id.Name)
	if !ok {
		if id.Position != nil && env.CurrentFile != "" {
			id.Position.Filename = env.CurrentFile
		}
		PanicWithStack(id.Position, "Undeclared variable: "+id.Name, env.callStack)
	}
	return val
}
