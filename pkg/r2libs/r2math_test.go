package r2libs

import (
	"math"
	"testing"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func TestMathFunctions(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterMath(env)

	mathModuleObj, ok := env.Get("math")
	if !ok {
		t.Fatal("math module not found")
	}
	mathModule := mathModuleObj.(map[string]interface{})

	sinFunc := mathModule["sin"].(r2core.BuiltinFunction)
	cosFunc := mathModule["cos"].(r2core.BuiltinFunction)
	tanFunc := mathModule["tan"].(r2core.BuiltinFunction)
	sqrtFunc := mathModule["sqrt"].(r2core.BuiltinFunction)
	powFunc := mathModule["pow"].(r2core.BuiltinFunction)
	logFunc := mathModule["log"].(r2core.BuiltinFunction)
	log10Func := mathModule["log10"].(r2core.BuiltinFunction)
	expFunc := mathModule["exp"].(r2core.BuiltinFunction)
	absFunc := mathModule["abs"].(r2core.BuiltinFunction)
	floorFunc := mathModule["floor"].(r2core.BuiltinFunction)
	ceilFunc := mathModule["ceil"].(r2core.BuiltinFunction)
	roundFunc := mathModule["round"].(r2core.BuiltinFunction)
	maxFunc := mathModule["max"].(r2core.BuiltinFunction)
	minFunc := mathModule["min"].(r2core.BuiltinFunction)
	hypotFunc := mathModule["hypot"].(r2core.BuiltinFunction)

	tests := []struct {
		name     string
		function r2core.BuiltinFunction
		args     []interface{}
		expected float64
	}{
		{
			"sin(PI/2)", sinFunc, []interface{}{math.Pi / 2}, 1.0,
		},
		{
			"cos(PI)", cosFunc, []interface{}{math.Pi}, -1.0,
		},
		{
			"tan(PI/4)", tanFunc, []interface{}{math.Pi / 4}, 1.0,
		},
		{
			"sqrt(9)", sqrtFunc, []interface{}{9.0}, 3.0,
		},
		{
			"pow(2, 3)", powFunc, []interface{}{2.0, 3.0}, 8.0,
		},
		{
			"log(E)", logFunc, []interface{}{math.E}, 1.0,
		},
		{
			"log10(100)", log10Func, []interface{}{100.0}, 2.0,
		},
		{
			"exp(1)", expFunc, []interface{}{1.0}, math.E,
		},
		{
			"abs(-5)", absFunc, []interface{}{-5.0}, 5.0,
		},
		{
			"floor(3.7)", floorFunc, []interface{}{3.7}, 3.0,
		},
		{
			"ceil(3.2)", ceilFunc, []interface{}{3.2}, 4.0,
		},
		{
			"round(3.7)", roundFunc, []interface{}{3.7}, 4.0,
		},
		{
			"max(5, 10)", maxFunc, []interface{}{5.0, 10.0}, 10.0,
		},
		{
			"min(5, 10)", minFunc, []interface{}{5.0, 10.0}, 5.0,
		},
		{
			"hypot(3, 4)", hypotFunc, []interface{}{3.0, 4.0}, 5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function(tt.args...)
			if resultFloat, ok := result.(float64); ok {
				if math.Abs(resultFloat-tt.expected) > 1e-9 {
					t.Errorf("expected %v, got %v", tt.expected, resultFloat)
				}
			} else {
				t.Errorf("expected float64, got %T", result)
			}
		})
	}
}

