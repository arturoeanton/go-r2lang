package r2core

import (
	"testing"
)

func TestIfStatement_Basic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			"simple if true",
			`if (true) { return 42; }`,
			float64(42),
		},
		{
			"simple if false",
			`if (false) { return 42; } return 0;`,
			float64(0),
		},
		{
			"if with else",
			`if (false) { return 42; } else { return 99; }`,
			float64(99),
		},
		{
			"if condition with expression",
			`if (5 > 3) { return "greater"; } else { return "less"; }`,
			"greater",
		},
		{
			"nested if statements",
			`if (true) { if (false) { return 1; } else { return 2; } }`,
			float64(2),
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

func TestIfStatement_ElseIf(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			"else if basic",
			`if (false) { return 1; } else if (true) { return 2; }`,
			float64(2),
		},
		{
			"else if chain",
			`if (false) { return 1; } else if (false) { return 2; } else if (true) { return 3; }`,
			float64(3),
		},
		{
			"else if with final else",
			`if (false) { return 1; } else if (false) { return 2; } else { return 3; }`,
			float64(3),
		},
		{
			"multiple else if conditions",
			`let x = 10; if (x < 5) { return "small"; } else if (x < 15) { return "medium"; } else { return "large"; }`,
			"medium",
		},
		{
			"else if with expressions",
			`let age = 25; if (age < 18) { return "child"; } else if (age < 65) { return "adult"; } else { return "senior"; }`,
			"adult",
		},
		{
			"nested else if",
			`if (true) { if (false) { return 1; } else if (true) { return 2; } } else { return 3; }`,
			float64(2),
		},
		{
			"else if with complex expressions",
			`let a = 5; let b = 3; if (a < b) { return "less"; } else if (a == b) { return "equal"; } else if (a > b) { return "greater"; }`,
			"greater",
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

func TestIfStatement_Boolean_Expressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			"boolean AND",
			`let a = true; let b = true; if (a && b) { return "both true"; } else { return "not both"; }`,
			"both true",
		},
		{
			"boolean OR",
			`let a = false; let b = true; if (a || b) { return "at least one"; } else { return "both false"; }`,
			"at least one",
		},
		{
			"boolean comparison",
			`let a = 5; let b = 3; if (a > b) { return "a greater"; } else { return "b greater"; }`,
			"a greater",
		},
		{
			"complex boolean expression",
			`if ((5 > 3) && (10 < 20)) { return "logical"; } else { return "illogical"; }`,
			"logical",
		},
		{
			"boolean with else if",
			`let a = false; let b = true; if (a && b) { return 1; } else if (a || b) { return 2; } else { return 3; }`,
			float64(2),
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

func TestIfStatement_Variables(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			"if with variable condition",
			`let condition = true; if (condition) { return "variable true"; } else { return "variable false"; }`,
			"variable true",
		},
		{
			"else if with variables",
			`let x = 5; let y = 10; if (x > y) { return "x greater"; } else if (y > x) { return "y greater"; } else { return "equal"; }`,
			"y greater",
		},
		{
			"variable assignment in if block",
			`let result = "default"; if (true) { result = "changed"; } return result;`,
			"changed",
		},
		{
			"local variable in if block",
			`let global = "global"; if (true) { let local = "local"; return local; } return global;`,
			"local",
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
