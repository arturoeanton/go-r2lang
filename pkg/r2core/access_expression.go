package r2core

import (
	"fmt"
	"sort"
)

type interfaceSlice []interface{}

type AccessExpression struct {
	Object Node
	Member string
}

func (ae *AccessExpression) Eval(env *Environment) interface{} {
	objVal := ae.Object.Eval(env)

	// Unwrap ReturnValue if necessary (recursively)
	for {
		if retVal, ok := objVal.(*ReturnValue); ok {
			objVal = retVal.Value
		} else {
			break
		}
	}

	// Debug: check what type we have after unwrapping
	if retVal, ok := objVal.(*ReturnValue); ok {
		panic(fmt.Sprintf("Still ReturnValue after unwrapping: %T - %v", objVal, retVal))
	}

	switch obj := objVal.(type) {
	case *ObjectInstance:
		return evalMemberAccess(obj, ae.Member)
	case map[string]interface{}:
		return evalMapAccess(obj, ae.Member)
	case interfaceSlice:
		return evalArrayAccess(obj, ae.Member, env)
	case []interface{}:
		return evalArrayAccess(obj, ae.Member, env)
	case *DSLDefinition:
		return evalDSLAccess(obj, ae.Member, env)
	case *DSLResult:
		return evalDSLResultAccess(obj, ae.Member, env)
	default:
		panic(fmt.Sprintf("access to property in unsupported type: %T", objVal))
	}
}

func evalMemberAccess(instance *ObjectInstance, member string) interface{} {
	val, exists := instance.Env.Get(member)
	if !exists {
		panic("The object does not have the property: " + member)
	}
	return val
}

func evalMapAccess(m map[string]interface{}, member string) interface{} {
	val, exists := m[member]
	if !exists {
		panic("The map does not have the key:" + member)
	}
	return val
}

func evalArrayAccess(arr interfaceSlice, member string, env *Environment) interface{} {
	switch member {
	case "len", "length", "size":
		return evalArrayLen(arr)
	case "delete", "remove", "pop", "del":
		return evalArrayDelete(arr)
	case "push", "append", "add", "insert":
		return evalArrayPush(arr)
	case "insert_at":
		return evalArrayInsertAt(arr)
	case "map", "each":
		return evalArrayMap(arr, env)
	case "filter":
		return evalArrayFilter(arr, env)
	case "reverse", "rev":
		return evalArrayReverse(arr)
	case "sort":
		return evalArraySort(arr, env)
	case "find", "index", "find_all", "indexes":
		return evalArrayFind(arr, member, env)
	case "reduce":
		return evalArrayReduce(arr, env)
	case "join":
		return evalArrayJoin(arr)
	default:
		panic("Array does not have property: " + member)
	}
}

func evalArrayLen(arr interfaceSlice) interface{} {
	return BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) != 0 {
			panic("len: only one argument is accepted")
		}
		// Object pool desactivado para operaciones simples
		return float64(len(arr))
	})
}

func evalArrayDelete(arr interfaceSlice) interface{} {
	return BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("delete: at least one argument is required")
		}
		var result interfaceSlice
		for i, v := range arr {
			flag := false
			for _, arg := range args {
				if v == arg {
					flag = true
					break
				}
			}
			if !flag {
				result = append(result, arr[i])
			}
		}
		return result
	})
}

func evalArrayPush(arr interfaceSlice) interface{} {
	return BuiltinFunction(func(args ...interface{}) interface{} {
		var result interfaceSlice
		result = make(interfaceSlice, len(arr))
		copy(result, arr)
		result = append(result, args...)
		return result
	})
}

func evalArrayInsertAt(arr interfaceSlice) interface{} {
	return BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("insert_in: at least two arguments are required")
		}
		index := int(toFloat(args[0]))
		if index < 0 || index >= len(arr) {
			panic("insert_in: index out of range")
		}
		var result interfaceSlice
		for i, v := range arr {
			if i == index {
				result = append(result, args[1:]...)
			}
			result = append(result, v)
		}
		return result
	})
}

func evalArrayMap(arr interfaceSlice, env *Environment) interface{} {
	return BuiltinFunction(func(args ...interface{}) interface{} {
		var result interfaceSlice
		result = make(interfaceSlice, len(arr))
		for i, v := range arr {
			if bf, ok := args[0].(BuiltinFunction); ok {
				result[i] = bf(v)
			}
			if uf, ok := args[0].(*UserFunction); ok {
				result[i] = uf.Call(v)
			}
			if fl, ok := args[0].(*FunctionLiteral); ok {
				result[i] = fl.Eval(env).(*UserFunction).Call(v)
			}
		}
		return result
	})
}

func evalArrayFilter(arr interfaceSlice, env *Environment) interface{} {
	return BuiltinFunction(func(args ...interface{}) interface{} {
		var result interfaceSlice
		for _, v := range arr {

			flag := false
			if bf, ok := args[0].(BuiltinFunction); ok {
				flag = bf(v).(bool)
			}
			if uf, ok := args[0].(*UserFunction); ok {
				flag = uf.Call(v).(bool)
			}
			if fl, ok := args[0].(*FunctionLiteral); ok {
				flag = fl.Eval(env).(*UserFunction).Call(v).(bool)
			}

			if flag {
				result = append(result, v)
			}
		}
		return result
	})
}

