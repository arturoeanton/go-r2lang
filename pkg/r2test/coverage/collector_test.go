package coverage

import (
	"testing"
	"time"
)

func TestNewCoverageCollector(t *testing.T) {
	collector := NewCoverageCollector("/test/path")

	if collector.basePath != "/test/path" {
		t.Errorf("Expected base path '/test/path', got '%s'", collector.basePath)
	}

	if !collector.enabled {
		t.Error("Expected collector to be enabled by default")
	}

	if collector.files == nil {
		t.Error("Expected files map to be initialized")
	}

	if len(collector.excludeGlobs) != 0 {
		t.Errorf("Expected empty exclude globs, got %d", len(collector.excludeGlobs))
	}
}

func TestCoverageCollectorEnableDisable(t *testing.T) {
	collector := NewCoverageCollector(".")

	// Should be enabled by default
	if !collector.IsEnabled() {
		t.Error("Expected collector to be enabled by default")
	}

	// Disable
	collector.Disable()
	if collector.IsEnabled() {
		t.Error("Expected collector to be disabled after Disable()")
	}

	// Enable again
	collector.Enable()
	if !collector.IsEnabled() {
		t.Error("Expected collector to be enabled after Enable()")
	}
}

func TestCoverageCollectorExcludeGlobs(t *testing.T) {
	collector := NewCoverageCollector(".")

	// Add exclude patterns
	collector.AddExcludeGlob("*_test.go")
	collector.AddExcludeGlob("vendor/*")

	if len(collector.excludeGlobs) != 2 {
		t.Errorf("Expected 2 exclude globs, got %d", len(collector.excludeGlobs))
	}

	// Test exclusion
	testCases := []struct {
		path     string
		excluded bool
	}{
		{"main.go", false},
		{"main_test.go", true},
		{"vendor/lib.go", true},
		{"src/util.go", false},
	}

	for _, tc := range testCases {
		result := collector.ShouldExcludeFile(tc.path)
		if result != tc.excluded {
			t.Errorf("Expected ShouldExcludeFile('%s') to be %t, got %t", tc.path, tc.excluded, result)
		}
	}
}

func TestCoverageCollectorRecordLineHit(t *testing.T) {
	collector := NewCoverageCollector(".")
	collector.Start()

	// Record some line hits
	collector.RecordLineHit("test.go", 1)
	collector.RecordLineHit("test.go", 2)
	collector.RecordLineHit("test.go", 1) // Hit line 1 again
	collector.RecordLineHit("other.go", 5)

	// Check file was created
	file, exists := collector.GetFileCoverage("test.go")
	if !exists {
		t.Error("Expected file coverage to exist for 'test.go'")
	}

	// Check line coverage
	line1, exists := file.Lines[1]
	if !exists {
		t.Error("Expected line 1 coverage to exist")
	}

	if line1.Hits != 2 {
		t.Errorf("Expected line 1 to have 2 hits, got %d", line1.Hits)
	}

	if !line1.IsHit {
		t.Error("Expected line 1 to be marked as hit")
	}

	line2, exists := file.Lines[2]
	if !exists {
		t.Error("Expected line 2 coverage to exist")
	}

	if line2.Hits != 1 {
		t.Errorf("Expected line 2 to have 1 hit, got %d", line2.Hits)
	}

	// Check other file
	otherFile, exists := collector.GetFileCoverage("other.go")
	if !exists {
		t.Error("Expected file coverage to exist for 'other.go'")
	}

	line5, exists := otherFile.Lines[5]
	if !exists {
		t.Error("Expected line 5 coverage to exist")
	}

	if line5.Hits != 1 {
		t.Errorf("Expected line 5 to have 1 hit, got %d", line5.Hits)
	}
}

func TestCoverageCollectorRecordDisabled(t *testing.T) {
	collector := NewCoverageCollector(".")
	collector.Disable()

	// Record hit while disabled
	collector.RecordLineHit("test.go", 1)

	// Should not create any files
	files := collector.GetAllFiles()
	if len(files) != 0 {
		t.Errorf("Expected no files when disabled, got %d", len(files))
	}
}

func TestCoverageCollectorRecordExcluded(t *testing.T) {
	collector := NewCoverageCollector(".")
	collector.AddExcludeGlob("*_test.go")

	// Record hit on excluded file
	collector.RecordLineHit("main_test.go", 1)

	// Should not create file
	_, exists := collector.GetFileCoverage("main_test.go")
	if exists {
		t.Error("Expected excluded file to not have coverage recorded")
	}
}

func TestCoverageCollectorAddStatement(t *testing.T) {
	collector := NewCoverageCollector(".")

	// Add statement
	collector.AddStatement("test.go", 1, 10, 15, 5, 20)

	file, exists := collector.GetFileCoverage("test.go")
	if !exists {
		t.Error("Expected file to be created when adding statement")
	}

	stmt, exists := file.Statements[1]
	if !exists {
		t.Error("Expected statement to be added")
	}

	if stmt.ID != 1 || stmt.StartLine != 10 || stmt.EndLine != 15 || stmt.StartCol != 5 || stmt.EndCol != 20 {
		t.Errorf("Statement data incorrect: got %+v", stmt)
	}

	if stmt.IsHit {
		t.Error("Expected statement to not be hit initially")
	}
}

