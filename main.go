package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/arturoeanton/go-r2lang/pkg/r2lang"
)

func main() {
	// Add error handling for R2Lang panics
	defer func() {
		if r := recover(); r != nil {
			errorStr := fmt.Sprintf("%v", r)
			// Check if this is an R2Lang error with call stack
			if strings.Contains(errorStr, "R2Lang call stack") {
				// Extract just the R2Lang error message and stack
				lines := strings.Split(errorStr, "\n")
				for _, line := range lines {
					if line == "" {
						break // Stop at first empty line (before Go stack trace)
					}
					fmt.Fprintln(os.Stderr, line)
				}
				os.Exit(1)
			} else {
				// For other errors, show full panic
				panic(r)
			}
		}
	}()

	filename := ""
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		if cmd == "-help" || cmd == "--help" {
			showHelp()
			os.Exit(0)
		}
		filename = cmd

	} else {
		// intentar main.r2
		if _, err := os.Stat("main.r2"); os.IsNotExist(err) {
			fmt.Println("Error: Debes pasar un archivo .r2 o tener main.r2 en el directorio actual.")
			fmt.Println("Use 'go run main.go -help' para más información.")
			os.Exit(1)
		}

		filename = "main.r2"
	}
	r2lang.RunCode(filename)
}

func showHelp() {
	fmt.Println("R2Lang - Lenguaje de Programación Dinámico")
	fmt.Println("===========================================")
	fmt.Println()
	fmt.Println("Uso:")
	fmt.Println("  go run main.go [archivo.r2]     Ejecutar un archivo R2Lang")
	fmt.Println("  go run main.go -help            Mostrar esta ayuda")
	fmt.Println()
	fmt.Println("Ejemplos:")
	fmt.Println("  go run main.go gold_test.r2")
	fmt.Println("  go run main.go examples/example1-if.r2")
	fmt.Println()
	fmt.Println("Comandos Especializados:")
	fmt.Println("  go run cmd/r2/main.go           Ejecutar con opciones avanzadas")
	fmt.Println("  go run cmd/repl/main.go         Iniciar REPL interactivo")
	fmt.Println("  go run cmd/r2test/main.go       Ejecutar framework de testing")
	fmt.Println()
	fmt.Println("Características del Lenguaje:")
	fmt.Println("  - Sintaxis similar a JavaScript")
	fmt.Println("  - Programación orientada a objetos")
	fmt.Println("  - Funciones y clases")
	fmt.Println("  - Concurrencia con goroutines")
	fmt.Println("  - Framework de testing integrado")
	fmt.Println("  - Librerías integradas extensas")
	fmt.Println()
	fmt.Println("Para más información, ver: https://github.com/arturoeanton/go-r2lang")
}
