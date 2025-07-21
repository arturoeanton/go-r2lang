// Test básico de compatibilidad JavaScript
print("=== Compatibilidad JavaScript -> R2Lang ===")

// Variables y constantes
let name = "Carlos"
const VERSION = "2.0"
var counter = 5
print("Variables:", name, VERSION, counter)

// Operadores
counter += 10
let isValid = !false
let result = 10 > 5
print("Operadores:", counter, isValid, result)

// Arrays y objetos
let nums = [1, 2, 3]
let user = {name: "Ana", age: 25}
print("Estructuras:", nums, user.name)

// Funciones
let add = (a, b) => a + b
func multiply(x, y = 1) { return x * y }
print("Funciones:", add(5, 3), multiply(7))

// Destructuring y spread
let [first, second] = [10, 20]
let newArr = [...nums, 4, 5]
print("ES6:", first, second, newArr)

// Optional chaining y null coalescing (P3)
let obj = {data: {value: 42}}
let val = obj?.data?.value
let defaultVal = null ?? "default"
print("P3 features:", val, defaultVal)

// Pattern matching (P3)
func classify(x) {
    return match x {
        case 0 => "zero"
        case n if n > 0 => "positive"
        case _ => "negative"
    }
}
print("Pattern matching:", classify(5), classify(-2))

// Comprehensions (P4)
let squares = [x * x for x in nums]
let lookup = {x: x * 2 for x in nums}
print("Comprehensions:", squares, lookup)

// Pipeline (P4)
let processNumber = x => x |> (n => n * 2) |> (n => n + 1)
print("Pipeline:", processNumber(5))

print("\n✅ JavaScript moderno altamente compatible con R2Lang!")