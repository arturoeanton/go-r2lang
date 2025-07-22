// Demo Sistema Contable LATAM - Versión Funcional
// POC para Siigo

console.log("🌍 Sistema Contable LATAM - Demo Siigo")
console.log("")

// Configuración regiones
let regiones = {
    "MX": { nombre: "México", moneda: "MXN", iva: 0.16 },
    "COL": { nombre: "Colombia", moneda: "COP", iva: 0.19 },
    "AR": { nombre: "Argentina", moneda: "ARS", iva: 0.21 },
    "CH": { nombre: "Chile", moneda: "CLP", iva: 0.19 },
    "UY": { nombre: "Uruguay", moneda: "UYU", iva: 0.22 },
    "EC": { nombre: "Ecuador", moneda: "USD", iva: 0.12 },
    "PE": { nombre: "Perú", moneda: "PEN", iva: 0.18 }
}

let transacciones = []

// Función procesar
func procesarTx(tipo, region, importe) {
    let config = regiones[region]
    let importeNum = std.parseFloat(importe)
    let iva = math.round((importeNum * config.iva) * 100) / 100
    let total = importeNum + iva
    
    let tx = {
        id: region + "-" + math.randomInt(9999),
        tipo: tipo,
        region: region,
        pais: config.nombre,
        importe: importeNum,
        iva: iva,
        total: total,
        moneda: config.moneda
    }
    
    transacciones[std.len(transacciones)] = tx
    return tx
}

// Demo consola
func demoConsola() {
    console.log("=== DEMO SISTEMA ===")
    
    let tx1 = procesarTx("venta", "MX", "100000")
    console.log("✓ Venta México: $" + tx1.total + " " + tx1.moneda)
    
    let tx2 = procesarTx("compra", "COL", "50000")
    console.log("✓ Compra Colombia: $" + tx2.total + " " + tx2.moneda)
    
    let tx3 = procesarTx("venta", "AR", "75000")
    console.log("✓ Venta Argentina: $" + tx3.total + " " + tx3.moneda)
    
    console.log("\nTotal transacciones: " + std.len(transacciones))
    console.log("\n🎯 VALUE PROPOSITION:")
    console.log("• 18 meses → 2 meses")
    console.log("• $500K → $150K")
    console.log("• ROI: 1,020%")
}

// Parsear params simple
func getParam(body, param) {
    if (!body || std.len(body) == 0) {
        return ""
    }
    
    let pairs = std.split(body, "&")
    let i = 0
    while (i < std.len(pairs)) {
        let pair = pairs[i]
        if (std.contains(pair, param + "=")) {
            let parts = std.split(pair, "=")
            if (std.len(parts) >= 2) {
                return std.replace(parts[1], "+", " ")
            }
        }
        i = i + 1
    }
    return ""
}

