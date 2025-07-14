package r2core

import (
	"fmt"
	"strconv"
	"sync"
)

var (
	// Cache para conversiones de string a float
	stringToFloatCache map[string]float64
	commonsStringCacheMu sync.RWMutex

	// Cache para números pequeños comunes
	intToFloatCache map[int]float64
)

func init() {
	// Inicializar caches
	stringToFloatCache = make(map[string]float64)
	intToFloatCache = make(map[int]float64)
	
	// Pre-poblar cache con números comunes
	for i := -1000; i <= 1000; i++ {
		intToFloatCache[i] = float64(i)
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
		// Buscar en cache primero
		commonsStringCacheMu.RLock()
		if cached, ok := stringToFloatCache[v]; ok {
			commonsStringCacheMu.RUnlock()
			return cached
		}
		commonsStringCacheMu.RUnlock()

		// Parsear y cachear
		f, err := strconv.ParseFloat(v, 64)
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

	if isNumeric(a) && isNumeric(b) {
		result := toFloat(a) + toFloat(b)
		// Usar object pool para números pequeños frecuentemente utilizados
		if IsSmallInteger(result) {
			return GetFloat64(result)
		}
		return result
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

	// Si uno es string => concatenar (optimizado)
	if sa, ok := a.(string); ok {
		sb := fmt.Sprint(b)
		return OptimizedStringConcat2(sa, sb)
	}
	if sb, ok := b.(string); ok {
		sa := fmt.Sprint(a)
		return OptimizedStringConcat2(sa, sb)
	}
	return toFloat(a) + toFloat(b)
}
func subValues(a, b interface{}) interface{} {
	result := toFloat(a) - toFloat(b)
	// Usar object pool para números pequeños frecuentemente utilizados
	if IsSmallInteger(result) {
		return GetFloat64(result)
	}
	return result
}
func mulValues(a, b interface{}) interface{} {
	result := toFloat(a) * toFloat(b)
	// Usar object pool para números pequeños frecuentemente utilizados
	if IsSmallInteger(result) {
		return GetFloat64(result)
	}
	return result
}
func divValues(a, b interface{}) interface{} {
	den := toFloat(b)
	if den == 0 {
		panic("Division by zero")
	}
	result := toFloat(a) / den
	// Usar object pool para números pequeños frecuentemente utilizados
	if IsSmallInteger(result) {
		return GetFloat64(result)
	}
	return result
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
		}
		container[idx] = newVal
		return newVal
	default:
		panic("Not a map or array to assign index")
	}
}

func isBinaryOp(op string) bool {
	ops := []string{"+", "-", "*", "/", "<", ">", "<=", ">=", "==", "!="}
	for _, o := range ops {
		if op == o {
			return true
		}
	}
	return false
}
