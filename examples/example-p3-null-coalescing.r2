// Ejemplo de Null Coalescing Operator ??
// Característica P3 - Valores por Defecto Inteligentes

std.print("❓ R2Lang P3: Null Coalescing Operator ??")
std.print("==================================================")

// 1. Operador de coalescencia nula básico
std.print("\n🔄 1. Funcionamiento Básico:")
let nullValue = nil
let validValue = "Hello World"

let result1 = nullValue ?? "Valor por defecto"
let result2 = validValue ?? "No se usará"

std.print("nil ?? 'default':", result1)        // "Valor por defecto"
std.print("'Hello' ?? 'default':", result2)    // "Hello World"

// 2. Diferencia con operador OR ||
std.print("\n⚖️  2. Diferencia con Operador ||:")
let zero = 0
let emptyString = ""
let falseValue = false

std.print("Comparación con valores 'falsy':")
std.print("0 || 10:", zero || 10)           // 10 (|| considera 0 como falsy)
std.print("0 ?? 10:", zero ?? 10)           // 0 (solo nil activa ??)

std.print("'' || 'default':", emptyString || "default")     // "default"
std.print("'' ?? 'default':", emptyString ?? "default")     // "" (string vacío es válido)

std.print("false || true:", falseValue || true)       // true
std.print("false ?? true:", falseValue ?? true)       // false (false no es nil)

// 3. Encadenamiento de coalescencia
std.print("\n🔗 3. Encadenamiento Multiple:")
let config1 = nil
let config2 = nil  
let config3 = nil
let defaultConfig = {theme: "dark", timeout: 5000}

let activeConfig = config1 ?? config2 ?? config3 ?? defaultConfig
std.print("Configuración activa:", activeConfig.theme)

// 4. Casos de uso prácticos
std.print("\n🛠️  4. Casos de Uso Prácticos:")

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

std.print("Usuario con preferencias:")
std.print("  Tema:", user1Settings.theme)         // "dark"
std.print("  Tamaño:", user1Settings.fontSize)    // 16
std.print("  Idioma:", user1Settings.language)    // "es" (por defecto)

std.print("Usuario sin preferencias:")
std.print("  Tema:", user2Settings.theme)         // "light"
std.print("  Tamaño:", user2Settings.fontSize)    // 14
std.print("  Notificaciones:", user2Settings.notifications)  // true

// 5. Combinando con Optional Chaining
std.print("\n🤝 5. Combinando ?. y ??:")
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
    std.print(user.name + ": " + bio)
}

// 6. Patrón común: inicialización de variables
std.print("\n🎯 6. Inicialización Robusta:")
let environmentConfig = nil  // Podría venir de variables de entorno
let userConfig = nil         // Podría venir de base de datos
let defaultConfig = {
    port: 3000,
    host: "localhost",
    debug: false
}

let appConfig = environmentConfig ?? userConfig ?? defaultConfig
std.print("Puerto de aplicación:", appConfig.port)
std.print("Host:", appConfig.host)

std.print("\n✅ Null coalescing implementado exitosamente!")
std.print("   - Solo actúa con valores nil/null")
std.print("   - Preserva valores falsy válidos (0, '', false)")
std.print("   - Permite encadenamiento múltiple")
std.print("   - Se combina perfectamente con optional chaining")