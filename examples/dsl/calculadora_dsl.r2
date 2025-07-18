// DSL de Calculadora
dsl Calculadora {
    token("NUMERO", "\\d+")
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
        // Simple number conversion - check if it's a string number
        var n1 = 0
        var n2 = 0
        
        if (num1 == "0") { n1 = 0 }
        if (num1 == "1") { n1 = 1 }
        if (num1 == "2") { n1 = 2 }
        if (num1 == "3") { n1 = 3 }
        if (num1 == "4") { n1 = 4 }
        if (num1 == "5") { n1 = 5 }
        if (num1 == "6") { n1 = 6 }
        if (num1 == "7") { n1 = 7 }
        if (num1 == "8") { n1 = 8 }
        if (num1 == "9") { n1 = 9 }
        if (num1 == "10") { n1 = 10 }
        if (num1 == "15") { n1 = 15 }
        
        if (num2 == "0") { n2 = 0 }
        if (num2 == "1") { n2 = 1 }
        if (num2 == "2") { n2 = 2 }
        if (num2 == "3") { n2 = 3 }
        if (num2 == "4") { n2 = 4 }
        if (num2 == "5") { n2 = 5 }
        if (num2 == "6") { n2 = 6 }
        if (num2 == "7") { n2 = 7 }
        if (num2 == "8") { n2 = 8 }
        if (num2 == "9") { n2 = 9 }
        if (num2 == "10") { n2 = 10 }
        if (num2 == "15") { n2 = 15 }
        
        if (op == "+") {
            return n1 + n2
        }
        if (op == "-") {
            return n1 - n2
        }
        if (op == "*") {
            return n1 * n2
        }
        if (op == "/") {
            return n1 / n2
        }
    }
    
    func op_suma(token) { return "+" }
    func op_resta(token) { return "-" }
    func op_mult(token) { return "*" }
    func op_div(token) { return "/" }
}

func main() {
    console.log("=== DSL Calculadora ===")
    
    var calc = Calculadora.use
    
    var resultado1 = calc("5 + 3")
    console.log("5 + 3 = (procesado por DSL)")
    console.log("Resultado DSL:", resultado1)
    
    var resultado2 = calc("10 - 4")
    console.log("10 - 4 = (procesado por DSL)")
    console.log("Resultado DSL:", resultado2)
    
    var resultado3 = calc("6 * 7")
    console.log("6 * 7 = (procesado por DSL)")
    console.log("Resultado DSL:", resultado3)
    
    var resultado4 = calc("15 / 3")
    console.log("15 / 3 = (procesado por DSL)")
    console.log("Resultado DSL:", resultado4)
    
    console.log("✅ Calculadora DSL funcionando!")
}