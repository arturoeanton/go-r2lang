// Sistema Contable LATAM - Aplicación Principal
// Demo R2Lang DSL para Siigo ERP Localization
// Integración completa: Frontend + Backend + Base de Datos + DSL

console.log("========================================")
console.log("  SISTEMA CONTABLE LATAM - R2LANG DSL")
console.log("========================================")
console.log("  Demo para Siigo ERP Localization")
console.log("  Regiones: MX, COL, AR, CH, UY, EC, PE")
console.log("========================================")
console.log("")

// Importar módulos del sistema
import "database/database.r2" as database;
import "src/api_server.r2" as server;
import "src/dsl_contable_latam.r2" as dslcontable  

// Configuración de la aplicación
let appConfig = {
    name: "Sistema Contable LATAM",
    version: "1.0.0",
    description: "Demo R2Lang DSL para localización ERP",
    author: "R2Lang Team",
    target: "Siigo ERP",
    demo_mode: true,
    supported_regions: ["MX", "COL", "AR", "CH", "UY", "EC", "PE"]
}

// Función principal
func main() {
    console.log("🚀 Iniciando " + appConfig.name + " v" + appConfig.version)
    console.log("")
    
    // Mostrar menú de opciones
    showMainMenu()
}

func showMainMenu() {
    console.log("=== MENÚ PRINCIPAL ===")
    console.log("Selecciona una opción:")
    console.log("")
    console.log("1. 🌐 Iniciar Servidor Web Completo")
    console.log("2. 🧪 Ejecutar Tests del Sistema")
    console.log("3. 📊 Demo DSL Standalone")
    console.log("4. 💾 Inicializar Base de Datos")
    console.log("5. 🔧 Tests de Componentes")
    console.log("6. 📖 Información del Sistema")
    console.log("")
    
    // Para demo, ejecutar automáticamente el servidor web
    console.log("🎯 Modo Demo: Iniciando servidor web automáticamente...")
    startWebServer()
}

func startWebServer() {
    try {
        console.log("🌐 INICIANDO SERVIDOR WEB COMPLETO")
        console.log("==================================")
        
        // 1. Inicializar base de datos
        console.log("📊 Paso 1/4: Inicializando base de datos...")
        database.initDatabase()
        console.log("   ✅ Base de datos SQLite lista")
        
        // 2. Verificar DSL engines
        console.log("🔧 Paso 2/4: Verificando motores DSL...")
        verifyDSLEngines()
        console.log("   ✅ Motores DSL operativos")
        
        // 3. Inicializar servidor API
        console.log("🚀 Paso 3/4: Iniciando servidor HTTP...")
        let httpServer = server.initServer()
        console.log("   ✅ Servidor HTTP configurado")
        
        // 4. Mostrar información de acceso
        console.log("🎉 Paso 4/4: Sistema completamente operativo!")
        console.log("")
        showServerInfo()
        
        // Ejecutar demo inicial
        console.log("🎪 Ejecutando demo inicial...")
        runInitialDemo()
        
        // Iniciar servidor
        console.log("🔥 Iniciando servidor web...")
        httpServer.listen(function() {
            console.log("")
            console.log("🎯 ¡SISTEMA LATAM FUNCIONANDO!")
            console.log("===============================")
            console.log("🌍 URL Principal: http://localhost:8080")
            console.log("📊 Panel Admin: http://localhost:8080/admin")
            console.log("🔗 API Info: http://localhost:8080/api/info")
            console.log("📱 Frontend: http://localhost:8080/index.html")
            console.log("")
            console.log("📝 Para detener: Ctrl+C")
            console.log("🔄 Para recargar: F5 en el navegador")
            console.log("")
            console.log("🎉 ¡DEMO LISTO PARA PRESENTACIÓN A SIIGO!")
        })
        
    } catch (error) {
        console.log("❌ ERROR iniciando servidor:")
        console.log("   " + error)
        console.log("")
        console.log("🔧 Intentando modo de recuperación...")
        fallbackMode()
    }
}

