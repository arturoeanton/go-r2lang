package r2libs

import (
	"strings"
	"sync"
	"testing"
	"time"
)

// TestMonitor_UnlockWithoutLock_PanicsCleanly is a regression test for a bug
// where calling Unlock() on a Monitor that was never locked reached the raw
// sync.Mutex.Unlock() and triggered "fatal error: sync: unlock of unlocked
// mutex" - a Go runtime fatal() that crashes the whole process and bypasses
// every defer/recover(), unlike a normal panic. Confirmed via a built binary
// running `goroutine.unlock(mon)` twice from an .r2 script before this fix
// (second call killed the process with exit code 2 and a runtime stack
// trace). After the fix, misuse must degrade to an ordinary, recoverable
// panic.
func TestMonitor_UnlockWithoutLock_PanicsCleanly(t *testing.T) {
	mon := NewMonitor()

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected a panic when unlocking a monitor that was never locked")
		}
		msg, ok := r.(string)
		if !ok || !strings.Contains(msg, "not locked") {
			t.Fatalf("expected a clean 'not locked' panic message, got: %v", r)
		}
	}()

	mon.Unlock()
}

// TestMonitor_DoubleUnlock_PanicsCleanly reproduces calling unlock() twice
// in a row (lock -> unlock -> unlock), which previously crashed the whole
// process with the same unrecoverable Go runtime fatal error.
func TestMonitor_DoubleUnlock_PanicsCleanly(t *testing.T) {
	mon := NewMonitor()
	mon.Lock()
	mon.Unlock()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected a panic on the second unlock")
		}
	}()
	mon.Unlock()
}

// TestMonitor_WaitWithoutLock_PanicsCleanly reproduces calling wait() before
// ever locking the monitor, which previously reached sync.Cond.Wait()'s
// internal c.L.Unlock() on an unlocked mutex and crashed the process the
// same way as the bare Unlock() case.
func TestMonitor_WaitWithoutLock_PanicsCleanly(t *testing.T) {
	mon := NewMonitor()

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected a panic when waiting on a monitor that was never locked")
		}
		msg, ok := r.(string)
		if !ok || !strings.Contains(msg, "must be locked") {
			t.Fatalf("expected a clean 'must be locked' panic message, got: %v", r)
		}
	}()

	mon.Wait()
}

// TestMonitor_NormalUsage_StillWorks is a functional regression test making
// sure the lock-state tracking added to fix the fatal-crash bugs does not
// break the legitimate lock/wait/signal/unlock producer-consumer pattern.
func TestMonitor_NormalUsage_StillWorks(t *testing.T) {
	mon := NewMonitor()
	done := false

	go func() {
		mon.Lock()
		done = true
		mon.Signal()
		mon.Unlock()
	}()

	mon.Lock()
	for !done {
		mon.Wait()
	}
	mon.Unlock()

	if !done {
		t.Fatal("expected done to be true after Wait returned")
	}
}

// TestMonitor_ConcurrentProducerConsumer_Race exercises the Monitor under
// -race with multiple goroutines performing correctly-paired lock/wait/
// signal/unlock sequences, guarding against the atomic lock-state tracking
// introducing a data race.
func TestMonitor_ConcurrentProducerConsumer_Race(t *testing.T) {
	mon := NewMonitor()
	queue := make([]int, 0, 100)
	const items = 50

	var wg sync.WaitGroup
	wg.Add(2)

	go func() { // producer
		defer wg.Done()
		for i := 0; i < items; i++ {
			mon.Lock()
			queue = append(queue, i)
			mon.Signal()
			mon.Unlock()
			time.Sleep(time.Microsecond)
		}
	}()

	go func() { // consumer
		defer wg.Done()
		consumed := 0
		for consumed < items {
			mon.Lock()
			for len(queue) == 0 {
				mon.Wait()
			}
			queue = queue[1:]
			consumed++
			mon.Unlock()
		}
	}()

	waitCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitCh)
	}()

	select {
	case <-waitCh:
	case <-time.After(10 * time.Second):
		t.Fatal("producer/consumer test timed out - possible deadlock")
	}
}

// TestSemaphore_BasicAcquireRelease is a sanity check for the Semaphore
// primitive living in the same file, unaffected by this fix but useful as a
// baseline race-safety check for the file as a whole.
func TestSemaphore_BasicAcquireRelease(t *testing.T) {
	sem := NewSemaphore(2)
	sem.Acquire()
	sem.Acquire()

	acquired := make(chan struct{})
	go func() {
		sem.Acquire()
		close(acquired)
	}()

	select {
	case <-acquired:
		t.Fatal("third Acquire should have blocked until a Release")
	case <-time.After(50 * time.Millisecond):
	}

	sem.Release()

	select {
	case <-acquired:
	case <-time.After(time.Second):
		t.Fatal("third Acquire should have unblocked after Release")
	}
}
