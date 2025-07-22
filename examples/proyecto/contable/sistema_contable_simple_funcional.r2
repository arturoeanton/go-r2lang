// Sistema Contable LATAM - Versión Simple Funcional
// POC para Siigo - Sin usar arrays en objetos

console.log("🌍 Sistema Contable LATAM - Simple y Funcional")
console.log("📊 Libro Diario + DSL + API")
console.log("✅ Solución alternativa sin arrays en objetos")
console.log("")

// Configuración regiones
let regiones = {
    "MX": { 
        nombre: "México", moneda: "MXN", iva: 0.16,
        clientes: "1201", ventas: "4101", ivaDebito: "2401",
        proveedores: "2101", compras: "5101", ivaCredito: "1401"
    },
    "COL": { 
        nombre: "Colombia", moneda: "COP", iva: 0.19,
        clientes: "130501", ventas: "413501", ivaDebito: "240801",
        proveedores: "220501", compras: "620501", ivaCredito: "240802"
    },
    "AR": { 
        nombre: "Argentina", moneda: "ARS", iva: 0.21,
        clientes: "1.1.2.01", ventas: "4.1.1.01", ivaDebito: "2.1.3.01",
        proveedores: "2.1.1.01", compras: "5.1.1.01", ivaCredito: "1.1.3.01"
    },
    "CH": { 
        nombre: "Chile", moneda: "CLP", iva: 0.19,
        clientes: "11030", ventas: "31010", ivaDebito: "21070",
        proveedores: "21010", compras: "41010", ivaCredito: "11070"
    },
    "UY": { 
        nombre: "Uruguay", moneda: "UYU", iva: 0.22,
        clientes: "1121", ventas: "4111", ivaDebito: "2141",
        proveedores: "2111", compras: "5111", ivaCredito: "1141"
    },
    "EC": { 
        nombre: "Ecuador", moneda: "USD", iva: 0.12,
        clientes: "102.01", ventas: "401.01", ivaDebito: "201.04",
        proveedores: "201.01", compras: "501.01", ivaCredito: "101.04"
    },
    "PE": { 
        nombre: "Perú", moneda: "PEN", iva: 0.18,
        clientes: "121", ventas: "701", ivaDebito: "401",
        proveedores: "421", compras: "601", ivaCredito: "401"
    }
}

// Arrays globales 
let transacciones = []
let asientosContables = []

// Arrays paralelos para los movimientos de cada asiento
let movimientosAsientos = []

// Función procesar transacción
func procesarTransaccion(tipo, region, importe) {
    console.log("\n[PROCESAR] Tipo: " + tipo + ", Región: " + region + ", Importe: " + importe)
    
    let config = regiones[region]
    if (!config) {
        console.log("[ERROR] Región no encontrada: " + region)
        return null
    }
    
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
    
    // Guardar transacción
    let indexTx = std.len(transacciones)
    transacciones[indexTx] = tx
    console.log("  ✅ Transacción creada: " + tx.id)
    
    // Crear asiento contable
    let asiento = {
        id: "AS-" + tx.id,
        fecha: tx.fecha,
        region: region,
        descripcion: std.toUpperCase(tipo) + " - " + config.nombre,
        indexMovimientos: std.len(movimientosAsientos)
    }
    
    // Crear movimientos como array separado
    let movimientos = []
    
    if (tipo == "ventas") {
        // Debe: Clientes
        movimientos[0] = {
            cuenta: config.clientes,
            descripcion: "Clientes",
            tipo: "DEBE",
            monto: total
        }
        // Haber: Ventas
        movimientos[1] = {
            cuenta: config.ventas,
            descripcion: "Ventas",
            tipo: "HABER",
            monto: importeNum
        }
        // Haber: IVA
        movimientos[2] = {
            cuenta: config.ivaDebito,
            descripcion: "IVA Débito",
            tipo: "HABER",
            monto: iva
        }
    } else {
        // Debe: Compras
        movimientos[0] = {
            cuenta: config.compras,
            descripcion: "Compras",
            tipo: "DEBE",
            monto: importeNum
        }
        // Debe: IVA
        movimientos[1] = {
            cuenta: config.ivaCredito,
            descripcion: "IVA Crédito",
            tipo: "DEBE",
            monto: iva
        }
        // Haber: Proveedores
        movimientos[2] = {
            cuenta: config.proveedores,
            descripcion: "Proveedores",
            tipo: "HABER",
            monto: total
        }
    }
    
    // Guardar movimientos
    movimientosAsientos[asiento.indexMovimientos] = movimientos
    
    // Guardar asiento
    let indexAs = std.len(asientosContables)
    asientosContables[indexAs] = asiento
    console.log("  ✅ Asiento creado: " + asiento.id + " con 3 movimientos")
    
    return tx
}

