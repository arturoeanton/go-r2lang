// Ejemplo DSL completamente funcional y claro
dsl EjemploFuncional {
    // Definir tokens básicos
    token("CREAR", "crear")
    token("USUARIO", "usuario")
    token("NOMBRE", "[a-zA-Z]+")
    
    // Definir regla simple
    rule("comando", ["CREAR", "USUARIO", "NOMBRE"], "ejecutar_creacion")
    
    // Función que procesa el comando
    func ejecutar_creacion(crear, usuario, nombre) {
        return "Usuario creado: " + nombre
    }
}

func main() {
    console.log("=== DSL Funcional ===")
    console.log("Ejemplo claro de DSL procesando comandos")
    console.log("")
    
    var procesador = EjemploFuncional.use
    
    console.log("Comando DSL: 'crear usuario juan'")
    var resultado = procesador("crear usuario juan")
    console.log("Resultado:", resultado)
    console.log("")
    
    console.log("Comando DSL: 'crear usuario maria'")
    var resultado2 = procesador("crear usuario maria")
    console.log("Resultado:", resultado2)
    console.log("")
    
    console.log("✅ DSL funcional procesando comandos correctamente!")
    console.log("El DSL ha interpretado los comandos y ejecutado las acciones")
}