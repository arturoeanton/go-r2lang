package r2core

import (
	"fmt"
	"sort"
)

type AccessExpression struct {
	Object Node
	Member string
}

func (ae *AccessExpression) Eval(env *Environment) interface{} {
	objVal := ae.Object.Eval(env)
	switch obj := objVal.(type) {
	case *ObjectInstance:
		return evalMemberAccess(obj, ae.Member)
	case map[string]interface{}:
		return evalMapAccess(obj, ae.Member)
	case []interface{}:
		return evalArrayAccess(obj, ae.Member, env)
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

func evalArrayAccess(arr []interface{}, member string, env *Environment) interface{} {
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

func evalArrayLen(arr []interface{}) interface{} {
	return BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) != 0 {
			panic("len: only one argument is accepted")
		}
		return float64(len(arr))
	})
}

func evalArrayDelete(arr []interface{}) interface{} {
	return BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("delete: at least one argument is required")
		}
		newArr := make([]interface{}, 0)
		for i, v := range arr {
			flag := false
			for _, arg := range args {
				if v == arg {
					flag = true
					break
				}
			}
			if flag {
				newArr = append(newArr, arr[i])
			}
		}
		return newArr
	})
}

func evalArrayPush(arr []interface{}) interface{} {
	return BuiltinFunction(func(args ...interface{}) interface{} {
		newArr := make([]interface{}, len(arr))
		copy(newArr, arr)
		newArr = append(newArr, args...)
		return newArr
	})
}

func evalArrayInsertAt(arr []interface{}) interface{} {
	return BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("insert_in: at least two arguments are required")
		}
		index := int(toFloat(args[0]))
		if index < 0 || index >= len(arr) {
			panic("insert_in: index out of range")
		}
		newArr := make([]interface{}, 0)
		for i, v := range arr {
			if i == index {
				newArr = append(newArr, args[1:]...)
			}
			newArr = append(newArr, v)
		}
		return newArr
	})
}

func evalArrayMap(arr []interface{}, env *Environment) interface{} {
	return BuiltinFunction(func(args ...interface{}) interface{} {
		newArr := make([]interface{}, len(arr))
		for i, v := range arr {
			if bf, ok := args[0].(BuiltinFunction); ok {
				newArr[i] = bf(v)
			}
			if uf, ok := args[0].(*UserFunction); ok {
				newArr[i] = uf.Call(v)
			}
			if fl, ok := args[0].(*FunctionLiteral); ok {
				newArr[i] = fl.Eval(env).(*UserFunction).Call(v)
			}
		}
		return newArr
	})
}

func evalArrayFilter(arr []interface{}, env *Environment) interface{} {
	return BuiltinFunction(func(args ...interface{}) interface{} {
		newArr := make([]interface{}, 0)
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
				newArr = append(newArr, v)
			}
		}
		return newArr
	})
}

func evalArrayReverse(arr []interface{}) interface{} {
	return BuiltinFunction(func(args ...interface{}) interface{} {
		newArr := make([]interface{}, len(arr))
		for i, v := range arr {
			newArr[len(arr)-1-i] = v
		}
		return newArr
	})
}

func evalArraySort(arr []interface{}, env *Environment) interface{} {
	return BuiltinFunction(func(args ...interface{}) interface{} {
		newArr := make([]interface{}, len(arr))
		copy(newArr, arr)
		if len(args) == 0 {
			sort.Slice(newArr, func(i, j int) bool {
				return toFloat(newArr[i]) < toFloat(newArr[j])
			})
			return newArr
		}

		if uf, ok := args[0].(*UserFunction); ok {

			sort.Slice(newArr, func(i, j int) bool {
				return uf.Call(newArr[i], newArr[j]).(bool)
			})

			return newArr
		}

		return nil
	})
}

func evalArrayFind(arr []interface{}, member string, env *Environment) interface{} {
	return BuiltinFunction(func(args ...interface{}) interface{} {
		flagAll := member == "find_all" || member == "indexes"
		out := make([]interface{}, 0)
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

func evalArrayReduce(arr []interface{}, env *Environment) interface{} {
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

func evalArrayJoin(arr []interface{}) interface{} {
	return BuiltinFunction(func(args ...interface{}) interface{} {
		sep := ""
		if len(args) == 1 {
			sep = args[0].(string)
		}
		var out string
		for i, v := range arr {
			if i > 0 {
				out += sep
			}
			out += fmt.Sprintf("%v", v)
		}
		return out
	})
}