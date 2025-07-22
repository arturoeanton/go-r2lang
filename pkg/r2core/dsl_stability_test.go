package r2core

import (
	"fmt"
	"testing"
	"time"
)

// TestDSLStabilityReproduction attempts to reproduce the intermittent DSL parsing issue
func TestDSLStabilityReproduction(t *testing.T) {
	// Create the DSL grammar that matches the LINQ example
	grammar := NewDSLGrammar()

	// Add tokens in order they appear in linq.r2
	grammar.AddToken("WORD", "[a-zA-Z_][a-zA-Z0-9_]*")
	grammar.AddToken("NUMBER", "[0-9]+")
	grammar.AddToken("OPERATOR", "[><=]+")

	// Add rules
	grammar.AddRule("query", []string{"WORD WORD WORD WORD WORD WORD OPERATOR NUMBER"}, "buildQuery")
	grammar.AddRule("query", []string{"WORD WORD WORD WORD"}, "buildSimpleQuery")

	// Add dummy actions
	grammar.AddAction("buildQuery", func(args []interface{}) interface{} {
		return fmt.Sprintf("Query with %d args", len(args))
	})
	grammar.AddAction("buildSimpleQuery", func(args []interface{}) interface{} {
		return fmt.Sprintf("Simple query with %d args", len(args))
	})

	// Test the same DSL code multiple times to detect intermittent behavior
	testCode := "select name from employees where salary > 50000"

	results := make([]interface{}, 0)
	errors := make([]error, 0)

	// Run the same parsing 50 times
	for i := 0; i < 50; i++ {
		parser := NewDSLParser(grammar)
		result, err := parser.Parse(testCode)

		if err != nil {
			errors = append(errors, err)
			t.Logf("Iteration %d: ERROR - %v", i+1, err)
		} else {
			results = append(results, result)
			t.Logf("Iteration %d: SUCCESS - %v", i+1, result)
		}

		// Add small delay to potentially trigger timing-related issues
		time.Sleep(1 * time.Millisecond)
	}

	// Analyze results
	successCount := len(results)
	errorCount := len(errors)

	t.Logf("Results after 50 iterations:")
	t.Logf("- Successes: %d", successCount)
	t.Logf("- Errors: %d", errorCount)

	// If we have both successes and errors, we've reproduced the intermittent issue
	if successCount > 0 && errorCount > 0 {
		t.Errorf("INTERMITTENT BEHAVIOR DETECTED: %d successes, %d errors", successCount, errorCount)

		// Show some example errors
		if len(errors) > 0 {
			t.Logf("Example error: %v", errors[0])
		}

		// Show unique error types
		errorTypes := make(map[string]int)
		for _, err := range errors {
			errorTypes[err.Error()]++
		}

		t.Logf("Error types:")
		for errType, count := range errorTypes {
			t.Logf("- %s: %d times", errType, count)
		}
	}

	// The test should be deterministic - either all pass or all fail
	if errorCount > 0 && successCount > 0 {
		t.Fatalf("DSL parsing is not deterministic! This indicates a serious stability issue.")
	}
}

// TestDSLTokenizationStability tests if tokenization is stable
func TestDSLTokenizationStability(t *testing.T) {
	grammar := NewDSLGrammar()
	grammar.AddToken("WORD", "[a-zA-Z_][a-zA-Z0-9_]*")
	grammar.AddToken("NUMBER", "[0-9]+")
	grammar.AddToken("OPERATOR", "[><=]+")

	testCode := "select name from employees where salary > 50000"

	var firstTokens []DSLTokenMatch

	// Test tokenization 20 times
	for i := 0; i < 20; i++ {
		parser := NewDSLParser(grammar)
		err := parser.Tokenize(testCode)

		if err != nil {
			t.Fatalf("Tokenization failed on iteration %d: %v", i+1, err)
		}

		if i == 0 {
			// Store first result as reference
			firstTokens = make([]DSLTokenMatch, len(parser.Tokens))
			copy(firstTokens, parser.Tokens)
		} else {
			// Compare with first result
			if len(parser.Tokens) != len(firstTokens) {
				t.Errorf("Token count mismatch on iteration %d: expected %d, got %d",
					i+1, len(firstTokens), len(parser.Tokens))
			}

			for j, token := range parser.Tokens {
				if j < len(firstTokens) {
					if token.Type != firstTokens[j].Type {
						t.Errorf("Token type mismatch on iteration %d, position %d: expected %s, got %s",
							i+1, j, firstTokens[j].Type, token.Type)
					}
					if token.Value != firstTokens[j].Value {
						t.Errorf("Token value mismatch on iteration %d, position %d: expected %s, got %s",
							i+1, j, firstTokens[j].Value, token.Value)
					}
				}
			}
		}
	}

	t.Logf("Tokenization is stable across 20 iterations")
	t.Logf("Expected tokens for '%s':", testCode)
	for i, token := range firstTokens {
		t.Logf("  %d: %s = '%s'", i, token.Type, token.Value)
	}
}

// TestDSLParsingStability tests if parsing (after tokenization) is stable
func TestDSLParsingStability(t *testing.T) {
	grammar := NewDSLGrammar()
	grammar.AddToken("WORD", "[a-zA-Z_][a-zA-Z0-9_]*")
	grammar.AddToken("NUMBER", "[0-9]+")
	grammar.AddToken("OPERATOR", "[><=]+")

	// Add rule that should match our test case
	grammar.AddRule("query", []string{"WORD WORD WORD WORD WORD WORD OPERATOR NUMBER"}, "buildQuery")
	grammar.AddAction("buildQuery", func(args []interface{}) interface{} {
		return fmt.Sprintf("Parsed with %d tokens", len(args))
	})

	testCode := "select name from employees where salary > 50000"

	// First, tokenize once and reuse tokens to isolate parsing issues
	parser := NewDSLParser(grammar)
	err := parser.Tokenize(testCode)
	if err != nil {
		t.Fatalf("Initial tokenization failed: %v", err)
	}

	originalTokens := make([]DSLTokenMatch, len(parser.Tokens))
	copy(originalTokens, parser.Tokens)

	// Test parsing 20 times with same tokens
	var firstResult interface{}

	for i := 0; i < 20; i++ {
		// Create new parser and copy tokens
		testParser := NewDSLParser(grammar)
		testParser.Tokens = make([]DSLTokenMatch, len(originalTokens))
		copy(testParser.Tokens, originalTokens)
		testParser.Pos = 0

		result, err := testParser.parseRule("query")

		if err != nil {
			t.Errorf("Parsing failed on iteration %d: %v", i+1, err)
		} else {
			if i == 0 {
				firstResult = result
			} else {
				// Check if result is consistent
				if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", firstResult) {
					t.Errorf("Result mismatch on iteration %d: expected %v, got %v",
						i+1, firstResult, result)
				}
			}
		}
	}

	t.Logf("Parsing is stable across 20 iterations with same tokens")
}
