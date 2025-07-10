# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

R2Lang is a custom programming language interpreter written in Go. It features a JavaScript-like syntax with support for functions, classes, objects, arrays, maps, and concurrency primitives.

## Key Commands

### Running R2Lang Programs
- `go run main.go filename.r2` - Execute a specific R2 file
- `go run main.go` - Execute `main.r2` in current directory (if exists)
- `go run main.go -repl` - Start interactive REPL mode
- `go run main.go -repl -no-output` - Start REPL without output

### Building and Testing
- `go build` - Build the interpreter binary
- `go test ./r2lang` - Run tests for the language implementation
- `go run main.go examples/example*.r2` - Run any example file

## Core Architecture

### Main Components

**main.go**: Entry point that handles command-line arguments and delegates to the interpreter.

**r2lang/r2lang.go**: Core interpreter containing:
- Lexer: Tokenizes R2 source code, handles numbers, strings, operators, keywords
- Parser: Builds AST from tokens using recursive descent parsing  
- AST Nodes: All language constructs implement the `Node` interface with `Eval(env *Environment)` method
- Environment: Variable scoping and function storage system
- Evaluation: Tree-walking interpreter that executes AST nodes

### Language Features

**Data Types**: Numbers (float64), strings, booleans, arrays, maps, objects, functions

**Control Flow**: if/else, while, for loops (including for-in iteration), try/catch/finally

**Objects & Classes**: 
- Class inheritance with `extends` keyword
- Method overriding and `super` calls
- Constructor functions
- Object instantiation

**Functions**: 
- Named functions with `func` keyword
- Anonymous functions/lambdas
- Method calls on objects
- Built-in function libraries

**Concurrency**: 
- `r2()` function for goroutines
- Built-in synchronization primitives

**Testing**: Built-in test case syntax with Given/When/Then/And steps

**Imports**: Module system with `import "file.r2" as alias` syntax

### Built-in Libraries

Located in `r2lang/r2*.go` files, each registers functions with the environment:

- **r2lib.go**: Core functions (print, goroutines)  
- **r2std.go**: Standard utilities (typeOf, len, sleep, parseInt)
- **r2string.go**: String manipulation
- **r2math.go**: Mathematical functions
- **r2io.go**: File I/O operations  
- **r2http.go**: HTTP server functionality
- **r2httpclient.go**: HTTP client requests
- **r2test.go**: Testing utilities
- **r2print.go**: Output formatting
- **r2os.go**: Operating system interface
- **r2collections.go**: Array/map operations
- **r2rand.go**: Random number generation

### Error Handling

The interpreter uses Go's panic/recover mechanism. Parser errors show line/column information. Runtime errors display the current function context.

## Development Patterns

### Adding New Language Features
1. Add token types in r2lang.go constants
2. Update lexer to recognize new syntax
3. Create AST node struct implementing Node interface  
4. Add parsing logic in parser methods
5. Implement evaluation in the node's Eval method

### Adding Built-in Functions
1. Create new r2*.go file in r2lang/ directory
2. Implement Register* function that adds BuiltinFunction entries to environment
3. Call the register function in RunCode() in r2lang.go

### File Structure
- `examples/`: R2 language example programs demonstrating features
- `r2-jetbrains-syntax/`: Syntax highlighting for IDEs  
- `r2lang/`: Core interpreter implementation
- `main.r2`: Example program showing class inheritance