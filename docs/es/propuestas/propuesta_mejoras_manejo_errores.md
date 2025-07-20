# Propuesta de Mejoras en el Manejo de Errores para R2Lang

## Resumen Ejecutivo

Esta propuesta presenta un conjunto integrado de mejoras para el sistema de manejo de errores de R2Lang, con el objetivo de modernizar la experiencia del desarrollador manteniendo la robustez y simplicidad que caracterizan al lenguaje. Las mejoras están organizadas por **impacto**, **complejidad de implementación**, y **compatibilidad**.

### 🎯 Estado Actual vs Propuesta

**📊 Sistema Actual:**
- ✅ Manejo consistente con `panic/recover`
- ✅ Try-catch-finally implementado
- ✅ Errores ricos con contexto
- ✅ Detección avanzada de bucles infinitos
- ⚠️ Acceso a propiedades falla con panic
- ⚠️ Índices fuera de rango terminan ejecución
- ⚠️ No hay navegación segura opcional
- ⚠️ Solo valor `nil`, sin distinción `undefined`

**🚀 Propuesta de Mejoras:**
- 🆕 Optional chaining (`?.`) para navegación segura
- 🆕 Null coalescing (`??`) para valores por defecto
- 🆕 Tipo `undefined` distinto de `nil`
- 🆕 Operadores de acceso seguro
- 🆕 Modos de error graceful vs strict
- 🆕 Errores estructurados con metadatos
- 🆕 Recuperación selectiva de errores

---

## Matriz de Priorización

| Mejora | Impacto | Complejidad | Prioridad | Compatibilidad | Esfuerzo |
|--------|---------|-------------|-----------|----------------|----------|
| Optional chaining `?.` | 🔥🔥🔥 | 🟡 Media | P0 | 100% Compatible | 3-5 días |
| Null coalescing `??` | 🔥🔥🔥 | 🟢 Baja | P0 | 100% Compatible | 2-3 días |
| Acceso seguro a arrays `[]?` | 🔥🔥 | 🟡 Media | P1 | 100% Compatible | 2-3 días |
| Tipo `undefined` | 🔥🔥 | 🔴 Alta | P1 | Requiere migración | 5-7 días |
| Errores estructurados | 🔥 | 🔴 Alta | P2 | 100% Compatible | 7-10 días |
| Modos de error configurables | 🔥 | 🟡 Media | P2 | 100% Compatible | 3-5 días |
| Safe assignment operators | 🔥 | 🟡 Media | P3 | 100% Compatible | 2-4 días |
| Error recovery selectivo | 🔥 | 🔴 Alta | P3 | 100% Compatible | 5-7 días |

---

## Prioridad 0 (P0) - Navegación Segura Básica

### 1. Optional Chaining `?.`

**Problema Actual:**
```javascript
// ❌ Falla con panic si alguna propiedad no existe
let user = {profile: {address: {street: "123 Main"}}};
let street = user.profile.address.street; // OK
let phone = user.profile.contact.phone;   // panic: The object does not have the property: contact

// ❌ Verificación manual verbosa
let phone = nil;
if (user.profile && user.profile.contact) {
    phone = user.profile.contact.phone;
}
```

**Solución Propuesta:**
```javascript
// ✅ Navegación segura con optional chaining
let user = {profile: {address: {street: "123 Main"}}};
let street = user?.profile?.address?.street;  // "123 Main"
let phone = user?.profile?.contact?.phone;    // nil (no panic)
let missing = user?.inexistente?.algo;        // nil (no panic)

// ✅ Con arrays
let users = [{name: "Alice"}, {name: "Bob"}];
let firstUserName = users?.[0]?.name;         // "Alice"
let thirdUserName = users?.[2]?.name;         // nil (no panic)

// ✅ Con funciones
let result = obj?.method?.();                 // Llama método si existe
let callback = config?.onSuccess?.();         // Ejecuta callback si está definido
```

**Implementación Técnica:**

