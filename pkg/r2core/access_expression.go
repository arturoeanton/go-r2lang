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

	// Manejar ObjectInstance
	if instance, ok := objVal.(*ObjectInstance); ok {
		val, exists := instance.Env.Get(ae.Member)
		if !exists {
			panic("The object does not have the property: " + ae.Member)
		}
		return val
	}

	// Manejar map[string]interface{}
	if m, ok := objVal.(map[string]interface{}); ok {
		val, exists := m[ae.Member]
		if !exists {
			panic("The map does not have the key:" + ae.Member)
		}
		return val
	}

	if arr, ok := objVal.([]interface{}); ok {
		if ae.Member == "len" || ae.Member == "length" || ae.Member == "size" {
			return BuiltinFunction(func(args ...interface{}) interface{} {
				if len(args) != 0 {
					panic("len: only one argument is accepted")
				}
				return float64(len(arr))
			})
		}

		if ae.Member == "delete" || ae.Member == "remove" || ae.Member == "pop" || ae.Member == "del" {
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

		if ae.Member == "push" || ae.Member == "append" || ae.Member == "add" || ae.Member == "insert" {
			return BuiltinFunction(func(args ...interface{}) interface{} {
				newArr := make([]interface{}, len(arr))
				copy(newArr, arr)
				newArr = append(newArr, args...)
				return newArr
			})
		}

		if ae.Member == "insert_at" {
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

		if ae.Member == "map" || ae.Member == "each" {
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

		if ae.Member == "filter" {
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

		if ae.Member == "reverse" || ae.Member == "rev" {
			return BuiltinFunction(func(args ...interface{}) interface{} {
				newArr := make([]interface{}, len(arr))
				for i, v := range arr {
					newArr[len(arr)-1-i] = v
				}
				return newArr
			})
		}

		if ae.Member == "sort" {
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

		if ae.Member == "find" || ae.Member == "index" || ae.Member == "find_all" || ae.Member == "indexes" {

			return BuiltinFunction(func(args ...interface{}) interface{} {
				flagAll := ae.Member == "find_all" || ae.Member == "indexes"
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

		if ae.Member == "reduce" {
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

		if ae.Member == "join" {
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

		panic("Array does not have property: " + ae.Member)
	}

	panic("ccess to property in unsupported type: " + fmt.Sprintf("%T", objVal))
}
