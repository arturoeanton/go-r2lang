// DSL básico para demostrar funcionalidad
dsl Basico {
    token("PALABRA", "[a-zA-Z]+")
    
    rule("texto", ["PALABRA"], "procesar_palabra")
    
    func procesar_palabra(palabra) {
        return "Procesé: " + palabra
    }
}

func main() {
    console.log("=== DSL Básico ===")
    console.log("Demuestra cómo un DSL puede procesar texto con sintaxis personalizada")
    
    var dslBasico = Basico.use
    
    console.log("Entrada: 'hola'")
    var resultado1 = dslBasico("hola")
    console.log("Salida DSL:", resultado1)
    
    console.log("Entrada: 'mundo'")
    var resultado2 = dslBasico("mundo")
    console.log("Salida DSL:", resultado2)
    
    console.log("✅ DSL Básico funciona correctamente!")
    console.log("El DSL ha procesado las palabras usando su gramática personalizada")
}