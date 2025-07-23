package r2core

import (
	"fmt"
)

// OptionalAccessExpression represents the optional chaining operator ?.
// It safely accesses properties without throwing errors if the object is nil/null.
type OptionalAccessExpression struct {
	Object Node
	Member string
}

func (oae *OptionalAccessExpression) Eval(env *Environment) interface{} {
	objVal := oae.Object.Eval(env)

	// Unwrap ReturnValue if necessary (recursively)
	for {
		if retVal, ok := objVal.(*ReturnValue); ok {
			objVal = retVal.Value
		} else {
			break
		}
	}

	// If object is nil/null, return nil instead of panicking
	if objVal == nil {
		return nil
	}

	// Debug: check what type we have after unwrapping
	if retVal, ok := objVal.(*ReturnValue); ok {
		panic(fmt.Sprintf("Still ReturnValue after unwrapping: %T - %v", objVal, retVal))
	}

	switch obj := objVal.(type) {
	case *ObjectInstance:
		return evalMemberAccessOptional(obj, oae.Member)
	case map[string]interface{}:
		return evalMapAccessOptional(obj, oae.Member)
	case InterfaceSlice:
		return evalArrayAccessOptional([]interface{}(obj), oae.Member, env)
	case []interface{}:
		return evalArrayAccessOptional(obj, oae.Member, env)
	default:
		// For optional chaining, return nil instead of panicking
		return nil
	}
}

// evalMemberAccessOptional safely accesses object members
func evalMemberAccessOptional(obj *ObjectInstance, member string) interface{} {
	if obj == nil || obj.Env == nil {
		return nil
	}

	if val, exists := obj.Env.Get(member); exists {
		return val
	}

	// Return nil instead of panicking for missing properties
	return nil
}

// evalMapAccessOptional safely accesses map properties
func evalMapAccessOptional(obj map[string]interface{}, member string) interface{} {
	if obj == nil {
		return nil
	}

	if val, exists := obj[member]; exists {
		return val
	}

	// Return nil instead of panicking for missing keys
	return nil
}

// evalArrayAccessOptional safely accesses array methods/properties
func evalArrayAccessOptional(obj []interface{}, member string, env *Environment) interface{} {
	if obj == nil {
		return nil
	}

	switch member {
	case "length":
		return float64(len(obj))
	case "push":
		return func(args ...interface{}) interface{} {
			for _, arg := range args {
				obj = append(obj, arg)
			}
			return float64(len(obj))
		}
	case "pop":
		return func() interface{} {
			if len(obj) == 0 {
				return nil
			}
			last := obj[len(obj)-1]
			obj = obj[:len(obj)-1]
			return last
		}
	case "join":
		return func(args ...interface{}) interface{} {
			separator := ""
			if len(args) > 0 {
				separator = toString(args[0])
			}
			result := ""
			for i, item := range obj {
				if i > 0 {
					result += separator
				}
				result += toString(item)
			}
			return result
		}
	default:
		// Return nil instead of panicking for unknown methods
		return nil
	}
}
