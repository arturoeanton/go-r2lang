package r2core

import (
	"fmt"
	"testing"
	"time"
)

// TestDSLStabilityReproduction attempts to reproduce the intermittent DSL parsing issue
// that used to affect the hand-rolled backtracking engine. go-dsl's deterministic
// tokenizer (priority > length > declaration order) plus its Packrat/AST-based
// parser should make repeated parses of the same input fully deterministic.
func TestDSLStabilityReproduction(t *testing.T) {
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

	testCode := "select name from employees where salary > 50000"

	var successCount, errorCount int
	var firstError error

	for i := 0; i < 50; i++ {
		result, err := grammar.Use(testCode, nil)
		if err != nil {
			errorCount++
			if firstError == nil {
				firstError = err
			}
			t.Logf("Iteration %d: ERROR - %v", i+1, err)
		} else {
			successCount++
			t.Logf("Iteration %d: SUCCESS - %v", i+1, result.Output)
		}

		time.Sleep(1 * time.Millisecond)
	}

	t.Logf("Results after 50 iterations: %d successes, %d errors", successCount, errorCount)

	if successCount > 0 && errorCount > 0 {
		t.Fatalf("INTERMITTENT BEHAVIOR DETECTED: %d successes, %d errors (example error: %v)",
			successCount, errorCount, firstError)
	}
}

// TestDSLTokenizationStability tests that tokenization is stable across repeated calls.
func TestDSLTokenizationStability(t *testing.T) {
	grammar := NewDSLGrammar()
	grammar.AddToken("WORD", "[a-zA-Z_][a-zA-Z0-9_]*")
	grammar.AddToken("NUMBER", "[0-9]+")
	grammar.AddToken("OPERATOR", "[><=]+")

	testCode := "select name from employees where salary > 50000"

	first, err := grammar.DebugTokens(testCode)
	if err != nil {
		t.Fatalf("Initial tokenization failed: %v", err)
	}

	for i := 1; i < 20; i++ {
		tokens, err := grammar.DebugTokens(testCode)
		if err != nil {
			t.Fatalf("Tokenization failed on iteration %d: %v", i+1, err)
		}

		if len(tokens) != len(first) {
			t.Errorf("Token count mismatch on iteration %d: expected %d, got %d",
				i+1, len(first), len(tokens))
			continue
		}

		for j, token := range tokens {
			if token.TokenType != first[j].TokenType {
				t.Errorf("Token type mismatch on iteration %d, position %d: expected %s, got %s",
					i+1, j, first[j].TokenType, token.TokenType)
			}
			if token.Value != first[j].Value {
				t.Errorf("Token value mismatch on iteration %d, position %d: expected %s, got %s",
					i+1, j, first[j].Value, token.Value)
			}
		}
	}

	t.Logf("Tokenization is stable across 20 iterations")
}

// TestDSLParsingStability tests that end-to-end parse+evaluate is stable across repeated calls.
func TestDSLParsingStability(t *testing.T) {
	grammar := NewDSLGrammar()
	grammar.AddToken("WORD", "[a-zA-Z_][a-zA-Z0-9_]*")
	grammar.AddToken("NUMBER", "[0-9]+")
	grammar.AddToken("OPERATOR", "[><=]+")

	grammar.AddRule("query", []string{"WORD WORD WORD WORD WORD WORD OPERATOR NUMBER"}, "buildQuery")
	grammar.AddAction("buildQuery", func(args []interface{}) interface{} {
		return fmt.Sprintf("Parsed with %d tokens", len(args))
	})

	testCode := "select name from employees where salary > 50000"

	var firstResult interface{}

	for i := 0; i < 20; i++ {
		result, err := grammar.Use(testCode, nil)
		if err != nil {
			t.Errorf("Parsing failed on iteration %d: %v", i+1, err)
			continue
		}

		if i == 0 {
			firstResult = result.Output
		} else if fmt.Sprintf("%v", result.Output) != fmt.Sprintf("%v", firstResult) {
			t.Errorf("Result mismatch on iteration %d: expected %v, got %v",
				i+1, firstResult, result.Output)
		}
	}

	t.Logf("Parsing is stable across 20 iterations")
}
