package main

import (
	"fmt"
	"os"

	"github.com/arturoeanton/go-r2lang/pkg/r2lang"
	"github.com/arturoeanton/go-r2lang/pkg/r2repl"
	"github.com/arturoeanton/go-r2lang/pkg/r2test"
	"github.com/arturoeanton/go-r2lang/pkg/r2test/core"
)

func main() {
	filename := ""
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		if cmd == "-repl" {
			outputFlag := true
			if len(os.Args) > 2 {
				if os.Args[2] == "-no-output" {
					outputFlag = false
				}
			}
			r2repl.Repl(outputFlag)
			os.Exit(0)
		}
		if cmd == "-test" {
			runTests()
			os.Exit(0)
		}
		if cmd == "-help" || cmd == "--help" {
			showHelp()
			os.Exit(0)
		}
		filename = cmd

	} else {
		// intentar main.r2
		if _, err := os.Stat("main.r2"); os.IsNotExist(err) {
			fmt.Println("Error: Debes pasar un archivo .r2 o tener main.r2 en el directorio actual.")
			os.Exit(1)
		}

		filename = "main.r2"
	}
	r2lang.RunCode(filename)
}

func runTests() {
	fmt.Println("🧪 R2Lang Test Runner")
	fmt.Println("=====================")
	fmt.Println()

	// Crear configuración por defecto
	config := core.DefaultConfig()

	// Ajustar configuración para el directorio actual
	config.TestDirs = []string{"./examples/testing", "./tests"}
	config.Verbose = true

	// Crear y ejecutar el runner de tests
	rt := r2test.New(config)

	results, err := rt.RunDiscoveredTests()
	if err != nil {
		fmt.Printf("❌ Error ejecutando tests: %v\n", err)
		os.Exit(1)
	}

	stats := results.GetStats()

	fmt.Println()
	fmt.Println("📊 Resumen de Resultados:")
	fmt.Printf("   Total: %d tests\n", stats.Total)
	if stats.Passed > 0 {
		fmt.Printf("   ✅ Pasaron: %d\n", stats.Passed)
	}
	if stats.Failed > 0 {
		fmt.Printf("   ❌ Fallaron: %d\n", stats.Failed)
	}
	if stats.Skipped > 0 {
		fmt.Printf("   ⏭️  Saltados: %d\n", stats.Skipped)
	}
	if stats.Timeout > 0 {
		fmt.Printf("   ⏱️  Timeout: %d\n", stats.Timeout)
	}

	fmt.Printf("   ⏱️  Duración: %v\n", results.Duration)
	fmt.Println()

	if stats.Failed > 0 || stats.Timeout > 0 {
		fmt.Println("❌ Algunos tests fallaron")
		os.Exit(1)
	} else if stats.Total == 0 {
		fmt.Println("⚠️  No se encontraron tests")
		fmt.Println("   Coloca archivos *_test.r2 en ./examples/testing/ o ./tests/")
		os.Exit(0)
	} else {
		fmt.Println("🎉 ¡Todos los tests pasaron exitosamente!")
		os.Exit(0)
	}
}

func showHelp() {
	fmt.Println("R2Lang - Lenguaje de Programación Dinámico")
	fmt.Println("===========================================")
	fmt.Println()
	fmt.Println("Uso:")
	fmt.Println("  go run main.go [archivo.r2]     Ejecutar un archivo R2Lang")
	fmt.Println("  go run main.go -repl            Iniciar modo REPL interactivo")
	fmt.Println("  go run main.go -repl -no-output Iniciar REPL sin salida")
	fmt.Println("  go run main.go -test            Ejecutar tests unitarios")
	fmt.Println("  go run main.go -help            Mostrar esta ayuda")
	fmt.Println()
	fmt.Println("Ejemplos:")
	fmt.Println("  go run main.go gold_test.r2")
	fmt.Println("  go run main.go examples/example1-if.r2")
	fmt.Println("  go run main.go -test")
	fmt.Println()
	fmt.Println("Framework de Testing:")
	fmt.Println("  Los archivos de test deben:")
	fmt.Println("  - Terminar en *_test.r2")
	fmt.Println("  - Estar en ./examples/testing/ o ./tests/")
	fmt.Println("  - Usar sintaxis describe() e it()")
	fmt.Println()
	fmt.Println("Para más información, ver: https://github.com/arturoeanton/go-r2lang")
}
