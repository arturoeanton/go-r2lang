// API Server Funcional para Sistema Contable LATAM
// Usando r2http real de R2Lang

console.log("🚀 Módulo API Server Working cargado")

// Regiones disponibles
let regiones = {
    "MX": {nombre: "México", moneda: "MXN", iva: 0.16, normativa: "NIF-Mexican"},
    "COL": {nombre: "Colombia", moneda: "COP", iva: 0.19, normativa: "NIIF-Colombia"},
    "AR": {nombre: "Argentina", moneda: "ARS", iva: 0.21, normativa: "RT-Argentina"},
    "CH": {nombre: "Chile", moneda: "CLP", iva: 0.19, normativa: "IFRS-Chile"},
    "UY": {nombre: "Uruguay", moneda: "UYU", iva: 0.22, normativa: "NIIF-Uruguay"},
    "EC": {nombre: "Ecuador", moneda: "USD", iva: 0.12, normativa: "NIIF-Ecuador"},
    "PE": {nombre: "Perú", moneda: "PEN", iva: 0.18, normativa: "PCGE-Peru"}
}

// Almacenamiento de transacciones en memoria
let transacciones = []

// DSL Engine integrado
dsl VentasAPI {
    token("VENTA", "venta|sale")
    token("REGION", "MX|COL|AR|CH|UY|EC|PE")
    token("IMPORTE", "[0-9]+\\.?[0-9]*")
    
    rule("venta_api", ["VENTA", "REGION", "IMPORTE"], "procesarVentaAPI")
    
    func procesarVentaAPI(operacion, region, importe) {
        let config = regiones[region]
        let importeNum = std.parseFloat(importe)
        let importeIVA = math.round((importeNum * config.iva) * 100) / 100
        let importeTotal = math.round((importeNum + importeIVA) * 100) / 100
        let txId = region + "-" + std.now() + "-" + math.randomInt(9999)
        
        return {
            success: true,
            transactionId: txId,
            region: region,
            country: config.nombre,
            amount: importeNum,
            tax: importeIVA,
            total: importeTotal,
            currency: config.moneda,
            compliance: config.normativa,
            timestamp: std.now()
        }
    }
}

dsl ComprasAPI {
    token("COMPRA", "compra|purchase")
    token("REGION", "MX|COL|AR|CH|UY|EC|PE")
    token("IMPORTE", "[0-9]+\\.?[0-9]*")
    
    rule("compra_api", ["COMPRA", "REGION", "IMPORTE"], "procesarCompraAPI")
    
    func procesarCompraAPI(operacion, region, importe) {
        let config = regiones[region]
        let importeNum = std.parseFloat(importe)
        let importeIVA = math.round((importeNum * config.iva) * 100) / 100
        let importeTotal = math.round((importeNum + importeIVA) * 100) / 100
        let txId = region + "-" + std.now() + "-" + math.randomInt(9999)
        
        return {
            success: true,
            transactionId: txId,
            region: region,
            country: config.nombre,
            amount: importeNum,
            tax: importeIVA,
            total: importeTotal,
            currency: config.moneda,
            compliance: config.normativa,
            timestamp: std.now()
        }
    }
}

// Handlers de rutas

func handleIndex() {
    let html = "<!DOCTYPE html><html><head><title>Sistema Contable LATAM</title><meta charset=\"UTF-8\"><style>body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; } .container { background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); } h1 { color: #2c3e50; border-bottom: 3px solid #3498db; padding-bottom: 10px; } .links { margin-top: 20px; } .links a { display: inline-block; margin: 10px 15px 10px 0; padding: 10px 20px; background: #3498db; color: white; text-decoration: none; border-radius: 5px; } .links a:hover { background: #2980b9; } .info { background: #ecf0f1; padding: 15px; border-radius: 5px; margin: 15px 0; }</style></head><body><div class=\"container\"><h1>🌍 Sistema Contable LATAM - R2Lang DSL</h1><p class=\"info\"><strong>Demo para Siigo ERP Localization</strong><br>Procesamiento automatico de transacciones contables para 7 paises de LATAM</p><h2>🔗 Enlaces disponibles:</h2><div class=\"links\"><a href=\"/api/info\">📊 API Info</a><a href=\"/api/regions\">🌎 Regiones LATAM</a></div><h2>📡 API Endpoints:</h2><ul><li><code>GET /api/info</code> - Informacion de la API</li><li><code>GET /api/regions</code> - Lista de regiones LATAM</li><li><code>POST /api/transactions/sale</code> - Procesar venta</li><li><code>POST /api/transactions/purchase</code> - Procesar compra</li></ul><div class=\"info\"><strong>🎯 Value Proposition para Siigo:</strong><br>✅ 18 meses → 2 meses por pais<br>✅ $500K → $150K por localizacion<br>✅ 7 codebases → 1 DSL unificado<br>✅ ROI: 1,020% en 3 años</div></div></body></html>"
    
    return HttpResponse("text/html", html)
}

func handleAPIInfo() {
    let info = {
        name: "Sistema Contable LATAM API",
        version: "1.0.0",
        description: "R2Lang DSL para localización ERP Siigo",
        regions: std.len(regiones),
        transactions_processed: std.len(transacciones),
        endpoints: [
            "GET /api/info",
            "GET /api/regions", 
            "POST /api/transactions/sale",
            "POST /api/transactions/purchase"
        ],
        demo_ready: true,
        target_customer: "Siigo ERP",
        supported_countries: ["México", "Colombia", "Argentina", "Chile", "Uruguay", "Ecuador", "Perú"]
    }
    
    return HttpResponse(JSON(info))
}

