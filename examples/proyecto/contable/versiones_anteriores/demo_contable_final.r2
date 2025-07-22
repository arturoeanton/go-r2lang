// POC Sistema Contable LATAM - Versión Final Completa
// Con Libro Diario (Debe/Haber), DSL Reportes y API funcional

console.log("🌍 Sistema Contable LATAM - POC Completo")
console.log("📊 Libro Diario + DSL Reportes + API")
console.log("")

// Configuración regiones con plan de cuentas
let regiones = {
    "MX": { 
        nombre: "México", moneda: "MXN", iva: 0.16,
        cuentas: {
            caja: "1101", clientes: "1201", proveedores: "2101",
            ventas: "4101", compras: "5101", ivaDebito: "2401", ivaCredito: "1401"
        }
    },
    "COL": { 
        nombre: "Colombia", moneda: "COP", iva: 0.19,
        cuentas: {
            caja: "110501", clientes: "130501", proveedores: "220501", 
            ventas: "413501", compras: "620501", ivaDebito: "240801", ivaCredito: "240802"
        }
    },
    "AR": { 
        nombre: "Argentina", moneda: "ARS", iva: 0.21,
        cuentas: {
            caja: "1.1.1.01", clientes: "1.1.2.01", proveedores: "2.1.1.01",
            ventas: "4.1.1.01", compras: "5.1.1.01", ivaDebito: "2.1.3.01", ivaCredito: "1.1.3.01"
        }
    },
    "CH": { 
        nombre: "Chile", moneda: "CLP", iva: 0.19,
        cuentas: {
            caja: "11010", clientes: "11030", proveedores: "21010",
            ventas: "31010", compras: "41010", ivaDebito: "21070", ivaCredito: "11070"
        }
    },
    "UY": { 
        nombre: "Uruguay", moneda: "UYU", iva: 0.22,
        cuentas: {
            caja: "1111", clientes: "1121", proveedores: "2111",
            ventas: "4111", compras: "5111", ivaDebito: "2141", ivaCredito: "1141"
        }
    },
    "EC": { 
        nombre: "Ecuador", moneda: "USD", iva: 0.12,
        cuentas: {
            caja: "101.01", clientes: "102.01", proveedores: "201.01",
            ventas: "401.01", compras: "501.01", ivaDebito: "201.04", ivaCredito: "101.04"
        }
    },
    "PE": { 
        nombre: "Perú", moneda: "PEN", iva: 0.18,
        cuentas: {
            caja: "101", clientes: "121", proveedores: "421",
            ventas: "701", compras: "601", ivaDebito: "401", ivaCredito: "401"
        }
    }
}

let transacciones = []
let asientosContables = []

