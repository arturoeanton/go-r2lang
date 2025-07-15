package r2core

import (
	"testing"
)

func TestLexer_NewLexer(t *testing.T) {
	input := "let x = 42"
	lexer := NewLexer(input)

	if lexer.input != input {
		t.Errorf("Expected input %q, got %q", input, lexer.input)
	}
	if lexer.pos != 0 {
		t.Errorf("Expected pos 0, got %d", lexer.pos)
	}
	if lexer.line != 1 {
		t.Errorf("Expected line 1, got %d", lexer.line)
	}
	if lexer.length != len(input) {
		t.Errorf("Expected length %d, got %d", len(input), lexer.length)
	}
}

func TestLexer_Numbers(t *testing.T) {
	tests := []struct {
		input    string
		expected []Token
	}{
		{
			input: "42",
			expected: []Token{
				{Type: TOKEN_NUMBER, Value: "42"},
			},
		},
		{
			input: "3.14",
			expected: []Token{
				{Type: TOKEN_NUMBER, Value: "3.14"},
			},
		},
		{
			input: "= -42",
			expected: []Token{
				{Type: TOKEN_SYMBOL, Value: "="},
				{Type: TOKEN_NUMBER, Value: "-42"},
			},
		},
		{
			input: "= +3.14",
			expected: []Token{
				{Type: TOKEN_SYMBOL, Value: "="},
				{Type: TOKEN_NUMBER, Value: "+3.14"},
			},
		},
		{
			input: "(-42)",
			expected: []Token{
				{Type: TOKEN_SYMBOL, Value: "("},
				{Type: TOKEN_NUMBER, Value: "-42"},
				{Type: TOKEN_SYMBOL, Value: ")"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			lexer := NewLexer(test.input)
			for i, expectedToken := range test.expected {
				token := lexer.NextToken()
				if token.Type != expectedToken.Type {
					t.Errorf("Token %d: expected type %q, got %q", i, expectedToken.Type, token.Type)
				}
				if token.Value != expectedToken.Value {
					t.Errorf("Token %d: expected value %q, got %q", i, expectedToken.Value, token.Value)
				}
			}
		})
	}
}

func TestLexer_Strings(t *testing.T) {
	tests := []struct {
		input    string
		expected []Token
	}{
		{
			input: `"hello"`,
			expected: []Token{
				{Type: TOKEN_STRING, Value: "hello"},
			},
		},
		{
			input: `'world'`,
			expected: []Token{
				{Type: TOKEN_STRING, Value: "world"},
			},
		},
		{
			input: `"hello world"`,
			expected: []Token{
				{Type: TOKEN_STRING, Value: "hello world"},
			},
		},
		{
			input: `""`,
			expected: []Token{
				{Type: TOKEN_STRING, Value: ""},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			lexer := NewLexer(test.input)
			for i, expectedToken := range test.expected {
				token := lexer.NextToken()
				if token.Type != expectedToken.Type {
					t.Errorf("Token %d: expected type %q, got %q", i, expectedToken.Type, token.Type)
				}
				if token.Value != expectedToken.Value {
					t.Errorf("Token %d: expected value %q, got %q", i, expectedToken.Value, token.Value)
				}
			}
		})
	}
}

func TestLexer_Identifiers(t *testing.T) {
	tests := []struct {
		input    string
		expected []Token
	}{
		{
			input: "variable",
			expected: []Token{
				{Type: TOKEN_IDENT, Value: "variable"},
			},
		},
		{
			input: "_private",
			expected: []Token{
				{Type: TOKEN_IDENT, Value: "_private"},
			},
		},
		{
			input: "$global",
			expected: []Token{
				{Type: TOKEN_IDENT, Value: "$global"},
			},
		},
		{
			input: "var123",
			expected: []Token{
				{Type: TOKEN_IDENT, Value: "var123"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			lexer := NewLexer(test.input)
			for i, expectedToken := range test.expected {
				token := lexer.NextToken()
				if token.Type != expectedToken.Type {
					t.Errorf("Token %d: expected type %q, got %q", i, expectedToken.Type, token.Type)
				}
				if token.Value != expectedToken.Value {
					t.Errorf("Token %d: expected value %q, got %q", i, expectedToken.Value, token.Value)
				}
			}
		})
	}
}

func TestLexer_Keywords(t *testing.T) {
	tests := []struct {
		input         string
		expectedType  string
		expectedValue string
	}{
		{"import", TOKEN_IMPORT, "import"},
		{"as", TOKEN_AS, "as"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			lexer := NewLexer(test.input)
			token := lexer.NextToken()
			if token.Type != test.expectedType {
				t.Errorf("Expected type %q, got %q", test.expectedType, token.Type)
			}
			if token.Value != test.expectedValue {
				t.Errorf("Expected value %q, got %q", test.expectedValue, token.Value)
			}
		})
	}
}

func TestLexer_Symbols(t *testing.T) {
	tests := []struct {
		input    string
		expected []Token
	}{
		{
			input: "x++",
			expected: []Token{
				{Type: TOKEN_IDENT, Value: "x"},
				{Type: TOKEN_SYMBOL, Value: "++"},
			},
		},
		{
			input: "x--",
			expected: []Token{
				{Type: TOKEN_IDENT, Value: "x"},
				{Type: TOKEN_SYMBOL, Value: "--"},
			},
		},
		{
			input: "=>",
			expected: []Token{
				{Type: TOKEN_ARROW, Value: "=>"},
			},
		},
		{
			input: "==",
			expected: []Token{
				{Type: TOKEN_SYMBOL, Value: "=="},
			},
		},
		{
			input: "!=",
			expected: []Token{
				{Type: TOKEN_SYMBOL, Value: "!="},
			},
		},
		{
			input: "<=",
			expected: []Token{
				{Type: TOKEN_SYMBOL, Value: "<="},
			},
		},
		{
			input: ">=",
			expected: []Token{
				{Type: TOKEN_SYMBOL, Value: ">="},
			},
		},
		{
			input: "(){}[]",
			expected: []Token{
				{Type: TOKEN_SYMBOL, Value: "("},
				{Type: TOKEN_SYMBOL, Value: ")"},
				{Type: TOKEN_SYMBOL, Value: "{"},
				{Type: TOKEN_SYMBOL, Value: "}"},
				{Type: TOKEN_SYMBOL, Value: "["},
				{Type: TOKEN_SYMBOL, Value: "]"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			lexer := NewLexer(test.input)
			for i, expectedToken := range test.expected {
				token := lexer.NextToken()
				if token.Type != expectedToken.Type {
					t.Errorf("Token %d: expected type %q, got %q", i, expectedToken.Type, token.Type)
				}
				if token.Value != expectedToken.Value {
					t.Errorf("Token %d: expected value %q, got %q", i, expectedToken.Value, token.Value)
				}
			}
		})
	}
}

func TestLexer_Comments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "single line comment",
			input: "let x = 42 // this is a comment\nlet y = 24",
			expected: []Token{
				{Type: TOKEN_IDENT, Value: "let"},
				{Type: TOKEN_IDENT, Value: "x"},
				{Type: TOKEN_SYMBOL, Value: "="},
				{Type: TOKEN_NUMBER, Value: "42"},
				{Type: TOKEN_SYMBOL, Value: "\n"},
				{Type: TOKEN_IDENT, Value: "let"},
				{Type: TOKEN_IDENT, Value: "y"},
				{Type: TOKEN_SYMBOL, Value: "="},
				{Type: TOKEN_NUMBER, Value: "24"},
			},
		},
		{
			name:  "block comment",
			input: "let x = 42 /* this is a \n block comment */ let y = 24",
			expected: []Token{
				{Type: TOKEN_IDENT, Value: "let"},
				{Type: TOKEN_IDENT, Value: "x"},
				{Type: TOKEN_SYMBOL, Value: "="},
				{Type: TOKEN_NUMBER, Value: "42"},
				{Type: TOKEN_IDENT, Value: "let"},
				{Type: TOKEN_IDENT, Value: "y"},
				{Type: TOKEN_SYMBOL, Value: "="},
				{Type: TOKEN_NUMBER, Value: "24"},
			},
		},
		{
			name:  "division vs comment",
			input: "x / y",
			expected: []Token{
				{Type: TOKEN_IDENT, Value: "x"},
				{Type: TOKEN_SYMBOL, Value: "/"},
				{Type: TOKEN_IDENT, Value: "y"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lexer := NewLexer(test.input)
			for i, expectedToken := range test.expected {
				token := lexer.NextToken()
				if token.Type != expectedToken.Type {
					t.Errorf("Token %d: expected type %q, got %q", i, expectedToken.Type, token.Type)
				}
				if token.Value != expectedToken.Value {
					t.Errorf("Token %d: expected value %q, got %q", i, expectedToken.Value, token.Value)
				}
			}
		})
	}
}

