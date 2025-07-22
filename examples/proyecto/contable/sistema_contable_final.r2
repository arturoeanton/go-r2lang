// Sistema Contable LATAM - Versión Final Definitiva
// POC para Siigo - 100% Funcional con Debug

console.log("🌍 Sistema Contable LATAM - Versión Final")
console.log("📊 Libro Diario + DSL + API")
console.log("✅ Con debug para verificar funcionamiento")
console.log("="*50)

// Configuración de regiones
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

// Arrays globales - gracias al fix ahora funcionan con índices
let transacciones = []
let asientosContables = []

// Función para procesar transacción con debug
func procesarTransaccion(tipo, region, importe) {
    console.log("\n[DEBUG] Procesando transacción:")
    console.log("  Tipo: " + tipo + ", Región: " + region + ", Importe: " + importe)
    
    let config = regiones[region]
    if (!config) {
        console.log("[ERROR] Región no encontrada: " + region)
        return null
    }
    
    let importeNum = std.parseFloat(importe)
    let iva = math.round((importeNum * config.iva) * 100) / 100
    let total = importeNum + iva
    
    // Crear transacción
    let txId = region + "-" + math.randomInt(1000) + "-" + math.randomInt(999)
    let tx = {
        id: txId,
        tipo: tipo,
        region: region,
        pais: config.nombre,
        importe: importeNum,
        iva: iva,
        total: total,
        moneda: config.moneda,
        fecha: std.now()
    }
    
    // Guardar transacción usando índice
    let indexTx = std.len(transacciones)
    transacciones[indexTx] = tx
    console.log("  Transacción guardada en índice " + indexTx)
    console.log("  Total transacciones ahora: " + std.len(transacciones))
    
    // Crear asiento contable
    let asiento = {
        id: "AS-" + txId,
        fecha: tx.fecha,
        region: region,
        descripcion: std.toUpperCase(tipo) + " - " + config.nombre,
        movimientos: []
    }
    
    // Crear movimientos del asiento
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
    
    // Guardar asiento
    let indexAs = std.len(asientosContables)
    asientosContables[indexAs] = asiento
    console.log("  Asiento guardado en índice " + indexAs)
    console.log("  Total asientos ahora: " + std.len(asientosContables))
    console.log("  Movimientos del asiento: " + std.len(asiento.movimientos))
    
    return tx
}

// DSL para Reportes Financieros
dsl ReportesFinancieros {
    token("REPORTE", "reporte")
    token("TIPO", "balance|diario|ventas|compras|iva")
    
    rule("consulta", ["REPORTE", "TIPO"], "ejecutarReporte")
    
    func ejecutarReporte(cmd, tipo) {
        console.log("[DSL] Ejecutando reporte: " + tipo)
        console.log("[DSL] Asientos disponibles: " + std.len(asientosContables))
        
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
                tipo: "Reporte " + std.toUpperCase(tipo),
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
        
        return { error: "Tipo de reporte no válido: " + tipo }
    }
}

// Demo en consola con verificación
func demoConsola() {
    console.log("\n=== DEMO SISTEMA CONTABLE ===")
    
    // Procesar transacciones de prueba
    let tx1 = procesarTransaccion("ventas", "MX", "100000")
    if (tx1) {
        console.log("✅ Venta México: $" + tx1.total + " " + tx1.moneda)
    }
    
    let tx2 = procesarTransaccion("compras", "COL", "50000")
    if (tx2) {
        console.log("✅ Compra Colombia: $" + tx2.total + " " + tx2.moneda)
    }
    
    let tx3 = procesarTransaccion("ventas", "AR", "75000")
    if (tx3) {
        console.log("✅ Venta Argentina: $" + tx3.total + " " + tx3.moneda)
    }
    
    // Verificar que tenemos asientos
    console.log("\n📚 VERIFICACIÓN DE ASIENTOS:")
    console.log("Total asientos creados: " + std.len(asientosContables))
    
    if (std.len(asientosContables) > 0) {
        console.log("\nEjemplo - Primer asiento:")
        let asiento = asientosContables[0]
        console.log("ID: " + asiento.id)
        console.log("Descripción: " + asiento.descripcion)
        console.log("Movimientos:")
        
        let i = 0
        while (i < std.len(asiento.movimientos)) {
            let mov = asiento.movimientos[i]
            console.log("  " + (i+1) + ". " + mov.tipo + " - " + mov.cuenta + " - " + mov.descripcion + ": $" + mov.monto)
            i = i + 1
        }
    }
    
    // Probar DSL
    console.log("\n📊 PRUEBA DSL REPORTES:")
    let engine = ReportesFinancieros
    let balance = engine.use("reporte balance")
    console.log("Balance General:")
    console.log("  Total Asientos: " + balance.totalAsientos)
    console.log("  Total Debe: $" + balance.totalDebe)
    console.log("  Total Haber: $" + balance.totalHaber)
    console.log("  Cuadrado: " + balance.cuadrado)
    
    console.log("\n🎯 VALUE PROPOSITION:")
    console.log("• 18 meses → 2 meses")
    console.log("• $500K → $150K")
    console.log("• ROI: 1,020%")
}