// Handler home
func handleHome(pathVars, method, body) {
    let html = "<!DOCTYPE html>\n<html>\n<head>\n"
    html = html + "<title>Sistema Contable LATAM</title>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<style>\n"
    html = html + "body{font-family:Arial;margin:20px;background:#f5f5f5;}\n"
    html = html + ".container{max-width:800px;margin:0 auto;background:white;padding:30px;border-radius:10px;}\n"
    html = html + "h1{color:#2c3e50;text-align:center;}\n"
    html = html + ".card{background:#f8f9fa;padding:20px;border-radius:8px;margin:15px 0;}\n"
    html = html + "input,select{width:100%;padding:8px;margin:5px 0;}\n"
    html = html + "button{background:#667eea;color:white;padding:10px;border:none;border-radius:4px;width:100%;}\n"
    html = html + ".value{background:#84fab0;padding:20px;border-radius:8px;margin:20px 0;}\n"
    html = html + "</style>\n</head>\n<body>\n"
    
    html = html + "<div class=\"container\">\n"
    html = html + "<h1>Sistema Contable LATAM</h1>\n"
    html = html + "<p style=\"text-align:center\">Demo Siigo - " + std.len(transacciones) + " transacciones</p>\n"
    
    html = html + "<div class=\"card\">\n"
    html = html + "<h3>Procesar Transacción</h3>\n"
    html = html + "<form action=\"/procesar\" method=\"POST\">\n"
    html = html + "<select name=\"tipo\">\n"
    html = html + "<option value=\"venta\">Venta</option>\n"
    html = html + "<option value=\"compra\">Compra</option>\n"
    html = html + "</select>\n"
    html = html + "<select name=\"region\">\n"
    html = html + "<option value=\"MX\">México (16%)</option>\n"
    html = html + "<option value=\"COL\">Colombia (19%)</option>\n"
    html = html + "<option value=\"AR\">Argentina (21%)</option>\n"
    html = html + "<option value=\"CH\">Chile (19%)</option>\n"
    html = html + "<option value=\"UY\">Uruguay (22%)</option>\n"
    html = html + "<option value=\"EC\">Ecuador (12%)</option>\n"
    html = html + "<option value=\"PE\">Perú (18%)</option>\n"
    html = html + "</select>\n"
    html = html + "<input type=\"number\" name=\"importe\" value=\"100000\">\n"
    html = html + "<button type=\"submit\">Procesar</button>\n"
    html = html + "</form>\n"
    html = html + "</div>\n"
    
    html = html + "<p style=\"text-align:center\">\n"
    html = html + "<a href=\"/demo\">Demo Auto</a> | "
    html = html + "<a href=\"/api\">Ver API</a>\n"
    html = html + "</p>\n"
    
    html = html + "<div class=\"value\">\n"
    html = html + "<h3>Value Proposition Siigo</h3>\n"
    html = html + "<p>18 meses → 2 meses | $500K → $150K | 7 sistemas → 1 DSL | ROI: 1,020%</p>\n"
    html = html + "</div>\n"
    
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// Handler procesar
func handleProcesar(pathVars, method, body) {
    let tipo = getParam(body, "tipo")
    let region = getParam(body, "region")
    let importe = getParam(body, "importe")
    
    if (tipo == "") { tipo = "venta" }
    if (region == "") { region = "COL" }
    if (importe == "") { importe = "100000" }
    
    let tx = procesarTx(tipo, region, importe)
    
    let html = "<!DOCTYPE html>\n<html>\n<head>\n"
    html = html + "<title>Comprobante</title>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<style>body{font-family:Arial;margin:20px;}.comp{max-width:600px;margin:0 auto;padding:30px;border:2px solid #667eea;border-radius:10px;}</style>\n"
    html = html + "</head>\n<body>\n"
    html = html + "<div class=\"comp\">\n"
    html = html + "<h1>COMPROBANTE</h1>\n"
    html = html + "<h2>" + tx.pais + "</h2>\n"
    html = html + "<p><strong>ID:</strong> " + tx.id + "</p>\n"
    html = html + "<p><strong>Tipo:</strong> " + tx.tipo + "</p>\n"
    html = html + "<p><strong>Importe:</strong> " + tx.moneda + " " + tx.importe + "</p>\n"
    html = html + "<p><strong>IVA:</strong> " + tx.moneda + " " + tx.iva + "</p>\n"
    html = html + "<hr>\n"
    html = html + "<p><strong>TOTAL:</strong> " + tx.moneda + " " + tx.total + "</p>\n"
    html = html + "<p style=\"text-align:center;color:green\">✅ PROCESADO</p>\n"
    html = html + "<p style=\"text-align:center\"><a href=\"/\">Volver</a></p>\n"
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// Handler demo
func handleDemo(pathVars, method, body) {
    procesarTx("venta", "MX", "100000")
    procesarTx("compra", "COL", "50000")
    procesarTx("venta", "AR", "75000")
    procesarTx("compra", "PE", "60000")
    
    let html = "<!DOCTYPE html>\n<html>\n<head>\n"
    html = html + "<title>Demo</title>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<style>body{font-family:Arial;margin:20px;}.container{max-width:600px;margin:0 auto;padding:30px;background:white;border-radius:10px;}</style>\n"
    html = html + "</head>\n<body>\n"
    html = html + "<div class=\"container\">\n"
    html = html + "<h1>Demo Automática</h1>\n"
    html = html + "<p>✅ Se procesaron 4 transacciones de diferentes regiones</p>\n"
    html = html + "<p>Total acumulado: " + std.len(transacciones) + " transacciones</p>\n"
    html = html + "<p style=\"text-align:center;margin-top:30px;\">\n"
    html = html + "<a href=\"/\">Volver</a> | <a href=\"/api\">Ver API</a>\n"
    html = html + "</p>\n"
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// Handler API
func handleAPI(pathVars, method, body) {
    return JSON({
        total: std.len(transacciones),
        regiones: 7,
        transacciones: transacciones
    })
}

// Main
func main() {
    // Demo consola
    demoConsola()
    
    console.log("\n✅ Iniciando servidor web...")
    
    // Rutas
    http.handler("GET", "/", handleHome)
    http.handler("POST", "/procesar", handleProcesar)
    http.handler("GET", "/demo", handleDemo)
    http.handler("GET", "/api", handleAPI)
    
    console.log("🌐 http://localhost:8080")
    console.log("")
    
    // Servidor
    http.serve(":8080")
}

main()