// Función para obtener movimientos de un asiento
func getMovimientos(asiento) {
    return movimientosAsientos[asiento.indexMovimientos]
}

// Demo en consola
func demoConsola() {
    console.log("=== DEMO SISTEMA CONTABLE ===")
    
    let tx1 = procesarTransaccion("ventas", "MX", "100000")
    console.log("Total Tx: " + std.len(transacciones) + ", Asientos: " + std.len(asientosContables))
    
    let tx2 = procesarTransaccion("compras", "COL", "50000")
    console.log("Total Tx: " + std.len(transacciones) + ", Asientos: " + std.len(asientosContables))
    
    let tx3 = procesarTransaccion("ventas", "AR", "75000")
    console.log("Total Tx: " + std.len(transacciones) + ", Asientos: " + std.len(asientosContables))
    
    // Verificar asientos
    console.log("\n📚 LIBRO DIARIO:")
    let i = 0
    while (i < std.len(asientosContables)) {
        let asiento = asientosContables[i]
        let movimientos = getMovimientos(asiento)
        
        console.log("\nAsiento: " + asiento.id + " - " + asiento.descripcion)
        
        let j = 0
        let totalDebe = 0
        let totalHaber = 0
        while (j < std.len(movimientos)) {
            let mov = movimientos[j]
            console.log("  " + mov.tipo + " - " + mov.cuenta + " - " + mov.descripcion + ": $" + mov.monto)
            if (mov.tipo == "DEBE") {
                totalDebe = totalDebe + mov.monto
            } else {
                totalHaber = totalHaber + mov.monto
            }
            j = j + 1
        }
        console.log("  Total Debe: $" + totalDebe + ", Total Haber: $" + totalHaber)
        i = i + 1
    }
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
                let value = std.replace(parts[1], "+", " ")
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
    let html = "<!DOCTYPE html><html><head>"
    html = html + "<title>Sistema Contable LATAM</title>"
    html = html + "<meta charset='UTF-8'>"
    html = html + "<style>"
    html = html + "body{font-family:Arial,sans-serif;margin:20px;background:#f5f5f5;}"
    html = html + ".container{max-width:1200px;margin:0 auto;background:white;padding:30px;border-radius:10px;box-shadow:0 2px 10px rgba(0,0,0,0.1);}"
    html = html + "h1{color:#2c3e50;text-align:center;}"
    html = html + ".grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(350px,1fr));gap:20px;}"
    html = html + ".card{background:#f8f9fa;padding:20px;border-radius:8px;border:1px solid #ddd;}"
    html = html + "input,select{width:100%;padding:8px;margin:5px 0;border:1px solid #ddd;border-radius:4px;box-sizing:border-box;}"
    html = html + "button{background:#667eea;color:white;padding:10px;border:none;border-radius:4px;width:100%;cursor:pointer;}"
    html = html + "button:hover{background:#5a67d8;}"
    html = html + ".stats{background:#e3f2fd;padding:10px;border-radius:8px;text-align:center;margin-bottom:20px;}"
    html = html + ".links{text-align:center;margin:20px 0;}"
    html = html + ".links a{color:#667eea;text-decoration:none;margin:0 10px;}"
    html = html + "</style></head><body>"
    
    html = html + "<div class='container'>"
    html = html + "<h1>🌍 Sistema Contable LATAM</h1>"
    html = html + "<div class='stats'>"
    html = html + "<strong>POC para Siigo</strong> | "
    html = html + "Transacciones: " + std.len(transacciones) + " | "
    html = html + "Asientos: " + std.len(asientosContables)
    html = html + "</div>"
    
    html = html + "<div class='grid'>"
    
    // Formulario transacción
    html = html + "<div class='card'>"
    html = html + "<h3>📝 Procesar Transacción</h3>"
    html = html + "<form action='/procesar' method='POST'>"
    html = html + "<label>Tipo:</label>"
    html = html + "<select name='tipo'>"
    html = html + "<option value='ventas'>Venta</option>"
    html = html + "<option value='compras'>Compra</option>"
    html = html + "</select>"
    html = html + "<label>Región:</label>"
    html = html + "<select name='region'>"
    html = html + "<option value='MX'>México (IVA 16%)</option>"
    html = html + "<option value='COL'>Colombia (IVA 19%)</option>"
    html = html + "<option value='AR'>Argentina (IVA 21%)</option>"
    html = html + "<option value='CH'>Chile (IVA 19%)</option>"
    html = html + "<option value='UY'>Uruguay (IVA 22%)</option>"
    html = html + "<option value='EC'>Ecuador (IVA 12%)</option>"
    html = html + "<option value='PE'>Perú (IVA 18%)</option>"
    html = html + "</select>"
    html = html + "<label>Importe (sin IVA):</label>"
    html = html + "<input type='number' name='importe' value='100000' required>"
    html = html + "<button type='submit'>Procesar</button>"
    html = html + "</form>"
    html = html + "</div>"
    
    html = html + "</div>"
    
    // Links
    html = html + "<div class='links'>"
    html = html + "<a href='/demo'>🚀 Demo Auto</a> | "
    html = html + "<a href='/libro'>📚 Libro Diario</a> | "
    html = html + "<a href='/api/transacciones'>📊 API JSON</a>"
    html = html + "</div>"
    
    html = html + "</div></body></html>"
    
    return html
}

