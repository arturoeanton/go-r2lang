package r2core

import (
	"fmt"
	"strings"
	"time"
)

// DateValue representa un valor de fecha en R2Lang
type DateValue struct {
	Time     time.Time
	Location *time.Location
}

func (dv *DateValue) String() string {
	return dv.Time.Format("2006-01-02T15:04:05Z07:00")
}

func (dv *DateValue) Type() string {
	return "Date"
}

func (dv *DateValue) Equals(other interface{}) bool {
	if otherDate, ok := other.(*DateValue); ok {
		return dv.Time.Equal(otherDate.Time)
	}
	return false
}

// Métodos de DateValue para acceso a componentes
func (dv *DateValue) Year() int {
	return dv.Time.Year()
}

func (dv *DateValue) Month() int {
	return int(dv.Time.Month())
}

func (dv *DateValue) Day() int {
	return dv.Time.Day()
}

func (dv *DateValue) Hour() int {
	return dv.Time.Hour()
}

func (dv *DateValue) Minute() int {
	return dv.Time.Minute()
}

func (dv *DateValue) Second() int {
	return dv.Time.Second()
}

func (dv *DateValue) Weekday() int {
	return int(dv.Time.Weekday())
}

func (dv *DateValue) YearDay() int {
	return dv.Time.YearDay()
}

// DurationValue representa una duración en R2Lang
type DurationValue struct {
	Duration time.Duration
}

func (dv *DurationValue) String() string {
	return dv.Duration.String()
}

func (dv *DurationValue) Type() string {
	return "Duration"
}

func (dv *DurationValue) Days() float64 {
	return dv.Duration.Hours() / 24
}

func (dv *DurationValue) Hours() float64 {
	return dv.Duration.Hours()
}

func (dv *DurationValue) Minutes() float64 {
	return dv.Duration.Minutes()
}

func (dv *DurationValue) Seconds() float64 {
	return dv.Duration.Seconds()
}

func (dv *DurationValue) Milliseconds() float64 {
	return float64(dv.Duration.Nanoseconds()) / 1e6
}

// NewDateValue crea un nuevo DateValue
func NewDateValue(t time.Time) *DateValue {
	return &DateValue{
		Time:     t,
		Location: t.Location(),
	}
}

// NewDateValueWithLocation crea un DateValue con zona horaria específica
func NewDateValueWithLocation(t time.Time, loc *time.Location) *DateValue {
	return &DateValue{
		Time:     t.In(loc),
		Location: loc,
	}
}

// NewDurationValue crea un nuevo DurationValue
func NewDurationValue(d time.Duration) *DurationValue {
	return &DurationValue{Duration: d}
}

// ParseDateLiteral parsea un literal de fecha como @2024-12-25
func ParseDateLiteral(dateStr string) (*DateValue, error) {
	// Remover @ si está presente
	if strings.HasPrefix(dateStr, "@") {
		dateStr = dateStr[1:]
	}
	
	// Remover comillas si están presentes
	dateStr = strings.Trim(dateStr, "\"'")
	
	// Intentar diferentes formatos
	formats := []string{
		"2006-01-02",                    // @2024-12-25
		"2006-01-02 15:04:05",          // @2024-12-25 14:30:00
		"2006-01-02T15:04:05",          // @2024-12-25T14:30:00
		"2006-01-02T15:04:05Z",         // @2024-12-25T14:30:00Z
		"2006-01-02T15:04:05Z07:00",    // @2024-12-25T14:30:00-05:00
		"2006-01-02T15:04:05.000Z",     // @2024-12-25T14:30:00.000Z
		"2006-01-02T15:04:05.000Z07:00", // @2024-12-25T14:30:00.000-05:00
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return NewDateValue(t), nil
		}
	}
	
	return nil, fmt.Errorf("unable to parse date: %s", dateStr)
}

// FormatDate formatea una fecha usando un patrón específico
func (dv *DateValue) Format(pattern string) string {
	// Convertir patrones comunes a formato Go
	goFormat := ConvertToGoFormat(pattern)
	return dv.Time.Format(goFormat)
}

// convertToGoFormat convierte patrones de fecha comunes a formato Go
func ConvertToGoFormat(pattern string) string {
	// Mapeo de patrones comunes a formato Go
	replacements := map[string]string{
		"YYYY": "2006",
		"YY":   "06",
		"MM":   "01",
		"DD":   "02",
		"HH":   "15",
		"mm":   "04",
		"ss":   "05",
		"SSS":  "000",
		"Z":    "Z07:00",
	}
	
	result := pattern
	for pattern, goPattern := range replacements {
		result = strings.ReplaceAll(result, pattern, goPattern)
	}
	
	// Manejar nombres de meses y días
	monthReplacements := map[string]string{
		"MMMM": "January",
		"MMM":  "Jan",
		"dddd": "Monday",
		"ddd":  "Mon",
	}
	
	for pattern, goPattern := range monthReplacements {
		result = strings.ReplaceAll(result, pattern, goPattern)
	}
	
	return result
}

// Add suma una duración a una fecha
func (dv *DateValue) Add(duration *DurationValue) *DateValue {
	return NewDateValue(dv.Time.Add(duration.Duration))
}

// Sub resta otra fecha o duración
func (dv *DateValue) Sub(other interface{}) interface{} {
	switch v := other.(type) {
	case *DateValue:
		// Fecha - Fecha = Duración
		diff := dv.Time.Sub(v.Time)
		return NewDurationValue(diff)
	case *DurationValue:
		// Fecha - Duración = Fecha
		return NewDateValue(dv.Time.Add(-v.Duration))
	default:
		panic("Cannot subtract non-date/duration from date")
	}
}

