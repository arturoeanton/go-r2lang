package r2core

import (
	"strings"
	"testing"
)

// TestOptionalChaining tests the optional chaining operator ?. (P3)
func TestOptionalChaining(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			name:     "basic optional chaining with existing property",
			input:    `let user = {name: "John", age: 30}; user?.name`,
			expected: "John",
		},
		{
			name:     "optional chaining with nil object",
			input:    `let user = nil; user?.name`,
			expected: nil,
		},
		{
			name:     "nested optional chaining",
			input:    `let user = {profile: {name: "Alice"}}; user?.profile?.name`,
			expected: "Alice",
		},
		{
			name:     "nested optional chaining with missing intermediate",
			input:    `let user = {name: "Bob"}; user?.profile?.name`,
			expected: nil,
		},
		{
			name:     "optional chaining with array-like access",
			input:    `let data = {items: [1, 2, 3]}; data?.items?.length`,
			expected: float64(3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			result := executeScript(tt.input, env)
			if result != tt.expected {
				t.Errorf("Expected %v (%T), got %v (%T)", tt.expected, tt.expected, result, result)
			}
		})
	}
}

// TestNullCoalescing tests the null coalescing operator ?? (P3)
func TestNullCoalescing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			name:     "null coalescing with nil left side",
			input:    `let x = nil; x ?? "default"`,
			expected: "default",
		},
		{
			name:     "null coalescing with non-nil left side",
			input:    `let x = "value"; x ?? "default"`,
			expected: "value",
		},
		{
			name:     "null coalescing with zero value",
			input:    `let x = 0; x ?? 10`,
			expected: float64(0), // 0 is not nil, should return 0
		},
		{
			name:     "null coalescing with empty string",
			input:    `let x = ""; x ?? "default"`,
			expected: "", // empty string is not nil
		},
		{
			name:     "null coalescing chain",
			input:    `let a = nil; let b = nil; let c = "result"; a ?? b ?? c`,
			expected: "result",
		},
		{
			name:     "null coalescing with expressions",
			input:    `let user = {}; user?.name ?? "Anonymous"`,
			expected: "Anonymous",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			result := executeScript(tt.input, env)
			if result != tt.expected {
				t.Errorf("Expected %v (%T), got %v (%T)", tt.expected, tt.expected, result, result)
			}
		})
	}
}

// TestPipelineOperator tests the pipeline operator |> (P4)
func TestPipelineOperator(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			name:     "simple pipeline with built-in function",
			input:    `func double(x) { return x * 2; } 5 |> double`,
			expected: float64(10),
		},
		{
			name:     "pipeline chain",
			input:    `func double(x) { return x * 2; } func add10(x) { return x + 10; } 5 |> double |> add10`,
			expected: float64(20),
		},
		{
			name:     "pipeline with arrow function",
			input:    `5 |> (x => x * 3)`,
			expected: float64(15),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			result := executeScript(tt.input, env)
			if result != tt.expected {
				t.Errorf("Expected %v (%T), got %v (%T)", tt.expected, tt.expected, result, result)
			}
		})
	}
}

// TestMatchExpression tests pattern matching (P3)
func TestMatchExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			name:     "literal pattern matching",
			input:    `match 5 { case 3 => "three" case 5 => "five" case _ => "other" }`,
			expected: "five",
		},
		{
			name:     "variable pattern binding",
			input:    `match 42 { case x => x * 2 }`,
			expected: float64(84),
		},
		{
			name:     "wildcard pattern",
			input:    `match "hello" { case "world" => "not found" case _ => "found" }`,
			expected: "found",
		},
		{
			name:     "array pattern matching",
			input:    `match [1, 2, 3] { case [a, b, c] => a + b + c case _ => 0 }`,
			expected: float64(6),
		},
		{
			name:     "object pattern matching",
			input:    `let user = {name: "John", age: 30}; match user { case {name, age} => name + " is " + age + " years old" case _ => "unknown" }`,
			expected: "John is 30 years old",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			result := executeScript(tt.input, env)
			if result != tt.expected {
				t.Errorf("Expected %v (%T), got %v (%T)", tt.expected, tt.expected, result, result)
			}
		})
	}
}