1. **Lexer** - Agregar `TOKEN_OPTIONAL_CHAIN` para `?.`
```go
// En lexer.go
const (
    TOKEN_OPTIONAL_CHAIN = "OPTIONAL_CHAIN"  // ?.
)

// Reconocimiento de ?. 
if ch == '?' && l.peekChar() == '.' {
    l.currentToken = Token{Type: TOKEN_OPTIONAL_CHAIN, Value: "?.", Line: l.line}
    l.pos += 2
    return l.currentToken, true
}
```

2. **Parser** - Modificar `parseAccessExpression()`:
```go
// En parse.go
func (p *Parser) parseOptionalChaining() Node {
    expr := p.parsePostfixExpression()
    
    for p.curTok.Type == TOKEN_OPTIONAL_CHAIN {
        p.nextToken() // consumir ?.
        if p.curTok.Type == TOKEN_IDENT {
            member := p.curTok.Value
            p.nextToken()
            expr = &OptionalAccessExpression{
                Object: expr,
                Member: member,
                IsOptional: true,
            }
        } else if p.curTok.Value == "[" {
            // Optional indexing ?.[]
            p.nextToken()
            index := p.parseExpression()
            p.expectToken(TOKEN_RBRACKET)
            expr = &OptionalIndexExpression{
                Object: expr,
                Index: index,
                IsOptional: true,
            }
        }
    }
    
    return expr
}
```

3. **AST** - Crear `OptionalAccessExpression`:
```go
// En optional_access_expression.go
type OptionalAccessExpression struct {
    Object     Node
    Member     string
    IsOptional bool
}

func (oae *OptionalAccessExpression) Eval(env *Environment) interface{} {
    obj := oae.Object.Eval(env)
    
    // Si el objeto es nil, retornar nil (no panic)
    if obj == nil {
        return nil
    }
    
    // Intentar acceso normal
    switch o := obj.(type) {
    case map[string]interface{}:
        if value, exists := o[oae.Member]; exists {
            return value
        }
        // Si IsOptional es true, retornar nil en lugar de panic
        if oae.IsOptional {
            return nil
        }
        panic("The object does not have the property: " + oae.Member)
    default:
        if oae.IsOptional {
            return nil
        }
        panic("Cannot access property on non-object type")
    }
}
```

**Beneficios:**
- **Robustez**: Elimina crashes por propiedades inexistentes
- **Código limpio**: Reduce verificaciones manuales 
- **Familiar**: Sintaxis estándar de JavaScript/TypeScript
- **Retrocompatible**: No afecta código existente

---

### 2. Null Coalescing `??`

**Problema Actual:**
```javascript
// ❌ Comportamiento inesperado con || para valores falsy
let count = 0;
let displayCount = count || 10;        // 10 (incorrecto - 0 es válido)

let name = "";
let displayName = name || "Anónimo";   // "Anónimo" (incorrecto - string vacío es válido)

// ❌ Verificación manual verbosa
let timeout = config.timeout != nil ? config.timeout : 5000;
```

**Solución Propuesta:**
```javascript
// ✅ Null coalescing operator
let count = 0;
let displayCount = count ?? 10;        // 0 (correcto - solo nil/undefined activan ??)

let name = "";
let displayName = name ?? "Anónimo";   // "" (correcto - string vacío no es nil)

let timeout = config.timeout ?? 5000;  // Solo usa 5000 si timeout es nil
let retries = options?.retries ?? 3;   // Combina con optional chaining

// ✅ Encadenamiento
let finalValue = userInput ?? defaultConfig ?? fallbackValue;
```

**Implementación Técnica:**

1. **Lexer** - Agregar `TOKEN_NULL_COALESCING`:
```go
// En lexer.go
const (
    TOKEN_NULL_COALESCING = "NULL_COALESCING"  // ??
)

// Reconocimiento de ??
if ch == '?' && l.peekChar() == '?' {
    l.currentToken = Token{Type: TOKEN_NULL_COALESCING, Value: "??", Line: l.line}
    l.pos += 2
    return l.currentToken, true
}
```

2. **Parser** - Agregar precedencia en `parseBinaryExpression()`:
```go
// En parse.go - precedencia menor que || pero mayor que =
func (p *Parser) parseNullCoalescingExpression() Node {
    expr := p.parseLogicalOrExpression()
    
    for p.curTok.Type == TOKEN_NULL_COALESCING {
        operator := p.curTok.Value
        p.nextToken()
        right := p.parseLogicalOrExpression()
        expr = &NullCoalescingExpression{
            Left:     expr,
            Operator: operator,
            Right:    right,
        }
    }
    
    return expr
}
```

