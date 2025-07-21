package r2libs

import (
	"fmt"
	"reflect"
)

// ArgumentError generates an improved argument error
func ArgumentError(function string, expected string, received int) {
	msg := fmt.Sprintf("Function '%s': expected %s, but received %d arguments",
		function, expected, received)
	panic(msg)
}

// TypeArgumentError generates an improved type argument error
func TypeArgumentError(function string, argIndex int, expected string, received interface{}) {
	receivedType := getTypeName(received)
	msg := fmt.Sprintf("Function '%s': argument %d must be %s, but got %s (value: %v)",
		function, argIndex+1, expected, receivedType, received)
	panic(msg)
}

// MathError generates specific mathematical errors
func MathError(function string, operation string, value interface{}) {
	msg := fmt.Sprintf("Math error in function '%s': %s (value: %v)",
		function, operation, value)
	panic(msg)
}

// getTypeName returns the name of a value's type
func getTypeName(value interface{}) string {
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
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Func:
		return "function"
	default:
		return t.String()
	}
}
