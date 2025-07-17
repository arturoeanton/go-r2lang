package r2libs

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func TestConvertToGoFormat(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
		name     string
	}{
		{"YYYY-MM-DD", "2006-01-02", "ISO date format"},
		{"DD/MM/YYYY", "02/01/2006", "European date format"},
		{"MM/DD/YYYY", "01/02/2006", "American date format"},
		{"YYYY-MM-DD HH:mm:ss", "2006-01-02 15:04:05", "DateTime format"},
		{"DD-MM-YY", "02-01-06", "Short year format"},
		{"YYYY", "2006", "Year only"},
		{"MM", "01", "Month only"},
		{"DD", "02", "Day only"},
		{"HH:mm:ss", "15:04:05", "Time only"},
		{"YYYY-MM-DD'T'HH:mm:ss", "2006-01-02T15:04:05", "ISO datetime with T"},
		{"YYYY-MM-DD'T'HH:mm:ss.SSS'Z'", "2006-01-02T15:04:05.000Z", "ISO datetime with milliseconds"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := r2core.ConvertToGoFormat(tc.input)
			if result != tc.expected {
				t.Errorf("ConvertToGoFormat(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestDateFormatComprehensive(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterDate(env)
	dateModuleObj, _ := env.Get("date")
	dateModule := dateModuleObj.(map[string]interface{})
	DateFunc := dateModule["Date"].(r2core.BuiltinFunction)
	dateObj := DateFunc().(map[string]interface{})
	createFunc := dateObj["create"].(r2core.BuiltinFunction)
	formatFunc := dateObj["format"].(r2core.BuiltinFunction)

	// Test date: 2024-07-15 14:30:25
	args := []interface{}{
		float64(2024),
		float64(7),
		float64(15),
		float64(14),
		float64(30),
		float64(25),
	}
	dateValue := createFunc(args...).(*r2core.DateValue)

	testCases := []struct {
		format   string
		expected string
		name     string
	}{
		{"YYYY-MM-DD", "2024-07-15", "ISO date"},
		{"DD/MM/YYYY", "15/07/2024", "European date"},
		{"MM/DD/YYYY", "07/15/2024", "American date"},
		{"YYYY-MM-DD HH:mm:ss", "2024-07-15 14:30:25", "Full datetime"},
		{"DD-MM-YY", "15-07-24", "Short year"},
		{"YYYY", "2024", "Year only"},
		{"MM", "07", "Month only"},
		{"DD", "15", "Day only"},
		{"HH:mm:ss", "14:30:25", "Time only"},
		{"YY-MM-DD", "24-07-15", "Short year first"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			formatArgs := []interface{}{
				dateValue,
				tc.format,
			}
			result := formatFunc(formatArgs...)
			formattedString, ok := result.(string)
			if !ok {
				t.Fatalf("format() should return a string, got %T", result)
			}
			if formattedString != tc.expected {
				t.Errorf("format(%q) = %q, want %q", tc.format, formattedString, tc.expected)
			}
		})
	}
}

func TestDateFormatEdgeCases(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterDate(env)
	dateModuleObj, _ := env.Get("date")
	dateModule := dateModuleObj.(map[string]interface{})
	DateFunc := dateModule["Date"].(r2core.BuiltinFunction)
	dateObj := DateFunc().(map[string]interface{})
	createFunc := dateObj["create"].(r2core.BuiltinFunction)
	formatFunc := dateObj["format"].(r2core.BuiltinFunction)

	// Test various years including edge cases
	testYears := []int{2024, 2025, 2020, 2000, 1999, 2030}

	for _, year := range testYears {
		t.Run(fmt.Sprintf("Year_%d", year), func(t *testing.T) {
			args := []interface{}{
				float64(year),
				float64(1),
				float64(1),
			}
			dateValue := createFunc(args...).(*r2core.DateValue)

			formatArgs := []interface{}{
				dateValue,
				"YYYY-MM-DD",
			}
			result := formatFunc(formatArgs...)
			formattedString, ok := result.(string)
			if !ok {
				t.Fatalf("format() should return a string, got %T", result)
			}

			expectedYear := fmt.Sprintf("%04d", year)
			if !strings.HasPrefix(formattedString, expectedYear) {
				t.Errorf("format(YYYY-MM-DD) for year %d = %q, should start with %q", year, formattedString, expectedYear)
			}
		})
	}
}

func TestDateFormatDirectConversion(t *testing.T) {
	// Test the ConvertToGoFormat function directly
	testTime := time.Date(2024, 7, 15, 14, 30, 25, 0, time.UTC)

	// Test that YYYY properly converts to 2006
	goFormat := r2core.ConvertToGoFormat("YYYY-MM-DD")
	result := testTime.Format(goFormat)
	expected := "2024-07-15"

	if result != expected {
		t.Errorf("Direct conversion test failed: got %q, want %q", result, expected)
		t.Errorf("Go format was: %q", goFormat)
	}
}
