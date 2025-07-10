# R2Lang Course - Module 1: Introduction and Fundamentals

## Welcome to R2Lang

### What is R2Lang?

R2Lang is a modern programming language that combines the simplicity of JavaScript with advanced features such as:

- **Integrated Testing**: Native BDD framework
- **Object-Oriented**: Classes with inheritance
- **Concurrency**: Simple primitives for parallel programming
- **Familiar Syntax**: Easy to learn for web developers

### Why Learn R2Lang?

1. **Simple Syntax**: Based on JavaScript, easy to read and write
2. **Testing-First**: Ideal for learning good testing practices
3. **Versatility**: Scripts, web servers, automation, prototypes
4. **Modern**: Includes features from contemporary languages

## Installation and Setup

### Prerequisites

```bash
# Verify that you have Go installed
go version
# Should show Go 1.23.4 or higher
```

### Installation

```bash
# 1. Clone the repository
git clone https://github.com/arturoeanton/go-r2lang.git
cd go-r2lang

# 2. Compile the interpreter
go build -o r2lang main.go

# 3. Verify installation
./r2lang --version
```

### Your First Program

Create a file named `hello.r2`:

```r2
func main() {
    print("Hello, R2Lang!")
}
```

Run it:

```bash
./r2lang hello.r2
```

You should see: `Hello, R2Lang!`

## Fundamental Concepts

### 1. Variables and Types

#### Variable Declaration

```r2
// Declaration with an initial value
let name = "John"
let age = 25
let active = true

// Declaration without an initial value (undefined)
let result
```

#### Basic Data Types

```r2
// Numbers (all are float64 internally)
let integer = 42
let decimal = 3.14159
let negative = -100

// Strings
let greeting = "Hello world"
let message = 'Single quotes can also be used'

// Booleans
let isTrue = true
let isFalse = false

// Null
let empty = null
let undefined = nil  // Synonym for null
```

#### Type Checking

```r2
func main() {
    let number = 42
    let text = "Hello"
    let flag = true
    
    print("Type of number:", typeOf(number))    // float64
    print("Type of text:", typeOf(text))      // string
    print("Type of flag:", typeOf(flag))  // bool
}
```

### 2. Operators

#### Arithmetic Operators

```r2
func main() {
    let a = 10
    let b = 3
    
    print("Sum:", a + b)        // 13
    print("Subtraction:", a - b)       // 7
    print("Multiplication:", a * b)  // 30
    print("Division:", a / b)    // 3.3333...
}
```

#### Comparison Operators

```r2
func main() {
    let x = 5
    let y = 10
    
    print("x == y:", x == y)     // false
    print("x != y:", x != y)     // true
    print("x < y:", x < y)       // true
    print("x > y:", x > y)       // false
    print("x <= y:", x <= y)     // true
    print("x >= y:", x >= y)     // false
}
```

#### Logical Operators

```r2
func main() {
    let a = true
    let b = false
    
    print("a && b:", a && b)     // false (logical AND)
    print("a || b:", a || b)     // true (logical OR)
    print("!a:", !a)             // false (logical NOT)
}
```

### 3. Strings

#### Basic Operations

```r2
func main() {
    let firstName = "John"
    let lastName = "Doe"
    
    // Concatenation
    let fullName = firstName + " " + lastName
    print("Full name:", fullName)
    
    // Length
    print("Length of name:", len(firstName))
    
    // Convert to uppercase/lowercase (using built-ins)
    print("Uppercase:", firstName.upper())
    print("Lowercase:", lastName.lower())
}
```

#### String Methods

```r2
func main() {
    let phrase = "Hello world from R2Lang"
    
    // Split into words
    let words = phrase.split(" ")
    print("Words:", words)
    
    // Check if it contains text
    let contains = phrase.contains("world")
    print("Contains 'world':", contains)
    
    // Length
    print("Length:", phrase.length())
}
```

### 4. Input and Output

#### print() Function

```r2
func main() {
    // Print multiple values
    print("Name:", "John", "Age:", 25)
    
    // Print variables
    let message = "Hello!"
    print(message)
    
    // Print result of operations
    print("2 + 2 =", 2 + 2)
}
```

## Practical Exercises

### Exercise 1: Variables and Operations

Create a program that:
1. Declares variables for your name, age, and city
2. Calculates your birth year
3. Prints a full introduction

```r2
func main() {
    // Your solution here
    let name = "Your Name"
    let age = 25
    let city = "Your City"
    
    let currentYear = 2024
    let birthYear = currentYear - age
    
    print("Hello, I'm", name)
    print("I am", age, "years old")
    print("I live in", city)
    print("I was born in", birthYear)
}
```

