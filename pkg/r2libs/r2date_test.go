package r2libs

import (
	"testing"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func TestDateConstructor(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterDate(env)
	dateModuleObj, _ := env.Get("date")
	dateModule := dateModuleObj.(map[string]interface{})
	DateFunc := dateModule["Date"].(r2core.BuiltinFunction)
	
	result := DateFunc()
	dateValue, ok := result.(*r2core.DateValue)
	if !ok {
		t.Fatalf("Date() should return a DateValue, got %T", result)
	}
	if time.Since(dateValue.Time) > time.Second {
		t.Errorf("Date() is not close to the current time")
	}
}

func TestDateConstructorWithArgs(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterDate(env)
	dateModuleObj, _ := env.Get("date")
	dateModule := dateModuleObj.(map[string]interface{})
	DateFunc := dateModule["Date"].(r2core.BuiltinFunction)
	
	args := []interface{}{
		float64(2024),
		float64(6),
		float64(15),
	}
	result := DateFunc(args...)
	dateValue, ok := result.(*r2core.DateValue)
	if !ok {
		t.Fatalf("Date() should return a DateValue, got %T", result)
	}
	if dateValue.Time.Year() != 2024 || dateValue.Time.Month() != 7 || dateValue.Time.Day() != 15 {
		t.Errorf("Date() did not create the correct date")
	}
}

func TestDateNow(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterDate(env)
	dateModuleObj, _ := env.Get("date")
	dateModule := dateModuleObj.(map[string]interface{})
	DateFunc := dateModule["Date"].(r2core.BuiltinFunction)
	dateObj := DateFunc().(map[string]interface{})
	nowFunc := dateObj["now"].(r2core.BuiltinFunction)
	result := nowFunc()
	timestamp, ok := result.(float64)
	if !ok {
		t.Fatalf("now() should return a timestamp, got %T", result)
	}
	if time.Since(time.UnixMilli(int64(timestamp))) > time.Second {
		t.Errorf("Date.now() is not close to the current time")
	}
}

func TestDateGetMethods(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterDate(env)
	dateModuleObj, _ := env.Get("date")
	dateModule := dateModuleObj.(map[string]interface{})
	DateFunc := dateModule["Date"].(r2core.BuiltinFunction)
	dateObj := DateFunc().(map[string]interface{})
	
	args := []interface{}{
		float64(2024),
		float64(6),
		float64(15),
		float64(14),
		float64(30),
		float64(25),
	}
	dateValue := DateFunc(args...).(*r2core.DateValue)
	
	getFullYearFunc := dateObj["getFullYear"].(r2core.BuiltinFunction)
	year := getFullYearFunc(dateValue).(float64)
	if year != 2024 {
		t.Errorf("getFullYear() returned %f, expected 2024", year)
	}
	
	getMonthFunc := dateObj["getMonth"].(r2core.BuiltinFunction)
	month := getMonthFunc(dateValue).(float64)
	if month != 6 {
		t.Errorf("getMonth() returned %f, expected 6", month)
	}
	
	getDateFunc := dateObj["getDate"].(r2core.BuiltinFunction)
	date := getDateFunc(dateValue).(float64)
	if date != 15 {
		t.Errorf("getDate() returned %f, expected 15", date)
	}
	
	getHoursFunc := dateObj["getHours"].(r2core.BuiltinFunction)
	hours := getHoursFunc(dateValue).(float64)
	if hours != 14 {
		t.Errorf("getHours() returned %f, expected 14", hours)
	}
	
	getMinutesFunc := dateObj["getMinutes"].(r2core.BuiltinFunction)
	minutes := getMinutesFunc(dateValue).(float64)
	if minutes != 30 {
		t.Errorf("getMinutes() returned %f, expected 30", minutes)
	}
	
	getSecondsFunc := dateObj["getSeconds"].(r2core.BuiltinFunction)
	seconds := getSecondsFunc(dateValue).(float64)
	if seconds != 25 {
		t.Errorf("getSeconds() returned %f, expected 25", seconds)
	}
}

func TestDateSetMethods(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterDate(env)
	dateModuleObj, _ := env.Get("date")
	dateModule := dateModuleObj.(map[string]interface{})
	DateFunc := dateModule["Date"].(r2core.BuiltinFunction)
	dateObj := DateFunc().(map[string]interface{})
	
	args := []interface{}{
		float64(2024),
		float64(6),
		float64(15),
	}
	dateValue := DateFunc(args...).(*r2core.DateValue)
	
	setFullYearFunc := dateObj["setFullYear"].(r2core.BuiltinFunction)
	setFullYearFunc(dateValue, float64(2025))
	
	getFullYearFunc := dateObj["getFullYear"].(r2core.BuiltinFunction)
	year := getFullYearFunc(dateValue).(float64)
	if year != 2025 {
		t.Errorf("After setFullYear(2025), getFullYear() returned %f, expected 2025", year)
	}
	
	setMonthFunc := dateObj["setMonth"].(r2core.BuiltinFunction)
	setMonthFunc(dateValue, float64(11))
	
	getMonthFunc := dateObj["getMonth"].(r2core.BuiltinFunction)
	month := getMonthFunc(dateValue).(float64)
	if month != 11 {
		t.Errorf("After setMonth(11), getMonth() returned %f, expected 11", month)
	}
}

func TestDateStringMethods(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterDate(env)
	dateModuleObj, _ := env.Get("date")
	dateModule := dateModuleObj.(map[string]interface{})
	DateFunc := dateModule["Date"].(r2core.BuiltinFunction)
	dateObj := DateFunc().(map[string]interface{})
	
	args := []interface{}{
		float64(2024),
		float64(6),
		float64(15),
	}
	dateValue := DateFunc(args...).(*r2core.DateValue)
	
	toISOStringFunc := dateObj["toISOString"].(r2core.BuiltinFunction)
	isoString := toISOStringFunc(dateValue).(string)
	if len(isoString) == 0 {
		t.Error("toISOString() returned empty string")
	}
	
	toDateStringFunc := dateObj["toDateString"].(r2core.BuiltinFunction)
	dateString := toDateStringFunc(dateValue).(string)
	if len(dateString) == 0 {
		t.Error("toDateString() returned empty string")
	}
	
	toTimeStringFunc := dateObj["toTimeString"].(r2core.BuiltinFunction)
	timeString := toTimeStringFunc(dateValue).(string)
	if len(timeString) == 0 {
		t.Error("toTimeString() returned empty string")
	}
}

func TestDateAddMethods(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterDate(env)
	dateModuleObj, _ := env.Get("date")
	dateModule := dateModuleObj.(map[string]interface{})
	DateFunc := dateModule["Date"].(r2core.BuiltinFunction)
	dateObj := DateFunc().(map[string]interface{})
	
	args := []interface{}{
		float64(2024),
		float64(6),
		float64(15),
	}
	dateValue := DateFunc(args...).(*r2core.DateValue)
	
	addDaysFunc := dateObj["addDays"].(r2core.BuiltinFunction)
	newDate := addDaysFunc(dateValue, float64(10)).(*r2core.DateValue)
	
	if newDate.Time.Day() != 25 {
		t.Errorf("After adding 10 days, expected day 25, got %d", newDate.Time.Day())
	}
	
	addMonthsFunc := dateObj["addMonths"].(r2core.BuiltinFunction)
	newDate2 := addMonthsFunc(dateValue, float64(2)).(*r2core.DateValue)
	
	if newDate2.Time.Month() != time.September {
		t.Errorf("After adding 2 months, expected September, got %s", newDate2.Time.Month())
	}
}

func TestDateDiff(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterDate(env)
	dateModuleObj, _ := env.Get("date")
	dateModule := dateModuleObj.(map[string]interface{})
	DateFunc := dateModule["Date"].(r2core.BuiltinFunction)
	dateObj := DateFunc().(map[string]interface{})
	
	date1 := DateFunc(float64(2024), float64(6), float64(15)).(*r2core.DateValue)
	date2 := DateFunc(float64(2024), float64(6), float64(10)).(*r2core.DateValue)
	
	diffFunc := dateObj["diff"].(r2core.BuiltinFunction)
	diff := diffFunc(date1, date2, "days").(float64)
	
	if diff != 5 {
		t.Errorf("Expected diff of 5 days, got %f", diff)
	}
}

func TestDateFormat(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterDate(env)
	dateModuleObj, _ := env.Get("date")
	dateModule := dateModuleObj.(map[string]interface{})
	DateFunc := dateModule["Date"].(r2core.BuiltinFunction)
	dateObj := DateFunc().(map[string]interface{})
	
	args := []interface{}{
		float64(2024),
		float64(6),
		float64(15),
	}
	dateValue := DateFunc(args...).(*r2core.DateValue)
	
	formatFunc := dateObj["format"].(r2core.BuiltinFunction)
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
	
	args := []interface{}{
		float64(2024),
		float64(6),
		float64(15),
	}
	dateValue := DateFunc(args...).(*r2core.DateValue)
	
	toTimezoneFunc := dateObj["toTimezone"].(r2core.BuiltinFunction)
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

func TestDateParse(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterDate(env)
	dateModuleObj, _ := env.Get("date")
	dateModule := dateModuleObj.(map[string]interface{})
	DateFunc := dateModule["Date"].(r2core.BuiltinFunction)
	dateObj := DateFunc().(map[string]interface{})
	
	parseFunc := dateObj["parse"].(r2core.BuiltinFunction)
	result := parseFunc("2024-07-15")
	dateValue, ok := result.(*r2core.DateValue)
	if !ok {
		t.Fatalf("parse() should return a DateValue, got %T", result)
	}
	if dateValue.Time.Year() != 2024 || dateValue.Time.Month() != 7 || dateValue.Time.Day() != 15 {
		t.Errorf("parse() did not parse the correct date")
	}
}

func TestDateUTCMethods(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterDate(env)
	dateModuleObj, _ := env.Get("date")
	dateModule := dateModuleObj.(map[string]interface{})
	DateFunc := dateModule["Date"].(r2core.BuiltinFunction)
	dateObj := DateFunc().(map[string]interface{})
	
	utcFunc := dateObj["UTC"].(r2core.BuiltinFunction)
	utcDate := utcFunc(float64(2024), float64(6), float64(15)).(*r2core.DateValue)
	
	if utcDate.Time.Location() != time.UTC {
		t.Error("UTC() did not create a UTC date")
	}
	
	getUTCFullYearFunc := dateObj["getUTCFullYear"].(r2core.BuiltinFunction)
	year := getUTCFullYearFunc(utcDate).(float64)
	if year != 2024 {
		t.Errorf("getUTCFullYear() returned %f, expected 2024", year)
	}
}
