package r2libs

import (
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func RegisterMath(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		// Basic math functions
		"sin": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("sin needs (number)")
			}
			x := toFloat(args[0])
			return math.Sin(x)
		}),

		"cos": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("cos needs (number)")
			}
			x := toFloat(args[0])
			return math.Cos(x)
		}),

		"tan": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("tan needs (number)")
			}
			x := toFloat(args[0])
			return math.Tan(x)
		}),

		"asin": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("asin needs (number)")
			}
			x := toFloat(args[0])
			return math.Asin(x)
		}),

		"acos": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("acos needs (number)")
			}
			x := toFloat(args[0])
			return math.Acos(x)
		}),

		"atan": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("atan needs (number)")
			}
			x := toFloat(args[0])
			return math.Atan(x)
		}),

		"atan2": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("atan2 needs (y, x)")
			}
			y := toFloat(args[0])
			x := toFloat(args[1])
			return math.Atan2(y, x)
		}),

		"sinh": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("sinh needs (number)")
			}
			x := toFloat(args[0])
			return math.Sinh(x)
		}),

		"cosh": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("cosh needs (number)")
			}
			x := toFloat(args[0])
			return math.Cosh(x)
		}),

		"tanh": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("tanh needs (number)")
			}
			x := toFloat(args[0])
			return math.Tanh(x)
		}),

		"sqrt": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("sqrt needs (number)")
			}
			x := toFloat(args[0])
			if x < 0 {
				panic("sqrt: could not calculate square root of negative number")
			}
			return math.Sqrt(x)
		}),

		"cbrt": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("cbrt needs (number)")
			}
			x := toFloat(args[0])
			return math.Cbrt(x)
		}),

		"pow": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("pow needs (base, exp)")
			}
			base := toFloat(args[0])
			exp := toFloat(args[1])
			return math.Pow(base, exp)
		}),

		"log": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("log needs (number)")
			}
			x := toFloat(args[0])
			if x <= 0 {
				panic("log: could not calculate log of zero or negative number")
			}
			return math.Log(x)
		}),

		"log10": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("log10 needs (number)")
			}
			x := toFloat(args[0])
			if x <= 0 {
				panic("log10: could not calculate log of zero or negative number")
			}
			return math.Log10(x)
		}),

		"log2": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("log2 needs (number)")
			}
			x := toFloat(args[0])
			if x <= 0 {
				panic("log2: could not calculate log of zero or negative number")
			}
			return math.Log2(x)
		}),

		"exp": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("exp needs (number)")
			}
			x := toFloat(args[0])
			return math.Exp(x)
		}),

		"exp2": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("exp2 needs (number)")
			}
			x := toFloat(args[0])
			return math.Exp2(x)
		}),

		"abs": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("abs needs (number)")
			}
			x := toFloat(args[0])
			return math.Abs(x)
		}),

		"floor": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("floor needs (number)")
			}
			x := toFloat(args[0])
			return math.Floor(x)
		}),

		"ceil": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("ceil needs (number)")
			}
			x := toFloat(args[0])
			return math.Ceil(x)
		}),

		"round": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("round needs (number)")
			}
			x := toFloat(args[0])
			return math.Round(x)
		}),

		"trunc": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("trunc needs (number)")
			}
			x := toFloat(args[0])
			return math.Trunc(x)
		}),

		"mod": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("mod needs (x, y)")
			}
			x := toFloat(args[0])
			y := toFloat(args[1])
			return math.Mod(x, y)
		}),

		"remainder": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("remainder needs (x, y)")
			}
			x := toFloat(args[0])
			y := toFloat(args[1])
			return math.Remainder(x, y)
		}),

		"max": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("max needs at least one number")
			}
			if len(args) == 1 {
				if arr, ok := args[0].([]interface{}); ok {
					return arrayMax(arr)
				}
			}

			max := toFloat(args[0])
			for _, arg := range args[1:] {
				val := toFloat(arg)
				if val > max {
					max = val
				}
			}
			return max
		}),

		"min": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("min needs at least one number")
			}
			if len(args) == 1 {
				if arr, ok := args[0].([]interface{}); ok {
					return arrayMin(arr)
				}
			}

			min := toFloat(args[0])
			for _, arg := range args[1:] {
				val := toFloat(arg)
				if val < min {
					min = val
				}
			}
			return min
		}),

		"clamp": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 3 {
				panic("clamp needs (value, min, max)")
			}
			value := toFloat(args[0])
			min := toFloat(args[1])
			max := toFloat(args[2])

			if value < min {
				return min
			}
			if value > max {
				return max
			}
			return value
		}),

		"hypot": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("hypot needs (x, y)")
			}
			x := toFloat(args[0])
			y := toFloat(args[1])
			return math.Hypot(x, y)
		}),

		"sign": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("sign needs (number)")
			}
			x := toFloat(args[0])
			if x > 0 {
				return 1.0
			} else if x < 0 {
				return -1.0
			}
			return 0.0
		}),

		"isNaN": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("isNaN needs (number)")
			}
			x := toFloat(args[0])
			return math.IsNaN(x)
		}),

		"isInf": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("isInf needs (number)")
			}
			x := toFloat(args[0])
			return math.IsInf(x, 0)
		}),

		"isFinite": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("isFinite needs (number)")
			}
			x := toFloat(args[0])
			return !math.IsNaN(x) && !math.IsInf(x, 0)
		}),

		// Array statistics functions
		"sum": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("sum needs (array)")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("sum: first argument must be array")
			}
			return arraySum(arr)
		}),

		"mean": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("mean needs (array)")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("mean: first argument must be array")
			}
			return arrayMean(arr)
		}),

		"median": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("median needs (array)")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("median: first argument must be array")
			}
			return arrayMedian(arr)
		}),

		"mode": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("mode needs (array)")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("mode: first argument must be array")
			}
			return arrayMode(arr)
		}),

		"variance": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("variance needs (array)")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("variance: first argument must be array")
			}
			return arrayVariance(arr)
		}),

		"stdDev": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("stdDev needs (array)")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("stdDev: first argument must be array")
			}
			return math.Sqrt(arrayVariance(arr))
		}),

		"range": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("range needs (array)")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("range: first argument must be array")
			}
			return arrayMax(arr) - arrayMin(arr)
		}),

		"percentile": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("percentile needs (array, p)")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("percentile: first argument must be array")
			}
			p := toFloat(args[1])
			return arrayPercentile(arr, p)
		}),

		"quartile": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("quartile needs (array, q)")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("quartile: first argument must be array")
			}
			q := toFloat(args[1])
			return arrayPercentile(arr, q*0.25)
		}),

		"correlation": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("correlation needs (array1, array2)")
			}
			arr1, ok1 := args[0].([]interface{})
			arr2, ok2 := args[1].([]interface{})
			if !ok1 || !ok2 {
				panic("correlation: both arguments must be arrays")
			}
			return arrayCorrelation(arr1, arr2)
		}),

		"covariance": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("covariance needs (array1, array2)")
			}
			arr1, ok1 := args[0].([]interface{})
			arr2, ok2 := args[1].([]interface{})
			if !ok1 || !ok2 {
				panic("covariance: both arguments must be arrays")
			}
			return arrayCovariance(arr1, arr2)
		}),

		"zscore": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("zscore needs (array)")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("zscore: first argument must be array")
			}
			return arrayZScore(arr)
		}),

		"normalize": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("normalize needs (array)")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("normalize: first argument must be array")
			}
			return arrayNormalize(arr)
		}),

		// Random number generation
		"random": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			return rand.Float64()
		}),

		"randomInt": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("randomInt needs (max)")
			}
			max := int(toFloat(args[0]))
			return float64(rand.Intn(max))
		}),

		"randomRange": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("randomRange needs (min, max)")
			}
			min := toFloat(args[0])
			max := toFloat(args[1])
			return min + rand.Float64()*(max-min)
		}),

		"randomSample": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("randomSample needs (array, count)")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("randomSample: first argument must be array")
			}
			count := int(toFloat(args[1]))
			return arrayRandomSample(arr, count)
		}),

		"shuffle": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("shuffle needs (array)")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("shuffle: first argument must be array")
			}
			return arrayShuffle(arr)
		}),

		"seed": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("seed needs (number)")
			}
			seed := int64(toFloat(args[0]))
			rand.Seed(seed)
			return nil
		}),

		// Distance functions
		"distance": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 4 {
				panic("distance needs (x1, y1, x2, y2)")
			}
			x1 := toFloat(args[0])
			y1 := toFloat(args[1])
			x2 := toFloat(args[2])
			y2 := toFloat(args[3])
			return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
		}),

		"manhattanDistance": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 4 {
				panic("manhattanDistance needs (x1, y1, x2, y2)")
			}
			x1 := toFloat(args[0])
			y1 := toFloat(args[1])
			x2 := toFloat(args[2])
			y2 := toFloat(args[3])
			return math.Abs(x2-x1) + math.Abs(y2-y1)
		}),

		// Angle functions
		"radToDeg": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("radToDeg needs (radians)")
			}
			rad := toFloat(args[0])
			return rad * 180 / math.Pi
		}),

		"degToRad": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("degToRad needs (degrees)")
			}
			deg := toFloat(args[0])
			return deg * math.Pi / 180
		}),

		// Factorial and combinations
		"factorial": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("factorial needs (n)")
			}
			n := int(toFloat(args[0]))
			return float64(factorial(n))
		}),

		"combination": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("combination needs (n, k)")
			}
			n := int(toFloat(args[0]))
			k := int(toFloat(args[1]))
			return float64(combination(n, k))
		}),

		"permutation": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("permutation needs (n, k)")
			}
			n := int(toFloat(args[0]))
			k := int(toFloat(args[1]))
			return float64(permutation(n, k))
		}),

		// GCD and LCM
		"gcd": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("gcd needs (a, b)")
			}
			a := int(toFloat(args[0]))
			b := int(toFloat(args[1]))
			return float64(gcd(a, b))
		}),

		"lcm": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("lcm needs (a, b)")
			}
			a := int(toFloat(args[0]))
			b := int(toFloat(args[1]))
			return float64(lcm(a, b))
		}),

		// Linear interpolation
		"lerp": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 3 {
				panic("lerp needs (a, b, t)")
			}
			a := toFloat(args[0])
			b := toFloat(args[1])
			t := toFloat(args[2])
			return a + t*(b-a)
		}),

		// Map value from one range to another
		"map": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 5 {
				panic("map needs (value, inMin, inMax, outMin, outMax)")
			}
			value := toFloat(args[0])
			inMin := toFloat(args[1])
			inMax := toFloat(args[2])
			outMin := toFloat(args[3])
			outMax := toFloat(args[4])
			return (value-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
		}),
	}

	// Add constants to the functions map as values
	constants := map[string]interface{}{
		"PI":       math.Pi,
		"E":        math.E,
		"PHI":      (1.0 + math.Sqrt(5.0)) / 2.0, // Golden ratio
		"SQRT2":    math.Sqrt2,
		"SQRT_E":   math.SqrtE,
		"SQRT_PI":  math.SqrtPi,
		"SQRT_PHI": math.SqrtPhi,
		"LN2":      math.Ln2,
		"LN10":     math.Ln10,
		"LOG2E":    math.Log2E,
	}

	RegisterModule(env, "math", functions)

	// Register constants within the math module
	mathModuleObj, _ := env.Get("math")
	mathModule := mathModuleObj.(map[string]interface{})
	for name, value := range constants {
		mathModule[name] = value
	}

	// Register constants globally as well for backwards compatibility
	env.Set("PI", math.Pi)
	env.Set("E", math.E)
	env.Set("PHI", (1.0+math.Sqrt(5.0))/2.0) // Golden ratio
	env.Set("SQRT2", math.Sqrt2)
	env.Set("SQRT_E", math.SqrtE)
	env.Set("SQRT_PI", math.SqrtPi)
	env.Set("SQRT_PHI", math.SqrtPhi)
	env.Set("LN2", math.Ln2)
	env.Set("LN10", math.Ln10)
	env.Set("LOG2E", math.Log2E)
	env.Set("LOG10E", math.Log10E)

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())
}

