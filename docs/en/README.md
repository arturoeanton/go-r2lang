# R2Lang - Modern Programming Language

## Overview

R2Lang is an interpreted programming language that combines the familiar syntax of JavaScript with modern features like object-oriented programming, native concurrency, and an integrated testing system. Designed to be simple to learn but powerful in its capabilities.

## Key Features

### 🚀 Intuitive Syntax
```r2
func main() {
    let message = "Hello, R2Lang!"
    print(message)
}
```

### 🎯 Object-Oriented Programming with Inheritance
```r2
class Person {
    let name
    let age
    
    constructor(name, age) {
        this.name = name
        this.age = age
    }
    
    greet() {
        print("Hello, I'm " + this.name)
    }
}

class Employee extends Person {
    let salary
    
    constructor(name, age, salary) {
        super.constructor(name, age)
        this.salary = salary
    }
    
    work() {
        print(this.name + " is working")
    }
}
```

### ⚡ Native Concurrency
```r2
func task() {
    print("Running in parallel")
    sleep(1)
    print("Task completed")
}

func main() {
    r2(task)  // Execute in goroutine
    r2(task)
    sleep(2)  // Wait for completion
}
```

### 🧪 Integrated Testing System
```r2
TestCase "Verify Addition" {
    Given func() { 
        setup()
        return "Preparing data"
    }
    When func() {
        let result = 2 + 3
        return "Executing operation"
    }
    Then func() {
        assertEqual(result, 5)
        return "Validating result"
    }
}
```

### 📦 Module System
```r2
import "math.r2" as math
import "./utils.r2" as utils

func main() {
    let result = math.sqrt(16)
    utils.log("Result: " + result)
}
```

## Installation

### Prerequisites
- Go 1.23.4 or higher
- Git

### Clone and Install
```bash
git clone https://github.com/arturoeanton/go-r2lang.git
cd go-r2lang
go build -o r2lang main.go
```

### Running Programs
```bash
# Run an R2 file
./r2lang program.r2

# Run main.r2 from current directory
./r2lang

# Interactive REPL mode
./r2lang -repl

# REPL without debug output
./r2lang -repl -no-output
```

## Code Examples

### Variables and Types
```r2
let number = 42
let text = "Hello world"
let flag = true
let array = [1, 2, 3, "four"]
let map = {
    name: "John",
    age: 30,
    active: true
}
```

### Control Flow
```r2
// Conditionals
if (age >= 18) {
    print("Adult")
} else {
    print("Minor")
}

// Loops
for (let i = 0; i < 10; i++) {
    print("Iteration: " + i)
}

// For-in for arrays
for (let item in array) {
    print("Item: " + item)
}

// While
let counter = 0
while (counter < 5) {
    print("Counter: " + counter)
    counter++
}
```

### Functions and Lambdas
```r2
func add(a, b) {
    return a + b
}

// Lambda function
let multiply = func(a, b) {
    return a * b
}

// Higher-order function
let numbers = [1, 2, 3, 4, 5]
let doubled = numbers.map(func(x) { return x * 2 })
```

### Error Handling
```r2
try {
    let result = divide(10, 0)
    print("Result: " + result)
} catch (error) {
    print("Error: " + error)
} finally {
    print("Cleanup completed")
}
```

### Arrays and Operations
```r2
let fruits = ["apple", "banana", "orange"]

// Array methods
fruits.push("grape")                    // Add element
let length = fruits.len()               // Get length
let found = fruits.find("banana")       // Search element
let filtered = fruits.filter(func(f) { 
    return f.length > 5 
})

// Functional operations
let numbers = [1, 2, 3, 4, 5]
let sum = numbers.reduce(func(acc, val) { 
    return acc + val 
})
let sorted = numbers.sort()
```

## Built-in Libraries

### Standard
- `print()` - Console output
- `len()` - Length of strings/arrays
- `typeOf()` - Variable type
- `sleep()` - Pause in seconds

### Mathematics
- `math.sqrt()`, `math.pow()`, `math.sin()`, etc.
- `rand.int()`, `rand.float()` - Random numbers

### I/O and System
- `io.readFile()`, `io.writeFile()` - File operations
- `os.getEnv()`, `os.exit()` - OS interaction

### HTTP
```r2
// HTTP Server
http.server(8080, func(req, res) {
    res.json({message: "Hello from R2Lang!"})
})

// HTTP Client
let response = httpClient.get("https://api.example.com")
print(response.body)
```

