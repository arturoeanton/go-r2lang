// API Server Simplificado para Sistema Contable LATAM
// Usando r2http nativo de R2Lang

console.log("🚀 Módulo API Server cargado")

// Configuración del servidor
let serverConfig = {
    port: 8080,
    host: "localhost"
}

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

// Inicializar servidor
func initServer() {
    console.log("🌐 Inicializando servidor HTTP...")
    console.log("Puerto: " + serverConfig.port)
    console.log("Host: " + serverConfig.host)
    
    // Crear servidor HTTP
    let server = createHTTPServer(serverConfig.port)
    
    // Configurar rutas
    setupRoutes(server)
    
    console.log("✅ Servidor HTTP configurado")
    return server
}

// Configurar rutas del servidor
func setupRoutes(server) {
    console.log("🔗 Configurando rutas API...")
    
    // Ruta principal
    server.route("/", "GET", handleIndex)
    
    // API Info
    server.route("/api/info", "GET", handleAPIInfo)
    
    // Regiones
    server.route("/api/regions", "GET", handleRegions)
    
    // Transacciones
    server.route("/api/transactions/sale", "POST", handleSale)
    server.route("/api/transactions/purchase", "POST", handlePurchase)
    
    // Archivos estáticos
    server.static("/static", "./static")
    
    console.log("✅ Rutas configuradas")
}

// Crear servidor HTTP usando r2http
func createHTTPServer(port) {
    console.log("🔧 Creando servidor r2http en puerto " + port)
    
    let server = http.server()
    server.port(port)
    
    return server
}

// Handlers de rutas

func handleIndex(request, response) {
    response.header("Content-Type", "text/html")
    response.send(`
        <!DOCTYPE html>
        <html>
        <head>
            <title>Sistema Contable LATAM</title>
        </head>
        <body>
            <h1>🌍 Sistema Contable LATAM - R2Lang DSL</h1>
            <p>Demo para Siigo ERP Localization</p>
            <h2>🔗 Enlaces:</h2>
            <ul>
                <li><a href="/api/info">📊 API Info</a></li>
                <li><a href="/api/regions">🌎 Regiones</a></li>
                <li><a href="/static/index.html">📱 Frontend</a></li>
            </ul>
        </body>
        </html>
    `)
}

func handleAPIInfo(request, response) {
    let info = {
        name: "Sistema Contable LATAM API",
        version: "1.0.0",
        description: "R2Lang DSL para localización ERP",
        regions: std.len(regiones),
        endpoints: [
            "GET /api/info",
            "GET /api/regions", 
            "POST /api/transactions/sale",
            "POST /api/transactions/purchase"
        ],
        demo_ready: true
    }
    
    response.json(info)
}

func handleRegions(request, response) {
    response.json(regiones)
}

func handleSale(request, response) {
    try {
        let body = request.body()
        let region = body.region
        let amount = body.amount
        
        let motorVentas = VentasAPI
        let result = motorVentas.use("venta " + region + " " + amount)
        
        response.json(result)
        
    } catch (error) {
        response.status(400)
        response.json({
            success: false,
            error: "Error procesando venta: " + error
        })
    }
}

func handlePurchase(request, response) {
    try {
        let body = request.body()
        let region = body.region
        let amount = body.amount
        
        let motorCompras = ComprasAPI
        let result = motorCompras.use("compra " + region + " " + amount)
        
        response.json(result)
        
    } catch (error) {
        response.status(400)
        response.json({
            success: false,
            error: "Error procesando compra: " + error
        })
    }
}

// Test del servidor
func testAPIs() {
    console.log("🧪 Testing API Server...")
    
    let motorVentas = VentasAPI
    let motorCompras = ComprasAPI
    
    let testVenta = motorVentas.use("venta COL 1000")
    let testCompra = motorCompras.use("compra MX 500")
    
    console.log("✅ Venta test: " + testVenta.transactionId)
    console.log("✅ Compra test: " + testCompra.transactionId)
    console.log("✅ API Server test completado")
    
    return true
}

console.log("🚀 API Server Simple ready")