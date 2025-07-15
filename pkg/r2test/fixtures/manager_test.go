package fixtures

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNewFixtureManager(t *testing.T) {
	fm := NewFixtureManager("./test_fixtures")

	if fm.basePath != "./test_fixtures" {
		t.Errorf("Expected base path './test_fixtures', got '%s'", fm.basePath)
	}

	if fm.fixtures == nil {
		t.Error("Expected fixtures map to be initialized")
	}

	if fm.loaders == nil {
		t.Error("Expected loaders map to be initialized")
	}

	// Check that default loaders are registered
	expectedExtensions := []string{".json", ".txt", ".md", ".r2", ".bin", ".dat", ".csv"}
	for _, ext := range expectedExtensions {
		if _, exists := fm.loaders[ext]; !exists {
			t.Errorf("Expected loader for extension '%s' to be registered", ext)
		}
	}
}

func TestFixtureManagerCreateTemporary(t *testing.T) {
	fm := NewFixtureManager("./test_fixtures")

	testData := map[string]interface{}{
		"name":  "test",
		"value": 123,
	}

	fixture, err := fm.CreateTemporary("temp_fixture", testData)
	if err != nil {
		t.Errorf("Unexpected error creating temporary fixture: %v", err)
	}

	if fixture.Name != "temp_fixture" {
		t.Errorf("Expected fixture name 'temp_fixture', got '%s'", fixture.Name)
	}

	if !fixture.isTemporary {
		t.Error("Expected fixture to be temporary")
	}

	// Compare the data using reflect.DeepEqual since we can't compare maps directly
	expectedName := testData["name"]
	expectedValue := testData["value"]

	fixtureData, ok := fixture.Data.(map[string]interface{})
	if !ok {
		t.Errorf("Expected fixture data to be map[string]interface{}, got %T", fixture.Data)
	} else {
		if fixtureData["name"] != expectedName {
			t.Errorf("Expected name %v, got %v", expectedName, fixtureData["name"])
		}
		if fixtureData["value"] != expectedValue {
			t.Errorf("Expected value %v, got %v", expectedValue, fixtureData["value"])
		}
	}

	// Verify it can be retrieved
	retrieved, exists := fm.Get("temp_fixture")
	if !exists {
		t.Error("Expected to be able to retrieve temporary fixture")
	}

	if retrieved != fixture {
		t.Error("Expected retrieved fixture to be the same instance")
	}
}

