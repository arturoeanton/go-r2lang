// main.r2

// Función que crea un adder (sumador) con un valor fijo
func makeAdder(x) {
    func add(y) {
        return x + y
    }
    return add
}
func main() {
    // Crear una closure que suma 5
    let add5 = makeAdder(5)

    // Usar la closure para sumar 10
    let result = add5(10)

    // Imprimir el resultado
    print(result) // Debería imprimir 15

    // Crear otra closure que suma 20
    let add20 = makeAdder(20)

    // Usar la segunda closure para sumar 30
    let result2 = add20(30)

    // Imprimir el segundo resultado
    print(result2) // Debería imprimir 50
}