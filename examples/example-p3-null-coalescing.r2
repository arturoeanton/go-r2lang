// Ejemplo de Null Coalescing Operator ??
// Característica P3 - Valores por Defecto Inteligentes

print("❓ R2Lang P3: Null Coalescing Operator ??")
print("=" * 50)

// 1. Operador de coalescencia nula básico
print("\n🔄 1. Funcionamiento Básico:")
let nullValue = nil
let validValue = "Hello World"

let result1 = nullValue ?? "Valor por defecto"
let result2 = validValue ?? "No se usará"

print("nil ?? 'default':", result1)        // "Valor por defecto"
print("'Hello' ?? 'default':", result2)    // "Hello World"

// 2. Diferencia con operador OR ||
print("\n⚖️  2. Diferencia con Operador ||:")
let zero = 0
let emptyString = ""
let falseValue = false

print("Comparación con valores 'falsy':")
print("0 || 10:", zero || 10)           // 10 (|| considera 0 como falsy)
print("0 ?? 10:", zero ?? 10)           // 0 (solo nil activa ??)

print("'' || 'default':", emptyString || "default")     // "default"
print("'' ?? 'default':", emptyString ?? "default")     // "" (string vacío es válido)

print("false || true:", falseValue || true)       // true
print("false ?? true:", falseValue ?? true)       // false (false no es nil)

// 3. Encadenamiento de coalescencia
print("\n🔗 3. Encadenamiento Multiple:")
let config1 = nil
let config2 = nil  
let config3 = nil
let defaultConfig = {theme: "dark", timeout: 5000}

let activeConfig = config1 ?? config2 ?? config3 ?? defaultConfig
print("Configuración activa:", activeConfig.theme)

// 4. Casos de uso prácticos
print("\n🛠️  4. Casos de Uso Prácticos:")

// Configuración de usuario con valores por defecto
func getUserSettings(userPrefs) {
    return {
        theme: userPrefs?.theme ?? "light",
        fontSize: userPrefs?.fontSize ?? 14,
        language: userPrefs?.language ?? "es",
        notifications: userPrefs?.notifications ?? true
    }
}

let user1Settings = getUserSettings({theme: "dark", fontSize: 16})
let user2Settings = getUserSettings(nil)

print("Usuario con preferencias:")
print("  Tema:", user1Settings.theme)         // "dark"
print("  Tamaño:", user1Settings.fontSize)    // 16
print("  Idioma:", user1Settings.language)    // "es" (por defecto)

print("Usuario sin preferencias:")
print("  Tema:", user2Settings.theme)         // "light"
print("  Tamaño:", user2Settings.fontSize)    // 14
print("  Notificaciones:", user2Settings.notifications)  // true

// 5. Combinando con Optional Chaining
print("\n🤝 5. Combinando ?. y ??:")
let apiData = {
    users: [
        {name: "Ana", profile: {bio: "Desarrolladora"}},
        {name: "Luis", profile: nil},
        {name: "Carmen"}
    ]
}

// Obtener biografías con valores por defecto seguros
for (let i = 0; i < 3; i++) {
    let user = apiData.users[i]
    let bio = user?.profile?.bio ?? "Sin biografía disponible"
    print(user.name + ": " + bio)
}

// 6. Patrón común: inicialización de variables
print("\n🎯 6. Inicialización Robusta:")
let environmentConfig = nil  // Podría venir de variables de entorno
let userConfig = nil         // Podría venir de base de datos
let defaultConfig = {
    port: 3000,
    host: "localhost",
    debug: false
}

let appConfig = environmentConfig ?? userConfig ?? defaultConfig
print("Puerto de aplicación:", appConfig.port)
print("Host:", appConfig.host)

print("\n✅ Null coalescing implementado exitosamente!")
print("   - Solo actúa con valores nil/null")
print("   - Preserva valores falsy válidos (0, '', false)")
print("   - Permite encadenamiento múltiple")
print("   - Se combina perfectamente con optional chaining")