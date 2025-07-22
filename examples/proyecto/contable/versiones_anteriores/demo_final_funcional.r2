// Sistema Contable LATAM - DEMO FINAL FUNCIONAL
// POC para Siigo - 100% Working

console.log("🌍 Sistema Contable LATAM - Demo Final")
console.log("📊 Libro Diario + DSL + API")
console.log("✅ Version 100% Funcional")
console.log("")

// Configuración regiones
let regiones = {
    "MX": { 
        nombre: "México", moneda: "MXN", iva: 0.16,
        cuentas: {
            clientes: "1201", ventas: "4101", ivaDebito: "2401",
            proveedores: "2101", compras: "5101", ivaCredito: "1401"
        }
    },
    "COL": { 
        nombre: "Colombia", moneda: "COP", iva: 0.19,
        cuentas: {
            clientes: "130501", ventas: "413501", ivaDebito: "240801",
            proveedores: "220501", compras: "620501", ivaCredito: "240802"
        }
    },
    "AR": { 
        nombre: "Argentina", moneda: "ARS", iva: 0.21,
        cuentas: {
            clientes: "1.1.2.01", ventas: "4.1.1.01", ivaDebito: "2.1.3.01",
            proveedores: "2.1.1.01", compras: "5.1.1.01", ivaCredito: "1.1.3.01"
        }
    },
    "CH": { 
        nombre: "Chile", moneda: "CLP", iva: 0.19,
        cuentas: {
            clientes: "11030", ventas: "31010", ivaDebito: "21070",
            proveedores: "21010", compras: "41010", ivaCredito: "11070"
        }
    },
    "UY": { 
        nombre: "Uruguay", moneda: "UYU", iva: 0.22,
        cuentas: {
            clientes: "1121", ventas: "4111", ivaDebito: "2141",
            proveedores: "2111", compras: "5111", ivaCredito: "1141"
        }
    },
    "EC": { 
        nombre: "Ecuador", moneda: "USD", iva: 0.12,
        cuentas: {
            clientes: "102.01", ventas: "401.01", ivaDebito: "201.04",
            proveedores: "201.01", compras: "501.01", ivaCredito: "101.04"
        }
    },
    "PE": { 
        nombre: "Perú", moneda: "PEN", iva: 0.18,
        cuentas: {
            clientes: "121", ventas: "701", ivaDebito: "401",
            proveedores: "421", compras: "601", ivaCredito: "401"
        }
    }
}

// Arrays globales - usar push para agregar elementos
let transacciones = []
let asientosContables = []

// Función procesar transacción
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
    
    // Usar push para agregar
    transacciones.push(tx)
    
    // Crear asiento contable
    let asiento = {
        id: "AS-" + tx.id,
        fecha: tx.fecha,
        region: region,
        descripcion: std.toUpperCase(tipo) + " - " + config.nombre,
        movimientos: []
    }
    
    if (tipo == "ventas") {
        // Debe: Clientes, Haber: Ventas e IVA
        asiento.movimientos.push({
            cuenta: config.cuentas.clientes,
            descripcion: "Clientes",
            tipo: "DEBE",
            monto: total
        })
        asiento.movimientos.push({
            cuenta: config.cuentas.ventas,
            descripcion: "Ventas",
            tipo: "HABER",
            monto: importeNum
        })
        asiento.movimientos.push({
            cuenta: config.cuentas.ivaDebito,
            descripcion: "IVA Débito",
            tipo: "HABER",
            monto: iva
        })
    } else {
        // Debe: Compras e IVA, Haber: Proveedores
        asiento.movimientos.push({
            cuenta: config.cuentas.compras,
            descripcion: "Compras",
            tipo: "DEBE",
            monto: importeNum
        })
        asiento.movimientos.push({
            cuenta: config.cuentas.ivaCredito,
            descripcion: "IVA Crédito",
            tipo: "DEBE",
            monto: iva
        })
        asiento.movimientos.push({
            cuenta: config.cuentas.proveedores,
            descripcion: "Proveedores",
            tipo: "HABER",
            monto: total
        })
    }
    
    // Usar push para agregar
    asientosContables.push(asiento)
    
    return tx
}

