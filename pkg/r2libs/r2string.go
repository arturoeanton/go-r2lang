package r2libs

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// r2string.go: Funciones nativas para manipulación de strings en R2

func RegisterString(env *r2core.Environment) {
	toUpper := r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("toUpper necesita (str)")
		}
		s, ok := args[0].(string)
		if !ok {
			panic("toUpper: argumento debe ser string")
		}
		return strings.ToUpper(s)
	})

	toLower := r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("toLower necesita (str)")
		}
		s, ok := args[0].(string)
		if !ok {
			panic("toLower: argumento debe ser string")
		}
		return strings.ToLower(s)
	})

	functions := map[string]r2core.BuiltinFunction{
		"toUpper": toUpper,
		// toUpperCase/toLowerCase: JS-style aliases for scripts written
		// against String.prototype naming conventions.
		"toUpperCase": toUpper,

		"toLower":     toLower,
		"toLowerCase": toLower,

		"trim": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("trim necesita (str)")
			}
			s, ok := args[0].(string)
			if !ok {
				panic("trim: argumento debe ser string")
			}
			return strings.TrimSpace(s)
		}),

		"substring": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
			runes := []rune(s)
			if start < 0 || length < 0 || start > len(runes) {
				return "" // o panic, a tu elección
			}
			end := start + length
			if end < start || end > len(runes) {
				end = len(runes)
			}
			return string(runes[start:end])
		}),

		"indexOf": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
		}),

		"lastIndexOf": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
		}),

		"startsWith": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("startsWith necesita (str, prefix)")
			}
			s, okS := args[0].(string)
			prefix, okP := args[1].(string)
			if !(okS && okP) {
				panic("startsWith: argumentos deben ser strings")
			}
			return strings.HasPrefix(s, prefix)
		}),

		"endsWith": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("endsWith necesita (str, suffix)")
			}
			s, okS := args[0].(string)
			suffix, okP := args[1].(string)
			if !(okS && okP) {
				panic("endsWith: argumentos deben ser strings")
			}
			return strings.HasSuffix(s, suffix)
		}),

		"replace": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
		}),

		"split": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
		}),

		"join": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
		}),

		"lengthOfString": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("lengthOfString necesita (str)")
			}
			s, ok := args[0].(string)
			if !ok {
				panic("lengthOfString: argumento debe ser string")
			}
			// length en caracteres (runas), no bytes
			return float64(utf8.RuneCountInString(s))
		}),

		"contains": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("contains necesita (str, sub)")
			}
			s, okS := args[0].(string)
			sub, okSub := args[1].(string)
			if !(okS && okSub) {
				panic("contains: argumentos deben ser strings")
			}
			return strings.Contains(s, sub)
		}),

		"repeat": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("repeat necesita (str, count)")
			}
			s, okS := args[0].(string)
			countF, okC := args[1].(float64)
			if !(okS && okC) {
				panic("repeat: (str, count) => str string y count numérico")
			}
			count := int(countF)
			if count < 0 {
				panic("repeat: count no puede ser negativo")
			}
			return strings.Repeat(s, count)
		}),

		"padStart": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 3 {
				panic("padStart necesita (str, targetLength, padStr)")
			}
			s, okS := args[0].(string)
			targetF, okT := args[1].(float64)
			pad, okP := args[2].(string)
			if !(okS && okT && okP) {
				panic("padStart: (str, targetLength, padStr) => str y padStr strings, targetLength numérico")
			}
			target := int(targetF)
			runes := utf8.RuneCountInString(s)
			if target <= runes || pad == "" {
				return s
			}
			padRunes := []rune(pad)
			needed := target - runes
			fill := make([]rune, needed)
			for i := 0; i < needed; i++ {
				fill[i] = padRunes[i%len(padRunes)]
			}
			return string(fill) + s
		}),

		"padEnd": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 3 {
				panic("padEnd necesita (str, targetLength, padStr)")
			}
			s, okS := args[0].(string)
			targetF, okT := args[1].(float64)
			pad, okP := args[2].(string)
			if !(okS && okT && okP) {
				panic("padEnd: (str, targetLength, padStr) => str y padStr strings, targetLength numérico")
			}
			target := int(targetF)
			runes := utf8.RuneCountInString(s)
			if target <= runes || pad == "" {
				return s
			}
			padRunes := []rune(pad)
			needed := target - runes
			fill := make([]rune, needed)
			for i := 0; i < needed; i++ {
				fill[i] = padRunes[i%len(padRunes)]
			}
			return s + string(fill)
		}),

		"trimStart": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("trimStart necesita (str)")
			}
			s, ok := args[0].(string)
			if !ok {
				panic("trimStart: argumento debe ser string")
			}
			return strings.TrimLeft(s, " \t\n\r")
		}),

		"trimEnd": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("trimEnd necesita (str)")
			}
			s, ok := args[0].(string)
			if !ok {
				panic("trimEnd: argumento debe ser string")
			}
			return strings.TrimRight(s, " \t\n\r")
		}),

		"capitalize": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("capitalize necesita (str)")
			}
			s, ok := args[0].(string)
			if !ok {
				panic("capitalize: argumento debe ser string")
			}
			runes := []rune(s)
			if len(runes) == 0 {
				return s
			}
			return strings.ToUpper(string(runes[0])) + string(runes[1:])
		}),

		"isBlank": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("isBlank necesita (str)")
			}
			s, ok := args[0].(string)
			if !ok {
				panic("isBlank: argumento debe ser string")
			}
			return strings.TrimSpace(s) == ""
		}),

		"reverse": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("reverse necesita (str)")
			}
			s, ok := args[0].(string)
			if !ok {
				panic("reverse: argumento debe ser string")
			}
			runes := []rune(s)
			for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
				runes[i], runes[j] = runes[j], runes[i]
			}
			return string(runes)
		}),
	}

	RegisterModule(env, "string", functions)
}
