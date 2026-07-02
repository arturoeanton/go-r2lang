package r2libs

import (
	"reflect"
	"testing"
	"time"

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

		// InterfaceSlice is the type .map()/.filter()/.sort()/.reverse()
		// return; it must be recognized as an array like a plain
		// []interface{} literal.
		if !isFunc(r2core.InterfaceSlice{1.0, 2.0}, "array").(bool) {
			t.Errorf("is: expected r2core.InterfaceSlice to be array")
		}
	})

	// deepCopy must not stack-overflow the process on a self-referential
	// array/map built via index assignment (e.g. a[0] = a); it should
	// instead produce a copy with the same cyclic shape.
	t.Run("deepCopy self-referential array", func(t *testing.T) {
		selfRefArray := []interface{}{1.0, 2.0, 3.0}
		selfRefArray[0] = selfRefArray

		copied := deepCopyFunc(selfRefArray).([]interface{})
		if copied[1].(float64) != 2.0 {
			t.Errorf("deepCopy: expected copied[1] == 2.0, got %v", copied[1])
		}
		copiedSelf, ok := copied[0].([]interface{})
		if !ok {
			t.Fatalf("deepCopy: expected copied[0] to be []interface{}, got %T", copied[0])
		}
		if copiedSelf[1].(float64) != 2.0 {
			t.Errorf("deepCopy: expected cycle to resolve back to the copy, got %v", copiedSelf[1])
		}
	})

	t.Run("deepCopy self-referential map", func(t *testing.T) {
		selfRefMap := map[string]interface{}{"a": 1.0}
		selfRefMap["self"] = selfRefMap

		copied := deepCopyFunc(selfRefMap).(map[string]interface{})
		if copied["a"].(float64) != 1.0 {
			t.Errorf("deepCopy: expected copied[\"a\"] == 1.0, got %v", copied["a"])
		}
		copiedSelf, ok := copied["self"].(map[string]interface{})
		if !ok {
			t.Fatalf("deepCopy: expected copied[\"self\"] to be map[string]interface{}, got %T", copied["self"])
		}
		if copiedSelf["a"].(float64) != 1.0 {
			t.Errorf("deepCopy: expected cycle to resolve back to the copy, got %v", copiedSelf["a"])
		}
	})

	// deepCopy must not reflect into pointer-to-struct interpreter-internal
	// types (functions, dates, durations, object instances): doing so used
	// to strip the pointer and rebuild a bare, uncallable/unusable struct
	// value that failed downstream type assertions like *r2core.UserFunction.
	t.Run("deepCopy preserves opaque pointer types", func(t *testing.T) {
		fnCode := "func identity(x) { return x }"
		fnParser := r2core.NewParser(fnCode)
		fnProgram := fnParser.ParseProgram()
		fnProgram.Eval(env)
		fnVal, ok := env.Get("identity")
		if !ok {
			t.Fatal("identity function not found in environment")
		}
		userFn := fnVal.(*r2core.UserFunction)

		copiedFn := deepCopyFunc(userFn)
		if _, ok := copiedFn.(*r2core.UserFunction); !ok {
			t.Errorf("deepCopy: expected copy of *r2core.UserFunction to remain *r2core.UserFunction, got %T", copiedFn)
		}

		date := r2core.NewDateValue(time.Now())
		copiedDate := deepCopyFunc(date)
		if _, ok := copiedDate.(*r2core.DateValue); !ok {
			t.Errorf("deepCopy: expected copy of *r2core.DateValue to remain *r2core.DateValue, got %T", copiedDate)
		}

		// Also inside a container, since that's the common real-world path.
		container := []interface{}{userFn, date}
		copiedContainer := deepCopyFunc(container).([]interface{})
		if _, ok := copiedContainer[0].(*r2core.UserFunction); !ok {
			t.Errorf("deepCopy: expected container[0] to remain *r2core.UserFunction, got %T", copiedContainer[0])
		}
		if _, ok := copiedContainer[1].(*r2core.DateValue); !ok {
			t.Errorf("deepCopy: expected container[1] to remain *r2core.DateValue, got %T", copiedContainer[1])
		}
	})

	// join must accept r2core.InterfaceSlice the same way it accepts a
	// plain []interface{} literal.
	t.Run("join accepts InterfaceSlice", func(t *testing.T) {
		joinFunc := stdModule["join"].(r2core.BuiltinFunction)
		result := joinFunc(r2core.InterfaceSlice{"a", "b", "c"}, "-")
		if result != "a-b-c" {
			t.Errorf("join: expected \"a-b-c\", got %v", result)
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
