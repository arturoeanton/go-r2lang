package reporters

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arturoeanton/go-r2lang/pkg/r2test/coverage"
)

// TestHTMLReporter_GenerateDoesNotPanic is a regression test: the main and
// file templates call {{sub .TotalLines .CoveredLines}}, but "sub" was never
// registered via template.Funcs, so html/template.Must panicked at parse
// time with "function \"sub\" not defined" every time a report was
// generated (e.g. `r2test -coverage`, which defaults to the html format).
func TestHTMLReporter_GenerateDoesNotPanic(t *testing.T) {
	collector := coverage.NewCoverageCollector(".")
	collector.RecordLineHit("example.r2", 1)
	collector.RecordLineHit("example.r2", 2)

	outputDir := filepath.Join(t.TempDir(), "coverage-html")
	reporter := NewHTMLReporter(outputDir)

	if err := reporter.Generate(collector, nil); err != nil {
		t.Fatalf("Generate should not error: %v", err)
	}

	indexPath := filepath.Join(outputDir, "index.html")
	if _, err := os.Stat(indexPath); err != nil {
		t.Fatalf("expected index.html to be created: %v", err)
	}
}
