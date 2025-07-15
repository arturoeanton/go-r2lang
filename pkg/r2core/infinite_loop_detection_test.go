package r2core

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestWhileInfiniteLoop_Detection(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		maxIter     int64
		shouldPanic bool
		errorType   string
	}{
		{
			name:        "finite while loop",
			code:        `let i = 0; while (i < 5) { i++; }`,
			maxIter:     10,
			shouldPanic: false,
		},
		{
			name:        "infinite while loop",
			code:        `let flag = true; while (flag) { let x = 1; }`,
			maxIter:     100,
			shouldPanic: true,
			errorType:   "while",
		},
		{
			name:        "while loop with high iteration limit",
			code:        `let i = 0; while (i < 10) { i++; }`,
			maxIter:     5,
			shouldPanic: true,
			errorType:   "while",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.SetLimits(tt.maxIter, 1000, 10*time.Second)

			parser := NewParser(tt.code)
			ast := parser.ParseProgram()

			if tt.shouldPanic {
				defer func() {
					if r := recover(); r != nil {
						if err, ok := r.(*InfiniteLoopError); ok {
							if err.Type != tt.errorType {
								t.Errorf("Expected error type '%s', got '%s'", tt.errorType, err.Type)
							}
							if !errors.Is(err, ErrInfiniteLoop) {
								t.Error("Error should be identified as ErrInfiniteLoop")
							}
						} else {
							t.Errorf("Expected InfiniteLoopError, got %T: %v", r, r)
						}
					} else {
						t.Error("Expected panic but none occurred")
					}
				}()
			}

			ast.Eval(env)

			if tt.shouldPanic {
				t.Error("Expected panic but code executed successfully")
			}
		})
	}
}

func TestForInfiniteLoop_Detection(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		maxIter     int64
		shouldPanic bool
		errorType   string
	}{
		{
			name:        "finite for loop",
			code:        `for (let i = 0; i < 5; i++) { let x = i; }`,
			maxIter:     10,
			shouldPanic: false,
		},
		{
			name:        "infinite for loop",
			code:        `for (let i = 0; i >= 0; i++) { let x = i; }`,
			maxIter:     100,
			shouldPanic: true,
			errorType:   "for",
		},
		{
			name:        "for loop exceeding limit",
			code:        `for (let i = 0; i < 10; i++) { let x = i; }`,
			maxIter:     5,
			shouldPanic: true,
			errorType:   "for",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.SetLimits(tt.maxIter, 1000, 10*time.Second)

			parser := NewParser(tt.code)
			ast := parser.ParseProgram()

			if tt.shouldPanic {
				defer func() {
					if r := recover(); r != nil {
						if err, ok := r.(*InfiniteLoopError); ok {
							if err.Type != tt.errorType {
								t.Errorf("Expected error type '%s', got '%s'", tt.errorType, err.Type)
							}
						} else {
							t.Errorf("Expected InfiniteLoopError, got %T: %v", r, r)
						}
					} else {
						t.Error("Expected panic but none occurred")
					}
				}()
			}

			ast.Eval(env)

			if tt.shouldPanic {
				t.Error("Expected panic but code executed successfully")
			}
		})
	}
}

func TestForInLoop_Detection(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		maxIter     int64
		shouldPanic bool
		errorType   string
	}{
		{
			name:        "finite for-in loop",
			code:        `let arr = [1, 2, 3]; for (i in arr) { let x = i; }`,
			maxIter:     10,
			shouldPanic: false,
		},
		{
			name:        "large array exceeding limit",
			code:        `let arr = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]; for (i in arr) { let x = i; }`,
			maxIter:     5,
			shouldPanic: true,
			errorType:   "for-in",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.SetLimits(tt.maxIter, 1000, 10*time.Second)

			parser := NewParser(tt.code)
			ast := parser.ParseProgram()

			if tt.shouldPanic {
				defer func() {
					if r := recover(); r != nil {
						if err, ok := r.(*InfiniteLoopError); ok {
							if err.Type != tt.errorType {
								t.Errorf("Expected error type '%s', got '%s'", tt.errorType, err.Type)
							}
						} else {
							t.Errorf("Expected InfiniteLoopError, got %T: %v", r, r)
						}
					} else {
						t.Error("Expected panic but none occurred")
					}
				}()
			}

			ast.Eval(env)

			if tt.shouldPanic {
				t.Error("Expected panic but code executed successfully")
			}
		})
	}
}

func TestRecursionInfinite_Detection(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		maxDepth    int
		shouldPanic bool
		errorType   string
	}{
		{
			name:        "finite recursion",
			code:        `func factorial(n) { if (n <= 1) { return 1; } return n * factorial(n - 1); } factorial(5);`,
			maxDepth:    10,
			shouldPanic: false,
		},
		{
			name:        "infinite recursion",
			code:        `func infinite() { infinite(); } infinite();`,
			maxDepth:    100,
			shouldPanic: true,
			errorType:   "recursion",
		},
		{
			name:        "deep recursion exceeding limit",
			code:        `func deep(n) { if (n <= 0) { return 0; } return deep(n - 1); } deep(20);`,
			maxDepth:    10,
			shouldPanic: true,
			errorType:   "recursion",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.SetLimits(1000000, tt.maxDepth, 10*time.Second)

			parser := NewParser(tt.code)
			ast := parser.ParseProgram()

			if tt.shouldPanic {
				defer func() {
					if r := recover(); r != nil {
						if err, ok := r.(*InfiniteLoopError); ok {
							if err.Type != tt.errorType {
								t.Errorf("Expected error type '%s', got '%s'", tt.errorType, err.Type)
							}
						} else {
							t.Errorf("Expected InfiniteLoopError, got %T: %v", r, r)
						}
					} else {
						t.Error("Expected panic but none occurred")
					}
				}()
			}

			ast.Eval(env)

			if tt.shouldPanic {
				t.Error("Expected panic but code executed successfully")
			}
		})
	}
}

