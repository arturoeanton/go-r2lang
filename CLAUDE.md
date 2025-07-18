# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

R2Lang is a comprehensive programming language interpreter written in Go featuring JavaScript-like syntax with advanced capabilities including functions, classes, objects, arrays, maps, concurrency primitives, built-in BDD testing, DSL support, and extensive library ecosystem.

## Key Commands

### Running R2Lang Programs
- `go run main.go filename.r2` - Execute a specific R2 file
- `go run main.go` - Execute `main.r2` in current directory (if exists)
- `go run main.go -repl` - Start interactive REPL mode
- `go run main.go -repl -no-output` - Start REPL without output
- `./cmd/r2/main.go filename.r2` - Alternative R2 executable
- `./cmd/r2test/main.go` - Dedicated test runner

### Building and Testing
- `go build` - Build the interpreter binary
- `go test ./pkg/...` - Run tests for all modules (122 Go files)
- `go test ./pkg/r2core/` - Run tests for core interpreter (58 files)
- `go test ./pkg/r2libs/` - Run tests for built-in libraries (42 files)
- `go run main.go examples/example*.r2` - Run any example file (60+ examples)

## Modular Architecture (Restructured)

R2Lang has been completely restructured from a monolithic design to a clean modular architecture using Go's `pkg/` pattern. This transformation eliminated the God Object anti-pattern and provides excellent separation of concerns.

### Module Structure

**main.go**: Entry point (delegating to pkg/r2lang/)

