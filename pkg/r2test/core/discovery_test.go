package core

import (
	"os"
	"path/filepath"
	"testing"
)

// writeTempTestFile writes content to a .r2 file inside a fresh temp
// directory and returns its path.
func writeTempTestFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample_test.r2")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp test file: %v", err)
	}
	return path
}

// TestParseTestFile_ActuallyExecutesTestBodies guards against the original
// bug this file replaced: ParseTestFile used to be a regex-based
// placeholder whose TestCase.Func was a no-op print statement — every
// discovered test "passed" regardless of what its body actually did. This
// confirms a genuinely failing assertion inside it()'s body is detected as
// a real failure by TestRunner, not silently reported as passing.
func TestParseTestFile_ActuallyExecutesTestBodies(t *testing.T) {
	path := writeTempTestFile(t, `
describe("real execution", func() {
    it("really fails", func() {
        assert.equals(1, 2);
    });
    it("really passes", func() {
        assert.equals(1, 1);
    });
});
`)

	td := NewTestDiscovery(DefaultConfig())
	suites, err := td.ParseTestFile(path)
	if err != nil {
		t.Fatalf("ParseTestFile returned an error: %v", err)
	}
	if len(suites) != 1 {
		t.Fatalf("expected 1 suite, got %d", len(suites))
	}
	if len(suites[0].Tests) != 2 {
		t.Fatalf("expected 2 tests, got %d", len(suites[0].Tests))
	}

	runner := NewTestRunner(DefaultConfig())
	runner.AddSuite(suites[0])
	results, err := runner.Run()
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	stats := results.GetStats()
	if stats.Passed != 1 || stats.Failed != 1 {
		t.Fatalf("expected 1 passed / 1 failed, got %d passed / %d failed", stats.Passed, stats.Failed)
	}
}

// TestParseTestFile_MultipleDescribeBlocks confirms a single file with
// several top-level describe() blocks (the normal shape of a real test
// file — see examples/testing/basic_test.r2) produces one TestSuite per
// describe(), not just the first/last one.
func TestParseTestFile_MultipleDescribeBlocks(t *testing.T) {
	path := writeTempTestFile(t, `
describe("suite one", func() {
    it("test one", func() { assert.true(true); });
});
describe("suite two", func() {
    it("test two", func() { assert.true(true); });
    it("test three", func() { assert.true(true); });
});
`)

	td := NewTestDiscovery(DefaultConfig())
	suites, err := td.ParseTestFile(path)
	if err != nil {
		t.Fatalf("ParseTestFile returned an error: %v", err)
	}
	if len(suites) != 2 {
		t.Fatalf("expected 2 suites, got %d", len(suites))
	}
	if suites[0].Name != "suite one" || len(suites[0].Tests) != 1 {
		t.Errorf("unexpected first suite: name=%q tests=%d", suites[0].Name, len(suites[0].Tests))
	}
	if suites[1].Name != "suite two" || len(suites[1].Tests) != 2 {
		t.Errorf("unexpected second suite: name=%q tests=%d", suites[1].Name, len(suites[1].Tests))
	}
}

// TestParseTestFile_HooksRunForReal confirms beforeEach/afterEach actually
// execute around each test body (not just get registered and ignored).
func TestParseTestFile_HooksRunForReal(t *testing.T) {
	path := writeTempTestFile(t, `
let log = [];
describe("hooks", func() {
    beforeEach(func() { log = log.push("before"); });
    afterEach(func() { log = log.push("after"); });
    it("first", func() { log = log.push("test"); });
});
`)

	td := NewTestDiscovery(DefaultConfig())
	suites, err := td.ParseTestFile(path)
	if err != nil {
		t.Fatalf("ParseTestFile returned an error: %v", err)
	}
	if len(suites) != 1 || len(suites[0].Tests) != 1 {
		t.Fatalf("expected 1 suite with 1 test, got %d suites", len(suites))
	}

	runner := NewTestRunner(DefaultConfig())
	runner.AddSuite(suites[0])
	results, err := runner.Run()
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}
	if results.GetStats().Failed != 0 {
		t.Fatalf("expected the test to pass, got failures: %+v", results.Results)
	}
}

// TestParseTestFile_MalformedFileReturnsError confirms a script that fails
// to parse/evaluate is reported as a discovery error rather than silently
// producing zero suites (which LoadTestSuites would otherwise treat as "no
// tests found here", masking a real problem in the file).
func TestParseTestFile_MalformedFileReturnsError(t *testing.T) {
	path := writeTempTestFile(t, `describe("broken", func() {`) // unclosed block
	td := NewTestDiscovery(DefaultConfig())
	if _, err := td.ParseTestFile(path); err == nil {
		t.Fatal("expected an error for a malformed test file, got none")
	}
}
