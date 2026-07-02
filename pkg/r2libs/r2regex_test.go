package r2libs

import (
	"reflect"
	"testing"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func regexModule(t *testing.T) map[string]interface{} {
	t.Helper()
	env := r2core.NewEnvironment()
	RegisterRegex(env)
	moduleObj, ok := env.Get("regex")
	if !ok {
		t.Fatal("regex module not found")
	}
	module, ok := moduleObj.(map[string]interface{})
	if !ok {
		t.Fatal("regex module has unexpected type")
	}
	return module
}

func TestRegexTest(t *testing.T) {
	module := regexModule(t)
	testFunc := module["test"].(r2core.BuiltinFunction)

	tests := []struct {
		name     string
		pattern  string
		str      string
		expected bool
	}{
		{"matches", `\d+`, "abc123", true},
		{"no match", `\d+`, "abcdef", false},
		{"empty string against optional pattern", `a*`, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testFunc(tt.pattern, tt.str)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestRegexMatch(t *testing.T) {
	module := regexModule(t)
	matchFunc := module["match"].(r2core.BuiltinFunction)

	t.Run("first match returned", func(t *testing.T) {
		result := matchFunc(`\d+`, "abc123def456")
		if result != "123" {
			t.Errorf("expected 123, got %v", result)
		}
	})

	t.Run("no match returns nil", func(t *testing.T) {
		result := matchFunc(`\d+`, "abcdef")
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})

	t.Run("empty match found", func(t *testing.T) {
		result := matchFunc(`x*`, "abc")
		if result != "" {
			t.Errorf("expected empty string match, got %v", result)
		}
	})
}

func TestRegexMatchAll(t *testing.T) {
	module := regexModule(t)
	matchAllFunc := module["matchAll"].(r2core.BuiltinFunction)

	t.Run("multiple matches", func(t *testing.T) {
		result := matchAllFunc(`\d+`, "abc123def456")
		expected := []interface{}{"123", "456"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("no matches returns empty array", func(t *testing.T) {
		result := matchAllFunc(`\d+`, "abcdef")
		arr, ok := result.([]interface{})
		if !ok {
			t.Fatalf("expected []interface{}, got %T", result)
		}
		if len(arr) != 0 {
			t.Errorf("expected empty array, got %v", arr)
		}
	})
}

func TestRegexGroups(t *testing.T) {
	module := regexModule(t)
	groupsFunc := module["groups"].(r2core.BuiltinFunction)

	t.Run("captures groups", func(t *testing.T) {
		result := groupsFunc(`(\w+)@(\w+)\.com`, "user@example.com")
		expected := []interface{}{"user@example.com", "user", "example"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("no match returns nil", func(t *testing.T) {
		result := groupsFunc(`(\d+)`, "abcdef")
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

func TestRegexReplace(t *testing.T) {
	module := regexModule(t)
	replaceFunc := module["replace"].(r2core.BuiltinFunction)

	t.Run("replaces first match only", func(t *testing.T) {
		result := replaceFunc(`\d+`, "a1b2c3", "X")
		if result != "aXb2c3" {
			t.Errorf("expected aXb2c3, got %v", result)
		}
	})

	t.Run("supports capture group references", func(t *testing.T) {
		result := replaceFunc(`(\w+)@(\w+)`, "user@host other@thing", "$2@$1")
		if result != "host@user other@thing" {
			t.Errorf("expected 'host@user other@thing', got %v", result)
		}
	})

	t.Run("no match returns original string", func(t *testing.T) {
		result := replaceFunc(`\d+`, "abcdef", "X")
		if result != "abcdef" {
			t.Errorf("expected abcdef, got %v", result)
		}
	})
}

func TestRegexReplaceAll(t *testing.T) {
	module := regexModule(t)
	replaceAllFunc := module["replaceAll"].(r2core.BuiltinFunction)

	t.Run("replaces every match", func(t *testing.T) {
		result := replaceAllFunc(`\d+`, "a1b2c3", "X")
		if result != "aXbXcX" {
			t.Errorf("expected aXbXcX, got %v", result)
		}
	})

	t.Run("no match returns original string", func(t *testing.T) {
		result := replaceAllFunc(`\d+`, "abcdef", "X")
		if result != "abcdef" {
			t.Errorf("expected abcdef, got %v", result)
		}
	})
}

func TestRegexSplit(t *testing.T) {
	module := regexModule(t)
	splitFunc := module["split"].(r2core.BuiltinFunction)

	t.Run("splits on pattern", func(t *testing.T) {
		result := splitFunc(`,\s*`, "a, b,c ,  d")
		expected := []interface{}{"a", "b", "c ", "d"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("no match returns whole string as single element", func(t *testing.T) {
		result := splitFunc(`;`, "abc")
		expected := []interface{}{"abc"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}

func TestRegexEscape(t *testing.T) {
	module := regexModule(t)
	escapeFunc := module["escape"].(r2core.BuiltinFunction)

	t.Run("escapes metacharacters", func(t *testing.T) {
		result := escapeFunc("a.b*c?")
		if result != `a\.b\*c\?` {
			t.Errorf("expected 'a\\.b\\*c\\?', got %v", result)
		}
	})

	t.Run("no metacharacters unchanged", func(t *testing.T) {
		result := escapeFunc("abc")
		if result != "abc" {
			t.Errorf("expected abc, got %v", result)
		}
	})
}

func TestRegexInvalidPatternPanics(t *testing.T) {
	module := regexModule(t)

	cases := []struct {
		name string
		fn   r2core.BuiltinFunction
		args []interface{}
	}{
		{"test", module["test"].(r2core.BuiltinFunction), []interface{}{"(", "abc"}},
		{"match", module["match"].(r2core.BuiltinFunction), []interface{}{"(", "abc"}},
		{"matchAll", module["matchAll"].(r2core.BuiltinFunction), []interface{}{"(", "abc"}},
		{"groups", module["groups"].(r2core.BuiltinFunction), []interface{}{"(", "abc"}},
		{"replace", module["replace"].(r2core.BuiltinFunction), []interface{}{"(", "abc", "x"}},
		{"replaceAll", module["replaceAll"].(r2core.BuiltinFunction), []interface{}{"(", "abc", "x"}},
		{"split", module["split"].(r2core.BuiltinFunction), []interface{}{"(", "abc"}},
	}

	for _, tt := range cases {
		t.Run(tt.name+" invalid pattern panics", func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("expected panic for invalid pattern")
				}
			}()
			tt.fn(tt.args...)
		})
	}
}
