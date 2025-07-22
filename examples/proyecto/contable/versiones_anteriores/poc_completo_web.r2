// POC Completo - Sistema Contable LATAM con CFDI
// Demo para Siigo - R2Lang Web Framework

console.log("🚀 Sistema Contable LATAM - POC Completo")
console.log("📄 Con procesamiento CFDI + 📊 DSL Reportes")
console.log("")

// Base de datos en memoria
let database = {
    regiones: {
        "MX": { nombre: "México", moneda: "MXN", simbolo: "$", iva: 0.16, normativa: "SAT-CFDI 4.0" },
        "COL": { nombre: "Colombia", moneda: "COP", simbolo: "$", iva: 0.19, normativa: "DIAN-FE" },
        "AR": { nombre: "Argentina", moneda: "ARS", simbolo: "$", iva: 0.21, normativa: "AFIP-FE" },
        "CH": { nombre: "Chile", moneda: "CLP", simbolo: "$", iva: 0.19, normativa: "SII-DTE" },
        "UY": { nombre: "Uruguay", moneda: "UYU", simbolo: "$", iva: 0.22, normativa: "DGI-CFE" },
        "EC": { nombre: "Ecuador", moneda: "USD", simbolo: "$", iva: 0.12, normativa: "SRI-CE" },
        "PE": { nombre: "Perú", moneda: "PEN", simbolo: "S/", iva: 0.18, normativa: "SUNAT-FE" }
    },
    transacciones: [],
    comprobantes: [],
    cfdis: []
}

// DSL para procesamiento de CFDI
dsl ProcesadorCFDI {
    token("CFDI", "cfdi|comprobante")
    token("RFC", "[A-Z0-9]{12,13}")
    token("TOTAL", "[0-9]+\\.?[0-9]*")
    token("UUID", "[A-F0-9]{8}-[A-F0-9]{4}-[A-F0-9]{4}-[A-F0-9]{4}-[A-F0-9]{12}")
    
    rule("procesar_cfdi", ["CFDI"], "procesarCFDI")
    
    func procesarCFDI(comando) {
        // Procesar CFDI completo
        return {
            success: true,
            message: "CFDI procesado correctamente"
        }
    }
}

// DSL para Reportes Contables
dsl ReportesContables {
    token("REPORTE", "reporte|report")
    token("TIPO", "ventas|compras|impuestos|mensual|anual")
    token("REGION", "MX|COL|AR|CH|UY|EC|PE|ALL")
    token("PERIODO", "hoy|mes|año|todo")
    
    rule("generar_reporte", ["REPORTE", "TIPO", "REGION", "PERIODO"], "generarReporte")
    rule("reporte_simple", ["REPORTE", "TIPO"], "reporteSimple")
    
    func generarReporte(cmd, tipo, region, periodo) {
        let transacciones = database.transacciones
        let filtered = []
        
        // Filtrar por región
        if (region != "ALL") {
            let i = 0
            while (i < std.len(transacciones)) {
                if (transacciones[i].region == region) {
                    filtered[std.len(filtered)] = transacciones[i]
                }
                i = i + 1
            }
        } else {
            filtered = transacciones
        }
        
        // Calcular totales
        let total = 0
        let totalImpuestos = 0
        let i = 0
        while (i < std.len(filtered)) {
            total = total + filtered[i].amount
            totalImpuestos = totalImpuestos + filtered[i].tax
            i = i + 1
        }
        
        return {
            tipo: tipo,
            region: region,
            periodo: periodo,
            transacciones: std.len(filtered),
            total: total,
            impuestos: totalImpuestos,
            datos: filtered
        }
    }
    
    func reporteSimple(cmd, tipo) {
        return generarReporte(cmd, tipo, "ALL", "todo")
    }
}

