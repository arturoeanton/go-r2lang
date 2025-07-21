package r2core

// Función con nombre
type FunctionDeclaration struct {
	BaseNode
	Name   string
	Args   []string    // For backward compatibility
	Params []Parameter // New parameter structure with default values
	Body   *BlockStatement
}

func (fd *FunctionDeclaration) Eval(env *Environment) interface{} {
	fn := &UserFunction{
		Args:     fd.Args,
		Params:   fd.Params,
		Body:     fd.Body,
		Env:      env,
		IsMethod: false,
		code:     fd.Name,
		position: fd.Position,
	}
	env.Set(fd.Name, fn)
	return nil
}
