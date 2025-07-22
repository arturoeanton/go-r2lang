// Debug script with proper regex escaping

dsl DebugDSL2 {
    token("NUMBER", "[0-9]+")
    token("PLUS", "\\+")
    token("MINUS", "\\-")      // Properly escaped
    token("MULTIPLY", "\\*")
    token("DIVIDE", "\\/")      // Properly escaped
    
    rule("expression", ["NUMBER", "operator", "NUMBER"], "calculate")
    rule("operator", ["PLUS"], "plus")
    rule("operator", ["MINUS"], "minus")
    rule("operator", ["MULTIPLY"], "multiply")
    rule("operator", ["DIVIDE"], "divide")
    
    func calculate(left, op, right) {
        return left + op + right
    }
    
    func plus(token) { return "+" }
    func minus(token) { return "-" }
    func multiply(token) { return "*" }
    func divide(token) { return "/" }
}

func main() {
    console.log("=== Debug Tokenization with Proper Escaping ===")
    
    console.log("Testing: 5 + 3")
    try {
        let result1 = DebugDSL2.use("5 + 3")
        console.log("SUCCESS: " + result1)
    } catch (e) {
        console.log("ERROR: " + e)
    }
    
    console.log("Testing: 10 - 4")
    try {
        let result2 = DebugDSL2.use("10 - 4")
        console.log("SUCCESS: " + result2)
    } catch (e) {
        console.log("ERROR: " + e)
    }
    
    console.log("Testing: 6 * 7")
    try {
        let result3 = DebugDSL2.use("6 * 7")
        console.log("SUCCESS: " + result3)
    } catch (e) {
        console.log("ERROR: " + e)
    }
    
    console.log("Testing: 15 / 3")
    try {
        let result4 = DebugDSL2.use("15 / 3")
        console.log("SUCCESS: " + result4)
    } catch (e) {
        console.log("ERROR: " + e)
    }
    
    console.log("Debug completed")
}