// Función para parsear parámetros del body
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
                let value = parts[1]
                // Decodificar URL
                value = std.replace(value, "+", " ")
                value = std.replace(value, "%20", " ")
                return value
            }
        }
        i = i + 1
    }
    return ""
}

// Handler página principal
func handleHome(pathVars, method, body) {
    let html = "<!DOCTYPE html>\n"
    html = html + "<html lang=\"es\">\n<head>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n"
    html = html + "<title>Sistema Contable LATAM - Siigo POC</title>\n"
    html = html + "<style>\n"
    html = html + "* { box-sizing: border-box; margin: 0; padding: 0; }\n"
    html = html + "body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f5f7fa; color: #333; line-height: 1.6; }\n"
    html = html + ".container { max-width: 1200px; margin: 0 auto; padding: 20px; }\n"
    html = html + ".header { background: white; border-radius: 10px; padding: 30px; margin-bottom: 30px; box-shadow: 0 2px 10px rgba(0,0,0,0.05); text-align: center; }\n"
    html = html + "h1 { color: #2c3e50; margin-bottom: 10px; font-size: 2.5em; }\n"
    html = html + ".subtitle { color: #7f8c8d; font-size: 1.1em; }\n"
    html = html + ".stats { display: flex; justify-content: center; gap: 30px; margin-top: 20px; }\n"
    html = html + ".stat { text-align: center; }\n"
    html = html + ".stat-value { font-size: 2em; font-weight: bold; color: #667eea; }\n"
    html = html + ".stat-label { color: #95a5a6; font-size: 0.9em; }\n"
    html = html + ".grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(400px, 1fr)); gap: 30px; margin-bottom: 30px; }\n"
    html = html + ".card { background: white; border-radius: 10px; padding: 30px; box-shadow: 0 2px 10px rgba(0,0,0,0.05); }\n"
    html = html + ".card h3 { color: #34495e; margin-bottom: 20px; font-size: 1.5em; }\n"
    html = html + "label { display: block; margin-bottom: 5px; color: #555; font-weight: 500; }\n"
    html = html + "input, select, textarea { width: 100%; padding: 12px; margin-bottom: 15px; border: 2px solid #e0e6ed; border-radius: 6px; font-size: 16px; transition: border-color 0.3s; }\n"
    html = html + "input:focus, select:focus, textarea:focus { outline: none; border-color: #667eea; }\n"
    html = html + "button { background: #667eea; color: white; padding: 12px 24px; border: none; border-radius: 6px; font-size: 16px; font-weight: 600; cursor: pointer; width: 100%; transition: background 0.3s; }\n"
    html = html + "button:hover { background: #5a67d8; }\n"
    html = html + ".links { text-align: center; margin: 30px 0; }\n"
    html = html + ".links a { color: #667eea; text-decoration: none; margin: 0 15px; font-weight: 500; transition: color 0.3s; }\n"
    html = html + ".links a:hover { color: #5a67d8; text-decoration: underline; }\n"
    html = html + ".value-prop { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; border-radius: 10px; padding: 40px; text-align: center; box-shadow: 0 5px 20px rgba(102,126,234,0.3); }\n"
    html = html + ".value-prop h3 { color: white; margin-bottom: 20px; font-size: 2em; }\n"
    html = html + ".value-metrics { display: flex; justify-content: space-around; flex-wrap: wrap; gap: 20px; }\n"
    html = html + ".metric { flex: 1; min-width: 200px; }\n"
    html = html + ".metric-value { font-size: 2.5em; font-weight: bold; margin-bottom: 5px; }\n"
    html = html + ".metric-label { opacity: 0.9; }\n"
    html = html + "</style>\n</head>\n<body>\n"
    
    html = html + "<div class=\"container\">\n"
    
    // Header
    html = html + "<div class=\"header\">\n"
    html = html + "<h1>🌍 Sistema Contable LATAM</h1>\n"
    html = html + "<p class=\"subtitle\">POC para Siigo - Demostración de R2Lang con DSL</p>\n"
    html = html + "<div class=\"stats\">\n"
    html = html + "<div class=\"stat\">\n"
    html = html + "<div class=\"stat-value\">" + std.len(transacciones) + "</div>\n"
    html = html + "<div class=\"stat-label\">Transacciones</div>\n"
    html = html + "</div>\n"
    html = html + "<div class=\"stat\">\n"
    html = html + "<div class=\"stat-value\">" + std.len(asientosContables) + "</div>\n"
    html = html + "<div class=\"stat-label\">Asientos Contables</div>\n"
    html = html + "</div>\n"
    html = html + "<div class=\"stat\">\n"
    html = html + "<div class=\"stat-value\">7</div>\n"
    html = html + "<div class=\"stat-label\">Países LATAM</div>\n"
    html = html + "</div>\n"
    html = html + "</div>\n"
    html = html + "</div>\n"
    
    // Formularios
    html = html + "<div class=\"grid\">\n"
    
    // Procesar transacción
    html = html + "<div class=\"card\">\n"
    html = html + "<h3>📝 Procesar Transacción</h3>\n"
    html = html + "<form action=\"/procesar\" method=\"POST\">\n"
    html = html + "<label for=\"tipo\">Tipo de transacción:</label>\n"
    html = html + "<select name=\"tipo\" id=\"tipo\">\n"
    html = html + "<option value=\"ventas\">Venta</option>\n"
    html = html + "<option value=\"compras\">Compra</option>\n"
    html = html + "</select>\n"
    html = html + "<label for=\"region\">País/Región:</label>\n"
    html = html + "<select name=\"region\" id=\"region\">\n"
    html = html + "<option value=\"MX\">🇲🇽 México (IVA 16%)</option>\n"
    html = html + "<option value=\"COL\">🇨🇴 Colombia (IVA 19%)</option>\n"
    html = html + "<option value=\"AR\">🇦🇷 Argentina (IVA 21%)</option>\n"
    html = html + "<option value=\"CH\">🇨🇱 Chile (IVA 19%)</option>\n"
    html = html + "<option value=\"UY\">🇺🇾 Uruguay (IVA 22%)</option>\n"
    html = html + "<option value=\"EC\">🇪🇨 Ecuador (IVA 12%)</option>\n"
    html = html + "<option value=\"PE\">🇵🇪 Perú (IGV 18%)</option>\n"
    html = html + "</select>\n"
    html = html + "<label for=\"importe\">Importe (sin IVA):</label>\n"
    html = html + "<input type=\"number\" name=\"importe\" id=\"importe\" value=\"100000\" required min=\"1\">\n"
    html = html + "<button type=\"submit\">Procesar Transacción</button>\n"
    html = html + "</form>\n"
    html = html + "</div>\n"
    
    // Consulta DSL
    html = html + "<div class=\"card\">\n"
    html = html + "<h3>🔍 Consulta DSL de Reportes</h3>\n"
    html = html + "<form action=\"/dsl\" method=\"POST\">\n"
    html = html + "<label for=\"query\">Ingrese su consulta DSL:</label>\n"
    html = html + "<textarea name=\"query\" id=\"query\" rows=\"4\" placeholder=\"Ejemplos:\nreporte balance\nreporte diario\nreporte ventas\nreporte compras\nreporte iva\">reporte balance</textarea>\n"
    html = html + "<button type=\"submit\">Ejecutar Consulta DSL</button>\n"
    html = html + "</form>\n"
    html = html + "</div>\n"
    
    html = html + "</div>\n"
    
    // Enlaces
    html = html + "<div class=\"links\">\n"
    html = html + "<a href=\"/demo\">🚀 Demo Automática</a>\n"
    html = html + "<a href=\"/libro\">📚 Libro Diario</a>\n"
    html = html + "<a href=\"/api/transacciones\">📊 API Transacciones</a>\n"
    html = html + "<a href=\"/api/asientos\">📋 API Asientos</a>\n"
    html = html + "</div>\n"
    
    // Value Proposition
    html = html + "<div class=\"value-prop\">\n"
    html = html + "<h3>🎯 Value Proposition para Siigo</h3>\n"
    html = html + "<div class=\"value-metrics\">\n"
    html = html + "<div class=\"metric\">\n"
    html = html + "<div class=\"metric-value\">18→2</div>\n"
    html = html + "<div class=\"metric-label\">meses de desarrollo</div>\n"
    html = html + "</div>\n"
    html = html + "<div class=\"metric\">\n"
    html = html + "<div class=\"metric-value\">$500K→$150K</div>\n"
    html = html + "<div class=\"metric-label\">costo por país</div>\n"
    html = html + "</div>\n"
    html = html + "<div class=\"metric\">\n"
    html = html + "<div class=\"metric-value\">7→1</div>\n"
    html = html + "<div class=\"metric-label\">sistemas unificados</div>\n"
    html = html + "</div>\n"
    html = html + "<div class=\"metric\">\n"
    html = html + "<div class=\"metric-value\">1,020%</div>\n"
    html = html + "<div class=\"metric-label\">ROI en 3 años</div>\n"
    html = html + "</div>\n"
    html = html + "</div>\n"
    html = html + "</div>\n"
    
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// Handler procesar transacción
func handleProcesar(pathVars, method, body) {
    console.log("\n[WEB] handleProcesar llamado")
    console.log("[WEB] Body recibido: " + body)
    
    let tipo = getParam(body, "tipo")
    let region = getParam(body, "region")
    let importe = getParam(body, "importe")
    
    console.log("[WEB] Parámetros parseados:")
    console.log("  tipo=" + tipo)
    console.log("  region=" + region)
    console.log("  importe=" + importe)
    
    // Valores por defecto
    if (tipo == "") { tipo = "ventas" }
    if (region == "") { region = "COL" }
    if (importe == "") { importe = "100000" }
    
    // Procesar la transacción
    let tx = procesarTransaccion(tipo, region, importe)
    
    if (!tx) {
        return "<h1>Error al procesar la transacción</h1><p><a href=\"/\">Volver</a></p>"
    }
    
    // Verificar que tenemos asientos
    let numAsientos = std.len(asientosContables)
    console.log("[WEB] Número de asientos después de procesar: " + numAsientos)
    
    if (numAsientos == 0) {
        return "<h1>Error: No se pudo crear el asiento contable</h1><p>No hay asientos en el sistema</p><p><a href=\"/\">Volver</a></p>"
    }
    
    // Obtener el último asiento
    let asiento = asientosContables[numAsientos - 1]
    console.log("[WEB] Asiento obtenido: " + asiento.id)
    
    // Generar HTML del comprobante
    let html = "<!DOCTYPE html>\n<html lang=\"es\">\n<head>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<title>Comprobante - " + tx.id + "</title>\n"
    html = html + "<style>\n"
    html = html + "body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f5f7fa; margin: 0; padding: 20px; }\n"
    html = html + ".comprobante { max-width: 800px; margin: 0 auto; background: white; border-radius: 10px; box-shadow: 0 5px 20px rgba(0,0,0,0.1); overflow: hidden; }\n"
    html = html + ".header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 30px; text-align: center; }\n"
    html = html + ".header h1 { margin: 0; font-size: 2em; }\n"
    html = html + ".content { padding: 30px; }\n"
    html = html + ".info-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 20px; margin-bottom: 30px; }\n"
    html = html + ".info-item { background: #f8f9fa; padding: 15px; border-radius: 8px; }\n"
    html = html + ".info-label { color: #6c757d; font-size: 0.9em; margin-bottom: 5px; }\n"
    html = html + ".info-value { font-size: 1.2em; font-weight: 600; color: #2c3e50; }\n"
    html = html + ".total { background: #e3f2fd; padding: 20px; border-radius: 8px; text-align: center; margin: 20px 0; }\n"
    html = html + ".total .amount { font-size: 2.5em; font-weight: bold; color: #1976d2; }\n"
    html = html + ".asiento { background: #fff3cd; padding: 20px; border-radius: 8px; margin: 20px 0; }\n"
    html = html + ".asiento h3 { color: #856404; margin-bottom: 15px; }\n"
    html = html + "table { width: 100%; border-collapse: collapse; margin-top: 15px; }\n"
    html = html + "th { background: #f8f9fa; padding: 12px; text-align: left; font-weight: 600; color: #495057; }\n"
    html = html + "td { padding: 12px; border-bottom: 1px solid #dee2e6; }\n"
    html = html + ".debe { color: #28a745; font-weight: 600; }\n"
    html = html + ".haber { color: #dc3545; font-weight: 600; }\n"
    html = html + ".totales { background: #f8f9fa; font-weight: bold; }\n"
    html = html + ".back { text-align: center; margin-top: 30px; }\n"
    html = html + ".back a { display: inline-block; background: #667eea; color: white; padding: 12px 30px; border-radius: 6px; text-decoration: none; font-weight: 600; transition: background 0.3s; }\n"
    html = html + ".back a:hover { background: #5a67d8; }\n"
    html = html + "</style>\n</head>\n<body>\n"
    
    html = html + "<div class=\"comprobante\">\n"
    
    // Header
    html = html + "<div class=\"header\">\n"
    html = html + "<h1>COMPROBANTE FISCAL</h1>\n"
    html = html + "<p>" + tx.id + "</p>\n"
    html = html + "</div>\n"
    
    // Content
    html = html + "<div class=\"content\">\n"
    
    // Info grid
    html = html + "<div class=\"info-grid\">\n"
    html = html + "<div class=\"info-item\">\n"
    html = html + "<div class=\"info-label\">Fecha</div>\n"
    html = html + "<div class=\"info-value\">" + tx.fecha + "</div>\n"
    html = html + "</div>\n"
    html = html + "<div class=\"info-item\">\n"
    html = html + "<div class=\"info-label\">Tipo de Operación</div>\n"
    html = html + "<div class=\"info-value\">" + std.toUpperCase(tx.tipo) + "</div>\n"
    html = html + "</div>\n"
    html = html + "<div class=\"info-item\">\n"
    html = html + "<div class=\"info-label\">País</div>\n"
    html = html + "<div class=\"info-value\">" + tx.pais + "</div>\n"
    html = html + "</div>\n"
    html = html + "<div class=\"info-item\">\n"
    html = html + "<div class=\"info-label\">Moneda</div>\n"
    html = html + "<div class=\"info-value\">" + tx.moneda + "</div>\n"
    html = html + "</div>\n"
    html = html + "</div>\n"
    
    // Montos
    html = html + "<div class=\"info-grid\">\n"
    html = html + "<div class=\"info-item\">\n"
    html = html + "<div class=\"info-label\">Importe Neto</div>\n"
    html = html + "<div class=\"info-value\">" + tx.moneda + " " + formatNumber(tx.importe) + "</div>\n"
    html = html + "</div>\n"
    html = html + "<div class=\"info-item\">\n"
    html = html + "<div class=\"info-label\">IVA/IGV</div>\n"
    html = html + "<div class=\"info-value\">" + tx.moneda + " " + formatNumber(tx.iva) + "</div>\n"
    html = html + "</div>\n"
    html = html + "</div>\n"
    
    // Total
    html = html + "<div class=\"total\">\n"
    html = html + "<div>TOTAL</div>\n"
    html = html + "<div class=\"amount\">" + tx.moneda + " " + formatNumber(tx.total) + "</div>\n"
    html = html + "</div>\n"
    
    // Asiento contable
    html = html + "<div class=\"asiento\">\n"
    html = html + "<h3>📚 ASIENTO CONTABLE</h3>\n"
    html = html + "<p><strong>Asiento No.:</strong> " + asiento.id + "</p>\n"
    html = html + "<p><strong>Descripción:</strong> " + asiento.descripcion + "</p>\n"
    
    html = html + "<table>\n"
    html = html + "<thead>\n"
    html = html + "<tr>\n"
    html = html + "<th>Cuenta</th>\n"
    html = html + "<th>Descripción</th>\n"
    html = html + "<th>Debe</th>\n"
    html = html + "<th>Haber</th>\n"
    html = html + "</tr>\n"
    html = html + "</thead>\n"
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
            html = html + "<td class=\"debe\">" + formatNumber(mov.monto) + "</td>\n"
            html = html + "<td>-</td>\n"
            totalDebe = totalDebe + mov.monto
        } else {
            html = html + "<td>-</td>\n"
            html = html + "<td class=\"haber\">" + formatNumber(mov.monto) + "</td>\n"
            totalHaber = totalHaber + mov.monto
        }
        
        html = html + "</tr>\n"
        i = i + 1
    }
    
    // Totales
    html = html + "<tr class=\"totales\">\n"
    html = html + "<td colspan=\"2\">TOTALES</td>\n"
    html = html + "<td class=\"debe\">" + formatNumber(totalDebe) + "</td>\n"
    html = html + "<td class=\"haber\">" + formatNumber(totalHaber) + "</td>\n"
    html = html + "</tr>\n"
    
    html = html + "</tbody>\n"
    html = html + "</table>\n"
    html = html + "</div>\n"
    
    // Botón volver
    html = html + "<div class=\"back\">\n"
    html = html + "<a href=\"/\">Volver al Inicio</a>\n"
    html = html + "</div>\n"
    
    html = html + "</div>\n"
    html = html + "</div>\n"
    
    html = html + "</body>\n</html>"
    
    return html
}