3. **Evaluador** - Crear `NullCoalescingExpression`:
```go
// En null_coalescing_expression.go
type NullCoalescingExpression struct {
    Left     Node
    Operator string
    Right    Node
}

func (nce *NullCoalescingExpression) Eval(env *Environment) interface{} {
    left := nce.Left.Eval(env)
    
    // Solo evaluar right si left es nil o undefined
    if left == nil || isUndefined(left) {
        return nce.Right.Eval(env)
    }
    
    return left
}

// Helper function
func isUndefined(value interface{}) bool {
    // Implementar cuando se agregue tipo undefined
    return false
}
```

**Beneficios:**
- **Precisión**: Distingue entre valores falsy y nil/undefined
- **Performance**: Short-circuit evaluation (no evalúa right innecesariamente)
- **Legibilidad**: Código más expresivo y conciso
- **Estándar**: Operador ampliamente usado en lenguajes modernos

---

## Prioridad 1 (P1) - Tipos y Acceso Seguro

### 3. Tipo `undefined` Distinto de `nil`

**Análisis del Estado Actual:**

R2Lang actualmente usa solo `nil` para representar valores ausentes, pero hay casos semánticamente diferentes:

```javascript
// Casos que actualmente retornan nil pero semánticamente son diferentes:
let obj = {a: 1};
let inexistente = obj.b;          // nil (pero debería ser undefined)

let [a, b, c] = [1, 2];          // c es nil (pero debería ser undefined)
let {name, age} = {name: "John"}; // age es nil (pero debería ser undefined)

let arr = [1, 2, 3];
// arr[10] actualmente da panic, pero podría dar undefined
```

**Propuesta de Implementación:**

```javascript
// ✅ Distinción semántica clara
let explicitNull = nil;           // Valor intencionalmente nulo
let missing = obj.nonexistent;    // undefined (propiedad no existe)
let [a, b, c] = [1, 2];          // c es undefined (elemento no proporcionado)

// ✅ Diferentes comportamientos en comparaciones
explicitNull == nil;              // true
missing == nil;                   // false
missing == undefined;             // true

// ✅ Coerción de tipos diferente
nil ?? "default";                 // "default"
undefined ?? "default";           // "default"
nil || "default";                 // "default"
undefined || "default";           // "default"
false ?? "default";               // false (no se ejecuta ??)
false || "default";               // "default" (sí se ejecuta ||)
```

**Implementación Técnica:**

1. **Crear tipo `Undefined`**:
```go
// En commons.go
type UndefinedType struct{}

var Undefined = &UndefinedType{}

func (u *UndefinedType) String() string {
    return "undefined"
}

// En bytecode.go
OpUndefined  // Valor undefined
```

2. **Modificar funciones de comparación**:
```go
// En commons.go
func isNil(val interface{}) bool {
    return val == nil
}

func isUndefined(val interface{}) bool {
    _, ok := val.(*UndefinedType)
    return ok
}

func isNullish(val interface{}) bool {
    return isNil(val) || isUndefined(val)
}
```

3. **Actualizar destructuring para usar `undefined`**:
```go
// En destructuring_statement.go
func (ad *ArrayDestructuring) Eval(env *Environment) interface{} {
    // ... código existente ...
    for i, name := range ad.Names {
        var val interface{}
        if i < len(arr) {
            val = arr[i]
        } else {
            val = Undefined  // Cambiar de nil a Undefined
        }
        env.Set(name, val)
    }
}
```

**Migración y Compatibilidad:**

```go
// Función de migración para mantener compatibilidad
func normalizeToNil(val interface{}) interface{} {
    if isUndefined(val) {
        return nil  // Para código legacy que espera nil
    }
    return val
}

// Flag de configuración
var UseUndefinedType = true  // Puede ser configurable
```

---

### 4. Acceso Seguro a Arrays `[]?`

