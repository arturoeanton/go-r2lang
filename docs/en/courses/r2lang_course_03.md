# R2Lang Course - Module 3: Object-Oriented Programming

## Introduction to Object-Oriented Programming

Object-Oriented Programming (OOP) is a paradigm that organizes code into **objects** containing data (properties) and code (methods). R2Lang supports OOP with classes, inheritance, and encapsulation.

### Fundamental Concepts

- **Class**: A template or blueprint for creating objects
- **Object**: An instance of a class
- **Properties**: Variables that belong to an object
- **Methods**: Functions that belong to an object
- **Constructor**: A special method that initializes an object
- **Inheritance**: The ability of a class to inherit from another

## Basic Classes and Objects

### 1. Defining a Simple Class

```r2
class Person {
    // Properties (declared with let)
    let name
    let age
    let email
    
    // Constructor: a special method that runs when an object is created
    constructor(name, age, email) {
        this.name = name
        this.age = age
        this.email = email
        print("New person created:", this.name)
    }
    
    // Methods: functions that belong to the class
    greet() {
        print("Hello, I am", this.name, "and I am", this.age, "years old")
    }
    
    birthday() {
        this.age++
        print(this.name, "is now", this.age, "years old")
    }
    
    changeEmail(newEmail) {
        let previousEmail = this.email
        this.email = newEmail
        print("Email of", this.name, "changed from", previousEmail, "to", newEmail)
    }
}

func main() {
    // Create objects (instantiate the class)
    let person1 = Person("Ana García", 25, "ana@email.com")
    let person2 = Person("Carlos López", 30, "carlos@email.com")
    
    // Use methods
    person1.greet()
    person2.greet()
    
    // Modify properties through methods
    person1.birthday()
    person2.changeEmail("carlos.lopez@gmail.com")
    
    // Access properties directly
    print("Name of person1:", person1.name)
    print("Age of person2:", person2.age)
}
```

### 2. Class with Advanced Methods

```r2
class BankAccount {
    let holder
    let accountNumber
    let balance
    let transactions
    
    constructor(holder, accountNumber, initialBalance) {
        this.holder = holder
        this.accountNumber = accountNumber
        this.balance = initialBalance
        this.transactions = []
        
        // Record initial transaction
        let initialTransaction = {
            type: "Opening",
            amount: initialBalance,
            date: "Today",
            resultingBalance: initialBalance
        }
        this.transactions = this.transactions.push(initialTransaction)
    }
    
    deposit(amount) {
        if (amount <= 0) {
            print("Error: The amount must be positive")
            return false
        }
        
        this.balance = this.balance + amount
        let transaction = {
            type: "Deposit",
            amount: amount,
            date: "Today",
            resultingBalance: this.balance
        }
        this.transactions = this.transactions.push(transaction)
        
        print("Deposit successful. New balance:", this.balance)
        return true
    }
    
    withdraw(amount) {
        if (amount <= 0) {
            print("Error: The amount must be positive")
            return false
        }
        
        if (amount > this.balance) {
            print("Error: Insufficient balance")
            return false
        }
        
        this.balance = this.balance - amount
        let transaction = {
            type: "Withdrawal",
            amount: amount,
            date: "Today",
            resultingBalance: this.balance
        }
        this.transactions = this.transactions.push(transaction)
        
        print("Withdrawal successful. New balance:", this.balance)
        return true
    }
    
    checkBalance() {
        print("Current balance of", this.holder + ":", this.balance)
        return this.balance
    }
    
    showTransactions() {
        print("=== TRANSACTIONS OF", this.holder, "===")
        print("Account:", this.accountNumber)
        
        for (let i = 0; i < this.transactions.length(); i++) {
            let mov = this.transactions[i]
            print("- " + mov.type + ": $" + mov.amount + " (Balance: $" + mov.resultingBalance + ")")
        }
        print("Total transactions:", this.transactions.length())
    }
    
    transfer(destinationAccount, amount) {
        if (this.withdraw(amount)) {
            if (destinationAccount.deposit(amount)) {
                print("Transfer successful from", this.holder, "to", destinationAccount.holder)
                return true
            } else {
                // Revert the withdrawal if the deposit fails
                this.deposit(amount)
                print("Error in transfer: deposit failed")
                return false
            }
        }
        return false
    }
}

func main() {
    // Create bank accounts
    let anaAccount = BankAccount("Ana García", "12345", 1000)
    let carlosAccount = BankAccount("Carlos López", "67890", 500)
    
    print()
    
    // Basic operations
    anaAccount.checkBalance()
    anaAccount.deposit(200)
    anaAccount.withdraw(150)
    
    print()
    
    // Transfer between accounts
    anaAccount.transfer(carlosAccount, 300)
    
    print()
    
    // Show transactions
    anaAccount.showTransactions()
    print()
    carlosAccount.showTransactions()
}
```

