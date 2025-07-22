// POC SIIGO - Sistema Contable LATAM Funcional
// Con CFDI y DSL Reportes

console.log("🌍 Sistema Contable LATAM - Demo Siigo")
console.log("📄 CFDI + 📊 DSL Reportes")
console.log("")

// Base de datos
let regiones = {
    "MX": { nombre: "México", moneda: "MXN", simbolo: "$", iva: 0.16, normativa: "SAT-CFDI" },
    "COL": { nombre: "Colombia", moneda: "COP", simbolo: "$", iva: 0.19, normativa: "DIAN" },
    "AR": { nombre: "Argentina", moneda: "ARS", simbolo: "$", iva: 0.21, normativa: "AFIP" },
    "CH": { nombre: "Chile", moneda: "CLP", simbolo: "$", iva: 0.19, normativa: "SII" },
    "UY": { nombre: "Uruguay", moneda: "UYU", simbolo: "$", iva: 0.22, normativa: "DGI" },
    "EC": { nombre: "Ecuador", moneda: "USD", simbolo: "$", iva: 0.12, normativa: "SRI" },
    "PE": { nombre: "Perú", moneda: "PEN", simbolo: "S/", iva: 0.18, normativa: "SUNAT" }
}

let transacciones = []
let cfdis = []

// DSL Reportes
dsl ReportesContables {
    token("REPORTE", "reporte")
    token("TIPO", "ventas|compras|todo")
    token("REGION", "MX|COL|AR|CH|UY|EC|PE|ALL")
    
    rule("generar", ["REPORTE", "TIPO", "REGION"], "generarReporte")
    
    func generarReporte(cmd, tipo, region) {
        let total = 0
        let count = 0
        let i = 0
        while (i < std.len(transacciones)) {
            let tx = transacciones[i]
            if ((region == "ALL" || tx.region == region) && (tipo == "todo" || tx.tipo == tipo)) {
                total = total + tx.total
                count = count + 1
            }
            i = i + 1
        }
        
        return {
            tipo: tipo,
            region: region,
            transacciones: count,
            total: math.round(total * 100) / 100
        }
    }
}

// Funciones auxiliares
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
                params[parts[0]] = std.replace(parts[1], "+", " ")
            }
        }
        i = i + 1
    }
    return params
}

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
        normativa: config.normativa,
        fecha: std.now()
    }
    
    transacciones[std.len(transacciones)] = tx
    return tx
}

func procesarCFDI(datos) {
    try {
        let cfdi = JSON.parse(datos)
        let comp = cfdi.Comprobante
        
        let registro = {
            uuid: comp.Complemento.TimbreFiscalDigital._UUID || "DEMO-" + std.now(),
            emisor: comp.Emisor._Nombre,
            receptor: comp.Receptor._Nombre,
            subtotal: std.parseFloat(comp._SubTotal),
            total: std.parseFloat(comp._Total),
            fecha: comp._Fecha,
            tipo: "CFDI"
        }
        
        cfdis[std.len(cfdis)] = registro
        
        // También crear transacción
        let tx = procesarTransaccion("cfdi", "MX", comp._SubTotal)
        tx.emisor = registro.emisor
        tx.uuid = registro.uuid
        
        return registro
    } catch {
        return {error: "CFDI inválido"}
    }
}