func TestLexer_LineAndColumn(t *testing.T) {
	input := "let x = 42\nlet y = 24"
	lexer := NewLexer(input)

	// let (line 1)
	token := lexer.NextToken()
	if token.Line != 1 {
		t.Errorf("Expected line 1, got %d", token.Line)
	}

	// x (line 1)
	token = lexer.NextToken()
	if token.Line != 1 {
		t.Errorf("Expected line 1, got %d", token.Line)
	}

	// = (line 1)
	token = lexer.NextToken()
	if token.Line != 1 {
		t.Errorf("Expected line 1, got %d", token.Line)
	}

	// 42 (line 1)
	token = lexer.NextToken()
	if token.Line != 1 {
		t.Errorf("Expected line 1, got %d", token.Line)
	}

	// \n (line 1)
	token = lexer.NextToken()
	if token.Line != 1 {
		t.Errorf("Expected line 1, got %d", token.Line)
	}

	// let (line 2)
	token = lexer.NextToken()
	if token.Line != 2 {
		t.Errorf("Expected line 2, got %d", token.Line)
	}
}

func TestLexer_EOF(t *testing.T) {
	input := "x"
	lexer := NewLexer(input)

	// x
	token := lexer.NextToken()
	if token.Type != TOKEN_IDENT {
		t.Errorf("Expected IDENT, got %q", token.Type)
	}

	// EOF
	token = lexer.NextToken()
	if token.Type != TOKEN_EOF {
		t.Errorf("Expected EOF, got %q", token.Type)
	}

	// Multiple EOF calls should return EOF
	token = lexer.NextToken()
	if token.Type != TOKEN_EOF {
		t.Errorf("Expected EOF, got %q", token.Type)
	}
}

