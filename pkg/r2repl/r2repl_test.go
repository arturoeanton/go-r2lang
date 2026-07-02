package r2repl

import "testing"

func TestIsIncomplete_BraceParenBracketNesting(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty", "", false},
		{"complete statement", "let x = 1\n", false},
		{"open brace", "func add(a, b) {\n", true},
		{"closed brace", "func add(a, b) {\n  return a + b\n}\n", false},
		{"open bracket (array literal)", "let arr = [1,\n", true},
		{"closed bracket", "let arr = [1,\n2,\n3]\n", false},
		{"open paren", "let y = (1 +\n", true},
		{"closed paren", "let y = (1 +\n2)\n", false},
		{"open brace map literal", "let m = {\n  a: 1,\n", true},
		{"closed brace map literal", "let m = {\n  a: 1,\n  b: 2\n}\n", false},
		{"unmatched close does not hang as incomplete", "}\n", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isIncomplete(tt.input)
			if got != tt.want {
				t.Errorf("isIncomplete(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// Regression test: a naive substring brace-counter (the REPL's original
// implementation) would miscount braces that appear inside string literals
// or comments, since it doesn't understand lexical context. The real-lexer
// based implementation must not have this problem.
func TestIsIncomplete_IgnoresBracesInsideStringsAndComments(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"brace inside double-quoted string", `let s = "hello { world"` + "\n"},
		{"unbalanced brackets inside string", `let s = "a [ b ) c"` + "\n"},
		{"brace inside line comment", "let z = 1 // comment with { unmatched brace\n"},
		{"brace inside block comment", "let z = 1 /* { unmatched */\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if isIncomplete(tt.input) {
				t.Errorf("isIncomplete(%q) = true, want false (braces are inside a string/comment, not real nesting)", tt.input)
			}
		})
	}
}

// Regression test: an unterminated template string (backtick) should be
// treated as needing more input, the same way an open brace does, so a
// multi-line template string can be typed across several REPL lines.
func TestIsIncomplete_UnterminatedTemplateStringNeedsMoreInput(t *testing.T) {
	if !isIncomplete("let msg = `hello\n") {
		t.Error("isIncomplete of an unterminated template string should be true (needs closing backtick)")
	}
	if isIncomplete("let msg = `hello\nworld`\n") {
		t.Error("isIncomplete of a closed multi-line template string should be false")
	}
}

// Regression test for a bug where pkg/r2repl's createR2Environment only
// called r2libs.RegisterLib (the "r2"/"go" concurrency builtins), so the
// REPL was missing nearly the entire standard library that scripts run via
// main.go/pkg/r2lang.RunCode get (std, math, json, io, os, string, ...).
func TestCreateR2Environment_RegistersStandardLibraryModules(t *testing.T) {
	env := createR2Environment()

	modules := []string{"std", "math", "json", "io", "os", "string", "console"}
	for _, name := range modules {
		if _, ok := env.Get(name); !ok {
			t.Errorf("expected module %q to be registered in the REPL environment, but it was not found", name)
		}
	}

	builtins := []string{"r2", "go"}
	for _, name := range builtins {
		if _, ok := env.Get(name); !ok {
			t.Errorf("expected builtin %q to be registered in the REPL environment, but it was not found", name)
		}
	}
}
