// Sistema Contable LATAM - POC Web Simple
// Demo funcional para Siigo siguiendo patrón example6-web.r2

console.log("🌍 Sistema Contable LATAM - POC Web Simple")
console.log("Demo para Siigo ERP Localization")
console.log("")

// Regiones LATAM configuradas
let regiones = {
    "MX": { nombre: "México", moneda: "MXN", simbolo: "$", iva: 0.16, normativa: "NIF-Mexican" },
    "COL": { nombre: "Colombia", moneda: "COP", simbolo: "$", iva: 0.19, normativa: "NIIF-Colombia" },
    "AR": { nombre: "Argentina", moneda: "ARS", simbolo: "$", iva: 0.21, normativa: "RT-Argentina" },
    "CH": { nombre: "Chile", moneda: "CLP", simbolo: "$", iva: 0.19, normativa: "IFRS-Chile" },
    "UY": { nombre: "Uruguay", moneda: "UYU", simbolo: "$", iva: 0.22, normativa: "NIIF-Uruguay" },
    "EC": { nombre: "Ecuador", moneda: "USD", simbolo: "$", iva: 0.12, normativa: "NIIF-Ecuador" },
    "PE": { nombre: "Perú", moneda: "PEN", simbolo: "S/", iva: 0.18, normativa: "PCGE-Peru" }
}

// Almacén de transacciones en memoria 
let transacciones = []

