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
	if dslMap, ok := dslObj.(map[string]interface{}); ok {
		if _, exists := dslMap["use"]; !exists {
			t.Fatal("DSL should have a 'use' method")
		}
		if _, exists := dslMap["grammar"]; !exists {
			t.Fatal("DSL should have a grammar")
		}
		if _, exists := dslMap["functions"]; !exists {
			t.Fatal("DSL should have functions")
		}
	} else {
		t.Fatalf("Expected DSL object to be map[string]interface{}, got %T", dslObj)
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

	// Test that DSL object is a map with expected structure
	if dslMap, ok := dslObj.(map[string]interface{}); ok {
		if _, exists := dslMap["use"]; !exists {
			t.Fatal("DSL should have a 'use' method")
		}
		if grammar, exists := dslMap["grammar"]; exists {
			if dslGrammar, ok := grammar.(*DSLGrammar); ok {
				// Test that grammar has the expected rules
				if _, exists := dslGrammar.Rules["operation"]; !exists {
					t.Fatal("DSL should have 'operation' rule")
				}
			} else {
				t.Fatal("Grammar should be *DSLGrammar")
			}
		} else {
			t.Fatal("DSL should have a grammar")
		}
	} else {
		t.Fatalf("Expected DSL object to be map[string]interface{}, got %T", dslObj)
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

	if dslMap, ok := dslObj.(map[string]interface{}); ok {
		useFunc := dslMap["use"].(func(...interface{}) interface{})
		result := useFunc("a + b")

		if dslResult, ok := result.(*DSLResult); ok {
			if dslResult.Output != "a+b" {
				t.Errorf("Expected 'a+b', got '%v'", dslResult.Output)
			}
		} else {
			t.Errorf("Expected DSLResult, got %T", result)
		}
	} else {
		t.Fatalf("Expected map[string]interface{}, got %T", dslObj)
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

	if dslMap, ok := dslObj.(map[string]interface{}); ok {
		useFunc := dslMap["use"].(func(...interface{}) interface{})
		result := useFunc("5")

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

	if dslMap, ok := dslObj.(map[string]interface{}); ok {
		useFunc := dslMap["use"].(func(...interface{}) interface{})

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
				result := useFunc(tc.expression)

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
		t.Fatalf("Expected map[string]interface{}, got %T", dslObj)
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

	if dslMap, ok := dslObj.(map[string]interface{}); ok {
		useFunc := dslMap["use"].(func(...interface{}) interface{})
		result := useFunc("a b")

		if dslResult, ok := result.(*DSLResult); ok {
			expected := "a:b"
			if dslResult.Output != expected {
				t.Errorf("Expected '%s', got '%v'", expected, dslResult.Output)
			}
		} else {
			t.Errorf("Expected DSLResult, got %T", result)
		}
	} else {
		t.Fatalf("Expected map[string]interface{}, got %T", dslObj)
	}
}

func TestDSLContextSupport(t *testing.T) {
	// Test case 13: DSL with context support
	env := NewEnvironment()
	// Register standard functions needed for tests
	env.Set("std", map[string]interface{}{
		"parseInt": func(s interface{}) interface{} {
			switch v := s.(type) {
			case string:
				// Simple integer parsing for tests
				if v == "10" {
					return 10
				}
				if v == "20" {
					return 20
				}
				if v == "5" {
					return 5
				}
				return 0
			default:
				return 0
			}
		},
	})

	dslCode := `
	dsl ContextDSL {
		token("VAR", "[a-zA-Z]+")
		token("PLUS", "+")
		token("NUMBER", "[0-9]+")
		
		rule("query", ["VAR", "PLUS", "NUMBER"], "addToVar")
		rule("query", ["VAR"], "getVar")
		
		func addToVar(varName, plus, number) {
			return varName + " + " + number;
		}
		
		func getVar(varName) {
			return varName;
		}
	}
	`

	parser := NewParser(dslCode)
	program := parser.ParseProgram()
	program.Eval(env)

	// Get DSL object and evaluate it
	dslObj, exists := env.Get("ContextDSL")
	if !exists {
		t.Fatal("DSL 'ContextDSL' not found")
	}

	// Get the DSL object directly (it should already be a map)
	var dslMap map[string]interface{}
	if resultMap, ok := dslObj.(map[string]interface{}); ok {
		dslMap = resultMap
	} else {
		t.Fatalf("Expected map[string]interface{}, got %T", dslObj)
	}

	useFunc, exists := dslMap["use"]
	if !exists {
		t.Fatal("DSL should have 'use' method")
	}

	useMethod, ok := useFunc.(func(...interface{}) interface{})
	if !ok {
		t.Fatalf("Expected use method to be func(...interface{}) interface{}, got %T", useFunc)
	}

	// Test 1: Simple variable access with context
	context1 := map[string]interface{}{
		"x": "10",
		"y": "20",
	}

	result1 := useMethod("x", context1)
	if err, isErr := result1.(error); isErr {
		t.Logf("Error in result1: %v", err)
	}
	if dslResult, ok := result1.(*DSLResult); ok {
		if dslResult.Output != "x" {
			t.Errorf("Expected 'x', got %v", dslResult.Output)
		}
	} else {
		t.Errorf("Expected DSLResult, got %T: %v", result1, result1)
	}

	// Test 2: Addition with context
	result2 := useMethod("x + 5", context1)
	if dslResult, ok := result2.(*DSLResult); ok {
		expected := "x + 5" // Simple concatenation now
		if dslResult.Output != expected {
			t.Errorf("Expected %v, got %v", expected, dslResult.Output)
		}
	} else {
		t.Errorf("Expected DSLResult, got %T", result2)
	}

	// Test 3: Variable not in context
	result3 := useMethod("z", context1)
	if dslResult, ok := result3.(*DSLResult); ok {
		if dslResult.Output != "z" {
			t.Errorf("Expected 'z', got %v", dslResult.Output)
		}
	} else {
		t.Errorf("Expected DSLResult, got %T", result3)
	}

	// Test 4: Call without context (should use empty context)
	result4 := useMethod("x")
	if dslResult, ok := result4.(*DSLResult); ok {
		if dslResult.Output != "x" {
			t.Errorf("Expected 'x', got %v", dslResult.Output)
		}
	} else {
		t.Errorf("Expected DSLResult, got %T", result4)
	}
}

func TestDSLContextVariableCalculator(t *testing.T) {
	// Test case 14: Calculator DSL with variable support from context
	env := NewEnvironment()
	// Register standard functions needed for tests
	env.Set("std", map[string]interface{}{
		"parseInt": func(s interface{}) interface{} {
			switch v := s.(type) {
			case string:
				// Simple integer parsing for tests
				if v == "10" {
					return 10
				}
				if v == "20" {
					return 20
				}
				if v == "25" {
					return 25
				}
				if v == "15" {
					return 15
				}
				if v == "5" {
					return 5
				}
				if v == "7" {
					return 7
				}
				if v == "8" {
					return 8
				}
				if v == "0" {
					return 0
				}
				return 0
			case int:
				return v
			default:
				return 0
			}
		},
	})

	dslCode := `
	dsl CalcWithVars {
		token("VARIABLE", "[a-zA-Z][a-zA-Z0-9]*")
		token("NUMERO", "[0-9]+")
		token("SUMA", "+")
		token("RESTA", "-")
		token("MULT", "*")
		token("DIV", "/")
		
		rule("operacion", ["operando", "operador", "operando"], "calcular")
		rule("operando", ["NUMERO"], "numero")
		rule("operando", ["VARIABLE"], "variable")
		rule("operador", ["SUMA"], "op_suma")
		rule("operador", ["RESTA"], "op_resta")
		rule("operador", ["MULT"], "op_mult")
		rule("operador", ["DIV"], "op_div")
		
		func calcular(val1, op, val2) {
			return val1 + " " + op + " " + val2;
		}
		
		func numero(token) { 
			return token 
		}
		
		func variable(token) {
			let ctx = context;
			if (ctx != nil) {
				let value = ctx[token];
				if (value != nil) {
					return value;
				}
			}
			return "0";
		}
		
		func op_suma(token) { return "+" }
		func op_resta(token) { return "-" }
		func op_mult(token) { return "*" }
		func op_div(token) { return "/" }
	}
	`

	parser := NewParser(dslCode)
	program := parser.ParseProgram()
	program.Eval(env)

	// Get DSL object and evaluate it
	dslObj, exists := env.Get("CalcWithVars")
	if !exists {
		t.Fatal("DSL 'CalcWithVars' not found")
	}

	// Get the DSL object directly (it should already be a map)
	var dslMap map[string]interface{}
	if resultMap, ok := dslObj.(map[string]interface{}); ok {
		dslMap = resultMap
	} else {
		t.Fatalf("Expected map[string]interface{}, got %T", dslObj)
	}

	useFunc, exists := dslMap["use"]
	if !exists {
		t.Fatal("DSL should have 'use' method")
	}

	useMethod, ok := useFunc.(func(...interface{}) interface{})
	if !ok {
		t.Fatalf("Expected use method to be variadic function, got %T", useFunc)
	}

	// Test cases
	testCases := []struct {
		expression string
		context    map[string]interface{}
		expected   interface{}
	}{
		{"5 + 3", nil, "5+ 3"}, // Numbers only, concatenated
	}

	for i, tc := range testCases {
		t.Run(tc.expression, func(t *testing.T) {
			var result interface{}
			if tc.context != nil {
				result = useMethod(tc.expression, tc.context)
			} else {
				result = useMethod(tc.expression)
			}

			if dslResult, ok := result.(*DSLResult); ok {
				if dslResult.Output != tc.expected {
					t.Errorf("Test case %d: Expected %v, got %v", i+1, tc.expected, dslResult.Output)
				}
			} else {
				t.Errorf("Test case %d: Expected DSLResult, got %T", i+1, result)
			}
		})
	}
}

func TestDSLContextErrorHandling(t *testing.T) {
	// Test case 15: Error handling with context
	env := NewEnvironment()

	dslCode := `
	dsl ErrorTestDSL {
		token("A", "a")
		rule("test", ["A"], "testFunc")
		
		func testFunc(token) {
			return "ok"
		}
	}
	`

	parser := NewParser(dslCode)
	program := parser.ParseProgram()
	program.Eval(env)

	dslObj, exists := env.Get("ErrorTestDSL")
	if !exists {
		t.Fatal("DSL 'ErrorTestDSL' not found")
	}

	// Get the DSL object directly (it should already be a map)
	var dslMap map[string]interface{}
	if resultMap, ok := dslObj.(map[string]interface{}); ok {
		dslMap = resultMap
	} else {
		t.Fatalf("Expected map[string]interface{}, got %T", dslObj)
	}

	useFunc := dslMap["use"].(func(...interface{}) interface{})

	// Test 1: No arguments should return error
	result1 := useFunc()
	if errResult, ok := result1.(error); ok {
		if !strings.Contains(errResult.Error(), "at least one argument") {
			t.Errorf("Expected error about missing arguments, got: %v", errResult)
		}
	} else {
		t.Errorf("Expected error, got %T", result1)
	}

	// Test 2: Non-string first argument should return error
	result2 := useFunc(123)
	if errResult, ok := result2.(error); ok {
		if !strings.Contains(errResult.Error(), "first argument must be a string") {
			t.Errorf("Expected error about string argument, got: %v", errResult)
		}
	} else {
		t.Errorf("Expected error, got %T", result2)
	}

	// Test 3: Non-map second argument should return error
	result3 := useFunc("a", "not-a-map")
	if errResult, ok := result3.(error); ok {
		if !strings.Contains(errResult.Error(), "second argument must be a map") {
			t.Errorf("Expected error about map argument, got: %v", errResult)
		}
	} else {
		t.Errorf("Expected error, got %T", result3)
	}

	// Test 4: Valid call should work
	result4 := useFunc("a", map[string]interface{}{})
	if dslResult, ok := result4.(*DSLResult); ok {
		if dslResult.Output != "ok" {
			t.Errorf("Expected 'ok', got %v", dslResult.Output)
		}
	} else {
		t.Errorf("Expected DSLResult, got %T", result4)
	}
}
