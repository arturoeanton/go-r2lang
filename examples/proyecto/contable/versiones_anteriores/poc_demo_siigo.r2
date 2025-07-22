// POC Demo Siigo - Sistema Contable LATAM
// Versión simplificada sin template literals

console.log("🌍 Sistema Contable LATAM - Demo Siigo")
console.log("📄 CFDI + 📊 DSL Reportes")
console.log("")

// Base de datos
let regiones = {
    "MX": { nombre: "México", moneda: "MXN", simbolo: "$", iva: 0.16, normativa: "SAT-CFDI" },
    "COL": { nombre: "Colombia", moneda: "COP", simbolo: "$", iva: 0.19, normativa: "DIAN" },
    "AR": { nombre: "Argentina", moneda: "ARS", simbolo: "$", iva: 0.21, normativa: "AFIP" },
    "CH": { nombre: "Chile", moneda: "CLP", simbolo: "$", iva: 0.19, normativa: "SII" },
    "UY": { nombre: "Uruguay", moneda: "UYU", simbolo: "$", iva: 0.22, normativa: "DGI" },
    "EC": { nombre: "Ecuador", moneda: "USD", simbolo: "$", iva: 0.12, normativa: "SRI" },
    "PE": { nombre: "Perú", moneda: "PEN", simbolo: "S/", iva: 0.18, normativa: "SUNAT" }
}

let transacciones = []
let cfdis = []

// DSL Reportes
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
            if ((region == "ALL" || tx.region == region) && (tipo == "todo" || tx.tipo == tipo)) {
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

// Funciones auxiliares
func procesarTransaccion(tipo, region, importe) {
    let config = regiones[region]
    let importeNum = std.parseFloat(importe)
    let iva = math.round((importeNum * config.iva) * 100) / 100
    let total = importeNum + iva
    
    let tx = {
        id: region + "-" + std.now() + "-" + math.randomInt(9999),
        tipo: tipo,
        region: region,
        pais: config.nombre,
        importe: importeNum,
        iva: iva,
        total: total,
        moneda: config.moneda,
        normativa: config.normativa,
        fecha: std.now()
    }
    
    transacciones[std.len(transacciones)] = tx
    return tx
}

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
                params[parts[0]] = std.replace(parts[1], "+", " ")
            }
        }
        i = i + 1
    }
    return params
}

// Demo en consola
func demostrarSistema() {
    console.log("\n=== DEMO SISTEMA CONTABLE LATAM ===\n")
    
    // Procesar algunas transacciones
    console.log("1. Procesando transacciones de ejemplo...")
    
    let tx1 = procesarTransaccion("ventas", "MX", "100000")
    console.log("   ✓ Venta México: $" + tx1.total + " " + tx1.moneda)
    
    let tx2 = procesarTransaccion("compras", "COL", "50000")
    console.log("   ✓ Compra Colombia: $" + tx2.total + " " + tx2.moneda)
    
    let tx3 = procesarTransaccion("ventas", "AR", "75000")
    console.log("   ✓ Venta Argentina: $" + tx3.total + " " + tx3.moneda)
    
    let tx4 = procesarTransaccion("compras", "PE", "60000")
    console.log("   ✓ Compra Perú: " + tx4.moneda + " " + tx4.total)
    
    console.log("\n2. Generando reportes con DSL...")
    
    let engine = ReportesContables
    
    let reporte1 = engine.use("reporte todo ALL")
    console.log("   ✓ Reporte general: " + reporte1.transacciones + " transacciones, Total: $" + reporte1.total)
    
    let reporte2 = engine.use("reporte ventas ALL")
    console.log("   ✓ Solo ventas: " + reporte2.transacciones + " transacciones, Total: $" + reporte2.total)
    
    let reporte3 = engine.use("reporte todo MX")
    console.log("   ✓ Todo México: " + reporte3.transacciones + " transacciones, Total: $" + reporte3.total)
    
    console.log("\n3. Procesando CFDI de ejemplo...")
    
    // Simular CFDI
    let cfdi = {
        uuid: "DEMO-123-456",
        emisor: "EMPRESA DEMO SA",
        receptor: "PUBLICO GENERAL",
        subtotal: 1000,
        iva: 160,
        total: 1160,
        fecha: std.now()
    }
    cfdis[std.len(cfdis)] = cfdi
    console.log("   ✓ CFDI procesado: UUID " + cfdi.uuid)
    
    console.log("\n=== RESUMEN ===")
    console.log("Transacciones totales: " + std.len(transacciones))
    console.log("CFDIs procesados: " + std.len(cfdis))
    console.log("Regiones activas: 7")
    
    console.log("\n🎯 VALUE PROPOSITION SIIGO:")
    console.log("   • Tiempo: 18 meses → 2 meses")
    console.log("   • Costo: $500K → $150K")
    console.log("   • Arquitectura: 7 sistemas → 1 DSL")
    console.log("   • ROI: 1,020% en 3 años")
}