// DSL Reportes Financieros
dsl ReportesFinancieros {
    token("REPORTE", "reporte|report")
    token("TIPO", "diario|mayor|balance|ventas|compras|iva")
    token("REGION", "MX|COL|AR|CH|UY|EC|PE|ALL")
    
    rule("generar_reporte", ["REPORTE", "TIPO"], "reporteSimple")
    rule("generar_reporte_region", ["REPORTE", "TIPO", "REGION"], "reporteConRegion")
    
    func reporteSimple(cmd, tipo) {
        if (tipo == "diario") {
            return generarLibroDiario("ALL")
        }
        if (tipo == "balance") {
            return generarBalance("ALL")
        }
        if (tipo == "ventas") {
            return generarReporteTipo("ventas", "ALL")
        }
        if (tipo == "compras") {
            return generarReporteTipo("compras", "ALL")
        }
        if (tipo == "iva") {
            return generarReporteIVA("ALL")
        }
        return { error: "Tipo de reporte no válido" }
    }
    
    func reporteConRegion(cmd, tipo, region) {
        if (tipo == "diario") {
            return generarLibroDiario(region)
        }
        if (tipo == "balance") {
            return generarBalance(region)
        }
        if (tipo == "ventas") {
            return generarReporteTipo("ventas", region)
        }
        if (tipo == "compras") {
            return generarReporteTipo("compras", region)
        }
        if (tipo == "iva") {
            return generarReporteIVA(region)
        }
        return { error: "Tipo de reporte no válido" }
    }
    
    func generarLibroDiario(region) {
        let asientos = []
        let i = 0
        while (i < std.len(asientosContables)) {
            let asiento = asientosContables[i]
            if (region == "ALL" || asiento.region == region) {
                asientos[std.len(asientos)] = asiento
            }
            i = i + 1
        }
        
        return {
            tipo: "Libro Diario",
            region: region,
            totalAsientos: std.len(asientos),
            asientos: asientos
        }
    }
    
    func generarBalance(region) {
        let totalDebe = 0
        let totalHaber = 0
        let i = 0
        
        while (i < std.len(asientosContables)) {
            let asiento = asientosContables[i]
            if (region == "ALL" || asiento.region == region) {
                let j = 0
                while (j < std.len(asiento.movimientos)) {
                    let mov = asiento.movimientos[j]
                    if (mov.tipo == "DEBE") {
                        totalDebe = totalDebe + mov.monto
                    } else {
                        totalHaber = totalHaber + mov.monto
                    }
                    j = j + 1
                }
            }
            i = i + 1
        }
        
        return {
            tipo: "Balance",
            region: region,
            totalDebe: math.round(totalDebe * 100) / 100,
            totalHaber: math.round(totalHaber * 100) / 100,
            cuadrado: math.round(totalDebe * 100) / 100 == math.round(totalHaber * 100) / 100
        }
    }
    
    func generarReporteTipo(tipo, region) {
        let total = 0
        let count = 0
        let i = 0
        
        while (i < std.len(transacciones)) {
            let tx = transacciones[i]
            if ((region == "ALL" || tx.region == region) && tx.tipo == tipo) {
                total = total + tx.total
                count = count + 1
            }
            i = i + 1
        }
        
        return {
            tipo: "Reporte " + tipo,
            region: region,
            transacciones: count,
            total: math.round(total * 100) / 100
        }
    }
    
    func generarReporteIVA(region) {
        let ivaDebito = 0
        let ivaCredito = 0
        let i = 0
        
        while (i < std.len(transacciones)) {
            let tx = transacciones[i]
            if (region == "ALL" || tx.region == region) {
                if (tx.tipo == "ventas") {
                    ivaDebito = ivaDebito + tx.iva
                } else {
                    ivaCredito = ivaCredito + tx.iva
                }
            }
            i = i + 1
        }
        
        return {
            tipo: "Reporte IVA",
            region: region,
            ivaDebito: math.round(ivaDebito * 100) / 100,
            ivaCredito: math.round(ivaCredito * 100) / 100,
            saldo: math.round((ivaDebito - ivaCredito) * 100) / 100
        }
    }
}

// Función procesar transacción con asiento contable
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
    
    // Crear asiento contable
    let asiento = {
        id: "AS-" + tx.id,
        fecha: tx.fecha,
        region: region,
        descripcion: std.toUpperCase(tipo) + " - " + config.nombre,
        movimientos: []
    }
    
    if (tipo == "ventas") {
        // Venta: Debe Clientes, Haber Ventas e IVA
        asiento.movimientos[0] = {
            cuenta: config.cuentas.clientes,
            descripcion: "Clientes",
            tipo: "DEBE",
            monto: total
        }
        asiento.movimientos[1] = {
            cuenta: config.cuentas.ventas,
            descripcion: "Ventas",
            tipo: "HABER",
            monto: importeNum
        }
        asiento.movimientos[2] = {
            cuenta: config.cuentas.ivaDebito,
            descripcion: "IVA Débito",
            tipo: "HABER",
            monto: iva
        }
    } else {
        // Compra: Debe Compras e IVA, Haber Proveedores
        asiento.movimientos[0] = {
            cuenta: config.cuentas.compras,
            descripcion: "Compras",
            tipo: "DEBE",
            monto: importeNum
        }
        asiento.movimientos[1] = {
            cuenta: config.cuentas.ivaCredito,
            descripcion: "IVA Crédito",
            tipo: "DEBE",
            monto: iva
        }
        asiento.movimientos[2] = {
            cuenta: config.cuentas.proveedores,
            descripcion: "Proveedores",
            tipo: "HABER",
            monto: total
        }
    }
    
    asientosContables[std.len(asientosContables)] = asiento
    
    return tx
}

