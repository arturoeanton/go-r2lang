package r2libs

import (
	"sync"
	"testing"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func syncModule(t *testing.T) map[string]interface{} {
	t.Helper()
	env := r2core.NewEnvironment()
	RegisterSync(env)
	raw, ok := env.Get("sync")
	if !ok {
		t.Fatal("sync module not registered")
	}
	mod, ok := raw.(map[string]interface{})
	if !ok {
		t.Fatalf("sync module has unexpected type %T", raw)
	}
	return mod
}

func callBuiltin(t *testing.T, mod map[string]interface{}, name string, args ...interface{}) interface{} {
	t.Helper()
	fnRaw, ok := mod[name]
	if !ok {
		t.Fatalf("sync.%s not found", name)
	}
	fn, ok := fnRaw.(r2core.BuiltinFunction)
	if !ok {
		t.Fatalf("sync.%s has unexpected type %T", name, fnRaw)
	}
	return fn(args...)
}

func method(t *testing.T, obj interface{}, name string) func(args ...interface{}) interface{} {
	t.Helper()
	getter, ok := obj.(interface {
		Getattr(string) (r2core.Node, bool)
	})
	if !ok {
		t.Fatalf("object %T does not implement Getattr", obj)
	}
	node, found := getter.Getattr(name)
	if !found {
		t.Fatalf("object %T has no method %q", obj, name)
	}
	nf, ok := node.(*NativeFunction)
	if !ok {
		t.Fatalf("Getattr(%q) returned unexpected node type %T", name, node)
	}
	return nf.Fn
}

// TestMutexConcurrentIncrement hammers a shared counter from many goroutines,
// each guarded by lock/unlock. Without correct mutual exclusion this would be
// racy and the final count would almost certainly land below N.
func TestMutexConcurrentIncrement(t *testing.T) {
	mod := syncModule(t)
	mu := callBuiltin(t, mod, "Mutex")
	lock := method(t, mu, "lock")
	unlock := method(t, mu, "unlock")

	const n = 500
	counter := 0
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			lock()
			counter++
			unlock()
		}()
	}
	wg.Wait()

	if counter != n {
		t.Fatalf("expected counter == %d, got %d (mutual exclusion violated)", n, counter)
	}
}

// TestMutexUnlockWithoutLockPanicsCleanly guards against a regression of a
// real crash bug: calling unlock() on a mutex that was never locked (or
// double-unlocking one) used to fall straight through to Go's
// sync.Mutex.Unlock, which calls the runtime's fatal() on that condition.
// fatal() is NOT a recoverable panic -- it kills the whole process even from
// inside a script's try/catch, so an R2Lang script could crash the entire
// interpreter with nothing but `sync.Mutex().unlock()`. MutexObject now
// tracks its own lock state and must turn this into an ordinary, recoverable
// Go panic instead of ever reaching the real Unlock() call in that state.
func TestMutexUnlockWithoutLockPanicsCleanly(t *testing.T) {
	mod := syncModule(t)
	mu := callBuiltin(t, mod, "Mutex")
	unlock := method(t, mu, "unlock")

	func() {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("expected unlock() of a never-locked mutex to panic")
			}
		}()
		unlock()
	}()
}

// TestMutexDoubleUnlockPanicsCleanly is the same regression guard but for the
// double-unlock case: lock() once, unlock() twice.
func TestMutexDoubleUnlockPanicsCleanly(t *testing.T) {
	mod := syncModule(t)
	mu := callBuiltin(t, mod, "Mutex")
	lock := method(t, mu, "lock")
	unlock := method(t, mu, "unlock")

	lock()
	unlock()

	func() {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("expected second unlock() to panic")
			}
		}()
		unlock()
	}()
}

// TestMutexTryLock checks that tryLock reports false while held and true once
// released.
func TestMutexTryLock(t *testing.T) {
	mod := syncModule(t)
	mu := callBuiltin(t, mod, "Mutex")
	lock := method(t, mu, "lock")
	unlock := method(t, mu, "unlock")
	tryLock := method(t, mu, "tryLock")

	if got := tryLock(); got != true {
		t.Fatalf("tryLock on free mutex = %v, want true", got)
	}
	unlock()

	lock()
	if got := tryLock(); got != false {
		t.Fatalf("tryLock on held mutex = %v, want false", got)
	}
	unlock()
}

