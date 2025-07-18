package r2core

import (
	"strings"
	"testing"
)

func TestDSLBasicFunctionality(t *testing.T) {
	// Test case 1: Simple command DSL
	env := NewEnvironment()

	// Create a simple DSL definition
	dslCode := `
	dsl TestDSL {
		token("HELLO", "hello")
		token("WORLD", "world")
		
		rule("greeting", ["HELLO", "WORLD"], "greet")
		
		func greet(h, w) {
			return h + " " + w
		}
	}
	`

	// Parse and execute the DSL definition
	parser := NewParser(dslCode)
	program := parser.ParseProgram()

	// Evaluate the program
	result := program.Eval(env)
	if result == nil {
		t.Fatal("Expected DSL evaluation to return a result")
	}

	// Get the DSL object from environment
	dslObj, exists := env.Get("TestDSL")
	if !exists {
		t.Fatal("DSL 'TestDSL' not found in environment")
	}

	// Test that DSL object has the expected structure
	if dslDef, ok := dslObj.(*DSLDefinition); ok {
		if dslDef.Name.Name != "TestDSL" {
			t.Fatalf("Expected DSL name 'TestDSL', got '%s'", dslDef.Name.Name)
		}
		if dslDef.Grammar == nil {
			t.Fatal("DSL should have a grammar")
		}
	} else {
		t.Fatalf("Expected DSL object to be *DSLDefinition, got %T", dslObj)
	}
}

func TestDSLParameterPassing(t *testing.T) {
	// Test case 2: Multiple parameters DSL (the bug we fixed)
	env := NewEnvironment()

	// Create a DSL that takes multiple parameters
	dslCode := `
	dsl CalcDSL {
		token("NUM", "[0-9]+")
		token("PLUS", "\\+")
		token("MINUS", "-")
		
		rule("operation", ["NUM", "PLUS", "NUM"], "add")
		rule("operation", ["NUM", "MINUS", "NUM"], "subtract")
		
		func add(n1, op, n2) {
			return n1 + "+" + n2
		}
		
		func subtract(n1, op, n2) {
			return n1 + "-" + n2
		}
	}
	`

	// Parse and execute the DSL definition
	parser := NewParser(dslCode)
	program := parser.ParseProgram()

	// Evaluate the program
	program.Eval(env)

	// Get the DSL object
	dslObj, exists := env.Get("CalcDSL")
	if !exists {
		t.Fatal("DSL 'CalcDSL' not found in environment")
	}

	// Test that DSL object is a DSLDefinition
	if dslDef, ok := dslObj.(*DSLDefinition); ok {
		if dslDef.Name.Name != "CalcDSL" {
			t.Fatalf("Expected DSL name 'CalcDSL', got '%s'", dslDef.Name.Name)
		}
		if dslDef.Grammar == nil {
			t.Fatal("DSL should have a grammar")
		}
		// Test that grammar has the expected rules
		if _, exists := dslDef.Grammar.Rules["operation"]; !exists {
			t.Fatal("DSL should have 'operation' rule")
		}
	} else {
		t.Fatalf("Expected DSL object to be *DSLDefinition, got %T", dslObj)
	}
}

func TestDSLTokenization(t *testing.T) {
	// Test case 3: Tokenization correctness
	grammar := NewDSLGrammar()

	// Add tokens
	err := grammar.AddToken("NUMBER", "[0-9]+")
	if err != nil {
		t.Fatalf("Failed to add NUMBER token: %v", err)
	}

	err = grammar.AddToken("PLUS", "\\+")
	if err != nil {
		t.Fatalf("Failed to add PLUS token: %v", err)
	}

	// Create parser and test tokenization
	parser := NewDSLParser(grammar)
	err = parser.Tokenize("123 + 456")
	if err != nil {
		t.Fatalf("Tokenization failed: %v", err)
	}

	// Check tokens
	if len(parser.Tokens) != 3 {
		t.Fatalf("Expected 3 tokens, got %d", len(parser.Tokens))
	}

	if parser.Tokens[0].Type != "NUMBER" || parser.Tokens[0].Value != "123" {
		t.Errorf("Expected first token to be NUMBER '123', got %s '%s'", parser.Tokens[0].Type, parser.Tokens[0].Value)
	}

	if parser.Tokens[1].Type != "PLUS" || parser.Tokens[1].Value != "+" {
		t.Errorf("Expected second token to be PLUS '+', got %s '%s'", parser.Tokens[1].Type, parser.Tokens[1].Value)
	}

	if parser.Tokens[2].Type != "NUMBER" || parser.Tokens[2].Value != "456" {
		t.Errorf("Expected third token to be NUMBER '456', got %s '%s'", parser.Tokens[2].Type, parser.Tokens[2].Value)
	}
}

