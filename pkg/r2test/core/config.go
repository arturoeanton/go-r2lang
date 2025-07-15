package core

import (
	"encoding/json"
	"os"
	"time"
)

// TestConfig holds configuration for the test framework
type TestConfig struct {
	// Test Discovery
	TestDirs       []string `json:"testDirs"`
	TestPatterns   []string `json:"patterns"`
	IgnorePatterns []string `json:"ignore"`

	// Execution
	DefaultTimeout time.Duration `json:"timeout"`
	Parallel       bool          `json:"parallel"`
	MaxWorkers     int           `json:"maxWorkers"`
	Bail           bool          `json:"bail"` // Stop on first failure

	// Filtering
	Grep       string   `json:"grep"`
	FilterTags []string `json:"tags"`
	Skip       []string `json:"skip"`
	Only       []string `json:"only"`

	// Coverage
	Coverage CoverageConfig `json:"coverage"`

	// Reporting
	Reporters []string `json:"reporters"`
	Verbose   bool     `json:"verbose"`
	Output    string   `json:"output"`

	// Advanced
	Retries   int  `json:"retries"`
	WatchMode bool `json:"watchMode"`

	// Snapshot testing
	Snapshot SnapshotConfig `json:"snapshot"`
}

// CoverageConfig holds coverage-related configuration
type CoverageConfig struct {
	Enabled   bool     `json:"enabled"`
	Threshold float64  `json:"threshold"`
	Output    string   `json:"output"`
	Formats   []string `json:"formats"`
}

// SnapshotConfig holds snapshot testing configuration
type SnapshotConfig struct {
	UpdateSnapshots bool   `json:"updateSnapshots"`
	SnapshotDir     string `json:"snapshotDir"`
}

// DefaultConfig returns a default test configuration
func DefaultConfig() *TestConfig {
	return &TestConfig{
		// Test Discovery
		TestDirs:       []string{"./tests", "./test"},
		TestPatterns:   []string{"*_test.r2", "*Test.r2", "test_*.r2"},
		IgnorePatterns: []string{"node_modules", "vendor", "*.tmp.r2", ".*"},

		// Execution
		DefaultTimeout: 30 * time.Second,
		Parallel:       false,
		MaxWorkers:     4,
		Bail:           false,

		// Filtering
		Grep:       "",
		FilterTags: []string{},
		Skip:       []string{},
		Only:       []string{},

		// Coverage
		Coverage: CoverageConfig{
			Enabled:   false,
			Threshold: 80.0,
			Output:    "./coverage",
			Formats:   []string{"console"},
		},

		// Reporting
		Reporters: []string{"console"},
		Verbose:   false,
		Output:    "",

		// Advanced
		Retries:   0,
		WatchMode: false,

		// Snapshot
		Snapshot: SnapshotConfig{
			UpdateSnapshots: false,
			SnapshotDir:     "./tests/__snapshots__",
		},
	}
}

// LoadConfig loads configuration from a file, merging with defaults
func LoadConfig(configPath string) (*TestConfig, error) {
	config := DefaultConfig()

	if configPath == "" {
		// Try to find config file automatically
		configPath = findConfigFile()
	}

	if configPath != "" {
		if err := config.loadFromFile(configPath); err != nil {
			return nil, err
		}
	}

	// Apply environment variable overrides
	config.applyEnvironmentOverrides()

	// Validate configuration
	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// loadFromFile loads configuration from a JSON file
func (tc *TestConfig) loadFromFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, tc)
}

// findConfigFile looks for configuration files in common locations
func findConfigFile() string {
	possiblePaths := []string{
		"r2test.config.json",
		"r2test.json",
		".r2testrc.json",
		"test/config.json",
		"tests/config.json",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

// applyEnvironmentOverrides applies environment variable overrides
func (tc *TestConfig) applyEnvironmentOverrides() {
	// R2TEST_TIMEOUT
	if timeoutStr := os.Getenv("R2TEST_TIMEOUT"); timeoutStr != "" {
		if timeout, err := time.ParseDuration(timeoutStr); err == nil {
			tc.DefaultTimeout = timeout
		}
	}

	// R2TEST_PARALLEL
	if parallelStr := os.Getenv("R2TEST_PARALLEL"); parallelStr == "true" {
		tc.Parallel = true
	} else if parallelStr == "false" {
		tc.Parallel = false
	}

	// R2TEST_VERBOSE
	if verboseStr := os.Getenv("R2TEST_VERBOSE"); verboseStr == "true" {
		tc.Verbose = true
	}

	// R2TEST_COVERAGE
	if coverageStr := os.Getenv("R2TEST_COVERAGE"); coverageStr == "true" {
		tc.Coverage.Enabled = true
	}

	// R2TEST_GREP
	if grep := os.Getenv("R2TEST_GREP"); grep != "" {
		tc.Grep = grep
	}

	// R2TEST_TAGS
	if tags := os.Getenv("R2TEST_TAGS"); tags != "" {
		tc.FilterTags = []string{tags} // Simple implementation, could be comma-separated
	}
}

// validate checks that the configuration is valid
func (tc *TestConfig) validate() error {
	// Ensure we have at least one test directory
	if len(tc.TestDirs) == 0 {
		tc.TestDirs = []string{"./tests"}
	}

	// Ensure we have test patterns
	if len(tc.TestPatterns) == 0 {
		tc.TestPatterns = []string{"*_test.r2"}
	}

	// Ensure timeout is reasonable
	if tc.DefaultTimeout <= 0 {
		tc.DefaultTimeout = 30 * time.Second
	}

	// Ensure max workers is reasonable
	if tc.MaxWorkers <= 0 {
		tc.MaxWorkers = 4
	}

	// Ensure coverage threshold is valid
	if tc.Coverage.Threshold < 0 || tc.Coverage.Threshold > 100 {
		tc.Coverage.Threshold = 80.0
	}

	// Create output directories if they don't exist
	if tc.Coverage.Output != "" {
		os.MkdirAll(tc.Coverage.Output, 0755)
	}

	if tc.Snapshot.SnapshotDir != "" {
		os.MkdirAll(tc.Snapshot.SnapshotDir, 0755)
	}

	return nil
}

// SaveConfig saves the configuration to a file
func (tc *TestConfig) SaveConfig(filePath string) error {
	data, err := json.MarshalIndent(tc, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// Merge merges another configuration into this one
func (tc *TestConfig) Merge(other *TestConfig) {
	if other == nil {
		return
	}

	// Merge test directories
	if len(other.TestDirs) > 0 {
		tc.TestDirs = other.TestDirs
	}

	// Merge patterns
	if len(other.TestPatterns) > 0 {
		tc.TestPatterns = other.TestPatterns
	}

	// Merge other fields as needed
	if other.DefaultTimeout > 0 {
		tc.DefaultTimeout = other.DefaultTimeout
	}

	if other.MaxWorkers > 0 {
		tc.MaxWorkers = other.MaxWorkers
	}

	// Override boolean flags
	tc.Parallel = other.Parallel
	tc.Bail = other.Bail
	tc.Verbose = other.Verbose
	tc.WatchMode = other.WatchMode

	// Merge coverage config
	if other.Coverage.Enabled {
		tc.Coverage = other.Coverage
	}

	// Merge snapshot config
	if other.Snapshot.SnapshotDir != "" {
		tc.Snapshot = other.Snapshot
	}
}

// Clone creates a copy of the configuration
func (tc *TestConfig) Clone() *TestConfig {
	data, _ := json.Marshal(tc)
	var clone TestConfig
	json.Unmarshal(data, &clone)
	return &clone
}
