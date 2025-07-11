package r2core

import "fmt"

type CallExpression struct {
	Callee Node
	Args   []Node
}

func (ce *CallExpression) Eval(env *Environment) interface{} {
	flagSuper := false
	if o := ce.Callee; o != nil {
		if ae, ok := o.(*AccessExpression); ok {
			if id, ok := ae.Object.(*Identifier); ok {
				if id.Name == "super" {
					flagSuper = true
				}
			}
		}
	}

	calleeVal := ce.Callee.Eval(env)

	var argVals []interface{}
	for _, a := range ce.Args {
		argVals = append(argVals, a.Eval(env))
	}
	switch cv := calleeVal.(type) {
	case BuiltinFunction:
		return cv(argVals...)
	case *UserFunction:
		if flagSuper {
			return cv.SuperCall(env, argVals...)
		}
		return cv.Call(argVals...)
	case map[string]interface{}:
		// Instanciar un blueprint
		return instantiateObject(env, cv, argVals)
	default:
		panic("Attempt to call something that is neither a function nor a blueprint [" + fmt.Sprintf("%T", ce.Callee) + "]")
	}
}
