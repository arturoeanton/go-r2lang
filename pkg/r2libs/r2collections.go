package r2libs

import "github.com/arturoeanton/go-r2lang/pkg/r2core"

// r2hack.go: Funciones de "seguridad", "forense" y "análisis" para R2.
// Enfoque didáctico, no pretende ser una suite de hacking real.

func RegisterCollections(env *r2core.Environment) {

	// funcion para generar rangos como range de golang
	env.Set("range", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
	}))

	env.Set("repeat", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
	}))

	// Función copy
	env.Set("copy", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) != 1 {
			panic("copy: solo se acepta un argumento")
		}

		arr := args[0].([]interface{})
		newArr := make([]interface{}, len(arr))
		copy(newArr, arr)

		return newArr
	}))

	// Función slice
	env.Set("slice", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
	}))

}
