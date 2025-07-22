package r2core

import (
	"testing"
)

func TestDSLTokenPriority(t *testing.T) {
	// Test that keyword tokens have higher priority than generic patterns
	grammar := NewDSLGrammar()

	// Add generic pattern first (lower priority)
	err := grammar.AddToken("IDENTIFIER", "[a-zA-Z]+")
	if err != nil {
		t.Fatalf("Failed to add IDENTIFIER token: %v", err)
	}

	// Add specific keyword (should have higher priority)
	err = grammar.AddKeywordToken("SELECT", "select")
	if err != nil {
		t.Fatalf("Failed to add SELECT keyword: %v", err)
	}

	// Test tokenization
	parser := NewDSLParser(grammar)
	err = parser.Tokenize("select")
	if err != nil {
		t.Fatalf("Tokenization failed: %v", err)
	}

	// Should match SELECT (keyword) not IDENTIFIER (generic)
	if len(parser.Tokens) != 1 {
		t.Fatalf("Expected 1 token, got %d", len(parser.Tokens))
	}

	token := parser.Tokens[0]
	if token.Type != "SELECT" {
		t.Errorf("Expected token type 'SELECT', got '%s'", token.Type)
	}

	if token.Value != "select" {
		t.Errorf("Expected token value 'select', got '%s'", token.Value)
	}
}

func TestDSLTokenPriorityMultiple(t *testing.T) {
	// Test multiple keywords with generic pattern
	grammar := NewDSLGrammar()

	// Add generic pattern
	grammar.AddToken("IDENTIFIER", "[a-zA-Z]+")

	// Add specific keywords
	grammar.AddKeywordToken("SELECT", "select")
	grammar.AddKeywordToken("FROM", "from")
	grammar.AddKeywordToken("WHERE", "where")

	// Test complex query
	parser := NewDSLParser(grammar)
	err := parser.Tokenize("select name from users where age")
	if err != nil {
		t.Fatalf("Tokenization failed: %v", err)
	}

	expected := []struct {
		Type  string
		Value string
	}{
		{"SELECT", "select"},
		{"IDENTIFIER", "name"},
		{"FROM", "from"},
		{"IDENTIFIER", "users"},
		{"WHERE", "where"},
		{"IDENTIFIER", "age"},
	}

	if len(parser.Tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(parser.Tokens))
	}

	for i, exp := range expected {
		if parser.Tokens[i].Type != exp.Type {
			t.Errorf("Token %d: expected type '%s', got '%s'", i, exp.Type, parser.Tokens[i].Type)
		}
		if parser.Tokens[i].Value != exp.Value {
			t.Errorf("Token %d: expected value '%s', got '%s'", i, exp.Value, parser.Tokens[i].Value)
		}
	}
}

func TestDSLKeywordDetection(t *testing.T) {
	// Test the keyword detection function
	dsl := &DSLDefinition{}

	testCases := []struct {
		pattern   string
		isKeyword bool
		desc      string
	}{
		{"select", true, "simple keyword"},
		{"from", true, "simple keyword"},
		{"where", true, "simple keyword"},
		{"[a-zA-Z]+", false, "regex pattern"},
		{"[0-9]+", false, "number pattern"},
		{"\\w+", false, "word pattern"},
		{"hello123", true, "alphanumeric keyword"},
		{"ORDER_BY", true, "keyword with underscore"},
		{">=", true, "operator keyword"},
		{"==", true, "equality operator"},
		{".", false, "regex metacharacter"},
		{"*", false, "regex quantifier"},
		{"(test)", false, "grouped pattern"},
	}

	for _, tc := range testCases {
		result := dsl.isKeywordToken(tc.pattern)
		if result != tc.isKeyword {
			t.Errorf("Pattern '%s' (%s): expected %v, got %v", tc.pattern, tc.desc, tc.isKeyword, result)
		}
	}
}

func TestDSLTokenPriorityLongestMatch(t *testing.T) {
	// Test that among same priority tokens, longest match wins
	grammar := NewDSLGrammar()

	// Add tokens with same priority but different lengths
	grammar.AddKeywordToken("ORDER", "order")
	grammar.AddKeywordToken("ORDER_BY", "order by")

	parser := NewDSLParser(grammar)
	err := parser.Tokenize("order by")
	if err != nil {
		t.Fatalf("Tokenization failed: %v", err)
	}

	// Should match ORDER_BY (longer) not ORDER
	if len(parser.Tokens) != 1 {
		t.Fatalf("Expected 1 token, got %d", len(parser.Tokens))
	}

	token := parser.Tokens[0]
	if token.Type != "ORDER_BY" {
		t.Errorf("Expected token type 'ORDER_BY', got '%s'", token.Type)
	}

	if token.Value != "order by" {
		t.Errorf("Expected token value 'order by', got '%s'", token.Value)
	}
}

