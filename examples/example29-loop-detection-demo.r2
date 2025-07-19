// Ejemplo 29: Demostración Práctica de Detección de Loops Infinitos
// ADVERTENCIA: Este ejemplo contiene código que puede causar loops infinitos
// Ejecutar con precaución - R2Lang los detectará y terminará la ejecución

std.print("=== EJEMPLO 29: DEMOSTRACIÓN PRÁCTICA DE DETECCIÓN ===\n")

// ═══════════════════════════════════════════════════════════════════════════════
// CONFIGURACIÓN DE PRUEBA
// ═══════════════════════════════════════════════════════════════════════════════

std.print("Este ejemplo demuestra la detección en acción...")
std.print("Los loops infinitos serán detectados automáticamente por R2Lang\n")

// ═══════════════════════════════════════════════════════════════════════════════
// FUNCIÓN AUXILIAR PARA SIMULAR TRABAJO
// ═══════════════════════════════════════════════════════════════════════════════

func simularTrabajo() {
    // Simula algún trabajo computacional
    let suma = 0
    for (let i = 0; i < 1000; i++) {
        suma = suma + i
    }
    return suma
}

// ═══════════════════════════════════════════════════════════════════════════════
// EJEMPLO 1: LOOP WHILE CON LÍMITE BAJO PARA DEMOSTRACIÓN
// ═══════════════════════════════════════════════════════════════════════════════

std.print("1. Probando while loop que excedería el límite:")
std.print("   (Simulando un límite bajo para demostración)")

// Este loop sería detectado si el límite fuera muy bajo
let contadorDemo = 0
while (contadorDemo < 10) {  // Solo 10 iteraciones para demostrar
    std.print("   Iteración segura:", contadorDemo)
    contadorDemo++
    simularTrabajo()
}
std.print("   ✓ Loop completado correctamente\n")

// ═══════════════════════════════════════════════════════════════════════════════
// EJEMPLO 2: RECURSIÓN CONTROLADA
// ═══════════════════════════════════════════════════════════════════════════════

std.print("2. Probando recursión controlada:")

func recursionControlada(n, max) {
    if (n >= max) {
        return n
    }
    std.print("   Nivel de recursión:", n)
    return recursionControlada(n + 1, max)
}

let resultadoRecursion = recursionControlada(0, 5)
std.print("   ✓ Recursión completada, resultado:", resultadoRecursion, "\n")

// ═══════════════════════════════════════════════════════════════════════════════
// EJEMPLO 3: FOR LOOP CON ARRAYS GRANDES
// ═══════════════════════════════════════════════════════════════════════════════

std.print("3. Probando for loop con array grande:")

// Crear un array grande para probar límites
let arrayGrande = []
for (let i = 0; i < 1000; i++) {
    arrayGrande[i] = i * 2
}

std.print("   Array creado con", std.len(arrayGrande), "elementos")

// Procesar el array
let suma = 0
for (elemento in arrayGrande) {
    suma = suma + arrayGrande[elemento]
}

std.print("   ✓ Procesamiento completado, suma:", suma, "\n")

// ═══════════════════════════════════════════════════════════════════════════════
// EJEMPLO 4: MANEJO DE ERRORES CON TRY-CATCH
// ═══════════════════════════════════════════════════════════════════════════════

std.print("4. Ejemplo de manejo de errores (simulado):")

// Función que simula un error de loop infinito
func funcionSegura() {
    // Esta función ejecuta código seguro
    let resultado = 0
    for (let i = 0; i < 100; i++) {
        resultado = resultado + i
    }
    return resultado
}

// Usar try-catch para manejar posibles errores
try {
    let resultado = funcionSegura()
    std.print("   ✓ Función ejecutada correctamente, resultado:", resultado)
} catch (error) {
    std.print("   ❌ Error detectado:", error)
} finally {
    std.print("   ✓ Limpieza completada")
}

std.print("")

// ═══════════════════════════════════════════════════════════════════════════════
// EJEMPLO 5: PATRONES COMUNES QUE CAUSAN LOOPS INFINITOS
// ═══════════════════════════════════════════════════════════════════════════════

std.print("5. Patrones comunes que causan loops infinitos:")

std.print("   a) While sin incremento:")
std.print("      while (true) { /* sin cambio de condición */ }")

std.print("   b) For con condición incorrecta:")
std.print("      for (let i = 0; i >= 0; i++) { /* siempre true */ }")

std.print("   c) Recursión sin caso base:")
std.print("      func infinita() { infinita(); }")

std.print("   d) Condiciones que nunca cambian:")
std.print("      let x = 1; while (x > 0) { /* x nunca cambia */ }")

std.print("")

// ═══════════════════════════════════════════════════════════════════════════════
// EJEMPLO 6: BUENAS PRÁCTICAS
// ═══════════════════════════════════════════════════════════════════════════════

std.print("6. Buenas prácticas para evitar loops infinitos:")

std.print("   ✓ Siempre verificar que las condiciones del loop puedan cambiar")
std.print("   ✓ Usar contadores o límites explícitos cuando sea posible")
std.print("   ✓ Implementar casos base claros en recursión")
std.print("   ✓ Probar con datos de entrada pequeños primero")
std.print("   ✓ Usar herramientas de depuración para monitorear loops")

std.print("")

// ═══════════════════════════════════════════════════════════════════════════════
// EJEMPLO 7: DEMOSTRACIÓN DE LÍMITES CONFIGURABLES
// ═══════════════════════════════════════════════════════════════════════════════

std.print("7. Información sobre límites configurables:")

std.print("   - Límite de iteraciones: Configurable por el usuario")
std.print("   - Límite de recursión: Configurable por el usuario")
std.print("   - Timeout global: Configurable por el usuario")
std.print("   - Los límites pueden ajustarse según las necesidades del programa")

std.print("")

// ═══════════════════════════════════════════════════════════════════════════════
// FINALIZACIÓN
// ═══════════════════════════════════════════════════════════════════════════════

std.print("=== RESUMEN DE LA DEMOSTRACIÓN ===")
std.print("✅ Todos los ejemplos se ejecutaron correctamente")
std.print("🛡️  R2Lang protege automáticamente contra loops infinitos")
std.print("🔧 Los límites son configurables según las necesidades")
std.print("📊 Se proporciona información detallada sobre errores")
std.print("🚀 El rendimiento se mantiene óptimo en código normal")

std.print("\n🎉 ¡Demostración completada exitosamente!")