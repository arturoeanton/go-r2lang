// POC Final Funcional - Sistema Contable LATAM con CFDI
// Demo completa para Siigo usando R2Lang existente

console.log("🚀 Sistema Contable LATAM - POC Final")
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
    cfdis: []
}

// DSL para procesamiento de CFDI
dsl ProcesadorCFDI {
    token("PROCESAR", "procesar|process")
    token("CFDI", "cfdi")
    token("UUID", "[A-F0-9-]+")
    
    rule("procesar_cfdi", ["PROCESAR", "CFDI"], "procesarCFDI")
    
    func procesarCFDI(cmd1, cmd2) {
        return {
            success: true,
            message: "CFDI procesado por DSL"
        }
    }
}

// DSL para Reportes Contables
dsl ReportesContables {
    token("REPORTE", "reporte|report")
    token("TIPO", "ventas|compras|impuestos|todo")
    token("REGION", "MX|COL|AR|CH|UY|EC|PE|ALL")
    
    rule("generar_reporte", ["REPORTE", "TIPO", "REGION"], "generarReporte")
    rule("reporte_simple", ["REPORTE", "TIPO"], "reporteSimple")
    
    func generarReporte(cmd, tipo, region) {
        let transacciones = database.transacciones
        let filtered = []
        
        // Filtrar por región si no es ALL
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
        
        // Filtrar por tipo si no es todo
        if (tipo != "todo") {
            let typeFiltered = []
            let i = 0
            while (i < std.len(filtered)) {
                if (filtered[i].tipo == tipo) {
                    typeFiltered[std.len(typeFiltered)] = filtered[i]
                }
                i = i + 1
            }
            filtered = typeFiltered
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
            transacciones: std.len(filtered),
            total: math.round(total * 100) / 100,
            impuestos: math.round(totalImpuestos * 100) / 100,
            montoTotal: math.round((total + totalImpuestos) * 100) / 100
        }
    }
    
    func reporteSimple(cmd, tipo) {
        return generarReporte(cmd, tipo, "ALL")
    }
}

// Función para procesar CFDI completo
func procesarCFDICompleto(cfdiJSON) {
    let cfdiData = JSON.parse(cfdiJSON)
    let comprobante = cfdiData.Comprobante
    
    // Extraer datos principales
    let emisor = comprobante.Emisor
    let receptor = comprobante.Receptor
    let timbre = comprobante.Complemento.TimbreFiscalDigital
    
    // Calcular totales
    let subtotal = std.parseFloat(comprobante._SubTotal)
    let total = std.parseFloat(comprobante._Total)
    let impuestos = 0
    
    if (comprobante.Impuestos && comprobante.Impuestos._TotalImpuestosTrasladados) {
        impuestos = std.parseFloat(comprobante.Impuestos._TotalImpuestosTrasladados)
    }
    
    // Crear registro
    let registro = {
        id: timbre._UUID,
        fecha: comprobante._Fecha,
        emisor: emisor._Nombre,
        emisorRFC: emisor._Rfc,
        receptor: receptor._Nombre,
        receptorRFC: receptor._Rfc,
        subtotal: subtotal,
        impuestos: impuestos,
        total: total,
        moneda: comprobante._Moneda,
        tipo: "CFDI-" + comprobante._TipoDeComprobante,
        region: "MX",
        procesado: std.now()
    }
    
    // Guardar CFDI
    database.cfdis[std.len(database.cfdis)] = registro
    
    // Crear transacción
    let tx = {
        transactionId: timbre._UUID,
        tipo: "cfdi",
        region: "MX",
        country: "México",
        amount: subtotal,
        tax: impuestos,
        total: total,
        currency: comprobante._Moneda,
        compliance: "SAT-CFDI 4.0",
        timestamp: comprobante._Fecha,
        emisor: emisor._Nombre,
        receptor: receptor._Nombre
    }
    
    database.transacciones[std.len(database.transacciones)] = tx
    
    return registro
}

// Función para parsear form data
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
                let key = parts[0]
                let value = parts[1]
                // Decodificar URL encoding básico
                value = std.replace(value, "+", " ")
                value = std.replace(value, "%20", " ")
                value = std.replace(value, "%2C", ",")
                params[key] = value
            }
        }
        i = i + 1
    }
    
    return params
}

