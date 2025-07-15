package r2libs

import (
	"testing"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func TestEvalFunction_Simple(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			"simple arithmetic",
			"eval(\"5 + 3\")",
			float64(8),
		},
		{
			"string concatenation",
			"eval(\"'Hello' + ' World'\")",
			"Hello World",
		},
		{
			"boolean expression",
			"eval(\"10 > 5\")",
			true,
		},
		{
			"variable access",
			"let x = 10; eval(\"x + 5\")",
			float64(15),
		},
		{
			"dynamic variable creation",
			"eval(\"let y = 42\"); y",
			float64(42),
		},
		{
			"ternary in eval",
			"eval(\"5 > 3 ? 100 : 200\")",
			float64(100),
		},
		{
			"multiple variables in eval",
			"let a = 5, b = 3; eval(\"a * b\")",
			float64(15),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := r2core.NewParser(tt.input)
			program := parser.ParseProgram()

			env := r2core.NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			RegisterStd(env)

			result := program.Eval(env)

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEvalFunction_Errors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			"no arguments",
			"eval()",
		},
		{
			"non-string argument",
			"eval(123)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("expected panic for input: %s", tt.input)
				}
			}()

			parser := r2core.NewParser(tt.input)
			program := parser.ParseProgram()

			env := r2core.NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			RegisterStd(env)

			program.Eval(env)
		})
	}
}
