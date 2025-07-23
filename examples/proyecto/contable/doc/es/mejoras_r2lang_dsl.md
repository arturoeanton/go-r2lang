# 🔧 Mejoras Propuestas para R2Lang y DSL Builder

## 🎯 Resumen Ejecutivo

Durante el desarrollo del sistema contable, se identificaron varias áreas donde R2Lang y su DSL Builder podrían mejorar significativamente. Estas mejoras no solo facilitarían el desarrollo, sino que posicionarían a R2Lang como una plataforma ideal para aplicaciones empresariales.

## 1. 🔄 Arrays y Estructuras de Datos ✅ COMPLETADO (2025-07-22)

### Problema Actual
```r2
// Esto NO funciona actualmente
let asiento = {
    movimientos: []
}
asiento.movimientos.push(item)  // ❌ No funciona
asiento.movimientos[0] = item   // ❌ No funciona
```

### Mejora Implementada ✅
```r2
// Arrays anidados ahora funcionan con patrón de reasignación
let asiento = {
    movimientos: []
}
// Workaround funcional:
asiento.movimientos = asiento.movimientos.push(item)  // ✅ Funciona
asiento.movimientos[0] = item   // ✅ Funciona

// Múltiples niveles también soportados
empresa.sucursales[0].empleados = empresa.sucursales[0].empleados.push(nuevoEmpleado)
```

### Implementación Realizada
- ✅ Modificado `GenericAssignStatement` para soportar asignación a propiedades de maps
- ✅ Actualizado `std.len()` para manejar tipo `InterfaceSlice`
- ✅ Arrays anidados funcionan con patrón de reasignación (push retorna nuevo array)
- ✅ Tests completos agregados en `tests/test_nested_arrays.r2`

## 2. 📝 Template Literals y Strings Multilínea ✅ COMPLETADO (2025-07-22)

### Problema Actual
```r2
// Construcción tediosa de HTML
let html = "<!DOCTYPE html>\n"
html = html + "<html>\n"
html = html + "<head>\n"
// ... cientos de líneas
```

### Mejora Implementada ✅
```r2
// Template literals con interpolación (YA FUNCIONAN!)
let html = `
<!DOCTYPE html>
<html>
<head>
    <title>${titulo}</title>
    <style>
        body { font-family: ${fuente}; }
    </style>
</head>
<body>
    <h1>Bienvenido ${usuario.nombre}</h1>
    <p>Saldo: ${formatearMoneda(saldo)}</p>
</body>
</html>
`

// Strings multilínea con template literals
let sql = `
    SELECT t.*, a.descripcion
    FROM transacciones t
    JOIN asientos a ON t.id = a.transaccion_id
    WHERE t.fecha BETWEEN ? AND ?
    ORDER BY t.fecha DESC
`
```

### Características Implementadas
- ✅ Interpolación con `${expresión}` totalmente funcional
- ✅ Preservación de indentación
- ✅ Soporte para expresiones complejas
- ✅ Strings multilínea con backticks
- ✅ Tests completos agregados en `tests/test_template_strings.r2`

## 3. 🏗️ DSL Builder Mejorado

### Estado Actual
```r2
dsl MiDSL {
    token("NUMERO", "[0-9]+")
    rule("suma", ["NUMERO", "+", "NUMERO"], "sumar")
    func sumar(a, op, b) {
        return parseFloat(a) + parseFloat(b)
    }
}
```

### Mejoras Propuestas

#### 3.1 Gramáticas más Expresivas
```r2
dsl ContabilidadDSL {
    // Tokens con nombres más descriptivos
    tokens {
        CUENTA: /[0-9]{4}(\.[0-9]{2})*/
        MONTO: /\$?[0-9]+(\.[0-9]{2})?/
        FECHA: /\d{4}-\d{2}-\d{2}/
        DEBE: "debe" | "DEBE" | "D"
        HABER: "haber" | "HABER" | "H"
    }
    
    // Reglas con sintaxis tipo BNF
    rules {
        asiento: fecha descripcion movimientos+
        movimiento: cuenta (DEBE | HABER) MONTO
        consulta: "balance" fecha? | "diario" periodo?
    }
    
    // Acciones semánticas inline
    asiento(fecha, desc, movs) => {
        validarBalance(movs)
        return crearAsiento(fecha, desc, movs)
    }
}
```

