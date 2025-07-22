package r2core

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

var (
	// Cache para conversiones de string a float
	stringToFloatCache   map[string]float64
	commonsStringCacheMu sync.RWMutex

	// Cache para números pequeños comunes
	intToFloatCache map[int]float64

	// Cache optimizado para números frecuentes
	frequentNumberCache map[string]float64
	frequentNumberMu    sync.RWMutex
)

func init() {
	// Inicializar caches
	stringToFloatCache = make(map[string]float64)
	intToFloatCache = make(map[int]float64)
	frequentNumberCache = make(map[string]float64)

	// Pre-poblar cache con números comunes
	for i := -1000; i <= 1000; i++ {
		intToFloatCache[i] = float64(i)
	}

	// Pre-poblar números frecuentes como strings
	frequentNumbers := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
		"100", "1000", "-1", "-2", "-3", "0.5", "1.5", "2.5"}
	for _, numStr := range frequentNumbers {
		if f, err := strconv.ParseFloat(numStr, 64); err == nil {
			frequentNumberCache[numStr] = f
		}
	}
}

func toFloat(val interface{}) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case int:
		// Usar cache para números pequeños
		if cached, ok := intToFloatCache[v]; ok {
			return cached
		}
		return float64(v)
	case bool:
		if v {
			return 1
		}
		return 0
	case nil:
		return 0
	case string:
		// Buscar primero en cache de números frecuentes (más rápido)
		frequentNumberMu.RLock()
		if cached, ok := frequentNumberCache[v]; ok {
			frequentNumberMu.RUnlock()
			return cached
		}
		frequentNumberMu.RUnlock()

		// Buscar en cache general
		commonsStringCacheMu.RLock()
		if cached, ok := stringToFloatCache[v]; ok {
			commonsStringCacheMu.RUnlock()
			return cached
		}
		commonsStringCacheMu.RUnlock()

		// Smart conversion: try to parse as number
		f, err := smartParseFloat(v)
		if err != nil {
			panic("Cannot convert string to number:" + v)
		}

		// Limitar tamaño del cache
		commonsStringCacheMu.Lock()
		if len(stringToFloatCache) < 10000 {
			stringToFloatCache[v] = f
		}
		commonsStringCacheMu.Unlock()
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
	// Fast path: evitar conversiones si ya son float64
	if af, ok := a.(float64); ok {
		if bf, ok := b.(float64); ok {
			return af + bf // Sin allocaciones extra
		}
	}

	// P5 Feature: Smart auto-conversion
	// Handle mixed types and smart numeric strings
	if shouldUseP5SmartConversion(a, b) {
		return smartToFloat(a) + smartToFloat(b)
	}

	if isNumeric(a) && isNumeric(b) {
		// Object pool desactivado para operaciones simples
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

	// P5 Feature: Smart string concatenation with fallback
	// Use smart conversion for better string handling
	return smartStringConcat(a, b)
}
func subValues(a, b interface{}) interface{} {
	// Fast path: evitar conversiones si ya son float64
	if af, ok := a.(float64); ok {
		if bf, ok := b.(float64); ok {
			return af - bf // Sin allocaciones extra
		}
	}
	// Object pool desactivado para operaciones simples
	return toFloat(a) - toFloat(b)
}
func mulValues(a, b interface{}) interface{} {
	// Fast path: evitar conversiones si ya son float64
	if af, ok := a.(float64); ok {
		if bf, ok := b.(float64); ok {
			return af * bf // Sin allocaciones extra
		}
	}
	// Object pool desactivado para operaciones simples
	return toFloat(a) * toFloat(b)
}
func divValues(a, b interface{}) interface{} {
	// Fast path: evitar conversiones si ya son float64
	if af, ok := a.(float64); ok {
		if bf, ok := b.(float64); ok {
			if bf == 0 {
				panic("Division by zero")
			}
			return af / bf // Sin allocaciones extra
		}
	}

	den := toFloat(b)
	if den == 0 {
		panic("Division by zero")
	}
	// Object pool desactivado para operaciones simples
	return toFloat(a) / den
}

