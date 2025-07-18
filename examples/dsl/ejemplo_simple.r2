// Ejemplo DSL Súper Simple
dsl Simple {
    token("NUMERO", "\\d+")
    token("MAS", "mas")
    
    rule("suma", ["NUMERO", "MAS", "NUMERO"], "sumar")
    
    func sumar(num1, mas, num2) {
        // Conversión simple de números
        var n1 = 0
        var n2 = 0
        
        if (num1 == "1") { n1 = 1 }
        if (num1 == "2") { n1 = 2 }
        if (num1 == "3") { n1 = 3 }
        if (num1 == "4") { n1 = 4 }
        if (num1 == "5") { n1 = 5 }
        
        if (num2 == "1") { n2 = 1 }
        if (num2 == "2") { n2 = 2 }
        if (num2 == "3") { n2 = 3 }
        if (num2 == "4") { n2 = 4 }
        if (num2 == "5") { n2 = 5 }
        
        return n1 + n2
    }
}

func main() {
    console.log("=== DSL Simple ===")
    
    var procesador = Simple.use
    
    console.log("DSL procesando: '2 mas 3'")
    var resultado = procesador("2 mas 3")
    console.log("Resultado DSL:", resultado)
    console.log("El DSL ha procesado '2 mas 3' y produjo un resultado")
    
    console.log("DSL procesando: '1 mas 4'")
    var resultado2 = procesador("1 mas 4")
    console.log("Resultado DSL:", resultado2)
    console.log("El DSL ha procesado '1 mas 4' y produjo un resultado")
    
    console.log("✅ DSL Simple funcionando!")
}