package r2libs

import (
	"fmt"
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

		// Advanced Data Analysis Functions
		"regression": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("regression needs (xArray, yArray)")
			}
			xArr, ok1 := args[0].([]interface{})
			yArr, ok2 := args[1].([]interface{})
			if !ok1 || !ok2 {
				panic("regression: both arguments must be arrays")
			}
			return linearRegression(xArr, yArr)
		}),

		"predict": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 3 {
				panic("predict needs (regressionResult, xValue, [order])")
			}
			regResult, ok1 := args[0].(map[string]interface{})
			xValue := toFloat(args[1])
			if !ok1 {
				panic("predict: first argument must be regression result")
			}

			order := 1 // linear by default
			if len(args) > 2 {
				order = int(toFloat(args[2]))
			}

			if slope, exists := regResult["slope"]; exists {
				if intercept, exists := regResult["intercept"]; exists {
					slopeVal := toFloat(slope)
					interceptVal := toFloat(intercept)
					if order == 1 {
						return slopeVal*xValue + interceptVal
					} else {
						// For higher order, use polynomial prediction (simplified)
						result := interceptVal
						for i := 1; i <= order; i++ {
							result += slopeVal * math.Pow(xValue, float64(i))
						}
						return result
					}
				}
			}
			return nil
		}),

		"movingAverage": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("movingAverage needs (array, windowSize)")
			}
			arr, ok := args[0].([]interface{})
			windowSize := int(toFloat(args[1]))
			if !ok {
				panic("movingAverage: first argument must be array")
			}
			return calculateMovingAverage(arr, windowSize)
		}),

		"exponentialSmoothing": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("exponentialSmoothing needs (array, alpha)")
			}
			arr, ok := args[0].([]interface{})
			alpha := toFloat(args[1])
			if !ok {
				panic("exponentialSmoothing: first argument must be array")
			}
			return exponentialSmoothing(arr, alpha)
		}),

		"differencing": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("differencing needs (array)")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("differencing: first argument must be array")
			}

			order := 1
			if len(args) > 1 {
				order = int(toFloat(args[1]))
			}

			return calculateDifferencing(arr, order)
		}),

		"autocorrelation": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("autocorrelation needs (array)")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("autocorrelation: first argument must be array")
			}

			maxLag := len(arr) - 1
			if len(args) > 1 {
				maxLag = int(toFloat(args[1]))
			}

			return calculateAutocorrelation(arr, maxLag)
		}),

		"seasonalDecompose": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("seasonalDecompose needs (array, period)")
			}
			arr, ok := args[0].([]interface{})
			period := int(toFloat(args[1]))
			if !ok {
				panic("seasonalDecompose: first argument must be array")
			}
			return seasonalDecomposition(arr, period)
		}),

		"outlierDetection": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("outlierDetection needs (array)")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("outlierDetection: first argument must be array")
			}

			method := "iqr"
			if len(args) > 1 {
				if m, ok := args[1].(string); ok {
					method = m
				}
			}

			return detectOutliers(arr, method)
		}),

		"histogram": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("histogram needs (array)")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("histogram: first argument must be array")
			}

			bins := 10
			if len(args) > 1 {
				bins = int(toFloat(args[1]))
			}

			return calculateHistogram(arr, bins)
		}),

		"frequency": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("frequency needs (array)")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("frequency: first argument must be array")
			}
			return calculateFrequency(arr)
		}),

		"cumulative": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("cumulative needs (array)")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("cumulative: first argument must be array")
			}
			return calculateCumulative(arr)
		}),

		"rollingStatistics": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 3 {
				panic("rollingStatistics needs (array, windowSize, statistic)")
			}
			arr, ok := args[0].([]interface{})
			windowSize := int(toFloat(args[1]))
			statistic, ok2 := args[2].(string)
			if !ok || !ok2 {
				panic("rollingStatistics: arguments must be (array, number, string)")
			}
			return calculateRollingStatistics(arr, windowSize, statistic)
		}),

		"trendAnalysis": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("trendAnalysis needs (array)")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("trendAnalysis: first argument must be array")
			}
			return analyzeTrend(arr)
		}),

		"dataQuality": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("dataQuality needs (array)")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				panic("dataQuality: first argument must be array")
			}
			return analyzeDataQuality(arr)
		}),

		"polynomialFit": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 3 {
				panic("polynomialFit needs (xArray, yArray, degree)")
			}
			xArr, ok1 := args[0].([]interface{})
			yArr, ok2 := args[1].([]interface{})
			degree := int(toFloat(args[2]))
			if !ok1 || !ok2 {
				panic("polynomialFit: x and y must be arrays")
			}
			return polynomialRegression(xArr, yArr, degree)
		}),

		"interpolate": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 3 {
				panic("interpolate needs (xArray, yArray, targetX)")
			}
			xArr, ok1 := args[0].([]interface{})
			yArr, ok2 := args[1].([]interface{})
			targetX := toFloat(args[2])
			if !ok1 || !ok2 {
				panic("interpolate: x and y must be arrays")
			}

			method := "linear"
			if len(args) > 3 {
				if m, ok := args[3].(string); ok {
					method = m
				}
			}

			return interpolateValue(xArr, yArr, targetX, method)
		}),

		"matrix": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("matrix needs (rows, cols)")
			}
			rows := int(toFloat(args[0]))
			cols := int(toFloat(args[1]))

			fillValue := 0.0
			if len(args) > 2 {
				fillValue = toFloat(args[2])
			}

			return createMatrix(rows, cols, fillValue)
		}),

		"matrixMultiply": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("matrixMultiply needs (matrixA, matrixB)")
			}
			matrixA, ok1 := args[0].([]interface{})
			matrixB, ok2 := args[1].([]interface{})
			if !ok1 || !ok2 {
				panic("matrixMultiply: both arguments must be matrices (arrays of arrays)")
			}
			return multiplyMatrices(matrixA, matrixB)
		}),

		"transpose": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("transpose needs (matrix)")
			}
			matrix, ok := args[0].([]interface{})
			if !ok {
				panic("transpose: argument must be matrix (array of arrays)")
			}
			return transposeMatrix(matrix)
		}),

		"determinant": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("determinant needs (matrix)")
			}
			matrix, ok := args[0].([]interface{})
			if !ok {
				panic("determinant: argument must be matrix (array of arrays)")
			}
			return calculateDeterminant(matrix)
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