// Helper functions for array operations
func arraySum(arr []interface{}) float64 {
	sum := 0.0
	for _, val := range arr {
		sum += toFloat(val)
	}
	return sum
}

func arrayMean(arr []interface{}) float64 {
	if len(arr) == 0 {
		return 0
	}
	return arraySum(arr) / float64(len(arr))
}

func arrayMin(arr []interface{}) float64 {
	if len(arr) == 0 {
		return 0
	}
	min := toFloat(arr[0])
	for _, val := range arr[1:] {
		if v := toFloat(val); v < min {
			min = v
		}
	}
	return min
}

func arrayMax(arr []interface{}) float64 {
	if len(arr) == 0 {
		return 0
	}
	max := toFloat(arr[0])
	for _, val := range arr[1:] {
		if v := toFloat(val); v > max {
			max = v
		}
	}
	return max
}

func arrayMedian(arr []interface{}) float64 {
	if len(arr) == 0 {
		return 0
	}

	vals := make([]float64, len(arr))
	for i, val := range arr {
		vals[i] = toFloat(val)
	}

	sort.Float64s(vals)

	n := len(vals)
	if n%2 == 0 {
		return (vals[n/2-1] + vals[n/2]) / 2
	}
	return vals[n/2]
}

func arrayMode(arr []interface{}) float64 {
	if len(arr) == 0 {
		return 0
	}

	counts := make(map[float64]int)
	for _, val := range arr {
		v := toFloat(val)
		counts[v]++
	}

	maxCount := 0
	var mode float64
	for val, count := range counts {
		if count > maxCount {
			maxCount = count
			mode = val
		}
	}

	return mode
}

