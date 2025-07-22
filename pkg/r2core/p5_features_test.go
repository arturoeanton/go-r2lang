package r2core

import (
	"testing"
)

// TestP5_StringInterpolationImproved tests enhanced string interpolation with formatting
func TestP5_StringInterpolationImproved(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "currency_formatting",
			code:     "let price = 123.456; return `Price: ${price:$,.2f}`;",
			expected: "Price: $123.46",
		},
		{
			name:     "percentage_formatting",
			code:     "let rate = 0.8534; return `Rate: ${rate:.1%}`;",
			expected: "Rate: 85.3%",
		},
		{
			name:     "float_formatting",
			code:     "let num = 3.14159; return `Number: ${num:.2f}`;",
			expected: "Number: 3.14",
		},
		{
			name:     "comma_formatting",
			code:     "let big = 1234567; return `Big: ${big:,}`;",
			expected: "Big: 1,234,567",
		},
		{
			name:     "string_upper",
			code:     "let text = \"hello\"; return `Text: ${text:upper}`;",
			expected: "Text: HELLO",
		},
		{
			name:     "string_lower",
			code:     "let text = \"WORLD\"; return `Text: ${text:lower}`;",
			expected: "Text: world",
		},
		{
			name:     "no_format_backward_compatibility",
			code:     "let name = \"R2Lang\"; return `Hello ${name}!`;",
			expected: "Hello R2Lang!",
		},
		{
			name:     "mixed_formatting",
			code:     "let price = 99.99; let name = \"item\"; return `${name:upper}: ${price:$,.2f}`;",
			expected: "ITEM: $99.99",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			parser := NewParser(tt.code)
			ast := parser.ParseProgram()
			result := ast.Eval(env)

			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestP5_SmartAutoConversion tests intelligent type conversion
func TestP5_SmartAutoConversion(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected interface{}
	}{
		{
			name:     "string_number_addition",
			code:     "return \"10\" + \"20\";",
			expected: float64(30),
		},
		{
			name:     "string_number_mixed",
			code:     "return \"10\" + 5;",
			expected: float64(15),
		},
		{
			name:     "boolean_number_addition",
			code:     "return true + 1;",
			expected: float64(2),
		},
		{
			name:     "boolean_false_addition",
			code:     "return false + 10;",
			expected: float64(10),
		},
		{
			name:     "comma_number_parsing",
			code:     "return \"1,000\" + \"2,000\";",
			expected: float64(3000),
		},
		{
			name:     "currency_parsing",
			code:     "return \"$100\" + \"$50\";",
			expected: float64(150),
		},
		{
			name:     "percentage_parsing",
			code:     "return \"50%\" + \"25%\";",
			expected: float64(0.75), // 0.5 + 0.25
		},
		{
			name:     "boolean_string_parsing",
			code:     "return \"true\" + \"false\";",
			expected: float64(1), // 1 + 0
		},
		{
			name:     "array_concatenation_preserved",
			code:     "return [1, 2] + [3, 4];",
			expected: []interface{}{float64(1), float64(2), float64(3), float64(4)},
		},
		{
			name:     "mixed_string_fallback",
			code:     "return \"hello\" + \" world\";",
			expected: "hello world",
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
			case string:
				if resultString, ok := result.(string); !ok || resultString != expected {
					t.Errorf("Expected %q, got %q (type: %T)", expected, result, result)
				}
			case []interface{}:
				if resultSlice, ok := result.([]interface{}); !ok {
					t.Errorf("Expected slice, got %T", result)
				} else {
					if len(resultSlice) != len(expected) {
						t.Errorf("Expected slice length %d, got %d", len(expected), len(resultSlice))
					} else {
						for i, exp := range expected {
							if resultSlice[i] != exp {
								t.Errorf("Expected slice[%d] = %v, got %v", i, exp, resultSlice[i])
							}
						}
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

// TestP5_SmartParseFloat tests the smart float parsing function
func TestP5_SmartParseFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		hasError bool
	}{
		{"123", 123, false},
		{"123.45", 123.45, false},
		{"1,000", 1000, false},
		{"1,234.56", 1234.56, false},
		{"$100", 100, false},
		{"$1,234.56", 1234.56, false},
		{"50%", 0.5, false},
		{"75.5%", 0.755, false},
		{"true", 1, false},
		{"false", 0, false},
		{"yes", 1, false},
		{"no", 0, false},
		{"on", 1, false},
		{"off", 0, false},
		{"  123  ", 123, false},
		{"", 0, false},
		{"abc", 0, true},
		{"12abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := smartParseFloat(tt.input)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error for input %q, but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input %q: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("For input %q, expected %v, got %v", tt.input, tt.expected, result)
				}
			}
		})
	}
}

// TestP5_FormatValue tests the formatting functions
func TestP5_FormatValue(t *testing.T) {
	tests := []struct {
		value    interface{}
		format   string
		expected string
	}{
		{123.456, "$,.2f", "$123.46"},
		{0.8534, ".1%", "85.3%"},
		{3.14159, ".2f", "3.14"},
		{1234567, ",", "1,234,567"},
		{"hello", "upper", "HELLO"},
		{"WORLD", "lower", "world"},
		{"test", "title", "Test"},
		{"  spaced  ", "trim", "spaced"},
		{42, "d", "42"},
		{3.14, "g", "3.14"},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			result := formatValue(tt.value, tt.format)
			if result != tt.expected {
				t.Errorf("formatValue(%v, %q) = %q, expected %q", tt.value, tt.format, result, tt.expected)
			}
		})
	}
}

// TestP5_BackwardCompatibility ensures P5 features don't break existing functionality
func TestP5_BackwardCompatibility(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected interface{}
	}{
		{
			name:     "existing_template_strings",
			code:     "let name = \"World\"; return `Hello ${name}!`;",
			expected: "Hello World!",
		},
		{
			name:     "existing_arithmetic",
			code:     "return 10 + 20;",
			expected: float64(30),
		},
		{
			name:     "existing_string_concat",
			code:     "return \"Hello\" + \" \" + \"World\";",
			expected: "Hello World",
		},
		{
			name:     "existing_array_concat",
			code:     "return [1, 2] + [3, 4];",
			expected: []interface{}{float64(1), float64(2), float64(3), float64(4)},
		},
		{
			name:     "existing_boolean_ops",
			code:     "return true && false;",
			expected: false,
		},
		{
			name:     "existing_comparison",
			code:     "return 5 > 3;",
			expected: true,
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
					t.Errorf("Expected %v, got %v", expected, result)
				}
			case string:
				if resultString, ok := result.(string); !ok || resultString != expected {
					t.Errorf("Expected %q, got %q", expected, result)
				}
			case bool:
				if resultBool, ok := result.(bool); !ok || resultBool != expected {
					t.Errorf("Expected %v, got %v", expected, result)
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
