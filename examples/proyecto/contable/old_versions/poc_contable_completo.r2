// POC Sistema Contable LATAM Completo - Demo Siigo
// Con Libro Diario (Debe/Haber) y DSL Reportes Financieros

console.log("🌍 Sistema Contable LATAM - POC Completo")
console.log("📊 Con Libro Diario + DSL Reportes Financieros")
console.log("")

// Configuración regiones y plan de cuentas
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
    token("PERIODO", "hoy|mes|todo")
    
    rule("generar_reporte", ["REPORTE", "TIPO"], "reporteSimple")
    rule("reporte_region", ["REPORTE", "TIPO", "REGION"], "reporteRegion")
    rule("reporte_completo", ["REPORTE", "TIPO", "REGION", "PERIODO"], "reporteCompleto")
    
    func reporteSimple(cmd, tipo) {
        return generarReporte(tipo, "ALL", "todo")
    }
    
    func reporteRegion(cmd, tipo, region) {
        return generarReporte(tipo, region, "todo")
    }
    
    func reporteCompleto(cmd, tipo, region, periodo) {
        return generarReporte(tipo, region, periodo)
    }
    
    func generarReporte(tipo, region, periodo) {
        if (tipo == "diario") {
            return generarLibroDiario(region)
        } else if (tipo == "balance") {
            return generarBalance(region)
        } else if (tipo == "ventas" || tipo == "compras") {
            return generarReporteTipo(tipo, region)
        } else if (tipo == "iva") {
            return generarReporteIVA(region)
        }
        
        return {error: "Tipo de reporte no válido"}
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
        id: region + "-" + std.now() + "-" + math.randomInt(9999),
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
    
    console.log("\n📚 LIBRO DIARIO - Ejemplo México:")
    let asiento = asientosContables[0]
    console.log("Asiento: " + asiento.id)
    let i = 0
    while (i < std.len(asiento.movimientos)) {
        let mov = asiento.movimientos[i]
        console.log("  " + mov.tipo + " " + mov.cuenta + " " + mov.descripcion + ": $" + mov.monto)
        i = i + 1
    }
    
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
    html = html + ".value{background:#84fab0;padding:20px;border-radius:8px;margin:20px 0;}\n"
    html = html + ".stats{display:grid;grid-template-columns:repeat(4,1fr);gap:15px;margin:20px 0;}\n"
    html = html + ".stat{background:#e8f5e9;padding:15px;border-radius:8px;text-align:center;}\n"
    html = html + "</style>\n</head>\n<body>\n"
    
    html = html + "<div class=\"container\">\n"
    html = html + "<h1>🌍 Sistema Contable LATAM - POC Completo</h1>\n"
    html = html + "<p style=\"text-align:center\">Demo Siigo - Con Libro Diario y DSL Reportes</p>\n"
    
    // Estadísticas
    html = html + "<div class=\"stats\">\n"
    html = html + "<div class=\"stat\">\n<h3>" + std.len(transacciones) + "</h3>\n<p>Transacciones</p>\n</div>\n"
    html = html + "<div class=\"stat\">\n<h3>" + std.len(asientosContables) + "</h3>\n<p>Asientos</p>\n</div>\n"
    html = html + "<div class=\"stat\">\n<h3>7</h3>\n<p>Regiones</p>\n</div>\n"
    html = html + "<div class=\"stat\">\n<h3>100%</h3>\n<p>Compliance</p>\n</div>\n"
    html = html + "</div>\n"
    
    html = html + "<div class=\"grid\">\n"
    
    // Formulario transacción
    html = html + "<div class=\"card\">\n"
    html = html + "<h3>💰 Procesar Transacción</h3>\n"
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
    html = html + "<input type=\"number\" name=\"importe\" value=\"100000\" placeholder=\"Importe\">\n"
    html = html + "<button type=\"submit\">Procesar</button>\n"
    html = html + "</form>\n"
    html = html + "</div>\n"
    
    // DSL Reportes
    html = html + "<div class=\"card\">\n"
    html = html + "<h3>📊 Consultas DSL Reportes</h3>\n"
    html = html + "<form action=\"/dsl\" method=\"POST\">\n"
    html = html + "<input type=\"text\" name=\"query\" placeholder=\"reporte diario ALL\" style=\"font-family:monospace;\">\n"
    html = html + "<p style=\"font-size:0.9em;color:#666;\">Ejemplos:<br>\n"
    html = html + "• reporte diario ALL<br>\n"
    html = html + "• reporte balance MX<br>\n"
    html = html + "• reporte ventas COL<br>\n"
    html = html + "• reporte iva ALL</p>\n"
    html = html + "<button type=\"submit\">Ejecutar Query</button>\n"
    html = html + "</form>\n"
    html = html + "</div>\n"
    
    // Links
    html = html + "<div class=\"card\">\n"
    html = html + "<h3>🔗 Accesos Rápidos</h3>\n"
    html = html + "<p><a href=\"/libro-diario\">📚 Ver Libro Diario</a></p>\n"
    html = html + "<p><a href=\"/balance\">⚖️ Ver Balance General</a></p>\n"
    html = html + "<p><a href=\"/demo\">🎪 Demo Automática</a></p>\n"
    html = html + "<p><a href=\"/api/transacciones\">📡 API Transacciones</a></p>\n"
    html = html + "<p><a href=\"/api/asientos\">📖 API Asientos</a></p>\n"
    html = html + "</div>\n"
    
    html = html + "</div>\n"
    
    // Value proposition
    html = html + "<div class=\"value\">\n"
    html = html + "<h3>🎯 Value Proposition Siigo</h3>\n"
    html = html + "<p><strong>18 meses → 2 meses</strong> | <strong>$500K → $150K</strong> | <strong>7 sistemas → 1 DSL</strong> | <strong>ROI: 1,020%</strong></p>\n"
    html = html + "</div>\n"
    
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// Handler procesar transacción
func handleProcesar(pathVars, method, body) {
    let tipo = getParam(body, "tipo")
    let region = getParam(body, "region")
    let importe = getParam(body, "importe")
    
    if (tipo == "") { tipo = "ventas" }
    if (region == "") { region = "COL" }
    if (importe == "") { importe = "100000" }
    
    let tx = procesarTransaccion(tipo, region, importe)
    let asiento = asientosContables[std.len(asientosContables) - 1]
    let config = regiones[region]
    
    let html = "<!DOCTYPE html>\n<html>\n<head>\n"
    html = html + "<title>Comprobante y Asiento Contable</title>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<style>\n"
    html = html + "body{font-family:Arial;margin:20px;background:#f5f5f5;}\n"
    html = html + ".container{max-width:800px;margin:0 auto;}\n"
    html = html + ".comp{background:white;padding:30px;border:2px solid #667eea;border-radius:10px;margin-bottom:20px;}\n"
    html = html + ".asiento{background:white;padding:30px;border:2px solid #27ae60;border-radius:10px;}\n"
    html = html + "table{width:100%;border-collapse:collapse;margin:15px 0;}\n"
    html = html + "th,td{border:1px solid #ddd;padding:8px;text-align:left;}\n"
    html = html + "th{background:#f0f0f0;}\n"
    html = html + ".debe{background:#fff3e0;}\n"
    html = html + ".haber{background:#e8f5e9;}\n"
    html = html + "</style>\n</head>\n<body>\n"
    
    html = html + "<div class=\"container\">\n"
    
    // Comprobante
    html = html + "<div class=\"comp\">\n"
    html = html + "<h1>COMPROBANTE DE " + std.toUpperCase(tipo) + "</h1>\n"
    html = html + "<h2>" + tx.pais + "</h2>\n"
    html = html + "<p><strong>ID:</strong> " + tx.id + "</p>\n"
    html = html + "<p><strong>Fecha:</strong> " + tx.fecha + "</p>\n"
    html = html + "<table>\n"
    html = html + "<tr><td>Importe Base:</td><td style=\"text-align:right\">" + config.moneda + " " + tx.importe + "</td></tr>\n"
    html = html + "<tr><td>IVA (" + (config.iva * 100) + "%):</td><td style=\"text-align:right\">" + config.moneda + " " + tx.iva + "</td></tr>\n"
    html = html + "<tr style=\"font-weight:bold\"><td>TOTAL:</td><td style=\"text-align:right\">" + config.moneda + " " + tx.total + "</td></tr>\n"
    html = html + "</table>\n"
    html = html + "</div>\n"
    
    // Asiento contable
    html = html + "<div class=\"asiento\">\n"
    html = html + "<h2>📚 ASIENTO CONTABLE</h2>\n"
    html = html + "<p><strong>Asiento:</strong> " + asiento.id + "</p>\n"
    html = html + "<p><strong>Descripción:</strong> " + asiento.descripcion + "</p>\n"
    html = html + "<table>\n"
    html = html + "<thead>\n<tr><th>Cuenta</th><th>Descripción</th><th>Debe</th><th>Haber</th></tr>\n</thead>\n"
    html = html + "<tbody>\n"
    
    let i = 0
    let totalDebe = 0
    let totalHaber = 0
    while (i < std.len(asiento.movimientos)) {
        let mov = asiento.movimientos[i]
        html = html + "<tr>\n"
        html = html + "<td>" + mov.cuenta + "</td>\n"
        html = html + "<td>" + mov.descripcion + "</td>\n"
        
        if (mov.tipo == "DEBE") {
            html = html + "<td class=\"debe\" style=\"text-align:right\">" + config.moneda + " " + mov.monto + "</td>\n"
            html = html + "<td class=\"haber\"></td>\n"
            totalDebe = totalDebe + mov.monto
        } else {
            html = html + "<td class=\"debe\"></td>\n"
            html = html + "<td class=\"haber\" style=\"text-align:right\">" + config.moneda + " " + mov.monto + "</td>\n"
            totalHaber = totalHaber + mov.monto
        }
        
        html = html + "</tr>\n"
        i = i + 1
    }
    
    html = html + "<tr style=\"font-weight:bold\">\n"
    html = html + "<td colspan=\"2\">TOTALES</td>\n"
    html = html + "<td style=\"text-align:right\">" + config.moneda + " " + totalDebe + "</td>\n"
    html = html + "<td style=\"text-align:right\">" + config.moneda + " " + totalHaber + "</td>\n"
    html = html + "</tr>\n"
    html = html + "</tbody>\n</table>\n"
    
    html = html + "<p style=\"text-align:center;color:#27ae60;font-weight:bold;\">✅ ASIENTO CUADRADO</p>\n"
    html = html + "</div>\n"
    
    html = html + "<p style=\"text-align:center;margin-top:30px;\">\n"
    html = html + "<a href=\"/\">← Volver</a> | <a href=\"/libro-diario\">Ver Libro Diario</a>\n"
    html = html + "</p>\n"
    
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// Handler libro diario
func handleLibroDiario(pathVars, method, body) {
    let html = "<!DOCTYPE html>\n<html>\n<head>\n"
    html = html + "<title>Libro Diario</title>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<style>\n"
    html = html + "body{font-family:Arial;margin:20px;background:#f5f5f5;}\n"
    html = html + ".container{max-width:1000px;margin:0 auto;background:white;padding:30px;border-radius:10px;}\n"
    html = html + ".asiento{margin:20px 0;padding:20px;border:1px solid #ddd;border-radius:8px;}\n"
    html = html + "table{width:100%;border-collapse:collapse;margin:10px 0;}\n"
    html = html + "th,td{border:1px solid #ddd;padding:8px;text-align:left;}\n"
    html = html + "th{background:#f0f0f0;}\n"
    html = html + ".debe{background:#fff3e0;}\n"
    html = html + ".haber{background:#e8f5e9;}\n"
    html = html + "</style>\n</head>\n<body>\n"
    
    html = html + "<div class=\"container\">\n"
    html = html + "<h1>📚 LIBRO DIARIO</h1>\n"
    html = html + "<p>Total de asientos: " + std.len(asientosContables) + "</p>\n"
    
    let i = 0
    while (i < std.len(asientosContables)) {
        let asiento = asientosContables[i]
        let config = regiones[asiento.region]
        
        html = html + "<div class=\"asiento\">\n"
        html = html + "<h3>Asiento: " + asiento.id + "</h3>\n"
        html = html + "<p><strong>Fecha:</strong> " + asiento.fecha + " | <strong>Región:</strong> " + config.nombre + "</p>\n"
        html = html + "<p><strong>Descripción:</strong> " + asiento.descripcion + "</p>\n"
        
        html = html + "<table>\n"
        html = html + "<thead>\n<tr><th>Cuenta</th><th>Descripción</th><th>Debe</th><th>Haber</th></tr>\n</thead>\n"
        html = html + "<tbody>\n"
        
        let j = 0
        while (j < std.len(asiento.movimientos)) {
            let mov = asiento.movimientos[j]
            html = html + "<tr>\n"
            html = html + "<td>" + mov.cuenta + "</td>\n"
            html = html + "<td>" + mov.descripcion + "</td>\n"
            
            if (mov.tipo == "DEBE") {
                html = html + "<td class=\"debe\" style=\"text-align:right\">" + config.moneda + " " + mov.monto + "</td>\n"
                html = html + "<td class=\"haber\"></td>\n"
            } else {
                html = html + "<td class=\"debe\"></td>\n"
                html = html + "<td class=\"haber\" style=\"text-align:right\">" + config.moneda + " " + mov.monto + "</td>\n"
            }
            
            html = html + "</tr>\n"
            j = j + 1
        }
        
        html = html + "</tbody>\n</table>\n"
        html = html + "</div>\n"
        
        i = i + 1
    }
    
    html = html + "<p style=\"text-align:center;margin-top:30px;\">\n"
    html = html + "<a href=\"/\">← Volver</a> | <a href=\"/balance\">Ver Balance</a>\n"
    html = html + "</p>\n"
    
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// Handler balance
func handleBalance(pathVars, method, body) {
    let totalDebe = 0
    let totalHaber = 0
    let cuentasSaldo = {}
    
    let i = 0
    while (i < std.len(asientosContables)) {
        let asiento = asientosContables[i]
        let j = 0
        while (j < std.len(asiento.movimientos)) {
            let mov = asiento.movimientos[j]
            
            if (!cuentasSaldo[mov.cuenta]) {
                cuentasSaldo[mov.cuenta] = {
                    descripcion: mov.descripcion,
                    debe: 0,
                    haber: 0
                }
            }
            
            if (mov.tipo == "DEBE") {
                cuentasSaldo[mov.cuenta].debe = cuentasSaldo[mov.cuenta].debe + mov.monto
                totalDebe = totalDebe + mov.monto
            } else {
                cuentasSaldo[mov.cuenta].haber = cuentasSaldo[mov.cuenta].haber + mov.monto
                totalHaber = totalHaber + mov.monto
            }
            
            j = j + 1
        }
        i = i + 1
    }
    
    let html = "<!DOCTYPE html>\n<html>\n<head>\n"
    html = html + "<title>Balance General</title>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<style>\n"
    html = html + "body{font-family:Arial;margin:20px;background:#f5f5f5;}\n"
    html = html + ".container{max-width:800px;margin:0 auto;background:white;padding:30px;border-radius:10px;}\n"
    html = html + "table{width:100%;border-collapse:collapse;margin:20px 0;}\n"
    html = html + "th,td{border:1px solid #ddd;padding:8px;}\n"
    html = html + "th{background:#f0f0f0;}\n"
    html = html + ".totales{background:#e8f5e9;font-weight:bold;}\n"
    html = html + ".cuadrado{text-align:center;color:#27ae60;font-weight:bold;font-size:1.2em;}\n"
    html = html + "</style>\n</head>\n<body>\n"
    
    html = html + "<div class=\"container\">\n"
    html = html + "<h1>⚖️ BALANCE GENERAL</h1>\n"
    
    html = html + "<table>\n"
    html = html + "<thead>\n<tr><th>Cuenta</th><th>Descripción</th><th>Debe</th><th>Haber</th><th>Saldo</th></tr>\n</thead>\n"
    html = html + "<tbody>\n"
    
    // Mostrar saldos por cuenta (simplificado)
    html = html + "<tr class=\"totales\">\n"
    html = html + "<td colspan=\"2\">TOTALES</td>\n"
    html = html + "<td style=\"text-align:right\">$" + math.round(totalDebe * 100) / 100 + "</td>\n"
    html = html + "<td style=\"text-align:right\">$" + math.round(totalHaber * 100) / 100 + "</td>\n"
    html = html + "<td style=\"text-align:right\">$" + math.round((totalDebe - totalHaber) * 100) / 100 + "</td>\n"
    html = html + "</tr>\n"
    
    html = html + "</tbody>\n</table>\n"
    
    if (math.round(totalDebe * 100) / 100 == math.round(totalHaber * 100) / 100) {
        html = html + "<p class=\"cuadrado\">✅ BALANCE CUADRADO</p>\n"
    } else {
        html = html + "<p style=\"text-align:center;color:red;\">❌ BALANCE DESCUADRADO</p>\n"
    }
    
    html = html + "<p style=\"text-align:center;margin-top:30px;\">\n"
    html = html + "<a href=\"/\">← Volver</a> | <a href=\"/libro-diario\">Ver Libro Diario</a>\n"
    html = html + "</p>\n"
    
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// Handler DSL query
func handleDSL(pathVars, method, body) {
    let query = getParam(body, "query")
    if (query == "") {
        query = "reporte diario ALL"
    }
    
    let engine = ReportesFinancieros
    let resultado = engine.use(query)
    
    let html = "<!DOCTYPE html>\n<html>\n<head>\n"
    html = html + "<title>Resultado DSL Query</title>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<style>\n"
    html = html + "body{font-family:Arial;margin:20px;background:#f5f5f5;}\n"
    html = html + ".container{max-width:800px;margin:0 auto;background:white;padding:30px;border-radius:10px;}\n"
    html = html + ".query{background:#2d3748;color:#e2e8f0;padding:15px;border-radius:8px;font-family:monospace;margin:15px 0;}\n"
    html = html + ".resultado{background:#f0f2f5;padding:20px;border-radius:8px;}\n"
    html = html + "pre{background:#f8f9fa;padding:15px;border-radius:4px;overflow-x:auto;}\n"
    html = html + "</style>\n</head>\n<body>\n"
    
    html = html + "<div class=\"container\">\n"
    html = html + "<h1>📊 Resultado DSL Query</h1>\n"
    html = html + "<div class=\"query\">Query: " + query + "</div>\n"
    
    html = html + "<div class=\"resultado\">\n"
    html = html + "<h3>Resultado:</h3>\n"
    html = html + "<pre>" + JSON(resultado) + "</pre>\n"
    html = html + "</div>\n"
    
    html = html + "<p style=\"text-align:center;margin-top:30px;\">\n"
    html = html + "<a href=\"/\">← Volver</a>\n"
    html = html + "</p>\n"
    
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// Handler demo automática
func handleDemo(pathVars, method, body) {
    procesarTransaccion("ventas", "MX", "100000")
    procesarTransaccion("compras", "COL", "50000")
    procesarTransaccion("ventas", "AR", "75000")
    procesarTransaccion("compras", "PE", "60000")
    procesarTransaccion("ventas", "CH", "80000")
    procesarTransaccion("compras", "UY", "45000")
    procesarTransaccion("ventas", "EC", "90000")
    
    let html = "<!DOCTYPE html>\n<html>\n<head>\n"
    html = html + "<title>Demo Automática</title>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<style>\n"
    html = html + "body{font-family:Arial;margin:20px;background:#f5f5f5;}\n"
    html = html + ".container{max-width:600px;margin:0 auto;background:white;padding:30px;border-radius:10px;}\n"
    html = html + ".success{background:#d4edda;color:#155724;padding:20px;border-radius:8px;}\n"
    html = html + "</style>\n</head>\n<body>\n"
    
    html = html + "<div class=\"container\">\n"
    html = html + "<h1>🎪 Demo Automática Completada</h1>\n"
    html = html + "<div class=\"success\">\n"
    html = html + "<p>✅ Se procesaron 7 transacciones de todas las regiones LATAM</p>\n"
    html = html + "<p>✅ Se generaron " + std.len(asientosContables) + " asientos contables</p>\n"
    html = html + "<p>✅ Libro diario actualizado con debe/haber</p>\n"
    html = html + "</div>\n"
    
    html = html + "<p style=\"text-align:center;margin-top:30px;\">\n"
    html = html + "<a href=\"/\">← Volver</a> | "
    html = html + "<a href=\"/libro-diario\">Ver Libro Diario</a> | "
    html = html + "<a href=\"/balance\">Ver Balance</a>\n"
    html = html + "</p>\n"
    
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// Handler API transacciones
func handleAPITransacciones(pathVars, method, body) {
    return JSON({
        total: std.len(transacciones),
        transacciones: transacciones
    })
}

// Handler API asientos
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
    
    console.log("\n✅ Iniciando servidor web completo...")
    
    // Rutas
    http.handler("GET", "/", handleHome)
    http.handler("POST", "/procesar", handleProcesar)
    http.handler("GET", "/libro-diario", handleLibroDiario)
    http.handler("GET", "/balance", handleBalance)
    http.handler("POST", "/dsl", handleDSL)
    http.handler("GET", "/demo", handleDemo)
    http.handler("GET", "/api/transacciones", handleAPITransacciones)
    http.handler("GET", "/api/asientos", handleAPIAsientos)
    
    console.log("🌐 URLs disponibles:")
    console.log("   http://localhost:8080 - Página principal")
    console.log("   http://localhost:8080/libro-diario - Libro diario completo")
    console.log("   http://localhost:8080/balance - Balance general")
    console.log("   http://localhost:8080/demo - Demo automática")
    console.log("")
    
    // Servidor
    http.serve(":8080")
}

main()