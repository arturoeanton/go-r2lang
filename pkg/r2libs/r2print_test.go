package r2libs

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"testing"
	"unicode/utf8"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// TestPrintHeaderMultibyteWidth guards against printHeader sizing its
// separator line by byte length instead of rune count, which misaligns the
// header box for any non-ASCII text.
func TestPrintHeaderMultibyteWidth(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterPrint(env)
	mod, _ := env.Get("r2printer")
	printHeaderFunc := mod.(map[string]interface{})["printHeader"].(r2core.BuiltinFunction)

	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	text := "café mañana"
	printHeaderFunc(text)

	w.Close()
	os.Stdout = origStdout

	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines (separator, text, separator), got %d: %#v", len(lines), lines)
	}

	wantWidth := utf8.RuneCountInString(text)
	for _, sepLine := range []string{lines[0], lines[2]} {
		if got := utf8.RuneCountInString(sepLine); got != wantWidth {
			t.Errorf("separator width = %d runes, want %d (text %q, separator %q)", got, wantWidth, text, sepLine)
		}
	}
}
