// DSL para comandos simples
dsl ComandoDSL {
    token("HACER", "hacer")
    token("CAFE", "cafe")
    token("TE", "te")
    
    rule("comando", ["HACER", "bebida"], "ejecutar_comando")
    rule("bebida", ["CAFE"], "elegir_cafe")
    rule("bebida", ["TE"], "elegir_te")
    
    func ejecutar_comando( hacer,bebida) {
        return "Ejecutando: " + hacer + " " + bebida
    }
    
    func elegir_cafe(cafe) {
        return "café"
    }
    
    func elegir_te(te) {
        return "té"
    }
}

func main() {
    console.log("=== DSL Comando Simple ===")
    
    var cmd = ComandoDSL.use
    
    console.log("Comando: 'hacer cafe'")
    var resultado1 = cmd("hacer cafe")
    console.log("Resultado:", resultado1)
    
    console.log("Comando: 'hacer te'")
    var resultado2 = cmd("hacer te")
    console.log("Resultado:", resultado2)
    
    console.log("✅ DSL Comando funcionando!")
}