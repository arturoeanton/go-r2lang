package reporters

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2test/core"
	"github.com/arturoeanton/go-r2lang/pkg/r2test/coverage"
)

// HTMLReporter generates HTML coverage reports
type HTMLReporter struct {
	OutputDir     string
	TemplateDir   string
	BasePath      string
	ReportTitle   string
	IncludeSource bool
}

// HTMLReportData represents data for HTML template
type HTMLReportData struct {
	Title          string
	GeneratedAt    time.Time
	Summary        *coverage.CoverageStats
	Files          []*HTMLFileData
	TestResults    *core.TestResults
	HasTestResults bool
	CSS            template.CSS
	JavaScript     template.JS
}

// HTMLFileData represents file data for HTML template
type HTMLFileData struct {
	*coverage.FileCoverage
	RelativePath string
	Percentage   float64
	StatusClass  string
	Lines        []*HTMLLineData
}

// HTMLLineData represents line data for HTML template
type HTMLLineData struct {
	*coverage.LineCoverage
	CSSClass string
	Content  string
}

// NewHTMLReporter creates a new HTML reporter
func NewHTMLReporter(outputDir string) *HTMLReporter {
	return &HTMLReporter{
		OutputDir:     outputDir,
		BasePath:      ".",
		ReportTitle:   "R2Lang Test Coverage Report",
		IncludeSource: true,
	}
}

