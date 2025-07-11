# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

R2Lang is a custom programming language interpreter written in Go. It features a JavaScript-like syntax with support for functions, classes, objects, arrays, maps, concurrency primitives, and built-in BDD testing capabilities.

## Key Commands

### Running R2Lang Programs
- `go run main.go filename.r2` - Execute a specific R2 file
- `go run main.go` - Execute `main.r2` in current directory (if exists)
- `go run main.go -repl` - Start interactive REPL mode
- `go run main.go -repl -no-output` - Start REPL without output

### Building and Testing
- `go build` - Build the interpreter binary
- `go test ./pkg/...` - Run tests for all modules
- `go test ./pkg/r2core/` - Run tests for core interpreter
- `go test ./pkg/r2libs/` - Run tests for built-in libraries
- `go run main.go examples/example*.r2` - Run any example file

## Modular Architecture (Restructured)

R2Lang has been completely restructured from a monolithic design to a clean modular architecture using Go's `pkg/` pattern. This transformation eliminated the God Object anti-pattern and provides excellent separation of concerns.

### Module Structure

**main.go**: Entry point (delegating to pkg/r2lang/)

**pkg/r2core/**: Core interpreter components (2,590 LOC, 30 files)
- `lexer.go`: Tokenizes R2 source code (330 LOC)
- `parse.go`: Builds AST from tokens using recursive descent parsing (678 LOC)
- `environment.go`: Variable scoping and function storage (98 LOC)
- 27 specialized AST node files: Each language construct in separate file
- `commons.go`: Shared utilities and error handling

**pkg/r2libs/**: Built-in function libraries (3,701 LOC, 18 files)
- `r2http.go`: HTTP server functionality (410 LOC)
- `r2httpclient.go`: HTTP client requests (324 LOC)
- `r2string.go`: String manipulation (194 LOC)
- `r2math.go`: Mathematical functions (87 LOC)
- `r2io.go`: File I/O operations (194 LOC)
- Plus 13 additional specialized libraries

**pkg/r2lang/**: High-level coordinator (45 LOC)
- Orchestrates core components
- Registers all built-in libraries
- Manages program lifecycle

**pkg/r2repl/**: Interactive REPL (185 LOC)
- Independent REPL implementation
- Advanced features: syntax highlighting, history, multiline support

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

### Built-in Libraries (pkg/r2libs/)

Each library file is self-contained and registers functions with the environment:

- **r2hack.go**: Cryptographic and security utilities (509 LOC)
- **r2http.go**: HTTP server with routing (410 LOC)
- **r2print.go**: Advanced output formatting (365 LOC)
- **r2httpclient.go**: HTTP client requests (324 LOC)
- **r2os.go**: Operating system interface (245 LOC)
- **r2goroutine.r2.go**: Concurrency primitives (237 LOC)
- **r2io.go**: File I/O operations (194 LOC)
- **r2string.go**: String manipulation (194 LOC)
- **r2std.go**: Standard utilities (typeOf, len, sleep, parseInt) (122 LOC)
- **r2math.go**: Mathematical functions (87 LOC)
- Plus 8 additional specialized libraries (collections, test, rand, etc.)

### Error Handling

The modular architecture enables consistent error handling patterns across modules. Each module handles errors appropriate to its domain, with standardized error propagation between modules.

## Development Patterns (Updated for Modular Architecture)

### Adding New Language Features
1. **Lexer changes**: Update `pkg/r2core/lexer.go` to recognize new tokens
2. **Parser changes**: Add parsing logic to `pkg/r2core/parse.go` or create new AST node file
3. **AST implementation**: Create new file in `pkg/r2core/` implementing Node interface
4. **Evaluation**: Implement `Eval(env *Environment)` method in the AST node
5. **Testing**: Add comprehensive tests in corresponding `*_test.go` file

### Adding Built-in Functions
1. **Create library file**: Add new `pkg/r2libs/r2newlib.go`
2. **Implement Register function**: Create `RegisterNewLib(env *Environment)` function
3. **Register in coordinator**: Call registration function in `pkg/r2lang/r2lang.go`
4. **Add tests**: Create `pkg/r2libs/r2newlib_test.go` with comprehensive coverage

### Contributing to Existing Modules
1. **pkg/r2core/**: Focus on specific AST node files or parser improvements
2. **pkg/r2libs/**: Each library is independent - pick one to enhance
3. **pkg/r2repl/**: REPL features and interactive improvements
4. **Documentation**: Update module-specific README files

### Module File Structure
```
R2Lang/
├── pkg/r2core/         # Core interpreter (30 files, highly modular)
├── pkg/r2libs/         # Built-in libraries (18 libraries, each specialized)
├── pkg/r2lang/         # High-level coordinator (minimal, clean)
├── pkg/r2repl/         # Interactive REPL (independent)
├── examples/           # R2 language example programs
├── docs/              # Comprehensive documentation (Spanish & English)
├── vscode_syntax_highlighting/  # VS Code extension
├── r2-jetbrains-syntax/         # IDE syntax highlighting
└── main.r2            # Example program

Architecture Benefits:
- Each module is independently testable
- Clear separation of concerns
- Easy onboarding for new contributors
- Parallel development possible
- Excellent maintainability (8.5/10 score)
```

### Quality Metrics (Post-Restructuring)
- **Code Quality**: 8.5/10 (vs 6.2/10 before)
- **Maintainability**: 8.5/10 (vs 2/10 before) 
- **Testability**: 9/10 (vs 1/10 before)
- **Technical Debt**: 79% reduction achieved