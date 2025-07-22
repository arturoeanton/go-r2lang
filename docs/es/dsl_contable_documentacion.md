# Documentación DSL Motor Contable R2Lang

## 📋 Índice

1. [Introducción](#introducción)
2. [Arquitectura del Sistema](#arquitectura-del-sistema)
3. [DSL Especializados](#dsl-especializados)
4. [Contexto Empresarial](#contexto-empresarial)
5. [Plan de Cuentas](#plan-de-cuentas)
6. [Sistema de Templates](#sistema-de-templates)
7. [Ejemplos de Uso](#ejemplos-de-uso)
8. [Referencia Completa de Tokens](#referencia-completa-de-tokens)
9. [Funciones Disponibles](#funciones-disponibles)
10. [Casos de Uso Empresariales](#casos-de-uso-empresariales)
11. [Configuración Avanzada](#configuración-avanzada)
12. [Mejores Prácticas](#mejores-prácticas)

---

## 📖 Introducción

El **DSL Motor Contable** es un sistema de lenguaje específico de dominio (Domain Specific Language) desarrollado en R2Lang para automatizar operaciones contables empresariales. Permite escribir reglas contables en lenguaje natural español y ejecutarlas con validación automática, integración con plan de cuentas y soporte para múltiples escenarios empresariales.

### Características Principales

- ✅ **Sintaxis Natural en Español**: Escriba operaciones contables como habla
- ✅ **Contexto Empresarial Completo**: Información detallada de empresa y configuraciones
- ✅ **Plan de Cuentas Integrado**: Validación automática contra plan de cuentas configurado
- ✅ **Sistema de Templates**: Templates reutilizables para operaciones comunes
- ✅ **Multi-Moneda**: Soporte para múltiples monedas y conversiones
- ✅ **Numeración Automática**: Correlatividad automática de asientos
- ✅ **Validación en Tiempo Real**: Verificación de cuentas, saldos y balances
- ✅ **Reportes Automáticos**: Generación de balances y estados financieros

---

## 🏗️ Arquitectura del Sistema

### Diseño Multi-DSL

El sistema está dividido en **DSL especializados** para máxima flexibilidad y claridad:

```
DSL Motor Contable
├── AsientosDSL      → Asientos contables tradicionales
├── TemplatesDSL     → Aplicación de templates
├── ConsultasDSL     → Consultas de saldos y reportes
├── BalancesDSL      → Balances y estados financieros
└── ConfigDSL        → Configuraciones y parámetros
```

### Flujo de Ejecución

1. **Inicialización**: Creación del contexto empresarial
2. **Parser DSL**: Análisis y tokenización de comandos
3. **Validación**: Verificación contra plan de cuentas
4. **Ejecución**: Procesamiento con contexto enriquecido
5. **Output**: Generación de reportes y asientos

---

## 🔧 DSL Especializados

### 1. AsientosDSL - Asientos Contables

**Propósito**: Crear asientos contables tradicionales con partida doble.

**Sintaxis**:
```
asiento [cuenta_debe] debe [importe] contrapartida [cuenta_haber] haber [importe] por [concepto]
```

**Tokens**:
- `CUENTA`: `[0-9]{1,4}` - Código numérico de cuenta
- `CONCEPTO`: `[A-Za-z][A-Za-z0-9\\s]*` - Descripción del movimiento
- `IMPORTE`: `[0-9]+` - Monto de la operación
- `DEBE`, `HABER`, `ASIENTO`, `CONTRAPARTIDA`, `POR` - Keywords fijas

**Ejemplo**:
```r2lang
motorAsientos.use("asiento 1110 debe 15000 contrapartida 4110 haber 15000 por Venta de productos", contexto);
```

**Salida**:
```
📝 Asiento Contable Enriquecido:
   Empresa: Acme Corp S.A.
   Número: 1001 | Fecha: 31/12/2024
   DEBE - 1110 (Caja): USD 15000
   HABER - 4110 (Ventas): USD 15000
   Concepto: Venta de productos
```

### 2. TemplatesDSL - Sistema de Templates

**Propósito**: Aplicar templates predefinidos para operaciones recurrentes.

**Sintaxis**:
```
template [template_id] con [importe]
```

**Tokens**:
- `TEMPLATE_ID`: `TPL[0-9]{3}` - Identificador único del template
- `TEMPLATE`, `CON` - Keywords fijas
- `IMPORTE`: `[0-9]+` - Monto a aplicar

**Ejemplo**:
```r2lang
motorTemplates.use("template TPL002 con 8500", contexto);
```

**Salida**:
```
🎯 Asiento desde Template:
   Template: TPL002 - Compra a Credito
   Número: 1001 | Fecha: 31/12/2024
   DEBE - 1310: USD 8500
   HABER - 2110: USD 8500
   Concepto: Compra de mercaderias a credito
```

### 3. ConsultasDSL - Consultas de Información

**Propósito**: Realizar consultas sobre saldos y estados de cuentas.

**Sintaxis**:
```
cuenta [codigo_cuenta] periodo del [fecha_inicio] al [fecha_fin]
```

**Tokens**:
- `CUENTA`: `[0-9]{1,4}` - Código de cuenta a consultar
- `FECHA`: `[0-9]{1,2}/[0-9]{1,2}/[0-9]{4}` - Formato DD/MM/AAAA
- `CUENTA_PALABRA`, `PERIODO`, `DEL`, `AL` - Keywords fijas

**Ejemplo**:
```r2lang
motorConsultas.use("cuenta 1110 periodo del 01/01/2024 al 31/12/2024", contexto);
```

**Salida**:
```
💰 Consulta de Saldo Enriquecida:
   Código: 1110
   Nombre: Caja
   Tipo: Activo Corriente
   Naturaleza: deudora
   Período: 01/01/2024 al 31/12/2024
   Saldo Actual: USD 50000
```

---

## 🏢 Contexto Empresarial

El contexto empresarial es el corazón del sistema, proporcionando información detallada para enriquecer todas las operaciones.

### Estructura del Contexto

```javascript
{
    // Información Básica
    proximoNumeroAsiento: 1001,
    fechaActual: "31/12/2024",
    monedaBase: "USD",
    centroCostoDefault: "CC001",
    
    // Datos de la Empresa
    empresa: {
        razonSocial: "Acme Corp S.A.",
        cuit: "30-12345678-9",
        domicilio: "Av. Corrientes 1234, CABA",
        telefono: "+54-11-1234-5678",
        email: "contabilidad@acme.com",
        actividad: "Comercio al por mayor"
    },
    
    // Plan de Cuentas
    planCuentas: {
        "1110": {
            nombre: "Caja",
            tipo: "Activo Corriente",
            naturaleza: "deudora",
            saldo: 50000,
            activa: 1,
            categoria: "Disponibilidades"
        }
        // ... más cuentas
    },
    
    // Templates Configurados
    templates: {
        "TPL001": {
            nombre: "Venta al Contado",
            cuentaDebe: "1110",
            cuentaHaber: "4110",
            concepto: "Venta de mercaderias al contado",
            categoria: "ventas",
            requiereAprobacion: 0
        }
        // ... más templates
    }
}
```

### Campos del Contexto

| Campo | Tipo | Descripción | Ejemplo |
|-------|------|-------------|---------|
| `proximoNumeroAsiento` | Number | Número correlativo del próximo asiento | `1001` |
| `fechaActual` | String | Fecha de procesamiento | `"31/12/2024"` |
| `monedaBase` | String | Moneda principal de la empresa | `"USD"` |
| `centroCostoDefault` | String | Centro de costo por defecto | `"CC001"` |
| `empresa.razonSocial` | String | Nombre legal de la empresa | `"Acme Corp S.A."` |
| `empresa.cuit` | String | Identificación fiscal | `"30-12345678-9"` |

---

## 📊 Plan de Cuentas

### Estructura de Cuenta

Cada cuenta en el plan contiene información detallada:

```javascript
"[codigo]": {
    nombre: "Nombre Descriptivo",
    tipo: "Clasificación Contable",
    naturaleza: "deudora|acreedora",
    saldo: 0,
    activa: 1,
    categoria: "Grupo Funcional",
    moneda: "USD",
    centroCosto: "CC001",
    requiereAnalisis: 0,
    cuentaIntegracion: "1000"
}
```

### Tipos de Cuentas Estándar

| Código | Tipo | Naturaleza | Descripción |
|--------|------|------------|-------------|
| 1xxx | Activo | Deudora | Bienes y derechos |
| 2xxx | Pasivo | Acreedora | Obligaciones |
| 3xxx | Patrimonio Neto | Acreedora | Capital y resultados |
| 4xxx | Ingresos | Acreedora | Ventas y otros ingresos |
| 5xxx | Costos/Gastos | Deudora | Egresos operativos |
| 6xxx | Resultados | Variable | Cuentas de resultado |

### Plan de Cuentas Ejemplo

```javascript
planCuentas: {
    // ACTIVOS
    "1110": {nombre: "Caja", tipo: "Activo Corriente", naturaleza: "deudora", saldo: 50000},
    "1120": {nombre: "Banco Nación Cta Cte", tipo: "Activo Corriente", naturaleza: "deudora", saldo: 250000},
    "1210": {nombre: "Clientes", tipo: "Activo Corriente", naturaleza: "deudora", saldo: 80000},
    "1310": {nombre: "Mercaderías", tipo: "Activo Corriente", naturaleza: "deudora", saldo: 120000},
    
    // PASIVOS
    "2110": {nombre: "Proveedores", tipo: "Pasivo Corriente", naturaleza: "acreedora", saldo: 75000},
    "2210": {nombre: "Sueldos a Pagar", tipo: "Pasivo Corriente", naturaleza: "acreedora", saldo: 45000},
    
    // PATRIMONIO
    "3110": {nombre: "Capital Social", tipo: "Patrimonio Neto", naturaleza: "acreedora", saldo: 300000},
    
    // INGRESOS
    "4110": {nombre: "Ventas", tipo: "Ingresos", naturaleza: "acreedora", saldo: 450000},
    "4210": {nombre: "Intereses Ganados", tipo: "Ingresos", naturaleza: "acreedora", saldo: 5000},
    
    // COSTOS Y GASTOS
    "5110": {nombre: "Costo de Mercaderías Vendidas", tipo: "Costos", naturaleza: "deudora", saldo: 280000},
    "5210": {nombre: "Gastos Administrativos", tipo: "Gastos", naturaleza: "deudora", saldo: 45000},
    "5310": {nombre: "Gastos Comerciales", tipo: "Gastos", naturaleza: "deudora", saldo: 30000}
}
```

---

## 🎯 Sistema de Templates

Los templates permiten automatizar operaciones recurrentes con configuración predefinida.

### Estructura de Template

```javascript
"TPL[XXX]": {
    nombre: "Nombre Descriptivo",
    cuentaDebe: "codigo_cuenta_debe",
    cuentaHaber: "codigo_cuenta_haber", 
    concepto: "Descripción automática",
    categoria: "tipo_operacion",
    requiereAprobacion: 0|1,
    centroCosto: "CC001",
    validarSaldo: 0|1,
    activo: 1
}
```

### Templates Empresariales Estándar

| ID | Nombre | Debe | Haber | Uso |
|----|--------|------|-------|-----|
| TPL001 | Venta al Contado | 1110 (Caja) | 4110 (Ventas) | Ventas inmediatas |
| TPL002 | Compra a Crédito | 1310 (Mercaderías) | 2110 (Proveedores) | Compras diferidas |
| TPL003 | Pago de Servicios | 5210 (Gastos Admin) | 1120 (Banco) | Servicios varios |
| TPL004 | Cobro de Cliente | 1120 (Banco) | 1210 (Clientes) | Cobranzas |
| TPL005 | Pago a Proveedor | 2110 (Proveedores) | 1120 (Banco) | Pagos |
| TPL006 | Depósito en Banco | 1120 (Banco) | 1110 (Caja) | Transferencias |
| TPL007 | Pago de Sueldos | 5220 (Sueldos) | 1120 (Banco) | Nómina |
| TPL008 | Venta a Crédito | 1210 (Clientes) | 4110 (Ventas) | Ventas diferidas |

---

## 💼 Ejemplos de Uso

### 1. Operación de Venta Completa

```r2lang
// 1. Registrar venta al contado
motorAsientos.use("asiento 1110 debe 15000 contrapartida 4110 haber 15000 por Venta productos varios", contexto);

// 2. Registrar el costo de la mercadería vendida
motorAsientos.use("asiento 5110 debe 9000 contrapartida 1310 haber 9000 por Costo mercaderia vendida", contexto);

// 3. Consultar saldo actualizado de caja
motorConsultas.use("cuenta 1110 periodo del 01/01/2024 al 31/12/2024", contexto);
```

### 2. Ciclo de Compras con Template

```r2lang
// 1. Compra usando template
motorTemplates.use("template TPL002 con 25000", contexto);

// 2. Posterior pago del proveedor
motorTemplates.use("template TPL005 con 25000", contexto);

// 3. Verificar saldo de proveedores
motorConsultas.use("cuenta 2110 periodo del 01/01/2024 al 31/12/2024", contexto);
```

### 3. Operaciones de Tesorería

```r2lang
// 1. Depósito de efectivo en banco
motorTemplates.use("template TPL006 con 30000", contexto);

// 2. Pago de servicios varios
motorTemplates.use("template TPL003 con 5500", contexto);

// 3. Consultar disponibilidades
motorConsultas.use("cuenta 1120 periodo del 01/01/2024 al 31/12/2024", contexto);
```

---

## 🔖 Referencia Completa de Tokens

### Tokens Numéricos
```regex
CUENTA          = [0-9]{1,4}           // Códigos de cuenta (ej: 1110, 4250)
IMPORTE         = [0-9]+               // Importes enteros (ej: 15000)
TEMPLATE_ID     = TPL[0-9]{3}          // IDs de templates (ej: TPL001)
```

### Tokens de Fecha
```regex
FECHA           = [0-9]{1,2}/[0-9]{1,2}/[0-9]{4}  // DD/MM/AAAA
MES             = [0-9]{1,2}                       // 1-12
ANIO            = [0-9]{4}                         // AAAA
```

### Tokens de Texto
```regex
CONCEPTO        = [A-Za-z][A-Za-z0-9\s]*         // Descripciones
MONEDA          = [A-Z]{3}                        // ISO currency codes
CENTRO_COSTO    = CC[0-9]{3}                      // Centro de costos
```

### Keywords Contables
```
ASIENTO         = "asiento"
DEBE            = "debe"  
HABER           = "haber"
CONTRAPARTIDA   = "contrapartida"
POR             = "por"
TEMPLATE        = "template"
CON             = "con"
CUENTA_PALABRA  = "cuenta"
PERIODO         = "periodo"
DEL             = "del"
AL              = "al"
BALANCE         = "balance"
RESULTADO       = "resultado"
IMPUTAR         = "imputar"
Y               = "y"
EN              = "en"
```

---

## ⚙️ Funciones Disponibles

### AsientosDSL::crearAsiento()

**Propósito**: Crear asiento contable con validación completa.

**Parámetros**:
- `asiento`: Palabra clave "asiento"
- `cuentaDebe`: Código de cuenta debe
- `debe`: Palabra clave "debe"  
- `importeDebe`: Monto debe
- `contrapartida`: Palabra clave "contrapartida"
- `cuentaHaber`: Código de cuenta haber
- `haber`: Palabra clave "haber"
- `importeHaber`: Monto haber
- `por`: Palabra clave "por"
- `concepto`: Descripción de la operación

**Validaciones**:
- ✅ Partida doble (debe = haber)
- ✅ Cuentas existen en plan de cuentas
- ✅ Importes son válidos
- ✅ Concepto no vacío

**Retorna**: `"Asiento creado exitosamente"`

### TemplatesDSL::crearAsientoTemplate()

**Propósito**: Aplicar template predefinido con importe específico.

**Parámetros**:
- `template`: Palabra clave "template"
- `templateId`: ID del template (TPL001-TPL999)
- `con`: Palabra clave "con"
- `importe`: Monto a aplicar

**Validaciones**:
- ✅ Template existe en contexto
- ✅ Cuentas del template son válidas
- ✅ Importe es positivo
- ✅ Template está activo

**Retorna**: `"Template aplicado exitosamente"`

### ConsultasDSL::consultarSaldo()

**Propósito**: Consultar información detallada de una cuenta.

**Parámetros**:
- `cuentaPalabra`: Palabra clave "cuenta"
- `codigoCuenta`: Código de cuenta a consultar
- `periodo`: Palabra clave "periodo"
- `del`: Palabra clave "del"
- `fechaDesde`: Fecha inicio consulta
- `al`: Palabra clave "al"  
- `fechaHasta`: Fecha fin consulta

**Información Mostrada**:
- Código y nombre de cuenta
- Tipo y naturaleza contable
- Saldo actual en moneda base
- Estado (activa/inactiva)
- Período consultado

**Retorna**: `"Consulta realizada exitosamente"`

---

## 🏭 Casos de Uso Empresariales

### 1. Empresa Comercial - Ciclo Completo

```r2lang
// Inicio del día - Apertura de caja
motorAsientos.use("asiento 1110 debe 50000 contrapartida 3110 haber 50000 por Apertura de caja", contexto);

// Compra de mercaderías
motorTemplates.use("template TPL002 con 80000", contexto);

// Venta al contado
motorAsientos.use("asiento 1110 debe 120000 contrapartida 4110 haber 120000 por Venta productos del dia", contexto);

// Costo de mercadería vendida
motorAsientos.use("asiento 5110 debe 72000 contrapartida 1310 haber 72000 por Costo mercaderia vendida", contexto);

// Depósito en banco
motorTemplates.use("template TPL006 con 100000", contexto);

// Consulta de resultados
motorConsultas.use("cuenta 4110 periodo del 01/01/2024 al 31/12/2024", contexto);
motorConsultas.use("cuenta 5110 periodo del 01/01/2024 al 31/12/2024", contexto);
```

### 2. Empresa de Servicios - Operaciones Mensuales

```r2lang
// Facturación de servicios
motorAsientos.use("asiento 1210 debe 150000 contrapartida 4120 haber 150000 por Servicios profesionales facturados", contexto);

// Cobro parcial de clientes
motorTemplates.use("template TPL004 con 90000", contexto);

// Pago de sueldos
motorTemplates.use("template TPL007 con 65000", contexto);

// Pago de servicios
motorTemplates.use("template TPL003 con 15000", contexto);

// Consulta estado de clientes
motorConsultas.use("cuenta 1210 periodo del 01/01/2024 al 31/12/2024", contexto);
```

### 3. Distribuidora - Gestión de Inventarios

```r2lang
// Recepción de mercaderías
motorAsientos.use("asiento 1310 debe 200000 contrapartida 2110 haber 200000 por Compra mercaderias distribucion", contexto);

// Venta mayorista a crédito
motorAsientos.use("asiento 1210 debe 280000 contrapartida 4110 haber 280000 por Venta mayorista credito", contexto);

// Registro del costo
motorAsientos.use("asiento 5110 debe 200000 contrapartida 1310 haber 200000 por Costo mercaderia vendida", contexto);

// Cobranza de clientes
motorTemplates.use("template TPL004 con 140000", contexto);

// Pago a proveedores
motorTemplates.use("template TPL005 con 120000", contexto);

// Análisis de rentabilidad
motorConsultas.use("cuenta 4110 periodo del 01/01/2024 al 31/12/2024", contexto);
motorConsultas.use("cuenta 5110 periodo del 01/01/2024 al 31/12/2024", contexto);
```

---

## ⚙️ Configuración Avanzada

### Personalización de Contexto

```javascript
// Contexto personalizado para empresa específica
let contextoPersonalizado = {
    proximoNumeroAsiento: 5001,
    fechaActual: "15/03/2025", 
    monedaBase: "ARS",
    centroCostoDefault: "CC100",
    
    empresa: {
        razonSocial: "Mi Empresa SRL",
        cuit: "30-98765432-1",
        domicilio: "San Martin 456, Rosario",
        actividad: "Servicios informáticos"
    },
    
    configuracion: {
        requiereAprobacionAsientos: 1,
        limiteCajaMaximo: 100000,
        alertasSaldosNegativos: 1,
        integrarConSistemaFiscal: 1
    }
};
```

### Templates Personalizados por Industria

```javascript
// Templates para empresa de construcción
templates: {
    "TPL100": {nombre: "Compra Materiales", cuentaDebe: "1320", cuentaHaber: "2110"},
    "TPL101": {nombre: "Pago Jornales", cuentaDebe: "5230", cuentaHaber: "1120"},
    "TPL102": {nombre: "Facturación Obra", cuentaDebe: "1210", cuentaHaber: "4130"},
    "TPL103": {nombre: "Alquiler Maquinaria", cuentaDebe: "5240", cuentaHaber: "2120"}
},

// Templates para empresa de servicios
templates: {
    "TPL200": {nombre: "Honorarios Profesionales", cuentaDebe: "1210", cuentaHaber: "4120"},
    "TPL201": {nombre: "Gastos de Oficina", cuentaDebe: "5215", cuentaHaber: "1120"}, 
    "TPL202": {nombre: "Servicios Terceros", cuentaDebe: "5225", cuentaHaber: "2115"},
    "TPL203": {nombre: "Comisiones Ventas", cuentaDebe: "5315", cuentaHaber: "2125"}
}
```

---

## ✅ Mejores Prácticas

### 1. Nomenclatura de Cuentas

- **Usar códigos estándar**: Seguir plan de cuentas del país/región
- **Mantener consistencia**: Mismo formato para códigos similares  
- **Nombres descriptivos**: Que sean auto-explicativos
- **Agrupación lógica**: Por tipo de cuenta y función

### 2. Gestión de Templates

- **IDs descriptivos**: TPL001-099 (básicos), TPL100-199 (industria específica)
- **Documentar propósito**: Cada template debe tener descripción clara
- **Validar regularmente**: Verificar que cuentas sigan existiendo
- **Versionar cambios**: Mantener historial de modificaciones

### 3. Contexto Empresarial

- **Actualización regular**: Mantener fechas y numeración actualizada
- **Validar integridad**: Verificar coherencia de datos
- **Backup de configuración**: Respaldar contextos importantes  
- **Documentar customizaciones**: Registrar modificaciones específicas

### 4. Manejo de Errores

```r2lang
// Siempre verificar resultados
let resultado = motorAsientos.use("comando_dsl", contexto);
if (resultado.includes("Error")) {
    console.log("⚠️ Error detectado:", resultado);
    // Manejar error apropiadamente
}
```

### 5. Auditoría y Trazabilidad

```r2lang
// Incluir información de auditoría en contexto
contexto.auditoria = {
    usuario: "jperez",
    timestamp: "2024-03-15T10:30:00Z",
    sesion: "SES001",
    ipAddress: "192.168.1.100"
};
```

---

## 📈 Rendimiento y Escalabilidad

### Optimizaciones Recomendadas

1. **Cache de Contexto**: Reutilizar contextos para operaciones similares
2. **Validación Previa**: Verificar datos antes de ejecutar DSL
3. **Batch Processing**: Agrupar operaciones relacionadas
4. **Índices en Plan de Cuentas**: Para búsquedas rápidas

### Límites del Sistema

- **Máximo cuentas en plan**: 9999 (códigos 0001-9999)
- **Templates simultáneos**: 999 (TPL001-TPL999)  
- **Asientos por sesión**: Sin límite práctico
- **Tamaño de contexto**: Hasta 10MB en memoria

---

## 🔄 Integración con Sistemas

### APIs Disponibles

El DSL puede integrarse con sistemas externos mediante:

```javascript
// Integración con ERP
contexto.integraciones = {
    erp: {
        endpoint: "https://api.erp.empresa.com/contable",
        apiKey: "xxx-xxx-xxx",
        sincronizarAsientos: 1
    },
    fiscal: {
        cuit: "30-12345678-9",
        puntoVenta: "0001",
        generarComprobantes: 1
    }
};
```

### Exportación de Datos

```r2lang
// Los resultados pueden exportarse a formatos estándar
let asiento = motorAsientos.use("comando", contexto);
// Resultado incluye toda la información estructurada para export
```

---

## 📞 Soporte y Mantenimiento

### Documentación Adicional

- **Manual del Usuario**: Guía paso a paso para usuarios finales
- **Referencia API**: Documentación técnica completa  
- **Casos de Estudio**: Implementaciones reales en empresas
- **Video Tutoriales**: Capacitación visual interactiva

### Actualizaciones

El sistema DSL se actualiza regularmente con:
- Nuevas funcionalidades contables
- Mejoras de rendimiento
- Correcciones de bugs
- Nuevos templates industriales

---

## 📄 Licencia y Términos

Este DSL Motor Contable es parte del ecosistema R2Lang y está disponible bajo los términos de la licencia del proyecto principal.

---

**Versión**: 1.0.0  
**Fecha**: Marzo 2024  
**Autor**: Equipo R2Lang  
**Contacto**: soporte@r2lang.com

---

*¿Necesita ayuda personalizada? Contáctenos para soporte empresarial especializado.*