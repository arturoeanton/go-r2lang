package r2libs

import (
	"reflect"
	"testing"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func TestStdFunctions(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterStd(env)

	stdModuleObj, ok := env.Get("std")
	if !ok {
		t.Fatal("std module not found")
	}
	stdModule := stdModuleObj.(map[string]interface{})

	deepCopyFunc := stdModule["deepCopy"].(r2core.BuiltinFunction)
	isFunc := stdModule["is"].(r2core.BuiltinFunction)

	// Test deepCopy
	t.Run("deepCopy", func(t *testing.T) {
		originalArray := []interface{}{1.0, "hello", []interface{}{true, false}}
		copiedArray := deepCopyFunc(originalArray).([]interface{})

		if !reflect.DeepEqual(originalArray, copiedArray) {
			t.Errorf("deepCopy: copied array is not deep equal to original")
		}

		// Modify copied array and check original
		copiedArray[0] = 99.0
		copiedArray[2].([]interface{})[0] = false

		if reflect.DeepEqual(originalArray, copiedArray) {
			t.Errorf("deepCopy: copied array is still deep equal after modification")
		}
		if originalArray[0].(float64) != 1.0 {
			t.Errorf("deepCopy: original array modified unexpectedly")
		}
		if originalArray[2].([]interface{})[0].(bool) != true {
			t.Errorf("deepCopy: original nested array modified unexpectedly")
		}
	})

	// Test is
	t.Run("is", func(t *testing.T) {
		if !isFunc(10.0, "float64").(bool) {
			t.Errorf("is: expected 10.0 to be float64")
		}
		if isFunc(10.0, "string").(bool) {
			t.Errorf("is: expected 10.0 not to be string")
		}
		if !isFunc("hello", "string").(bool) {
			t.Errorf("is: expected \"hello\" to be string")
		}
		if !isFunc(true, "bool").(bool) {
			t.Errorf("is: expected true to be bool")
		}
		if !isFunc([]interface{}{}, "array").(bool) {
			t.Errorf("is: expected []interface{} to be array")
		}
		if !isFunc(map[string]interface{}{}, "map").(bool) {
			t.Errorf("is: expected map[string]interface{} to be map")
		}
		if !isFunc(nil, "nil").(bool) {
			t.Errorf("is: expected nil to be nil")
		}
		// Test with r2core.UserFunction
		identityFuncCode := "func identity(x) { return x }"
		identityParser := r2core.NewParser(identityFuncCode)
		identityProgram := identityParser.ParseProgram()
		identityProgram.Eval(env) // Evaluate to register the function
		var identityUserFunc *r2core.UserFunction
		{
			val, ok := env.Get("identity")
			if !ok {
				t.Fatal("identity function not found in environment")
			}
			identityUserFunc = val.(*r2core.UserFunction)
		}

		if !isFunc(identityUserFunc, "function").(bool) {
			t.Errorf("is: expected UserFunction to be function")
		}
	})

	// Test contains
	t.Run("contains", func(t *testing.T) {
		containsFunc := stdModule["contains"].(r2core.BuiltinFunction)

		// Test basic functionality
		if !containsFunc("hello world", "world").(bool) {
			t.Errorf("contains: expected 'hello world' to contain 'world'")
		}
		if containsFunc("hello world", "foo").(bool) {
			t.Errorf("contains: expected 'hello world' not to contain 'foo'")
		}
		if !containsFunc("test=value", "=").(bool) {
			t.Errorf("contains: expected 'test=value' to contain '='")
		}

		// Test edge cases
		if !containsFunc("", "").(bool) {
			t.Errorf("contains: expected empty string to contain empty string")
		}
		if !containsFunc("abc", "").(bool) {
			t.Errorf("contains: expected 'abc' to contain empty string")
		}
		if containsFunc("", "a").(bool) {
			t.Errorf("contains: expected empty string not to contain 'a'")
		}
	})
}
