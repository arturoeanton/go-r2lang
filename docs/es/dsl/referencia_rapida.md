# Referencia Rápida DSL R2Lang

## Plantilla Básica

```r2
dsl NombreDSL {
    // Tokens
    token("NOMBRE", "regex")
    
    // Reglas
    rule("regla", ["TOKEN1", "TOKEN2"], "accion")
    
    // Acciones
    func accion(param1, param2) {
        return resultado
    }
}

// Uso
var dsl = NombreDSL.use
var resultado = dsl("input")
console.log(resultado.Output)
```

## Tokens Comunes

```r2
// Números
token("NUMERO", "[0-9]+")
token("DECIMAL", "[0-9]+\\.[0-9]+")

// Texto
token("ID", "[a-zA-Z][a-zA-Z0-9_]*")
token("STRING", "\"[^\"]*\"")

// Operadores
token("SUMA", "\\+")
token("RESTA", "-")
token("MULT", "\\*")
token("DIV", "/")
token("IGUAL", "=")

// Separadores
token("COMA", ",")
token("PUNTO", "\\.")
token("PAREN_IZQ", "\\(")
token("PAREN_DER", "\\)")
```

## Patrones Regex

| Patrón | Descripción |
|--------|-------------|
| `[0-9]+` | Números enteros |
| `[a-zA-Z]+` | Solo letras |
| `[a-zA-Z0-9_]*` | Identificadores |
| `\\+` | Símbolo + |
| `\\*` | Símbolo * |
| `"[^"]*"` | String con comillas |
| `\\s+` | Espacios |
| `true\|false` | Booleanos |

## Ejemplos Rápidos

### Calculadora Simple
```r2
dsl Calc {
    token("NUM", "[0-9]+")
    token("PLUS", "\\+")
    rule("sum", ["NUM", "PLUS", "NUM"], "add")
    func add(a, op, b) { return parseFloat(a) + parseFloat(b) }
}
```

### Validador Email
```r2
dsl Email {
    token("USER", "[a-zA-Z0-9._%+-]+")
    token("AT", "@")
    token("DOMAIN", "[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}")
    rule("email", ["USER", "AT", "DOMAIN"], "validate")
    func validate(u, at, d) { return u + "@" + d }
}
```

### Comandos Simple
```r2
dsl Cmd {
    token("CREATE", "create")
    token("DELETE", "delete")
    token("NAME", "[a-zA-Z0-9_]+")
    rule("command", ["CREATE", "NAME"], "create")
    rule("command", ["DELETE", "NAME"], "delete")
    func create(cmd, name) { return "Creating " + name }
    func delete(cmd, name) { return "Deleting " + name }
}
```

## Acceso a Resultados

```r2
var resultado = dsl("input")

// Propiedades
resultado.Output    // Resultado final
resultado.Code      // Código original
resultado.AST       // AST interno

// Métodos
resultado.GetResult()  // Obtener resultado
resultado.toString()   // Representación string
```

## Errores Comunes

### Regex sin escapar
```r2
// ❌ Incorrecto
token("PLUS", "+")

// ✅ Correcto
token("PLUS", "\\+")
```

### Parámetros incorrectos
```r2
// ❌ Incorrecto
rule("sum", ["A", "B", "C"], "add")
func add(x, y) { return x + y }  // Solo 2 parámetros

// ✅ Correcto
func add(x, y, z) { return x + y + z }  // 3 parámetros
```

### Tokens no definidos
```r2
// ❌ Incorrecto
rule("sum", ["NUM", "PLUS", "NUM"], "add")  // PLUS no definido

// ✅ Correcto
token("PLUS", "\\+")
rule("sum", ["NUM", "PLUS", "NUM"], "add")
```

## Debugging

```r2
// Ver resultado completo
console.log(resultado)  // "DSL[input] -> output"

// Ver solo resultado
console.log(resultado.Output)

// Ver código original
console.log(resultado.Code)
```

## Mejores Prácticas

1. **Nombres descriptivos**: `NUMERO` mejor que `N`
2. **Escape regex**: Siempre escapar caracteres especiales
3. **Validación**: Verificar parámetros en acciones
4. **Documentar**: Comentar tokens y reglas complejas
5. **Testear**: Probar con múltiples inputs

## Casos de Uso Típicos

- 📊 Calculadoras especializadas
- 🔧 Sistemas de configuración
- ✅ Validadores de datos
- 🎮 Lenguajes de comandos
- 📝 Procesadores de texto
- 🔍 Analizadores de logs