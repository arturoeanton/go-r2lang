// Sistema Contable LATAM - Versión Actualizada con Arrays Anidados
// POC para Siigo - Actualizado 2025-07-22
// Demuestra:
// 1. Manejo correcto de arrays anidados con patrón de reasignación
// 2. Template literals para construcción de HTML
// 3. Sistema multi-región completo
// 4. DSL para contabilidad

console.log("🌍 Sistema Contable LATAM - Versión Actualizada")
console.log("=============================================")

// Configuración multiregión
let configuraciones = {
    MX: {
        nombre: "México",
        iva: 0.16,
        moneda: "MXN",
        formato: "DD/MM/YYYY",
        cuentas: {
            clientes: "1105",
            ventas: "4135",
            ivaDebito: "2408",
            compras: "5195",
            ivaCredito: "2406",
            proveedores: "2205"
        }
    },
    COL: {
        nombre: "Colombia",
        iva: 0.19,
        moneda: "COP",
        formato: "DD/MM/YYYY",
        cuentas: {
            clientes: "1305",
            ventas: "4135",
            ivaDebito: "2408",
            compras: "6205",
            ivaCredito: "2408",
            proveedores: "2205"
        }
    }
}

// Arrays globales
let transacciones = []
let asientosContables = []

// Función principal mejorada con arrays anidados
func procesarTransaccion(tipo, importe, region) {
    let config = configuraciones[region]
    if (!config) {
        console.error("Región no soportada: " + region)
        return null
    }
    
    let importeNum = std.parseFloat(importe)
    let iva = importeNum * config.iva
    let total = importeNum + iva
    
    // Crear transacción
    let tx = {
        id: "TX-" + math.randomInt(99999),
        tipo: tipo,
        fecha: std.now(),
        importe: importeNum,
        iva: iva,
        total: total,
        region: region,
        moneda: config.moneda
    }
    
    // Usar reasignación para push
    transacciones = transacciones.push(tx)
    console.log(`  ✅ Transacción creada: ${tx.id} - Total: ${tx.moneda} ${tx.total}`)
    
    // Crear asiento contable con movimientos vacíos
    let asiento = {
        id: "AS-" + tx.id,
        fecha: tx.fecha,
        region: region,
        descripcion: std.toUpperCase(tipo) + " - " + config.nombre,
        movimientos: []
    }
    
    // IMPORTANTE: Usar patrón de reasignación para push
    if (tipo == "ventas") {
        // Debe: Clientes, Haber: Ventas e IVA
        asiento.movimientos = asiento.movimientos.push({
            cuenta: config.cuentas.clientes,
            descripcion: "Clientes",
            tipo: "DEBE",
            monto: total
        })
        asiento.movimientos = asiento.movimientos.push({
            cuenta: config.cuentas.ventas,
            descripcion: "Ventas",
            tipo: "HABER",
            monto: importeNum
        })
        asiento.movimientos = asiento.movimientos.push({
            cuenta: config.cuentas.ivaDebito,
            descripcion: "IVA Débito",
            tipo: "HABER",
            monto: iva
        })
    } else {
        // Debe: Compras e IVA, Haber: Proveedores
        asiento.movimientos = asiento.movimientos.push({
            cuenta: config.cuentas.compras,
            descripcion: "Compras",
            tipo: "DEBE",
            monto: importeNum
        })
        asiento.movimientos = asiento.movimientos.push({
            cuenta: config.cuentas.ivaCredito,
            descripcion: "IVA Crédito",
            tipo: "DEBE",
            monto: iva
        })
        asiento.movimientos = asiento.movimientos.push({
            cuenta: config.cuentas.proveedores,
            descripcion: "Proveedores",
            tipo: "HABER",
            monto: total
        })
    }
    
    // Usar reasignación para push
    asientosContables = asientosContables.push(asiento)
    console.log(`  ✅ Asiento creado: ${asiento.id} con ${std.len(asiento.movimientos)} movimientos`)
    
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
        
        return "Tipo de reporte no implementado: " + tipo
    }
}

// DSL Procesamiento Transacciones
dsl ProcesadorTransacciones {
    token("TIPO", "venta|compra")
    token("MONTO", "[0-9]+(\.[0-9]+)?")
    token("REGION", "MX|COL|AR|CH|UY|EC|PE")
    
    rule("transaccion", ["TIPO", "MONTO", "REGION"], "procesarDSL")
    
    func procesarDSL(tipo, monto, region) {
        if (tipo == "venta") {
            return procesarTransaccion("ventas", monto, region)
        } else {
            return procesarTransaccion("compras", monto, region)
        }
    }
}

