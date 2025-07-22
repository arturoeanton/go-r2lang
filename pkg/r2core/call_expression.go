package r2core

import "fmt"

type CallExpression struct {
	BaseNode
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
	// Expandir spreads en argumentos si los hay
	argVals = ExpandSpreadInFunctionCall(argVals)

	// P6 Feature: Check for placeholders in arguments
	if hasPlaceholders(argVals) {
		// Create a partial function if placeholders are found
		switch cv := calleeVal.(type) {
		case BuiltinFunction, *UserFunction, func(...interface{}) interface{}:
			return createPartialFunction(cv, argVals)
		default:
			if ce.Position != nil && env.CurrentFile != "" {
				ce.Position.Filename = env.CurrentFile
			}
			PanicWithStack(ce.Position, "Cannot create partial application with non-function value ["+fmt.Sprintf("%T", ce.Callee)+"]", env.callStack)
			return nil
		}
	}

	switch cv := calleeVal.(type) {
	case BuiltinFunction:
		return cv(argVals...)
	case *UserFunction:
		if flagSuper {
			return cv.SuperCall(env, argVals...)
		}
		return cv.Call(argVals...)
	case *PartialFunction:
		// P6 Feature: Apply arguments to partial function
		return cv.Apply(argVals...)
	case *CurriedFunction:
		// P6 Feature: Apply arguments to curried function
		var result interface{} = cv
		for _, arg := range argVals {
			if curriedResult, ok := result.(*CurriedFunction); ok {
				result = curriedResult.Apply(arg)
			} else {
				// Function is fully applied, return the result
				return result
			}
		}
		return result
	case map[string]interface{}:
		// Instanciar un blueprint
		return instantiateObject(env, cv, argVals)
	case func() interface{}:
		// Handle Go native functions with no args
		return cv()
	case func(...interface{}) interface{}:
		// Handle Go native functions with variable args
		return cv(argVals...)
	case func(string) interface{}:
		// Handle DSL use function with string argument
		if len(argVals) > 0 {
			if str, ok := argVals[0].(string); ok {
				return cv(str)
			}
		}
		return cv("")
	default:
		if ce.Position != nil && env.CurrentFile != "" {
			ce.Position.Filename = env.CurrentFile
		}
		PanicWithStack(ce.Position, "Attempt to call something that is neither a function nor a blueprint ["+fmt.Sprintf("%T", ce.Callee)+"]", env.callStack)
		return nil
	}
}