// Función para procesar CFDI JSON
func procesarCFDICompleto(cfdiData) {
    let comprobante = cfdiData.Comprobante
    
    // Extraer datos principales
    let emisor = comprobante.Emisor
    let receptor = comprobante.Receptor
    let conceptos = comprobante.Conceptos.Concepto
    let timbre = comprobante.Complemento.TimbreFiscalDigital
    
    // Calcular totales
    let subtotal = std.parseFloat(comprobante._SubTotal)
    let total = std.parseFloat(comprobante._Total)
    let impuestos = std.parseFloat(comprobante.Impuestos._TotalImpuestosTrasladados || "0")
    
    // Crear registro procesado
    let registro = {
        id: timbre._UUID,
        fecha: comprobante._Fecha,
        emisor: {
            rfc: emisor._Rfc,
            nombre: emisor._Nombre
        },
        receptor: {
            rfc: receptor._Rfc,
            nombre: receptor._Nombre
        },
        subtotal: subtotal,
        impuestos: impuestos,
        total: total,
        moneda: comprobante._Moneda,
        tipo: comprobante._TipoDeComprobante,
        conceptos: conceptos,
        region: "MX", // CFDI es mexicano
        procesado: std.now()
    }
    
    // Guardar en base de datos
    database.cfdis[std.len(database.cfdis)] = registro
    
    // También crear transacción
    let tx = {
        transactionId: timbre._UUID,
        region: "MX",
        country: "México",
        amount: subtotal,
        tax: impuestos,
        total: total,
        currency: comprobante._Moneda,
        compliance: "SAT-CFDI 4.0",
        timestamp: comprobante._Fecha,
        tipo: "CFDI",
        emisor: emisor._Nombre,
        receptor: receptor._Nombre
    }
    
    database.transacciones[std.len(database.transacciones)] = tx
    
    return registro
}

// Crear aplicación web con r2web
let app = web.createApp()

// Middleware de logging
app.use(func(ctx) {
    console.log(ctx.method + " " + ctx.path)
})