// Handler procesar
func handleProcesar(pathVars, method, body) {
    console.log("\n[HTTP] /procesar - Body: " + body)
    
    let tipo = getParam(body, "tipo")
    let region = getParam(body, "region")
    let importe = getParam(body, "importe")
    
    if (tipo == "") { tipo = "ventas" }
    if (region == "") { region = "COL" }
    if (importe == "") { importe = "100000" }
    
    let tx = procesarTransaccion(tipo, region, importe)
    
    if (!tx) {
        return "<h1>Error al procesar transacción</h1><p><a href='/'>Volver</a></p>"
    }
    
    // Obtener último asiento
    let numAsientos = std.len(asientosContables)
    if (numAsientos == 0) {
        return "<h1>Error: No se pudo crear el asiento</h1><p><a href='/'>Volver</a></p>"
    }
    
    let asiento = asientosContables[numAsientos - 1]
    let movimientos = getMovimientos(asiento)
    
    let html = "<!DOCTYPE html><html><head>"
    html = html + "<title>Comprobante</title>"
    html = html + "<meta charset='UTF-8'>"
    html = html + "<style>"
    html = html + "body{font-family:Arial;margin:20px;background:#f5f5f5;}"
    html = html + ".comp{max-width:800px;margin:0 auto;padding:30px;background:white;border-radius:10px;box-shadow:0 2px 10px rgba(0,0,0,0.1);}"
    html = html + "h1{color:#2c3e50;text-align:center;}"
    html = html + ".info{background:#f8f9fa;padding:20px;border-radius:8px;margin:20px 0;}"
    html = html + ".asiento{background:#e3f2fd;padding:20px;border-radius:8px;margin:20px 0;}"
    html = html + "table{width:100%;border-collapse:collapse;margin-top:15px;}"
    html = html + "th,td{padding:10px;text-align:left;border-bottom:1px solid #ddd;}"
    html = html + "th{background:#667eea;color:white;}"
    html = html + ".debe{color:#28a745;font-weight:bold;}"
    html = html + ".haber{color:#dc3545;font-weight:bold;}"
    html = html + "</style></head><body>"
    
    html = html + "<div class='comp'>"
    html = html + "<h1>COMPROBANTE FISCAL</h1>"
    
    html = html + "<div class='info'>"
    html = html + "<p><strong>ID:</strong> " + tx.id + "</p>"
    html = html + "<p><strong>Fecha:</strong> " + tx.fecha + "</p>"
    html = html + "<p><strong>País:</strong> " + tx.pais + "</p>"
    html = html + "<p><strong>Tipo:</strong> " + std.toUpperCase(tx.tipo) + "</p>"
    html = html + "<p><strong>Importe:</strong> " + tx.moneda + " " + tx.importe + "</p>"
    html = html + "<p><strong>IVA:</strong> " + tx.moneda + " " + tx.iva + "</p>"
    html = html + "<p style='font-size:20px;'><strong>TOTAL:</strong> " + tx.moneda + " " + tx.total + "</p>"
    html = html + "</div>"
    
    // Asiento contable
    html = html + "<div class='asiento'>"
    html = html + "<h3>📚 ASIENTO CONTABLE</h3>"
    html = html + "<p><strong>Asiento:</strong> " + asiento.id + "</p>"
    html = html + "<p><strong>Descripción:</strong> " + asiento.descripcion + "</p>"
    
    html = html + "<table>"
    html = html + "<tr><th>Cuenta</th><th>Descripción</th><th>Debe</th><th>Haber</th></tr>"
    
    let i = 0
    let totalDebe = 0
    let totalHaber = 0
    while (i < std.len(movimientos)) {
        let mov = movimientos[i]
        html = html + "<tr>"
        html = html + "<td>" + mov.cuenta + "</td>"
        html = html + "<td>" + mov.descripcion + "</td>"
        if (mov.tipo == "DEBE") {
            html = html + "<td class='debe'>" + mov.monto + "</td>"
            html = html + "<td>-</td>"
            totalDebe = totalDebe + mov.monto
        } else {
            html = html + "<td>-</td>"
            html = html + "<td class='haber'>" + mov.monto + "</td>"
            totalHaber = totalHaber + mov.monto
        }
        html = html + "</tr>"
        i = i + 1
    }
    
    html = html + "<tr style='border-top:2px solid #667eea;'>"
    html = html + "<td colspan='2'><strong>TOTALES</strong></td>"
    html = html + "<td class='debe'><strong>" + totalDebe + "</strong></td>"
    html = html + "<td class='haber'><strong>" + totalHaber + "</strong></td>"
    html = html + "</tr>"
    
    html = html + "</table>"
    html = html + "</div>"
    
    html = html + "<p style='text-align:center;'><a href='/'>Volver</a> | <a href='/libro'>Ver Libro Diario</a></p>"
    html = html + "</div></body></html>"
    
    return html
}

