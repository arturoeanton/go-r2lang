func main() {
    print("=== R2Lang For-In Loop Examples ===")
    print("Demonstrating multiline map literals and for-in loops")
    
    // 1. Classic array iteration (existing functionality)
    print("\n--- Array Iteration ---")
    let numbers = [10, 20, 30, 40, 50]
    print("Numbers array:", numbers)
    print("Array length:", len(numbers))
    
    for (i in numbers) {
        print("Index:", $k + " -> Value:", $v)
    }
    
    // 2. Map iteration with MULTILINE map literals (NEW!)
    print("\n--- Map Iteration with Multiline Syntax ---")
    let user = {
        name: "Alice",
        age: 28,
        city: "Barcelona",
        active: true
    }
    print("User map:", user)
    print("Map length:", len(user))
    
    let userKeys = keys(user)
    print("Using keys() function to iterate:")
    for (i in userKeys) {
        let key = userKeys[$k]
        print("Key:", key + " -> Value:", user[key])
    }
    
    // 3. Using $k and $v variables directly (with arrays)
    print("\n--- Using $k and $v Variables ---")
    let configArray = ["host", "port", "ssl", "timeout"]
    
    for (i in configArray) {
        print("Property $k=" + $k + " has value $v=" + $v)
    }
    
    // 4. Complex map with different value types
    print("\n--- Complex Map Iteration ---")
    let mixed = {
        string: "Hello World",
        number: 42,
        boolean: true,
        negative: false,
        zero: 0
    }
    
    let mixedKeys = keys(mixed)
    for (i in mixedKeys) {
        let key = mixedKeys[$k]
        print("Key '" + key + "' (" + typeOf(mixed[key]) + "): " + mixed[key])
    }
    
    // 5. Simplified nested map demo
    print("\n--- Nested Map Iteration ---")
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
    
    print("Company:", company.name)
    let empKeys = keys(company.employees)
    for (i in empKeys) {
        let empId = empKeys[$k]
        let emp = company.employees[empId]
        print("Employee " + empId + ": " + emp.position + " (salary: " + emp.salary + ")")
    }
    
    // 6. Using keys() function with for-in
    print("\n--- Using keys() Function ---")
    let scores = {
        alice: 95,
        bob: 87,
        charlie: 92,
        diana: 88
    }
    let studentKeys = keys(scores)
    
    print("Students:", studentKeys)
    print("Total students:", len(studentKeys))
    
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
        
        print("Student " + student + ": " + score + " points -> Grade " + grade)
    }
    
    // 7. Filtering during iteration
    print("\n--- Filtering During Iteration ---")
    let products = {
        laptop: 1200,
        mouse: 25,
        keyboard: 80,
        monitor: 300,
        cable: 15
    }
    
    print("Products over $50:")
    let productKeys = keys(products)
    for (i in productKeys) {
        let product = productKeys[$k]
        let price = products[product]
        
        // Demonstrating 'else if' for price categorization
        if (price > 1000) {
            print("- " + product + ": $" + price + " (PREMIUM)")
        } else if (price > 100) {
            print("- " + product + ": $" + price + " (EXPENSIVE)")
        } else if (price > 50) {
            print("- " + product + ": $" + price + " (MODERATE)")
        } else {
            print("- " + product + ": $" + price + " (BUDGET)")
        }
    }
    
    // 8. Building arrays from map iteration
    print("\n--- Building Arrays from Maps ---")
    let inventory = {
        apples: 50,
        bananas: 30,
        oranges: 25,
        grapes: 40
    }
    let lowStockCount = 0
    
    let inventoryKeys = keys(inventory)
    print("All products:", inventoryKeys)
    
    for (i in inventoryKeys) {
        let item = inventoryKeys[$k]
        let stock = inventory[item]
        
        // Stock level analysis using 'else if'
        if (stock >= 50) {
            print("✓ " + item + ": " + stock + " (Good stock)")
        } else if (stock >= 35) {
            print("⚠ " + item + ": " + stock + " (Medium stock)")
            lowStockCount = lowStockCount + 1
        } else if (stock >= 20) {
            print("⚠ " + item + ": " + stock + " (Low stock)")
            lowStockCount = lowStockCount + 1
        } else {
            print("❌ " + item + ": " + stock + " (Critical stock!)")
            lowStockCount = lowStockCount + 1
        }
    }
    
    print("Low stock count:", lowStockCount)
    
    // 9. Break and continue in array iteration
    print("\n--- Break and Continue Examples ---")
    let dataArray = [1, 2, 3, 4, 5, 6]
    
    print("Finding first value > 3:")
    for (i in dataArray) {
        if (dataArray[$k] > 3) {
            print("Found at index " + $k + ": " + dataArray[$k])
            break
        }
        print("Checking index " + $k + ": " + dataArray[$k])
    }
    
    print("Skipping even values:")
    for (i in dataArray) {
        if (dataArray[$k] % 2 == 0) {
            continue
        }
        print("Odd value at " + $k + ": " + dataArray[$k])
    }
    
    // 10. String array iteration (comparison)
    print("\n--- String Array vs Map Comparison ---")
    let strArray = ["first", "second", "third"]
    let strMap = {
        0: "first",
        1: "second", 
        2: "third"
    }
    
    print("Array iteration:")
    for (i in strArray) {
        print("Index " + $k + ": " + strArray[$k])
    }
    
    print("Map iteration:")
    let strMapKeys = keys(strMap)
    for (i in strMapKeys) {
        let key = strMapKeys[$k]
        print("Key " + key + ": " + strMap[key])
    }
    
    // 11. Performance comparison demo
    print("\n--- Performance Demo ---")
    let largeMap = {}
    
    // Build a larger map
    for (let i = 0; i < 10; i++) {
        largeMap["key" + i] = "value" + i
    }
    
    print("Large map size:", len(largeMap))
    print("Large map keys:", keys(largeMap))
    
    let count = 0
    let largeMapKeys = keys(largeMap)
    for (i in largeMapKeys) {
        count = count + 1
    }
    print("Iterated through " + count + " items")
    
    // 12. Advanced multiline map with nested structures and mixed separators
    print("\n--- Advanced Multiline Map Example ---")
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
    
    print("Complex config structure created with multiline syntax!")
    print("Server config:", complexConfig.server)
    print("Database type:", complexConfig.database.type)
    print("DB credentials timeout:", complexConfig.database.credentials.timeout)
    
    // Demonstrate iteration through nested structures
    let mainKeys = keys(complexConfig)
    print("\nMain configuration sections:", mainKeys)
    
    for (i in mainKeys) {
        let section = mainKeys[$k]
        let propCount = len(keys(complexConfig[section]))
        
        // Demonstrating 'else if' with modulo operator '%' (NEW!)
        if (propCount % 4 == 0) {
            print("Section '" + section + "' has " + propCount + " properties (divisible by 4)")
        } else if (propCount % 3 == 0) {
            print("Section '" + section + "' has " + propCount + " properties (divisible by 3)")
        } else if (propCount % 2 == 0) {
            print("Section '" + section + "' has " + propCount + " properties (even number)")
        } else {
            print("Section '" + section + "' has " + propCount + " properties (odd number)")
        }
    }
    
    // 13. Advanced 'else if' examples with modulo operator
    print("\n--- Advanced 'else if' and Modulo Examples ---")
    let testNumbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 15, 20]
    
    for (i in testNumbers) {
        let num = testNumbers[$k]
        
        // Complex 'else if' chain with modulo operations
        if (num % 15 == 0) {
            print("Number " + num + " is divisible by both 3 and 5 (FizzBuzz!)")
        } else if (num % 5 == 0) {
            print("Number " + num + " is divisible by 5 (Buzz)")
        } else if (num % 3 == 0) {
            print("Number " + num + " is divisible by 3 (Fizz)")
        } else if (num % 2 == 0) {
            print("Number " + num + " is even")
        } else {
            print("Number " + num + " is odd and not divisible by 3 or 5")
        }
    }

    print("\n=== R2Lang New Features Summary ===")
    print("Array iteration: ✓ (existing)")
    print("Map iteration: ✓ (NEW)")
    print("Multiline map literals: ✓ (NEW)")
    print("Mixed comma/newline separators: ✓ (NEW)")
    print("Nested multiline maps: ✓ (NEW)")
    print("'else if' syntax: ✓ (NEW)")
    print("Modulo operator '%': ✓ (NEW)")
    print("Complex 'else if' chains: ✓ (NEW)")
    print("$k and $v variables: ✓")
    print("Break/continue support: ✓")
    print("Nested iteration: ✓")
    print("keys() function integration: ✓")
}