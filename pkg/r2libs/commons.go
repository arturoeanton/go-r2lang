package r2libs

import (
	"fmt"
	"strconv"
)

func toFloat(val interface{}) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case bool:
		if v {
			return 1
		}
		return 0
	case nil:
		return 0
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			panic("Cannot convert string to number:" + v)
		}
		return f
	}
	panic("Cannot convert value to number")
}
func toBool(val interface{}) bool {
	if val == nil {
		return false
	}
	switch v := val.(type) {
	case bool:
		return v
	case float64:
		return v != 0
	case int:
		return v != 0
	case string:
		return v != ""
	}
	return true
}

// Para unificar la lógica numérica en "=="
func isNumeric(v interface{}) bool {
	switch v.(type) {
	case float64, int:
		return true
	}
	return false
}

// Corrige la comparación "=="
func equals(a, b interface{}) bool {
	// Si ambos son numéricos, compare con toFloat
	if isNumeric(a) && isNumeric(b) {
		return toFloat(a) == toFloat(b)
	}
	// sino comparamos string/bool/nil
	switch aa := a.(type) {
	case string:
		if bb, ok := b.(string); ok {
			return aa == bb
		}
	case bool:
		if bb, ok := b.(bool); ok {
			return aa == bb
		}
	case nil:
		return b == nil
	}
	return false
}

func addValues(a, b interface{}) interface{} {

	if isNumeric(a) && isNumeric(b) {
		return toFloat(a) + toFloat(b)
	}

	if aa, ok := a.([]interface{}); ok {
		if bb, ok := b.([]interface{}); ok {
			return append(aa, bb...)
		}
		return append(aa, b)
	}

	if ab, ok := b.([]interface{}); ok {
		return append([]interface{}{a}, ab...)
	}

	// Si uno es string => concatenar
	if sa, ok := a.(string); ok {
		return sa + fmt.Sprint(b)
	}
	if sb, ok := b.(string); ok {
		return fmt.Sprint(a) + sb
	}
	return toFloat(a) + toFloat(b)
}
func subValues(a, b interface{}) interface{} {
	return toFloat(a) - toFloat(b)
}
func mulValues(a, b interface{}) interface{} {
	return toFloat(a) * toFloat(b)
}
func divValues(a, b interface{}) interface{} {
	den := toFloat(b)
	if den == 0 {
		panic("Division by zero")
	}
	return toFloat(a) / den
}
