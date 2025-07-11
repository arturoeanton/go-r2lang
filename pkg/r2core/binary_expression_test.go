package r2core

import (
	"testing"
)

func TestBinaryExpression_Arithmetic(t *testing.T) {
	env := NewEnvironment()

	tests := []struct {
		name     string
		left     Node
		op       string
		right    Node
		expected interface{}
	}{
		// Addition
		{
			name:     "number addition",
			left:     &NumberLiteral{Value: 5},
			op:       "+",
			right:    &NumberLiteral{Value: 3},
			expected: 8.0,
		},
		{
			name:     "string concatenation",
			left:     &StringLiteral{Value: "hello"},
			op:       "+",
			right:    &StringLiteral{Value: " world"},
			expected: "hello world",
		},
		{
			name:     "string and number concatenation",
			left:     &StringLiteral{Value: "count: "},
			op:       "+",
			right:    &NumberLiteral{Value: 42},
			expected: "count: 42",
		},
		{
			name:     "number and string concatenation",
			left:     &NumberLiteral{Value: 42},
			op:       "+",
			right:    &StringLiteral{Value: " items"},
			expected: "42 items",
		},
		// Subtraction
		{
			name:     "number subtraction",
			left:     &NumberLiteral{Value: 10},
			op:       "-",
			right:    &NumberLiteral{Value: 3},
			expected: 7.0,
		},
		{
			name:     "negative result subtraction",
			left:     &NumberLiteral{Value: 3},
			op:       "-",
			right:    &NumberLiteral{Value: 10},
			expected: -7.0,
		},
		// Multiplication
		{
			name:     "number multiplication",
			left:     &NumberLiteral{Value: 6},
			op:       "*",
			right:    &NumberLiteral{Value: 7},
			expected: 42.0,
		},
		{
			name:     "multiplication by zero",
			left:     &NumberLiteral{Value: 42},
			op:       "*",
			right:    &NumberLiteral{Value: 0},
			expected: 0.0,
		},
		// Division
		{
			name:     "number division",
			left:     &NumberLiteral{Value: 15},
			op:       "/",
			right:    &NumberLiteral{Value: 3},
			expected: 5.0,
		},
		{
			name:     "division with decimal result",
			left:     &NumberLiteral{Value: 7},
			op:       "/",
			right:    &NumberLiteral{Value: 2},
			expected: 3.5,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			be := &BinaryExpression{
				Left:  test.left,
				Op:    test.op,
				Right: test.right,
			}

			result := be.Eval(env)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestBinaryExpression_Comparison(t *testing.T) {
	env := NewEnvironment()

	tests := []struct {
		name     string
		left     Node
		op       string
		right    Node
		expected bool
	}{
		// Less than
		{
			name:     "less than true",
			left:     &NumberLiteral{Value: 3},
			op:       "<",
			right:    &NumberLiteral{Value: 5},
			expected: true,
		},
		{
			name:     "less than false",
			left:     &NumberLiteral{Value: 5},
			op:       "<",
			right:    &NumberLiteral{Value: 3},
			expected: false,
		},
		{
			name:     "less than equal",
			left:     &NumberLiteral{Value: 5},
			op:       "<",
			right:    &NumberLiteral{Value: 5},
			expected: false,
		},
		// Greater than
		{
			name:     "greater than true",
			left:     &NumberLiteral{Value: 5},
			op:       ">",
			right:    &NumberLiteral{Value: 3},
			expected: true,
		},
		{
			name:     "greater than false",
			left:     &NumberLiteral{Value: 3},
			op:       ">",
			right:    &NumberLiteral{Value: 5},
			expected: false,
		},
		// Less than or equal
		{
			name:     "less than or equal (less)",
			left:     &NumberLiteral{Value: 3},
			op:       "<=",
			right:    &NumberLiteral{Value: 5},
			expected: true,
		},
		{
			name:     "less than or equal (equal)",
			left:     &NumberLiteral{Value: 5},
			op:       "<=",
			right:    &NumberLiteral{Value: 5},
			expected: true,
		},
		{
			name:     "less than or equal (greater)",
			left:     &NumberLiteral{Value: 7},
			op:       "<=",
			right:    &NumberLiteral{Value: 5},
			expected: false,
		},
		// Greater than or equal
		{
			name:     "greater than or equal (greater)",
			left:     &NumberLiteral{Value: 7},
			op:       ">=",
			right:    &NumberLiteral{Value: 5},
			expected: true,
		},
		{
			name:     "greater than or equal (equal)",
			left:     &NumberLiteral{Value: 5},
			op:       ">=",
			right:    &NumberLiteral{Value: 5},
			expected: true,
		},
		{
			name:     "greater than or equal (less)",
			left:     &NumberLiteral{Value: 3},
			op:       ">=",
			right:    &NumberLiteral{Value: 5},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			be := &BinaryExpression{
				Left:  test.left,
				Op:    test.op,
				Right: test.right,
			}

			result := be.Eval(env)
			if result != test.expected {
				t.Errorf("Expected %t, got %v", test.expected, result)
			}
		})
	}
}

func TestBinaryExpression_Equality(t *testing.T) {
	env := NewEnvironment()

	tests := []struct {
		name     string
		left     Node
		op       string
		right    Node
		expected bool
	}{
		// Equal
		{
			name:     "numbers equal",
			left:     &NumberLiteral{Value: 5},
			op:       "==",
			right:    &NumberLiteral{Value: 5},
			expected: true,
		},
		{
			name:     "numbers not equal",
			left:     &NumberLiteral{Value: 5},
			op:       "==",
			right:    &NumberLiteral{Value: 3},
			expected: false,
		},
		{
			name:     "strings equal",
			left:     &StringLiteral{Value: "hello"},
			op:       "==",
			right:    &StringLiteral{Value: "hello"},
			expected: true,
		},
		{
			name:     "strings not equal",
			left:     &StringLiteral{Value: "hello"},
			op:       "==",
			right:    &StringLiteral{Value: "world"},
			expected: false,
		},
		{
			name:     "booleans equal",
			left:     &BooleanLiteral{Value: true},
			op:       "==",
			right:    &BooleanLiteral{Value: true},
			expected: true,
		},
		{
			name:     "booleans not equal",
			left:     &BooleanLiteral{Value: true},
			op:       "==",
			right:    &BooleanLiteral{Value: false},
			expected: false,
		},
		// Not equal
		{
			name:     "numbers not equal operator",
			left:     &NumberLiteral{Value: 5},
			op:       "!=",
			right:    &NumberLiteral{Value: 3},
			expected: true,
		},
		{
			name:     "numbers equal with not equal operator",
			left:     &NumberLiteral{Value: 5},
			op:       "!=",
			right:    &NumberLiteral{Value: 5},
			expected: false,
		},
		{
			name:     "strings not equal operator",
			left:     &StringLiteral{Value: "hello"},
			op:       "!=",
			right:    &StringLiteral{Value: "world"},
			expected: true,
		},
		{
			name:     "strings equal with not equal operator",
			left:     &StringLiteral{Value: "hello"},
			op:       "!=",
			right:    &StringLiteral{Value: "hello"},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			be := &BinaryExpression{
				Left:  test.left,
				Op:    test.op,
				Right: test.right,
			}

			result := be.Eval(env)
			if result != test.expected {
				t.Errorf("Expected %t, got %v", test.expected, result)
			}
		})
	}
}

func TestBinaryExpression_ArrayConcatenation(t *testing.T) {
	env := NewEnvironment()

	// Test array + array concatenation
	leftArray := &ArrayLiteral{Elements: []Node{
		&NumberLiteral{Value: 1},
		&NumberLiteral{Value: 2},
	}}
	rightArray := &ArrayLiteral{Elements: []Node{
		&NumberLiteral{Value: 3},
		&NumberLiteral{Value: 4},
	}}

	be := &BinaryExpression{
		Left:  leftArray,
		Op:    "+",
		Right: rightArray,
	}

	result := be.Eval(env)
	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("Expected []interface{}, got %T", result)
	}

	expected := []interface{}{1.0, 2.0, 3.0, 4.0}
	if len(slice) != len(expected) {
		t.Fatalf("Expected length %d, got %d", len(expected), len(slice))
	}

	for i, expectedVal := range expected {
		if slice[i] != expectedVal {
			t.Errorf("Expected slice[%d] = %v, got %v", i, expectedVal, slice[i])
		}
	}
}

