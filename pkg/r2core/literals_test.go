package r2core

import (
	"testing"
)

func TestNumberLiteral_Eval(t *testing.T) {
	env := NewEnvironment()

	tests := []struct {
		name     string
		value    float64
		expected float64
	}{
		{"positive integer", 42, 42},
		{"negative integer", -42, -42},
		{"positive float", 3.14, 3.14},
		{"negative float", -3.14, -3.14},
		{"zero", 0, 0},
		{"large number", 1e6, 1e6},
		{"small number", 1e-6, 1e-6},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nl := &NumberLiteral{Value: test.value}
			result := nl.Eval(env)

			if result != test.expected {
				t.Errorf("Expected %f, got %v", test.expected, result)
			}

			// Verify type is preserved
			if _, ok := result.(float64); !ok {
				t.Errorf("Expected result to be float64, got %T", result)
			}
		})
	}
}

func TestStringLiteral_Eval(t *testing.T) {
	env := NewEnvironment()

	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"simple string", "hello", "hello"},
		{"empty string", "", ""},
		{"string with spaces", "hello world", "hello world"},
		{"string with special chars", "hello\nworld\t!", "hello\nworld\t!"},
		{"unicode string", "héllo wørld 🌍", "héllo wørld 🌍"},
		{"long string", "Lorem ipsum dolor sit amet, consectetur adipiscing elit.", "Lorem ipsum dolor sit amet, consectetur adipiscing elit."},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sl := &StringLiteral{Value: test.value}
			result := sl.Eval(env)

			if result != test.expected {
				t.Errorf("Expected %q, got %v", test.expected, result)
			}

			// Verify type is preserved
			if _, ok := result.(string); !ok {
				t.Errorf("Expected result to be string, got %T", result)
			}
		})
	}
}

func TestBooleanLiteral_Eval(t *testing.T) {
	env := NewEnvironment()

	tests := []struct {
		name     string
		value    bool
		expected bool
	}{
		{"true value", true, true},
		{"false value", false, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bl := &BooleanLiteral{Value: test.value}
			result := bl.Eval(env)

			if result != test.expected {
				t.Errorf("Expected %t, got %v", test.expected, result)
			}

			// Verify type is preserved
			if _, ok := result.(bool); !ok {
				t.Errorf("Expected result to be bool, got %T", result)
			}
		})
	}
}

func TestArrayLiteral_Eval_Empty(t *testing.T) {
	env := NewEnvironment()
	al := &ArrayLiteral{Elements: []Node{}}

	result := al.Eval(env)

	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("Expected []interface{}, got %T", result)
	}

	if len(slice) != 0 {
		t.Errorf("Expected empty array, got length %d", len(slice))
	}
}

func TestArrayLiteral_Eval_WithElements(t *testing.T) {
	env := NewEnvironment()

	// Create array [42, "hello", true]
	elements := []Node{
		&NumberLiteral{Value: 42},
		&StringLiteral{Value: "hello"},
		&BooleanLiteral{Value: true},
	}

	al := &ArrayLiteral{Elements: elements}
	result := al.Eval(env)

	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("Expected []interface{}, got %T", result)
	}

	if len(slice) != 3 {
		t.Fatalf("Expected array length 3, got %d", len(slice))
	}

	// Check first element (number)
	if slice[0] != 42.0 {
		t.Errorf("Expected first element to be 42.0, got %v", slice[0])
	}

	// Check second element (string)
	if slice[1] != "hello" {
		t.Errorf("Expected second element to be 'hello', got %v", slice[1])
	}

	// Check third element (boolean)
	if slice[2] != true {
		t.Errorf("Expected third element to be true, got %v", slice[2])
	}
}

func TestArrayLiteral_Eval_NestedArrays(t *testing.T) {
	env := NewEnvironment()

	// Create nested array [[1, 2], [3, 4]]
	innerArray1 := &ArrayLiteral{Elements: []Node{
		&NumberLiteral{Value: 1},
		&NumberLiteral{Value: 2},
	}}

	innerArray2 := &ArrayLiteral{Elements: []Node{
		&NumberLiteral{Value: 3},
		&NumberLiteral{Value: 4},
	}}

	outerArray := &ArrayLiteral{Elements: []Node{
		innerArray1,
		innerArray2,
	}}

	result := outerArray.Eval(env)

	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("Expected []interface{}, got %T", result)
	}

	if len(slice) != 2 {
		t.Fatalf("Expected outer array length 2, got %d", len(slice))
	}

	// Check first inner array
	inner1, ok := slice[0].([]interface{})
	if !ok {
		t.Fatalf("Expected first element to be []interface{}, got %T", slice[0])
	}
	if len(inner1) != 2 {
		t.Errorf("Expected first inner array length 2, got %d", len(inner1))
	}
	if inner1[0] != 1.0 || inner1[1] != 2.0 {
		t.Errorf("Expected first inner array [1, 2], got [%v, %v]", inner1[0], inner1[1])
	}

	// Check second inner array
	inner2, ok := slice[1].([]interface{})
	if !ok {
		t.Fatalf("Expected second element to be []interface{}, got %T", slice[1])
	}
	if len(inner2) != 2 {
		t.Errorf("Expected second inner array length 2, got %d", len(inner2))
	}
	if inner2[0] != 3.0 || inner2[1] != 4.0 {
		t.Errorf("Expected second inner array [3, 4], got [%v, %v]", inner2[0], inner2[1])
	}
}

