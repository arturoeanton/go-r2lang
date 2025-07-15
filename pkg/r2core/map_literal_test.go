package r2core

import (
	"testing"
)

func TestMapLiteral_BasicKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]interface{}
	}{
		{
			"string keys",
			`{"name": "Juan", "age": 30}`,
			map[string]interface{}{"name": "Juan", "age": float64(30)},
		},
		{
			"identifier keys",
			`{name: "Juan", age: 30}`,
			map[string]interface{}{"name": "Juan", "age": float64(30)},
		},
		{
			"mixed keys",
			`{"name": "Juan", age: 30, "active": true}`,
			map[string]interface{}{"name": "Juan", "age": float64(30), "active": true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			program := parser.ParseProgram()

			if len(program.Statements) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ExprStatement)
			if !ok {
				t.Fatalf("expected ExprStatement, got %T", program.Statements[0])
			}

			env := NewEnvironment()
			// Registrar constantes básicas
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			result := stmt.Eval(env)

			resultMap, ok := result.(map[string]interface{})
			if !ok {
				t.Fatalf("expected map[string]interface{}, got %T", result)
			}

			if len(resultMap) != len(tt.expected) {
				t.Errorf("expected %d keys, got %d", len(tt.expected), len(resultMap))
			}

			for key, expectedValue := range tt.expected {
				actualValue, exists := resultMap[key]
				if !exists {
					t.Errorf("key %s not found in result", key)
					continue
				}

				if actualValue != expectedValue {
					t.Errorf("for key %s: expected %v, got %v", key, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestMapLiteral_ExpressionKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]interface{}
	}{
		{
			"string concatenation key",
			`{("key" + "1"): "value1"}`,
			map[string]interface{}{"key1": "value1"},
		},
		{
			"number concatenation key",
			`{("item" + 42): "value42"}`,
			map[string]interface{}{"item42": "value42"},
		},
		{
			"expression with variables",
			`let prefix = "user"; {(prefix + "_" + 123): "data"}`,
			map[string]interface{}{"user_123": "data"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			program := parser.ParseProgram()

			env := NewEnvironment()
			// Registrar constantes básicas
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			result := program.Eval(env)

			resultMap, ok := result.(map[string]interface{})
			if !ok {
				t.Fatalf("expected map[string]interface{}, got %T", result)
			}

			if len(resultMap) != len(tt.expected) {
				t.Errorf("expected %d keys, got %d", len(tt.expected), len(resultMap))
			}

			for key, expectedValue := range tt.expected {
				actualValue, exists := resultMap[key]
				if !exists {
					t.Errorf("key %s not found in result", key)
					continue
				}

				if actualValue != expectedValue {
					t.Errorf("for key %s: expected %v, got %v", key, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestMapLiteral_ComplexExpressionKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]interface{}
	}{
		{
			"arithmetic expression key",
			`{(10 + 20): "thirty"}`,
			map[string]interface{}{"30": "thirty"},
		},
		{
			"boolean expression key",
			`{(5 > 3): "true result"}`,
			map[string]interface{}{"true": "true result"},
		},
		{
			"multiple expression keys",
			`{(1 + 1): "two", ("hello" + "world"): "greeting", (true): "boolean"}`,
			map[string]interface{}{"2": "two", "helloworld": "greeting", "true": "boolean"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			program := parser.ParseProgram()

			env := NewEnvironment()
			// Registrar constantes básicas
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			result := program.Eval(env)

			resultMap, ok := result.(map[string]interface{})
			if !ok {
				t.Fatalf("expected map[string]interface{}, got %T", result)
			}

			if len(resultMap) != len(tt.expected) {
				t.Errorf("expected %d keys, got %d", len(tt.expected), len(resultMap))
			}

			for key, expectedValue := range tt.expected {
				actualValue, exists := resultMap[key]
				if !exists {
					t.Errorf("key %s not found in result", key)
					continue
				}

				if actualValue != expectedValue {
					t.Errorf("for key %s: expected %v, got %v", key, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestStringConcatenation_Optimization(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"string + string",
			`"Hello" + " " + "World"`,
			"Hello World",
		},
		{
			"string + number",
			`"Value: " + 42`,
			"Value: 42",
		},
		{
			"number + string",
			`123 + " items"`,
			"123 items",
		},
		{
			"string + boolean",
			`"Result: " + true`,
			"Result: true",
		},
		{
			"complex concatenation",
			`"User " + 1 + " has " + 5 + " points"`,
			"User 1 has 5 points",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			program := parser.ParseProgram()

			if len(program.Statements) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ExprStatement)
			if !ok {
				t.Fatalf("expected ExprStatement, got %T", program.Statements[0])
			}

			env := NewEnvironment()
			// Registrar constantes básicas
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			result := stmt.Eval(env)

			resultStr, ok := result.(string)
			if !ok {
				t.Fatalf("expected string, got %T: %v", result, result)
			}

			if resultStr != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, resultStr)
			}
		})
	}
}

func TestToString_Function(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"string", "hello", "hello"},
		{"integer", float64(42), "42"},
		{"float", float64(3.14), "3.14"},
		{"boolean true", true, "true"},
		{"boolean false", false, "false"},
		{"nil", nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toString(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestToStringOptimized_Function(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"string", "hello", "hello"},
		{"small integer", float64(42), "42"},
		{"large integer", float64(1000), "1000"},
		{"float", float64(3.14), "3.14"},
		{"boolean true", true, "true"},
		{"boolean false", false, "false"},
		{"nil", nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toStringOptimized(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestMapLiteral_StringConcatOptimization(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]interface{}
	}{
		{
			"optimized large string concat key",
			`{("very_long_prefix_name_for_key_" + "optimization_test_suffix"): "optimized"}`,
			map[string]interface{}{"very_long_prefix_name_for_key_optimization_test_suffix": "optimized"},
		},
		{
			"multiple optimized concat keys",
			`{("prefix1_" + "suffix1"): "val1", ("prefix2_" + "suffix2"): "val2"}`,
			map[string]interface{}{"prefix1_suffix1": "val1", "prefix2_suffix2": "val2"},
		},
		{
			"mixed small and large concat keys",
			`{("a" + "b"): "small", ("very_long_key_prefix_" + "very_long_key_suffix"): "large"}`,
			map[string]interface{}{"ab": "small", "very_long_key_prefix_very_long_key_suffix": "large"},
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

			resultMap, ok := result.(map[string]interface{})
			if !ok {
				t.Fatalf("expected map[string]interface{}, got %T", result)
			}

			if len(resultMap) != len(tt.expected) {
				t.Errorf("expected %d keys, got %d", len(tt.expected), len(resultMap))
			}

			for key, expectedValue := range tt.expected {
				actualValue, exists := resultMap[key]
				if !exists {
					t.Errorf("key %s not found in result", key)
					continue
				}

				if actualValue != expectedValue {
					t.Errorf("for key %s: expected %v, got %v", key, expectedValue, actualValue)
				}
			}
		})
	}
}
