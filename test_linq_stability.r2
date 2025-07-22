// Test script to reproduce DSL stability issues with LINQ DSL
// This will test the same DSL code multiple times to catch intermittent failures

func from(source) { return source }
func where(data, predicate) { return data.filter(predicate) }
func select(data, projection) { return data.map(projection) }

// DSL estilo LINQ declarativo
dsl LinqQuery {
    token("WORD", "[a-zA-Z_][a-zA-Z0-9_]*")
    token("NUMBER", "[0-9]+")
    token("OPERATOR", "[><=]+")
    
    rule("query", ["WORD", "WORD", "WORD", "WORD", "WORD", "WORD", "OPERATOR", "NUMBER"], "buildQuery")
    rule("query", ["WORD", "WORD", "WORD", "WORD"], "buildSimpleQuery")

    func buildQuery(selectKw, selectField, fromKw, sourceName, whereKw, conditionField, operator, conditionValue) {
        
        let source = context[sourceName];
        if (!source) {
            throw "Fuente de datos '" + sourceName + "' no encontrada en contexto";
        }
        
        // Simula la obtención de datos desde una fuente
        let data = from(source); 
        
        // Filtra los datos según la condición
        let filteredData = where(data, x => {
            if (operator == ">") {
                return x[conditionField] > std.parseInt(conditionValue);
            } else if (operator == "<") {       
                return x[conditionField] < std.parseInt(conditionValue);
            } else if (operator == "=") {
                return x[conditionField] == std.parseInt(conditionValue);
            } else {
                throw "Operador no soportado: " + operator;
            }
        });
        
        
        // Selecciona el campo especificado
        let result = select(filteredData, x => x[selectField]);
        
        return result;
    }
    
    func buildSimpleQuery(selectKw, selectField, fromKw, sourceName) {        
        let source = context[sourceName];
        if (!source) {
            throw "Fuente de datos '" + sourceName + "' no encontrada en contexto";
        }
        
        let data = from(source);
        let result = select(data, x => x[selectField]);
        
        return result;
    }
}

func main() {
    console.log("=== Testing DSL LINQ Stability ===")
    
    let employees = [{name: "John", salary: 50000, department: "IT"}, {name: "Jane", salary: 60000, department: "HR"}, {name: "Mike", salary: 45000, department: "IT"}];
    let context = {employees: employees};
    
    // Test the same query multiple times in rapid succession
    for (let i = 1; i <= 30; i = i + 1) {
        console.log("Run " + i + ":");
        
        try {
            let result = LinqQuery.use("select name from employees where salary > 50000", context);
            console.log("SUCCESS: " + result);
            if (result && result.Output) {
                console.log("Output: " + result.Output);
            } else {
                console.log("No Output property found");
            }
        } catch (e) {
            console.log("ERROR: " + e);
        }
        
        // Also test simple queries
        try {
            let result2 = LinqQuery.use("select department from employees", context);
            console.log("Simple query: " + result2);
        } catch (e) {
            console.log("Simple query ERROR: " + e);
        }
        
        console.log("---");
    }
    
    console.log("LINQ stability test completed")
}