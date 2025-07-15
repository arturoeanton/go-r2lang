# Propuesta: String Templates y Strings Multilínea en R2Lang

## Resumen Ejecutivo

Esta propuesta describe la implementación de **String Templates** (interpolación de strings) y **Strings Multilínea** en R2Lang, dos características fundamentales que mejorarán significativamente la expresividad y usabilidad del lenguaje.

## 1. String Templates (Interpolación de Strings)

### 1.1 Motivación

Actualmente, R2Lang requiere concatenación manual para combinar strings con variables:

```r2
// Actual - Concatenación manual
let nombre = "Juan";
let edad = 30;
let mensaje = "Hola " + nombre + ", tienes " + edad + " años";
```

Los string templates proporcionarían una sintaxis más limpia y legible:

```r2
// Propuesto - String Templates
let nombre = "Juan";
let edad = 30;
let mensaje = `Hola ${nombre}, tienes ${edad} años`;
```

### 1.2 Sintaxis Propuesta

#### 1.2.1 Sintaxis Básica
- **Delimitadores**: Backticks (`` ` ``) para delimitar string templates
- **Interpolación**: `${expresión}` para insertar valores dinámicos
- **Compatibilidad**: Mantener compatibilidad total con strings tradicionales (`"` y `'`)

#### 1.2.2 Ejemplos de Uso

```r2
// Variables simples
let usuario = "Ana";
let saludo = `Bienvenida ${usuario}`;

// Expresiones aritméticas
let a = 10, b = 20;
let resultado = `La suma de ${a} + ${b} = ${a + b}`;

// Llamadas a funciones
func obtenerFecha() { return "2025-07-15"; }
let mensaje = `Hoy es ${obtenerFecha()}`;

// Expresiones complejas
let objeto = {nombre: "Pedro", edad: 25};
let info = `Usuario: ${objeto.nombre}, Edad: ${objeto.edad}`;

// Anidamiento básico
let nivel = 2;
let estado = `Estado: ${nivel > 1 ? "Avanzado" : "Principiante"}`;
```

### 1.3 Características Técnicas

#### 1.3.1 Evaluación
- **Evaluación inmediata**: Las expresiones se evalúan al momento de crear el string template
- **Contexto de scope**: Las variables se resuelven en el scope actual donde se define el template
- **Optimización**: Usar el sistema de optimización de concatenación existente

#### 1.3.2 Escape de Caracteres
```r2
// Escape de backticks y ${
let literal = `Este es un backtick \` y esto es \${no_una_variable}`;
// Resultado: "Este es un backtick ` y esto es ${no_una_variable}"
```

### 1.4 Implementación Técnica

#### 1.4.1 Lexer (pkg/r2core/lexer.go)
```go
// Nuevo token para string templates
TOKEN_TEMPLATE_STRING

// Lógica de parsing para backticks y ${...}
func (l *Lexer) readTemplateString() string {
    // Implementar parsing de template strings
    // Manejar escape sequences
    // Identificar interpolaciones ${...}
}
```

#### 1.4.2 Parser (pkg/r2core/parse.go)
```go
// Nuevo nodo AST para template strings
type TemplateString struct {
    Parts []TemplatePart
}

type TemplatePart struct {
    IsExpression bool
    Content      string    // Para texto literal
    Expression   Node      // Para ${expresión}
}
```

#### 1.4.3 Evaluación
```go
func (ts *TemplateString) Eval(env *Environment) interface{} {
    var result strings.Builder
    
    for _, part := range ts.Parts {
        if part.IsExpression {
            value := part.Expression.Eval(env)
            result.WriteString(toStringOptimized(value))
        } else {
            result.WriteString(part.Content)
        }
    }
    
    return result.String()
}
```

## 2. Strings Multilínea

### 2.1 Motivación

Actualmente, crear strings de múltiples líneas requiere concatenación o caracteres de escape:

```r2
// Actual - Concatenación tediosa
let html = "<div>" +
           "  <h1>Título</h1>" +
           "  <p>Contenido</p>" +
           "</div>";
```

Los strings multilínea proporcionarían una sintaxis más natural:

```r2
// Propuesto - String Multilínea
let html = `<div>
  <h1>Título</h1>
  <p>Contenido</p>
</div>`;
```

### 2.2 Sintaxis Propuesta

#### 2.2.1 Usando Backticks
Los strings templates automáticamente soportarían multilínea:

```r2
let consulta_sql = `
    SELECT u.nombre, u.email, p.titulo
    FROM usuarios u
    JOIN publicaciones p ON u.id = p.usuario_id
    WHERE u.activo = true
    ORDER BY p.fecha_creacion DESC
`;
```

#### 2.2.2 Preservación de Formato
- **Saltos de línea**: Se preservan tal como aparecen en el código
- **Indentación**: Se respeta la indentación original
- **Espacios**: Se mantienen espacios en blanco significativos

### 2.3 Funcionalidades Avanzadas

#### 2.3.1 Indentación Inteligente
```r2
func generarConfig() {
    let servidor = "localhost";
    let puerto = 8080;
    
    return `
        servidor: ${servidor}
        puerto: ${puerto}
        ssl: false
        debug: true
    `.trim(); // Función trim() para eliminar espacios iniciales/finales
}
```

#### 2.3.2 Heredoc-style (Opcional)
Para casos donde se necesite control total sobre el formato:

```r2
let documento = `>>>END
Este es un documento
que preserva exactamente
    la indentación
        y formato
END`;
```

## 3. Compatibilidad y Migración

### 3.1 Compatibilidad Hacia Atrás
- **100% Compatible**: Todo código existente seguirá funcionando
- **Sin Breaking Changes**: Los strings con `"` y `'` mantienen su comportamiento actual
- **Coexistencia**: Los tres tipos de strings pueden usarse juntos

### 3.2 Migración Gradual
```r2
// Código existente (sigue funcionando)
let mensaje1 = "Hola " + nombre;

// Nueva sintaxis (opcional)
let mensaje2 = `Hola ${nombre}`;

// Ambas pueden coexistir
let combinado = mensaje1 + " y " + mensaje2;
```

## 4. Impacto en el Rendimiento

### 4.1 Optimizaciones Existentes
- **Reutilización**: Aprovechar el sistema de `StringBuilderPool` existente
- **Cache**: Usar `OptimizedStringConcat2` para concatenaciones eficientes
- **Pool de objetos**: Minimizar allocaciones durante la interpolación

### 4.2 Nuevas Optimizaciones
```go
// Cache para templates frecuentes
type TemplateCache struct {
    cache map[string]*CompiledTemplate
    mutex sync.RWMutex
}

// Template compilado para reutilización
type CompiledTemplate struct {
    StaticParts   []string
    ExprPositions []int
    Expressions   []Node
}
```

## 5. Casos de Uso Principales

### 5.1 Generación de HTML/XML
```r2
func crearPagina(titulo, contenido) {
    return `<!DOCTYPE html>
<html>
<head>
    <title>${titulo}</title>
</head>
<body>
    <h1>${titulo}</h1>
    <div class="contenido">
        ${contenido}
    </div>
</body>
</html>`;
}
```

### 5.2 Consultas SQL Dinámicas
```r2
func buscarUsuarios(nombre, activo) {
    let query = `
        SELECT id, nombre, email
        FROM usuarios
        WHERE 1=1
        ${nombre ? `AND nombre LIKE '%${nombre}%'` : ''}
        ${activo !== null ? `AND activo = ${activo}` : ''}
        ORDER BY nombre
    `;
    return ejecutarQuery(query);
}
```

### 5.3 Configuración y Logs
```r2
func logError(usuario, accion, error) {
    let timestamp = obtenerTimestamp();
    let logEntry = `[${timestamp}] ERROR - Usuario: ${usuario}
    Acción: ${accion}
    Error: ${error}
    Stack: ${obtenerStack()}`;
    
    escribirLog(logEntry);
}
```

### 5.4 APIs y JSON
```r2
func crearRespuestaAPI(datos, mensaje) {
    return `{
        "status": "success",
        "message": "${mensaje}",
        "timestamp": "${obtenerTimestamp()}",
        "data": ${JSON.stringify(datos)}
    }`;
}
```

## 6. Implementación por Fases

### Fase 1: String Templates Básicos (Sprint 1)
- [x] ✅ **Completado**: Concatenación optimizada existente
- [ ] Lexer: Reconocimiento de backticks y `${}`
- [ ] Parser: AST para TemplateString
- [ ] Evaluador: Interpolación básica
- [ ] Tests: Casos fundamentales

### Fase 2: Strings Multilínea (Sprint 2) 
- [ ] Soporte para saltos de línea en templates
- [ ] Preservación de formato e indentación
- [ ] Funciones auxiliares (trim, stripIndent)
- [ ] Tests: Casos multilínea complejos

### Fase 3: Optimizaciones (Sprint 3)
- [ ] Cache de templates compilados
- [ ] Optimización de memoria para templates grandes
- [ ] Profiling y benchmarks
- [ ] Documentación y ejemplos

### Fase 4: Características Avanzadas (Sprint 4)
- [ ] Escape sequences mejorados
- [ ] Anidamiento de templates
- [ ] Integración con herramientas de desarrollo
- [ ] Syntax highlighting actualizado

## 7. Testing y Validación

### 7.1 Test Suite Propuesta
```r2
// Tests básicos de interpolación
TestCase "String Template Básico" {
    Given let nombre = "Ana"
    When let resultado = `Hola ${nombre}`
    Then resultado == "Hola Ana"
}

// Tests de expresiones complejas
TestCase "Expresiones en Templates" {
    Given let a = 5, b = 3
    When let resultado = `${a} + ${b} = ${a + b}`
    Then resultado == "5 + 3 = 8"
}

// Tests multilínea
TestCase "String Multilínea" {
    Given let template = `Línea 1
Línea 2
Línea 3`
    When let lineas = template.split("\n")
    Then lineas.length == 3
}
```

### 7.2 Benchmarks de Rendimiento
- Comparar templates vs concatenación manual
- Medir uso de memoria
- Profiling de casos de uso reales

## 8. Documentación y Adopción

### 8.1 Documentación Propuesta
- **Guía de migración**: Cómo actualizar código existente
- **Mejores prácticas**: Cuándo usar cada tipo de string
- **Ejemplos reales**: Casos de uso comunes
- **Referencia técnica**: Sintaxis completa y limitaciones

### 8.2 Herramientas de Desarrollo
- **VS Code Extension**: Syntax highlighting para templates
- **JetBrains Plugin**: Soporte en IDEs de JetBrains
- **Linter Rules**: Reglas para uso consistente

## 9. Conclusión

La implementación de String Templates y Strings Multilínea en R2Lang representa una mejora significativa que:

1. **Mejora la legibilidad** del código considerablemente
2. **Mantiene compatibilidad total** con código existente
3. **Aprovecha optimizaciones existentes** del motor de strings
4. **Sigue estándares modernos** de otros lenguajes
5. **Reduce errores** en concatenación compleja

Esta propuesta posiciona a R2Lang como un lenguaje más expresivo y moderno, manteniendo su simplicidad y rendimiento característicos.

---

**Estado del Proyecto**: ✅ Concatenación optimizada implementada y probada
**Próximo Paso**: Implementar parsing básico de string templates en lexer.go
**Estimación**: 3-4 sprints para implementación completa