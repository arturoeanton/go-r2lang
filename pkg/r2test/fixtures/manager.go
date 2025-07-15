package fixtures

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// FixtureManager manages test fixtures and their lifecycle
type FixtureManager struct {
	basePath string
	fixtures map[string]*Fixture
	mu       sync.RWMutex
	loaders  map[string]FixtureLoader
	cleanups []func() error
}

// Fixture represents a single test fixture
type Fixture struct {
	Name        string
	Path        string
	Data        interface{}
	LoadedAt    time.Time
	Type        string
	Size        int64
	Metadata    map[string]interface{}
	isTemporary bool
}

// FixtureLoader interface for loading different fixture types
type FixtureLoader interface {
	Load(path string) (interface{}, error)
	Extensions() []string
}

// JSONLoader loads JSON fixtures
type JSONLoader struct{}

func (jl *JSONLoader) Load(path string) (interface{}, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON fixture %s: %w", path, err)
	}

	var result interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON fixture %s: %w", path, err)
	}

	return result, nil
}

func (jl *JSONLoader) Extensions() []string {
	return []string{".json"}
}

// TextLoader loads plain text fixtures
type TextLoader struct{}

func (tl *TextLoader) Load(path string) (interface{}, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read text fixture %s: %w", path, err)
	}
	return string(data), nil
}

func (tl *TextLoader) Extensions() []string {
	return []string{".txt", ".md", ".r2"}
}

// BinaryLoader loads binary fixtures
type BinaryLoader struct{}

func (bl *BinaryLoader) Load(path string) (interface{}, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read binary fixture %s: %w", path, err)
	}
	return data, nil
}

func (bl *BinaryLoader) Extensions() []string {
	return []string{".bin", ".dat"}
}

// CSVLoader loads CSV fixtures as arrays of maps
type CSVLoader struct{}

func (cl *CSVLoader) Load(path string) (interface{}, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV fixture %s: %w", path, err)
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) < 2 {
		return []map[string]string{}, nil
	}

	headers := strings.Split(lines[0], ",")
	for i := range headers {
		headers[i] = strings.TrimSpace(headers[i])
	}

	var result []map[string]string
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		values := strings.Split(line, ",")
		row := make(map[string]string)

		for j, header := range headers {
			if j < len(values) {
				row[header] = strings.TrimSpace(values[j])
			} else {
				row[header] = ""
			}
		}
		result = append(result, row)
	}

	return result, nil
}

func (cl *CSVLoader) Extensions() []string {
	return []string{".csv"}
}

// NewFixtureManager creates a new fixture manager
func NewFixtureManager(basePath string) *FixtureManager {
	fm := &FixtureManager{
		basePath: basePath,
		fixtures: make(map[string]*Fixture),
		loaders:  make(map[string]FixtureLoader),
		cleanups: make([]func() error, 0),
	}

	// Register default loaders
	fm.RegisterLoader(&JSONLoader{})
	fm.RegisterLoader(&TextLoader{})
	fm.RegisterLoader(&BinaryLoader{})
	fm.RegisterLoader(&CSVLoader{})

	return fm
}

// RegisterLoader registers a fixture loader for specific file extensions
func (fm *FixtureManager) RegisterLoader(loader FixtureLoader) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	for _, ext := range loader.Extensions() {
		fm.loaders[ext] = loader
	}
}

// Load loads a fixture by name
func (fm *FixtureManager) Load(name string) (*Fixture, error) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	// Check if already loaded
	if fixture, exists := fm.fixtures[name]; exists {
		return fixture, nil
	}

	// Find fixture file
	fixturePath := filepath.Join(fm.basePath, name)

	// Try different extensions if no extension provided
	if filepath.Ext(name) == "" {
		for ext := range fm.loaders {
			testPath := fixturePath + ext
			if _, err := os.Stat(testPath); err == nil {
				fixturePath = testPath
				break
			}
		}
	}

	// Check if file exists
	fileInfo, err := os.Stat(fixturePath)
	if err != nil {
		return nil, fmt.Errorf("fixture file not found: %s", fixturePath)
	}

	// Determine loader by extension
	ext := filepath.Ext(fixturePath)
	loader, exists := fm.loaders[ext]
	if !exists {
		return nil, fmt.Errorf("no loader registered for extension: %s", ext)
	}

	// Load fixture data
	data, err := loader.Load(fixturePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load fixture %s: %w", name, err)
	}

	// Create fixture
	fixture := &Fixture{
		Name:     name,
		Path:     fixturePath,
		Data:     data,
		LoadedAt: time.Now(),
		Type:     ext,
		Size:     fileInfo.Size(),
		Metadata: make(map[string]interface{}),
	}

	fm.fixtures[name] = fixture
	return fixture, nil
}

// LoadFromPath loads a fixture from a specific path
func (fm *FixtureManager) LoadFromPath(path string) (*Fixture, error) {
	name := filepath.Base(path)
	return fm.loadFromFullPath(name, path)
}