func TestCoverageCollectorRecordStatementHit(t *testing.T) {
	collector := NewCoverageCollector(".")

	// Add and hit statement
	collector.AddStatement("test.go", 1, 10, 15, 5, 20)
	collector.RecordStatementHit("test.go", 1)

	file, exists := collector.GetFileCoverage("test.go")
	if !exists {
		t.Error("Expected file to exist")
	}

	stmt, exists := file.Statements[1]
	if !exists {
		t.Error("Expected statement to exist")
	}

	if stmt.Hits != 1 {
		t.Errorf("Expected statement to have 1 hit, got %d", stmt.Hits)
	}

	if !stmt.IsHit {
		t.Error("Expected statement to be marked as hit")
	}
}

func TestCoverageCollectorAddFunction(t *testing.T) {
	collector := NewCoverageCollector(".")

	// Add function
	collector.AddFunction("test.go", "testFunc", 10, 20)

	file, exists := collector.GetFileCoverage("test.go")
	if !exists {
		t.Error("Expected file to be created when adding function")
	}

	function, exists := file.Functions["testFunc"]
	if !exists {
		t.Error("Expected function to be added")
	}

	if function.Name != "testFunc" || function.StartLine != 10 || function.EndLine != 20 {
		t.Errorf("Function data incorrect: got %+v", function)
	}

	if function.IsHit {
		t.Error("Expected function to not be hit initially")
	}
}

func TestCoverageCollectorRecordFunctionHit(t *testing.T) {
	collector := NewCoverageCollector(".")

	// Add and hit function
	collector.AddFunction("test.go", "testFunc", 10, 20)
	collector.RecordFunctionHit("test.go", "testFunc")

	file, exists := collector.GetFileCoverage("test.go")
	if !exists {
		t.Error("Expected file to exist")
	}

	function, exists := file.Functions["testFunc"]
	if !exists {
		t.Error("Expected function to exist")
	}

	if function.Hits != 1 {
		t.Errorf("Expected function to have 1 hit, got %d", function.Hits)
	}

	if !function.IsHit {
		t.Error("Expected function to be marked as hit")
	}
}

func TestCoverageCollectorGetStats(t *testing.T) {
	collector := NewCoverageCollector(".")
	collector.Start()

	// Add some coverage data
	collector.RecordLineHit("file1.go", 1)
	collector.RecordLineHit("file1.go", 2)
	collector.RecordLineHit("file1.go", 3)
	collector.RecordLineHit("file2.go", 1)

	// Add statements
	collector.AddStatement("file1.go", 1, 1, 1, 0, 10)
	collector.AddStatement("file1.go", 2, 2, 2, 0, 10)
	collector.RecordStatementHit("file1.go", 1)

	// Add functions
	collector.AddFunction("file1.go", "func1", 1, 5)
	collector.AddFunction("file2.go", "func2", 1, 3)
	collector.RecordFunctionHit("file1.go", "func1")

	stats := collector.GetStats()

	if stats.TotalFiles != 2 {
		t.Errorf("Expected 2 total files, got %d", stats.TotalFiles)
	}

	if stats.TotalLines != 4 {
		t.Errorf("Expected 4 total lines, got %d", stats.TotalLines)
	}

	if stats.CoveredLines != 4 {
		t.Errorf("Expected 4 covered lines, got %d", stats.CoveredLines)
	}

	if stats.LinePercentage != 100.0 {
		t.Errorf("Expected 100%% line coverage, got %.2f", stats.LinePercentage)
	}

	if stats.TotalStatements != 2 {
		t.Errorf("Expected 2 total statements, got %d", stats.TotalStatements)
	}

	if stats.CoveredStatements != 1 {
		t.Errorf("Expected 1 covered statement, got %d", stats.CoveredStatements)
	}

	if stats.StatementPercentage != 50.0 {
		t.Errorf("Expected 50%% statement coverage, got %.2f", stats.StatementPercentage)
	}

	if stats.TotalFunctions != 2 {
		t.Errorf("Expected 2 total functions, got %d", stats.TotalFunctions)
	}

	if stats.CoveredFunctions != 1 {
		t.Errorf("Expected 1 covered function, got %d", stats.CoveredFunctions)
	}

	if stats.FunctionPercentage != 50.0 {
		t.Errorf("Expected 50%% function coverage, got %.2f", stats.FunctionPercentage)
	}

	if stats.Duration <= 0 {
		t.Error("Expected duration to be positive")
	}
}