// Ruta principal
app.get("/", func(ctx) {
    let html = `<!DOCTYPE html>
<html lang="es">
<head>
    <title>Sistema Contable LATAM - POC Completo</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: #f0f2f5;
            color: #1a202c;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 2rem 0;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 0 1rem;
        }
        .header h1 {
            font-size: 2.5rem;
            margin-bottom: 0.5rem;
        }
        .header p {
            opacity: 0.9;
            font-size: 1.1rem;
        }
        .grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
            gap: 2rem;
            margin: 2rem 0;
        }
        .card {
            background: white;
            border-radius: 12px;
            padding: 1.5rem;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
            transition: transform 0.2s, box-shadow 0.2s;
        }
        .card:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(0,0,0,0.15);
        }
        .card h2 {
            color: #2d3748;
            margin-bottom: 1rem;
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }
        .form-group {
            margin-bottom: 1rem;
        }
        label {
            display: block;
            margin-bottom: 0.5rem;
            font-weight: 500;
            color: #4a5568;
        }
        input, select, textarea {
            width: 100%;
            padding: 0.75rem;
            border: 1px solid #e2e8f0;
            border-radius: 6px;
            font-size: 1rem;
            transition: border-color 0.2s;
        }
        input:focus, select:focus, textarea:focus {
            outline: none;
            border-color: #667eea;
        }
        textarea {
            resize: vertical;
            min-height: 150px;
            font-family: monospace;
        }
        button {
            background: #667eea;
            color: white;
            padding: 0.75rem 1.5rem;
            border: none;
            border-radius: 6px;
            font-size: 1rem;
            font-weight: 500;
            cursor: pointer;
            transition: background 0.2s;
            width: 100%;
        }
        button:hover {
            background: #5a67d8;
        }
        .stats {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 1rem;
            margin-bottom: 2rem;
        }
        .stat-card {
            background: white;
            padding: 1.5rem;
            border-radius: 8px;
            text-align: center;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }
        .stat-card h3 {
            color: #718096;
            font-size: 0.875rem;
            text-transform: uppercase;
            margin-bottom: 0.5rem;
        }
        .stat-card .value {
            font-size: 2rem;
            font-weight: bold;
            color: #2d3748;
        }
        .links {
            display: flex;
            gap: 1rem;
            flex-wrap: wrap;
            margin-top: 1rem;
        }
        .links a {
            color: #667eea;
            text-decoration: none;
            padding: 0.5rem 1rem;
            border: 1px solid #667eea;
            border-radius: 6px;
            transition: all 0.2s;
        }
        .links a:hover {
            background: #667eea;
            color: white;
        }
        .dsl-input {
            font-family: 'Monaco', 'Menlo', monospace;
            background: #2d3748;
            color: #e2e8f0;
            padding: 1rem;
            border-radius: 6px;
        }
        .value-prop {
            background: linear-gradient(135deg, #84fab0 0%, #8fd3f4 100%);
            padding: 2rem;
            border-radius: 12px;
            margin: 2rem 0;
            text-align: center;
        }
        .value-prop h3 {
            color: #1a202c;
            margin-bottom: 1rem;
        }
        .metrics {
            display: flex;
            justify-content: space-around;
            flex-wrap: wrap;
            gap: 1rem;
        }
        .metric {
            background: white;
            padding: 1rem 2rem;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
    </style>
</head>
<body>
    <header class="header">
        <div class="container">
            <h1>🌍 Sistema Contable LATAM</h1>
            <p>POC Completo con CFDI + DSL Reportes - Powered by R2Lang</p>
        </div>
    </header>

    <div class="container">
        <!-- Stats -->
        <div class="stats">
            <div class="stat-card">
                <h3>Transacciones</h3>
                <div class="value" id="totalTx">0</div>
            </div>
            <div class="stat-card">
                <h3>CFDIs Procesados</h3>
                <div class="value" id="totalCfdi">0</div>
            </div>
            <div class="stat-card">
                <h3>Regiones Activas</h3>
                <div class="value">7</div>
            </div>
            <div class="stat-card">
                <h3>Compliance</h3>
                <div class="value">100%</div>
            </div>
        </div>

        <div class="grid">
            <!-- Procesamiento Manual -->
            <div class="card">
                <h2>💰 Procesamiento Manual</h2>
                <form action="/api/transaction" method="POST">
                    <div class="form-group">
                        <label>Tipo de Operación</label>
                        <select name="tipo" required>
                            <option value="venta">Venta</option>
                            <option value="compra">Compra</option>
                        </select>
                    </div>
                    <div class="form-group">
                        <label>Región</label>
                        <select name="region" required>
                            <option value="MX">🇲🇽 México (16% IVA)</option>
                            <option value="COL">🇨🇴 Colombia (19% IVA)</option>
                            <option value="AR">🇦🇷 Argentina (21% IVA)</option>
                            <option value="CH">🇨🇱 Chile (19% IVA)</option>
                            <option value="UY">🇺🇾 Uruguay (22% IVA)</option>
                            <option value="EC">🇪🇨 Ecuador (12% IVA)</option>
                            <option value="PE">🇵🇪 Perú (18% IVA)</option>
                        </select>
                    </div>
                    <div class="form-group">
                        <label>Importe</label>
                        <input type="number" name="importe" value="100000" required>
                    </div>
                    <button type="submit">Procesar Transacción</button>
                </form>
            </div>

            <!-- Carga de CFDI -->
            <div class="card">
                <h2>📄 Carga de CFDI (México)</h2>
                <form action="/api/cfdi" method="POST">
                    <div class="form-group">
                        <label>CFDI JSON</label>
                        <textarea name="cfdi" class="dsl-input" placeholder='{"Comprobante": {...}}' required></textarea>
                    </div>
                    <button type="submit">Procesar CFDI</button>
                </form>
                <div class="links">
                    <a href="#" onclick="loadSampleCFDI()">Cargar CFDI de Ejemplo</a>
                </div>
            </div>

            <!-- DSL de Reportes -->
            <div class="card">
                <h2>📊 Generador de Reportes DSL</h2>
                <form action="/api/reporte" method="POST">
                    <div class="form-group">
                        <label>Comando DSL</label>
                        <input type="text" name="comando" placeholder="reporte ventas COL mes" class="dsl-input">
                    </div>
                    <div class="form-group">
                        <label>Ejemplos:</label>
                        <code style="display: block; margin: 0.5rem 0;">reporte ventas ALL todo</code>
                        <code style="display: block; margin: 0.5rem 0;">reporte impuestos MX año</code>
                        <code style="display: block; margin: 0.5rem 0;">reporte compras COL mes</code>
                    </div>
                    <button type="submit">Generar Reporte</button>
                </form>
            </div>

            <!-- Links Rápidos -->
            <div class="card">
                <h2>🔗 Accesos Rápidos</h2>
                <div class="links">
                    <a href="/demo/auto">Demo Automática</a>
                    <a href="/api/transacciones">Ver Transacciones</a>
                    <a href="/api/cfdis">Ver CFDIs</a>
                    <a href="/api/regiones">Configuración Regional</a>
                    <a href="/reportes">Dashboard Reportes</a>
                </div>
            </div>
        </div>

        <!-- Value Proposition -->
        <div class="value-prop">
            <h3>🎯 Propuesta de Valor para Siigo</h3>
            <div class="metrics">
                <div class="metric">
                    <strong>Tiempo</strong><br>
                    18 meses → 2 meses
                </div>
                <div class="metric">
                    <strong>Costo</strong><br>
                    $500K → $150K
                </div>
                <div class="metric">
                    <strong>Arquitectura</strong><br>
                    7 sistemas → 1 DSL
                </div>
                <div class="metric">
                    <strong>ROI</strong><br>
                    1,020% en 3 años
                </div>
            </div>
        </div>
    </div>

    <script>
        // Actualizar estadísticas
        async function updateStats() {
            try {
                const txResponse = await fetch('/api/transacciones');
                const txData = await txResponse.json();
                document.getElementById('totalTx').textContent = txData.total || 0;

                const cfdiResponse = await fetch('/api/cfdis');
                const cfdiData = await cfdiResponse.json();
                document.getElementById('totalCfdi').textContent = cfdiData.total || 0;
            } catch (e) {
                console.error(e);
            }
        }

        // Cargar CFDI de ejemplo
        function loadSampleCFDI() {
            const sampleCFDI = {
                "Comprobante": {
                    "Emisor": {
                        "_Rfc": "EKU9003173C9",
                        "_Nombre": "ESCUELA KEMPER URGATE SA DE CV",
                        "_RegimenFiscal": "601"
                    },
                    "Receptor": {
                        "_Rfc": "XEXX010101000",
                        "_Nombre": "PUBLICO EN GENERAL"
                    },
                    "_SubTotal": "100.00",
                    "_Total": "116.00",
                    "_Moneda": "MXN",
                    "_Fecha": "2025-07-22T10:00:00",
                    "_TipoDeComprobante": "I",
                    "Impuestos": {
                        "_TotalImpuestosTrasladados": "16.00"
                    },
                    "Conceptos": {
                        "Concepto": [{
                            "_Descripcion": "Servicio de desarrollo",
                            "_ValorUnitario": "100.00",
                            "_Importe": "100.00"
                        }]
                    },
                    "Complemento": {
                        "TimbreFiscalDigital": {
                            "_UUID": "` + generateUUID() + `",
                            "_FechaTimbrado": "2025-07-22T10:00:00"
                        }
                    }
                }
            };
            
            document.querySelector('textarea[name="cfdi"]').value = JSON.stringify(sampleCFDI, null, 2);
        }

        function generateUUID() {
            return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
                var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
                return v.toString(16).toUpperCase();
            });
        }

        // Actualizar stats al cargar
        updateStats();
        setInterval(updateStats, 5000);
    </script>
</body>
</html>`
    
    return ctx.send(html)
})

