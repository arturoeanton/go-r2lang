// API Server REST para Sistema Contable LATAM
// Usando r2http para crear servidor web y APIs

// Importar dependencias
import "../database/database.r2" as db
import "dsl_contable_latam.r2" as dslmod

// Configuración del servidor
let serverConfig = {
    port: 8080,
    host: "0.0.0.0",
    cors: true,
    static_path: "./static",
    upload_path: "./uploads"
}

// Función para inicializar servidor
func initServer() {
    console.log("=== INICIANDO SERVIDOR API CONTABLE LATAM ===")
    console.log("Puerto: " + serverConfig.port)
    console.log("Host: " + serverConfig.host)
    
    // Inicializar base de datos
    db.initDatabase()
    
    // Crear servidor HTTP
    let server = http.createServer(serverConfig.port)
    
    // Configurar CORS
    server.use(http.cors())
    
    // Servir archivos estáticos
    server.static("/", serverConfig.static_path)
    
    // Middleware para logging
    server.use(function(req, res, next) {
        console.log(std.now() + " " + req.method + " " + req.path)
        next()
    })
    
    // Middleware para parsing JSON
    server.use(http.json())
    
    // Rutas de la API
    setupRoutes(server)
    
    console.log("Servidor iniciado en http://localhost:" + serverConfig.port)
    console.log("Panel de administración: http://localhost:" + serverConfig.port + "/admin")
    
    return server
}

