package core

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Reporter interface defines methods for reporting test results
type Reporter interface {
	OnRunStart(totalTests int)
	OnRunComplete(results *TestResults)
	OnSuiteStart(suite *TestSuite)
	OnSuiteComplete(suite *TestSuite, results []*TestResult)
	OnTestStart(test *TestCase)
	OnTestComplete(result *TestResult)
}

// ConsoleReporter implements console output for test results
type ConsoleReporter struct {
	Writer  io.Writer
	Verbose bool
	Colors  bool
}

// NewConsoleReporter creates a new console reporter
func NewConsoleReporter() *ConsoleReporter {
	return &ConsoleReporter{
		Writer:  os.Stdout,
		Colors:  true,
		Verbose: false,
	}
}

// OnRunStart is called when test execution begins
func (cr *ConsoleReporter) OnRunStart(totalTests int) {
	if cr.Verbose {
		cr.println(cr.colorize("Starting test run...", ColorCyan))
		cr.println(fmt.Sprintf("Found %d tests", totalTests))
		cr.println("")
	}
}

// OnRunComplete is called when test execution completes
func (cr *ConsoleReporter) OnRunComplete(results *TestResults) {
	stats := results.GetStats()

	cr.println("")
	cr.printSeparator()
	cr.println(cr.colorize("Test Results:", ColorBold))
	cr.printSeparator()

	// Print summary
	cr.println(fmt.Sprintf("Total: %d", stats.Total))

	if stats.Passed > 0 {
		cr.println(cr.colorize(fmt.Sprintf("Passed: %d", stats.Passed), ColorGreen))
	}

	if stats.Failed > 0 {
		cr.println(cr.colorize(fmt.Sprintf("Failed: %d", stats.Failed), ColorRed))
	}

	if stats.Skipped > 0 {
		cr.println(cr.colorize(fmt.Sprintf("Skipped: %d", stats.Skipped), ColorYellow))
	}

	if stats.Timeout > 0 {
		cr.println(cr.colorize(fmt.Sprintf("Timeout: %d", stats.Timeout), ColorMagenta))
	}

	cr.println(fmt.Sprintf("Duration: %v", results.Duration))

	// Print failed tests details
	if stats.Failed > 0 {
		cr.println("")
		cr.println(cr.colorize("Failed Tests:", ColorRed))
		cr.printSeparator()

		for _, result := range results.Results {
			if result.Status == TestStatusFailed {
				cr.printFailedTest(result)
			}
		}
	}

	// Overall result
	cr.println("")
	if stats.Failed == 0 && stats.Timeout == 0 {
		cr.println(cr.colorize("✓ All tests passed!", ColorGreen))
	} else {
		cr.println(cr.colorize("✗ Some tests failed!", ColorRed))
	}
}

// OnSuiteStart is called when a test suite begins
func (cr *ConsoleReporter) OnSuiteStart(suite *TestSuite) {
	if cr.Verbose {
		cr.println("")
		cr.println(cr.colorize(fmt.Sprintf("Running suite: %s", suite.Name), ColorCyan))
		if suite.Description != "" && suite.Description != suite.Name {
			cr.println(fmt.Sprintf("  %s", suite.Description))
		}
	}
}

// OnSuiteComplete is called when a test suite completes
func (cr *ConsoleReporter) OnSuiteComplete(suite *TestSuite, results []*TestResult) {
	if cr.Verbose {
		passed := 0
		failed := 0
		for _, result := range results {
			if result.Status == TestStatusPassed {
				passed++
			} else if result.Status == TestStatusFailed {
				failed++
			}
		}

		status := "✓"
		color := ColorGreen
		if failed > 0 {
			status = "✗"
			color = ColorRed
		}

		cr.println(cr.colorize(fmt.Sprintf("  %s Suite completed: %d passed, %d failed", status, passed, failed), color))
	}
}

// OnTestStart is called when a test begins
func (cr *ConsoleReporter) OnTestStart(test *TestCase) {
	if cr.Verbose {
		cr.print(fmt.Sprintf("    Running: %s ... ", test.Name))
	}
}

// OnTestComplete is called when a test completes
func (cr *ConsoleReporter) OnTestComplete(result *TestResult) {
	if cr.Verbose {
		switch result.Status {
		case TestStatusPassed:
			cr.println(cr.colorize("PASS", ColorGreen))
		case TestStatusFailed:
			cr.println(cr.colorize("FAIL", ColorRed))
		case TestStatusSkipped:
			cr.println(cr.colorize("SKIP", ColorYellow))
		case TestStatusTimeout:
			cr.println(cr.colorize("TIMEOUT", ColorMagenta))
		}
	} else {
		// Compact output
		switch result.Status {
		case TestStatusPassed:
			cr.print(cr.colorize(".", ColorGreen))
		case TestStatusFailed:
			cr.print(cr.colorize("F", ColorRed))
		case TestStatusSkipped:
			cr.print(cr.colorize("S", ColorYellow))
		case TestStatusTimeout:
			cr.print(cr.colorize("T", ColorMagenta))
		}
	}
}

