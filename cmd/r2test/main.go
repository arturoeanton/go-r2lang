package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2test"
	"github.com/arturoeanton/go-r2lang/pkg/r2test/core"
	"github.com/arturoeanton/go-r2lang/pkg/r2test/coverage"
	"github.com/arturoeanton/go-r2lang/pkg/r2test/fixtures"
	"github.com/arturoeanton/go-r2lang/pkg/r2test/mocking"
	"github.com/arturoeanton/go-r2lang/pkg/r2test/reporters"
)

const version = "1.0.0"

type Config struct {
	// Test Discovery
	TestDirs []string `json:"testDirs"`
	Patterns []string `json:"patterns"`
	Ignore   []string `json:"ignore"`

	// Execution
	Timeout    string `json:"timeout"`
	Parallel   bool   `json:"parallel"`
	MaxWorkers int    `json:"maxWorkers"`
	Bail       bool   `json:"bail"`

	// Filtering
	Grep string   `json:"grep"`
	Tags []string `json:"tags"`
	Skip []string `json:"skip"`
	Only []string `json:"only"`

	// Coverage
	Coverage CoverageConfig `json:"coverage"`

	// Reporting
	Reporters []string `json:"reporters"`
	OutputDir string   `json:"outputDir"`
	Verbose   bool     `json:"verbose"`

	// Advanced
	Retries   int  `json:"retries"`
	WatchMode bool `json:"watchMode"`

	// Fixtures
	FixtureDir string `json:"fixtureDir"`
}

type CoverageConfig struct {
	Enabled   bool     `json:"enabled"`
	Threshold float64  `json:"threshold"`
	Output    string   `json:"output"`
	Formats   []string `json:"formats"`
	Exclude   []string `json:"exclude"`
}

func defaultConfig() *Config {
	return &Config{
		TestDirs:   []string{"./tests", "./test"},
		Patterns:   []string{"*_test.r2", "*Test.r2", "test_*.r2"},
		Ignore:     []string{"node_modules", "vendor", "*.tmp.r2"},
		Timeout:    "30s",
		Parallel:   false,
		MaxWorkers: 4,
		Bail:       false,
		Coverage: CoverageConfig{
			Enabled:   false,
			Threshold: 80.0,
			Output:    "./coverage",
			Formats:   []string{"html"},
			Exclude:   []string{"*_test.r2", "vendor/*", "node_modules/*"},
		},
		Reporters:  []string{"console"},
		OutputDir:  "./test-results",
		Verbose:    false,
		Retries:    0,
		WatchMode:  false,
		FixtureDir: "./fixtures",
	}
}

