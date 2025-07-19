func main() {
    std.print("=== R2Lang For-In Loop Examples ===")
    std.print("Demonstrating multiline map literals and for-in loops")
    
    // 1. Classic array iteration (existing functionality)
    std.print("\n--- Array Iteration ---")
    let numbers = [10, 20, 30, 40, 50]
    std.print("Numbers array:", numbers)
    std.print("Array length:", std.len(numbers))
    
    for (i in numbers) {
        std.print("Index:", $k + " -> Value:", $v)
    }
    
    // 2. Map iteration with MULTILINE map literals (NEW!)
    std.print("\n--- Map Iteration with Multiline Syntax ---")
    let user = {
        name: "Alice",
        age: 28,
        city: "Barcelona",
        active: true
    }
    std.print("User map:", user)
    std.print("Map length:", std.len(user))
    
    let userKeys = std.keys(user)
    std.print("Using std.keys() function to iterate:")
    for (i in userKeys) {
        let key = userKeys[$k]
        std.print("Key:", key + " -> Value:", user[key])
    }
    
    // 3. Using $k and $v variables directly (with arrays)
    std.print("\n--- Using $k and $v Variables ---")
    let configArray = ["host", "port", "ssl", "timeout"]
    
    for (i in configArray) {
        std.print("Property $k=" + $k + " has value $v=" + $v)
    }
    
    // 4. Complex map with different value types
    std.print("\n--- Complex Map Iteration ---")
    let mixed = {
        string: "Hello World",
        number: 42,
        boolean: true,
        negative: false,
        zero: 0
    }
    
    let mixedKeys = std.keys(mixed)
    for (i in mixedKeys) {
        let key = mixedKeys[$k]
        std.print("Key '" + key + "' (" + std.typeOf(mixed[key]) + "): " + mixed[key])
    }
    
    // 5. Simplified nested map demo
    std.print("\n--- Nested Map Iteration ---")
    let john = {
        position: "developer",
        salary: 50000
    }
    let jane = {
        position: "designer", 
        salary: 45000
    }
    let bob = {
        position: "manager",
        salary: 60000
    }
    let employees = {
        john: john,
        jane: jane,
        bob: bob
    }
    let company = {
        name: "TechCorp",
        employees: employees,
        founded: 2020,
        active: true
    }
    
    std.print("Company:", company.name)
    let empKeys = std.keys(company.employees)
    for (i in empKeys) {
        let empId = empKeys[$k]
        let emp = company.employees[empId]
        std.print("Employee " + empId + ": " + emp.position + " (salary: " + emp.salary + ")")
    }
    
    // 6. Using std.keys() function with for-in
    std.print("\n--- Using std.keys() Function ---")
    let scores = {
        alice: 95,
        bob: 87,
        charlie: 92,
        diana: 88
    }
    let studentKeys = std.keys(scores)
    
    std.print("Students:", studentKeys)
    std.print("Total students:", std.len(studentKeys))
    
    for (i in studentKeys) {
        let student = studentKeys[$k]
        let score = scores[student]
        let grade = "F"
        
        // Using new 'else if' syntax (improved readability!)
        if (score >= 90) {
            grade = "A"
        } else if (score >= 80) {
            grade = "B"
        } else if (score >= 70) {
            grade = "C"
        } else if (score >= 60) {
            grade = "D"
        } else {
            grade = "F"
        }
        
        std.print("Student " + student + ": " + score + " points -> Grade " + grade)
    }
    
    // 7. Filtering during iteration
    std.print("\n--- Filtering During Iteration ---")
    let products = {
        laptop: 1200,
        mouse: 25,
        keyboard: 80,
        monitor: 300,
        cable: 15
    }
    
    std.print("Products over $50:")
    let productKeys = std.keys(products)
    for (i in productKeys) {
        let product = productKeys[$k]
        let price = products[product]
        
        // Demonstrating 'else if' for price categorization
        if (price > 1000) {
            std.print("- " + product + ": $" + price + " (PREMIUM)")
        } else if (price > 100) {
            std.print("- " + product + ": $" + price + " (EXPENSIVE)")
        } else if (price > 50) {
            std.print("- " + product + ": $" + price + " (MODERATE)")
        } else {
            std.print("- " + product + ": $" + price + " (BUDGET)")
        }
    }
    
    // 8. Building arrays from map iteration
    std.print("\n--- Building Arrays from Maps ---")
    let inventory = {
        apples: 50,
        bananas: 30,
        oranges: 25,
        grapes: 40
    }
    let lowStockCount = 0
    
    let inventoryKeys = std.keys(inventory)
    std.print("All products:", inventoryKeys)
    
    for (i in inventoryKeys) {
        let item = inventoryKeys[$k]
        let stock = inventory[item]
        
        // Stock level analysis using 'else if'
        if (stock >= 50) {
            std.print("✓ " + item + ": " + stock + " (Good stock)")
        } else if (stock >= 35) {
            std.print("⚠ " + item + ": " + stock + " (Medium stock)")
            lowStockCount = lowStockCount + 1
        } else if (stock >= 20) {
            std.print("⚠ " + item + ": " + stock + " (Low stock)")
            lowStockCount = lowStockCount + 1
        } else {
            std.print("❌ " + item + ": " + stock + " (Critical stock!)")
            lowStockCount = lowStockCount + 1
        }
    }
    
    std.print("Low stock count:", lowStockCount)
    
    // 9. Break and continue in array iteration
    std.print("\n--- Break and Continue Examples ---")
    let dataArray = [1, 2, 3, 4, 5, 6]
    
    std.print("Finding first value > 3:")
    for (i in dataArray) {
        if (dataArray[$k] > 3) {
            std.print("Found at index " + $k + ": " + dataArray[$k])
            break
        }
        std.print("Checking index " + $k + ": " + dataArray[$k])
    }
    
    std.print("Skipping even values:")
    for (i in dataArray) {
        if (dataArray[$k] % 2 == 0) {
            continue
        }
        std.print("Odd value at " + $k + ": " + dataArray[$k])
    }
    
    // 10. String array iteration (comparison)
    std.print("\n--- String Array vs Map Comparison ---")
    let strArray = ["first", "second", "third"]
    let strMap = {
        0: "first",
        1: "second", 
        2: "third"
    }
    
    std.print("Array iteration:")
    for (i in strArray) {
        std.print("Index " + $k + ": " + strArray[$k])
    }
    
    std.print("Map iteration:")
    let strMapKeys = std.keys(strMap)
    for (i in strMapKeys) {
        let key = strMapKeys[$k]
        std.print("Key " + key + ": " + strMap[key])
    }
    
    // 11. Performance comparison demo
    std.print("\n--- Performance Demo ---")
    let largeMap = {}
    
    // Build a larger map
    for (let i = 0; i < 10; i++) {
        largeMap["key" + i] = "value" + i
    }
    
    std.print("Large map size:", std.len(largeMap))
    std.print("Large map keys:", std.keys(largeMap))
    
    let count = 0
    let largeMapKeys = std.keys(largeMap)
    for (i in largeMapKeys) {
        count = count + 1
    }
    std.print("Iterated through " + count + " items")
    
    // 12. Advanced multiline map with nested structures and mixed separators
    std.print("\n--- Advanced Multiline Map Example ---")
    let complexConfig = {
        server: {
            host: "localhost"
            port: 8080,
            ssl: true
        },
        database: {
            type: "postgresql",
            host: "db.example.com"
            port: 5432
            credentials: {
                username: "admin"
                password: "secret",
                timeout: 30
            }
        }
        features: {
            logging: true,
            caching: false
            monitoring: true,
            debug: false
        }
    }
    
    std.print("Complex config structure created with multiline syntax!")
    std.print("Server config:", complexConfig.server)
    std.print("Database type:", complexConfig.database.type)
    std.print("DB credentials timeout:", complexConfig.database.credentials.timeout)
    
    // Demonstrate iteration through nested structures
    let mainKeys = std.keys(complexConfig)
    std.print("\nMain configuration sections:", mainKeys)
    
    for (i in mainKeys) {
        let section = mainKeys[$k]
        let propCount = std.len(std.keys(complexConfig[section]))
        
        // Demonstrating 'else if' with modulo operator '%' (NEW!)
        if (propCount % 4 == 0) {
            std.print("Section '" + section + "' has " + propCount + " properties (divisible by 4)")
        } else if (propCount % 3 == 0) {
            std.print("Section '" + section + "' has " + propCount + " properties (divisible by 3)")
        } else if (propCount % 2 == 0) {
            std.print("Section '" + section + "' has " + propCount + " properties (even number)")
        } else {
            std.print("Section '" + section + "' has " + propCount + " properties (odd number)")
        }
    }
    
    // 13. Advanced 'else if' examples with modulo operator
    std.print("\n--- Advanced 'else if' and Modulo Examples ---")
    let testNumbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 15, 20]
    
    for (i in testNumbers) {
        let num = testNumbers[$k]
        
        // Complex 'else if' chain with modulo operations
        if (num % 15 == 0) {
            std.print("Number " + num + " is divisible by both 3 and 5 (FizzBuzz!)")
        } else if (num % 5 == 0) {
            std.print("Number " + num + " is divisible by 5 (Buzz)")
        } else if (num % 3 == 0) {
            std.print("Number " + num + " is divisible by 3 (Fizz)")
        } else if (num % 2 == 0) {
            std.print("Number " + num + " is even")
        } else {
            std.print("Number " + num + " is odd and not divisible by 3 or 5")
        }
    }

    std.print("\n=== R2Lang New Features Summary ===")
    std.print("Array iteration: ✓ (existing)")
    std.print("Map iteration: ✓ (NEW)")
    std.print("Multiline map literals: ✓ (NEW)")
    std.print("Mixed comma/newline separators: ✓ (NEW)")
    std.print("Nested multiline maps: ✓ (NEW)")
    std.print("'else if' syntax: ✓ (NEW)")
    std.print("Modulo operator '%': ✓ (NEW)")
    std.print("Complex 'else if' chains: ✓ (NEW)")
    std.print("$k and $v variables: ✓")
    std.print("Break/continue support: ✓")
    std.print("Nested iteration: ✓")
    std.print("std.keys() function integration: ✓")
}