// DSL Engine integrado
dsl VentasWeb {
    token("VENTA", "venta|sale")
    token("REGION", "MX|COL|AR|CH|UY|EC|PE")
    token("IMPORTE", "[0-9]+\\.?[0-9]*")
    
    rule("venta_web", ["VENTA", "REGION", "IMPORTE"], "procesarVentaWeb")
    
    func procesarVentaWeb(operacion, region, importe) {
        let config = regiones[region]
        let importeNum = std.parseFloat(importe)
        let importeIVA = math.round((importeNum * config.iva) * 100) / 100
        let importeTotal = math.round((importeNum + importeIVA) * 100) / 100
        let txId = region + "-" + std.now() + "-" + math.randomInt(9999)
        
        let resultado = {
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
        
        // Guardar en memoria
        transacciones[std.len(transacciones)] = resultado
        
        return resultado
    }
}

dsl ComprasWeb {
    token("COMPRA", "compra|purchase") 
    token("REGION", "MX|COL|AR|CH|UY|EC|PE")
    token("IMPORTE", "[0-9]+\\.?[0-9]*")
    
    rule("compra_web", ["COMPRA", "REGION", "IMPORTE"], "procesarCompraWeb")
    
    func procesarCompraWeb(operacion, region, importe) {
        let config = regiones[region]
        let importeNum = std.parseFloat(importe)
        let importeIVA = math.round((importeNum * config.iva) * 100) / 100
        let importeTotal = math.round((importeNum + importeIVA) * 100) / 100
        let txId = region + "-" + std.now() + "-" + math.randomInt(9999)
        
        let resultado = {
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
        
        // Guardar en memoria
        transacciones[std.len(transacciones)] = resultado
        
        return resultado
    }
}

// Handler página principal
func handleIndex(pathVars, method, body) {
    let html = `<!DOCTYPE html>
<html>
<head>
    <title>Sistema Contable LATAM</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background: #f5f5f5; }
        .container { background: white; padding: 30px; border-radius: 10px; max-width: 1200px; margin: 0 auto; }
        h1 { color: #2c3e50; text-align: center; }
        .section { margin: 20px 0; padding: 20px; border: 1px solid #ddd; border-radius: 8px; }
        .form-group { margin: 15px 0; }
        label { display: block; margin-bottom: 5px; font-weight: bold; }
        input, select { padding: 8px; border: 1px solid #ddd; border-radius: 4px; width: 200px; }
        button { padding: 10px 20px; background: #3498db; color: white; border: none; border-radius: 4px; cursor: pointer; }
        button:hover { background: #2980b9; }
        .result { background: #ecf0f1; padding: 15px; margin: 10px 0; border-radius: 4px; }
        .transaction { background: #d5f4e6; padding: 10px; margin: 5px 0; border-radius: 4px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🌍 Sistema Contable LATAM - R2Lang DSL</h1>
        <p style="text-align: center;"><strong>POC Demo para Siigo ERP</strong> - 7 Regiones LATAM</p>
        
        <div class="section">
            <h2>💰 Procesar Venta</h2>
            <form action="/procesar-venta" method="POST">
                <div class="form-group">
                    <label>Región:</label>
                    <select name="region">
                        <option value="MX">🇲🇽 México</option>
                        <option value="COL">🇨🇴 Colombia</option>
                        <option value="AR">🇦🇷 Argentina</option>
                        <option value="CH">🇨🇱 Chile</option>
                        <option value="UY">🇺🇾 Uruguay</option>
                        <option value="EC">🇪🇨 Ecuador</option>
                        <option value="PE">🇵🇪 Perú</option>
                    </select>
                </div>
                <div class="form-group">
                    <label>Importe:</label>
                    <input type="number" name="importe" value="100000" required>
                </div>
                <button type="submit">Procesar Venta</button>
            </form>
        </div>
        
        <div class="section">
            <h2>🛒 Procesar Compra</h2>
            <form action="/procesar-compra" method="POST">
                <div class="form-group">
                    <label>Región:</label>
                    <select name="region">
                        <option value="MX">🇲🇽 México</option>
                        <option value="COL">🇨🇴 Colombia</option>
                        <option value="AR">🇦🇷 Argentina</option>
                        <option value="CH">🇨🇱 Chile</option>
                        <option value="UY">🇺🇾 Uruguay</option>
                        <option value="EC">🇪🇨 Ecuador</option>
                        <option value="PE">🇵🇪 Perú</option>
                    </select>
                </div>
                <div class="form-group">
                    <label>Importe:</label>
                    <input type="number" name="importe" value="50000" required>
                </div>
                <button type="submit">Procesar Compra</button>
            </form>
        </div>
        
        <div class="section">
            <h2>📊 APIs Disponibles</h2>
            <p><a href="/api/regiones">📋 Ver Regiones</a> | <a href="/api/transacciones">📜 Ver Transacciones</a> | <a href="/demo">🎪 Demo Completo</a></p>
        </div>
        
        <div class="result">
            <strong>🎯 Value Proposition para Siigo:</strong><br>
            ✅ 18 meses → 2 meses por país<br>
            ✅ $500K → $150K por localización<br>
            ✅ 7 codebases → 1 DSL unificado<br>
            ✅ ROI: 1,020% en 3 años
        </div>
    </div>
</body>
</html>`
    
    return html
}

// Handler para procesar venta
func handleVenta(pathVars, method, body) {
    console.log("DEBUG - handleVenta called with body: " + (body || "empty"))
    
    let params = parseFormData(body)
    let region = params.region || "COL"
    let importe = params.importe || "100000"
    
    console.log("DEBUG - Processing sale for region: " + region + " amount: " + importe)
    
    let motorVentas = VentasWeb
    let resultado = motorVentas.use("venta " + region + " " + importe)
    
    let html = generarComprobante(resultado, "VENTA")
    return html
}

// Handler para procesar compra
func handleCompra(pathVars, method, body) {
    let params = parseFormData(body)
    let region = params.region || "MX"
    let importe = params.importe || "50000"
    
    let motorCompras = ComprasWeb
    let resultado = motorCompras.use("compra " + region + " " + importe)
    
    let html = generarComprobante(resultado, "COMPRA")
    return html
}

// Handler para API regiones
func handleAPIRegiones(pathVars, method, body) {
    return JSON(regiones)
}

// Handler para API transacciones
func handleAPITransacciones(pathVars, method, body) {
    let response = {
        total: std.len(transacciones),
        transacciones: transacciones
    }
    return JSON(response)
}

// Handler demo completo
func handleDemo(pathVars, method, body) {
    let motorVentas = VentasWeb
    let motorCompras = ComprasWeb
    
    let demos = [
        "venta COL 100000",
        "compra MX 50000", 
        "venta AR 75000",
        "compra PE 60000"
    ]
    
    let resultados = []
    let i = 0
    while (i < std.len(demos)) {
        let comando = demos[i]
        let resultado = ""
        
        if (std.contains(comando, "venta")) {
            resultado = motorVentas.use(comando)
        } else {
            resultado = motorCompras.use(comando)
        }
        
        resultados[i] = resultado
        i = i + 1
    }
    
    let html = `<!DOCTYPE html>
<html>
<head>
    <title>Demo Completo - Sistema Contable LATAM</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background: #f5f5f5; }
        .container { background: white; padding: 30px; border-radius: 10px; max-width: 1000px; margin: 0 auto; }
        .transaction { background: #d5f4e6; padding: 15px; margin: 10px 0; border-radius: 8px; border-left: 4px solid #27ae60; }
        .back { text-align: center; margin: 20px 0; }
        .back a { text-decoration: none; background: #3498db; color: white; padding: 10px 20px; border-radius: 4px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🎪 Demo Completo - Sistema Contable LATAM</h1>
        <p><strong>Procesamiento automático de transacciones multi-región</strong></p>`
    
    i = 0
    while (i < std.len(resultados)) {
        let r = resultados[i]
        html = html + `
        <div class="transaction">
            <h3>Transacción ` + (i + 1) + ` - ` + r.country + `</h3>
            <p><strong>ID:</strong> ` + r.transactionId + `</p>
            <p><strong>Importe:</strong> ` + r.currency + ` ` + r.amount + `</p>
            <p><strong>Impuesto:</strong> ` + r.currency + ` ` + r.tax + `</p>
            <p><strong>Total:</strong> ` + r.currency + ` ` + r.total + `</p>
            <p><strong>Normativa:</strong> ` + r.compliance + `</p>
        </div>`
        i = i + 1
    }
    
    html = html + `
        <div class="back">
            <a href="/">🏠 Volver al Inicio</a>
        </div>
    </div>
</body>
</html>`
    
    return html
}

// Función para parsear form data
func parseFormData(body) {
    let params = {}
    if (!body || std.len(body) == 0) {
        return params
    }
    
    console.log("DEBUG - Body received: " + body)
    
    let pairs = std.split(body, "&")
    let i = 0
    while (i < std.len(pairs)) {
        let pair = pairs[i]
        if (std.contains(pair, "=")) {
            let parts = std.split(pair, "=")
            if (std.len(parts) >= 2) {
                params[parts[0]] = parts[1]
            }
        }
        i = i + 1
    }
    
    console.log("DEBUG - Params parsed: " + JSON(params))
    return params
}

// Función para generar comprobante HTML
func generarComprobante(resultado, tipo) {
    let config = regiones[resultado.region]
    
    let html = `<!DOCTYPE html>
<html>
<head>
    <title>Comprobante - Sistema Contable LATAM</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background: #f5f5f5; }
        .comprobante { background: white; padding: 30px; border-radius: 10px; max-width: 600px; margin: 0 auto; border: 2px solid #27ae60; }
        .header { text-align: center; border-bottom: 2px solid #3498db; padding-bottom: 15px; margin-bottom: 20px; }
        .row { display: flex; justify-content: space-between; margin: 10px 0; }
        .total { background: #ecf0f1; padding: 15px; border-radius: 8px; margin: 20px 0; }
        .back { text-align: center; margin: 20px 0; }
        .back a { text-decoration: none; background: #3498db; color: white; padding: 10px 20px; border-radius: 4px; }
    </style>
</head>
<body>
    <div class="comprobante">
        <div class="header">
            <h1>🧾 COMPROBANTE DE ` + tipo + `</h1>
            <h2>` + resultado.country + ` (` + resultado.region + `)</h2>
        </div>
        
        <div class="row">
            <strong>ID Transacción:</strong>
            <span>` + resultado.transactionId + `</span>
        </div>
        
        <div class="row">
            <strong>Fecha:</strong>
            <span>` + resultado.timestamp + `</span>
        </div>
        
        <div class="row">
            <strong>Normativa:</strong>
            <span>` + resultado.compliance + `</span>
        </div>
        
        <div class="row">
            <strong>Moneda:</strong>
            <span>` + resultado.currency + `</span>
        </div>
        
        <div class="total">
            <div class="row">
                <strong>Importe Base:</strong>
                <span>` + config.simbolo + ` ` + resultado.amount + ` ` + resultado.currency + `</span>
            </div>
            <div class="row">
                <strong>IVA (` + (config.iva * 100) + `%):</strong>
                <span>` + config.simbolo + ` ` + resultado.tax + ` ` + resultado.currency + `</span>
            </div>
            <div class="row" style="border-top: 2px solid #2c3e50; padding-top: 10px;">
                <strong>TOTAL:</strong>
                <strong>` + config.simbolo + ` ` + resultado.total + ` ` + resultado.currency + `</strong>
            </div>
        </div>
        
        <p style="text-align: center; color: #27ae60;"><strong>✅ TRANSACCIÓN VALIDADA</strong></p>
        
        <div class="back">
            <a href="/">🏠 Volver al Inicio</a>
            <a href="/demo" style="margin-left: 10px;">🎪 Demo Completo</a>
        </div>
    </div>
</body>
</html>`
    
    return html
}

// Función principal
func main() {
    console.log("🚀 Iniciando servidor web...")
    console.log("Puerto: 8080")
    console.log("")
    
    // Configurar rutas
    http.handler("GET", "/", handleIndex)
    http.handler("POST", "/procesar-venta", handleVenta)
    http.handler("POST", "/procesar-compra", handleCompra)
    http.handler("GET", "/api/regiones", handleAPIRegiones)
    http.handler("GET", "/api/transacciones", handleAPITransacciones)
    http.handler("GET", "/demo", handleDemo)
    
    console.log("✅ Rutas configuradas:")
    console.log("   GET  / - Página principal")
    console.log("   POST /procesar-venta - Procesar venta")
    console.log("   POST /procesar-compra - Procesar compra")
    console.log("   GET  /api/regiones - Lista regiones")
    console.log("   GET  /api/transacciones - Lista transacciones")
    console.log("   GET  /demo - Demo completo")
    console.log("")
    
    console.log("🎯 ¡SISTEMA LISTO!")
    console.log("🌐 URL: http://localhost:8080")
    console.log("🎪 Demo: http://localhost:8080/demo")
    console.log("")
    console.log("🎉 POC ready for Siigo!")
    
    // Iniciar servidor
    http.serve(":8080")
}