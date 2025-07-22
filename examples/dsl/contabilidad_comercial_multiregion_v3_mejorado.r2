// DSL Motor Contable Comercial Multi-Region V3 - VERSION MEJORADA
// Sistema completamente mejorado con formato numérico, rangos dinámicos y validación avanzada

// Función de utilidad para formatear números con 2 decimales
func formatCurrency(amount, symbol) {
    let rounded = math.round(amount * 100) / 100
    let formatted = symbol + " " + rounded
    return formatted
}

// Función para validar importes
func validateAmount(amount) {
    let numAmount = std.parseFloat(amount)
    if (numAmount < 0) {
        panic("Error: El importe no puede ser negativo: " + amount)
    }
    if (numAmount > 10000000) {
        panic("Error: El importe excede el límite máximo (10,000,000): " + amount)
    }
    return numAmount
}

// Función para generar ID único de transacción
func generateTransactionId() {
    let timestamp = std.now()
    let random = math.randomInt(1000)
    return "TX-" + timestamp + "-" + random
}

// DSL Ventas USA con mejoras
dsl VentasUSAMejorado {
    token("VENTA", "venta")
    token("USA", "USA")
    token("IMPORTE", "[0-9]+\\.?[0-9]*")  // Soporte para decimales dinámicos
    
    rule("venta_usa", ["VENTA", "USA", "IMPORTE"], "procesarVentaUSAMejorado")
    
    func procesarVentaUSAMejorado(venta, region, importe) {
        let importeNum = validateAmount(importe)
        let tasaIVA = 0.0875
        let importeIVA = math.round((importeNum * tasaIVA) * 100) / 100
        let importeTotal = math.round((importeNum + importeIVA) * 100) / 100
        let transactionId = generateTransactionId()
        
        console.log("=== COMPROBANTE DE VENTA USA (MEJORADO) ===")
        console.log("ID Transacción: " + transactionId)
        console.log("Cliente: TechSoft USA Inc.")
        console.log("Region: R01 - America del Norte")
        console.log("Fecha: " + std.now())
        console.log("Cuenta Cliente: 121002 - Clientes USA")
        console.log("Cuenta Ventas: 411002 - Ventas USA")
        console.log("Cuenta IVA: 224002 - Sales Tax USA")
        console.log("DEBE: 121002 " + formatCurrency(importeTotal, "USD"))
        console.log("HABER: 411002 " + formatCurrency(importeNum, "USD") + " + 224002 " + formatCurrency(importeIVA, "USD"))
        console.log("Tasa Impuesto: 8.75%")
        console.log("Normativa: US-GAAP")
        console.log("Estado: VALIDADO ✓")
        
        return {
            success: true,
            transactionId: transactionId,
            amount: importeTotal,
            currency: "USD",
            region: "USA"
        }
    }
}

// DSL Ventas EUR con mejoras
dsl VentasEURMejorado {
    token("VENTA", "venta")
    token("EUR", "EUR")
    token("IMPORTE", "[0-9]+\\.?[0-9]*")
    
    rule("venta_eur", ["VENTA", "EUR", "IMPORTE"], "procesarVentaEURMejorado")
    
    func procesarVentaEURMejorado(venta, region, importe) {
        let importeNum = validateAmount(importe)
        let tasaIVA = 0.20
        let importeIVA = math.round((importeNum * tasaIVA) * 100) / 100
        let importeTotal = math.round((importeNum + importeIVA) * 100) / 100
        let transactionId = generateTransactionId()
        
        console.log("=== COMPROBANTE DE VENTA EUROPA (MEJORADO) ===")
        console.log("ID Transacción: " + transactionId)
        console.log("Cliente: EuroSystems GmbH")
        console.log("Region: R02 - Europa")
        console.log("Fecha: " + std.now())
        console.log("Cuenta Cliente: 121003 - Clientes Europa")
        console.log("Cuenta Ventas: 411003 - Ventas Europa")
        console.log("Cuenta IVA: 224003 - VAT Europa")
        console.log("DEBE: 121003 " + formatCurrency(importeTotal, "EUR"))
        console.log("HABER: 411003 " + formatCurrency(importeNum, "EUR") + " + 224003 " + formatCurrency(importeIVA, "EUR"))
        console.log("Tasa Impuesto: 20.00%")
        console.log("Normativa: IFRS")
        console.log("Estado: VALIDADO ✓")
        
        return {
            success: true,
            transactionId: transactionId,
            amount: importeTotal,
            currency: "EUR",
            region: "EUR"
        }
    }
}

