# Manual Completo de DSL en R2Lang

## Fecha de Actualización
2025-01-18

## Tabla de Contenidos

1. [Introducción](#introducción)
2. [Conceptos Básicos](#conceptos-básicos)
3. [Sintaxis del DSL](#sintaxis-del-dsl)
4. [Tokens](#tokens)
5. [Reglas](#reglas)
6. [Acciones Semánticas](#acciones-semánticas)
7. [Uso del DSL](#uso-del-dsl)
8. [Ejemplos Prácticos](#ejemplos-prácticos)
9. [Mejores Prácticas](#mejores-prácticas)
10. [Solución de Problemas](#solución-de-problemas)
11. [Referencia API](#referencia-api)

## Introducción

El sistema DSL (Domain-Specific Language) de R2Lang permite crear lenguajes específicos de dominio de manera simple y elegante. Con una sintaxis declarativa y integración nativa, puedes crear parsers personalizados en minutos.

### ¿Qué es un DSL?

Un DSL es un lenguaje de programación o especificación dedicado a un dominio específico. Ejemplos comunes incluyen:
- SQL para bases de datos
- HTML para páginas web
- CSS para estilos
- Regex para patrones de texto

### Ventajas del DSL de R2Lang

- **Sintaxis Simple**: Declarativa y fácil de entender
- **Integración Nativa**: Funciona perfectamente con R2Lang
- **Desarrollo Rápido**: Desde concepto hasta implementación en minutos
- **Resultados Estructurados**: Acceso fácil a resultados y metadatos
- **Debugging Natural**: Usa las herramientas de R2Lang

## Conceptos Básicos

### Componentes de un DSL

Un DSL en R2Lang consta de tres componentes principales:

1. **Tokens**: Elementos léxicos (palabras, símbolos, números)
2. **Reglas**: Gramática que define cómo se combinan los tokens
3. **Acciones**: Funciones que procesan los tokens matched

### Flujo de Procesamiento

```
Código DSL → Tokenización → Parsing → Acciones Semánticas → Resultado
```

## Sintaxis del DSL

### Estructura Básica

```r2
dsl NombreDSL {
    // Definición de tokens
    token("NOMBRE_TOKEN", "patrón_regex")
    
    // Definición de reglas
    rule("nombre_regla", ["token1", "token2"], "acción")
    
    // Funciones de acción
    func acción(param1, param2) {
        // Lógica de procesamiento
        return resultado
    }
}
```

### Ejemplo Mínimo

```r2
dsl Saludo {
    token("HOLA", "hola")
    token("NOMBRE", "[a-zA-Z]+")
    
    rule("saludo", ["HOLA", "NOMBRE"], "procesar_saludo")
    
    func procesar_saludo(hola, nombre) {
        return "¡Hola, " + nombre + "!"
    }
}
```

## Tokens

Los tokens son los elementos básicos del lenguaje. Cada token tiene un nombre y un patrón de expresión regular.

### Sintaxis de Tokens

```r2
token("NOMBRE_TOKEN", "patrón_regex")
```

### Ejemplos de Tokens

```r2
// Números enteros
token("NUMERO", "[0-9]+")

// Números decimales
token("DECIMAL", "[0-9]+\\.[0-9]+")

// Identificadores
token("ID", "[a-zA-Z][a-zA-Z0-9]*")

// Strings
token("STRING", "\"[^\"]*\"")

// Operadores
token("SUMA", "\\+")
token("RESTA", "-")
token("MULT", "\\*")
token("DIV", "/")

// Palabras clave
token("IF", "if")
token("ELSE", "else")
token("WHILE", "while")

// Espacios en blanco (usualmente ignorados)
token("ESPACIO", "\\s+")
```

### Patrones Regex Comunes

| Patrón | Descripción | Ejemplo |
|--------|-------------|---------|
| `[0-9]+` | Números enteros | `123`, `456` |
| `[0-9]*` | Números (incluyendo vacío) | `123`, `` |
| `[a-zA-Z]+` | Letras | `hello`, `World` |
| `[a-zA-Z0-9]*` | Alfanumérico | `hello123` |
| `\\+` | Símbolo + (escapado) | `+` |
| `\\*` | Símbolo * (escapado) | `*` |
| `"[^"]*"` | String con comillas | `"hello world"` |
| `\\s+` | Espacios en blanco | ` `, `\t`, `\n` |

### Consejos para Tokens

1. **Escape de Caracteres**: Los caracteres especiales deben ser escapados
2. **Orden Importante**: Los tokens más específicos primero
3. **Nombres Descriptivos**: Usa nombres claros y consistentes
4. **Regex Optimizado**: Usa `+` en lugar de `*` cuando sea posible

## Reglas

Las reglas definen la gramática del DSL, especificando cómo se combinan los tokens.

### Sintaxis de Reglas

```r2
rule("nombre_regla", ["token1", "token2", "token3"], "acción")
```

### Ejemplos de Reglas

```r2
// Regla simple
rule("suma", ["NUMERO", "SUMA", "NUMERO"], "sumar")

// Regla con múltiples tokens
rule("asignacion", ["ID", "IGUAL", "NUMERO"], "asignar")

// Regla recursiva (usando otras reglas)
rule("expresion", ["NUMERO", "OPERADOR", "expresion"], "evaluar")
```

### Reglas Múltiples

Puedes definir múltiples reglas para el mismo nombre (alternativas):

```r2
rule("operador", ["SUMA"], "op_suma")
rule("operador", ["RESTA"], "op_resta")
rule("operador", ["MULT"], "op_mult")
rule("operador", ["DIV"], "op_div")
```

### Reglas Anidadas

```r2
rule("expresion", ["termino", "SUMA", "expresion"], "sumar")
rule("expresion", ["termino", "RESTA", "expresion"], "restar")
rule("expresion", ["termino"], "termino_simple")

rule("termino", ["factor", "MULT", "termino"], "multiplicar")
rule("termino", ["factor", "DIV", "termino"], "dividir")
rule("termino", ["factor"], "factor_simple")

rule("factor", ["NUMERO"], "numero")
rule("factor", ["PAREN_IZQ", "expresion", "PAREN_DER"], "parentesis")
```

## Acciones Semánticas

Las acciones semánticas son funciones R2Lang que procesan los tokens matched.

### Sintaxis de Acciones

```r2
func nombre_accion(param1, param2, param3) {
    // Procesamiento
    return resultado
}
```

### Parámetros de Acciones

Los parámetros corresponden a los elementos de la regla en orden:

```r2
rule("suma", ["NUMERO", "SUMA", "NUMERO"], "sumar")

func sumar(num1, operador, num2) {
    // num1 = valor del primer NUMERO
    // operador = "+"
    // num2 = valor del segundo NUMERO
    return num1 + num2
}
```

### Ejemplos de Acciones

```r2
// Acción simple
func numero(n) {
    return n
}

// Acción con lógica
func calcular(izq, op, der) {
    if (op == "+") {
        return izq + der
    }
    if (op == "-") {
        return izq - der
    }
    return 0
}

// Acción con validación
func asignar(variable, igual, valor) {
    if (variable == "pi") {
        return "Error: no se puede asignar a constante"
    }
    return variable + " = " + valor
}
```

## Uso del DSL

### Crear Instancia del DSL

```r2
var mi_dsl = MiDSL.use
```

### Procesar Código

```r2
var resultado = mi_dsl("código a parsear")
```

### Acceder a Resultados

```r2
var resultado = mi_dsl("5 + 3")

// Resultado completo
console.log(resultado)  // "DSL[5 + 3] -> 8"

// Solo el resultado
console.log(resultado.Output)  // 8

// Código original
console.log(resultado.Code)  // "5 + 3"

// AST interno
console.log(resultado.AST)  // Estructura interna
```

## Ejemplos Prácticos

### 1. Calculadora Básica

```r2
dsl Calculadora {
    token("NUMERO", "[0-9]+")
    token("SUMA", "\\+")
    token("RESTA", "-")
    token("MULT", "\\*")
    token("DIV", "/")
    
    rule("expresion", ["NUMERO", "operador", "NUMERO"], "calcular")
    rule("operador", ["SUMA"], "op_suma")
    rule("operador", ["RESTA"], "op_resta")
    rule("operador", ["MULT"], "op_mult")
    rule("operador", ["DIV"], "op_div")
    
    func calcular(izq, op, der) {
        var a = parseFloat(izq)
        var b = parseFloat(der)
        
        if (op == "+") return a + b
        if (op == "-") return a - b
        if (op == "*") return a * b
        if (op == "/") return a / b
        
        return 0
    }
    
    func op_suma(token) { return "+" }
    func op_resta(token) { return "-" }
    func op_mult(token) { return "*" }
    func op_div(token) { return "/" }
}

func main() {
    var calc = Calculadora.use
    
    var resultado = calc("15 + 25")
    console.log(resultado.Output)  // 40
}
```

### 2. Sistema de Comandos

```r2
dsl ComandosSistema {
    token("CREAR", "crear")
    token("ELIMINAR", "eliminar")
    token("LISTAR", "listar")
    token("ARCHIVO", "archivo")
    token("DIRECTORIO", "directorio")
    token("NOMBRE", "[a-zA-Z0-9_]+")
    
    rule("comando", ["accion", "tipo", "NOMBRE"], "ejecutar")
    rule("comando", ["LISTAR", "tipo"], "listar")
    
    rule("accion", ["CREAR"], "crear")
    rule("accion", ["ELIMINAR"], "eliminar")
    rule("tipo", ["ARCHIVO"], "archivo")
    rule("tipo", ["DIRECTORIO"], "directorio")
    
    func ejecutar(accion, tipo, nombre) {
        return accion + " " + tipo + ": " + nombre
    }
    
    func listar(accion, tipo) {
        return "Listando " + tipo + "s"
    }
    
    func crear(token) { return "creando" }
    func eliminar(token) { return "eliminando" }
    func archivo(token) { return "archivo" }
    func directorio(token) { return "directorio" }
}

func main() {
    var cmd = ComandosSistema.use
    
    console.log(cmd("crear archivo config.txt"))  // "creando archivo: config.txt"
    console.log(cmd("listar directorio"))         // "Listando directorios"
}
```

### 3. Validador de Emails

```r2
dsl ValidadorEmail {
    token("USUARIO", "[a-zA-Z0-9._%+-]+")
    token("ARROBA", "@")
    token("DOMINIO", "[a-zA-Z0-9.-]+")
    token("PUNTO", "\\.")
    token("EXTENSION", "[a-zA-Z]{2,}")
    
    rule("email", ["USUARIO", "ARROBA", "dominio_completo"], "validar_email")
    rule("dominio_completo", ["DOMINIO", "PUNTO", "EXTENSION"], "validar_dominio")
    
    func validar_email(usuario, arroba, dominio) {
        if (usuario.length < 3) {
            return "Error: usuario muy corto"
        }
        return "Email válido: " + usuario + "@" + dominio
    }
    
    func validar_dominio(dominio, punto, extension) {
        if (extension.length < 2) {
            return "Error: extensión muy corta"
        }
        return dominio + "." + extension
    }
}

func main() {
    var validador = ValidadorEmail.use
    
    var resultado = validador("usuario@dominio.com")
    console.log(resultado.Output)  // "Email válido: usuario@dominio.com"
}
```

### 4. Lenguaje de Configuración

```r2
dsl ConfigDSL {
    token("ID", "[a-zA-Z][a-zA-Z0-9_]*")
    token("IGUAL", "=")
    token("STRING", "\"[^\"]*\"")
    token("NUMERO", "[0-9]+")
    token("BOOL", "true|false")
    
    rule("configuracion", ["ID", "IGUAL", "valor"], "asignar")
    rule("valor", ["STRING"], "string_val")
    rule("valor", ["NUMERO"], "numero_val")
    rule("valor", ["BOOL"], "bool_val")
    
    func asignar(nombre, igual, valor) {
        return {
            "nombre": nombre,
            "valor": valor,
            "tipo": typeof(valor)
        }
    }
    
    func string_val(s) {
        return s.substring(1, s.length - 1)  // Quitar comillas
    }
    
    func numero_val(n) {
        return parseFloat(n)
    }
    
    func bool_val(b) {
        return b == "true"
    }
}

func main() {
    var config = ConfigDSL.use
    
    var resultado = config('host = "localhost"')
    console.log(resultado.Output)  // {"nombre": "host", "valor": "localhost", "tipo": "string"}
}
```

## Mejores Prácticas

### 1. Diseño de Tokens

```r2
// ✅ Bueno: Tokens específicos y claros
token("NUMERO_ENTERO", "[0-9]+")
token("NUMERO_DECIMAL", "[0-9]+\\.[0-9]+")
token("IDENTIFICADOR", "[a-zA-Z][a-zA-Z0-9_]*")

// ❌ Malo: Tokens ambiguos
token("COSA", ".*")
token("TEXTO", ".+")
```

### 2. Organización de Reglas

```r2
// ✅ Bueno: Reglas jerárquicas
rule("programa", ["declaraciones"], "procesar_programa")
rule("declaraciones", ["declaracion", "declaraciones"], "lista_declaraciones")
rule("declaracion", ["variable", "IGUAL", "expresion"], "asignar_variable")

// ❌ Malo: Reglas planas y confusas
rule("todo", ["ID", "IGUAL", "NUMERO", "PUNTO", "ID"], "hacer_algo")
```

### 3. Acciones Semánticas

```r2
// ✅ Bueno: Acciones específicas y descriptivas
func sumar_numeros(izq, op, der) {
    return parseFloat(izq) + parseFloat(der)
}

func crear_variable(nombre, igual, valor) {
    return {
        "tipo": "variable",
        "nombre": nombre,
        "valor": valor
    }
}

// ❌ Malo: Acciones genéricas
func procesar(a, b, c) {
    return a + b + c
}
```

### 4. Manejo de Errores

```r2
func validar_numero(n) {
    if (n == "") {
        return "Error: número vacío"
    }
    
    var num = parseFloat(n)
    if (isNaN(num)) {
        return "Error: no es un número válido"
    }
    
    return num
}
```

### 5. Documentación

```r2
dsl MiDSL {
    // Token para números enteros positivos
    token("NUMERO", "[0-9]+")
    
    // Token para operador suma
    token("SUMA", "\\+")
    
    // Regla para expresiones de suma simple
    rule("suma", ["NUMERO", "SUMA", "NUMERO"], "sumar")
    
    // Función que realiza la suma de dos números
    func sumar(a, op, b) {
        return parseFloat(a) + parseFloat(b)
    }
}
```

## Solución de Problemas

### Problema: Tokens no reconocidos

```r2
// Error: unexpected character at position 5: '+'
// Solución: Escapar caracteres especiales
token("SUMA", "\\+")  // ✅ Correcto
token("SUMA", "+")    // ❌ Incorrecto
```

### Problema: Reglas no coinciden

```r2
// Error: no alternative matched for rule 'expresion'
// Verificar que los tokens y reglas coincidan exactamente
rule("expresion", ["NUMERO", "SUMA", "NUMERO"], "sumar")
// Asegurar que existan tokens NUMERO y SUMA
```

### Problema: Parámetros incorrectos

```r2
// Error: función recibe parámetros incorrectos
rule("suma", ["A", "B", "C"], "sumar")

func sumar(x, y, z) {  // ✅ Correcto: 3 parámetros
    return x + y + z
}

func sumar(x, y) {     // ❌ Incorrecto: 2 parámetros
    return x + y
}
```

### Problema: Regex inválido

```r2
// Error: regex compilation failed
token("NUMERO", "[0-9")  // ❌ Incorrecto: bracket no cerrado
token("NUMERO", "[0-9]+")  // ✅ Correcto
```

### Problema: Resultado no accesible

```r2
// Para acceder al resultado
var resultado = dsl.use("input")
console.log(resultado.Output)  // ✅ Correcto
console.log(resultado)         // ❌ Muestra objeto completo
```

## Referencia API

### DSLResult

El objeto resultado que retorna el DSL:

```r2
interface DSLResult {
    Output: any        // Resultado final de la evaluación
    Code: string       // Código original procesado
    AST: any          // Árbol de sintaxis abstracta
    GetResult(): any   // Método para obtener el resultado
}
```

### Métodos Disponibles

```r2
// Crear instancia
var dsl = MiDSL.use

// Procesar código
var resultado = dsl("código")

// Acceder a propiedades
resultado.Output    // Resultado final
resultado.Code      // Código original
resultado.AST       // AST interno
resultado.GetResult()  // Método getter
```

### Representación String

```r2
console.log(resultado.toString())  // "DSL[código] -> resultado"
```

## Casos de Uso Avanzados

### 1. DSL con Estado

```r2
dsl ContadorDSL {
    token("INCREMENTAR", "\\+\\+")
    token("DECREMENTAR", "--")
    token("RESET", "reset")
    
    rule("comando", ["INCREMENTAR"], "incrementar")
    rule("comando", ["DECREMENTAR"], "decrementar")
    rule("comando", ["RESET"], "reset")
    
    func incrementar(op) {
        // Usar variables globales para estado
        contador = contador + 1
        return contador
    }
    
    func decrementar(op) {
        contador = contador - 1
        return contador
    }
    
    func reset(op) {
        contador = 0
        return contador
    }
}
```

### 2. DSL con Validación

```r2
dsl ValidadorDSL {
    token("EMAIL", "[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}")
    token("TELEFONO", "\\+?[0-9]{10,}")
    token("VALIDAR", "validar")
    
    rule("validacion", ["VALIDAR", "tipo"], "validar")
    rule("tipo", ["EMAIL"], "email")
    rule("tipo", ["TELEFONO"], "telefono")
    
    func validar(cmd, valor) {
        if (valor.includes("@")) {
            return validar_email(valor)
        } else {
            return validar_telefono(valor)
        }
    }
    
    func validar_email(email) {
        // Lógica de validación
        return email.includes("@") && email.includes(".")
    }
    
    func validar_telefono(tel) {
        // Lógica de validación
        return tel.length >= 10
    }
}
```

### 3. DSL con Transformaciones

```r2
dsl TransformadorDSL {
    token("TEXTO", "[a-zA-Z ]+")
    token("MAYUSCULAS", "upper")
    token("MINUSCULAS", "lower")
    token("TITULO", "title")
    
    rule("transformacion", ["comando", "TEXTO"], "transformar")
    rule("comando", ["MAYUSCULAS"], "upper")
    rule("comando", ["MINUSCULAS"], "lower")
    rule("comando", ["TITULO"], "title")
    
    func transformar(cmd, texto) {
        if (cmd == "upper") {
            return texto.toUpperCase()
        }
        if (cmd == "lower") {
            return texto.toLowerCase()
        }
        if (cmd == "title") {
            return texto.charAt(0).toUpperCase() + texto.slice(1).toLowerCase()
        }
        return texto
    }
}
```

## Conclusión

El sistema DSL de R2Lang proporciona una herramienta poderosa y fácil de usar para crear lenguajes específicos de dominio. Con su sintaxis declarativa, integración nativa, y resultados estructurados, es ideal para:

- Prototipado rápido de lenguajes
- Procesamiento de configuraciones
- Validación de datos
- Sistemas de comandos
- Transformación de texto

La combinación de simplicidad y potencia hace que el DSL de R2Lang sea una excelente opción para desarrolladores que necesitan crear parsers personalizados de manera eficiente.

---

**Actualizado por**: Claude Code  
**Fecha**: 2025-01-18  
**Versión**: Manual DSL v2.0