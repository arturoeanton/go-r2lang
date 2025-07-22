// Demo Final Sistema Contable LATAM para Siigo
// POC funcional con DSL y procesamiento multi-región

console.log("🌍 Sistema Contable LATAM - Demo Siigo")
console.log("📊 DSL Reportes + Multi-región")
console.log("")

// Configuración de regiones
let regiones = {}
regiones["MX"] = { nombre: "México", moneda: "MXN", simbolo: "$", iva: 0.16 }
regiones["COL"] = { nombre: "Colombia", moneda: "COP", simbolo: "$", iva: 0.19 }
regiones["AR"] = { nombre: "Argentina", moneda: "ARS", simbolo: "$", iva: 0.21 }
regiones["CH"] = { nombre: "Chile", moneda: "CLP", simbolo: "$", iva: 0.19 }
regiones["UY"] = { nombre: "Uruguay", moneda: "UYU", simbolo: "$", iva: 0.22 }
regiones["EC"] = { nombre: "Ecuador", moneda: "USD", simbolo: "$", iva: 0.12 }
regiones["PE"] = { nombre: "Perú", moneda: "PEN", simbolo: "S/", iva: 0.18 }

let transacciones = []

// DSL para reportes
dsl ReportesContables {
    token("REPORTE", "reporte")
    token("TIPO", "ventas|compras|todo")
    token("REGION", "MX|COL|AR|CH|UY|EC|PE|ALL")
    
    rule("generar", ["REPORTE", "TIPO", "REGION"], "generarReporte")
    
    func generarReporte(cmd, tipo, region) {
        let total = 0
        let count = 0
        let i = 0
        while (i < std.len(transacciones)) {
            let tx = transacciones[i]
            let incluir = true
            
            if (region != "ALL" && tx.region != region) {
                incluir = false
            }
            
            if (tipo != "todo" && tx.tipo != tipo) {
                incluir = false
            }
            
            if (incluir) {
                total = total + tx.total
                count = count + 1
            }
            
            i = i + 1
        }
        
        return {
            tipo: tipo,
            region: region,
            transacciones: count,
            total: math.round(total * 100) / 100
        }
    }
}

// Función para procesar transacciones
func procesarTransaccion(tipo, region, importe) {
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
        moneda: config.moneda,
        fecha: std.now()
    }
    
    transacciones[std.len(transacciones)] = tx
    return tx
}

// Demo en consola
func ejecutarDemo() {
    console.log("\n=== DEMO SISTEMA CONTABLE ===\n")
    
    console.log("1. Procesando transacciones...")
    
    let tx1 = procesarTransaccion("ventas", "MX", "100000")
    console.log("   ✓ Venta México: $" + tx1.total + " " + tx1.moneda)
    
    let tx2 = procesarTransaccion("compras", "COL", "50000")
    console.log("   ✓ Compra Colombia: $" + tx2.total + " " + tx2.moneda)
    
    let tx3 = procesarTransaccion("ventas", "AR", "75000")
    console.log("   ✓ Venta Argentina: $" + tx3.total + " " + tx3.moneda)
    
    let tx4 = procesarTransaccion("compras", "PE", "60000")
    console.log("   ✓ Compra Perú: " + tx4.moneda + " " + tx4.total)
    
    console.log("\n2. Generando reportes DSL...")
    
    let engine = ReportesContables
    
    let rep1 = engine.use("reporte todo ALL")
    console.log("   ✓ Todas las transacciones: " + rep1.transacciones + " ops, Total: $" + rep1.total)
    
    let rep2 = engine.use("reporte ventas ALL")
    console.log("   ✓ Solo ventas: " + rep2.transacciones + " ops, Total: $" + rep2.total)
    
    let rep3 = engine.use("reporte todo MX")
    console.log("   ✓ México completo: " + rep3.transacciones + " ops, Total: $" + rep3.total)
    
    console.log("\n=== RESUMEN ===")
    console.log("Total transacciones: " + std.len(transacciones))
    console.log("Regiones activas: 7")
    
    console.log("\n🎯 VALUE PROPOSITION:")
    console.log("• Tiempo: 18 meses → 2 meses")
    console.log("• Costo: $500K → $150K")
    console.log("• ROI: 1,020%")
}

// Parsear form data
func parseFormData(body) {
    let params = {}
    if (!body || std.len(body) == 0) {
        return params
    }
    
    let pairs = std.split(body, "&")
    let i = 0
    while (i < std.len(pairs)) {
        let pair = pairs[i]
        if (std.contains(pair, "=")) {
            let parts = std.split(pair, "=")
            if (std.len(parts) >= 2) {
                let key = parts[0]
                let value = parts[1]
                value = std.replace(value, "+", " ")
                params[key] = value
            }
        }
        i = i + 1
    }
    return params
}