func TestLexer_WhitespaceHandling(t *testing.T) {
	input := "  \t  let \n\r\t x  =  42  "
	expected := []Token{
		{Type: TOKEN_IDENT, Value: "let"},
		{Type: TOKEN_SYMBOL, Value: "\n"},
		{Type: TOKEN_IDENT, Value: "x"},
		{Type: TOKEN_SYMBOL, Value: "="},
		{Type: TOKEN_NUMBER, Value: "42"},
	}

	lexer := NewLexer(input)
	for i, expectedToken := range expected {
		token := lexer.NextToken()
		if token.Type != expectedToken.Type {
			t.Errorf("Token %d: expected type %q, got %q", i, expectedToken.Type, token.Type)
		}
		if token.Value != expectedToken.Value {
			t.Errorf("Token %d: expected value %q, got %q", i, expectedToken.Value, token.Value)
		}
	}
}

func TestLexer_CompleteExpression(t *testing.T) {
	input := `let factorial = func(n) {
		if (n <= 1) {
			return 1
		}
		return n * factorial(n - 1)
	}`

	expected := []Token{
		{Type: TOKEN_IDENT, Value: "let"},
		{Type: TOKEN_IDENT, Value: "factorial"},
		{Type: TOKEN_SYMBOL, Value: "="},
		{Type: TOKEN_IDENT, Value: "func"},
		{Type: TOKEN_SYMBOL, Value: "("},
		{Type: TOKEN_IDENT, Value: "n"},
		{Type: TOKEN_SYMBOL, Value: ")"},
		{Type: TOKEN_SYMBOL, Value: "{"},
		{Type: TOKEN_SYMBOL, Value: "\n"},
		{Type: TOKEN_IDENT, Value: "if"},
		{Type: TOKEN_SYMBOL, Value: "("},
		{Type: TOKEN_IDENT, Value: "n"},
		{Type: TOKEN_SYMBOL, Value: "<="},
		{Type: TOKEN_NUMBER, Value: "1"},
		{Type: TOKEN_SYMBOL, Value: ")"},
		{Type: TOKEN_SYMBOL, Value: "{"},
		{Type: TOKEN_SYMBOL, Value: "\n"},
		{Type: TOKEN_IDENT, Value: "return"},
		{Type: TOKEN_NUMBER, Value: "1"},
		{Type: TOKEN_SYMBOL, Value: "\n"},
		{Type: TOKEN_SYMBOL, Value: "}"},
		{Type: TOKEN_SYMBOL, Value: "\n"},
		{Type: TOKEN_IDENT, Value: "return"},
		{Type: TOKEN_IDENT, Value: "n"},
		{Type: TOKEN_SYMBOL, Value: "*"},
		{Type: TOKEN_IDENT, Value: "factorial"},
		{Type: TOKEN_SYMBOL, Value: "("},
		{Type: TOKEN_IDENT, Value: "n"},
		{Type: TOKEN_SYMBOL, Value: "-"},
		{Type: TOKEN_NUMBER, Value: "1"},
		{Type: TOKEN_SYMBOL, Value: ")"},
		{Type: TOKEN_SYMBOL, Value: "\n"},
		{Type: TOKEN_SYMBOL, Value: "}"},
	}

	lexer := NewLexer(input)
	for i, expectedToken := range expected {
		token := lexer.NextToken()
		if token.Type != expectedToken.Type {
			t.Errorf("Token %d: expected type %q, got %q", i, expectedToken.Type, token.Type)
		}
		if token.Value != expectedToken.Value {
			t.Errorf("Token %d: expected value %q, got %q", i, expectedToken.Value, token.Value)
		}
	}
}

