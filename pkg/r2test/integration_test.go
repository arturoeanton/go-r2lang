package r2test

import (
	"testing"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2test/assertions"
	"github.com/arturoeanton/go-r2lang/pkg/r2test/core"
)

func TestR2Test_Integration(t *testing.T) {
	// Create a new R2Test instance
	rt := NewWithDefaults()

	// Test data
	var testResults []string
	var setupCount, teardownCount int

	// Create a test suite with hooks
	rt.Describe("Integration Test Suite", func() {
		rt.BeforeAll(func() {
			testResults = append(testResults, "beforeAll")
		})

		rt.AfterAll(func() {
			testResults = append(testResults, "afterAll")
		})

		rt.BeforeEach(func() {
			setupCount++
			testResults = append(testResults, "beforeEach")
		})

		rt.AfterEach(func() {
			teardownCount++
			testResults = append(testResults, "afterEach")
		})

		rt.It("first test", func() {
			testResults = append(testResults, "test1")
		})

		rt.It("second test", func() {
			testResults = append(testResults, "test2")
		})
	})

	// Run the tests
	results, err := rt.Run()
	if err != nil {
		t.Fatalf("Failed to run tests: %v", err)
	}

	// Verify results
	stats := results.GetStats()
	if stats.Total != 2 {
		t.Errorf("Expected 2 tests, got %d", stats.Total)
	}

	if stats.Passed != 2 {
		t.Errorf("Expected 2 passed tests, got %d", stats.Passed)
	}

	if stats.Failed != 0 {
		t.Errorf("Expected 0 failed tests, got %d", stats.Failed)
	}

	// Verify hook execution order
	expectedOrder := []string{
		"beforeAll",
		"beforeEach", "test1", "afterEach",
		"beforeEach", "test2", "afterEach",
		"afterAll",
	}

	if len(testResults) != len(expectedOrder) {
		t.Errorf("Expected %d execution steps, got %d", len(expectedOrder), len(testResults))
	}

	for i, expected := range expectedOrder {
		if i >= len(testResults) || testResults[i] != expected {
			t.Errorf("Expected execution order %v, got %v", expectedOrder, testResults)
			break
		}
	}

	// Verify setup/teardown counts
	if setupCount != 2 {
		t.Errorf("Expected beforeEach to run 2 times, got %d", setupCount)
	}

	if teardownCount != 2 {
		t.Errorf("Expected afterEach to run 2 times, got %d", teardownCount)
	}
}

func TestR2Test_FailingTest(t *testing.T) {
	rt := NewWithDefaults()

	rt.Describe("Failing Test Suite", func() {
		rt.It("should pass", func() {
			// This test passes
		})

		rt.It("should fail", func() {
			panic("Test failure")
		})

		rt.It("should also pass", func() {
			// This test passes
		})
	})

	results, err := rt.Run()
	if err != nil {
		t.Fatalf("Failed to run tests: %v", err)
	}

	stats := results.GetStats()
	if stats.Total != 3 {
		t.Errorf("Expected 3 tests, got %d", stats.Total)
	}

	if stats.Passed != 2 {
		t.Errorf("Expected 2 passed tests, got %d", stats.Passed)
	}

	if stats.Failed != 1 {
		t.Errorf("Expected 1 failed test, got %d", stats.Failed)
	}
}

func TestR2Test_SkippedTest(t *testing.T) {
	rt := NewWithDefaults()

	rt.Describe("Skipped Test Suite", func() {
		rt.It("should run", func() {
			// This test runs
		})
	})

	// Add a skipped test manually
	if len(rt.runner.Suites) > 0 && len(rt.runner.Suites[0].Tests) > 0 {
		skippedTest := &core.TestCase{
			Name: "should be skipped",
			Skip: true,
			Func: func() {
				t.Error("Skipped test should not execute")
			},
			Suite: rt.runner.Suites[0],
		}
		rt.runner.Suites[0].Tests = append(rt.runner.Suites[0].Tests, skippedTest)
	}

	results, err := rt.Run()
	if err != nil {
		t.Fatalf("Failed to run tests: %v", err)
	}

	stats := results.GetStats()
	if stats.Total != 2 {
		t.Errorf("Expected 2 tests, got %d", stats.Total)
	}

	if stats.Passed != 1 {
		t.Errorf("Expected 1 passed test, got %d", stats.Passed)
	}

	if stats.Skipped != 1 {
		t.Errorf("Expected 1 skipped test, got %d", stats.Skipped)
	}
}

