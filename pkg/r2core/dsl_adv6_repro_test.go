package r2core

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestDSLReentrantUseDeadlock checks that an action function calling .use()
// on the SAME dsl object it belongs to (re-entrant use while useMu is
// already held by the outer call, same goroutine) does not deadlock, and
// that the outer call's remaining state (currentExecutionEnv) is restored
// correctly once the nested call returns.
func TestDSLReentrantUseDeadlock(t *testing.T) {
	env := NewEnvironment()

	dslCode := `
	dsl ReDSL {
		token("A", "a")
		token("B", "b")

		rule("start", ["A"], "recurse")
		rule("start", ["B"], "base")

		func recurse(tok) {
			let inner = ReDSL.use("b")
			return "outer(" + inner.Output + ")"
		}

		func base(tok) {
			return "base-" + tok
		}
	}
	`

	parser := NewParser(dslCode)
	program := parser.ParseProgram()
	program.Eval(env)

	dslObj, exists := env.Get("ReDSL")
	if !exists {
		t.Fatal("DSL 'ReDSL' not found")
	}
	dslMap := dslObj.(map[string]interface{})
	useFunc := dslMap["use"].(func(...interface{}) interface{})

	done := make(chan interface{}, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- fmt.Errorf("panic: %v", r)
			}
		}()
		done <- useFunc("a")
	}()

	select {
	case result := <-done:
		dslResult, ok := result.(*DSLResult)
		if !ok {
			t.Fatalf("expected *DSLResult, got %T (%v)", result, result)
		}
		if dslResult.Output != "outer(base-b)" {
			t.Errorf("expected 'outer(base-b)', got %q", dslResult.Output)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("reentrant .use() call from within an action deadlocked (useMu is not reentrant)")
	}
}

// TestReentrantMutexStillExcludesOtherGoroutines confirms that making useMu
// reentrant (to fix same-goroutine re-entry) did not also make it a no-op
// for real cross-goroutine mutual exclusion.
func TestReentrantMutexStillExcludesOtherGoroutines(t *testing.T) {
	var mu reentrantMutex
	var counter int
	var wg sync.WaitGroup

	const n = 200
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			counter++
			local := counter
			if local != counter {
				t.Errorf("race detected: counter changed under lock")
			}
		}()
	}
	wg.Wait()

	if counter != n {
		t.Errorf("expected counter == %d, got %d", n, counter)
	}
}

// TestReentrantMutexSameGoroutineNestsCleanly exercises Lock/Unlock nesting
// directly, independent of the DSL layer.
func TestReentrantMutexSameGoroutineNestsCleanly(t *testing.T) {
	var mu reentrantMutex
	mu.Lock()
	mu.Lock()
	mu.Lock()
	mu.Unlock()
	mu.Unlock()

	done := make(chan struct{})
	go func() {
		mu.Lock()
		close(done)
		mu.Unlock()
	}()

	select {
	case <-done:
		t.Fatal("another goroutine acquired the lock while the owner still held a nested Lock()")
	case <-time.After(200 * time.Millisecond):
	}

	mu.Unlock()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("other goroutine never acquired the lock after the owner fully unlocked")
	}
}

// TestDSLInvalidTokenSurfacesError covers extractToken's fallback path: when
// a token's regex pattern is invalid (or improveRegexPattern's "improved"
// variant is), BOTH AddToken attempts previously failed silently (the
// second call's error was never checked), leaving the token unregistered
// with no diagnostic at all. This must now fail loudly at DSL-definition
// time with the real regex-compile error, whether or not any rule
// references the broken token.
func TestDSLInvalidTokenSurfacesError(t *testing.T) {
	cases := []struct {
		name    string
		dslCode string
		wantSub string
	}{
		{
			name: "referenced by a rule",
			dslCode: `
			dsl BadTokRef {
				token("A", "a")
				token("BAD", "[0-9")
				rule("start", ["A", "BAD"], "act")
				func act(a, b) { return a + b }
			}
			`,
			wantSub: "invalid token",
		},
		{
			name: "unreferenced by any rule",
			dslCode: `
			dsl BadTokUnref {
				token("A", "a")
				token("BAD", "[0-9")
				rule("start", ["A"], "act")
				func act(a) { return a }
			}
			`,
			wantSub: "invalid token",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			env := NewEnvironment()
			parser := NewParser(tc.dslCode)
			program := parser.ParseProgram()

			defer func() {
				r := recover()
				if r == nil {
					t.Fatal("expected a panic reporting the invalid regex pattern, got none (error was silently swallowed)")
				}
				msg := fmt.Sprintf("%v", r)
				if !strings.Contains(msg, tc.wantSub) {
					t.Errorf("expected panic message to contain %q, got: %s", tc.wantSub, msg)
				}
			}()
			program.Eval(env)
		})
	}
}

