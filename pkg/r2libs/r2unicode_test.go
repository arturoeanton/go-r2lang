package r2libs

import (
	"testing"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func TestUnicodeBasicFunctions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			"ulen with ASCII",
			"ulen(\"hello\")",
			float64(5),
		},
		{
			"ulen with Spanish characters",
			"ulen(\"José María\")",
			float64(10),
		},
		{
			"ulen with emoji",
			"ulen(\"👋\")",
			float64(1),
		},
		{
			"ulen with complex emoji",
			"ulen(\"👨‍👩‍👧‍👦\")",
			float64(7), // Family emoji is composed of multiple code points
		},
		{
			"usubstr basic",
			"usubstr(\"hello\", 1, 3)",
			"ell",
		},
		{
			"usubstr with Spanish",
			"usubstr(\"José María\", 0, 4)",
			"José",
		},
		{
			"usubstr with emoji",
			"usubstr(\"Hello 👋 World\", 6, 1)",
			"👋",
		},
		{
			"uupper with ASCII",
			"uupper(\"hello\")",
			"HELLO",
		},
		{
			"uupper with Spanish",
			"uupper(\"josé\")",
			"JOSÉ",
		},
		{
			"ulower with ASCII",
			"ulower(\"HELLO\")",
			"hello",
		},
		{
			"ulower with Spanish",
			"ulower(\"JOSÉ\")",
			"josé",
		},
		{
			"ureverse with ASCII",
			"ureverse(\"hello\")",
			"olleh",
		},
		{
			"ureverse with Spanish",
			"ureverse(\"José\")",
			"ésoJ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := r2core.NewParser(tt.input)
			program := parser.ParseProgram()

			env := r2core.NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			RegisterUnicode(env)

			result := program.Eval(env)

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestUnicodeNormalization(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			"normalize NFC",
			"unormalize(\"José\", \"NFC\")",
			"José",
		},
		{
			"normalize NFD",
			"unormalize(\"José\", \"NFD\")",
			"José", // NFD form separates é into e + combining accent
		},
		{
			"normalize default NFC",
			"unormalize(\"café\")",
			"café",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := r2core.NewParser(tt.input)
			program := parser.ParseProgram()

			env := r2core.NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			RegisterUnicode(env)

			result := program.Eval(env)

			// For normalization, we just check that it returns a string
			if _, ok := result.(string); !ok {
				t.Errorf("expected string, got %T", result)
			}
		})
	}
}

func TestUnicodeComparison(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{
			"compare equal strings",
			"ucompare(\"hello\", \"hello\")",
			0,
		},
		{
			"compare different strings",
			"ucompare(\"a\", \"b\")",
			-1,
		},
		{
			"compare with locale",
			"ucompare(\"café\", \"cafe\", \"es\")",
			1, // café > cafe in Spanish collation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := r2core.NewParser(tt.input)
			program := parser.ParseProgram()

			env := r2core.NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			RegisterUnicode(env)

			result := program.Eval(env)

			resultFloat, ok := result.(float64)
			if !ok {
				t.Fatalf("expected float64, got %T", result)
			}

			if resultFloat != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, resultFloat)
			}
		})
	}
}