// Compare compara esta fecha con otra
func (dv *DateValue) Compare(other *DateValue) int {
	if dv.Time.Before(other.Time) {
		return -1
	} else if dv.Time.After(other.Time) {
		return 1
	}
	return 0
}

// ToTimezone convierte la fecha a otra zona horaria
func (dv *DateValue) ToTimezone(locationName string) (*DateValue, error) {
	loc, err := time.LoadLocation(locationName)
	if err != nil {
		return nil, err
	}
	
	return NewDateValueWithLocation(dv.Time, loc), nil
}

// IsWeekend verifica si la fecha es fin de semana
func (dv *DateValue) IsWeekend() bool {
	weekday := dv.Time.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// StartOf retorna el inicio de la unidad de tiempo especificada
func (dv *DateValue) StartOf(unit string) *DateValue {
	t := dv.Time
	switch strings.ToLower(unit) {
	case "day":
		t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	case "week":
		// Ir al lunes de la semana
		days := int(t.Weekday()) - 1
		if days < 0 {
			days = 6 // Domingo
		}
		t = t.AddDate(0, 0, -days)
		t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	case "month":
		t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	case "year":
		t = time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
	case "hour":
		t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
	case "minute":
		t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
	}
	return NewDateValue(t)
}

// EndOf retorna el final de la unidad de tiempo especificada
func (dv *DateValue) EndOf(unit string) *DateValue {
	t := dv.Time
	switch strings.ToLower(unit) {
	case "day":
		t = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
	case "week":
		// Ir al domingo de la semana
		days := 7 - int(t.Weekday())
		if days == 7 {
			days = 0
		}
		t = t.AddDate(0, 0, days)
		t = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
	case "month":
		// Último día del mes
		nextMonth := t.AddDate(0, 1, 0)
		lastDay := nextMonth.AddDate(0, 0, -1).Day()
		t = time.Date(t.Year(), t.Month(), lastDay, 23, 59, 59, 999999999, t.Location())
	case "year":
		t = time.Date(t.Year(), 12, 31, 23, 59, 59, 999999999, t.Location())
	case "hour":
		t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 59, 59, 999999999, t.Location())
	case "minute":
		t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 59, 999999999, t.Location())
	}
	return NewDateValue(t)
}

// AddYears suma años a la fecha
func (dv *DateValue) AddYears(years int) *DateValue {
	return NewDateValue(dv.Time.AddDate(years, 0, 0))
}

// AddMonths suma meses a la fecha
func (dv *DateValue) AddMonths(months int) *DateValue {
	return NewDateValue(dv.Time.AddDate(0, months, 0))
}

// AddDays suma días a la fecha
func (dv *DateValue) AddDays(days int) *DateValue {
	return NewDateValue(dv.Time.AddDate(0, 0, days))
}

// IsLeapYear verifica si el año de la fecha es bisiesto
func (dv *DateValue) IsLeapYear() bool {
	year := dv.Time.Year()
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// GetQuarter retorna el trimestre de la fecha (1-4)
func (dv *DateValue) GetQuarter() int {
	month := int(dv.Time.Month())
	return ((month - 1) / 3) + 1
}

// Unix retorna el timestamp Unix
func (dv *DateValue) Unix() int64 {
	return dv.Time.Unix()
}

// UnixMilli retorna el timestamp Unix en milisegundos
func (dv *DateValue) UnixMilli() int64 {
	return dv.Time.UnixMilli()
}

// ToISOString retorna la fecha en formato ISO 8601
func (dv *DateValue) ToISOString() string {
	return dv.Time.Format(time.RFC3339)
}

// ToDateString retorna solo la parte de fecha
func (dv *DateValue) ToDateString() string {
	return dv.Time.Format("2006-01-02")
}

// ToTimeString retorna solo la parte de tiempo
func (dv *DateValue) ToTimeString() string {
	return dv.Time.Format("15:04:05")
}

// CreateDurationFromDays crea una duración desde días
func CreateDurationFromDays(days float64) *DurationValue {
	return NewDurationValue(time.Duration(days * 24 * float64(time.Hour)))
}

// CreateDurationFromHours crea una duración desde horas
func CreateDurationFromHours(hours float64) *DurationValue {
	return NewDurationValue(time.Duration(hours * float64(time.Hour)))
}

// CreateDurationFromMinutes crea una duración desde minutos
func CreateDurationFromMinutes(minutes float64) *DurationValue {
	return NewDurationValue(time.Duration(minutes * float64(time.Minute)))
}

// CreateDurationFromSeconds crea una duración desde segundos
func CreateDurationFromSeconds(seconds float64) *DurationValue {
	return NewDurationValue(time.Duration(seconds * float64(time.Second)))
}

// ParseDurationString parsea una duración desde string como "1h30m"
func ParseDurationString(s string) (*DurationValue, error) {
	duration, err := time.ParseDuration(s)
	if err != nil {
		return nil, err
	}
	return NewDurationValue(duration), nil
}

// CreateDate crea una fecha desde componentes
func CreateDate(year, month, day int, hour, minute, second int, location *time.Location) *DateValue {
	if location == nil {
		location = time.Local
	}
	t := time.Date(year, time.Month(month), day, hour, minute, second, 0, location)
	return NewDateValue(t)
}

// CreateDateFromTimestamp crea una fecha desde timestamp Unix
func CreateDateFromTimestamp(timestamp int64) *DateValue {
	t := time.Unix(timestamp, 0)
	return NewDateValue(t)
}

// CreateDateFromMilliTimestamp crea una fecha desde timestamp en milisegundos
func CreateDateFromMilliTimestamp(timestamp int64) *DateValue {
	t := time.UnixMilli(timestamp)
	return NewDateValue(t)
}