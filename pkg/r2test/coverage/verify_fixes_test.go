package coverage

import (
	"math"
	"testing"
)

// Regression test: GetSortedFiles must not produce NaN for files with
// TotalLines == 0 (e.g. a function/statement registered but never hit by
// any recorded line), and must agree with GetFilePercentage() (0%), not
// silently sort them inconsistently via a raw 0/0 division.
func TestBugfix_GetSortedFilesHandlesZeroTotalLines(t *testing.T) {
	cc := NewCoverageCollector(".")
	cc.AddFunction("empty.go", "f", 1, 2) // no RecordLineHit -> TotalLines stays 0
	cc.RecordLineHit("normal.go", 1)

	files := cc.GetSortedFiles()
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}

	for _, f := range files {
		pct := f.GetFilePercentage()
		if math.IsNaN(pct) {
			t.Fatalf("file %q has NaN percentage", f.Path)
		}
	}

	// The empty file (0%) must sort before the fully covered file (100%).
	if files[0].Path != "empty.go" {
		t.Fatalf("expected empty.go (0%%) to sort first, got %q", files[0].Path)
	}
}
