package r2core

import (
	"reflect"
	"sort"
	"testing"
)

// TestArrayComprehension_ChainedGenerators guards against a bug where a
// later "for" clause's iterator expression couldn't see variables bound by
// an earlier "for" clause. generateElements evaluated each generator's
// Iterator against the raw outer env instead of an env carrying the
// bindings accumulated so far, so "for row in matrix for y in row" panicked
// with "Undeclared variable: row" instead of flattening the matrix.
func TestArrayComprehension_ChainedGenerators(t *testing.T) {
	input := `let matrix = [[1, 2], [3, 4]]; [y for row in matrix for y in row]`

	parser := NewParser(input)
	program := parser.ParseProgram()
	env := NewEnvironment()

	var result interface{}
	for _, stmt := range program.Statements {
		result = stmt.Eval(env)
	}

	arr, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected []interface{}, got %T: %v", result, result)
	}

	expected := []interface{}{float64(1), float64(2), float64(3), float64(4)}
	if !reflect.DeepEqual(arr, expected) {
		t.Errorf("expected %v, got %v", expected, arr)
	}
}

// TestArrayComprehension_ChainedGeneratorsWithFilter is the same scenario
// exercised through the P4 example file (flatten + filter), which is the
// realistic use case that surfaced the bug.
func TestArrayComprehension_ChainedGeneratorsWithFilter(t *testing.T) {
	input := `let matrix = [[1, 2], [3, 4]]; [y for row in matrix for y in row if y > 1]`

	parser := NewParser(input)
	program := parser.ParseProgram()
	env := NewEnvironment()

	var result interface{}
	for _, stmt := range program.Statements {
		result = stmt.Eval(env)
	}

	arr, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected []interface{}, got %T: %v", result, result)
	}

	expected := []interface{}{float64(2), float64(3), float64(4)}
	if !reflect.DeepEqual(arr, expected) {
		t.Errorf("expected %v, got %v", expected, arr)
	}
}

// TestArrayComprehension_MapSourceBindsKey guards against a bug where
// iterating a map[string]interface{} source in a comprehension bound the
// loop variable to the map's VALUE instead of its KEY, unlike the regular
// for-in loop (ForStatement.evalForIn) which binds to the key. This made
// "[k for k in someMap]" silently return values where the (very similarly
// spelled) variable name strongly implied keys were expected.
func TestArrayComprehension_MapSourceBindsKey(t *testing.T) {
	input := `let obj = {a: 1, b: 2}; [k for k in obj]`

	parser := NewParser(input)
	program := parser.ParseProgram()
	env := NewEnvironment()

	var result interface{}
	for _, stmt := range program.Statements {
		result = stmt.Eval(env)
	}

	arr, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected []interface{}, got %T: %v", result, result)
	}

	got := make([]string, 0, len(arr))
	for _, v := range arr {
		s, ok := v.(string)
		if !ok {
			t.Fatalf("expected string keys, got %T: %v", v, v)
		}
		got = append(got, s)
	}
	sort.Strings(got)

	expected := []string{"a", "b"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected keys %v, got %v", expected, got)
	}
}

// TestObjectComprehension_MapSourceBindsKey is the ObjectComprehension
// counterpart of TestArrayComprehension_MapSourceBindsKey.
func TestObjectComprehension_MapSourceBindsKey(t *testing.T) {
	input := `let obj = {a: 1, b: 2}; {k: k for k in obj}`

	parser := NewParser(input)
	program := parser.ParseProgram()
	env := NewEnvironment()

	var result interface{}
	for _, stmt := range program.Statements {
		result = stmt.Eval(env)
	}

	m, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map[string]interface{}, got %T: %v", result, result)
	}

	expected := map[string]interface{}{"a": "a", "b": "b"}
	if !reflect.DeepEqual(m, expected) {
		t.Errorf("expected %v, got %v", expected, m)
	}
}

func TestArrayComprehension_Basic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []interface{}
	}{
		{
			name:     "simple transform",
			input:    "let nums = [1, 2, 3]; [x * x for x in nums]",
			expected: []interface{}{float64(1), float64(4), float64(9)},
		},
		{
			name:     "with filter",
			input:    "let nums = [1, 2, 3, 4, 5]; [x for x in nums if x % 2 == 0]",
			expected: []interface{}{float64(2), float64(4)},
		},
		{
			name:     "empty source array",
			input:    "let nums = []; [x for x in nums]",
			expected: []interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			program := parser.ParseProgram()
			env := NewEnvironment()

			var result interface{}
			for _, stmt := range program.Statements {
				result = stmt.Eval(env)
			}

			arr, ok := result.([]interface{})
			if !ok {
				t.Fatalf("expected []interface{}, got %T: %v", result, result)
			}

			if len(arr) != len(tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, arr)
			}
			for i := range arr {
				if arr[i] != tt.expected[i] {
					t.Errorf("index %d: expected %v, got %v", i, tt.expected[i], arr[i])
				}
			}
		})
	}
}

func TestArrayComprehension_NonIterableSourcePanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic when comprehending over a non-iterable value")
		}
	}()

	input := `let n = 42; [x for x in n]`
	parser := NewParser(input)
	program := parser.ParseProgram()
	env := NewEnvironment()

	for _, stmt := range program.Statements {
		stmt.Eval(env)
	}
}
