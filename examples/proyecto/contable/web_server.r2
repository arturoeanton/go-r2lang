// Sistema Contable LATAM - Servidor Web Funcional
// Demo R2Lang DSL para Siigo ERP con servidor HTTP real

console.log("========================================")
console.log("  SISTEMA CONTABLE LATAM - R2LANG DSL")
console.log("========================================")
console.log("  Demo Web para Siigo ERP Localization")
console.log("  Regiones: MX, COL, AR, CH, UY, EC, PE")
console.log("========================================")
console.log("")

// Importar servidor API funcional
import "src/api_server_working.r2" as apiserver

// Importar base de datos simple
import "database/database_simple.r2" as database

// Configuración de la aplicación
let appConfig = {
    name: "Sistema Contable LATAM",
    version: "1.0.0",
    description: "Demo R2Lang DSL para localización ERP",
    author: "R2Lang Team",
    target: "Siigo ERP",
    demo_mode: true,
    port: 8080,
    supported_regions: ["MX", "COL", "AR", "CH", "UY", "EC", "PE"]
}

// Función principal
func main() {
    console.log("🚀 Iniciando " + appConfig.name + " v" + appConfig.version)
    console.log("🎯 Objetivo: " + appConfig.target)
    console.log("")
    
    startWebApplication()
}

// Iniciar aplicación web completa
func startWebApplication() {
    console.log("🌐 INICIANDO APLICACIÓN WEB COMPLETA")
    console.log("===================================")
    console.log("")
    
    try {
        // Paso 1: Inicializar base de datos
        console.log("📊 Paso 1/3: Inicializando base de datos...")
        database.initDatabase()
        console.log("   ✅ Base de datos lista")
        console.log("")
        
        // Paso 2: Test DSL engines
        console.log("🔧 Paso 2/3: Verificando motores DSL...")
        apiserver.testAPIs()
        console.log("   ✅ Motores DSL operativos")
        console.log("")
        
        // Paso 3: Iniciar servidor web
        console.log("🚀 Paso 3/3: Iniciando servidor web...")
        console.log("")
        showWelcomeInfo()
        
        // Ejecutar servidor (esto bloquea el hilo principal)
        apiserver.startServer()
        
    } catch (error) {
        console.log("❌ ERROR iniciando aplicación: " + error)
        console.log("")
        runFallbackDemo()
    }
}

// Mostrar información de bienvenida
func showWelcomeInfo() {
    console.log("📋 INFORMACIÓN DE LA APLICACIÓN WEB")
    console.log("==================================")
    console.log("Aplicación: " + appConfig.name)
    console.log("Versión: " + appConfig.version)
    console.log("Puerto: " + appConfig.port)
    console.log("Modo: Demo para " + appConfig.target)
    console.log("")
    console.log("🌎 REGIONES LATAM SOPORTADAS:")
    let i = 0
    while (i < std.len(appConfig.supported_regions)) {
        let region = appConfig.supported_regions[i]
        console.log("   ✓ " + region + " - " + getRegionName(region))
        i = i + 1
    }
    console.log("")
    console.log("🔧 TECNOLOGÍAS:")
    console.log("   • R2Lang DSL Engine")
    console.log("   • r2http (Servidor Web nativo)")
    console.log("   • Almacenamiento en memoria")
    console.log("   • HTML5 + JavaScript")
    console.log("")
    console.log("💡 VALUE PROPOSITION PARA SIIGO:")
    console.log("   ✅ Reducción de 18 meses a 2 meses por país")
    console.log("   ✅ Ahorro de $500K a $150K por localización")
    console.log("   ✅ De 7 codebases separados a 1 DSL unificado")
    console.log("   ✅ ROI proyectado: 1,020% en 3 años")
    console.log("")
}

// Obtener nombre de región
func getRegionName(code) {
    let names = {
        "MX": "México",
        "COL": "Colombia", 
        "AR": "Argentina",
        "CH": "Chile",
        "UY": "Uruguay",
        "EC": "Ecuador",
        "PE": "Perú"
    }
    return names[code] || code
}

// Demo de fallback en caso de error
func runFallbackDemo() {
    console.log("🔧 MODO FALLBACK")
    console.log("================")
    console.log("Servidor web no pudo iniciarse.")
    console.log("Ejecutando demo básico...")
    console.log("")
    
    // Demo básico de DSL
    console.log("🎪 DEMO DSL BÁSICO")
    console.log("==================")
    console.log("• Sistema: " + appConfig.name)
    console.log("• Versión: " + appConfig.version)
    console.log("• Regiones: " + std.len(appConfig.supported_regions))
    console.log("• Target: " + appConfig.target)
    console.log("")
    
    // Test rápido
    try {
        apiserver.runQuickDemo()
    } catch (error) {
        console.log("⚠️ Error en demo básico: " + error)
    }
    
    console.log("")
    console.log("✅ Demo fallback completado")
    console.log("")
    console.log("💡 Sugerencias:")
    console.log("1. Verificar puerto " + appConfig.port + " libre")
    console.log("2. Reintentar: go run main.go examples/proyecto/contable/web_server.r2")
    console.log("3. Usar demo_completo.r2 como alternativa")
}

// Información del sistema
func showSystemInfo() {
    console.log("📖 INFORMACIÓN DETALLADA DEL SISTEMA")
    console.log("====================================")
    console.log("Nombre completo: " + appConfig.name)
    console.log("Versión: " + appConfig.version)
    console.log("Descripción: " + appConfig.description)
    console.log("Desarrollado por: " + appConfig.author)
    console.log("Cliente objetivo: " + appConfig.target)
    console.log("Modo demo: " + (appConfig.demo_mode ? "ACTIVO" : "INACTIVO"))
    console.log("Puerto HTTP: " + appConfig.port)
    console.log("")
    console.log("🎯 CARACTERÍSTICAS PRINCIPALES:")
    console.log("   • Procesamiento automático de transacciones")
    console.log("   • Localización multi-región LATAM")
    console.log("   • Cálculo automático de impuestos por país")
    console.log("   • Generación de asientos contables")
    console.log("   • API REST completa")
    console.log("   • Cumplimiento normativo por región")
    console.log("")
    console.log("🚀 READY FOR SIIGO DEMO!")
}

// Ejecutar aplicación
console.log("🎯 Preparando aplicación web...")
main()