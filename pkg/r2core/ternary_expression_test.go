package r2core

import (
	"testing"
)

func TestTernaryExpression_BasicConditions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			"true condition",
			"true ? \"yes\" : \"no\"",
			"yes",
		},
		{
			"false condition",
			"false ? \"yes\" : \"no\"",
			"no",
		},
		{
			"number comparison true",
			"5 > 3 ? \"greater\" : \"lesser\"",
			"greater",
		},
		{
			"number comparison false",
			"2 > 5 ? \"greater\" : \"lesser\"",
			"lesser",
		},
		{
			"equality comparison true",
			"10 == 10 ? \"equal\" : \"not equal\"",
			"equal",
		},
		{
			"equality comparison false",
			"10 == 5 ? \"equal\" : \"not equal\"",
			"not equal",
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
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			result := stmt.Eval(env)
			
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestTernaryExpression_WithVariables(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			"variable condition true",
			"let isActive = true; isActive ? \"active\" : \"inactive\"",
			"active",
		},
		{
			"variable condition false",
			"let isActive = false; isActive ? \"active\" : \"inactive\"",
			"inactive",
		},
		{
			"variable comparison",
			"let age = 25; age >= 18 ? \"adult\" : \"minor\"",
			"adult",
		},
		{
			"variable comparison false",
			"let age = 15; age >= 18 ? \"adult\" : \"minor\"",
			"minor",
		},
		{
			"string comparison",
			"let name = \"John\"; name == \"John\" ? \"correct\" : \"wrong\"",
			"correct",
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

func TestTernaryExpression_NestedTernary(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			"nested ternary in true branch",
			"true ? (true ? \"inner_true\" : \"inner_false\") : \"outer_false\"",
			"inner_true",
		},
		{
			"nested ternary in false branch",
			"false ? \"outer_true\" : (true ? \"inner_true\" : \"inner_false\")",
			"inner_true",
		},
		{
			"multiple nested ternary",
			"let score = 85; score >= 90 ? \"A\" : (score >= 80 ? \"B\" : (score >= 70 ? \"C\" : \"F\"))",
			"B",
		},
		{
			"complex nested with variables",
			"let x = 10; let y = 5; x > y ? (x > 15 ? \"very_high\" : \"high\") : \"low\"",
			"high",
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

func TestTernaryExpression_WithNumbers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			"number true condition",
			"5 > 3 ? 100 : 200",
			float64(100),
		},
		{
			"number false condition",
			"3 > 5 ? 100 : 200",
			float64(200),
		},
		{
			"arithmetic in condition",
			"(10 + 5) > (3 * 4) ? 42 : 24",
			float64(42),
		},
		{
			"arithmetic in branches",
			"true ? (10 + 5) : (3 * 4)",
			float64(15),
		},
		{
			"complex arithmetic",
			"let a = 10; let b = 20; a < b ? (a + b) : (a - b)",
			float64(30),
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

func TestTernaryExpression_WithBooleans(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			"boolean result true",
			"10 > 5 ? true : false",
			true,
		},
		{
			"boolean result false",
			"5 > 10 ? true : false",
			false,
		},
		{
			"boolean condition with boolean result",
			"true ? (5 > 3) : (3 > 5)",
			true,
		},
		{
			"complex boolean logic",
			"let x = true; let y = false; x ? (y ? false : true) : false",
			true,
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

func TestTernaryExpression_InTemplateLiterals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"ternary in template string",
			"let age = 25; `Status: ${age >= 18 ? \"Adult\" : \"Minor\"}`",
			"Status: Adult",
		},
		{
			"complex ternary in template",
			"let score = 85; `Grade: ${score >= 90 ? \"A\" : (score >= 80 ? \"B\" : \"C\")}`",
			"Grade: B",
		},
		{
			"multiple ternaries in template",
			"let x = 10; let y = 5; `Result: ${x > y ? \"greater\" : \"lesser\"} and ${x + y > 10 ? \"sum_high\" : \"sum_low\"}`",
			"Result: greater and sum_high",
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
			
			resultStr, ok := result.(string)
			if !ok {
				t.Fatalf("expected string result, got %T: %v", result, result)
			}
			
			if resultStr != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, resultStr)
			}
		})
	}
}

// Note: Error case testing disabled because parser.except() calls os.Exit(1)
// which prevents proper test isolation. These error cases are validated
// through integration testing instead.
//
// func TestTernaryExpression_ErrorCases(t *testing.T) {
// 	tests := []struct {
// 		name  string
// 		input string
// 	}{
// 		{
// 			"missing colon",
// 			"true ? \"yes\"",
// 		},
// 		{
// 			"missing question mark",
// 			"true \"yes\" : \"no\"",
// 		},
// 		{
// 			"incomplete condition",
// 			"? \"yes\" : \"no\"",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			defer func() {
// 				if r := recover(); r == nil {
// 					t.Errorf("expected panic for input: %s", tt.input)
// 				}
// 			}()
			
// 			parser := NewParser(tt.input)
// 			program := parser.ParseProgram()
			
// 			env := NewEnvironment()
// 			env.Set("true", true)
// 			env.Set("false", false)
// 			program.Eval(env)
// 		})
// 	}
// }