// Procesamiento de ejemplo
console.log("\n📊 Procesando transacciones de ejemplo...")

// México
procesarTransaccion("ventas", "1000", "MX")
procesarTransaccion("compras", "500", "MX")

// Colombia
procesarTransaccion("ventas", "2000", "COL")
procesarTransaccion("compras", "800", "COL")

// Usando DSL
console.log("\n🤖 Procesando con DSL...")
let dslResult1 = ProcesadorTransacciones.use("venta 1500 MX")
let dslResult2 = ProcesadorTransacciones.use("compra 750 COL")

// Reportes
console.log("\n📈 Generando reportes...")
let balanceGeneral = ReportesFinancieros.use("reporte balance")
console.log("Balance General:")
console.log("  Total Debe: " + balanceGeneral.totalDebe)
console.log("  Total Haber: " + balanceGeneral.totalHaber)
console.log("  Cuadrado: " + (balanceGeneral.cuadrado ? "✅ SI" : "❌ NO"))

let reporteIVA = ReportesFinancieros.use("reporte iva")
console.log("\nReporte IVA:")
console.log("  IVA Débito: " + reporteIVA.ivaDebito)
console.log("  IVA Crédito: " + reporteIVA.ivaCredito)
console.log("  Saldo: " + reporteIVA.saldo)

// Generar HTML con template literals
func generarHTML() {
    let html = `<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <title>Sistema Contable LATAM</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .container { max-width: 1200px; margin: 0 auto; }
        h1 { color: #333; }
        .stats { background: #f5f5f5; padding: 10px; border-radius: 5px; }
        .grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; margin: 20px 0; }
        .card { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        table { width: 100%; border-collapse: collapse; }
        th, td { padding: 8px; text-align: left; border-bottom: 1px solid #ddd; }
        th { background: #f0f0f0; }
        .debe { color: #d32f2f; }
        .haber { color: #388e3c; }
        .balance { font-weight: bold; font-size: 1.2em; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🌍 Sistema Contable LATAM</h1>
        <div class="stats">
            <strong>POC para Siigo</strong> | 
            Transacciones: ${std.len(transacciones)} | 
            Asientos: ${std.len(asientosContables)}
        </div>
        
        <div class="grid">
            <div class="card">
                <h2>📊 Balance General</h2>
                <p>Total Debe: <span class="debe">$${balanceGeneral.totalDebe}</span></p>
                <p>Total Haber: <span class="haber">$${balanceGeneral.totalHaber}</span></p>
                <p class="balance">Estado: ${balanceGeneral.cuadrado ? "✅ CUADRADO" : "❌ DESCUADRADO"}</p>
            </div>
            
            <div class="card">
                <h2>💰 Reporte IVA</h2>
                <p>IVA Débito: $${reporteIVA.ivaDebito}</p>
                <p>IVA Crédito: $${reporteIVA.ivaCredito}</p>
                <p class="balance">Saldo: $${reporteIVA.saldo}</p>
            </div>
        </div>
        
        <div class="card">
            <h2>📝 Últimos Asientos</h2>
            <table>
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Fecha</th>
                        <th>Descripción</th>
                        <th>Cuenta</th>
                        <th>Tipo</th>
                        <th>Monto</th>
                    </tr>
                </thead>
                <tbody>`
    
    // Agregar últimos 5 asientos
    let start = math.max(0, std.len(asientosContables) - 5)
    let i = start
    while (i < std.len(asientosContables)) {
        let asiento = asientosContables[i]
        let j = 0
        while (j < std.len(asiento.movimientos)) {
            let mov = asiento.movimientos[j]
            html = html + `
                    <tr>
                        <td>${asiento.id}</td>
                        <td>${asiento.fecha}</td>
                        <td>${asiento.descripcion}</td>
                        <td>${mov.cuenta}</td>
                        <td class="${mov.tipo == "DEBE" ? "debe" : "haber"}">${mov.tipo}</td>
                        <td>$${mov.monto}</td>
                    </tr>`
            j = j + 1
        }
        i = i + 1
    }
    
    html = html + `
                </tbody>
            </table>
        </div>
    </div>
    
    <script>
        console.log('Sistema Contable LATAM cargado');
        console.log('Actualizado con arrays anidados y template literals');
    </script>
</body>
</html>`
    
    return html
}

// Guardar HTML
let htmlContent = generarHTML()
io.writeFile("reporte_contable_actualizado.html", htmlContent)
console.log("\n✅ Reporte HTML generado: reporte_contable_actualizado.html")
console.log("✅ Sistema actualizado con arrays anidados y template literals")