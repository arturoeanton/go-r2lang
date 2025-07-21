// Demostración simple de P3 y P4 features

print("R2Lang P3 y P4 - Demostración")

// P3.1: Optional Chaining
print("P3.1: Optional Chaining")
let user = {name: "Alice", profile: {city: "Madrid"}}
let city = user?.profile?.city
let missing = user?.profile?.phone
print("City:", city)
print("Phone:", missing)

// P3.2: Null Coalescing
print("P3.2: Null Coalescing") 
let config = nil
let port = config ?? 8080
print("Port:", port)

// P3.3: Pattern Matching
print("P3.3: Pattern Matching")
func getStatus(code) {
    return match code {
        case 200 => "OK"
        case 404 => "Not Found"
        case _ => "Other"
    }
}
print("200:", getStatus(200))
print("404:", getStatus(404))

// P4.1: Array Comprehensions
print("P4.1: Array Comprehensions")
let nums = [1, 2, 3, 4, 5]
let doubled = [x * 2 for x in nums]
print("Doubled:", doubled)

// P4.2: Object Comprehensions  
print("P4.2: Object Comprehensions")
let words = ["cat", "dog", "bird"]
let lengths = {word: word.length for word in words}
print("Lengths:", lengths)

// P4.3: Pipeline Operator
print("P4.3: Pipeline Operator")
func triple(x) { return x * 3 }
func addFive(x) { return x + 5 }
let result = 10 |> triple |> addFive
print("10 |> triple |> addFive =", result)

print("Todas las características implementadas!")