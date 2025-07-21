// Ejemplo de Array y Object Comprehensions
// Característica P4 - Transformaciones Expresivas

print("📦 R2Lang P4: Array y Object Comprehensions")
print("=" * 55)

// 1. Array comprehensions básicas
print("\n🔄 1. Array Comprehensions Básicas:")
let numbers = [1, 2, 3, 4, 5]
let squares = [x * x for x in numbers]
print("Números:", numbers)
print("Cuadrados:", squares)

// 2. Comprehensions con filtros
print("\n🔍 2. Comprehensions con Filtros:")
let evens = [x for x in numbers if x % 2 == 0]
let evenSquares = [x * x for x in numbers if x % 2 == 0]
print("Pares:", evens)
print("Cuadrados pares:", evenSquares)

// 3. Transformaciones complejas
print("\n🎯 3. Transformaciones Complejas:")
let words = ["hello", "world", "r2lang", "comprehension"]
let lengths = [word.length for word in words if word.length > 5]
let upperLong = [word.toUpperCase() for word in words if word.length > 5]

print("Palabras largas (longitud):", lengths)
print("Palabras largas (mayúsculas):", upperLong)

// 4. Comprehensions anidadas (simuladas)
print("\n🎭 4. Múltiples Generadores:")
let matrix = []
// Simulación de comprensión anidada: [[i+j for j in range(3)] for i in range(3)]
for (let i = 1; i <= 3; i++) {
    let row = [i + j for j in [1, 2, 3]]
    matrix.push(row)
}
print("Matriz 3x3:", matrix)

// 5. Object comprehensions básicas
print("\n🏗️  5. Object Comprehensions:")
let fruits = ["apple", "banana", "cherry"]
let fruitLengths = {fruit: fruit.length for fruit in fruits}
print("Longitudes de frutas:", fruitLengths)

// 6. Object comprehensions con transformaciones
print("\n🔧 6. Object Comprehensions Avanzadas:")
let users = [
    {id: 1, name: "Ana", age: 25},
    {id: 2, name: "Luis", age: 30},
    {id: 3, name: "Carmen", age: 28}
]

// Crear lookup por ID
let userLookup = {user.id: user.name for user in users}
print("Lookup de usuarios:", userLookup)

// Filtrar usuarios activos (simulado)
let activeUsers = {user.name: user.age for user in users if user.age >= 28}
print("Usuarios activos (≥28):", activeUsers)

// 7. Casos de uso prácticos
print("\n💼 7. Casos de Uso Prácticos:")

// Procesamiento de datos de ventas
let sales = [
    {product: "Laptop", price: 1200, quantity: 2},
    {product: "Mouse", price: 25, quantity: 10},
    {product: "Keyboard", price: 75, quantity: 5}
]

// Calcular totales por producto
let totals = {sale.product: sale.price * sale.quantity for sale in sales}
print("Totales por producto:", totals)

// Productos caros (>$100 total)
let expensiveProducts = {sale.product: sale.price * sale.quantity for sale in sales if sale.price * sale.quantity > 100}
print("Productos caros:", expensiveProducts)

// 8. Transformación de datos API
print("\n🌐 8. Transformación de Datos de API:")
let apiResponse = [
    {userId: 1, userName: "alice", email: "alice@example.com", active: true},
    {userId: 2, userName: "bob", email: "bob@example.com", active: false},
    {userId: 3, userName: "charlie", email: "charlie@example.com", active: true}
]

// Extraer solo usuarios activos con email
let activeEmails = [user.email for user in apiResponse if user.active]
print("Emails de usuarios activos:", activeEmails)

// Crear mapa de username -> email para usuarios activos
let activeUserEmails = {user.userName: user.email for user in apiResponse if user.active}
print("Mapa username->email:", activeUserEmails)

// 9. Procesamiento de strings
print("\n📝 9. Procesamiento de Strings:")
let sentence = "The quick brown fox jumps"
let wordsList = sentence.split(" ")
let wordCounts = {word: word.length for word in wordsList}
let longWords = [word for word in wordsList if word.length > 4]

print("Frase original:", sentence)
print("Conteo de caracteres:", wordCounts)
print("Palabras largas:", longWords)

// 10. Comparación: antes vs después
print("\n⚖️  10. Antes vs Después de Comprehensions:")

// ANTES - Código imperativo verboso
print("❌ ANTES (imperativo):")
let evenSquaresBefore = []
for (let i = 0; i < numbers.length; i++) {
    if (numbers[i] % 2 == 0) {
        evenSquaresBefore.push(numbers[i] * numbers[i])
    }
}
print("  Cuadrados pares:", evenSquaresBefore)

// DESPUÉS - Código declarativo conciso  
print("✅ DESPUÉS (declarativo):")
let evenSquaresAfter = [x * x for x in numbers if x % 2 == 0]
print("  Cuadrados pares:", evenSquaresAfter)

print("\n✅ Array/Object Comprehensions implementadas exitosamente!")
print("   - Sintaxis expresiva y concisa")
print("   - Soporte para filtros con 'if'")
print("   - Transformaciones complejas en una línea")
print("   - Código más legible y mantenible")
print("   - Compatible con arrays y objetos")