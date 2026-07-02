package r2libs

import (
	"testing"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func TestUnicodeBasicFunctions(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterUnicode(env)

	unicodeModuleObj, ok := env.Get("unicode")
	if !ok {
		t.Fatal("unicode module not found")
	}
	unicodeModule := unicodeModuleObj.(map[string]interface{})

	ulenFunc := unicodeModule["ulen"].(r2core.BuiltinFunction)
	usubstrFunc := unicodeModule["usubstr"].(r2core.BuiltinFunction)
	uupperFunc := unicodeModule["uupper"].(r2core.BuiltinFunction)
	ulowerFunc := unicodeModule["ulower"].(r2core.BuiltinFunction)
	ureverseFunc := unicodeModule["ureverse"].(r2core.BuiltinFunction)

	tests := []struct {
		name     string
		function r2core.BuiltinFunction
		args     []interface{}
		expected interface{}
	}{
		{
			"ulen with ASCII",
			ulenFunc,
			[]interface{}{"hello"},
			float64(5),
		},
		{
			"ulen with Spanish characters",
			ulenFunc,
			[]interface{}{"José María"},
			float64(10),
		},
		{
			"ulen with emoji",
			ulenFunc,
			[]interface{}{"👋"},
			float64(1),
		},
		{
			"usubstr basic",
			usubstrFunc,
			[]interface{}{"hello", float64(1), float64(3)},
			"ell",
		},
		{
			"usubstr with Spanish",
			usubstrFunc,
			[]interface{}{"José María", float64(0), float64(4)},
			"José",
		},
		{
			"usubstr with emoji",
			usubstrFunc,
			[]interface{}{"Hello 👋 World", float64(6), float64(1)},
			"👋",
		},
		{
			"uupper with ASCII",
			uupperFunc,
			[]interface{}{"hello"},
			"HELLO",
		},
		{
			"uupper with Spanish",
			uupperFunc,
			[]interface{}{"josé"},
			"JOSÉ",
		},
		{
			"ulower with ASCII",
			ulowerFunc,
			[]interface{}{"HELLO"},
			"hello",
		},
		{
			"ulower with Spanish",
			ulowerFunc,
			[]interface{}{"JOSÉ"},
			"josé",
		},
		{
			"ureverse with ASCII",
			ureverseFunc,
			[]interface{}{"hello"},
			"olleh",
		},
		{
			"ureverse with Spanish",
			ureverseFunc,
			[]interface{}{"José"},
			"ésoJ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function(tt.args...)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestUnicodeGetCategory_Deterministic is a regression test for a bug where
// ugetCategory() iterated Go's unicode.Categories map (whose iteration order
// is randomized) and returned the first category table that matched,
// meaning the very same rune could non-deterministically return either its
// general category ("N", "L", ...) or its more specific subcategory ("Nd",
// "Lu", ...) depending on map iteration order - even across repeated calls
// within a single process run. Verified with a direct Go loop over
// unicode.Categories before the fix (seen both "N" and "Nd" for '5' within
// one run). The fix always returns the most specific (longest) matching
// category name, with alphabetical tie-breaking, so the result is stable.
func TestUnicodeGetCategory_Deterministic(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterUnicode(env)
	unicodeModuleObj, _ := env.Get("unicode")
	unicodeModule := unicodeModuleObj.(map[string]interface{})
	getCategoryFunc := unicodeModule["ugetCategory"].(r2core.BuiltinFunction)

	cases := map[string]string{
		"5": "Nd", // decimal digit: matches both "N" and "Nd"
		"A": "Lu", // uppercase letter: matches both "L" and "Lu"
		"a": "Ll", // lowercase letter: matches both "L" and "Ll"
		" ": "Zs", // space separator: matches both "Z" and "Zs"
	}

	for input, want := range cases {
		var last interface{}
		for i := 0; i < 50; i++ {
			got := getCategoryFunc(input)
			if last != nil && got != last {
				t.Fatalf("ugetCategory(%q) is nondeterministic: got %v then %v across calls", input, last, got)
			}
			last = got
		}
		if last != want {
			t.Errorf("ugetCategory(%q) = %v, want most-specific category %v", input, last, want)
		}
	}
}

// TestUnicodeTitle_PerWordCasing is a regression test for a bug where
// unicodeTitle() used strings.ToTitle, which performs a per-character
// Unicode title-case mapping equivalent to uppercasing the entire string for
// ordinary text (e.g. "hello world" -> "HELLO WORLD"), rather than the
// per-word "Hello World" capitalization implied by a function named
// utitle()/"title case". The fix uses golang.org/x/text/cases.Title, which
// is already an indirect dependency of this file via collate/language.
func TestUnicodeTitle_PerWordCasing(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterUnicode(env)
	unicodeModuleObj, _ := env.Get("unicode")
	unicodeModule := unicodeModuleObj.(map[string]interface{})
	utitleFunc := unicodeModule["utitle"].(r2core.BuiltinFunction)

	tests := []struct {
		input    string
		expected string
	}{
		{"hello world", "Hello World"},
		{"josé maría", "José María"},
		{"already Title Case", "Already Title Case"},
	}

	for _, tt := range tests {
		got := utitleFunc(tt.input)
		if got != tt.expected {
			t.Errorf("utitle(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}
