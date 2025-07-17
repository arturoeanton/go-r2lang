package r2libs

import (
	"testing"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func TestDateNow(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterDate(env)
	dateModuleObj, _ := env.Get("date")
	dateModule := dateModuleObj.(map[string]interface{})
	DateFunc := dateModule["Date"].(r2core.BuiltinFunction)
	dateObj := DateFunc().(map[string]interface{})
	nowFunc := dateObj["now"].(r2core.BuiltinFunction)
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
	dateModuleObj, _ := env.Get("date")
	dateModule := dateModuleObj.(map[string]interface{})
	DateFunc := dateModule["Date"].(r2core.BuiltinFunction)
	dateObj := DateFunc().(map[string]interface{})
	createFunc := dateObj["create"].(r2core.BuiltinFunction)
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
	dateModuleObj, _ := env.Get("date")
	dateModule := dateModuleObj.(map[string]interface{})
	DateFunc := dateModule["Date"].(r2core.BuiltinFunction)
	dateObj := DateFunc().(map[string]interface{})
	createFunc := dateObj["create"].(r2core.BuiltinFunction)
	formatFunc := dateObj["format"].(r2core.BuiltinFunction)
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
	dateModuleObj, _ := env.Get("date")
	dateModule := dateModuleObj.(map[string]interface{})
	DateFunc := dateModule["Date"].(r2core.BuiltinFunction)
	dateObj := DateFunc().(map[string]interface{})
	timezoneFunc := dateObj["timezone"].(r2core.BuiltinFunction)
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
	dateModuleObj, _ := env.Get("date")
	dateModule := dateModuleObj.(map[string]interface{})
	DateFunc := dateModule["Date"].(r2core.BuiltinFunction)
	dateObj := DateFunc().(map[string]interface{})
	createFunc := dateObj["create"].(r2core.BuiltinFunction)
	toTimezoneFunc := dateObj["toTimezone"].(r2core.BuiltinFunction)
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