// TestSemaphoreLimitsConcurrency verifies the semaphore never allows more
// than N holders in its critical section at once, by tracking the observed
// concurrency high-water mark under its own mutex.
func TestSemaphoreLimitsConcurrency(t *testing.T) {
	mod := syncModule(t)
	const permits = 3
	sem := callBuiltin(t, mod, "Semaphore", float64(permits))
	acquire := method(t, sem, "acquire")
	release := method(t, sem, "release")

	var mu sync.Mutex
	current := 0
	maxObserved := 0

	const workers = 20
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			acquire()
			mu.Lock()
			current++
			if current > maxObserved {
				maxObserved = current
			}
			mu.Unlock()

			mu.Lock()
			current--
			mu.Unlock()
			release()
		}()
	}
	wg.Wait()

	if maxObserved > permits {
		t.Fatalf("observed %d concurrent holders, want <= %d", maxObserved, permits)
	}
	if maxObserved == 0 {
		t.Fatal("no concurrency was ever observed; test is not exercising the semaphore")
	}
}

// TestWaitGroupWaitsForAll confirms wait() blocks until every done() call has
// landed, using a value only written right before done() to catch a wait()
// that returns early.
func TestWaitGroupWaitsForAll(t *testing.T) {
	mod := syncModule(t)
	wgObj := callBuiltin(t, mod, "WaitGroup")
	add := method(t, wgObj, "add")
	done := method(t, wgObj, "done")
	wait := method(t, wgObj, "wait")

	const n = 200
	add(float64(n))

	var mu sync.Mutex
	completed := 0
	for i := 0; i < n; i++ {
		go func() {
			mu.Lock()
			completed++
			mu.Unlock()
			done()
		}()
	}

	wait()

	mu.Lock()
	defer mu.Unlock()
	if completed != n {
		t.Fatalf("wait() returned before all goroutines finished: completed=%d, want %d", completed, n)
	}
}

// TestOnceRunsExactlyOnce races many goroutines calling do() and checks the
// callback body executed exactly once.
func TestOnceRunsExactlyOnce(t *testing.T) {
	mod := syncModule(t)
	onceObj := callBuiltin(t, mod, "Once")
	do := method(t, onceObj, "do")

	env := r2core.NewEnvironment()
	var mu sync.Mutex
	runs := 0
	env.Set("record", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		mu.Lock()
		runs++
		mu.Unlock()
		return nil
	}))

	parser := r2core.NewParser("func bump() { record() }")
	parser.ParseProgram().Eval(env)
	fnRaw, ok := env.Get("bump")
	if !ok {
		t.Fatal("bump function not found in environment")
	}
	fn := fnRaw.(*r2core.UserFunction)

	const workers = 50
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			do(fn)
		}()
	}
	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	if runs != 1 {
		t.Fatalf("Once.do executed %d times, want exactly 1", runs)
	}
}

// TestMutexHammerFromManyGoroutines drives lock/unlock/tryLock on a single
// shared MutexObject from many goroutines in a tight loop, the kind of abuse
// a script fanning out via r2()/go() could produce. Run with -race; this is
// meant to catch data races in MutexObject's own bookkeeping (the `locked`
// field added to detect unlock-of-unlocked-mutex), not just in the
// underlying sync.Mutex.
func TestMutexHammerFromManyGoroutines(t *testing.T) {
	mod := syncModule(t)
	mu := callBuiltin(t, mod, "Mutex")
	lock := method(t, mu, "lock")
	unlock := method(t, mu, "unlock")
	tryLock := method(t, mu, "tryLock")

	const workers = 50
	const iterations = 200
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				lock()
				unlock()
				if tryLock().(bool) {
					unlock()
				}
			}
		}()
	}
	wg.Wait()
}
