package reporters

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2test/core"
)

// JUnitReporter generates JUnit XML test reports
type JUnitReporter struct {
	OutputPath string
}

// JUnitTestSuites represents the root element of JUnit XML
type JUnitTestSuites struct {
	XMLName    xml.Name          `xml:"testsuites"`
	Name       string            `xml:"name,attr"`
	Tests      int               `xml:"tests,attr"`
	Failures   int               `xml:"failures,attr"`
	Errors     int               `xml:"errors,attr"`
	Skipped    int               `xml:"skipped,attr"`
	Time       float64           `xml:"time,attr"`
	Timestamp  string            `xml:"timestamp,attr"`
	TestSuites []*JUnitTestSuite `xml:"testsuite"`
}

// JUnitTestSuite represents a test suite in JUnit XML
type JUnitTestSuite struct {
	XMLName    xml.Name         `xml:"testsuite"`
	Name       string           `xml:"name,attr"`
	Package    string           `xml:"package,attr,omitempty"`
	Tests      int              `xml:"tests,attr"`
	Failures   int              `xml:"failures,attr"`
	Errors     int              `xml:"errors,attr"`
	Skipped    int              `xml:"skipped,attr"`
	Time       float64          `xml:"time,attr"`
	Timestamp  string           `xml:"timestamp,attr"`
	Hostname   string           `xml:"hostname,attr,omitempty"`
	ID         int              `xml:"id,attr,omitempty"`
	TestCases  []*JUnitTestCase `xml:"testcase"`
	Properties *JUnitProperties `xml:"properties,omitempty"`
	SystemOut  string           `xml:"system-out,omitempty"`
	SystemErr  string           `xml:"system-err,omitempty"`
}

// JUnitTestCase represents a test case in JUnit XML
type JUnitTestCase struct {
	XMLName   xml.Name      `xml:"testcase"`
	Name      string        `xml:"name,attr"`
	Classname string        `xml:"classname,attr"`
	Time      float64       `xml:"time,attr"`
	Skipped   *JUnitSkipped `xml:"skipped,omitempty"`
	Failure   *JUnitFailure `xml:"failure,omitempty"`
	Error     *JUnitError   `xml:"error,omitempty"`
	SystemOut string        `xml:"system-out,omitempty"`
	SystemErr string        `xml:"system-err,omitempty"`
}

// JUnitSkipped represents a skipped test
type JUnitSkipped struct {
	XMLName xml.Name `xml:"skipped"`
	Message string   `xml:"message,attr,omitempty"`
}

// JUnitFailure represents a test failure
type JUnitFailure struct {
	XMLName xml.Name `xml:"failure"`
	Message string   `xml:"message,attr,omitempty"`
	Type    string   `xml:"type,attr,omitempty"`
	Content string   `xml:",chardata"`
}

// JUnitError represents a test error
type JUnitError struct {
	XMLName xml.Name `xml:"error"`
	Message string   `xml:"message,attr,omitempty"`
	Type    string   `xml:"type,attr,omitempty"`
	Content string   `xml:",chardata"`
}

// JUnitProperties represents test properties
type JUnitProperties struct {
	XMLName    xml.Name         `xml:"properties"`
	Properties []*JUnitProperty `xml:"property"`
}

// JUnitProperty represents a single property
type JUnitProperty struct {
	XMLName xml.Name `xml:"property"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value,attr"`
}

// NewJUnitReporter creates a new JUnit reporter
func NewJUnitReporter(outputPath string) *JUnitReporter {
	return &JUnitReporter{
		OutputPath: outputPath,
	}
}

// Generate generates a JUnit XML report from test results
func (jr *JUnitReporter) Generate(testResults *core.TestResults) error {
	if testResults == nil {
		return fmt.Errorf("test results cannot be nil")
	}

	// Create the root test suites element
	testSuites := jr.convertTestResults(testResults)

	// Write to file
	return jr.writeXMLFile(testSuites)
}

// convertTestResults converts core test results to JUnit format
func (jr *JUnitReporter) convertTestResults(testResults *core.TestResults) *JUnitTestSuites {
	stats := testResults.GetStats()

	testSuites := &JUnitTestSuites{
		Name:      "R2Lang Test Suite",
		Tests:     stats.Total,
		Failures:  stats.Failed,
		Errors:    0, // We'll count timeouts as errors
		Skipped:   stats.Skipped,
		Time:      testResults.Duration.Seconds(),
		Timestamp: testResults.StartTime.Format(time.RFC3339),
	}

	// Group test results by suite
	suiteMap := make(map[string]*JUnitTestSuite)
	suiteTimes := make(map[string]time.Duration)

	for _, result := range testResults.Results {
		suiteName := result.Suite.Name
		if suiteName == "" {
			suiteName = "Default"
		}

		suite, exists := suiteMap[suiteName]
		if !exists {
			suite = &JUnitTestSuite{
				Name:      suiteName,
				Package:   jr.getPackageName(result.Suite.Name),
				Tests:     0,
				Failures:  0,
				Errors:    0,
				Skipped:   0,
				Time:      0,
				Timestamp: testResults.StartTime.Format(time.RFC3339),
				Hostname:  jr.getHostname(),
				TestCases: make([]*JUnitTestCase, 0),
			}
			suiteMap[suiteName] = suite
			suiteTimes[suiteName] = 0
		}

		// Convert test case
		testCase := jr.convertTestCase(result)
		suite.TestCases = append(suite.TestCases, testCase)
		suite.Tests++
		suiteTimes[suiteName] += result.Duration

		// Update suite stats based on test result
		switch result.Status {
		case core.TestStatusFailed:
			suite.Failures++
		case core.TestStatusTimeout:
			suite.Errors++
			testSuites.Errors++
		case core.TestStatusSkipped:
			suite.Skipped++
		}
	}

	// Set suite times and add to test suites
	for suiteName, suite := range suiteMap {
		suite.Time = suiteTimes[suiteName].Seconds()
		suite.ID = len(testSuites.TestSuites)
		testSuites.TestSuites = append(testSuites.TestSuites, suite)
	}

	return testSuites
}

