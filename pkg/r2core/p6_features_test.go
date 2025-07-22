package r2core

import (
	"fmt"
	"testing"
)

// TestP6_PlaceholderHandling tests that '_' is recognized as a placeholder
func TestP6_PlaceholderHandling(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "placeholder_identifier",
			code:     "let ph = _; return ph;",
			expected: "*r2core.Placeholder",
		},
		{
			name:     "multiple_placeholders",
			code:     "return [_, _, _];",
			expected: "placeholders_array",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			// Set up basic built-in functions for testing
			env.Set("typeOf", BuiltinFunction(func(args ...interface{}) interface{} {
				if len(args) < 1 {
					return "nil"
				}
				return fmt.Sprintf("%T", args[0])
			}))
			parser := NewParser(tt.code)
			ast := parser.ParseProgram()
			result := ast.Eval(env)

			if tt.expected == "placeholders_array" {
				if arr, ok := result.([]interface{}); ok {
					if len(arr) != 3 {
						t.Errorf("Expected array of 3 placeholders, got length %d", len(arr))
					}
					for i, item := range arr {
						if _, ok := item.(*Placeholder); !ok {
							t.Errorf("Expected placeholder at index %d, got %T", i, item)
						}
					}
				} else {
					t.Errorf("Expected array, got %T", result)
				}
			} else if tt.expected == "*r2core.Placeholder" {
				if _, ok := result.(*Placeholder); !ok {
					t.Errorf("Expected placeholder, got %T", result)
				}
			} else {
				if str, ok := result.(string); !ok || str != tt.expected {
					t.Errorf("Expected %q, got %q (type: %T)", tt.expected, result, result)
				}
			}
		})
	}
}

// TestP6_PartialApplication tests partial application with placeholders
func TestP6_PartialApplication(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected interface{}
	}{
		{
			name:     "simple_partial_application",
			code:     `func add(a, b) { return a + b; } let addFive = add(5, _); return addFive(10);`,
			expected: float64(15),
		},
		{
			name:     "multiple_placeholders",
			code:     `func multiply(a, b, c) { return a * b * c; } let partialMult = multiply(2, _, _); return partialMult(3, 4);`,
			expected: float64(24),
		},
		{
			name:     "mixed_arguments_and_placeholders",
			code:     `func subtract(a, b, c) { return a - b - c; } let partialSub = subtract(_, 2, _); return partialSub(10, 3);`,
			expected: float64(5), // 10 - 2 - 3
		},
		{
			name:     "nested_partial_application",
			code:     `func add3(a, b, c) { return a + b + c; } let addTen = add3(10, _, _); return addTen(5, 3);`,
			expected: float64(18), // 10 + 5 + 3
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			parser := NewParser(tt.code)
			ast := parser.ParseProgram()
			result := ast.Eval(env)

			switch expected := tt.expected.(type) {
			case float64:
				if resultFloat, ok := result.(float64); !ok || resultFloat != expected {
					t.Errorf("Expected %v, got %v (type: %T)", expected, result, result)
				}
			default:
				if result != expected {
					t.Errorf("Expected %v, got %v", expected, result)
				}
			}
		})
	}
}

// TestP6_ExplicitPartialFunction tests the explicit partial() function
func TestP6_ExplicitPartialFunction(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected interface{}
	}{
		{
			name:     "explicit_partial_with_partial_function",
			code:     `func multiply(a, b) { return a * b; } let doubler = partial(multiply, 2); return doubler(5);`,
			expected: float64(10),
		},
		{
			name:     "explicit_partial_with_multiple_args",
			code:     `func add3(a, b, c) { return a + b + c; } let addTenFive = partial(add3, 10, 5); return addTenFive(3);`,
			expected: float64(18),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			env.Set("partial", BuiltinFunction(PartialBuiltin))
			parser := NewParser(tt.code)
			ast := parser.ParseProgram()
			result := ast.Eval(env)

			switch expected := tt.expected.(type) {
			case float64:
				if resultFloat, ok := result.(float64); !ok || resultFloat != expected {
					t.Errorf("Expected %v, got %v (type: %T)", expected, result, result)
				}
			default:
				if result != expected {
					t.Errorf("Expected %v, got %v", expected, result)
				}
			}
		})
	}
}