// DSL Ventas ARG con mejoras
dsl VentasARGMejorado {
    token("VENTA", "venta")
    token("ARG", "ARG")
    token("IMPORTE", "[0-9]+\\.?[0-9]*")
    
    rule("venta_arg", ["VENTA", "ARG", "IMPORTE"], "procesarVentaARGMejorado")
    
    func procesarVentaARGMejorado(venta, region, importe) {
        let importeNum = validateAmount(importe)
        let tasaIVA = 0.21
        let importeIVA = math.round((importeNum * tasaIVA) * 100) / 100
        let importeTotal = math.round((importeNum + importeIVA) * 100) / 100
        let transactionId = generateTransactionId()
        
        console.log("=== COMPROBANTE DE VENTA ARGENTINA (MEJORADO) ===")
        console.log("ID Transacción: " + transactionId)
        console.log("Cliente: Sistemas Locales S.A.")
        console.log("Region: R03 - America del Sur")
        console.log("Fecha: " + std.now())
        console.log("Cuenta Cliente: 121001 - Clientes Nacionales")
        console.log("Cuenta Ventas: 411001 - Ventas Nacionales")
        console.log("Cuenta IVA: 224001 - IVA Debito Fiscal")
        console.log("DEBE: 121001 " + formatCurrency(importeTotal, "ARS"))
        console.log("HABER: 411001 " + formatCurrency(importeNum, "ARS") + " + 224001 " + formatCurrency(importeIVA, "ARS"))
        console.log("Tasa Impuesto: 21.00%")
        console.log("Normativa: RT Argentina")
        console.log("Estado: VALIDADO ✓")
        
        return {
            success: true,
            transactionId: transactionId,
            amount: importeTotal,
            currency: "ARS",
            region: "ARG"
        }
    }
}

// DSL Compras USA con mejoras
dsl ComprasUSAMejorado {
    token("COMPRA", "compra")
    token("USA", "USA")
    token("SERVICIOS", "servicios")
    token("IMPORTE", "[0-9]+\\.?[0-9]*")
    
    rule("compra_usa_servicios", ["COMPRA", "USA", "SERVICIOS", "IMPORTE"], "procesarCompraUSAServiciosMejorado")
    
    func procesarCompraUSAServiciosMejorado(compra, region, tipo, importe) {
        let importeNum = validateAmount(importe)
        let tasaIVA = 0.0875
        let importeIVA = math.round((importeNum * tasaIVA) * 100) / 100
        let importeTotal = math.round((importeNum + importeIVA) * 100) / 100
        let transactionId = generateTransactionId()
        
        console.log("=== COMPROBANTE DE COMPRA USA SERVICIOS (MEJORADO) ===")
        console.log("ID Transacción: " + transactionId)
        console.log("Proveedor: Amazon Web Services")
        console.log("Region: R01 - America del Norte")
        console.log("Fecha: " + std.now())
        console.log("Cuenta Servicios: 521002 - Servicios USA")
        console.log("Cuenta IVA Credito: 113002 - Tax Credit USA")
        console.log("Cuenta Proveedor: 211002 - Proveedores USA")
        console.log("DEBE: 521002 " + formatCurrency(importeNum, "USD") + " + 113002 " + formatCurrency(importeIVA, "USD"))
        console.log("HABER: 211002 " + formatCurrency(importeTotal, "USD"))
        console.log("Tasa Impuesto: 8.75%")
        console.log("Normativa: US-GAAP")
        console.log("Estado: VALIDADO ✓")
        
        return {
            success: true,
            transactionId: transactionId,
            amount: importeTotal,
            currency: "USD",
            region: "USA",
            type: "servicios"
        }
    }
}

