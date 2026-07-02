package r2libs

import (
	"fmt"
	"sort"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// sortableSlice implements sort.Interface for []interface{}
type sortableSlice struct {
	slice []interface{}
	comp  *r2core.UserFunction // Optional comparison function
}

func (s *sortableSlice) Len() int {
	return len(s.slice)
}

func (s *sortableSlice) Swap(i, j int) {
	s.slice[i], s.slice[j] = s.slice[j], s.slice[i]
}

func (s *sortableSlice) Less(i, j int) bool {
	if s.comp != nil {
		// Call the user-provided comparison function
		result := s.comp.Call(s.slice[i], s.slice[j])
		// Assuming comparison function returns a boolean (true if i < j)
		if b, ok := result.(bool); ok {
			return b
		}
		// If it returns a number, interpret < 0 as true
		if f, ok := result.(float64); ok {
			return f < 0
		}
		panic("Comparison function must return boolean or number")
	}

	// Default comparison (numeric then string)
	valI := s.slice[i]
	valJ := s.slice[j]

	// Try to compare as numbers
	if fI, okI := valI.(float64); okI {
		if fJ, okJ := valJ.(float64); okJ {
			return fI < fJ
		}
	}

	// Fallback to string comparison
	return fmt.Sprintf("%v", valI) < fmt.Sprintf("%v", valJ)
}

// r2hack.go: Funciones de "seguridad", "forense" y "análisis" para R2.
// Enfoque didáctico, no pretende ser una suite de hacking real.

func RegisterCollections(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		"range": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) != 2 {
				panic("range: solo se aceptan 2 argumentos")
			}
			startF, ok1 := args[0].(float64)
			endF, ok2 := args[1].(float64)
			if !ok1 || !ok2 {
				panic("range: los argumentos deben ser numéricos")
			}
			start := int(startF)
			end := int(endF)
			if start > end {
				panic("range: start no puede ser mayor que end")
			}
			arr := make([]interface{}, end-start)
			for i := start; i < end; i++ {
				arr[i-start] = i
			}
			return arr
		}),

		"repeat": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) != 2 {
				panic("repeat: solo se aceptan 2 argumentos")
			}
			endF, ok := args[0].(float64)
			if !ok {
				panic("repeat: el primer argumento debe ser numérico")
			}
			start := 0
			end := int(endF)
			if start > end {
				panic("range: start no puede ser mayor que end")
			}
			arr := make([]interface{}, end-start)
			for i := start; i < end; i++ {
				arr[i-start] = args[1]
			}
			return arr
		}),

		"copy": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) != 1 {
				panic("copy: solo se acepta un argumento")
			}

			arr, ok := args[0].([]interface{})
			if !ok {
				panic("copy: el argumento debe ser un array")
			}
			newArr := make([]interface{}, len(arr))
			copy(newArr, arr)

			return newArr
		}),

		"slice": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) != 3 {
				panic("slice: solo se aceptan 3 argumentos")
			}

			arr, ok := args[0].([]interface{})
			if !ok {
				panic("slice: el primer argumento debe ser un array")
			}
			startF, ok1 := args[1].(float64)
			endF, ok2 := args[2].(float64)
			if !ok1 || !ok2 {
				panic("slice: start y end deben ser numéricos")
			}
			start := int(startF)
			end := int(endF)

			if start < 0 || start >= len(arr) {
				panic("slice: start fuera de rango")
			}

			if end < 0 || end >= len(arr) {
				panic("slice: end fuera de rango")
			}

			return arr[start:end]
		}),
		"map": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) != 2 {
				panic("map: se aceptan 2 argumentos (array, funcion)")
			}
			arr, ok1 := args[0].([]interface{})
			fn, ok2 := args[1].(*r2core.UserFunction)
			if !ok1 || !ok2 {
				panic("map: los argumentos deben ser (array, funcion)")
			}
			newArr := make([]interface{}, len(arr))
			for i, v := range arr {
				newArr[i] = fn.Call(v)
			}
			return newArr
		}),
		"filter": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) != 2 {
				panic("filter: se aceptan 2 argumentos (array, funcion)")
			}
			arr, ok1 := args[0].([]interface{})
			fn, ok2 := args[1].(*r2core.UserFunction)
			if !ok1 || !ok2 {
				panic("filter: los argumentos deben ser (array, funcion)")
			}
			newArr := make([]interface{}, 0)
			for _, v := range arr {
				if toBool(fn.Call(v)) {
					newArr = append(newArr, v)
				}
			}
			return newArr
		}),
		"reduce": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) != 3 {
				panic("reduce: se aceptan 3 argumentos (array, funcion, inicial)")
			}
			arr, ok1 := args[0].([]interface{})
			fn, ok2 := args[1].(*r2core.UserFunction)
			if !ok1 || !ok2 {
				panic("reduce: los argumentos deben ser (array, funcion, inicial)")
			}
			acc := args[2]
			for _, v := range arr {
				acc = fn.Call(acc, v)
			}
			return acc
		}),
		"sort": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 || len(args) > 2 {
				panic("sort: se acepta 1 o 2 argumentos (array, [funcion_comparacion])")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("sort: el primer argumento debe ser un array")
			}

			s := &sortableSlice{slice: arr}
			if len(args) == 2 {
				compFunc, isUserFunc := args[1].(*r2core.UserFunction)
				if !isUserFunc {
					panic("sort: el segundo argumento debe ser una función de comparación")
				}
				s.comp = compFunc
			}

			sort.Sort(s)
			return s.slice
		}),
		"find": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) != 2 {
				panic("find: se aceptan 2 argumentos (array, funcion)")
			}
			arr, ok1 := args[0].([]interface{})
			fn, ok2 := args[1].(*r2core.UserFunction)
			if !ok1 || !ok2 {
				panic("find: los argumentos deben ser (array, funcion)")
			}
			for _, v := range arr {
				if toBool(fn.Call(v)) {
					return v
				}
			}
			return nil
		}),
		"contains": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) != 2 {
				panic("contains: se aceptan 2 argumentos (array, valor)")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("contains: el primer argumento debe ser un array")
			}
			val := args[1]
			for _, v := range arr {
				if equals(v, val) {
					return true
				}
			}
			return false
		}),

		"indexOf": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) != 2 {
				panic("indexOf: se aceptan 2 argumentos (array, valor)")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("indexOf: el primer argumento debe ser un array")
			}
			val := args[1]
			for i, v := range arr {
				if equals(v, val) {
					return float64(i)
				}
			}
			return float64(-1)
		}),

		"unique": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) != 1 {
				panic("unique: solo se acepta un argumento")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("unique: el argumento debe ser un array")
			}
			newArr := make([]interface{}, 0, len(arr))
			for _, v := range arr {
				found := false
				for _, u := range newArr {
					if equals(u, v) {
						found = true
						break
					}
				}
				if !found {
					newArr = append(newArr, v)
				}
			}
			return newArr
		}),

		"compact": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) != 1 {
				panic("compact: solo se acepta un argumento")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("compact: el argumento debe ser un array")
			}
			newArr := make([]interface{}, 0, len(arr))
			for _, v := range arr {
				if v != nil && toBool(v) {
					newArr = append(newArr, v)
				}
			}
			return newArr
		}),

		"flatten": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 || len(args) > 2 {
				panic("flatten: se acepta 1 o 2 argumentos (array, [depth])")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("flatten: el primer argumento debe ser un array")
			}
			depth := 1
			if len(args) == 2 {
				d, ok := args[1].(float64)
				if !ok {
					panic("flatten: depth debe ser numérico")
				}
				depth = int(d)
			}
			return flattenArray(arr, depth)
		}),

		"chunk": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) != 2 {
				panic("chunk: se aceptan 2 argumentos (array, size)")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("chunk: el primer argumento debe ser un array")
			}
			sizeF, ok := args[1].(float64)
			if !ok {
				panic("chunk: size debe ser numérico")
			}
			size := int(sizeF)
			if size <= 0 {
				panic("chunk: size debe ser mayor que cero")
			}
			result := make([]interface{}, 0, (len(arr)+size-1)/size)
			for i := 0; i < len(arr); i += size {
				end := i + size
				if end > len(arr) {
					end = len(arr)
				}
				piece := make([]interface{}, end-i)
				copy(piece, arr[i:end])
				result = append(result, piece)
			}
			return result
		}),

		"partition": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) != 2 {
				panic("partition: se aceptan 2 argumentos (array, funcion)")
			}
			arr, ok1 := args[0].([]interface{})
			fn, ok2 := args[1].(*r2core.UserFunction)
			if !ok1 || !ok2 {
				panic("partition: los argumentos deben ser (array, funcion)")
			}
			matched := make([]interface{}, 0)
			rest := make([]interface{}, 0)
			for _, v := range arr {
				if toBool(fn.Call(v)) {
					matched = append(matched, v)
				} else {
					rest = append(rest, v)
				}
			}
			return []interface{}{matched, rest}
		}),

		"zip": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) != 2 {
				panic("zip: se aceptan 2 argumentos (array1, array2)")
			}
			arr1, ok1 := args[0].([]interface{})
			arr2, ok2 := args[1].([]interface{})
			if !ok1 || !ok2 {
				panic("zip: ambos argumentos deben ser arrays")
			}
			n := len(arr1)
			if len(arr2) < n {
				n = len(arr2)
			}
			result := make([]interface{}, n)
			for i := 0; i < n; i++ {
				result[i] = []interface{}{arr1[i], arr2[i]}
			}
			return result
		}),

		"groupBy": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) != 2 {
				panic("groupBy: se aceptan 2 argumentos (array, funcion)")
			}
			arr, ok1 := args[0].([]interface{})
			fn, ok2 := args[1].(*r2core.UserFunction)
			if !ok1 || !ok2 {
				panic("groupBy: los argumentos deben ser (array, funcion)")
			}
			groups := make(map[string]interface{})
			order := make([]string, 0)
			for _, v := range arr {
				key := fmt.Sprintf("%v", fn.Call(v))
				existing, found := groups[key]
				if !found {
					order = append(order, key)
					groups[key] = []interface{}{v}
					continue
				}
				groups[key] = append(existing.([]interface{}), v)
			}
			result := make(map[string]interface{}, len(groups))
			for _, key := range order {
				result[key] = groups[key]
			}
			return result
		}),

		"sortBy": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) != 2 {
				panic("sortBy: se aceptan 2 argumentos (array, funcion)")
			}
			arr, ok1 := args[0].([]interface{})
			fn, ok2 := args[1].(*r2core.UserFunction)
			if !ok1 || !ok2 {
				panic("sortBy: los argumentos deben ser (array, funcion)")
			}
			type keyedValue struct {
				key interface{}
				val interface{}
			}
			pairs := make([]keyedValue, len(arr))
			for i, v := range arr {
				pairs[i] = keyedValue{key: fn.Call(v), val: v}
			}
			sort.SliceStable(pairs, func(i, j int) bool {
				ki, kj := pairs[i].key, pairs[j].key
				if fi, ok := ki.(float64); ok {
					if fj, ok := kj.(float64); ok {
						return fi < fj
					}
				}
				return fmt.Sprintf("%v", ki) < fmt.Sprintf("%v", kj)
			})
			newArr := make([]interface{}, len(pairs))
			for i, p := range pairs {
				newArr[i] = p.val
			}
			return newArr
		}),
	}

	RegisterModule(env, "collections", functions)
}

func flattenArray(arr []interface{}, depth int) []interface{} {
	if depth <= 0 {
		result := make([]interface{}, len(arr))
		copy(result, arr)
		return result
	}
	result := make([]interface{}, 0, len(arr))
	for _, v := range arr {
		if inner, ok := v.([]interface{}); ok {
			result = append(result, flattenArray(inner, depth-1)...)
		} else {
			result = append(result, v)
		}
	}
	return result
}
