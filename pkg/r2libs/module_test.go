package r2libs

import (
	"github.com/arturoeanton/go-r2lang/pkg/r2core"
	"testing"
)

func TestModuleSystem(t *testing.T) {
	env := r2core.NewEnvironment()
	env.Set("true", true)
	env.Set("false", false)
	env.Set("nil", nil)

	// Register all modules
	RegisterStd(env)
	RegisterPrint(env)
	RegisterMath(env)
	RegisterString(env)
	RegisterIO(env)
	RegisterOS(env)
	RegisterRand(env)
	RegisterTest(env)

	tests := []struct {
		name     string
		code     string
		expected interface{}
	}{
		{
			name:     "std.print function exists",
			code:     `std.print("hello"); return std.typeOf(std.print);`,
			expected: "r2core.BuiltinFunction",
		},
		{
			name:     "std.len function works",
			code:     `return std.len("hello");`,
			expected: float64(5),
		},
		{
			name:     "std.typeOf function works",
			code:     `return std.typeOf(123);`,
			expected: "float64",
		},
		{
			name:     "std.parseInt function works",
			code:     `return std.parseInt("42");`,
			expected: float64(42),
		},
		{
			name:     "std.keys function works",
			code:     `let obj = {"a": 1, "b": 2}; return std.len(std.keys(obj));`,
			expected: float64(2),
		},
		{
			name:     "math module functions exist",
			code:     `return std.typeOf(math.sin);`,
			expected: "r2core.BuiltinFunction",
		},
		{
			name:     "string module functions exist",
			code:     `return std.typeOf(string.toUpper);`,
			expected: "r2core.BuiltinFunction",
		},
		{
			name:     "print module functions exist",
			code:     `return std.typeOf(print.printBox);`,
			expected: "r2core.BuiltinFunction",
		},
		{
			name:     "io module functions exist",
			code:     `return std.typeOf(io.readFile);`,
			expected: "r2core.BuiltinFunction",
		},
		{
			name:     "os module functions exist",
			code:     `return std.typeOf(os.getEnv);`,
			expected: "r2core.BuiltinFunction",
		},
		{
			name:     "rand module functions exist",
			code:     `return std.typeOf(rand.randFloat);`,
			expected: "r2core.BuiltinFunction",
		},
		{
			name:     "test module functions exist",
			code:     `return std.typeOf(test.assertEq);`,
			expected: "r2core.BuiltinFunction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := r2core.NewParser(tt.code)
			program := parser.ParseProgram()
			result := program.Eval(env)

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestModuleNamespaceIsolation(t *testing.T) {
	env := r2core.NewEnvironment()
	env.Set("true", true)
	env.Set("false", false)
	env.Set("nil", nil)

	// Register all modules
	RegisterStd(env)
	RegisterMath(env)
	RegisterString(env)

	// Test that functions are properly namespaced
	tests := []struct {
		name      string
		code      string
		shouldErr bool
	}{
		{
			name:      "std.print works",
			code:      `std.print("hello"); return "ok";`,
			shouldErr: false,
		},
		{
			name:      "math.sin works",
			code:      `return math.sin(0);`,
			shouldErr: false,
		},
		{
			name:      "string.toUpper works",
			code:      `return string.toUpper("hello");`,
			shouldErr: false,
		},
		{
			name:      "old print function should not exist",
			code:      `print("hello"); return "ok";`,
			shouldErr: true,
		},
		{
			name:      "old len function should not exist",
			code:      `return len("hello");`,
			shouldErr: true,
		},
		{
			name:      "old typeOf function should not exist",
			code:      `return typeOf(123);`,
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.shouldErr {
						t.Errorf("Expected no error, but got panic: %v", r)
					}
				}
			}()

			parser := r2core.NewParser(tt.code)
			program := parser.ParseProgram()
			result := program.Eval(env)

			if tt.shouldErr {
				t.Errorf("Expected error, but got result: %v", result)
			}
		})
	}
}

func TestModuleAccess(t *testing.T) {
	env := r2core.NewEnvironment()
	env.Set("true", true)
	env.Set("false", false)
	env.Set("nil", nil)

	RegisterStd(env)

	// Test accessing module as object
	code := `
		let stdModule = std;
		return stdModule.len("hello");
	`

	parser := r2core.NewParser(code)
	program := parser.ParseProgram()
	result := program.Eval(env)

	if result != float64(5) {
		t.Errorf("Expected 5, got %v", result)
	}
}
