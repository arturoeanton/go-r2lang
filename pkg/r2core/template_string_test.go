package r2core

import (
	"math"
	"testing"
)

func TestTemplateString_BasicInterpolation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"simple variable interpolation",
			"let name = \"Juan\"; `Hello ${name}`",
			"Hello Juan",
		},
		{
			"number interpolation",
			"let age = 30; `Age: ${age}`",
			"Age: 30",
		},
		{
			"boolean interpolation",
			"let active = true; `Status: ${active}`",
			"Status: true",
		},
		{
			"expression interpolation",
			"let a = 5; let b = 3; `Result: ${a + b}`",
			"Result: 8",
		},
		{
			"multiple interpolations",
			"let name = \"Ana\"; let age = 25; `Name: ${name}, Age: ${age}`",
			"Name: Ana, Age: 25",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			program := parser.ParseProgram()

			env := NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			result := program.Eval(env)

			resultStr, ok := result.(string)
			if !ok {
				t.Fatalf("expected string result, got %T: %v", result, result)
			}

			if resultStr != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, resultStr)
			}
		})
	}
}

func TestTemplateString_MultilineSupport(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"simple multiline",
			"`Line 1\nLine 2\nLine 3`",
			"Line 1\nLine 2\nLine 3",
		},
		{
			"multiline with interpolation",
			"let title = \"Report\"; `${title}\n=======\nContent here`",
			"Report\n=======\nContent here",
		},
		{
			"multiline HTML template",
			"let title = \"Test Page\"; `<html>\n<head><title>${title}</title></head>\n<body>Hello World</body>\n</html>`",
			"<html>\n<head><title>Test Page</title></head>\n<body>Hello World</body>\n</html>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			program := parser.ParseProgram()

			env := NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			result := program.Eval(env)

			resultStr, ok := result.(string)
			if !ok {
				t.Fatalf("expected string result, got %T: %v", result, result)
			}

			if resultStr != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, resultStr)
			}
		})
	}
}

func TestTemplateString_EscapeSequences(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"escaped backtick",
			"`This is a backtick: \\` `",
			"This is a backtick: ` ",
		},
		{
			"escaped dollar",
			"`Price: \\$100`",
			"Price: $100",
		},
		{
			"escaped interpolation",
			"`Not interpolated: \\${variable}`",
			"Not interpolated: ${variable}",
		},
		{
			"escaped newline and tab",
			"`Line1\\nLine2\\tTabbed`",
			"Line1\nLine2\tTabbed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			program := parser.ParseProgram()

			env := NewEnvironment()
			result := program.Eval(env)

			resultStr, ok := result.(string)
			if !ok {
				t.Fatalf("expected string result, got %T: %v", result, result)
			}

			if resultStr != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, resultStr)
			}
		})
	}
}

func TestTemplateString_ComplexExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"function call interpolation",
			"func getName() { return \"Function Result\"; } `Hello ${getName()}`",
			"Hello Function Result",
		},
		{
			"object property access",
			"let user = {name: \"Pedro\", age: 35}; `User: ${user.name}, Age: ${user.age}`",
			"User: Pedro, Age: 35",
		},
		{
			"array access",
			"let items = [\"apple\", \"banana\", \"cherry\"]; `First: ${items[0]}, Second: ${items[1]}`",
			"First: apple, Second: banana",
		},
		{
			"conditional expression",
			"let score = 85; `Grade: ${score >= 80}`",
			"Grade: true",
		},
		{
			"nested template strings",
			"let inner = `inner template`; `Outer: ${inner}`",
			"Outer: inner template",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			program := parser.ParseProgram()

			env := NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			result := program.Eval(env)

			resultStr, ok := result.(string)
			if !ok {
				t.Fatalf("expected string result, got %T: %v", result, result)
			}

			if resultStr != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, resultStr)
			}
		})
	}
}

func TestTemplateString_Performance(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"empty template",
			"``",
			"",
		},
		{
			"plain text only",
			"`Just plain text without interpolation`",
			"Just plain text without interpolation",
		},
		{
			"single interpolation",
			"let value = \"test\"; `${value}`",
			"test",
		},
		{
			"large template with multiple interpolations",
			"let a = 1; let b = 2; let c = 3; `${a} ${b} ${c} ${a+b} ${b+c} ${a+c} ${a+b+c}`",
			"1 2 3 3 5 4 6",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			program := parser.ParseProgram()

			env := NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)

			// Run multiple times to test performance
			for i := 0; i < 100; i++ {
				result := program.Eval(env)

				resultStr, ok := result.(string)
				if !ok {
					t.Fatalf("expected string result, got %T: %v", result, result)
				}

				if resultStr != tt.expected {
					t.Errorf("expected %q, got %q", tt.expected, resultStr)
				}
			}
		})
	}
}

