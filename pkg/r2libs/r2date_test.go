package r2libs

import (
	"testing"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func TestDateNow(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterDate(env)
	dateObj, _ := env.Get("Date")
	nowFunc := dateObj.(map[string]interface{})["now"].(r2core.BuiltinFunction)
	result := nowFunc()
	dateValue, ok := result.(*r2core.DateValue)
	if !ok {
		t.Fatalf("now() should return a DateValue, got %T", result)
	}
	if time.Since(dateValue.Time) > time.Second {
		t.Errorf("Date.now() is not close to the current time")
	}
}

func TestDateCreate(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterDate(env)
	dateObj, _ := env.Get("Date")
	createFunc := dateObj.(map[string]interface{})["create"].(r2core.BuiltinFunction)
	args := []interface{}{
		float64(2024),
		float64(7),
		float64(15),
	}
	result := createFunc(args...)
	dateValue, ok := result.(*r2core.DateValue)
	if !ok {
		t.Fatalf("create() should return a DateValue, got %T", result)
	}
	if dateValue.Time.Year() != 2024 || dateValue.Time.Month() != 7 || dateValue.Time.Day() != 15 {
		t.Errorf("Date.create() did not create the correct date")
	}
}

func TestDateFormat(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterDate(env)
	dateObj, _ := env.Get("Date")
	createFunc := dateObj.(map[string]interface{})["create"].(r2core.BuiltinFunction)
	formatFunc := dateObj.(map[string]interface{})["format"].(r2core.BuiltinFunction)
	args := []interface{}{
		float64(2024),
		float64(7),
		float64(15),
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
	if formattedString != "2024-07-15" {
		t.Errorf("format() returned incorrect string, got %s, want %s", formattedString, "2024-07-15")
	}
}

func TestDateTimezone(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterDate(env)
	dateObj, _ := env.Get("Date")
	timezoneFunc := dateObj.(map[string]interface{})["timezone"].(r2core.BuiltinFunction)
	args := []interface{}{
		"America/New_York",
		float64(2024),
		float64(7),
		float64(15),
	}
	result := timezoneFunc(args...)
	dateValue, ok := result.(*r2core.DateValue)
	if !ok {
		t.Fatalf("timezone() should return a DateValue, got %T", result)
	}
	if dateValue.Time.Year() != 2024 || dateValue.Time.Month() != 7 || dateValue.Time.Day() != 15 {
		t.Errorf("timezone() did not create the correct date")
	}
	if dateValue.Time.Location().String() != "America/New_York" {
		t.Errorf("timezone() did not set the correct timezone")
	}
}

func TestDateToTimezone(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterDate(env)
	dateObj, _ := env.Get("Date")
	createFunc := dateObj.(map[string]interface{})["create"].(r2core.BuiltinFunction)
	toTimezoneFunc := dateObj.(map[string]interface{})["toTimezone"].(r2core.BuiltinFunction)
	args := []interface{}{
		float64(2024),
		float64(7),
		float64(15),
	}
	dateValue := createFunc(args...).(*r2core.DateValue)
	toTimezoneArgs := []interface{}{
		dateValue,
		"America/New_York",
	}
	result := toTimezoneFunc(toTimezoneArgs...)
	convertedDate, ok := result.(*r2core.DateValue)
	if !ok {
		t.Fatalf("toTimezone() should return a DateValue, got %T", result)
	}
	if convertedDate.Time.Location().String() != "America/New_York" {
		t.Errorf("toTimezone() did not convert to the correct timezone")
	}
}