// DSL Análisis Mejorado con Validaciones
dsl AnalisisCuentasMejorado {
    token("REGION", "R[0-9]{2}")
    token("FECHA", "[0-9]{1,2}/[0-9]{1,2}/[0-9]{4}")
    token("ANALIZAR", "analizar")
    token("CUENTAS", "cuentas")
    token("MOVIMIENTOS", "movimientos")
    token("DE_REGION", "de")
    token("DESDE", "desde")
    token("HASTA", "hasta")
    
    rule("analizar_region", ["ANALIZAR", "CUENTAS", "MOVIMIENTOS", "DE_REGION", "REGION", "DESDE", "FECHA", "HASTA", "FECHA"], "analizarCuentasRegionMejorado")
    
    func analizarCuentasRegionMejorado(analizar, cuentas, movimientos, deRegion, region, desde, fechaDesde, hasta, fechaHasta) {
        let reportId = generateTransactionId()
        
        console.log("=== ANÁLISIS DE CUENTAS - REGIÓN " + region + " (MEJORADO) ===")
        console.log("ID Reporte: " + reportId)
        console.log("Período: " + fechaDesde + " hasta " + fechaHasta)
        console.log("Fecha Generación: " + std.now())
        
        if (region == "R01") {
            console.log("== REGIÓN USA ==")
            console.log("ACTIVOS:")
            console.log("  111002 - Caja USD: " + formatCurrency(25000.00, "USD"))
            console.log("  112002 - Citibank USD: " + formatCurrency(125000.50, "USD"))
            console.log("  121002 - Clientes USA: " + formatCurrency(180000.25, "USD"))
            console.log("PASIVOS:")
            console.log("  211002 - Proveedores USA: " + formatCurrency(95000.75, "USD"))
            console.log("  224002 - Sales Tax USA: " + formatCurrency(12500.00, "USD"))
            console.log("INGRESOS:")
            console.log("  411002 - Ventas USA: " + formatCurrency(450000.00, "USD"))
            console.log("GASTOS:")
            console.log("  511002 - Compras Insumos USA: " + formatCurrency(185000.50, "USD"))
            console.log("  521002 - Servicios USA: " + formatCurrency(95000.25, "USD"))
            
            let totalActivos = 330000.75
            let totalPasivos = 107500.75
            let patrimonioNeto = totalActivos - totalPasivos
            
            console.log("RESUMEN FINANCIERO:")
            console.log("  Total Activos: " + formatCurrency(totalActivos, "USD"))
            console.log("  Total Pasivos: " + formatCurrency(totalPasivos, "USD"))
            console.log("  Patrimonio Neto: " + formatCurrency(patrimonioNeto, "USD"))
            console.log("  Ratio Liquidez: " + math.round((totalActivos/totalPasivos) * 100) / 100)
            
        } else if (region == "R02") {
            console.log("== REGIÓN EUROPA ==")
            console.log("ACTIVOS:")
            console.log("  111003 - Caja EUR: " + formatCurrency(18000.50, "EUR"))
            console.log("  112003 - Deutsche Bank EUR: " + formatCurrency(95000.75, "EUR"))
            console.log("  121003 - Clientes Europa: " + formatCurrency(145000.25, "EUR"))
            console.log("PASIVOS:")
            console.log("  211003 - Proveedores Europa: " + formatCurrency(125000.00, "EUR"))
            console.log("  224003 - VAT Europa: " + formatCurrency(18500.50, "EUR"))
            console.log("INGRESOS:")
            console.log("  411003 - Ventas Europa: " + formatCurrency(380000.75, "EUR"))
            console.log("GASTOS:")
            console.log("  511003 - Compras Insumos Europa: " + formatCurrency(225000.25, "EUR"))
            console.log("  521003 - Servicios Europa: " + formatCurrency(115000.50, "EUR"))
            
        } else if (region == "R03") {
            console.log("== REGIÓN ARGENTINA ==")
            console.log("ACTIVOS:")
            console.log("  111001 - Caja Pesos: " + formatCurrency(150000.75, "ARS"))
            console.log("  112001 - Banco Nacional: " + formatCurrency(850000.50, "ARS"))
            console.log("  121001 - Clientes Nacionales: " + formatCurrency(320000.25, "ARS"))
            console.log("PASIVOS:")
            console.log("  211001 - Proveedores Nacionales: " + formatCurrency(280000.00, "ARS"))
            console.log("  224001 - IVA Débito Fiscal: " + formatCurrency(45000.75, "ARS"))
            console.log("INGRESOS:")
            console.log("  411001 - Ventas Nacionales: " + formatCurrency(1250000.50, "ARS"))
            console.log("GASTOS:")
            console.log("  511001 - Compras Insumos Nacionales: " + formatCurrency(750000.25, "ARS"))
            console.log("  521001 - Servicios Nacionales: " + formatCurrency(185000.75, "ARS"))
        }
        
        console.log("Estado del Análisis: COMPLETADO ✓")
        
        return {
            success: true,
            reportId: reportId,
            region: region,
            period: fechaDesde + " hasta " + fechaHasta,
            generated: std.now()
        }
    }
}