// Advanced Data Analysis Helper Functions

func linearRegression(xArr, yArr []interface{}) map[string]interface{} {
	if len(xArr) != len(yArr) || len(xArr) < 2 {
		return map[string]interface{}{
			"slope":       0.0,
			"intercept":   0.0,
			"correlation": 0.0,
			"r_squared":   0.0,
			"error":       "Invalid data or insufficient points",
		}
	}

	n := float64(len(xArr))
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumXX := 0.0
	sumYY := 0.0

	for i := 0; i < len(xArr); i++ {
		x := toFloat(xArr[i])
		y := toFloat(yArr[i])
		sumX += x
		sumY += y
		sumXY += x * y
		sumXX += x * x
		sumYY += y * y
	}

	slope := (n*sumXY - sumX*sumY) / (n*sumXX - sumX*sumX)
	intercept := (sumY - slope*sumX) / n

	// Calculate correlation coefficient
	numerator := n*sumXY - sumX*sumY
	denominator := math.Sqrt((n*sumXX - sumX*sumX) * (n*sumYY - sumY*sumY))
	correlation := numerator / denominator
	rSquared := correlation * correlation

	return map[string]interface{}{
		"slope":       slope,
		"intercept":   intercept,
		"correlation": correlation,
		"r_squared":   rSquared,
		"error":       nil,
	}
}

func calculateMovingAverage(arr []interface{}, windowSize int) []interface{} {
	if windowSize <= 0 || windowSize > len(arr) {
		return []interface{}{}
	}

	var result []interface{}
	for i := windowSize - 1; i < len(arr); i++ {
		sum := 0.0
		for j := i - windowSize + 1; j <= i; j++ {
			sum += toFloat(arr[j])
		}
		result = append(result, sum/float64(windowSize))
	}

	return result
}

