package r2lang

import "sort"

// r2hack.go: Funciones de "seguridad", "forense" y "análisis" para R2.
// Enfoque didáctico, no pretende ser una suite de hacking real.

func RegisterCollections(env *Environment) {

	// 1. Función append
	env.Set("append", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("append: necesita al menos 2 argumentos")
		}

		return append(args[0].([]interface{}), args[1:]...)
	}))

	// Funcion count
	env.Set("count", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) != 2 {
			panic("count: solo se aceptan 2 argumentos")
		}
		count := 0
		arr := args[0].([]interface{})
		for _, v := range arr {
			if v == args[1] {
				count++
			}
		}
		return count

	}))

	// funcion para calcular cuantoes elementos distintos tengo
	env.Set("distinct", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) != 1 {
			panic("count: solo se aceptan 1 argumentos")
		}
		count := 0
		mapCounter := make(map[interface{}]int)
		for _, v := range args[0].([]interface{}) {
			mapCounter[v.(interface{})]++
		}
		count = len(mapCounter)
		return count
	}))

	// Función len
	env.Set("len", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) != 1 {
			panic("len: solo se acepta un argumento")
		}
		v, ok := args[0].([]interface{})
		if ok {
			return len(v)
		}
		v1, ok := args[0].(string)
		if ok {
			return len(v1)
		}

		v2 := args[0].(map[string]interface{})

		return len(v2)
	}))

	// funcions para calcular index
	env.Set("index", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) != 2 {
			panic("index: solo se aceptan 2 argumentos")
		}
		index := -1
		arr := args[0].([]interface{})
		for i, v := range arr {
			if v == args[1] {
				index = i
				break
			}
		}
		return index
	}))

	// funcions para calcular index
	env.Set("indexes", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) != 2 {
			panic("index: solo se aceptan 2 argumentos")
		}
		index := -1
		arr := args[0].([]interface{})
		indexes := make([]int, 0)
		for i, v := range arr {
			if v == args[1] {
				index = i
				indexes = append(indexes, index)
			}
		}
		return indexes
	}))

	// funcion para generar rangos como range de golang
	env.Set("range", BuiltinFunction(func(args ...interface{}) interface{} {
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

	env.Set("repeat", BuiltinFunction(func(args ...interface{}) interface{} {
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

	// funcion insert
	env.Set("insert", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 3 {
			panic("insert: necesita al menos 3 argumentos")
		}
		arr := args[0].([]interface{})
		index := int(args[2].(float64))
		if index < 0 || index > len(arr) {
			panic("insert: index fuera de rango")
		}
		value := args[1]
		arr = append(arr[:index], append([]interface{}{value}, arr[index:]...)...)
		return arr
	}))

	// 3. Función delete
	env.Set("delete", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("delete: necesita al menos 2 argumentos")
		}

		arr := args[0].([]interface{})
		indexes := make([]int, len(args)-1)
		for i := 1; i < len(args); i++ {
			index := int(args[i].(float64))
			if index < 0 || index >= len(arr) {
				continue
			}
			indexes[i-1] = index
		}
		if len(indexes) == 0 {
			return arr
		}
		if len(indexes) == 1 {
			arr = append(arr[:indexes[0]], arr[indexes[0]+1:]...)
			return arr
		}

		sort.Ints(indexes)
		fix := 0
		for i, index := range indexes {
			if i > 0 {
				if index == indexes[i-1] {
					continue
				}
			}
			pos := index - fix
			if pos < 0 || pos >= len(arr) {
				panic("delete: index fuera de rango")
			}
			arr = append(arr[:pos], arr[pos+1:]...)
			fix++
		}
		return arr
	}))

	// Función copy
	env.Set("copy", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) != 1 {
			panic("copy: solo se acepta un argumento")
		}

		arr := args[0].([]interface{})
		newArr := make([]interface{}, len(arr))
		copy(newArr, arr)

		return newArr
	}))

	// Función slice
	env.Set("slice", BuiltinFunction(func(args ...interface{}) interface{} {
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

	// Función reverse
	env.Set("reverse", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) != 1 {
			panic("reverse: solo se acepta un argumento")
		}

		arr := args[0].([]interface{})
		newArr := make([]interface{}, len(arr))

		for i := 0; i < len(arr); i++ {
			newArr[i] = arr[len(arr)-1-i]
		}

		return newArr
	}))

	// Función sort
	env.Set("sort", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) != 1 {
			panic("sort: solo se acepta un argumento")
		}

		arr := args[0].([]interface{})
		newArr := make([]interface{}, len(arr))
		copy(newArr, arr)

		quickSort(newArr, 0, len(newArr)-1)

		return newArr
	}))

}

func quickSort(arr []interface{}, low, high int) {
	if low < high {
		p := partition(arr, low, high)
		quickSort(arr, low, p-1)
		quickSort(arr, p+1, high)
	}
}

func partition(arr []interface{}, low, high int) int {
	pivot := int(arr[high].(float64))
	i := low

	for j := low; j < high; j++ {
		if int(arr[j].(float64)) < pivot {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}

	arr[i], arr[high] = arr[high], arr[i]

	return i
}
