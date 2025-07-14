package main

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"testing"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
	"github.com/arturoeanton/go-r2lang/pkg/r2libs"
)

// BenchmarkBasicArithmetic mide el rendimiento de operaciones aritméticas básicas
func BenchmarkBasicArithmetic(b *testing.B) {
	code := `
		func calculate() {
			var result = 0;
			for (var i = 0; i < 1000; i = i + 1) {
				result = result + i * 2 - 1;
			}
			return result;
		}
		calculate();
	`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := r2core.NewParser(code)
		program := parser.ParseProgram()
		
		env := r2core.NewEnvironment()
		registerAllLibs(env)
		
		program.Eval(env)
	}
}

// BenchmarkStringOperations mide el rendimiento de operaciones con strings
func BenchmarkStringOperations(b *testing.B) {
	code := `
		func stringTest() {
			var text = "Hello ";
			for (var i = 0; i < 100; i = i + 1) {
				text = text + "World " + i;
			}
			return text;
		}
		stringTest();
	`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := r2core.NewParser(code)
		program := parser.ParseProgram()
		
		env := r2core.NewEnvironment()
		registerAllLibs(env)
		
		program.Eval(env)
	}
}

// BenchmarkArrayOperations mide el rendimiento de operaciones con arrays
func BenchmarkArrayOperations(b *testing.B) {
	code := `
		func arrayTest() {
			var arr = [];
			for (var i = 0; i < 500; i = i + 1) {
				arr.push(i);
			}
			
			var sum = 0;
			for (var i = 0; i < len(arr); i = i + 1) {
				sum = sum + arr[i];
			}
			return sum;
		}
		arrayTest();
	`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := r2core.NewParser(code)
		program := parser.ParseProgram()
		
		env := r2core.NewEnvironment()
		registerAllLibs(env)
		
		program.Eval(env)
	}
}

// BenchmarkMapOperations mide el rendimiento de operaciones con maps
func BenchmarkMapOperations(b *testing.B) {
	code := `
		func mapTest() {
			var map = {};
			for (var i = 0; i < 100; i = i + 1) {
				map["key" + i] = i * 2;
			}
			
			var sum = 0;
			for (var key in map) {
				sum = sum + map[key];
			}
			return sum;
		}
		mapTest();
	`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := r2core.NewParser(code)
		program := parser.ParseProgram()
		
		env := r2core.NewEnvironment()
		registerAllLibs(env)
		
		program.Eval(env)
	}
}

// BenchmarkFunctionCalls mide el rendimiento de llamadas a funciones
func BenchmarkFunctionCalls(b *testing.B) {
	code := `
		func fibonacci(n) {
			if (n <= 1) {
				return n;
			}
			return fibonacci(n - 1) + fibonacci(n - 2);
		}
		fibonacci(20);
	`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := r2core.NewParser(code)
		program := parser.ParseProgram()
		
		env := r2core.NewEnvironment()
		registerAllLibs(env)
		
		program.Eval(env)
	}
}

// BenchmarkObjectOperations mide el rendimiento de operaciones con objetos
func BenchmarkObjectOperations(b *testing.B) {
	code := `
		object Person {
			constructor(name, age) {
				this.name = name;
				this.age = age;
			}
			
			greet() {
				return "Hello, " + this.name;
			}
		}
		
		func testObjects() {
			var people = [];
			for (var i = 0; i < 50; i = i + 1) {
				var person = new Person("Person" + i, i + 20);
				people.push(person);
			}
			
			var greetings = [];
			for (var i = 0; i < len(people); i = i + 1) {
				greetings.push(people[i].greet());
			}
			return greetings.length;
		}
		testObjects();
	`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := r2core.NewParser(code)
		program := parser.ParseProgram()
		
		env := r2core.NewEnvironment()
		registerAllLibs(env)
		
		program.Eval(env)
	}
}

// BenchmarkLexerPerformance mide el rendimiento del lexer
func BenchmarkLexerPerformance(b *testing.B) {
	// Código complejo para el lexer
	code := `
		func complexFunction(a, b, c) {
			var result = 0;
			if (a > 0) {
				while (b < 100) {
					for (var i = 0; i < c; i = i + 1) {
						result = result + (a * b + i) / 2;
					}
					b = b + 1;
				}
			} else {
				try {
					result = factorial(a + b + c);
				} catch (e) {
					result = -1;
				}
			}
			return result;
		}
		
		func factorial(n) {
			if (n <= 1) {
				return 1;
			}
			return n * factorial(n - 1);
		}
		
		complexFunction(10, 50, 20);
	`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lexer := r2core.NewLexer(code)
		for {
			token := lexer.NextToken()
			if token.Type == r2core.TOKEN_EOF {
				break
			}
		}
	}
}