func exponentialSmoothing(arr []interface{}, alpha float64) []interface{} {
	if len(arr) == 0 || alpha <= 0 || alpha > 1 {
		return []interface{}{}
	}

	result := make([]interface{}, len(arr))
	result[0] = toFloat(arr[0])

	for i := 1; i < len(arr); i++ {
		current := toFloat(arr[i])
		previous := toFloat(result[i-1])
		result[i] = alpha*current + (1-alpha)*previous
	}

	return result
}

func calculateDifferencing(arr []interface{}, order int) []interface{} {
	if order <= 0 || len(arr) <= order {
		return []interface{}{}
	}

	current := make([]float64, len(arr))
	for i, v := range arr {
		current[i] = toFloat(v)
	}

	for d := 0; d < order; d++ {
		next := make([]float64, len(current)-1)
		for i := 1; i < len(current); i++ {
			next[i-1] = current[i] - current[i-1]
		}
		current = next
	}

	result := make([]interface{}, len(current))
	for i, v := range current {
		result[i] = v
	}

	return result
}

func calculateAutocorrelation(arr []interface{}, maxLag int) []interface{} {
	n := len(arr)
	if maxLag >= n {
		maxLag = n - 1
	}

	// Convert to float64 slice
	data := make([]float64, n)
	mean := 0.0
	for i, v := range arr {
		data[i] = toFloat(v)
		mean += data[i]
	}
	mean /= float64(n)

	// Center the data
	for i := range data {
		data[i] -= mean
	}

	// Calculate autocorrelations
	result := make([]interface{}, maxLag+1)

	// Variance (lag 0)
	variance := 0.0
	for _, v := range data {
		variance += v * v
	}

	for lag := 0; lag <= maxLag; lag++ {
		correlation := 0.0
		for i := 0; i < n-lag; i++ {
			correlation += data[i] * data[i+lag]
		}
		result[lag] = correlation / variance
	}

	return result
}

func seasonalDecomposition(arr []interface{}, period int) map[string]interface{} {
	n := len(arr)
	if period <= 0 || n < 2*period {
		return map[string]interface{}{
			"trend":    []interface{}{},
			"seasonal": []interface{}{},
			"residual": []interface{}{},
			"error":    "Invalid period or insufficient data",
		}
	}

	data := make([]float64, n)
	for i, v := range arr {
		data[i] = toFloat(v)
	}

	// Calculate trend using moving average
	trend := make([]float64, n)
	for i := 0; i < n; i++ {
		if i < period/2 || i >= n-period/2 {
			trend[i] = data[i] // Use original value for edges
		} else {
			sum := 0.0
			for j := i - period/2; j <= i+period/2; j++ {
				sum += data[j]
			}
			trend[i] = sum / float64(period+1)
		}
	}

	// Calculate seasonal component
	seasonal := make([]float64, n)
	seasonalSums := make([]float64, period)
	seasonalCounts := make([]int, period)

	for i := 0; i < n; i++ {
		detrended := data[i] - trend[i]
		seasonIndex := i % period
		seasonalSums[seasonIndex] += detrended
		seasonalCounts[seasonIndex]++
	}

	// Average seasonal components
	seasonalAvg := make([]float64, period)
	for i := 0; i < period; i++ {
		if seasonalCounts[i] > 0 {
			seasonalAvg[i] = seasonalSums[i] / float64(seasonalCounts[i])
		}
	}

	for i := 0; i < n; i++ {
		seasonal[i] = seasonalAvg[i%period]
	}

	// Calculate residual
	residual := make([]float64, n)
	for i := 0; i < n; i++ {
		residual[i] = data[i] - trend[i] - seasonal[i]
	}

	// Convert back to interface slices
	trendResult := make([]interface{}, n)
	seasonalResult := make([]interface{}, n)
	residualResult := make([]interface{}, n)

	for i := 0; i < n; i++ {
		trendResult[i] = trend[i]
		seasonalResult[i] = seasonal[i]
		residualResult[i] = residual[i]
	}

	return map[string]interface{}{
		"trend":    trendResult,
		"seasonal": seasonalResult,
		"residual": residualResult,
		"error":    nil,
	}
}

