// Sistema Contable LATAM - Version Funcionando
// Demo para Siigo

console.log("🌍 Sistema Contable LATAM")
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

// Base de datos - Usar un objeto con contador
let database = {
    transacciones: [],
    asientos: [],
    countTx: 0,
    countAs: 0
}

// Función agregar transacción
func addTransaccion(tx) {
    database.transacciones[database.countTx] = tx
    database.countTx = database.countTx + 1
}

// Función agregar asiento
func addAsiento(asiento) {
    database.asientos[database.countAs] = asiento
    database.countAs = database.countAs + 1
}

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
        moneda: config.moneda,
        fecha: std.now()
    }
    
    addTransaccion(tx)
    
    // Crear asiento
    let asiento = {
        id: "AS-" + tx.id,
        fecha: tx.fecha,
        region: region,
        descripcion: std.toUpperCase(tipo) + " - " + config.nombre,
        movimientos: []
    }
    
    if (tipo == "ventas") {
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
    
    addAsiento(asiento)
    
    return tx
}

// DSL Reportes
dsl ReportesFinancieros {
    token("REPORTE", "reporte")
    token("TIPO", "balance|diario|ventas|compras|iva")
    
    rule("query", ["REPORTE", "TIPO"], "ejecutar")
    
    func ejecutar(cmd, tipo) {
        if (tipo == "balance") {
            let totalDebe = 0
            let totalHaber = 0
            let i = 0
            while (i < database.countAs) {
                let asiento = database.asientos[i]
                let j = 0
                while (j < 3) {
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
                totalAsientos: database.countAs,
                totalDebe: math.round(totalDebe * 100) / 100,
                totalHaber: math.round(totalHaber * 100) / 100,
                cuadrado: math.round(totalDebe * 100) / 100 == math.round(totalHaber * 100) / 100
            }
        }
        
        if (tipo == "diario") {
            return {
                tipo: "Libro Diario",
                totalAsientos: database.countAs,
                asientos: database.asientos
            }
        }
        
        if (tipo == "ventas" || tipo == "compras") {
            let total = 0
            let count = 0
            let i = 0
            while (i < database.countTx) {
                let tx = database.transacciones[i]
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
            while (i < database.countTx) {
                let tx = database.transacciones[i]
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
        
        return { error: "Tipo no válido" }
    }
}

// Demo
func demoConsola() {
    console.log("=== DEMO ===")
    
    let tx1 = procesarTx("ventas", "MX", "100000")
    console.log("✓ Venta México: $" + tx1.total + " " + tx1.moneda)
    
    let tx2 = procesarTx("compras", "COL", "50000")
    console.log("✓ Compra Colombia: $" + tx2.total + " " + tx2.moneda)
    
    let tx3 = procesarTx("ventas", "AR", "75000")
    console.log("✓ Venta Argentina: $" + tx3.total + " " + tx3.moneda)
    
    console.log("\nTotal transacciones: " + database.countTx)
    console.log("Total asientos: " + database.countAs)
    
    let engine = ReportesFinancieros
    let balance = engine.use("reporte balance")
    console.log("\nBalance: Debe=$" + balance.totalDebe + " Haber=$" + balance.totalHaber)
}

// Parse params
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
    html = html + ".container{max-width:1200px;margin:0 auto;background:white;padding:30px;border-radius:10px;}\n"
    html = html + "h1{color:#2c3e50;text-align:center;}\n"
    html = html + ".grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(350px,1fr));gap:20px;}\n"
    html = html + ".card{background:#f8f9fa;padding:20px;border-radius:8px;border:1px solid #ddd;}\n"
    html = html + "input,select,textarea{width:100%;padding:8px;margin:5px 0;}\n"
    html = html + "button{background:#667eea;color:white;padding:10px;border:none;border-radius:4px;width:100%;cursor:pointer;}\n"
    html = html + "</style>\n</head>\n<body>\n"
    
    html = html + "<div class=\"container\">\n"
    html = html + "<h1>🌍 Sistema Contable LATAM</h1>\n"
    html = html + "<p style=\"text-align:center\">POC - " + database.countTx + " transacciones - " + database.countAs + " asientos</p>\n"
    
    html = html + "<div class=\"grid\">\n"
    
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
    html = html + "<input type=\"number\" name=\"importe\" value=\"100000\">\n"
    html = html + "<button type=\"submit\">Procesar</button>\n"
    html = html + "</form>\n"
    html = html + "</div>\n"
    
    html = html + "<div class=\"card\">\n"
    html = html + "<h3>DSL Query</h3>\n"
    html = html + "<form action=\"/dsl\" method=\"POST\">\n"
    html = html + "<textarea name=\"query\" rows=\"3\">reporte balance</textarea>\n"
    html = html + "<button type=\"submit\">Ejecutar</button>\n"
    html = html + "</form>\n"
    html = html + "</div>\n"
    
    html = html + "</div>\n"
    
    html = html + "<p style=\"text-align:center;margin-top:20px;\">\n"
    html = html + "<a href=\"/demo\">Demo</a> | "
    html = html + "<a href=\"/libro\">Libro Diario</a> | "
    html = html + "<a href=\"/api/transacciones\">API TX</a> | "
    html = html + "<a href=\"/api/asientos\">API AS</a>\n"
    html = html + "</p>\n"
    
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
    
    let tx = procesarTx(tipo, region, importe)
    
    if (database.countAs == 0) {
        return "<h1>Error</h1><p><a href=\"/\">Volver</a></p>"
    }
    
    let asiento = database.asientos[database.countAs - 1]
    
    let html = "<!DOCTYPE html>\n<html>\n<head>\n"
    html = html + "<title>Comprobante</title>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "<style>\n"
    html = html + "body{font-family:Arial;margin:20px;}\n"
    html = html + ".comp{max-width:800px;margin:0 auto;padding:30px;border:2px solid #667eea;}\n"
    html = html + "table{width:100%;border-collapse:collapse;}\n"
    html = html + "th,td{padding:8px;text-align:left;border-bottom:1px solid #ddd;}\n"
    html = html + "</style>\n</head>\n<body>\n"
    
    html = html + "<div class=\"comp\">\n"
    html = html + "<h1>COMPROBANTE</h1>\n"
    html = html + "<p>ID: " + tx.id + "</p>\n"
    html = html + "<p>País: " + tx.pais + "</p>\n"
    html = html + "<p>Importe: " + tx.moneda + " " + tx.importe + "</p>\n"
    html = html + "<p>IVA: " + tx.moneda + " " + tx.iva + "</p>\n"
    html = html + "<p>TOTAL: " + tx.moneda + " " + tx.total + "</p>\n"
    
    html = html + "<h3>ASIENTO CONTABLE</h3>\n"
    html = html + "<table>\n"
    html = html + "<tr><th>Cuenta</th><th>Descripción</th><th>Debe</th><th>Haber</th></tr>\n"
    
    let i = 0
    while (i < 3) {
        let mov = asiento.movimientos[i]
        html = html + "<tr>\n"
        html = html + "<td>" + mov.cuenta + "</td>\n"
        html = html + "<td>" + mov.descripcion + "</td>\n"
        if (mov.tipo == "DEBE") {
            html = html + "<td>" + mov.monto + "</td><td></td>\n"
        } else {
            html = html + "<td></td><td>" + mov.monto + "</td>\n"
        }
        html = html + "</tr>\n"
        i = i + 1
    }
    
    html = html + "</table>\n"
    html = html + "<p><a href=\"/\">Volver</a></p>\n"
    html = html + "</div>\n</body>\n</html>"
    
    return html
}

// Handler libro
func handleLibro(pathVars, method, body) {
    let html = "<!DOCTYPE html>\n<html>\n<head>\n"
    html = html + "<title>Libro Diario</title>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "</head>\n<body>\n"
    html = html + "<h1>LIBRO DIARIO</h1>\n"
    html = html + "<p>Total: " + database.countAs + " asientos</p>\n"
    
    let i = 0
    while (i < database.countAs) {
        let asiento = database.asientos[i]
        html = html + "<h3>Asiento: " + asiento.id + "</h3>\n"
        html = html + "<p>" + asiento.descripcion + "</p>\n"
        html = html + "<table border=\"1\">\n"
        html = html + "<tr><th>Cuenta</th><th>Descripción</th><th>Debe</th><th>Haber</th></tr>\n"
        
        let j = 0
        while (j < 3) {
            let mov = asiento.movimientos[j]
            html = html + "<tr>\n"
            html = html + "<td>" + mov.cuenta + "</td>\n"
            html = html + "<td>" + mov.descripcion + "</td>\n"
            if (mov.tipo == "DEBE") {
                html = html + "<td>" + mov.monto + "</td><td></td>\n"
            } else {
                html = html + "<td></td><td>" + mov.monto + "</td>\n"
            }
            html = html + "</tr>\n"
            j = j + 1
        }
        
        html = html + "</table>\n"
        i = i + 1
    }
    
    html = html + "<p><a href=\"/\">Volver</a></p>\n"
    html = html + "</body>\n</html>"
    
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
    html = html + "<title>DSL Result</title>\n"
    html = html + "<meta charset=\"UTF-8\">\n"
    html = html + "</head>\n<body>\n"
    html = html + "<h1>Resultado DSL</h1>\n"
    html = html + "<p>Query: " + query + "</p>\n"
    html = html + "<pre>" + JSON.stringify(resultado) + "</pre>\n"
    html = html + "<p><a href=\"/\">Volver</a></p>\n"
    html = html + "</body>\n</html>"
    
    return html
}

// Handler demo
func handleDemo(pathVars, method, body) {
    procesarTx("ventas", "MX", "100000")
    procesarTx("compras", "COL", "50000")
    procesarTx("ventas", "AR", "75000")
    procesarTx("compras", "PE", "60000")
    
    return "<h1>Demo</h1><p>4 transacciones procesadas</p><p>Total: " + database.countTx + "</p><p><a href=\"/\">Volver</a></p>"
}

// APIs
func handleAPITransacciones(pathVars, method, body) {
    return JSON({
        total: database.countTx,
        transacciones: database.transacciones
    })
}

func handleAPIAsientos(pathVars, method, body) {
    return JSON({
        total: database.countAs,
        asientos: database.asientos
    })
}

// Main
func main() {
    demoConsola()
    
    console.log("\n✅ Servidor...")
    
    http.handler("GET", "/", handleHome)
    http.handler("POST", "/procesar", handleProcesar)
    http.handler("GET", "/libro", handleLibro)
    http.handler("POST", "/dsl", handleDSL)
    http.handler("GET", "/demo", handleDemo)
    http.handler("GET", "/api/transacciones", handleAPITransacciones)
    http.handler("GET", "/api/asientos", handleAPIAsientos)
    
    console.log("🌐 http://localhost:8080")
    
    http.serve(":8080")
}

main()