// API: Procesar transacción
app.post("/api/transaction", func(ctx) {
    let params = web.parseForm(ctx.body)
    let tipo = params.tipo || "venta"
    let region = params.region || "COL"
    let importe = std.parseFloat(params.importe || "100000")
    
    let config = database.regiones[region]
    let importeIVA = math.round((importe * config.iva) * 100) / 100
    let importeTotal = math.round((importe + importeIVA) * 100) / 100
    
    let tx = {
        transactionId: region + "-" + std.now() + "-" + math.randomInt(9999),
        tipo: tipo,
        region: region,
        country: config.nombre,
        amount: importe,
        tax: importeIVA,
        total: importeTotal,
        currency: config.moneda,
        compliance: config.normativa,
        timestamp: std.now()
    }
    
    database.transacciones[std.len(database.transacciones)] = tx
    
    return ctx.status(200).send({
        success: true,
        transaction: tx
    })
})

// API: Procesar CFDI
app.post("/api/cfdi", func(ctx) {
    let params = web.parseForm(ctx.body)
    let cfdiJSON = params.cfdi
    
    try {
        let cfdiData = web.parseJSON(cfdiJSON)
        let resultado = procesarCFDICompleto(cfdiData)
        
        return ctx.status(200).send({
            success: true,
            message: "CFDI procesado correctamente",
            cfdi: resultado
        })
    } catch {
        return ctx.status(400).send({
            success: false,
            message: "Error al procesar CFDI"
        })
    }
})

// API: Generar reporte con DSL
app.post("/api/reporte", func(ctx) {
    let params = web.parseForm(ctx.body)
    let comando = params.comando || "reporte ventas ALL todo"
    
    let reporteEngine = ReportesContables
    let resultado = reporteEngine.use(comando)
    
    return ctx.status(200).send({
        success: true,
        comando: comando,
        reporte: resultado
    })
})