func main() {
	var (
		// Basic flags
		helpFlag    = flag.Bool("help", false, "Show help information")
		versionFlag = flag.Bool("version", false, "Show version information")
		configFile  = flag.String("config", "", "Path to configuration file (JSON)")

		// Test discovery
		testDirs = flag.String("dirs", "", "Comma-separated list of test directories")
		patterns = flag.String("patterns", "", "Comma-separated list of test file patterns")
		ignore   = flag.String("ignore", "", "Comma-separated list of patterns to ignore")

		// Execution
		timeout    = flag.String("timeout", "", "Test timeout (e.g., 30s, 5m)")
		parallel   = flag.Bool("parallel", false, "Run tests in parallel")
		maxWorkers = flag.Int("workers", 0, "Maximum number of parallel workers")
		bail       = flag.Bool("bail", false, "Stop on first test failure")

		// Filtering
		grep = flag.String("grep", "", "Run only tests matching pattern")
		tags = flag.String("tags", "", "Run only tests with specified tags")
		skip = flag.String("skip", "", "Skip tests matching pattern")
		only = flag.String("only", "", "Run only tests matching pattern (exclusive)")

		// Coverage
		coverage          = flag.Bool("coverage", false, "Enable coverage collection")
		coverageThreshold = flag.Float64("coverage-threshold", 0, "Coverage threshold percentage")
		coverageOutput    = flag.String("coverage-output", "", "Coverage output directory")
		coverageFormats   = flag.String("coverage-formats", "", "Coverage report formats (html,json)")
		coverageExclude   = flag.String("coverage-exclude", "", "Coverage exclude patterns")

		// Reporting
		reporters = flag.String("reporters", "", "Comma-separated list of reporters")
		outputDir = flag.String("output", "", "Output directory for reports")
		verbose   = flag.Bool("verbose", false, "Verbose output")
		quiet     = flag.Bool("quiet", false, "Quiet output (errors only)")

		// Advanced
		retries = flag.Int("retries", 0, "Number of retries for failed tests")
		watch   = flag.Bool("watch", false, "Watch mode - rerun tests on file changes")

		// Fixtures
		fixtureDir = flag.String("fixtures", "", "Fixture directory")

		// Mock and isolation
		cleanupMocks = flag.Bool("cleanup-mocks", true, "Cleanup mocks after tests")
		isolation    = flag.Bool("isolation", false, "Run tests in isolation contexts")

		// Debugging
		debug  = flag.Bool("debug", false, "Enable debug output")
		dryRun = flag.Bool("dry-run", false, "Show what would be executed without running")
	)

	flag.Usage = func() {
		showHelp()
	}

	flag.Parse()

	if *helpFlag {
		showHelp()
		return
	}

	if *versionFlag {
		showVersion()
		return
	}

	// Load configuration
	config := defaultConfig()
	if *configFile != "" {
		if err := loadConfigFile(*configFile, config); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config file: %v\n", err)
			os.Exit(1)
		}
	}

	// Override config with command line flags
	overrideConfig(config, &configOverrides{
		testDirs:          *testDirs,
		patterns:          *patterns,
		ignore:            *ignore,
		timeout:           *timeout,
		parallel:          *parallel,
		maxWorkers:        *maxWorkers,
		bail:              *bail,
		grep:              *grep,
		tags:              *tags,
		skip:              *skip,
		only:              *only,
		coverage:          *coverage,
		coverageThreshold: *coverageThreshold,
		coverageOutput:    *coverageOutput,
		coverageFormats:   *coverageFormats,
		coverageExclude:   *coverageExclude,
		reporters:         *reporters,
		outputDir:         *outputDir,
		verbose:           *verbose,
		quiet:             *quiet,
		retries:           *retries,
		watch:             *watch,
		fixtureDir:        *fixtureDir,
		cleanupMocks:      *cleanupMocks,
		isolation:         *isolation,
		debug:             *debug,
		dryRun:            *dryRun,
	})

	if *dryRun {
		showDryRun(config)
		return
	}

	// Get test files/directories from arguments
	args := flag.Args()
	if len(args) > 0 {
		config.TestDirs = args
	}

	// Run tests
	exitCode := runTests(config)
	os.Exit(exitCode)
}

type configOverrides struct {
	testDirs          string
	patterns          string
	ignore            string
	timeout           string
	parallel          bool
	maxWorkers        int
	bail              bool
	grep              string
	tags              string
	skip              string
	only              string
	coverage          bool
	coverageThreshold float64
	coverageOutput    string
	coverageFormats   string
	coverageExclude   string
	reporters         string
	outputDir         string
	verbose           bool
	quiet             bool
	retries           int
	watch             bool
	fixtureDir        string
	cleanupMocks      bool
	isolation         bool
	debug             bool
	dryRun            bool
}

func overrideConfig(config *Config, overrides *configOverrides) {
	if overrides.testDirs != "" {
		config.TestDirs = strings.Split(overrides.testDirs, ",")
	}
	if overrides.patterns != "" {
		config.Patterns = strings.Split(overrides.patterns, ",")
	}
	if overrides.ignore != "" {
		config.Ignore = strings.Split(overrides.ignore, ",")
	}
	if overrides.timeout != "" {
		config.Timeout = overrides.timeout
	}
	if overrides.parallel {
		config.Parallel = true
	}
	if overrides.maxWorkers > 0 {
		config.MaxWorkers = overrides.maxWorkers
	}
	if overrides.bail {
		config.Bail = true
	}
	if overrides.grep != "" {
		config.Grep = overrides.grep
	}
	if overrides.tags != "" {
		config.Tags = strings.Split(overrides.tags, ",")
	}
	if overrides.skip != "" {
		config.Skip = strings.Split(overrides.skip, ",")
	}
	if overrides.only != "" {
		config.Only = strings.Split(overrides.only, ",")
	}
	if overrides.coverage {
		config.Coverage.Enabled = true
	}
	if overrides.coverageThreshold > 0 {
		config.Coverage.Threshold = overrides.coverageThreshold
	}
	if overrides.coverageOutput != "" {
		config.Coverage.Output = overrides.coverageOutput
	}
	if overrides.coverageFormats != "" {
		config.Coverage.Formats = strings.Split(overrides.coverageFormats, ",")
	}
	if overrides.coverageExclude != "" {
		config.Coverage.Exclude = strings.Split(overrides.coverageExclude, ",")
	}
	if overrides.reporters != "" {
		config.Reporters = strings.Split(overrides.reporters, ",")
	}
	if overrides.outputDir != "" {
		config.OutputDir = overrides.outputDir
	}
	if overrides.verbose {
		config.Verbose = true
	}
	if overrides.retries > 0 {
		config.Retries = overrides.retries
	}
	if overrides.watch {
		config.WatchMode = true
	}
	if overrides.fixtureDir != "" {
		config.FixtureDir = overrides.fixtureDir
	}
}