// TestDSLInvalidLiteralSurfacesError covers extractLiteral, which
// previously ignored DSLGrammar.AddLiteral's error return entirely. An
// empty literal() text quotes down to an empty-matching token, which
// go-dsl's AddToken correctly rejects; that rejection must not be dropped.
func TestDSLInvalidLiteralSurfacesError(t *testing.T) {
	env := NewEnvironment()
	dslCode := `
	dsl BadLit {
		literal("EMPTY", "")
		token("A", "a")
		rule("start", ["A"], "act")
		func act(a) { return a }
	}
	`
	parser := NewParser(dslCode)
	program := parser.ParseProgram()

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected a panic reporting the invalid literal pattern, got none (error was silently swallowed)")
		}
		msg := fmt.Sprintf("%v", r)
		if !strings.Contains(msg, "invalid literal") {
			t.Errorf("expected panic message to contain 'invalid literal', got: %s", msg)
		}
	}()
	program.Eval(env)
}

// TestDSLActionInfiniteRecursionIsCaught checks that a DSL action calling
// .use() on its own DSL with no terminating base case is caught by the
// same recursion-depth limiter normal R2Lang function calls get, instead
// of hanging forever (callDSLFunction previously bypassed UserFunction.Call
// entirely, so it never registered with the ExecutionLimiter).
func TestDSLActionInfiniteRecursionIsCaught(t *testing.T) {
	env := NewEnvironment()
	// Use a much smaller recursion-depth ceiling than the default (1000)
	// so this test's own runtime is bounded by an iteration COUNT rather
	// than wall-clock time. Each recursive .use() call here re-runs a full
	// DSL parse+evaluate cycle, which is real work (~2ms/level locally) —
	// waiting to hit the default 1000-level limit took ~2s on this
	// machine, and reliably exceeded a 5s wall-clock ceiling on slower/
	// shared CI hardware, making the test flake there despite the actual
	// recursion-limiting mechanism working correctly. A smaller depth
	// limit makes the test fast and deterministic regardless of hardware.
	limiter := NewExecutionLimiter()
	limiter.MaxRecursionDepth = 50
	env.SetLimiter(limiter)
	dslCode := `
	dsl Loopy {
		token("WORD", "[a-zA-Z]+")
		rule("wrap", ["WORD"], "wrapIt")
		func wrapIt(w) {
			let inner = Loopy.use(w);
			return "wrapped(" + inner.Output + ")";
		}
	}
	`
	parser := NewParser(dslCode)
	program := parser.ParseProgram()
	program.Eval(env)

	dslObj, _ := env.Get("Loopy")
	dslMap := dslObj.(map[string]interface{})
	useFunc := dslMap["use"].(func(...interface{}) interface{})

	done := make(chan interface{}, 1)
	go func() {
		defer func() { done <- recover() }()
		useFunc("hello")
	}()

	select {
	case r := <-done:
		if r == nil {
			t.Fatal("expected a recursion-limit panic, got none")
		}
		msg := fmt.Sprintf("%v", r)
		if !strings.Contains(msg, "recursion") && !strings.Contains(msg, "Loop infinito") {
			t.Errorf("expected a recursion-limit error, got: %s", msg)
		}
	case <-time.After(20 * time.Second):
		t.Fatal("infinite recursion through a DSL action was not caught within 20s (hung instead of hitting the recursion limit)")
	}
}
