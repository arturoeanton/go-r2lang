package r2libs

import (
	"strings"
	"testing"
)

func TestArgumentError(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			msg := r.(string)
			if !strings.Contains(msg, "Function 'sin'") {
				t.Errorf("Expected function name in error, got: %s", msg)
			}
			if !strings.Contains(msg, "expected 1 argument") {
				t.Errorf("Expected argument expectation, got: %s", msg)
			}
			if !strings.Contains(msg, "received 0 arguments") {
				t.Errorf("Expected received count, got: %s", msg)
			}
		} else {
			t.Error("Expected panic but didn't get one")
		}
	}()

	ArgumentError("sin", "1 argument (number)", 0)
}

func TestTypeArgumentError(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			msg := r.(string)
			if !strings.Contains(msg, "Function 'sqrt'") {
				t.Errorf("Expected function name in error, got: %s", msg)
			}
			if !strings.Contains(msg, "argument 1 must be number") {
				t.Errorf("Expected argument type requirement, got: %s", msg)
			}
			if !strings.Contains(msg, "got string") {
				t.Errorf("Expected received type, got: %s", msg)
			}
			if !strings.Contains(msg, "value: hello") {
				t.Errorf("Expected value, got: %s", msg)
			}
		} else {
			t.Error("Expected panic but didn't get one")
		}
	}()

	TypeArgumentError("sqrt", 0, "number", "hello")
}

func TestMathError(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			msg := r.(string)
			if !strings.Contains(msg, "Math error in function 'sqrt'") {
				t.Errorf("Expected math error for sqrt, got: %s", msg)
			}
			if !strings.Contains(msg, "cannot calculate square root") {
				t.Errorf("Expected specific math operation error, got: %s", msg)
			}
			if !strings.Contains(msg, "value: -1") {
				t.Errorf("Expected value in error, got: %s", msg)
			}
		} else {
			t.Error("Expected panic but didn't get one")
		}
	}()

	MathError("sqrt", "cannot calculate square root of negative number", -1)
}

func TestGetTypeName(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected string
	}{
		{nil, "nil"},
		{"hello", "string"},
		{42, "integer"},
		{3.14, "number"},
		{true, "boolean"},
		{[]interface{}{1, 2, 3}, "array"},
		{map[string]interface{}{"key": "value"}, "object"},
		{func() {}, "function"},
	}

	for _, test := range tests {
		result := getTypeName(test.value)
		if result != test.expected {
			t.Errorf("getTypeName(%v) = %s, expected %s", test.value, result, test.expected)
		}
	}
}