func detectOutliers(arr []interface{}, method string) map[string]interface{} {
	n := len(arr)
	if n < 4 {
		return map[string]interface{}{
			"outliers": []interface{}{},
			"indices":  []interface{}{},
			"bounds":   map[string]interface{}{},
		}
	}

	data := make([]float64, n)
	for i, v := range arr {
		data[i] = toFloat(v)
	}

	var outliers []interface{}
	var indices []interface{}
	var bounds map[string]interface{}

	switch method {
	case "iqr":
		// Sort data for percentile calculation
		sorted := make([]float64, n)
		copy(sorted, data)
		sort.Float64s(sorted)

		q1 := sorted[n/4]
		q3 := sorted[3*n/4]
		iqr := q3 - q1
		lowerBound := q1 - 1.5*iqr
		upperBound := q3 + 1.5*iqr

		bounds = map[string]interface{}{
			"lower": lowerBound,
			"upper": upperBound,
			"q1":    q1,
			"q3":    q3,
			"iqr":   iqr,
		}

		for i, v := range data {
			if v < lowerBound || v > upperBound {
				outliers = append(outliers, v)
				indices = append(indices, float64(i))
			}
		}

	case "zscore":
		mean := arrayMean(arr)
		stdDev := math.Sqrt(arrayVariance(arr))
		threshold := 3.0

		bounds = map[string]interface{}{
			"mean":      mean,
			"std_dev":   stdDev,
			"threshold": threshold,
		}

		for i, v := range data {
			zscore := math.Abs((v - mean) / stdDev)
			if zscore > threshold {
				outliers = append(outliers, v)
				indices = append(indices, float64(i))
			}
		}
	}

	return map[string]interface{}{
		"outliers": outliers,
		"indices":  indices,
		"bounds":   bounds,
		"method":   method,
	}
}

func calculateHistogram(arr []interface{}, bins int) map[string]interface{} {
	if bins <= 0 || len(arr) == 0 {
		return map[string]interface{}{
			"bins":   []interface{}{},
			"counts": []interface{}{},
			"edges":  []interface{}{},
		}
	}

	min := arrayMin(arr)
	max := arrayMax(arr)
	binWidth := (max - min) / float64(bins)

	counts := make([]int, bins)
	edges := make([]float64, bins+1)

	for i := 0; i <= bins; i++ {
		edges[i] = min + float64(i)*binWidth
	}

	for _, v := range arr {
		val := toFloat(v)
		binIndex := int((val - min) / binWidth)
		if binIndex >= bins {
			binIndex = bins - 1
		}
		if binIndex < 0 {
			binIndex = 0
		}
		counts[binIndex]++
	}

	// Convert to interface slices
	countsResult := make([]interface{}, bins)
	edgesResult := make([]interface{}, bins+1)
	binsResult := make([]interface{}, bins)

	for i := 0; i < bins; i++ {
		countsResult[i] = float64(counts[i])
		binsResult[i] = (edges[i] + edges[i+1]) / 2 // bin center
	}
	for i := 0; i <= bins; i++ {
		edgesResult[i] = edges[i]
	}

	return map[string]interface{}{
		"bins":   binsResult,
		"counts": countsResult,
		"edges":  edgesResult,
	}
}

func calculateFrequency(arr []interface{}) map[string]interface{} {
	frequency := make(map[string]int)
	total := len(arr)

	for _, v := range arr {
		key := fmt.Sprintf("%v", v)
		frequency[key]++
	}

	// Convert to result format
	values := make([]interface{}, 0, len(frequency))
	counts := make([]interface{}, 0, len(frequency))
	percentages := make([]interface{}, 0, len(frequency))

	for value, count := range frequency {
		values = append(values, value)
		counts = append(counts, float64(count))
		percentages = append(percentages, float64(count)/float64(total)*100)
	}

	return map[string]interface{}{
		"values":      values,
		"counts":      counts,
		"percentages": percentages,
		"total":       float64(total),
	}
}

func calculateCumulative(arr []interface{}) []interface{} {
	result := make([]interface{}, len(arr))
	cumSum := 0.0

	for i, v := range arr {
		cumSum += toFloat(v)
		result[i] = cumSum
	}

	return result
}