// Test helper functions
func TestHelperFunctions(t *testing.T) {
	// Test isWhitespace
	if !isWhitespace(' ') {
		t.Error("Expected space to be whitespace")
	}
	if !isWhitespace('\t') {
		t.Error("Expected tab to be whitespace")
	}
	if !isWhitespace('\r') {
		t.Error("Expected carriage return to be whitespace")
	}
	if isWhitespace('a') {
		t.Error("Expected 'a' to not be whitespace")
	}

	// Test isLetter
	if !isLetter('a') {
		t.Error("Expected 'a' to be letter")
	}
	if !isLetter('Z') {
		t.Error("Expected 'Z' to be letter")
	}
	if !isLetter('_') {
		t.Error("Expected '_' to be letter")
	}
	if !isLetter('$') {
		t.Error("Expected '$' to be letter")
	}
	if isLetter('1') {
		t.Error("Expected '1' to not be letter")
	}

	// Test isDigit
	if !isDigit('0') {
		t.Error("Expected '0' to be digit")
	}
	if !isDigit('9') {
		t.Error("Expected '9' to be digit")
	}
	if isDigit('a') {
		t.Error("Expected 'a' to not be digit")
	}
}

// Benchmark tests
func BenchmarkLexer_Numbers(b *testing.B) {
	input := "123 456.789 -42 +3.14"
	for i := 0; i < b.N; i++ {
		lexer := NewLexer(input)
		for {
			token := lexer.NextToken()
			if token.Type == TOKEN_EOF {
				break
			}
		}
	}
}

func BenchmarkLexer_Identifiers(b *testing.B) {
	input := "variable func let while for if else return"
	for i := 0; i < b.N; i++ {
		lexer := NewLexer(input)
		for {
			token := lexer.NextToken()
			if token.Type == TOKEN_EOF {
				break
			}
		}
	}
}

func BenchmarkLexer_ComplexCode(b *testing.B) {
	input := `
	func factorial(n) {
		if (n <= 1) {
			return 1
		}
		return n * factorial(n - 1)
	}
	let result = factorial(10)
	print("Result: " + result)
	`
	for i := 0; i < b.N; i++ {
		lexer := NewLexer(input)
		for {
			token := lexer.NextToken()
			if token.Type == TOKEN_EOF {
				break
			}
		}
	}
}
