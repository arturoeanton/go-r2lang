// Base de Datos SQLite para Sistema Contable LATAM
// Usando r2db para gestión de datos

func initDatabase() {
    console.log("Inicializando base de datos SQLite...")
    
    // Conectar a base de datos
    let conn = db.connect("sqlite", "./database/contable_latam.db")
    
    // Crear tabla de transacciones
    let createTransactions = `
        CREATE TABLE IF NOT EXISTS transactions (
            id TEXT PRIMARY KEY,
            region TEXT NOT NULL,
            country TEXT NOT NULL,
            type TEXT NOT NULL,
            amount REAL NOT NULL,
            tax REAL NOT NULL,
            total REAL NOT NULL,
            currency TEXT NOT NULL,
            compliance TEXT NOT NULL,
            accounts TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            status TEXT DEFAULT 'ACTIVE'
        )
    `
    
    // Crear tabla de regiones
    let createRegions = `
        CREATE TABLE IF NOT EXISTS regions (
            code TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            currency TEXT NOT NULL,
            symbol TEXT NOT NULL,
            tax_rate REAL NOT NULL,
            compliance_standard TEXT NOT NULL,
            chart_of_accounts TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    `
    
    // Crear tabla de audit log
    let createAuditLog = `
        CREATE TABLE IF NOT EXISTS audit_log (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            transaction_id TEXT,
            action TEXT NOT NULL,
            details TEXT,
            user_agent TEXT,
            ip_address TEXT,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (transaction_id) REFERENCES transactions(id)
        )
    `
    
    // Crear tabla de configuraciones
    let createConfigurations = `
        CREATE TABLE IF NOT EXISTS configurations (
            key TEXT PRIMARY KEY,
            value TEXT NOT NULL,
            description TEXT,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    `
    
    // Ejecutar creación de tablas
    db.exec(conn, createTransactions)
    db.exec(conn, createRegions)  
    db.exec(conn, createAuditLog)
    db.exec(conn, createConfigurations)
    
    console.log("Tablas creadas exitosamente")
    
    // Insertar datos iniciales de regiones
    seedRegionalData(conn)
    
    db.close(conn)
    console.log("Base de datos inicializada correctamente")
    
    return true
}

func seedRegionalData(conn) {
    console.log("Insertando datos iniciales de regiones...")
    
    // Datos de regiones LATAM
    let regiones = [
        {
            code: "MX",
            name: "México", 
            currency: "MXN",
            symbol: "$",
            tax_rate: 0.16,
            compliance: "NIF-Mexican",
            accounts: `{"cliente":"110001","ventas":"400001","iva_debito":"210001","proveedor":"200001","compras":"500001","iva_credito":"110010"}`
        },
        {
            code: "COL",
            name: "Colombia",
            currency: "COP", 
            symbol: "$",
            tax_rate: 0.19,
            compliance: "NIIF-Colombia",
            accounts: `{"cliente":"130501","ventas":"413501","iva_debito":"240801","proveedor":"220501","compras":"613501","iva_credito":"135516"}`
        },
        {
            code: "AR", 
            name: "Argentina",
            currency: "ARS",
            symbol: "$",
            tax_rate: 0.21,
            compliance: "RT-Argentina", 
            accounts: `{"cliente":"112001","ventas":"401001","iva_debito":"213001","proveedor":"201001","compras":"501001","iva_credito":"118001"}`
        },
        {
            code: "CH",
            name: "Chile",
            currency: "CLP",
            symbol: "$", 
            tax_rate: 0.19,
            compliance: "IFRS-Chile",
            accounts: `{"cliente":"113001","ventas":"411001","iva_debito":"214001","proveedor":"202001","compras":"511001","iva_credito":"119001"}`
        },
        {
            code: "UY",
            name: "Uruguay",
            currency: "UYU",
            symbol: "$",
            tax_rate: 0.22,
            compliance: "NIIF-Uruguay",
            accounts: `{"cliente":"114001","ventas":"421001","iva_debito":"215001","proveedor":"203001","compras":"521001","iva_credito":"120001"}`
        },
        {
            code: "EC",
            name: "Ecuador", 
            currency: "USD",
            symbol: "$",
            tax_rate: 0.12,
            compliance: "NIIF-Ecuador",
            accounts: `{"cliente":"115001","ventas":"431001","iva_debito":"216001","proveedor":"204001","compras":"531001","iva_credito":"121001"}`
        },
        {
            code: "PE",
            name: "Perú",
            currency: "PEN",
            symbol: "S/",
            tax_rate: 0.18, 
            compliance: "PCGE-Peru",
            accounts: `{"cliente":"116001","ventas":"441001","iva_debito":"217001","proveedor":"205001","compras":"541001","iva_credito":"122001"}`
        }
    ]
    
    // Insertar cada región
    regiones.forEach(function(region) {
        let insertSQL = `
            INSERT OR REPLACE INTO regions 
            (code, name, currency, symbol, tax_rate, compliance_standard, chart_of_accounts)
            VALUES (?, ?, ?, ?, ?, ?, ?)
        `
        
        db.exec(conn, insertSQL, [
            region.code,
            region.name,
            region.currency,
            region.symbol, 
            region.tax_rate,
            region.compliance,
            region.accounts
        ])
    })
    
    // Insertar configuraciones iniciales
    let configs = [
        {key: "app_name", value: "Sistema Contable LATAM", description: "Nombre de la aplicación"},
        {key: "version", value: "1.0.0", description: "Versión del sistema"},
        {key: "demo_mode", value: "true", description: "Modo demostración activo"},
        {key: "supported_regions", value: "MX,COL,AR,CH,UY,EC,PE", description: "Regiones soportadas"},
        {key: "default_language", value: "es", description: "Idioma por defecto"}
    ]
    
    configs.forEach(function(config) {
        let configSQL = `
            INSERT OR REPLACE INTO configurations (key, value, description)
            VALUES (?, ?, ?)
        `
        db.exec(conn, configSQL, [config.key, config.value, config.description])
    })
    
    console.log("Datos iniciales insertados correctamente")
}