// TestArrayComprehension tests array comprehensions (P4)
func TestArrayComprehension(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []interface{}
	}{
		{
			name:     "simple array comprehension",
			input:    `[x * 2 for x in [1, 2, 3, 4, 5]]`,
			expected: []interface{}{float64(2), float64(4), float64(6), float64(8), float64(10)},
		},
		{
			name:     "array comprehension with condition",
			input:    `[x for x in [1, 2, 3, 4, 5] if x % 2 == 0]`,
			expected: []interface{}{float64(2), float64(4)},
		},
		{
			name:     "array comprehension with expression",
			input:    `[x * x for x in [1, 2, 3] if x > 1]`,
			expected: []interface{}{float64(4), float64(9)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			result := executeScript(tt.input, env)

			// Convert result to []interface{} for comparison
			if arr, ok := result.([]interface{}); ok {
				if len(arr) != len(tt.expected) {
					t.Errorf("Expected array length %d, got %d", len(tt.expected), len(arr))
					return
				}
				for i, expected := range tt.expected {
					if arr[i] != expected {
						t.Errorf("At index %d: expected %v (%T), got %v (%T)", i, expected, expected, arr[i], arr[i])
					}
				}
			} else {
				t.Errorf("Expected array result, got %v (%T)", result, result)
			}
		})
	}
}

// TestObjectComprehension tests object comprehensions (P4)
func TestObjectComprehension(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]interface{}
	}{
		{
			name:     "simple object comprehension",
			input:    `{x: x * x for x in [1, 2, 3]}`,
			expected: map[string]interface{}{"1": float64(1), "2": float64(4), "3": float64(9)},
		},
		{
			name:     "object comprehension with condition",
			input:    `{x: x * 2 for x in [1, 2, 3, 4] if x % 2 == 0}`,
			expected: map[string]interface{}{"2": float64(4), "4": float64(8)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			result := executeScript(tt.input, env)

			// Convert result to map[string]interface{} for comparison
			if obj, ok := result.(map[string]interface{}); ok {
				if len(obj) != len(tt.expected) {
					t.Errorf("Expected object with %d keys, got %d keys", len(tt.expected), len(obj))
					return
				}
				for key, expected := range tt.expected {
					if actual, exists := obj[key]; exists {
						if actual != expected {
							t.Errorf("For key %s: expected %v (%T), got %v (%T)", key, expected, expected, actual, actual)
						}
					} else {
						t.Errorf("Expected key %s not found in result", key)
					}
				}
			} else {
				t.Errorf("Expected object result, got %v (%T)", result, result)
			}
		})
	}
}

// TestCombinedP3P4Features tests combinations of P3 and P4 features
func TestCombinedP3P4Features(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			name:     "optional chaining with null coalescing",
			input:    `let user = {}; user?.profile?.name ?? "Anonymous"`,
			expected: "Anonymous",
		},
		{
			name:     "pipeline with simple function",
			input:    `func double(x) { return x * 2; } 5 |> double`,
			expected: float64(10),
		},
		{
			name:     "match with array comprehension",
			input:    `let data = [1, 2, 3]; match 3 { case 3 => [x * x for x in data] case _ => [] }`,
			expected: []interface{}{float64(1), float64(4), float64(9)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			result := executeScript(tt.input, env)

			switch expected := tt.expected.(type) {
			case []interface{}:
				if arr, ok := result.([]interface{}); ok {
					if len(arr) != len(expected) {
						t.Errorf("Expected array length %d, got %d", len(expected), len(arr))
						return
					}
					for i, exp := range expected {
						if arr[i] != exp {
							t.Errorf("At index %d: expected %v (%T), got %v (%T)", i, exp, exp, arr[i], arr[i])
						}
					}
				} else {
					t.Errorf("Expected array result, got %v (%T)", result, result)
				}
			default:
				if result != expected {
					t.Errorf("Expected %v (%T), got %v (%T)", expected, expected, result, result)
				}
			}
		})
	}
}

// executeScript is a helper function to execute R2Lang scripts for testing
func executeScript(script string, env *Environment) interface{} {
	parser := NewParser(script)
	program := parser.ParseProgram()

	// Execute statements and return the last result
	var result interface{}
	for _, stmt := range program.Statements {
		result = stmt.Eval(env)
	}

	return result
}

// TestErrorHandling tests error cases for P3 and P4 features
func TestP3P4ErrorHandling(t *testing.T) {
	errorTests := []struct {
		name          string
		input         string
		shouldPanic   bool
		expectedError string
	}{
		{
			name:          "match without matching case",
			input:         `match 5 { case 1 => "one" case 2 => "two" }`,
			shouldPanic:   true,
			expectedError: "No matching case found",
		},
		{
			name:        "pipeline with invalid function",
			input:       `5 |> notExist`,
			shouldPanic: true,
		},
		{
			name:        "comprehension with invalid iterator",
			input:       `[x for x in 42]`,
			shouldPanic: true,
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Expected panic, but none occurred")
					} else if tt.expectedError != "" {
						errorStr := r.(string)
						if !strings.Contains(errorStr, tt.expectedError) {
							t.Errorf("Expected error containing '%s', got '%s'", tt.expectedError, errorStr)
						}
					}
				}()
			}

			env := NewEnvironment()
			executeScript(tt.input, env)

			if tt.shouldPanic {
				t.Errorf("Expected panic but execution completed normally")
			}
		})
	}
}