// Handler página principal
func handleHome(pathVars, method, body) {
    let html = "<!DOCTYPE html><html><head>"
    html = html + "<title>Sistema Contable LATAM - Demo Siigo</title>"
    html = html + "<meta charset='UTF-8'>"
    html = html + "<style>"
    html = html + "body{font-family:Arial,sans-serif;margin:20px;background:#f5f5f5;}"
    html = html + ".container{max-width:800px;margin:0 auto;background:white;padding:30px;border-radius:10px;}"
    html = html + "h1{color:#2c3e50;text-align:center;}"
    html = html + ".card{background:#f8f9fa;padding:20px;border-radius:8px;margin:15px 0;}"
    html = html + "input,select{width:100%;padding:8px;margin:5px 0;}"
    html = html + "button{background:#667eea;color:white;padding:10px;border:none;border-radius:4px;width:100%;cursor:pointer;}"
    html = html + ".stats{background:#e8f5e9;padding:15px;border-radius:8px;text-align:center;}"
    html = html + ".value{background:#84fab0;padding:20px;border-radius:8px;margin:20px 0;}"
    html = html + "</style></head><body>"
    
    html = html + "<div class='container'>"
    html = html + "<h1>🌍 Sistema Contable LATAM</h1>"
    html = html + "<p style='text-align:center'>Demo para Siigo - R2Lang DSL</p>"
    
    html = html + "<div class='stats'>"
    html = html + "<h3>Estadísticas</h3>"
    html = html + "<p>Transacciones: <strong>" + std.len(transacciones) + "</strong> | "
    html = html + "CFDIs: <strong>" + std.len(cfdis) + "</strong> | "
    html = html + "Regiones: <strong>7</strong></p>"
    html = html + "</div>"
    
    // Formulario de transacción
    html = html + "<div class='card'>"
    html = html + "<h3>💰 Procesar Transacción</h3>"
    html = html + "<form action='/procesar' method='POST'>"
    html = html + "<select name='tipo'>"
    html = html + "<option value='ventas'>Venta</option>"
    html = html + "<option value='compras'>Compra</option>"
    html = html + "</select>"
    html = html + "<select name='region'>"
    html = html + "<option value='MX'>México (16%)</option>"
    html = html + "<option value='COL'>Colombia (19%)</option>"
    html = html + "<option value='AR'>Argentina (21%)</option>"
    html = html + "<option value='CH'>Chile (19%)</option>"
    html = html + "<option value='UY'>Uruguay (22%)</option>"
    html = html + "<option value='EC'>Ecuador (12%)</option>"
    html = html + "<option value='PE'>Perú (18%)</option>"
    html = html + "</select>"
    html = html + "<input type='number' name='importe' value='100000' required>"
    html = html + "<button type='submit'>Procesar</button>"
    html = html + "</form>"
    html = html + "</div>"
    
    // Generador de reportes
    html = html + "<div class='card'>"
    html = html + "<h3>📊 Generar Reporte DSL</h3>"
    html = html + "<form action='/reporte' method='POST'>"
    html = html + "<input type='text' name='comando' value='reporte todo ALL' style='font-family:monospace;'>"
    html = html + "<p>Ejemplos: reporte ventas MX | reporte todo COL</p>"
    html = html + "<button type='submit'>Generar</button>"
    html = html + "</form>"
    html = html + "</div>"
    
    // Links
    html = html + "<div style='text-align:center;margin:20px 0;'>"
    html = html + "<a href='/demo'>Demo Automática</a> | "
    html = html + "<a href='/api/transacciones'>API Transacciones</a>"
    html = html + "</div>"
    
    // Value proposition
    html = html + "<div class='value'>"
    html = html + "<h3>🎯 Value Proposition Siigo</h3>"
    html = html + "<p><strong>18 meses → 2 meses</strong> | "
    html = html + "<strong>$500K → $150K</strong> | "
    html = html + "<strong>7 sistemas → 1 DSL</strong> | "
    html = html + "<strong>ROI: 1,020%</strong></p>"
    html = html + "</div>"
    
    html = html + "</div></body></html>"
    
    return html
}

