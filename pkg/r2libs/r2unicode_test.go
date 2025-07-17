package r2libs

import (
	"testing"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func TestUnicodeBasicFunctions(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterUnicode(env)

	unicodeModuleObj, ok := env.Get("unicode")
	if !ok {
		t.Fatal("unicode module not found")
	}
	unicodeModule := unicodeModuleObj.(map[string]interface{})

	ulenFunc := unicodeModule["ulen"].(r2core.BuiltinFunction)
	usubstrFunc := unicodeModule["usubstr"].(r2core.BuiltinFunction)
	uupperFunc := unicodeModule["uupper"].(r2core.BuiltinFunction)
	ulowerFunc := unicodeModule["ulower"].(r2core.BuiltinFunction)
	ureverseFunc := unicodeModule["ureverse"].(r2core.BuiltinFunction)

	tests := []struct {
		name     string
		function r2core.BuiltinFunction
		args     []interface{}
		expected interface{}
	}{
		{
			"ulen with ASCII",
			ulenFunc,
			[]interface{}{"hello"},
			float64(5),
		},
		{
			"ulen with Spanish characters",
			ulenFunc,
			[]interface{}{"José María"},
			float64(10),
		},
		{
			"ulen with emoji",
			ulenFunc,
			[]interface{}{"👋"},
			float64(1),
		},
		{
			"usubstr basic",
			usubstrFunc,
			[]interface{}{"hello", float64(1), float64(3)},
			"ell",
		},
		{
			"usubstr with Spanish",
			usubstrFunc,
			[]interface{}{"José María", float64(0), float64(4)},
			"José",
		},
		{
			"usubstr with emoji",
			usubstrFunc,
			[]interface{}{"Hello 👋 World", float64(6), float64(1)},
			"👋",
		},
		{
			"uupper with ASCII",
			uupperFunc,
			[]interface{}{"hello"},
			"HELLO",
		},
		{
			"uupper with Spanish",
			uupperFunc,
			[]interface{}{"josé"},
			"JOSÉ",
		},
		{
			"ulower with ASCII",
			ulowerFunc,
			[]interface{}{"HELLO"},
			"hello",
		},
		{
			"ulower with Spanish",
			ulowerFunc,
			[]interface{}{"JOSÉ"},
			"josé",
		},
		{
			"ureverse with ASCII",
			ureverseFunc,
			[]interface{}{"hello"},
			"olleh",
		},
		{
			"ureverse with Spanish",
			ureverseFunc,
			[]interface{}{"José"},
			"ésoJ",
		},
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
