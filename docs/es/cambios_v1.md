# Cambios v1 - Soporte para true/false y Map Literales Estilo JavaScript

## Resumen

Este documento describe las mejoras implementadas en R2Lang v1, que incluyen soporte nativo para literales booleanos `true` y `false`, map literales estilo JavaScript, y soporte completo para iteración `for-in` sobre maps además de arrays.

## Cambios Implementados

### 1. Literales Booleanos `true` y `false`

#### Tokens Añadidos al Lexer
- **TRUE**: Reconoce `true` (insensible a mayúsculas/minúsculas)
- **FALSE**: Reconoce `false` (insensible a mayúsculas/minúsculas)

#### Parsing de Booleanos
- Añadido soporte para `TOKEN_TRUE` y `TOKEN_FALSE` en el parser
- Los literales `true` y `false` ahora se parsean como `BooleanLiteral` nodes
- Funciona con variaciones: `true`, `TRUE`, `false`, `FALSE`

#### Ejemplo de Uso
```r2
let isActive = true
let isDisabled = false
let result = TRUE  // También válido
let flag = FALSE   // También válido

print("Active: " + isActive)   // Output: Active: true
print("Disabled: " + isDisabled) // Output: Disabled: false
```

### 2. Map Literales Estilo JavaScript

#### Sintaxis Soportada
- **Claves string**: `{"name": "Juan", "age": 30}`
- **Claves identificador**: `{name: "Juan", age: 30}` (equivale a `{"name": "Juan", "age": 30}`)
- **Claves numéricas**: `{123: "valor", 456: "otro"}` 
- **Claves computadas**: `{[expression]: valor}`
- **Claves mixtas**: Combinación de todos los tipos anteriores

#### Ejemplos de Map Literales
```r2
// Sintaxis básica
let usuario = {name: "Juan", age: 30, active: true}

// Claves string explícitas
let config = {"host": "localhost", "port": 8080}

// Claves numéricas
let códigos = {200: "OK", 404: "Not Found", 500: "Server Error"}

// Claves computadas
let prefix = "user"
let data = {[prefix + "_id"]: 123, [prefix + "_name"]: "Juan"}

// Maps con valores booleanos
let permisos = {read: true, write: false, admin: true}
```

### 3. For-In Loop Mejorado

#### Soporte para Maps
El loop `for-in` ahora funciona tanto con arrays como con maps:

```r2
// Con arrays (funcionaba antes)
let arr = ["a", "b", "c"]
for (let i in arr) {
    print("Índice: " + i + ", Valor: " + arr[i])
}
// Output: Índice: 0, Valor: a
//         Índice: 1, Valor: b  
//         Índice: 2, Valor: c

// Con maps (NUEVO)
let user = {name: "Juan", age: 30, city: "Madrid"}
for (let key in user) {
    print("Clave: " + key + ", Valor: " + user[key])
}
// Output: Clave: name, Valor: Juan
//         Clave: age, Valor: 30
//         Clave: city, Valor: Madrid
```

#### Variables Especiales $k y $v
Durante la iteración, están disponibles las variables especiales:
- **$k**: La clave actual (índice para arrays, string para maps)
- **$v**: El valor actual

```r2
let data = {x: 10, y: 20, z: 30}
for (let key in data) {
    print($k + " = " + $v)  // $k es la clave, $v es el valor
}
// Output: x = 10
//         y = 20
//         z = 30
```

## Mejoras Técnicas

### 1. Lexer (pkg/r2core/lexer.go)
- Añadidos tokens `TOKEN_TRUE` y `TOKEN_FALSE`
- Reconocimiento case-insensitive de `true` y `false`
- Parsing correcto en el switch de identificadores

### 2. Parser (pkg/r2core/parse.go)
- Soporte para parsing de `TOKEN_TRUE` y `TOKEN_FALSE` como `BooleanLiteral`
- Mejoras en `parseMapLiteral()` para soportar:
  - Claves numéricas
  - Claves computadas con `[expression]: value`
  - Sintaxis estilo JavaScript

### 3. For Statement (pkg/r2core/for_statement.go)
- La función `evalForIn()` ya soportaba maps
- Verificación de tipo tanto para `[]interface{}` como `map[string]interface{}`
- Iteración correcta sobre claves de maps

### 4. Literals (pkg/r2core/literals.go)
- `BooleanLiteral` ya existía y funciona correctamente
- Evaluación directa de valores booleanos

## Tests Implementados

### 1. Tests de Literales Booleanos
```go
func TestBooleanLiteralFromTokens(t *testing.T)
```
- Verifica parsing de `true`, `false`, `TRUE`, `FALSE`
- Confirma evaluación correcta de literales booleanos

### 2. Tests de Map Literales
```go
func TestMapLiteral_BasicKeys(t *testing.T)
func TestMapLiteral_ComputedKeys(t *testing.T)
```
- Claves string, identificador, numéricas
- Valores booleanos con keywords `true`/`false`
- Claves computadas con expresiones

### 3. Tests de For-In
- Iteración sobre maps y arrays
- Variables especiales `$k` y `$v`
- Break y continue en iteración
- Maps anidados

## Compatibilidad

### Retro-compatibilidad
- ✅ Todo el código R2Lang existente sigue funcionando
- ✅ Arrays con for-in siguen funcionando igual
- ✅ Maps existentes siguen funcionando
- ✅ No hay breaking changes

### Nuevas Características
- ✅ `true` y `false` como literales nativos
- ✅ Map literales con sintaxis JavaScript
- ✅ For-in funciona con maps
- ✅ Claves computadas en maps
- ✅ Soporte case-insensitive para booleanos

## Ejemplos Prácticos

### Configuración de Aplicación
```r2
let config = {
    server: {
        host: "localhost",
        port: 8080,
        ssl: true
    },
    database: {
        type: "postgres", 
        host: "db.example.com",
        port: 5432,
        ssl: true
    },
    features: {
        auth: true,
        logging: true,
        debug: false
    }
}

// Iterar sobre características
for (let feature in config.features) {
    if (config.features[feature]) {
        print("Característica habilitada: " + feature)
    }
}
```

### Validación de Datos
```r2
let user = {name: "Juan", email: "juan@email.com", active: true}
let requiredFields = ["name", "email"]
let isValid = true

for (let field in requiredFields) {
    let fieldName = requiredFields[field]
    if (!user[fieldName]) {
        print("Campo requerido faltante: " + fieldName)
        isValid = false
    }
}

if (isValid && user.active) {
    print("Usuario válido y activo")
} else {
    print("Usuario inválido o inactivo")
}
```

### Transformación de Datos
```r2
let scores = {alice: 95, bob: 87, charlie: 92, diana: 88}
let grades = {}

for (let student in scores) {
    let score = scores[student]
    let grade = "F"
    
    if (score >= 90) {
        grade = "A"
    } else if (score >= 80) {
        grade = "B"
    } else if (score >= 70) {
        grade = "C"
    }
    
    grades[student] = grade
}

// Mostrar resultados
for (let student in grades) {
    print(student + ": " + scores[student] + " -> " + grades[student])
}
```

## Conclusión

R2Lang v1 introduce soporte completo para booleanos nativos y map literales estilo JavaScript, manteniendo 100% de compatibilidad hacia atrás. Estas mejoras hacen que R2Lang sea más intuitivo para desarrolladores familiarizados con JavaScript, mientras mantiene su rendimiento y simplicidad característicos.

Los tests confirman que todas las características funcionan correctamente y que no hay regresiones en funcionalidad existente. El soporte para for-in sobre maps extiende significativamente las capacidades de iteración del lenguaje.