// Handler procesar transacción
func handleProcesar(pathVars, method, body) {
    let params = parseFormData(body)
    let tx = procesarTransaccion(
        params.tipo || "ventas",
        params.region || "COL",
        params.importe || "100000"
    )
    
    let html = "<!DOCTYPE html><html><head>"
    html = html + "<title>Comprobante</title>"
    html = html + "<meta charset='UTF-8'>"
    html = html + "<style>"
    html = html + "body{font-family:Arial,sans-serif;margin:20px;}"
    html = html + ".comprobante{max-width:600px;margin:0 auto;background:white;padding:30px;border:2px solid #667eea;border-radius:10px;}"
    html = html + ".row{display:flex;justify-content:space-between;margin:10px 0;}"
    html = html + ".total{background:#f5f5f5;padding:15px;border-radius:5px;margin:20px 0;}"
    html = html + "</style></head><body>"
    
    html = html + "<div class='comprobante'>"
    html = html + "<h1 style='text-align:center;'>COMPROBANTE</h1>"
    html = html + "<h2 style='text-align:center;'>" + tx.pais + "</h2>"
    
    html = html + "<div class='row'><strong>ID:</strong><span>" + tx.id + "</span></div>"
    html = html + "<div class='row'><strong>Tipo:</strong><span>" + tx.tipo + "</span></div>"
    html = html + "<div class='row'><strong>Fecha:</strong><span>" + tx.fecha + "</span></div>"
    
    html = html + "<div class='total'>"
    html = html + "<div class='row'><strong>Importe:</strong><span>" + tx.moneda + " " + tx.importe + "</span></div>"
    html = html + "<div class='row'><strong>IVA:</strong><span>" + tx.moneda + " " + tx.iva + "</span></div>"
    html = html + "<div class='row' style='border-top:2px solid #333;padding-top:10px;'>"
    html = html + "<strong>TOTAL:</strong><strong>" + tx.moneda + " " + tx.total + "</strong></div>"
    html = html + "</div>"
    
    html = html + "<p style='text-align:center;color:#27ae60;'><strong>✅ TRANSACCIÓN PROCESADA</strong></p>"
    html = html + "<p style='text-align:center;'><a href='/'>← Volver</a></p>"
    html = html + "</div></body></html>"
    
    return html
}

