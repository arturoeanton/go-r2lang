package r2libs

import (
	"math/rand"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// r2rand.go: Funciones de generación aleatoria para R2

func RegisterRand(env *r2core.Environment) {
	// randInit([seed])
	// Si se pasa un argumento numérico => se usa como semilla,
	// si no => se usa time.Now().
	env.Set("randInit", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) == 0 {
			// seed con time
			rand.Seed(time.Now().UnixNano())
		} else {
			seed := toFloat(args[0])
			rand.Seed(int64(seed))
		}
		return nil
	}))

	// randInt(min, max) => random int en [min, max]
	env.Set("randInt", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("randInt needs (min, max)")
		}
		min := int(toFloat(args[0]))
		max := int(toFloat(args[1]))
		if max < min {
			panic("randInt: max < min")
		}
		return float64(rand.Intn(max-min+1) + min)
	}))

	// randFloat() => float64 en [0, 1)
	env.Set("randFloat", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		return rand.Float64()
	}))

	// randChoice(array) => elige un elemento al azar
	env.Set("randChoice", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("randChoice needs (array)")
		}
		arr, ok := args[0].([]interface{})
		if !ok {
			panic("randChoice: first arg should be native []array")
		}
		if len(arr) == 0 {
			panic("randChoice: array empty. No elements to choose from")
		}
		idx := rand.Intn(len(arr))
		return arr[idx]
	}))

	// shuffle(array) => baraja el array in-place
	env.Set("shuffle", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("shuffle needs (array)")
		}
		arr, ok := args[0].([]interface{})
		if !ok {
			panic("shuffle: first arg should be native []array")
		}
		n := len(arr)
		for i := n - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			arr[i], arr[j] = arr[j], arr[i]
		}
		return nil // barajado in-place
	}))

	// sample(array, n) => devuelve un nuevo array con n elementos aleatorios sin reemplazo
	env.Set("sample", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("sample needs (array, n)")
		}
		arr, ok := args[0].([]interface{})
		if !ok {
			panic("sample: first arg should be native []array")
		}
		n := int(toFloat(args[1]))
		if n < 0 {
			panic("sample: n < 0")
		}
		if n > len(arr) {
			// podría pánico, o ajustarlo
			panic("sample: n > length of array")
		}
		// clonamos arr para no modificar el original
		cloned := make([]interface{}, len(arr))
		copy(cloned, arr)
		// barajamos
		for i := len(cloned) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			cloned[i], cloned[j] = cloned[j], cloned[i]
		}
		// tomamos los primeros n
		result := cloned[:n]
		// convertimos a []interface{}
		res := make([]interface{}, n)
		for i := 0; i < n; i++ {
			res[i] = result[i]
		}
		return res
	}))
}
