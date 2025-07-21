// Basic P3/P4 features test
print("Testing P3/P4 features")

// P3.1: Optional chaining
let obj = {a: {b: "hello"}}
let result1 = obj?.a?.b
print("Optional chaining:", result1)

// P3.2: Null coalescing
let config = nil
let port = config ?? 8080
print("Null coalescing:", port)

// P3.3: Pattern matching
func testMatch(x) {
    return match x {
        case 1 => "one"
        case 2 => "two"
        case _ => "other"
    }
}
print("Pattern matching:", testMatch(1))
print("Pattern matching:", testMatch(3))

// P4.1: Array comprehensions
let nums = [1, 2, 3, 4, 5]
let doubled = [x * 2 for x in nums]
print("Array comprehension:", doubled)

// P4.2: Object comprehensions
let words = ["cat", "dog"]
let lengths = {word: word.length for word in words}
print("Object comprehension:", lengths)

// P4.3: Pipeline operator
func double(x) { return x * 2 }
func addTen(x) { return x + 10 }
let result2 = 5 |> double |> addTen
print("Pipeline:", result2)

print("All P3/P4 features working!")