// DSL de Calculadora con soporte para contexto
dsl Calculadora {
    token("VARIABLE", "[a-zA-Z][a-zA-Z0-9]*")
    token("NUMERO", "[0-9]+")
    token("SUMA", "\\+")
    token("RESTA", "-")
    token("MULT", "\\*")
    token("DIV", "/")
    
    rule("operacion", ["operando", "operador", "operando"], "calcular")
    rule("operando", ["NUMERO"], "numero")
    rule("operando", ["VARIABLE"], "variable")
    rule("operador", ["SUMA"], "op_suma")
    rule("operador", ["RESTA"], "op_resta")
    rule("operador", ["MULT"], "op_mult")
    rule("operador", ["DIV"], "op_div")
    
    func calcular(val1, op, val2) {
        std.print("Calculando: " + val1 + " " + op + " " + val2)
        // Convertir a número, manejando tanto strings como números
        var n1 = (std.typeOf(val1) == "string") ? std.parseInt(val1) : val1
        var n2 = (std.typeOf(val2) == "string") ? std.parseInt(val2) : val2
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
    
    func numero(token) { 
        return token 
    }
    
    func variable(token) {
        // Obtener valor del contexto
        if (context[token] != nil) {
            std.print("Variable " + token + " = " + context[token])
            return context[token]
        } else {
            std.print("Variable " + token + " no encontrada en contexto")
            return "0"
        }
    }
    
    func op_suma(token) { return "+" }
    func op_resta(token) { return "-" }
    func op_mult(token) { return "*" }
    func op_div(token) { return "/" }
}

func main() {
    console.log("=== DSL Calculadora con Contexto ===")
    
    var calc = Calculadora
    var useMethod = calc.use
    
    // Ejemplo 1: Suma simple sin contexto
    console.log("\n1. Operación: '5 + 3'")
    var result1 = useMethod("5 + 3")
    console.log("Resultado completo:", result1)
    if (result1 && result1.Output) {
        console.log("Solo el resultado:", result1.Output)
    }
    
    // Ejemplo 2: Usando variables con contexto
    console.log("\n2. Operación con variables: 'a + b'")
    var context1 = {a: 10, b: 20}
    var result2 = useMethod("a + b", context1)
    console.log("Resultado completo:", result2)
    if (result2 && result2.Output) {
        console.log("Solo el resultado:", result2.Output)
    }
    
    // Ejemplo 3: Multiplicación con contexto
    console.log("\n3. Operación: 'x * y'")
    var context2 = {x: 7, y: 8}
    var result3 = useMethod("x * y", context2)
    console.log("Resultado completo:", result3)
    if (result3 && result3.Output) {
        console.log("Solo el resultado:", result3.Output)
    }
    
    // Ejemplo 4: Operaciones mixtas
    console.log("\n4. Operación mixta: 'a + 15'")
    var context3 = {a: 25}
    var result4 = useMethod("a + 15", context3)
    console.log("Resultado completo:", result4)
    if (result4 && result4.Output) {
        console.log("Solo el resultado:", result4.Output)
    }
    
    console.log("\n✅ DSL Calculadora con contexto funcionando!")
}