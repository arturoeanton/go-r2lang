package assertions

import (
	"fmt"
	"reflect"
	"strings"
)

// Assert provides basic assertion functionality
type Assert struct {
	testName string
}

// NewAssert creates a new assertion context
func NewAssert(testName string) *Assert {
	return &Assert{
		testName: testName,
	}
}

// AssertionError represents an assertion failure
type AssertionError struct {
	Message  string
	Expected interface{}
	Actual   interface{}
	TestName string
}

func (ae *AssertionError) Error() string {
	return ae.Message
}

// fail creates an assertion error and panics with it
func (a *Assert) fail(message string, expected, actual interface{}) {
	err := &AssertionError{
		Message:  message,
		Expected: expected,
		Actual:   actual,
		TestName: a.testName,
	}
	panic(err)
}

// Equals asserts that two values are equal
func (a *Assert) Equals(actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		a.fail(
			fmt.Sprintf("Expected %v to equal %v", actual, expected),
			expected,
			actual,
		)
	}
}

// NotEquals asserts that two values are not equal
func (a *Assert) NotEquals(actual, expected interface{}) {
	if reflect.DeepEqual(actual, expected) {
		a.fail(
			fmt.Sprintf("Expected %v to not equal %v", actual, expected),
			expected,
			actual,
		)
	}
}

// True asserts that a value is true
func (a *Assert) True(actual interface{}) {
	if !isTruthy(actual) {
		a.fail(
			fmt.Sprintf("Expected %v to be true", actual),
			true,
			actual,
		)
	}
}

// False asserts that a value is false
func (a *Assert) False(actual interface{}) {
	if isTruthy(actual) {
		a.fail(
			fmt.Sprintf("Expected %v to be false", actual),
			false,
			actual,
		)
	}
}

// Nil asserts that a value is nil
func (a *Assert) Nil(actual interface{}) {
	if actual != nil {
		a.fail(
			fmt.Sprintf("Expected %v to be nil", actual),
			nil,
			actual,
		)
	}
}

// NotNil asserts that a value is not nil
func (a *Assert) NotNil(actual interface{}) {
	if actual == nil {
		a.fail(
			fmt.Sprintf("Expected value to not be nil"),
			"not nil",
			nil,
		)
	}
}

// Contains asserts that a string contains a substring
func (a *Assert) Contains(haystack, needle string) {
	if !strings.Contains(haystack, needle) {
		a.fail(
			fmt.Sprintf("Expected '%s' to contain '%s'", haystack, needle),
			needle,
			haystack,
		)
	}
}

// NotContains asserts that a string does not contain a substring
func (a *Assert) NotContains(haystack, needle string) {
	if strings.Contains(haystack, needle) {
		a.fail(
			fmt.Sprintf("Expected '%s' to not contain '%s'", haystack, needle),
			fmt.Sprintf("not containing '%s'", needle),
			haystack,
		)
	}
}

// Greater asserts that a value is greater than another
func (a *Assert) Greater(actual, expected interface{}) {
	if !isGreater(actual, expected) {
		a.fail(
			fmt.Sprintf("Expected %v to be greater than %v", actual, expected),
			fmt.Sprintf("> %v", expected),
			actual,
		)
	}
}

// GreaterOrEqual asserts that a value is greater than or equal to another
func (a *Assert) GreaterOrEqual(actual, expected interface{}) {
	if !isGreaterOrEqual(actual, expected) {
		a.fail(
			fmt.Sprintf("Expected %v to be greater than or equal to %v", actual, expected),
			fmt.Sprintf(">= %v", expected),
			actual,
		)
	}
}

// Less asserts that a value is less than another
func (a *Assert) Less(actual, expected interface{}) {
	if !isLess(actual, expected) {
		a.fail(
			fmt.Sprintf("Expected %v to be less than %v", actual, expected),
			fmt.Sprintf("< %v", expected),
			actual,
		)
	}
}

