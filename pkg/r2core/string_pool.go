package r2core

import (
	"strings"
	"sync"
)

// StringBuilderPool es un pool de StringBuilder para concatenación eficiente
type StringBuilderPool struct {
	pool sync.Pool
}

// StringBuilder encapsula un strings.Builder para reutilización
type StringBuilder struct {
	builder strings.Builder
}

var (
	// Pool global para string builders
	globalStringBuilderPool = &StringBuilderPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &StringBuilder{}
			},
		},
	}

	// Cache para concatenaciones frecuentes (hasta 512 bytes)
	stringConcatCache map[string]string
	poolStringCacheMu sync.RWMutex
	maxCacheSize      = 1000 // Limitar el tamaño del cache
	currentCacheSize  = 0
)

func init() {
	stringConcatCache = make(map[string]string)
}

// GetStringBuilder obtiene un StringBuilder desde el pool
func GetStringBuilder() *StringBuilder {
	sb := globalStringBuilderPool.pool.Get().(*StringBuilder)
	sb.builder.Reset() // Limpiar cualquier contenido previo
	return sb
}

// PutStringBuilder devuelve un StringBuilder al pool para reutilización
func PutStringBuilder(sb *StringBuilder) {
	if sb.builder.Len() < 64*1024 { // Solo reutilizar si no es muy grande
		globalStringBuilderPool.pool.Put(sb)
	}
}

// WriteString escribe una cadena al builder
func (sb *StringBuilder) WriteString(s string) {
	sb.builder.WriteString(s)
}

// String devuelve la cadena concatenada
func (sb *StringBuilder) String() string {
	return sb.builder.String()
}

// Len devuelve la longitud actual del builder
func (sb *StringBuilder) Len() int {
	return sb.builder.Len()
}

// OptimizedStringConcat concatena múltiples strings de manera eficiente
func OptimizedStringConcat(parts ...string) string {
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return parts[0]
	}
	if len(parts) == 2 {
		return OptimizedStringConcat2(parts[0], parts[1])
	}

	// Para más de 2 partes, usar StringBuilder del pool
	sb := GetStringBuilder()
	defer PutStringBuilder(sb)

	for _, part := range parts {
		sb.WriteString(part)
	}

	return sb.String()
}

// OptimizedStringConcat2 concatena dos strings de manera eficiente
func OptimizedStringConcat2(a, b string) string {
	// Para strings muy pequeños, usar concatenación directa
	if len(a)+len(b) < 32 {
		return a + b
	}

	// Para strings medianos, verificar cache
	if len(a)+len(b) < 512 {
		cacheKey := a + "\x00" + b // Usar null separator para evitar colisiones

		poolStringCacheMu.RLock()
		if result, exists := stringConcatCache[cacheKey]; exists {
			poolStringCacheMu.RUnlock()
			return result
		}
		poolStringCacheMu.RUnlock()

		result := a + b

		// Agregar al cache si no está lleno
		poolStringCacheMu.Lock()
		if currentCacheSize < maxCacheSize {
			stringConcatCache[cacheKey] = result
			currentCacheSize++
		}
		poolStringCacheMu.Unlock()

		return result
	}

	// Para strings grandes, usar StringBuilder del pool
	sb := GetStringBuilder()
	defer PutStringBuilder(sb)

	sb.WriteString(a)
	sb.WriteString(b)

	return sb.String()
}

// IsStringSmall verifica si una cadena es pequeña para concatenación directa
func IsStringSmall(s string) bool {
	return len(s) < 32
}

// ClearStringCache limpia el cache de strings (para testing)
func ClearStringCache() {
	poolStringCacheMu.Lock()
	defer poolStringCacheMu.Unlock()
	stringConcatCache = make(map[string]string)
	currentCacheSize = 0
}