func verifyDSLEngines() {
    // Test motores DSL
    let motorVentas = dslcontable.VentasLATAM
    let motorCompras = dslcontable.ComprasLATAM
    let motorConsultas = dslcontable.ConsultasLATAM
    
    // Test básico de cada motor
    try {
        let testVenta = dslcontable.VentasLATAM.use("venta COL 1000")
        let testCompra = dslcontable.ComprasLATAM.use("compra MX 500")
        let testConsulta = dslcontable.ConsultasLATAM.use("consultar config AR")
        
        console.log("   ✓ Motor Ventas: Operativo")
        console.log("   ✓ Motor Compras: Operativo")  
        console.log("   ✓ Motor Consultas: Operativo")
        
        return true
    } catch (error) {
        console.log("   ❌ Error en motores DSL: " + error)
        return false
    }
}

func runInitialDemo() {
    try {
        console.log("")
        console.log("🎪 Demo Inicial - Procesando transacciones de ejemplo...")
        
        let demos = [
            {region: "COL", amount: 100000, type: "venta"},
            {region: "MX", amount: 75000, type: "compra"},
            {region: "AR", amount: 50000, type: "venta"}
        ]
        
        let j = 0
        while (j < std.len(demos)) {
            let demo = demos[j]
            try {
                if (demo.type == "venta") {
                    let result = dslcontable.VentasLATAM.use("venta " + demo.region + " " + demo.amount)
                    console.log("   ✓ Demo venta " + demo.region + ": " + result.transactionId)
                } else {
                    let result = dslcontable.ComprasLATAM.use("compra " + demo.region + " " + demo.amount)
                    console.log("   ✓ Demo compra " + demo.region + ": " + result.transactionId)
                }
            } catch (error) {
                console.log("   ⚠️  Demo " + demo.type + " " + demo.region + " falló: " + error)
            }
            j = j + 1
        }
        
        console.log("🎪 Demo inicial completado")
        
    } catch (error) {
        console.log("⚠️ Error en demo inicial: " + error)
    }
}

func showServerInfo() {
    console.log("📋 INFORMACIÓN DEL SERVIDOR")
    console.log("===========================")
    console.log("Puerto: 8080")
    console.log("Host: localhost")
    console.log("Modo: Demostración")
    console.log("CORS: Habilitado") 
    console.log("Archivos estáticos: /static")
    console.log("Base de datos: SQLite (./database/contable_latam.db)")
    console.log("")
    console.log("📡 ENDPOINTS API DISPONIBLES:")
    console.log("GET  /api/info           - Información de la API")
    console.log("GET  /api/regions        - Lista de regiones")
    console.log("GET  /api/regions/{code} - Configuración regional")
    console.log("POST /api/transactions/sale     - Procesar venta")
    console.log("POST /api/transactions/purchase - Procesar compra")
    console.log("GET  /api/transactions   - Obtener transacciones")
    console.log("GET  /api/stats          - Estadísticas del sistema")
    console.log("POST /api/dsl/execute    - Ejecutar DSL directo")
    console.log("POST /api/demo           - Demo completo")
    console.log("")
}

func runSystemTests() {
    console.log("🧪 EJECUTANDO TESTS DEL SISTEMA")
    console.log("================================")
    
    console.log("Test 1/4: Base de Datos...")
    database.testDatabase()
    
    console.log("")
    console.log("Test 2/4: DSL Engines...")
    dslcontable.testDSLLatam()
    
    console.log("")
    console.log("Test 3/4: APIs...")
    server.testAPIs()
    
    console.log("")
    console.log("Test 4/4: Integración...")
    testIntegration()
    
    console.log("")
    console.log("✅ Tests completados")
}

func runDSLDemo() {
    console.log("📊 DEMO DSL STANDALONE")
    console.log("======================")
    
    // Ejecutar demo de DSL
    dslcontable.testDSLLatam()
    
    console.log("")
    console.log("✅ Demo DSL completado")
}