// loadFromFullPath loads a fixture from a full path
func (fm *FixtureManager) loadFromFullPath(name, fullPath string) (*Fixture, error) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	// Check if file exists
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("fixture file not found: %s", fullPath)
	}

	// Determine loader by extension
	ext := filepath.Ext(fullPath)
	loader, exists := fm.loaders[ext]
	if !exists {
		return nil, fmt.Errorf("no loader registered for extension: %s", ext)
	}

	// Load fixture data
	data, err := loader.Load(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load fixture %s: %w", name, err)
	}

	// Create fixture
	fixture := &Fixture{
		Name:     name,
		Path:     fullPath,
		Data:     data,
		LoadedAt: time.Now(),
		Type:     ext,
		Size:     fileInfo.Size(),
		Metadata: make(map[string]interface{}),
	}

	fm.fixtures[name] = fixture
	return fixture, nil
}

// Get retrieves a loaded fixture by name
func (fm *FixtureManager) Get(name string) (*Fixture, bool) {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	fixture, exists := fm.fixtures[name]
	return fixture, exists
}

// GetData retrieves fixture data by name, loading if necessary
func (fm *FixtureManager) GetData(name string) (interface{}, error) {
	fixture, err := fm.Load(name)
	if err != nil {
		return nil, err
	}
	return fixture.Data, nil
}

// CreateTemporary creates a temporary fixture with given data
func (fm *FixtureManager) CreateTemporary(name string, data interface{}) (*Fixture, error) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	fixture := &Fixture{
		Name:        name,
		Path:        "", // No file path for temporary fixtures
		Data:        data,
		LoadedAt:    time.Now(),
		Type:        "temp",
		Size:        0,
		Metadata:    make(map[string]interface{}),
		isTemporary: true,
	}

	fm.fixtures[name] = fixture
	return fixture, nil
}

// SaveTemporary saves a temporary fixture to disk
func (fm *FixtureManager) SaveTemporary(name string, subPath string) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	fixture, exists := fm.fixtures[name]
	if !exists {
		return fmt.Errorf("fixture not found: %s", name)
	}

	if !fixture.isTemporary {
		return fmt.Errorf("fixture %s is not temporary", name)
	}

	// Create directory if it doesn't exist
	fullPath := filepath.Join(fm.basePath, subPath)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Save based on data type
	var err error
	switch data := fixture.Data.(type) {
	case string:
		err = ioutil.WriteFile(fullPath, []byte(data), 0644)
	case []byte:
		err = ioutil.WriteFile(fullPath, data, 0644)
	default:
		// Try to marshal as JSON
		jsonData, marshalErr := json.MarshalIndent(data, "", "  ")
		if marshalErr != nil {
			return fmt.Errorf("failed to marshal fixture data: %w", marshalErr)
		}
		err = ioutil.WriteFile(fullPath, jsonData, 0644)
	}

	if err != nil {
		return fmt.Errorf("failed to save fixture to %s: %w", fullPath, err)
	}

	// Update fixture
	fixture.Path = fullPath
	fixture.isTemporary = false
	fixture.Type = filepath.Ext(fullPath)

	return nil
}

// List returns all loaded fixture names
func (fm *FixtureManager) List() []string {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	names := make([]string, 0, len(fm.fixtures))
	for name := range fm.fixtures {
		names = append(names, name)
	}
	return names
}

// Unload removes a fixture from memory
func (fm *FixtureManager) Unload(name string) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	delete(fm.fixtures, name)
}

// Clear removes all fixtures from memory
func (fm *FixtureManager) Clear() {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	fm.fixtures = make(map[string]*Fixture)
}

// RegisterCleanup registers a cleanup function to be called when cleaning up
func (fm *FixtureManager) RegisterCleanup(cleanup func() error) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	fm.cleanups = append(fm.cleanups, cleanup)
}

// Cleanup runs all registered cleanup functions
func (fm *FixtureManager) Cleanup() error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	for _, cleanup := range fm.cleanups {
		if err := cleanup(); err != nil {
			return err
		}
	}

	fm.cleanups = make([]func() error, 0)
	return nil
}

// SetBasePath sets the base path for fixtures
func (fm *FixtureManager) SetBasePath(path string) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	fm.basePath = path
}

// GetBasePath returns the current base path
func (fm *FixtureManager) GetBasePath() string {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	return fm.basePath
}

// Global fixture manager
var GlobalFixtureManager = NewFixtureManager("./fixtures")

// Global convenience functions

// LoadFixture loads a fixture using the global manager
func LoadFixture(name string) (*Fixture, error) {
	return GlobalFixtureManager.Load(name)
}

// GetFixture retrieves a fixture using the global manager
func GetFixture(name string) (*Fixture, bool) {
	return GlobalFixtureManager.Get(name)
}

// GetFixtureData retrieves fixture data using the global manager
func GetFixtureData(name string) (interface{}, error) {
	return GlobalFixtureManager.GetData(name)
}

// CreateTemporaryFixture creates a temporary fixture using the global manager
func CreateTemporaryFixture(name string, data interface{}) (*Fixture, error) {
	return GlobalFixtureManager.CreateTemporary(name, data)
}

// SetFixtureBasePath sets the base path for the global manager
func SetFixtureBasePath(path string) {
	GlobalFixtureManager.SetBasePath(path)
}

// ClearFixtures clears all fixtures from the global manager
func ClearFixtures() {
	GlobalFixtureManager.Clear()
}