#### 3.2 Composición de DSLs
```r2
// DSL base reutilizable
dsl BaseFiscal {
    tokens {
        RFC: /[A-Z]{3,4}[0-9]{6}[A-Z0-9]{3}/
        MONTO: /\$?[0-9]+(\.[0-9]{2})?/
    }
    
    validarRFC(rfc) => {
        // Validación del RFC
    }
}

// DSL que extiende el base
dsl FacturacionCFDI extends BaseFiscal {
    tokens {
        UUID: /[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/
        USO_CFDI: "G01" | "G02" | "G03" // etc
    }
    
    rules {
        factura: emisor receptor conceptos impuestos
    }
}
```

#### 3.3 Validación y Mensajes de Error
```r2
dsl ValidacionDSL {
    // Validaciones personalizadas
    constraints {
        monto: must be > 0 "El monto debe ser positivo"
        fecha: must match /\d{4}-\d{2}-\d{2}/ "Formato: YYYY-MM-DD"
        cuenta: must exist in catalogoCuentas "Cuenta no existe"
    }
    
    // Mensajes de error contextuales
    onError {
        typeMismatch: "Se esperaba ${expected} pero se encontró ${found}"
        unexpectedToken: "Token inesperado '${token}' en línea ${line}"
        customError(msg): "Error: ${msg}"
    }
}
```

## 4. 🔗 Sistema de Tipos Opcional

### Propuesta
```r2
// Tipos opcionales para mejor tooling
type Transaccion = {
    id: string,
    fecha: Date,
    monto: number,
    tipo: "venta" | "compra",
    region: Region
}

type Region = "MX" | "COL" | "AR" | "CH" | "UY" | "EC" | "PE"

func procesarTransaccion(tx: Transaccion): Asiento {
    // El IDE puede ofrecer autocompletado
    // El compilador puede detectar errores
}

// Inferencia de tipos
let total = 0  // Inferido como number
total = "texto"  // ⚠️ Warning: tipo incompatible
```

### Beneficios
- Autocompletado en IDEs
- Detección temprana de errores
- Documentación implícita
- Refactoring más seguro
- Opcional: no rompe código existente

## 5. 🎭 Macros y Metaprogramación

### Propuesta
```r2
// Macros para código repetitivo
macro defineHandler(path, method) {
    http.handler(method, path, handle${path}${method})
    
    func handle${path}${method}(pathVars, method, body) {
        console.log("[${method}] ${path}")
        // Código generado
    }
}

// Uso
defineHandler("/usuarios", "GET")
defineHandler("/usuarios", "POST")
defineHandler("/usuarios/:id", "GET")

// Genera automáticamente:
// - handleUsuariosGET()
// - handleUsuariosPOST()
// - handleUsuariosIdGET()
```

## 6. 🧩 Módulos y Namespaces

### Propuesta
```r2
// Definir módulo
module Contabilidad {
    export type Cuenta = {
        codigo: string,
        nombre: string,
        tipo: "activo" | "pasivo" | "capital"
    }
    
    export func crearAsiento(fecha, descripcion) {
        // Implementación
    }
    
    // Privado del módulo
    let configuracionInterna = {}
}

// Usar módulo
import { Cuenta, crearAsiento } from Contabilidad

// O importar todo
import * as Cont from Contabilidad
let cuenta: Cont.Cuenta = { ... }
```

## 7. 🔄 Async/Await Nativo

### Propuesta
```r2
// Soporte nativo para operaciones asíncronas
async func obtenerDatosExternos(url) {
    try {
        let response = await http.get(url)
        let data = await response.json()
        return data
    } catch (error) {
        console.error("Error: " + error)
        return null
    }
}

// Uso
let datos = await obtenerDatosExternos("https://api.example.com/data")
```

## 8. 🎨 Decoradores y Anotaciones

### Propuesta
```r2
// Decoradores para cross-cutting concerns
@authenticate
@rateLimit(100)
@cache(ttl: 300)
func handleAPI(pathVars, method, body) {
    // La función está automáticamente:
    // - Protegida por autenticación
    // - Limitada a 100 requests por minuto
    // - Cacheada por 5 minutos
}

// Decorador personalizado
decorator measureTime(func) {
    return func(...args) {
        let start = std.now()
        let result = func(...args)
        let duration = std.now() - start
        console.log(`${func.name} tomó ${duration}ms`)
        return result
    }
}

@measureTime
func procesoLento() {
    // Se medirá automáticamente
}
```