// Handlers HTTP
func handleHome(pathVars, method, body) {
    return `<!DOCTYPE html>
<html>
<head>
    <title>Sistema Contable LATAM - Siigo Demo</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 30px; border-radius: 10px; }
        h1 { color: #2c3e50; text-align: center; }
        .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; margin: 20px 0; }
        .card { background: #f8f9fa; padding: 20px; border-radius: 8px; border: 1px solid #ddd; }
        input, select, textarea { width: 100%; padding: 8px; margin: 5px 0; border: 1px solid #ddd; border-radius: 4px; }
        button { background: #667eea; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; width: 100%; }
        button:hover { background: #5a67d8; }
        .stats { background: #e8f5e9; padding: 15px; border-radius: 8px; margin: 20px 0; text-align: center; }
        a { color: #667eea; text-decoration: none; margin: 0 10px; }
        .value { background: linear-gradient(135deg, #84fab0 0%, #8fd3f4 100%); padding: 20px; border-radius: 8px; margin: 20px 0; }
        textarea { min-height: 100px; font-family: monospace; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🌍 Sistema Contable LATAM - Demo Siigo</h1>
        
        <div class="stats">
            <h3>Estadísticas</h3>
            <p>Transacciones: <strong>` + std.len(transacciones) + `</strong> | 
               CFDIs: <strong>` + std.len(cfdis) + `</strong> | 
               Regiones: <strong>7</strong></p>
        </div>
        
        <div class="grid">
            <div class="card">
                <h3>💰 Procesar Transacción</h3>
                <form action="/procesar" method="POST">
                    <label>Tipo:</label>
                    <select name="tipo">
                        <option value="ventas">Venta</option>
                        <option value="compras">Compra</option>
                    </select>
                    
                    <label>Región:</label>
                    <select name="region">
                        <option value="MX">México (16%)</option>
                        <option value="COL">Colombia (19%)</option>
                        <option value="AR">Argentina (21%)</option>
                        <option value="CH">Chile (19%)</option>
                        <option value="UY">Uruguay (22%)</option>
                        <option value="EC">Ecuador (12%)</option>
                        <option value="PE">Perú (18%)</option>
                    </select>
                    
                    <label>Importe:</label>
                    <input type="number" name="importe" value="100000" required>
                    
                    <button type="submit">Procesar</button>
                </form>
            </div>
            
            <div class="card">
                <h3>📄 Cargar CFDI</h3>
                <form action="/cfdi" method="POST">
                    <label>JSON CFDI:</label>
                    <textarea name="cfdi" placeholder='{"Comprobante": {...}}'></textarea>
                    <button type="submit">Procesar CFDI</button>
                    <p><a href="#" onclick="cargarEjemplo()">Cargar ejemplo</a></p>
                </form>
            </div>
            
            <div class="card">
                <h3>📊 Generar Reporte DSL</h3>
                <form action="/reporte" method="POST">
                    <label>Comando:</label>
                    <input type="text" name="comando" value="reporte todo ALL" style="font-family: monospace;">
                    <p>Ejemplos:<br>
                    <code>reporte ventas MX</code><br>
                    <code>reporte compras ALL</code><br>
                    <code>reporte todo COL</code></p>
                    <button type="submit">Generar</button>
                </form>
            </div>
        </div>
        
        <div style="text-align: center; margin: 20px 0;">
            <a href="/demo">Demo Auto</a> |
            <a href="/api/transacciones">API Transacciones</a> |
            <a href="/api/regiones">API Regiones</a>
        </div>
        
        <div class="value">
            <h3>🎯 Value Proposition Siigo</h3>
            <p><strong>18 meses → 2 meses</strong> | 
               <strong>$500K → $150K</strong> | 
               <strong>7 sistemas → 1 DSL</strong> | 
               <strong>ROI: 1,020%</strong></p>
        </div>
    </div>
    
    <script>
        function cargarEjemplo() {
            const cfdi = {
                Comprobante: {
                    Emisor: {_Nombre: "EMPRESA DEMO SA", _Rfc: "DEMO010101XXX"},
                    Receptor: {_Nombre: "PUBLICO GENERAL", _Rfc: "XAXX010101000"},
                    _SubTotal: "1000.00",
                    _Total: "1160.00",
                    _Fecha: "2025-07-22T12:00:00",
                    Complemento: {
                        TimbreFiscalDigital: {
                            _UUID: "DEMO-" + Date.now()
                        }
                    }
                }
            };
            document.querySelector('textarea[name="cfdi"]').value = JSON.stringify(cfdi, null, 2);
        }
    </script>
</body>
</html>`
}

func handleProcesar(pathVars, method, body) {
    let params = parseFormData(body)
    let tx = procesarTransaccion(
        params.tipo || "ventas",
        params.region || "COL", 
        params.importe || "100000"
    )
    
    return `<!DOCTYPE html>
<html>
<head>
    <title>Comprobante</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .comprobante { max-width: 600px; margin: 0 auto; background: white; padding: 30px; border: 2px solid #667eea; border-radius: 10px; }
        .row { display: flex; justify-content: space-between; margin: 10px 0; }
        .total { background: #f5f5f5; padding: 15px; border-radius: 5px; margin: 20px 0; }
        .success { color: #27ae60; text-align: center; font-weight: bold; }
    </style>
</head>
<body>
    <div class="comprobante">
        <h1 style="text-align: center;">COMPROBANTE</h1>
        <h2 style="text-align: center;">` + tx.pais + `</h2>
        
        <div class="row">
            <strong>ID:</strong> <span>` + tx.id + `</span>
        </div>
        <div class="row">
            <strong>Tipo:</strong> <span>` + tx.tipo + `</span>
        </div>
        <div class="row">
            <strong>Fecha:</strong> <span>` + tx.fecha + `</span>
        </div>
        <div class="row">
            <strong>Normativa:</strong> <span>` + tx.normativa + `</span>
        </div>
        
        <div class="total">
            <div class="row">
                <strong>Importe:</strong> <span>` + tx.moneda + ` ` + tx.importe + `</span>
            </div>
            <div class="row">
                <strong>IVA:</strong> <span>` + tx.moneda + ` ` + tx.iva + `</span>
            </div>
            <div class="row" style="border-top: 2px solid #333; padding-top: 10px;">
                <strong>TOTAL:</strong> <strong>` + tx.moneda + ` ` + tx.total + `</strong>
            </div>
        </div>
        
        <p class="success">✅ TRANSACCIÓN PROCESADA</p>
        
        <p style="text-align: center;"><a href="/">← Volver</a></p>
    </div>
</body>
</html>`
}

