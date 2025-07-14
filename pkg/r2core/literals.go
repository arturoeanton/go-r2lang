package r2core

type NumberLiteral struct {
	Value float64
}

func (nl *NumberLiteral) Eval(env *Environment) interface{} {
	// Object pool desactivado para operaciones simples - usar solo valor directo
	return nl.Value
}

type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) Eval(env *Environment) interface{} {
	return sl.Value
}

type BooleanLiteral struct {
	Value bool
}

func (b *BooleanLiteral) Eval(env *Environment) interface{} {
	return b.Value
}

type ArrayLiteral struct {
	Elements []Node
}

func (al *ArrayLiteral) Eval(env *Environment) interface{} {
	var result []interface{}
	for _, e := range al.Elements {
		result = append(result, e.Eval(env))
	}
	return result
}

// Queremos un "FunctionLiteral" para soportar func(...) { ... } anónimas
type FunctionLiteral struct {
	Args []string
	Body *BlockStatement
}

func (fl *FunctionLiteral) Eval(env *Environment) interface{} {
	fn := &UserFunction{
		Args:     fl.Args,
		Body:     fl.Body,
		Env:      env, // closure
		IsMethod: false,
	}
	return fn
}