// Generate generates an HTML coverage report
func (hr *HTMLReporter) Generate(collector *coverage.CoverageCollector, testResults *core.TestResults) error {
	// Create output directory
	if err := os.MkdirAll(hr.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Prepare report data
	data := &HTMLReportData{
		Title:          hr.ReportTitle,
		GeneratedAt:    time.Now(),
		Summary:        collector.GetStats(),
		TestResults:    testResults,
		HasTestResults: testResults != nil,
		CSS:            template.CSS(hr.getCSS()),
		JavaScript:     template.JS(hr.getJavaScript()),
	}

	// Prepare file data
	files := collector.GetSortedFiles()
	data.Files = make([]*HTMLFileData, len(files))

	for i, file := range files {
		fileData := &HTMLFileData{
			FileCoverage: file,
			RelativePath: file.GetRelativePath(hr.BasePath),
			Percentage:   file.GetFilePercentage(),
		}

		// Set status class based on coverage
		fileData.StatusClass = hr.getStatusClass(fileData.Percentage)

		// Prepare line data if including source
		if hr.IncludeSource {
			fileData.Lines = hr.prepareLineData(file)
		}

		data.Files[i] = fileData
	}

	// Generate main report
	if err := hr.generateMainReport(data); err != nil {
		return fmt.Errorf("failed to generate main report: %w", err)
	}

	// Generate individual file reports
	if hr.IncludeSource {
		for _, file := range data.Files {
			if err := hr.generateFileReport(file, data); err != nil {
				return fmt.Errorf("failed to generate file report for %s: %w", file.Path, err)
			}
		}
	}

	return nil
}

// generateMainReport generates the main HTML report
func (hr *HTMLReporter) generateMainReport(data *HTMLReportData) error {
	tmpl := template.Must(template.New("main").Parse(hr.getMainTemplate()))

	outputFile := filepath.Join(hr.OutputDir, "index.html")
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create main report file: %w", err)
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// generateFileReport generates an HTML report for a specific file
func (hr *HTMLReporter) generateFileReport(file *HTMLFileData, globalData *HTMLReportData) error {
	tmpl := template.Must(template.New("file").Parse(hr.getFileTemplate()))

	// Create file-specific data
	fileData := struct {
		*HTMLReportData
		CurrentFile *HTMLFileData
	}{
		HTMLReportData: globalData,
		CurrentFile:    file,
	}

	// Create safe filename
	safeFilename := filepath.Base(file.Path) + ".html"
	outputFile := filepath.Join(hr.OutputDir, safeFilename)

	fileHandle, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create file report: %w", err)
	}
	defer fileHandle.Close()

	return tmpl.Execute(fileHandle, fileData)
}

// prepareLineData prepares line data with CSS classes and content
func (hr *HTMLReporter) prepareLineData(file *coverage.FileCoverage) []*HTMLLineData {
	// Get all line numbers
	var lineNumbers []int
	for lineNum := range file.Lines {
		lineNumbers = append(lineNumbers, lineNum)
	}
	sort.Ints(lineNumbers)

	lines := make([]*HTMLLineData, len(lineNumbers))
	for i, lineNum := range lineNumbers {
		line := file.Lines[lineNum]

		htmlLine := &HTMLLineData{
			LineCoverage: line,
			Content:      line.Source,
		}

		// Set CSS class based on coverage
		switch {
		case line.Type != coverage.LineTypeCode:
			htmlLine.CSSClass = "line-ignored"
		case line.IsHit:
			htmlLine.CSSClass = "line-covered"
		default:
			htmlLine.CSSClass = "line-uncovered"
		}

		lines[i] = htmlLine
	}

	return lines
}

// getStatusClass returns CSS class based on coverage percentage
func (hr *HTMLReporter) getStatusClass(percentage float64) string {
	switch {
	case percentage >= 90:
		return "status-high"
	case percentage >= 70:
		return "status-medium"
	case percentage >= 50:
		return "status-low"
	default:
		return "status-critical"
	}
}

// getMainTemplate returns the main HTML template
func (hr *HTMLReporter) getMainTemplate() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>{{.CSS}}</style>
</head>
<body>
    <header>
        <h1>{{.Title}}</h1>
        <p>Generated on {{.GeneratedAt.Format "2006-01-02 15:04:05"}}</p>
    </header>

    <main>
        <section class="summary">
            <h2>Coverage Summary</h2>
            <div class="metrics">
                <div class="metric">
                    <div class="metric-value">{{printf "%.1f%%" .Summary.LinePercentage}}</div>
                    <div class="metric-label">Line Coverage</div>
                    <div class="metric-detail">{{.Summary.CoveredLines}}/{{.Summary.TotalLines}} lines</div>
                </div>
                <div class="metric">
                    <div class="metric-value">{{printf "%.1f%%" .Summary.StatementPercentage}}</div>
                    <div class="metric-label">Statement Coverage</div>
                    <div class="metric-detail">{{.Summary.CoveredStatements}}/{{.Summary.TotalStatements}} statements</div>
                </div>
                <div class="metric">
                    <div class="metric-value">{{printf "%.1f%%" .Summary.BranchPercentage}}</div>
                    <div class="metric-label">Branch Coverage</div>
                    <div class="metric-detail">{{.Summary.CoveredBranches}}/{{.Summary.TotalBranches}} branches</div>
                </div>
                <div class="metric">
                    <div class="metric-value">{{printf "%.1f%%" .Summary.FunctionPercentage}}</div>
                    <div class="metric-label">Function Coverage</div>
                    <div class="metric-detail">{{.Summary.CoveredFunctions}}/{{.Summary.TotalFunctions}} functions</div>
                </div>
            </div>
        </section>

        {{if .HasTestResults}}
        <section class="test-results">
            <h2>Test Results</h2>
            {{$stats := .TestResults.GetStats}}
            <div class="test-summary">
                <span class="test-stat">Total: {{$stats.Total}}</span>
                <span class="test-stat passed">Passed: {{$stats.Passed}}</span>
                <span class="test-stat failed">Failed: {{$stats.Failed}}</span>
                <span class="test-stat skipped">Skipped: {{$stats.Skipped}}</span>
                <span class="test-stat timeout">Timeout: {{$stats.Timeout}}</span>
            </div>
        </section>
        {{end}}

        <section class="files">
            <h2>File Coverage</h2>
            <table class="file-table">
                <thead>
                    <tr>
                        <th>File</th>
                        <th>Coverage</th>
                        <th>Lines</th>
                        <th>Uncovered Lines</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .Files}}
                    <tr class="{{.StatusClass}}">
                        <td class="file-name">
                            <a href="{{.RelativePath}}.html">{{.RelativePath}}</a>
                        </td>
                        <td class="coverage-percentage">{{printf "%.1f%%" .Percentage}}</td>
                        <td class="line-count">{{.CoveredLines}}/{{.TotalLines}}</td>
                        <td class="uncovered-count">{{sub .TotalLines .CoveredLines}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </section>
    </main>

    <footer>
        <p>Generated by R2Lang Test Framework</p>
    </footer>

    <script>{{.JavaScript}}</script>
</body>
</html>`
}

// getFileTemplate returns the file-specific HTML template
func (hr *HTMLReporter) getFileTemplate() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - {{.CurrentFile.RelativePath}}</title>
    <style>{{.CSS}}</style>
</head>
<body>
    <header>
        <h1><a href="index.html">{{.Title}}</a></h1>
        <h2>{{.CurrentFile.RelativePath}}</h2>
        <p>Coverage: {{printf "%.1f%%" .CurrentFile.Percentage}} ({{.CurrentFile.CoveredLines}}/{{.CurrentFile.TotalLines}} lines)</p>
    </header>

    <main>
        <section class="source-code">
            <table class="source-table">
                <thead>
                    <tr>
                        <th class="line-number">Line</th>
                        <th class="line-hits">Hits</th>
                        <th class="line-source">Source</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .CurrentFile.Lines}}
                    <tr class="{{.CSSClass}}">
                        <td class="line-number">{{.LineNumber}}</td>
                        <td class="line-hits">{{if .IsHit}}{{.Hits}}{{else}}-{{end}}</td>
                        <td class="line-source"><pre>{{.Content}}</pre></td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </section>
    </main>

    <footer>
        <p><a href="index.html">← Back to Summary</a></p>
    </footer>

    <script>{{.JavaScript}}</script>
</body>
</html>`
}

