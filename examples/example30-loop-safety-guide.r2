// Ejemplo 30: Guía de Seguridad para Loops en R2Lang
// Una guía práctica para escribir loops seguros y eficientes

print("=== EJEMPLO 30: GUÍA DE SEGURIDAD PARA LOOPS ===\n")

// ═══════════════════════════════════════════════════════════════════════════════
// LOOPS SEGUROS: MEJORES PRÁCTICAS
// ═══════════════════════════════════════════════════════════════════════════════

print("1. LOOPS WHILE SEGUROS:")
print("   ✓ Siempre modificar la variable de control dentro del loop")

// Ejemplo correcto de while loop
let contador = 0
let limite = 5
while (contador < limite) {
    print("   Contador seguro:", contador)
    contador++  // ¡IMPORTANTE! Siempre incrementar/modificar la variable
}
print("   ✓ While loop completado correctamente\n")

// ═══════════════════════════════════════════════════════════════════════════════

print("2. LOOPS FOR SEGUROS:")
print("   ✓ Usar condiciones claras y verificar límites")

// Ejemplo correcto de for loop
for (let i = 0; i < 5; i++) {
    print("   For seguro:", i)
}
print("   ✓ For loop completado correctamente\n")

// ═══════════════════════════════════════════════════════════════════════════════

print("3. LOOPS FOR-IN SEGUROS:")
print("   ✓ Verificar que el array/map no sea modificado durante la iteración")

let datos = ["a", "b", "c", "d", "e"]
for (indice in datos) {
    print("   Elemento", indice, ":", datos[indice])
    // ¡PRECAUCIÓN! No modificar 'datos' aquí dentro
}
print("   ✓ For-in loop completado correctamente\n")

// ═══════════════════════════════════════════════════════════════════════════════
// RECURSIÓN SEGURA
// ═══════════════════════════════════════════════════════════════════════════════

print("4. RECURSIÓN SEGURA:")
print("   ✓ Siempre incluir un caso base claro")

func fibonacci(n) {
    // Caso base: evita recursión infinita
    if (n <= 1) {
        return n
    }
    // Recursión controlada
    return fibonacci(n - 1) + fibonacci(n - 2)
}

print("   fibonacci(6) =", fibonacci(6))
print("   ✓ Recursión completada correctamente\n")

// ═══════════════════════════════════════════════════════════════════════════════
// PATRONES DE CONTROL DE LÍMITES
// ═══════════════════════════════════════════════════════════════════════════════

print("5. PATRONES DE CONTROL DE LÍMITES:")

// Patrón 1: Contador con límite máximo
print("   - Patrón 1: Contador con límite máximo")
let maxIteraciones = 10
let iteracion = 0
while (iteracion < maxIteraciones) {
    print("     Iteración protegida:", iteracion)
    iteracion++
    
    // Condición de salida adicional
    if (iteracion >= maxIteraciones) {
        print("     ✓ Límite alcanzado, saliendo...")
        break
    }
}
print("   ✓ Patrón 1 completado\n")

// Patrón 2: Timeout simulado
print("   - Patrón 2: Loop con timeout simulado")
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
        print("     ✓ Límite de operaciones alcanzado")
        break
    }
}
print("   ✓ Patrón 2 completado,", operaciones, "operaciones\n")

// ═══════════════════════════════════════════════════════════════════════════════
// DETECCIÓN MANUAL DE PROBLEMAS
// ═══════════════════════════════════════════════════════════════════════════════

print("6. TÉCNICAS DE DETECCIÓN MANUAL:")

// Técnica 1: Logging de progreso
print("   - Técnica 1: Logging de progreso")
let progreso = 0
let totalTrabajo = 8
while (progreso < totalTrabajo) {
    print("     Progreso:", progreso, "/", totalTrabajo)
    progreso++
    
    // Simular trabajo
    let calculo = progreso * progreso
}
print("   ✓ Técnica 1 completada\n")

// Técnica 2: Validación de entrada
print("   - Técnica 2: Validación de entrada")
func procesarArray(arr) {
    // Validar entrada
    if (arr == nil) {
        print("     ❌ Error: Array es nil")
        return false
    }
    
    let longitud = len(arr)
    if (longitud == 0) {
        print("     ⚠️  Advertencia: Array vacío")
        return true
    }
    
    if (longitud > 1000) {
        print("     ⚠️  Advertencia: Array muy grande (", longitud, "elementos)")
    }
    
    // Procesar con límites
    let procesados = 0
    for (item in arr) {
        procesados++
        if (procesados > 100) {  // Límite de procesamiento
            print("     ✓ Límite de procesamiento alcanzado")
            break
        }
    }
    
    return true
}

let testArray = [1, 2, 3, 4, 5]
procesarArray(testArray)
print("   ✓ Técnica 2 completada\n")

// ═══════════════════════════════════════════════════════════════════════════════
// CONSEJOS AVANZADOS
// ═══════════════════════════════════════════════════════════════════════════════

print("7. CONSEJOS AVANZADOS:")

print("   ✓ Usar variables de control claras y descriptivas")
print("   ✓ Implementar múltiples condiciones de salida")
print("   ✓ Documentar el propósito y límites de cada loop")
print("   ✓ Probar con casos extremos (arrays vacíos, valores grandes)")
print("   ✓ Considerar alternativas como map(), filter(), reduce()")

print("")

// ═══════════════════════════════════════════════════════════════════════════════
// EJEMPLO DE REFACTORIZACIÓN
// ═══════════════════════════════════════════════════════════════════════════════

print("8. EJEMPLO DE REFACTORIZACIÓN:")
print("   Transformar un loop potencialmente problemático en uno seguro")

// Versión mejorada de procesamiento
func procesarDatosSeguros(datos) {
    print("   Procesando", len(datos), "elementos...")
    
    let procesados = 0
    let errores = 0
    let maxProcesar = 1000  // Límite de seguridad
    
    for (i in datos) {
        if (procesados >= maxProcesar) {
            print("   ⚠️  Límite de procesamiento alcanzado")
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
            print("   ❌ Error procesando elemento", i, ":", error)
        }
    }
    
    print("   ✓ Procesamiento completado:")
    print("     - Procesados:", procesados)
    print("     - Errores:", errores)
    
    return procesados
}

let datosTest = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
procesarDatosSeguros(datosTest)

print("")

// ═══════════════════════════════════════════════════════════════════════════════
// RESUMEN DE MEJORES PRÁCTICAS
// ═══════════════════════════════════════════════════════════════════════════════

print("=== RESUMEN DE MEJORES PRÁCTICAS ===")
print("1. 🎯 Siempre definir condiciones de salida claras")
print("2. 🔍 Verificar que las variables de control cambien")
print("3. 📊 Implementar límites máximos explícitos")
print("4. 🛡️  Usar try-catch para manejar errores")
print("5. 📝 Documentar el propósito y límites de loops complejos")
print("6. 🧪 Probar con casos extremos y datos inválidos")
print("7. 🚀 Considerar alternativas más funcionales cuando sea apropiado")
print("8. 📈 Monitorear el rendimiento en loops intensivos")

print("\n🎉 ¡Guía de seguridad completada!")
print("🛡️  R2Lang te protege automáticamente contra loops infinitos")
print("💡 Usar estas prácticas hace tu código más robusto y mantenible")