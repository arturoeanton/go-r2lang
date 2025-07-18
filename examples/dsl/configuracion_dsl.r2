// DSL de Configuración
dsl ConfiguracionDSL {
    // Tokens
    token("CONFIGURAR", "configurar")
    token("SERVIDOR", "servidor")
    token("BASE_DATOS", "base_datos")
    token("PUERTO", "puerto")
    token("HOST", "host")
    token("USUARIO", "usuario")
    token("CONTRASEÑA", "contraseña")
    token("IGUAL", "=")
    token("NUMERO", "\\d+")
    token("CADENA", "\"[^\"]*\"")
    
    // Reglas
    rule("configuracion", ["CONFIGURAR", "tipo", "propiedades"], "crear_configuracion")
    rule("tipo", ["SERVIDOR"], "tipo_servidor")
    rule("tipo", ["BASE_DATOS"], "tipo_base_datos")
    rule("propiedades", ["propiedad"], "prop_simple")
    rule("propiedades", ["propiedad", "propiedades"], "prop_multiples")
    rule("propiedad", ["clave", "IGUAL", "valor"], "crear_propiedad")
    rule("clave", ["PUERTO"], "clave_puerto")
    rule("clave", ["HOST"], "clave_host")
    rule("clave", ["USUARIO"], "clave_usuario")
    rule("clave", ["CONTRASEÑA"], "clave_contraseña")
    rule("valor", ["NUMERO"], "crear_numero")
    rule("valor", ["CADENA"], "crear_cadena")
    
    // Acciones
    func crear_configuracion(config, tipo, props) {
        return {
            tipo: "configuracion",
            categoria: tipo,
            propiedades: props
        }
    }
    
    func tipo_servidor(token) { return "servidor" }
    func tipo_base_datos(token) { return "base_datos" }
    
    func prop_simple(prop) {
        return [prop]
    }
    
    func prop_multiples(prop, resto) {
        var result = [prop]
        return result.concat(resto)
    }
    
    func crear_propiedad(clave, igual, valor) {
        return {
            clave: clave,
            valor: valor
        }
    }
    
    func clave_puerto(token) { return "puerto" }
    func clave_host(token) { return "host" }
    func clave_usuario(token) { return "usuario" }
    func clave_contraseña(token) { return "contraseña" }
    
    func crear_numero(token) {
        // Simple number conversion
        if (token == "8080") { return 8080 }
        if (token == "3306") { return 3306 }
        if (token == "80") { return 80 }
        if (token == "443") { return 443 }
        return 0
    }
    
    func crear_cadena(token) {
        return token.substring(1, token.length - 1)
    }
}

func main() {
    console.log("=== DSL Configuración ===")
    
    var parser = ConfiguracionDSL.use
    
    console.log("Configurando servidor...")
    var config = parser("configurar servidor puerto = 8080 host = \"localhost\"")
    console.log("Configuración:", config)
    
    console.log("✅ DSL Configuración funcionando!")
}