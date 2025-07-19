// Ejemplo 28: Detección de Loops Infinitos en R2Lang
// Este ejemplo demuestra cómo R2Lang detecta y previene loops infinitos

std.print("=== EJEMPLO 28: DETECCIÓN DE LOOPS INFINITOS ===\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 1. EJEMPLO DE WHILE LOOP INFINITO
// ═══════════════════════════════════════════════════════════════════════════════

std.print("1. Detectando while loop infinito:")
std.print("   Código: while (true) { print('Iterando...'); }")

// Nota: Este código sería detectado como loop infinito
// while (true) {
//     std.print("Iterando...")
// }

std.print("   ✓ R2Lang detectaría este loop infinito después de cierto número de iteraciones\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 2. EJEMPLO DE FOR LOOP INFINITO
// ═══════════════════════════════════════════════════════════════════════════════

std.print("2. Detectando for loop infinito:")
std.print("   Código: for (let i = 0; i >= 0; i++) { print('Contando:', i); }")

// Nota: Este código sería detectado como loop infinito
// for (let i = 0; i >= 0; i++) {
//     std.print("Contando:", i)
// }

std.print("   ✓ R2Lang detectaría este loop infinito por exceder el límite de iteraciones\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 3. EJEMPLO DE RECURSIÓN INFINITA
// ═══════════════════════════════════════════════════════════════════════════════

std.print("3. Detectando recursión infinita:")
std.print("   Código: func infinita() { infinita(); }")

// Función que causaría recursión infinita
func recursionInfinita() {
    std.print("Llamada recursiva...")
    // recursionInfinita() // Esta línea causaría recursión infinita
}

std.print("   ✓ R2Lang detectaría esta recursión infinita por exceder la profundidad máxima\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 4. EJEMPLOS DE LOOPS FINITOS (CORRECTOS)
// ═══════════════════════════════════════════════════════════════════════════════

std.print("4. Ejemplos de loops finitos correctos:")

// While loop finito
std.print("   - While loop finito:")
let contador = 0
while (contador < 3) {
    std.print("     Iteración:", contador)
    contador++
}
std.print("   ✓ Completado correctamente\n")

// For loop finito
std.print("   - For loop finito:")
for (let i = 0; i < 3; i++) {
    std.print("     For iteración:", i)
}
std.print("   ✓ Completado correctamente\n")

// For-in loop finito
std.print("   - For-in loop finito:")
let frutas = ["manzana", "banana", "naranja"]
for (fruta in frutas) {
    std.print("     Fruta:", frutas[fruta])
}
std.print("   ✓ Completado correctamente\n")

// Recursión finita
std.print("   - Recursión finita (factorial):")
func factorial(n) {
    if (n <= 1) {
        return 1
    }
    return n * factorial(n - 1)
}

let resultado = factorial(5)
std.print("     factorial(5) =", resultado)
std.print("   ✓ Completado correctamente\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 5. CONFIGURACIÓN DE LÍMITES
// ═══════════════════════════════════════════════════════════════════════════════

std.print("5. Configuración de límites:")
std.print("   R2Lang utiliza límites configurables para detectar loops infinitos:")
std.print("   - Límite de iteraciones: Por defecto 1,000,000 iteraciones")
std.print("   - Límite de recursión: Por defecto 1,000 niveles de profundidad")
std.print("   - Timeout global: Tiempo máximo de ejecución")
std.print("   - Cancelación por contexto: Permite interrumpir ejecución\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 6. TIPOS DE ERRORES
// ═══════════════════════════════════════════════════════════════════════════════

std.print("6. Tipos de errores de loop infinito:")
std.print("   - InfiniteLoopError: Para loops while y for infinitos")
std.print("   - RecursionError: Para recursión infinita")
std.print("   - TimeoutError: Para ejecución que excede el tiempo límite")
std.print("   - ContextCancelError: Para ejecución cancelada por contexto\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 7. BENEFICIOS DE LA DETECCIÓN
// ═══════════════════════════════════════════════════════════════════════════════

std.print("7. Beneficios de la detección de loops infinitos:")
std.print("   ✓ Previene que el programa se cuelgue indefinidamente")
std.print("   ✓ Proporciona mensajes de error informativos")
std.print("   ✓ Permite depuración más fácil")
std.print("   ✓ Mejora la estabilidad del intérprete")
std.print("   ✓ Protege recursos del sistema\n")

std.print("=== FIN DEL EJEMPLO ===")
std.print("✅ R2Lang ahora incluye detección automática de loops infinitos")
std.print("🛡️  Tu código está protegido contra loops infinitos accidentales")