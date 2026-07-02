package r2libs

import (
	"reflect"
	"strings"
	"testing"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// testGoStruct is a fixture struct used to exercise goSetField/goGetField/
// goCallMethod through reflection, mirroring how a real Go embedder would
// register a struct with goStructRegistry.
type testGoStruct struct {
	Count int
	Name  string
	Ptr   *int
}

func (t *testGoStruct) Greet(suffix string) string {
	return t.Name + suffix
}

func setupGoInterOpEnv(t *testing.T) *r2core.Environment {
	t.Helper()
	env := r2core.NewEnvironment()
	RegisterGoInterOp(env)
	goStructRegistry["testGoStruct"] = func() interface{} { return &testGoStruct{} }
	goFuncRegistry["testAdd"] = reflect.ValueOf(func(a, b int) int { return a + b })
	goFuncRegistry["testNilable"] = reflect.ValueOf(func(s *int) string {
		if s == nil {
			return "nil"
		}
		return "not-nil"
	})
	return env
}

func getBuiltin(t *testing.T, env *r2core.Environment, name string) r2core.BuiltinFunction {
	t.Helper()
	moduleVal, ok := env.Get("native")
	if !ok {
		t.Fatalf("module \"native\" not registered")
	}
	module, ok := moduleVal.(map[string]interface{})
	if !ok {
		t.Fatalf("\"native\" is not a module map")
	}
	v, ok := module[name]
	if !ok {
		t.Fatalf("builtin %q not registered under native", name)
	}
	fn, ok := v.(r2core.BuiltinFunction)
	if !ok {
		t.Fatalf("native.%q is not a BuiltinFunction", name)
	}
	return fn
}

// TestGoSetField_FloatToIntConversion confirms that setting a Go int field
// with an R2Lang float64 (R2Lang numbers are ALWAYS float64) no longer
// panics with an opaque low-level reflect message
// ("reflect.Set: value of type float64 is not assignable to type int") and
// instead performs the natural numeric conversion, matching real usage from
// R2Lang scripts where every number literal is a float64.
func TestGoSetField_FloatToIntConversion(t *testing.T) {
	env := setupGoInterOpEnv(t)
	nativeNew := getBuiltin(t, env, "new")
	nativeSetField := getBuiltin(t, env, "setField")
	nativeGetField := getBuiltin(t, env, "getField")

	obj := nativeNew("testGoStruct")
	nativeSetField(obj, "Count", float64(42))

	got := nativeGetField(obj, "Count")
	if got != 42 {
		t.Fatalf("expected Count=42 (int), got %v (%T)", got, got)
	}
}

// TestGoSetField_IncompatibleType verifies that a genuinely incompatible
// assignment (e.g. a string into an int field) fails with a clean,
// descriptive panic instead of leaking a raw reflect error, and that it is
// a normal recoverable panic (not a process crash).
func TestGoSetField_IncompatibleType(t *testing.T) {
	env := setupGoInterOpEnv(t)
	nativeNew := getBuiltin(t, env, "new")
	nativeSetField := getBuiltin(t, env, "setField")

	obj := nativeNew("testGoStruct")

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for incompatible field assignment, got none")
		}
		msg, ok := r.(string)
		if !ok {
			t.Fatalf("expected string panic message, got %T: %v", r, r)
		}
		if !strings.Contains(msg, "native.setField") {
			t.Fatalf("expected panic message to mention native.setField, got: %s", msg)
		}
	}()

	// Name is a string field; a struct value (the testGoStruct pointer
	// itself) is not assignable nor convertible to string.
	nativeSetField(obj, "Name", obj)
}

// TestGoSetField_NilValue ensures assigning nil produces a clean panic
// rather than a raw reflect panic from fieldVal.Set(reflect.ValueOf(nil)).
func TestGoSetField_NilValue(t *testing.T) {
	env := setupGoInterOpEnv(t)
	nativeNew := getBuiltin(t, env, "new")
	nativeSetField := getBuiltin(t, env, "setField")

	obj := nativeNew("testGoStruct")

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when assigning nil, got none")
		}
	}()
	nativeSetField(obj, "Name", nil)
}

// TestCallGoFunc_NilArgument confirms that passing nil as an argument to a
// registered Go function with a pointer parameter no longer panics with the
// opaque "reflect: Call using zero Value argument" message; nil is
// translated to the zero value of the expected parameter type.
func TestCallGoFunc_NilArgument(t *testing.T) {
	env := setupGoInterOpEnv(t)
	nativeCallFunc := getBuiltin(t, env, "callFunc")

	result := nativeCallFunc("testNilable", nil)
	if result != "nil" {
		t.Fatalf("expected 'nil', got %v", result)
	}
}

// TestCallGoFunc_Basic is a sanity check that ordinary calls still work
// after introducing buildCallArgs.
func TestCallGoFunc_Basic(t *testing.T) {
	env := setupGoInterOpEnv(t)
	nativeCallFunc := getBuiltin(t, env, "callFunc")

	result := nativeCallFunc("testAdd", 2, 3)
	if result != 5 {
		t.Fatalf("expected 5, got %v", result)
	}
}

// TestGoCallMethod_Basic exercises goCallMethod end to end through the
// buildCallArgs helper.
func TestGoCallMethod_Basic(t *testing.T) {
	env := setupGoInterOpEnv(t)
	nativeNew := getBuiltin(t, env, "new")
	nativeSetField := getBuiltin(t, env, "setField")
	nativeCallMethod := getBuiltin(t, env, "callMethod")

	obj := nativeNew("testGoStruct")
	nativeSetField(obj, "Name", "hello")

	result := nativeCallMethod(obj, "Greet", "!")
	if result != "hello!" {
		t.Fatalf("expected 'hello!', got %v", result)
	}
}

// TestCallGoFunc_Float64ArgsToIntParams confirms that calling a registered
// Go function whose parameters are int (or any other numeric kind) with the
// float64 arguments R2Lang scripts always produce no longer panics with the
// opaque "reflect: Call using float64 as type int" message. Found via a real
// embedding repro (a host Go program registering `func Add(a, b int) int`
// and calling native.callFunc("add", 3, 4) from an R2Lang script), not just
// TestCallGoFunc_Basic, which calls the builtin directly with Go int
// literals and therefore never exercised this path.
func TestCallGoFunc_Float64ArgsToIntParams(t *testing.T) {
	env := setupGoInterOpEnv(t)
	nativeCallFunc := getBuiltin(t, env, "callFunc")

	result := nativeCallFunc("testAdd", float64(2), float64(3))
	if result != 5 {
		t.Fatalf("expected 5, got %v (%T)", result, result)
	}
}

// TestGoSetField_NumericToStringNotSilentlyConverted guards against the
// broader reflect.Type.ConvertibleTo/Convert behavior, which treats an int
// as a Unicode code point when converting to string (e.g. float64(65) would
// silently become "A"). Assigning a numeric value to a string field must
// fail with a clear panic instead of silently producing a nonsensical
// string.
func TestGoSetField_NumericToStringNotSilentlyConverted(t *testing.T) {
	env := setupGoInterOpEnv(t)
	nativeNew := getBuiltin(t, env, "new")
	nativeSetField := getBuiltin(t, env, "setField")

	obj := nativeNew("testGoStruct")

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic when assigning a float64 to a string field, got none")
		}
	}()
	nativeSetField(obj, "Name", float64(65))
}