// Handler reporte
func handleReporte(pathVars, method, body) {
    let params = parseFormData(body)
    let comando = params.comando || "reporte todo ALL"
    
    let engine = ReportesContables
    let resultado = engine.use(comando)
    
    let html = "<!DOCTYPE html><html><head>"
    html = html + "<title>Reporte</title>"
    html = html + "<meta charset='UTF-8'>"
    html = html + "<style>"
    html = html + "body{font-family:Arial,sans-serif;margin:20px;}"
    html = html + ".container{max-width:800px;margin:0 auto;background:white;padding:30px;}"
    html = html + ".resultado{background:#f0f2f5;padding:20px;border-radius:8px;}"
    html = html + "</style></head><body>"
    
    html = html + "<div class='container'>"
    html = html + "<h1>📊 Reporte Generado</h1>"
    html = html + "<p><strong>Comando DSL:</strong> <code>" + comando + "</code></p>"
    
    html = html + "<div class='resultado'>"
    html = html + "<h3>Resumen:</h3>"
    html = html + "<p><strong>Tipo:</strong> " + resultado.tipo + "</p>"
    html = html + "<p><strong>Región:</strong> " + resultado.region + "</p>"
    html = html + "<p><strong>Transacciones:</strong> " + resultado.transacciones + "</p>"
    html = html + "<p><strong>Total:</strong> $" + resultado.total + "</p>"
    html = html + "</div>"
    
    html = html + "<p style='text-align:center;margin-top:30px;'><a href='/'>← Volver</a></p>"
    html = html + "</div></body></html>"
    
    return html
}

// Handler demo
func handleDemo(pathVars, method, body) {
    // Generar transacciones de demo
    procesarTransaccion("ventas", "MX", "100000")
    procesarTransaccion("compras", "COL", "50000")
    procesarTransaccion("ventas", "AR", "75000")
    procesarTransaccion("compras", "PE", "60000")
    
    let html = "<!DOCTYPE html><html><head>"
    html = html + "<title>Demo Automática</title>"
    html = html + "<meta charset='UTF-8'>"
    html = html + "<style>"
    html = html + "body{font-family:Arial,sans-serif;margin:20px;background:#f5f5f5;}"
    html = html + ".container{max-width:600px;margin:0 auto;background:white;padding:30px;border-radius:10px;}"
    html = html + ".success{background:#d4edda;color:#155724;padding:20px;border-radius:8px;}"
    html = html + "</style></head><body>"
    
    html = html + "<div class='container'>"
    html = html + "<h1>🎪 Demo Automática</h1>"
    html = html + "<div class='success'>"
    html = html + "<h3>✅ Completado</h3>"
    html = html + "<p>Se procesaron 4 transacciones de diferentes regiones</p>"
    html = html + "</div>"
    html = html + "<p>Total de transacciones: <strong>" + std.len(transacciones) + "</strong></p>"
    html = html + "<p style='text-align:center;margin-top:30px;'>"
    html = html + "<a href='/'>← Volver</a> | "
    html = html + "<a href='/api/transacciones'>Ver todas</a>"
    html = html + "</p>"
    html = html + "</div></body></html>"
    
    return html
}

// Handler API transacciones
func handleAPITransacciones(pathVars, method, body) {
    return JSON({
        total: std.len(transacciones),
        transacciones: transacciones
    })
}

// Handler API regiones  
func handleAPIRegiones(pathVars, method, body) {
    return JSON(regiones)
}

// Función principal
func main() {
    console.log("✅ Sistema Contable LATAM - POC Siigo")
    console.log("")
    
    // Primero ejecutar demo en consola
    demostrarSistema()
    
    console.log("\n✅ Configurando servidor web...")
    
    // Configurar rutas
    http.handler("GET", "/", handleHome)
    http.handler("POST", "/procesar", handleProcesar)
    http.handler("POST", "/reporte", handleReporte)
    http.handler("GET", "/demo", handleDemo)
    http.handler("GET", "/api/transacciones", handleAPITransacciones)
    http.handler("GET", "/api/regiones", handleAPIRegiones)
    
    console.log("🌐 Servidor listo en http://localhost:8080")
    console.log("🎪 Demo auto en http://localhost:8080/demo")
    console.log("")
    
    // Iniciar servidor
    http.serve(":8080")
}

// Ejecutar
main()