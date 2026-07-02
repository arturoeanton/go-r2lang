// Ejemplo de Array y Object Comprehensions
// Característica P4 - Transformaciones Expresivas

std.print("📦 R2Lang P4: Array y Object Comprehensions")
std.print("=======================================================")

// 1. Array comprehensions básicas
std.print("\n🔄 1. Array Comprehensions Básicas:")
let numbers = [1, 2, 3, 4, 5]
let squares = [x * x for x in numbers]
std.print("Números:", numbers)
std.print("Cuadrados:", squares)

// 2. Comprehensions con filtros
std.print("\n🔍 2. Comprehensions con Filtros:")
let evens = [x for x in numbers if x % 2 == 0]
let evenSquares = [x * x for x in numbers if x % 2 == 0]
std.print("Pares:", evens)
std.print("Cuadrados pares:", evenSquares)

// 3. Transformaciones complejas
std.print("\n🎯 3. Transformaciones Complejas:")
let words = ["hello", "world", "r2lang", "comprehension"]
let lengths = [word.length for word in words if word.length > 5]
let upperLong = [string.toUpper(word) for word in words if word.length > 5]

std.print("Palabras largas (longitud):", lengths)
std.print("Palabras largas (mayúsculas):", upperLong)

// 4. Comprehensions anidadas (simuladas)
std.print("\n🎭 4. Múltiples Generadores:")
let matrix = []
// Simulación de comprensión anidada: [[i+j for j in range(3)] for i in range(3)]
for (let i = 1; i <= 3; i++) {
    let row = [i + j for j in [1, 2, 3]]
    matrix.push(row)
}
std.print("Matriz 3x3:", matrix)

// 5. Object comprehensions básicas
std.print("\n🏗️  5. Object Comprehensions:")
let fruits = ["apple", "banana", "cherry"]
let fruitLengths = {fruit: fruit.length for fruit in fruits}
std.print("Longitudes de frutas:", fruitLengths)

// 6. Object comprehensions con transformaciones
std.print("\n🔧 6. Object Comprehensions Avanzadas:")
let users = [
    {id: 1, name: "Ana", age: 25},
    {id: 2, name: "Luis", age: 30},
    {id: 3, name: "Carmen", age: 28}
]

// Crear lookup por ID
let userLookup = {user.id: user.name for user in users}
std.print("Lookup de usuarios:", userLookup)

// Filtrar usuarios activos (simulado)
let activeUsers = {user.name: user.age for user in users if user.age >= 28}
std.print("Usuarios activos (≥28):", activeUsers)

// 7. Casos de uso prácticos
std.print("\n💼 7. Casos de Uso Prácticos:")

// Procesamiento de datos de ventas
let sales = [
    {product: "Laptop", price: 1200, quantity: 2},
    {product: "Mouse", price: 25, quantity: 10},
    {product: "Keyboard", price: 75, quantity: 5}
]

// Calcular totales por producto
let totals = {sale.product: sale.price * sale.quantity for sale in sales}
std.print("Totales por producto:", totals)

// Productos caros (>$100 total)
let expensiveProducts = {sale.product: sale.price * sale.quantity for sale in sales if sale.price * sale.quantity > 100}
std.print("Productos caros:", expensiveProducts)

// 8. Transformación de datos API
std.print("\n🌐 8. Transformación de Datos de API:")
let apiResponse = [
    {userId: 1, userName: "alice", email: "alice@example.com", active: true},
    {userId: 2, userName: "bob", email: "bob@example.com", active: false},
    {userId: 3, userName: "charlie", email: "charlie@example.com", active: true}
]

// Extraer solo usuarios activos con email
let activeEmails = [user.email for user in apiResponse if user.active]
std.print("Emails de usuarios activos:", activeEmails)

// Crear mapa de username -> email para usuarios activos
let activeUserEmails = {user.userName: user.email for user in apiResponse if user.active}
std.print("Mapa username->email:", activeUserEmails)

// 9. Procesamiento de strings
std.print("\n📝 9. Procesamiento de Strings:")
let sentence = "The quick brown fox jumps"
let wordsList = string.split(sentence, " ")
let wordCounts = {word: word.length for word in wordsList}
let longWords = [word for word in wordsList if word.length > 4]

std.print("Frase original:", sentence)
std.print("Conteo de caracteres:", wordCounts)
std.print("Palabras largas:", longWords)

// 10. Comparación: antes vs después
std.print("\n⚖️  10. Antes vs Después de Comprehensions:")

// ANTES - Código imperativo verboso
std.print("❌ ANTES (imperativo):")
let evenSquaresBefore = []
for (let i = 0; i < numbers.length; i++) {
    if (numbers[i] % 2 == 0) {
        evenSquaresBefore.push(numbers[i] * numbers[i])
    }
}
std.print("  Cuadrados pares:", evenSquaresBefore)

// DESPUÉS - Código declarativo conciso  
std.print("✅ DESPUÉS (declarativo):")
let evenSquaresAfter = [x * x for x in numbers if x % 2 == 0]
std.print("  Cuadrados pares:", evenSquaresAfter)

std.print("\n✅ Array/Object Comprehensions implementadas exitosamente!")
std.print("   - Sintaxis expresiva y concisa")
std.print("   - Soporte para filtros con 'if'")
std.print("   - Transformaciones complejas en una línea")
std.print("   - Código más legible y mantenible")
std.print("   - Compatible con arrays y objetos")