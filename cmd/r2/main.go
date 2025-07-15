package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/arturoeanton/go-r2lang/pkg/r2lang"
)

const version = "1.0.0"

func main() {
	var (
		helpFlag    = flag.Bool("help", false, "Show help information")
		versionFlag = flag.Bool("version", false, "Show version information")
		verbose     = flag.Bool("verbose", false, "Enable verbose output")
		debug       = flag.Bool("debug", false, "Enable debug mode")
		optimize    = flag.Bool("optimize", false, "Enable code optimization")
		profile     = flag.String("profile", "", "Enable profiling (cpu, memory, trace)")
		output      = flag.String("output", "", "Output file for compilation")
		workDir     = flag.String("workdir", "", "Working directory for execution")
		args        = flag.String("args", "", "Arguments to pass to R2Lang program")
		env         = flag.String("env", "", "Environment variables (key=value,key2=value2)")
		timeout     = flag.String("timeout", "", "Execution timeout (e.g., 30s, 5m)")
		maxMemory   = flag.String("max-memory", "", "Maximum memory usage (e.g., 100MB, 1GB)")
		interactive = flag.Bool("interactive", false, "Enable interactive mode")
		check       = flag.Bool("check", false, "Check syntax only, don't execute")
		format      = flag.Bool("format", false, "Format R2Lang code")
		compile     = flag.Bool("compile", false, "Compile to bytecode")
		bytecode    = flag.Bool("bytecode", false, "Execute bytecode file")
	)

	flag.Usage = func() {
		showHelp()
	}

	flag.Parse()

	if *helpFlag {
		showHelp()
		return
	}

	if *versionFlag {
		showVersion()
		return
	}

	// Get the filename from arguments
	argsv := flag.Args()
	filename := ""

	if len(argsv) > 0 {
		filename = argsv[0]
	} else {
		// Try to find main.r2 in current directory
		if _, err := os.Stat("main.r2"); os.IsNotExist(err) {
			fmt.Println("Error: You must provide a .r2 file or have main.r2 in the current directory.")
			fmt.Println("Use 'r2 -help' for usage information.")
			os.Exit(1)
		}
		filename = "main.r2"
	}

	// Validate file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' does not exist.\n", filename)
		os.Exit(1)
	}

	// Validate file extension
	if !strings.HasSuffix(filename, ".r2") {
		fmt.Printf("Error: File '%s' is not a .r2 file.\n", filename)
		os.Exit(1)
	}

	// Set working directory if specified
	if *workDir != "" {
		if err := os.Chdir(*workDir); err != nil {
			fmt.Printf("Error: Cannot change to working directory '%s': %v\n", *workDir, err)
			os.Exit(1)
		}
	}

	// Handle different modes
	if *check {
		if *verbose {
			fmt.Printf("Checking syntax of '%s'...\n", filename)
		}
		checkSyntax(filename)
		return
	}

	if *format {
		if *verbose {
			fmt.Printf("Formatting '%s'...\n", filename)
		}
		formatCode(filename)
		return
	}

	if *compile {
		if *verbose {
			fmt.Printf("Compiling '%s'...\n", filename)
		}
		compileCode(filename, *output)
		return
	}

	if *bytecode {
		if *verbose {
			fmt.Printf("Executing bytecode '%s'...\n", filename)
		}
		executeBytecode(filename)
		return
	}

	// Configure execution environment
	if *env != "" {
		envPairs := strings.Split(*env, ",")
		for _, pair := range envPairs {
			if kv := strings.Split(pair, "="); len(kv) == 2 {
				os.Setenv(kv[0], kv[1])
			}
		}
	}

	// Debug information
	if *debug {
		fmt.Printf("R2Lang v%s - Debug Mode\n", version)
		fmt.Printf("Configuration:\n")
		fmt.Printf("  File: %s\n", filename)
		fmt.Printf("  Verbose: %t\n", *verbose)
		fmt.Printf("  Optimize: %t\n", *optimize)
		fmt.Printf("  Interactive: %t\n", *interactive)
		if *profile != "" {
			fmt.Printf("  Profile: %s\n", *profile)
		}
		if *timeout != "" {
			fmt.Printf("  Timeout: %s\n", *timeout)
		}
		if *maxMemory != "" {
			fmt.Printf("  Max Memory: %s\n", *maxMemory)
		}
		if *workDir != "" {
			fmt.Printf("  Working Dir: %s\n", *workDir)
		}
		if *args != "" {
			fmt.Printf("  Arguments: %s\n", *args)
		}
		fmt.Println()
	}

	// Execute the R2Lang program
	if *verbose {
		fmt.Printf("Executing '%s'...\n", filename)
	}

	r2lang.RunCode(filename)
}

