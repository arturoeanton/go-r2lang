package core

import (
	"testing"
	"time"
)

func TestTestRunner_NewTestRunner(t *testing.T) {
	config := DefaultConfig()
	runner := NewTestRunner(config)

	if runner == nil {
		t.Error("NewTestRunner should return a valid runner")
	}

	if runner.Config != config {
		t.Error("Runner should use the provided config")
	}

	if len(runner.Suites) != 0 {
		t.Error("New runner should have no suites initially")
	}
}

func TestTestRunner_AddSuite(t *testing.T) {
	config := DefaultConfig()
	runner := NewTestRunner(config)

	suite := &TestSuite{
		Name:  "Test Suite",
		Tests: make([]*TestCase, 0),
	}

	runner.AddSuite(suite)

	if len(runner.Suites) != 1 {
		t.Error("Runner should have one suite after adding")
	}

	if runner.Suites[0] != suite {
		t.Error("Runner should contain the added suite")
	}
}

func TestTestRunner_RunEmptySuite(t *testing.T) {
	config := DefaultConfig()
	runner := NewTestRunner(config)

	suite := &TestSuite{
		Name:  "Empty Suite",
		Tests: make([]*TestCase, 0),
	}

	runner.AddSuite(suite)

	results, err := runner.Run()
	if err != nil {
		t.Errorf("Running empty suite should not error: %v", err)
	}

	if results == nil {
		t.Error("Results should not be nil")
	}

	stats := results.GetStats()
	if stats.Total != 0 {
		t.Errorf("Empty suite should have 0 total tests, got %d", stats.Total)
	}
}

func TestTestRunner_RunPassingTest(t *testing.T) {
	config := DefaultConfig()
	runner := NewTestRunner(config)

	testPassed := false
	test := &TestCase{
		Name: "Passing Test",
		Func: func() {
			testPassed = true
		},
	}

	suite := &TestSuite{
		Name:  "Test Suite",
		Tests: []*TestCase{test},
	}
	test.Suite = suite

	runner.AddSuite(suite)

	results, err := runner.Run()
	if err != nil {
		t.Errorf("Running passing test should not error: %v", err)
	}

	if !testPassed {
		t.Error("Test function should have been executed")
	}

	stats := results.GetStats()
	if stats.Total != 1 {
		t.Errorf("Expected 1 total test, got %d", stats.Total)
	}

	if stats.Passed != 1 {
		t.Errorf("Expected 1 passed test, got %d", stats.Passed)
	}

	if stats.Failed != 0 {
		t.Errorf("Expected 0 failed tests, got %d", stats.Failed)
	}
}

func TestTestRunner_RunFailingTest(t *testing.T) {
	config := DefaultConfig()
	runner := NewTestRunner(config)

	test := &TestCase{
		Name: "Failing Test",
		Func: func() {
			panic("Test failed")
		},
	}

	suite := &TestSuite{
		Name:  "Test Suite",
		Tests: []*TestCase{test},
	}
	test.Suite = suite

	runner.AddSuite(suite)

	results, err := runner.Run()
	if err != nil {
		t.Errorf("Running failing test should not error: %v", err)
	}

	stats := results.GetStats()
	if stats.Total != 1 {
		t.Errorf("Expected 1 total test, got %d", stats.Total)
	}

	if stats.Passed != 0 {
		t.Errorf("Expected 0 passed tests, got %d", stats.Passed)
	}

	if stats.Failed != 1 {
		t.Errorf("Expected 1 failed test, got %d", stats.Failed)
	}
}

func TestTestRunner_RunSkippedTest(t *testing.T) {
	config := DefaultConfig()
	runner := NewTestRunner(config)

	test := &TestCase{
		Name: "Skipped Test",
		Skip: true,
		Func: func() {
			t.Error("Skipped test should not execute")
		},
	}

	suite := &TestSuite{
		Name:  "Test Suite",
		Tests: []*TestCase{test},
	}
	test.Suite = suite

	runner.AddSuite(suite)

	results, err := runner.Run()
	if err != nil {
		t.Errorf("Running skipped test should not error: %v", err)
	}

	stats := results.GetStats()
	if stats.Total != 1 {
		t.Errorf("Expected 1 total test, got %d", stats.Total)
	}

	if stats.Skipped != 1 {
		t.Errorf("Expected 1 skipped test, got %d", stats.Skipped)
	}
}

func TestTestRunner_RunTestWithTimeout(t *testing.T) {
	config := DefaultConfig()
	config.DefaultTimeout = 100 * time.Millisecond
	runner := NewTestRunner(config)

	test := &TestCase{
		Name: "Timeout Test",
		Func: func() {
			time.Sleep(200 * time.Millisecond) // Sleep longer than timeout
		},
	}

	suite := &TestSuite{
		Name:  "Test Suite",
		Tests: []*TestCase{test},
	}
	test.Suite = suite

	runner.AddSuite(suite)

	results, err := runner.Run()
	if err != nil {
		t.Errorf("Running timeout test should not error: %v", err)
	}

	stats := results.GetStats()
	if stats.Timeout != 1 {
		t.Errorf("Expected 1 timeout test, got %d", stats.Timeout)
	}
}

func TestTestRunner_BeforeEachAfterEach(t *testing.T) {
	config := DefaultConfig()
	runner := NewTestRunner(config)

	var execOrder []string

	test1 := &TestCase{
		Name: "Test 1",
		Func: func() {
			execOrder = append(execOrder, "test1")
		},
	}

	test2 := &TestCase{
		Name: "Test 2",
		Func: func() {
			execOrder = append(execOrder, "test2")
		},
	}

	suite := &TestSuite{
		Name:  "Test Suite",
		Tests: []*TestCase{test1, test2},
		BeforeEach: func() {
			execOrder = append(execOrder, "beforeEach")
		},
		AfterEach: func() {
			execOrder = append(execOrder, "afterEach")
		},
	}
	test1.Suite = suite
	test2.Suite = suite

	runner.AddSuite(suite)

	_, err := runner.Run()
	if err != nil {
		t.Errorf("Running tests with hooks should not error: %v", err)
	}

	expected := []string{
		"beforeEach", "test1", "afterEach",
		"beforeEach", "test2", "afterEach",
	}

	if len(execOrder) != len(expected) {
		t.Errorf("Expected %d executions, got %d", len(expected), len(execOrder))
	}

	for i, exp := range expected {
		if i >= len(execOrder) || execOrder[i] != exp {
			t.Errorf("Expected execution order %v, got %v", expected, execOrder)
			break
		}
	}
}

func TestTestStatus_String(t *testing.T) {
	tests := []struct {
		status   TestStatus
		expected string
	}{
		{TestStatusPending, "pending"},
		{TestStatusRunning, "running"},
		{TestStatusPassed, "passed"},
		{TestStatusFailed, "failed"},
		{TestStatusSkipped, "skipped"},
		{TestStatusTimeout, "timeout"},
		{TestStatus(999), "unknown"},
	}

	for _, test := range tests {
		result := test.status.String()
		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}
}
