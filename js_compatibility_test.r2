// Test de compatibilidad JavaScript -> R2Lang
print("=== Prueba de Compatibilidad JavaScript -> R2Lang ===")

// 1. Variables y tipos básicos (✅ Compatible)
print("\n1. Variables y tipos:")
let name = "Juan"
const PI = 3.14159
var counter = 0
let isActive = true
let data = null
print("✅ Variables:", name, PI, counter, isActive, data)

// 2. Operadores básicos (✅ Compatible) 
print("\n2. Operadores básicos:")
let sum = 5 + 3
let diff = 10 - 4
let prod = 6 * 2
let div = 15 / 3
let mod = 17 % 5
print("✅ Aritméticos:", sum, diff, prod, div, mod)

let eq = (5 == 5)
let neq = (3 != 4)  
let gt = (7 > 5)
let gte = (8 >= 8)
print("✅ Comparación:", eq, neq, gt, gte)

// 3. Operadores de asignación (✅ Compatible)
print("\n3. Asignación compuesta:")
counter += 5
counter *= 2
counter -= 3
print("✅ Asignación:", counter)

// 4. Negación lógica (✅ Compatible)
print("\n4. Negación lógica:")
let active = true
let inactive = !active
print("✅ Negación:", active, inactive)

// 5. Arrays (✅ Compatible)
print("\n5. Arrays:")
let numbers = [1, 2, 3, 4, 5]
let mixed = ["hello", 42, true, null]
print("✅ Arrays:", numbers, mixed)

// 6. Objetos/Maps (✅ Compatible)
print("\n6. Objetos:")
let user = {
    name: "Alice",
    age: 30,
    active: true
}
print("✅ Objetos:", user.name, user.age)

// 7. Funciones (✅ Compatible)
print("\n7. Funciones:")
func add(a, b) {
    return a + b
}
let result = add(10, 20)
print("✅ Funciones:", result)

// 8. Arrow functions (✅ Compatible)
print("\n8. Arrow functions:")
let multiply = (x, y) => x * y
let square = x => x * x
let getValue = () => 42
print("✅ Arrow functions:", multiply(3, 4), square(5), getValue())

// 9. Parámetros por defecto (✅ Compatible)
print("\n9. Parámetros por defecto:")
func greet(name = "World") {
    return "Hello " + name
}
print("✅ Default params:", greet(), greet("Alice"))

// 10. Destructuring (✅ Compatible)
print("\n10. Destructuring:")
let [first, second] = [10, 20]
let {name, age} = {name: "Bob", age: 25}
print("✅ Destructuring:", first, second, name, age)

// 11. Spread operator (✅ Compatible)
print("\n11. Spread operator:")
let arr1 = [1, 2, 3]
let arr2 = [4, 5, 6]
let combined = [...arr1, ...arr2]
print("✅ Spread:", combined)

// 12. Template strings (✅ Compatible)
print("\n12. Template strings:")
let person = "Carlos"
let age_person = 28
let message = `Hello ${person}, you are ${age_person} years old`
print("✅ Template strings:", message)

// 13. Optional chaining (✅ Compatible - P3)
print("\n13. Optional chaining:")
let profile = {user: {details: {city: "Madrid"}}}
let city = profile?.user?.details?.city
let missing = profile?.user?.phone?.number
print("✅ Optional chaining:", city, missing)

// 14. Null coalescing (✅ Compatible - P3)
print("\n14. Null coalescing:")
let config_value = null
let timeout = config_value ?? 5000
print("✅ Null coalescing:", timeout)

// 15. If/else (✅ Compatible)
print("\n15. Control de flujo:")
if (numbers.length > 0) {
    print("✅ If/else: Array no vacío")
} else {
    print("❌ Array vacío")
}

// 16. For loops (✅ Compatible)
print("\n16. Loops:")
for (let i = 0; i < 3; i++) {
    print("✅ For loop:", i)
}

// 17. Operadores bitwise (✅ Compatible - P2) 
print("\n17. Operadores bitwise:")
let bit_and = 5 & 3
let bit_or = 5 | 3  
let bit_xor = 5 ^ 3
print("✅ Bitwise:", bit_and, bit_or, bit_xor)

print("\n=== Resumen de Compatibilidad ===")
print("✅ COMPATIBLE: Variables, operadores, funciones, arrays, objetos")
print("✅ COMPATIBLE: Arrow functions, destructuring, spread, template strings")
print("✅ COMPATIBLE: Optional chaining, null coalescing")
print("✅ COMPATIBLE: Control de flujo, loops básicos")
print("✅ COMPATIBLE: Parámetros por defecto, asignación compuesta")
print("✅ COMPATIBLE: Operadores bitwise, negación lógica")