## Inheritance

### 1. Basic Inheritance

```r2
// Parent class (superclass)
class Animal {
    let name
    let species
    let age
    
    constructor(name, species, age) {
        this.name = name
        this.species = species
        this.age = age
    }
    
    makeSound() {
        print(this.name, "makes a sound")
    }
    
    sleep() {
        print(this.name, "is sleeping")
    }
    
    eat() {
        print(this.name, "is eating")
    }
    
    showInfo() {
        print("Name:", this.name)
        print("Species:", this.species)
        print("Age:", this.age, "years")
    }
}

// Child class (subclass)
class Dog extends Animal {
    let breed
    let owner
    
    constructor(name, age, breed, owner) {
        // Call the parent class constructor
        super.constructor(name, "Canine", age)
        this.breed = breed
        this.owner = owner
    }
    
    // Override parent method
    makeSound() {
        print(this.name, "barks: Woof woof!")
    }
    
    // Dog-specific methods
    fetchBall() {
        print(this.name, "is fetching the ball")
    }
    
    wagTail() {
        print(this.name, "wags its tail happily")
    }
    
    // Override method to show specific information
    showInfo() {
        super.showInfo()  // Call parent method
        print("Breed:", this.breed)
        print("Owner:", this.owner)
    }
}

class Cat extends Animal {
    let furColor
    let isIndependent
    
    constructor(name, age, furColor) {
        super.constructor(name, "Feline", age)
        this.furColor = furColor
        this.isIndependent = true
    }
    
    makeSound() {
        print(this.name, "meows: Meow!")
    }
    
    purr() {
        print(this.name, "purrs contentedly")
    }
    
    sharpenClaws() {
        print(this.name, "is sharpening its claws")
    }
    
    showInfo() {
        super.showInfo()
        print("Fur color:", this.furColor)
        print("Is independent:", this.isIndependent)
    }
}

func main() {
    // Create instances
    let myDog = Dog("Max", 3, "Labrador", "Juan")
    let myCat = Cat("Luna", 2, "Gray")
    
    print("=== PET INFORMATION ===")
    print()
    
    print("🐕 DOG:")
    myDog.showInfo()
    print()
    
    print("🐱 CAT:")
    myCat.showInfo()
    print()
    
    print("=== ACTIONS ===")
    // Inherited methods
    myDog.eat()
    myCat.sleep()
    
    // Overridden methods
    myDog.makeSound()
    myCat.makeSound()
    
    // Specific methods
    myDog.fetchBall()
    myDog.wagTail()
    
    myCat.purr()
    myCat.sharpenClaws()
}
```

### 2. Multi-level Inheritance