func TestTemplateString_ErrorCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			"unclosed template string",
			"`unclosed template",
		},
		{
			"unclosed interpolation",
			"`Hello ${name`",
		},
		{
			"invalid expression in interpolation",
			"`Hello ${invalid syntax here}`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("expected panic for input: %s", tt.input)
				}
			}()

			parser := NewParser(tt.input)
			program := parser.ParseProgram()

			env := NewEnvironment()
			program.Eval(env)
		})
	}
}

// TestTemplateString_FormatSpecifierEdgeCases exercises formatValue's
// printf-style, currency, percentage and comma-grouping code paths with
// inputs that previously produced wrong or garbled output (bugs found
// during adversarial review of template_string.go):
//   - printf verbs that require an integer/string type (%d, %x, %X, %s)
//     were fed the raw float64 R2Lang number, which Go's fmt package
//     doesn't coerce, so it silently emitted noise like "%!d(float64=5)"
//     or formatted the float's bit pattern for %x instead of an integer.
//   - malformed/negative precision (e.g. ":$.-2f") produced garbled output
//     like "$%!-(float64=3)2f" instead of falling back sanely.
//   - the "," (comma grouping) specifier used %g internally, which
//     switches to scientific notation for large non-integer numbers and
//     for Inf, corrupting the comma-grouped result (e.g. "+,Inf" or
//     "-1.234567891e+06" instead of "-1,234,567.891").
func TestTemplateString_FormatSpecifierEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"printf decimal verb on float64", "let x = 255.0; `${x:d}`", "255"},
		{"printf lowercase hex verb on float64", "let x = 255.0; `${x:x}`", "ff"},
		{"printf uppercase hex verb on float64", "let x = 255.0; `${x:X}`", "FF"},
		{"printf string verb on float64", "let x = 255.0; `${x:s}`", "255"},
		{"printf %g verb still works with raw float", "let x = 255.0; `${x:g}`", "255"},
		{"negative precision falls back to default precision", "let x = 3.14159; `${x:$.-2f}`", "$3.14"},
		{"comma grouping on large non-integer number", "let x = -1234567.891; `${x:,}`", "-1,234,567.891"},
		{"comma grouping on ordinary large integer", "let x = 1234567.0; `${x:,}`", "1,234,567"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			program := parser.ParseProgram()
			env := NewEnvironment()

			result := program.Eval(env)
			resultStr, ok := result.(string)
			if !ok {
				t.Fatalf("expected string result, got %T: %v", result, result)
			}
			if resultStr != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, resultStr)
			}
		})
	}
}

// TestFormatNumberWithCommas_SpecialValues locks down addCommas/
// formatNumberWithCommas behavior for +Inf/-Inf/NaN: these must be
// returned unmodified rather than having commas spliced into their
// letters (previously "+Inf" became "+,Inf").
func TestFormatNumberWithCommas_SpecialValues(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected string
	}{
		{"positive infinity", math.Inf(1), "+Inf"},
		{"negative infinity", math.Inf(-1), "-Inf"},
		{"NaN", math.NaN(), "NaN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatNumberWithCommas(tt.value, ",")
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestLexer_TemplateString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			"simple template",
			"`hello world`",
			[]Token{
				{Type: TOKEN_TEMPLATE_STRING, Value: "TEXT:hello world"},
			},
		},
		{
			"template with interpolation",
			"`hello ${name}`",
			[]Token{
				{Type: TOKEN_TEMPLATE_STRING, Value: "TEXT:hello \x00EXPR:name"},
			},
		},
		{
			"template with multiple parts",
			"`${a} + ${b} = ${a + b}`",
			[]Token{
				{Type: TOKEN_TEMPLATE_STRING, Value: "EXPR:a\x00TEXT: + \x00EXPR:b\x00TEXT: = \x00EXPR:a + b"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)

			for i, expected := range tt.expected {
				token := lexer.NextToken()

				if token.Type != expected.Type {
					t.Errorf("token %d: expected type %q, got %q", i, expected.Type, token.Type)
				}

				if token.Value != expected.Value {
					t.Errorf("token %d: expected value %q, got %q", i, expected.Value, token.Value)
				}
			}

			// Should end with EOF
			eofToken := lexer.NextToken()
			if eofToken.Type != TOKEN_EOF {
				t.Errorf("expected EOF token, got %q", eofToken.Type)
			}
		})
	}
}