func loadConfigFile(filename string, config *Config) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, config)
}

func runTests(config *Config) int {
	if config.Verbose {
		fmt.Printf("R2Test v%s - Running tests...\n", version)
	}

	// Setup coverage if enabled
	if config.Coverage.Enabled {
		coverage.EnableCoverage()
		coverage.StartCoverage()
		coverage.SetCoverageBasePath(".")

		for _, pattern := range config.Coverage.Exclude {
			coverage.GlobalCoverageCollector.AddExcludeGlob(pattern)
		}
	}

	// Setup fixtures
	fixtures.SetFixtureBasePath(config.FixtureDir)

	// Create test configuration
	testConfig := &core.TestConfig{
		TestDirs:       config.TestDirs,
		TestPatterns:   config.Patterns,
		IgnorePatterns: config.Ignore,
		Parallel:       config.Parallel,
		MaxWorkers:     config.MaxWorkers,
		Bail:           config.Bail,
		FilterTags:     config.Tags,
		Verbose:        config.Verbose,
	}

	// Parse timeout
	if config.Timeout != "" {
		if timeout, err := time.ParseDuration(config.Timeout); err == nil {
			testConfig.DefaultTimeout = timeout
		} else {
			fmt.Fprintf(os.Stderr, "Invalid timeout format: %s\n", config.Timeout)
			return 1
		}
	}

	// Create and run tests
	r2test := r2test.New(testConfig)

	var results *core.TestResults
	var err error

	if config.WatchMode {
		fmt.Println("Watch mode not yet implemented")
		return 1
	} else {
		results, err = r2test.RunDiscoveredTests()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running tests: %v\n", err)
		return 1
	}

	// Generate reports
	if err := generateReports(config, results); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating reports: %v\n", err)
		return 1
	}

	// Check coverage threshold
	if config.Coverage.Enabled {
		stats := coverage.GetCoverageStats()
		if config.Verbose {
			fmt.Printf("\nCoverage Statistics:\n")
			fmt.Printf("  Lines: %.2f%% (%d/%d)\n", stats.LinePercentage, stats.CoveredLines, stats.TotalLines)
			fmt.Printf("  Statements: %.2f%% (%d/%d)\n", stats.StatementPercentage, stats.CoveredStatements, stats.TotalStatements)
			fmt.Printf("  Functions: %.2f%% (%d/%d)\n", stats.FunctionPercentage, stats.CoveredFunctions, stats.TotalFunctions)
		}

		if stats.LinePercentage < config.Coverage.Threshold {
			fmt.Fprintf(os.Stderr, "Coverage threshold not met: %.2f%% < %.2f%%\n",
				stats.LinePercentage, config.Coverage.Threshold)
			return 1
		}
	}

	// Cleanup
	fixtures.ClearFixtures()
	mocking.ResetAllMocks()
	mocking.RestoreAllMocks()
	mocking.ResetAllSpies()
	mocking.RestoreAllSpies()

	// Print summary
	testStats := results.GetStats()
	if !config.Verbose && testStats.Failed == 0 {
		fmt.Printf("✓ %d tests passed", testStats.Passed)
		if config.Coverage.Enabled {
			stats := coverage.GetCoverageStats()
			fmt.Printf(" (%.1f%% coverage)", stats.LinePercentage)
		}
		fmt.Println()
	} else {
		fmt.Printf("\nTest Results:\n")
		fmt.Printf("  Total: %d\n", testStats.Total)
		fmt.Printf("  Passed: %d\n", testStats.Passed)
		fmt.Printf("  Failed: %d\n", testStats.Failed)
		fmt.Printf("  Skipped: %d\n", testStats.Skipped)
		fmt.Printf("  Duration: %v\n", results.Duration)
	}

	if testStats.Failed > 0 {
		return 1
	}

	return 0
}

