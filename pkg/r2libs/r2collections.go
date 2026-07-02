package r2libs

import (
	"fmt"
	"reflect"
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

			arr, ok := toGenericSlice(args[0])
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

			// end (and start) are exclusive upper bounds, matching Go slice-
			// expression semantics: len(arr) is a valid value meaning
			// "through the last element". The previous "end >= len(arr)"
			// check rejected end == len(arr) too, making it impossible to
			// ever slice through the last element of the array.
			if start < 0 || start > len(arr) {
				panic("slice: start fuera de rango")
			}

			if end < 0 || end > len(arr) {
				panic("slice: end fuera de rango")
			}

			if start > end {
				panic("slice: start no puede ser mayor que end")
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
			if depth > maxFlattenDepth {
				panic("flatten: depth demasiado grande (máximo permitido: 10000)")
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
		"deepEqual": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) != 2 {
				panic("deepEqual: se aceptan 2 argumentos (a, b)")
			}
			return deepEqualValues(args[0], args[1], make(map[[2]uintptr]bool))
		}),

		"deepClone": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) != 1 {
				panic("deepClone: solo se acepta un argumento")
			}
			return deepCopy(args[0])
		}),
	}

	RegisterModule(env, "collections", functions)
}

func toGenericSlice(v interface{}) ([]interface{}, bool) {
	switch s := v.(type) {
	case []interface{}:
		return s, true
	case r2core.InterfaceSlice:
		return []interface{}(s), true
	}
	return nil, false
}

// deepEqualValues compares two R2Lang values structurally. seen tracks the
// active ancestor chain of (a, b) pointer pairs currently being compared, so
// that R2Lang arrays/maps made self-referential via index assignment (e.g.
// a[0] = a) cannot recurse forever and crash the process with an
// unrecoverable Go stack overflow, the same bug class fixed in
// collections.flatten. Revisiting the same pair mid-chain means both sides
// cycle in lockstep, so that branch is treated as equal.
func deepEqualValues(a, b interface{}, seen map[[2]uintptr]bool) bool {
	if aMap, ok := a.(map[string]interface{}); ok {
		bMap, ok := b.(map[string]interface{})
		if !ok || len(aMap) != len(bMap) {
			return false
		}
		if len(aMap) == 0 {
			return true
		}
		key := [2]uintptr{reflect.ValueOf(aMap).Pointer(), reflect.ValueOf(bMap).Pointer()}
		if seen[key] {
			return true
		}
		seen[key] = true
		defer delete(seen, key)
		for k, av := range aMap {
			bv, ok := bMap[k]
			if !ok || !deepEqualValues(av, bv, seen) {
				return false
			}
		}
		return true
	}

	if aArr, ok := toGenericSlice(a); ok {
		bArr, ok := toGenericSlice(b)
		if !ok || len(aArr) != len(bArr) {
			return false
		}
		if len(aArr) == 0 {
			return true
		}
		key := [2]uintptr{reflect.ValueOf(aArr).Pointer(), reflect.ValueOf(bArr).Pointer()}
		if seen[key] {
			return true
		}
		seen[key] = true
		defer delete(seen, key)
		for i := range aArr {
			if !deepEqualValues(aArr[i], bArr[i], seen) {
				return false
			}
		}
		return true
	}

	return equals(a, b)
}

// maxFlattenDepth bounds recursion in flattenArray so that a deeply nested
// (or self-referential, e.g. built via `a[0] = a`) array cannot trigger an
// unrecoverable Go "stack overflow" fatal error, which panic/recover cannot
// catch.
const maxFlattenDepth = 10000

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
