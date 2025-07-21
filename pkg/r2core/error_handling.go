package r2core

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// ErrorContext contiene información de contexto para errores mejorados
type ErrorContext struct {
	Function   string
	Line       int
	Column     int
	SourceFile string
	Expected   string
	Received   string
	Operation  string
	Value      interface{}
}

// ErrorFormatter provee métodos para formatear errores consistentemente
type ErrorFormatter struct{}

// NewErrorFormatter crea una nueva instancia de ErrorFormatter
func NewErrorFormatter() *ErrorFormatter {
	return &ErrorFormatter{}
}

// FormatTypeError formats type errors with VS Code compatible format
func (ef *ErrorFormatter) FormatTypeError(ctx ErrorContext) string {
	if ctx.Line > 0 && ctx.Column > 0 {
		// VS Code format: file:line:col: error: message
		base := fmt.Sprintf("line %d, column %d: type error in %s", ctx.Line, ctx.Column, ctx.Function)
		if ctx.Expected != "" && ctx.Received != "" {
			base += fmt.Sprintf(": expected %s, got %s", ctx.Expected, ctx.Received)
		}
		if ctx.Value != nil {
			base += fmt.Sprintf(" (value: %v)", ctx.Value)
		}
		return base
	}

	// Fallback format without position
	base := fmt.Sprintf("Type error in %s", ctx.Function)
	if ctx.Expected != "" && ctx.Received != "" {
		base += fmt.Sprintf(": expected %s, got %s", ctx.Expected, ctx.Received)
	}
	if ctx.Value != nil {
		base += fmt.Sprintf(" (value: %v)", ctx.Value)
	}

	return base
}

// FormatOperationError formats operation errors with VS Code compatible format
func (ef *ErrorFormatter) FormatOperationError(ctx ErrorContext) string {
	if ctx.Line > 0 && ctx.Column > 0 {
		base := fmt.Sprintf("line %d, column %d: operation error", ctx.Line, ctx.Column)
		if ctx.Function != "" {
			base += fmt.Sprintf(" in %s", ctx.Function)
		}
		if ctx.Operation != "" {
			base += fmt.Sprintf(" (%s)", ctx.Operation)
		}
		if ctx.Expected != "" {
			base += fmt.Sprintf(": %s", ctx.Expected)
		}
		if ctx.Value != nil {
			base += fmt.Sprintf(" (value: %v)", ctx.Value)
		}
		return base
	}

	// Fallback format
	base := fmt.Sprintf("Operation error")
	if ctx.Function != "" {
		base += fmt.Sprintf(" in %s", ctx.Function)
	}
	if ctx.Operation != "" {
		base += fmt.Sprintf(" (%s)", ctx.Operation)
	}
	if ctx.Expected != "" {
		base += fmt.Sprintf(": %s", ctx.Expected)
	}
	if ctx.Value != nil {
		base += fmt.Sprintf(" (value: %v)", ctx.Value)
	}

	return base
}

// FormatArgumentError formats argument errors with VS Code compatible format
func (ef *ErrorFormatter) FormatArgumentError(ctx ErrorContext) string {
	if ctx.Line > 0 && ctx.Column > 0 {
		base := fmt.Sprintf("line %d, column %d: argument error in %s", ctx.Line, ctx.Column, ctx.Function)
		if ctx.Expected != "" {
			base += fmt.Sprintf(": %s", ctx.Expected)
		}
		if ctx.Received != "" {
			base += fmt.Sprintf(" (received: %s)", ctx.Received)
		}
		return base
	}

	// Fallback format
	base := fmt.Sprintf("Argument error in %s", ctx.Function)
	if ctx.Expected != "" {
		base += fmt.Sprintf(": %s", ctx.Expected)
	}
	if ctx.Received != "" {
		base += fmt.Sprintf(" (received: %s)", ctx.Received)
	}

	return base
}

// FormatRuntimeError formats runtime errors with VS Code compatible format
func (ef *ErrorFormatter) FormatRuntimeError(ctx ErrorContext) string {
	if ctx.Line > 0 && ctx.Column > 0 {
		base := fmt.Sprintf("line %d, column %d: runtime error", ctx.Line, ctx.Column)
		if ctx.Function != "" {
			base += fmt.Sprintf(" in %s", ctx.Function)
		}
		if ctx.Operation != "" {
			base += fmt.Sprintf(": %s", ctx.Operation)
		}
		if ctx.Value != nil {
			base += fmt.Sprintf(" (value: %v)", ctx.Value)
		}
		return base
	}

	// Fallback format
	base := fmt.Sprintf("Runtime error")
	if ctx.Function != "" {
		base += fmt.Sprintf(" in %s", ctx.Function)
	}
	if ctx.Operation != "" {
		base += fmt.Sprintf(": %s", ctx.Operation)
	}
	if ctx.Value != nil {
		base += fmt.Sprintf(" (value: %v)", ctx.Value)
	}

	return base
}

// GetCallerInfo obtiene información del contexto de llamada
func GetCallerInfo() (string, string, int) {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return "", "", 0
	}

	funcName := runtime.FuncForPC(pc).Name()

	// Extraer solo el nombre de la función sin el paquete completo
	if idx := strings.LastIndex(funcName, "."); idx >= 0 {
		funcName = funcName[idx+1:]
	}

	// Extraer solo el nombre del archivo
	if idx := strings.LastIndex(file, "/"); idx >= 0 {
		file = file[idx+1:]
	}

	return funcName, file, line
}

// PanicWithContext genera un panic con contexto mejorado
func PanicWithContext(errorType string, ctx ErrorContext) {
	ef := NewErrorFormatter()
	var message string

	switch errorType {
	case "type":
		message = ef.FormatTypeError(ctx)
	case "operation":
		message = ef.FormatOperationError(ctx)
	case "argument":
		message = ef.FormatArgumentError(ctx)
	case "runtime":
		message = ef.FormatRuntimeError(ctx)
	default:
		message = fmt.Sprintf("Error desconocido en %s: %s", ctx.Function, ctx.Operation)
	}

	panic(message)
}

// Helper functions para crear contextos rápidamente

// TypeErrorContext crea un contexto para errores de tipo
func TypeErrorContext(function, expected, received string, value interface{}) ErrorContext {
	funcName, _, _ := GetCallerInfo()
	if function == "" {
		function = funcName
	}

	return ErrorContext{
		Function: function,
		Expected: expected,
		Received: received,
		Value:    value,
	}
}

// OperationErrorContext crea un contexto para errores de operación
func OperationErrorContext(operation, description string, value interface{}) ErrorContext {
	funcName, _, _ := GetCallerInfo()

	return ErrorContext{
		Function:  funcName,
		Operation: operation,
		Expected:  description,
		Value:     value,
	}
}

// ArgumentErrorContext crea un contexto para errores de argumentos
func ArgumentErrorContext(function, expected, received string) ErrorContext {
	funcName, _, _ := GetCallerInfo()
	if function == "" {
		function = funcName
	}

	return ErrorContext{
		Function: function,
		Expected: expected,
		Received: received,
	}
}

// typeof retorna el tipo de un valor como string
func typeof(value interface{}) string {
	if value == nil {
		return "nil"
	}

	t := reflect.TypeOf(value)
	switch t.Kind() {
	case reflect.Slice:
		return "array"
	case reflect.Map:
		return "object"
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "int"
	case reflect.Float32, reflect.Float64:
		return "float64"
	case reflect.Bool:
		return "bool"
	case reflect.Func:
		return "function"
	default:
		return t.String()
	}
}
