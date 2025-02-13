package r2lang

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// r2std.go: Librería estándar con diversas funciones auxiliares

func RegisterStd(env *Environment) {

	env.Set("typeOf", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			return "nil"
		}
		val := args[0]
		return fmt.Sprintf("%T", val)
	}))

	env.Set("len", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("len needs 1 argument")
		}
		switch v := args[0].(type) {
		case string:
			return float64(len(v))
		case []interface{}:
			return float64(len(v))
		default:
			panic("len: Expected string or []interface{}")
		}
	}))

	env.Set("sleep", BuiltinFunction(func(args ...interface{}) interface{} {
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

	env.Set("parseInt", BuiltinFunction(func(args ...interface{}) interface{} {
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

	env.Set("parseFloat", BuiltinFunction(func(args ...interface{}) interface{} {
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
	env.Set("toString", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("toString needs 1 argument")
		}
		return fmt.Sprint(args[0])
	}))

	env.Set("range", BuiltinFunction(func(args ...interface{}) interface{} {
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

	env.Set("now", BuiltinFunction(func(args ...interface{}) interface{} {
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

	env.Set("join", BuiltinFunction(func(args ...interface{}) interface{} {
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

	env.Set("split", BuiltinFunction(func(args ...interface{}) interface{} {
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

}
