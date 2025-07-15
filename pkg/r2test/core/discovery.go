package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TestDiscovery handles finding and loading test files
type TestDiscovery struct {
	Config *TestConfig
}

// NewTestDiscovery creates a new test discovery instance
func NewTestDiscovery(config *TestConfig) *TestDiscovery {
	return &TestDiscovery{
		Config: config,
	}
}

// DiscoverTests finds all test files matching the patterns
func (td *TestDiscovery) DiscoverTests() ([]string, error) {
	var testFiles []string

	for _, testDir := range td.Config.TestDirs {
		files, err := td.discoverInDirectory(testDir)
		if err != nil {
			return nil, fmt.Errorf("error discovering tests in %s: %w", testDir, err)
		}
		testFiles = append(testFiles, files...)
	}

	return testFiles, nil
}

// discoverInDirectory recursively finds test files in a directory
func (td *TestDiscovery) discoverInDirectory(dir string) ([]string, error) {
	var testFiles []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if it's a directory
		if info.IsDir() {
			// Check if this directory should be ignored
			if td.shouldIgnoreDirectory(path) {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if this file matches our test patterns
		if td.isTestFile(path) && !td.shouldIgnoreFile(path) {
			testFiles = append(testFiles, path)
		}

		return nil
	})

	return testFiles, err
}

// isTestFile checks if a file matches the test patterns
func (td *TestDiscovery) isTestFile(filePath string) bool {
	fileName := filepath.Base(filePath)

	// Check file extension
	if !strings.HasSuffix(fileName, ".r2") {
		return false
	}

	// Check patterns
	for _, pattern := range td.Config.TestPatterns {
		matched, err := filepath.Match(pattern, fileName)
		if err == nil && matched {
			return true
		}
	}

	return false
}

// shouldIgnoreDirectory checks if a directory should be ignored
func (td *TestDiscovery) shouldIgnoreDirectory(dirPath string) bool {
	dirName := filepath.Base(dirPath)

	for _, ignorePattern := range td.Config.IgnorePatterns {
		matched, err := filepath.Match(ignorePattern, dirName)
		if err == nil && matched {
			return true
		}
	}

	return false
}

// shouldIgnoreFile checks if a file should be ignored
func (td *TestDiscovery) shouldIgnoreFile(filePath string) bool {
	fileName := filepath.Base(filePath)

	for _, ignorePattern := range td.Config.IgnorePatterns {
		matched, err := filepath.Match(ignorePattern, fileName)
		if err == nil && matched {
			return true
		}
	}

	return false
}

// ParseTestFile parses a test file and extracts test suites
// This will be implemented to integrate with R2Lang parser
func (td *TestDiscovery) ParseTestFile(filePath string) (*TestSuite, error) {
	// TODO: Integrate with R2Lang parser to parse test files
	// For now, return a placeholder implementation

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read test file %s: %w", filePath, err)
	}

	suite := &TestSuite{
		Name:        extractSuiteName(filePath),
		Description: fmt.Sprintf("Test suite from %s", filePath),
		Tests:       make([]*TestCase, 0),
		Tags:        make([]string, 0),
	}

	// Basic parsing - this will be replaced with proper R2Lang integration
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)

		// Look for describe() function calls
		if strings.Contains(line, "describe(") {
			// Extract suite information
			// This is a simplified implementation
			if name := extractQuotedString(line); name != "" {
				suite.Description = name
			}
		}

		// Look for it() function calls
		if strings.Contains(line, "it(") {
			// Extract test case information
			if name := extractQuotedString(line); name != "" {
				test := &TestCase{
					Name:        name,
					Description: name,
					Tags:        make([]string, 0),
					Suite:       suite,
					// Func will be set when integrating with R2Lang interpreter
					Func: func() {
						// Placeholder - will execute R2Lang code
						fmt.Printf("Executing test: %s (line %d)\n", name, i+1)
					},
				}
				suite.Tests = append(suite.Tests, test)
			}
		}
	}

	return suite, nil
}

// extractSuiteName extracts a suite name from file path
func extractSuiteName(filePath string) string {
	fileName := filepath.Base(filePath)
	// Remove .r2 extension
	if strings.HasSuffix(fileName, ".r2") {
		fileName = fileName[:len(fileName)-3]
	}

	// Convert underscores and dashes to spaces and title case
	name := strings.ReplaceAll(fileName, "_", " ")
	name = strings.ReplaceAll(name, "-", " ")

	return strings.Title(name)
}

// extractQuotedString extracts a string from quotes (basic implementation)
func extractQuotedString(line string) string {
	// Look for content between quotes
	start := strings.Index(line, "\"")
	if start == -1 {
		start = strings.Index(line, "'")
	}
	if start == -1 {
		return ""
	}

	quote := line[start]
	end := strings.Index(line[start+1:], string(quote))
	if end == -1 {
		return ""
	}

	return line[start+1 : start+1+end]
}

// LoadTestSuites discovers and loads all test suites
func (td *TestDiscovery) LoadTestSuites() ([]*TestSuite, error) {
	testFiles, err := td.DiscoverTests()
	if err != nil {
		return nil, err
	}

	var suites []*TestSuite

	for _, filePath := range testFiles {
		suite, err := td.ParseTestFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to parse test file %s: %w", filePath, err)
		}

		if len(suite.Tests) > 0 { // Only add suites that have tests
			suites = append(suites, suite)
		}
	}

	return suites, nil
}
