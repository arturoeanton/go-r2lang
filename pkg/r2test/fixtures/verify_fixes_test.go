package fixtures

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// Regression test: Load(name) with an ambiguous, extension-less name must
// deterministically resolve to the same file across repeated calls/manager
// instances, instead of depending on Go's randomized map iteration order
// over the registered loader extensions.
func TestBugfix_LoadExtensionResolutionIsDeterministic(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "fixture_ambig")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	if err := ioutil.WriteFile(filepath.Join(tempDir, "data.json"), []byte(`{"from":"json"}`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(filepath.Join(tempDir, "data.txt"), []byte("from txt"), 0644); err != nil {
		t.Fatal(err)
	}

	var firstType string
	for i := 0; i < 30; i++ {
		fm := NewFixtureManager(tempDir)
		fixture, err := fm.Load("data")
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}
		if i == 0 {
			firstType = fixture.Type
		} else if fixture.Type != firstType {
			t.Fatalf("Load(\"data\") resolved inconsistently: run 0 got %q, run %d got %q", firstType, i, fixture.Type)
		}
	}
}
