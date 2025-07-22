// Test Simple de Endpoints - Verificar que funcionen todos

console.log("🧪 TEST ENDPOINTS - Sistema Contable LATAM")
console.log("=========================================")

// Test básico de handlers
func testIndex(pathVars, method, body) {
    return "✅ Página principal funcionando"
}

func testVenta(pathVars, method, body) {
    return "✅ Endpoint venta funcionando - Body: " + (body || "empty")
}

func testCompra(pathVars, method, body) {
    return "✅ Endpoint compra funcionando - Body: " + (body || "empty") 
}

func testAPI(pathVars, method, body) {
    return JSON({status: "ok", message: "API funcionando", timestamp: std.now()})
}

func testDemo(pathVars, method, body) {
    return `<!DOCTYPE html>
<html>
<head><title>Demo Test</title></head>
<body>
    <h1>✅ Demo endpoint funcionando</h1>
    <p>Todos los endpoints están operativos</p>
    <a href="/">Volver al inicio</a>
</body>
</html>`
}

func main() {
    console.log("🚀 Iniciando test server...")
    
    // Configurar rutas básicas
    http.handler("GET", "/", testIndex)
    http.handler("POST", "/procesar-venta", testVenta)
    http.handler("POST", "/procesar-compra", testCompra)
    http.handler("GET", "/api/test", testAPI)
    http.handler("GET", "/demo", testDemo)
    
    console.log("✅ Rutas de test configuradas")
    console.log("🌐 URL: http://localhost:8080")
    console.log("🧪 Test endpoints:")
    console.log("   GET  / - Test página principal")
    console.log("   POST /procesar-venta - Test venta")
    console.log("   POST /procesar-compra - Test compra")
    console.log("   GET  /api/test - Test API")
    console.log("   GET  /demo - Test demo")
    console.log("")
    console.log("Listening...")
    
    // Iniciar servidor
    http.serve(":8080")
}