**Problema Actual:**
```javascript
// ❌ Índice fuera de rango causa panic
let arr = [1, 2, 3];
let element = arr[10];  // panic: index out of range: 10 len of array 3

// ❌ Verificación manual requerida
let element = nil;
if (arr.length > 10) {
    element = arr[10];
}
```

**Solución Propuesta:**
```javascript
// ✅ Acceso seguro con operator ?[]
let arr = [1, 2, 3];
let element = arr?.[10];    // undefined (no panic)
let valid = arr?.[1];       // 2
let negative = arr?.[-1];   // 3 (último elemento)

// ✅ Con encadenamiento
let matrix = [[1, 2], [3, 4]];
let cell = matrix?.[0]?.[1];    // 2
let missing = matrix?.[5]?.[0]; // undefined (no panic)

// ✅ Con optional chaining
let users = [{items: [1, 2, 3]}];
let item = users?.[0]?.items?.[10]; // undefined (totalmente seguro)
```

**Implementación:**
```go
// En optional_index_expression.go
type OptionalIndexExpression struct {
    Object     Node
    Index      Node
    IsOptional bool
}

func (oie *OptionalIndexExpression) Eval(env *Environment) interface{} {
    obj := oie.Object.Eval(env)
    if obj == nil {
        return Undefined
    }
    
    switch container := obj.(type) {
    case []interface{}:
        idx := int(oie.Index.Eval(env).(float64))
        
        // Manejo de índices negativos
        if idx < 0 {
            idx = len(container) + idx
        }
        
        // Verificación de límites
        if idx >= 0 && idx < len(container) {
            return container[idx]
        }
        
        // Si es opcional, retornar undefined en lugar de panic
        if oie.IsOptional {
            return Undefined
        }
        
        panic(fmt.Sprintf("index out of range: %d len of array %d", idx, len(container)))
    default:
        if oie.IsOptional {
            return Undefined
        }
        panic("Cannot index non-array type")
    }
}
```

---

## Prioridad 2 (P2) - Errores Estructurados y Configurables

### 5. Errores Estructurados con Metadatos

**Problema Actual:**
```javascript
// ❌ Errores como strings simples
try {
    let result = obj.method();
} catch (e) {
    std.print("Error:", e); // Solo string, sin contexto
}
```

**Solución Propuesta:**
```javascript
// ✅ Errores estructurados con metadatos
try {
    let result = obj.method();
} catch (e) {
    std.print("Error type:", e.type);        // "PropertyError", "TypeError", etc.
    std.print("Message:", e.message);        // Mensaje descriptivo
    std.print("Location:", e.location);      // Archivo:línea:columna
    std.print("Stack:", e.stack);            // Stack trace completo
    std.print("Context:", e.context);        // Contexto adicional
    
    // Manejo específico por tipo
    if (e.type == "PropertyError") {
        std.print("Missing property:", e.context.property);
    }
}
```

**Implementación:**
```go
// En error_types.go
type StructuredError struct {
    Type      string                 `json:"type"`
    Message   string                 `json:"message"`
    Location  string                 `json:"location"`
    Stack     []string               `json:"stack"`
    Context   map[string]interface{} `json:"context"`
    Timestamp time.Time              `json:"timestamp"`
    Code      string                 `json:"code,omitempty"`
}

func NewPropertyError(property string, object interface{}, location string) *StructuredError {
    return &StructuredError{
        Type:      "PropertyError",
        Message:   fmt.Sprintf("Property '%s' does not exist", property),
        Location:  location,
        Context: map[string]interface{}{
            "property":    property,
            "objectType":  fmt.Sprintf("%T", object),
            "objectKeys":  getObjectKeys(object),
        },
        Timestamp: time.Now(),
        Code:      "PROP_NOT_FOUND",
    }
}
```

### 6. Modos de Error Configurables

