// Sistema Contable LATAM - Demo Completa Final
// POC para Siigo con Libro Diario + DSL Reportes

console.log("🌍 Sistema Contable LATAM - POC Final")
console.log("📊 Libro Diario + DSL + API")
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

// Variables globales - usar objetos para asegurar persistencia
let datos = {
    transacciones: [],
    asientos: []
}

// Función procesar transacción
func procesarTransaccion(tipo, region, importe) {
    let config = regiones[region]
    if (!config) {
        console.log("ERROR: Región no válida: " + region)
        return null
    }
    
    let importeNum = std.parseFloat(importe)
    let iva = math.round((importeNum * config.iva) * 100) / 100
    let total = importeNum + iva
    
    // Crear transacción
    let tx = {
        id: region + "-" + math.randomInt(1000) + "-" + math.randomInt(999),
        tipo: tipo,
        region: region,
        pais: config.nombre,
        importe: importeNum,
        iva: iva,
        total: total,
        moneda: config.moneda,
        fecha: std.now()
    }
    
    // Agregar transacción
    datos.transacciones[std.len(datos.transacciones)] = tx
    
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
        // Debe: Compras e IVA, Haber: Proveedores
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
    
    // Agregar asiento
    datos.asientos[std.len(datos.asientos)] = asiento
    
    console.log("✅ Procesado: " + tipo + " " + region + " - Total: " + total)
    
    return tx
}