**pkg/r2core/**: Core interpreter components (11,964 LOC, 58 files)
- `lexer.go`: Tokenizes R2 source code with advanced features
- `parse.go`: Builds AST from tokens using recursive descent parsing
- `environment.go`: Variable scoping and function storage
- `dsl_definition.go`, `dsl_grammar.go`, `dsl_usage.go`: Domain-Specific Language support
- `bytecode.go`: Bytecode compilation and optimization
- `execution_limiter.go`: Infinite loop detection and prevention
- `template_string.go`: Multi-line string template support
- `date_literal.go`, `date_value.go`: Native date handling
- 50+ specialized AST node files: Each language construct in separate file
- Comprehensive test coverage with dedicated test files

**pkg/r2libs/**: Built-in function libraries (17,610 LOC, 42 files)
- `r2http.go`: HTTP server functionality with routing
- `r2httpclient.go`: HTTP client requests
- `r2requests.go`: Advanced HTTP request handling
- `r2soap.go`: SOAP web service support
- `r2string.go`: String manipulation and Unicode support
- `r2math.go`: Mathematical functions
- `r2io.go`: File I/O operations
- `r2date.go`: Advanced date/time handling
- `r2json.go`: JSON processing
- `r2xml.go`: XML processing
- `r2csv.go`: CSV file handling
- `r2db.go`: Database connectivity
- `r2jwt.go`: JWT token handling
- `r2hack.go`: Cryptographic utilities
- `r2collections.go`: Advanced data structures
- `r2unicode.go`: Unicode text processing
- `r2os.go`: Operating system interface
- `r2print.go`: Advanced output formatting
- `r2console.go`: Console interaction
- `r2goroutine.r2.go`: Concurrency primitives
- `r2test.go`: Testing framework
- `r2std.go`: Standard utilities
- `r2rand.go`: Random number generation
- Plus comprehensive test coverage for all libraries

**pkg/r2lang/**: High-level coordinator
- Orchestrates core components
- Registers all built-in libraries
- Manages program lifecycle

**pkg/r2repl/**: Interactive REPL
- Independent REPL implementation
- Advanced features: syntax highlighting, history, multiline support

**pkg/r2test/**: Advanced Testing Framework
- `assertions/`: Test assertion library
- `core/`: Test discovery, configuration, and execution
- `coverage/`: Code coverage analysis
- `fixtures/`: Test data management
- `mocking/`: Mock objects and test isolation
- `reporters/`: HTML, JSON, and JUnit reporting

### Language Features

**Data Types**: Numbers (float64), strings, booleans, arrays, maps, objects, functions, dates

**Control Flow**: 
- if/else, while, for loops (including for-in iteration)
- try/catch/finally with proper error handling
- break/continue statements for loop control
- Infinite loop detection and prevention

**Objects & Classes**: 
- Class inheritance with `extends` keyword
- Method overriding and `super` calls
- Constructor functions
- Object instantiation and method calls

**Functions**: 
- Named functions with `func` keyword
- Anonymous functions/lambdas
- Method calls on objects
- Extensive built-in function libraries

**Advanced Features**:
- **DSL Support**: Domain-Specific Language creation and execution
- **Template Strings**: Multi-line string templates with interpolation
- **Date Handling**: Native date/time support with formatting
- **Unicode Support**: Full Unicode text processing
- **Bytecode Compilation**: Performance optimization through bytecode
- **JIT Loop Optimization**: Just-in-time compilation for loops

**Concurrency**: 
- `r2()` function for goroutines
- Built-in synchronization primitives
- Thread-safe operations

**Testing**: 
- Built-in test case syntax with Given/When/Then/And steps
- Advanced testing framework with assertions
- Code coverage analysis
- Mock objects and test isolation
- Multiple reporting formats (HTML, JSON, JUnit)

**Imports**: Module system with `import "file.r2" as alias` syntax

**Web & Network**:
- HTTP server and client functionality
- SOAP web service support
- RESTful API development
- Request/response handling

### Built-in Libraries (pkg/r2libs/)

Each library file is self-contained and registers functions with the environment:

**Core Libraries**:
- **r2std.go**: Standard utilities (typeOf, len, sleep, parseInt, eval)
- **r2math.go**: Mathematical functions and operations
- **r2string.go**: String manipulation and Unicode support
- **r2collections.go**: Advanced data structures and algorithms
- **r2rand.go**: Random number generation

**I/O & Data Processing**:
- **r2io.go**: File I/O operations and file system access
- **r2json.go**: JSON parsing and serialization
- **r2xml.go**: XML processing and manipulation
- **r2csv.go**: CSV file handling
- **r2unicode.go**: Unicode text processing and normalization

**Network & Web**:
- **r2http.go**: HTTP server with routing and middleware
- **r2httpclient.go**: HTTP client requests
- **r2requests.go**: Advanced HTTP request handling
- **r2soap.go**: SOAP web service support

**System & OS**:
- **r2os.go**: Operating system interface
- **r2db.go**: Database connectivity (MySQL, PostgreSQL, SQLite)
- **r2console.go**: Console interaction and input handling

**Security & Cryptography**:
- **r2hack.go**: Cryptographic and security utilities
- **r2jwt.go**: JWT token creation and validation

**Development & Testing**:
- **r2test.go**: Built-in testing framework
- **r2print.go**: Advanced output formatting and debugging

**Concurrency & Performance**:
- **r2goroutine.r2.go**: Concurrency primitives and goroutine management
- **r2go.go**: Go language interoperability

**Date & Time**:
- **r2date.go**: Advanced date/time handling and formatting

**Data Visualization**:
- **r2lang_graph.go**: Graph and chart generation capabilities

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
├── main.go                      # Primary entry point
├── cmd/                         # Command-line tools
│   ├── r2/                     # Alternative R2 executable
│   ├── r2test/                 # Dedicated test runner
│   └── repl/                   # Standalone REPL
├── pkg/r2core/                 # Core interpreter (58 files, 11,964 LOC)
├── pkg/r2libs/                 # Built-in libraries (42 files, 17,610 LOC)
├── pkg/r2lang/                 # High-level coordinator (minimal, clean)
├── pkg/r2repl/                 # Interactive REPL (independent)
├── pkg/r2test/                 # Advanced testing framework
│   ├── assertions/             # Test assertion library
│   ├── core/                   # Test discovery and execution
│   ├── coverage/               # Code coverage analysis
│   ├── fixtures/               # Test data management
│   ├── mocking/                # Mock objects and isolation
│   └── reporters/              # HTML, JSON, JUnit reporting
├── examples/                   # R2 language example programs (60+ files)
│   ├── dsl/                    # Domain-Specific Language examples
│   ├── testing/                # Testing framework examples
│   └── unit_testing/           # Unit testing examples
├── docs/                       # Comprehensive documentation
│   ├── en/                     # English documentation
│   │   ├── courses/            # Learning courses
│   │   └── *.md               # Technical documentation
│   └── es/                     # Spanish documentation
│       ├── cursos/             # Cursos de aprendizaje
│       ├── dsl/                # DSL documentation
│       ├── fixes/              # Bug fixes and improvements
│       └── reportes/           # Quality and performance reports
├── vscode_syntax_highlighting/ # VS Code extension
├── r2-jetbrains-syntax/        # IDE syntax highlighting
└── proposals/                  # Future feature proposals

Architecture Benefits:
- Each module is independently testable (122 total Go files)
- Clear separation of concerns across 5 main packages
- Easy onboarding for new contributors
- Parallel development possible
- Excellent maintainability and code quality
- Comprehensive testing and documentation
- Enterprise-ready with advanced features
```

### Quality Metrics (Current State)
- **Code Quality**: 9/10 (significantly improved)
- **Maintainability**: 9/10 (excellent modular structure)
- **Testability**: 9.5/10 (comprehensive test coverage)
- **Technical Debt**: Minimal (well-structured codebase)
- **Total Lines of Code**: 29,574+ (core + libraries)
- **Test Coverage**: Extensive across all modules
- **Documentation**: Comprehensive (English + Spanish)

## Current Development Status

### Recently Implemented Features
- **DSL Support**: Complete Domain-Specific Language framework
- **Advanced Testing**: Professional testing framework with coverage analysis
- **Bytecode Compilation**: Performance optimization through bytecode
- **Infinite Loop Detection**: Automatic detection and prevention
- **Template Strings**: Multi-line string templates with interpolation
- **Native Date Support**: Full date/time handling with formatting
- **Unicode Support**: Complete Unicode text processing
- **Database Integration**: MySQL, PostgreSQL, SQLite support
- **SOAP Web Services**: Enterprise web service support
- **JWT Authentication**: Token-based authentication
- **Graph Visualization**: Data visualization capabilities

### Command-Line Tools
- **r2**: Primary R2Lang interpreter
- **r2test**: Dedicated test runner with advanced features
- **r2repl**: Interactive REPL with syntax highlighting

### IDE Support
- **VS Code Extension**: Complete syntax highlighting and language support
- **JetBrains Plugin**: IDE integration for IntelliJ family
- **Syntax Highlighting**: tmLanguage grammar definitions

### Documentation & Learning
- **Comprehensive Docs**: English and Spanish documentation
- **Learning Courses**: 10+ structured learning modules
- **60+ Examples**: Practical examples covering all features
- **Technical Reports**: Performance analysis and quality metrics

### Enterprise Features
- **Professional Testing**: Assertions, mocking, coverage analysis
- **Multi-format Reporting**: HTML, JSON, JUnit test reports
- **Database Connectivity**: Enterprise database support
- **Web Services**: HTTP servers, SOAP services, RESTful APIs
- **Security**: Cryptographic utilities and JWT support
- **Performance**: Bytecode compilation and JIT optimization