```r2
// Base class
class Vehicle {
    let brand
    let model
    let year
    let maxSpeed
    
    constructor(brand, model, year, maxSpeed) {
        this.brand = brand
        this.model = model
        this.year = year
        this.maxSpeed = maxSpeed
    }
    
    accelerate() {
        print(this.brand, this.model, "is accelerating")
    }
    
    brake() {
        print(this.brand, this.model, "is braking")
    }
    
    showInfo() {
        print("Vehicle:", this.brand, this.model, "(" + this.year + ")")
        print("Max speed:", this.maxSpeed, "km/h")
    }
}

// Intermediate class
class Car extends Vehicle {
    let numDoors
    let fuelType
    
    constructor(brand, model, year, maxSpeed, numDoors, fuelType) {
        super.constructor(brand, model, year, maxSpeed)
        this.numDoors = numDoors
        this.fuelType = fuelType
    }
    
    startEngine() {
        print("Engine of", this.brand, this.model, "started")
    }
    
    showInfo() {
        super.showInfo()
        print("Doors:", this.numDoors)
        print("Fuel:", this.fuelType)
    }
}

// Specific class
class ElectricCar extends Car {
    let batteryCapacity
    let range
    let chargeLevel
    
    constructor(brand, model, year, maxSpeed, numDoors, batteryCapacity, range) {
        super.constructor(brand, model, year, maxSpeed, numDoors, "Electric")
        this.batteryCapacity = batteryCapacity
        this.range = range
        this.chargeLevel = 100  // Starts with full charge
    }
    
    chargeBattery(percentage) {
        this.chargeLevel = this.chargeLevel + percentage
        if (this.chargeLevel > 100) {
            this.chargeLevel = 100
        }
        print("Battery charged to", this.chargeLevel + "%")
    }
    
    checkRange() {
        let currentRange = (this.range * this.chargeLevel) / 100
        print("Current range:", currentRange, "km")
        return currentRange
    }
    
    ecoMode() {
        print("Activating eco mode for", this.brand, this.model)
        print("Speed limited and consumption optimized")
    }
    
    showInfo() {
        super.showInfo()
        print("Battery capacity:", this.batteryCapacity, "kWh")
        print("Maximum range:", this.range, "km")
        print("Charge level:", this.chargeLevel + "%")
    }
}

func main() {
    let tesla = ElectricCar("Tesla", "Model 3", 2023, 250, 4, 75, 500)
    
    print("=== ELECTRIC CAR ===")
    tesla.showInfo()
    print()
    
    print("=== ACTIONS ===")
    tesla.startEngine()
    tesla.accelerate()
    tesla.checkRange()
    tesla.ecoMode()
    tesla.chargeBattery(20)  // Try to charge more
    tesla.brake()
}
```

## Maps and Advanced Objects

### 1. Maps (Dictionaries)

```r2
func main() {
    // Create an empty map
    let person = {}
    
    // Create a map with initial data
    let student = {
        name: "Ana García",
        age: 20,
        major: "Software Engineering",
        semester: 5,
        subjects: ["Programming", "Mathematics", "Physics"]
    }
    
    print("Student:", student.name)
    print("Major:", student.major)
    print("Subjects:", student.subjects)
    
    // Add new properties
    student.average = 8.5
    student.active = true
    
    // Modify existing properties
    student.semester = 6
    
    print("Average:", student.average)
    print("Current semester:", student.semester)
}
```

### 2. Dynamic Maps

```r2
func createInventory() {
    let inventory = {}
    
    // Function to add a product
    let addProduct = func(code, name, price, quantity) {
        inventory[code] = {
            name: name,
            price: price,
            quantity: quantity,
            totalValue: price * quantity
        }
        print("Product added:", name)
    }
    
    // Function to show inventory
    let showInventory = func() {
        print("=== INVENTORY ===")
        let total = 0
        
        // We can't iterate over maps directly in R2Lang yet
        // We will simulate with known codes
        if (inventory["001"] != null) {
            let prod = inventory["001"]
            print("001 -", prod.name, "- Price:", prod.price, "- Quantity:", prod.quantity)
            total = total + prod.totalValue
        }
        
        if (inventory["002"] != null) {
            let prod = inventory["002"]
            print("002 -", prod.name, "- Price:", prod.price, "- Quantity:", prod.quantity)
            total = total + prod.totalValue
        }
        
        if (inventory["003"] != null) {
            let prod = inventory["003"]
            print("003 -", prod.name, "- Price:", prod.price, "- Quantity:", prod.quantity)
            total = total + prod.totalValue
        }
        
        print("Total inventory value:", total)
    }
    
    return {
        add: addProduct,
        show: showInventory,
        get: func(code) {
            return inventory[code]
        }
    }
}

func main() {
    let inv = createInventory()
    
    // Add products
    inv.add("001", "Laptop", 1500, 10)
    inv.add("002", "Mouse", 25, 50)
    inv.add("003", "Keyboard", 75, 30)
    
    print()
    inv.show()
    
    print()
    let laptop = inv.get("001")
    print("Laptop details:", laptop.name, "- Stock:", laptop.quantity)
}
```

### 3. Object Composition

