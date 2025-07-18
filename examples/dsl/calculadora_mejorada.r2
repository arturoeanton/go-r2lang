// DSL de Calculadora con mejores resultados
dsl Calculadora {
    token("NUMERO", "[0-9]+")
    token("SUMA", "\\+")
    token("RESTA", "-")
    token("MULT", "\\*")
    token("DIV", "/")
    
    rule("operacion", ["NUMERO", "operador", "NUMERO"], "calcular")
    rule("operador", ["SUMA"], "op_suma")
    rule("operador", ["RESTA"], "op_resta")
    rule("operador", ["MULT"], "op_mult")
    rule("operador", ["DIV"], "op_div")
    
    func calcular(num1, op, num2) {
        std.print("Calculando: " + num1 + " " + op + " " + num2)
        var n1 = std.parseInt(num1)
        var n2 = std.parseInt(num2)
        var resultado = 0
        
        if (op == "+") {
            resultado = n1 + n2
        }
        if (op == "-") {
            resultado = n1 - n2
        }
        if (op == "*") {
            resultado = n1 * n2
        }
        if (op == "/") {
            resultado = n1 / n2
        }
        
        std.print("Resultado: " + resultado)
        return resultado
    }
    
    func op_suma(token) { return "+" }
    func op_resta(token) { return "-" }
    func op_mult(token) { return "*" }
    func op_div(token) { return "/" }
}

func main() {
    console.log("=== DSL Calculadora Mejorada ===")
    
    var calc = Calculadora.use
    
    // Ejemplo 1: Suma
    console.log("Operación: '5 + 3'")
    var result1 = calc("5 + 3")
    console.log("Resultado completo:", result1)
    if (result1 && result1.Output) {
        console.log("Solo el resultado:", result1.Output)
    }
    
    // Ejemplo 2: Multiplicación
    console.log("Operación: '12 * 7'")
    var result2 = calc("12 * 7")
    console.log("Resultado completo:", result2)
    if (result2 && result2.Output) {
        console.log("Solo el resultado:", result2.Output)
    }
    
    // Ejemplo 3: División
    console.log("Operación: '100 / 5'")
    var result3 = calc("100 / 5")
    console.log("Resultado completo:", result3)
    if (result3 && result3.Output) {
        console.log("Solo el resultado:", result3.Output)
    }
    
    console.log("✅ DSL Calculadora mejorada funcionando!")
}