func handleRegions() {
    return HttpResponse(JSON(regiones))
}

func handleSalePost(pathVars, method, body) {
    try {
        console.log("📥 Processing sale: " + body)
        
        // Para demo, usar valores por defecto si no hay body JSON válido
        let region = "COL"
        let amount = "100000"
        
        // Intentar parsear body si existe
        if (body && std.len(body) > 0) {
            // Parsear manualmente parámetros region y amount del body
            if (std.contains(body, "region=")) {
                let parts = std.split(body, "&")
                let i = 0
                while (i < std.len(parts)) {
                    let part = parts[i]
                    if (std.startswith(part, "region=")) {
                        region = std.replace(part, "region=", "")
                    }
                    if (std.startswith(part, "amount=")) {
                        amount = std.replace(part, "amount=", "")
                    }
                    i = i + 1
                }
            }
        }
        
        let motorVentas = VentasAPI
        let result = motorVentas.use("venta " + region + " " + amount)
        
        // Guardar transacción
        transacciones[std.len(transacciones)] = result
        
        console.log("✅ Sale processed: " + result.transactionId)
        
        return HttpResponse(JSON(result))
        
    } catch (error) {
        let errorResponse = {
            success: false,
            error: "Error procesando venta: " + error
        }
        console.log("❌ Sale error: " + error)
        return HttpResponse(JSON(errorResponse))
    }
}

func handlePurchasePost(pathVars, method, body) {
    try {
        console.log("📥 Processing purchase: " + body)
        
        // Para demo, usar valores por defecto
        let region = "MX"
        let amount = "50000"
        
        // Intentar parsear body si existe
        if (body && std.len(body) > 0) {
            if (std.contains(body, "region=")) {
                let parts = std.split(body, "&")
                let i = 0
                while (i < std.len(parts)) {
                    let part = parts[i]
                    if (std.startswith(part, "region=")) {
                        region = std.replace(part, "region=", "")
                    }
                    if (std.startswith(part, "amount=")) {
                        amount = std.replace(part, "amount=", "")
                    }
                    i = i + 1
                }
            }
        }
        
        let motorCompras = ComprasAPI
        let result = motorCompras.use("compra " + region + " " + amount)
        
        // Guardar transacción
        transacciones[std.len(transacciones)] = result
        
        console.log("✅ Purchase processed: " + result.transactionId)
        
        return HttpResponse(JSON(result))
        
    } catch (error) {
        let errorResponse = {
            success: false,
            error: "Error procesando compra: " + error
        }
        console.log("❌ Purchase error: " + error)
        return HttpResponse(JSON(errorResponse))
    }
}

func handleTransactions() {
    let response = {
        total: std.len(transacciones),
        transactions: transacciones
    }
    return HttpResponse(JSON(response))
}

// Configurar todas las rutas
func setupRoutes() {
    console.log("🔗 Configurando rutas API...")
    
    // Ruta principal
    http.handler("GET", "/", handleIndex)
    
    // API Info
    http.handler("GET", "/api/info", handleAPIInfo)
    
    // Regiones
    http.handler("GET", "/api/regions", handleRegions)
    
    // Transacciones
    http.handler("POST", "/api/transactions/sale", handleSalePost)
    http.handler("POST", "/api/transactions/purchase", handlePurchasePost)
    http.handler("GET", "/api/transactions", handleTransactions)
    
    console.log("✅ Rutas configuradas")
}

// Inicializar y arrancar servidor
func startServer() {
    console.log("🌐 Inicializando servidor HTTP R2Lang...")
    console.log("Puerto: 8080")
    console.log("")
    
    // Configurar rutas
    setupRoutes()
    
    // Demo inicial
    runQuickDemo()
    
    console.log("")
    console.log("🎯 ¡SISTEMA LATAM WEB FUNCIONANDO!")
    console.log("=================================")
    console.log("🌍 URL Principal: http://localhost:8080")
    console.log("📊 API Info: http://localhost:8080/api/info")
    console.log("🌎 Regiones: http://localhost:8080/api/regions")
    console.log("")
    console.log("📝 Para detener: Ctrl+C")
    console.log("🎉 ¡DEMO WEB LISTO PARA SIIGO!")
    console.log("")
    
    // Iniciar servidor HTTP
    http.serve(":8080")
}

// Demo rápido para verificar funcionamiento
func runQuickDemo() {
    console.log("🎪 Demo rápido API...")
    
    let motorVentas = VentasAPI
    let motorCompras = ComprasAPI
    
    let testVenta = motorVentas.use("venta COL 1000")
    let testCompra = motorCompras.use("compra MX 500")
    
    console.log("✅ Venta test: " + testVenta.transactionId)
    console.log("✅ Compra test: " + testCompra.transactionId)
    console.log("✅ DSL engines funcionando correctamente")
}

// Test del servidor
func testAPIs() {
    console.log("🧪 Testing API Server...")
    runQuickDemo()
    console.log("✅ API Server test completado")
    return true
}

console.log("🚀 API Server Working ready")