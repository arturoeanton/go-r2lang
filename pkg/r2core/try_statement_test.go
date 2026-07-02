package r2core

import (
	"testing"
)

func evalTryCode(t *testing.T, code string) (result interface{}, panicVal interface{}) {
	t.Helper()
	env := NewEnvironment()
	env.Set("true", true)
	env.Set("false", false)

	defer func() {
		panicVal = recover()
	}()

	parser := NewParser(code)
	ast := parser.ParseProgram()
	result = ast.Eval(env)
	return
}

// TestTryStatement_FinallyRunsWhenCatchThrows guards against a regression
// where a panic raised inside the catch block skipped the finally block
// entirely, because the catch block ran without its own recover() and its
// panic unwound straight past the "if ts.FinallyBlock != nil" check.
func TestTryStatement_FinallyRunsWhenCatchThrows(t *testing.T) {
	env := NewEnvironment()
	env.Set("true", true)
	env.Set("false", false)

	code := `
		let ranFinally = false;
		try {
			try {
				throw "original";
			} catch (e) {
				throw "from-catch";
			} finally {
				ranFinally = true;
			}
		} catch (outer) {
			outer;
		}
	`

	parser := NewParser(code)
	ast := parser.ParseProgram()
	result := ast.Eval(env)

	ranFinally, ok := env.Get("ranFinally")
	if !ok || ranFinally != true {
		t.Fatalf("expected finally block to run when catch throws, ranFinally=%v", ranFinally)
	}
	if result != "from-catch" {
		t.Fatalf("expected the catch block's exception to propagate, got %v", result)
	}
}

// TestTryStatement_FinallyExceptionReplacesCatchException checks that when
// both the catch block and the finally block throw, the finally block's
// exception is the one that ultimately propagates (matching JS/Java/Python).
func TestTryStatement_FinallyExceptionReplacesCatchException(t *testing.T) {
	code := `
		try {
			throw "A";
		} catch (e) {
			throw "B";
		} finally {
			throw "C";
		}
	`
	_, panicVal := evalTryCode(t, code)
	if panicVal != "C" {
		t.Fatalf("expected finally's exception 'C' to win, got %v", panicVal)
	}
}

// TestTryStatement_NestedFinallyReturnOverridesException checks that a
// return inside an inner finally overrides a pending exception from the
// same inner try, and that the outer catch never runs.
func TestTryStatement_NestedFinallyReturnOverridesException(t *testing.T) {
	env := NewEnvironment()
	env.Set("true", true)
	env.Set("false", false)

	code := `
		func f() {
			try {
				try {
					throw "inner-original";
				} finally {
					return "inner-finally-return";
				}
			} catch (e) {
				return "outer-caught:" + e;
			}
		}
		f();
	`
	parser := NewParser(code)
	ast := parser.ParseProgram()
	result := ast.Eval(env)

	if result != "inner-finally-return" {
		t.Fatalf("expected inner finally's return to win, got %v", result)
	}
}

// TestTryStatement_ContinueInFinallyInsideLoop checks that a continue
// statement inside a finally block correctly skips the rest of the
// enclosing loop iteration.
func TestTryStatement_ContinueInFinallyInsideLoop(t *testing.T) {
	env := NewEnvironment()
	env.Set("true", true)
	env.Set("false", false)

	code := `
		let afterTryCount = 0;
		let i = 0;
		while (i < 5) {
			i = i + 1;
			try {
			} finally {
				if (i == 3) {
					continue;
				}
			}
			afterTryCount = afterTryCount + 1;
		}
		afterTryCount;
	`
	parser := NewParser(code)
	ast := parser.ParseProgram()
	result := ast.Eval(env)

	// i=1,2,4,5 should reach "afterTryCount++" (4 times); i=3 is skipped by continue.
	if result.(float64) != 4 {
		t.Fatalf("expected afterTryCount=4, got %v", result)
	}
}

// TestTryStatement_TryInsideCatchInsideFinally exercises deep nesting to
// make sure control flow and exception propagation compose correctly.
func TestTryStatement_TryInsideCatchInsideFinally(t *testing.T) {
	env := NewEnvironment()
	env.Set("true", true)
	env.Set("false", false)

	code := `
		let trace = [];
		try {
			throw "L1";
		} catch (e1) {
			try {
				throw "L2";
			} catch (e2) {
				trace = trace + [e2];
			} finally {
				trace = trace + ["inner-finally"];
			}
			trace = trace + [e1];
		} finally {
			trace = trace + ["outer-finally"];
		}
		trace;
	`
	parser := NewParser(code)
	ast := parser.ParseProgram()
	result := ast.Eval(env)

	trace, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected trace to be a slice, got %T: %v", result, result)
	}
	expected := []interface{}{"L2", "inner-finally", "L1", "outer-finally"}
	if len(trace) != len(expected) {
		t.Fatalf("expected trace %v, got %v", expected, trace)
	}
	for i := range expected {
		if trace[i] != expected[i] {
			t.Fatalf("expected trace %v, got %v", expected, trace)
		}
	}
}