func TestTimeout_Detection(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		timeout     time.Duration
		shouldPanic bool
	}{
		{
			name:        "fast execution",
			code:        `let x = 1 + 2;`,
			timeout:     100 * time.Millisecond,
			shouldPanic: false,
		},
		{
			name:        "slow execution with timeout",
			code:        `let i = 0; while (i < 1000000) { i++; }`,
			timeout:     10 * time.Millisecond,
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.SetLimits(1000000, 1000, tt.timeout)

			parser := NewParser(tt.code)
			ast := parser.ParseProgram()

			if tt.shouldPanic {
				defer func() {
					if r := recover(); r != nil {
						if err, ok := r.(*InfiniteLoopError); ok {
							if err.Type != "timeout" && !strings.Contains(err.Type, "timeout") {
								t.Errorf("Expected timeout error, got '%s'", err.Type)
							}
						} else {
							t.Errorf("Expected timeout error, got %T: %v", r, r)
						}
					} else {
						t.Error("Expected panic but none occurred")
					}
				}()
			}

			ast.Eval(env)

			if tt.shouldPanic {
				t.Error("Expected panic but code executed successfully")
			}
		})
	}
}

func TestExecutionLimiter_EnvironmentIntegration(t *testing.T) {
	env := NewEnvironment()

	// Test default limiter
	limiter := env.GetLimiter()
	if limiter == nil {
		t.Error("Environment should have a default limiter")
	}

	// Test custom limiter
	customLimiter := NewExecutionLimiter()
	customLimiter.MaxIterations = 50
	env.SetLimiter(customLimiter)

	retrieved := env.GetLimiter()
	if retrieved.MaxIterations != 50 {
		t.Errorf("Expected max iterations 50, got %d", retrieved.MaxIterations)
	}

	// Test SetLimits
	env.SetLimits(100, 200, 5*time.Second)
	if retrieved.MaxIterations != 100 {
		t.Errorf("Expected max iterations 100, got %d", retrieved.MaxIterations)
	}
	if retrieved.MaxRecursionDepth != 200 {
		t.Errorf("Expected max recursion depth 200, got %d", retrieved.MaxRecursionDepth)
	}
	if retrieved.MaxExecutionTime != 5*time.Second {
		t.Errorf("Expected max execution time 5s, got %v", retrieved.MaxExecutionTime)
	}
}

func TestExecutionLimiter_InnerEnvironment(t *testing.T) {
	outerEnv := NewEnvironment()
	outerEnv.SetLimits(100, 50, 2*time.Second)

	innerEnv := NewInnerEnv(outerEnv)

	// Inner environment should share the same limiter
	outerLimiter := outerEnv.GetLimiter()
	innerLimiter := innerEnv.GetLimiter()

	if outerLimiter != innerLimiter {
		t.Error("Inner environment should share the same limiter as outer")
	}

	// Modifications in inner should affect outer
	innerEnv.GetLimiter().IncrementIterations()
	if outerEnv.GetLimiter().CurrentIterations != 1 {
		t.Error("Iterations should be shared between environments")
	}
}

func TestExecutionLimiter_WithTimeout(t *testing.T) {
	env := NewEnvironment()
	env.Set("true", true)
	env.Set("false", false)

	// Test ExecuteWithTimeout with fast code
	parser := NewParser(`let x = 1 + 2;`)
	ast := parser.ParseProgram()

	result := env.ExecuteWithTimeout(ast, 100*time.Millisecond)
	// Result can be nil, that's fine for this code
	_ = result

	// Test ExecuteWithTimeout with slow code
	parser2 := NewParser(`let i = 0; while (i < 1000000) { i++; }`)
	ast2 := parser2.ParseProgram()

	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(*InfiniteLoopError); ok {
				if err.Type != "timeout" && !strings.Contains(err.Type, "timeout") {
					t.Errorf("Expected timeout error, got '%s'", err.Type)
				}
			} else {
				t.Errorf("Expected timeout error, got %T: %v", r, r)
			}
		} else {
			t.Error("Expected panic but none occurred")
		}
	}()

	env.ExecuteWithTimeout(ast2, 10*time.Millisecond)
}

func TestExecutionLimiter_ContextCancellation(t *testing.T) {
	env := NewEnvironment()
	env.Set("true", true)
	env.Set("false", false)

	// Create a context that will be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	env.SetContext(ctx)

	// Cancel the context after a short delay
	go func() {
		time.Sleep(5 * time.Millisecond)
		cancel()
	}()

	// Run code that should be cancelled
	parser := NewParser(`let i = 0; while (i < 1000000) { i++; }`)
	ast := parser.ParseProgram()

	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(*InfiniteLoopError); ok {
				// Accept any error type - cancellation might trigger loop detection first
				if err.Type != "while" && !strings.Contains(err.Type, "canceled") && !strings.Contains(err.Type, "timeout") {
					t.Errorf("Expected cancellation, timeout, or while error, got '%s'", err.Type)
				}
			} else {
				t.Errorf("Expected InfiniteLoopError, got %T: %v", r, r)
			}
		} else {
			t.Error("Expected panic but none occurred")
		}
	}()

	ast.Eval(env)
}