### Exercise 2: Basic Calculator

Create a program that performs mathematical operations:

```r2
func main() {
    let a = 15
    let b = 4
    
    print("Numbers:", a, "and", b)
    print("Sum:", a + b)
    print("Subtraction:", a - b)
    print("Multiplication:", a * b)
    print("Division:", a / b)
    
    // Bonus: Check if a is greater than b
    if (a > b) {
        print(a, "is greater than", b)
    } else {
        print(a, "is not greater than", b)
    }
}
```

### Exercise 3: String Manipulation

Create a program that works with strings:

```r2
func main() {
    let phrase = "R2Lang is a modern language"
    
    print("Original phrase:", phrase)
    print("Length:", len(phrase))
    print("Uppercase:", phrase.upper())
    print("Lowercase:", phrase.lower())
    
    // Split words
    let words = phrase.split(" ")
    print("Words:", words)
    print("First word:", words[0])
    print("Last word:", words[words.length() - 1])
}
```

## Best Practices

### 1. Descriptive Names

```r2
// ❌ Bad
let x = 25
let y = "John"

// ✅ Good
let age = 25
let userName = "John"
```

### 2. Useful Comments

```r2
func main() {
    // Initial setup
    let price = 100
    let discount = 0.15  // 15% discount
    
    // Calculate final price
    let finalPrice = price * (1 - discount)
    print("Final price:", finalPrice)
}
```

### 3. Code Organization

```r2
func main() {
    // 1. Declare variables
    let base = 10
    let height = 5
    
    // 2. Perform calculations
    let area = base * height
    
    // 3. Display results
    print("Area of the rectangle:", area)
}
```

## Common Errors

### 1. Undeclared Variables

```r2
func main() {
    print(name)  // ❌ Error: undeclared variable
}

// Solution:
func main() {
    let name = "John"  // ✅ Declare first
    print(name)
}
```

### 2. Incompatible Types

```r2
func main() {
    let number = 5
    let text = "10"
    
    // ❌ May cause unexpected behavior
    print(number + text)  // Concatenation: "510"
    
    // ✅ Better to be explicit
    print("Number:", number, "Text:", text)
}
```

### 3. Division by Zero

```r2
func main() {
    let a = 10
    let b = 0
    
    // ❌ Will cause a runtime error
    print(a / b)
    
    // ✅ Check before dividing
    if (b != 0) {
        print("Division:", a / b)
    } else {
        print("Error: cannot divide by zero")
    }
}
```

## Module Project

### Personal Calculator

Create a program that functions as a personal calculator with the following features:

```r2
func main() {
    // Personal information
    let name = "Your Name"
    let monthlySalary = 2500
    let fixedExpenses = 1200
    let monthlySavings = monthlySalary - fixedExpenses
    
    print("=== PERSONAL CALCULATOR ===")
    print("User:", name)
    print()
    
    // Financial calculations
    print("MONTHLY FINANCES:")
    print("Salary:", monthlySalary)
    print("Fixed expenses:", fixedExpenses)
    print("Monthly savings:", monthlySavings)
    print()
    
    // Annual projections
    let annualSavings = monthlySavings * 12
    print("ANNUAL PROJECTION:")
    print("Annual savings:", annualSavings)
    
    // Percentage analysis
    let savingsPercentage = (monthlySavings / monthlySalary) * 100
    print("Savings percentage:", savingsPercentage + "%")
    
    // Advice based on savings
    if (savingsPercentage >= 20) {
        print("Excellent! You have a good level of savings.")
    } else if (savingsPercentage >= 10) {
        print("Good job, but you could save a little more.")
    } else {
        print("Consider reviewing your expenses to save more.")
    }
}
```

## Module Summary

In this first module, you have learned:

### Key Concepts
- ✅ What R2Lang is and its features
- ✅ Installation and setup
- ✅ Variables and data types
- ✅ Arithmetic, comparison, and logical operators
- ✅ Basic string manipulation
- ✅ Input and output with print()

### Skills Developed
- ✅ Create and run basic R2Lang programs
- ✅ Declare and use variables of different types
- ✅ Perform mathematical and logical operations
- ✅ Work with strings
- ✅ Write clean and well-commented code

### Next Module

In **Module 2**, you will learn:
- Control flow (if/else, while, for)
- User-defined functions
- Arrays and their manipulation
- Basic error handling

### Additional Resources

- **Examples**: Check the `examples/` folder in the repository
- **Documentation**: Read the files in `docs/en/`
- **Practice**: Experiment by modifying the given examples

Congratulations on completing Module 1! You are ready to continue with more advanced control structures.
