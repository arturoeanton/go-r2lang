// Ejemplo 30: Guía de Seguridad para Loops en R2Lang
// Una guía práctica para escribir loops seguros y eficientes

std.print("=== EJEMPLO 30: GUÍA DE SEGURIDAD PARA LOOPS ===\n")

// ═══════════════════════════════════════════════════════════════════════════════
// LOOPS SEGUROS: MEJORES PRÁCTICAS
// ═══════════════════════════════════════════════════════════════════════════════

std.print("1. LOOPS WHILE SEGUROS:")
std.print("   ✓ Siempre modificar la variable de control dentro del loop")

// Ejemplo correcto de while loop
let contador = 0
let limite = 5
while (contador < limite) {
    std.print("   Contador seguro:", contador)
    contador++  // ¡IMPORTANTE! Siempre incrementar/modificar la variable
}
std.print("   ✓ While loop completado correctamente\n")

// ═══════════════════════════════════════════════════════════════════════════════

std.print("2. LOOPS FOR SEGUROS:")
std.print("   ✓ Usar condiciones claras y verificar límites")

// Ejemplo correcto de for loop
for (let i = 0; i < 5; i++) {
    std.print("   For seguro:", i)
}
std.print("   ✓ For loop completado correctamente\n")

// ═══════════════════════════════════════════════════════════════════════════════

std.print("3. LOOPS FOR-IN SEGUROS:")
std.print("   ✓ Verificar que el array/map no sea modificado durante la iteración")

let datos = ["a", "b", "c", "d", "e"]
for (indice in datos) {
    std.print("   Elemento", indice, ":", datos[indice])
    // ¡PRECAUCIÓN! No modificar 'datos' aquí dentro
}
std.print("   ✓ For-in loop completado correctamente\n")

// ═══════════════════════════════════════════════════════════════════════════════
// RECURSIÓN SEGURA
// ═══════════════════════════════════════════════════════════════════════════════

std.print("4. RECURSIÓN SEGURA:")
std.print("   ✓ Siempre incluir un caso base claro")

func fibonacci(n) {
    // Caso base: evita recursión infinita
    if (n <= 1) {
        return n
    }
    // Recursión controlada
    return fibonacci(n - 1) + fibonacci(n - 2)
}

std.print("   fibonacci(6) =", fibonacci(6))
std.print("   ✓ Recursión completada correctamente\n")

// ═══════════════════════════════════════════════════════════════════════════════
// PATRONES DE CONTROL DE LÍMITES
// ═══════════════════════════════════════════════════════════════════════════════

std.print("5. PATRONES DE CONTROL DE LÍMITES:")

// Patrón 1: Contador con límite máximo
std.print("   - Patrón 1: Contador con límite máximo")
let maxIteraciones = 10
let iteracion = 0
while (iteracion < maxIteraciones) {
    std.print("     Iteración protegida:", iteracion)
    iteracion++
    
    // Condición de salida adicional
    if (iteracion >= maxIteraciones) {
        std.print("     ✓ Límite alcanzado, saliendo...")
        break
    }
}
std.print("   ✓ Patrón 1 completado\n")

// Patrón 2: Timeout simulado
std.print("   - Patrón 2: Loop con timeout simulado")
let operaciones = 0

while (true) {
    operaciones++
    
    // Simular algún trabajo
    let suma = 0
    for (let i = 0; i < 1000; i++) {
        suma = suma + i
    }
    
    // Verificar timeout simulado
    if (operaciones > 100) {  // Límite de operaciones como proxy para tiempo
        std.print("     ✓ Límite de operaciones alcanzado")
        break
    }
}
std.print("   ✓ Patrón 2 completado,", operaciones, "operaciones\n")

// ═══════════════════════════════════════════════════════════════════════════════
// DETECCIÓN MANUAL DE PROBLEMAS
// ═══════════════════════════════════════════════════════════════════════════════

std.print("6. TÉCNICAS DE DETECCIÓN MANUAL:")