func TestMathExtendedFunctions(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterMath(env)

	mathModuleObj, ok := env.Get("math")
	if !ok {
		t.Fatal("math module not found")
	}
	mathModule := mathModuleObj.(map[string]interface{})

	roundToFunc := mathModule["roundTo"].(r2core.BuiltinFunction)
	isPrimeFunc := mathModule["isPrime"].(r2core.BuiltinFunction)
	isEvenFunc := mathModule["isEven"].(r2core.BuiltinFunction)
	isOddFunc := mathModule["isOdd"].(r2core.BuiltinFunction)
	nthRootFunc := mathModule["nthRoot"].(r2core.BuiltinFunction)

	t.Run("roundTo rounds to given decimals", func(t *testing.T) {
		result := roundToFunc(3.14159, 2.0)
		if result != 3.14 {
			t.Errorf("expected 3.14, got %v", result)
		}
	})

	t.Run("roundTo zero decimals behaves like round", func(t *testing.T) {
		result := roundToFunc(3.6, 0.0)
		if result != 4.0 {
			t.Errorf("expected 4, got %v", result)
		}
	})

	primeTests := []struct {
		n        float64
		expected bool
	}{
		{2, true},
		{17, true},
		{1, false},
		{0, false},
		{-7, false},
		{15, false},
	}
	for _, pt := range primeTests {
		t.Run("isPrime", func(t *testing.T) {
			result := isPrimeFunc(pt.n)
			if result != pt.expected {
				t.Errorf("isPrime(%v): expected %v, got %v", pt.n, pt.expected, result)
			}
		})
	}

	t.Run("isEven true", func(t *testing.T) {
		if isEvenFunc(4.0) != true {
			t.Errorf("expected true")
		}
	})

	t.Run("isEven false", func(t *testing.T) {
		if isEvenFunc(3.0) != false {
			t.Errorf("expected false")
		}
	})

	t.Run("isOdd true", func(t *testing.T) {
		if isOddFunc(3.0) != true {
			t.Errorf("expected true")
		}
	})

	t.Run("isOdd false", func(t *testing.T) {
		if isOddFunc(4.0) != false {
			t.Errorf("expected false")
		}
	})

	t.Run("nthRoot basic", func(t *testing.T) {
		result := nthRootFunc(27.0, 3.0).(float64)
		if math.Abs(result-3.0) > 1e-9 {
			t.Errorf("expected 3, got %v", result)
		}
	})

	t.Run("nthRoot negative odd root", func(t *testing.T) {
		result := nthRootFunc(-8.0, 3.0).(float64)
		if math.Abs(result-(-2.0)) > 1e-9 {
			t.Errorf("expected -2, got %v", result)
		}
	})

	t.Run("nthRoot negative even root panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic")
			}
		}()
		nthRootFunc(-4.0, 2.0)
	})
}

func TestMathConstants(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterMath(env)

	piVal, ok := env.Get("PI")
	if !ok {
		t.Fatal("PI constant not found")
	}
	if pi, ok := piVal.(float64); !ok || math.Abs(pi-math.Pi) > 1e-9 {
		t.Errorf("Expected PI to be %v, got %v", math.Pi, pi)
	}

	eVal, ok := env.Get("E")
	if !ok {
		t.Fatal("E constant not found")
	}
	if e, ok := eVal.(float64); !ok || math.Abs(e-math.E) > 1e-9 {
		t.Errorf("Expected E to be %v, got %v", math.E, e)
	}
}

func TestMathErrorHandling(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterMath(env)

	mathModuleObj, ok := env.Get("math")
	if !ok {
		t.Fatal("math module not found")
	}
	mathModule := mathModuleObj.(map[string]interface{})

	sqrtFunc := mathModule["sqrt"].(r2core.BuiltinFunction)
	logFunc := mathModule["log"].(r2core.BuiltinFunction)
	log10Func := mathModule["log10"].(r2core.BuiltinFunction)

	// Test sqrt with negative number
	t.Run("sqrt(-1) panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for sqrt(-1)")
			}
		}()
		sqrtFunc(-1.0)
	})

	// Test log with zero
	t.Run("log(0) panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for log(0)")
			}
		}()
		logFunc(0.0)
	})

	// Test log10 with negative number
	t.Run("log10(-1) panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for log10(-1)")
			}
		}()
		log10Func(-1.0)
	})
}
