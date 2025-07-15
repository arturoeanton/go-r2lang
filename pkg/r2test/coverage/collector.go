package coverage

import (
	"fmt"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// CoverageCollector collects code coverage data during test execution
type CoverageCollector struct {
	enabled      bool
	files        map[string]*FileCoverage
	mu           sync.RWMutex
	basePath     string
	startTime    time.Time
	excludeGlobs []string
}

// FileCoverage represents coverage data for a single file
type FileCoverage struct {
	Path         string
	Lines        map[int]*LineCoverage
	TotalLines   int
	CoveredLines int
	Statements   map[int]*StatementCoverage
	Branches     map[int]*BranchCoverage
	Functions    map[string]*FunctionCoverage
}

// LineCoverage represents coverage data for a single line
type LineCoverage struct {
	LineNumber int
	Hits       int
	IsHit      bool
	Source     string
	Type       LineType
}

// StatementCoverage represents coverage data for a statement
type StatementCoverage struct {
	ID        int
	StartLine int
	EndLine   int
	StartCol  int
	EndCol    int
	Hits      int
	IsHit     bool
}

// BranchCoverage represents coverage data for a branch
type BranchCoverage struct {
	ID           int
	LineNumber   int
	BranchNumber int
	Taken        int
	NotTaken     int
	Type         BranchType
}

// FunctionCoverage represents coverage data for a function
type FunctionCoverage struct {
	Name      string
	StartLine int
	EndLine   int
	Hits      int
	IsHit     bool
}

// LineType represents the type of a line of code
type LineType int

const (
	LineTypeCode LineType = iota
	LineTypeComment
	LineTypeEmpty
	LineTypeDeclaration
)

// BranchType represents the type of a branch
type BranchType int

const (
	BranchTypeIf BranchType = iota
	BranchTypeElse
	BranchTypeSwitch
	BranchTypeCase
	BranchTypeLoop
	BranchTypeTernary
)

// CoverageStats represents aggregated coverage statistics
type CoverageStats struct {
	TotalFiles          int
	TotalLines          int
	CoveredLines        int
	LinePercentage      float64
	TotalStatements     int
	CoveredStatements   int
	StatementPercentage float64
	TotalBranches       int
	CoveredBranches     int
	BranchPercentage    float64
	TotalFunctions      int
	CoveredFunctions    int
	FunctionPercentage  float64
	Duration            time.Duration
}

// NewCoverageCollector creates a new coverage collector
func NewCoverageCollector(basePath string) *CoverageCollector {
	return &CoverageCollector{
		enabled:      true,
		files:        make(map[string]*FileCoverage),
		basePath:     basePath,
		excludeGlobs: make([]string, 0),
	}
}

// Enable enables coverage collection
func (cc *CoverageCollector) Enable() {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.enabled = true
}

// Disable disables coverage collection
func (cc *CoverageCollector) Disable() {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.enabled = false
}

// IsEnabled returns whether coverage collection is enabled
func (cc *CoverageCollector) IsEnabled() bool {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	return cc.enabled
}

// Start starts coverage collection
func (cc *CoverageCollector) Start() {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.startTime = time.Now()
}

// AddExcludeGlob adds a glob pattern to exclude from coverage
func (cc *CoverageCollector) AddExcludeGlob(pattern string) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.excludeGlobs = append(cc.excludeGlobs, pattern)
}

// ShouldExcludeFile checks if a file should be excluded from coverage
func (cc *CoverageCollector) ShouldExcludeFile(filePath string) bool {
	for _, pattern := range cc.excludeGlobs {
		if matched, _ := filepath.Match(pattern, filePath); matched {
			return true
		}
		if matched, _ := filepath.Match(pattern, filepath.Base(filePath)); matched {
			return true
		}
	}
	return false
}

// RecordLineHit records a hit on a specific line
func (cc *CoverageCollector) RecordLineHit(filePath string, lineNumber int) {
	if !cc.IsEnabled() {
		return
	}

	if cc.ShouldExcludeFile(filePath) {
		return
	}

	cc.mu.Lock()
	defer cc.mu.Unlock()

	file := cc.getOrCreateFile(filePath)
	line := cc.getOrCreateLine(file, lineNumber)
	line.Hits++
	line.IsHit = true

	// Update file coverage stats
	cc.updateFileCoverage(file)
}

// RecordStatementHit records a hit on a specific statement
func (cc *CoverageCollector) RecordStatementHit(filePath string, statementID int) {
	if !cc.IsEnabled() {
		return
	}

	if cc.ShouldExcludeFile(filePath) {
		return
	}

	cc.mu.Lock()
	defer cc.mu.Unlock()

	file := cc.getOrCreateFile(filePath)
	statement, exists := file.Statements[statementID]
	if !exists {
		return
	}

	statement.Hits++
	statement.IsHit = true
}