func TestDSLRuleSequenceParsing(t *testing.T) {
	// Test case 4: Rule sequence parsing (the specific bug we fixed)
	grammar := NewDSLGrammar()

	// Add tokens
	grammar.AddToken("A", "a")
	grammar.AddToken("B", "b")
	grammar.AddToken("C", "c")

	// Add rule with sequence
	grammar.AddRule("sequence", []string{"A B C"}, "action")

	// Check that the rule was added correctly
	rule, exists := grammar.Rules["sequence"]
	if !exists {
		t.Fatal("Rule 'sequence' not found")
	}

	if len(rule.Alternatives) != 1 {
		t.Fatalf("Expected 1 alternative, got %d", len(rule.Alternatives))
	}

	alt := rule.Alternatives[0]
	if len(alt.Sequence) != 3 {
		t.Fatalf("Expected sequence of 3 elements, got %d", len(alt.Sequence))
	}

	if alt.Sequence[0] != "A" || alt.Sequence[1] != "B" || alt.Sequence[2] != "C" {
		t.Errorf("Expected sequence [A B C], got %v", alt.Sequence)
	}

	if alt.Action != "action" {
		t.Errorf("Expected action 'action', got '%s'", alt.Action)
	}
}

func TestDSLMultipleRuleAlternatives(t *testing.T) {
	// Test case 5: Multiple alternatives for the same rule
	grammar := NewDSLGrammar()

	// Add tokens
	grammar.AddToken("X", "x")
	grammar.AddToken("Y", "y")

	// Add multiple rules with same name (should create alternatives)
	grammar.AddRule("choice", []string{"X"}, "action1")
	grammar.AddRule("choice", []string{"Y"}, "action2")

	// Check that the rule has both alternatives
	rule, exists := grammar.Rules["choice"]
	if !exists {
		t.Fatal("Rule 'choice' not found")
	}

	if len(rule.Alternatives) != 2 {
		t.Fatalf("Expected 2 alternatives, got %d", len(rule.Alternatives))
	}

	// Check first alternative
	alt1 := rule.Alternatives[0]
	if len(alt1.Sequence) != 1 || alt1.Sequence[0] != "X" {
		t.Errorf("Expected first alternative to be [X], got %v", alt1.Sequence)
	}
	if alt1.Action != "action1" {
		t.Errorf("Expected first alternative action to be 'action1', got '%s'", alt1.Action)
	}

	// Check second alternative
	alt2 := rule.Alternatives[1]
	if len(alt2.Sequence) != 1 || alt2.Sequence[0] != "Y" {
		t.Errorf("Expected second alternative to be [Y], got %v", alt2.Sequence)
	}
	if alt2.Action != "action2" {
		t.Errorf("Expected second alternative action to be 'action2', got '%s'", alt2.Action)
	}
}

func TestDSLRegexTokens(t *testing.T) {
	// Test case 6: Regex validation for tokens
	grammar := NewDSLGrammar()

	// Test valid regex
	err := grammar.AddToken("VALID", "[0-9]+")
	if err != nil {
		t.Errorf("Valid regex should not cause error: %v", err)
	}

	// Test invalid regex
	err = grammar.AddToken("INVALID", "[0-9")
	if err == nil {
		t.Error("Invalid regex should cause error")
	}

	// Test that valid token works
	token, exists := grammar.Tokens["VALID"]
	if !exists {
		t.Error("VALID token should exist")
	}

	if token.Regex == nil {
		t.Error("Token should have compiled regex")
	}

	// Test regex matching
	matches := token.Regex.FindStringIndex("123abc")
	if matches == nil || matches[0] != 0 || matches[1] != 3 {
		t.Errorf("Regex should match '123' at beginning, got %v", matches)
	}
}

