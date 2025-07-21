// Demostración simple de P3 y P4 features

std.print("R2Lang P3 y P4 - Demostración")

// P3.1: Optional Chaining
std.print("P3.1: Optional Chaining")
let user = {name: "Alice", profile: {city: "Madrid"}}
let city = user?.profile?.city
let missing = user?.profile?.phone
std.print("City:", city)
std.print("Phone:", missing)

// P3.2: Null Coalescing
std.print("P3.2: Null Coalescing") 
let config = nil
let port = config ?? 8080
std.print("Port:", port)

// P3.3: Pattern Matching
std.print("P3.3: Pattern Matching")
func getStatus(code) {
    return match code {
        case x if x >= 100 && x < 200 => "Informational"
        case x if x >= 200 && x < 300 => "OK"
        case x if x >= 300 && x < 400 => "Redirect"
        case x if x >= 400 && x < 500 => "Client Error"
        case x if x >= 500 => "Server Error"
        case _ => "Other"
    }
}
std.print("200:", getStatus(200))
std.print("300:", getStatus(300))
std.print("500:", getStatus(500))
std.print("Client Error:", getStatus(404))
std.print("Server Error:", getStatus(500))
std.print("");
std.print("Other:", getStatus(100)) 
std.print("404:", getStatus(404))

// P4.1: Array Comprehensions
std.print("P4.1: Array Comprehensions")
let nums = [1, 2, 3, 4, 5]
let doubled = [x * 2 for x in nums]
std.print("Doubled:", doubled)

// P4.2: Object Comprehensions  
std.print("P4.2: Object Comprehensions")
let words = ["cat", "dog", "bird"]
let lengths = {word: word.length for word in words}
std.print("Lengths:", lengths)

// P4.3: Pipeline Operator
std.print("P4.3: Pipeline Operator")
func triple(x) { return x * 3 }
func addFive(x) { return x + 5 }
let result = 10 |> triple |> addFive
std.print("10 |> triple |> addFive =", result)

std.print("Todas las características implementadas!")