// RecordBranchHit records a hit on a specific branch
func (cc *CoverageCollector) RecordBranchHit(filePath string, branchID int, taken bool) {
	if !cc.IsEnabled() {
		return
	}

	if cc.ShouldExcludeFile(filePath) {
		return
	}

	cc.mu.Lock()
	defer cc.mu.Unlock()

	file := cc.getOrCreateFile(filePath)
	branch, exists := file.Branches[branchID]
	if !exists {
		return
	}

	if taken {
		branch.Taken++
	} else {
		branch.NotTaken++
	}
}

// RecordFunctionHit records a hit on a specific function
func (cc *CoverageCollector) RecordFunctionHit(filePath string, functionName string) {
	if !cc.IsEnabled() {
		return
	}

	if cc.ShouldExcludeFile(filePath) {
		return
	}

	cc.mu.Lock()
	defer cc.mu.Unlock()

	file := cc.getOrCreateFile(filePath)
	function, exists := file.Functions[functionName]
	if !exists {
		return
	}

	function.Hits++
	function.IsHit = true
}

// getOrCreateFile gets or creates a file coverage record
func (cc *CoverageCollector) getOrCreateFile(filePath string) *FileCoverage {
	file, exists := cc.files[filePath]
	if !exists {
		file = &FileCoverage{
			Path:       filePath,
			Lines:      make(map[int]*LineCoverage),
			Statements: make(map[int]*StatementCoverage),
			Branches:   make(map[int]*BranchCoverage),
			Functions:  make(map[string]*FunctionCoverage),
		}
		cc.files[filePath] = file
	}
	return file
}

// getOrCreateLine gets or creates a line coverage record
func (cc *CoverageCollector) getOrCreateLine(file *FileCoverage, lineNumber int) *LineCoverage {
	line, exists := file.Lines[lineNumber]
	if !exists {
		line = &LineCoverage{
			LineNumber: lineNumber,
			Hits:       0,
			IsHit:      false,
			Type:       LineTypeCode,
		}
		file.Lines[lineNumber] = line
	}
	return line
}

// updateFileCoverage updates the coverage statistics for a file
func (cc *CoverageCollector) updateFileCoverage(file *FileCoverage) {
	totalLines := 0
	coveredLines := 0

	for _, line := range file.Lines {
		if line.Type == LineTypeCode {
			totalLines++
			if line.IsHit {
				coveredLines++
			}
		}
	}

	file.TotalLines = totalLines
	file.CoveredLines = coveredLines
}

// AddStatement adds a statement to track for coverage
func (cc *CoverageCollector) AddStatement(filePath string, statementID, startLine, endLine, startCol, endCol int) {
	if cc.ShouldExcludeFile(filePath) {
		return
	}

	cc.mu.Lock()
	defer cc.mu.Unlock()

	file := cc.getOrCreateFile(filePath)
	file.Statements[statementID] = &StatementCoverage{
		ID:        statementID,
		StartLine: startLine,
		EndLine:   endLine,
		StartCol:  startCol,
		EndCol:    endCol,
		Hits:      0,
		IsHit:     false,
	}
}

// AddBranch adds a branch to track for coverage
func (cc *CoverageCollector) AddBranch(filePath string, branchID, lineNumber, branchNumber int, branchType BranchType) {
	if cc.ShouldExcludeFile(filePath) {
		return
	}

	cc.mu.Lock()
	defer cc.mu.Unlock()

	file := cc.getOrCreateFile(filePath)
	file.Branches[branchID] = &BranchCoverage{
		ID:           branchID,
		LineNumber:   lineNumber,
		BranchNumber: branchNumber,
		Taken:        0,
		NotTaken:     0,
		Type:         branchType,
	}
}

// AddFunction adds a function to track for coverage
func (cc *CoverageCollector) AddFunction(filePath string, functionName string, startLine, endLine int) {
	if cc.ShouldExcludeFile(filePath) {
		return
	}

	cc.mu.Lock()
	defer cc.mu.Unlock()

	file := cc.getOrCreateFile(filePath)
	file.Functions[functionName] = &FunctionCoverage{
		Name:      functionName,
		StartLine: startLine,
		EndLine:   endLine,
		Hits:      0,
		IsHit:     false,
	}
}