func TestDSLErrorHandling(t *testing.T) {
	// Test case 7: Error handling in DSL parsing
	grammar := NewDSLGrammar()
	grammar.AddToken("A", "a")

	parser := NewDSLParser(grammar)

	// Test parsing with invalid character
	err := parser.Tokenize("ax")
	if err == nil {
		t.Error("Should get error for invalid character")
	}

	if !strings.Contains(err.Error(), "unexpected character") {
		t.Errorf("Error should mention unexpected character, got: %v", err)
	}
}

func TestDSLParameterFormatting(t *testing.T) {
	// Test case 8: Parameter formatting (fix for {+} issue)
	env := NewEnvironment()

	dslCode := `
	dsl ParamTestDSL {
		token("A", "a")
		token("B", "b")
		token("PLUS", "\\+")
		
		rule("expr", ["A", "PLUS", "B"], "combine")
		
		func combine(left, op, right) {
			// Test that parameters are not wrapped in {}
			if (op == "+") {
				return left + op + right
			}
			return "failed"
		}
	}
	`

	parser := NewParser(dslCode)
	program := parser.ParseProgram()
	program.Eval(env)

	// Get the DSL and use it
	dslObj, exists := env.Get("ParamTestDSL")
	if !exists {
		t.Fatal("DSL 'ParamTestDSL' not found")
	}

	if dslDef, ok := dslObj.(*DSLDefinition); ok {
		result := dslDef.evaluateDSLCode("a + b", env)

		if dslResult, ok := result.(*DSLResult); ok {
			if dslResult.Output != "a+b" {
				t.Errorf("Expected 'a+b', got '%v'", dslResult.Output)
			}
		} else {
			t.Errorf("Expected DSLResult, got %T", result)
		}
	} else {
		t.Fatalf("Expected DSLDefinition, got %T", dslObj)
	}
}

func TestDSLResultAccess(t *testing.T) {
	// Test case 9: DSL result access and properties
	env := NewEnvironment()

	dslCode := `
	dsl ResultTestDSL {
		token("NUM", "[0-9]+")
		rule("number", ["NUM"], "double")
		
		func double(n) {
			return n + n
		}
	}
	`

	parser := NewParser(dslCode)
	program := parser.ParseProgram()
	program.Eval(env)

	dslObj, exists := env.Get("ResultTestDSL")
	if !exists {
		t.Fatal("DSL 'ResultTestDSL' not found")
	}

	if dslDef, ok := dslObj.(*DSLDefinition); ok {
		result := dslDef.evaluateDSLCode("5", env)

		if dslResult, ok := result.(*DSLResult); ok {
			// Test Output property
			if dslResult.Output != "55" {
				t.Errorf("Expected '55', got %v", dslResult.Output)
			}

			// Test Code property
			if dslResult.Code != "5" {
				t.Errorf("Expected '5', got '%s'", dslResult.Code)
			}

			// Test GetResult method
			if dslResult.GetResult() != "55" {
				t.Errorf("Expected '55' from GetResult(), got %v", dslResult.GetResult())
			}

			// Test String representation
			expected := "DSL[5] -> \"55\""
			if dslResult.String() != expected {
				t.Errorf("Expected '%s', got '%s'", expected, dslResult.String())
			}
		} else {
			t.Errorf("Expected DSLResult, got %T", result)
		}
	} else {
		t.Fatalf("Expected DSLDefinition, got %T", dslObj)
	}
}