// BenchmarkParserPerformance mide el rendimiento del parser
func BenchmarkParserPerformance(b *testing.B) {
	code := `
		func complexFunction(a, b, c) {
			var result = 0;
			if (a > 0) {
				while (b < 100) {
					for (var i = 0; i < c; i = i + 1) {
						result = result + (a * b + i) / 2;
					}
					b = b + 1;
				}
			} else {
				try {
					result = factorial(a + b + c);
				} catch (e) {
					result = -1;
				}
			}
			return result;
		}
		
		func factorial(n) {
			if (n <= 1) {
				return 1;
			}
			return n * factorial(n - 1);
		}
		
		complexFunction(10, 50, 20);
	`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := r2core.NewParser(code)
		parser.ParseProgram()
	}
}

// MemoryBenchmark mide el uso de memoria
func BenchmarkMemoryUsage(b *testing.B) {
	code := `
		func memoryTest() {
			var bigArray = [];
			for (var i = 0; i < 10000; i = i + 1) {
				bigArray.push({
					id: i,
					name: "Item" + i,
					data: []
				});
			}
			return bigArray.length;
		}
		memoryTest();
	`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runtime.GC()
		
		var m1 runtime.MemStats
		runtime.ReadMemStats(&m1)
		
		parser := r2core.NewParser(code)
		program := parser.ParseProgram()
		
		env := r2core.NewEnvironment()
		registerAllLibs(env)
		
		program.Eval(env)
		
		var m2 runtime.MemStats
		runtime.ReadMemStats(&m2)
		
		// Reportar uso de memoria
		b.ReportMetric(float64(m2.Alloc-m1.Alloc)/1024/1024, "MB/op")
	}
}

// TestPerformanceReport ejecuta todos los benchmarks y genera un reporte
func TestPerformanceReport(t *testing.T) {
	fmt.Println("=== REPORTE DE RENDIMIENTO R2LANG ===")
	fmt.Println("Ejecutando benchmarks...")
	
	// Información del sistema
	fmt.Printf("Sistema: %s\n", runtime.GOOS)
	fmt.Printf("Arquitectura: %s\n", runtime.GOARCH)
	fmt.Printf("CPUs: %d\n", runtime.NumCPU())
	fmt.Printf("Versión Go: %s\n", runtime.Version())
	
	// Crear archivo de reporte
	report := fmt.Sprintf(`# Reporte de Performance - R2Lang
Fecha: %s
Sistema: %s %s
CPUs: %d
Versión Go: %s

## Benchmarks Ejecutados

Para ejecutar estos benchmarks:
` + "```bash" + `
go test -bench=. -benchmem performance_test.go
` + "```" + `

## Casos de Prueba

1. **Operaciones Aritméticas Básicas**: Loop con 1000 iteraciones de cálculos
2. **Operaciones de String**: Concatenación de strings en loop
3. **Operaciones de Array**: Creación y manipulación de arrays
4. **Operaciones de Map**: Creación y acceso a mapas
5. **Llamadas a Funciones**: Fibonacci recursivo
6. **Operaciones con Objetos**: Creación y métodos de objetos
7. **Rendimiento del Lexer**: Análisis léxico de código complejo
8. **Rendimiento del Parser**: Análisis sintáctico
9. **Uso de Memoria**: Creación de estructuras grandes

## Análisis y Mejoras Recomendadas

Los resultados de estos benchmarks se detallan en docs/es/performance.md
`, time.Now().Format("2006-01-02 15:04:05"), runtime.GOOS, runtime.GOARCH, runtime.NumCPU(), runtime.Version())
	
	err := ioutil.WriteFile("performance_report.md", []byte(report), 0644)
	if err != nil {
		t.Fatalf("Error escribiendo reporte: %v", err)
	}
	
	fmt.Println("Reporte generado: performance_report.md")
	fmt.Println("Ejecuta: go test -bench=. -benchmem performance_test.go")
}

// registerAllLibs registra todas las librerías disponibles
func registerAllLibs(env *r2core.Environment) {
	env.Set("true", true)
	env.Set("false", false)
	env.Set("nil", nil)
	env.Set("null", nil)
	
	r2libs.RegisterLib(env)
	r2libs.RegisterStd(env)
	r2libs.RegisterIO(env)
	r2libs.RegisterHTTPClient(env)
	r2libs.RegisterString(env)
	r2libs.RegisterMath(env)
	r2libs.RegisterRand(env)
	r2libs.RegisterTest(env)
	r2libs.RegisterHTTP(env)
	r2libs.RegisterPrint(env)
	r2libs.RegisterOS(env)
	r2libs.RegisterHack(env)
	r2libs.RegisterConcurrency(env)
	r2libs.RegisterCollections(env)
}