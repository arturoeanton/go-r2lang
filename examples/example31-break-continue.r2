// Ejemplo 31: Uso de break y continue en R2Lang
// Demuestra las nuevas funcionalidades de control de flujo

std.print("=== EJEMPLO 31: BREAK Y CONTINUE EN R2LANG ===\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 1. BREAK EN WHILE LOOPS
// ═══════════════════════════════════════════════════════════════════════════════

std.print("1. BREAK EN WHILE LOOPS:")

// Ejemplo básico de break
std.print("   - Ejemplo básico de break:")
let i = 0
while (true) {
    if (i >= 5) {
        std.print("     Rompiendo el loop en i =", i)
        break
    }
    std.print("     i =", i)
    i++
}
std.print("   ✓ Loop terminado\n")

// Break con condición específica
std.print("   - Break con condición específica:")
let contador = 0
while (contador < 20) {
    if (contador == 3) {
        std.print("     Encontrado el valor 3, saliendo...")
        break
    }
    std.print("     contador =", contador)
    contador++
}
std.print("   ✓ Loop terminado por break\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 2. CONTINUE EN WHILE LOOPS
// ═══════════════════════════════════════════════════════════════════════════════

std.print("2. CONTINUE EN WHILE LOOPS:")

// Ejemplo básico de continue
std.print("   - Ejemplo básico de continue:")
let j = 0
while (j < 6) {
    j++
    if (j == 3) {
        std.print("     Saltando j =", j)
        continue
    }
    std.print("     Procesando j =", j)
}
std.print("   ✓ Loop completado\n")

// Continue con múltiples condiciones
std.print("   - Continue con múltiples condiciones:")
let num = 0
while (num < 10) {
    num++
    if (num == 2) {
        std.print("     Saltando", num, "(par)")
        continue
    }
    if (num == 4) {
        std.print("     Saltando", num, "(par)")
        continue
    }
    if (num == 6) {
        std.print("     Saltando", num, "(par)")
        continue
    }
    std.print("     Procesando", num)
}
std.print("   ✓ Loop completado\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 3. BREAK Y CONTINUE EN FOR LOOPS
// ═══════════════════════════════════════════════════════════════════════════════

std.print("3. BREAK Y CONTINUE EN FOR LOOPS:")

// Break en for loop
std.print("   - Break en for loop:")
for (let k = 0; k < 10; k++) {
    if (k == 4) {
        std.print("     Break en k =", k)
        break
    }
    std.print("     k =", k)
}
std.print("   ✓ For loop terminado\n")

// Continue en for loop
std.print("   - Continue en for loop:")
for (let m = 0; m < 8; m++) {
    if (m == 2) {
        std.print("     Continue en m =", m)
        continue
    }
    if (m == 5) {
        std.print("     Continue en m =", m)
        continue
    }
    std.print("     m =", m)
}
std.print("   ✓ For loop completado\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 4. BREAK Y CONTINUE EN FOR-IN LOOPS
// ═══════════════════════════════════════════════════════════════════════════════

std.print("4. BREAK Y CONTINUE EN FOR-IN LOOPS:")

// Break en for-in loop
std.print("   - Break en for-in loop:")
let frutas = ["manzana", "banana", "naranja", "uva", "kiwi"]
for (indice in frutas) {
    if (frutas[indice] == "naranja") {
        std.print("     Encontrada naranja, saliendo...")
        break
    }
    std.print("     Fruta:", frutas[indice])
}
std.print("   ✓ For-in loop terminado\n")

// Continue en for-in loop
std.print("   - Continue en for-in loop:")
let numeros = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
for (idx in numeros) {
    let valor = numeros[idx]
    if (valor == 3) {
        std.print("     Saltando el número", valor)
        continue
    }
    if (valor == 7) {
        std.print("     Saltando el número", valor)
        continue
    }
    std.print("     Procesando número:", valor)
}
std.print("   ✓ For-in loop completado\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 5. LOOPS ANIDADOS CON BREAK Y CONTINUE
// ═══════════════════════════════════════════════════════════════════════════════

std.print("5. LOOPS ANIDADOS CON BREAK Y CONTINUE:")

// Break en loop anidado (solo afecta al loop interno)
std.print("   - Break en loop anidado:")
for (let outer = 0; outer < 3; outer++) {
    std.print("     Outer loop:", outer)
    for (let inner = 0; inner < 5; inner++) {
        if (inner == 2) {
            std.print("       Break en inner =", inner)
            break
        }
        std.print("       Inner:", inner)
    }
    std.print("     Outer loop", outer, "completado")
}
std.print("   ✓ Loops anidados completados\n")

// Continue en loop anidado
std.print("   - Continue en loop anidado:")
for (let a = 0; a < 3; a++) {
    std.print("     Loop A:", a)
    for (let b = 0; b < 4; b++) {
        if (b == 1) {
            std.print("       Continue en b =", b)
            continue
        }
        std.print("       B:", b)
    }
    std.print("     Loop A", a, "completado")
}
std.print("   ✓ Loops anidados completados\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 6. CASOS PRÁCTICOS DE USO
// ═══════════════════════════════════════════════════════════════════════════════

std.print("6. CASOS PRÁCTICOS DE USO:")

// Búsqueda con break
std.print("   - Búsqueda con break:")
let palabras = ["hola", "mundo", "R2Lang", "programación", "break"]
let buscada = "R2Lang"
let encontrada = false
for (p in palabras) {
    if (palabras[p] == buscada) {
        std.print("     ✓ Encontrada:", palabras[p])
        encontrada = true
        break
    }
    std.print("     Revisando:", palabras[p])
}
if (encontrada == false) {
    std.print("     ❌ No encontrada")
}
std.print("   ✓ Búsqueda completada\n")

// Filtrado con continue
std.print("   - Filtrado con continue:")
let valores = [1, 2, -3, 4, -5, 6, -7, 8, 9, -10]
let suma = 0
let procesados = 0
for (v in valores) {
    let valor = valores[v]
    if (valor < 0) {
        std.print("     Saltando valor negativo:", valor)
        continue
    }
    suma = suma + valor
    procesados++
    std.print("     Sumando valor positivo:", valor, "- suma actual:", suma)
}
std.print("   ✓ Procesados", procesados, "valores positivos, suma total:", suma, "\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 7. MEJORES PRÁCTICAS
// ═══════════════════════════════════════════════════════════════════════════════

std.print("7. MEJORES PRÁCTICAS:")

std.print("   ✓ Usar break para salir de loops cuando se cumple una condición")
std.print("   ✓ Usar continue para saltar iteraciones específicas")
std.print("   ✓ Break y continue solo afectan al loop más interno")
std.print("   ✓ Documentar claramente por qué se usa break o continue")
std.print("   ✓ Considerar alternativas como funciones auxiliares para lógica compleja")
std.print("   ✓ Evitar break y continue anidados profundos para mantener legibilidad")

std.print("")

// ═══════════════════════════════════════════════════════════════════════════════
// FINALIZACIÓN
// ═══════════════════════════════════════════════════════════════════════════════

std.print("=== RESUMEN DEL EJEMPLO ===")
std.print("✅ Break: Termina el loop inmediatamente")
std.print("✅ Continue: Salta a la siguiente iteración")
std.print("✅ Funcionan en while, for, y for-in loops")
std.print("✅ Solo afectan al loop más interno en casos anidados")
std.print("✅ Mejoran el control de flujo y legibilidad del código")

std.print("\n🎉 ¡Ejemplo de break y continue completado!")
std.print("🚀 Ahora puedes usar estas potentes herramientas de control de flujo en R2Lang")