// getCSS returns the CSS styles for the HTML report
func (hr *HTMLReporter) getCSS() string {
	return `
body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    margin: 0;
    padding: 20px;
    background-color: #f5f5f5;
    color: #333;
}

header {
    background: white;
    padding: 20px;
    border-radius: 8px;
    margin-bottom: 20px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

header h1 {
    margin: 0 0 10px 0;
    color: #2c3e50;
}

header h1 a {
    color: inherit;
    text-decoration: none;
}

header h2 {
    margin: 0 0 10px 0;
    color: #7f8c8d;
    font-weight: normal;
}

.summary, .test-results, .files {
    background: white;
    padding: 20px;
    border-radius: 8px;
    margin-bottom: 20px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.metrics {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 20px;
    margin-top: 20px;
}

.metric {
    text-align: center;
    padding: 20px;
    border-radius: 8px;
    background: #f8f9fa;
}

.metric-value {
    font-size: 2.5em;
    font-weight: bold;
    color: #2c3e50;
}

.metric-label {
    font-weight: bold;
    margin: 10px 0 5px 0;
    color: #7f8c8d;
}

.metric-detail {
    font-size: 0.9em;
    color: #95a5a6;
}

.test-summary {
    display: flex;
    gap: 20px;
    flex-wrap: wrap;
}

.test-stat {
    padding: 8px 16px;
    border-radius: 4px;
    font-weight: bold;
    background: #ecf0f1;
}

.test-stat.passed { background: #d5e8d4; color: #27ae60; }
.test-stat.failed { background: #f8d7da; color: #e74c3c; }
.test-stat.skipped { background: #fff3cd; color: #f39c12; }
.test-stat.timeout { background: #f4cccc; color: #c0392b; }

.file-table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 20px;
}

.file-table th,
.file-table td {
    padding: 12px;
    text-align: left;
    border-bottom: 1px solid #ecf0f1;
}

.file-table th {
    background: #f8f9fa;
    font-weight: bold;
    color: #2c3e50;
}

.file-table tr:hover {
    background: #f8f9fa;
}

.file-name a {
    color: #3498db;
    text-decoration: none;
}

.file-name a:hover {
    text-decoration: underline;
}

.coverage-percentage {
    font-weight: bold;
    text-align: right;
}

.line-count,
.uncovered-count {
    text-align: right;
    font-family: monospace;
}

.status-high { border-left: 4px solid #27ae60; }
.status-medium { border-left: 4px solid #f39c12; }
.status-low { border-left: 4px solid #e67e22; }
.status-critical { border-left: 4px solid #e74c3c; }

.source-table {
    width: 100%;
    border-collapse: collapse;
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    font-size: 14px;
}

.source-table th,
.source-table td {
    padding: 4px 8px;
    border-bottom: 1px solid #ecf0f1;
}

.source-table th {
    background: #f8f9fa;
    font-weight: bold;
    position: sticky;
    top: 0;
}

.line-number {
    text-align: right;
    color: #95a5a6;
    background: #f8f9fa;
    user-select: none;
    width: 50px;
}

.line-hits {
    text-align: right;
    width: 60px;
    font-weight: bold;
}

.line-source {
    white-space: nowrap;
}

.line-source pre {
    margin: 0;
    padding: 0;
    background: none;
    border: none;
    font: inherit;
    white-space: pre;
}

.line-covered {
    background: #d5e8d4;
}

.line-covered .line-hits {
    color: #27ae60;
}

.line-uncovered {
    background: #f8d7da;
}

.line-uncovered .line-hits {
    color: #e74c3c;
}

.line-ignored {
    color: #95a5a6;
    background: #fafafa;
}

footer {
    text-align: center;
    margin-top: 40px;
    color: #95a5a6;
}

footer a {
    color: #3498db;
    text-decoration: none;
}

footer a:hover {
    text-decoration: underline;
}

@media (max-width: 768px) {
    body {
        padding: 10px;
    }
    
    .metrics {
        grid-template-columns: 1fr;
    }
    
    .test-summary {
        flex-direction: column;
    }
    
    .source-table {
        font-size: 12px;
    }
}
`
}

// getJavaScript returns JavaScript for the HTML report
func (hr *HTMLReporter) getJavaScript() string {
	return `
// Add interactive features
document.addEventListener('DOMContentLoaded', function() {
    // Add tooltips for coverage metrics
    const metrics = document.querySelectorAll('.metric');
    metrics.forEach(metric => {
        metric.addEventListener('mouseenter', function() {
            this.style.transform = 'scale(1.05)';
            this.style.transition = 'transform 0.2s ease';
        });
        
        metric.addEventListener('mouseleave', function() {
            this.style.transform = 'scale(1)';
        });
    });

    // Add click-to-highlight for source lines
    const sourceLines = document.querySelectorAll('.source-table tr');
    sourceLines.forEach(line => {
        line.addEventListener('click', function() {
            // Remove previous highlights
            sourceLines.forEach(l => l.classList.remove('highlighted'));
            // Add highlight to clicked line
            this.classList.add('highlighted');
        });
    });
});

// Add CSS for highlighted lines
const style = document.createElement('style');
style.textContent = '.highlighted { outline: 2px solid #3498db; outline-offset: -2px; }';
document.head.appendChild(style);
`
}

// GenerateReport generates an HTML coverage report using the global collector
func GenerateHTMLReport(outputDir string, testResults *core.TestResults) error {
	reporter := NewHTMLReporter(outputDir)
	return reporter.Generate(coverage.GlobalCoverageCollector, testResults)
}