func TestBinaryExpression_ArrayElementConcatenation(t *testing.T) {
	env := NewEnvironment()

	// Test array + element
	leftArray := &ArrayLiteral{Elements: []Node{
		&NumberLiteral{Value: 1},
		&NumberLiteral{Value: 2},
	}}

	be := &BinaryExpression{
		Left:  leftArray,
		Op:    "+",
		Right: &NumberLiteral{Value: 3},
	}

	result := be.Eval(env)
	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("Expected []interface{}, got %T", result)
	}

	expected := []interface{}{1.0, 2.0, 3.0}
	if len(slice) != len(expected) {
		t.Fatalf("Expected length %d, got %d", len(expected), len(slice))
	}

	for i, expectedVal := range expected {
		if slice[i] != expectedVal {
			t.Errorf("Expected slice[%d] = %v, got %v", i, expectedVal, slice[i])
		}
	}
}

func TestBinaryExpression_DivisionByZero(t *testing.T) {
	env := NewEnvironment()

	be := &BinaryExpression{
		Left:  &NumberLiteral{Value: 10},
		Op:    "/",
		Right: &NumberLiteral{Value: 0},
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for division by zero")
		} else if r != "Division by zero" {
			t.Errorf("Expected 'Division by zero' panic, got %v", r)
		}
	}()

	be.Eval(env)
}

