package reporters

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2test/core"
)

// Regression test: GenerateJUnitReportWithProperties previously called the
// base Generate() and then did nothing else, silently discarding the
// caller-supplied properties while still returning nil (success). It must
// actually embed the properties into each <testsuite> element.
func TestBugfix_JUnitPropertiesAreEmbedded(t *testing.T) {
	suite := &core.TestSuite{Name: "MySuite"}
	testCase := &core.TestCase{Name: "does a thing", Suite: suite}
	results := &core.TestResults{
		Results: []*core.TestResult{
			{
				TestCase:  testCase,
				Suite:     suite,
				Status:    core.TestStatusPassed,
				StartTime: time.Now(),
				EndTime:   time.Now(),
			},
		},
	}

	outPath := filepath.Join(t.TempDir(), "out.xml")
	if err := GenerateJUnitReportWithProperties(outPath, results, map[string]string{"ci": "true", "branch": "main"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}

	xml := string(data)
	if !strings.Contains(xml, `name="ci"`) || !strings.Contains(xml, `value="true"`) {
		t.Fatalf("expected property 'ci=true' embedded in JUnit XML, got:\n%s", xml)
	}
	if !strings.Contains(xml, `name="branch"`) || !strings.Contains(xml, `value="main"`) {
		t.Fatalf("expected property 'branch=main' embedded in JUnit XML, got:\n%s", xml)
	}
}
