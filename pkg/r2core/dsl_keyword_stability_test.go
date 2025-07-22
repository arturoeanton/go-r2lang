package r2core

import (
	"fmt"
	"testing"
)

// TestDSLWithSpecificKeywords tests DSL with specific keywords like SELECT, FROM, WHERE
func TestDSLWithSpecificKeywords(t *testing.T) {
	grammar := NewDSLGrammar()

	// Add keywords first (high priority)
	grammar.AddKeywordToken("SELECT", "select")
	grammar.AddKeywordToken("FROM", "from")
	grammar.AddKeywordToken("WHERE", "where")

	// Add generic patterns (lower priority)
	grammar.AddToken("IDENTIFIER", "[a-zA-Z_][a-zA-Z0-9_]*")
	grammar.AddToken("NUMBER", "[0-9]+")
	grammar.AddToken("OPERATOR", "[><=]+")

	// Add rules
	grammar.AddRule("query", []string{"SELECT IDENTIFIER FROM IDENTIFIER WHERE IDENTIFIER OPERATOR NUMBER"}, "buildQuery")
	grammar.AddRule("query", []string{"SELECT IDENTIFIER FROM IDENTIFIER"}, "buildSimpleQuery")

	// Add actions
	grammar.AddAction("buildQuery", func(args []interface{}) interface{} {
		return map[string]interface{}{
			"type": "complex_query",
			"args": args,
		}
	})
	grammar.AddAction("buildSimpleQuery", func(args []interface{}) interface{} {
		return map[string]interface{}{
			"type": "simple_query",
			"args": args,
		}
	})

	// Test cases with keyword-specific queries
	testCases := []struct {
		name          string
		code          string
		shouldSucceed bool
	}{
		{"Complex query with keywords", "select name from users where age > 25", true},
		{"Simple query with keywords", "select name from users", true},
		{"Case insensitive keywords", "SELECT name FROM users WHERE age > 30", true},
		{"Mixed case keywords", "Select name From users Where age > 40", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test this query multiple times to ensure stability
			for i := 0; i < 10; i++ {
				parser := NewDSLParser(grammar)
				result, err := parser.Parse(tc.code)

				if tc.shouldSucceed {
					if err != nil {
						t.Errorf("Iteration %d: Expected success, got error: %v", i+1, err)
					} else {
						// Verify result structure
						if resultMap, ok := result.(map[string]interface{}); ok {
							if resultMap["type"] == nil {
								t.Errorf("Iteration %d: Missing 'type' in result", i+1)
							}
							if resultMap["args"] == nil {
								t.Errorf("Iteration %d: Missing 'args' in result", i+1)
							}
						} else {
							t.Errorf("Iteration %d: Expected map result, got %T", i+1, result)
						}
					}
				} else {
					if err == nil {
						t.Errorf("Iteration %d: Expected error, got success: %v", i+1, result)
					}
				}
			}
		})
	}
}

// TestDSLTokenOrderInsensitivity tests that DSL tokens work regardless of order
func TestDSLTokenOrderInsensitivity(t *testing.T) {
	// Create two grammars with tokens added in different orders
	grammar1 := NewDSLGrammar()
	grammar2 := NewDSLGrammar()

	// Grammar 1: Add generic patterns first, then keywords
	grammar1.AddToken("IDENTIFIER", "[a-zA-Z_][a-zA-Z0-9_]*")
	grammar1.AddToken("NUMBER", "[0-9]+")
	grammar1.AddKeywordToken("SELECT", "select")
	grammar1.AddKeywordToken("FROM", "from")
	grammar1.AddKeywordToken("WHERE", "where")
	grammar1.AddToken("OPERATOR", "[><=]+")

	// Grammar 2: Add keywords first, then generic patterns
	grammar2.AddKeywordToken("SELECT", "select")
	grammar2.AddKeywordToken("FROM", "from")
	grammar2.AddKeywordToken("WHERE", "where")
	grammar2.AddToken("IDENTIFIER", "[a-zA-Z_][a-zA-Z0-9_]*")
	grammar2.AddToken("NUMBER", "[0-9]+")
	grammar2.AddToken("OPERATOR", "[><=]+")

	// Add same rules to both grammars
	rule := []string{"SELECT IDENTIFIER FROM IDENTIFIER WHERE IDENTIFIER OPERATOR NUMBER"}
	grammar1.AddRule("query", rule, "buildQuery")
	grammar2.AddRule("query", rule, "buildQuery")

	// Add same actions to both grammars
	action := func(args []interface{}) interface{} {
		return len(args) // Just return the number of args
	}
	grammar1.AddAction("buildQuery", action)
	grammar2.AddAction("buildQuery", action)

	// Test the same query with both grammars
	testCode := "select name from users where age > 30"

	// Test multiple times with both grammars
	for i := 0; i < 10; i++ {
		parser1 := NewDSLParser(grammar1)
		parser2 := NewDSLParser(grammar2)

		result1, err1 := parser1.Parse(testCode)
		result2, err2 := parser2.Parse(testCode)

		// Both should succeed
		if err1 != nil {
			t.Errorf("Grammar 1, iteration %d: %v", i+1, err1)
		}
		if err2 != nil {
			t.Errorf("Grammar 2, iteration %d: %v", i+1, err2)
		}

		// Both should produce the same result
		if err1 == nil && err2 == nil {
			if result1 != result2 {
				t.Errorf("Results differ between grammars on iteration %d: %v vs %v", i+1, result1, result2)
			}
		}

		// Check tokenization is the same
		parser1.Tokens = []DSLTokenMatch{}
		parser2.Tokens = []DSLTokenMatch{}

		err1 = parser1.Tokenize(testCode)
		err2 = parser2.Tokenize(testCode)

		if err1 != nil || err2 != nil {
			t.Errorf("Tokenization failed: grammar1=%v, grammar2=%v", err1, err2)
			continue
		}

		// Compare tokens
		if len(parser1.Tokens) != len(parser2.Tokens) {
			t.Errorf("Different token counts: grammar1=%d, grammar2=%d", len(parser1.Tokens), len(parser2.Tokens))
			continue
		}

		for j, token1 := range parser1.Tokens {
			token2 := parser2.Tokens[j]
			if token1.Type != token2.Type {
				t.Errorf("Token %d type differs: grammar1=%s, grammar2=%s", j, token1.Type, token2.Type)
			}
			if token1.Value != token2.Value {
				t.Errorf("Token %d value differs: grammar1=%s, grammar2=%s", j, token1.Value, token2.Value)
			}
		}
	}
}

