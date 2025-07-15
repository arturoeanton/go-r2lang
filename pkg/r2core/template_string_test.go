package r2core

import (
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