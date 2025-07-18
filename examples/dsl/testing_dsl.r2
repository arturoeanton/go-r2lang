// DSL para Testing
dsl TestingDSL {
    token("DADO", "dado")
    token("CUANDO", "cuando")
    token("ENTONCES", "entonces")
    token("Y", "y")
    token("QUE", "que")
    token("ACCION", "[a-zA-Z_][a-zA-Z0-9_]*")
    token("CADENA", "\"[^\"]*\"")
    
    rule("escenario", ["pasos"], "crear_escenario")
    rule("pasos", ["paso"], "paso_simple")
    rule("pasos", ["paso", "pasos"], "pasos_multiples")
    rule("paso", ["DADO", "QUE", "condicion"], "paso_dado")
    rule("paso", ["CUANDO", "accion"], "paso_cuando")
    rule("paso", ["ENTONCES", "expectativa"], "paso_entonces")
    rule("paso", ["Y", "condicion"], "paso_y")
    rule("condicion", ["CADENA"], "crear_condicion")
    rule("accion", ["ACCION", "CADENA"], "crear_accion")
    rule("expectativa", ["CADENA"], "crear_expectativa")
    
    func crear_escenario(pasos) {
        return {
            tipo: "escenario",
            pasos: pasos
        }
    }
    
    func paso_simple(paso) {
        return [paso]
    }
    
    func pasos_multiples(paso, resto) {
        var result = [paso]
        return result.concat(resto)
    }
    
    func paso_dado(dado, que, condicion) {
        return {
            tipo: "dado",
            descripcion: condicion
        }
    }
    
    func paso_cuando(cuando, accion) {
        return {
            tipo: "cuando",
            accion: accion
        }
    }
    
    func paso_entonces(entonces, expectativa) {
        return {
            tipo: "entonces",
            expectativa: expectativa
        }
    }
    
    func paso_y(y, condicion) {
        return {
            tipo: "y",
            descripcion: condicion
        }
    }
    
    func crear_condicion(texto) {
        return texto.substring(1, texto.length - 1)
    }
    
    func crear_accion(nombre, parametros) {
        return {
            nombre: nombre,
            parametros: parametros.substring(1, parametros.length - 1)
        }
    }
    
    func crear_expectativa(texto) {
        return texto.substring(1, texto.length - 1)
    }
}

func main() {
    console.log("=== DSL Testing ===")
    
    var parser = TestingDSL.use
    
    console.log("Escenario de prueba:")
    var escenario = parser("dado que \"usuario está logueado\" cuando login \"usuario123\" entonces \"dashboard visible\"")
    console.log("Escenario:", escenario)
    
    console.log("✅ DSL Testing funcionando!")
}