func TestFunctionLiteral_Eval(t *testing.T) {
	env := NewEnvironment()

	// Create a simple function literal: func(x, y) { return x + y }
	args := []string{"x", "y"}
	body := &BlockStatement{Statements: []Node{
		// We'll just test that the function structure is created correctly
		// The actual body evaluation would require more complex setup
	}}

	fl := &FunctionLiteral{
		Args: args,
		Body: body,
	}

	result := fl.Eval(env)

	userFunc, ok := result.(*UserFunction)
	if !ok {
		t.Fatalf("Expected *UserFunction, got %T", result)
	}

	// Check arguments
	if len(userFunc.Args) != 2 {
		t.Errorf("Expected 2 arguments, got %d", len(userFunc.Args))
	}
	if userFunc.Args[0] != "x" {
		t.Errorf("Expected first argument 'x', got %q", userFunc.Args[0])
	}
	if userFunc.Args[1] != "y" {
		t.Errorf("Expected second argument 'y', got %q", userFunc.Args[1])
	}

	// Check body
	if userFunc.Body != body {
		t.Error("Expected function body to match original")
	}

	// Check environment (closure)
	if userFunc.Env != env {
		t.Error("Expected function environment to be the evaluation environment")
	}

	// Check method flag
	if userFunc.IsMethod {
		t.Error("Expected IsMethod to be false for function literal")
	}
}

func TestFunctionLiteral_Eval_NoArgs(t *testing.T) {
	env := NewEnvironment()

	// Create a function with no arguments: func() { }
	fl := &FunctionLiteral{
		Args: []string{},
		Body: &BlockStatement{Statements: []Node{}},
	}

	result := fl.Eval(env)

	userFunc, ok := result.(*UserFunction)
	if !ok {
		t.Fatalf("Expected *UserFunction, got %T", result)
	}

	if len(userFunc.Args) != 0 {
		t.Errorf("Expected 0 arguments, got %d", len(userFunc.Args))
	}
}

func TestFunctionLiteral_Eval_ClosureCapture(t *testing.T) {
	outerEnv := NewEnvironment()
	outerEnv.Set("capturedVar", "captured")

	innerEnv := NewInnerEnv(outerEnv)

	fl := &FunctionLiteral{
		Args: []string{"param"},
		Body: &BlockStatement{Statements: []Node{}},
	}

	result := fl.Eval(innerEnv)

	userFunc, ok := result.(*UserFunction)
	if !ok {
		t.Fatalf("Expected *UserFunction, got %T", result)
	}

	// The function should capture the environment it was created in
	if userFunc.Env != innerEnv {
		t.Error("Expected function to capture the inner environment")
	}

	// Verify that the captured environment has access to outer variables
	val, exists := userFunc.Env.Get("capturedVar")
	if !exists {
		t.Error("Expected function environment to have access to captured variables")
	}
	if val != "captured" {
		t.Errorf("Expected captured value 'captured', got %v", val)
	}
}

// Benchmark tests
func BenchmarkNumberLiteral_Eval(b *testing.B) {
	env := NewEnvironment()
	nl := &NumberLiteral{Value: 42.0}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nl.Eval(env)
	}
}

func BenchmarkStringLiteral_Eval(b *testing.B) {
	env := NewEnvironment()
	sl := &StringLiteral{Value: "hello world"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.Eval(env)
	}
}

func BenchmarkArrayLiteral_Eval_Small(b *testing.B) {
	env := NewEnvironment()
	al := &ArrayLiteral{Elements: []Node{
		&NumberLiteral{Value: 1},
		&NumberLiteral{Value: 2},
		&NumberLiteral{Value: 3},
	}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		al.Eval(env)
	}
}

func BenchmarkArrayLiteral_Eval_Large(b *testing.B) {
	env := NewEnvironment()
	elements := make([]Node, 100)
	for i := 0; i < 100; i++ {
		elements[i] = &NumberLiteral{Value: float64(i)}
	}
	al := &ArrayLiteral{Elements: elements}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		al.Eval(env)
	}
}

func BenchmarkFunctionLiteral_Eval(b *testing.B) {
	env := NewEnvironment()
	fl := &FunctionLiteral{
		Args: []string{"x", "y"},
		Body: &BlockStatement{Statements: []Node{}},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fl.Eval(env)
	}
}
