// DSL Contable LATAM - Motor para Localización ERP Siigo
// Regiones: México, Colombia, Argentina, Chile, Uruguay, Ecuador, Perú

// Configuración de regiones LATAM
let regiones = {
    "MX": {
        nombre: "México",
        moneda: "MXN",
        simbolo: "$",
        iva: 0.16,
        normativa: "NIF-Mexican",
        cuentas: {
            cliente: "110001",
            ventas: "400001", 
            iva_debito: "210001",
            proveedor: "200001",
            compras: "500001",
            iva_credito: "110010"
        }
    },
    "COL": {
        nombre: "Colombia", 
        moneda: "COP",
        simbolo: "$",
        iva: 0.19,
        normativa: "NIIF-Colombia",
        cuentas: {
            cliente: "130501",
            ventas: "413501",
            iva_debito: "240801", 
            proveedor: "220501",
            compras: "613501",
            iva_credito: "135516"
        }
    },
    "AR": {
        nombre: "Argentina",
        moneda: "ARS", 
        simbolo: "$",
        iva: 0.21,
        normativa: "RT-Argentina",
        cuentas: {
            cliente: "112001",
            ventas: "401001",
            iva_debito: "213001",
            proveedor: "201001", 
            compras: "501001",
            iva_credito: "118001"
        }
    },
    "CH": {
        nombre: "Chile",
        moneda: "CLP",
        simbolo: "$",
        iva: 0.19,
        normativa: "IFRS-Chile",
        cuentas: {
            cliente: "113001", 
            ventas: "411001",
            iva_debito: "214001",
            proveedor: "202001",
            compras: "511001",
            iva_credito: "119001"
        }
    },
    "UY": {
        nombre: "Uruguay",
        moneda: "UYU",
        simbolo: "$",
        iva: 0.22,
        normativa: "NIIF-Uruguay", 
        cuentas: {
            cliente: "114001",
            ventas: "421001", 
            iva_debito: "215001",
            proveedor: "203001",
            compras: "521001",
            iva_credito: "120001"
        }
    },
    "EC": {
        nombre: "Ecuador",
        moneda: "USD",
        simbolo: "$",
        iva: 0.12,
        normativa: "NIIF-Ecuador",
        cuentas: {
            cliente: "115001",
            ventas: "431001",
            iva_debito: "216001", 
            proveedor: "204001", 
            compras: "531001",
            iva_credito: "121001"
        }
    },
    "PE": {
        nombre: "Perú",
        moneda: "PEN", 
        simbolo: "S/",
        iva: 0.18,
        normativa: "PCGE-Peru",
        cuentas: {
            cliente: "116001",
            ventas: "441001",
            iva_debito: "217001",
            proveedor: "205001",
            compras: "541001", 
            iva_credito: "122001"
        }
    }
}

// Función utilitaria para formato monetario
func formatMoney(amount, region) {
    let config = regiones[region]
    let rounded = math.round(amount * 100) / 100
    return config.simbolo + " " + rounded + " " + config.moneda
}

// Función para validar región
func validateRegion(region) {
    if (!regiones[region]) {
        panic("Región no soportada: " + region + ". Regiones disponibles: MX, COL, AR, CH, UY, EC, PE")
    }
    return true
}

// Función para generar ID de transacción
func generateTxId(region) {
    let timestamp = std.now()
    let random = math.randomInt(9999)
    return region + "-" + timestamp + "-" + random
}

// DSL Ventas LATAM
dsl VentasLATAM {
    token("VENTA", "venta|sale")
    token("REGION", "MX|COL|AR|CH|UY|EC|PE")
    token("IMPORTE", "[0-9]+\\.?[0-9]*")
    
    rule("venta_latam", ["VENTA", "REGION", "IMPORTE"], "procesarVentaLATAM")
    
    func procesarVentaLATAM(operacion, region, importe) {
        validateRegion(region)
        let config = regiones[region]
        let importeNum = std.parseFloat(importe)
        
        if (importeNum <= 0) {
            panic("El importe debe ser mayor a cero")
        }
        
        let importeIVA = math.round((importeNum * config.iva) * 100) / 100
        let importeTotal = math.round((importeNum + importeIVA) * 100) / 100
        let txId = generateTxId(region)
        
        console.log("=== COMPROBANTE DE VENTA " + config.nombre + " ===")
        console.log("ID Transacción: " + txId)
        console.log("Región: " + region + " - " + config.nombre)
        console.log("Fecha: " + std.now())
        console.log("Normativa: " + config.normativa)
        console.log("")
        console.log("ASIENTO CONTABLE:")
        console.log("DEBE:")
        console.log("  " + config.cuentas.cliente + " - Clientes: " + formatMoney(importeTotal, region))
        console.log("HABER:")
        console.log("  " + config.cuentas.ventas + " - Ventas: " + formatMoney(importeNum, region))
        console.log("  " + config.cuentas.iva_debito + " - IVA Débito: " + formatMoney(importeIVA, region))
        console.log("")
        console.log("Tasa IVA: " + (config.iva * 100) + "%")
        console.log("Estado: VALIDADO ✓")
        
        return {
            success: true,
            transactionId: txId,
            region: region,
            amount: importeNum,
            tax: importeIVA, 
            total: importeTotal,
            currency: config.moneda,
            accounts: {
                debit: [{account: config.cuentas.cliente, amount: importeTotal}],
                credit: [
                    {account: config.cuentas.ventas, amount: importeNum},
                    {account: config.cuentas.iva_debito, amount: importeIVA}
                ]
            },
            compliance: config.normativa
        }
    }
}

