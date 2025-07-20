package r2core

import (
	"testing"
)

func TestArrayDestructuring_Basic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]interface{}
	}{
		{
			name:  "basic array destructuring",
			input: "let [a, b, c] = [1, 2, 3]; [a, b, c]",
			expected: map[string]interface{}{
				"a": float64(1),
				"b": float64(2),
				"c": float64(3),
			},
		},
		{
			name:  "more variables than elements",
			input: "let [a, b, c, d] = [1, 2]; [a, b, c, d]",
			expected: map[string]interface{}{
				"a": float64(1),
				"b": float64(2),
				"c": nil,
				"d": nil,
			},
		},
		{
			name:  "fewer variables than elements",
			input: "let [a, b] = [1, 2, 3, 4]; [a, b]",
			expected: map[string]interface{}{
				"a": float64(1),
				"b": float64(2),
			},
		},
		{
			name:  "empty array destructuring",
			input: "let [a, b] = []; [a, b]",
			expected: map[string]interface{}{
				"a": nil,
				"b": nil,
			},
		},
		{
			name:  "mixed types",
			input: "let [a, b, c] = [\"hello\", true, 42]; [a, b, c]",
			expected: map[string]interface{}{
				"a": "hello",
				"b": true,
				"c": float64(42),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			program := parser.ParseProgram()
			env := NewEnvironment()

			for _, stmt := range program.Statements {
				stmt.Eval(env)
			}

			// Verificar que las variables fueron asignadas correctamente
			for varName, expectedVal := range tt.expected {
				actualVal, exists := env.Get(varName)
				if !exists {
					t.Errorf("Variable %s was not found in environment", varName)
					continue
				}
				if actualVal != expectedVal {
					t.Errorf("Variable %s: expected %v (%T), got %v (%T)",
						varName, expectedVal, expectedVal, actualVal, actualVal)
				}
			}
		})
	}
}

func TestObjectDestructuring_Basic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]interface{}
	}{
		{
			name:  "basic object destructuring",
			input: "let {name, age} = {name: \"John\", age: 30}; [name, age]",
			expected: map[string]interface{}{
				"name": "John",
				"age":  float64(30),
			},
		},
		{
			name:  "missing properties",
			input: "let {name, age, city} = {name: \"John\"}; [name, age, city]",
			expected: map[string]interface{}{
				"name": "John",
				"age":  nil,
				"city": nil,
			},
		},
		{
			name:  "extra properties",
			input: "let {name} = {name: \"John\", age: 30, city: \"NYC\"}; name",
			expected: map[string]interface{}{
				"name": "John",
			},
		},
		{
			name:  "empty object destructuring",
			input: "let {name, age} = {}; [name, age]",
			expected: map[string]interface{}{
				"name": nil,
				"age":  nil,
			},
		},
		{
			name:  "nested values",
			input: "let {user, active} = {user: {name: \"John\"}, active: true}; [user, active]",
			expected: map[string]interface{}{
				"user":   map[string]interface{}{"name": "John"},
				"active": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			program := parser.ParseProgram()
			env := NewEnvironment()

			for _, stmt := range program.Statements {
				stmt.Eval(env)
			}

			// Verificar que las variables fueron asignadas correctamente
			for varName, expectedVal := range tt.expected {
				actualVal, exists := env.Get(varName)
				if !exists {
					t.Errorf("Variable %s was not found in environment", varName)
					continue
				}

				// Comparación especial para mapas
				if expectedMap, ok := expectedVal.(map[string]interface{}); ok {
					actualMap, ok := actualVal.(map[string]interface{})
					if !ok {
						t.Errorf("Variable %s: expected map, got %T", varName, actualVal)
						continue
					}

					for k, v := range expectedMap {
						if actualMap[k] != v {
							t.Errorf("Variable %s[%s]: expected %v, got %v",
								varName, k, v, actualMap[k])
						}
					}
				} else if actualVal != expectedVal {
					t.Errorf("Variable %s: expected %v (%T), got %v (%T)",
						varName, expectedVal, expectedVal, actualVal, actualVal)
				}
			}
		})
	}
}

func TestDestructuring_ErrorHandling(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		panics bool
	}{
		{
			name:   "array destructuring with non-array",
			input:  "let [a, b] = \"not an array\"",
			panics: true,
		},
		{
			name:   "object destructuring with non-object",
			input:  "let {name} = \"not an object\"",
			panics: true,
		},
		{
			name:   "array destructuring with number",
			input:  "let [a] = 42",
			panics: true,
		},
		{
			name:   "object destructuring with array",
			input:  "let {prop} = [1, 2, 3]",
			panics: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tt.panics && r == nil {
					t.Errorf("Expected panic but didn't get one")
				} else if !tt.panics && r != nil {
					t.Errorf("Unexpected panic: %v", r)
				}
			}()

			parser := NewParser(tt.input)
			program := parser.ParseProgram()
			env := NewEnvironment()

			for _, stmt := range program.Statements {
				stmt.Eval(env)
			}
		})
	}
}
