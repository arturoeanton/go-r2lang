# R2Lang Course - Module 2: Control Flow and Functions

## Introduction

In this module, you will learn to control the execution flow of your programs and create reusable functions. These are the fundamental tools for creating more complex and organized programs.

## Control Flow

### 1. Conditionals (if/else)

#### Basic Syntax

```r2
func main() {
    let age = 18
    
    if (age >= 18) {
        print("You are of legal age")
    } else {
        print("You are a minor")
    }
}
```

#### Multiple Conditionals (else if)

```r2
func main() {
    let grade = 85
    
    if (grade >= 90) {
        print("Excellent!")
    } else if (grade >= 80) {
        print("Very good!")
    } else if (grade >= 70) {
        print("Good")
    } else if (grade >= 60) {
        print("Sufficient")
    } else {
        print("You need to study more")
    }
}
```

#### Logical Operators in Conditionals

```r2
func main() {
    let age = 25
    let hasLicense = true
    let hasExperience = false
    
    // AND (&&)
    if (age >= 18 && hasLicense) {
        print("Can drive")
    }
    
    // OR (||)
    if (hasLicense || hasExperience) {
        print("Has some qualification")
    }
    
    // NOT (!)
    if (!hasExperience) {
        print("Needs practice")
    }
    
    // Complex combinations
    if ((age >= 21 && hasLicense) || hasExperience) {
        print("Can drive commercial vehicles")
    }
}
```

### 2. Loops

#### While Loop

```r2
func main() {
    let counter = 1
    
    while (counter <= 5) {
        print("Counter:", counter)
        counter++  // Increment by 1
    }
    
    print("Loop finished")
}
```

#### Practical Example: Sum of Numbers

```r2
func main() {
    let number = 1
    let sum = 0
    
    while (number <= 10) {
        sum = sum + number
        print("Adding", number, "- Total:", sum)
        number++
    }
    
    print("Total sum from 1 to 10:", sum)
}
```

#### Traditional For Loop

```r2
func main() {
    // Syntax: for (initialization; condition; increment)
    for (let i = 1; i <= 5; i++) {
        print("Iteration", i)
    }
    
    // Countdown
    for (let i = 10; i >= 1; i--) {
        print("Countdown:", i)
    }
    print("Liftoff!")
}
```

#### For Loop with Arrays (for-in)

```r2
func main() {
    let fruits = ["apple", "banana", "orange", "grape"]
    
    // Iterate over elements
    for (let fruit in fruits) {
        print("Fruit:", fruit)
    }
    
    // Iterate with indices
    for (let i = 0; i < fruits.length(); i++) {
        print("Position", i, ":", fruits[i])
    }
}
```

### 3. Loop Control

#### Break - Exiting the Loop

```r2
func main() {
    let number = 1
    
    while (true) {  // Infinite loop
        if (number > 5) {
            break  // Exit the loop
        }
        print("Number:", number)
        number++
    }
    print("Loop finished with break")
}
```

#### Continue - Skipping an Iteration

```r2
func main() {
    for (let i = 1; i <= 10; i++) {
        if (i % 2 == 0) {  // If it's even
            continue  // Skip to the next one
        }
        print("Odd number:", i)
    }
}
```

## Functions

### 1. Function Definition and Calling

#### Simple Function

```r2
func greet() {
    print("Hello from a function!")
}

func main() {
    greet()  // Call the function
    greet()  // Call it again
}
```

#### Functions with Parameters

```r2
func greetPerson(name) {
    print("Hello", name + "!")
}

func add(a, b) {
    let result = a + b
    print(a, "+", b, "=", result)
}

func main() {
    greetPerson("Ana")
    greetPerson("Carlos")
    
    add(5, 3)
    add(10, 25)
}
```

#### Functions with Return Values

```r2
func multiply(a, b) {
    return a * b
}

func isOfLegalAge(age) {
    return age >= 18
}

func getMessage(name, age) {
    if (isOfLegalAge(age)) {
        return name + " is of legal age"
    } else {
        return name + " is a minor"
    }
}

func main() {
    let result = multiply(6, 7)
    print("6 × 7 =", result)
    
    let message = getMessage("Laura", 22)
    print(message)
    
    // Use function directly in a conditional
    if (isOfLegalAge(16)) {
        print("Can vote")
    } else {
        print("Cannot vote")
    }
}
```

### 2. Variable Scope

#### Local vs. Global Variables

```r2
let globalVariable = "I am global"

func myFunction() {
    let localVariable = "I am local"
    print("Inside function:")
    print("- Global:", globalVariable)
    print("- Local:", localVariable)
}

func main() {
    print("In main:")
    print("- Global:", globalVariable)
    // print("- Local:", localVariable)  // ❌ Error: does not exist here
    
    myFunction()
}
```