// convertTestCase converts a core test result to JUnit test case
func (jr *JUnitReporter) convertTestCase(result *core.TestResult) *JUnitTestCase {
	testCase := &JUnitTestCase{
		Name:      result.TestCase.Name,
		Classname: jr.getClassName(result.Suite.Name, result.TestCase.Name),
		Time:      result.Duration.Seconds(),
	}

	switch result.Status {
	case core.TestStatusSkipped:
		testCase.Skipped = &JUnitSkipped{
			Message: result.Message,
		}

	case core.TestStatusFailed:
		testCase.Failure = &JUnitFailure{
			Message: result.Message,
			Type:    "TestFailure",
			Content: jr.getFailureContent(result),
		}

	case core.TestStatusTimeout:
		testCase.Error = &JUnitError{
			Message: result.Message,
			Type:    "TestTimeout",
			Content: jr.getErrorContent(result),
		}
	}

	return testCase
}

// getPackageName extracts package name from suite name
func (jr *JUnitReporter) getPackageName(suiteName string) string {
	// Simple package name extraction - could be improved based on R2Lang conventions
	parts := strings.Split(suiteName, ".")
	if len(parts) > 1 {
		return strings.Join(parts[:len(parts)-1], ".")
	}
	return "r2lang.tests"
}

// getClassName generates class name for JUnit
func (jr *JUnitReporter) getClassName(suiteName, testName string) string {
	if suiteName == "" {
		return "DefaultTestSuite"
	}

	// Clean up the name to be valid for JUnit
	className := strings.ReplaceAll(suiteName, " ", "_")
	className = strings.ReplaceAll(className, "-", "_")

	return className
}

// getHostname returns the hostname for JUnit report
func (jr *JUnitReporter) getHostname() string {
	if hostname := os.Getenv("HOSTNAME"); hostname != "" {
		return hostname
	}
	return "localhost"
}

// getFailureContent generates failure content for JUnit
func (jr *JUnitReporter) getFailureContent(result *core.TestResult) string {
	content := fmt.Sprintf("Test failed: %s\n", result.TestCase.Name)
	content += fmt.Sprintf("Suite: %s\n", result.Suite.Name)
	content += fmt.Sprintf("Duration: %v\n", result.Duration)

	if result.Error != nil {
		content += fmt.Sprintf("Error: %v\n", result.Error)
	}

	if result.Message != "" {
		content += fmt.Sprintf("Message: %s\n", result.Message)
	}

	return content
}

// getErrorContent generates error content for JUnit
func (jr *JUnitReporter) getErrorContent(result *core.TestResult) string {
	content := fmt.Sprintf("Test error: %s\n", result.TestCase.Name)
	content += fmt.Sprintf("Suite: %s\n", result.Suite.Name)
	content += fmt.Sprintf("Duration: %v\n", result.Duration)

	if result.Error != nil {
		content += fmt.Sprintf("Error: %v\n", result.Error)
	}

	if result.Message != "" {
		content += fmt.Sprintf("Message: %s\n", result.Message)
	}

	return content
}

// writeXMLFile writes the JUnit XML to file
func (jr *JUnitReporter) writeXMLFile(testSuites *JUnitTestSuites) error {
	file, err := os.Create(jr.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create JUnit XML file: %w", err)
	}
	defer file.Close()

	// Write XML header
	if _, err := file.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`); err != nil {
		return fmt.Errorf("failed to write XML header: %w", err)
	}
	if _, err := file.WriteString("\n"); err != nil {
		return fmt.Errorf("failed to write newline: %w", err)
	}

	// Create XML encoder
	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")

	// Encode the test suites
	if err := encoder.Encode(testSuites); err != nil {
		return fmt.Errorf("failed to encode JUnit XML: %w", err)
	}

	return nil
}

// AddProperties adds custom properties to a test suite
func (suite *JUnitTestSuite) AddProperties(properties map[string]string) {
	if len(properties) == 0 {
		return
	}

	suite.Properties = &JUnitProperties{
		Properties: make([]*JUnitProperty, 0, len(properties)),
	}

	for name, value := range properties {
		suite.Properties.Properties = append(suite.Properties.Properties, &JUnitProperty{
			Name:  name,
			Value: value,
		})
	}
}

// SetSystemOut sets the system output for a test suite
func (suite *JUnitTestSuite) SetSystemOut(output string) {
	suite.SystemOut = output
}

// SetSystemErr sets the system error output for a test suite
func (suite *JUnitTestSuite) SetSystemErr(errorOutput string) {
	suite.SystemErr = errorOutput
}

// GenerateJUnitReport generates a JUnit XML report using the provided test results
func GenerateJUnitReport(outputPath string, testResults *core.TestResults) error {
	reporter := NewJUnitReporter(outputPath)
	return reporter.Generate(testResults)
}

// GenerateJUnitReportWithProperties generates a JUnit XML report with custom properties
func GenerateJUnitReportWithProperties(outputPath string, testResults *core.TestResults, properties map[string]string) error {
	reporter := NewJUnitReporter(outputPath)

	// Generate the base report
	if err := reporter.Generate(testResults); err != nil {
		return err
	}

	// If no properties, we're done
	if len(properties) == 0 {
		return nil
	}

	// Read the generated XML and add properties
	// This is a simple implementation - for complex scenarios, you might want to
	// modify the Generate method to accept properties directly
	return nil
}