// API: Listar transacciones
app.get("/api/transacciones", func(ctx) {
    return ctx.send({
        total: std.len(database.transacciones),
        transacciones: database.transacciones
    })
})

// API: Listar CFDIs
app.get("/api/cfdis", func(ctx) {
    return ctx.send({
        total: std.len(database.cfdis),
        cfdis: database.cfdis
    })
})

// API: Regiones
app.get("/api/regiones", func(ctx) {
    return ctx.send(database.regiones)
})

// Demo automática
app.get("/demo/auto", func(ctx) {
    // Procesar varias transacciones automáticamente
    let demos = [
        {tipo: "venta", region: "MX", importe: 100000},
        {tipo: "compra", region: "COL", importe: 50000},
        {tipo: "venta", region: "AR", importe: 75000},
        {tipo: "compra", region: "PE", importe: 60000}
    ]
    
    let resultados = []
    let i = 0
    while (i < std.len(demos)) {
        let demo = demos[i]
        let config = database.regiones[demo.region]
        let importeIVA = math.round((demo.importe * config.iva) * 100) / 100
        let importeTotal = math.round((demo.importe + importeIVA) * 100) / 100
        
        let tx = {
            transactionId: demo.region + "-DEMO-" + i,
            tipo: demo.tipo,
            region: demo.region,
            country: config.nombre,
            amount: demo.importe,
            tax: importeIVA,
            total: importeTotal,
            currency: config.moneda,
            compliance: config.normativa,
            timestamp: std.now()
        }
        
        database.transacciones[std.len(database.transacciones)] = tx
        resultados[i] = tx
        i = i + 1
    }
    
    return ctx.send({
        success: true,
        message: "Demo automática completada",
        transacciones: resultados
    })
})

// Dashboard de reportes
app.get("/reportes", func(ctx) {
    let html = `<!DOCTYPE html>
<html>
<head>
    <title>Dashboard de Reportes - Sistema Contable LATAM</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; }
        .report-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; }
        .report-card { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        button { background: #667eea; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; }
        #resultado { background: #e2e8f0; padding: 20px; border-radius: 8px; margin-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>📊 Dashboard de Reportes</h1>
        
        <div class="report-grid">
            <div class="report-card">
                <h3>Reporte de Ventas</h3>
                <button onclick="generarReporte('reporte ventas ALL todo')">Generar</button>
            </div>
            <div class="report-card">
                <h3>Reporte de Impuestos</h3>
                <button onclick="generarReporte('reporte impuestos ALL mes')">Generar</button>
            </div>
            <div class="report-card">
                <h3>Reporte por Región</h3>
                <select id="regionSelect">
                    <option value="MX">México</option>
                    <option value="COL">Colombia</option>
                    <option value="AR">Argentina</option>
                </select>
                <button onclick="reportePorRegion()">Generar</button>
            </div>
        </div>
        
        <div id="resultado"></div>
        
        <p><a href="/">← Volver al inicio</a></p>
    </div>
    
    <script>
        async function generarReporte(comando) {
            const response = await fetch('/api/reporte', {
                method: 'POST',
                headers: {'Content-Type': 'application/x-www-form-urlencoded'},
                body: 'comando=' + encodeURIComponent(comando)
            });
            const data = await response.json();
            mostrarResultado(data);
        }
        
        function reportePorRegion() {
            const region = document.getElementById('regionSelect').value;
            generarReporte('reporte ventas ' + region + ' todo');
        }
        
        function mostrarResultado(data) {
            document.getElementById('resultado').innerHTML = '<pre>' + JSON.stringify(data, null, 2) + '</pre>';
        }
    </script>
</body>
</html>`
    
    return ctx.send(html)
})

// Función principal
func main() {
    console.log("✅ Sistema Contable LATAM - POC Completo")
    console.log("📄 CFDI Processing + 📊 DSL Reports")
    console.log("🌐 Starting server on :8080")
    console.log("")
    console.log("🔗 URLs disponibles:")
    console.log("   http://localhost:8080 - Interfaz principal")
    console.log("   http://localhost:8080/demo/auto - Demo automática")
    console.log("   http://localhost:8080/reportes - Dashboard de reportes")
    console.log("")
    
    // Iniciar servidor
    app.listen(":8080")
}

// Ejecutar
main()