#### Parameters are Local

```r2
func modifyParameter(number) {
    number = number + 10
    print("Inside function:", number)
    return number
}

func main() {
    let myNumber = 5
    print("Before:", myNumber)
    
    let newNumber = modifyParameter(myNumber)
    print("After:", myNumber)      // Still 5
    print("Returned:", newNumber) // Is 15
}
```

### 3. Advanced Functions

#### Functions as Variables

```r2
func add(a, b) {
    return a + b
}

func subtract(a, b) {
    return a - b
}

func main() {
    // Assign function to a variable
    let operation = add
    let result = operation(10, 5)
    print("Result:", result)
    
    // Change the operation
    operation = subtract
    result = operation(10, 5)
    print("Result:", result)
}
```

#### Anonymous Functions (Lambda)

```r2
func main() {
    // Anonymous function assigned to a variable
    let double = func(x) {
        return x * 2
    }
    
    print("Double 7:", double(7))
    
    // Direct anonymous function
    let result = func(a, b) {
        return a * b + 10
    }(5, 3)
    
    print("Result:", result)  // (5*3)+10 = 25
}
```

## Arrays and Collections

### 1. Declaration and Access

```r2
func main() {
    // Create an empty array
    let numbers = []
    
    // Create an array with elements
    let fruits = ["apple", "banana", "orange"]
    let ages = [25, 30, 18, 45]
    
    // Access elements
    print("First fruit:", fruits[0])
    print("Last age:", ages[ages.length() - 1])
    
    // Modify elements
    fruits[1] = "strawberry"
    print("Modified fruits:", fruits)
}
```

### 2. Array Methods

```r2
func main() {
    let numbers = [1, 2, 3]
    
    // Add elements
    numbers = numbers.push(4)
    numbers = numbers.push(5)
    print("After push:", numbers)
    
    // Length
    print("Length:", numbers.length())
    
    // Find element
    let position = numbers.find(3)
    print("Position of 3:", position)
    
    // Check if it contains
    let contains = numbers.find(10)
    if (contains != null) {
        print("Contains 10")
    } else {
        print("Does not contain 10")
    }
}
```

### 3. Iterating over Arrays

```r2
func printArray(arr, name) {
    print("=== " + name + " ===")
    for (let i = 0; i < arr.length(); i++) {
        print("Position", i, ":", arr[i])
    }
}

func main() {
    let colors = ["red", "green", "blue"]
    let numbers = [10, 20, 30, 40]
    
    printArray(colors, "COLORS")
    printArray(numbers, "NUMBERS")
}
```

## Practical Exercises

### Exercise 1: Average Calculator

```r2
func calculateAverage(numbers) {
    let sum = 0
    let count = numbers.length()
    
    for (let i = 0; i < count; i++) {
        sum = sum + numbers[i]
    }
    
    return sum / count
}

func main() {
    let grades = [85, 92, 78, 90, 88]
    let average = calculateAverage(grades)
    
    print("Grades:", grades)
    print("Average:", average)
    
    if (average >= 90) {
        print("Grade: Excellent")
    } else if (average >= 80) {
        print("Grade: Very Good")
    } else if (average >= 70) {
        print("Grade: Good")
    } else {
        print("Grade: Needs Improvement")
    }
}
```

### Exercise 2: Even and Odd Numbers

```r2
func classifyNumbers(numbers) {
    let evens = []
    let odds = []
    
    for (let number in numbers) {
        if (number % 2 == 0) {
            evens = evens.push(number)
        } else {
            odds = odds.push(number)
        }
    }
    
    print("Original numbers:", numbers)
    print("Even numbers:", evens)
    print("Odd numbers:", odds)
}

func main() {
    let numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
    classifyNumbers(numbers)
}
```

### Exercise 3: Word Finder

```r2
func findWord(words, searchedWord) {
    let found = []
    
    for (let i = 0; i < words.length(); i++) {
        let word = words[i]
        if (word.contains(searchedWord)) {
            found = found.push(word)
        }
    }
    
    return found
}

func main() {
    let dictionary = ["programming", "program", "code", "development", "programmer"]
    let search = "program"
    
    let results = findWord(dictionary, search)
    
    print("Searching for words containing:", search)
    print("Results found:", results)
    print("Total found:", results.length())
}
```

## Basic Error Handling

### 1. Parameter Validation

```r2
func divide(a, b) {
    if (b == 0) {
        print("Error: Cannot divide by zero")
        return null
    }
    return a / b
}

func main() {
    let result1 = divide(10, 2)
    let result2 = divide(10, 0)
    
    if (result1 != null) {
        print("10 ÷ 2 =", result1)
    }
    
    if (result2 != null) {
        print("10 ÷ 0 =", result2)
    } else {
        print("Division by zero is not valid")
    }
}
```

