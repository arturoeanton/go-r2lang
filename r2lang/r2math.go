package r2lang

import (
	"math"
)

// r2math.go: Funciones matemáticas para R2

func RegisterMath(env *Environment) {
	// sin(x) => float64
	env.Set("sin", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("sin necesita (number)")
		}
		x := toFloat(args[0])
		return math.Sin(x)
	}))

	// cos(x)
	env.Set("cos", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("cos necesita (number)")
		}
		x := toFloat(args[0])
		return math.Cos(x)
	}))

	// tan(x)
	env.Set("tan", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("tan necesita (number)")
		}
		x := toFloat(args[0])
		return math.Tan(x)
	}))

	// sqrt(x)
	env.Set("sqrt", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("sqrt necesita (number)")
		}
		x := toFloat(args[0])
		if x < 0 {
			panic("sqrt: no se puede raíz de número negativo (sin complejos)")
		}
		return math.Sqrt(x)
	}))

	// pow(base, exp)
	env.Set("pow", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("pow necesita (base, exp)")
		}
		base := toFloat(args[0])
		exp := toFloat(args[1])
		return math.Pow(base, exp)
	}))

	// log(x) => log natural (base e)
	env.Set("log", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("log necesita (number)")
		}
		x := toFloat(args[0])
		if x <= 0 {
			panic("log: no se puede log de cero o negativo")
		}
		return math.Log(x)
	}))

	// log10(x) => log base 10
	env.Set("log10", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("log10 necesita (number)")
		}
		x := toFloat(args[0])
		if x <= 0 {
			panic("log10: no se puede log10 de cero o negativo")
		}
		return math.Log10(x)
	}))

	// exp(x) => e^x
	env.Set("exp", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("exp necesita (number)")
		}
		x := toFloat(args[0])
		return math.Exp(x)
	}))

	// abs(x) => valor absoluto
	env.Set("abs", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("abs necesita (number)")
		}
		x := toFloat(args[0])
		return math.Abs(x)
	}))

	// floor(x) => float64
	env.Set("floor", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("floor necesita (number)")
		}
		x := toFloat(args[0])
		return math.Floor(x)
	}))

	// ceil(x) => float64
	env.Set("ceil", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("ceil necesita (number)")
		}
		x := toFloat(args[0])
		return math.Ceil(x)
	}))

	// round(x) => float64
	env.Set("round", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("round necesita (number)")
		}
		x := toFloat(args[0])
		return math.Round(x)
	}))

	// max(a, b) => float64
	env.Set("max", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("max necesita (a, b)")
		}
		a := toFloat(args[0])
		b := toFloat(args[1])
		return math.Max(a, b)
	}))

	// min(a, b) => float64
	env.Set("min", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("min necesita (a, b)")
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
	env.Set("hypot", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("hypot necesita (x, y)")
		}
		x := toFloat(args[0])
		y := toFloat(args[1])
		return math.Hypot(x, y)
	}))
}
