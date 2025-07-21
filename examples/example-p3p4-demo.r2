// Demostración de características P3 y P4 de R2Lang
// Incluye: Optional Chaining, Null Coalescing, Pattern Matching, 
// Comprehensions y Pipeline Operator

print("🚀 R2Lang - Nuevas características P3 y P4")
print("=" * 50)

// ===== P3.1: OPTIONAL CHAINING ?. =====
print("\n🔗 P3.1: Optional Chaining ?.") 
let user = {
    name: "Alice",
    profile: {
        address: {
            street: "Main St"
        }
    }
}

let street = user?.profile?.address?.street
let missing = user?.profile?.phone?.number

print("Street:", street)    // "Main St"
print("Missing:", missing)  // nil

// ===== P3.2: NULL COALESCING ?? =====
print("\n❓ P3.2: Null Coalescing ??")
let config = nil
let defaultPort = 8080
let port = config?.port ?? defaultPort

print("Port:", port)        // 8080

// ===== P3.3: PATTERN MATCHING =====
print("\n🎯 P3.3: Pattern Matching")
func processStatus(status) {
    return match status {
        case 200 => "Success"
        case 404 => "Not Found"
        case 500 => "Server Error"
        case x => "Other: " + x
    }
}

print("Status 200:", processStatus(200))
print("Status 404:", processStatus(404))
print("Status 418:", processStatus(418))

// ===== P4.1: ARRAY COMPREHENSIONS =====
print("\n📦 P4.1: Array Comprehensions")
let numbers = [1, 2, 3, 4, 5]
let squares = [x * x for x in numbers]
let evenSquares = [x * x for x in numbers if x % 2 == 0]

print("Numbers:", numbers)
print("Squares:", squares)
print("Even squares:", evenSquares)

// ===== P4.2: OBJECT COMPREHENSIONS =====
print("\n🏗️  P4.2: Object Comprehensions")
let fruits = ["apple", "banana", "cherry"]
let lengths = {fruit: fruit.length for fruit in fruits}

print("Fruit lengths:", lengths)

// ===== P4.3: PIPELINE OPERATOR |> =====
print("\n🚀 P4.3: Pipeline Operator |>")
func double(x) { return x * 2 }
func addTen(x) { return x + 10 }

let result = 5 |> double |> addTen
let lambdaResult = 10 |> (x => x * 3)

print("5 |> double |> addTen =", result)       // 20
print("10 |> (x => x * 3) =", lambdaResult)   // 30

// ===== COMBINACIÓN DE CARACTERÍSTICAS =====
print("\n🤝 Combinando P3 y P4:")

// Optional chaining + null coalescing + comprehensions
let data = {
    users: [
        {name: "Ana", age: 25},
        {name: "Luis", age: 30},
        {name: "Carmen"}
    ]
}

let ages = [user?.age ?? 0 for user in data?.users ?? []]
print("Ages with defaults:", ages)

// Pattern matching + pipeline
func categorize(age) {
    return match age {
        case a if a < 18 => "Minor"
        case a if a < 65 => "Adult" 
        case _ => "Senior"
    }
}

// Categories using simple function
func mapCategories(arr) {
    return [categorize(age) for age in arr]
}

let categories = ages |> mapCategories
print("Categories:", categories)

print("\n✅ Todas las características P3 y P4 funcionando!")
print("   P3: Optional chaining, null coalescing, pattern matching")
print("   P4: Array/object comprehensions, pipeline operator")
print("   100% compatible con sintaxis existente")