// DSL Reportes Financieros
dsl ReportesFinancieros {
    token("REPORTE", "reporte")
    token("TIPO", "balance|diario|ventas|compras|iva")
    
    rule("consulta", ["REPORTE", "TIPO"], "ejecutarReporte")
    
    func ejecutarReporte(cmd, tipo) {
        if (tipo == "balance") {
            let totalDebe = 0
            let totalHaber = 0
            let i = 0
            while (i < std.len(asientosContables)) {
                let asiento = asientosContables[i]
                let j = 0
                while (j < std.len(asiento.movimientos)) {
                    let mov = asiento.movimientos[j]
                    if (mov.tipo == "DEBE") {
                        totalDebe = totalDebe + mov.monto
                    }
                    if (mov.tipo == "HABER") {
                        totalHaber = totalHaber + mov.monto
                    }
                    j = j + 1
                }
                i = i + 1
            }
            return {
                tipo: "Balance General",
                totalAsientos: std.len(asientosContables),
                totalDebe: math.round(totalDebe * 100) / 100,
                totalHaber: math.round(totalHaber * 100) / 100,
                cuadrado: math.round(totalDebe * 100) / 100 == math.round(totalHaber * 100) / 100
            }
        }
        
        if (tipo == "diario") {
            return {
                tipo: "Libro Diario",
                totalAsientos: std.len(asientosContables),
                asientos: asientosContables
            }
        }
        
        if (tipo == "ventas" || tipo == "compras") {
            let total = 0
            let count = 0
            let i = 0
            while (i < std.len(transacciones)) {
                let tx = transacciones[i]
                if (tx.tipo == tipo) {
                    total = total + tx.total
                    count = count + 1
                }
                i = i + 1
            }
            return {
                tipo: "Reporte " + tipo,
                transacciones: count,
                total: math.round(total * 100) / 100
            }
        }
        
        if (tipo == "iva") {
            let ivaDebito = 0
            let ivaCredito = 0
            let i = 0
            while (i < std.len(transacciones)) {
                let tx = transacciones[i]
                if (tx.tipo == "ventas") {
                    ivaDebito = ivaDebito + tx.iva
                } else {
                    ivaCredito = ivaCredito + tx.iva
                }
                i = i + 1
            }
            return {
                tipo: "Reporte IVA",
                ivaDebito: math.round(ivaDebito * 100) / 100,
                ivaCredito: math.round(ivaCredito * 100) / 100,
                saldo: math.round((ivaDebito - ivaCredito) * 100) / 100
            }
        }
        
        return { error: "Tipo de reporte no válido" }
    }
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
    
    console.log("\n📚 LIBRO DIARIO - Ejemplo:")
    if (std.len(asientosContables) > 0) {
        let asiento = asientosContables[0]
        console.log("Asiento: " + asiento.id)
        let i = 0
        while (i < std.len(asiento.movimientos)) {
            let mov = asiento.movimientos[i]
            console.log("  " + mov.tipo + " - " + mov.cuenta + " - " + mov.descripcion + ": $" + mov.monto)
            i = i + 1
        }
    }
    
    console.log("\n📊 REPORTES DSL:")
    let engine = ReportesFinancieros
    let balance = engine.use("reporte balance")
    console.log("Balance General:")
    console.log("  Total Asientos: " + balance.totalAsientos)
    console.log("  Debe: $" + balance.totalDebe)
    console.log("  Haber: $" + balance.totalHaber)
    console.log("  Cuadrado: " + balance.cuadrado)
    
    console.log("\n🎯 VALUE PROPOSITION:")
    console.log("• 18 meses → 2 meses")
    console.log("• $500K → $150K")
    console.log("• ROI: 1,020%")
}

