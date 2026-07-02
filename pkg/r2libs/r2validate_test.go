package r2libs

import (
	"testing"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func validateModuleForTest(t *testing.T) map[string]interface{} {
	t.Helper()
	env := r2core.NewEnvironment()
	RegisterValidate(env)
	modObj, ok := env.Get("validate")
	if !ok {
		t.Fatal("validate module not found")
	}
	return modObj.(map[string]interface{})
}

func TestValidateIsEmail(t *testing.T) {
	mod := validateModuleForTest(t)
	isEmail := mod["isEmail"].(r2core.BuiltinFunction)

	cases := []struct {
		input    string
		expected bool
	}{
		{"user@example.com", true},
		{"first.last+tag@sub.example.co", true},
		{"not-an-email", false},
		{"missing-domain@", false},
		{"@missing-local.com", false},
		{"", false},
		{"spaces in@email.com", false},
	}
	for _, tc := range cases {
		if got := isEmail(tc.input).(bool); got != tc.expected {
			t.Errorf("isEmail(%q) = %v, want %v", tc.input, got, tc.expected)
		}
	}
}

func TestValidateIsURL(t *testing.T) {
	mod := validateModuleForTest(t)
	isURL := mod["isURL"].(r2core.BuiltinFunction)

	cases := []struct {
		input    string
		expected bool
	}{
		{"https://example.com", true},
		{"http://example.com/path?query=1", true},
		{"ftp://files.example.com", true},
		{"not a url", false},
		{"example.com", false},
		{"", false},
		{"http://", false},
	}
	for _, tc := range cases {
		if got := isURL(tc.input).(bool); got != tc.expected {
			t.Errorf("isURL(%q) = %v, want %v", tc.input, got, tc.expected)
		}
	}
}

func TestValidateIsIP(t *testing.T) {
	mod := validateModuleForTest(t)
	isIP := mod["isIP"].(r2core.BuiltinFunction)

	cases := []struct {
		input    string
		expected bool
	}{
		{"127.0.0.1", true},
		{"192.168.1.1", true},
		{"::1", true},
		{"2001:db8::1", true},
		{"999.999.999.999", false},
		{"not-an-ip", false},
		{"", false},
	}
	for _, tc := range cases {
		if got := isIP(tc.input).(bool); got != tc.expected {
			t.Errorf("isIP(%q) = %v, want %v", tc.input, got, tc.expected)
		}
	}
}