// TestDSLEnvironmentIsolation tests that DSL executions don't interfere with each other
func TestDSLEnvironmentIsolation(t *testing.T) {
	// This test simulates the real DSL usage with contexts
	env := NewEnvironment()

	// Create DSL definition
	dsl := &DSLDefinition{
		Name:      &Identifier{Name: "TestDSL"},
		Grammar:   NewDSLGrammar(),
		Functions: make(map[string]*FunctionDeclaration),
		IsActive:  true,
		GlobalEnv: env,
	}

	// Add tokens and rules
	dsl.Grammar.AddKeywordToken("SELECT", "select")
	dsl.Grammar.AddKeywordToken("FROM", "from")
	dsl.Grammar.AddToken("IDENTIFIER", "[a-zA-Z_][a-zA-Z0-9_]*")
	dsl.Grammar.AddRule("query", []string{"SELECT IDENTIFIER FROM IDENTIFIER"}, "buildQuery")
	dsl.Grammar.AddAction("buildQuery", func(args []interface{}) interface{} {
		return "Query result: " + args[1].(string) + " from " + args[3].(string)
	})

	// Create the DSL object with use function
	dslObject := map[string]interface{}{
		"use": func(args ...interface{}) interface{} {
			if len(args) < 1 {
				return fmt.Errorf("DSL use requires at least one argument")
			}

			code, ok := args[0].(string)
			if !ok {
				return fmt.Errorf("DSL use: first argument must be a string")
			}

			// Create isolated environment for this execution
			execEnv := NewInnerEnv(env)

			// Handle context if provided
			if len(args) > 1 {
				if contextMap, ok := args[1].(map[string]interface{}); ok {
					execEnv.Set("context", contextMap)
				}
			}

			// Parse and return result
			parser := NewDSLParserWithContext(dsl.Grammar, dsl)
			result, err := parser.Parse(code)
			if err != nil {
				return err
			}
			return result
		},
		"grammar": dsl.Grammar,
	}

	// Test multiple concurrent DSL uses
	testCases := []struct {
		code    string
		context map[string]interface{}
	}{
		{"select name from users", map[string]interface{}{"table": "users"}},
		{"select id from products", map[string]interface{}{"table": "products"}},
		{"select title from articles", map[string]interface{}{"table": "articles"}},
	}

	// Run multiple DSL uses concurrently-like (sequentially but rapidly)
	for iteration := 0; iteration < 20; iteration++ {
		for i, tc := range testCases {
			useFunc := dslObject["use"].(func(args ...interface{}) interface{})
			result := useFunc(tc.code, tc.context)

			// Check that result is valid
			if err, isError := result.(error); isError {
				t.Errorf("Iteration %d, case %d: DSL use failed: %v", iteration, i, err)
			} else {
				// Result should be consistent for the same input
				expected := fmt.Sprintf("Query result: %s from %s",
					[]string{"name", "id", "title"}[i],
					[]string{"users", "products", "articles"}[i])
				if result != expected {
					t.Errorf("Iteration %d, case %d: Expected '%s', got '%s'",
						iteration, i, expected, result)
				}
			}
		}
	}
}
