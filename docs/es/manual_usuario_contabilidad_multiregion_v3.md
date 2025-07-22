# Manual de Usuario - Sistema Contable Comercial Multi-Región V3

## 🎯 Introducción

El **Sistema Contable Comercial Multi-Región V3** es una solución empresarial avanzada construida con R2Lang DSL que automatiza el procesamiento de comprobantes contables para múltiples regiones geográficas con diferentes normativas fiscales.

### ✨ Características Principales V3

- **Formato Numérico Avanzado**: Redondeo automático a 2 decimales
- **Rangos Dinámicos**: Soporte para importes variables con validación
- **Validación Robusta**: Controles de entrada y límites de seguridad
- **Trazabilidad Completa**: ID único por transacción
- **Multi-Regional**: USA, Europa y Argentina
- **Multi-Moneda**: USD, EUR, ARS
- **Cumplimiento Normativo**: US-GAAP, IFRS, RT Argentina

## 📋 Instalación y Configuración

### Requisitos del Sistema

- **R2Lang Runtime**: Versión 2.0 o superior
- **Memoria**: Mínimo 512MB RAM
- **Almacenamiento**: 50MB espacio libre
- **SO**: Linux, Windows, macOS

### Instalación

```bash
# Clonar el repositorio
git clone https://github.com/tu-empresa/go-r2lang.git

# Navegar al directorio
cd go-r2lang

# Compilar el sistema
go build -o r2lang main.go

# Verificar instalación
./r2lang --version
```

### Configuración Inicial

```bash
# Ejecutar el sistema contable
./r2lang examples/dsl/contabilidad_comercial_multiregion_v3_mejorado.r2
```

## 🚀 Guía de Uso

### 1. Procesamiento de Ventas

#### Venta en Estados Unidos
```r2
venta USA 85250.75
```

**Resultado Esperado:**
```
=== COMPROBANTE DE VENTA USA (MEJORADO) ===
ID Transacción: TX-1642873456-456
Cliente: TechSoft USA Inc.
Region: R01 - America del Norte
Fecha: 2025-01-22 14:30:15
DEBE: 121002 USD 92,647.19
HABER: 411002 USD 85,250.75 + 224002 USD 7,396.44
Tasa Impuesto: 8.75%
Estado: VALIDADO ✓
```

#### Venta en Europa
```r2
venta EUR 15000.50
```

**Resultado Esperado:**
```
=== COMPROBANTE DE VENTA EUROPA (MEJORADO) ===
ID Transacción: TX-1642873457-789
Cliente: EuroSystems GmbH
Region: R02 - Europa
DEBE: 121003 EUR 18,000.60
HABER: 411003 EUR 15,000.50 + 224003 EUR 3,000.10
Tasa Impuesto: 20.00%
Estado: VALIDADO ✓
```

#### Venta en Argentina
```r2
venta ARG 120750.25
```

**Resultado Esperado:**
```
=== COMPROBANTE DE VENTA ARGENTINA (MEJORADO) ===
ID Transacción: TX-1642873458-123
Cliente: Sistemas Locales S.A.
Region: R03 - America del Sur
DEBE: 121001 ARS 146,107.80
HABER: 411001 ARS 120,750.25 + 224001 ARS 25,357.55
Tasa Impuesto: 21.00%
Estado: VALIDADO ✓
```

### 2. Procesamiento de Compras

#### Compra de Servicios USA
```r2
compra USA servicios 25000.50
```

**Resultado Esperado:**
```
=== COMPROBANTE DE COMPRA USA SERVICIOS (MEJORADO) ===
ID Transacción: TX-1642873459-456
Proveedor: Amazon Web Services
Region: R01 - America del Norte
DEBE: 521002 USD 25,000.50 + 113002 USD 2,187.54
HABER: 211002 USD 27,188.04
Estado: VALIDADO ✓
```

### 3. Análisis de Cuentas

#### Análisis Regional
```r2
analizar cuentas movimientos de R01 desde 01/01/2025 hasta 31/01/2025
```

**Resultado Esperado:**
```
=== ANÁLISIS DE CUENTAS - REGIÓN R01 (MEJORADO) ===
ID Reporte: TX-1642873460-789
Período: 01/01/2025 hasta 31/01/2025
Fecha Generación: 2025-01-22 14:35:20

== REGIÓN USA ==
ACTIVOS:
  111002 - Caja USD: USD 25,000.00
  112002 - Citibank USD: USD 125,000.50
  121002 - Clientes USA: USD 180,000.25

RESUMEN FINANCIERO:
  Total Activos: USD 330,000.75
  Total Pasivos: USD 107,500.75
  Patrimonio Neto: USD 222,500.00
  Ratio Liquidez: 3.07
Estado del Análisis: COMPLETADO ✓
```

## ⚙️ Configuraciones Avanzadas

### Límites de Validación

El sistema incluye validaciones automáticas:

- **Importe Mínimo**: $0.01 (cualquier moneda)
- **Importe Máximo**: $10,000,000.00
- **Decimales**: Máximo 2 posiciones
- **Formato**: Números positivos únicamente

### Personalización de Tasas

Las tasas de impuestos están configuradas por región:

```r2
// En el código DSL
let tasaIVA_USA = 0.0875    // 8.75%
let tasaIVA_EUR = 0.20      // 20.00%
let tasaIVA_ARG = 0.21      // 21.00%
```

### ID de Transacciones

Formato automático: `TX-[timestamp]-[random]`
- Ejemplo: `TX-1642873456-456`
- Garantiza unicidad y trazabilidad

