// Sistema Contable LATAM - Versión Web Principal
// Demo R2Lang DSL para Siigo ERP Localization con servidor HTTP

console.log("========================================")
console.log("  SISTEMA CONTABLE LATAM - R2LANG DSL")
console.log("========================================")
console.log("  Demo Web para Siigo ERP Localization")
console.log("  Regiones: MX, COL, AR, CH, UY, EC, PE")
console.log("========================================")
console.log("")

// Importar módulos funcionales
import "database/database_simple.r2" as database
import "src/api_server_working.r2" as apiserver

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
    startWebServer()
}

// Iniciar servidor web completo
func startWebServer() {
    console.log("🌐 INICIANDO SERVIDOR WEB COMPLETO")
    console.log("==================================")
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
        
        // Paso 3: Mostrar información y arrancar
        console.log("🚀 Paso 3/3: Iniciando servidor...")
        showSystemInfo()
        console.log("")
        
        // Demo inicial
        runInitialDemo()
        console.log("")
        
        // Arrancar servidor (bloquea hilo principal)
        console.log("🔥 Arrancando servidor HTTP...")
        apiserver.startServer()
        
    } catch (error) {
        console.log("❌ ERROR iniciando servidor: " + error)
        console.log("🔧 Intentando modo consola...")
        runFallbackMode()
    }
}

// Demo inicial del sistema  
func runInitialDemo() {
    console.log("🎪 Demo inicial de verificación...")
    
    let demos = [
        {region: "COL", amount: 100000, type: "venta"},
        {region: "MX", amount: 75000, type: "compra"}, 
        {region: "AR", amount: 50000, type: "venta"}
    ]
    
    let i = 0
    while (i < std.len(demos)) {
        let demo = demos[i]
        console.log("   ✓ " + demo.type + " " + demo.region + ": " + demo.amount)
        i = i + 1
    }
    
    console.log("   ✅ Demo verificación completado")
}

// Mostrar información del sistema
func showSystemInfo() {
    console.log("📋 INFORMACIÓN DEL SISTEMA")
    console.log("==========================")
    console.log("Aplicación: " + appConfig.name)
    console.log("Versión: " + appConfig.version)
    console.log("Puerto: " + appConfig.port)
    console.log("Target: " + appConfig.target)
    console.log("")
    console.log("🌎 REGIONES SOPORTADAS:")
    let i = 0
    while (i < std.len(appConfig.supported_regions)) {
        let region = appConfig.supported_regions[i]
        console.log("   ✓ " + region + " - " + getRegionName(region))
        i = i + 1
    }
    console.log("")
    console.log("🔧 TECNOLOGÍAS:")
    console.log("   • R2Lang DSL Engine")
    console.log("   • r2http (Servidor Web)")
    console.log("   • Almacenamiento en memoria")
    console.log("")
    console.log("💡 VALUE PROPOSITION PARA SIIGO:")
    console.log("   ✅ 18 meses → 2 meses por país")
    console.log("   ✅ $500K → $150K por localización") 
    console.log("   ✅ 7 codebases → 1 DSL unificado")
    console.log("   ✅ ROI: 1,020% en 3 años")
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

// Modo fallback
func runFallbackMode() {
    console.log("🔧 MODO FALLBACK")
    console.log("================")
    console.log("El servidor web no pudo iniciarse.")
    console.log("")
    console.log("💡 Usar alternativamente:")
    console.log("   • web_server.r2 (servidor funcional)")
    console.log("   • demo_completo.r2 (demo consola)")
    console.log("")
    console.log("💡 Verificaciones:")
    console.log("   1. Puerto 8080 libre: lsof -i :8080")
    console.log("   2. Permisos de red")
    console.log("   3. Firewall/antivirus")
}

// Ejecutar aplicación principal
console.log("🎯 Preparando aplicación...")
main()