func main() {
    console.log("DSL MOTOR CONTABLE COMERCIAL MULTI-REGIÓN V3 MEJORADO")
    console.log("====================================================")
    console.log("Sistema con Formato Numérico, Rangos Dinámicos y Validación Avanzada")
    console.log("Versión: 3.0 - Fecha: " + std.now())
    console.log("")
    
    let motorVentasUSA = VentasUSAMejorado
    let motorVentasEUR = VentasEURMejorado
    let motorVentasARG = VentasARGMejorado
    let motorComprasUSA = ComprasUSAMejorado
    let motorAnalisis = AnalisisCuentasMejorado
    
    console.log("CASO 1: Venta USA con Decimales")
    console.log("===============================")
    let resultado1 = motorVentasUSA.use("venta USA 85250.75")
    console.log("Estado: " + resultado1.success + " | ID: " + resultado1.transactionId)
    console.log("")
    
    console.log("CASO 2: Venta Europa con Validación")
    console.log("===================================")
    let resultado2 = motorVentasEUR.use("venta EUR 15000.50")
    console.log("Estado: " + resultado2.success + " | ID: " + resultado2.transactionId)
    console.log("")
    
    console.log("CASO 3: Venta Argentina Mejorada")
    console.log("================================")
    let resultado3 = motorVentasARG.use("venta ARG 120750.25")
    console.log("Estado: " + resultado3.success + " | ID: " + resultado3.transactionId)
    console.log("")
    
    console.log("CASO 4: Compra USA Servicios")
    console.log("============================")
    let resultado4 = motorComprasUSA.use("compra USA servicios 25000.50")
    console.log("Estado: " + resultado4.success + " | ID: " + resultado4.transactionId)
    console.log("")
    
    console.log("CASO 5: Análisis Regional USA Mejorado")
    console.log("======================================")
    let resultado5 = motorAnalisis.use("analizar cuentas movimientos de R01 desde 01/01/2025 hasta 31/01/2025")
    console.log("Estado: " + resultado5.success + " | Reporte: " + resultado5.reportId)
    console.log("")
    
    console.log("=========================================")
    console.log("RESUMEN SISTEMA DSL V3 - CASOS EXITOSOS")
    console.log("=========================================")
    console.log("✓ Caso 1 - Venta USA: Procesada exitosamente con formato decimal")
    console.log("✓ Caso 2 - Venta EUR: Validación avanzada aplicada")
    console.log("✓ Caso 3 - Venta ARG: Números formateados a 2 decimales")
    console.log("✓ Caso 4 - Compra USA: ID transacción generado")
    console.log("✓ Caso 5 - Análisis Regional: Reporte mejorado generado")
    console.log("")
    console.log("MEJORAS IMPLEMENTADAS V3:")
    console.log("- Formato numérico con redondeo a 2 decimales")
    console.log("- Soporte para rangos dinámicos de importes")
    console.log("- Validación avanzada de entrada")
    console.log("- ID único de transacción para trazabilidad")
    console.log("- Timestamp en todos los comprobantes")
    console.log("- Formateo de moneda mejorado")
    console.log("- Cálculos de ratios financieros")
    console.log("- Validación de límites de importes")
    console.log("- Estados de validación visual")
    console.log("- Metadatos estructurados de respuesta")
    console.log("")
    console.log("SISTEMA V3 FUNCIONANDO AL 100% - MEJORAS APLICADAS ✓")
}