func TestBinaryExpression_UnsupportedOperator(t *testing.T) {
	env := NewEnvironment()

	be := &BinaryExpression{
		Left:  &NumberLiteral{Value: 5},
		Op:    "**", // Unsupported operator
		Right: &NumberLiteral{Value: 2},
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for unsupported operator")
		} else {
			expectedMsg := "Unsupported binary operator: **"
			if r != expectedMsg {
				t.Errorf("Expected %q panic, got %v", expectedMsg, r)
			}
		}
	}()

	be.Eval(env)
}

func TestBinaryExpression_NestedExpressions(t *testing.T) {
	env := NewEnvironment()

	// Test (2 + 3) * 4 = 20
	innerExpr := &BinaryExpression{
		Left:  &NumberLiteral{Value: 2},
		Op:    "+",
		Right: &NumberLiteral{Value: 3},
	}

	outerExpr := &BinaryExpression{
		Left:  innerExpr,
		Op:    "*",
		Right: &NumberLiteral{Value: 4},
	}

	result := outerExpr.Eval(env)
	expected := 20.0
	if result != expected {
		t.Errorf("Expected %f, got %v", expected, result)
	}
}

func TestBinaryExpression_ComplexNestedExpressions(t *testing.T) {
	env := NewEnvironment()

	// Test ((10 - 5) * 2) + (8 / 4) = 12
	leftSide := &BinaryExpression{
		Left: &BinaryExpression{
			Left:  &NumberLiteral{Value: 10},
			Op:    "-",
			Right: &NumberLiteral{Value: 5},
		},
		Op:    "*",
		Right: &NumberLiteral{Value: 2},
	}

	rightSide := &BinaryExpression{
		Left:  &NumberLiteral{Value: 8},
		Op:    "/",
		Right: &NumberLiteral{Value: 4},
	}

	finalExpr := &BinaryExpression{
		Left:  leftSide,
		Op:    "+",
		Right: rightSide,
	}

	result := finalExpr.Eval(env)
	expected := 12.0
	if result != expected {
		t.Errorf("Expected %f, got %v", expected, result)
	}
}

// Benchmark tests
func BenchmarkBinaryExpression_Addition(b *testing.B) {
	env := NewEnvironment()
	be := &BinaryExpression{
		Left:  &NumberLiteral{Value: 5},
		Op:    "+",
		Right: &NumberLiteral{Value: 3},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		be.Eval(env)
	}
}

func BenchmarkBinaryExpression_StringConcatenation(b *testing.B) {
	env := NewEnvironment()
	be := &BinaryExpression{
		Left:  &StringLiteral{Value: "hello"},
		Op:    "+",
		Right: &StringLiteral{Value: " world"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		be.Eval(env)
	}
}

func BenchmarkBinaryExpression_Comparison(b *testing.B) {
	env := NewEnvironment()
	be := &BinaryExpression{
		Left:  &NumberLiteral{Value: 5},
		Op:    "<",
		Right: &NumberLiteral{Value: 10},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		be.Eval(env)
	}
}

func BenchmarkBinaryExpression_NestedArithmetic(b *testing.B) {
	env := NewEnvironment()

	// (10 + 5) * (20 - 5)
	leftExpr := &BinaryExpression{
		Left:  &NumberLiteral{Value: 10},
		Op:    "+",
		Right: &NumberLiteral{Value: 5},
	}

	rightExpr := &BinaryExpression{
		Left:  &NumberLiteral{Value: 20},
		Op:    "-",
		Right: &NumberLiteral{Value: 5},
	}

	be := &BinaryExpression{
		Left:  leftExpr,
		Op:    "*",
		Right: rightExpr,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		be.Eval(env)
	}
}
