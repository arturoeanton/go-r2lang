// Ejemplo 28: Detección de Loops Infinitos en R2Lang
// Este ejemplo demuestra cómo R2Lang detecta y previene loops infinitos

print("=== EJEMPLO 28: DETECCIÓN DE LOOPS INFINITOS ===\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 1. EJEMPLO DE WHILE LOOP INFINITO
// ═══════════════════════════════════════════════════════════════════════════════

print("1. Detectando while loop infinito:")
print("   Código: while (true) { print('Iterando...'); }")

// Nota: Este código sería detectado como loop infinito
// while (true) {
//     print("Iterando...")
// }

print("   ✓ R2Lang detectaría este loop infinito después de cierto número de iteraciones\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 2. EJEMPLO DE FOR LOOP INFINITO
// ═══════════════════════════════════════════════════════════════════════════════

print("2. Detectando for loop infinito:")
print("   Código: for (let i = 0; i >= 0; i++) { print('Contando:', i); }")

// Nota: Este código sería detectado como loop infinito
// for (let i = 0; i >= 0; i++) {
//     print("Contando:", i)
// }

print("   ✓ R2Lang detectaría este loop infinito por exceder el límite de iteraciones\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 3. EJEMPLO DE RECURSIÓN INFINITA
// ═══════════════════════════════════════════════════════════════════════════════

print("3. Detectando recursión infinita:")
print("   Código: func infinita() { infinita(); }")

// Función que causaría recursión infinita
func recursionInfinita() {
    print("Llamada recursiva...")
    // recursionInfinita() // Esta línea causaría recursión infinita
}

print("   ✓ R2Lang detectaría esta recursión infinita por exceder la profundidad máxima\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 4. EJEMPLOS DE LOOPS FINITOS (CORRECTOS)
// ═══════════════════════════════════════════════════════════════════════════════

print("4. Ejemplos de loops finitos correctos:")

// While loop finito
print("   - While loop finito:")
let contador = 0
while (contador < 3) {
    print("     Iteración:", contador)
    contador++
}
print("   ✓ Completado correctamente\n")

// For loop finito
print("   - For loop finito:")
for (let i = 0; i < 3; i++) {
    print("     For iteración:", i)
}
print("   ✓ Completado correctamente\n")

// For-in loop finito
print("   - For-in loop finito:")
let frutas = ["manzana", "banana", "naranja"]
for (fruta in frutas) {
    print("     Fruta:", frutas[fruta])
}
print("   ✓ Completado correctamente\n")

// Recursión finita
print("   - Recursión finita (factorial):")
func factorial(n) {
    if (n <= 1) {
        return 1
    }
    return n * factorial(n - 1)
}

let resultado = factorial(5)
print("     factorial(5) =", resultado)
print("   ✓ Completado correctamente\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 5. CONFIGURACIÓN DE LÍMITES
// ═══════════════════════════════════════════════════════════════════════════════

print("5. Configuración de límites:")
print("   R2Lang utiliza límites configurables para detectar loops infinitos:")
print("   - Límite de iteraciones: Por defecto 1,000,000 iteraciones")
print("   - Límite de recursión: Por defecto 1,000 niveles de profundidad")
print("   - Timeout global: Tiempo máximo de ejecución")
print("   - Cancelación por contexto: Permite interrumpir ejecución\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 6. TIPOS DE ERRORES
// ═══════════════════════════════════════════════════════════════════════════════

print("6. Tipos de errores de loop infinito:")
print("   - InfiniteLoopError: Para loops while y for infinitos")
print("   - RecursionError: Para recursión infinita")
print("   - TimeoutError: Para ejecución que excede el tiempo límite")
print("   - ContextCancelError: Para ejecución cancelada por contexto\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 7. BENEFICIOS DE LA DETECCIÓN
// ═══════════════════════════════════════════════════════════════════════════════

print("7. Beneficios de la detección de loops infinitos:")
print("   ✓ Previene que el programa se cuelgue indefinidamente")
print("   ✓ Proporciona mensajes de error informativos")
print("   ✓ Permite depuración más fácil")
print("   ✓ Mejora la estabilidad del intérprete")
print("   ✓ Protege recursos del sistema\n")

print("=== FIN DEL EJEMPLO ===")
print("✅ R2Lang ahora incluye detección automática de loops infinitos")
print("🛡️  Tu código está protegido contra loops infinitos accidentales")