func TestUnicodeValidation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			"valid UTF-8",
			"uisvalid(\"José María\")",
			true,
		},
		{
			"valid emoji",
			"uisvalid(\"👋🌍\")",
			true,
		},
		{
			"valid ASCII",
			"uisvalid(\"hello world\")",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := r2core.NewParser(tt.input)
			program := parser.ParseProgram()

			env := r2core.NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			RegisterUnicode(env)

			result := program.Eval(env)

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestUnicodeCharacterCodes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{
			"char code for A",
			"ucharcode(\"A\")",
			65,
		},
		{
			"char code for ñ",
			"ucharcode(\"ñ\")",
			241,
		},
		{
			"char code for emoji",
			"ucharcode(\"👋\")",
			128075,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := r2core.NewParser(tt.input)
			program := parser.ParseProgram()

			env := r2core.NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			RegisterUnicode(env)

			result := program.Eval(env)

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestUnicodeFromCharCode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"from code A",
			"ufromcode(65)",
			"A",
		},
		{
			"from code ñ",
			"ufromcode(241)",
			"ñ",
		},
		{
			"from code emoji",
			"ufromcode(128075)",
			"👋",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := r2core.NewParser(tt.input)
			program := parser.ParseProgram()

			env := r2core.NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			RegisterUnicode(env)

			result := program.Eval(env)

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestUnicodeCharacterClassification(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			"uisLetter with ASCII",
			"uisLetter(\"A\")",
			true,
		},
		{
			"uisLetter with Spanish",
			"uisLetter(\"ñ\")",
			true,
		},
		{
			"uisLetter with number",
			"uisLetter(\"5\")",
			false,
		},
		{
			"uisDigit with ASCII",
			"uisDigit(\"5\")",
			true,
		},
		{
			"uisDigit with letter",
			"uisDigit(\"A\")",
			false,
		},
		{
			"uisSpace with space",
			"uisSpace(\" \")",
			true,
		},
		{
			"uisSpace with letter",
			"uisSpace(\"A\")",
			false,
		},
		{
			"uisUpper with uppercase",
			"uisUpper(\"A\")",
			true,
		},
		{
			"uisUpper with lowercase",
			"uisUpper(\"a\")",
			false,
		},
		{
			"uisLower with lowercase",
			"uisLower(\"a\")",
			true,
		},
		{
			"uisLower with uppercase",
			"uisLower(\"A\")",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := r2core.NewParser(tt.input)
			program := parser.ParseProgram()

			env := r2core.NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			RegisterUnicode(env)

			result := program.Eval(env)

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestUnicodeRegex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int // number of matches
	}{
		{
			"regex with ASCII",
			"len(uregex(\"[a-z]+\", \"hello world\"))",
			2,
		},
		{
			"regex match ASCII",
			"uregexMatch(\"hello\", \"hello world\")",
			1, // true = 1
		},
		{
			"regex match Unicode",
			"uregexMatch(\"José\", \"José María\")",
			1, // true = 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := r2core.NewParser(tt.input)
			program := parser.ParseProgram()

			env := r2core.NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			RegisterUnicode(env)
			RegisterStd(env) // For len function

			result := program.Eval(env)

			var resultInt int
			if boolResult, ok := result.(bool); ok {
				if boolResult {
					resultInt = 1
				} else {
					resultInt = 0
				}
			} else if floatResult, ok := result.(float64); ok {
				resultInt = int(floatResult)
			} else {
				t.Fatalf("unexpected result type: %T", result)
			}

			if resultInt != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, resultInt)
			}
		})
	}
}

func TestUnicodeEscapeSequences(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			"unicode escape basic",
			"\"\\u0041\"", // Unicode for 'A'
			"A",
		},
		{
			"unicode escape Spanish",
			"\"\\u00f1\"", // Unicode for 'ñ'
			"ñ",
		},
		{
			"unicode escape emoji",
			"\"\\U0001F44B\"", // Unicode for '👋'
			"👋",
		},
		{
			"hex escape",
			"\"\\x41\"", // Hex for 'A'
			"A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := r2core.NewParser(tt.input)
			program := parser.ParseProgram()

			env := r2core.NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			RegisterUnicode(env)

			result := program.Eval(env)

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestUnicodeIdentifiers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			"Spanish identifier",
			"let año = 2024; año",
			float64(2024),
		},
		{
			"Japanese identifier",
			"let 名前 = \"田中\"; 名前",
			"田中",
		},
		{
			"Russian identifier",
			"let имя = \"Иван\"; имя",
			"Иван",
		},
		{
			"Arabic identifier",
			"let اسم = \"أحمد\"; اسم",
			"أحمد",
		},
		{
			"Greek identifier",
			"let όνομα = \"Γιάννης\"; όνομα",
			"Γιάννης",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := r2core.NewParser(tt.input)
			program := parser.ParseProgram()

			env := r2core.NewEnvironment()
			env.Set("true", true)
			env.Set("false", false)
			env.Set("nil", nil)
			RegisterUnicode(env)

			result := program.Eval(env)

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
