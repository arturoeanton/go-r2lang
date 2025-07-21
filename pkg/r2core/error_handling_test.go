package r2core

import (
	"strings"
	"testing"
)

func TestErrorFormatter(t *testing.T) {
	ef := NewErrorFormatter()

	t.Run("FormatTypeError", func(t *testing.T) {
		ctx := ErrorContext{
			Function: "testFunction",
			Line:     10,
			Column:   5,
			Expected: "number",
			Received: "string",
			Value:    "hello",
		}

		result := ef.FormatTypeError(ctx)

		if !strings.Contains(result, "type error in testFunction") {
			t.Errorf("Expected type error message, got: %s", result)
		}
		if !strings.Contains(result, "line 10, column 5") {
			t.Errorf("Expected line and column info, got: %s", result)
		}
		if !strings.Contains(result, "expected number, got string") {
			t.Errorf("Expected expected/received info, got: %s", result)
		}
		if !strings.Contains(result, "value: hello") {
			t.Errorf("Expected value info, got: %s", result)
		}
	})

	t.Run("FormatOperationError", func(t *testing.T) {
		ctx := ErrorContext{
			Function:  "binaryOperation",
			Operation: "division",
			Expected:  "division by zero not allowed",
			Value:     "10 / 0",
		}

		result := ef.FormatOperationError(ctx)

		if !strings.Contains(result, "Operation error in binaryOperation") {
			t.Errorf("Expected operation error message, got: %s", result)
		}
		if !strings.Contains(result, "(division)") {
			t.Errorf("Expected operation name, got: %s", result)
		}
		if !strings.Contains(result, "division by zero not allowed") {
			t.Errorf("Expected operation description, got: %s", result)
		}
	})

	t.Run("FormatArgumentError", func(t *testing.T) {
		ctx := ErrorContext{
			Function: "sin",
			Expected: "1 argument (number)",
			Received: "0 arguments",
		}

		result := ef.FormatArgumentError(ctx)

		if !strings.Contains(result, "Argument error in sin") {
			t.Errorf("Expected argument error message, got: %s", result)
		}
		if !strings.Contains(result, "1 argument (number)") {
			t.Errorf("Expected argument requirements, got: %s", result)
		}
	})

	t.Run("FormatRuntimeError", func(t *testing.T) {
		ctx := ErrorContext{
			Function:  "pipeline",
			Operation: "function not found",
			Value:     "unknownFunction",
		}

		result := ef.FormatRuntimeError(ctx)

		if !strings.Contains(result, "Runtime error in pipeline") {
			t.Errorf("Expected runtime error message, got: %s", result)
		}
		if !strings.Contains(result, "function not found") {
			t.Errorf("Expected operation description, got: %s", result)
		}
	})
}

func TestTypeof(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected string
	}{
		{nil, "nil"},
		{"hello", "string"},
		{42, "int"},
		{3.14, "float64"},
		{true, "bool"},
		{[]interface{}{1, 2, 3}, "array"},
		{map[string]interface{}{"key": "value"}, "object"},
		{func() {}, "function"},
	}

	for _, test := range tests {
		result := typeof(test.value)
		if result != test.expected {
			t.Errorf("typeof(%v) = %s, expected %s", test.value, result, test.expected)
		}
	}
}

func TestErrorContextHelpers(t *testing.T) {
	t.Run("TypeErrorContext", func(t *testing.T) {
		ctx := TypeErrorContext("testFunc", "number", "string", "hello")

		if ctx.Function != "testFunc" {
			t.Errorf("Expected function name 'testFunc', got '%s'", ctx.Function)
		}
		if ctx.Expected != "number" {
			t.Errorf("Expected 'number', got '%s'", ctx.Expected)
		}
		if ctx.Received != "string" {
			t.Errorf("Expected 'string', got '%s'", ctx.Received)
		}
		if ctx.Value != "hello" {
			t.Errorf("Expected 'hello', got %v", ctx.Value)
		}
	})

	t.Run("OperationErrorContext", func(t *testing.T) {
		ctx := OperationErrorContext("division", "division by zero", "10/0")

		if ctx.Operation != "division" {
			t.Errorf("Expected operation 'division', got '%s'", ctx.Operation)
		}
		if ctx.Expected != "division by zero" {
			t.Errorf("Expected 'division by zero', got '%s'", ctx.Expected)
		}
		if ctx.Value != "10/0" {
			t.Errorf("Expected '10/0', got %v", ctx.Value)
		}
	})

	t.Run("ArgumentErrorContext", func(t *testing.T) {
		ctx := ArgumentErrorContext("sin", "1 argument", "0 arguments")

		if ctx.Function != "sin" {
			t.Errorf("Expected function 'sin', got '%s'", ctx.Function)
		}
		if ctx.Expected != "1 argument" {
			t.Errorf("Expected '1 argument', got '%s'", ctx.Expected)
		}
		if ctx.Received != "0 arguments" {
			t.Errorf("Expected '0 arguments', got '%s'", ctx.Received)
		}
	})
}

// Test panic behavior (these should panic with improved messages)
func TestImprovedPanicMessages(t *testing.T) {
	t.Run("PanicWithContext_Type", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				msg := r.(string)
				if !strings.Contains(msg, "Type error") {
					t.Errorf("Expected type error message, got: %s", msg)
				}
			} else {
				t.Error("Expected panic but didn't get one")
			}
		}()

		ctx := TypeErrorContext("test", "number", "string", "hello")
		PanicWithContext("type", ctx)
	})

	t.Run("PanicWithContext_Operation", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				msg := r.(string)
				if !strings.Contains(msg, "Operation error") {
					t.Errorf("Expected operation error message, got: %s", msg)
				}
			} else {
				t.Error("Expected panic but didn't get one")
			}
		}()

		ctx := OperationErrorContext("division", "division by zero", "10/0")
		PanicWithContext("operation", ctx)
	})
}
