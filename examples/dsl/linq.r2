// Funciones helper para LINQ
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
    console.log("=== DSL LINQ con Contexto ===")
    
    // Datos de prueba - mapas multilinea
    let data = [
        {name: "Alice", age: 30, team: "A"},
        {name: "Bob", age: 25, team: "B"}, 
        {name: "Carol", age: 35, team: "A"},
        {name: "David", age: 28, team: "B"},
        {name: "Eve", age: 32, team: "A"}
    ];
    
    // Ejemplo 1: Uso funcional tradicional con pipeline (una línea)
    console.log("\n1. Uso funcional tradicional:")
    let result1 = data 
        |> (x => where(x, person => person.age > 30)) 
        |> (x => select(x, person => person.name));
    console.log("Resultado pipeline:", result1);

    // Ejemplo 2: DSL LINQ con contexto
    console.log("\n2. DSL LINQ con contexto:")
    let context1 = {data: data};
    let linq = LinqQuery;
    let result2 = linq.use("select name from data where age > 30", context1);
    console.log("Resultado DSL LINQ:", result2);
    if (result2 && result2.Output) {
        console.log("Solo el resultado:", result2.Output);
    }
    
    // Ejemplo 3: Query simple sin WHERE
    console.log("\n3. Query simple sin filtro:")
    let result3 = linq.use("select team from data", context1);
    console.log("Resultado DSL simple:", result3);
    if (result3 && result3.Output) {
        console.log("Solo el resultado:", result3.Output);
    }
    
    // Ejemplo 4: Diferentes contextos
    console.log("\n4. Contexto diferente:")
    let employees = [{name: "John", salary: 50000, department: "IT"}, {name: "Jane", salary: 60000, department: "HR"}, {name: "Mike", salary: 45000, department: "IT"}];
    
    let context2 = {employees: employees};
    let result4 = linq.use("select name from employees where salary > 50000", context2);
    console.log("Resultado contexto 2:", result4);
    if (result4 && result4.Output) {
        console.log("Solo el resultado:", result4.Output);
    }
    
    console.log("\n✅ DSL LINQ con contexto funcionando!")
}