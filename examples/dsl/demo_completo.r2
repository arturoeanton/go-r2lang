// Demostración completa de DSL en R2Lang
dsl DemoCompleto {
    // Tokens para acciones
    token("ENVIAR", "enviar")
    token("CREAR", "crear")
    token("ELIMINAR", "eliminar")
    
    // Tokens para objetos
    token("EMAIL", "email")
    token("ARCHIVO", "archivo")
    token("USUARIO", "usuario")
    
    // Tokens para valores
    token("NOMBRE", "[a-zA-Z]+")
    token("A", "a")
    
    // Reglas principales
    rule("comando", ["accion", "objeto", "A", "NOMBRE"], "ejecutar_comando")
    rule("accion", ["ENVIAR"], "accion_enviar")
    rule("accion", ["CREAR"], "accion_crear")
    rule("accion", ["ELIMINAR"], "accion_eliminar")
    rule("objeto", ["EMAIL"], "objeto_email")
    rule("objeto", ["ARCHIVO"], "objeto_archivo")
    rule("objeto", ["USUARIO"], "objeto_usuario")
    
    // Funciones semánticas
    func ejecutar_comando(accion, objeto, a, nombre) {
        return accion + " " + objeto + " para " + nombre
    }
    
    func accion_enviar(token) { return "Enviando" }
    func accion_crear(token) { return "Creando" }
    func accion_eliminar(token) { return "Eliminando" }
    
    func objeto_email(token) { return "email" }
    func objeto_archivo(token) { return "archivo" }
    func objeto_usuario(token) { return "usuario" }
}

func main() {
    console.log("=== Demo Completo DSL ===")
    console.log("Demuestra la capacidad de los DSL para procesar comandos complejos")
    console.log("")
    
    var procesador = DemoCompleto.use
    
    console.log("Comando 1: 'enviar email a juan'")
    var resultado1 = procesador("enviar email a juan")
    console.log("Resultado:", resultado1)
    console.log("")
    
    console.log("Comando 2: 'crear archivo a maria'")
    var resultado2 = procesador("crear archivo a maria")
    console.log("Resultado:", resultado2)
    console.log("")
    
    console.log("Comando 3: 'eliminar usuario a pedro'")
    var resultado3 = procesador("eliminar usuario a pedro")
    console.log("Resultado:", resultado3)
    console.log("")
    
    console.log("✅ Demo completo DSL funcionando!")
    console.log("El DSL ha procesado todos los comandos usando gramática personalizada")
}