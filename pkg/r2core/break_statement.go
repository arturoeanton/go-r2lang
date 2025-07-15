package r2core

// BreakStatement representa una declaración break
type BreakStatement struct{}

// BreakValue representa el valor de control que se devuelve cuando se ejecuta break
type BreakValue struct{}

func (bs *BreakStatement) Eval(env *Environment) interface{} {
	return BreakValue{}
}
