package r2core

// Función con nombre
type FunctionDeclaration struct {
	Name string
	Args []string
	Body *BlockStatement
}

func (fd *FunctionDeclaration) Eval(env *Environment) interface{} {
	fn := &UserFunction{
		Args:     fd.Args,
		Body:     fd.Body,
		Env:      env,
		IsMethod: false,
		code:     fd.Name,
	}
	env.Set(fd.Name, fn)
	return nil
}
