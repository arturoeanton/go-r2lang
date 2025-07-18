// Ejemplo del Inicio Rápido
dsl MiPrimerDSL {
    // Definir tokens (palabras clave)
    token("SALUDO", "hola")
    token("NOMBRE", "[a-zA-Z]+")
    
    // Definir reglas (sintaxis)
    rule("mensaje", ["SALUDO", "NOMBRE"], "crear_saludo")
    
    // Definir acción (qué hacer)
    func crear_saludo(saludo, nombre) {
        return "¡" + saludo + " " + nombre + "!"
    }
}

func main() {
    console.log("=== Mi Primer DSL ===")
    console.log("Este DSL procesa saludos con sintaxis personalizada")
    
    // Usar el DSL
    var metodo = MiPrimerDSL.use
    
    console.log("Procesando: 'hola mundo'")
    var resultado = metodo("hola mundo")
    console.log("Resultado DSL:", resultado)
    console.log("El DSL ha convertido 'hola mundo' en un saludo")
    
    console.log("✅ Primer DSL funcionando!")
}