// Handler página principal
func handleHome(pathVars, method, body) {
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
        .stat-value {
            font-size: 2rem;
            font-weight: bold;
            color: #2d3748;
        }
        .value-prop {
            background: linear-gradient(135deg, #84fab0 0%, #8fd3f4 100%);
            padding: 2rem;
            border-radius: 12px;
            margin: 2rem 0;
        }
        .links a {
            color: #667eea;
            text-decoration: none;
            margin-right: 1rem;
        }
        pre { background: #f4f4f4; padding: 1rem; border-radius: 4px; overflow-x: auto; }
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
        <div class="stats">
            <div class="stat-card">
                <h3>Transacciones</h3>
                <div class="stat-value">` + std.len(database.transacciones) + `</div>
            </div>
            <div class="stat-card">
                <h3>CFDIs</h3>
                <div class="stat-value">` + std.len(database.cfdis) + `</div>
            </div>
            <div class="stat-card">
                <h3>Regiones</h3>
                <div class="stat-value">7</div>
            </div>
        </div>

        <div class="grid">
            <!-- Procesamiento Manual -->
            <div class="card">
                <h2>💰 Procesamiento Manual</h2>
                <form action="/procesar-transaccion" method="POST">
                    <div class="form-group">
                        <label>Tipo</label>
                        <select name="tipo">
                            <option value="ventas">Venta</option>
                            <option value="compras">Compra</option>
                        </select>
                    </div>
                    <div class="form-group">
                        <label>Región</label>
                        <select name="region">
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
                    <button type="submit">Procesar</button>
                </form>
            </div>

            <!-- Carga de CFDI -->
            <div class="card">
                <h2>📄 Carga de CFDI (México)</h2>
                <form action="/procesar-cfdi" method="POST">
                    <div class="form-group">
                        <label>CFDI JSON</label>
                        <textarea name="cfdi" placeholder='{"Comprobante": {...}}'></textarea>
                    </div>
                    <button type="submit">Procesar CFDI</button>
                    <p style="margin-top: 1rem;">
                        <a href="#" onclick="cargarEjemploCFDI(); return false;">Cargar ejemplo</a>
                    </p>
                </form>
            </div>

            <!-- DSL de Reportes -->
            <div class="card">
                <h2>📊 Generador de Reportes DSL</h2>
                <form action="/generar-reporte" method="POST">
                    <div class="form-group">
                        <label>Comando DSL</label>
                        <input type="text" name="comando" value="reporte todo ALL" style="font-family: monospace;">
                    </div>
                    <div class="form-group">
                        <label>Ejemplos:</label>
                        <code>reporte ventas MX</code><br>
                        <code>reporte compras COL</code><br>
                        <code>reporte impuestos ALL</code><br>
                        <code>reporte todo PE</code>
                    </div>
                    <button type="submit">Generar Reporte</button>
                </form>
            </div>

            <!-- Links -->
            <div class="card">
                <h2>🔗 Accesos Rápidos</h2>
                <div class="links">
                    <a href="/demo-auto">Demo Automática</a><br><br>
                    <a href="/api/transacciones">Ver Transacciones</a><br><br>
                    <a href="/api/cfdis">Ver CFDIs</a><br><br>
                    <a href="/api/regiones">Configuración Regional</a><br><br>
                    <a href="/dashboard">Dashboard Reportes</a>
                </div>
            </div>
        </div>

        <div class="value-prop">
            <h3>🎯 Propuesta de Valor para Siigo</h3>
            <p><strong>Tiempo:</strong> 18 meses → 2 meses | 
               <strong>Costo:</strong> $500K → $150K | 
               <strong>Arquitectura:</strong> 7 sistemas → 1 DSL | 
               <strong>ROI:</strong> 1,020% en 3 años</p>
        </div>
    </div>

    <script>
        function cargarEjemploCFDI() {
            const ejemplo = {
                "Comprobante": {
                    "Emisor": {
                        "_Rfc": "EKU9003173C9",
                        "_Nombre": "EMPRESA DEMO SA DE CV",
                        "_RegimenFiscal": "601"
                    },
                    "Receptor": {
                        "_Rfc": "XAXX010101000",
                        "_Nombre": "PUBLICO EN GENERAL"
                    },
                    "_SubTotal": "1000.00",
                    "_Total": "1160.00",
                    "_Moneda": "MXN",
                    "_Fecha": "2025-07-22T12:00:00",
                    "_TipoDeComprobante": "I",
                    "Impuestos": {
                        "_TotalImpuestosTrasladados": "160.00"
                    },
                    "Complemento": {
                        "TimbreFiscalDigital": {
                            "_UUID": "` + generarUUID() + `",
                            "_FechaTimbrado": "2025-07-22T12:00:00"
                        }
                    }
                }
            };
            document.querySelector('textarea[name="cfdi"]').value = JSON.stringify(ejemplo, null, 2);
        }
        
        function generarUUID() {
            return 'XXXXXXXX-XXXX-4XXX-YXXX-XXXXXXXXXXXX'.replace(/[XY]/g, function(c) {
                var r = Math.random() * 16 | 0, v = c == 'X' ? r : (r & 0x3 | 0x8);
                return v.toString(16).toUpperCase();
            });
        }
    </script>
</body>
</html>`
    
    return html
}

// Handler procesar transacción
func handleProcesarTransaccion(pathVars, method, body) {
    let params = parseFormData(body)
    let tipo = params.tipo || "ventas"
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
    
    // Generar comprobante HTML
    let html = `<!DOCTYPE html>
<html>
<head>
    <title>Comprobante - Sistema Contable LATAM</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background: #f5f5f5; }
        .comprobante { background: white; padding: 30px; border-radius: 10px; max-width: 600px; margin: 0 auto; }
        .header { text-align: center; border-bottom: 2px solid #667eea; padding-bottom: 20px; margin-bottom: 20px; }
        .row { display: flex; justify-content: space-between; margin: 10px 0; }
        .total { background: #f0f2f5; padding: 15px; border-radius: 8px; margin: 20px 0; }
        .back { text-align: center; margin-top: 20px; }
        a { color: #667eea; text-decoration: none; }
    </style>
</head>
<body>
    <div class="comprobante">
        <div class="header">
            <h1>COMPROBANTE DE ` + std.toUpperCase(tipo) + `</h1>
            <h2>` + config.nombre + ` (` + region + `)</h2>
        </div>
        
        <div class="row">
            <strong>ID Transacción:</strong>
            <span>` + tx.transactionId + `</span>
        </div>
        
        <div class="row">
            <strong>Fecha:</strong>
            <span>` + tx.timestamp + `</span>
        </div>
        
        <div class="row">
            <strong>Normativa:</strong>
            <span>` + tx.compliance + `</span>
        </div>
        
        <div class="total">
            <div class="row">
                <strong>Importe Base:</strong>
                <span>` + config.simbolo + ` ` + importe + ` ` + config.moneda + `</span>
            </div>
            <div class="row">
                <strong>IVA (` + (config.iva * 100) + `%):</strong>
                <span>` + config.simbolo + ` ` + importeIVA + ` ` + config.moneda + `</span>
            </div>
            <div class="row" style="border-top: 2px solid #333; padding-top: 10px;">
                <strong>TOTAL:</strong>
                <strong>` + config.simbolo + ` ` + importeTotal + ` ` + config.moneda + `</strong>
            </div>
        </div>
        
        <p style="text-align: center; color: #27ae60;"><strong>✅ TRANSACCIÓN PROCESADA</strong></p>
        
        <div class="back">
            <a href="/">← Volver al inicio</a>
        </div>
    </div>
</body>
</html>`
    
    return html
}

// Handler procesar CFDI
func handleProcesarCFDI(pathVars, method, body) {
    let params = parseFormData(body)
    let cfdiJSON = params.cfdi || "{}"
    
    try {
        let resultado = procesarCFDICompleto(cfdiJSON)
        
        let html = `<!DOCTYPE html>
<html>
<head>
    <title>CFDI Procesado - Sistema Contable LATAM</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background: #f5f5f5; }
        .container { background: white; padding: 30px; border-radius: 10px; max-width: 800px; margin: 0 auto; }
        h1 { color: #2c3e50; }
        .info { background: #e8f5e9; padding: 15px; border-radius: 8px; margin: 20px 0; }
        .back { text-align: center; margin-top: 20px; }
        a { color: #667eea; text-decoration: none; }
        pre { background: #f4f4f4; padding: 15px; border-radius: 4px; overflow-x: auto; }
    </style>
</head>
<body>
    <div class="container">
        <h1>✅ CFDI Procesado Exitosamente</h1>
        
        <div class="info">
            <p><strong>UUID:</strong> ` + resultado.id + `</p>
            <p><strong>Fecha:</strong> ` + resultado.fecha + `</p>
            <p><strong>Emisor:</strong> ` + resultado.emisor + ` (` + resultado.emisorRFC + `)</p>
            <p><strong>Receptor:</strong> ` + resultado.receptor + ` (` + resultado.receptorRFC + `)</p>
            <p><strong>Total:</strong> $ ` + resultado.total + ` ` + resultado.moneda + `</p>
        </div>
        
        <h3>Datos Completos:</h3>
        <pre>` + JSON(resultado) + `</pre>
        
        <div class="back">
            <a href="/">← Volver al inicio</a>
        </div>
    </div>
</body>
</html>`
        
        return html
    } catch {
        return `<h1>Error al procesar CFDI</h1><p>Por favor verifica el formato JSON</p><a href="/">Volver</a>`
    }
}

// Handler generar reporte
func handleGenerarReporte(pathVars, method, body) {
    let params = parseFormData(body)
    let comando = params.comando || "reporte todo ALL"
    
    let reporteEngine = ReportesContables
    let resultado = reporteEngine.use(comando)
    
    let html = `<!DOCTYPE html>
<html>
<head>
    <title>Reporte - Sistema Contable LATAM</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background: #f5f5f5; }
        .container { background: white; padding: 30px; border-radius: 10px; max-width: 800px; margin: 0 auto; }
        h1 { color: #2c3e50; }
        .metric { background: #f0f2f5; padding: 15px; border-radius: 8px; margin: 10px 0; }
        .back { text-align: center; margin-top: 20px; }
        a { color: #667eea; text-decoration: none; }
        pre { background: #f4f4f4; padding: 15px; border-radius: 4px; overflow-x: auto; }
    </style>
</head>
<body>
    <div class="container">
        <h1>📊 Reporte Generado</h1>
        <p><strong>Comando DSL:</strong> <code>` + comando + `</code></p>
        
        <div class="metric">
            <h3>Resumen</h3>
            <p><strong>Tipo:</strong> ` + resultado.tipo + `</p>
            <p><strong>Región:</strong> ` + resultado.region + `</p>
            <p><strong>Transacciones:</strong> ` + resultado.transacciones + `</p>
            <p><strong>Total Base:</strong> $ ` + resultado.total + `</p>
            <p><strong>Total Impuestos:</strong> $ ` + resultado.impuestos + `</p>
            <p><strong>Monto Total:</strong> $ ` + resultado.montoTotal + `</p>
        </div>
        
        <h3>Datos Completos:</h3>
        <pre>` + JSON(resultado) + `</pre>
        
        <div class="back">
            <a href="/">← Volver al inicio</a> | 
            <a href="/dashboard">Ver Dashboard</a>
        </div>
    </div>
</body>
</html>`
    
    return html
}

// Handler demo automática
func handleDemoAuto(pathVars, method, body) {
    // Procesar varias transacciones automáticamente
    let demos = [
        {tipo: "ventas", region: "MX", importe: 100000},
        {tipo: "compras", region: "COL", importe: 50000},
        {tipo: "ventas", region: "AR", importe: 75000},
        {tipo: "compras", region: "PE", importe: 60000},
        {tipo: "ventas", region: "CH", importe: 80000},
        {tipo: "compras", region: "UY", importe: 45000},
        {tipo: "ventas", region: "EC", importe: 90000}
    ]
    
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
        i = i + 1
    }
    
    // Procesar algunos CFDIs de ejemplo
    let cfdiEjemplo = `{
        "Comprobante": {
            "Emisor": {"_Rfc": "DEMO010101ABC", "_Nombre": "DEMO EMPRESA SA"},
            "Receptor": {"_Rfc": "XAXX010101000", "_Nombre": "PUBLICO GENERAL"},
            "_SubTotal": "5000.00",
            "_Total": "5800.00",
            "_Moneda": "MXN",
            "_Fecha": "2025-07-22T15:00:00",
            "_TipoDeComprobante": "I",
            "Impuestos": {"_TotalImpuestosTrasladados": "800.00"},
            "Complemento": {
                "TimbreFiscalDigital": {
                    "_UUID": "DEMO-` + std.now() + `",
                    "_FechaTimbrado": "2025-07-22T15:00:00"
                }
            }
        }
    }`
    
    procesarCFDICompleto(cfdiEjemplo)
    
    return `<!DOCTYPE html>
<html>
<head>
    <title>Demo Automática - Sistema Contable LATAM</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background: #f5f5f5; }
        .container { background: white; padding: 30px; border-radius: 10px; max-width: 800px; margin: 0 auto; }
        .success { background: #d4edda; color: #155724; padding: 15px; border-radius: 8px; margin: 20px 0; }
        .stats { display: grid; grid-template-columns: repeat(3, 1fr); gap: 15px; margin: 20px 0; }
        .stat { background: #f0f2f5; padding: 15px; border-radius: 8px; text-align: center; }
        a { color: #667eea; text-decoration: none; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🎪 Demo Automática Completada</h1>
        
        <div class="success">
            <h3>✅ Procesamiento Exitoso</h3>
            <p>Se han procesado automáticamente transacciones de las 7 regiones LATAM</p>
        </div>
        
        <div class="stats">
            <div class="stat">
                <h3>` + std.len(demos) + `</h3>
                <p>Transacciones</p>
            </div>
            <div class="stat">
                <h3>7</h3>
                <p>Regiones</p>
            </div>
            <div class="stat">
                <h3>1</h3>
                <p>CFDI Demo</p>
            </div>
        </div>
        
        <p>Las transacciones han sido agregadas a la base de datos y están disponibles para reportes.</p>
        
        <p style="text-align: center; margin-top: 30px;">
            <a href="/">← Volver al inicio</a> | 
            <a href="/api/transacciones">Ver todas las transacciones</a> | 
            <a href="/dashboard">Ver Dashboard</a>
        </p>
    </div>
</body>
</html>`
}

// Handler API transacciones
func handleAPITransacciones(pathVars, method, body) {
    return JSON({
        total: std.len(database.transacciones),
        transacciones: database.transacciones
    })
}

// Handler API CFDIs
func handleAPICFDIs(pathVars, method, body) {
    return JSON({
        total: std.len(database.cfdis),
        cfdis: database.cfdis
    })
}

// Handler API regiones
func handleAPIRegiones(pathVars, method, body) {
    return JSON(database.regiones)
}

// Handler Dashboard
func handleDashboard(pathVars, method, body) {
    // Calcular estadísticas
    let totalTransacciones = std.len(database.transacciones)
    let totalBase = 0
    let totalImpuestos = 0
    let transaccionesPorRegion = {}
    
    let i = 0
    while (i < totalTransacciones) {
        let tx = database.transacciones[i]
        totalBase = totalBase + tx.amount
        totalImpuestos = totalImpuestos + tx.tax
        
        if (!transaccionesPorRegion[tx.region]) {
            transaccionesPorRegion[tx.region] = 0
        }
        transaccionesPorRegion[tx.region] = transaccionesPorRegion[tx.region] + 1
        
        i = i + 1
    }
    
    let html = `<!DOCTYPE html>
<html>
<head>
    <title>Dashboard - Sistema Contable LATAM</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; }
        h1 { color: #2c3e50; }
        .stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 20px; margin: 20px 0; }
        .stat-card { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .stat-value { font-size: 2rem; font-weight: bold; color: #667eea; }
        .chart { background: white; padding: 20px; border-radius: 8px; margin: 20px 0; }
        .report-buttons { display: flex; gap: 10px; flex-wrap: wrap; margin: 20px 0; }
        button { background: #667eea; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; }
        a { color: #667eea; text-decoration: none; }
        #reportResult { background: white; padding: 20px; border-radius: 8px; margin: 20px 0; display: none; }
        pre { background: #f4f4f4; padding: 10px; border-radius: 4px; overflow-x: auto; }
    </style>
</head>
<body>
    <div class="container">
        <h1>📊 Dashboard Sistema Contable LATAM</h1>
        
        <div class="stats">
            <div class="stat-card">
                <h3>Total Transacciones</h3>
                <div class="stat-value">` + totalTransacciones + `</div>
            </div>
            <div class="stat-card">
                <h3>Total Base</h3>
                <div class="stat-value">$` + math.round(totalBase) + `</div>
            </div>
            <div class="stat-card">
                <h3>Total Impuestos</h3>
                <div class="stat-value">$` + math.round(totalImpuestos) + `</div>
            </div>
            <div class="stat-card">
                <h3>CFDIs Procesados</h3>
                <div class="stat-value">` + std.len(database.cfdis) + `</div>
            </div>
        </div>
        
        <div class="chart">
            <h2>Reportes Rápidos DSL</h2>
            <div class="report-buttons">
                <button onclick="generarReporte('reporte todo ALL')">Reporte General</button>
                <button onclick="generarReporte('reporte ventas ALL')">Solo Ventas</button>
                <button onclick="generarReporte('reporte compras ALL')">Solo Compras</button>
                <button onclick="generarReporte('reporte ventas MX')">Ventas México</button>
                <button onclick="generarReporte('reporte ventas COL')">Ventas Colombia</button>
                <button onclick="generarReporte('reporte todo AR')">Todo Argentina</button>
            </div>
            
            <div id="reportResult"></div>
        </div>
        
        <p style="text-align: center; margin-top: 30px;">
            <a href="/">← Volver al inicio</a> | 
            <a href="/demo-auto">Ejecutar Demo Automática</a>
        </p>
    </div>
    
    <script>
        async function generarReporte(comando) {
            const formData = new URLSearchParams();
            formData.append('comando', comando);
            
            const response = await fetch('/generar-reporte', {
                method: 'POST',
                headers: {'Content-Type': 'application/x-www-form-urlencoded'},
                body: formData
            });
            
            const html = await response.text();
            // Extraer el JSON del pre tag
            const match = html.match(/<pre>(.*?)<\/pre>/s);
            if (match) {
                const resultDiv = document.getElementById('reportResult');
                resultDiv.style.display = 'block';
                resultDiv.innerHTML = '<h3>Resultado: ' + comando + '</h3><pre>' + match[1] + '</pre>';
            }
        }
    </script>
</body>
</html>`
    
    return html
}

// Función principal
func main() {
    console.log("✅ Configurando rutas...")
    
    // Rutas principales
    http.handler("GET", "/", handleHome)
    http.handler("POST", "/procesar-transaccion", handleProcesarTransaccion)
    http.handler("POST", "/procesar-cfdi", handleProcesarCFDI)
    http.handler("POST", "/generar-reporte", handleGenerarReporte)
    http.handler("GET", "/demo-auto", handleDemoAuto)
    http.handler("GET", "/dashboard", handleDashboard)
    
    // APIs
    http.handler("GET", "/api/transacciones", handleAPITransacciones)
    http.handler("GET", "/api/cfdis", handleAPICFDIs)
    http.handler("GET", "/api/regiones", handleAPIRegiones)
    
    console.log("✅ Sistema listo!")
    console.log("🌐 URL: http://localhost:8080")
    console.log("🎪 Demo automática: http://localhost:8080/demo-auto")
    console.log("📊 Dashboard: http://localhost:8080/dashboard")
    console.log("")
    console.log("🎯 POC Lista para Siigo!")
    
    // Iniciar servidor
    http.serve(":8080")
}

// Ejecutar
main()