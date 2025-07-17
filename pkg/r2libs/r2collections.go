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
			start := int(args[0].(float64))
			end := int(args[1].(float64))
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
			start := 0
			end := int(args[0].(float64))
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

			arr := args[0].([]interface{})
			newArr := make([]interface{}, len(arr))
			copy(newArr, arr)

			return newArr
		}),

		"slice": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) != 3 {
				panic("slice: solo se aceptan 3 argumentos")
			}

			arr := args[0].([]interface{})
			start := args[1].(int)
			end := args[2].(int)

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
			arr := args[0].([]interface{})
			fn := args[1].(*r2core.UserFunction)
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
			arr := args[0].([]interface{})
			fn := args[1].(*r2core.UserFunction)
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
			arr := args[0].([]interface{})
			fn := args[1].(*r2core.UserFunction)
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
			arr := args[0].([]interface{})
			fn := args[1].(*r2core.UserFunction)
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
			arr := args[0].([]interface{})
			val := args[1]
			for _, v := range arr {
				if equals(v, val) {
					return true
				}
			}
			return false
		}),
	}

	RegisterModule(env, "collections", functions)
}
