package r2libs

import (
	"math/rand"
	"sync"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

var (
	randMu    sync.Mutex
	localRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func RegisterRand(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		"randInit": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			randMu.Lock()
			defer randMu.Unlock()
			if len(args) == 0 {
				localRand = rand.New(rand.NewSource(time.Now().UnixNano()))
			} else {
				seed := toFloat(args[0])
				localRand = rand.New(rand.NewSource(int64(seed)))
			}
			return nil
		}),

		"randInt": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("randInt needs (min, max)")
			}
			minF := toFloat(args[0])
			maxF := toFloat(args[1])
			if maxF < minF {
				panic("randInt: max < min")
			}
			// keep bounds well inside int64 so max-min+1 below can't wrap around and panic in Intn
			const randIntBound = 1e15
			if minF < -randIntBound || maxF > randIntBound {
				panic("randInt: min/max out of supported range")
			}
			min := int(minF)
			max := int(maxF)
			randMu.Lock()
			defer randMu.Unlock()
			return float64(localRand.Intn(max-min+1) + min)
		}),

		"randFloat": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			randMu.Lock()
			defer randMu.Unlock()
			return localRand.Float64()
		}),

		"randChoice": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("randChoice needs (array)")
			}
			arr, ok := toGenericSlice(args[0])
			if !ok {
				panic("randChoice: first arg should be native []array")
			}
			if len(arr) == 0 {
				panic("randChoice: array empty. No elements to choose from")
			}
			randMu.Lock()
			idx := localRand.Intn(len(arr))
			randMu.Unlock()
			return arr[idx]
		}),

		"shuffle": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("shuffle needs (array)")
			}
			arr, ok := toGenericSlice(args[0])
			if !ok {
				panic("shuffle: first arg should be native []array")
			}
			// Copy-then-shuffle, returning the new array, instead of
			// shuffling arr in place and returning nil: every other array
			// builtin in this codebase (.push/.filter/.sort/.reverse/...,
			// including math.shuffle, a second, independently-implemented
			// shuffle in this same language) is immutable/functional —
			// mutating in place here was a silent inconsistency that made
			// `let shuffled = rand.shuffle(arr)` assign nil while quietly
			// mutating the caller's original array instead.
			n := len(arr)
			result := make([]interface{}, n)
			copy(result, arr)
			randMu.Lock()
			defer randMu.Unlock()
			for i := n - 1; i > 0; i-- {
				j := localRand.Intn(i + 1)
				result[i], result[j] = result[j], result[i]
			}
			return result
		}),

		"sample": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("sample needs (array, n)")
			}
			arr, ok := toGenericSlice(args[0])
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
			randMu.Lock()
			for i := len(cloned) - 1; i > 0; i-- {
				j := localRand.Intn(i + 1)
				cloned[i], cloned[j] = cloned[j], cloned[i]
			}
			randMu.Unlock()
			return cloned[:n]
		}),
	}

	RegisterModule(env, "rand", functions)
}