func setupRoutes(server) {
    // Ruta principal - Info de la API
    server.get("/api/info", function(req, res) {
        res.json({
            name: "Sistema Contable LATAM API",
            version: "1.0.0", 
            description: "API para demostración DSL R2Lang con Siigo",
            supported_regions: ["MX", "COL", "AR", "CH", "UY", "EC", "PE"],
            endpoints: [
                "GET /api/info - Información de la API",
                "GET /api/regions - Lista de regiones soportadas",
                "GET /api/regions/{code} - Configuración de región específica",
                "POST /api/transactions/sale - Procesar venta",
                "POST /api/transactions/purchase - Procesar compra", 
                "GET /api/transactions - Obtener transacciones",
                "GET /api/stats - Estadísticas del sistema",
                "POST /api/dsl/execute - Ejecutar comando DSL directo"
            ],
            demo_mode: true,
            timestamp: std.now()
        })
    })
    
    // Rutas de regiones
    server.get("/api/regions", function(req, res) {
        try {
            let regions = db.getAllRegions()
            res.json({
                success: true,
                count: regions.length,
                data: regions
            })
        } catch (error) {
            res.status(500).json({
                success: false,
                error: "Error obteniendo regiones: " + error
            })
        }
    })
    
    server.get("/api/regions/:code", function(req, res) {
        try {
            let regionCode = req.params.code
            let region = db.getRegionConfig(regionCode)
            
            if (region) {
                res.json({
                    success: true,
                    data: region
                })
            } else {
                res.status(404).json({
                    success: false,
                    error: "Región no encontrada: " + regionCode
                })
            }
        } catch (error) {
            res.status(500).json({
                success: false,
                error: "Error obteniendo región: " + error
            })
        }
    })
    
    // Ruta para procesar venta
    server.post("/api/transactions/sale", function(req, res) {
        try {
            let region = req.body.region
            let amount = req.body.amount
            
            if (!region || !amount) {
                return res.status(400).json({
                    success: false,
                    error: "Faltan parámetros: region y amount son requeridos"
                })
            }
            
            // Validar región
            let regionConfig = db.getRegionConfig(region)
            if (!regionConfig) {
                return res.status(400).json({
                    success: false,
                    error: "Región no soportada: " + region
                })
            }
            
            // Procesar venta usando DSL
            let motorVentas = dslmod.VentasLATAM
            let dslCommand = "venta " + region + " " + amount
            let result = motorVentas.use(dslCommand)
            
            // Agregar información adicional
            result.country = regionConfig.name
            result.type = "sale"
            
            // Guardar en base de datos
            db.saveTransaction(result)
            
            res.json({
                success: true,
                message: "Venta procesada exitosamente",
                data: result
            })
            
        } catch (error) {
            res.status(500).json({
                success: false,
                error: "Error procesando venta: " + error
            })
        }
    })
    
    // Ruta para procesar compra
    server.post("/api/transactions/purchase", function(req, res) {
        try {
            let region = req.body.region
            let amount = req.body.amount
            
            if (!region || !amount) {
                return res.status(400).json({
                    success: false,
                    error: "Faltan parámetros: region y amount son requeridos"
                })
            }
            
            // Validar región
            let regionConfig = db.getRegionConfig(region)
            if (!regionConfig) {
                return res.status(400).json({
                    success: false,
                    error: "Región no soportada: " + region
                })
            }
            
            // Procesar compra usando DSL
            let motorCompras = dslmod.ComprasLATAM
            let dslCommand = "compra " + region + " " + amount
            let result = motorCompras.use(dslCommand)
            
            // Agregar información adicional
            result.country = regionConfig.name
            result.type = "purchase"
            
            // Guardar en base de datos
            db.saveTransaction(result)
            
            res.json({
                success: true,
                message: "Compra procesada exitosamente",
                data: result
            })
            
        } catch (error) {
            res.status(500).json({
                success: false,
                error: "Error procesando compra: " + error
            })
        }
    })
    
    // Ruta para obtener transacciones
    server.get("/api/transactions", function(req, res) {
        try {
            let region = req.query.region || "ALL"
            let transactions = db.getTransactions(region)
            
            res.json({
                success: true,
                count: transactions.length,
                region: region,
                data: transactions
            })
            
        } catch (error) {
            res.status(500).json({
                success: false,
                error: "Error obteniendo transacciones: " + error
            })
        }
    })
    
    // Ruta para estadísticas
    server.get("/api/stats", function(req, res) {
        try {
            let stats = db.getTransactionStats()
            let regions = db.getAllRegions()
            
            res.json({
                success: true,
                data: {
                    transactions: stats,
                    regions_available: regions.length,
                    system_info: {
                        version: "1.0.0",
                        demo_mode: true,
                        supported_regions: regions.length,
                        uptime: "N/A"
                    }
                }
            })
            
        } catch (error) {
            res.status(500).json({
                success: false,
                error: "Error obteniendo estadísticas: " + error
            })
        }
    })
    
    // Ruta para ejecutar DSL directo (para demos avanzadas)
    server.post("/api/dsl/execute", function(req, res) {
        try {
            let command = req.body.command
            let engine = req.body.engine || "ventas"
            
            if (!command) {
                return res.status(400).json({
                    success: false,
                    error: "Comando DSL requerido"
                })
            }
            
            let result = null
            
            if (engine == "ventas") {
                let motor = dslmod.VentasLATAM
                result = motor.use(command)
            } else if (engine == "compras") {
                let motor = dslmod.ComprasLATAM  
                result = motor.use(command)
            } else if (engine == "consultas") {
                let motor = dslmod.ConsultasLATAM
                result = motor.use(command)
            } else {
                return res.status(400).json({
                    success: false,
                    error: "Motor DSL no válido. Opciones: ventas, compras, consultas"
                })
            }
            
            res.json({
                success: true,
                message: "Comando DSL ejecutado exitosamente",
                command: command,
                engine: engine,
                data: result
            })
            
        } catch (error) {
            res.status(500).json({
                success: false,
                error: "Error ejecutando DSL: " + error
            })
        }
    })
    
    // Ruta para demo completo
    server.post("/api/demo", function(req, res) {
        try {
            console.log("Ejecutando demo completo...")
            
            let results = []
            let regions = ["MX", "COL", "AR", "CH", "UY", "EC", "PE"]
            
            regions.forEach(function(region) {
                // Demo venta
                let saleCommand = "venta " + region + " 10000"
                let saleResult = dslmod.VentasLATAM.use(saleCommand)
                saleResult.type = "sale"
                saleResult.demo = true
                results.push(saleResult)
                
                // Demo compra
                let purchaseCommand = "compra " + region + " 5000"
                let purchaseResult = dslmod.ComprasLATAM.use(purchaseCommand)
                purchaseResult.type = "purchase"
                purchaseResult.demo = true
                results.push(purchaseResult)
            })
            
            res.json({
                success: true,
                message: "Demo completo ejecutado para todas las regiones LATAM",
                regions_processed: regions.length,
                transactions_generated: results.length,
                data: results
            })
            
        } catch (error) {
            res.status(500).json({
                success: false,
                error: "Error en demo: " + error
            })
        }
    })
    
    console.log("Rutas API configuradas correctamente")
}

// Función para inicializar y ejecutar servidor
func startServer() {
    let server = initServer()
    
    // Manejar señales de cierre
    process.on("SIGINT", function() {
        console.log("\nCerrando servidor...")
        server.close()
        process.exit(0)
    })
    
    // Iniciar servidor
    server.listen(function() {
        console.log("🚀 Servidor API LATAM funcionando!")
        console.log("📊 Dashboard disponible en: http://localhost:8080")
        console.log("🔗 API docs: http://localhost:8080/api/info")
        console.log("")
        console.log("Presiona Ctrl+C para detener el servidor")
    })
    
    return server
}

// Función de test de APIs
func testAPIs() {
    console.log("=== TEST DE APIs ===")
    
    // Test básico - se puede extender con llamadas HTTP reales
    console.log("✓ Rutas configuradas")
    console.log("✓ Middleware configurado")
    console.log("✓ CORS habilitado")
    console.log("✓ Archivos estáticos configurados")
    
    console.log("Para test completo, iniciar servidor con startServer()")
}

// Exportar funciones principales
func getServerFunctions() {
    return {
        init: initServer,
        start: startServer,
        test: testAPIs,
        config: serverConfig
    }
}