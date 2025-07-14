package r2core

import (
	"sync"
)

// NumberPool es un pool de objetos para reutilizar números frecuentemente utilizados
type NumberPool struct {
	pool sync.Pool
}

// NumberWrapper encapsula un valor numérico para reutilización
type NumberWrapper struct {
	Value float64
}

var (
	// Pool global para números
	globalNumberPool = &NumberPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &NumberWrapper{}
			},
		},
	}
	
	// Cache para números pequeños frecuentemente utilizados
	smallNumberCache map[int]*NumberWrapper
	smallNumberMu    sync.RWMutex
)

func init() {
	// Inicializar cache
	smallNumberCache = make(map[int]*NumberWrapper)
	
	// Pre-poblar cache con números pequeños (-100 a 100)
	for i := -100; i <= 100; i++ {
		smallNumberCache[i] = &NumberWrapper{Value: float64(i)}
	}
}

// GetNumber obtiene un NumberWrapper desde el pool o cache
func GetNumber(value float64) *NumberWrapper {
	// Para números enteros pequeños, usar cache
	if value >= -100 && value <= 100 && value == float64(int(value)) {
		intVal := int(value)
		smallNumberMu.RLock()
		if wrapper, ok := smallNumberCache[intVal]; ok {
			smallNumberMu.RUnlock()
			return wrapper
		}
		smallNumberMu.RUnlock()
	}
	
	// Para otros números, usar pool
	wrapper := globalNumberPool.pool.Get().(*NumberWrapper)
	wrapper.Value = value
	return wrapper
}

// PutNumber devuelve un NumberWrapper al pool para reutilización
func PutNumber(wrapper *NumberWrapper) {
	// Solo devolver al pool si no es un número del cache
	if wrapper.Value < -100 || wrapper.Value > 100 || wrapper.Value != float64(int(wrapper.Value)) {
		globalNumberPool.pool.Put(wrapper)
	}
}

// GetFloat64 obtiene directamente el valor float64 optimizado
func GetFloat64(value float64) float64 {
	// Para números enteros pequeños, usar cache directo
	if value >= -100 && value <= 100 && value == float64(int(value)) {
		return value // Go optimiza esto automáticamente
	}
	return value
}

// IsSmallInteger verifica si un número es un entero pequeño
func IsSmallInteger(value float64) bool {
	return value >= -100 && value <= 100 && value == float64(int(value))
}