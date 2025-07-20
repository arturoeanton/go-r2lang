package r2core

import (
	"reflect"
	"testing"
)

func TestSpreadExpression_ArraySpread(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []interface{}
	}{
		{
			name:     "basic array spread",
			input:    "let arr1 = [1, 2]; let arr2 = [...arr1, 3, 4]; arr2",
			expected: []interface{}{float64(1), float64(2), float64(3), float64(4)},
		},
		{
			name:     "multiple spreads",
			input:    "let arr1 = [1, 2]; let arr2 = [5, 6]; let arr3 = [...arr1, 3, 4, ...arr2]; arr3",
			expected: []interface{}{float64(1), float64(2), float64(3), float64(4), float64(5), float64(6)},
		},
		{
			name:     "spread at beginning",
			input:    "let arr1 = [2, 3]; let arr2 = [...arr1, 1]; arr2",
			expected: []interface{}{float64(2), float64(3), float64(1)},
		},
		{
			name:     "spread in middle",
			input:    "let arr1 = [2, 3]; let arr2 = [1, ...arr1, 4]; arr2",
			expected: []interface{}{float64(1), float64(2), float64(3), float64(4)},
		},
		{
			name:     "empty array spread",
			input:    "let arr1 = []; let arr2 = [1, ...arr1, 2]; arr2",
			expected: []interface{}{float64(1), float64(2)},
		},
		{
			name:     "nested array spread",
			input:    "let arr1 = [[1, 2], 3]; let arr2 = [...arr1, 4]; arr2",
			expected: []interface{}{[]interface{}{float64(1), float64(2)}, float64(3), float64(4)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			program := parser.ParseProgram()
			env := NewEnvironment()

			var result interface{}
			for _, stmt := range program.Statements {
				result = stmt.Eval(env)
			}

			actualArray, ok := result.([]interface{})
			if !ok {
				t.Fatalf("Expected array result, got %T", result)
			}

			if !reflect.DeepEqual(actualArray, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, actualArray)
			}
		})
	}
}

func TestSpreadExpression_ObjectSpread(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]interface{}
	}{
		{
			name:  "basic object spread",
			input: "let obj1 = {a: 1, b: 2}; let obj2 = {...obj1, c: 3}; obj2",
			expected: map[string]interface{}{
				"a": float64(1),
				"b": float64(2),
				"c": float64(3),
			},
		},
		{
			name:  "override properties",
			input: "let obj1 = {a: 1, b: 2}; let obj2 = {...obj1, b: 3, c: 4}; obj2",
			expected: map[string]interface{}{
				"a": float64(1),
				"b": float64(3),
				"c": float64(4),
			},
		},
		{
			name:  "multiple spreads",
			input: "let obj1 = {a: 1}; let obj2 = {b: 2}; let obj3 = {...obj1, ...obj2, c: 3}; obj3",
			expected: map[string]interface{}{
				"a": float64(1),
				"b": float64(2),
				"c": float64(3),
			},
		},
		{
			name:  "empty object spread",
			input: "let obj1 = {}; let obj2 = {...obj1, a: 1}; obj2",
			expected: map[string]interface{}{
				"a": float64(1),
			},
		},
		{
			name:  "spread with complex values",
			input: "let obj1 = {name: \"John\", data: [1, 2, 3]}; let obj2 = {...obj1, age: 30}; obj2",
			expected: map[string]interface{}{
				"name": "John",
				"data": []interface{}{float64(1), float64(2), float64(3)},
				"age":  float64(30),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			program := parser.ParseProgram()
			env := NewEnvironment()

			var result interface{}
			for _, stmt := range program.Statements {
				result = stmt.Eval(env)
			}

			actualObj, ok := result.(map[string]interface{})
			if !ok {
				t.Fatalf("Expected object result, got %T", result)
			}

			for key, expectedVal := range tt.expected {
				actualVal, exists := actualObj[key]
				if !exists {
					t.Errorf("Missing key %s in result", key)
					continue
				}

				// Handle array comparison specially
				if expectedArr, ok := expectedVal.([]interface{}); ok {
					actualArr, ok := actualVal.([]interface{})
					if !ok {
						t.Errorf("Key %s: expected array, got %T", key, actualVal)
						continue
					}
					if !reflect.DeepEqual(actualArr, expectedArr) {
						t.Errorf("Key %s: expected %v, got %v", key, expectedArr, actualArr)
					}
				} else if actualVal != expectedVal {
					t.Errorf("Key %s: expected %v, got %v", key, expectedVal, actualVal)
				}
			}

			// Check for extra keys
			for key := range actualObj {
				if _, exists := tt.expected[key]; !exists {
					t.Errorf("Unexpected key %s in result", key)
				}
			}
		})
	}
}

func TestSpreadExpression_FunctionCall(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			name:     "spread array as function arguments",
			input:    "func sum(a, b, c) { return a + b + c; } let args = [1, 2, 3]; sum(...args)",
			expected: float64(6),
		},
		{
			name:     "spread with additional arguments",
			input:    "func sum(a, b, c, d) { return a + b + c + d; } let args = [1, 2]; sum(...args, 3, 4)",
			expected: float64(10),
		},
		{
			name:     "mixed spread and regular arguments",
			input:    "func sum(a, b, c, d) { return a + b + c + d; } let args = [2, 3]; sum(1, ...args, 4)",
			expected: float64(10),
		},
		{
			name:     "empty array spread",
			input:    "func sum() { return 42; } let args = []; sum(...args)",
			expected: float64(42),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			program := parser.ParseProgram()
			env := NewEnvironment()

			var result interface{}
			for _, stmt := range program.Statements {
				result = stmt.Eval(env)
			}

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSpreadValue_HelperFunctions(t *testing.T) {
	t.Run("IsSpreadValue", func(t *testing.T) {
		// Test with regular value
		regularVal := "not a spread"
		if sv, isSpread := IsSpreadValue(regularVal); isSpread {
			t.Errorf("Regular value incorrectly identified as spread: %v", sv)
		}

		// Test with spread value
		spreadVal := &SpreadValue{Value: []interface{}{1, 2, 3}}
		if sv, isSpread := IsSpreadValue(spreadVal); !isSpread {
			t.Errorf("Spread value not identified correctly")
		} else if sv.Value == nil {
			t.Errorf("Spread value content is nil")
		}
	})

	t.Run("ExpandSpreadInArray", func(t *testing.T) {
		// Test with mixed spread and regular elements
		arr := []interface{}{
			float64(1),
			&SpreadValue{Value: []interface{}{float64(2), float64(3)}},
			float64(4),
		}

		result := ExpandSpreadInArray(arr)
		expected := []interface{}{float64(1), float64(2), float64(3), float64(4)}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("ExpandSpreadInFunctionCall", func(t *testing.T) {
		// Test function call argument expansion
		args := []interface{}{
			float64(1),
			&SpreadValue{Value: []interface{}{float64(2), float64(3)}},
			float64(4),
		}

		result := ExpandSpreadInFunctionCall(args)
		expected := []interface{}{float64(1), float64(2), float64(3), float64(4)}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}

func TestSpreadExpression_ErrorCases(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		panics bool
	}{
		{
			name:   "spread non-array in function call",
			input:  "func test(a, b) { return a + b; } test(...\"not array\")",
			panics: false, // Should handle gracefully by treating as single argument
		},
		{
			name:   "spread non-object in object literal",
			input:  "let myobj = {...\"not object\", a: 1}; myobj",
			panics: false, // Should handle by treating as regular property
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