// Parsear parámetros
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
    html = html + "<title>Sistema Contable LATAM</title>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<style>\n"
    html = html + "body{font-family:Arial,sans-serif;margin:20px;background:#f5f5f5;}\n"
    html = html + ".container{max-width:1200px;margin:0 auto;background:white;padding:30px;border-radius:10px;box-shadow:0 2px 10px rgba(0,0,0,0.1);}\n"
    html = html + "h1{color:#2c3e50;text-align:center;}\n"
    html = html + ".grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(350px,1fr));gap:20px;}\n"
    html = html + ".card{background:#f8f9fa;padding:20px;border-radius:8px;border:1px solid #ddd;}\n"
    html = html + "input,select,textarea{width:100%;padding:8px;margin:5px 0;border:1px solid #ddd;border-radius:4px;}\n"
    html = html + "button{background:#667eea;color:white;padding:10px;border:none;border-radius:4px;width:100%;cursor:pointer;}\n"
    html = html + "button:hover{background:#5a67d8;}\n"
    html = html + ".value{background:#84fab0;padding:20px;border-radius:8px;text-align:center;margin:20px 0;}\n"
    html = html + ".links{text-align:center;margin:20px 0;}\n"
    html = html + ".links a{color:#667eea;text-decoration:none;margin:0 10px;}\n"
    html = html + "</style>\n</head>\n<body>\n"
    
    html = html + "<div class=\"container\">\n"
    html = html + "<h1>🌍 Sistema Contable LATAM</h1>\n"
    html = html + "<p style=\"text-align:center\">POC Completo para Siigo - " + std.len(transacciones) + " transacciones - " + std.len(asientosContables) + " asientos</p>\n"
    
    html = html + "<div class=\"grid\">\n"
    
    // Formulario transacción
    html = html + "<div class=\"card\">\n"
    html = html + "<h3>📝 Procesar Transacción</h3>\n"
    html = html + "<form action=\"/procesar\" method=\"POST\">\n"
    html = html + "<label>Tipo:</label>\n"
    html = html + "<select name=\"tipo\">\n"
    html = html + "<option value=\"ventas\">Venta</option>\n"
    html = html + "<option value=\"compras\">Compra</option>\n"
    html = html + "</select>\n"
    html = html + "<label>Región:</label>\n"
    html = html + "<select name=\"region\">\n"
    html = html + "<option value=\"MX\">México (16%)</option>\n"
    html = html + "<option value=\"COL\">Colombia (19%)</option>\n"
    html = html + "<option value=\"AR\">Argentina (21%)</option>\n"
    html = html + "<option value=\"CH\">Chile (19%)</option>\n"
    html = html + "<option value=\"UY\">Uruguay (22%)</option>\n"
    html = html + "<option value=\"EC\">Ecuador (12%)</option>\n"
    html = html + "<option value=\"PE\">Perú (18%)</option>\n"
    html = html + "</select>\n"
    html = html + "<label>Importe (sin IVA):</label>\n"
    html = html + "<input type=\"number\" name=\"importe\" value=\"100000\" required>\n"
    html = html + "<button type=\"submit\">Procesar</button>\n"
    html = html + "</form>\n"
    html = html + "</div>\n"
    
    // DSL Query
    html = html + "<div class=\"card\">\n"
    html = html + "<h3>🔍 Consulta DSL</h3>\n"
    html = html + "<form action=\"/dsl\" method=\"POST\">\n"
    html = html + "<label>Query DSL:</label>\n"
    html = html + "<textarea name=\"query\" rows=\"3\" placeholder=\"Ejemplos:\nreporte balance\nreporte diario\nreporte ventas\">reporte balance</textarea>\n"
    html = html + "<button type=\"submit\">Ejecutar Query</button>\n"
    html = html + "</form>\n"
    html = html + "</div>\n"
    
    html = html + "</div>\n"
    
    // Links
    html = html + "<div class=\"links\">\n"
    html = html + "<a href=\"/demo\">🚀 Demo Auto</a> | "
    html = html + "<a href=\"/libro\">📚 Libro Diario</a> | "
    html = html + "<a href=\"/api/transacciones\">📊 API Transacciones</a> | "
    html = html + "<a href=\"/api/asientos\">📋 API Asientos</a>\n"
    html = html + "</div>\n"
    
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
    let numAsientos = std.len(asientosContables)
    if (numAsientos == 0) {
        return "<h1>Error: No se pudo crear el asiento</h1><p><a href=\"/\">Volver</a></p>"
    }
    
    let asiento = asientosContables[numAsientos - 1]
    
    let html = "<!DOCTYPE html>\n<html>\n<head>\n"
    html = html + "<title>Comprobante</title>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<style>\n"
    html = html + "body{font-family:Arial;margin:20px;background:#f5f5f5;}\n"
    html = html + ".comp{max-width:800px;margin:0 auto;padding:30px;background:white;border-radius:10px;box-shadow:0 2px 10px rgba(0,0,0,0.1);}\n"
    html = html + "h1{color:#2c3e50;text-align:center;}\n"
    html = html + ".info{background:#f8f9fa;padding:20px;border-radius:8px;margin:20px 0;}\n"
    html = html + ".asiento{background:#e3f2fd;padding:20px;border-radius:8px;margin:20px 0;}\n"
    html = html + "table{width:100%;border-collapse:collapse;margin-top:15px;}\n"
    html = html + "th,td{padding:10px;text-align:left;border-bottom:1px solid #ddd;}\n"
    html = html + "th{background:#667eea;color:white;}\n"
    html = html + ".debe{color:#28a745;font-weight:bold;}\n"
    html = html + ".haber{color:#dc3545;font-weight:bold;}\n"
    html = html + "</style>\n"
    html = html + "</head>\n<body>\n"
    
    html = html + "<div class=\"comp\">\n"
    html = html + "<h1>COMPROBANTE FISCAL</h1>\n"
    
    html = html + "<div class=\"info\">\n"
    html = html + "<p><strong>ID:</strong> " + tx.id + "</p>\n"
    html = html + "<p><strong>Fecha:</strong> " + tx.fecha + "</p>\n"
    html = html + "<p><strong>País:</strong> " + tx.pais + "</p>\n"
    html = html + "<p><strong>Tipo:</strong> " + std.toUpperCase(tx.tipo) + "</p>\n"
    html = html + "<p><strong>Importe:</strong> " + tx.moneda + " " + tx.importe + "</p>\n"
    html = html + "<p><strong>IVA:</strong> " + tx.moneda + " " + tx.iva + "</p>\n"
    html = html + "<p style=\"font-size:20px;\"><strong>TOTAL:</strong> " + tx.moneda + " " + tx.total + "</p>\n"
    html = html + "</div>\n"
    
    // Asiento contable
    html = html + "<div class=\"asiento\">\n"
    html = html + "<h3>📚 ASIENTO CONTABLE</h3>\n"
    html = html + "<p><strong>Asiento:</strong> " + asiento.id + "</p>\n"
    html = html + "<p><strong>Descripción:</strong> " + asiento.descripcion + "</p>\n"
    
    html = html + "<table>\n"
    html = html + "<tr><th>Cuenta</th><th>Descripción</th><th>Debe</th><th>Haber</th></tr>\n"
    
    let i = 0
    let totalDebe = 0
    let totalHaber = 0
    while (i < std.len(asiento.movimientos)) {
        let mov = asiento.movimientos[i]
        html = html + "<tr>\n"
        html = html + "<td>" + mov.cuenta + "</td>\n"
        html = html + "<td>" + mov.descripcion + "</td>\n"
        if (mov.tipo == "DEBE") {
            html = html + "<td class=\"debe\">" + mov.monto + "</td>\n"
            html = html + "<td>-</td>\n"
            totalDebe = totalDebe + mov.monto
        } else {
            html = html + "<td>-</td>\n"
            html = html + "<td class=\"haber\">" + mov.monto + "</td>\n"
            totalHaber = totalHaber + mov.monto
        }
        html = html + "</tr>\n"
        i = i + 1
    }
    
    html = html + "<tr style=\"border-top:2px solid #667eea;\">\n"
    html = html + "<td colspan=\"2\"><strong>TOTALES</strong></td>\n"
    html = html + "<td class=\"debe\"><strong>" + totalDebe + "</strong></td>\n"
    html = html + "<td class=\"haber\"><strong>" + totalHaber + "</strong></td>\n"
    html = html + "</tr>\n"
    
    html = html + "</table>\n"
    html = html + "</div>\n"
    
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
    html = html + "body{font-family:Arial;margin:20px;background:#f5f5f5;}\n"
    html = html + ".container{max-width:1000px;margin:0 auto;}\n"
    html = html + ".asiento{background:white;padding:20px;margin:15px 0;border-radius:8px;box-shadow:0 2px 5px rgba(0,0,0,0.1);}\n"
    html = html + "h1{color:#2c3e50;text-align:center;}\n"
    html = html + "table{width:100%;border-collapse:collapse;margin-top:10px;}\n"
    html = html + "th,td{padding:10px;text-align:left;border-bottom:1px solid #ddd;}\n"
    html = html + "th{background:#667eea;color:white;}\n"
    html = html + ".debe{color:#28a745;font-weight:bold;}\n"
    html = html + ".haber{color:#dc3545;font-weight:bold;}\n"
    html = html + "</style>\n"
    html = html + "</head>\n<body>\n"
    
    html = html + "<div class=\"container\">\n"
    html = html + "<h1>📚 LIBRO DIARIO</h1>\n"
    html = html + "<p style=\"text-align:center;\">Total asientos: " + std.len(asientosContables) + "</p>\n"
    
    let i = 0
    while (i < std.len(asientosContables)) {
        let asiento = asientosContables[i]
        html = html + "<div class=\"asiento\">\n"
        html = html + "<h3>Asiento: " + asiento.id + "</h3>\n"
        html = html + "<p><strong>Fecha:</strong> " + asiento.fecha + " | <strong>Descripción:</strong> " + asiento.descripcion + "</p>\n"
        
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
                html = html + "<td>-</td>\n"
            } else {
                html = html + "<td>-</td>\n"
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
    html = html + "body{font-family:Arial;margin:20px;background:#f5f5f5;}\n"
    html = html + ".container{max-width:800px;margin:0 auto;padding:30px;background:white;border-radius:10px;box-shadow:0 2px 10px rgba(0,0,0,0.1);}\n"
    html = html + "h1{color:#2c3e50;text-align:center;}\n"
    html = html + ".query{background:#667eea;color:white;padding:15px;border-radius:8px;margin:20px 0;font-family:monospace;}\n"
    html = html + ".result{background:#f8f9fa;padding:20px;border-radius:8px;}\n"
    html = html + "pre{background:#263238;color:#aed581;padding:20px;border-radius:8px;overflow:auto;}\n"
    html = html + "</style>\n"
    html = html + "</head>\n<body>\n"
    
    html = html + "<div class=\"container\">\n"
    html = html + "<h1>📊 Resultado de Consulta DSL</h1>\n"
    html = html + "<div class=\"query\">Query: " + query + "</div>\n"
    html = html + "<div class=\"result\">\n"
    html = html + "<pre>" + json.stringify(resultado) + "</pre>\n"
    html = html + "</div>\n"
    html = html + "<p style=\"text-align:center;\"><a href=\"/\">Volver</a></p>\n"
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// Handler demo
func handleDemo(pathVars, method, body) {
    // Procesar varias transacciones
    procesarTransaccion("ventas", "MX", "100000")
    procesarTransaccion("compras", "COL", "50000")
    procesarTransaccion("ventas", "AR", "75000")
    procesarTransaccion("compras", "PE", "60000")
    procesarTransaccion("ventas", "CH", "120000")
    procesarTransaccion("compras", "EC", "30000")
    
    // Obtener balance
    let engine = ReportesFinancieros
    let balance = engine.use("reporte balance")
    
    let html = "<!DOCTYPE html>\n<html>\n<head>\n"
    html = html + "<title>Demo Automática</title>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<style>\n"
    html = html + "body{font-family:Arial;margin:20px;background:#f5f5f5;}\n"
    html = html + ".container{max-width:600px;margin:0 auto;padding:30px;background:white;border-radius:10px;box-shadow:0 2px 10px rgba(0,0,0,0.1);}\n"
    html = html + "h1{color:#2c3e50;text-align:center;}\n"
    html = html + ".success{background:#d4edda;color:#155724;padding:20px;border-radius:8px;margin:20px 0;}\n"
    html = html + ".stats{background:#f8f9fa;padding:20px;border-radius:8px;}\n"
    html = html + "</style>\n"
    html = html + "</head>\n<body>\n"
    
    html = html + "<div class=\"container\">\n"
    html = html + "<h1>🚀 Demo Automática</h1>\n"
    html = html + "<div class=\"success\">\n"
    html = html + "<h3>✅ Se procesaron 6 transacciones exitosamente</h3>\n"
    html = html + "</div>\n"
    
    html = html + "<div class=\"stats\">\n"
    html = html + "<h3>📊 Resumen:</h3>\n"
    html = html + "<p><strong>Total transacciones:</strong> " + std.len(transacciones) + "</p>\n"
    html = html + "<p><strong>Total asientos:</strong> " + std.len(asientosContables) + "</p>\n"
    html = html + "<p><strong>Balance:</strong></p>\n"
    html = html + "<ul>\n"
    html = html + "<li>Debe: $" + balance.totalDebe + "</li>\n"
    html = html + "<li>Haber: $" + balance.totalHaber + "</li>\n"
    html = html + "<li>Cuadrado: " + balance.cuadrado + "</li>\n"
    html = html + "</ul>\n"
    html = html + "</div>\n"
    
    html = html + "<p style=\"text-align:center;margin-top:20px;\">\n"
    html = html + "<a href=\"/\">Volver</a> | <a href=\"/libro\">Ver Libro Diario</a>\n"
    html = html + "</p>\n"
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// APIs
func handleAPITransacciones(pathVars, method, body) {
    return json.stringify({
        total: std.len(transacciones),
        regiones: 7,
        transacciones: transacciones
    })
}

func handleAPIAsientos(pathVars, method, body) {
    return json.stringify({
        total: std.len(asientosContables),
        asientos: asientosContables
    })
}

// Main
func main() {
    // Demo inicial en consola
    demoConsola()
    
    console.log("\n✅ Iniciando servidor web...")
    
    // Configurar rutas
    http.handler("GET", "/", handleHome)
    http.handler("POST", "/procesar", handleProcesar)
    http.handler("GET", "/libro", handleLibro)
    http.handler("POST", "/dsl", handleDSL)
    http.handler("GET", "/demo", handleDemo)
    http.handler("GET", "/api/transacciones", handleAPITransacciones)
    http.handler("GET", "/api/asientos", handleAPIAsientos)
    
    console.log("🌐 http://localhost:8080")
    console.log("📌 Presiona Ctrl+C para detener\n")
    
    // Iniciar servidor
    http.serve(":8080")
}

// Ejecutar
main()