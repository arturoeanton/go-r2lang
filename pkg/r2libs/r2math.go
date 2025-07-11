package r2libs

import (
	"math"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// r2math.go: Funciones matemáticas para R2

func RegisterMath(env *r2core.Environment) {
	// sin(x) => float64
	env.Set("sin", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("sin needs (number)")
		}
		x := toFloat(args[0])
		return math.Sin(x)
	}))

	// cos(x)
	env.Set("cos", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("cos  needs (number)")
		}
		x := toFloat(args[0])
		return math.Cos(x)
	}))

	// tan(x)
	env.Set("tan", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("tan needs (number)")
		}
		x := toFloat(args[0])
		return math.Tan(x)
	}))

	// sqrt(x)
	env.Set("sqrt", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("sqrt needs (number)")
		}
		x := toFloat(args[0])
		if x < 0 {
			panic("sqrt: could not calculate square root of negative number")
		}
		return math.Sqrt(x)
	}))

	// pow(base, exp)
	env.Set("pow", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("pow needs (base, exp)")
		}
		base := toFloat(args[0])
		exp := toFloat(args[1])
		return math.Pow(base, exp)
	}))

	// log(x) => log natural (base e)
	env.Set("log", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("log needs (number)")
		}
		x := toFloat(args[0])
		if x <= 0 {
			panic("log: could not calculate log of zero or negative number")
		}
		return math.Log(x)
	}))

	// log10(x) => log base 10
	env.Set("log10", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("log10 needs (number)")
		}
		x := toFloat(args[0])
		if x <= 0 {
			panic("log10: could not calculate log of zero or negative number")
		}
		return math.Log10(x)
	}))

	// exp(x) => e^x
	env.Set("exp", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("exp needs (number)")
		}
		x := toFloat(args[0])
		return math.Exp(x)
	}))

	// abs(x) => valor absoluto
	env.Set("abs", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("abs needs (number)")
		}
		x := toFloat(args[0])
		return math.Abs(x)
	}))

	// floor(x) => float64
	env.Set("floor", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("floor needs (number)")
		}
		x := toFloat(args[0])
		return math.Floor(x)
	}))

	// ceil(x) => float64
	env.Set("ceil", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("ceil needs (number)")
		}
		x := toFloat(args[0])
		return math.Ceil(x)
	}))

	// round(x) => float64
	env.Set("round", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("round needs (number)")
		}
		x := toFloat(args[0])
		return math.Round(x)
	}))

	// max(a, b) => float64
	env.Set("max", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("max needs (a, b)")
		}
		a := toFloat(args[0])
		b := toFloat(args[1])
		return math.Max(a, b)
	}))

	// min(a, b) => float64
	env.Set("min", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("min needs (a, b)")
		}
		a := toFloat(args[0])
		b := toFloat(args[1])
		return math.Min(a, b)
	}))

	// random pi
	// Podrías simplemente exponer pi
	env.Set("pi", float64(math.Pi))

	// podrías exponer E => math.E, etc.
	env.Set("e", float64(math.E))

	// Ejemplo: hypot(x, y) => sqrt(x^2 + y^2)
	env.Set("hypot", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("hypot needs (x, y)")
		}
		x := toFloat(args[0])
		y := toFloat(args[1])
		return math.Hypot(x, y)
	}))

}
