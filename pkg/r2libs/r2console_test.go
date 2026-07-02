package r2libs

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"regexp"
	"strconv"
	"sync"
	"testing"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// TestConsoleLogConcurrentWritesAreNotTorn guards against console.log
// interleaving output from concurrent callers: each call must produce a
// single, uninterrupted line on stdout.
func TestConsoleLogConcurrentWritesAreNotTorn(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterConsole(env)
	mod, _ := env.Get("console")
	logFunc := mod.(map[string]interface{})["log"].(r2core.BuiltinFunction)

	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	const goroutines = 20
	const callsPerGoroutine = 50

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		go func(id int) {
			defer wg.Done()
			marker := "ID" + strconv.Itoa(id)
			for i := 0; i < callsPerGoroutine; i++ {
				logFunc(marker, marker, marker, marker, marker, marker)
			}
		}(g)
	}
	wg.Wait()

	w.Close()
	os.Stdout = origStdout

	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	lineRe := regexp.MustCompile(`^\[\d\d:\d\d:\d\d\] (ID\d+)( ID\d+){5}$`)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	total := 0
	for scanner.Scan() {
		line := scanner.Text()
		total++
		m := lineRe.FindStringSubmatch(line)
		if m == nil {
			t.Fatalf("torn/interleaved console.log output detected: %q", line)
		}
		expected := m[1]
		if !allTokensMatch(line, expected) {
			t.Fatalf("mixed markers from different goroutines in one line: %q", line)
		}
	}
	if total != goroutines*callsPerGoroutine {
		t.Fatalf("expected %d lines, got %d", goroutines*callsPerGoroutine, total)
	}
}

func allTokensMatch(line, expected string) bool {
	re := regexp.MustCompile(`ID\d+`)
	for _, tok := range re.FindAllString(line, -1) {
		if tok != expected {
			return false
		}
	}
	return true
}