// Demo en consola
func demoConsola() {
    console.log("=== DEMO SISTEMA CONTABLE ===")
    
    let tx1 = procesarTransaccion("ventas", "MX", "100000")
    console.log("✓ Venta México: $" + tx1.total + " " + tx1.moneda)
    
    let tx2 = procesarTransaccion("compras", "COL", "50000")
    console.log("✓ Compra Colombia: $" + tx2.total + " " + tx2.moneda)
    
    let tx3 = procesarTransaccion("ventas", "AR", "75000")
    console.log("✓ Venta Argentina: $" + tx3.total + " " + tx3.moneda)
    
    // Mostrar ejemplo de asiento contable
    if (std.len(asientosContables) > 0) {
        console.log("\n📚 LIBRO DIARIO - Ejemplo México:")
        let asiento = asientosContables[0]
        console.log("Asiento: " + asiento.id)
        let i = 0
        while (i < std.len(asiento.movimientos)) {
            let mov = asiento.movimientos[i]
            console.log("  " + mov.tipo + " " + mov.cuenta + " " + mov.descripcion + ": $" + mov.monto)
            i = i + 1
        }
    }
    
    // Ejemplo DSL
    console.log("\n📊 EJEMPLOS DSL:")
    let engine = ReportesFinancieros
    let balance = engine.use("reporte balance")
    console.log("Balance: Debe=$" + balance.totalDebe + " Haber=$" + balance.totalHaber)
    
    console.log("\n🎯 VALUE PROPOSITION:")
    console.log("• 18 meses → 2 meses")
    console.log("• $500K → $150K")
    console.log("• ROI: 1,020%")
}

// Función para obtener parámetro
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