func evalArrayReverse(arr interfaceSlice) interface{} {
	return BuiltinFunction(func(args ...interface{}) interface{} {
		var result interfaceSlice
		result = make(interfaceSlice, len(arr))
		for i, v := range arr {
			result[len(arr)-1-i] = v
		}
		return result
	})
}

func evalArraySort(arr interfaceSlice, env *Environment) interface{} {
	return BuiltinFunction(func(args ...interface{}) interface{} {
		var result interfaceSlice
		result = make(interfaceSlice, len(arr))
		copy(result, arr)
		if len(args) == 0 {
			sort.Slice(result, func(i, j int) bool {
				return toFloat(result[i]) < toFloat(result[j])
			})
			return result
		}

		if uf, ok := args[0].(*UserFunction); ok {

			sort.Slice(result, func(i, j int) bool {
				return uf.Call(result[i], result[j]).(bool)
			})

			return result
		}

		return nil
	})
}

func evalArrayFind(arr interfaceSlice, member string, env *Environment) interface{} {
	return BuiltinFunction(func(args ...interface{}) interface{} {
		flagAll := member == "find_all" || member == "indexes"
		var out interfaceSlice
		if len(arr) == 0 {
			panic("find: at least one argument is required and optionally a function find([fx], elem)")
		}

		if len(args) == 1 {
			isFx := false
			if bf, ok := args[0].(BuiltinFunction); ok {
				isFx = true
				for idx, v := range arr {
					if bf(v).(bool) {
						if flagAll {
							out = append(out, idx)
							continue
						}
						return idx
					}
				}
			}

			if uf, ok := args[0].(*UserFunction); ok {
				isFx = true
				for idx, v := range arr {
					if uf.Call(v).(bool) {
						if flagAll {
							out = append(out, idx)
							continue
						}
						return idx
					}
				}
			}

			if fl, ok := args[0].(*FunctionLiteral); ok {
				isFx = true
				for idx, v := range arr {
					if fl.Eval(env).(*UserFunction).Call(v).(bool) {
						if flagAll {
							out = append(out, idx)
							continue
						}
						return idx
					}
				}
			}

			if isFx {
				if flagAll {
					return out
				}
				return nil
			}

			for idx, v := range arr {
				if v == args[0] {
					if flagAll {
						out = append(out, idx)
						continue
					}
					return idx
				}
			}
			if flagAll {
				return out
			}
			return nil
		}

		if len(args) == 2 {
			elem := args[1]
			for idx, v := range arr {
				flag := false
				if bf, ok := args[0].(BuiltinFunction); ok {
					flag = bf(v, elem).(bool)
				}
				if uf, ok := args[0].(*UserFunction); ok {
					flag = uf.Call(v, elem).(bool)
				}
				if fl, ok := args[0].(*FunctionLiteral); ok {
					flag = fl.Eval(env).(*UserFunction).Call(v, elem).(bool)
				}
				if flag {
					if flagAll {
						out = append(out, idx)
						continue
					}
					return idx
				}
			}
			if flagAll {
				return out
			}
			return nil
		}

		panic("find: at least one argument is required and optionally a function find([fx], elem)")

	})
}

func evalArrayReduce(arr interfaceSlice, env *Environment) interface{} {
	return BuiltinFunction(func(args ...interface{}) interface{} {
		if len(arr) == 0 {
			return nil
		}
		var acc interface{}
		for _, v := range arr {
			if bf, ok := args[0].(BuiltinFunction); ok {
				acc = bf(acc, v)
			}
			if uf, ok := args[0].(*UserFunction); ok {
				acc = uf.Call(acc, v)
			}
			if fl, ok := args[0].(*FunctionLiteral); ok {
				acc = fl.Eval(env).(*UserFunction).Call(acc, v)
			}
		}
		return acc
	})
}

func evalArrayJoin(arr interfaceSlice) interface{} {
	return BuiltinFunction(func(args ...interface{}) interface{} {
		sep := ""
		if len(args) == 1 {
			sep = args[0].(string)
		}
		// Solo usar StringBuilder para arrays grandes, concatenación simple para pequeños
		if len(arr) <= 10 {
			// Concatenación simple para arrays pequeños
			var result string
			for i, v := range arr {
				if i > 0 {
					result += sep
				}
				result += fmt.Sprintf("%v", v)
			}
			return result
		}

		// StringBuilder para arrays grandes
		sb := GetStringBuilder()
		defer PutStringBuilder(sb)

		for i, v := range arr {
			if i > 0 {
				sb.WriteString(sep)
			}
			sb.WriteString(fmt.Sprintf("%v", v))
		}
		return sb.String()
	})
}

func evalDSLAccess(dsl *DSLDefinition, member string, env *Environment) interface{} {
	switch member {
	case "use":
		return func(code string) interface{} {
			return dsl.evaluateDSLCode(code, env)
		}
	default:
		panic("DSL does not have property: " + member)
	}
}

func evalDSLResultAccess(dslResult *DSLResult, member string, env *Environment) interface{} {
	switch member {
	case "AST":
		return dslResult.AST
	case "Code":
		return dslResult.Code
	case "Output":
		return dslResult.Output
	case "GetResult":
		return func([]interface{}) interface{} {
			return dslResult.GetResult()
		}
	default:
		return nil
	}
}
