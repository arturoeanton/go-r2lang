
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

**R2Lang** is a dynamic, interpreted programming language written in Go. It's designed to be simple, intuitive, and powerful, with a syntax heavily inspired by JavaScript and first-class support for **comprehensive unit testing**.

Whether you're writing automation scripts, building a web API, or creating a robust testing suite, R2Lang provides the tools you need in a clean and readable package.

## ✨ Key Features

| Feature                 | Description                                                                                                 | Example                                                              |
| ----------------------- | ----------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------- |
| **🧪 Built-in Unit Testing** | Comprehensive testing framework with `describe()` and `it()` syntax. No external frameworks needed.        | `describe("User Login", func() { it("should authenticate", func() { ... }) })`      |
| **🚀 Simple & Familiar**    | If you know JavaScript, you'll feel right at home. This makes it incredibly easy to pick up and start coding. | `let message = "Hello, World!"; print(message);`                     |
| **⚡ Concurrent**          | Leverage the power of Go's goroutines with a simple `r2()` function to run code in parallel.                | `r2(myFunction, "arg1");`                                            |
| **🧱 Object-Oriented**     | Use classes, inheritance (`extends`), and `this` to structure your code in a clean, object-oriented way.    | `class User extends Person { ... }`                                  |
| **🌐 Web Ready**            | Create web servers and REST APIs with a built-in `http` library that feels like Express.js.                 | `http.get("/users", func(req, res) { res.json(...) });`               |
| **🧩 Easily Extendable**   | Written in Go, R2Lang can be easily extended with new native functions and libraries.                       | `env.Set("myNativeFunc", r2lang.NewBuiltinFunction(...));`            |

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
        print("Hello, R2Lang! 🚀");
    }
    ```

4.  **Run it!**
    ```bash
    # Using the basic interpreter
    ./r2lang hello.r2
    
    # Using the advanced R2 command
    ./r2 hello.r2
    
    # Or with Go directly
    go run main.go hello.r2
    
    # Output: Hello, R2Lang! 🚀
    ```

---

## 🛠️ Command Line Tools

R2Lang provides specialized command-line tools for different development workflows:

### Available Commands

| Command | Purpose | Usage |
|---------|---------|-------|
| `r2lang` | Basic interpreter | `r2lang script.r2` |
| `r2` | Advanced R2 execution with optimizations | `r2 [options] script.r2` |
| `r2test` | Comprehensive testing framework | `r2test [options] [test-dirs...]` |
| `r2repl` | Interactive REPL shell | `r2repl [options]` |

### Basic Usage Examples

```bash
# Run a script with basic interpreter
r2lang script.r2

# Run a script with advanced options
r2 -verbose -optimize script.r2

# Start interactive REPL
r2repl

# Run tests with coverage
r2test -coverage -verbose

# Run tests in specific directory
r2test ./tests

# Run tests with detailed reporting
r2test -coverage -reporters json,junit -output ./reports
```

### Advanced R2 Command Options

The `r2` command provides extensive options for advanced usage:

```bash
# Performance and optimization
r2 -optimize -profile cpu script.r2

# Execution control
r2 -timeout 30s -max-memory 100MB script.r2

# Development tools
r2 -check script.r2           # Syntax check only
r2 -format script.r2          # Format code
r2 -compile script.r2         # Compile to bytecode (future)

# Debug and verbose output
r2 -debug -verbose script.r2
```

### Testing Framework Options

The `r2test` command supports comprehensive testing workflows:

```bash
# Basic test execution
r2test                        # Run all tests
r2test ./tests               # Run tests in specific directory
r2test -verbose              # Verbose output

# Test filtering
r2test -grep "Calculator"    # Run tests matching pattern
r2test -tags unit,integration # Run tests with specific tags

# Coverage and reporting
r2test -coverage -coverage-threshold 80
r2test -reporters json,junit -output ./reports

# Performance and parallel execution
r2test -parallel -workers 4 -timeout 60s
```

### REPL Options

The `r2repl` command provides an interactive development environment:

```bash
# Basic REPL
r2repl

# REPL with options
r2repl -quiet               # Minimal startup output
r2repl -no-output          # Disable output display
r2repl -debug              # Debug mode with configuration info
```

---

## 🧪 Advanced Unit Testing Framework

R2Lang features a professional-grade testing framework with enterprise-level capabilities including **mocking**, **fixtures**, **coverage reporting**, and **test isolation**:

### Basic Test Structure
```r2
describe("Calculator", func() {
    it("should add numbers correctly", func() {
        let result = 2 + 3;
        assert.equals(result, 5);
    });
    
    it("should handle subtraction", func() {
        let result = 10 - 4;
        assert.equals(result, 6);
    });
    
    it("should handle multiplication", func() {
        let result = 3 * 4;
        assert.equals(result, 12);
    });
});
```

### Available Assertions
```r2
// Basic comparisons
assert.equals(actual, expected);
assert.notEquals(actual, expected);
assert.true(value);
assert.false(value);

