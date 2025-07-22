package r2core

import (
	"fmt"
	"testing"
)

func TestDSLDebug(t *testing.T) {
	env := NewEnvironment()

	// Simple DSL that just returns the input
	dslCode := `
	dsl DebugDSL {
		token("WORD", "[a-zA-Z]+")
		rule("simple", ["WORD"], "returnWord")
		
		func returnWord(word) {
			return word;
		}
	}
	`

	parser := NewParser(dslCode)
	program := parser.ParseProgram()
	program.Eval(env)

	// Get DSL object
	dslObj, exists := env.Get("DebugDSL")
	if !exists {
		t.Fatal("DSL 'DebugDSL' not found")
	}

	// Get the DSL object directly
	dslMap, ok := dslObj.(map[string]interface{})
	if !ok {
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

	// Test simple usage
	result := useMethod("hello")

	// Print the actual result to understand what's happening
	fmt.Printf("Result type: %T\n", result)
	fmt.Printf("Result value: %v\n", result)

	if err, isErr := result.(error); isErr {
		t.Errorf("Got error: %v", err)
	}

	if dslResult, ok := result.(*DSLResult); ok {
		t.Logf("DSL Result Output: %v", dslResult.Output)
		t.Logf("DSL Result Code: %s", dslResult.Code)
		if dslResult.Output != "hello" {
			t.Errorf("Expected 'hello', got %v", dslResult.Output)
		}
	} else {
		t.Errorf("Expected DSLResult, got %T: %v", result, result)
	}
}
