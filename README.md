
<div align="center">
  <br />
  <h1>R2Lang</h1>
  <p>
    <b>Write elegant tests, scripts, and applications with a language that blends simplicity and power.</b>
  </p>
  <br />
</div>

<div align="center">

[![Go Report Card](https://goreportcard.com/badge/github.com/arturoeanton/go-r2lang)](https://goreportcard.com/report/github.com/arturoeanton/go-r2lang)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GitHub stars](https://img.shields.io/github/stars/arturoeanton/go-r2lang.svg?style=social&label=Star)](https://github.com/arturoeanton/go-r2lang)
[![GitHub forks](https://img.shields.io/github/forks/arturoeanton/go-r2lang.svg?style=social&label=Fork)](https://github.com/arturoeanton/go-r2lang)
[![GitHub issues](https://img.shields.io/github/issues/arturoeanton/go-r2lang.svg)](https://github.com/arturoeanton/go-r2lang/issues)
[![Contributors](https://img.shields.io/github/contributors/arturoeanton/go-r2lang.svg)](https://github.com/arturoeanton/go-r2lang/graphs/contributors)

</div>

---

**R2Lang** is a modern, dynamic programming language written in Go that combines the simplicity of JavaScript with powerful features for modern development. With its clean modular architecture and comprehensive built-in capabilities, R2Lang 2025 offers everything you need for scripting, web development, testing, and application building.

Whether you're writing automation scripts, building web APIs, creating robust test suites, or developing concurrent applications, R2Lang provides professional-grade tools in an elegant, readable package.

## ✨ Key Features

| Feature | Description | Example |
|---------|-------------|---------|
| **🎯 Modern Language (2025)** | JavaScript-style syntax with `true`/`false` literals, multiline maps, `else if`, and `%` operator | `if (x % 2 == 0) { ... } else if (x % 3 == 0) { ... }` |
| **🔧 DSL Builder** | **🌟 ORIGINAL** - Create custom Domain-Specific Languages with simple syntax | `dsl Calculator { token("NUM", "[0-9]+"); rule("sum", ["NUM", "+", "NUM"], "add"); }` |
| **🗺️ Advanced Maps** | JavaScript-style map literals with multiline support and mixed separators | `let config = { host: "localhost", port: 8080, ssl: true }` |
| **📝 Template Strings** | String interpolation with backticks and multiple variables | `let msg = \`Hello ${name}, you have ${count} messages\`` |
| **🌍 Unicode Support** | Full international character support in strings and identifiers | `let año = 2025; let текст = "Привет мир"; let 名前 = "田中"` |
| **📅 Native Dates** | Built-in date types with comprehensive formatting | `let date = @2024-12-25; Date.format(date, "YYYY-MM-DD")` |
| **🧪 Testing Ready** | Comprehensive test suite with 416+ test cases validating all features | `go test ./pkg/r2core/ && go run main.go gold_test.r2` |
| **⚡ Concurrent** | Goroutines with simple `r2()` function and synchronization primitives | `r2(processData, userData); sleep(100);` |
| **🏗️ Object-Oriented** | Classes with inheritance, method overriding, and `super` calls | `class Manager extends Employee { super(); this.team = []; }` |
| **🌐 Web Ready** | Built-in HTTP server and client with REST API support | `http.get("/api/users", handleUsers); request.post(url, data)` |
| **⚡ gRPC Support** | Dynamic gRPC client without code generation, all streaming types | `let client = grpc.grpcClient("service.proto", "localhost:9090")` |
| **🛡️ Secure & Safe** | Infinite loop protection, timeout controls, and error handling | `try { risky(); } catch (e) { log(e); } finally { cleanup(); }` |

---

## 🌟 DSL Builder - Our Most Original Feature

**R2Lang's DSL Builder is our most innovative feature**, allowing you to create custom Domain-Specific Languages with an elegantly simple syntax. This sets R2Lang apart from other languages by making parser creation as easy as writing a function.

### Why DSL Builder is Special

Unlike complex parser generators like ANTLR or Lex/Yacc, R2Lang's DSL system is:
- **Native Integration**: DSL code runs directly in R2Lang
- **Zero Setup**: No external tools or code generation
- **Intuitive Syntax**: Declarative and readable
- **Instant Results**: From concept to working parser in minutes

### Quick DSL Example

```r2
// Define a simple calculator DSL
dsl Calculator {
    token("NUMBER", "[0-9]+")
    token("PLUS", "\\+")
    token("MINUS", "-")
    token("MULTIPLY", "\\*")
    token("DIVIDE", "/")
    
    rule("operation", ["NUMBER", "operator", "NUMBER"], "calculate")
    rule("operator", ["PLUS"], "plus")
    rule("operator", ["MINUS"], "minus")
    rule("operator", ["MULTIPLY"], "multiply")
    rule("operator", ["DIVIDE"], "divide")
    
    func calculate(left, op, right) {
        var l = parseFloat(left)
        var r = parseFloat(right)
        
        if (op == "+") return l + r
        if (op == "-") return l - r
        if (op == "*") return l * r
        if (op == "/") return l / r
    }
    
    func plus(token) { return "+" }
    func minus(token) { return "-" }
    func multiply(token) { return "*" }
    func divide(token) { return "/" }
}

// Use your DSL
func main() {
    var calc = Calculator.use
    
    var result = calc("15 + 25")
    console.log(result.Output)  // 40
    
    console.log(result)  // "DSL[15 + 25] -> 40"
}
```

### DSL vs Traditional Parsers

| Feature | R2Lang DSL | ANTLR | Lex/Yacc |
|---------|------------|-------|-----------|
| **Setup Time** | Minutes | Hours | Days |
| **Code Generation** | None | Required | Required |
| **Learning Curve** | Minimal | Steep | Very Steep |
| **Integration** | Native | External | External |
| **Debugging** | R2Lang tools | Specialized | Complex |
| **Result Access** | `result.Output` | Generated code | Generated code |

### DSL Use Cases

- **Configuration Languages**: Custom config file formats
- **Command Systems**: Domain-specific command languages
- **Data Validators**: Custom validation rules
- **Text Processors**: Specialized text parsing
- **Business Rules**: Domain-specific business logic

### Learn More

- [**Complete DSL Documentation**](./docs/es/dsl/) - Full guide and examples
- [**DSL Examples**](./examples/dsl/) - Working calculator and command examples
- [**DSL API Reference**](./docs/es/dsl/referencia_rapida.md) - Quick reference guide

---

## 🚀 Getting Started

### Prerequisites

- **Go**: Version 1.23 or higher.

### Installation & "Hello, World!"

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/arturoeanton/go-r2lang.git
    cd go-r2lang
    ```

2.  **Build the tools:**
    ```bash
    # Build main R2Lang interpreter
    go build -o r2lang main.go
    
    # Build specialized commands
    go build -o r2 cmd/r2/main.go
    go build -o r2test cmd/r2test/main.go
    go build -o r2repl cmd/repl/main.go
    ```

3.  **Create your first R2Lang file (`hello.r2`):**
    ```r2
    func main() {
        let name = "R2Lang 2025";
        let version = "1.0";
        
        // Modern template strings with interpolation
        let message = `Welcome to ${name} v${version}! 🚀`;
        print(message);
        
        // Multiline maps with modern syntax
        let features = {
            unicode: true,
            dates: true,
            templates: true
            multiline_maps: true,
            else_if_syntax: true,
            modulo_operator: true
        };
        
        print("New features:", len(features));
        
        // Modern else if syntax
        if (len(features) > 5) {
            print("🎉 Fully featured language!");
        } else if (len(features) > 3) {
            print("👍 Good feature set");
        } else {
            print("📝 Basic features");
        }
    }
    ```

4.  **Run it!**
    ```bash
    # Using Go directly (recommended)
    go run main.go hello.r2
    
    # Or start the REPL for interactive exploration
    go run main.go -repl
    
    # Output: 
    # Welcome to R2Lang 2025 v1.0! 🚀
    # New features: 6
    # 🎉 Fully featured language!
    ```

---

## 🛠️ Quick Start Examples

### Basic Usage
```bash
# Run a script
go run main.go script.r2

# Start interactive REPL
go run main.go -repl

# Run without output (useful for testing)
go run main.go -repl -no-output

# Run examples
go run main.go examples/example01-variables.r2
go run main.go examples/example37-map-literals.r2
go run main.go examples/example38-for-in-loops.r2
```

### Advanced Examples
```bash
# Run the comprehensive gold test (validates all features)
go run main.go gold_test.r2

# Test specific modules
go test ./pkg/r2core/
go test ./pkg/r2libs/

# Try our innovative DSL examples
go run main.go examples/dsl/calculadora_dsl.r2
go run main.go examples/dsl/comando_simple.r2

# Try gRPC examples (terminal 1: server, terminal 2: client)  
cd examples/grpc/example1 && go run simple_grpc_server.go
go run main.go examples/grpc/example1/introspection_demo.r2

# Build the interpreter
go build -o r2lang main.go
./r2lang script.r2
```

## 🎯 R2Lang 2025 Language Features

### Modern Syntax Examples

```r2
// Boolean literals and modern conditionals
let isValid = true;
let count = 15;

if (count % 3 == 0) {
    print("Divisible by 3");
} else if (count % 2 == 0) {
    print("Even number");
} else {
    print("Odd number");
}

// Multiline maps with mixed separators
let config = {
    database: {
        host: "localhost"
        port: 5432,
        ssl: true,
        timeout: 30
    },
    api: {
        version: "v1",
        rate_limit: 1000
        auth: true
    }
};

// Template strings with interpolation
let user = "Alice";
let age = 30;
let message = `User ${user} is ${age} years old and has ${len(config)} config sections`;

// Unicode support in identifiers and strings
let año = 2025;
let 用户名 = "张三";
let emoji = "🚀 R2Lang rocks! 🎉";

// Native date formatting
let birthday = @1990-05-15;
let formatted = Date.format(birthday, "DD/MM/YYYY");
print(`Birthday: ${formatted}`);

// For-in loops with modern syntax
let scores = { alice: 95, bob: 87, charlie: 92 };
let studentKeys = keys(scores);

for (i in studentKeys) {
    let student = studentKeys[$k];
    let score = scores[student];
    
    if (score >= 90) {
        print(`${student}: A grade (${score})`);
    } else if (score >= 80) {
        print(`${student}: B grade (${score})`);
    } else {
        print(`${student}: C grade (${score})`);
    }
}
```

---

## 🏗️ Architecture & Quality

### Modular Design (2025)
R2Lang features a clean, modular architecture that eliminated the previous monolithic design:

- **pkg/r2core/**: Core interpreter (2,590 LOC, 30 specialized files)
- **pkg/r2libs/**: Built-in libraries (3,701 LOC, 18 specialized modules)  
- **pkg/r2lang/**: High-level coordinator (45 LOC)
- **pkg/r2repl/**: Interactive REPL (185 LOC)

### Quality Metrics
- **Code Quality**: 8.5/10 (up from 6.2/10)
- **Maintainability**: 8.5/10 (up from 2/10)  
- **Testability**: 9/10 (up from 1/10)
- **Technical Debt Reduction**: 79%

### Comprehensive Testing
```bash
# Run all Go tests (416+ test cases)
go test ./pkg/...

# Run R2Lang gold test (validates all language features)
go run main.go gold_test.r2

# Test specific modules
go test ./pkg/r2core/    # Core interpreter tests
go test ./pkg/r2libs/    # Built-in library tests
```

### Built-in Libraries
- **r2hack.go**: Cryptographic utilities (509 LOC)
- **r2http.go**: HTTP server with routing (410 LOC)
- **r2httpclient.go**: HTTP client requests (324 LOC)
- **r2grpc.go**: Dynamic gRPC client without code generation (1,467 LOC)
- **r2print.go**: Advanced output formatting (365 LOC)
- **r2os.go**: Operating system interface (245 LOC)
- **r2string.go**: String manipulation (194 LOC)
- **r2math.go**: Mathematical functions (87 LOC)
- **r2std.go**: Standard utilities (122 LOC)
- Plus 10 additional specialized libraries

---

## 📚 Documentation & Full Course

Ready to dive deeper? We have a complete, module-by-module course to take you from beginner to expert.

-   [**Read the Full Course (English)**](./docs/en/README.md)
-   [**Leer el Curso Completo (Español)**](./docs/es/README.md)

The documentation covers everything from basic syntax to advanced topics like concurrency, error handling, and web development.

---

## 💖 Contributing

**We are actively looking for contributors!** Whether you're a seasoned developer, a documentation writer, or just enthusiastic about new programming languages, we'd love your help.

Here’s how you can contribute:

1.  **Find an issue:** Check out our [**Issues**](https://github.com/arturoeanton/go-r2lang/issues) and look for `good first issue` or `help wanted` tags.
2.  **Explore the Roadmap:** See our [**Technical Roadmap**](./docs/en/roadmap.md) for long-term goals and big features we need help with.
3.  **Improve Documentation:** Found a typo or a section that could be clearer? Let us know!
4.  **Submit a Pull Request:**
    -   Fork the repository.
    -   Create a new branch (`git checkout -b feature/my-awesome-feature`).
    -   Commit your changes.
    -   Open a Pull Request!

We believe in a welcoming and supportive community. No contribution is too small!

---

## 🗺️ Project Roadmap

We have big plans for R2Lang! Our goal is to make it a fast, reliable, and feature-rich language for a wide range of applications.

Key areas of focus include:

-   **🚀 Performance Revolution:** Implementing a bytecode VM and eventually a JIT compiler for significant speed boosts.
-   **🧠 Advanced Features:** Adding pattern matching, a more sophisticated type system, and advanced concurrency models.
-   **🛠️ Richer Standard Library:** Expanding the built-in libraries for databases, file systems, and more.
-   **📦 Package Manager:** Creating a dedicated package manager for sharing and reusing R2Lang code.

For a detailed look at our plans, check out the [**Technical Roadmap**](./docs/en/roadmap.md) and our [**TODO List**](./TODO.md).

---

## 🤝 Community

-   **Report a Bug:** Found something wrong? Open an [**Issue**](https://github.com/arturoeanton/go-r2lang/issues/new).
-   **Request a Feature:** Have a great idea? Let's discuss it in the [**Issues**](https://github.com/arturoeanton/go-r2lang/issues).
-   **Ask a Question:** Don't hesitate to open an issue for questions and discussions.

---

## 📜 License

R2Lang is licensed under the **Apache License 2.0**. See the [LICENSE](./LICENSE) file for details.
