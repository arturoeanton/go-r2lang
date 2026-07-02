package r2core

import (
	"fmt"
	"sync"
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
				result, err := grammar.Use(tc.code, nil)

				if tc.shouldSucceed {
					if err != nil {
						t.Errorf("Iteration %d: Expected success, got error: %v", i+1, err)
					} else if resultMap, ok := result.Output.(map[string]interface{}); ok {
						if resultMap["type"] == nil {
							t.Errorf("Iteration %d: Missing 'type' in result", i+1)
						}
						if resultMap["args"] == nil {
							t.Errorf("Iteration %d: Missing 'args' in result", i+1)
						}
					} else {
						t.Errorf("Iteration %d: Expected map result, got %T", i+1, result.Output)
					}
				} else if err == nil {
					t.Errorf("Iteration %d: Expected error, got success: %v", i+1, result.Output)
				}
			}
		})
	}
}

// TestDSLTokenOrderInsensitivity tests that DSL tokens work regardless of declaration order
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

	testCode := "select name from users where age > 30"

	for i := 0; i < 10; i++ {
		result1, err1 := grammar1.Use(testCode, nil)
		result2, err2 := grammar2.Use(testCode, nil)

		if err1 != nil {
			t.Errorf("Grammar 1, iteration %d: %v", i+1, err1)
		}
		if err2 != nil {
			t.Errorf("Grammar 2, iteration %d: %v", i+1, err2)
		}

		if err1 == nil && err2 == nil && result1.Output != result2.Output {
			t.Errorf("Results differ between grammars on iteration %d: %v vs %v",
				i+1, result1.Output, result2.Output)
		}

		// Compare tokenization
		tokens1, err1 := grammar1.DebugTokens(testCode)
		tokens2, err2 := grammar2.DebugTokens(testCode)

		if err1 != nil || err2 != nil {
			t.Errorf("Tokenization failed: grammar1=%v, grammar2=%v", err1, err2)
			continue
		}

		if len(tokens1) != len(tokens2) {
			t.Errorf("Different token counts: grammar1=%d, grammar2=%d", len(tokens1), len(tokens2))
			continue
		}

		for j, token1 := range tokens1 {
			token2 := tokens2[j]
			if token1.TokenType != token2.TokenType {
				t.Errorf("Token %d type differs: grammar1=%s, grammar2=%s", j, token1.TokenType, token2.TokenType)
			}
			if token1.Value != token2.Value {
				t.Errorf("Token %d value differs: grammar1=%s, grammar2=%s", j, token1.Value, token2.Value)
			}
		}
	}
}

// TestDSLEnvironmentIsolation tests that DSL executions don't interfere with each other
func TestDSLEnvironmentIsolation(t *testing.T) {
	env := NewEnvironment()

	dsl := &DSLDefinition{
		Name:      &Identifier{Name: "TestDSL"},
		Grammar:   NewDSLGrammar(),
		Functions: make(map[string]*FunctionDeclaration),
		IsActive:  true,
		GlobalEnv: env,
	}

	dsl.Grammar.AddKeywordToken("SELECT", "select")
	dsl.Grammar.AddKeywordToken("FROM", "from")
	dsl.Grammar.AddToken("IDENTIFIER", "[a-zA-Z_][a-zA-Z0-9_]*")
	dsl.Grammar.AddRule("query", []string{"SELECT IDENTIFIER FROM IDENTIFIER"}, "buildQuery")
	dsl.Grammar.AddAction("buildQuery", func(args []interface{}) interface{} {
		return "Query result: " + args[1].(string) + " from " + args[3].(string)
	})

	use := func(code string, context map[string]interface{}) interface{} {
		result, err := dsl.Grammar.Use(code, context)
		if err != nil {
			return err
		}
		return result.Output
	}

	testCases := []struct {
		code    string
		context map[string]interface{}
	}{
		{"select name from users", map[string]interface{}{"table": "users"}},
		{"select id from products", map[string]interface{}{"table": "products"}},
		{"select title from articles", map[string]interface{}{"table": "articles"}},
	}

	// Run multiple DSL uses in quick succession to check for state bleed
	for iteration := 0; iteration < 20; iteration++ {
		for i, tc := range testCases {
			result := use(tc.code, tc.context)

			if err, isError := result.(error); isError {
				t.Errorf("Iteration %d, case %d: DSL use failed: %v", iteration, i, err)
				continue
			}

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

// TestDSLConcurrentUseIsolation calls .use() on a single shared DSLDefinition
// from many goroutines at once (as R2Lang scripts can via `go`/`r2`) and
// checks each call sees its own environment, not another goroutine's. Run
// with -race to also confirm there's no data race on currentExecutionEnv.
func TestDSLConcurrentUseIsolation(t *testing.T) {
	env := NewEnvironment()

	dsl := &DSLDefinition{
		Name:      &Identifier{Name: "ConcurrentDSL"},
		Grammar:   NewDSLGrammar(),
		Functions: make(map[string]*FunctionDeclaration),
		IsActive:  true,
		GlobalEnv: env,
	}

	dsl.Grammar.AddToken("WORD", "[a-zA-Z]+")
	dsl.Grammar.AddRule("say", []string{"WORD"}, "echo")
	dsl.Grammar.AddAction("echo", func(args []interface{}) interface{} {
		return args[0]
	})

	words := []string{"alpha", "beta", "gamma", "delta", "epsilon"}

	var wg sync.WaitGroup
	errCh := make(chan string, len(words)*20)

	for round := 0; round < 20; round++ {
		for _, w := range words {
			wg.Add(1)
			go func(word string) {
				defer wg.Done()

				execEnv := NewInnerEnv(env)
				execEnv.Set("context", map[string]interface{}{})
				result := dsl.evaluateDSLCode(word, execEnv)

				dslResult, ok := result.(*DSLResult)
				if !ok {
					errCh <- fmt.Sprintf("expected *DSLResult for %q, got %T (%v)", word, result, result)
					return
				}
				if dslResult.Output != word {
					errCh <- fmt.Sprintf("expected output %q, got %q", word, dslResult.Output)
				}
			}(w)
		}
	}

	wg.Wait()
	close(errCh)

	for msg := range errCh {
		t.Error(msg)
	}
}
