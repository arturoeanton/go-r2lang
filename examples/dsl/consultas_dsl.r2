// DSL de Consultas
dsl ConsultasDSL {
    // Tokens
    token("SELECCIONAR", "seleccionar")
    token("DE", "de")
    token("DONDE", "donde")
    token("IDENTIFICADOR", "[a-zA-Z_][a-zA-Z0-9_]*")
    token("COMA", ",")
    token("MAYOR", ">")
    token("MENOR", "<")
    token("IGUAL", "==")
    token("NUMERO", "\\d+")
    token("CADENA", "\"[^\"]*\"")
    
    // Reglas
    rule("consulta", ["SELECCIONAR", "campos", "DE", "tabla"], "crear_consulta_simple")
    rule("consulta", ["SELECCIONAR", "campos", "DE", "tabla", "DONDE", "condicion"], "crear_consulta_con_filtro")
    rule("campos", ["IDENTIFICADOR"], "campo_simple")
    rule("campos", ["IDENTIFICADOR", "COMA", "campos"], "campos_multiples")
    rule("tabla", ["IDENTIFICADOR"], "crear_tabla")
    rule("condicion", ["IDENTIFICADOR", "operador", "valor"], "crear_condicion")
    rule("operador", ["MAYOR"], "op_mayor")
    rule("operador", ["MENOR"], "op_menor")
    rule("operador", ["IGUAL"], "op_igual")
    rule("valor", ["NUMERO"], "crear_numero")
    rule("valor", ["CADENA"], "crear_cadena")
    
    // Acciones
    func crear_consulta_simple(sel, campos, de, tabla) {
        return {
            tipo: "consulta",
            campos: campos,
            tabla: tabla,
            sql: "SELECT " + campos + " FROM " + tabla
        }
    }
    
    func crear_consulta_con_filtro(sel, campos, de, tabla, donde, condicion) {
        return {
            tipo: "consulta",
            campos: campos,
            tabla: tabla,
            condicion: condicion,
            sql: "SELECT " + campos + " FROM " + tabla + " WHERE " + condicion
        }
    }
    
    func campo_simple(nombre) {
        return nombre
    }
    
    func campos_multiples(campo, coma, resto) {
        return campo + ", " + resto
    }
    
    func crear_tabla(nombre) {
        return nombre
    }
    
    func crear_condicion(campo, op, valor) {
        return campo + " " + op + " " + valor
    }
    
    func op_mayor(token) { return ">" }
    func op_menor(token) { return "<" }
    func op_igual(token) { return "=" }
    
    func crear_numero(token) {
        return token
    }
    
    func crear_cadena(token) {
        return token
    }
}

func main() {
    console.log("=== DSL Consultas ===")
    
    var parser = ConsultasDSL.use
    
    console.log("Consulta simple:")
    var consulta1 = parser("seleccionar nombre de usuarios")
    console.log("Resultado DSL:", consulta1)
    
    console.log("Consulta con filtro:")
    var consulta2 = parser("seleccionar nombre, edad de usuarios donde edad > 18")
    console.log("Resultado DSL:", consulta2)
    
    console.log("✅ DSL Consultas funcionando!")
}