// Página principal
func handleHome(pathVars, method, body) {
    let html = "<!DOCTYPE html>\n"
    html = html + "<html>\n<head>\n"
    html = html + "<title>Sistema Contable LATAM</title>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<style>\n"
    html = html + "body{font-family:Arial,sans-serif;margin:20px;background:#f5f5f5;}\n"
    html = html + ".container{max-width:800px;margin:0 auto;background:white;padding:30px;border-radius:10px;}\n"
    html = html + "h1{color:#2c3e50;text-align:center;}\n"
    html = html + ".card{background:#f8f9fa;padding:20px;border-radius:8px;margin:15px 0;}\n"
    html = html + "input,select{width:100%;padding:8px;margin:5px 0;border:1px solid #ddd;border-radius:4px;}\n"
    html = html + "button{background:#667eea;color:white;padding:10px;border:none;border-radius:4px;width:100%;cursor:pointer;}\n"
    html = html + "</style>\n</head>\n<body>\n"
    
    html = html + "<div class=\"container\">\n"
    html = html + "<h1>Sistema Contable LATAM</h1>\n"
    html = html + "<p style=\"text-align:center\">POC para Siigo - " + std.len(transacciones) + " transacciones</p>\n"
    
    // Formulario
    html = html + "<div class=\"card\">\n"
    html = html + "<h3>Procesar Transacción</h3>\n"
    html = html + "<form action=\"/procesar\" method=\"POST\">\n"
    html = html + "<select name=\"tipo\">\n"
    html = html + "<option value=\"ventas\">Venta</option>\n"
    html = html + "<option value=\"compras\">Compra</option>\n"
    html = html + "</select>\n"
    html = html + "<select name=\"region\">\n"
    html = html + "<option value=\"MX\">México</option>\n"
    html = html + "<option value=\"COL\">Colombia</option>\n"
    html = html + "<option value=\"AR\">Argentina</option>\n"
    html = html + "<option value=\"CH\">Chile</option>\n"
    html = html + "<option value=\"UY\">Uruguay</option>\n"
    html = html + "<option value=\"EC\">Ecuador</option>\n"
    html = html + "<option value=\"PE\">Perú</option>\n"
    html = html + "</select>\n"
    html = html + "<input type=\"number\" name=\"importe\" value=\"100000\" required>\n"
    html = html + "<button type=\"submit\">Procesar</button>\n"
    html = html + "</form>\n"
    html = html + "</div>\n"
    
    // Links
    html = html + "<p style=\"text-align:center\">\n"
    html = html + "<a href=\"/demo\">Demo Auto</a> | "
    html = html + "<a href=\"/api/transacciones\">API</a>\n"
    html = html + "</p>\n"
    
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// Procesar transacción
func handleProcesar(pathVars, method, body) {
    let params = parseFormData(body)
    let tipo = params.tipo || "ventas"
    let region = params.region || "COL" 
    let importe = params.importe || "100000"
    
    let tx = procesarTransaccion(tipo, region, importe)
    
    let html = "<!DOCTYPE html>\n<html>\n<head>\n"
    html = html + "<title>Comprobante</title>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<style>body{font-family:Arial;margin:20px;}.comp{max-width:600px;margin:0 auto;padding:30px;border:2px solid #667eea;}</style>\n"
    html = html + "</head>\n<body>\n"
    html = html + "<div class=\"comp\">\n"
    html = html + "<h1>COMPROBANTE</h1>\n"
    html = html + "<p>ID: " + tx.id + "</p>\n"
    html = html + "<p>País: " + tx.pais + "</p>\n"
    html = html + "<p>Importe: " + tx.moneda + " " + tx.importe + "</p>\n"
    html = html + "<p>IVA: " + tx.moneda + " " + tx.iva + "</p>\n"
    html = html + "<p><strong>TOTAL: " + tx.moneda + " " + tx.total + "</strong></p>\n"
    html = html + "<p><a href=\"/\">Volver</a></p>\n"
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// Demo automática
func handleDemo(pathVars, method, body) {
    procesarTransaccion("ventas", "MX", "100000")
    procesarTransaccion("compras", "COL", "50000")
    procesarTransaccion("ventas", "AR", "75000")
    procesarTransaccion("compras", "PE", "60000")
    
    let html = "<!DOCTYPE html>\n<html>\n<head>\n"
    html = html + "<title>Demo</title>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "</head>\n<body>\n"
    html = html + "<h1>Demo Completada</h1>\n"
    html = html + "<p>Se procesaron 4 transacciones</p>\n"
    html = html + "<p>Total: " + std.len(transacciones) + " transacciones</p>\n"
    html = html + "<p><a href=\"/\">Volver</a></p>\n"
    html = html + "</body>\n</html>"
    
    return html
}

// API
func handleAPI(pathVars, method, body) {
    return JSON({
        total: std.len(transacciones),
        transacciones: transacciones
    })
}

// Main
func main() {
    // Ejecutar demo en consola
    ejecutarDemo()
    
    console.log("\n✅ Configurando servidor web...")
    
    // Rutas
    http.handler("GET", "/", handleHome)
    http.handler("POST", "/procesar", handleProcesar)
    http.handler("GET", "/demo", handleDemo)
    http.handler("GET", "/api/transacciones", handleAPI)
    
    console.log("🌐 http://localhost:8080")
    console.log("")
    
    // Servidor
    http.serve(":8080")
}

main()