// LessOrEqual asserts that a value is less than or equal to another
func (a *Assert) LessOrEqual(actual, expected interface{}) {
	if !isLessOrEqual(actual, expected) {
		a.fail(
			fmt.Sprintf("Expected %v to be less than or equal to %v", actual, expected),
			fmt.Sprintf("<= %v", expected),
			actual,
		)
	}
}

// HasLength asserts that a collection has a specific length
func (a *Assert) HasLength(collection interface{}, expectedLength int) {
	length := getLength(collection)
	if length != expectedLength {
		a.fail(
			fmt.Sprintf("Expected collection to have length %d, but got %d", expectedLength, length),
			expectedLength,
			length,
		)
	}
}

// Empty asserts that a collection is empty
func (a *Assert) Empty(collection interface{}) {
	if getLength(collection) != 0 {
		a.fail(
			fmt.Sprintf("Expected collection to be empty, but got length %d", getLength(collection)),
			0,
			getLength(collection),
		)
	}
}

// NotEmpty asserts that a collection is not empty
func (a *Assert) NotEmpty(collection interface{}) {
	if getLength(collection) == 0 {
		a.fail(
			fmt.Sprintf("Expected collection to not be empty"),
			"> 0",
			0,
		)
	}
}

// Panics asserts that a function panics
func (a *Assert) Panics(fn func()) {
	defer func() {
		if r := recover(); r == nil {
			a.fail(
				"Expected function to panic, but it didn't",
				"panic",
				"no panic",
			)
		}
	}()
	fn()
}

// NotPanics asserts that a function does not panic
func (a *Assert) NotPanics(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			a.fail(
				fmt.Sprintf("Expected function to not panic, but it panicked with: %v", r),
				"no panic",
				fmt.Sprintf("panic: %v", r),
			)
		}
	}()
	fn()
}

// Helper functions

// isTruthy determines if a value is considered true in R2Lang context
func isTruthy(value interface{}) bool {
	if value == nil {
		return false
	}

	switch v := value.(type) {
	case bool:
		return v
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(v).Int() != 0
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(v).Uint() != 0
	case float32, float64:
		return reflect.ValueOf(v).Float() != 0.0
	case string:
		return v != ""
	case []interface{}:
		return len(v) > 0
	case map[string]interface{}:
		return len(v) > 0
	default:
		// Use reflection for other types
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array, reflect.Map:
			return rv.Len() > 0
		case reflect.Ptr, reflect.Interface:
			return !rv.IsNil()
		default:
			return true // Non-zero values are truthy
		}
	}
}

// getLength returns the length of a collection
func getLength(collection interface{}) int {
	if collection == nil {
		return 0
	}

	switch v := collection.(type) {
	case string:
		return len(v)
	case []interface{}:
		return len(v)
	case map[string]interface{}:
		return len(v)
	default:
		// Use reflection for other types
		rv := reflect.ValueOf(collection)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array, reflect.Map, reflect.String:
			return rv.Len()
		default:
			return 0
		}
	}
}

// Comparison helper functions
func isGreater(a, b interface{}) bool {
	return compareNumbers(a, b) > 0
}

func isGreaterOrEqual(a, b interface{}) bool {
	return compareNumbers(a, b) >= 0
}

func isLess(a, b interface{}) bool {
	return compareNumbers(a, b) < 0
}

func isLessOrEqual(a, b interface{}) bool {
	return compareNumbers(a, b) <= 0
}

// compareNumbers compares two numeric values, returning -1, 0, or 1
func compareNumbers(a, b interface{}) int {
	aFloat := toFloat64(a)
	bFloat := toFloat64(b)

	if aFloat < bFloat {
		return -1
	} else if aFloat > bFloat {
		return 1
	}
	return 0
}

// toFloat64 converts a numeric value to float64
func toFloat64(value interface{}) float64 {
	switch v := value.(type) {
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case float32:
		return float64(v)
	case float64:
		return v
	default:
		// Try to use reflection
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return float64(rv.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return float64(rv.Uint())
		case reflect.Float32, reflect.Float64:
			return rv.Float()
		default:
			return 0.0
		}
	}
}