## 🔧 Solución de Problemas

### Errores Comunes

#### Error: "El importe no puede ser negativo"
**Causa**: Se ingresó un valor negativo
**Solución**: Usar solo valores positivos
```r2
// ❌ Incorrecto
venta USA -1000

// ✅ Correcto  
venta USA 1000
```

#### Error: "El importe excede el límite máximo"
**Causa**: Valor superior a 10,000,000
**Solución**: Fraccionar en múltiples transacciones
```r2
// ❌ Incorrecto
venta USA 15000000

// ✅ Correcto
venta USA 9999999
```

#### Error de Formato de Fecha
**Causa**: Formato de fecha incorrecto en análisis
**Solución**: Usar formato DD/MM/YYYY
```r2
// ❌ Incorrecto
analizar cuentas movimientos de R01 desde 2025-01-01 hasta 2025-01-31

// ✅ Correcto
analizar cuentas movimientos de R01 desde 01/01/2025 hasta 31/01/2025
```

### Logs y Depuración

Para habilitar logs detallados:
```bash
./r2lang --debug examples/dsl/contabilidad_comercial_multiregion_v3_mejorado.r2
```

## 📊 Casos de Uso Empresariales

### Caso 1: Empresa Multinacional
**Escenario**: Procesar ventas simultáneas en 3 regiones
```r2
venta USA 50000.00
venta EUR 35000.50  
venta ARG 85000.75
```

### Caso 2: Análisis Consolidado Mensual
**Escenario**: Generar reportes de todas las regiones
```r2
analizar cuentas movimientos de R01 desde 01/01/2025 hasta 31/01/2025
analizar cuentas movimientos de R02 desde 01/01/2025 hasta 31/01/2025
analizar cuentas movimientos de R03 desde 01/01/2025 hasta 31/01/2025
```

### Caso 3: Procesamiento de Compras de Servicios
**Escenario**: Registro de gastos en servicios cloud
```r2
compra USA servicios 15000.00
compra USA servicios 25000.50
compra USA servicios 8750.25
```

## 🔐 Seguridad y Cumplimiento

### Normativas Soportadas

- **US-GAAP** (Estados Unidos): Principios contables estadounidenses
- **IFRS** (Europa): Estándares internacionales de información financiera  
- **RT Argentina**: Resoluciones técnicas argentinas

### Auditoría y Trazabilidad

Cada transacción genera:
- ID único de transacción
- Timestamp de procesamiento
- Región y normativa aplicada
- Estados de validación
- Metadatos estructurados

### Respaldo de Datos

Los resultados se pueden exportar en formato estructurado:
```json
{
  "success": true,
  "transactionId": "TX-1642873456-456",
  "amount": 92647.19,
  "currency": "USD",
  "region": "USA"
}
```

## 📈 Métricas y Reportes

### Indicadores Financieros Automáticos

El sistema calcula automáticamente:
- **Ratio de Liquidez**: Activos / Pasivos
- **Patrimonio Neto**: Activos - Pasivos
- **Totales por Categoría**: Sumas automáticas
- **Porcentajes de Impuestos**: Cálculos precisos

### Exportación de Datos

Los resultados pueden integrarse con:
- Sistemas ERP existentes
- Herramientas de BI
- Plataformas de auditoría
- Software contable tradicional

## 🎓 Capacitación y Soporte

### Recursos de Aprendizaje

1. **Documentación Técnica**: `/docs/es/`
2. **Ejemplos Prácticos**: `/examples/dsl/`
3. **Casos de Prueba**: Ver archivo de documentación V2
4. **Videos Tutoriales**: (Disponibles en portal corporativo)

### Soporte Técnico

- **Email**: soporte-contabilidad@tu-empresa.com
- **Teléfono**: +1-800-CONTABLE
- **Portal Web**: https://support.tu-empresa.com
- **Horarios**: Lunes a Viernes 8:00-18:00

### Actualizaciones

El sistema se actualiza automáticamente. Para forzar actualización:
```bash
./r2lang --update
./r2lang --version
```

## ⚡ Mejores Prácticas

### Recomendaciones de Uso

1. **Validar Importes**: Siempre verificar rangos antes del procesamiento
2. **Usar Decimales**: Especificar centavos para precisión contable
3. **Generar Reportes Regulares**: Análisis mensual mínimo
4. **Mantener Trazabilidad**: Conservar IDs de transacción
5. **Respaldar Datos**: Exportar resultados regularmente

### Optimización de Performance

- Procesar hasta 1000 transacciones por lote
- Usar análisis regional por separado para grandes volúmenes
- Programar procesamiento en horarios de baja carga

---

## 📝 Changelog V3

### Nuevas Características
- ✅ Formato numérico con redondeo a 2 decimales
- ✅ Soporte para rangos dinámicos de importes  
- ✅ Validación avanzada de entrada
- ✅ ID único de transacción
- ✅ Timestamp en comprobantes
- ✅ Formateo mejorado de moneda
- ✅ Cálculos de ratios financieros
- ✅ Estados de validación visual
- ✅ Metadatos estructurados

### Mejoras de Funcionalidad
- Eliminado límite fijo de importes predefinidos
- Agregada validación de límites máximos y mínimos
- Implementado sistema de trazabilidad completo
- Mejorada precisión en cálculos de impuestos

---

**© 2025 Tu Empresa - Sistema Contable Multi-Región V3**  
**Versión del Manual**: 3.0.1 - Fecha: 22/01/2025