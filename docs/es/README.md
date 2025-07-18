# R2Lang - Lenguaje de Programación Moderno

## Descripción General

R2Lang es un lenguaje de programación interpretado que combina la sintaxis familiar de JavaScript con características modernas como orientación a objetos, concurrencia nativa, sistema de testing integrado, y **nuestro innovador sistema DSL** para crear lenguajes específicos de dominio de forma elegante y sencilla.

## ✨ Características Principales

### 🌟 Constructor de DSL - Nuestra Característica Más Original

**El sistema DSL de R2Lang es nuestra característica más innovadora**, permitiendo crear lenguajes específicos de dominio con una sintaxis elegante y simple. Esto diferencia a R2Lang de otros lenguajes al hacer que la creación de parsers sea tan fácil como escribir una función.

#### ¿Por qué nuestro DSL es especial?

A diferencia de generadores de parsers complejos como ANTLR o Lex/Yacc, el sistema DSL de R2Lang es:
- **Integración Nativa**: El código DSL se ejecuta directamente en R2Lang
- **Configuración Cero**: Sin herramientas externas o generación de código
- **Sintaxis Intuitiva**: Declarativa y legible
- **Resultados Instantáneos**: De la idea a un parser funcional en minutos

#### Ejemplo Rápido de DSL

```r2
// Definir un DSL de calculadora simple
dsl Calculadora {
    token("NUMERO", "[0-9]+")
    token("SUMA", "\\+")
    token("RESTA", "-")
    token("MULT", "\\*")
    token("DIV", "/")
    
    rule("operacion", ["NUMERO", "operador", "NUMERO"], "calcular")
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
    }
    
    func op_suma(token) { return "+" }
    func op_resta(token) { return "-" }
    func op_mult(token) { return "*" }
    func op_div(token) { return "/" }
}

// Usar tu DSL
func main() {
    var calc = Calculadora.use
    
    var resultado = calc("15 + 25")
    console.log(resultado.Output)  // 40
    
    console.log(resultado)  // "DSL[15 + 25] -> 40"
}
```

#### DSL vs Parsers Tradicionales

| Característica | R2Lang DSL | ANTLR | Lex/Yacc |
|----------------|------------|-------|-----------|
| **Tiempo de Setup** | Minutos | Horas | Días |
| **Generación de Código** | Ninguna | Requerida | Requerida |
| **Curva de Aprendizaje** | Mínima | Empinada | Muy Empinada |
| **Integración** | Nativa | Externa | Externa |
| **Debugging** | Herramientas R2Lang | Especializado | Complejo |
| **Acceso a Resultados** | `resultado.Output` | Código generado | Código generado |

#### Casos de Uso del DSL

- **Lenguajes de Configuración**: Formatos de archivos de configuración personalizados
- **Sistemas de Comandos**: Lenguajes de comandos específicos de dominio
- **Validadores de Datos**: Reglas de validación personalizadas
- **Procesadores de Texto**: Parseo de texto especializado
- **Reglas de Negocio**: Lógica de negocio específica de dominio

#### Aprende Más

- [**Documentación Completa del DSL**](./dsl/) - Guía completa y ejemplos
- [**Ejemplos de DSL**](../../examples/dsl/) - Ejemplos de calculadora y comandos
- [**Referencia Rápida DSL**](./dsl/referencia_rapida.md) - Guía de referencia rápida

---

### 🚀 Sintaxis Intuitiva
```r2
func main() {
    let mensaje = "¡Hola, R2Lang!"
    print(mensaje)
}
```

### 🎯 Orientación a Objetos con Herencia
```r2
class Persona {
    let nombre
    let edad
    
    constructor(nombre, edad) {
        this.nombre = nombre
        this.edad = edad
    }
    
    saludar() {
        print("Hola, soy " + this.nombre)
    }
}

class Empleado extends Persona {
    let salario
    
    constructor(nombre, edad, salario) {
        super.constructor(nombre, edad)
        this.salario = salario
    }
    
    trabajar() {
        print(this.nombre + " está trabajando")
    }
}
```

### ⚡ Concurrencia Nativa
```r2
func tarea() {
    print("Ejecutando en paralelo")
    sleep(1)
    print("Tarea completada")
}

func main() {
    r2(tarea)  // Ejecuta en goroutine
    r2(tarea)
    sleep(2)   // Esperar completación
}
```