func generateReports(config *Config, results *core.TestResults) error {
	// Ensure output directory exists
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return err
	}

	for _, reporter := range config.Reporters {
		switch reporter {
		case "console":
			// Console output is handled by the test runner

		case "json":
			jsonPath := filepath.Join(config.OutputDir, "test-results.json")
			if err := reporters.GenerateJSONReport(jsonPath, results); err != nil {
				return fmt.Errorf("failed to generate JSON report: %w", err)
			}
			if config.Verbose {
				fmt.Printf("JSON report generated: %s\n", jsonPath)
			}

		case "junit":
			junitPath := filepath.Join(config.OutputDir, "junit.xml")
			if err := reporters.GenerateJUnitReport(junitPath, results); err != nil {
				return fmt.Errorf("failed to generate JUnit report: %w", err)
			}
			if config.Verbose {
				fmt.Printf("JUnit report generated: %s\n", junitPath)
			}
		}
	}

	// Generate coverage reports if enabled
	if config.Coverage.Enabled {
		coverageDir := config.Coverage.Output
		if err := os.MkdirAll(coverageDir, 0755); err != nil {
			return err
		}

		for _, format := range config.Coverage.Formats {
			switch format {
			case "html":
				htmlDir := filepath.Join(coverageDir, "html")
				if err := reporters.GenerateHTMLReport(htmlDir, results); err != nil {
					return fmt.Errorf("failed to generate HTML coverage report: %w", err)
				}
				if config.Verbose {
					fmt.Printf("HTML coverage report generated: %s\n", htmlDir)
				}

			case "json":
				jsonPath := filepath.Join(coverageDir, "coverage.json")
				if err := reporters.GenerateCoverageOnlyJSON(jsonPath); err != nil {
					return fmt.Errorf("failed to generate JSON coverage report: %w", err)
				}
				if config.Verbose {
					fmt.Printf("JSON coverage report generated: %s\n", jsonPath)
				}
			}
		}
	}

	return nil
}

