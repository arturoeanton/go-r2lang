package r2libs

import (
	"fmt"
	"regexp"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func RegisterRegex(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		"test": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("regex.test necesita (pattern, str)")
			}
			pattern, okP := args[0].(string)
			str, okS := args[1].(string)
			if !(okP && okS) {
				panic("regex.test: (pattern, str) => strings")
			}
			re, err := regexp.Compile(pattern)
			if err != nil {
				panic(fmt.Sprintf("regex.test: invalid pattern: %v", err))
			}
			return re.MatchString(str)
		}),

		"match": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("regex.match necesita (pattern, str)")
			}
			pattern, okP := args[0].(string)
			str, okS := args[1].(string)
			if !(okP && okS) {
				panic("regex.match: (pattern, str) => strings")
			}
			re, err := regexp.Compile(pattern)
			if err != nil {
				panic(fmt.Sprintf("regex.match: invalid pattern: %v", err))
			}
			loc := re.FindStringIndex(str)
			if loc == nil {
				return nil
			}
			return str[loc[0]:loc[1]]
		}),

		"matchAll": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("regex.matchAll necesita (pattern, str)")
			}
			pattern, okP := args[0].(string)
			str, okS := args[1].(string)
			if !(okP && okS) {
				panic("regex.matchAll: (pattern, str) => strings")
			}
			re, err := regexp.Compile(pattern)
			if err != nil {
				panic(fmt.Sprintf("regex.matchAll: invalid pattern: %v", err))
			}
			matches := re.FindAllString(str, -1)
			result := make([]interface{}, len(matches))
			for i, m := range matches {
				result[i] = m
			}
			return result
		}),

		"groups": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("regex.groups necesita (pattern, str)")
			}
			pattern, okP := args[0].(string)
			str, okS := args[1].(string)
			if !(okP && okS) {
				panic("regex.groups: (pattern, str) => strings")
			}
			re, err := regexp.Compile(pattern)
			if err != nil {
				panic(fmt.Sprintf("regex.groups: invalid pattern: %v", err))
			}
			match := re.FindStringSubmatch(str)
			if match == nil {
				return nil
			}
			result := make([]interface{}, len(match))
			for i, g := range match {
				result[i] = g
			}
			return result
		}),

		// replace substitutes only the first match. The replacement string
		// follows regexp.ReplaceAllString semantics: $1, $2, ... (or
		// ${name}) refer to capture groups from the pattern.
		"replace": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 3 {
				panic("regex.replace necesita (pattern, str, replacement)")
			}
			pattern, okP := args[0].(string)
			str, okS := args[1].(string)
			replacement, okR := args[2].(string)
			if !(okP && okS && okR) {
				panic("regex.replace: (pattern, str, replacement) => strings")
			}
			re, err := regexp.Compile(pattern)
			if err != nil {
				panic(fmt.Sprintf("regex.replace: invalid pattern: %v", err))
			}
			loc := re.FindStringIndex(str)
			if loc == nil {
				return str
			}
			replaced := re.ReplaceAllString(str[loc[0]:loc[1]], replacement)
			return str[:loc[0]] + replaced + str[loc[1]:]
		}),

		// replaceAll substitutes every match. See regex.replace for the
		// $1, $2, ... capture group syntax supported in replacement.
		"replaceAll": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 3 {
				panic("regex.replaceAll necesita (pattern, str, replacement)")
			}
			pattern, okP := args[0].(string)
			str, okS := args[1].(string)
			replacement, okR := args[2].(string)
			if !(okP && okS && okR) {
				panic("regex.replaceAll: (pattern, str, replacement) => strings")
			}
			re, err := regexp.Compile(pattern)
			if err != nil {
				panic(fmt.Sprintf("regex.replaceAll: invalid pattern: %v", err))
			}
			return re.ReplaceAllString(str, replacement)
		}),

		"split": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("regex.split necesita (pattern, str)")
			}
			pattern, okP := args[0].(string)
			str, okS := args[1].(string)
			if !(okP && okS) {
				panic("regex.split: (pattern, str) => strings")
			}
			re, err := regexp.Compile(pattern)
			if err != nil {
				panic(fmt.Sprintf("regex.split: invalid pattern: %v", err))
			}
			parts := re.Split(str, -1)
			result := make([]interface{}, len(parts))
			for i, p := range parts {
				result[i] = p
			}
			return result
		}),

		"escape": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("regex.escape necesita (str)")
			}
			str, ok := args[0].(string)
			if !ok {
				panic("regex.escape: argumento debe ser string")
			}
			return regexp.QuoteMeta(str)
		}),
	}

	RegisterModule(env, "regex", functions)
}