```r2
class Engine {
    let cylinders
    let horsepower
    let fuel
    let isOn
    
    constructor(cylinders, horsepower, fuel) {
        this.cylinders = cylinders
        this.horsepower = horsepower
        this.fuel = fuel
        this.isOn = false
    }
    
    start() {
        if (!this.isOn) {
            this.isOn = true
            print("Engine started -", this.horsepower, "HP")
        } else {
            print("Engine is already on")
        }
    }
    
    stop() {
        if (this.isOn) {
            this.isOn = false
            print("Engine stopped")
        } else {
            print("Engine is already off")
        }
    }
    
    getInfo() {
        return {
            cylinders: this.cylinders,
            horsepower: this.horsepower,
            fuel: this.fuel,
            status: this.isOn ? "On" : "Off"
        }
    }
}

class System {
    let name
    let isActive
    
    constructor(name) {
        this.name = name
        this.isActive = false
    }
    
    activate() {
        this.isActive = true
        print("System", this.name, "activated")
    }
    
    deactivate() {
        this.isActive = false
        print("System", this.name, "deactivated")
    }
}

class Car {
    let brand
    let model
    let engine
    let systems
    
    constructor(brand, model, engine) {
        this.brand = brand
        this.model = model
        this.engine = engine
        this.systems = {
            airConditioning: System("Air Conditioning"),
            gpsNavigation: System("GPS Navigation"),
            soundSystem: System("Sound System")
        }
    }
    
    start() {
        print("Starting", this.brand, this.model)
        this.engine.start()
        
        // Activate basic systems
        this.systems.gpsNavigation.activate()
        print("Car ready to drive")
    }
    
    stop() {
        print("Stopping", this.brand, this.model)
        
        // Deactivate systems
        this.systems.airConditioning.deactivate()
        this.systems.gpsNavigation.deactivate()
        this.systems.soundSystem.deactivate()
        
        this.engine.stop()
        print("Car completely off")
    }
    
    activateAirConditioning() {
        this.systems.airConditioning.activate()
    }
    
    activateSoundSystem() {
        this.systems.soundSystem.activate()
    }
    
    showStatus() {
        print("=== STATUS OF", this.brand, this.model, "===")
        
        let engineInfo = this.engine.getInfo()
        print("Engine:", engineInfo.horsepower, "HP -", engineInfo.status)
        
        print("Air conditioning:", this.systems.airConditioning.isActive ? "ON" : "OFF")
        print("Navigation:", this.systems.gpsNavigation.isActive ? "ON" : "OFF")
        print("Sound:", this.systems.soundSystem.isActive ? "ON" : "OFF")
    }
}

func main() {
    // Create engine
    let v6Engine = Engine(6, 300, "Gasoline")
    
    // Create car with composition
    let myCar = Car("Toyota", "Camry", v6Engine)
    
    print("=== CAR TEST ===")
    myCar.showStatus()
    print()
    
    myCar.start()
    print()
    
    myCar.activateAirConditioning()
    myCar.activateSoundSystem()
    print()
    
    myCar.showStatus()
    print()
    
    myCar.stop()
}
```

## Practical Exercises

### Exercise 1: Library System

```r2
class Book {
    let title
    let author
    let isbn
    let isAvailable
    let loanDate
    
    constructor(title, author, isbn) {
        this.title = title
        this.author = author
        this.isbn = isbn
        this.isAvailable = true
        this.loanDate = null
    }
    
    loan() {
        if (this.isAvailable) {
            this.isAvailable = false
            this.loanDate = "Today"
            print("Book", this.title, "loaned successfully")
            return true
        } else {
            print("Book", this.title, "is not available")
            return false
        }
    }
    
    returnBook() {
        if (!this.isAvailable) {
            this.isAvailable = true
            this.loanDate = null
            print("Book", this.title, "returned successfully")
            return true
        } else {
            print("Book", this.title, "is already available")
            return false
        }
    }
    
    showInfo() {
        print("📖", this.title, "by", this.author)
        print("   ISBN:", this.isbn)
        print("   Status:", this.isAvailable ? "Available" : "Loaned")
        if (!this.isAvailable) {
            print("   Loaned since:", this.loanDate)
        }
    }
}

class Library {
    let name
    let books
    
    constructor(name) {
        this.name = name
        this.books = []
    }
    
    addBook(book) {
        this.books = this.books.push(book)
        print("Book added to", this.name)
    }
    
    findByTitle(title) {
        for (let i = 0; i < this.books.length(); i++) {
            let book = this.books[i]
            if (book.title.contains(title)) {
                return book
            }
        }
        return null
    }
    
    showCatalog() {
        print("=== CATALOG OF", this.name, "===")
        print("Total books:", this.books.length())
        print()
        
        for (let i = 0; i < this.books.length(); i++) {
            this.books[i].showInfo()
            print()
        }
    }
    
    availableBooks() {
        let available = []
        for (let i = 0; i < this.books.length(); i++) {
            let book = this.books[i]
            if (book.isAvailable) {
                available = available.push(book)
            }
        }
        return available
    }
}

func main() {
    let library = Library("Central Library")
    
    // Create books
    let book1 = Book("Don Quixote", "Cervantes", "123456")
    let book2 = Book("One Hundred Years of Solitude", "García Márquez", "789012")
    let book3 = Book("1984", "George Orwell", "345678")
    
    // Add to library
    library.addBook(book1)
    library.addBook(book2)
    library.addBook(book3)
    
    print()
    library.showCatalog()
    
    // Make loans
    print("=== LOANS ===")
    book1.loan()
    book2.loan()
    
    print()
    print("Available books:", library.availableBooks().length())
    
    // Return
    print()
    print("=== RETURNS ===")
    book1.returnBook()
    
    print()
    let foundBook = library.findByTitle("1984")
    if (foundBook != null) {
        print("Book found:")
        foundBook.showInfo()
    }
}
```

