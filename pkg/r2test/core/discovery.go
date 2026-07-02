package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
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

// ParseTestFile parses a test file by actually running it through the
// R2Lang interpreter (see r2bridge.go): the file's top-level describe()
// calls execute for real, and each one's it()/beforeEach()/etc. calls
// build a real TestSuite/TestCase whose Func actually invokes the R2Lang
// test body — not a static regex-extracted stub. A single file commonly
// contains multiple describe() blocks (see examples/testing/*.r2), so this
// returns every suite built while evaluating the file, not just one.
//
// Note: this executes describe()'s own callback bodies immediately (to
// collect their it()/hook calls), matching this framework's existing
// Go-level Describe()/It() semantics — but it() bodies themselves are NOT
// run here, only registered; they run later via TestRunner.
func (td *TestDiscovery) ParseTestFile(filePath string) ([]*TestSuite, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read test file %s: %w", filePath, err)
	}

	ctx := &r2BridgeContext{}
	env := newR2Environment(ctx)

	if err := runTestFileProgram(env, filePath, string(content)); err != nil {
		return nil, fmt.Errorf("failed to parse/run test file %s: %w", filePath, err)
	}

	return ctx.suites, nil
}

// runTestFileProgram parses and evaluates code as an R2Lang program in env,
// converting any panic (parser error, or a describe()/it() registration
// call failing outside a test body) into a Go error instead of letting it
// propagate — a malformed test FILE should be reported as a discovery
// error, distinct from an individual test's body panicking (which
// TestRunner already handles per-test via its own recover()).
func runTestFileProgram(env *r2core.Environment, filePath, code string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	parser := r2core.NewParserWithFile(code, filePath)
	program := parser.ParseProgram()
	program.Eval(env)
	return nil
}

// LoadTestSuites discovers and loads all test suites
func (td *TestDiscovery) LoadTestSuites() ([]*TestSuite, error) {
	testFiles, err := td.DiscoverTests()
	if err != nil {
		return nil, err
	}

	var suites []*TestSuite

	for _, filePath := range testFiles {
		fileSuites, err := td.ParseTestFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to parse test file %s: %w", filePath, err)
		}

		for _, suite := range fileSuites {
			if len(suite.Tests) > 0 { // Only add suites that have tests
				suites = append(suites, suite)
			}
		}
	}

	return suites, nil
}
