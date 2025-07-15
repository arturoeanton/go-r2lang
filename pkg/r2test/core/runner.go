package core

import (
	"fmt"
	"time"
)

// TestRunner executes test suites and manages test execution
type TestRunner struct {
	Suites   []*TestSuite
	Reporter Reporter
	Config   *TestConfig
}

// TestSuite represents a collection of related tests
type TestSuite struct {
	Name        string
	Description string
	Tests       []*TestCase
	BeforeEach  func()
	AfterEach   func()
	BeforeAll   func()
	AfterAll    func()
	Tags        []string
	Environment interface{} // Will be *r2core.Environment when integrated
}

// TestCase represents an individual test
type TestCase struct {
	Name        string
	Description string
	Func        func()
	Tags        []string
	Timeout     time.Duration
	Skip        bool
	Only        bool
	Suite       *TestSuite
}

// TestResult represents the result of a test execution
type TestResult struct {
	TestCase  *TestCase
	Suite     *TestSuite
	Status    TestStatus
	Duration  time.Duration
	Error     error
	Message   string
	StartTime time.Time
	EndTime   time.Time
}

// TestStatus represents the status of a test
type TestStatus int

const (
	TestStatusPending TestStatus = iota
	TestStatusRunning
	TestStatusPassed
	TestStatusFailed
	TestStatusSkipped
	TestStatusTimeout
)

func (s TestStatus) String() string {
	switch s {
	case TestStatusPending:
		return "pending"
	case TestStatusRunning:
		return "running"
	case TestStatusPassed:
		return "passed"
	case TestStatusFailed:
		return "failed"
	case TestStatusSkipped:
		return "skipped"
	case TestStatusTimeout:
		return "timeout"
	default:
		return "unknown"
	}
}

// NewTestRunner creates a new test runner
func NewTestRunner(config *TestConfig) *TestRunner {
	return &TestRunner{
		Suites:   make([]*TestSuite, 0),
		Config:   config,
		Reporter: NewConsoleReporter(),
	}
}

// AddSuite adds a test suite to the runner
func (tr *TestRunner) AddSuite(suite *TestSuite) {
	tr.Suites = append(tr.Suites, suite)
}

// Run executes all test suites
func (tr *TestRunner) Run() (*TestResults, error) {
	results := &TestResults{
		StartTime: time.Now(),
		Results:   make([]*TestResult, 0),
	}

	tr.Reporter.OnRunStart(len(tr.getAllTests()))

	for _, suite := range tr.Suites {
		suiteResults := tr.runSuite(suite)
		results.Results = append(results.Results, suiteResults...)
	}

	results.EndTime = time.Now()
	results.Duration = results.EndTime.Sub(results.StartTime)

	tr.Reporter.OnRunComplete(results)

	return results, nil
}

// runSuite executes a single test suite
func (tr *TestRunner) runSuite(suite *TestSuite) []*TestResult {
	var results []*TestResult

	tr.Reporter.OnSuiteStart(suite)

	// Run BeforeAll hook
	if suite.BeforeAll != nil {
		suite.BeforeAll()
	}

	for _, test := range suite.Tests {
		if tr.shouldSkipTest(test) {
			result := &TestResult{
				TestCase:  test,
				Suite:     suite,
				Status:    TestStatusSkipped,
				StartTime: time.Now(),
				EndTime:   time.Now(),
				Message:   "Test skipped",
			}
			results = append(results, result)
			tr.Reporter.OnTestComplete(result)
			continue
		}

		result := tr.runTest(test, suite)
		results = append(results, result)
	}

	// Run AfterAll hook
	if suite.AfterAll != nil {
		suite.AfterAll()
	}

	tr.Reporter.OnSuiteComplete(suite, results)

	return results
}

// runTest executes a single test
func (tr *TestRunner) runTest(test *TestCase, suite *TestSuite) *TestResult {
	result := &TestResult{
		TestCase:  test,
		Suite:     suite,
		Status:    TestStatusRunning,
		StartTime: time.Now(),
	}

	tr.Reporter.OnTestStart(test)

	// Run BeforeEach hook
	if suite.BeforeEach != nil {
		suite.BeforeEach()
	}

	// Execute test with timeout
	done := make(chan bool, 1)
	var testError error

	go func() {
		defer func() {
			if r := recover(); r != nil {
				testError = fmt.Errorf("test panicked: %v", r)
			}
			done <- true
		}()

		test.Func()
	}()

	timeout := test.Timeout
	if timeout == 0 {
		timeout = tr.Config.DefaultTimeout
	}

	select {
	case <-done:
		if testError != nil {
			result.Status = TestStatusFailed
			result.Error = testError
			result.Message = testError.Error()
		} else {
			result.Status = TestStatusPassed
			result.Message = "Test passed"
		}
	case <-time.After(timeout):
		result.Status = TestStatusTimeout
		result.Error = fmt.Errorf("test timed out after %v", timeout)
		result.Message = fmt.Sprintf("Test timed out after %v", timeout)
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	// Run AfterEach hook
	if suite.AfterEach != nil {
		suite.AfterEach()
	}

	tr.Reporter.OnTestComplete(result)

	return result
}

// shouldSkipTest determines if a test should be skipped
func (tr *TestRunner) shouldSkipTest(test *TestCase) bool {
	if test.Skip {
		return true
	}

	// Check if there are "only" tests and this isn't one of them
	if tr.hasOnlyTests() && !test.Only {
		return true
	}

	// Check tag filtering
	if len(tr.Config.FilterTags) > 0 && !tr.hasMatchingTags(test.Tags, tr.Config.FilterTags) {
		return true
	}

	return false
}

// hasOnlyTests checks if there are any tests marked with "only"
func (tr *TestRunner) hasOnlyTests() bool {
	for _, suite := range tr.Suites {
		for _, test := range suite.Tests {
			if test.Only {
				return true
			}
		}
	}
	return false
}

// hasMatchingTags checks if test tags match filter tags
func (tr *TestRunner) hasMatchingTags(testTags, filterTags []string) bool {
	for _, filterTag := range filterTags {
		for _, testTag := range testTags {
			if testTag == filterTag {
				return true
			}
		}
	}
	return false
}

// getAllTests returns all tests across all suites
func (tr *TestRunner) getAllTests() []*TestCase {
	var allTests []*TestCase
	for _, suite := range tr.Suites {
		allTests = append(allTests, suite.Tests...)
	}
	return allTests
}

// TestResults aggregates all test results
type TestResults struct {
	Results   []*TestResult
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
}

// GetStats returns test execution statistics
func (tr *TestResults) GetStats() TestStats {
	stats := TestStats{}

	for _, result := range tr.Results {
		stats.Total++
		switch result.Status {
		case TestStatusPassed:
			stats.Passed++
		case TestStatusFailed:
			stats.Failed++
		case TestStatusSkipped:
			stats.Skipped++
		case TestStatusTimeout:
			stats.Timeout++
		}
	}

	return stats
}

// TestStats holds test execution statistics
type TestStats struct {
	Total   int
	Passed  int
	Failed  int
	Skipped int
	Timeout int
}