func showHelp() {
	fmt.Printf("R2Test v%s - Advanced Testing Framework for R2Lang\n\n", version)
	fmt.Println("USAGE:")
	fmt.Println("  r2test [OPTIONS] [DIRECTORIES...]")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  r2test                              # Run tests in default directories")
	fmt.Println("  r2test ./tests                      # Run tests in specific directory")
	fmt.Println("  r2test -coverage -verbose ./tests   # Run with coverage and verbose output")
	fmt.Println("  r2test -grep \"Calculator\" ./tests   # Run only tests matching pattern")
	fmt.Println("  r2test -parallel -workers 4         # Run tests in parallel")
	fmt.Println("  r2test -config r2test.json          # Use configuration file")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println()
	fmt.Println("Basic Options:")
	fmt.Println("  -help                    Show this help message")
	fmt.Println("  -version                 Show version information")
	fmt.Println("  -config FILE             Load configuration from JSON file")
	fmt.Println("  -verbose                 Enable verbose output")
	fmt.Println("  -quiet                   Quiet output (errors only)")
	fmt.Println("  -debug                   Enable debug output")
	fmt.Println("  -dry-run                 Show what would be executed without running")
	fmt.Println()
	fmt.Println("Test Discovery:")
	fmt.Println("  -dirs DIRS               Comma-separated test directories (default: ./tests,./test)")
	fmt.Println("  -patterns PATTERNS       Test file patterns (default: *_test.r2,*Test.r2,test_*.r2)")
	fmt.Println("  -ignore PATTERNS         Patterns to ignore (default: node_modules,vendor)")
	fmt.Println()
	fmt.Println("Test Execution:")
	fmt.Println("  -timeout DURATION        Test timeout (e.g., 30s, 5m) (default: 30s)")
	fmt.Println("  -parallel               Run tests in parallel")
	fmt.Println("  -workers N              Maximum parallel workers (default: 4)")
	fmt.Println("  -bail                   Stop on first test failure")
	fmt.Println("  -retries N              Number of retries for failed tests")
	fmt.Println()
	fmt.Println("Test Filtering:")
	fmt.Println("  -grep PATTERN           Run only tests matching pattern")
	fmt.Println("  -tags TAGS              Run only tests with specified tags")
	fmt.Println("  -skip PATTERN           Skip tests matching pattern")
	fmt.Println("  -only PATTERN           Run only tests matching pattern (exclusive)")
	fmt.Println()
	fmt.Println("Coverage Options:")
	fmt.Println("  -coverage               Enable coverage collection")
	fmt.Println("  -coverage-threshold N   Coverage threshold percentage (default: 80)")
	fmt.Println("  -coverage-output DIR    Coverage output directory (default: ./coverage)")
	fmt.Println("  -coverage-formats LIST  Coverage formats: html,json (default: html)")
	fmt.Println("  -coverage-exclude LIST  Coverage exclude patterns")
	fmt.Println()
	fmt.Println("Reporting Options:")
	fmt.Println("  -reporters LIST         Reporters: console,json,junit (default: console)")
	fmt.Println("  -output DIR             Output directory for reports (default: ./test-results)")
	fmt.Println()
	fmt.Println("Advanced Options:")
	fmt.Println("  -watch                  Watch mode - rerun tests on file changes")
	fmt.Println("  -fixtures DIR           Fixture directory (default: ./fixtures)")
	fmt.Println("  -cleanup-mocks          Cleanup mocks after tests (default: true)")
	fmt.Println("  -isolation              Run tests in isolation contexts")
	fmt.Println()
	fmt.Println("Configuration File:")
	fmt.Println("  Use -config to specify a JSON configuration file with the following structure:")
	fmt.Println(`  {
    "testDirs": ["./tests"],
    "patterns": ["*_test.r2"],
    "timeout": "30s",
    "parallel": true,
    "coverage": {
      "enabled": true,
      "threshold": 80,
      "formats": ["html", "json"]
    }
  }`)
	fmt.Println()
	fmt.Println("For more information, visit: https://github.com/arturoeanton/go-r2lang")
}

func showVersion() {
	fmt.Printf("R2Test v%s\n", version)
	fmt.Println("Advanced Testing Framework for R2Lang")
	fmt.Println("Copyright (c) 2025 R2Lang Contributors")
	fmt.Println("Licensed under Apache License 2.0")
}

func showDryRun(config *Config) {
	fmt.Printf("R2Test v%s - Dry Run Mode\n\n", version)
	fmt.Println("Configuration that would be used:")
	fmt.Printf("  Test Directories: %v\n", config.TestDirs)
	fmt.Printf("  Patterns: %v\n", config.Patterns)
	fmt.Printf("  Ignore: %v\n", config.Ignore)
	fmt.Printf("  Timeout: %s\n", config.Timeout)
	fmt.Printf("  Parallel: %t\n", config.Parallel)
	if config.Parallel {
		fmt.Printf("  Max Workers: %d\n", config.MaxWorkers)
	}
	fmt.Printf("  Bail on Failure: %t\n", config.Bail)

	if config.Grep != "" {
		fmt.Printf("  Grep Pattern: %s\n", config.Grep)
	}
	if len(config.Tags) > 0 {
		fmt.Printf("  Tags: %v\n", config.Tags)
	}

	fmt.Printf("  Coverage Enabled: %t\n", config.Coverage.Enabled)
	if config.Coverage.Enabled {
		fmt.Printf("  Coverage Threshold: %.1f%%\n", config.Coverage.Threshold)
		fmt.Printf("  Coverage Output: %s\n", config.Coverage.Output)
		fmt.Printf("  Coverage Formats: %v\n", config.Coverage.Formats)
	}

	fmt.Printf("  Reporters: %v\n", config.Reporters)
	fmt.Printf("  Output Directory: %s\n", config.OutputDir)
	fmt.Printf("  Fixture Directory: %s\n", config.FixtureDir)
	fmt.Printf("  Verbose: %t\n", config.Verbose)

	fmt.Println("\nNo tests would be executed in dry-run mode.")
}