func modValues(a, b interface{}) interface{} {
	// Fast path: evitar conversiones si ya son float64
	if af, ok := a.(float64); ok {
		if bf, ok := b.(float64); ok {
			if bf == 0 {
				panic("Modulo by zero")
			}
			return float64(int(af) % int(bf))
		}
	}

	den := toFloat(b)
	if den == 0 {
		panic("Modulo by zero")
	}
	return float64(int(toFloat(a)) % int(den))
}

// Asignación en map/array
func assignIndexExpression(idxExpr *IndexExpression, newVal interface{}, env *Environment) interface{} {
	leftVal := idxExpr.Left.Eval(env)
	indexVal := idxExpr.Index.Eval(env)

	switch container := leftVal.(type) {
	case map[string]interface{}:
		key, ok := indexVal.(string)
		if !ok {
			panic("assignIndexExpression: index for map must be a string")
		}
		container[key] = newVal
		return newVal
	case []interface{}:
		idxF, ok := indexVal.(float64)
		if !ok {
			panic("assignIndexExpression: array index must be a number")
		}
		idx := int(idxF)
		if idx < 0 {
			idx = len(container) + idx
		}
		// auto-extender
		if idx >= len(container) {
			for len(container) <= idx {
				container = append(container, nil)
			}
			// Need to update the variable that holds the array
			// since append might have created a new slice
			updateArrayInEnv(idxExpr.Left, container, env)
		}
		container[idx] = newVal
		return newVal
	default:
		panic("Not a map or array to assign index")
	}
}

// Helper function to update array in environment after extending
func updateArrayInEnv(node Node, newArray []interface{}, env *Environment) {
	switch n := node.(type) {
	case *Identifier:
		env.Update(n.Name, newArray)
	case *AccessExpression:
		objVal := n.Object.Eval(env)
		instance, ok := objVal.(*ObjectInstance)
		if ok {
			instance.Env.Set(n.Member, newArray)
		}
	case *IndexExpression:
		// Handle nested array access
		leftVal := n.Left.Eval(env)
		indexVal := n.Index.Eval(env)
		switch container := leftVal.(type) {
		case []interface{}:
			if idxF, ok := indexVal.(float64); ok {
				idx := int(idxF)
				if idx >= 0 && idx < len(container) {
					container[idx] = newArray
				}
			}
		}
	}
}

func isBinaryOp(op string) bool {
	ops := []string{"+", "-", "*", "/", "%", "<", ">", "<=", ">=", "==", "!=", "&&", "||", "&", "|", "^", "<<", ">>", "??", "|>"}
	for _, o := range ops {
		if op == o {
			return true
		}
	}
	return false
}

// toString convierte cualquier valor a string
func toString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case float64:
		// Si es un número entero, no mostrar decimales
		if v == float64(int64(v)) {
			return fmt.Sprintf("%.0f", v)
		}
		return fmt.Sprintf("%g", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}

// toStringOptimized convierte valores a string de manera optimizada para concatenación
func toStringOptimized(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case float64:
		// Cache común para números pequeños enteros
		if v >= 0 && v <= 100 && v == float64(int64(v)) {
			// Usar cache para números comunes
			return fmt.Sprintf("%.0f", v)
		}
		// Si es un número entero, no mostrar decimales
		if v == float64(int64(v)) {
			return fmt.Sprintf("%.0f", v)
		}
		return fmt.Sprintf("%g", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}

// smartParseFloat intelligently parses strings to float with enhanced conversion
func smartParseFloat(s string) (float64, error) {
	// Trim whitespace
	s = strings.TrimSpace(s)

	// Handle empty string
	if s == "" {
		return 0, nil
	}

	// Handle boolean-like strings
	switch strings.ToLower(s) {
	case "true", "yes", "on", "1":
		return 1, nil
	case "false", "no", "off", "0":
		return 0, nil
	}

	// Remove common non-numeric characters for smart conversion
	cleaned := strings.ReplaceAll(s, ",", "")      // Remove commas from numbers like "1,000"
	cleaned = strings.ReplaceAll(cleaned, "$", "") // Remove currency symbols
	cleaned = strings.ReplaceAll(cleaned, " ", "") // Remove spaces

	// Handle percentage
	if strings.HasSuffix(cleaned, "%") {
		cleaned = cleaned[:len(cleaned)-1]
		if f, err := strconv.ParseFloat(cleaned, 64); err == nil {
			return f / 100, nil
		}
	}

	// Try standard parsing
	return strconv.ParseFloat(cleaned, 64)
}

// smartAddValues performs intelligent addition with enhanced type coercion
func smartAddValues(a, b interface{}) interface{} {
	// Fast path: both already float64
	if af, ok := a.(float64); ok {
		if bf, ok := b.(float64); ok {
			return af + bf
		}
	}

	// Smart numeric conversion
	if isSmartNumeric(a) && isSmartNumeric(b) {
		return smartToFloat(a) + smartToFloat(b)
	}

	// Array concatenation
	if aa, ok := a.([]interface{}); ok {
		if bb, ok := b.([]interface{}); ok {
			return append(aa, bb...)
		}
		return append(aa, b)
	}

	if ab, ok := b.([]interface{}); ok {
		return append([]interface{}{a}, ab...)
	}

	// String concatenation with smart conversion
	return smartStringConcat(a, b)
}

// isSmartNumeric checks if a value can be intelligently converted to a number
func isSmartNumeric(val interface{}) bool {
	switch v := val.(type) {
	case float64, int, bool:
		return true
	case string:
		_, err := smartParseFloat(v)
		return err == nil
	default:
		return false
	}
}

// smartToFloat converts values to float with intelligent parsing
func smartToFloat(val interface{}) float64 {
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
	case string:
		if f, err := smartParseFloat(v); err == nil {
			return f
		}
		// Fallback to 0 for non-numeric strings in arithmetic context
		return 0
	default:
		return 0
	}
}

