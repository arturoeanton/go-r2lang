// Ejemplo de Pattern Matching con match
// Característica P3 - Lógica Condicional Expresiva

print("🎯 R2Lang P3: Pattern Matching con match")
print("=" * 50)

// 1. Pattern matching básico con literales
print("\n🔢 1. Matching con Literales:")
func describeNumber(n) {
    return match n {
        case 0 => "Zero"
        case 1 => "Uno"
        case 2 => "Dos"
        case _ => "Otro número"
    }
}

print("Número 0:", describeNumber(0))    // "Zero"
print("Número 1:", describeNumber(1))    // "Uno"  
print("Número 5:", describeNumber(5))    // "Otro número"

// 2. Pattern matching con variables (binding)
print("\n🏷️  2. Binding de Variables:")
func processValue(x) {
    return match x {
        case n => "El valor es: " + n
    }
}

print("Binding:", processValue("Hello"))   // "El valor es: Hello"
print("Binding:", processValue(42))        // "El valor es: 42"

// 3. Pattern matching con arrays
print("\n📋 3. Destructuring de Arrays:")
func analyzeArray(arr) {
    return match arr {
        case [a, b, c] => "Tres elementos: " + a + ", " + b + ", " + c
        case [x, y] => "Dos elementos: " + x + ", " + y
        case [single] => "Un elemento: " + single
        case [] => "Array vacío"
        case _ => "Array con más elementos"
    }
}

print("Array [1,2,3]:", analyzeArray([1, 2, 3]))
print("Array [10,20]:", analyzeArray([10, 20]))
print("Array [42]:", analyzeArray([42]))
print("Array []:", analyzeArray([]))

// 4. Pattern matching con objetos
print("\n🏢 4. Destructuring de Objetos:")
func processUser(user) {
    return match user {
        case {name, age} => name + " tiene " + age + " años"
        case {name} => name + " (edad desconocida)"
        case _ => "Usuario inválido"
    }
}

let user1 = {name: "Ana", age: 28}
let user2 = {name: "Luis"}
let user3 = {id: 123}

print("Usuario completo:", processUser(user1))
print("Usuario sin edad:", processUser(user2))
print("Usuario inválido:", processUser(user3))

// 5. Pattern matching con guards (condiciones)
print("\n🛡️  5. Guards (Condiciones):")
func categorizeAge(person) {
    return match person {
        case {age: a} if a < 13 => "Niño"
        case {age: a} if a < 20 => "Adolescente"
        case {age: a} if a < 60 => "Adulto"
        case {age: a} if a >= 60 => "Adulto mayor"
        case _ => "Edad no especificada"
    }
}

let person1 = {name: "María", age: 10}
let person2 = {name: "Carlos", age: 17}
let person3 = {name: "Elena", age: 35}
let person4 = {name: "Roberto", age: 65}

print("María (10):", categorizeAge(person1))
print("Carlos (17):", categorizeAge(person2))
print("Elena (35):", categorizeAge(person3))
print("Roberto (65):", categorizeAge(person4))

// 6. Casos de uso avanzados - HTTP Status
print("\n🌐 6. Procesamiento de Respuestas HTTP:")
func handleHttpResponse(response) {
    return match response.status {
        case 200 => "✅ Éxito: " + response.data
        case 201 => "✅ Creado: " + response.data
        case 400 => "❌ Solicitud incorrecta"
        case 401 => "🔐 No autorizado"
        case 404 => "🔍 No encontrado"
        case s if s >= 500 => "💥 Error del servidor: " + s
        case s => "⚠️  Estado desconocido: " + s
    }
}

let responses = [
    {status: 200, data: "Usuario creado"},
    {status: 404, data: nil},
    {status: 500, data: nil},
    {status: 418, data: nil}  // I'm a teapot!
]

for (let i = 0; i < responses.length; i++) {
    let response = responses[i]
    print("Status " + response.status + ":", handleHttpResponse(response))
}

// 7. Pattern matching anidado
print("\n🎭 7. Matching Anidado Complejo:")
func analyzeRequest(req) {
    return match req {
        case {method: "GET", url: path} => "Obteniendo: " + path
        case {method: "POST", body: {name, email}} => "Creando usuario: " + name + " (" + email + ")"
        case {method: "PUT", url: path, body: data} if path.startsWith("/users/") => "Actualizando usuario en " + path
        case {method: m} => "Método no soportado: " + m
        case _ => "Solicitud inválida"
    }
}

let requests = [
    {method: "GET", url: "/api/users"},
    {method: "POST", body: {name: "Juan", email: "juan@example.com"}},
    {method: "PUT", url: "/users/123", body: {name: "Juan Carlos"}},
    {method: "DELETE", url: "/users/123"},
    {invalid: "request"}
]

for (let i = 0; i < requests.length; i++) {
    let req = requests[i]
    print("Request " + (i + 1) + ":", analyzeRequest(req))
}

print("\n✅ Pattern matching implementado exitosamente!")
print("   - Matching con literales, variables y wildcards")
print("   - Destructuring de arrays y objetos")
print("   - Guards para condiciones complejas")
print("   - Código más expresivo y mantenible")