// Handler libro diario
func handleLibro(pathVars, method, body) {
    let html = "<!DOCTYPE html><html><head>"
    html = html + "<title>Libro Diario</title>"
    html = html + "<meta charset='UTF-8'>"
    html = html + "<style>"
    html = html + "body{font-family:Arial;margin:20px;background:#f5f5f5;}"
    html = html + ".container{max-width:1200px;margin:0 auto;}"
    html = html + ".asiento{background:white;padding:20px;margin:15px 0;border-radius:8px;box-shadow:0 2px 5px rgba(0,0,0,0.1);}"
    html = html + "h1{color:#2c3e50;text-align:center;}"
    html = html + "table{width:100%;border-collapse:collapse;margin-top:10px;}"
    html = html + "th,td{padding:10px;text-align:left;border-bottom:1px solid #ddd;}"
    html = html + "th{background:#667eea;color:white;}"
    html = html + ".debe{color:#28a745;font-weight:bold;}"
    html = html + ".haber{color:#dc3545;font-weight:bold;}"
    html = html + "</style></head><body>"
    
    html = html + "<div class='container'>"
    html = html + "<h1>📚 LIBRO DIARIO</h1>"
    html = html + "<p style='text-align:center;'>Total asientos: " + std.len(asientosContables) + "</p>"
    
    let i = 0
    while (i < std.len(asientosContables)) {
        let asiento = asientosContables[i]
        let movimientos = getMovimientos(asiento)
        
        html = html + "<div class='asiento'>"
        html = html + "<h3>Asiento #" + (i + 1) + ": " + asiento.id + "</h3>"
        html = html + "<p><strong>Fecha:</strong> " + asiento.fecha + " | <strong>Descripción:</strong> " + asiento.descripcion + "</p>"
        
        html = html + "<table>"
        html = html + "<tr><th>Cuenta</th><th>Descripción</th><th>Debe</th><th>Haber</th></tr>"
        
        let j = 0
        let totalDebe = 0
        let totalHaber = 0
        while (j < std.len(movimientos)) {
            let mov = movimientos[j]
            html = html + "<tr>"
            html = html + "<td>" + mov.cuenta + "</td>"
            html = html + "<td>" + mov.descripcion + "</td>"
            if (mov.tipo == "DEBE") {
                html = html + "<td class='debe'>" + mov.monto + "</td>"
                html = html + "<td>-</td>"
                totalDebe = totalDebe + mov.monto
            } else {
                html = html + "<td>-</td>"
                html = html + "<td class='haber'>" + mov.monto + "</td>"
                totalHaber = totalHaber + mov.monto
            }
            html = html + "</tr>"
            j = j + 1
        }
        
        html = html + "<tr style='border-top:2px solid #667eea;'>"
        html = html + "<td colspan='2'><strong>TOTALES</strong></td>"
        html = html + "<td class='debe'><strong>" + totalDebe + "</strong></td>"
        html = html + "<td class='haber'><strong>" + totalHaber + "</strong></td>"
        html = html + "</tr>"
        
        html = html + "</table>"
        html = html + "</div>"
        i = i + 1
    }
    
    html = html + "<p style='text-align:center;'><a href='/'>Volver</a></p>"
    html = html + "</div></body></html>"
    
    return html
}

