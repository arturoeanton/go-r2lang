package reporters

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2test/core"
	"github.com/arturoeanton/go-r2lang/pkg/r2test/coverage"
)

// JSONReporter generates JSON coverage and test reports
type JSONReporter struct {
	OutputPath string
	Pretty     bool
}

// JSONReport represents the complete JSON report structure
type JSONReport struct {
	Metadata    JSONMetadata        `json:"metadata"`
	Coverage    *JSONCoverageReport `json:"coverage,omitempty"`
	TestResults *JSONTestResults    `json:"testResults,omitempty"`
}

// JSONMetadata contains report metadata
type JSONMetadata struct {
	GeneratedAt time.Time `json:"generatedAt"`
	Generator   string    `json:"generator"`
	Version     string    `json:"version"`
	Format      string    `json:"format"`
}

// JSONCoverageReport represents coverage data in JSON format
type JSONCoverageReport struct {
	Summary JSONCoverageSummary          `json:"summary"`
	Files   map[string]*JSONFileCoverage `json:"files"`
}

// JSONCoverageSummary represents aggregated coverage statistics
type JSONCoverageSummary struct {
	TotalFiles          int     `json:"totalFiles"`
	TotalLines          int     `json:"totalLines"`
	CoveredLines        int     `json:"coveredLines"`
	LinePercentage      float64 `json:"linePercentage"`
	TotalStatements     int     `json:"totalStatements"`
	CoveredStatements   int     `json:"coveredStatements"`
	StatementPercentage float64 `json:"statementPercentage"`
	TotalBranches       int     `json:"totalBranches"`
	CoveredBranches     int     `json:"coveredBranches"`
	BranchPercentage    float64 `json:"branchPercentage"`
	TotalFunctions      int     `json:"totalFunctions"`
	CoveredFunctions    int     `json:"coveredFunctions"`
	FunctionPercentage  float64 `json:"functionPercentage"`
	Duration            string  `json:"duration"`
}

// JSONFileCoverage represents coverage data for a single file
type JSONFileCoverage struct {
	Path         string                            `json:"path"`
	TotalLines   int                               `json:"totalLines"`
	CoveredLines int                               `json:"coveredLines"`
	Percentage   float64                           `json:"percentage"`
	Lines        map[string]*JSONLineCoverage      `json:"lines"`
	Statements   map[string]*JSONStatementCoverage `json:"statements,omitempty"`
	Branches     map[string]*JSONBranchCoverage    `json:"branches,omitempty"`
	Functions    map[string]*JSONFunctionCoverage  `json:"functions,omitempty"`
}

// JSONLineCoverage represents line coverage data
type JSONLineCoverage struct {
	LineNumber int    `json:"lineNumber"`
	Hits       int    `json:"hits"`
	IsHit      bool   `json:"isHit"`
	Type       string `json:"type"`
}

// JSONStatementCoverage represents statement coverage data
type JSONStatementCoverage struct {
	ID        int  `json:"id"`
	StartLine int  `json:"startLine"`
	EndLine   int  `json:"endLine"`
	StartCol  int  `json:"startCol"`
	EndCol    int  `json:"endCol"`
	Hits      int  `json:"hits"`
	IsHit     bool `json:"isHit"`
}

// JSONBranchCoverage represents branch coverage data
type JSONBranchCoverage struct {
	ID           int    `json:"id"`
	LineNumber   int    `json:"lineNumber"`
	BranchNumber int    `json:"branchNumber"`
	Taken        int    `json:"taken"`
	NotTaken     int    `json:"notTaken"`
	Type         string `json:"type"`
}

// JSONFunctionCoverage represents function coverage data
type JSONFunctionCoverage struct {
	Name      string `json:"name"`
	StartLine int    `json:"startLine"`
	EndLine   int    `json:"endLine"`
	Hits      int    `json:"hits"`
	IsHit     bool   `json:"isHit"`
}

// JSONTestResults represents test execution results
type JSONTestResults struct {
	Summary JSONTestSummary  `json:"summary"`
	Suites  []*JSONTestSuite `json:"suites"`
}