// Técnica 1: Logging de progreso
std.print("   - Técnica 1: Logging de progreso")
let progreso = 0
let totalTrabajo = 8
while (progreso < totalTrabajo) {
    std.print("     Progreso:", progreso, "/", totalTrabajo)
    progreso++
    
    // Simular trabajo
    let calculo = progreso * progreso
}
std.print("   ✓ Técnica 1 completada\n")

// Técnica 2: Validación de entrada
std.print("   - Técnica 2: Validación de entrada")
func procesarArray(arr) {
    // Validar entrada
    if (arr == nil) {
        std.print("     ❌ Error: Array es nil")
        return false
    }
    
    let longitud = std.len(arr)
    if (longitud == 0) {
        std.print("     ⚠️  Advertencia: Array vacío")
        return true
    }
    
    if (longitud > 1000) {
        std.print("     ⚠️  Advertencia: Array muy grande (", longitud, "elementos)")
    }
    
    // Procesar con límites
    let procesados = 0
    for (item in arr) {
        procesados++
        if (procesados > 100) {  // Límite de procesamiento
            std.print("     ✓ Límite de procesamiento alcanzado")
            break
        }
    }
    
    return true
}

let testArray = [1, 2, 3, 4, 5]
procesarArray(testArray)
std.print("   ✓ Técnica 2 completada\n")

// ═══════════════════════════════════════════════════════════════════════════════
// CONSEJOS AVANZADOS
// ═══════════════════════════════════════════════════════════════════════════════

std.print("7. CONSEJOS AVANZADOS:")

std.print("   ✓ Usar variables de control claras y descriptivas")
std.print("   ✓ Implementar múltiples condiciones de salida")
std.print("   ✓ Documentar el propósito y límites de cada loop")
std.print("   ✓ Probar con casos extremos (arrays vacíos, valores grandes)")
std.print("   ✓ Considerar alternativas como map(), filter(), reduce()")

std.print("")

// ═══════════════════════════════════════════════════════════════════════════════
// EJEMPLO DE REFACTORIZACIÓN
// ═══════════════════════════════════════════════════════════════════════════════

std.print("8. EJEMPLO DE REFACTORIZACIÓN:")
std.print("   Transformar un loop potencialmente problemático en uno seguro")

// Versión mejorada de procesamiento
func procesarDatosSeguros(datos) {
    std.print("   Procesando", std.len(datos), "elementos...")
    
    let procesados = 0
    let errores = 0
    let maxProcesar = 1000  // Límite de seguridad
    
    for (i in datos) {
        if (procesados >= maxProcesar) {
            std.print("   ⚠️  Límite de procesamiento alcanzado")
            break
        }
        
        // Procesar elemento
        try {
            let elemento = datos[i]
            // Simular procesamiento
            if (elemento != nil) {
                procesados++
            }
        } catch (error) {
            errores++
            std.print("   ❌ Error procesando elemento", i, ":", error)
        }
    }
    
    std.print("   ✓ Procesamiento completado:")
    std.print("     - Procesados:", procesados)
    std.print("     - Errores:", errores)
    
    return procesados
}

let datosTest = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
procesarDatosSeguros(datosTest)

std.print("")

// ═══════════════════════════════════════════════════════════════════════════════
// RESUMEN DE MEJORES PRÁCTICAS
// ═══════════════════════════════════════════════════════════════════════════════

std.print("=== RESUMEN DE MEJORES PRÁCTICAS ===")
std.print("1. 🎯 Siempre definir condiciones de salida claras")
std.print("2. 🔍 Verificar que las variables de control cambien")
std.print("3. 📊 Implementar límites máximos explícitos")
std.print("4. 🛡️  Usar try-catch para manejar errores")
std.print("5. 📝 Documentar el propósito y límites de loops complejos")
std.print("6. 🧪 Probar con casos extremos y datos inválidos")
std.print("7. 🚀 Considerar alternativas más funcionales cuando sea apropiado")
std.print("8. 📈 Monitorear el rendimiento en loops intensivos")

std.print("\n🎉 ¡Guía de seguridad completada!")
std.print("🛡️  R2Lang te protege automáticamente contra loops infinitos")
std.print("💡 Usar estas prácticas hace tu código más robusto y mantenible")