// Ejemplo de Optional Chaining (Encadenamiento Opcional) ?.
// Característica P3 - Navegación Segura

print("🔗 R2Lang P3: Optional Chaining Operator ?.")
print("=" * 50)

// 1. Acceso seguro a propiedades anidadas
print("\n📋 1. Acceso Seguro a Propiedades:")
let user = {
    name: "Alice",
    profile: {
        address: {
            street: "123 Main St",
            city: "Madrid"
        }
    }
}

// Acceso tradicional (puede fallar)
let street1 = user.profile.address.street
print("Calle tradicional:", street1)

// Acceso con optional chaining (seguro)
let street2 = user?.profile?.address?.street
print("Calle opcional:", street2)

// 2. Navegación segura con objetos nil
print("\n🚫 2. Navegación con Objetos Nulos:")
let emptyUser = nil
let safeName = emptyUser?.name
print("Nombre seguro:", safeName)  // nil en lugar de error

// 3. Propiedades inexistentes
print("\n❓ 3. Propiedades Inexistentes:")
let userWithoutPhone = {name: "Bob", email: "bob@example.com"}
let phone = userWithoutPhone?.phone?.number
print("Teléfono:", phone)  // nil en lugar de error

// 4. Combinando con objetos complejos
print("\n🏢 4. Casos de Uso Reales:")
let apiResponse = {
    status: "success",
    data: {
        users: [
            {id: 1, name: "Juan", profile: {avatar: "juan.jpg"}},
            {id: 2, name: "María", profile: nil},
            {id: 3, name: "Carlos"}
        ]
    }
}

// Acceso seguro a avatares de usuarios
let avatar1 = apiResponse?.data?.users?.[0]?.profile?.avatar
let avatar2 = apiResponse?.data?.users?.[1]?.profile?.avatar
let avatar3 = apiResponse?.data?.users?.[2]?.profile?.avatar

print("Avatar de Juan:", avatar1)      // "juan.jpg"
print("Avatar de María:", avatar2)     // nil
print("Avatar de Carlos:", avatar3)    // nil

// 5. Evitando errores comunes
print("\n⚠️  5. Comparación Sin/Con Optional Chaining:")

// SIN optional chaining - propenso a errores
let config = {server: {port: 8080}}
// let unsafeTimeout = config.server.timeout.value  // ❌ Esto fallaría

// CON optional chaining - seguro
let safeTimeout = config?.server?.timeout?.value
print("Timeout seguro:", safeTimeout)  // nil en lugar de error

print("\n✅ Optional chaining implementado exitosamente!")
print("   - Navegación segura sin errores")
print("   - Retorna nil en lugar de panic")
print("   - Sintaxis familiar para desarrolladores JS/TS")