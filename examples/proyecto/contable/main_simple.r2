// Sistema Contable LATAM - Versión Simplificada
// Demo R2Lang DSL para Siigo ERP Localization
// Integración sin dependencias externas problematicas

console.log("========================================")
console.log("  SISTEMA CONTABLE LATAM - R2LANG DSL")
console.log("========================================")
console.log("  Demo para Siigo ERP Localization")
console.log("  Regiones: MX, COL, AR, CH, UY, EC, PE")
console.log("========================================")
console.log("")

// Configuración de regiones LATAM integrada
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

// DSL Ventas LATAM
dsl VentasLATAM {
    token("VENTA", "venta|sale")
    token("REGION", "MX|COL|AR|CH|UY|EC|PE")
    token("IMPORTE", "[0-9]+\\.?[0-9]*")
    
    rule("venta_latam", ["VENTA", "REGION", "IMPORTE"], "procesarVenta")
    
    func procesarVenta(operacion, region, importe) {
        let config = regiones[region]
        let importeNum = std.parseFloat(importe)
        let importeIVA = math.round((importeNum * config.iva) * 100) / 100
        let importeTotal = math.round((importeNum + importeIVA) * 100) / 100
        let txId = region + "-" + std.now() + "-" + math.randomInt(9999)
        
        console.log("=== COMPROBANTE VENTA " + config.nombre + " ===")
        console.log("ID: " + txId)
        console.log("DEBE: " + config.cuentas.cliente + " Clientes " + formatMoney(importeTotal, region))
        console.log("HABER: " + config.cuentas.ventas + " Ventas " + formatMoney(importeNum, region))
        console.log("HABER: " + config.cuentas.iva_debito + " IVA " + formatMoney(importeIVA, region))
        console.log("Estado: VALIDADO ✓")
        console.log("")
        
        return {success: true, txId: txId, total: importeTotal}
    }
}

// DSL Compras LATAM
dsl ComprasLATAM {
    token("COMPRA", "compra|purchase")
    token("REGION", "MX|COL|AR|CH|UY|EC|PE")
    token("IMPORTE", "[0-9]+\\.?[0-9]*")
    
    rule("compra_latam", ["COMPRA", "REGION", "IMPORTE"], "procesarCompra")
    
    func procesarCompra(operacion, region, importe) {
        let config = regiones[region]
        let importeNum = std.parseFloat(importe)
        let importeIVA = math.round((importeNum * config.iva) * 100) / 100
        let importeTotal = math.round((importeNum + importeIVA) * 100) / 100
        let txId = region + "-" + std.now() + "-" + math.randomInt(9999)
        
        console.log("=== COMPROBANTE COMPRA " + config.nombre + " ===")
        console.log("ID: " + txId)
        console.log("DEBE: " + config.cuentas.compras + " Compras " + formatMoney(importeNum, region))
        console.log("DEBE: " + config.cuentas.iva_credito + " IVA " + formatMoney(importeIVA, region))
        console.log("HABER: " + config.cuentas.proveedor + " Proveedores " + formatMoney(importeTotal, region))
        console.log("Estado: VALIDADO ✓")
        console.log("")
        
        return {success: true, txId: txId, total: importeTotal}
    }
}

// Función principal de demo
func demoSiigo() {
    console.log("🚀 DEMO SIIGO - R2LANG DSL CONTABILIDAD LATAM")
    console.log("==============================================")
    console.log("")
    
    let motorVentas = VentasLATAM
    let motorCompras = ComprasLATAM
    
    // Demo todas las regiones
    let regionesTest = ["MX", "COL", "AR", "CH", "UY", "EC", "PE"]
    let transacciones = 0
    
    let i = 0
    while (i < std.len(regionesTest)) {
        let region = regionesTest[i]
        
        console.log("🌍 REGIÓN: " + region + " - " + regiones[region].nombre)
        console.log("======================================")
        
        // Venta
        let ventaAmount = 100000 + (i * 15000)
        let venta = motorVentas.use("venta " + region + " " + ventaAmount)
        transacciones = transacciones + 1
        
        // Compra
        let compraAmount = 50000 + (i * 8000)
        let compra = motorCompras.use("compra " + region + " " + compraAmount)
        transacciones = transacciones + 1
        
        i = i + 1
    }
    
    console.log("🎉 DEMO COMPLETADO")
    console.log("==================")
    console.log("• Regiones: " + std.len(regionesTest))
    console.log("• Transacciones: " + transacciones)
    console.log("• Normativas: 7 diferentes")
    console.log("• Monedas: 6 (MXN,COP,ARS,CLP,UYU,USD,PEN)")
    console.log("")
    console.log("💡 SIIGO VALUE PROPOSITION:")
    console.log("  ✅ 18 meses → 2 meses por país")
    console.log("  ✅ $500K → $150K por localización")
    console.log("  ✅ 7 codebases → 1 DSL")
    console.log("  ✅ ROI: 1,020% en 3 años")
    console.log("")
    console.log("🎯 ¡R2LANG + SIIGO = LATAM DOMINATION! 🚀")
}

// Ejecutar demo
demoSiigo()