// GetStats returns aggregated coverage statistics
func (cc *CoverageCollector) GetStats() *CoverageStats {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	stats := &CoverageStats{
		Duration: time.Since(cc.startTime),
	}

	for _, file := range cc.files {
		stats.TotalFiles++
		stats.TotalLines += file.TotalLines
		stats.CoveredLines += file.CoveredLines

		// Count statements
		for _, statement := range file.Statements {
			stats.TotalStatements++
			if statement.IsHit {
				stats.CoveredStatements++
			}
		}

		// Count branches
		for _, branch := range file.Branches {
			stats.TotalBranches++
			if branch.Taken > 0 {
				stats.CoveredBranches++
			}
		}

		// Count functions
		for _, function := range file.Functions {
			stats.TotalFunctions++
			if function.IsHit {
				stats.CoveredFunctions++
			}
		}
	}

	// Calculate percentages
	if stats.TotalLines > 0 {
		stats.LinePercentage = float64(stats.CoveredLines) / float64(stats.TotalLines) * 100
	}
	if stats.TotalStatements > 0 {
		stats.StatementPercentage = float64(stats.CoveredStatements) / float64(stats.TotalStatements) * 100
	}
	if stats.TotalBranches > 0 {
		stats.BranchPercentage = float64(stats.CoveredBranches) / float64(stats.TotalBranches) * 100
	}
	if stats.TotalFunctions > 0 {
		stats.FunctionPercentage = float64(stats.CoveredFunctions) / float64(stats.TotalFunctions) * 100
	}

	return stats
}

// GetFileCoverage returns coverage data for a specific file
func (cc *CoverageCollector) GetFileCoverage(filePath string) (*FileCoverage, bool) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	file, exists := cc.files[filePath]
	return file, exists
}

// GetAllFiles returns all files with coverage data
func (cc *CoverageCollector) GetAllFiles() map[string]*FileCoverage {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	result := make(map[string]*FileCoverage)
	for path, file := range cc.files {
		result[path] = file
	}
	return result
}

// GetSortedFiles returns all files sorted by coverage percentage
func (cc *CoverageCollector) GetSortedFiles() []*FileCoverage {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	files := make([]*FileCoverage, 0, len(cc.files))
	for _, file := range cc.files {
		files = append(files, file)
	}

	sort.Slice(files, func(i, j int) bool {
		percentageI := float64(files[i].CoveredLines) / float64(files[i].TotalLines) * 100
		percentageJ := float64(files[j].CoveredLines) / float64(files[j].TotalLines) * 100
		return percentageI < percentageJ
	})

	return files
}

// Clear clears all coverage data
func (cc *CoverageCollector) Clear() {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	cc.files = make(map[string]*FileCoverage)
	cc.startTime = time.Now()
}

// MeetsThreshold checks if coverage meets the specified threshold
func (cc *CoverageCollector) MeetsThreshold(threshold float64) bool {
	stats := cc.GetStats()
	return stats.LinePercentage >= threshold
}

// GetUncoveredLines returns lines that are not covered
func (cc *CoverageCollector) GetUncoveredLines(filePath string) []int {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	file, exists := cc.files[filePath]
	if !exists {
		return nil
	}

	var uncovered []int
	for lineNum, line := range file.Lines {
		if line.Type == LineTypeCode && !line.IsHit {
			uncovered = append(uncovered, lineNum)
		}
	}

	sort.Ints(uncovered)
	return uncovered
}

// String returns a string representation of coverage stats
func (cs *CoverageStats) String() string {
	return fmt.Sprintf(
		"Coverage: %.2f%% lines (%d/%d), %.2f%% statements (%d/%d), %.2f%% branches (%d/%d), %.2f%% functions (%d/%d)",
		cs.LinePercentage, cs.CoveredLines, cs.TotalLines,
		cs.StatementPercentage, cs.CoveredStatements, cs.TotalStatements,
		cs.BranchPercentage, cs.CoveredBranches, cs.TotalBranches,
		cs.FunctionPercentage, cs.CoveredFunctions, cs.TotalFunctions,
	)
}

// GetFilePercentage returns the coverage percentage for a file
func (fc *FileCoverage) GetFilePercentage() float64 {
	if fc.TotalLines == 0 {
		return 0
	}
	return float64(fc.CoveredLines) / float64(fc.TotalLines) * 100
}

// GetRelativePath returns the relative path of a file
func (fc *FileCoverage) GetRelativePath(basePath string) string {
	rel, err := filepath.Rel(basePath, fc.Path)
	if err != nil {
		return fc.Path
	}
	return rel
}

// Global coverage collector
var GlobalCoverageCollector = NewCoverageCollector(".")

// Global convenience functions

// EnableCoverage enables coverage collection using the global collector
func EnableCoverage() {
	GlobalCoverageCollector.Enable()
}

// DisableCoverage disables coverage collection using the global collector
func DisableCoverage() {
	GlobalCoverageCollector.Disable()
}

// StartCoverage starts coverage collection using the global collector
func StartCoverage() {
	GlobalCoverageCollector.Start()
}

// RecordHit records a line hit using the global collector
func RecordHit(filePath string, lineNumber int) {
	GlobalCoverageCollector.RecordLineHit(filePath, lineNumber)
}

// GetCoverageStats returns coverage statistics using the global collector
func GetCoverageStats() *CoverageStats {
	return GlobalCoverageCollector.GetStats()
}

// SetCoverageBasePath sets the base path for the global collector
func SetCoverageBasePath(basePath string) {
	GlobalCoverageCollector.basePath = basePath
}
