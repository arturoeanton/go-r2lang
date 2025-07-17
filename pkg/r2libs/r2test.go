package r2libs

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// Estructura interna para guardar el resultado de cada test
type testResult struct {
	name    string
	passed  bool
	message string
	elapsed time.Duration
}

func RegisterTest(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		"assertEq": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("assertEq necesita (actual, expected, [msg])")
			}
			actual := args[0]
			expected := args[1]
			var msg string
			if len(args) >= 3 {
				if m, ok := args[2].(string); ok {
					msg = m
				}
			}
			if !reflect.DeepEqual(actual, expected) {
				if msg == "" {
					msg = fmt.Sprintf("assertEq fallo: actual=%v, expected=%v", actual, expected)
				} else {
					msg = fmt.Sprintf("%s (actual=%v, expected=%v)", msg, actual, expected)
				}
				panic(msg)
			}
			return nil
		}),

		"assertTrue": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("assertTrue necesita (cond, [msg])")
			}
			cond := toBool(args[0])
			var msg string
			if len(args) >= 2 {
				if m, ok := args[1].(string); ok {
					msg = m
				}
			}
			if !cond {
				if msg == "" {
					msg = "assertTrue fallo: la condición es falsa"
				}
				panic(msg)
			}
			return nil
		}),

		"runAllTests": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			// Recolectar nombres de funciones que empiecen con "test" (insensitive o no)
			var testNames []string
			for k, v := range env.GetStore() {
				// Podrías chequear en e.outer también si deseas
				if strings.HasPrefix(strings.ToLower(k), "test") {
					// Chequeamos si es *UserFunction
					if _, isFunc := v.(*r2core.UserFunction); isFunc {
						testNames = append(testNames, k)
					}
				}
			}
			if len(testNames) == 0 {
				fmt.Println("No se encontraron funciones test* en este script.")
				return nil
			}
			// Ordenar alfabéticamente, si quieres
			// sort.Strings(testNames)

			var results []testResult
			startGlobal := time.Now()

			for _, tName := range testNames {
				start := time.Now()
				testPassed := true
				testMsg := ""

				// Capturar panic
				func() {
					defer func() {
						if r := recover(); r != nil {
							testPassed = false
							testMsg = fmt.Sprint(r)
						}
					}()
					// Llamar la función
					val, _ := env.Get(tName) // ya sabemos que es *UserFunction
					fn, _ := val.(*r2core.UserFunction)
					fn.Call() // sin args
				}()
				elapsed := time.Since(start)
				results = append(results, testResult{
					name:    tName,
					passed:  testPassed,
					message: testMsg,
					elapsed: elapsed,
				})
			}
			endGlobal := time.Since(startGlobal)

			// Mostrar reporte
			passedCount := 0
			for _, r := range results {
				status := "PASSED"
				if !r.passed {
					status = "FAILED"
				}
				fmt.Printf("[%s] %s (%.2f ms)\n", status, r.name, float64(r.elapsed.Microseconds())/1000.0)
				if !r.passed {
					// indent el message
					lines := strings.Split(r.message, "\n")
					for _, ln := range lines {
						fmt.Printf("   %s\n", ln)
					}
				} else {
					passedCount++
				}
			}
			total := len(results)
			failedCount := total - passedCount
			fmt.Printf("\nResumen: %d PASSED, %d FAILED, %d TOTAL (%.2f ms)\n",
				passedCount, failedCount, total, float64(endGlobal.Microseconds())/1000.0)

			return nil
		}),

		"printStep": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			for _, arg := range args {
				fmt.Print(arg, " ")
			}
			fmt.Println()
			return nil
		}),

		"assertEqual": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) != 2 {
				panic("assertEqual requiere exactamente 2 argumentos")
			}
			a, b := args[0], args[1]
			if !equals(a, b) {
				panic(fmt.Sprintf("Assertion Failed: %v != %v", a, b))
			}
			return nil
		}),
	}

	RegisterModule(env, "test", functions)
}
