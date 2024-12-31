package r2lang

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// r2std.go: Librería estándar con diversas funciones auxiliares

func RegisterStd(env *Environment) {

	// 1) typeOf(value): retorna un string con el tipo Go subyacente
	env.Set("typeOf", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			return "nil"
		}
		val := args[0]
		// Por simplicidad, retornamos reflect.TypeOf(val).String() o similar
		return fmt.Sprintf("%T", val)
	}))

	// 2) len(value): si es string => largo de la cadena; si es slice nativo => largo
	env.Set("len", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("len necesita 1 argumento")
		}
		switch v := args[0].(type) {
		case string:
			return float64(len(v))
		case []interface{}:
			return float64(len(v))
		default:
			panic("len: se esperaba string o array nativa")
		}
	}))

	// 3) sleep(segundos): duerme la ejecución `segundos`
	env.Set("sleep", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("sleep necesita 1 argumento (segundos)")
		}
		secs, ok := args[0].(float64)
		if !ok {
			panic("sleep: argumento debe ser un número (segundos)")
		}
		time.Sleep(time.Duration(secs) * time.Second)
		return nil
	}))

	// 4) parseInt(str): convierte un string a int => float64 (en R2)
	env.Set("parseInt", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("parseInt necesita 1 argumento (string)")
		}
		s, ok := args[0].(string)
		if !ok {
			panic("parseInt: argumento debe ser string")
		}
		i, err := strconv.Atoi(s)
		if err != nil {
			panic("parseInt: no pudo convertir '" + s + "' a int")
		}
		return float64(i)
	}))

	// 5) parseFloat(str): convierte un string a float64
	env.Set("parseFloat", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("parseFloat necesita 1 argumento (string)")
		}
		s, ok := args[0].(string)
		if !ok {
			panic("parseFloat: argumento debe ser string")
		}
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic("parseFloat: no pudo convertir '" + s + "' a float")
		}
		return f
	}))

	// 6) toString(value): convierte un valor a string
	env.Set("toString", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("toString necesita 1 argumento")
		}
		return fmt.Sprint(args[0])
	}))

	// 7) vars(map, key): obtiene map[key] (o nil si no existe)
	env.Set("vars", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("vars necesita 2 argumentos: (map, key)")
		}
		theMap, ok1 := args[0].(map[string]interface{})
		theKey, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			panic("vars: primer argumento debe ser un map<string, interface{}>, segundo un string")
		}
		val, found := theMap[theKey]
		if !found {
			return nil
		}
		return val
	}))

	// 8) varsSet(map, key, value): asigna map[key] = value
	env.Set("varsSet", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 3 {
			panic("varsSet necesita 3 argumentos: (map, key, value)")
		}
		theMap, ok1 := args[0].(map[string]interface{})
		theKey, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			panic("varsSet: primer arg debe ser map, segundo un string")
		}
		theMap[theKey] = args[2]
		return nil
	}))

	// 9) range(start, end): retorna un array nativo de floats [start, start+1, ..., end-1]
	// (Muy simplificado, sin step)
	env.Set("range", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("range necesita 2 argumentos: (start, end)")
		}
		start, ok1 := args[0].(float64)
		end, ok2 := args[1].(float64)
		if !ok1 || !ok2 {
			panic("range: argumentos deben ser numéricos")
		}
		arr := []interface{}{}
		for i := int(start); float64(i) < end; i++ {
			arr = append(arr, float64(i))
		}
		return arr
	}))

	// 10) now(): retorna la fecha/hora actual como string
	env.Set("now", BuiltinFunction(func(args ...interface{}) interface{} {
		return time.Now().Format("2006-01-02 15:04:05")
	}))

	// 11) join(array, separator): une un array de strings con el separador
	env.Set("join", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("join necesita 2 argumentos: (array, separator)")
		}
		arr, ok := args[0].([]interface{})
		sep, ok2 := args[1].(string)
		if !ok || !ok2 {
			panic("join: (array, separator) con array = []interface{} y separator = string")
		}
		// Convertir cada elemento a string
		strArr := make([]string, len(arr))
		for i, v := range arr {
			strArr[i] = fmt.Sprint(v)
		}
		return strings.Join(strArr, sep)
	}))

	// 12) split(str, separator): retorna un array nativo de strings
	env.Set("split", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("split necesita 2 argumentos: (str, separator)")
		}
		s, ok1 := args[0].(string)
		sep, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			panic("split: argumentos deben ser string, string")
		}
		parts := strings.Split(s, sep)
		arr := make([]interface{}, len(parts))
		for i, p := range parts {
			arr[i] = p
		}
		return arr
	}))

	// Listo. Puedes agregar más funciones según necesites
}
