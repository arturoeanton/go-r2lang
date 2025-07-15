package r2core

// ContinueStatement representa una declaración continue
type ContinueStatement struct{}

// ContinueValue representa el valor de control que se devuelve cuando se ejecuta continue
type ContinueValue struct{}

func (cs *ContinueStatement) Eval(env *Environment) interface{} {
	return ContinueValue{}
}