func checkSyntax(filename string) {
	// This would integrate with the parser to check syntax
	fmt.Printf("Syntax check for '%s' - Feature not yet implemented\n", filename)
	fmt.Println("This feature will be available in a future version.")
}

func formatCode(filename string) {
	// This would integrate with a code formatter
	fmt.Printf("Code formatting for '%s' - Feature not yet implemented\n", filename)
	fmt.Println("This feature will be available in a future version.")
}

func compileCode(filename, output string) {
	// This would compile R2Lang to bytecode
	outputFile := output
	if outputFile == "" {
		outputFile = strings.TrimSuffix(filename, filepath.Ext(filename)) + ".r2c"
	}

	fmt.Printf("Compiling '%s' to '%s' - Feature not yet implemented\n", filename, outputFile)
	fmt.Println("This feature will be available in a future version.")
}

func executeBytecode(filename string) {
	// This would execute compiled bytecode
	fmt.Printf("Executing bytecode '%s' - Feature not yet implemented\n", filename)
	fmt.Println("This feature will be available in a future version.")
}

func showHelp() {
	fmt.Printf("R2Lang v%s - Dynamic Programming Language\n\n", version)
	fmt.Println("USAGE:")
	fmt.Println("  r2 [OPTIONS] [FILE]")
	fmt.Println()
	fmt.Println("DESCRIPTION:")
	fmt.Println("  R2Lang is a dynamic programming language with JavaScript-like syntax.")
	fmt.Println("  It supports functions, classes, objects, arrays, concurrency, and more.")
	fmt.Println()
	fmt.Println("ARGUMENTS:")
	fmt.Println("  FILE                    R2Lang file to execute (default: main.r2)")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println()
	fmt.Println("Basic Options:")
	fmt.Println("  -help                   Show this help message")
	fmt.Println("  -version                Show version information")
	fmt.Println("  -verbose                Enable verbose output")
	fmt.Println("  -debug                  Enable debug mode")
	fmt.Println("  -interactive            Enable interactive mode")
	fmt.Println()
	fmt.Println("Execution Options:")
	fmt.Println("  -workdir DIR            Set working directory")
	fmt.Println("  -args STRING            Arguments to pass to R2Lang program")
	fmt.Println("  -env KEY=VALUE,...      Environment variables")
	fmt.Println("  -timeout DURATION       Execution timeout (e.g., 30s, 5m)")
	fmt.Println("  -max-memory SIZE        Maximum memory usage (e.g., 100MB, 1GB)")
	fmt.Println()
	fmt.Println("Code Processing:")
	fmt.Println("  -check                  Check syntax only, don't execute")
	fmt.Println("  -format                 Format R2Lang code")
	fmt.Println("  -optimize               Enable code optimization")
	fmt.Println("  -compile                Compile to bytecode")
	fmt.Println("  -output FILE            Output file for compilation")
	fmt.Println("  -bytecode               Execute bytecode file")
	fmt.Println()
	fmt.Println("Performance:")
	fmt.Println("  -profile TYPE           Enable profiling (cpu, memory, trace)")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  r2                              # Execute main.r2")
	fmt.Println("  r2 script.r2                    # Execute script.r2")
	fmt.Println("  r2 -verbose script.r2           # Execute with verbose output")
	fmt.Println("  r2 -debug script.r2             # Execute with debug information")
	fmt.Println("  r2 -check script.r2             # Check syntax only")
	fmt.Println("  r2 -format script.r2            # Format code")
	fmt.Println("  r2 -compile script.r2           # Compile to bytecode")
	fmt.Println("  r2 -compile -output app.r2c script.r2  # Compile with custom output")
	fmt.Println("  r2 -bytecode app.r2c            # Execute bytecode")
	fmt.Println("  r2 -timeout 30s script.r2       # Execute with timeout")
	fmt.Println("  r2 -optimize script.r2          # Execute with optimizations")
	fmt.Println("  r2 -profile cpu script.r2       # Execute with CPU profiling")
	fmt.Println()
	fmt.Println("LANGUAGE FEATURES:")
	fmt.Println("  Variables:              let x = 10;")
	fmt.Println("  Functions:              func add(a, b) { return a + b; }")
	fmt.Println("  Classes:                class Person extends Object { ... }")
	fmt.Println("  Objects:                let obj = {name: \"John\", age: 30};")
	fmt.Println("  Arrays:                 let arr = [1, 2, 3, 4, 5];")
	fmt.Println("  Maps:                   let map = {\"key\": \"value\"};")
	fmt.Println("  Control Flow:           if/else, while, for, try/catch/finally")
	fmt.Println("  Imports:                import \"module.r2\" as mod;")
	fmt.Println("  Concurrency:            r2() for goroutines")
	fmt.Println("  Testing:                Built-in describe()/it() test framework")
	fmt.Println()
	fmt.Println("BUILT-IN LIBRARIES:")
	fmt.Println("  String manipulation:    r2string library")
	fmt.Println("  Math functions:         r2math library")
	fmt.Println("  File I/O:              r2io library")
	fmt.Println("  HTTP server/client:     r2http, r2httpclient libraries")
	fmt.Println("  Operating system:       r2os library")
	fmt.Println("  Cryptography:           r2hack library")
	fmt.Println("  Date/time:             r2date library")
	fmt.Println("  Collections:           r2collections library")
	fmt.Println("  Unicode support:        r2unicode library")
	fmt.Println("  Concurrency:           r2goroutine library")
	fmt.Println()
	fmt.Println("FILE EXTENSIONS:")
	fmt.Println("  .r2                    R2Lang source files")
	fmt.Println("  .r2c                   R2Lang compiled bytecode (future)")
	fmt.Println("  *_test.r2              R2Lang test files")
	fmt.Println()
	fmt.Println("RELATED COMMANDS:")
	fmt.Println("  r2repl                 Start interactive REPL")
	fmt.Println("  r2test                 Run R2Lang tests")
	fmt.Println()
	fmt.Println("For more information, visit: https://github.com/arturoeanton/go-r2lang")
}

func showVersion() {
	fmt.Printf("R2Lang v%s\n", version)
	fmt.Println("Dynamic Programming Language with JavaScript-like Syntax")
	fmt.Println("Copyright (c) 2025 R2Lang Contributors")
	fmt.Println("Licensed under Apache License 2.0")
	fmt.Println()
	fmt.Println("Architecture: Modular Design")
	fmt.Println("Components:")
	fmt.Println("  - r2core: Core interpreter with lexer, parser, and AST")
	fmt.Println("  - r2libs: Built-in function libraries")
	fmt.Println("  - r2repl: Interactive REPL")
	fmt.Println("  - r2test: Testing framework")
	fmt.Println()
	fmt.Println("Features:")
	fmt.Println("  - JavaScript-like syntax")
	fmt.Println("  - Object-oriented programming")
	fmt.Println("  - Functional programming support")
	fmt.Println("  - Concurrency with goroutines")
	fmt.Println("  - Built-in testing framework")
	fmt.Println("  - Extensive standard library")
	fmt.Println("  - Import/module system")
	fmt.Println("  - Error handling with try/catch")
}
