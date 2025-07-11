package r2libs

import (
	"math/rand"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

var localRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func RegisterRand(env *r2core.Environment) {
	// randInit([seed])
	env.Set("randInit", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) == 0 {
			localRand = rand.New(rand.NewSource(time.Now().UnixNano()))
		} else {
			seed := toFloat(args[0])
			localRand = rand.New(rand.NewSource(int64(seed)))
		}
		return nil
	}))

	// randInt(min, max)
	env.Set("randInt", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("randInt needs (min, max)")
		}
		min := int(toFloat(args[0]))
		max := int(toFloat(args[1]))
		if max < min {
			panic("randInt: max < min")
		}
		return float64(localRand.Intn(max-min+1) + min)
	}))

	// randFloat()
	env.Set("randFloat", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		return localRand.Float64()
	}))

	// randChoice(array)
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
		idx := localRand.Intn(len(arr))
		return arr[idx]
	}))

	// shuffle(array)
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
			j := localRand.Intn(i + 1)
			arr[i], arr[j] = arr[j], arr[i]
		}
		return nil
	}))

	// sample(array, n)
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
			panic("sample: n > length of array")
		}
		cloned := make([]interface{}, len(arr))
		copy(cloned, arr)
		for i := len(cloned) - 1; i > 0; i-- {
			j := localRand.Intn(i + 1)
			cloned[i], cloned[j] = cloned[j], cloned[i]
		}
		return cloned[:n]
	}))
}