### 2. Array Validation

```r2
func getElement(array, index) {
    if (array.length() == 0) {
        print("Error: Array is empty")
        return null
    }
    
    if (index < 0 || index >= array.length()) {
        print("Error: Index out of range")
        return null
    }
    
    return array[index]
}

func main() {
    let numbers = [10, 20, 30]
    let empty = []
    
    print("Valid element:", getElement(numbers, 1))
    print("Invalid index:", getElement(numbers, 5))
    print("Empty array:", getElement(empty, 0))
}
```

## Module Project: Student Management System

```r2
// Simple student management system

func createStudent(name, age, grades) {
    return {
        name: name,
        age: age,
        grades: grades
    }
}

func calculateAverage(grades) {
    if (grades.length() == 0) {
        return 0
    }
    
    let sum = 0
    for (let grade in grades) {
        sum = sum + grade
    }
    
    return sum / grades.length()
}

func getGrade(average) {
    if (average >= 90) {
        return "A"
    } else if (average >= 80) {
        return "B"
    } else if (average >= 70) {
        return "C"
    } else if (average >= 60) {
        return "D"
    } else {
        return "F"
    }
}

func showStudent(student) {
    let average = calculateAverage(student.grades)
    let grade = getGrade(average)
    
    print("=== STUDENT ===")
    print("Name:", student.name)
    print("Age:", student.age)
    print("Grades:", student.grades)
    print("Average:", average)
    print("Grade:", grade)
    print()
}

func findBestStudent(students) {
    if (students.length() == 0) {
        return null
    }
    
    let best = students[0]
    let bestAverage = calculateAverage(best.grades)
    
    for (let i = 1; i < students.length(); i++) {
        let current = students[i]
        let currentAverage = calculateAverage(current.grades)
        
        if (currentAverage > bestAverage) {
            best = current
            bestAverage = currentAverage
        }
    }
    
    return best
}

func main() {
    // Create students
    let student1 = createStudent("Ana García", 20, [85, 92, 88, 90])
    let student2 = createStudent("Carlos López", 19, [78, 85, 82, 89])
    let student3 = createStudent("María Rodríguez", 21, [95, 98, 92, 96])
    
    let students = [student1, student2, student3]
    
    // Show all students
    print("STUDENT REPORT")
    print("======================")
    
    for (let student in students) {
        showStudent(student)
    }
    
    // Find the best student
    let best = findBestStudent(students)
    if (best != null) {
        print("🏆 BEST STUDENT:")
        showStudent(best)
    }
    
    // General statistics
    let totalStudents = students.length()
    let sumOfAverages = 0
    
    for (let student in students) {
        sumOfAverages = sumOfAverages + calculateAverage(student.grades)
    }
    
    let generalAverage = sumOfAverages / totalStudents
    
    print("GENERAL STATISTICS:")
    print("Total students:", totalStudents)
    print("General average:", generalAverage)
    print("General grade:", getGrade(generalAverage))
}
```

## Patterns and Best Practices

### 1. Small and Specific Functions

```r2
// ❌ Function that does too much
func processData(data) {
    // Validate, process, calculate, format, print...
    // 50 lines of code
}

// ✅ Specific functions
func validateData(data) {
    return data != null && data.length() > 0
}

func calculateStatistics(data) {
    // Only calculate
}

func formatResults(statistics) {
    // Only format
}
```

### 2. Descriptive Names

```r2
// ❌ Unclear names
func calc(x, y) {
    return x * y * 0.15
}

// ✅ Descriptive names
func calculateTax(price, quantity) {
    let tax = 0.15
    return price * quantity * tax
}
```

### 3. Input Validation

```r2
func createPerson(name, age) {
    // Validate parameters
    if (name == null || name == "") {
        print("Error: Name is required")
        return null
    }
    
    if (age < 0 || age > 150) {
        print("Error: Invalid age")
        return null
    }
    
    return {
        name: name,
        age: age
    }
}
```

## Module Summary

### Concepts Learned
- ✅ Conditionals (if/else/else if)
- ✅ Loops (while, for, for-in)
- ✅ Loop control (break, continue)
- ✅ Function definition and usage
- ✅ Parameters and return values
- ✅ Variable scope
- ✅ Arrays and their basic methods
- ✅ Anonymous functions
- ✅ Basic error handling

### Skills Developed
- ✅ Create programs with conditional logic
- ✅ Implement efficient loops
- ✅ Write reusable functions
- ✅ Work with collections of data
- ✅ Validate input and handle errors
- ✅ Organize code into small functions

### Next Module

In **Module 3**, you will learn:
- Object-oriented programming (classes and objects)
- Inheritance and polymorphism
- Maps and advanced objects
- Basic file handling

Excellent work completing Module 2! You now have the fundamental tools to create structured and functional programs.
