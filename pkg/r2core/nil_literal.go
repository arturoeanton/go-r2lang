package r2core

// NilLiteral represents a nil value
type NilLiteral struct{}

func (nl *NilLiteral) Eval(env *Environment) interface{} {
	return nil
}