func handleCFDI(pathVars, method, body) {
    let params = parseFormData(body)
    let resultado = procesarCFDI(params.cfdi || "{}")
    
    if (resultado.error) {
        return `<h1>Error</h1><p>` + resultado.error + `</p><a href="/">Volver</a>`
    }
    
    return `<!DOCTYPE html>
<html>
<head>
    <title>CFDI Procesado</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 30px; }
        .info { background: #e8f5e9; padding: 20px; border-radius: 8px; }
        pre { background: #f5f5f5; padding: 15px; border-radius: 5px; overflow-x: auto; }
    </style>
</head>
<body>
    <div class="container">
        <h1>✅ CFDI Procesado</h1>
        
        <div class="info">
            <p><strong>UUID:</strong> ` + resultado.uuid + `</p>
            <p><strong>Emisor:</strong> ` + resultado.emisor + `</p>
            <p><strong>Receptor:</strong> ` + resultado.receptor + `</p>
            <p><strong>Total:</strong> $` + resultado.total + `</p>
            <p><strong>Fecha:</strong> ` + resultado.fecha + `</p>
        </div>
        
        <h3>Datos completos:</h3>
        <pre>` + JSON(resultado) + `</pre>
        
        <p style="text-align: center;"><a href="/">← Volver</a></p>
    </div>
</body>
</html>`
}

func handleReporte(pathVars, method, body) {
    let params = parseFormData(body)
    let comando = params.comando || "reporte todo ALL"
    
    let engine = ReportesContables
    let resultado = engine.use(comando)
    
    return `<!DOCTYPE html>
<html>
<head>
    <title>Reporte</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 30px; }
        .resultado { background: #f0f2f5; padding: 20px; border-radius: 8px; }
        pre { background: #f5f5f5; padding: 15px; border-radius: 5px; overflow-x: auto; }
    </style>
</head>
<body>
    <div class="container">
        <h1>📊 Reporte Generado</h1>
        <p><strong>Comando DSL:</strong> <code>` + comando + `</code></p>
        
        <div class="resultado">
            <h3>Resumen:</h3>
            <p><strong>Tipo:</strong> ` + resultado.tipo + `</p>
            <p><strong>Región:</strong> ` + resultado.region + `</p>
            <p><strong>Transacciones:</strong> ` + resultado.transacciones + `</p>
            <p><strong>Total:</strong> $` + resultado.total + `</p>
        </div>
        
        <h3>JSON:</h3>
        <pre>` + JSON(resultado) + `</pre>
        
        <p style="text-align: center;"><a href="/">← Volver</a></p>
    </div>
</body>
</html>`
}

func handleDemo(pathVars, method, body) {
    // Generar transacciones de demo
    let demos = [
        ["ventas", "MX", "100000"],
        ["compras", "COL", "50000"],
        ["ventas", "AR", "75000"],
        ["compras", "PE", "60000"]
    ]
    
    let i = 0
    while (i < std.len(demos)) {
        let d = demos[i]
        procesarTransaccion(d[0], d[1], d[2])
        i = i + 1
    }
    
    // CFDI demo
    let cfdiDemo = `{
        "Comprobante": {
            "Emisor": {"_Nombre": "DEMO SA", "_Rfc": "DEMO123"},
            "Receptor": {"_Nombre": "PUBLICO", "_Rfc": "XAXX010101000"},
            "_SubTotal": "5000.00",
            "_Total": "5800.00",
            "_Fecha": "2025-07-22",
            "Complemento": {
                "TimbreFiscalDigital": {"_UUID": "DEMO-AUTO"}
            }
        }
    }`
    procesarCFDI(cfdiDemo)
    
    return `<!DOCTYPE html>
<html>
<head>
    <title>Demo Automática</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background: #f5f5f5; }
        .container { max-width: 600px; margin: 0 auto; background: white; padding: 30px; border-radius: 10px; }
        .success { background: #d4edda; color: #155724; padding: 20px; border-radius: 8px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🎪 Demo Automática</h1>
        
        <div class="success">
            <h3>✅ Completado</h3>
            <p>Se procesaron:</p>
            <ul>
                <li>4 transacciones de diferentes regiones</li>
                <li>1 CFDI de ejemplo</li>
            </ul>
        </div>
        
        <p>Total de transacciones en el sistema: <strong>` + std.len(transacciones) + `</strong></p>
        
        <p style="text-align: center; margin-top: 30px;">
            <a href="/">← Volver</a> | 
            <a href="/api/transacciones">Ver todas</a>
        </p>
    </div>
</body>
</html>`
}

func handleAPITransacciones(pathVars, method, body) {
    return JSON({
        total: std.len(transacciones),
        transacciones: transacciones
    })
}

func handleAPIRegiones(pathVars, method, body) {
    return JSON(regiones)
}

// Configurar servidor
console.log("✅ Configurando rutas...")

http.handler("GET", "/", handleHome)
http.handler("POST", "/procesar", handleProcesar)
http.handler("POST", "/cfdi", handleCFDI)
http.handler("POST", "/reporte", handleReporte)
http.handler("GET", "/demo", handleDemo)
http.handler("GET", "/api/transacciones", handleAPITransacciones)
http.handler("GET", "/api/regiones", handleAPIRegiones)

console.log("🌐 Servidor listo en http://localhost:8080")
console.log("🎪 Demo auto en http://localhost:8080/demo")
console.log("")

// Iniciar servidor
http.serve(":8080")