## Module Project: School Management System

```r2
// Base class for people
class Person {
    let name
    let age
    let id
    
    constructor(name, age, id) {
        this.name = name
        this.age = age
        this.id = id
    }
    
    showInfo() {
        print("Name:", this.name)
        print("Age:", this.age)
        print("ID:", this.id)
    }
}

// Student inherits from Person
class Student extends Person {
    let grade
    let grades
    let enrolledSubjects
    
    constructor(name, age, id, grade) {
        super.constructor(name, age, id)
        this.grade = grade
        this.grades = {}
        this.enrolledSubjects = []
    }
    
    enrollSubject(subject) {
        this.enrolledSubjects = this.enrolledSubjects.push(subject)
        this.grades[subject] = []
        print(this.name, "enrolled in", subject)
    }
    
    addGrade(subject, grade) {
        if (this.grades[subject] != null) {
            let grades = this.grades[subject]
            grades = grades.push(grade)
            this.grades[subject] = grades
            print("Grade", grade, "added in", subject, "for", this.name)
        } else {
            print("Error:", this.name, "is not enrolled in", subject)
        }
    }
    
    calculateAverage(subject) {
        if (this.grades[subject] != null) {
            let grades = this.grades[subject]
            if (grades.length() == 0) {
                return 0
            }
            
            let sum = 0
            for (let i = 0; i < grades.length(); i++) {
                sum = sum + grades[i]
            }
            return sum / grades.length()
        }
        return 0
    }
    
    showInfo() {
        super.showInfo()
        print("Grade:", this.grade)
        print("Enrolled subjects:", this.enrolledSubjects.length())
        
        for (let i = 0; i < this.enrolledSubjects.length(); i++) {
            let subject = this.enrolledSubjects[i]
            let average = this.calculateAverage(subject)
            print("-", subject + ":", average)
        }
    }
}

// Teacher inherits from Person
class Teacher extends Person {
    let specialtySubject
    let students
    let salary
    
    constructor(name, age, id, specialtySubject, salary) {
        super.constructor(name, age, id)
        this.specialtySubject = specialtySubject
        this.students = []
        this.salary = salary
    }
    
    assignStudent(student) {
        this.students = this.students.push(student)
        print("Student", student.name, "assigned to teacher", this.name)
    }
    
    gradeStudent(student, grade) {
        student.addGrade(this.specialtySubject, grade)
    }
    
    showStudents() {
        print("Students of teacher", this.name, "(" + this.specialtySubject + "):")
        for (let i = 0; i < this.students.length(); i++) {
            let st = this.students[i]
            let average = st.calculateAverage(this.specialtySubject)
            print("-", st.name, "(Average:", average + ")")
        }
    }
    
    showInfo() {
        super.showInfo()
        print("Specialty:", this.specialtySubject)
        print("Salary:", this.salary)
        print("Students in charge:", this.students.length())
    }
}

// School as a container
class School {
    let name
    let students
    let teachers
    
    constructor(name) {
        this.name = name
        this.students = []
        this.teachers = []
    }
    
    addStudent(student) {
        this.students = this.students.push(student)
        print("Student", student.name, "added to", this.name)
    }
    
    addTeacher(teacher) {
        this.teachers = this.teachers.push(teacher)
        print("Teacher", teacher.name, "added to", this.name)
    }
    
    showSummary() {
        print("=== SCHOOL", this.name, "===")
        print("Total students:", this.students.length())
        print("Total teachers:", this.teachers.length())
        print()
        
        print("STUDENTS:")
        for (let i = 0; i < this.students.length(); i++) {
            let st = this.students[i]
            print("-", st.name, "(" + st.grade + "th grade)")
        }
        print()
        
        print("TEACHERS:")
        for (let i = 0; i < this.teachers.length(); i++) {
            let prof = this.teachers[i]
            print("-", prof.name, "(" + prof.specialtySubject + ")")
        }
    }
}

func main() {
    // Create school
    let school = School("Technological Institute")
    
    // Create students
    let ana = Student("Ana García", 16, "STU001", 10)
    let carlos = Student("Carlos López", 15, "STU002", 9)
    let maria = Student("María Rodríguez", 17, "STU003", 11)
    
    // Create teachers
    let profMath = Teacher("Dr. González", 45, "PROF001", "Mathematics", 3500)
    let profPhysics = Teacher("Dra. Martínez", 38, "PROF002", "Physics", 3200)
    
    // Add to school
    school.addStudent(ana)
    school.addStudent(carlos)
    school.addStudent(maria)
    school.addTeacher(profMath)
    school.addTeacher(profPhysics)
    
    print()
    
    // Enroll students in subjects
    ana.enrollSubject("Mathematics")
    ana.enrollSubject("Physics")
    carlos.enrollSubject("Mathematics")
    maria.enrollSubject("Physics")
    
    print()
    
    // Assign students to teachers
    profMath.assignStudent(ana)
    profMath.assignStudent(carlos)
    profPhysics.assignStudent(ana)
    profPhysics.assignStudent(maria)
    
    print()
    
    // Grade students
    profMath.gradeStudent(ana, 95)
    profMath.gradeStudent(ana, 88)
    profMath.gradeStudent(carlos, 78)
    profMath.gradeStudent(carlos, 85)
    
    profPhysics.gradeStudent(ana, 92)
    profPhysics.gradeStudent(maria, 89)
    
    print()
    
    // Show information
    school.showSummary()
    print()
    
    print("=== STUDENT DETAILS ===")
    ana.showInfo()
    print()
    carlos.showInfo()
    print()
    
    print("=== STUDENTS BY TEACHER ===")
    profMath.showStudents()
    print()
    profPhysics.showStudents()
}
```

