package r2core

import "fmt"

// Placeholder represents the '_' placeholder in partial application
type Placeholder struct {
	BaseNode
}

func (p *Placeholder) Eval(env *Environment) interface{} {
	return p // Return the placeholder itself
}

func (p *Placeholder) String() string {
	return "_"
}

// PartialFunction represents a partially applied function
type PartialFunction struct {
	BaseNode
	Function  interface{}   // Original function (UserFunction, BuiltinFunction, etc.)
	Arguments []interface{} // Mix of actual values and placeholders
	Arity     int           // Expected number of arguments
	IsPartial bool          // Whether this is a partial application
}

func (pf *PartialFunction) Eval(env *Environment) interface{} {
	return pf // Return the partial function itself
}

func (pf *PartialFunction) String() string {
	return fmt.Sprintf("PartialFunction(arity: %d)", pf.Arity)
}

// Apply applies arguments to a partial function
func (pf *PartialFunction) Apply(args ...interface{}) interface{} {
	// If there are existing arguments (from explicit partial), combine them
	if len(pf.Arguments) > 0 && !hasPlaceholders(pf.Arguments) {
		// This is an explicit partial (partial(func, arg1, arg2, ...))
		// Just append the new arguments
		combinedArgs := make([]interface{}, len(pf.Arguments)+len(args))
		copy(combinedArgs, pf.Arguments)
		copy(combinedArgs[len(pf.Arguments):], args)

		// If we have enough arguments, call the function
		if len(combinedArgs) >= pf.Arity {
			return callFunction(pf.Function, combinedArgs[:pf.Arity]...)
		}

		// Otherwise, return a new partial function
		return &PartialFunction{
			Function:  pf.Function,
			Arguments: combinedArgs,
			Arity:     pf.Arity,
			IsPartial: true,
		}
	}

	// This is a placeholder-based partial (func(a, _, b))
	// Create a new arguments slice with placeholders filled
	newArgs := make([]interface{}, len(pf.Arguments))
	copy(newArgs, pf.Arguments)

	argIndex := 0
	filledCount := 0

	for i, arg := range newArgs {
		if _, isPlaceholder := arg.(*Placeholder); isPlaceholder {
			if argIndex < len(args) {
				newArgs[i] = args[argIndex]
				argIndex++
				filledCount++
			}
		} else {
			filledCount++
		}
	}

	// If all arguments are filled, call the function
	if filledCount == pf.Arity {
		return callFunction(pf.Function, newArgs...)
	}

	// Otherwise, return a new partial function
	return &PartialFunction{
		Function:  pf.Function,
		Arguments: newArgs,
		Arity:     pf.Arity,
		IsPartial: true,
	}
}

// CurriedFunction represents a curried function
type CurriedFunction struct {
	BaseNode
	Function      interface{}   // Original function
	Arity         int           // Total number of parameters expected
	Arguments     []interface{} // Arguments collected so far
	OriginalArity int           // Original function arity
}

func (cf *CurriedFunction) Eval(env *Environment) interface{} {
	return cf
}

func (cf *CurriedFunction) String() string {
	return fmt.Sprintf("CurriedFunction(arity: %d, collected: %d)", cf.Arity, len(cf.Arguments))
}

// Apply applies an argument to a curried function
func (cf *CurriedFunction) Apply(arg interface{}) interface{} {
	newArgs := make([]interface{}, len(cf.Arguments)+1)
	copy(newArgs, cf.Arguments)
	newArgs[len(cf.Arguments)] = arg

	// If we have all arguments, call the function
	if len(newArgs) == cf.Arity {
		return callFunction(cf.Function, newArgs...)
	}

	// Otherwise, return a new curried function with the additional argument
	return &CurriedFunction{
		Function:      cf.Function,
		Arity:         cf.Arity,
		Arguments:     newArgs,
		OriginalArity: cf.OriginalArity,
	}
}

// callFunction is a utility to call different types of functions
func callFunction(fn interface{}, args ...interface{}) interface{} {
	switch f := fn.(type) {
	case *UserFunction:
		return f.Call(args...)
	case BuiltinFunction:
		return f(args...)
	case func(...interface{}) interface{}:
		return f(args...)
	default:
		panic(fmt.Sprintf("Cannot call function of type %T", fn))
	}
}

// isPlaceholder checks if a value is a placeholder
func isPlaceholder(value interface{}) bool {
	_, ok := value.(*Placeholder)
	return ok
}

// hasPlaceholders checks if any of the arguments are placeholders
func hasPlaceholders(args []interface{}) bool {
	for _, arg := range args {
		if isPlaceholder(arg) {
			return true
		}
	}
	return false
}

// countPlaceholders counts the number of placeholders in arguments
func countPlaceholders(args []interface{}) int {
	count := 0
	for _, arg := range args {
		if isPlaceholder(arg) {
			count++
		}
	}
	return count
}

// getFunctionArity returns the arity (number of parameters) of a function
func getFunctionArity(fn interface{}) int {
	switch f := fn.(type) {
	case *UserFunction:
		if len(f.Params) > 0 {
			return len(f.Params)
		}
		return len(f.Args)
	case BuiltinFunction:
		// For built-in functions, we'll assume variable arity
		// This could be enhanced with metadata in the future
		return -1 // Variable arity
	default:
		return -1 // Unknown arity
	}
}

// createPartialFunction creates a partial function from a regular function and arguments with placeholders
func createPartialFunction(fn interface{}, args []interface{}) *PartialFunction {
	arity := getFunctionArity(fn)
	if arity == -1 {
		// For variable arity functions, use the number of provided arguments
		arity = len(args)
	}

	return &PartialFunction{
		Function:  fn,
		Arguments: args,
		Arity:     arity,
		IsPartial: true,
	}
}

// CurryFunction - P6 Feature: curry function implementation
func CurryFunction(args ...interface{}) interface{} {
	if len(args) < 1 {
		panic("curry: at least one argument required (function to curry)")
	}

	fn := args[0]
	arity := getFunctionArity(fn)

	if arity == -1 {
		// For variable arity functions, require an explicit arity
		if len(args) < 2 {
			panic("curry: variable arity function requires explicit arity as second argument")
		}
		if arityVal, ok := args[1].(float64); ok {
			arity = int(arityVal)
		} else {
			panic("curry: arity must be a number")
		}
	}

	return &CurriedFunction{
		Function:      fn,
		Arity:         arity,
		Arguments:     []interface{}{},
		OriginalArity: arity,
	}
}

// PartialBuiltin - P6 Feature: partial function implementation for explicit partial application
func PartialBuiltin(args ...interface{}) interface{} {
	if len(args) < 1 {
		panic("partial: at least one argument required (function to partially apply)")
	}

	fn := args[0]
	partialArgs := args[1:]

	return createPartialFunction(fn, partialArgs)
}