func calculateRollingStatistics(arr []interface{}, windowSize int, statistic string) []interface{} {
	if windowSize <= 0 || windowSize > len(arr) {
		return []interface{}{}
	}

	var result []interface{}
	for i := windowSize - 1; i < len(arr); i++ {
		window := arr[i-windowSize+1 : i+1]

		var value float64
		switch statistic {
		case "mean":
			value = arrayMean(window)
		case "sum":
			value = arraySum(window)
		case "min":
			value = arrayMin(window)
		case "max":
			value = arrayMax(window)
		case "std":
			value = math.Sqrt(arrayVariance(window))
		case "var":
			value = arrayVariance(window)
		default:
			value = arrayMean(window) // default to mean
		}

		result = append(result, value)
	}

	return result
}

func analyzeTrend(arr []interface{}) map[string]interface{} {
	n := len(arr)
	if n < 2 {
		return map[string]interface{}{
			"direction": "unknown",
			"slope":     0.0,
			"strength":  0.0,
		}
	}

	// Create x values (indices)
	xArr := make([]interface{}, n)
	for i := 0; i < n; i++ {
		xArr[i] = float64(i)
	}

	// Perform linear regression
	regression := linearRegression(xArr, arr)
	slope := toFloat(regression["slope"])
	correlation := toFloat(regression["correlation"])

	direction := "stable"
	if slope > 0.01 {
		direction = "increasing"
	} else if slope < -0.01 {
		direction = "decreasing"
	}

	strength := math.Abs(correlation)

	return map[string]interface{}{
		"direction":   direction,
		"slope":       slope,
		"strength":    strength,
		"correlation": correlation,
		"r_squared":   toFloat(regression["r_squared"]),
	}
}

func analyzeDataQuality(arr []interface{}) map[string]interface{} {
	n := len(arr)
	if n == 0 {
		return map[string]interface{}{
			"total_count":     0,
			"missing_count":   0,
			"missing_percent": 0.0,
			"unique_count":    0,
			"duplicates":      0,
		}
	}

	missing := 0
	unique := make(map[string]bool)
	total := float64(n)

	for _, v := range arr {
		if v == nil {
			missing++
		} else {
			key := fmt.Sprintf("%v", v)
			unique[key] = true
		}
	}

	uniqueCount := len(unique)
	duplicates := n - uniqueCount
	missingPercent := float64(missing) / total * 100

	return map[string]interface{}{
		"total_count":     total,
		"missing_count":   float64(missing),
		"missing_percent": missingPercent,
		"unique_count":    float64(uniqueCount),
		"duplicates":      float64(duplicates),
		"completeness":    100.0 - missingPercent,
	}
}

func polynomialRegression(xArr, yArr []interface{}, degree int) map[string]interface{} {
	// Simplified polynomial regression - for full implementation would need matrix operations
	if degree > 3 || len(xArr) < degree+1 {
		return map[string]interface{}{
			"coefficients": []interface{}{},
			"error":        "Degree too high or insufficient data",
		}
	}

	// For simplicity, return linear regression for degree 1
	if degree == 1 {
		return linearRegression(xArr, yArr)
	}

	// Simplified quadratic for degree 2
	return map[string]interface{}{
		"coefficients": []interface{}{0.0, 1.0, 0.0}, // placeholder
		"degree":       float64(degree),
		"error":        "Polynomial regression simplified implementation",
	}
}

func interpolateValue(xArr, yArr []interface{}, targetX float64, method string) interface{} {
	n := len(xArr)
	if n != len(yArr) || n < 2 {
		return nil
	}

	// Find surrounding points
	var x1, x2, y1, y2 float64
	found := false

	for i := 0; i < n-1; i++ {
		x1 = toFloat(xArr[i])
		x2 = toFloat(xArr[i+1])
		if targetX >= x1 && targetX <= x2 {
			y1 = toFloat(yArr[i])
			y2 = toFloat(yArr[i+1])
			found = true
			break
		}
	}

	if !found {
		return nil // Target X is outside the range
	}

	switch method {
	case "linear":
		if x2 == x1 {
			return y1
		}
		return y1 + (y2-y1)*(targetX-x1)/(x2-x1)
	case "nearest":
		if math.Abs(targetX-x1) < math.Abs(targetX-x2) {
			return y1
		}
		return y2
	default:
		// Default to linear
		if x2 == x1 {
			return y1
		}
		return y1 + (y2-y1)*(targetX-x1)/(x2-x1)
	}
}

