package r2test

import (
	"fmt"

	"github.com/arturoeanton/go-r2lang/pkg/r2test/assertions"
	"github.com/arturoeanton/go-r2lang/pkg/r2test/core"
)

// R2Test provides the main testing interface
type R2Test struct {
	config       *core.TestConfig
	runner       *core.TestRunner
	discovery    *core.TestDiscovery
	currentSuite *core.TestSuite
}

// New creates a new R2Test instance
func New(config *core.TestConfig) *R2Test {
	if config == nil {
		config = core.DefaultConfig()
	}

	return &R2Test{
		config:    config,
		runner:    core.NewTestRunner(config),
		discovery: core.NewTestDiscovery(config),
	}
}

// NewWithDefaults creates a new R2Test instance with default configuration
func NewWithDefaults() *R2Test {
	return New(core.DefaultConfig())
}

// Describe creates a new test suite
func (rt *R2Test) Describe(name string, fn func()) *R2Test {
	suite := &core.TestSuite{
		Name:        name,
		Description: name,
		Tests:       make([]*core.TestCase, 0),
		Tags:        make([]string, 0),
	}

	// Set current suite context
	previousSuite := rt.currentSuite
	rt.currentSuite = suite

	// Execute the describe function to collect tests
	fn()

	// Restore previous suite context
	rt.currentSuite = previousSuite

	// Add suite to runner
	rt.runner.AddSuite(suite)

	return rt
}

// It creates a new test case within the current suite
func (rt *R2Test) It(name string, fn func()) *R2Test {
	if rt.currentSuite == nil {
		panic("It() must be called within a Describe() block")
	}

	test := &core.TestCase{
		Name:        name,
		Description: name,
		Func:        fn,
		Tags:        make([]string, 0),
		Suite:       rt.currentSuite,
	}

	rt.currentSuite.Tests = append(rt.currentSuite.Tests, test)

	return rt
}

// BeforeEach sets a function to run before each test in the current suite
func (rt *R2Test) BeforeEach(fn func()) *R2Test {
	if rt.currentSuite == nil {
		panic("BeforeEach() must be called within a Describe() block")
	}

	rt.currentSuite.BeforeEach = fn
	return rt
}

// AfterEach sets a function to run after each test in the current suite
func (rt *R2Test) AfterEach(fn func()) *R2Test {
	if rt.currentSuite == nil {
		panic("AfterEach() must be called within a Describe() block")
	}

	rt.currentSuite.AfterEach = fn
	return rt
}

// BeforeAll sets a function to run before all tests in the current suite
func (rt *R2Test) BeforeAll(fn func()) *R2Test {
	if rt.currentSuite == nil {
		panic("BeforeAll() must be called within a Describe() block")
	}

	rt.currentSuite.BeforeAll = fn
	return rt
}

// AfterAll sets a function to run after all tests in the current suite
func (rt *R2Test) AfterAll(fn func()) *R2Test {
	if rt.currentSuite == nil {
		panic("AfterAll() must be called within a Describe() block")
	}

	rt.currentSuite.AfterAll = fn
	return rt
}

// Run executes all tests and returns results
func (rt *R2Test) Run() (*core.TestResults, error) {
	return rt.runner.Run()
}

// RunDiscoveredTests discovers and runs all tests in configured directories
func (rt *R2Test) RunDiscoveredTests() (*core.TestResults, error) {
	suites, err := rt.discovery.LoadTestSuites()
	if err != nil {
		return nil, fmt.Errorf("failed to discover tests: %w", err)
	}

	for _, suite := range suites {
		rt.runner.AddSuite(suite)
	}

	return rt.runner.Run()
}

// Global testing functions for direct use in R2Lang

var globalR2Test *R2Test

// init initializes the global R2Test instance
func init() {
	globalR2Test = NewWithDefaults()
}

// Describe creates a test suite (global function)
func Describe(name string, fn func()) {
	globalR2Test.Describe(name, fn)
}

// It creates a test case (global function)
func It(name string, fn func()) {
	globalR2Test.It(name, fn)
}

// BeforeEach sets up before each test (global function)
func BeforeEach(fn func()) {
	globalR2Test.BeforeEach(fn)
}

// AfterEach cleans up after each test (global function)
func AfterEach(fn func()) {
	globalR2Test.AfterEach(fn)
}

// BeforeAll sets up before all tests (global function)
func BeforeAll(fn func()) {
	globalR2Test.BeforeAll(fn)
}

// AfterAll cleans up after all tests (global function)
func AfterAll(fn func()) {
	globalR2Test.AfterAll(fn)
}

// Assert creates a new assertion context (global function)
func Assert(testName string) *assertions.Assert {
	return assertions.NewAssert(testName)
}

// RunTests runs all global tests
func RunTests() (*core.TestResults, error) {
	return globalR2Test.Run()
}

// RunDiscoveredTests discovers and runs tests
func RunDiscoveredTests() (*core.TestResults, error) {
	return globalR2Test.RunDiscoveredTests()
}

// SetConfig sets the global configuration
func SetConfig(config *core.TestConfig) {
	globalR2Test.config = config
	globalR2Test.runner = core.NewTestRunner(config)
	globalR2Test.discovery = core.NewTestDiscovery(config)
}

// GetConfig returns the current configuration
func GetConfig() *core.TestConfig {
	return globalR2Test.config
}
