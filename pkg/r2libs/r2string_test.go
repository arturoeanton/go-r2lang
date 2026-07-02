package r2libs

import (
	"testing"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func TestStringFunctions(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterString(env)

	stringModuleObj, ok := env.Get("string")
	if !ok {
		t.Fatal("string module not found")
	}
	stringModule := stringModuleObj.(map[string]interface{})

	containsFunc := stringModule["contains"].(r2core.BuiltinFunction)
	repeatFunc := stringModule["repeat"].(r2core.BuiltinFunction)
	padStartFunc := stringModule["padStart"].(r2core.BuiltinFunction)
	padEndFunc := stringModule["padEnd"].(r2core.BuiltinFunction)
	trimStartFunc := stringModule["trimStart"].(r2core.BuiltinFunction)
	trimEndFunc := stringModule["trimEnd"].(r2core.BuiltinFunction)
	capitalizeFunc := stringModule["capitalize"].(r2core.BuiltinFunction)
	isBlankFunc := stringModule["isBlank"].(r2core.BuiltinFunction)
	reverseFunc := stringModule["reverse"].(r2core.BuiltinFunction)

	tests := []struct {
		name     string
		function r2core.BuiltinFunction
		args     []interface{}
		expected interface{}
	}{
		{"contains found", containsFunc, []interface{}{"hello world", "world"}, true},
		{"contains not found", containsFunc, []interface{}{"hello world", "xyz"}, false},
		{"repeat basic", repeatFunc, []interface{}{"ab", 3.0}, "ababab"},
		{"repeat zero", repeatFunc, []interface{}{"ab", 0.0}, ""},
		{"padStart basic", padStartFunc, []interface{}{"5", 3.0, "0"}, "005"},
		{"padStart no-op when long enough", padStartFunc, []interface{}{"555", 2.0, "0"}, "555"},
		{"padStart multi-char pad", padStartFunc, []interface{}{"x", 5.0, "ab"}, "ababx"},
		{"padEnd basic", padEndFunc, []interface{}{"5", 3.0, "0"}, "500"},
		{"padEnd no-op when long enough", padEndFunc, []interface{}{"555", 2.0, "0"}, "555"},
		{"trimStart basic", trimStartFunc, []interface{}{"   hi  "}, "hi  "},
		{"trimEnd basic", trimEndFunc, []interface{}{"   hi  "}, "   hi"},
		{"capitalize basic", capitalizeFunc, []interface{}{"hello"}, "Hello"},
		{"capitalize empty", capitalizeFunc, []interface{}{""}, ""},
		{"isBlank true", isBlankFunc, []interface{}{"   "}, true},
		{"isBlank false", isBlankFunc, []interface{}{" a "}, false},
		{"reverse basic", reverseFunc, []interface{}{"hello"}, "olleh"},
		{"reverse empty", reverseFunc, []interface{}{""}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function(tt.args...)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestStringFunctionsPanics(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterString(env)
	stringModuleObj, _ := env.Get("string")
	stringModule := stringModuleObj.(map[string]interface{})

	repeatFunc := stringModule["repeat"].(r2core.BuiltinFunction)

	t.Run("repeat negative count panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic for negative count")
			}
		}()
		repeatFunc("ab", -1.0)
	})
}