func TestFixtureManagerLoadJSON(t *testing.T) {
	// Create temporary directory for test fixtures
	tempDir, err := ioutil.TempDir("", "fixture_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	fm := NewFixtureManager(tempDir)

	// Create a test JSON file
	jsonContent := `{"name": "test", "value": 123, "nested": {"key": "value"}}`
	jsonFile := filepath.Join(tempDir, "test.json")
	err = ioutil.WriteFile(jsonFile, []byte(jsonContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test JSON file: %v", err)
	}

	// Load the fixture
	fixture, err := fm.Load("test.json")
	if err != nil {
		t.Errorf("Unexpected error loading JSON fixture: %v", err)
	}

	if fixture.Name != "test.json" {
		t.Errorf("Expected fixture name 'test.json', got '%s'", fixture.Name)
	}

	if fixture.Type != ".json" {
		t.Errorf("Expected fixture type '.json', got '%s'", fixture.Type)
	}

	// Verify the data was parsed correctly
	data, ok := fixture.Data.(map[string]interface{})
	if !ok {
		t.Errorf("Expected fixture data to be map[string]interface{}, got %T", fixture.Data)
	}

	if data["name"] != "test" {
		t.Errorf("Expected name to be 'test', got '%v'", data["name"])
	}

	if data["value"].(float64) != 123 {
		t.Errorf("Expected value to be 123, got '%v'", data["value"])
	}
}

func TestFixtureManagerLoadText(t *testing.T) {
	// Create temporary directory for test fixtures
	tempDir, err := ioutil.TempDir("", "fixture_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	fm := NewFixtureManager(tempDir)

	// Create a test text file
	textContent := "This is a test text file\nwith multiple lines\nfor testing."
	textFile := filepath.Join(tempDir, "test.txt")
	err = ioutil.WriteFile(textFile, []byte(textContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test text file: %v", err)
	}

	// Load the fixture
	fixture, err := fm.Load("test.txt")
	if err != nil {
		t.Errorf("Unexpected error loading text fixture: %v", err)
	}

	if fixture.Type != ".txt" {
		t.Errorf("Expected fixture type '.txt', got '%s'", fixture.Type)
	}

	// Verify the data was loaded correctly
	data, ok := fixture.Data.(string)
	if !ok {
		t.Errorf("Expected fixture data to be string, got %T", fixture.Data)
	}

	if data != textContent {
		t.Errorf("Expected text content to match")
	}
}

func TestFixtureManagerLoadCSV(t *testing.T) {
	// Create temporary directory for test fixtures
	tempDir, err := ioutil.TempDir("", "fixture_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	fm := NewFixtureManager(tempDir)

	// Create a test CSV file
	csvContent := "name,age,city\nJohn,25,New York\nJane,30,Los Angeles\nBob,35,Chicago"
	csvFile := filepath.Join(tempDir, "test.csv")
	err = ioutil.WriteFile(csvFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}

	// Load the fixture
	fixture, err := fm.Load("test.csv")
	if err != nil {
		t.Errorf("Unexpected error loading CSV fixture: %v", err)
	}

	if fixture.Type != ".csv" {
		t.Errorf("Expected fixture type '.csv', got '%s'", fixture.Type)
	}

	// Verify the data was parsed correctly
	data, ok := fixture.Data.([]map[string]string)
	if !ok {
		t.Errorf("Expected fixture data to be []map[string]string, got %T", fixture.Data)
	}

	if len(data) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(data))
	}

	// Check first row
	if data[0]["name"] != "John" || data[0]["age"] != "25" || data[0]["city"] != "New York" {
		t.Errorf("First row data incorrect: %v", data[0])
	}

	// Check second row
	if data[1]["name"] != "Jane" || data[1]["age"] != "30" || data[1]["city"] != "Los Angeles" {
		t.Errorf("Second row data incorrect: %v", data[1])
	}
}

func TestFixtureManagerLoadNonexistent(t *testing.T) {
	fm := NewFixtureManager("./nonexistent")

	_, err := fm.Load("nonexistent.json")
	if err == nil {
		t.Error("Expected error loading nonexistent fixture")
	}

	expectedMsg := "fixture file not found"
	if !contains(err.Error(), expectedMsg) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestFixtureManagerGetData(t *testing.T) {
	fm := NewFixtureManager("./test_fixtures")

	// Create temporary fixture
	testData := "test data"
	fm.CreateTemporary("test_data", testData)

	// Get data
	data, err := fm.GetData("test_data")
	if err != nil {
		t.Errorf("Unexpected error getting fixture data: %v", err)
	}

	if data != testData {
		t.Errorf("Expected data '%s', got '%v'", testData, data)
	}
}

func TestFixtureManagerList(t *testing.T) {
	fm := NewFixtureManager("./test_fixtures")

	// Initially should be empty
	list := fm.List()
	if len(list) != 0 {
		t.Errorf("Expected empty list initially, got %d items", len(list))
	}

	// Add some fixtures
	fm.CreateTemporary("fixture1", "data1")
	fm.CreateTemporary("fixture2", "data2")
	fm.CreateTemporary("fixture3", "data3")

	// Should now have 3 items
	list = fm.List()
	if len(list) != 3 {
		t.Errorf("Expected 3 items in list, got %d", len(list))
	}

	// Check that all fixture names are present
	expectedNames := []string{"fixture1", "fixture2", "fixture3"}
	for _, expected := range expectedNames {
		found := false
		for _, name := range list {
			if name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find fixture '%s' in list", expected)
		}
	}
}

func TestFixtureManagerClear(t *testing.T) {
	fm := NewFixtureManager("./test_fixtures")

	// Add some fixtures
	fm.CreateTemporary("fixture1", "data1")
	fm.CreateTemporary("fixture2", "data2")

	// Verify they exist
	list := fm.List()
	if len(list) != 2 {
		t.Errorf("Expected 2 fixtures before clear, got %d", len(list))
	}

	// Clear
	fm.Clear()

	// Verify they're gone
	list = fm.List()
	if len(list) != 0 {
		t.Errorf("Expected 0 fixtures after clear, got %d", len(list))
	}
}

func TestFixtureManagerUnload(t *testing.T) {
	fm := NewFixtureManager("./test_fixtures")

	// Add a fixture
	fm.CreateTemporary("test_fixture", "test_data")

	// Verify it exists
	_, exists := fm.Get("test_fixture")
	if !exists {
		t.Error("Expected fixture to exist before unload")
	}

	// Unload
	fm.Unload("test_fixture")

	// Verify it's gone
	_, exists = fm.Get("test_fixture")
	if exists {
		t.Error("Expected fixture to be gone after unload")
	}
}

func TestGlobalFixtureManager(t *testing.T) {
	// Test global convenience functions

	// Create temporary fixture
	fixture, err := CreateTemporaryFixture("global_test", "global_data")
	if err != nil {
		t.Errorf("Unexpected error creating global temporary fixture: %v", err)
	}

	if fixture.Name != "global_test" {
		t.Errorf("Expected fixture name 'global_test', got '%s'", fixture.Name)
	}

	// Retrieve using global function
	retrieved, exists := GetFixture("global_test")
	if !exists {
		t.Error("Expected to find global fixture")
	}

	if retrieved != fixture {
		t.Error("Expected retrieved fixture to be the same instance")
	}

	// Get data using global function
	data, err := GetFixtureData("global_test")
	if err != nil {
		t.Errorf("Unexpected error getting global fixture data: %v", err)
	}

	if data != "global_data" {
		t.Errorf("Expected data 'global_data', got '%v'", data)
	}

	// Clean up
	ClearFixtures()
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
