package r2core

import (
	"testing"
)

func TestMultipleLetStatement_BasicDeclarations(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]interface{}
	}{
		{
			"two variables with values",
			"let a = 1, b = 2;",
			map[string]interface{}{"a": float64(1), "b": float64(2)},
		},
		{
			"three variables with values",
			"let x = 10, y = 20, z = 30;",
			map[string]interface{}{"x": float64(10), "y": float64(20), "z": float64(30)},
		},
		{
			"mixed types",
			"let name = \"John\", age = 25, active = true;",
			map[string]interface{}{"name": "John", "age": float64(25), "active": true},
		},
		{
			"some without values",
			"let a = 1, b, c = 3;",
			map[string]interface{}{"a": float64(1), "b": nil, "c": float64(3)},
		},
		{
			"all without values",
			"let x, y, z;",
			map[string]interface{}{"x": nil, "y": nil, "z": nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			program := parser.ParseProgram()

			if len(program.Statements) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(program.Statements))
			}

			env := NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			program.Eval(env)

			// Verificar que todas las variables esperadas están definidas
			for name, expectedValue := range tt.expected {
				actualValue, exists := env.Get(name)
				if !exists {
					t.Errorf("variable %s not found in environment", name)
					continue
				}

				if actualValue != expectedValue {
					t.Errorf("for variable %s: expected %v, got %v", name, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestMultipleLetStatement_WithExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]interface{}
	}{
		{
			"arithmetic expressions",
			"let a = 5 + 3, b = 10 - 2, c = 4 * 2;",
			map[string]interface{}{"a": float64(8), "b": float64(8), "c": float64(8)},
		},
		{
			"string concatenation",
			"let greeting = \"Hello\" + \" World\", name = \"R2\" + \"Lang\";",
			map[string]interface{}{"greeting": "Hello World", "name": "R2Lang"},
		},
		{
			"boolean expressions",
			"let isTrue = 5 > 3, isFalse = 2 > 5, isEqual = 10 == 10;",
			map[string]interface{}{"isTrue": true, "isFalse": false, "isEqual": true},
		},
		{
			"mixed expressions",
			"let sum = 1 + 2, text = \"Value: \" + 42, valid = true;",
			map[string]interface{}{"sum": float64(3), "text": "Value: 42", "valid": true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			program := parser.ParseProgram()

			env := NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			program.Eval(env)

			for name, expectedValue := range tt.expected {
				actualValue, exists := env.Get(name)
				if !exists {
					t.Errorf("variable %s not found in environment", name)
					continue
				}

				if actualValue != expectedValue {
					t.Errorf("for variable %s: expected %v, got %v", name, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestMultipleLetStatement_WithDependencies(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]interface{}
	}{
		{
			"using previously declared variable",
			"let a = 10; let b = a + 5, c = b * 2;",
			map[string]interface{}{"a": float64(10), "b": float64(15), "c": float64(30)},
		},
		{
			"chain of dependencies",
			"let x = 1; let y = x + 1, z = y + 1, w = z + 1;",
			map[string]interface{}{"x": float64(1), "y": float64(2), "z": float64(3), "w": float64(4)},
		},
		{
			"mixed single and multiple declarations",
			"let base = 100; let increment = 10, total = base + increment;",
			map[string]interface{}{"base": float64(100), "increment": float64(10), "total": float64(110)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			program := parser.ParseProgram()

			env := NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			program.Eval(env)

			for name, expectedValue := range tt.expected {
				actualValue, exists := env.Get(name)
				if !exists {
					t.Errorf("variable %s not found in environment", name)
					continue
				}

				if actualValue != expectedValue {
					t.Errorf("for variable %s: expected %v, got %v", name, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestMultipleLetStatement_BackwardCompatibility(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			"single declaration still works",
			"let x = 42;",
			float64(42),
		},
		{
			"single declaration without value",
			"let y;",
			nil,
		},
		{
			"single declaration with expression",
			"let result = 10 + 5 * 2;",
			float64(30),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			program := parser.ParseProgram()

			if len(program.Statements) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(program.Statements))
			}

			env := NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			program.Eval(env)

			// Para declaraciones simples, verificar la primera variable definida
			var varName string
			if tt.name == "single declaration still works" {
				varName = "x"
			} else if tt.name == "single declaration without value" {
				varName = "y"
			} else if tt.name == "single declaration with expression" {
				varName = "result"
			}

			actualValue, exists := env.Get(varName)
			if !exists {
				t.Errorf("variable %s not found in environment", varName)
				return
			}

			if actualValue != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, actualValue)
			}
		})
	}
}

func TestMultipleLetStatement_WithTernaryOperator(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]interface{}
	}{
		{
			"ternary in multiple declarations",
			"let condition = true, result = condition ? \"yes\" : \"no\", value = 10 > 5 ? 100 : 0;",
			map[string]interface{}{"condition": true, "result": "yes", "value": float64(100)},
		},
		{
			"complex ternary expressions",
			"let age = 25, status = age >= 18 ? \"adult\" : \"minor\", discount = age > 65 ? 0.2 : 0.0;",
			map[string]interface{}{"age": float64(25), "status": "adult", "discount": float64(0.0)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			program := parser.ParseProgram()

			env := NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			program.Eval(env)

			for name, expectedValue := range tt.expected {
				actualValue, exists := env.Get(name)
				if !exists {
					t.Errorf("variable %s not found in environment", name)
					continue
				}

				if actualValue != expectedValue {
					t.Errorf("for variable %s: expected %v, got %v", name, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestMultipleLetStatement_WithTemplateLiterals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]interface{}
	}{
		{
			"template literals in multiple declarations",
			"let name = \"World\", greeting = `Hello ${name}!`, count = 5, message = `Count: ${count}`;",
			map[string]interface{}{"name": "World", "greeting": "Hello World!", "count": float64(5), "message": "Count: 5"},
		},
		{
			"complex template expressions",
			"let x = 10, y = 20, sum = x + y, report = `Sum of ${x} and ${y} is ${sum}`;",
			map[string]interface{}{"x": float64(10), "y": float64(20), "sum": float64(30), "report": "Sum of 10 and 20 is 30"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			program := parser.ParseProgram()

			env := NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			program.Eval(env)

			for name, expectedValue := range tt.expected {
				actualValue, exists := env.Get(name)
				if !exists {
					t.Errorf("variable %s not found in environment", name)
					continue
				}

				if actualValue != expectedValue {
					t.Errorf("for variable %s: expected %v, got %v", name, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestMultipleLetStatement_InFunctions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			"multiple declarations in function",
			`func test() { let a = 1, b = 2, c = a + b; return c; } test();`,
			float64(3),
		},
		{
			"multiple declarations with function parameters",
			`func calculate(x, y) { let sum = x + y, product = x * y, result = sum + product; return result; } calculate(3, 4);`,
			float64(19), // (3+4) + (3*4) = 7 + 12 = 19
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			program := parser.ParseProgram()

			env := NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			result := program.Eval(env)

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