func TestDSLStringFormatting(t *testing.T) {
	// Test case 10: String formatting for different types
	testCases := []struct {
		name     string
		output   interface{}
		code     string
		expected string
	}{
		{"integer", 42, "test", "DSL[test] -> 42"},
		{"float", 3.14, "pi", "DSL[pi] -> 3.14"},
		{"string", "hello", "greeting", "DSL[greeting] -> \"hello\""},
		{"boolean true", true, "condition", "DSL[condition] -> true"},
		{"boolean false", false, "condition", "DSL[condition] -> false"},
		{"nil", nil, "empty", "DSL[empty] -> <no result>"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := &DSLResult{
				Output: tc.output,
				Code:   tc.code,
			}

			if result.String() != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result.String())
			}
		})
	}
}

func TestDSLComplexCalculation(t *testing.T) {
	// Test case 11: Complex calculation DSL (integration test)
	env := NewEnvironment()

	dslCode := `
	dsl CalculatorDSL {
		token("NUMBER", "[0-9]+")
		token("PLUS", "\\+")
		token("MINUS", "-")
		token("MULTIPLY", "\\*")
		token("DIVIDE", "/")
		
		rule("expression", ["NUMBER", "operator", "NUMBER"], "calculate")
		rule("operator", ["PLUS"], "plus")
		rule("operator", ["MINUS"], "minus")
		rule("operator", ["MULTIPLY"], "multiply")
		rule("operator", ["DIVIDE"], "divide")
		
		func calculate(left, op, right) {
			// Simple string-based calculation for testing
			if (op == "+") {
				return left + "+" + right
			}
			if (op == "-") {
				return left + "-" + right
			}
			if (op == "*") {
				return left + "*" + right
			}
			if (op == "/") {
				return left + "/" + right
			}
			return "error"
		}
		
		func plus(token) { return "+" }
		func minus(token) { return "-" }
		func multiply(token) { return "*" }
		func divide(token) { return "/" }
	}
	`

	parser := NewParser(dslCode)
	program := parser.ParseProgram()
	program.Eval(env)

	dslObj, exists := env.Get("CalculatorDSL")
	if !exists {
		t.Fatal("DSL 'CalculatorDSL' not found")
	}

	if dslDef, ok := dslObj.(*DSLDefinition); ok {
		testCases := []struct {
			expression string
			expected   string
		}{
			{"5 + 3", "5+3"},
			{"10 - 4", "10-4"},
			{"6 * 7", "6*7"},
			{"15 / 3", "15/3"},
		}

		for _, tc := range testCases {
			t.Run(tc.expression, func(t *testing.T) {
				result := dslDef.evaluateDSLCode(tc.expression, env)

				if dslResult, ok := result.(*DSLResult); ok {
					if dslResult.Output != tc.expected {
						t.Errorf("Expected %v, got %v", tc.expected, dslResult.Output)
					}
				} else {
					t.Errorf("Expected DSLResult, got %T", result)
				}
			})
		}
	} else {
		t.Fatalf("Expected DSLDefinition, got %T", dslObj)
	}
}

func TestDSLReturnValueUnpacking(t *testing.T) {
	// Test case 12: Ensure ReturnValue objects are properly unpacked
	env := NewEnvironment()

	dslCode := `
	dsl UnpackTestDSL {
		token("A", "a")
		token("B", "b")
		
		rule("pair", ["A", "B"], "combine")
		
		func combine(first, second) {
			// This tests that parameters are unpacked correctly
			return first + ":" + second
		}
	}
	`

	parser := NewParser(dslCode)
	program := parser.ParseProgram()
	program.Eval(env)

	dslObj, exists := env.Get("UnpackTestDSL")
	if !exists {
		t.Fatal("DSL 'UnpackTestDSL' not found")
	}

	if dslDef, ok := dslObj.(*DSLDefinition); ok {
		result := dslDef.evaluateDSLCode("a b", env)

		if dslResult, ok := result.(*DSLResult); ok {
			expected := "a:b"
			if dslResult.Output != expected {
				t.Errorf("Expected '%s', got '%v'", expected, dslResult.Output)
			}
		} else {
			t.Errorf("Expected DSLResult, got %T", result)
		}
	} else {
		t.Fatalf("Expected DSLDefinition, got %T", dslObj)
	}
}