func TestR2Test_TimeoutTest(t *testing.T) {
	config := core.DefaultConfig()
	config.DefaultTimeout = 100 * time.Millisecond
	rt := New(config)

	rt.Describe("Timeout Test Suite", func() {
		rt.It("should timeout", func() {
			time.Sleep(200 * time.Millisecond) // Sleep longer than timeout
		})
	})

	results, err := rt.Run()
	if err != nil {
		t.Fatalf("Failed to run tests: %v", err)
	}

	stats := results.GetStats()
	if stats.Total != 1 {
		t.Errorf("Expected 1 test, got %d", stats.Total)
	}

	if stats.Timeout != 1 {
		t.Errorf("Expected 1 timeout test, got %d", stats.Timeout)
	}
}

func TestR2Test_NestedDescribe(t *testing.T) {
	rt := NewWithDefaults()

	var executionOrder []string

	rt.Describe("Outer Suite", func() {
		rt.BeforeEach(func() {
			executionOrder = append(executionOrder, "outer-beforeEach")
		})

		rt.It("outer test", func() {
			executionOrder = append(executionOrder, "outer-test")
		})

		rt.Describe("Inner Suite", func() {
			rt.BeforeEach(func() {
				executionOrder = append(executionOrder, "inner-beforeEach")
			})

			rt.It("inner test", func() {
				executionOrder = append(executionOrder, "inner-test")
			})
		})
	})

	results, err := rt.Run()
	if err != nil {
		t.Fatalf("Failed to run tests: %v", err)
	}

	stats := results.GetStats()
	if stats.Total != 2 {
		t.Errorf("Expected 2 tests, got %d", stats.Total)
	}

	if stats.Passed != 2 {
		t.Errorf("Expected 2 passed tests, got %d", stats.Passed)
	}

	// Note: This test demonstrates nested describe behavior
	// The actual implementation may vary based on how nesting is handled
}

func TestR2Test_Configuration(t *testing.T) {
	config := core.DefaultConfig()
	config.Verbose = true
	config.DefaultTimeout = 5 * time.Second
	config.Parallel = false

	rt := New(config)

	if rt.config.Verbose != true {
		t.Error("Configuration should be applied to R2Test instance")
	}

	if rt.config.DefaultTimeout != 5*time.Second {
		t.Error("Timeout configuration should be applied")
	}

	// Test setting configuration globally
	SetConfig(config)
	globalConfig := GetConfig()

	if globalConfig.Verbose != true {
		t.Error("Global configuration should be set correctly")
	}
}

func TestR2Test_EmptyTestSuite(t *testing.T) {
	rt := NewWithDefaults()

	rt.Describe("Empty Suite", func() {
		// No tests in this suite
	})

	results, err := rt.Run()
	if err != nil {
		t.Fatalf("Failed to run empty test suite: %v", err)
	}

	stats := results.GetStats()
	if stats.Total != 0 {
		t.Errorf("Expected 0 tests in empty suite, got %d", stats.Total)
	}
}

func TestGlobalFunctions(t *testing.T) {
	// Test global functions
	var executed bool

	Describe("Global Test Suite", func() {
		It("should use global functions", func() {
			executed = true
		})
	})

	results, err := RunTests()
	if err != nil {
		t.Fatalf("Failed to run global tests: %v", err)
	}

	if !executed {
		t.Error("Global test should have executed")
	}

	stats := results.GetStats()
	if stats.Total < 1 {
		t.Error("Global test should be counted")
	}
}

func TestAssertIntegration(t *testing.T) {
	// Test that assertions work with the framework
	assert := Assert("integration test")

	// These should not panic
	assert.Equals(5, 5)
	assert.True(true)
	assert.False(false)
	assert.Nil(nil)
	assert.NotNil("something")

	// Test that assertion errors are properly formatted
	defer func() {
		if r := recover(); r != nil {
			if ae, ok := r.(*assertions.AssertionError); ok {
				if ae.TestName != "integration test" {
					t.Errorf("Expected test name 'integration test', got '%s'", ae.TestName)
				}
			} else {
				t.Errorf("Expected AssertionError, got %T", r)
			}
		} else {
			t.Error("Expected assertion to fail")
		}
	}()

	// This should panic
	assert.Equals(5, 10)
}