## Design Patterns in OOP

### 1. Factory Pattern

```r2
class VehicleFactory {
    createVehicle(type, brand, model) {
        if (type == "car") {
            return Car(brand, model, 4)
        } else if (type == "motorcycle") {
            return Motorcycle(brand, model, 2)
        } else if (type == "truck") {
            return Truck(brand, model, 6)
        } else {
            print("Vehicle type not supported")
            return null
        }
    }
}
```

### 2. Observer Pattern (Simulated)

```r2
class Notifier {
    let observers
    
    constructor() {
        this.observers = []
    }
    
    subscribe(observer) {
        this.observers = this.observers.push(observer)
    }
    
    notify(message) {
        for (let i = 0; i < this.observers.length(); i++) {
            this.observers[i].update(message)
        }
    }
}
```

## Module Summary

### Concepts Learned
- ✅ Basic classes and objects
- ✅ Constructors and methods
- ✅ Class properties
- ✅ Inheritance with `extends`
- ✅ Using `super` to call parent methods
- ✅ Method overriding
- ✅ Maps and dynamic objects
- ✅ Object composition
- ✅ Basic design patterns

### Skills Developed
- ✅ Designing effective classes
- ✅ Implementing appropriate inheritance
- ✅ Creating object hierarchies
- ✅ Using composition vs. inheritance
- ✅ Managing collections of objects
- ✅ Applying basic OOP principles

### Next Module

In **Module 4**, you will learn:
- Concurrency and parallel programming
- Advanced error handling (try/catch/finally)
- Working with files
- Advanced built-in libraries

Excellent work! You now master the fundamental concepts of object-oriented programming in R2Lang.