// TestP6_Currying tests automatic currying functionality
func TestP6_Currying(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected interface{}
	}{
		{
			name:     "simple_currying",
			code:     `func add(a, b) { return a + b; } let curriedAdd = curry(add); return curriedAdd(5)(10);`,
			expected: float64(15),
		},
		{
			name:     "three_parameter_currying",
			code:     `func add3(a, b, c) { return a + b + c; } let curriedAdd3 = curry(add3); return curriedAdd3(1)(2)(3);`,
			expected: float64(6),
		},
		{
			name:     "currying_with_intermediate_functions",
			code:     `func multiply(a, b, c) { return a * b * c; } let curriedMult = curry(multiply); let multBy2 = curriedMult(2); let multBy2And3 = multBy2(3); return multBy2And3(4);`,
			expected: float64(24),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			env.Set("curry", BuiltinFunction(CurryFunction))
			parser := NewParser(tt.code)
			ast := parser.ParseProgram()
			result := ast.Eval(env)

			switch expected := tt.expected.(type) {
			case float64:
				if resultFloat, ok := result.(float64); !ok || resultFloat != expected {
					t.Errorf("Expected %v, got %v (type: %T)", expected, result, result)
				}
			default:
				if result != expected {
					t.Errorf("Expected %v, got %v", expected, result)
				}
			}
		})
	}
}

// TestP6_PartialFunctionTypes tests that partial functions return correct types
func TestP6_PartialFunctionTypes(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		checkFunc func(interface{}) bool
	}{
		{
			name: "partial_function_type",
			code: `func add(a, b) { return a + b; } return add(5, _);`,
			checkFunc: func(result interface{}) bool {
				_, ok := result.(*PartialFunction)
				return ok
			},
		},
		{
			name: "curried_function_type",
			code: `func add(a, b) { return a + b; } let curriedAdd = curry(add); return curriedAdd(5);`,
			checkFunc: func(result interface{}) bool {
				_, ok := result.(*CurriedFunction)
				return ok
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			env.Set("curry", BuiltinFunction(CurryFunction))
			parser := NewParser(tt.code)
			ast := parser.ParseProgram()
			result := ast.Eval(env)

			if !tt.checkFunc(result) {
				t.Errorf("Type check failed for test %s, got type %T", tt.name, result)
			}
		})
	}
}

// TestP6_BackwardCompatibility ensures P6 features don't break existing functionality
func TestP6_BackwardCompatibility(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected interface{}
	}{
		{
			name:     "normal_function_calls_unchanged",
			code:     `func add(a, b) { return a + b; } return add(5, 10);`,
			expected: float64(15),
		},
		{
			name:     "existing_identifiers_unchanged",
			code:     `let x = 42; return x;`,
			expected: float64(42),
		},
		{
			name:     "underscore_in_strings_unchanged",
			code:     `return "hello_world";`,
			expected: "hello_world",
		},
		{
			name:     "normal_arrays_unchanged",
			code:     `return [1, 2, 3];`,
			expected: []interface{}{float64(1), float64(2), float64(3)},
		},
		{
			name:     "existing_built_ins_unchanged",
			code:     `return len("hello");`,
			expected: float64(5),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			env.Set("len", BuiltinFunction(func(args ...interface{}) interface{} {
				if len(args) < 1 {
					panic("len needs 1 argument")
				}
				switch v := args[0].(type) {
				case string:
					return float64(len(v))
				case []interface{}:
					return float64(len(v))
				default:
					panic("len: Expected string or []interface{}")
				}
			}))
			parser := NewParser(tt.code)
			ast := parser.ParseProgram()
			result := ast.Eval(env)

			switch expected := tt.expected.(type) {
			case float64:
				if resultFloat, ok := result.(float64); !ok || resultFloat != expected {
					t.Errorf("Expected %v, got %v", expected, result)
				}
			case string:
				if resultString, ok := result.(string); !ok || resultString != expected {
					t.Errorf("Expected %q, got %q", expected, result)
				}
			case []interface{}:
				if resultSlice, ok := result.([]interface{}); !ok {
					t.Errorf("Expected slice, got %T", result)
				} else {
					if len(resultSlice) != len(expected) {
						t.Errorf("Expected slice length %d, got %d", len(expected), len(resultSlice))
					}
				}
			default:
				if result != expected {
					t.Errorf("Expected %v, got %v", expected, result)
				}
			}
		})
	}
}