### 🧪 Sistema de Testing Integrado
```r2
TestCase "Verificar Suma" {
    Given func() { 
        setup()
        return "Preparando datos"
    }
    When func() {
        let resultado = 2 + 3
        return "Ejecutando operación"
    }
    Then func() {
        assertEqual(resultado, 5)
        return "Validando resultado"
    }
}
```

### 📦 Sistema de Módulos
```r2
import "math.r2" as math
import "./utils.r2" as utils

func main() {
    let resultado = math.sqrt(16)
    utils.log("Resultado: " + resultado)
}
```

## Instalación

### Requisitos Previos
- Go 1.23.4 o superior
- Git

### Clonar e Instalar
```bash
git clone https://github.com/arturoeanton/go-r2lang.git
cd go-r2lang
go build -o r2lang main.go
```

### Ejecutar Programas
```bash
# Ejecutar un archivo R2
./r2lang programa.r2

# Ejecutar main.r2 del directorio actual
./r2lang

# Modo interactivo REPL
./r2lang -repl

# REPL sin salida de debug
./r2lang -repl -no-output
```

## Ejemplos de Código

### Variables y Tipos
```r2
let numero = 42
let texto = "Hola mundo"
let bandera = true
let arreglo = [1, 2, 3, "cuatro"]
let mapa = {
    nombre: "Juan",
    edad: 30,
    activo: true
}
```

### Control de Flujo
```r2
// Condicionales
if (edad >= 18) {
    print("Mayor de edad")
} else {
    print("Menor de edad")
}

// Bucles
for (let i = 0; i < 10; i++) {
    print("Iteración: " + i)
}

// For-in para arreglos
for (let item in arreglo) {
    print("Item: " + item)
}

// While
let contador = 0
while (contador < 5) {
    print("Contador: " + contador)
    contador++
}
```

### Funciones y Lambdas
```r2
func sumar(a, b) {
    return a + b
}

// Función lambda
let multiplicar = func(a, b) {
    return a * b
}

// Función de orden superior
let numeros = [1, 2, 3, 4, 5]
let duplicados = numeros.map(func(x) { return x * 2 })
```

### Manejo de Errores
```r2
try {
    let resultado = dividir(10, 0)
    print("Resultado: " + resultado)
} catch (error) {
    print("Error: " + error)
} finally {
    print("Limpieza completada")
}
```

### Arrays y Operaciones
```r2
let frutas = ["manzana", "banana", "naranja"]

// Métodos de arrays
frutas.push("uva")                    // Agregar elemento
let longitud = frutas.len()           // Obtener longitud
let encontrado = frutas.find("banana") // Buscar elemento
let filtradas = frutas.filter(func(f) { 
    return f.length > 5 
})

// Operaciones funcionales
let numeros = [1, 2, 3, 4, 5]
let suma = numeros.reduce(func(acc, val) { 
    return acc + val 
})
let ordenados = numeros.sort()
```

## Bibliotecas Integradas

### Estándar
- `print()` - Salida por consola
- `len()` - Longitud de strings/arrays
- `typeOf()` - Tipo de variable
- `sleep()` - Pausa en segundos

### Matemáticas
- `math.sqrt()`, `math.pow()`, `math.sin()`, etc.
- `rand.int()`, `rand.float()` - Números aleatorios

### I/O y Sistema
- `io.readFile()`, `io.writeFile()` - Operaciones de archivo
- `os.getEnv()`, `os.exit()` - Interacción con SO

### HTTP
```r2
// Servidor HTTP
http.server(8080, func(req, res) {
    res.json({mensaje: "¡Hola desde R2Lang!"})
})

// Cliente HTTP
let respuesta = httpClient.get("https://api.ejemplo.com")
print(respuesta.body)
```

### Strings
```r2
let texto = "Hola Mundo"
let mayusculas = texto.upper()
let palabras = texto.split(" ")
let contiene = texto.contains("Mundo")
```

## Arquitectura del Intérprete

### Componentes Principales

1. **Lexer** (`r2lang/r2lang.go:139-321`)
   - Tokenización del código fuente
   - Manejo de números, strings, operadores
   - Soporte para comentarios de línea y bloque

2. **Parser** (`r2lang/r2lang.go:1662-2331`)
   - Análisis sintáctico recursivo descendente
   - Construcción del AST (Abstract Syntax Tree)
   - Manejo de precedencia de operadores

