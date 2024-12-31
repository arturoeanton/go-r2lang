// examplePrint.r2

func main() {
    print("=== r2print Demo ===");

    // 1. println
    println("Este es un mensaje con println.")

    // 2. printf
    printf("Este es un mensaje con printf: %d, %s\n", 42, "respuesta")

    // 3. printError, printWarning, printSuccess
    printError("Este es un mensaje de error.")
    printWarning("Este es un mensaje de advertencia.")
    printSuccess("Este es un mensaje de éxito.")

    // 4. printJSON
    let datos = { nombre: "Juan", edad: 30, hobbies: ["leer", "correr"] }
    printJSON(datos)

    // 5. clearScreen
    // Descomenta la siguiente línea para limpiar la pantalla
    // clearScreen()

    // 6. printTimestamp
    printTimestamp()

    // 7. printHeader
    printHeader("Encabezado Principal")

    // 8. printSeparator
    printSeparator(50)

    // Uso de printBox
    printBox("Caja de Texto", 20)

    // Uso de printAlign
    printAlign("Texto Izquierda", "left", 30)
    printAlign("Texto Derecha", "right", 30)
    printAlign("Texto Centro", "center", 30)

    // Uso de printProgress
    printProgress("Cargando", 10, 200)

    // Uso de printTable
    let tabla = [
        ["Nombre", "Edad", "Ciudad"],
        ["Ana", 25, "Madrid"],
        ["Luis", 28, "Barcelona"],
        ["Marta", 22, "Valencia"]
    ]
    printTable(tabla)

print("=== Fin ===");
}