func arrayVariance(arr []interface{}) float64 {
	if len(arr) == 0 {
		return 0
	}

	mean := arrayMean(arr)
	sum := 0.0

	for _, val := range arr {
		diff := toFloat(val) - mean
		sum += diff * diff
	}

	return sum / float64(len(arr))
}

func arrayPercentile(arr []interface{}, p float64) float64 {
	if len(arr) == 0 {
		return 0
	}

	vals := make([]float64, len(arr))
	for i, val := range arr {
		vals[i] = toFloat(val)
	}

	sort.Float64s(vals)

	index := p * float64(len(vals)-1)
	if index == float64(int(index)) {
		return vals[int(index)]
	}

	lower := vals[int(index)]
	upper := vals[int(index)+1]
	return lower + (upper-lower)*(index-float64(int(index)))
}

func arrayCorrelation(arr1, arr2 []interface{}) float64 {
	if len(arr1) != len(arr2) || len(arr1) == 0 {
		return 0
	}

	mean1 := arrayMean(arr1)
	mean2 := arrayMean(arr2)

	numerator := 0.0
	sum1 := 0.0
	sum2 := 0.0

	for i := 0; i < len(arr1); i++ {
		diff1 := toFloat(arr1[i]) - mean1
		diff2 := toFloat(arr2[i]) - mean2

		numerator += diff1 * diff2
		sum1 += diff1 * diff1
		sum2 += diff2 * diff2
	}

	denominator := math.Sqrt(sum1 * sum2)
	if denominator == 0 {
		return 0
	}

	return numerator / denominator
}

