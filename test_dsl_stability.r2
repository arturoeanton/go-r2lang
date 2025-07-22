// Test script to reproduce DSL stability issues
// This will test the same DSL code multiple times to catch intermittent failures

func main() {
    console.log("=== Testing DSL Stability ===")
    
    // Datos de prueba
    let employees = [
        {name: "Alice", salary: 60000}, 
        {name: "Bob", salary: 45000}, 
        {name: "Jane", salary: 55000}
    ];
    
    // DSL definition - similar to the failing one
    dsl QueryDSL {
        token("WORD", "[a-zA-Z_][a-zA-Z0-9_]*")
        token("NUMBER", "[0-9]+")
        token("OPERATOR", "[><=]+")
        
        rule("query", ["WORD", "WORD", "WORD", "WORD", "WORD", "WORD", "OPERATOR", "NUMBER"], "buildQuery")
        
        func buildQuery(selectKw, selectField, fromKw, sourceName, whereKw, conditionField, operator, conditionValue) {
            let source = context[sourceName];
            if (!source) {
                return "Source not found";
            }
            
            let result = [];
            for (let i = 0; i < source.length; i = i + 1) {
                let item = source[i];
                let condValue = std.parseInt(conditionValue);
                
                if (operator == ">" && item[conditionField] > condValue) {
                    result.push(item[selectField]);
                }
            }
            
            return result;
        }
    }
    
    let context = {employees: employees};
    
    // Run the same DSL code multiple times to check for stability
    for (let i = 1; i <= 20; i = i + 1) {
        console.log("Run " + i + ":");
        
        try {
            let result = QueryDSL.use("select name from employees where salary > 50000", context);
            console.log("SUCCESS: " + result);
            if (result && result.Output) {
                console.log("Output: " + result.Output);
            }
        } catch (e) {
            console.log("ERROR: " + e);
        }
        
        console.log("---");
    }
    
    console.log("Stability test completed")
}