// DSL de Calculadora
dsl Calculadora {
    token("NUMERO", "\\\\d+")
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
        std.print( "Calculando: " + num1 + " " + op + " " + num2);
        var n1 = std.parseInt(num1)
        var n2 =  std.parseInt(num2)
        // Perform the operation based on the operator
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