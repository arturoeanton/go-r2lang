// Base de Datos Simplificada para Sistema Contable LATAM
// Simulación sin dependencias externas

console.log("📊 Módulo Database Simple cargado")

// Almacenamiento en memoria para demo
let transactions = []
let regions = {
    "MX": {code: "MX", name: "México", currency: "MXN", tax: 0.16},
    "COL": {code: "COL", name: "Colombia", currency: "COP", tax: 0.19},
    "AR": {code: "AR", name: "Argentina", currency: "ARS", tax: 0.21},
    "CH": {code: "CH", name: "Chile", currency: "CLP", tax: 0.19},
    "UY": {code: "UY", name: "Uruguay", currency: "UYU", tax: 0.22},
    "EC": {code: "EC", name: "Ecuador", currency: "USD", tax: 0.12},
    "PE": {code: "PE", name: "Perú", currency: "PEN", tax: 0.18}
}

// Inicializar base de datos (simulado)
func initDatabase() {
    console.log("✅ Base de datos en memoria inicializada")
    console.log("📋 Regiones configuradas: " + std.len(regions))
    return true
}

// Guardar transacción
func saveTransaction(transaction) {
    let id = std.len(transactions) + 1
    transaction.id = id
    transaction.timestamp = std.now()
    
    transactions[std.len(transactions)] = transaction
    
    console.log("💾 Transacción guardada: " + transaction.transactionId)
    return transaction
}

// Obtener todas las transacciones
func getTransactions() {
    return transactions
}

// Obtener estadísticas
func getTransactionStats() {
    let stats = {
        total_transactions: std.len(transactions),
        regions_active: std.len(regions),
        last_update: std.now()
    }
    
    return stats
}

// Obtener configuración regional
func getRegionConfig(regionCode) {
    if (regions[regionCode]) {
        return regions[regionCode]
    }
    return false
}

// Test de base de datos
func testDatabase() {
    console.log("🧪 Testing Database Module...")
    
    initDatabase()
    
    let testTx = {
        transactionId: "TEST-001",
        region: "COL", 
        amount: 1000,
        type: "sale"
    }
    
    saveTransaction(testTx)
    let stats = getTransactionStats()
    
    console.log("✅ Database test completado - " + stats.total_transactions + " transacciones")
    return true
}

console.log("📊 Database Simple ready")