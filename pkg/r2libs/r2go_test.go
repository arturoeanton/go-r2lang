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
	v, ok := env.Get(name)
	if !ok {
		t.Fatalf("builtin %q not registered", name)
	}
	fn, ok := v.(r2core.BuiltinFunction)
	if !ok {
		t.Fatalf("%q is not a BuiltinFunction", name)
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
	goNew := getBuiltin(t, env, "goNew")
	goSetField := getBuiltin(t, env, "goSetField")
	goGetField := getBuiltin(t, env, "goGetField")

	obj := goNew("testGoStruct")
	goSetField(obj, "Count", float64(42))

	got := goGetField(obj, "Count")
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
	goNew := getBuiltin(t, env, "goNew")
	goSetField := getBuiltin(t, env, "goSetField")

	obj := goNew("testGoStruct")

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for incompatible field assignment, got none")
		}
		msg, ok := r.(string)
		if !ok {
			t.Fatalf("expected string panic message, got %T: %v", r, r)
		}
		if !strings.Contains(msg, "goSetField") {
			t.Fatalf("expected panic message to mention goSetField, got: %s", msg)
		}
	}()

	// Name is a string field; a struct value (the testGoStruct pointer
	// itself) is not assignable nor convertible to string.
	goSetField(obj, "Name", obj)
}

// TestGoSetField_NilValue ensures assigning nil produces a clean panic
// rather than a raw reflect panic from fieldVal.Set(reflect.ValueOf(nil)).
func TestGoSetField_NilValue(t *testing.T) {
	env := setupGoInterOpEnv(t)
	goNew := getBuiltin(t, env, "goNew")
	goSetField := getBuiltin(t, env, "goSetField")

	obj := goNew("testGoStruct")

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when assigning nil, got none")
		}
	}()
	goSetField(obj, "Name", nil)
}

// TestCallGoFunc_NilArgument confirms that passing nil as an argument to a
// registered Go function with a pointer parameter no longer panics with the
// opaque "reflect: Call using zero Value argument" message; nil is
// translated to the zero value of the expected parameter type.
func TestCallGoFunc_NilArgument(t *testing.T) {
	env := setupGoInterOpEnv(t)
	callGoFunc := getBuiltin(t, env, "callGoFunc")

	result := callGoFunc("testNilable", nil)
	if result != "nil" {
		t.Fatalf("expected 'nil', got %v", result)
	}
}

// TestCallGoFunc_Basic is a sanity check that ordinary calls still work
// after introducing buildCallArgs.
func TestCallGoFunc_Basic(t *testing.T) {
	env := setupGoInterOpEnv(t)
	callGoFunc := getBuiltin(t, env, "callGoFunc")

	result := callGoFunc("testAdd", 2, 3)
	if result != 5 {
		t.Fatalf("expected 5, got %v", result)
	}
}

// TestGoCallMethod_Basic exercises goCallMethod end to end through the
// buildCallArgs helper.
func TestGoCallMethod_Basic(t *testing.T) {
	env := setupGoInterOpEnv(t)
	goNew := getBuiltin(t, env, "goNew")
	goSetField := getBuiltin(t, env, "goSetField")
	goCallMethod := getBuiltin(t, env, "goCallMethod")

	obj := goNew("testGoStruct")
	goSetField(obj, "Name", "hello")

	result := goCallMethod(obj, "Greet", "!")
	if result != "hello!" {
		t.Fatalf("expected 'hello!', got %v", result)
	}
}