// Numeric comparisons
assert.greater(10, 5);
assert.greaterOrEqual(10, 10);
assert.less(5, 10);
assert.lessOrEqual(5, 5);

// String operations
assert.contains("Hello World", "World");
assert.notContains("Hello", "Goodbye");

// Collection operations
assert.hasLength(array, 5);
assert.empty([]);
assert.notEmpty([1, 2, 3]);

// Null checks
assert.nil(null);
assert.notNil("value");

// Error handling
assert.panics(func() { throw "error"; });
assert.notPanics(func() { return "ok"; });
```

### Test Lifecycle Hooks
```r2
test.describe("Database Tests", func() {
    test.beforeAll(func() {
        // Setup database connection
    });
    
    test.beforeEach(func() {
        // Reset test data before each test
    });
    
    test.afterEach(func() {
        // Cleanup after each test
    });
    
    test.afterAll(func() {
        // Close database connection
    });
    
    test.it("should save user data", func() {
        // Test implementation
    });
});
```

### Comprehensive Assertion Library
- **Basic**: `.equals()`, `.isTrue()`, `.isFalse()`, `.isNull()`, `.isNotNull()`
- **Types**: `.isNumber()`, `.isString()`, `.isBoolean()`, `.isArray()`, `.isObject()`
- **Comparisons**: `.isGreaterThan()`, `.isLessThan()`, `.isGreaterThanOrEqual()`
- **Strings**: `.contains()`, `.startsWith()`, `.endsWith()`, `.matches()`
- **Collections**: `.hasProperty()`, `.hasLength()`, `.isEmpty()`, `.isNotEmpty()`
- **Exceptions**: `.throws()`, `.doesNotThrow()`, `.withMessage()`

### Mocking and Spying
```r2
import "r2test/mocking" as mock;

test.describe("HTTP Service Tests", func() {
    test.it("should mock API calls", func() {
        let httpMock = mock.createMock("httpService");
        httpMock.when("get", "/api/users").returns({users: []});
        
        let result = httpMock.call("get", "/api/users");
        test.assert("API result").that(result.users).isArray();
        test.assert("Mock called").that(httpMock.wasCalled("get")).isTrue();
    });
    
    test.it("should spy on methods", func() {
        let spy = mock.spyOn("userService.validate", userService.validate);
        spy.callThrough();
        
        userService.validate("test@email.com");
        
        test.assert("Spy called").that(spy.wasCalledWith("test@email.com")).isTrue();
    });
});
```

### Fixture Management
```r2
import "r2test/fixtures" as fixtures;

test.describe("Data Tests", func() {
    test.it("should load JSON fixtures", func() {
        let userData = fixtures.load("user_data.json");
        test.assert("User name").that(userData.name).equals("John Doe");
    });
    
    test.it("should create temporary fixtures", func() {
        fixtures.createTemporary("test_data", {id: 1, name: "Test"});
        let data = fixtures.getData("test_data");
        test.assert("Temp data").that(data.id).equals(1);
    });
});
```

### Coverage Reporting
```r2
import "r2test/coverage" as coverage;
import "r2test/reporters" as reporters;

// Enable coverage collection
coverage.enable();
coverage.start();

// Run tests and generate reports
test.runTests();

// Generate multiple report formats
reporters.generateHTMLReport("./coverage/html");    // Interactive HTML report
reporters.generateJSONReport("./coverage/data.json"); // Machine-readable JSON
reporters.generateJUnitReport("./coverage/junit.xml"); // CI/CD integration

// Check coverage thresholds
let stats = coverage.getStats();
if (stats.linePercentage < 80) {
    throw "Coverage below 80% threshold";
}
```

### Test Isolation
```r2
import "r2test/mocking" as mock;

test.describe("Isolated Tests", func() {
    test.it("should run in isolation", func() {
        mock.runInIsolation("isolated test", func(context) {
            let mockDb = context.createMock("database");
            mockDb.when("save").returns(123);
            
            // Test runs in complete isolation
            let result = mockDb.call("save", {data: "test"});
            test.assert("Isolated result").that(result).equals(123);
        });
    });
});
```

### Running Tests
```bash
# Run specific test files
go run main.go examples/unit_testing/basic_test_example.r2
go run main.go examples/unit_testing/mocking_example.r2
go run main.go examples/unit_testing/coverage_example.r2

# Run all tests in directory
go run main.go -test ./tests/

# Generate coverage reports
go run main.go -test -coverage ./tests/
```

### Advanced Features
- **Test Discovery**: Automatic test file detection with patterns
- **Parallel Execution**: Run tests in parallel for faster execution
- **Watch Mode**: Re-run tests when files change
- **Custom Reporters**: Pluggable reporting system
- **CI/CD Integration**: JUnit XML and JSON output for build systems

For complete examples and advanced features, see [examples/unit_testing/](./examples/unit_testing/).

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
