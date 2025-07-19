# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Quick Commands

### Building and Running
```bash
# Main R2Lang interpreter
go run main.go script.r2
go build -o r2lang main.go

# Specialized commands
go build -o r2 cmd/r2/main.go        # Advanced CLI with flags
go build -o r2test cmd/r2test/main.go # Testing framework
go build -o r2repl cmd/repl/main.go   # Interactive REPL

# Quick test execution
go run main.go gold_test.r2           # Comprehensive language test
```

### Testing
```bash
# Go unit tests (416+ test cases)
go test ./pkg/r2core/
go test ./pkg/r2libs/
go test ./pkg/...

# R2Lang script testing
go run cmd/r2test/main.go ./examples/testing/
go run cmd/r2test/main.go -coverage -verbose ./tests
```

### REPL and Interactive Development
```bash
go run cmd/repl/main.go               # Start interactive REPL
go run main.go -repl                  # Alternative REPL entry
```

## Architecture Overview

### Core Components
- **pkg/r2core/**: Core interpreter (2,590 LOC) - lexer, parser, AST, bytecode
- **pkg/r2libs/**: Built-in libraries (3,701 LOC) - 18 specialized modules
- **pkg/r2lang/**: High-level coordinator (45 LOC) - main execution entry
- **pkg/r2repl/**: Interactive REPL (185 LOC)
- **pkg/r2test/**: Testing framework with coverage, mocking, and reporting

### Built-in Library Registration Pattern
All r2libs modules are auto-registered in `pkg/r2lang/r2lang.go`:
```go
r2libs.RegisterStd(env)     // Standard utilities
r2libs.RegisterIO(env)      // File I/O operations  
r2libs.RegisterOS(env)      // Operating system interface
r2libs.RegisterHTTP(env)    // HTTP server/client
r2libs.RegisterGRPC(env)    // Dynamic gRPC client
r2libs.RegisterJSON(env)    // JSON operations
// ... and 12 more modules
```

### Key Library Modules
- **r2io.go**: Comprehensive file I/O with Path objects and FileStream
- **r2os.go**: OS interface with Command objects and process management
- **r2grpc.go**: Dynamic gRPC client without code generation (1,467 LOC)
- **r2http.go**: HTTP server with routing (410 LOC)
- **r2hack.go**: Cryptographic utilities (509 LOC)
- **r2print.go**: Advanced output formatting (365 LOC)

## R2Lang Language Specifics

### Modern Syntax Features (2025)
- Boolean literals: `true`/`false` 
- Multiline maps: `{ key: value, another: value }`
- Template strings: `let msg = \`Hello ${name}\``
- Native dates: `let date = @2024-12-25`
- Modulo operator: `if (x % 2 == 0)`
- Enhanced else-if: `if (x) { } else if (y) { } else { }`

### Built-in Function Access
All r2libs functions are available globally without import:
```r2
// File operations
writeFile("test.txt", "content")
let content = readFile("test.txt")
let exists = exists("test.txt")

// OS operations  
let cmd = Command("echo hello")
cmd.run()
let output = cmd.stdout()

// HTTP operations
let response = request.get("https://api.example.com")
```

### DSL Builder (Unique Feature)
R2Lang includes a native DSL builder for creating custom parsers:
```r2
dsl Calculator {
    token("NUMBER", "[0-9]+")
    rule("operation", ["NUMBER", "+", "NUMBER"], "add")
    func add(left, op, right) { return parseFloat(left) + parseFloat(right) }
}
```

## Development Workflow

### File Structure Conventions
- `.r2` files: R2Lang source code
- `*_test.r2`: R2Lang test files  
- `*.r2c`: Compiled bytecode (future feature)
- Examples in `examples/` directory with numbered convention (`example1-if.r2`, etc.)

### Adding New Built-in Functions
1. Add functions to appropriate module in `pkg/r2libs/`
2. Register module in `pkg/r2lang/r2lang.go`
3. Functions become globally available in R2Lang scripts
4. Follow the `RegisterModule(env, "moduleName", functions)` pattern

### Testing Strategy
- **Go tests**: Core interpreter functionality (`pkg/r2core/`, `pkg/r2libs/`)
- **Gold test**: Comprehensive language validation (`gold_test.r2`)  
- **Example tests**: Individual feature demonstrations (`examples/`)
- **R2Test framework**: Advanced testing with coverage and mocking

### Error Handling
R2Lang uses panic-based error handling in built-in functions:
```go
if err != nil {
    panic(fmt.Sprintf("functionName: error description: %v", err))
}
```

## Performance and Quality
- Code Quality: 8.5/10 (improved from 6.2/10)
- Technical Debt Reduction: 79%
- 416+ test cases covering all features
- Modular architecture eliminates previous monolithic design

## Common Pitfalls
- Built-in functions use panic for errors, not Go error returns
- All r2libs modules are auto-imported; no manual import needed
- Use `print()` for output, not `println()` or `console.log()` (unless using console module)
- File paths in examples should be relative to project root
- DSL definitions require specific token/rule/function structure