func createMatrix(rows, cols int, fillValue float64) []interface{} {
	matrix := make([]interface{}, rows)
	for i := 0; i < rows; i++ {
		row := make([]interface{}, cols)
		for j := 0; j < cols; j++ {
			row[j] = fillValue
		}
		matrix[i] = row
	}
	return matrix
}

func multiplyMatrices(matrixA, matrixB []interface{}) interface{} {
	rowsA := len(matrixA)
	if rowsA == 0 {
		return nil
	}

	firstRowA, ok := matrixA[0].([]interface{})
	if !ok {
		return nil
	}
	colsA := len(firstRowA)

	rowsB := len(matrixB)
	if rowsB == 0 || rowsB != colsA {
		return nil
	}

	firstRowB, ok := matrixB[0].([]interface{})
	if !ok {
		return nil
	}
	colsB := len(firstRowB)

	// Create result matrix
	result := make([]interface{}, rowsA)
	for i := 0; i < rowsA; i++ {
		row := make([]interface{}, colsB)
		for j := 0; j < colsB; j++ {
			sum := 0.0
			for k := 0; k < colsA; k++ {
				rowA, ok1 := matrixA[i].([]interface{})
				rowB, ok2 := matrixB[k].([]interface{})
				if ok1 && ok2 && k < len(rowA) && j < len(rowB) {
					sum += toFloat(rowA[k]) * toFloat(rowB[j])
				}
			}
			row[j] = sum
		}
		result[i] = row
	}

	return result
}

func transposeMatrix(matrix []interface{}) []interface{} {
	rows := len(matrix)
	if rows == 0 {
		return []interface{}{}
	}

	firstRow, ok := matrix[0].([]interface{})
	if !ok {
		return []interface{}{}
	}
	cols := len(firstRow)

	result := make([]interface{}, cols)
	for j := 0; j < cols; j++ {
		row := make([]interface{}, rows)
		for i := 0; i < rows; i++ {
			if matrixRow, ok := matrix[i].([]interface{}); ok && j < len(matrixRow) {
				row[i] = matrixRow[j]
			} else {
				row[i] = 0.0
			}
		}
		result[j] = row
	}

	return result
}

func calculateDeterminant(matrix []interface{}) interface{} {
	n := len(matrix)
	if n == 0 {
		return 0.0
	}

	// Check if it's a square matrix
	for i := 0; i < n; i++ {
		if row, ok := matrix[i].([]interface{}); !ok || len(row) != n {
			return nil // Not a square matrix
		}
	}

	// Convert to float64 matrix for calculation
	floatMatrix := make([][]float64, n)
	for i := 0; i < n; i++ {
		floatMatrix[i] = make([]float64, n)
		row := matrix[i].([]interface{})
		for j := 0; j < n; j++ {
			floatMatrix[i][j] = toFloat(row[j])
		}
	}

	// Calculate determinant using LU decomposition (simplified for small matrices)
	if n == 1 {
		return floatMatrix[0][0]
	} else if n == 2 {
		return floatMatrix[0][0]*floatMatrix[1][1] - floatMatrix[0][1]*floatMatrix[1][0]
	} else {
		// For larger matrices, use a simplified approach
		det := 0.0
		for j := 0; j < n; j++ {
			// Create submatrix
			subMatrix := make([][]float64, n-1)
			for i := 1; i < n; i++ {
				subMatrix[i-1] = make([]float64, n-1)
				colIndex := 0
				for k := 0; k < n; k++ {
					if k != j {
						subMatrix[i-1][colIndex] = floatMatrix[i][k]
						colIndex++
					}
				}
			}

			// Convert submatrix back to interface{} for recursive call
			subMatrixInterface := make([]interface{}, n-1)
			for i := 0; i < n-1; i++ {
				row := make([]interface{}, n-1)
				for k := 0; k < n-1; k++ {
					row[k] = subMatrix[i][k]
				}
				subMatrixInterface[i] = row
			}

			subDet := toFloat(calculateDeterminant(subMatrixInterface))
			sign := 1.0
			if j%2 == 1 {
				sign = -1.0
			}
			det += sign * floatMatrix[0][j] * subDet
		}
		return det
	}
}