func TestDSLTokenCaseInsensitive(t *testing.T) {
	// Test that keyword tokens are case insensitive
	grammar := NewDSLGrammar()

	grammar.AddKeywordToken("SELECT", "select")

	testCases := []string{"select", "SELECT", "Select", "sElEcT"}

	for _, testCase := range testCases {
		parser := NewDSLParser(grammar)
		err := parser.Tokenize(testCase)
		if err != nil {
			t.Fatalf("Tokenization failed for '%s': %v", testCase, err)
		}

		if len(parser.Tokens) != 1 {
			t.Fatalf("Expected 1 token for '%s', got %d", testCase, len(parser.Tokens))
		}

		token := parser.Tokens[0]
		if token.Type != "SELECT" {
			t.Errorf("Case '%s': expected token type 'SELECT', got '%s'", testCase, token.Type)
		}
	}
}

func TestDSLTokenWithNumbers(t *testing.T) {
	// Test keywords with numbers vs generic patterns
	grammar := NewDSLGrammar()

	grammar.AddToken("NUMBER", "[0-9]+")
	grammar.AddKeywordToken("TOP", "top")
	grammar.AddToken("IDENTIFIER", "[a-zA-Z]+")

	parser := NewDSLParser(grammar)
	err := parser.Tokenize("top 10 users")
	if err != nil {
		t.Fatalf("Tokenization failed: %v", err)
	}

	expected := []struct {
		Type  string
		Value string
	}{
		{"TOP", "top"},
		{"NUMBER", "10"},
		{"IDENTIFIER", "users"},
	}

	if len(parser.Tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(parser.Tokens))
	}

	for i, exp := range expected {
		if parser.Tokens[i].Type != exp.Type {
			t.Errorf("Token %d: expected type '%s', got '%s'", i, exp.Type, parser.Tokens[i].Type)
		}
		if parser.Tokens[i].Value != exp.Value {
			t.Errorf("Token %d: expected value '%s', got '%s'", i, exp.Value, parser.Tokens[i].Value)
		}
	}
}

func TestDSLComplexTokenPriority(t *testing.T) {
	// Test complex scenario with mixed token types
	grammar := NewDSLGrammar()

	// Generic patterns (low priority)
	grammar.AddToken("IDENTIFIER", "[a-zA-Z_][a-zA-Z0-9_]*")
	grammar.AddToken("NUMBER", "[0-9]+")
	grammar.AddToken("OPERATOR", "[><=]+")

	// Specific keywords (high priority)
	grammar.AddKeywordToken("SELECT", "select")
	grammar.AddKeywordToken("FROM", "from")
	grammar.AddKeywordToken("WHERE", "where")
	grammar.AddKeywordToken("GROUP", "group")
	grammar.AddKeywordToken("BY", "by")

	// Test complex query
	parser := NewDSLParser(grammar)
	err := parser.Tokenize("select user_id from users where age > 25 group by department")
	if err != nil {
		t.Fatalf("Tokenization failed: %v", err)
	}

	expected := []struct {
		Type  string
		Value string
	}{
		{"SELECT", "select"},
		{"IDENTIFIER", "user_id"},
		{"FROM", "from"},
		{"IDENTIFIER", "users"},
		{"WHERE", "where"},
		{"IDENTIFIER", "age"},
		{"OPERATOR", ">"},
		{"NUMBER", "25"},
		{"GROUP", "group"},
		{"BY", "by"},
		{"IDENTIFIER", "department"},
	}

	if len(parser.Tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(parser.Tokens))
	}

	for i, exp := range expected {
		if parser.Tokens[i].Type != exp.Type {
			t.Errorf("Token %d: expected type '%s', got '%s'", i, exp.Type, parser.Tokens[i].Type)
		}
		if parser.Tokens[i].Value != exp.Value {
			t.Errorf("Token %d: expected value '%s', got '%s'", i, exp.Value, parser.Tokens[i].Value)
		}
	}
}
