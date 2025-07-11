package r2core

import (
	"testing"
)

func TestIdentifier_Eval_ExistingVariable(t *testing.T) {
	env := NewEnvironment()

	// Test different types of variables
	testCases := []struct {
		name          string
		variableName  string
		variableValue interface{}
		expected      interface{}
	}{
		{"string variable", "name", "John", "John"},
		{"number variable", "age", 42.0, 42.0},
		{"boolean variable", "active", true, true},
		{"nil variable", "empty", nil, nil},
		{"array variable", "items", []interface{}{1, 2, 3}, []interface{}{1, 2, 3}},
		{"map variable", "config", map[string]interface{}{"key": "value"}, map[string]interface{}{"key": "value"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set the variable in environment
			env.Set(tc.variableName, tc.variableValue)

			// Create identifier and evaluate
			identifier := &Identifier{Name: tc.variableName}
			result := identifier.Eval(env)

			// Check result
			switch expected := tc.expected.(type) {
			case []interface{}:
				actual, ok := result.([]interface{})
				if !ok {
					t.Errorf("Expected []interface{}, got %T", result)
				} else if len(actual) != len(expected) {
					t.Errorf("Expected array length %d, got %d", len(expected), len(actual))
				} else {
					for i, v := range expected {
						if actual[i] != v {
							t.Errorf("Expected array[%d] = %v, got %v", i, v, actual[i])
						}
					}
				}
			case map[string]interface{}:
				actual, ok := result.(map[string]interface{})
				if !ok {
					t.Errorf("Expected map[string]interface{}, got %T", result)
				} else if len(actual) != len(expected) {
					t.Errorf("Expected map length %d, got %d", len(expected), len(actual))
				} else {
					for k, v := range expected {
						if actual[k] != v {
							t.Errorf("Expected map[%q] = %v, got %v", k, v, actual[k])
						}
					}
				}
			default:
				if result != tc.expected {
					t.Errorf("Expected %v, got %v", tc.expected, result)
				}
			}
		})
	}
}

func TestIdentifier_Eval_UndeclaredVariable(t *testing.T) {
	env := NewEnvironment()

	identifier := &Identifier{Name: "undeclaredVariable"}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for undeclared variable")
		} else {
			expectedMsg := "Undeclared variable: undeclaredVariable"
			if r != expectedMsg {
				t.Errorf("Expected %q panic, got %v", expectedMsg, r)
			}
		}
	}()

	identifier.Eval(env)
}

func TestIdentifier_Eval_NestedScopes(t *testing.T) {
	outerEnv := NewEnvironment()
	innerEnv := NewInnerEnv(outerEnv)

	// Set variable in outer scope
	outerEnv.Set("outerVar", "outerValue")

	// Set variable in inner scope
	innerEnv.Set("innerVar", "innerValue")

	// Test accessing outer variable from inner scope
	outerIdentifier := &Identifier{Name: "outerVar"}
	result := outerIdentifier.Eval(innerEnv)
	if result != "outerValue" {
		t.Errorf("Expected 'outerValue', got %v", result)
	}

	// Test accessing inner variable from inner scope
	innerIdentifier := &Identifier{Name: "innerVar"}
	result = innerIdentifier.Eval(innerEnv)
	if result != "innerValue" {
		t.Errorf("Expected 'innerValue', got %v", result)
	}

	// Test that inner variable is not accessible from outer scope
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when accessing inner variable from outer scope")
		}
	}()

	innerIdentifier.Eval(outerEnv)
}

func TestIdentifier_Eval_VariableShadowing(t *testing.T) {
	outerEnv := NewEnvironment()
	innerEnv := NewInnerEnv(outerEnv)

	// Set variable in outer scope
	outerEnv.Set("sharedVar", "outerValue")

	// Shadow the variable in inner scope
	innerEnv.Set("sharedVar", "innerValue")

	identifier := &Identifier{Name: "sharedVar"}

	// Test that inner scope returns shadowed value
	result := identifier.Eval(innerEnv)
	if result != "innerValue" {
		t.Errorf("Expected 'innerValue' (shadowed), got %v", result)
	}

	// Test that outer scope returns original value
	result = identifier.Eval(outerEnv)
	if result != "outerValue" {
		t.Errorf("Expected 'outerValue', got %v", result)
	}
}

func TestIdentifier_Eval_MultipleNestedScopes(t *testing.T) {
	// Create chain: global -> middle -> inner
	global := NewEnvironment()
	middle := NewInnerEnv(global)
	inner := NewInnerEnv(middle)

	// Set variables at different levels
	global.Set("globalVar", "globalValue")
	middle.Set("middleVar", "middleValue")
	inner.Set("innerVar", "innerValue")

	// Test accessing from innermost scope
	tests := []struct {
		name     string
		varName  string
		expected string
	}{
		{"global variable from inner", "globalVar", "globalValue"},
		{"middle variable from inner", "middleVar", "middleValue"},
		{"inner variable from inner", "innerVar", "innerValue"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			identifier := &Identifier{Name: test.varName}
			result := identifier.Eval(inner)
			if result != test.expected {
				t.Errorf("Expected %q, got %v", test.expected, result)
			}
		})
	}
}