// Handler página principal
func handleHome(pathVars, method, body) {
    let html = "<!DOCTYPE html>\n<html>\n<head>\n"
    html = html + "<title>Sistema Contable LATAM - POC Completo</title>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<style>\n"
    html = html + "body{font-family:Arial;margin:20px;background:#f5f5f5;}\n"
    html = html + ".container{max-width:1200px;margin:0 auto;background:white;padding:30px;border-radius:10px;}\n"
    html = html + "h1{color:#2c3e50;text-align:center;}\n"
    html = html + ".grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(350px,1fr));gap:20px;}\n"
    html = html + ".card{background:#f8f9fa;padding:20px;border-radius:8px;border:1px solid #ddd;}\n"
    html = html + "input,select,textarea{width:100%;padding:8px;margin:5px 0;}\n"
    html = html + "button{background:#667eea;color:white;padding:10px;border:none;border-radius:4px;width:100%;cursor:pointer;}\n"
    html = html + ".value{background:#84fab0;padding:20px;border-radius:8px;text-align:center;}\n"
    html = html + ".asiento{background:#fff3cd;padding:15px;border-radius:5px;margin:10px 0;font-family:monospace;}\n"
    html = html + "</style>\n</head>\n<body>\n"
    
    html = html + "<div class=\"container\">\n"
    html = html + "<h1>🌍 Sistema Contable LATAM</h1>\n"
    html = html + "<p style=\"text-align:center\">POC Completo para Siigo - " + std.len(transacciones) + " transacciones - " + std.len(asientosContables) + " asientos</p>\n"
    
    html = html + "<div class=\"grid\">\n"
    
    // Formulario transacción
    html = html + "<div class=\"card\">\n"
    html = html + "<h3>Procesar Transacción</h3>\n"
    html = html + "<form action=\"/procesar\" method=\"POST\">\n"
    html = html + "<select name=\"tipo\">\n"
    html = html + "<option value=\"ventas\">Venta</option>\n"
    html = html + "<option value=\"compras\">Compra</option>\n"
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
    html = html + "<input type=\"number\" name=\"importe\" value=\"100000\" required>\n"
    html = html + "<button type=\"submit\">Procesar</button>\n"
    html = html + "</form>\n"
    html = html + "</div>\n"
    
    // DSL Query
    html = html + "<div class=\"card\">\n"
    html = html + "<h3>Consulta DSL Reportes</h3>\n"
    html = html + "<form action=\"/dsl\" method=\"POST\">\n"
    html = html + "<textarea name=\"query\" rows=\"3\" placeholder=\"Ejemplos:\nreporte diario\nreporte balance MX\nreporte iva ALL\">reporte balance</textarea>\n"
    html = html + "<button type=\"submit\">Ejecutar Query</button>\n"
    html = html + "</form>\n"
    html = html + "</div>\n"
    
    html = html + "</div>\n"
    
    // Links
    html = html + "<p style=\"text-align:center;margin-top:20px;\">\n"
    html = html + "<a href=\"/demo\">Demo Auto</a> | "
    html = html + "<a href=\"/libro\">Libro Diario</a> | "
    html = html + "<a href=\"/api/transacciones\">API Transacciones</a> | "
    html = html + "<a href=\"/api/asientos\">API Asientos</a>\n"
    html = html + "</p>\n"
    
    // Value proposition
    html = html + "<div class=\"value\">\n"
    html = html + "<h3>🎯 Value Proposition Siigo</h3>\n"
    html = html + "<p><strong>18 meses → 2 meses</strong> | <strong>$500K → $150K</strong> | <strong>7 sistemas → 1 DSL</strong> | <strong>ROI: 1,020%</strong></p>\n"
    html = html + "</div>\n"
    
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// Handler procesar
func handleProcesar(pathVars, method, body) {
    let tipo = getParam(body, "tipo")
    let region = getParam(body, "region")
    let importe = getParam(body, "importe")
    
    if (tipo == "") { tipo = "ventas" }
    if (region == "") { region = "COL" }
    if (importe == "") { importe = "100000" }
    
    let tx = procesarTransaccion(tipo, region, importe)
    
    // Obtener último asiento
    let asiento = asientosContables[std.len(asientosContables) - 1]
    
    let html = "<!DOCTYPE html>\n<html>\n<head>\n"
    html = html + "<title>Comprobante y Asiento</title>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<style>\n"
    html = html + "body{font-family:Arial;margin:20px;}\n"
    html = html + ".comp{max-width:800px;margin:0 auto;padding:30px;border:2px solid #667eea;border-radius:10px;}\n"
    html = html + ".asiento{background:#f8f9fa;padding:20px;margin:20px 0;border-radius:8px;}\n"
    html = html + ".debe{color:#28a745;font-weight:bold;}\n"
    html = html + ".haber{color:#dc3545;font-weight:bold;}\n"
    html = html + "</style>\n"
    html = html + "</head>\n<body>\n"
    
    html = html + "<div class=\"comp\">\n"
    html = html + "<h1>COMPROBANTE</h1>\n"
    html = html + "<h2>" + tx.pais + "</h2>\n"
    html = html + "<p><strong>ID:</strong> " + tx.id + "</p>\n"
    html = html + "<p><strong>Tipo:</strong> " + std.toUpperCase(tx.tipo) + "</p>\n"
    html = html + "<p><strong>Importe:</strong> " + tx.moneda + " " + tx.importe + "</p>\n"
    html = html + "<p><strong>IVA:</strong> " + tx.moneda + " " + tx.iva + "</p>\n"
    html = html + "<hr>\n"
    html = html + "<p><strong>TOTAL:</strong> " + tx.moneda + " " + tx.total + "</p>\n"
    
    // Mostrar asiento contable
    html = html + "<div class=\"asiento\">\n"
    html = html + "<h3>📚 ASIENTO CONTABLE</h3>\n"
    html = html + "<p><strong>Asiento:</strong> " + asiento.id + "</p>\n"
    html = html + "<p><strong>Descripción:</strong> " + asiento.descripcion + "</p>\n"
    html = html + "<table width=\"100%\">\n"
    html = html + "<tr><th>Cuenta</th><th>Descripción</th><th>Debe</th><th>Haber</th></tr>\n"
    
    let i = 0
    while (i < std.len(asiento.movimientos)) {
        let mov = asiento.movimientos[i]
        html = html + "<tr>\n"
        html = html + "<td>" + mov.cuenta + "</td>\n"
        html = html + "<td>" + mov.descripcion + "</td>\n"
        if (mov.tipo == "DEBE") {
            html = html + "<td class=\"debe\">" + mov.monto + "</td>\n"
            html = html + "<td></td>\n"
        } else {
            html = html + "<td></td>\n"
            html = html + "<td class=\"haber\">" + mov.monto + "</td>\n"
        }
        html = html + "</tr>\n"
        i = i + 1
    }
    
    html = html + "</table>\n"
    html = html + "</div>\n"
    
    html = html + "<p style=\"text-align:center;color:green;\">✅ PROCESADO CON ÉXITO</p>\n"
    html = html + "<p style=\"text-align:center;\"><a href=\"/\">Volver</a></p>\n"
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// Handler libro diario
func handleLibro(pathVars, method, body) {
    let html = "<!DOCTYPE html>\n<html>\n<head>\n"
    html = html + "<title>Libro Diario</title>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<style>\n"
    html = html + "body{font-family:Arial;margin:20px;}\n"
    html = html + ".container{max-width:1000px;margin:0 auto;}\n"
    html = html + ".asiento{background:#f8f9fa;padding:15px;margin:15px 0;border-radius:8px;}\n"
    html = html + "table{width:100%;border-collapse:collapse;margin:10px 0;}\n"
    html = html + "th,td{padding:8px;text-align:left;border-bottom:1px solid #ddd;}\n"
    html = html + ".debe{color:#28a745;font-weight:bold;}\n"
    html = html + ".haber{color:#dc3545;font-weight:bold;}\n"
    html = html + "</style>\n"
    html = html + "</head>\n<body>\n"
    
    html = html + "<div class=\"container\">\n"
    html = html + "<h1>📚 LIBRO DIARIO</h1>\n"
    html = html + "<p>Total asientos: " + std.len(asientosContables) + "</p>\n"
    
    let i = 0
    while (i < std.len(asientosContables)) {
        let asiento = asientosContables[i]
        html = html + "<div class=\"asiento\">\n"
        html = html + "<h3>Asiento: " + asiento.id + "</h3>\n"
        html = html + "<p><strong>Fecha:</strong> " + asiento.fecha + " | <strong>Región:</strong> " + asiento.region + "</p>\n"
        html = html + "<p><strong>Descripción:</strong> " + asiento.descripcion + "</p>\n"
        
        html = html + "<table>\n"
        html = html + "<tr><th>Cuenta</th><th>Descripción</th><th>Debe</th><th>Haber</th></tr>\n"
        
        let j = 0
        while (j < std.len(asiento.movimientos)) {
            let mov = asiento.movimientos[j]
            html = html + "<tr>\n"
            html = html + "<td>" + mov.cuenta + "</td>\n"
            html = html + "<td>" + mov.descripcion + "</td>\n"
            if (mov.tipo == "DEBE") {
                html = html + "<td class=\"debe\">" + mov.monto + "</td>\n"
                html = html + "<td></td>\n"
            } else {
                html = html + "<td></td>\n"
                html = html + "<td class=\"haber\">" + mov.monto + "</td>\n"
            }
            html = html + "</tr>\n"
            j = j + 1
        }
        
        html = html + "</table>\n"
        html = html + "</div>\n"
        i = i + 1
    }
    
    html = html + "<p style=\"text-align:center;\"><a href=\"/\">Volver</a></p>\n"
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// Handler DSL
func handleDSL(pathVars, method, body) {
    let query = getParam(body, "query")
    if (query == "") {
        query = "reporte balance"
    }
    
    let engine = ReportesFinancieros
    let resultado = engine.use(query)
    
    let html = "<!DOCTYPE html>\n<html>\n<head>\n"
    html = html + "<title>Resultado DSL</title>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<style>\n"
    html = html + "body{font-family:Arial;margin:20px;}\n"
    html = html + ".container{max-width:800px;margin:0 auto;padding:30px;background:#f8f9fa;border-radius:10px;}\n"
    html = html + "pre{background:white;padding:20px;border-radius:5px;overflow:auto;}\n"
    html = html + "</style>\n"
    html = html + "</head>\n<body>\n"
    
    html = html + "<div class=\"container\">\n"
    html = html + "<h1>📊 Resultado DSL</h1>\n"
    html = html + "<p><strong>Query:</strong> " + query + "</p>\n"
    html = html + "<pre>" + JSON.stringify(resultado) + "</pre>\n"
    html = html + "<p style=\"text-align:center;\"><a href=\"/\">Volver</a></p>\n"
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// Handler demo
func handleDemo(pathVars, method, body) {
    procesarTransaccion("ventas", "MX", "100000")
    procesarTransaccion("compras", "COL", "50000")
    procesarTransaccion("ventas", "AR", "75000")
    procesarTransaccion("compras", "PE", "60000")
    procesarTransaccion("ventas", "CH", "120000")
    procesarTransaccion("compras", "EC", "30000")
    
    let html = "<!DOCTYPE html>\n<html>\n<head>\n"
    html = html + "<title>Demo Automática</title>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<style>body{font-family:Arial;margin:20px;}.container{max-width:600px;margin:0 auto;padding:30px;background:white;border-radius:10px;}</style>\n"
    html = html + "</head>\n<body>\n"
    html = html + "<div class=\"container\">\n"
    html = html + "<h1>Demo Automática</h1>\n"
    html = html + "<p>✅ Se procesaron 6 transacciones de diferentes regiones</p>\n"
    html = html + "<p>Total acumulado: " + std.len(transacciones) + " transacciones</p>\n"
    html = html + "<p>Total asientos: " + std.len(asientosContables) + " asientos contables</p>\n"
    
    // Mostrar balance
    let engine = ReportesFinancieros
    let balance = engine.use("reporte balance")
    html = html + "<p><strong>Balance General:</strong></p>\n"
    html = html + "<p>Debe: $" + balance.totalDebe + " | Haber: $" + balance.totalHaber + "</p>\n"
    html = html + "<p>Cuadrado: " + balance.cuadrado + "</p>\n"
    
    html = html + "<p style=\"text-align:center;margin-top:30px;\">\n"
    html = html + "<a href=\"/\">Volver</a> | <a href=\"/libro\">Ver Libro Diario</a>\n"
    html = html + "</p>\n"
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// APIs
func handleAPITransacciones(pathVars, method, body) {
    return JSON({
        total: std.len(transacciones),
        regiones: 7,
        transacciones: transacciones
    })
}

func handleAPIAsientos(pathVars, method, body) {
    return JSON({
        total: std.len(asientosContables),
        asientos: asientosContables
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
    http.handler("GET", "/libro", handleLibro)
    http.handler("POST", "/dsl", handleDSL)
    http.handler("GET", "/demo", handleDemo)
    http.handler("GET", "/api/transacciones", handleAPITransacciones)
    http.handler("GET", "/api/asientos", handleAPIAsientos)
    
    console.log("🌐 http://localhost:8080")
    console.log("")
    
    // Servidor
    http.serve(":8080")
}

main()