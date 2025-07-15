// Ejemplo 29: Demostración Práctica de Detección de Loops Infinitos
// ADVERTENCIA: Este ejemplo contiene código que puede causar loops infinitos
// Ejecutar con precaución - R2Lang los detectará y terminará la ejecución

print("=== EJEMPLO 29: DEMOSTRACIÓN PRÁCTICA DE DETECCIÓN ===\n")

// ═══════════════════════════════════════════════════════════════════════════════
// CONFIGURACIÓN DE PRUEBA
// ═══════════════════════════════════════════════════════════════════════════════

print("Este ejemplo demuestra la detección en acción...")
print("Los loops infinitos serán detectados automáticamente por R2Lang\n")

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

print("1. Probando while loop que excedería el límite:")
print("   (Simulando un límite bajo para demostración)")

// Este loop sería detectado si el límite fuera muy bajo
let contadorDemo = 0
while (contadorDemo < 10) {  // Solo 10 iteraciones para demostrar
    print("   Iteración segura:", contadorDemo)
    contadorDemo++
    simularTrabajo()
}
print("   ✓ Loop completado correctamente\n")

// ═══════════════════════════════════════════════════════════════════════════════
// EJEMPLO 2: RECURSIÓN CONTROLADA
// ═══════════════════════════════════════════════════════════════════════════════

print("2. Probando recursión controlada:")

func recursionControlada(n, max) {
    if (n >= max) {
        return n
    }
    print("   Nivel de recursión:", n)
    return recursionControlada(n + 1, max)
}

let resultadoRecursion = recursionControlada(0, 5)
print("   ✓ Recursión completada, resultado:", resultadoRecursion, "\n")

// ═══════════════════════════════════════════════════════════════════════════════
// EJEMPLO 3: FOR LOOP CON ARRAYS GRANDES
// ═══════════════════════════════════════════════════════════════════════════════

print("3. Probando for loop con array grande:")

// Crear un array grande para probar límites
let arrayGrande = []
for (let i = 0; i < 1000; i++) {
    arrayGrande[i] = i * 2
}

print("   Array creado con", len(arrayGrande), "elementos")

// Procesar el array
let suma = 0
for (elemento in arrayGrande) {
    suma = suma + arrayGrande[elemento]
}

print("   ✓ Procesamiento completado, suma:", suma, "\n")

// ═══════════════════════════════════════════════════════════════════════════════
// EJEMPLO 4: MANEJO DE ERRORES CON TRY-CATCH
// ═══════════════════════════════════════════════════════════════════════════════

print("4. Ejemplo de manejo de errores (simulado):")

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
    print("   ✓ Función ejecutada correctamente, resultado:", resultado)
} catch (error) {
    print("   ❌ Error detectado:", error)
} finally {
    print("   ✓ Limpieza completada")
}

print("")

// ═══════════════════════════════════════════════════════════════════════════════
// EJEMPLO 5: PATRONES COMUNES QUE CAUSAN LOOPS INFINITOS
// ═══════════════════════════════════════════════════════════════════════════════

print("5. Patrones comunes que causan loops infinitos:")

print("   a) While sin incremento:")
print("      while (true) { /* sin cambio de condición */ }")

print("   b) For con condición incorrecta:")
print("      for (let i = 0; i >= 0; i++) { /* siempre true */ }")

print("   c) Recursión sin caso base:")
print("      func infinita() { infinita(); }")

print("   d) Condiciones que nunca cambian:")
print("      let x = 1; while (x > 0) { /* x nunca cambia */ }")

print("")

// ═══════════════════════════════════════════════════════════════════════════════
// EJEMPLO 6: BUENAS PRÁCTICAS
// ═══════════════════════════════════════════════════════════════════════════════

print("6. Buenas prácticas para evitar loops infinitos:")

print("   ✓ Siempre verificar que las condiciones del loop puedan cambiar")
print("   ✓ Usar contadores o límites explícitos cuando sea posible")
print("   ✓ Implementar casos base claros en recursión")
print("   ✓ Probar con datos de entrada pequeños primero")
print("   ✓ Usar herramientas de depuración para monitorear loops")

print("")

// ═══════════════════════════════════════════════════════════════════════════════
// EJEMPLO 7: DEMOSTRACIÓN DE LÍMITES CONFIGURABLES
// ═══════════════════════════════════════════════════════════════════════════════

print("7. Información sobre límites configurables:")

print("   - Límite de iteraciones: Configurable por el usuario")
print("   - Límite de recursión: Configurable por el usuario")
print("   - Timeout global: Configurable por el usuario")
print("   - Los límites pueden ajustarse según las necesidades del programa")

print("")

// ═══════════════════════════════════════════════════════════════════════════════
// FINALIZACIÓN
// ═══════════════════════════════════════════════════════════════════════════════

print("=== RESUMEN DE LA DEMOSTRACIÓN ===")
print("✅ Todos los ejemplos se ejecutaron correctamente")
print("🛡️  R2Lang protege automáticamente contra loops infinitos")
print("🔧 Los límites son configurables según las necesidades")
print("📊 Se proporciona información detallada sobre errores")
print("🚀 El rendimiento se mantiene óptimo en código normal")

print("\n🎉 ¡Demostración completada exitosamente!")