func TestIdentifier_Eval_EmptyVariableName(t *testing.T) {
	env := NewEnvironment()

	identifier := &Identifier{Name: ""}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for empty variable name")
		} else {
			expectedMsg := "Undeclared variable: "
			if r != expectedMsg {
				t.Errorf("Expected %q panic, got %v", expectedMsg, r)
			}
		}
	}()

	identifier.Eval(env)
}

func TestIdentifier_Eval_SpecialCharacterNames(t *testing.T) {
	env := NewEnvironment()

	// Test identifiers with special characters that are valid
	testCases := []struct {
		name  string
		value string
	}{
		{"_private", "private value"},
		{"$global", "global value"},
		{"var123", "numbered variable"},
		{"camelCase", "camel case value"},
		{"snake_case", "snake case value"},
		{"MixedCase123", "mixed case value"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			env.Set(tc.name, tc.value)

			identifier := &Identifier{Name: tc.name}
			result := identifier.Eval(env)

			if result != tc.value {
				t.Errorf("Expected %q, got %v", tc.value, result)
			}
		})
	}
}

func TestIdentifier_Eval_CaseSensitivity(t *testing.T) {
	env := NewEnvironment()

	// Set variables with different cases
	env.Set("Variable", "uppercase")
	env.Set("variable", "lowercase")
	env.Set("VARIABLE", "allcaps")

	tests := []struct {
		name     string
		expected string
	}{
		{"Variable", "uppercase"},
		{"variable", "lowercase"},
		{"VARIABLE", "allcaps"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			identifier := &Identifier{Name: test.name}
			result := identifier.Eval(env)

			if result != test.expected {
				t.Errorf("Expected %q, got %v", test.expected, result)
			}
		})
	}

	// Test that non-existent case variations fail
	nonExistentIdentifier := &Identifier{Name: "VaRiAbLe"}
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for non-existent case variation")
		}
	}()

	nonExistentIdentifier.Eval(env)
}

func TestIdentifier_Eval_FunctionValues(t *testing.T) {
	env := NewEnvironment()

	// Create a simple function literal
	fl := &FunctionLiteral{
		Args: []string{"x"},
		Body: &BlockStatement{Statements: []Node{}},
	}

	// Store function in environment
	functionValue := fl.Eval(env)
	env.Set("myFunction", functionValue)

	// Retrieve function through identifier
	identifier := &Identifier{Name: "myFunction"}
	result := identifier.Eval(env)

	// Verify it's the same function
	userFunc, ok := result.(*UserFunction)
	if !ok {
		t.Fatalf("Expected *UserFunction, got %T", result)
	}

	if len(userFunc.Args) != 1 || userFunc.Args[0] != "x" {
		t.Errorf("Function not preserved correctly")
	}
}

func TestIdentifier_Eval_ModificationAfterRetrieval(t *testing.T) {
	env := NewEnvironment()

	// Set a map in the environment
	originalMap := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	env.Set("myMap", originalMap)

	// Retrieve map through identifier
	identifier := &Identifier{Name: "myMap"}
	result := identifier.Eval(env)

	// Modify the retrieved map
	retrievedMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got %T", result)
	}

	retrievedMap["key3"] = "value3"

	// Verify that the original map in environment is also modified (reference semantics)
	result2 := identifier.Eval(env)
	retrievedMap2, ok := result2.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got %T", result2)
	}

	if retrievedMap2["key3"] != "value3" {
		t.Error("Expected map modifications to persist (reference semantics)")
	}
}

// Benchmark tests
func BenchmarkIdentifier_Eval_LocalScope(b *testing.B) {
	env := NewEnvironment()
	env.Set("testVar", "testValue")

	identifier := &Identifier{Name: "testVar"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		identifier.Eval(env)
	}
}

func BenchmarkIdentifier_Eval_NestedScope(b *testing.B) {
	// Create deep nesting (10 levels)
	env := NewEnvironment()
	current := env
	for i := 0; i < 10; i++ {
		current = NewInnerEnv(current)
	}

	env.Set("deepVar", "deepValue")
	identifier := &Identifier{Name: "deepVar"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		identifier.Eval(current)
	}
}

func BenchmarkIdentifier_Eval_ShadowedVariable(b *testing.B) {
	// Create environment with shadowed variables
	outer := NewEnvironment()
	inner := NewInnerEnv(outer)

	outer.Set("shadowed", "outerValue")
	inner.Set("shadowed", "innerValue")

	identifier := &Identifier{Name: "shadowed"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		identifier.Eval(inner)
	}
}

func BenchmarkIdentifier_Eval_LongVariableName(b *testing.B) {
	env := NewEnvironment()

	// Create a very long variable name
	longName := "veryLongVariableNameThatIsReallyReallyLongAndShouldTestPerformanceWithLongIdentifierNames"
	env.Set(longName, "value")

	identifier := &Identifier{Name: longName}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		identifier.Eval(env)
	}
}
