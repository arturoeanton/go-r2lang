package r2core

import (
	"testing"
)

func TestMatchExpression_Basic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			name:     "literal pattern match",
			input:    `let x = 2; match x { case 1 => "one" case 2 => "two" case _ => "other" }`,
			expected: "two",
		},
		{
			name:     "wildcard fallback",
			input:    `let x = 99; match x { case 1 => "one" case _ => "other" }`,
			expected: "other",
		},
		{
			name:     "variable pattern binds value",
			input:    `let x = 42; match x { case y => y + 1 }`,
			expected: float64(43),
		},
		{
			name:     "array pattern destructures",
			input:    `let arr = [1, 2]; match arr { case [a, b] => a + b }`,
			expected: float64(3),
		},
		{
			name:     "object pattern destructures",
			input:    `let obj = {name: "Ana", age: 30}; match obj { case {name} => name }`,
			expected: "Ana",
		},
		{
			name:     "guard clause",
			input:    `let x = 10; match x { case y if y > 5 => "big" case _ => "small" }`,
			expected: "big",
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
				t.Errorf("expected %v (%T), got %v (%T)", tt.expected, tt.expected, result, result)
			}
		})
	}
}

func TestMatchExpression_NoMatchPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic when no case matches")
		}
	}()

	input := `let x = 5; match x { case 1 => "one" case 2 => "two" }`
	parser := NewParser(input)
	program := parser.ParseProgram()
	env := NewEnvironment()

	for _, stmt := range program.Statements {
		stmt.Eval(env)
	}
}

// TestObjectPattern_VariableMap covers matching against a
// map[string]*Variable, the representation used for `import "x" as alias`
// (see import_statement.go / AccessExpression.evalVariableMapAccess).
// Before the fix, ObjectPattern.MatchValue only recognized
// map[string]interface{} and silently failed to match (returning
// false, nil) against an imported module alias, so a `match m { case
// {add} => ... }` would always fall through to the wildcard case instead
// of matching the module's exported symbols.
func TestObjectPattern_VariableMap(t *testing.T) {
	env := NewEnvironment()
	moduleObj := map[string]*Variable{
		"add": {Value: float64(7)},
	}

	pattern := &ObjectPattern{
		Fields: map[string]Pattern{
			"add": &VariablePattern{Name: "add"},
		},
	}

	matched, bindings := pattern.MatchValue(moduleObj, env)
	if !matched {
		t.Fatalf("expected ObjectPattern to match a map[string]*Variable value")
	}
	if bindings["add"] != float64(7) {
		t.Errorf("expected bound value 7, got %v", bindings["add"])
	}

	// Field that doesn't exist in the module should fail to match.
	missingPattern := &ObjectPattern{
		Fields: map[string]Pattern{
			"doesNotExist": &VariablePattern{Name: "doesNotExist"},
		},
	}
	matched, _ = missingPattern.MatchValue(moduleObj, env)
	if matched {
		t.Errorf("expected no match for a field absent from the module")
	}
}