### Strings
```r2
let text = "Hello World"
let uppercase = text.upper()
let words = text.split(" ")
let contains = text.contains("World")
```

## Interpreter Architecture

### Main Components

1. **Lexer** (`r2lang/r2lang.go:139-321`)
   - Source code tokenization
   - Handling numbers, strings, operators
   - Support for line and block comments

2. **Parser** (`r2lang/r2lang.go:1662-2331`)
   - Recursive descent parsing
   - AST (Abstract Syntax Tree) construction
   - Operator precedence handling

3. **AST and Evaluation** (`r2lang/r2lang.go:327-1657`)
   - AST nodes implement `Node` interface
   - Tree-walking interpreter
   - Lazy expression evaluation

4. **Environment** (`r2lang/r2lang.go:1429-1507`)
   - Scoping system with nested environments
   - Variable and function management
   - Closure support

5. **Native Libraries** (`r2lang/r2*.go`)
   - Built-in functions in Go
   - Modular language extension

## Use Cases

### Scripting and Automation
```r2
// File processing
let content = io.readFile("data.txt")
let lines = content.split("\n")
let processed = lines.map(func(line) {
    return line.trim().upper()
})
io.writeFile("output.txt", processed.join("\n"))
```

### APIs and Microservices
```r2
func handleUsers(req, res) {
    if (req.method == "GET") {
        res.json(getUsers())
    } else if (req.method == "POST") {
        let user = req.body
        createUser(user)
        res.status(201).json({message: "User created"})
    }
}

http.server(3000, handleUsers)
```

### Testing and QA
```r2
TestCase "User API" {
    Given func() {
        cleanDatabase()
        return "Database cleaned"
    }
    When func() {
        let response = httpClient.post("/api/users", {
            name: "John",
            email: "john@example.com"
        })
        return "User created via API"
    }
    Then func() {
        assertEqual(response.status, 201)
        assertTrue(response.body.id != null)
        return "Valid response"
    }
}
```

## Roadmap

### Version 1.0 (Current)
- ✅ Functional basic interpreter
- ✅ Object-oriented programming with inheritance
- ✅ Concurrency with goroutines
- ✅ Integrated testing system
- ✅ Basic libraries (I/O, HTTP, Math)

### Version 1.1 (Q1 2024)
- 🔄 Enhanced type system
- 🔄 Robust error handling
- 🔄 Integrated debugger
- 🔄 Performance optimizations

### Version 1.5 (Q2 2024)
- 📋 Generics
- 📋 Pattern matching
- 📋 Native async/await
- 📋 Package manager

### Version 2.0 (Q3 2024)
- 📋 JIT compilation
- 📋 Language Server Protocol
- 📋 WebAssembly target
- 📋 Complete standard library

## Contributing

### Environment Setup
```bash
git clone https://github.com/arturoeanton/go-r2lang.git
cd go-r2lang
go mod tidy
go test ./r2lang
```

### Project Structure
```
├── main.go              # Entry point
├── main.r2             # Main example
├── r2lang/             # Interpreter core
│   ├── r2lang.go       # Lexer, Parser, AST
│   ├── r2std.go        # Standard library
│   ├── r2http.go       # HTTP functions
│   └── ...             # Other libraries
├── examples/           # Code examples
└── docs/              # Documentation
```

### Adding New Features

1. **New Tokens**: Update constants in `r2lang.go:17-54`
2. **Syntax**: Modify lexer in `r2lang.go:139-321`
3. **AST Nodes**: Create structs implementing `Node` interface
4. **Parsing**: Add logic to parser `r2lang.go:1662-2331`
5. **Evaluation**: Implement `Eval()` method in the node
6. **Testing**: Create examples in `examples/`

### Adding Native Libraries
```go
// r2lang/r2new.go
func RegisterNew(env *Environment) {
    env.Set("newFunction", BuiltinFunction(func(args ...interface{}) interface{} {
        // Implementation
        return result
    }))
}

// In r2lang.go:RunCode(), add:
RegisterNew(env)
```

## License

This project is under the MIT license. See [LICENSE](../../LICENSE) for more details.

## Contact and Support

- **Issues**: [GitHub Issues](https://github.com/arturoeanton/go-r2lang/issues)
- **Documentation**: See `docs/` folder
- **Examples**: See `examples/` folder

---

*R2Lang - Simplicity meets Power* 🚀