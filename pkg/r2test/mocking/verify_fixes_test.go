package mocking

import (
	"sync"
	"testing"
)

// Regression test: Times(0)/AtMost(0) must reject any call, not silently
// allow unlimited calls. Previously the MaxCalls guard used `> 0`, which
// treated an explicit "exactly zero" bound the same as "unlimited" (-1).
func TestBugfix_TimesZeroRejectsCalls(t *testing.T) {
	mock := NewMock("m")
	mock.When("f").Times(0)

	if _, err := mock.Call("f"); err == nil {
		t.Fatal("expected error calling f when Times(0) was set")
	}

	// AtMost(0) should behave the same way.
	mock2 := NewMock("m2")
	mock2.When("g").AtMost(0)
	if _, err := mock2.Call("g"); err == nil {
		t.Fatal("expected error calling g when AtMost(0) was set")
	}
}

// Regression test: CreateContext must never hand out an ID that collides
// with a still-live context, even after other contexts have been cleaned
// up in between. Previously the ID was derived from len(ti.contexts),
// which can repeat once contexts are removed from the map.
func TestBugfix_ContextIDsNeverCollideAfterCleanup(t *testing.T) {
	ti := NewTestIsolation()

	a := ti.CreateContext("x")
	b := ti.CreateContext("x")

	if err := ti.CleanupContext(a.ID); err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}

	c := ti.CreateContext("x")

	if c.ID == b.ID {
		t.Fatalf("new context ID %q collides with live context b's ID %q", c.ID, b.ID)
	}
	if _, ok := ti.GetContext(b.ID); !ok {
		t.Fatal("context b should still be retrievable and untouched by c's creation")
	}
}

// Regression test: IsolationContext's maps must be safe for concurrent
// access, since IsolationContext exists specifically to support isolated
// concurrent test execution.
func TestBugfix_IsolationContextConcurrentAccessSafe(t *testing.T) {
	ctx := CreateIsolationContext("race")
	defer CleanupIsolationContext(ctx.ID)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(3)
		go func(n int) {
			defer wg.Done()
			ctx.SetGlobalVariable("k", n)
		}(i)
		go func() {
			defer wg.Done()
			ctx.CreateMock("m")
		}()
		go func() {
			defer wg.Done()
			_, _ = ctx.GetGlobalVariable("k")
		}()
	}
	wg.Wait()
}

// Regression test: Spy.Call must not panic when the call-through target
// receives a nil argument (a nil interface{}, as commonly produced when
// bridging "null" from R2Lang).
func TestBugfix_SpyCallThroughNilArgDoesNotPanic(t *testing.T) {
	orig := func(a interface{}) interface{} { return a }
	spy := NewSpy("f", orig)
	spy.CallThrough()

	results, err := spy.Call(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0] != nil {
		t.Fatalf("expected [nil], got %v", results)
	}
}

// Regression test: Spy.Call's call-through must not panic on an arity
// mismatch either; it should return a descriptive error instead.
func TestBugfix_SpyCallThroughArityMismatchReturnsError(t *testing.T) {
	orig := func(a, b int) int { return a + b }
	spy := NewSpy("f", orig)
	spy.CallThrough()

	if _, err := spy.Call(1); err == nil {
		t.Fatal("expected error for arity mismatch, got nil")
	}
}
