package r2core

// DateLiteral representa un literal de fecha en el AST
type DateLiteral struct {
	Value *DateValue
}

func (dl *DateLiteral) Eval(env *Environment) interface{} {
	return dl.Value
}

func (dl *DateLiteral) String() string {
	return dl.Value.String()
}