// DSL Compras LATAM  
dsl ComprasLATAM {
    token("COMPRA", "compra|purchase")
    token("REGION", "MX|COL|AR|CH|UY|EC|PE")
    token("IMPORTE", "[0-9]+\\.?[0-9]*")
    
    rule("compra_latam", ["COMPRA", "REGION", "IMPORTE"], "procesarCompraLATAM")
    
    func procesarCompraLATAM(operacion, region, importe) {
        validateRegion(region)
        let config = regiones[region]
        let importeNum = std.parseFloat(importe)
        
        if (importeNum <= 0) {
            panic("El importe debe ser mayor a cero")
        }
        
        let importeIVA = math.round((importeNum * config.iva) * 100) / 100
        let importeTotal = math.round((importeNum + importeIVA) * 100) / 100
        let txId = generateTxId(region)
        
        console.log("=== COMPROBANTE DE COMPRA " + config.nombre + " ===")
        console.log("ID Transacción: " + txId)
        console.log("Región: " + region + " - " + config.nombre)  
        console.log("Fecha: " + std.now())
        console.log("Normativa: " + config.normativa)
        console.log("")
        console.log("ASIENTO CONTABLE:")
        console.log("DEBE:")
        console.log("  " + config.cuentas.compras + " - Compras: " + formatMoney(importeNum, region))
        console.log("  " + config.cuentas.iva_credito + " - IVA Crédito: " + formatMoney(importeIVA, region))
        console.log("HABER:")
        console.log("  " + config.cuentas.proveedor + " - Proveedores: " + formatMoney(importeTotal, region))
        console.log("")
        console.log("Tasa IVA: " + (config.iva * 100) + "%")
        console.log("Estado: VALIDADO ✓")
        
        return {
            success: true,
            transactionId: txId,
            region: region,
            amount: importeNum,
            tax: importeIVA,
            total: importeTotal, 
            currency: config.moneda,
            accounts: {
                debit: [
                    {account: config.cuentas.compras, amount: importeNum},
                    {account: config.cuentas.iva_credito, amount: importeIVA}
                ],
                credit: [{account: config.cuentas.proveedor, amount: importeTotal}]
            },
            compliance: config.normativa
        }
    }
}

// DSL Consultas LATAM
dsl ConsultasLATAM {
    token("CONSULTAR", "consultar|query")
    token("REGION", "MX|COL|AR|CH|UY|EC|PE")
    token("TIPO", "config|configuration|plan")
    
    rule("consulta_region", ["CONSULTAR", "TIPO", "REGION"], "consultarRegion")
    
    func consultarRegion(consultar, tipo, region) {
        validateRegion(region)
        let config = regiones[region]
        
        console.log("=== CONFIGURACIÓN REGIONAL " + config.nombre + " ===")
        console.log("Región: " + region)
        console.log("País: " + config.nombre)
        console.log("Moneda: " + config.moneda + " (" + config.simbolo + ")")
        console.log("Tasa IVA: " + (config.iva * 100) + "%")
        console.log("Normativa: " + config.normativa)
        console.log("")
        console.log("PLAN DE CUENTAS:")
        console.log("  Clientes: " + config.cuentas.cliente)
        console.log("  Ventas: " + config.cuentas.ventas) 
        console.log("  IVA Débito: " + config.cuentas.iva_debito)
        console.log("  Proveedores: " + config.cuentas.proveedor)
        console.log("  Compras: " + config.cuentas.compras)
        console.log("  IVA Crédito: " + config.cuentas.iva_credito)
        console.log("")
        console.log("Estado: CONFIGURACIÓN ACTIVA ✓")
        
        return {
            success: true,
            region: region,
            country: config.nombre,
            currency: config.moneda,
            taxRate: config.iva,
            compliance: config.normativa,
            chartOfAccounts: config.cuentas
        }
    }
}

// Función principal para testing  
func testDSLLatam() {
    console.log("=== DEMO DSL CONTABLE LATAM ===")
    console.log("Motor de Localización ERP para Siigo")
    console.log("")
    
    let motorVentas = VentasLATAM
    let motorCompras = ComprasLATAM  
    let motorConsultas = ConsultasLATAM
    
    // Test todas las regiones
    let regionesTest = ["MX", "COL", "AR", "CH", "UY", "EC", "PE"]
    
    let i = 0
    while (i < std.len(regionesTest)) {
        let region = regionesTest[i]
        console.log("--- TESTING REGIÓN " + region + " ---")
        
        // Test venta
        let venta = motorVentas.use("venta " + region + " 10000")
        
        // Test compra
        let compra = motorCompras.use("compra " + region + " 5000")
        
        // Test configuración
        let config = motorConsultas.use("consultar config " + region)
        
        console.log("")
        i = i + 1
    }
    
    console.log("=== DEMO COMPLETADO ===")
    console.log("Todas las regiones LATAM funcionando correctamente")
}

// Para ejecutar test, usar: testDSLLatam()
// testDSLLatam()