// Handler demo
func handleDemo(pathVars, method, body) {
    console.log("\n[DEMO] Ejecutando demo automática...")
    
    // Procesar varias transacciones
    procesarTransaccion("ventas", "MX", "100000")
    procesarTransaccion("compras", "COL", "50000")
    procesarTransaccion("ventas", "AR", "75000")
    procesarTransaccion("compras", "PE", "60000")
    procesarTransaccion("ventas", "CH", "120000")
    procesarTransaccion("compras", "EC", "30000")
    
    return "<html><body><h1>Demo Ejecutada</h1><p>Se procesaron 6 transacciones</p><p><a href='/libro'>Ver Libro Diario</a> | <a href='/'>Volver</a></p></body></html>"
}

// API
func handleAPI(pathVars, method, body) {
    let result = {
        transacciones: transacciones,
        asientos: asientosContables,
        totalTx: std.len(transacciones),
        totalAs: std.len(asientosContables)
    }
    return json.stringify(result)
}

// Main
func main() {
    // Demo inicial
    demoConsola()
    
    console.log("\n✅ Iniciando servidor web...")
    console.log("🔧 Usando arrays paralelos para evitar limitaciones")
    
    // Rutas
    http.handler("GET", "/", handleHome)
    http.handler("POST", "/procesar", handleProcesar)
    http.handler("GET", "/libro", handleLibro)
    http.handler("GET", "/demo", handleDemo)
    http.handler("GET", "/api/transacciones", handleAPI)
    
    console.log("🌐 http://localhost:8080")
    console.log("📌 Presiona Ctrl+C para detener\n")
    
    // Iniciar servidor
    http.serve(":8080")
}

// Ejecutar
main()