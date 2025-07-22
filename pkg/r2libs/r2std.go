package r2libs

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// r2std.go: Librería estándar con diversas funciones auxiliares

func RegisterStd(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		"typeOf": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				return "nil"
			}
			val := args[0]
			return fmt.Sprintf("%T", val)
		}),

		"len": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
		}),

		"sleep": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("sleep needs 1 argument (seconds)")
			}
			secs, ok := args[0].(float64)
			if !ok {
				panic("sleep: arg should be a number")
			}
			time.Sleep(time.Duration(secs) * time.Second)
			return nil
		}),

		"parseInt": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
		}),

		"parseFloat": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
		}),

		"toString": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("toString needs 1 argument")
			}
			return fmt.Sprint(args[0])
		}),

		"deepCopy": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("deepCopy needs 1 argument")
			}
			return deepCopy(args[0])
		}),

		"is": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("is needs 2 arguments (value, typeString)")
			}
			val := args[0]
			typeStr, ok := args[1].(string)
			if !ok {
				panic("is: second argument must be a string")
			}
			return isType(val, typeStr)
		}),

		"range": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
		}),

		"now": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
		}),

		"join": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
		}),

		"split": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
		}),

		"contains": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("contains needs 2 arguments: (str, substr)")
			}
			s, ok1 := args[0].(string)
			substr, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("contains: args should be string, string")
			}
			return strings.Contains(s, substr)
		}),

		"replace": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 3 {
				panic("replace needs 3 arguments: (str, old, new)")
			}
			s, ok1 := args[0].(string)
			old, ok2 := args[1].(string)
			new, ok3 := args[2].(string)
			if !ok1 || !ok2 || !ok3 {
				panic("replace: args should be string, string, string")
			}
			return strings.ReplaceAll(s, old, new)
		}),

		"toUpperCase": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("toUpperCase needs 1 argument")
			}
			s, ok := args[0].(string)
			if !ok {
				panic("toUpperCase: arg should be string")
			}
			return strings.ToUpper(s)
		}),

		"toLowerCase": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("toLowerCase needs 1 argument")
			}
			s, ok := args[0].(string)
			if !ok {
				panic("toLowerCase: arg should be string")
			}
			return strings.ToLower(s)
		}),

		"eval": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
		}),

		"keys": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
		}),

		"print": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			for i, arg := range args {
				if i > 0 {
					fmt.Print(" ")
				}
				fmt.Print(arg)
			}
			fmt.Println()
			return nil
		}),

		// P6 Feature: Partial Application and Currying
		"curry":   r2core.BuiltinFunction(r2core.CurryFunction),
		"partial": r2core.BuiltinFunction(r2core.PartialBuiltin),
	}

	RegisterModule(env, "std", functions)
}

// Helper for deepCopy (recursive)
func deepCopy(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Ptr:
		// Dereference pointer and copy recursively
		return deepCopy(v.Elem().Interface())
	case reflect.Map:
		newMap := make(map[string]interface{}, v.Len())
		for _, key := range v.MapKeys() {
			strKey := fmt.Sprint(key.Interface())
			newMap[strKey] = deepCopy(v.MapIndex(key).Interface())
		}
		return newMap
	case reflect.Slice:
		// Handle byte slices separately to avoid treating them as generic arrays
		if v.Type().Elem().Kind() == reflect.Uint8 {
			newSlice := make([]byte, v.Len())
			reflect.Copy(reflect.ValueOf(newSlice), v)
			return newSlice
		}
		newSlice := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			newSlice[i] = deepCopy(v.Index(i).Interface())
		}
		return newSlice
	case reflect.Array:
		newArray := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			newArray[i] = deepCopy(v.Index(i).Interface())
		}
		return newArray
	default:
		return value // Return value directly for primitive types
	}
}

// Helper for is
func isType(value interface{}, typeString string) bool {
	switch typeString {
	case "number", "float", "float64":
		_, ok := value.(float64)
		return ok
	case "string":
		_, ok := value.(string)
		return ok
	case "bool", "boolean":
		_, ok := value.(bool)
		return ok
	case "array":
		_, ok := value.([]interface{})
		return ok
	case "map", "object":
		_, ok := value.(map[string]interface{})
		return ok
	case "function":
		_, ok := value.(*r2core.UserFunction)
		return ok
	case "nil", "null":
		return value == nil
	case "date":
		_, ok := value.(*r2core.DateValue)
		return ok
	case "duration":
		_, ok := value.(*r2core.DurationValue)
		return ok
	default:
		return false
	}
}