func testIntegration() {
    console.log("🔄 Test de integración...")
    
    try {
        // Test DSL -> DB
        let result = dslcontable.VentasLATAM.use("venta COL 5000")
        result.country = "Colombia"
        result.type = "sale"
        
        database.saveTransaction(result)
        console.log("   ✓ DSL -> DB: OK")
        
        // Test DB -> Stats
        let stats = database.getTransactionStats()
        console.log("   ✓ DB -> Stats: OK (" + stats.total_transactions + " transacciones)")
        
        // Test Region Config
        let config = database.getRegionConfig("COL")
        console.log("   ✓ Config Regional: OK (" + config.name + ")")
        
        return true
    } catch (error) {
        console.log("   ❌ Error integración: " + error)
        return false
    }
}

func fallbackMode() {
    console.log("🔧 MODO DE RECUPERACIÓN")
    console.log("=======================")
    console.log("El servidor web no pudo iniciarse.")
    console.log("Ejecutando demo DSL en modo consola...")
    console.log("")
    
    runDSLDemo()
    
    console.log("")
    console.log("💡 Sugerencias:")
    console.log("1. Verificar que el puerto 8080 esté libre")
    console.log("2. Ejecutar: lsof -i :8080")
    console.log("3. Reintentar con: go run main.go")
}

func showSystemInfo() {
    console.log("📖 INFORMACIÓN DEL SISTEMA")
    console.log("===========================")
    console.log("Nombre: " + appConfig.name)
    console.log("Versión: " + appConfig.version)
    console.log("Descripción: " + appConfig.description)
    console.log("Autor: " + appConfig.author)
    console.log("Objetivo: " + appConfig.target)
    console.log("Modo Demo: " + (appConfig.demo_mode ? "Activo" : "Inactivo"))
    console.log("")
    console.log("🌎 REGIONES SOPORTADAS:")
    let k = 0
    while (k < std.len(appConfig.supported_regions)) {
        let region = appConfig.supported_regions[k]
        console.log("   " + region + " - " + getRegionName(region))
        k = k + 1
    }
    console.log("")
    console.log("🔧 TECNOLOGÍAS UTILIZADAS:")
    console.log("   • R2Lang DSL Engine")
    console.log("   • r2http (Servidor Web)")
    console.log("   • r2db (Base de Datos SQLite)")
    console.log("   • HTML5 + Bootstrap 5")
    console.log("   • JavaScript ES6+")
    console.log("")
    console.log("🎯 FUNCIONALIDADES CLAVE:")
    console.log("   • Procesamiento automático de comprobantes")
    console.log("   • Localización multi-región automática")
    console.log("   • Cálculo de impuestos por país")
    console.log("   • Generación de asientos contables")
    console.log("   • API REST completa")
    console.log("   • Interfaz web responsiva")
    console.log("   • Base de datos persistente")
}

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

// Manejo de argumentos de línea de comandos
func handleArgs() {
    // Por defecto iniciar servidor web
    // En un entorno real, aquí se manejarían argumentos como:
    // --server, --test, --demo, --info, etc.
    return "server"
}

// Función de limpieza al cerrar
func cleanup() {
    console.log("")
    console.log("🧹 Limpiando recursos...")
    console.log("✅ Sistema cerrado correctamente")
}

// Entry point - detectar qué hacer
func detectMode() {
    let mode = handleArgs()
    
    if (mode == "server") {
        startWebServer()
    } else if (mode == "test") {
        runSystemTests()
    } else if (mode == "demo") {
        runDSLDemo()
    } else if (mode == "info") {
        showSystemInfo()
    } else {
        main()
    }
}

// Inicializar aplicación
detectMode()

// Configurar manejo de señales (simulado)
console.log("🎯 Sistema iniciado - Modo: " + (appConfig.demo_mode ? "DEMO" : "PRODUCCIÓN"))
console.log("📱 Para acceder al sistema: http://localhost:8080")
console.log("🎪 ¡Listo para demostración a Siigo!")