package main

import (
	"flag"
	"fmt"

	"github.com/arturoeanton/go-r2lang/pkg/r2repl"
)

const version = "1.0.0"

func main() {
	var (
		helpFlag    = flag.Bool("help", false, "Show help information")
		versionFlag = flag.Bool("version", false, "Show version information")
		noOutput    = flag.Bool("no-output", false, "Disable output display")
		quiet       = flag.Bool("quiet", false, "Quiet mode - minimal output")
		debug       = flag.Bool("debug", false, "Enable debug mode")
		multiline   = flag.Bool("multiline", false, "Enable multiline input mode")
		history     = flag.String("history", "", "History file path (default: ~/.r2lang_history)")
		prompt      = flag.String("prompt", "", "Custom prompt string")
		maxHistory  = flag.Int("max-history", 1000, "Maximum number of history entries")
		syntax      = flag.Bool("syntax", true, "Enable syntax highlighting")
		completion  = flag.Bool("completion", true, "Enable auto-completion")
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

	// Configure REPL options
	outputFlag := !*noOutput

	if *debug {
		fmt.Printf("R2Lang REPL v%s - Debug Mode\n", version)
		fmt.Printf("Configuration:\n")
		fmt.Printf("  Output: %t\n", outputFlag)
		fmt.Printf("  Quiet: %t\n", *quiet)
		fmt.Printf("  Multiline: %t\n", *multiline)
		fmt.Printf("  Syntax Highlighting: %t\n", *syntax)
		fmt.Printf("  Auto-completion: %t\n", *completion)
		fmt.Printf("  Max History: %d\n", *maxHistory)
		if *history != "" {
			fmt.Printf("  History File: %s\n", *history)
		}
		if *prompt != "" {
			fmt.Printf("  Custom Prompt: %s\n", *prompt)
		}
		fmt.Println()
	}

	if !*quiet {
		fmt.Printf("R2Lang REPL v%s\n", version)
		fmt.Println("Type 'exit' or 'quit' to exit, 'help' for help")
		fmt.Println("=========================================")
		fmt.Println()
	}

	// Start REPL with configured options
	r2repl.Repl(outputFlag)
}

func showHelp() {
	fmt.Printf("R2Lang REPL v%s - Interactive R2Lang Shell\n\n", version)
	fmt.Println("USAGE:")
	fmt.Println("  r2repl [OPTIONS]")
	fmt.Println()
	fmt.Println("DESCRIPTION:")
	fmt.Println("  R2Lang REPL provides an interactive shell for R2Lang programming.")
	fmt.Println("  You can execute R2Lang code line by line, define functions,")
	fmt.Println("  variables, and experiment with the language features.")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  -help                    Show this help message")
	fmt.Println("  -version                 Show version information")
	fmt.Println("  -no-output              Disable output display")
	fmt.Println("  -quiet                  Quiet mode - minimal startup output")
	fmt.Println("  -debug                  Enable debug mode with configuration info")
	fmt.Println("  -multiline              Enable multiline input mode")
	fmt.Println("  -history FILE           History file path (default: ~/.r2lang_history)")
	fmt.Println("  -prompt STRING          Custom prompt string")
	fmt.Println("  -max-history N          Maximum number of history entries (default: 1000)")
	fmt.Println("  -syntax                 Enable syntax highlighting (default: true)")
	fmt.Println("  -completion             Enable auto-completion (default: true)")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  r2repl                          # Start REPL with default settings")
	fmt.Println("  r2repl -quiet                   # Start REPL in quiet mode")
	fmt.Println("  r2repl -no-output              # Start REPL without output display")
	fmt.Println("  r2repl -debug                   # Start REPL with debug information")
	fmt.Println("  r2repl -prompt \"R2> \"           # Start REPL with custom prompt")
	fmt.Println("  r2repl -multiline               # Start REPL with multiline support")
	fmt.Println()
	fmt.Println("REPL COMMANDS:")
	fmt.Println("  help                    Show REPL help")
	fmt.Println("  exit, quit              Exit the REPL")
	fmt.Println("  clear                   Clear the screen")
	fmt.Println("  history                 Show command history")
	fmt.Println("  reset                   Reset the REPL environment")
	fmt.Println("  load FILE               Load and execute R2Lang file")
	fmt.Println("  save FILE               Save current session to file")
	fmt.Println()
	fmt.Println("LANGUAGE FEATURES:")
	fmt.Println("  Variables:              let x = 10;")
	fmt.Println("  Functions:              func add(a, b) { return a + b; }")
	fmt.Println("  Classes:                class Person { ... }")
	fmt.Println("  Objects:                let obj = {name: \"John\", age: 30};")
	fmt.Println("  Arrays:                 let arr = [1, 2, 3, 4, 5];")
	fmt.Println("  Maps:                   let map = {\"key\": \"value\"};")
	fmt.Println("  Control Flow:           if, while, for, try/catch")
	fmt.Println("  Concurrency:            r2() for goroutines")
	fmt.Println("  Built-in Libraries:     String, Math, I/O, HTTP, etc.")
	fmt.Println()
	fmt.Println("SHORTCUTS:")
	fmt.Println("  Ctrl+C                  Interrupt current input")
	fmt.Println("  Ctrl+D                  Exit REPL")
	fmt.Println("  Up/Down Arrow           Navigate command history")
	fmt.Println("  Tab                     Auto-completion")
	fmt.Println("  Ctrl+L                  Clear screen")
	fmt.Println()
	fmt.Println("For more information, visit: https://github.com/arturoeanton/go-r2lang")
}

func showVersion() {
	fmt.Printf("R2Lang REPL v%s\n", version)
	fmt.Println("Interactive R2Lang Shell")
	fmt.Println("Part of R2Lang - A Dynamic Programming Language")
	fmt.Println("Copyright (c) 2025 R2Lang Contributors")
	fmt.Println("Licensed under Apache License 2.0")
	fmt.Println()
	fmt.Println("Features:")
	fmt.Println("  - Interactive code execution")
	fmt.Println("  - Syntax highlighting")
	fmt.Println("  - Auto-completion")
	fmt.Println("  - Command history")
	fmt.Println("  - Multiline input support")
	fmt.Println("  - Built-in help system")
}
