package r2libs

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// r2string.go: Funciones nativas para manipulación de strings en R2

func RegisterString(env *r2core.Environment) {
	// toUpper(str) => STR
	env.Set("toUpper", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("toUpper necesita (str)")
		}
		s, ok := args[0].(string)
		if !ok {
			panic("toUpper: argumento debe ser string")
		}
		return strings.ToUpper(s)
	}))

	// toLower(str) => STR
	env.Set("toLower", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("toLower necesita (str)")
		}
		s, ok := args[0].(string)
		if !ok {
			panic("toLower: argumento debe ser string")
		}
		return strings.ToLower(s)
	}))

	// trim(str) => STR  (quita espacios al inicio/fin)
	env.Set("trim", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("trim necesita (str)")
		}
		s, ok := args[0].(string)
		if !ok {
			panic("trim: argumento debe ser string")
		}
		return strings.TrimSpace(s)
	}))

	// substring(str, start, length) => STR
	env.Set("substring", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 3 {
			panic("substring necesita (str, start, length)")
		}
		s, okS := args[0].(string)
		startF, ok1 := args[1].(float64)
		lengthF, ok2 := args[2].(float64)
		if !(okS && ok1 && ok2) {
			panic("substring: (str, start, length) => str y numéricos")
		}
		start := int(startF)
		length := int(lengthF)
		if start < 0 || length < 0 || start > len(s) {
			return "" // o panic, a tu elección
		}
		end := start + length
		if end > len(s) {
			end = len(s)
		}
		return s[start:end]
	}))

	// indexOf(str, sub) => number (primera ocurrencia) o -1
	env.Set("indexOf", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("indexOf necesita (str, sub)")
		}
		s, okS := args[0].(string)
		sub, okSub := args[1].(string)
		if !(okS && okSub) {
			panic("indexOf: argumentos deben ser strings")
		}
		idx := strings.Index(s, sub)
		if idx < 0 {
			return float64(-1)
		}
		return float64(idx)
	}))

	// lastIndexOf(str, sub) => number (última ocurrencia) o -1
	env.Set("lastIndexOf", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("lastIndexOf necesita (str, sub)")
		}
		s, okS := args[0].(string)
		sub, okSub := args[1].(string)
		if !(okS && okSub) {
			panic("lastIndexOf: argumentos deben ser strings")
		}
		idx := strings.LastIndex(s, sub)
		if idx < 0 {
			return float64(-1)
		}
		return float64(idx)
	}))

	// startsWith(str, prefix) => bool
	env.Set("startsWith", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("startsWith necesita (str, prefix)")
		}
		s, okS := args[0].(string)
		prefix, okP := args[1].(string)
		if !(okS && okP) {
			panic("startsWith: argumentos deben ser strings")
		}
		return strings.HasPrefix(s, prefix)
	}))

	// endsWith(str, suffix) => bool
	env.Set("endsWith", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("endsWith necesita (str, suffix)")
		}
		s, okS := args[0].(string)
		suffix, okP := args[1].(string)
		if !(okS && okP) {
			panic("endsWith: argumentos deben ser strings")
		}
		return strings.HasSuffix(s, suffix)
	}))

	// replace(str, old, new) => STR  (reemplaza todas las ocurrencias)
	env.Set("replace", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 3 {
			panic("replace necesita (str, old, new)")
		}
		s, sOk := args[0].(string)
		oldS, oOk := args[1].(string)
		newS, nOk := args[2].(string)
		if !(sOk && oOk && nOk) {
			panic("replace: (str, old, new) => strings")
		}
		return strings.ReplaceAll(s, oldS, newS)
	}))

	// split(str, sep) => array nativa
	env.Set("split", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("split necesita (str, sep)")
		}
		s, ok1 := args[0].(string)
		sep, ok2 := args[1].(string)
		if !(ok1 && ok2) {
			panic("split: argumentos deben ser strings")
		}
		parts := strings.Split(s, sep)
		arr := make([]interface{}, len(parts))
		for i, p := range parts {
			arr[i] = p
		}
		return arr
	}))

	// join(array, sep) => STR
	env.Set("join", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("join necesita (array, sep)")
		}
		arr, ok1 := args[0].([]interface{})
		sep, ok2 := args[1].(string)
		if !(ok1 && ok2) {
			panic("join: primer argumento array nativo, segundo un string")
		}
		strArr := make([]string, len(arr))
		for i, v := range arr {
			strArr[i] = fmt.Sprint(v) // conviertes a string
		}
		return strings.Join(strArr, sep)
	}))

	// lengthOfString(str) => number (cuenta runas /caracteres/ en UTF-8)
	env.Set("lengthOfString", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("lengthOfString necesita (str)")
		}
		s, ok := args[0].(string)
		if !ok {
			panic("lengthOfString: argumento debe ser string")
		}
		// length en caracteres (runas), no bytes
		return float64(utf8.RuneCountInString(s))
	}))
}
