// Debug script to understand tokenization issues

dsl DebugDSL {
    token("NUMBER", "[0-9]+")
    token("PLUS", "\\+")
    token("MINUS", "-")
    token("MULTIPLY", "\\*")
    token("DIVIDE", "/")
    
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
    console.log("=== Debug Tokenization ===")
    
    console.log("Testing: 5 + 3")
    try {
        let result1 = DebugDSL.use("5 + 3")
        console.log("SUCCESS: " + result1)
    } catch (e) {
        console.log("ERROR: " + e)
    }
    
    console.log("Testing: 10 - 4")
    try {
        let result2 = DebugDSL.use("10 - 4")
        console.log("SUCCESS: " + result2)
    } catch (e) {
        console.log("ERROR: " + e)
    }
    
    console.log("Testing: 6 * 7")
    try {
        let result3 = DebugDSL.use("6 * 7")
        console.log("SUCCESS: " + result3)
    } catch (e) {
        console.log("ERROR: " + e)
    }
    
    console.log("Testing: 15 / 3")
    try {
        let result4 = DebugDSL.use("15 / 3")
        console.log("SUCCESS: " + result4)
    } catch (e) {
        console.log("ERROR: " + e)
    }
    
    console.log("Debug completed")
}