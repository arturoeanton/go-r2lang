package r2libs

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// r2std.go: Librería estándar con diversas funciones auxiliares

func RegisterStd(env *r2core.Environment) {

	env.Set("typeOf", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			return "nil"
		}
		val := args[0]
		return fmt.Sprintf("%T", val)
	}))

	env.Set("len", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("len needs 1 argument")
		}
		switch v := args[0].(type) {
		case string:
			return float64(len(v))
		case []interface{}:
			return float64(len(v))
		case map[string]interface{}:
			return float64(len(v))
		default:
			panic("len: Expected string, []interface{}, or map[string]interface{}")
		}
	}))

	env.Set("sleep", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("sleep needs 1 argument (seconds)")
		}
		secs, ok := args[0].(float64)
		if !ok {
			panic("sleep: arg should be a number")
		}
		time.Sleep(time.Duration(secs) * time.Second)
		return nil
	}))

	env.Set("parseInt", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("parseInt needs 1 argument (string)")
		}
		s, ok := args[0].(string)
		if !ok {
			panic("parseInt: arg should be string")
		}
		i, err := strconv.Atoi(s)
		if err != nil {
			panic("parseInt: could not convert '" + s + "' to int")
		}
		return float64(i)
	}))

	env.Set("parseFloat", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("parseFloat needs 1 argument (string)")
		}
		s, ok := args[0].(string)
		if !ok {
			panic("parseFloat: arg should be string")
		}
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic("parseFloat: could not convert '" + s + "' to float")
		}
		return f
	}))

	// 6) toString(value): convierte un valor a string
	env.Set("toString", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("toString needs 1 argument")
		}
		return fmt.Sprint(args[0])
	}))

	env.Set("range", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("range needs 2 arguments: (start, end)")
		}
		start, ok1 := args[0].(float64)
		end, ok2 := args[1].(float64)
		if !ok1 || !ok2 {
			panic("range: arg should be number, number")
		}
		arr := []interface{}{}
		for i := int(start); float64(i) < end; i++ {
			arr = append(arr, float64(i))
		}
		return arr
	}))

	env.Set("now", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) == 0 {
			return time.Now().Format("2006-01-02 15:04:05")
		}
		if len(args) == 1 {
			format, ok := args[0].(string)
			if !ok {
				panic("now: arg should be string")
			}
			return time.Now().Format(format)
		}
		panic("now: too many arguments")
	}))

	env.Set("join", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("join needs 2 arguments: (array, separator)")
		}
		arr, ok := args[0].([]interface{})
		sep, ok2 := args[1].(string)
		if !ok || !ok2 {
			panic("join: args should be Array, string")
		}
		// Convertir cada elemento a string
		strArr := make([]string, len(arr))
		for i, v := range arr {
			strArr[i] = fmt.Sprint(v)
		}
		return strings.Join(strArr, sep)
	}))

	env.Set("split", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("split needs 2 arguments: (str, separator)")
		}
		s, ok1 := args[0].(string)
		sep, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			panic("split: args should be string, string")
		}
		parts := strings.Split(s, sep)
		arr := make([]interface{}, len(parts))
		for i, p := range parts {
			arr[i] = p
		}
		return arr
	}))

	env.Set("eval", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("eval needs 1 argument (code string)")
		}
		code, ok := args[0].(string)
		if !ok {
			panic("eval: argument should be a string")
		}

		// Crear un nuevo parser para el código dinámico
		parser := r2core.NewParser(code)
		program := parser.ParseProgram()

		// Evaluar en el contexto del entorno actual
		return program.Eval(env)
	}))

	env.Set("keys", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("keys needs 1 argument (map)")
		}
		mapVal, ok := args[0].(map[string]interface{})
		if !ok {
			panic("keys: argument should be a map")
		}

		// Crear array con todas las claves
		keys := make([]interface{}, 0, len(mapVal))
		for key := range mapVal {
			keys = append(keys, key)
		}
		return keys
	}))

}
