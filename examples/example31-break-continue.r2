// Ejemplo 31: Uso de break y continue en R2Lang
// Demuestra las nuevas funcionalidades de control de flujo

print("=== EJEMPLO 31: BREAK Y CONTINUE EN R2LANG ===\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 1. BREAK EN WHILE LOOPS
// ═══════════════════════════════════════════════════════════════════════════════

print("1. BREAK EN WHILE LOOPS:")

// Ejemplo básico de break
print("   - Ejemplo básico de break:")
let i = 0
while (true) {
    if (i >= 5) {
        print("     Rompiendo el loop en i =", i)
        break
    }
    print("     i =", i)
    i++
}
print("   ✓ Loop terminado\n")

// Break con condición específica
print("   - Break con condición específica:")
let contador = 0
while (contador < 20) {
    if (contador == 3) {
        print("     Encontrado el valor 3, saliendo...")
        break
    }
    print("     contador =", contador)
    contador++
}
print("   ✓ Loop terminado por break\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 2. CONTINUE EN WHILE LOOPS
// ═══════════════════════════════════════════════════════════════════════════════

print("2. CONTINUE EN WHILE LOOPS:")

// Ejemplo básico de continue
print("   - Ejemplo básico de continue:")
let j = 0
while (j < 6) {
    j++
    if (j == 3) {
        print("     Saltando j =", j)
        continue
    }
    print("     Procesando j =", j)
}
print("   ✓ Loop completado\n")

// Continue con múltiples condiciones
print("   - Continue con múltiples condiciones:")
let num = 0
while (num < 10) {
    num++
    if (num == 2) {
        print("     Saltando", num, "(par)")
        continue
    }
    if (num == 4) {
        print("     Saltando", num, "(par)")
        continue
    }
    if (num == 6) {
        print("     Saltando", num, "(par)")
        continue
    }
    print("     Procesando", num)
}
print("   ✓ Loop completado\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 3. BREAK Y CONTINUE EN FOR LOOPS
// ═══════════════════════════════════════════════════════════════════════════════

print("3. BREAK Y CONTINUE EN FOR LOOPS:")

// Break en for loop
print("   - Break en for loop:")
for (let k = 0; k < 10; k++) {
    if (k == 4) {
        print("     Break en k =", k)
        break
    }
    print("     k =", k)
}
print("   ✓ For loop terminado\n")

// Continue en for loop
print("   - Continue en for loop:")
for (let m = 0; m < 8; m++) {
    if (m == 2) {
        print("     Continue en m =", m)
        continue
    }
    if (m == 5) {
        print("     Continue en m =", m)
        continue
    }
    print("     m =", m)
}
print("   ✓ For loop completado\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 4. BREAK Y CONTINUE EN FOR-IN LOOPS
// ═══════════════════════════════════════════════════════════════════════════════

print("4. BREAK Y CONTINUE EN FOR-IN LOOPS:")

// Break en for-in loop
print("   - Break en for-in loop:")
let frutas = ["manzana", "banana", "naranja", "uva", "kiwi"]
for (indice in frutas) {
    if (frutas[indice] == "naranja") {
        print("     Encontrada naranja, saliendo...")
        break
    }
    print("     Fruta:", frutas[indice])
}
print("   ✓ For-in loop terminado\n")

// Continue en for-in loop
print("   - Continue en for-in loop:")
let numeros = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
for (idx in numeros) {
    let valor = numeros[idx]
    if (valor == 3) {
        print("     Saltando el número", valor)
        continue
    }
    if (valor == 7) {
        print("     Saltando el número", valor)
        continue
    }
    print("     Procesando número:", valor)
}
print("   ✓ For-in loop completado\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 5. LOOPS ANIDADOS CON BREAK Y CONTINUE
// ═══════════════════════════════════════════════════════════════════════════════

print("5. LOOPS ANIDADOS CON BREAK Y CONTINUE:")

// Break en loop anidado (solo afecta al loop interno)
print("   - Break en loop anidado:")
for (let outer = 0; outer < 3; outer++) {
    print("     Outer loop:", outer)
    for (let inner = 0; inner < 5; inner++) {
        if (inner == 2) {
            print("       Break en inner =", inner)
            break
        }
        print("       Inner:", inner)
    }
    print("     Outer loop", outer, "completado")
}
print("   ✓ Loops anidados completados\n")

// Continue en loop anidado
print("   - Continue en loop anidado:")
for (let a = 0; a < 3; a++) {
    print("     Loop A:", a)
    for (let b = 0; b < 4; b++) {
        if (b == 1) {
            print("       Continue en b =", b)
            continue
        }
        print("       B:", b)
    }
    print("     Loop A", a, "completado")
}
print("   ✓ Loops anidados completados\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 6. CASOS PRÁCTICOS DE USO
// ═══════════════════════════════════════════════════════════════════════════════

print("6. CASOS PRÁCTICOS DE USO:")

// Búsqueda con break
print("   - Búsqueda con break:")
let palabras = ["hola", "mundo", "R2Lang", "programación", "break"]
let buscada = "R2Lang"
let encontrada = false
for (p in palabras) {
    if (palabras[p] == buscada) {
        print("     ✓ Encontrada:", palabras[p])
        encontrada = true
        break
    }
    print("     Revisando:", palabras[p])
}
if (encontrada == false) {
    print("     ❌ No encontrada")
}
print("   ✓ Búsqueda completada\n")

// Filtrado con continue
print("   - Filtrado con continue:")
let valores = [1, 2, -3, 4, -5, 6, -7, 8, 9, -10]
let suma = 0
let procesados = 0
for (v in valores) {
    let valor = valores[v]
    if (valor < 0) {
        print("     Saltando valor negativo:", valor)
        continue
    }
    suma = suma + valor
    procesados++
    print("     Sumando valor positivo:", valor, "- suma actual:", suma)
}
print("   ✓ Procesados", procesados, "valores positivos, suma total:", suma, "\n")

// ═══════════════════════════════════════════════════════════════════════════════
// 7. MEJORES PRÁCTICAS
// ═══════════════════════════════════════════════════════════════════════════════

print("7. MEJORES PRÁCTICAS:")

print("   ✓ Usar break para salir de loops cuando se cumple una condición")
print("   ✓ Usar continue para saltar iteraciones específicas")
print("   ✓ Break y continue solo afectan al loop más interno")
print("   ✓ Documentar claramente por qué se usa break o continue")
print("   ✓ Considerar alternativas como funciones auxiliares para lógica compleja")
print("   ✓ Evitar break y continue anidados profundos para mantener legibilidad")

print("")

// ═══════════════════════════════════════════════════════════════════════════════
// FINALIZACIÓN
// ═══════════════════════════════════════════════════════════════════════════════

print("=== RESUMEN DEL EJEMPLO ===")
print("✅ Break: Termina el loop inmediatamente")
print("✅ Continue: Salta a la siguiente iteración")
print("✅ Funcionan en while, for, y for-in loops")
print("✅ Solo afectan al loop más interno en casos anidados")
print("✅ Mejoran el control de flujo y legibilidad del código")

print("\n🎉 ¡Ejemplo de break y continue completado!")
print("🚀 Ahora puedes usar estas potentes herramientas de control de flujo en R2Lang")