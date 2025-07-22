package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/arturoeanton/go-r2lang/pkg/r2lang"
)

const version = "1.0.0"

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

	// Handle command line arguments
	if len(os.Args) > 1 {
		cmd := os.Args[1]

		switch cmd {
		case "-help", "--help", "help":
			showHelp()
			os.Exit(0)
		case "-version", "--version", "version":
			fmt.Printf("R2Lang version %s\n", version)
			os.Exit(0)
		default:
			// Run R2Lang file
			filename := cmd
			r2lang.RunCode(filename)
		}
	} else {
		// Try to run main.r2 in current directory
		if _, err := os.Stat("main.r2"); os.IsNotExist(err) {
			fmt.Println("Error: Provide a .r2 file or have main.r2 in current directory.")
			fmt.Println("Use 'r2lang --help' for more information.")
			os.Exit(1)
		}
		r2lang.RunCode("main.r2")
	}
}

func showHelp() {
	fmt.Println("R2Lang - Dynamic Programming Language")
	fmt.Println("=====================================")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  r2lang [file.r2]              Run a R2Lang file")
	fmt.Println("  r2lang --help                 Show this help")
	fmt.Println("  r2lang --version              Show version")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  r2lang main.r2                Run main.r2")
	fmt.Println("  r2lang examples/example1.r2   Run example file")
	fmt.Println()
	fmt.Println("Other Commands:")
	fmt.Println("  r2           Advanced CLI with flags (cmd/r2/)")
	fmt.Println("  r2repl       Interactive REPL (cmd/repl/)")
	fmt.Println("  r2test       Testing framework (cmd/r2test/)")
	fmt.Println()
	fmt.Println("Installation:")
	fmt.Println("  go install github.com/arturoeliasanton/go-r2lang/cmd/r2lang@latest")
	fmt.Println("  go install github.com/arturoeliasanton/go-r2lang/cmd/r2@latest")
	fmt.Println("  go install github.com/arturoeliasanton/go-r2lang/cmd/r2repl@latest")
	fmt.Println("  go install github.com/arturoeliasanton/go-r2lang/cmd/r2test@latest")
	fmt.Println()
	fmt.Println("Language Features:")
	fmt.Println("  - JavaScript-like syntax")
	fmt.Println("  - Object-oriented programming")
	fmt.Println("  - Functions and classes")
	fmt.Println("  - Concurrency with goroutines")
	fmt.Println("  - Integrated testing framework")
	fmt.Println("  - Extensive built-in libraries")
	fmt.Println("  - DSL builder capabilities")
	fmt.Println()
	fmt.Println("For more information: https://github.com/arturoeliasanton/go-r2lang")
}