func TestCoverageCollectorMeetsThreshold(t *testing.T) {
	collector := NewCoverageCollector(".")

	// Add coverage that results in 75% line coverage
	collector.RecordLineHit("test.go", 1)
	collector.RecordLineHit("test.go", 2)
	collector.RecordLineHit("test.go", 3)

	// Create an uncovered line by adding it to the file but not hitting it
	file := collector.getOrCreateFile("test.go")
	file.Lines[4] = &LineCoverage{
		LineNumber: 4,
		Hits:       0,
		IsHit:      false,
		Type:       LineTypeCode,
	}
	collector.updateFileCoverage(file)

	// Should meet 70% threshold
	if !collector.MeetsThreshold(70.0) {
		t.Error("Expected to meet 70% threshold")
	}

	// Should not meet 80% threshold
	if collector.MeetsThreshold(80.0) {
		t.Error("Expected to not meet 80% threshold")
	}
}

func TestCoverageCollectorGetUncoveredLines(t *testing.T) {
	collector := NewCoverageCollector(".")

	// Add mixed coverage
	collector.RecordLineHit("test.go", 1)
	collector.RecordLineHit("test.go", 3)

	// Add uncovered lines
	file := collector.getOrCreateFile("test.go")
	file.Lines[2] = &LineCoverage{
		LineNumber: 2,
		Hits:       0,
		IsHit:      false,
		Type:       LineTypeCode,
	}
	file.Lines[4] = &LineCoverage{
		LineNumber: 4,
		Hits:       0,
		IsHit:      false,
		Type:       LineTypeCode,
	}

	uncovered := collector.GetUncoveredLines("test.go")

	if len(uncovered) != 2 {
		t.Errorf("Expected 2 uncovered lines, got %d", len(uncovered))
	}

	if uncovered[0] != 2 || uncovered[1] != 4 {
		t.Errorf("Expected uncovered lines [2, 4], got %v", uncovered)
	}
}

func TestCoverageCollectorClear(t *testing.T) {
	collector := NewCoverageCollector(".")
	collector.Start()

	// Add some data
	collector.RecordLineHit("test.go", 1)

	// Verify data exists
	files := collector.GetAllFiles()
	if len(files) != 1 {
		t.Errorf("Expected 1 file before clear, got %d", len(files))
	}

	// Clear
	oldStartTime := collector.startTime
	time.Sleep(1 * time.Millisecond) // Ensure time difference
	collector.Clear()

	// Verify data is cleared
	files = collector.GetAllFiles()
	if len(files) != 0 {
		t.Errorf("Expected 0 files after clear, got %d", len(files))
	}

	// Verify start time was reset
	if !collector.startTime.After(oldStartTime) {
		t.Error("Expected start time to be reset after clear")
	}
}

func TestFileCoverageGetFilePercentage(t *testing.T) {
	file := &FileCoverage{
		TotalLines:   10,
		CoveredLines: 7,
	}

	percentage := file.GetFilePercentage()
	expected := 70.0

	if percentage != expected {
		t.Errorf("Expected percentage %.1f, got %.1f", expected, percentage)
	}

	// Test zero total lines
	emptyFile := &FileCoverage{
		TotalLines:   0,
		CoveredLines: 0,
	}

	percentage = emptyFile.GetFilePercentage()
	if percentage != 0.0 {
		t.Errorf("Expected 0%% for empty file, got %.1f", percentage)
	}
}

func TestCoverageStatsString(t *testing.T) {
	stats := &CoverageStats{
		TotalLines:          100,
		CoveredLines:        85,
		LinePercentage:      85.0,
		TotalStatements:     50,
		CoveredStatements:   40,
		StatementPercentage: 80.0,
		TotalBranches:       20,
		CoveredBranches:     15,
		BranchPercentage:    75.0,
		TotalFunctions:      10,
		CoveredFunctions:    8,
		FunctionPercentage:  80.0,
	}

	str := stats.String()

	expectedContains := []string{
		"85.00% lines (85/100)",
		"80.00% statements (40/50)",
		"75.00% branches (15/20)",
		"80.00% functions (8/10)",
	}

	for _, expected := range expectedContains {
		if !contains(str, expected) {
			t.Errorf("Expected string to contain '%s', got '%s'", expected, str)
		}
	}
}

func TestGlobalCoverageCollector(t *testing.T) {
	// Test global convenience functions

	// Enable and start
	EnableCoverage()
	StartCoverage()

	if !GlobalCoverageCollector.IsEnabled() {
		t.Error("Expected global collector to be enabled")
	}

	// Record hit
	RecordHit("global_test.go", 1)

	// Check stats
	stats := GetCoverageStats()
	if stats.TotalLines != 1 {
		t.Errorf("Expected 1 total line in global stats, got %d", stats.TotalLines)
	}

	if stats.CoveredLines != 1 {
		t.Errorf("Expected 1 covered line in global stats, got %d", stats.CoveredLines)
	}

	// Clean up
	GlobalCoverageCollector.Clear()
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