// JSONTestSummary represents test execution summary
type JSONTestSummary struct {
	Total     int    `json:"total"`
	Passed    int    `json:"passed"`
	Failed    int    `json:"failed"`
	Skipped   int    `json:"skipped"`
	Timeout   int    `json:"timeout"`
	Duration  string `json:"duration"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// JSONTestSuite represents a test suite
type JSONTestSuite struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Tests       []*JSONTestCase `json:"tests"`
	Tags        []string        `json:"tags,omitempty"`
}

// JSONTestCase represents a test case
type JSONTestCase struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	Duration    string   `json:"duration"`
	Error       string   `json:"error,omitempty"`
	Message     string   `json:"message,omitempty"`
	StartTime   string   `json:"startTime"`
	EndTime     string   `json:"endTime"`
	Tags        []string `json:"tags,omitempty"`
}

// NewJSONReporter creates a new JSON reporter
func NewJSONReporter(outputPath string) *JSONReporter {
	return &JSONReporter{
		OutputPath: outputPath,
		Pretty:     true,
	}
}

// SetPretty sets whether to format JSON output prettily
func (jr *JSONReporter) SetPretty(pretty bool) {
	jr.Pretty = pretty
}

// Generate generates a JSON report
func (jr *JSONReporter) Generate(collector *coverage.CoverageCollector, testResults *core.TestResults) error {
	report := &JSONReport{
		Metadata: JSONMetadata{
			GeneratedAt: time.Now(),
			Generator:   "R2Lang Test Framework",
			Version:     "1.0.0",
			Format:      "json",
		},
	}

	// Add coverage data if collector is provided
	if collector != nil {
		report.Coverage = jr.convertCoverageData(collector)
	}

	// Add test results if provided
	if testResults != nil {
		report.TestResults = jr.convertTestResults(testResults)
	}

	// Write to file
	return jr.writeJSONFile(report)
}

// convertCoverageData converts coverage data to JSON format
func (jr *JSONReporter) convertCoverageData(collector *coverage.CoverageCollector) *JSONCoverageReport {
	stats := collector.GetStats()
	files := collector.GetAllFiles()

	jsonCoverage := &JSONCoverageReport{
		Summary: JSONCoverageSummary{
			TotalFiles:          stats.TotalFiles,
			TotalLines:          stats.TotalLines,
			CoveredLines:        stats.CoveredLines,
			LinePercentage:      stats.LinePercentage,
			TotalStatements:     stats.TotalStatements,
			CoveredStatements:   stats.CoveredStatements,
			StatementPercentage: stats.StatementPercentage,
			TotalBranches:       stats.TotalBranches,
			CoveredBranches:     stats.CoveredBranches,
			BranchPercentage:    stats.BranchPercentage,
			TotalFunctions:      stats.TotalFunctions,
			CoveredFunctions:    stats.CoveredFunctions,
			FunctionPercentage:  stats.FunctionPercentage,
			Duration:            stats.Duration.String(),
		},
		Files: make(map[string]*JSONFileCoverage),
	}

	// Convert file coverage data
	for path, file := range files {
		jsonFile := &JSONFileCoverage{
			Path:         path,
			TotalLines:   file.TotalLines,
			CoveredLines: file.CoveredLines,
			Percentage:   file.GetFilePercentage(),
			Lines:        make(map[string]*JSONLineCoverage),
			Statements:   make(map[string]*JSONStatementCoverage),
			Branches:     make(map[string]*JSONBranchCoverage),
			Functions:    make(map[string]*JSONFunctionCoverage),
		}

		// Convert line coverage
		for lineNum, line := range file.Lines {
			jsonFile.Lines[fmt.Sprintf("%d", lineNum)] = &JSONLineCoverage{
				LineNumber: line.LineNumber,
				Hits:       line.Hits,
				IsHit:      line.IsHit,
				Type:       jr.lineTypeToString(line.Type),
			}
		}

		// Convert statement coverage
		for stmtID, stmt := range file.Statements {
			jsonFile.Statements[fmt.Sprintf("%d", stmtID)] = &JSONStatementCoverage{
				ID:        stmt.ID,
				StartLine: stmt.StartLine,
				EndLine:   stmt.EndLine,
				StartCol:  stmt.StartCol,
				EndCol:    stmt.EndCol,
				Hits:      stmt.Hits,
				IsHit:     stmt.IsHit,
			}
		}

		// Convert branch coverage
		for branchID, branch := range file.Branches {
			jsonFile.Branches[fmt.Sprintf("%d", branchID)] = &JSONBranchCoverage{
				ID:           branch.ID,
				LineNumber:   branch.LineNumber,
				BranchNumber: branch.BranchNumber,
				Taken:        branch.Taken,
				NotTaken:     branch.NotTaken,
				Type:         jr.branchTypeToString(branch.Type),
			}
		}

		// Convert function coverage
		for funcName, function := range file.Functions {
			jsonFile.Functions[funcName] = &JSONFunctionCoverage{
				Name:      function.Name,
				StartLine: function.StartLine,
				EndLine:   function.EndLine,
				Hits:      function.Hits,
				IsHit:     function.IsHit,
			}
		}

		jsonCoverage.Files[path] = jsonFile
	}

	return jsonCoverage
}

// convertTestResults converts test results to JSON format
func (jr *JSONReporter) convertTestResults(testResults *core.TestResults) *JSONTestResults {
	stats := testResults.GetStats()

	jsonResults := &JSONTestResults{
		Summary: JSONTestSummary{
			Total:     stats.Total,
			Passed:    stats.Passed,
			Failed:    stats.Failed,
			Skipped:   stats.Skipped,
			Timeout:   stats.Timeout,
			Duration:  testResults.Duration.String(),
			StartTime: testResults.StartTime.Format(time.RFC3339),
			EndTime:   testResults.EndTime.Format(time.RFC3339),
		},
		Suites: make([]*JSONTestSuite, 0),
	}

	// Group results by suite
	suiteMap := make(map[string]*JSONTestSuite)

	for _, result := range testResults.Results {
		suiteName := result.Suite.Name

		suite, exists := suiteMap[suiteName]
		if !exists {
			suite = &JSONTestSuite{
				Name:        result.Suite.Name,
				Description: result.Suite.Description,
				Tests:       make([]*JSONTestCase, 0),
				Tags:        result.Suite.Tags,
			}
			suiteMap[suiteName] = suite
		}

		testCase := &JSONTestCase{
			Name:        result.TestCase.Name,
			Description: result.TestCase.Description,
			Status:      result.Status.String(),
			Duration:    result.Duration.String(),
			StartTime:   result.StartTime.Format(time.RFC3339),
			EndTime:     result.EndTime.Format(time.RFC3339),
			Tags:        result.TestCase.Tags,
			Message:     result.Message,
		}

		if result.Error != nil {
			testCase.Error = result.Error.Error()
		}

		suite.Tests = append(suite.Tests, testCase)
	}

	// Convert map to slice
	for _, suite := range suiteMap {
		jsonResults.Suites = append(jsonResults.Suites, suite)
	}

	return jsonResults
}

// writeJSONFile writes the JSON report to file
func (jr *JSONReporter) writeJSONFile(report *JSONReport) error {
	file, err := os.Create(jr.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create JSON report file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if jr.Pretty {
		encoder.SetIndent("", "  ")
	}

	if err := encoder.Encode(report); err != nil {
		return fmt.Errorf("failed to encode JSON report: %w", err)
	}

	return nil
}

// lineTypeToString converts LineType to string
func (jr *JSONReporter) lineTypeToString(lineType coverage.LineType) string {
	switch lineType {
	case coverage.LineTypeCode:
		return "code"
	case coverage.LineTypeComment:
		return "comment"
	case coverage.LineTypeEmpty:
		return "empty"
	case coverage.LineTypeDeclaration:
		return "declaration"
	default:
		return "unknown"
	}
}

// branchTypeToString converts BranchType to string
func (jr *JSONReporter) branchTypeToString(branchType coverage.BranchType) string {
	switch branchType {
	case coverage.BranchTypeIf:
		return "if"
	case coverage.BranchTypeElse:
		return "else"
	case coverage.BranchTypeSwitch:
		return "switch"
	case coverage.BranchTypeCase:
		return "case"
	case coverage.BranchTypeLoop:
		return "loop"
	case coverage.BranchTypeTernary:
		return "ternary"
	default:
		return "unknown"
	}
}

// GenerateJSONReport generates a JSON report using the global collector
func GenerateJSONReport(outputPath string, testResults *core.TestResults) error {
	reporter := NewJSONReporter(outputPath)
	return reporter.Generate(coverage.GlobalCoverageCollector, testResults)
}

// GenerateCoverageOnlyJSON generates a JSON report with only coverage data
func GenerateCoverageOnlyJSON(outputPath string) error {
	reporter := NewJSONReporter(outputPath)
	return reporter.Generate(coverage.GlobalCoverageCollector, nil)
}

// GenerateTestOnlyJSON generates a JSON report with only test results
func GenerateTestOnlyJSON(outputPath string, testResults *core.TestResults) error {
	reporter := NewJSONReporter(outputPath)
	return reporter.Generate(nil, testResults)
}