func arrayCovariance(arr1, arr2 []interface{}) float64 {
	if len(arr1) != len(arr2) || len(arr1) == 0 {
		return 0
	}

	mean1 := arrayMean(arr1)
	mean2 := arrayMean(arr2)

	sum := 0.0
	for i := 0; i < len(arr1); i++ {
		sum += (toFloat(arr1[i]) - mean1) * (toFloat(arr2[i]) - mean2)
	}

	return sum / float64(len(arr1))
}

func arrayZScore(arr []interface{}) []interface{} {
	if len(arr) == 0 {
		return []interface{}{}
	}

	mean := arrayMean(arr)
	stdDev := math.Sqrt(arrayVariance(arr))

	result := make([]interface{}, len(arr))
	for i, val := range arr {
		if stdDev == 0 {
			result[i] = 0.0
		} else {
			result[i] = (toFloat(val) - mean) / stdDev
		}
	}

	return result
}

func arrayNormalize(arr []interface{}) []interface{} {
	if len(arr) == 0 {
		return []interface{}{}
	}

	min := arrayMin(arr)
	max := arrayMax(arr)

	result := make([]interface{}, len(arr))
	for i, val := range arr {
		if max == min {
			result[i] = 0.0
		} else {
			result[i] = (toFloat(val) - min) / (max - min)
		}
	}

	return result
}

func arrayRandomSample(arr []interface{}, count int) []interface{} {
	if count >= len(arr) {
		return arrayShuffle(arr)
	}

	shuffled := arrayShuffle(arr)
	return shuffled[:count]
}

func arrayShuffle(arr []interface{}) []interface{} {
	result := make([]interface{}, len(arr))
	copy(result, arr)

	for i := len(result) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		result[i], result[j] = result[j], result[i]
	}

	return result
}

func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

func combination(n, k int) int {
	if k > n || k < 0 {
		return 0
	}
	return factorial(n) / (factorial(k) * factorial(n-k))
}

func permutation(n, k int) int {
	if k > n || k < 0 {
		return 0
	}
	return factorial(n) / factorial(n-k)
}

func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

func lcm(a, b int) int {
	return a * b / gcd(a, b)
}