func saveTransaction(transaction) {
    let conn = db.connect("sqlite", "./database/contable_latam.db")
    
    let insertSQL = `
        INSERT INTO transactions 
        (id, region, country, type, amount, tax, total, currency, compliance, accounts)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
    
    let accountsJSON = JSON.stringify(transaction.accounts)
    
    db.exec(conn, insertSQL, [
        transaction.transactionId,
        transaction.region,
        transaction.country,
        transaction.type,
        transaction.amount,
        transaction.tax,
        transaction.total,
        transaction.currency,
        transaction.compliance,
        accountsJSON
    ])
    
    // Log de auditoría
    logAuditEvent(conn, transaction.transactionId, "CREATED", "Transaction created via DSL", "", "127.0.0.1")
    
    db.close(conn)
    
    console.log("Transacción guardada: " + transaction.transactionId)
    return true
}

func getTransactions(region) {
    let conn = db.connect("sqlite", "./database/contable_latam.db")
    
    let sql = "SELECT * FROM transactions WHERE region = ? ORDER BY created_at DESC"
    if (region == "ALL") {
        sql = "SELECT * FROM transactions ORDER BY created_at DESC"
    }
    
    let results = []
    if (region == "ALL") {
        results = db.query(conn, "SELECT * FROM transactions ORDER BY created_at DESC")
    } else {
        results = db.query(conn, sql, [region])
    }
    
    db.close(conn)
    return results
}

func getRegionConfig(regionCode) {
    let conn = db.connect("sqlite", "./database/contable_latam.db")
    
    let sql = "SELECT * FROM regions WHERE code = ?"
    let results = db.query(conn, sql, [regionCode])
    
    db.close(conn)
    
    if (results.length > 0) {
        return results[0] 
    }
    return null
}

func getAllRegions() {
    let conn = db.connect("sqlite", "./database/contable_latam.db")
    
    let results = db.query(conn, "SELECT * FROM regions ORDER BY name")
    
    db.close(conn)
    return results
}

func getTransactionStats() {
    let conn = db.connect("sqlite", "./database/contable_latam.db")
    
    let stats = {
        total_transactions: 0,
        total_amount: 0,
        by_region: {},
        by_type: {}
    }
    
    // Total de transacciones
    let totalResult = db.query(conn, "SELECT COUNT(*) as total FROM transactions")
    stats.total_transactions = totalResult[0].total
    
    // Total de importes
    let amountResult = db.query(conn, "SELECT SUM(amount) as total FROM transactions") 
    stats.total_amount = amountResult[0].total || 0
    
    // Por región
    let regionResults = db.query(conn, "SELECT region, COUNT(*) as count, SUM(amount) as total FROM transactions GROUP BY region")
    regionResults.forEach(function(row) {
        stats.by_region[row.region] = {
            count: row.count,
            total: row.total
        }
    })
    
    // Por tipo
    let typeResults = db.query(conn, "SELECT type, COUNT(*) as count FROM transactions GROUP BY type")
    typeResults.forEach(function(row) {
        stats.by_type[row.type] = row.count
    })
    
    db.close(conn)
    return stats
}

func logAuditEvent(conn, transactionId, action, details, userAgent, ipAddress) {
    let auditSQL = `
        INSERT INTO audit_log (transaction_id, action, details, user_agent, ip_address)
        VALUES (?, ?, ?, ?, ?)
    `
    
    db.exec(conn, auditSQL, [transactionId, action, details, userAgent, ipAddress])
}

func getAuditLog(limit) {
    let conn = db.connect("sqlite", "./database/contable_latam.db")
    
    let sql = "SELECT * FROM audit_log ORDER BY created_at DESC"
    if (limit) {
        sql = sql + " LIMIT " + limit
    }
    
    let results = db.query(conn, sql)
    db.close(conn)
    
    return results
}

// Función de test de base de datos
func testDatabase() {
    console.log("=== TEST DE BASE DE DATOS ===")
    
    // Inicializar
    initDatabase()
    
    // Test obtener regiones
    let regiones = getAllRegions()
    console.log("Regiones cargadas: " + regiones.length)
    
    // Test obtener configuración de región
    let configCOL = getRegionConfig("COL")
    console.log("Config Colombia: " + configCOL.name)
    
    // Test stats iniciales
    let stats = getTransactionStats()
    console.log("Stats iniciales - Total transacciones: " + stats.total_transactions)
    
    console.log("Test de base de datos completado ✓")
}

// Exportar funciones principales
func getDatabaseFunctions() {
    return {
        init: initDatabase,
        saveTransaction: saveTransaction,
        getTransactions: getTransactions,
        getRegionConfig: getRegionConfig, 
        getAllRegions: getAllRegions,
        getStats: getTransactionStats,
        test: testDatabase
    }
}