// smartStringConcat intelligently concatenates values as strings
func smartStringConcat(a, b interface{}) interface{} {
	sa := smartToString(a)
	sb := smartToString(b)

	// Use optimized concatenation for larger strings
	if len(sa)+len(sb) > 32 {
		return OptimizedStringConcat(sa, sb)
	}
	return sa + sb
}

// smartToString converts values to string with intelligent formatting
func smartToString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case float64:
		// Smart number formatting: avoid unnecessary decimals
		if v == float64(int64(v)) {
			return fmt.Sprintf("%.0f", v)
		}
		return fmt.Sprintf("%g", v)
	case int:
		return strconv.Itoa(v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}

// shouldUseP5SmartConversion determines if P5 smart conversion should be applied
func shouldUseP5SmartConversion(a, b interface{}) bool {
	// Handle string + numeric or boolean + numeric cases
	if isSmartNumeric(a) && isSmartNumeric(b) {
		// Both can be converted to numbers

		// Special case: avoid converting simple single-digit strings that might be DSL tokens
		sa, aIsString := a.(string)
		sb, bIsString := b.(string)

		if aIsString && bIsString {
			// If both are simple single-digit strings, be conservative
			if len(sa) == 1 && len(sb) == 1 &&
				isNumericChar(sa[0]) && isNumericChar(sb[0]) {
				return false // Let it concatenate for DSL compatibility
			}
		}

		return true
	}

	return false
}

// isNumericChar checks if a byte is a digit
func isNumericChar(b byte) bool {
	return b >= '0' && b <= '9'
}

// isObviouslyNumericString checks if a string is obviously meant to be a number
func isObviouslyNumericString(s string) bool {
	if s == "" {
		return false
	}

	// Simple numeric strings (pure digits) are considered numeric for P5
	// unless they look like they could be identifiers or other non-numeric context
	trimmed := strings.TrimSpace(s)

	// Check for obvious numeric formatting first
	if strings.Contains(trimmed, ",") || // Has commas like "1,000"
		strings.HasPrefix(trimmed, "$") || // Currency like "$100"
		strings.HasSuffix(trimmed, "%") || // Percentage like "50%"
		strings.Contains(trimmed, ".") || // Decimal like "123.45"
		strings.ToLower(trimmed) == "true" || // Boolean strings
		strings.ToLower(trimmed) == "false" ||
		strings.ToLower(trimmed) == "yes" ||
		strings.ToLower(trimmed) == "no" {
		return true
	}

	// Simple numeric strings (pure digits or floating point) are also considered numeric
	// but with a length check to avoid treating long identifiers as numbers
	if len(trimmed) <= 10 { // Reasonable limit for numeric strings
		_, err := strconv.ParseFloat(trimmed, 64)
		return err == nil
	}

	return false
}