## 9. 🔍 Pattern Matching

### Propuesta
```r2
// Pattern matching al estilo funcional
let resultado = match transaccion.tipo {
    "venta" => procesarVenta(transaccion),
    "compra" => procesarCompra(transaccion),
    "devolucion" if transaccion.monto < 0 => procesarDevolucion(transaccion),
    _ => lanzarError("Tipo desconocido")
}

// Destructuring en patterns
match response {
    { status: 200, data: { id, nombre } } => {
        console.log(`Usuario ${nombre} creado con ID ${id}`)
    },
    { status: 404 } => console.error("No encontrado"),
    { status, error } => console.error(`Error ${status}: ${error}`)
}
```

## 10. 🛠️ Herramientas de Desarrollo

### R2Lang LSP (Language Server Protocol)
```json
{
    "r2lang.lsp": {
        "features": [
            "autocompletado",
            "go-to-definition",
            "find-references",
            "rename-symbol",
            "format-document",
            "code-actions"
        ]
    }
}
```

### R2Lang Formatter
```bash
# Formatear código automáticamente
r2fmt --style=siigo archivo.r2

# Configuración .r2fmt.json
{
    "indentSize": 4,
    "maxLineLength": 100,
    "alignAssignments": true,
    "sortImports": true
}
```

### R2Lang Test Framework
```r2
// Framework de testing integrado
test "calcular IVA México" {
    let resultado = calcularIVA(100, "MX")
    assert(resultado == 16, "IVA de México debe ser 16%")
}

test "balance debe cuadrar" {
    let asiento = crearAsiento(...)
    let sumaDebe = asiento.movimientos
        .filter(m => m.tipo == "DEBE")
        .map(m => m.monto)
        .reduce((a, b) => a + b, 0)
    
    let sumaHaber = asiento.movimientos
        .filter(m => m.tipo == "HABER")
        .map(m => m.monto)
        .reduce((a, b) => a + b, 0)
    
    assert(sumaDebe == sumaHaber, "El asiento no está balanceado")
}

// Ejecutar tests
r2test --coverage --watch
```

## 11. 🌐 DSL para Configuración

### Propuesta
```r2
// DSL para configurar aplicaciones
config AppConfig {
    server {
        port: getEnv("PORT", 8080)
        host: "0.0.0.0"
        timeout: 30s
    }
    
    database {
        driver: "sqlite"
        path: "./data/contabilidad.db"
        
        pool {
            min: 2
            max: 10
            idle: 5m
        }
    }
    
    regiones {
        MX {
            nombre: "México"
            iva: 16%
            moneda: "MXN"
        }
        // ... más regiones
    }
}

// Uso
let config = AppConfig.parse()
http.serve(`:${config.server.port}`)
```

## 12. 🔐 Seguridad Integrada

### Propuesta
```r2
// Validación automática de entrada
@validate
func procesarTransaccion(
    tipo: "venta" | "compra",
    region: string @matches(/^[A-Z]{2,3}$/),
    monto: number @min(0) @max(1000000)
) {
    // Los parámetros ya están validados
}

// Sanitización automática
let html = sanitize.html(userInput, {
    allowedTags: ["p", "br", "strong", "em"],
    allowedAttributes: {}
})

// Prevención de inyección SQL
let resultado = db.query(
    sql`SELECT * FROM transacciones WHERE region = ${region}`,
    // Los parámetros se escapan automáticamente
)
```

## 🎯 Conclusión

Estas mejoras posicionarían a R2Lang como una plataforma moderna y poderosa para el desarrollo de aplicaciones empresariales. La combinación de:

1. **Mejores estructuras de datos** (arrays anidados funcionales)
2. **DSL Builder más potente** (gramáticas expresivas, composición)
3. **Sistema de tipos opcional** (mejor tooling sin perder flexibilidad)
4. **Características modernas** (async/await, pattern matching)
5. **Herramientas profesionales** (LSP, formatter, testing)

Haría de R2Lang una opción atractiva para empresas como Siigo que buscan:
- Desarrollo rápido
- Código mantenible
- Localización eficiente
- Reducción de costos

La implementación gradual de estas mejoras, priorizando las más críticas (arrays anidados, template literals), permitiría una adopción suave mientras se mantiene la compatibilidad con código existente.