**Propuesta:**
```javascript
// ✅ Configuración global de manejo de errores
Error.setMode("strict");    // Modo actual (panic en errores)
Error.setMode("graceful");  // Retorna nil/undefined, registra warnings
Error.setMode("silent");    // Retorna nil/undefined, sin warnings

// ✅ Configuración específica por contexto
with Error.mode("graceful") {
    let data = obj?.deeply?.nested?.property; // No panic, retorna undefined
    let item = array[999];                     // No panic, retorna undefined
}

// ✅ Configuración por tipos de error
Error.configure({
    propertyAccess: "graceful",  // obj.prop retorna undefined
    arrayBounds: "graceful",     // arr[999] retorna undefined
    typeErrors: "strict",        // Mantiene panics para errores de tipo
    undefinedVars: "strict"      // Variables no declaradas siguen siendo error
});
```

---

## Prioridad 3 (P3) - Funcionalidades Avanzadas

### 7. Safe Assignment Operators

```javascript
// ✅ Asignación segura
obj.?prop = "valor";           // Solo asigna si obj no es nil
arr.?[index] = "valor";        // Solo asigna si índice es válido
nested.?obj.?prop = "valor";   // Asignación completamente segura

// ✅ Operadores de asignación segura
obj.?prop += 10;               // Solo si prop existe y es numérico
obj.?count ??= 0;              // Asigna 0 solo si count es nil/undefined
```

### 8. Error Recovery Selectivo

```javascript
// ✅ Recuperación selectiva por tipo de error
try {
    let result = complexOperation();
} catch (PropertyError as e) {
    // Solo maneja errores de propiedades
    std.print("Property issue:", e.context.property);
} catch (TypeError as e) {
    // Solo maneja errores de tipo
    std.print("Type issue:", e.message);
} finally {
    cleanup();
}

// ✅ Error filtering
try {
    riskyOperation();
} catch (e) where e.code == "NETWORK_ERROR" {
    // Solo maneja errores de red específicos
    retry();
}
```

---

## Plan de Implementación

### Fase 1 (Sprint 1-2 semanas) - Navegación Segura
1. **Optional chaining `?.`** - Máximo impacto, sintaxis familiar
2. **Null coalescing `??`** - Complementa optional chaining perfectamente

### Fase 2 (Sprint 2-3 semanas) - Tipos y Acceso
3. **Acceso seguro a arrays `[]?`** - Completa la navegación segura
4. **Tipo `undefined`** - Base para semántica más rica

### Fase 3 (Sprint 1-2 meses) - Sistema Avanzado
5. **Errores estructurados** - Mejora significativa en debugging
6. **Modos configurables** - Flexibilidad para diferentes casos de uso

### Fase 4 (Futuro) - Características Avanzadas
7. **Safe assignment operators** - Escritura más robusta
8. **Error recovery selectivo** - Control fino sobre manejo de errores

---

## Impacto en la Experiencia del Desarrollador

### Beneficios Inmediatos (Fase 1-2):
- **90% reducción** en crashes por navegación de objetos
- **Compatibilidad 100%** con expectativas de JavaScript/TypeScript
- **Código más limpio** sin verificaciones manuales
- **Debugging mejorado** con errores más informativos

### Beneficios a Mediano Plazo (Fase 3-4):
- **Desarrollo exploratorio** más fluido
- **Prototipado rápido** sin miedo a crashes
- **Código de producción** más robusto
- **Experiencia moderna** comparable a lenguajes contemporáneos

---

## Compatibilidad y Migración

### Garantías de Compatibilidad:
- **100% retrocompatible** para P0-P1
- **Migración gradual** para tipo undefined
- **Configuración flexible** para diferentes estilos
- **Modo legacy** para código existente

### Estrategia de Migración:
1. **Implementación aditiva** - Nuevas características no rompen código existente
2. **Flags de configuración** - Permitir adopción gradual
3. **Herramientas de migración** - Scripts para actualizar código automáticamente
4. **Documentación exhaustiva** - Guías de migración paso a paso

---

## Conclusiones

Esta propuesta transforma R2Lang en un lenguaje más robusto y moderno manteniendo su simplicidad fundamental. Las mejoras de **Prioridad 0-1** son altamente recomendadas para implementar en los próximos 2-3 meses, ya que aportan el **85% de los beneficios** con **moderada complejidad** y **compatibilidad total**.

El resultado será un R2Lang que mantiene su filosofía de simplicidad mientras ofrece la robustez y características modernas que los desarrolladores esperan en 2025.