// DSL Reportes Financieros - versión mejorada
dsl ReportesFinancieros {
    token("REPORTE", "reporte")
    token("TIPO", "balance|diario|ventas|compras|iva")
    
    rule("consulta", ["REPORTE", "TIPO"], "ejecutarReporte")
    
    func ejecutarReporte(cmd, tipo) {
        console.log("DSL: Ejecutando reporte " + tipo)
        
        if (tipo == "balance") {
            let totalDebe = 0
            let totalHaber = 0
            let i = 0
            while (i < std.len(datos.asientos)) {
                let asiento = datos.asientos[i]
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
                totalAsientos: std.len(datos.asientos),
                totalDebe: math.round(totalDebe * 100) / 100,
                totalHaber: math.round(totalHaber * 100) / 100,
                cuadrado: math.round(totalDebe * 100) / 100 == math.round(totalHaber * 100) / 100
            }
        }
        
        if (tipo == "diario") {
            return {
                tipo: "Libro Diario",
                totalAsientos: std.len(datos.asientos),
                asientos: datos.asientos
            }
        }
        
        if (tipo == "ventas" || tipo == "compras") {
            let total = 0
            let count = 0
            let i = 0
            while (i < std.len(datos.transacciones)) {
                let tx = datos.transacciones[i]
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
            while (i < std.len(datos.transacciones)) {
                let tx = datos.transacciones[i]
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
        
        return { error: "Tipo de reporte no válido: " + tipo }
    }
}

// Demo en consola
func demoConsola() {
    console.log("=== DEMO SISTEMA CONTABLE ===")
    
    let tx1 = procesarTransaccion("ventas", "MX", "100000")
    let tx2 = procesarTransaccion("compras", "COL", "50000")
    let tx3 = procesarTransaccion("ventas", "AR", "75000")
    
    console.log("\n📚 LIBRO DIARIO:")
    console.log("Total asientos: " + std.len(datos.asientos))
    
    if (std.len(datos.asientos) > 0) {
        let asiento = datos.asientos[0]
        console.log("\nEjemplo - Asiento: " + asiento.id)
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
    console.log("  Debe: $" + balance.totalDebe)
    console.log("  Haber: $" + balance.totalHaber)
    console.log("  Cuadrado: " + balance.cuadrado)
    
    console.log("\n🎯 VALUE PROPOSITION:")
    console.log("• 18 meses → 2 meses")
    console.log("• $500K → $150K")
    console.log("• ROI: 1,020%")
}

// Parsear parámetros del body
func parseBody(body) {
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
                let value = std.replace(parts[1], "+", " ")
                params[key] = value
            }
        }
        i = i + 1
    }
    return params
}

// Handler página principal
func handleHome(pathVars, method, body) {
    let html = "<!DOCTYPE html>\n"
    html = html + "<html>\n<head>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<title>Sistema Contable LATAM</title>\n"
    html = html + "<style>\n"
    html = html + "body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background: #f5f5f5; }\n"
    html = html + ".container { max-width: 1200px; margin: 0 auto; background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }\n"
    html = html + "h1 { color: #2c3e50; text-align: center; margin-bottom: 10px; }\n"
    html = html + ".subtitle { text-align: center; color: #7f8c8d; margin-bottom: 30px; }\n"
    html = html + ".grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(400px, 1fr)); gap: 20px; margin-bottom: 30px; }\n"
    html = html + ".card { background: #f8f9fa; padding: 25px; border-radius: 8px; border: 1px solid #e9ecef; }\n"
    html = html + ".card h3 { color: #495057; margin-top: 0; }\n"
    html = html + "input, select, textarea { width: 100%; padding: 10px; margin: 8px 0; border: 1px solid #ced4da; border-radius: 4px; box-sizing: border-box; }\n"
    html = html + "button { background: #667eea; color: white; padding: 12px 24px; border: none; border-radius: 4px; cursor: pointer; width: 100%; font-size: 16px; }\n"
    html = html + "button:hover { background: #5a67d8; }\n"
    html = html + ".links { text-align: center; margin: 20px 0; }\n"
    html = html + ".links a { color: #667eea; text-decoration: none; margin: 0 10px; }\n"
    html = html + ".links a:hover { text-decoration: underline; }\n"
    html = html + ".value-prop { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 30px; border-radius: 8px; text-align: center; }\n"
    html = html + ".value-prop h3 { color: white; margin-top: 0; }\n"
    html = html + ".stats { display: flex; justify-content: space-around; margin-top: 20px; }\n"
    html = html + ".stat { text-align: center; }\n"
    html = html + ".stat-value { font-size: 24px; font-weight: bold; }\n"
    html = html + ".stat-label { font-size: 14px; opacity: 0.9; }\n"
    html = html + "</style>\n"
    html = html + "</head>\n<body>\n"
    
    html = html + "<div class=\"container\">\n"
    html = html + "<h1>🌍 Sistema Contable LATAM</h1>\n"
    html = html + "<p class=\"subtitle\">POC Completo para Siigo - " + std.len(datos.transacciones) + " transacciones - " + std.len(datos.asientos) + " asientos contables</p>\n"
    
    html = html + "<div class=\"grid\">\n"
    
    // Formulario de transacción
    html = html + "<div class=\"card\">\n"
    html = html + "<h3>📝 Procesar Transacción</h3>\n"
    html = html + "<form action=\"/procesar\" method=\"POST\">\n"
    html = html + "<label>Tipo de transacción:</label>\n"
    html = html + "<select name=\"tipo\" required>\n"
    html = html + "<option value=\"ventas\">Venta</option>\n"
    html = html + "<option value=\"compras\">Compra</option>\n"
    html = html + "</select>\n"
    html = html + "<label>Región:</label>\n"
    html = html + "<select name=\"region\" required>\n"
    html = html + "<option value=\"MX\">🇲🇽 México (IVA 16%)</option>\n"
    html = html + "<option value=\"COL\">🇨🇴 Colombia (IVA 19%)</option>\n"
    html = html + "<option value=\"AR\">🇦🇷 Argentina (IVA 21%)</option>\n"
    html = html + "<option value=\"CH\">🇨🇱 Chile (IVA 19%)</option>\n"
    html = html + "<option value=\"UY\">🇺🇾 Uruguay (IVA 22%)</option>\n"
    html = html + "<option value=\"EC\">🇪🇨 Ecuador (IVA 12%)</option>\n"
    html = html + "<option value=\"PE\">🇵🇪 Perú (IVA 18%)</option>\n"
    html = html + "</select>\n"
    html = html + "<label>Importe (sin IVA):</label>\n"
    html = html + "<input type=\"number\" name=\"importe\" value=\"100000\" required>\n"
    html = html + "<button type=\"submit\">Procesar Transacción</button>\n"
    html = html + "</form>\n"
    html = html + "</div>\n"
    
    // DSL Query
    html = html + "<div class=\"card\">\n"
    html = html + "<h3>🔍 Consulta DSL de Reportes</h3>\n"
    html = html + "<form action=\"/dsl\" method=\"POST\">\n"
    html = html + "<label>Ingrese su consulta DSL:</label>\n"
    html = html + "<textarea name=\"query\" rows=\"4\" placeholder=\"Ejemplos:\nreporte balance\nreporte diario\nreporte ventas\nreporte compras\nreporte iva\" required>reporte balance</textarea>\n"
    html = html + "<button type=\"submit\">Ejecutar Consulta DSL</button>\n"
    html = html + "</form>\n"
    html = html + "</div>\n"
    
    html = html + "</div>\n"
    
    // Enlaces
    html = html + "<div class=\"links\">\n"
    html = html + "<a href=\"/demo\">🚀 Demo Automática</a> | "
    html = html + "<a href=\"/libro\">📚 Libro Diario</a> | "
    html = html + "<a href=\"/api/transacciones\">📊 API Transacciones</a> | "
    html = html + "<a href=\"/api/asientos\">📋 API Asientos</a>\n"
    html = html + "</div>\n"
    
    // Value proposition
    html = html + "<div class=\"value-prop\">\n"
    html = html + "<h3>🎯 Value Proposition para Siigo</h3>\n"
    html = html + "<div class=\"stats\">\n"
    html = html + "<div class=\"stat\">\n"
    html = html + "<div class=\"stat-value\">18→2</div>\n"
    html = html + "<div class=\"stat-label\">meses de desarrollo</div>\n"
    html = html + "</div>\n"
    html = html + "<div class=\"stat\">\n"
    html = html + "<div class=\"stat-value\">$500K→$150K</div>\n"
    html = html + "<div class=\"stat-label\">costo por país</div>\n"
    html = html + "</div>\n"
    html = html + "<div class=\"stat\">\n"
    html = html + "<div class=\"stat-value\">7→1</div>\n"
    html = html + "<div class=\"stat-label\">sistemas unificados</div>\n"
    html = html + "</div>\n"
    html = html + "<div class=\"stat\">\n"
    html = html + "<div class=\"stat-value\">1,020%</div>\n"
    html = html + "<div class=\"stat-label\">ROI en 3 años</div>\n"
    html = html + "</div>\n"
    html = html + "</div>\n"
    html = html + "</div>\n"
    
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// Handler procesar transacción
func handleProcesar(pathVars, method, body) {
    let params = parseBody(body)
    let tipo = params.tipo || "ventas"
    let region = params.region || "COL"
    let importe = params.importe || "100000"
    
    console.log("Procesando: " + tipo + " en " + region + " por " + importe)
    
    let tx = procesarTransaccion(tipo, region, importe)
    
    if (!tx) {
        return "<h1>Error al procesar la transacción</h1><p><a href=\"/\">Volver</a></p>"
    }
    
    // Obtener el último asiento
    let ultimoIndice = std.len(datos.asientos) - 1
    if (ultimoIndice < 0) {
        return "<h1>Error: No se pudo crear el asiento contable</h1><p><a href=\"/\">Volver</a></p>"
    }
    
    let asiento = datos.asientos[ultimoIndice]
    
    let html = "<!DOCTYPE html>\n<html>\n<head>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<title>Comprobante</title>\n"
    html = html + "<style>\n"
    html = html + "body { font-family: Arial; margin: 20px; background: #f5f5f5; }\n"
    html = html + ".comprobante { max-width: 800px; margin: 0 auto; background: white; padding: 40px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }\n"
    html = html + "h1 { color: #2c3e50; text-align: center; }\n"
    html = html + ".info { background: #f8f9fa; padding: 20px; border-radius: 8px; margin: 20px 0; }\n"
    html = html + ".info p { margin: 10px 0; }\n"
    html = html + ".asiento { background: #e3f2fd; padding: 20px; border-radius: 8px; margin: 20px 0; }\n"
    html = html + "table { width: 100%; border-collapse: collapse; margin-top: 15px; }\n"
    html = html + "th, td { padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }\n"
    html = html + "th { background: #667eea; color: white; }\n"
    html = html + ".debe { color: #28a745; font-weight: bold; }\n"
    html = html + ".haber { color: #dc3545; font-weight: bold; }\n"
    html = html + ".total { font-size: 24px; font-weight: bold; color: #667eea; text-align: center; margin: 20px 0; }\n"
    html = html + ".back { text-align: center; margin-top: 30px; }\n"
    html = html + ".back a { background: #667eea; color: white; padding: 10px 20px; text-decoration: none; border-radius: 4px; }\n"
    html = html + "</style>\n"
    html = html + "</head>\n<body>\n"
    
    html = html + "<div class=\"comprobante\">\n"
    html = html + "<h1>COMPROBANTE FISCAL</h1>\n"
    
    html = html + "<div class=\"info\">\n"
    html = html + "<p><strong>ID Transacción:</strong> " + tx.id + "</p>\n"
    html = html + "<p><strong>Fecha:</strong> " + tx.fecha + "</p>\n"
    html = html + "<p><strong>País:</strong> " + tx.pais + "</p>\n"
    html = html + "<p><strong>Tipo:</strong> " + std.toUpperCase(tx.tipo) + "</p>\n"
    html = html + "<p><strong>Importe:</strong> " + tx.moneda + " " + tx.importe + "</p>\n"
    html = html + "<p><strong>IVA:</strong> " + tx.moneda + " " + tx.iva + "</p>\n"
    html = html + "</div>\n"
    
    html = html + "<div class=\"total\">TOTAL: " + tx.moneda + " " + tx.total + "</div>\n"
    
    // Asiento contable
    html = html + "<div class=\"asiento\">\n"
    html = html + "<h3>📚 ASIENTO CONTABLE</h3>\n"
    html = html + "<p><strong>Asiento No.:</strong> " + asiento.id + "</p>\n"
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
    
    html = html + "<tr style=\"border-top: 2px solid #667eea;\">\n"
    html = html + "<td colspan=\"2\"><strong>TOTALES</strong></td>\n"
    html = html + "<td class=\"debe\"><strong>" + totalDebe + "</strong></td>\n"
    html = html + "<td class=\"haber\"><strong>" + totalHaber + "</strong></td>\n"
    html = html + "</tr>\n"
    
    html = html + "</table>\n"
    html = html + "</div>\n"
    
    html = html + "<div class=\"back\">\n"
    html = html + "<a href=\"/\">Volver al inicio</a>\n"
    html = html + "</div>\n"
    
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// Handler libro diario
func handleLibro(pathVars, method, body) {
    let html = "<!DOCTYPE html>\n<html>\n<head>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<title>Libro Diario</title>\n"
    html = html + "<style>\n"
    html = html + "body { font-family: Arial; margin: 20px; background: #f5f5f5; }\n"
    html = html + ".container { max-width: 1200px; margin: 0 auto; }\n"
    html = html + "h1 { color: #2c3e50; text-align: center; }\n"
    html = html + ".asiento { background: white; padding: 20px; margin: 15px 0; border-radius: 8px; box-shadow: 0 2px 5px rgba(0,0,0,0.1); }\n"
    html = html + ".asiento-header { background: #f8f9fa; padding: 15px; margin: -20px -20px 15px -20px; border-radius: 8px 8px 0 0; }\n"
    html = html + "table { width: 100%; border-collapse: collapse; }\n"
    html = html + "th, td { padding: 10px; text-align: left; border-bottom: 1px solid #e9ecef; }\n"
    html = html + "th { background: #667eea; color: white; }\n"
    html = html + ".debe { color: #28a745; font-weight: bold; }\n"
    html = html + ".haber { color: #dc3545; font-weight: bold; }\n"
    html = html + ".back { text-align: center; margin: 30px 0; }\n"
    html = html + ".back a { background: #667eea; color: white; padding: 10px 20px; text-decoration: none; border-radius: 4px; }\n"
    html = html + "</style>\n"
    html = html + "</head>\n<body>\n"
    
    html = html + "<div class=\"container\">\n"
    html = html + "<h1>📚 LIBRO DIARIO</h1>\n"
    html = html + "<p style=\"text-align: center;\">Total de asientos contables: " + std.len(datos.asientos) + "</p>\n"
    
    let i = 0
    while (i < std.len(datos.asientos)) {
        let asiento = datos.asientos[i]
        html = html + "<div class=\"asiento\">\n"
        html = html + "<div class=\"asiento-header\">\n"
        html = html + "<strong>Asiento No. " + asiento.id + "</strong> | "
        html = html + asiento.fecha + " | "
        html = html + asiento.descripcion + "\n"
        html = html + "</div>\n"
        
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
    
    html = html + "<div class=\"back\">\n"
    html = html + "<a href=\"/\">Volver al inicio</a>\n"
    html = html + "</div>\n"
    
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// Handler DSL
func handleDSL(pathVars, method, body) {
    let params = parseBody(body)
    let query = params.query || "reporte balance"
    
    console.log("Ejecutando DSL query: " + query)
    
    let engine = ReportesFinancieros
    let resultado = engine.use(query)
    
    let html = "<!DOCTYPE html>\n<html>\n<head>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<title>Resultado DSL</title>\n"
    html = html + "<style>\n"
    html = html + "body { font-family: Arial; margin: 20px; background: #f5f5f5; }\n"
    html = html + ".container { max-width: 800px; margin: 0 auto; background: white; padding: 40px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }\n"
    html = html + "h1 { color: #2c3e50; text-align: center; }\n"
    html = html + ".query { background: #667eea; color: white; padding: 15px; border-radius: 8px; margin: 20px 0; font-family: monospace; font-size: 16px; }\n"
    html = html + ".result { background: #f8f9fa; padding: 20px; border-radius: 8px; margin: 20px 0; }\n"
    html = html + "pre { background: #263238; color: #aed581; padding: 20px; border-radius: 8px; overflow: auto; }\n"
    html = html + ".back { text-align: center; margin-top: 30px; }\n"
    html = html + ".back a { background: #667eea; color: white; padding: 10px 20px; text-decoration: none; border-radius: 4px; }\n"
    html = html + "</style>\n"
    html = html + "</head>\n<body>\n"
    
    html = html + "<div class=\"container\">\n"
    html = html + "<h1>📊 Resultado de Consulta DSL</h1>\n"
    html = html + "<div class=\"query\">Query: " + query + "</div>\n"
    
    html = html + "<div class=\"result\">\n"
    html = html + "<h3>Resultado:</h3>\n"
    html = html + "<pre>" + JSON.stringify(resultado) + "</pre>\n"
    html = html + "</div>\n"
    
    html = html + "<div class=\"back\">\n"
    html = html + "<a href=\"/\">Volver al inicio</a>\n"
    html = html + "</div>\n"
    
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// Handler demo automática
func handleDemo(pathVars, method, body) {
    console.log("Ejecutando demo automática...")
    
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
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<title>Demo Automática</title>\n"
    html = html + "<style>\n"
    html = html + "body { font-family: Arial; margin: 20px; background: #f5f5f5; }\n"
    html = html + ".container { max-width: 800px; margin: 0 auto; background: white; padding: 40px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }\n"
    html = html + "h1 { color: #2c3e50; text-align: center; }\n"
    html = html + ".success { background: #d4edda; color: #155724; padding: 20px; border-radius: 8px; margin: 20px 0; }\n"
    html = html + ".stats { background: #f8f9fa; padding: 20px; border-radius: 8px; margin: 20px 0; }\n"
    html = html + ".stat-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 15px; }\n"
    html = html + ".stat-item { background: white; padding: 15px; border-radius: 8px; border: 1px solid #e9ecef; }\n"
    html = html + ".stat-label { color: #6c757d; font-size: 14px; }\n"
    html = html + ".stat-value { font-size: 24px; font-weight: bold; color: #495057; }\n"
    html = html + ".back { text-align: center; margin-top: 30px; }\n"
    html = html + ".back a { background: #667eea; color: white; padding: 10px 20px; text-decoration: none; border-radius: 4px; margin: 0 5px; }\n"
    html = html + "</style>\n"
    html = html + "</head>\n<body>\n"
    
    html = html + "<div class=\"container\">\n"
    html = html + "<h1>🚀 Demo Automática Completada</h1>\n"
    
    html = html + "<div class=\"success\">\n"
    html = html + "<h3>✅ Se procesaron 6 transacciones exitosamente</h3>\n"
    html = html + "<p>Se generaron transacciones de venta y compra en diferentes países de LATAM</p>\n"
    html = html + "</div>\n"
    
    html = html + "<div class=\"stats\">\n"
    html = html + "<h3>📊 Resumen del Sistema</h3>\n"
    html = html + "<div class=\"stat-grid\">\n"
    
    html = html + "<div class=\"stat-item\">\n"
    html = html + "<div class=\"stat-label\">Total Transacciones</div>\n"
    html = html + "<div class=\"stat-value\">" + std.len(datos.transacciones) + "</div>\n"
    html = html + "</div>\n"
    
    html = html + "<div class=\"stat-item\">\n"
    html = html + "<div class=\"stat-label\">Total Asientos</div>\n"
    html = html + "<div class=\"stat-value\">" + std.len(datos.asientos) + "</div>\n"
    html = html + "</div>\n"
    
    html = html + "<div class=\"stat-item\">\n"
    html = html + "<div class=\"stat-label\">Total Debe</div>\n"
    html = html + "<div class=\"stat-value\">$" + balance.totalDebe + "</div>\n"
    html = html + "</div>\n"
    
    html = html + "<div class=\"stat-item\">\n"
    html = html + "<div class=\"stat-label\">Total Haber</div>\n"
    html = html + "<div class=\"stat-value\">$" + balance.totalHaber + "</div>\n"
    html = html + "</div>\n"
    
    html = html + "</div>\n"
    
    html = html + "<p style=\"text-align: center; margin-top: 20px;\">\n"
    html = html + "<strong>Balance Cuadrado:</strong> " + balance.cuadrado + "\n"
    html = html + "</p>\n"
    
    html = html + "</div>\n"
    
    html = html + "<div class=\"back\">\n"
    html = html + "<a href=\"/\">Volver al inicio</a>\n"
    html = html + "<a href=\"/libro\">Ver Libro Diario</a>\n"
    html = html + "</div>\n"
    
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// APIs JSON
func handleAPITransacciones(pathVars, method, body) {
    return JSON({
        total: std.len(datos.transacciones),
        regiones: 7,
        transacciones: datos.transacciones
    })
}

func handleAPIAsientos(pathVars, method, body) {
    return JSON({
        total: std.len(datos.asientos),
        asientos: datos.asientos
    })
}

// Main
func main() {
    // Demo inicial en consola
    demoConsola()
    
    console.log("\n✅ Configurando servidor web...")
    
    // Configurar rutas
    http.handler("GET", "/", handleHome)
    http.handler("POST", "/procesar", handleProcesar)
    http.handler("GET", "/libro", handleLibro)
    http.handler("POST", "/dsl", handleDSL)
    http.handler("GET", "/demo", handleDemo)
    http.handler("GET", "/api/transacciones", handleAPITransacciones)
    http.handler("GET", "/api/asientos", handleAPIAsientos)
    
    console.log("🌐 Servidor listo en: http://localhost:8080")
    console.log("📌 Presiona Ctrl+C para detener\n")
    
    // Iniciar servidor
    http.serve(":8080")
}

// Ejecutar
main()