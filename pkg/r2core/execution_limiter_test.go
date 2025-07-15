package r2core

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestExecutionLimiter_Basic(t *testing.T) {
	limiter := NewExecutionLimiter()

	// Test initial state
	if limiter.CurrentIterations != 0 {
		t.Errorf("Expected initial iterations to be 0, got %d", limiter.CurrentIterations)
	}

	if len(limiter.CallStack) != 0 {
		t.Errorf("Expected initial call stack to be empty, got %d", len(limiter.CallStack))
	}

	if !limiter.Enabled {
		t.Error("Expected limiter to be enabled by default")
	}
}

func TestExecutionLimiter_IterationLimit(t *testing.T) {
	limiter := NewExecutionLimiter()
	limiter.MaxIterations = 5

	// Test under limit
	for i := 0; i < 4; i++ {
		limiter.IncrementIterations()
		if limiter.CheckIterationLimit() {
			t.Errorf("Should not hit limit at iteration %d", i)
		}
	}

	// Test at limit
	limiter.IncrementIterations()
	if !limiter.CheckIterationLimit() {
		t.Error("Should hit limit at max iterations")
	}
}

func TestExecutionLimiter_RecursionDepth(t *testing.T) {
	limiter := NewExecutionLimiter()
	limiter.MaxRecursionDepth = 3

	// Test under limit
	limiter.EnterFunction("func1")
	limiter.EnterFunction("func2")
	if limiter.CheckRecursionDepth() {
		t.Error("Should not hit recursion limit")
	}

	// Test at limit
	limiter.EnterFunction("func3")
	if !limiter.CheckRecursionDepth() {
		t.Error("Should hit recursion limit")
	}

	// Test exit function
	limiter.ExitFunction()
	if limiter.CheckRecursionDepth() {
		t.Error("Should not hit limit after exit")
	}
}

func TestExecutionLimiter_TimeLimit(t *testing.T) {
	limiter := NewExecutionLimiter()
	limiter.MaxExecutionTime = 50 * time.Millisecond

	// Test under time limit
	if limiter.CheckTimeLimit() {
		t.Error("Should not hit time limit immediately")
	}

	// Wait and test over time limit
	time.Sleep(60 * time.Millisecond)
	if !limiter.CheckTimeLimit() {
		t.Error("Should hit time limit after waiting")
	}
}

func TestExecutionLimiter_Context(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	limiter := &ExecutionLimiter{
		Enabled: true,
		Context: ctx,
	}

	// Test before timeout
	if limiter.CheckContext() {
		t.Error("Should not hit context limit immediately")
	}

	// Wait for timeout
	time.Sleep(60 * time.Millisecond)
	if !limiter.CheckContext() {
		t.Error("Should hit context limit after timeout")
	}
}

func TestExecutionLimiter_Disable(t *testing.T) {
	limiter := NewExecutionLimiter()
	limiter.MaxIterations = 1
	limiter.MaxRecursionDepth = 1
	limiter.MaxExecutionTime = 1 * time.Nanosecond

	// Disable limiter
	limiter.Disable()

	// None of the limits should trigger
	for i := 0; i < 10; i++ {
		limiter.IncrementIterations()
		limiter.EnterFunction("test")
	}

	if limiter.CheckIterationLimit() {
		t.Error("Disabled limiter should not check iteration limit")
	}

	if limiter.CheckRecursionDepth() {
		t.Error("Disabled limiter should not check recursion depth")
	}

	if limiter.CheckTimeLimit() {
		t.Error("Disabled limiter should not check time limit")
	}
}

func TestExecutionLimiter_Reset(t *testing.T) {
	limiter := NewExecutionLimiter()

	// Add some state
	limiter.CurrentIterations = 100
	limiter.EnterFunction("test1")
	limiter.EnterFunction("test2")

	// Reset
	limiter.Reset()

	// Check state is reset
	if limiter.CurrentIterations != 0 {
		t.Errorf("Expected iterations to be reset to 0, got %d", limiter.CurrentIterations)
	}

	if len(limiter.CallStack) != 0 {
		t.Errorf("Expected call stack to be reset to empty, got %d", len(limiter.CallStack))
	}
}

func TestInfiniteLoopError(t *testing.T) {
	ctx := &LoopContext{
		Type:       "while",
		Iterations: 1000,
		StartTime:  time.Now().Add(-5 * time.Second),
		Location:   "test.r2:10",
	}

	err := NewInfiniteLoopError("while", ctx)

	if err.Type != "while" {
		t.Errorf("Expected error type 'while', got '%s'", err.Type)
	}

	if err.Iterations != 1000 {
		t.Errorf("Expected 1000 iterations, got %d", err.Iterations)
	}

	if err.Location != "test.r2:10" {
		t.Errorf("Expected location 'test.r2:10', got '%s'", err.Location)
	}

	if !errors.Is(err, ErrInfiniteLoop) {
		t.Error("Error should be identified as ErrInfiniteLoop")
	}

	// Test error message
	message := err.Error()
	if message == "" {
		t.Error("Error message should not be empty")
	}

	// Test short message
	shortMsg := err.ShortMessage()
	if shortMsg == "" {
		t.Error("Short message should not be empty")
	}
}

func TestRecursionError(t *testing.T) {
	err := NewRecursionError("max_depth", 1000)

	if err.Type != "recursion" {
		t.Errorf("Expected error type 'recursion', got '%s'", err.Type)
	}

	if !errors.Is(err, ErrRecursionLimit) {
		t.Error("Error should be identified as ErrRecursionLimit")
	}
}

func TestTimeoutError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := NewTimeoutError("execution_timeout", ctx)

	if err.Type != "timeout" {
		t.Errorf("Expected error type 'timeout', got '%s'", err.Type)
	}

	if !errors.Is(err, ErrTimeout) {
		t.Error("Error should be identified as ErrTimeout")
	}
}

func TestExecutionLimiter_SetLimits(t *testing.T) {
	limiter := NewExecutionLimiter()

	// Set custom limits
	limiter.SetLimits(500, 100, 10*time.Second)

	if limiter.MaxIterations != 500 {
		t.Errorf("Expected max iterations 500, got %d", limiter.MaxIterations)
	}

	if limiter.MaxRecursionDepth != 100 {
		t.Errorf("Expected max recursion depth 100, got %d", limiter.MaxRecursionDepth)
	}

	if limiter.MaxExecutionTime != 10*time.Second {
		t.Errorf("Expected max execution time 10s, got %v", limiter.MaxExecutionTime)
	}
}