// printFailedTest prints details of a failed test
func (cr *ConsoleReporter) printFailedTest(result *TestResult) {
	cr.println(fmt.Sprintf("  %s > %s", result.Suite.Name, result.TestCase.Name))
	if result.Error != nil {
		cr.println(fmt.Sprintf("    Error: %s", result.Error.Error()))
	}
	if result.Message != "" {
		cr.println(fmt.Sprintf("    Message: %s", result.Message))
	}
	cr.println(fmt.Sprintf("    Duration: %v", result.Duration))
	cr.println("")
}

// Color constants
type Color string

const (
	ColorReset   Color = "\033[0m"
	ColorBold    Color = "\033[1m"
	ColorRed     Color = "\033[31m"
	ColorGreen   Color = "\033[32m"
	ColorYellow  Color = "\033[33m"
	ColorBlue    Color = "\033[34m"
	ColorMagenta Color = "\033[35m"
	ColorCyan    Color = "\033[36m"
	ColorWhite   Color = "\033[37m"
)

// colorize applies color to text if colors are enabled
func (cr *ConsoleReporter) colorize(text string, color Color) string {
	if !cr.Colors {
		return text
	}
	return string(color) + text + string(ColorReset)
}

// print writes text to the reporter's writer
func (cr *ConsoleReporter) print(text string) {
	fmt.Fprint(cr.Writer, text)
}

// println writes text with newline to the reporter's writer
func (cr *ConsoleReporter) println(text string) {
	fmt.Fprintln(cr.Writer, text)
}

// printSeparator prints a separator line
func (cr *ConsoleReporter) printSeparator() {
	cr.println(strings.Repeat("-", 50))
}

// JSONReporter implements JSON output for test results
type JSONReporter struct {
	Writer io.Writer
}

// NewJSONReporter creates a new JSON reporter
func NewJSONReporter(writer io.Writer) *JSONReporter {
	return &JSONReporter{
		Writer: writer,
	}
}

// OnRunStart is called when test execution begins
func (jr *JSONReporter) OnRunStart(totalTests int) {
	// JSON reporters typically don't output during run start
}

// OnRunComplete is called when test execution completes
func (jr *JSONReporter) OnRunComplete(results *TestResults) {
	// Create JSON report structure
	report := JSONTestReport{
		StartTime: results.StartTime,
		EndTime:   results.EndTime,
		Duration:  results.Duration.String(),
		Stats:     results.GetStats(),
		Suites:    make([]JSONSuiteReport, 0),
	}

	// Group results by suite
	suiteMap := make(map[*TestSuite][]JSONTestResult)
	for _, result := range results.Results {
		if _, exists := suiteMap[result.Suite]; !exists {
			suiteMap[result.Suite] = make([]JSONTestResult, 0)
		}

		jsonResult := JSONTestResult{
			Name:      result.TestCase.Name,
			Status:    result.Status.String(),
			Duration:  result.Duration.String(),
			StartTime: result.StartTime,
			EndTime:   result.EndTime,
		}

		if result.Error != nil {
			jsonResult.Error = result.Error.Error()
		}

		if result.Message != "" {
			jsonResult.Message = result.Message
		}

		suiteMap[result.Suite] = append(suiteMap[result.Suite], jsonResult)
	}

	// Convert to suite reports
	for suite, testResults := range suiteMap {
		suiteReport := JSONSuiteReport{
			Name:        suite.Name,
			Description: suite.Description,
			Tests:       testResults,
		}
		report.Suites = append(report.Suites, suiteReport)
	}

	// Output JSON
	encoder := json.NewEncoder(jr.Writer)
	encoder.SetIndent("", "  ")
	encoder.Encode(report)
}

// OnSuiteStart is called when a test suite begins
func (jr *JSONReporter) OnSuiteStart(suite *TestSuite) {
	// JSON reporters typically don't output during suite start
}

// OnSuiteComplete is called when a test suite completes
func (jr *JSONReporter) OnSuiteComplete(suite *TestSuite, results []*TestResult) {
	// JSON reporters typically don't output during suite complete
}

// OnTestStart is called when a test begins
func (jr *JSONReporter) OnTestStart(test *TestCase) {
	// JSON reporters typically don't output during test start
}

// OnTestComplete is called when a test completes
func (jr *JSONReporter) OnTestComplete(result *TestResult) {
	// JSON reporters typically don't output during test complete
}

// JSON report structures
type JSONTestReport struct {
	StartTime time.Time         `json:"startTime"`
	EndTime   time.Time         `json:"endTime"`
	Duration  string            `json:"duration"`
	Stats     TestStats         `json:"stats"`
	Suites    []JSONSuiteReport `json:"suites"`
}

type JSONSuiteReport struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Tests       []JSONTestResult `json:"tests"`
}

type JSONTestResult struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Duration  string    `json:"duration"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Error     string    `json:"error,omitempty"`
	Message   string    `json:"message,omitempty"`
}