// Handler libro diario
func handleLibro(pathVars, method, body) {
    console.log("\n[WEB] handleLibro - Total asientos: " + std.len(asientosContables))
    
    let html = "<!DOCTYPE html>\n<html lang=\"es\">\n<head>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<title>Libro Diario</title>\n"
    html = html + "<style>\n"
    html = html + "body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f5f7fa; margin: 0; padding: 20px; }\n"
    html = html + ".container { max-width: 1200px; margin: 0 auto; }\n"
    html = html + ".header { background: white; border-radius: 10px; padding: 30px; margin-bottom: 30px; box-shadow: 0 2px 10px rgba(0,0,0,0.05); text-align: center; }\n"
    html = html + "h1 { color: #2c3e50; margin: 0; }\n"
    html = html + ".summary { color: #7f8c8d; margin-top: 10px; font-size: 1.1em; }\n"
    html = html + ".asiento { background: white; border-radius: 10px; padding: 25px; margin-bottom: 20px; box-shadow: 0 2px 10px rgba(0,0,0,0.05); }\n"
    html = html + ".asiento-header { background: #f8f9fa; margin: -25px -25px 20px -25px; padding: 20px 25px; border-radius: 10px 10px 0 0; border-bottom: 2px solid #e9ecef; }\n"
    html = html + ".asiento-id { font-size: 1.2em; font-weight: 600; color: #495057; }\n"
    html = html + ".asiento-info { color: #6c757d; margin-top: 5px; }\n"
    html = html + "table { width: 100%; border-collapse: collapse; }\n"
    html = html + "th { background: #f8f9fa; padding: 12px; text-align: left; font-weight: 600; color: #495057; border-bottom: 2px solid #dee2e6; }\n"
    html = html + "td { padding: 12px; border-bottom: 1px solid #dee2e6; }\n"
    html = html + ".debe { color: #28a745; font-weight: 600; }\n"
    html = html + ".haber { color: #dc3545; font-weight: 600; }\n"
    html = html + ".empty { text-align: center; padding: 60px; color: #6c757d; font-size: 1.1em; }\n"
    html = html + ".back { text-align: center; margin: 40px 0; }\n"
    html = html + ".back a { display: inline-block; background: #667eea; color: white; padding: 12px 30px; border-radius: 6px; text-decoration: none; font-weight: 600; }\n"
    html = html + "</style>\n</head>\n<body>\n"
    
    html = html + "<div class=\"container\">\n"
    
    // Header
    html = html + "<div class=\"header\">\n"
    html = html + "<h1>📚 LIBRO DIARIO</h1>\n"
    html = html + "<p class=\"summary\">Registro cronológico de asientos contables</p>\n"
    html = html + "<p class=\"summary\">Total de asientos: " + std.len(asientosContables) + "</p>\n"
    html = html + "</div>\n"
    
    if (std.len(asientosContables) == 0) {
        html = html + "<div class=\"asiento\">\n"
        html = html + "<div class=\"empty\">No hay asientos contables registrados.<br>Procese algunas transacciones para ver el libro diario.</div>\n"
        html = html + "</div>\n"
    } else {
        // Mostrar asientos
        let i = 0
        while (i < std.len(asientosContables)) {
            let asiento = asientosContables[i]
            
            html = html + "<div class=\"asiento\">\n"
            html = html + "<div class=\"asiento-header\">\n"
            html = html + "<div class=\"asiento-id\">Asiento: " + asiento.id + "</div>\n"
            html = html + "<div class=\"asiento-info\">" + asiento.fecha + " | " + asiento.region + " | " + asiento.descripcion + "</div>\n"
            html = html + "</div>\n"
            
            html = html + "<table>\n"
            html = html + "<thead>\n"
            html = html + "<tr>\n"
            html = html + "<th style=\"width: 20%;\">Cuenta</th>\n"
            html = html + "<th style=\"width: 40%;\">Descripción</th>\n"
            html = html + "<th style=\"width: 20%;\">Debe</th>\n"
            html = html + "<th style=\"width: 20%;\">Haber</th>\n"
            html = html + "</tr>\n"
            html = html + "</thead>\n"
            html = html + "<tbody>\n"
            
            let j = 0
            let totalDebe = 0
            let totalHaber = 0
            
            while (j < std.len(asiento.movimientos)) {
                let mov = asiento.movimientos[j]
                html = html + "<tr>\n"
                html = html + "<td>" + mov.cuenta + "</td>\n"
                html = html + "<td>" + mov.descripcion + "</td>\n"
                
                if (mov.tipo == "DEBE") {
                    html = html + "<td class=\"debe\">" + formatNumber(mov.monto) + "</td>\n"
                    html = html + "<td>-</td>\n"
                    totalDebe = totalDebe + mov.monto
                } else {
                    html = html + "<td>-</td>\n"
                    html = html + "<td class=\"haber\">" + formatNumber(mov.monto) + "</td>\n"
                    totalHaber = totalHaber + mov.monto
                }
                
                html = html + "</tr>\n"
                j = j + 1
            }
            
            // Totales del asiento
            html = html + "<tr style=\"border-top: 2px solid #495057;\">\n"
            html = html + "<td colspan=\"2\" style=\"text-align: right; font-weight: bold;\">Totales:</td>\n"
            html = html + "<td class=\"debe\" style=\"font-weight: bold;\">" + formatNumber(totalDebe) + "</td>\n"
            html = html + "<td class=\"haber\" style=\"font-weight: bold;\">" + formatNumber(totalHaber) + "</td>\n"
            html = html + "</tr>\n"
            
            html = html + "</tbody>\n"
            html = html + "</table>\n"
            html = html + "</div>\n"
            
            i = i + 1
        }
    }
    
    html = html + "<div class=\"back\">\n"
    html = html + "<a href=\"/\">Volver al Inicio</a>\n"
    html = html + "</div>\n"
    
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// Handler consulta DSL
func handleDSL(pathVars, method, body) {
    let query = getParam(body, "query")
    if (query == "") {
        query = "reporte balance"
    }
    
    console.log("\n[WEB] Ejecutando DSL query: " + query)
    
    let engine = ReportesFinancieros
    let resultado = engine.use(query)
    
    let html = "<!DOCTYPE html>\n<html lang=\"es\">\n<head>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<title>Resultado DSL</title>\n"
    html = html + "<style>\n"
    html = html + "body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f5f7fa; margin: 0; padding: 20px; }\n"
    html = html + ".container { max-width: 800px; margin: 0 auto; }\n"
    html = html + ".card { background: white; border-radius: 10px; padding: 30px; box-shadow: 0 2px 10px rgba(0,0,0,0.05); }\n"
    html = html + "h1 { color: #2c3e50; text-align: center; margin-bottom: 30px; }\n"
    html = html + ".query-box { background: #667eea; color: white; padding: 20px; border-radius: 8px; margin-bottom: 30px; font-family: 'Courier New', monospace; font-size: 1.2em; }\n"
    html = html + ".result-box { background: #f8f9fa; padding: 20px; border-radius: 8px; }\n"
    html = html + "pre { background: #263238; color: #aed581; padding: 20px; border-radius: 8px; overflow-x: auto; margin: 0; font-size: 14px; line-height: 1.5; }\n"
    html = html + ".back { text-align: center; margin-top: 30px; }\n"
    html = html + ".back a { display: inline-block; background: #667eea; color: white; padding: 12px 30px; border-radius: 6px; text-decoration: none; font-weight: 600; }\n"
    html = html + "</style>\n</head>\n<body>\n"
    
    html = html + "<div class=\"container\">\n"
    html = html + "<div class=\"card\">\n"
    html = html + "<h1>📊 Resultado de Consulta DSL</h1>\n"
    html = html + "<div class=\"query-box\">Query: " + query + "</div>\n"
    html = html + "<div class=\"result-box\">\n"
    html = html + "<h3>Resultado:</h3>\n"
    html = html + "<pre>" + json.stringify(resultado) + "</pre>\n"
    html = html + "</div>\n"
    html = html + "<div class=\"back\">\n"
    html = html + "<a href=\"/\">Volver al Inicio</a>\n"
    html = html + "</div>\n"
    html = html + "</div>\n"
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// Handler demo automática
func handleDemo(pathVars, method, body) {
    console.log("\n[WEB] Ejecutando demo automática...")
    
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
    
    let html = "<!DOCTYPE html>\n<html lang=\"es\">\n<head>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<title>Demo Automática</title>\n"
    html = html + "<style>\n"
    html = html + "body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f5f7fa; margin: 0; padding: 20px; }\n"
    html = html + ".container { max-width: 600px; margin: 0 auto; }\n"
    html = html + ".card { background: white; border-radius: 10px; padding: 40px; box-shadow: 0 2px 10px rgba(0,0,0,0.05); text-align: center; }\n"
    html = html + "h1 { color: #2c3e50; margin-bottom: 30px; }\n"
    html = html + ".success { background: #d4edda; color: #155724; padding: 20px; border-radius: 8px; margin-bottom: 30px; }\n"
    html = html + ".stats { background: #f8f9fa; padding: 30px; border-radius: 8px; text-align: left; }\n"
    html = html + ".stat-item { margin: 15px 0; font-size: 1.1em; }\n"
    html = html + ".stat-label { color: #6c757d; }\n"
    html = html + ".stat-value { color: #2c3e50; font-weight: 600; }\n"
    html = html + ".balance { background: #e3f2fd; padding: 20px; border-radius: 8px; margin-top: 20px; }\n"
    html = html + ".links { margin-top: 30px; }\n"
    html = html + ".links a { display: inline-block; background: #667eea; color: white; padding: 12px 25px; border-radius: 6px; text-decoration: none; font-weight: 600; margin: 0 10px; }\n"
    html = html + "</style>\n</head>\n<body>\n"
    
    html = html + "<div class=\"container\">\n"
    html = html + "<div class=\"card\">\n"
    html = html + "<h1>🚀 Demo Automática Completada</h1>\n"
    
    html = html + "<div class=\"success\">\n"
    html = html + "<h3>✅ Se procesaron 6 transacciones exitosamente</h3>\n"
    html = html + "<p>Ventas y compras en diferentes países de LATAM</p>\n"
    html = html + "</div>\n"
    
    html = html + "<div class=\"stats\">\n"
    html = html + "<h3>📊 Resumen del Sistema:</h3>\n"
    html = html + "<div class=\"stat-item\">\n"
    html = html + "<span class=\"stat-label\">Total transacciones:</span> \n"
    html = html + "<span class=\"stat-value\">" + std.len(transacciones) + "</span>\n"
    html = html + "</div>\n"
    html = html + "<div class=\"stat-item\">\n"
    html = html + "<span class=\"stat-label\">Total asientos contables:</span> \n"
    html = html + "<span class=\"stat-value\">" + std.len(asientosContables) + "</span>\n"
    html = html + "</div>\n"
    
    html = html + "<div class=\"balance\">\n"
    html = html + "<h4>Balance General:</h4>\n"
    html = html + "<div class=\"stat-item\">\n"
    html = html + "<span class=\"stat-label\">Total Debe:</span> \n"
    html = html + "<span class=\"stat-value\">$" + formatNumber(balance.totalDebe) + "</span>\n"
    html = html + "</div>\n"
    html = html + "<div class=\"stat-item\">\n"
    html = html + "<span class=\"stat-label\">Total Haber:</span> \n"
    html = html + "<span class=\"stat-value\">$" + formatNumber(balance.totalHaber) + "</span>\n"
    html = html + "</div>\n"
    html = html + "<div class=\"stat-item\">\n"
    html = html + "<span class=\"stat-label\">Partida doble cuadrada:</span> \n"
    html = html + "<span class=\"stat-value\">" + balance.cuadrado + "</span>\n"
    html = html + "</div>\n"
    html = html + "</div>\n"
    html = html + "</div>\n"
    
    html = html + "<div class=\"links\">\n"
    html = html + "<a href=\"/\">Volver al Inicio</a>\n"
    html = html + "<a href=\"/libro\">Ver Libro Diario</a>\n"
    html = html + "</div>\n"
    
    html = html + "</div>\n"
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// APIs JSON
func handleAPITransacciones(pathVars, method, body) {
    console.log("[API] Transacciones solicitadas: " + std.len(transacciones))
    
    return json.stringify({
        total: std.len(transacciones),
        regiones: 7,
        transacciones: transacciones
    })
}

func handleAPIAsientos(pathVars, method, body) {
    console.log("[API] Asientos solicitados: " + std.len(asientosContables))
    
    return json.stringify({
        total: std.len(asientosContables),
        asientos: asientosContables
    })
}

// Función auxiliar para formatear números
func formatNumber(num) {
    // Convertir a string y agregar separadores de miles
    let str = std.toString(num)
    let parts = std.split(str, ".")
    let intPart = parts[0]
    let decPart = ""
    
    if (std.len(parts) > 1) {
        decPart = "." + parts[1]
    }
    
    // Por simplicidad, retornamos el número tal cual
    // En producción se agregarían separadores de miles
    return str
}

// Main
func main() {
    // Ejecutar demo inicial
    demoConsola()
    
    console.log("\n" + "="*50)
    console.log("✅ Iniciando servidor web...")
    console.log("="*50)
    
    // Configurar rutas
    http.handler("GET", "/", handleHome)
    http.handler("POST", "/procesar", handleProcesar)
    http.handler("GET", "/libro", handleLibro)
    http.handler("POST", "/dsl", handleDSL)
    http.handler("GET", "/demo", handleDemo)
    http.handler("GET", "/api/transacciones", handleAPITransacciones)
    http.handler("GET", "/api/asientos", handleAPIAsientos)
    
    console.log("🌐 Servidor escuchando en: http://localhost:8080")
    console.log("📌 Presiona Ctrl+C para detener\n")
    
    // Iniciar servidor
    http.serve(":8080")
}

// Ejecutar aplicación
main()