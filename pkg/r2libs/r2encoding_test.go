package r2libs

import (
	"testing"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func setupEncodingEnv(t *testing.T) (map[string]interface{}, map[string]interface{}) {
	t.Helper()
	env := r2core.NewEnvironment()
	RegisterEncoding(env)

	encodingObj, ok := env.Get("encoding")
	if !ok {
		t.Fatal("encoding module not found")
	}
	uuidObj, ok := env.Get("uuid")
	if !ok {
		t.Fatal("uuid module not found")
	}
	return encodingObj.(map[string]interface{}), uuidObj.(map[string]interface{})
}

func TestEncodingRoundTrips(t *testing.T) {
	encoding, _ := setupEncodingEnv(t)

	base64Encode := encoding["base64Encode"].(r2core.BuiltinFunction)
	base64Decode := encoding["base64Decode"].(r2core.BuiltinFunction)
	base64UrlEncode := encoding["base64UrlEncode"].(r2core.BuiltinFunction)
	base64UrlDecode := encoding["base64UrlDecode"].(r2core.BuiltinFunction)
	hexEncode := encoding["hexEncode"].(r2core.BuiltinFunction)
	hexDecode := encoding["hexDecode"].(r2core.BuiltinFunction)
	urlEncode := encoding["urlEncode"].(r2core.BuiltinFunction)
	urlDecode := encoding["urlDecode"].(r2core.BuiltinFunction)

	tests := []struct {
		name   string
		encode r2core.BuiltinFunction
		decode r2core.BuiltinFunction
		input  string
	}{
		{"base64 simple", base64Encode, base64Decode, "hello world"},
		{"base64 empty", base64Encode, base64Decode, ""},
		{"base64 special chars", base64Encode, base64Decode, "R2Lang!@#$%^&*()_+ áéíóú"},
		{"base64url simple", base64UrlEncode, base64UrlDecode, "hello world"},
		{"base64url binary-ish", base64UrlEncode, base64UrlDecode, "\x00\x01\xff\xfe"},
		{"hex simple", hexEncode, hexDecode, "hello world"},
		{"hex empty", hexEncode, hexDecode, ""},
		{"url query simple", urlEncode, urlDecode, "hello world & friends = 1"},
		{"url query special", urlEncode, urlDecode, "a=b&c=d?e#f"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := tt.encode(tt.input).(string)
			decoded := tt.decode(encoded).(string)
			if decoded != tt.input {
				t.Errorf("round trip failed: input %q, encoded %q, decoded %q", tt.input, encoded, decoded)
			}
		})
	}
}

func TestEncodingInvalidInputPanics(t *testing.T) {
	encoding, _ := setupEncodingEnv(t)

	tests := []struct {
		name    string
		decode  r2core.BuiltinFunction
		invalid string
	}{
		{"base64Decode invalid", encoding["base64Decode"].(r2core.BuiltinFunction), "not-valid-base64!!"},
		{"base64UrlDecode invalid", encoding["base64UrlDecode"].(r2core.BuiltinFunction), "not valid base64url!!"},
		{"hexDecode invalid", encoding["hexDecode"].(r2core.BuiltinFunction), "zzz-not-hex"},
		{"urlDecode invalid", encoding["urlDecode"].(r2core.BuiltinFunction), "%zz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("expected panic for invalid input %q", tt.invalid)
				}
			}()
			tt.decode(tt.invalid)
		})
	}
}

func TestEncodingUrlParse(t *testing.T) {
	encoding, _ := setupEncodingEnv(t)
	urlParse := encoding["urlParse"].(r2core.BuiltinFunction)

	result := urlParse("https://example.com/path/to/thing?foo=bar&baz=qux").(map[string]interface{})

	if result["scheme"] != "https" {
		t.Errorf("expected scheme https, got %v", result["scheme"])
	}
	if result["host"] != "example.com" {
		t.Errorf("expected host example.com, got %v", result["host"])
	}
	if result["path"] != "/path/to/thing" {
		t.Errorf("expected path /path/to/thing, got %v", result["path"])
	}
	query, ok := result["query"].(map[string]interface{})
	if !ok {
		t.Fatal("expected query to be a map")
	}
	if query["foo"] != "bar" {
		t.Errorf("expected query foo=bar, got %v", query["foo"])
	}
	if query["baz"] != "qux" {
		t.Errorf("expected query baz=qux, got %v", query["baz"])
	}
}

func TestUUIDV4(t *testing.T) {
	_, uuid := setupEncodingEnv(t)
	v4 := uuid["v4"].(r2core.BuiltinFunction)
	isValid := uuid["isValid"].(r2core.BuiltinFunction)

	seen := make(map[string]bool)
	for i := 0; i < 20; i++ {
		id := v4().(string)
		if seen[id] {
			t.Errorf("uuid.v4 generated a duplicate: %s", id)
		}
		seen[id] = true

		if !isValid(id).(bool) {
			t.Errorf("uuid.v4 generated %q which failed uuid.isValid", id)
		}
	}
}

func TestUUIDIsValid(t *testing.T) {
	_, uuid := setupEncodingEnv(t)
	isValid := uuid["isValid"].(r2core.BuiltinFunction)

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid uuid", "550e8400-e29b-41d4-a716-446655440000", true},
		{"valid uuid uppercase", "550E8400-E29B-41D4-A716-446655440000", true},
		{"missing dashes", "550e8400e29b41d4a716446655440000", false},
		{"too short", "550e8400-e29b-41d4-a716", false},
		{"garbage", "not-a-uuid-at-all", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValid(tt.input).(bool)
			if result != tt.expected {
				t.Errorf("isValid(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