3. **AST y Evaluación** (`r2lang/r2lang.go:327-1657`)
   - Nodos del AST implementan interfaz `Node`
   - Tree-walking interpreter
   - Evaluación lazy de expresiones

4. **Environment** (`r2lang/r2lang.go:1429-1507`)
   - Sistema de scoping con environments anidados
   - Gestión de variables y funciones
   - Closure support

5. **Bibliotecas Nativas** (`r2lang/r2*.go`)
   - Funciones integradas en Go
   - Extensión modular del lenguaje

## Casos de Uso

### Scripting y Automatización
```r2
// Procesamiento de archivos
let contenido = io.readFile("datos.txt")
let lineas = contenido.split("\n")
let procesadas = lineas.map(func(linea) {
    return linea.trim().upper()
})
io.writeFile("salida.txt", procesadas.join("\n"))
```

### APIs y Microservicios
```r2
func manejarUsuarios(req, res) {
    if (req.method == "GET") {
        res.json(obtenerUsuarios())
    } else if (req.method == "POST") {
        let usuario = req.body
        crearUsuario(usuario)
        res.status(201).json({mensaje: "Usuario creado"})
    }
}

http.server(3000, manejarUsuarios)
```

### Testing y QA
```r2
TestCase "API de Usuarios" {
    Given func() {
        limpiarBaseDatos()
        return "Base de datos limpia"
    }
    When func() {
        let respuesta = httpClient.post("/api/usuarios", {
            nombre: "Juan",
            email: "juan@ejemplo.com"
        })
        return "Usuario creado via API"
    }
    Then func() {
        assertEqual(respuesta.status, 201)
        assertTrue(respuesta.body.id != null)
        return "Respuesta válida"
    }
}
```

## Roadmap

### Version 1.0 (Actual)
- ✅ Intérprete básico funcional
- ✅ Orientación a objetos con herencia
- ✅ Concurrencia con goroutines
- ✅ Sistema de testing integrado
- ✅ Bibliotecas básicas (I/O, HTTP, Math)

### Version 1.1 (Q1 2024)
- 🔄 Sistema de tipos mejorado
- 🔄 Manejo de errores robusto
- 🔄 Debugger integrado
- 🔄 Optimizaciones de rendimiento

### Version 1.5 (Q2 2024)
- 📋 Generics
- 📋 Pattern matching
- 📋 Async/await nativo
- 📋 Package manager

### Version 2.0 (Q3 2024)
- 📋 Compilación JIT
- 📋 Language Server Protocol
- 📋 WebAssembly target
- 📋 Standard library completa

## Contribuir

### Configuración del Entorno
```bash
git clone https://github.com/arturoeanton/go-r2lang.git
cd go-r2lang
go mod tidy
go test ./r2lang
```

### Estructura del Proyecto
```
├── main.go              # Punto de entrada
├── main.r2             # Ejemplo principal
├── r2lang/             # Core del intérprete
│   ├── r2lang.go       # Lexer, Parser, AST
│   ├── r2std.go        # Biblioteca estándar
│   ├── r2http.go       # Funciones HTTP
│   └── ...             # Otras bibliotecas
├── examples/           # Ejemplos de código
└── docs/              # Documentación
```

### Añadir Nuevas Características

1. **Nuevos Tokens**: Actualizar constantes en `r2lang.go:17-54`
2. **Sintaxis**: Modificar lexer en `r2lang.go:139-321`
3. **AST Nodes**: Crear structs que implementen `Node` interface
4. **Parsing**: Añadir lógica en parser `r2lang.go:1662-2331`
5. **Evaluación**: Implementar método `Eval()` en el nodo
6. **Testing**: Crear ejemplos en `examples/`

### Añadir Bibliotecas Nativas
```go
// r2lang/r2nueva.go
func RegisterNueva(env *Environment) {
    env.Set("nuevaFuncion", BuiltinFunction(func(args ...interface{}) interface{} {
        // Implementación
        return resultado
    }))
}

// En r2lang.go:RunCode(), añadir:
RegisterNueva(env)
```

## Licencia

Este proyecto está bajo la licencia MIT. Ver [LICENSE](../../LICENSE) para más detalles.

## Contacto y Soporte

- **Issues**: [GitHub Issues](https://github.com/arturoeanton/go-r2lang/issues)
- **Documentación**: Ver carpeta `docs/`
- **Ejemplos**: Ver carpeta `examples/`

---

*R2Lang - Simplicity meets Power* 🚀