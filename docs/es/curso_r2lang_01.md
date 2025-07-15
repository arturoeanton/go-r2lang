# Curso R2Lang - Módulo 1: Introducción y Fundamentos

## Bienvenido a R2Lang

### ¿Qué es R2Lang?

R2Lang es un lenguaje de programación moderno que combina la simplicidad de JavaScript con características avanzadas como:

- **Testing integrado**: Framework BDD nativo
- **Orientación a objetos**: Clases con herencia
- **Concurrencia**: Primitivas simples para programación paralela
- **Sintaxis familiar**: Fácil de aprender para desarrolladores web

### ¿Por qué Aprender R2Lang?

1. **Sintaxis Simple**: Basada en JavaScript, fácil de leer y escribir
2. **Testing-First**: Ideal para aprender buenas prácticas de testing
3. **Versatilidad**: Scripts, web servers, automatización, prototipos
4. **Moderno**: Incluye características de lenguajes contemporáneos

## Instalación y Configuración

### Requisitos Previos

```bash
# Verificar que tienes Go instalado
go version
# Debe mostrar Go 1.23.4 o superior
```

### Instalación

```bash
# 1. Clonar el repositorio
git clone https://github.com/arturoeanton/go-r2lang.git
cd go-r2lang

# 2. Compilar el intérprete (Nueva arquitectura modular)
go build -o r2lang main.go

# 3. Verificar instalación
./r2lang --version
```

### Arquitectura Modular v2

R2Lang ha sido completamente reestructurado con una arquitectura modular que separa las responsabilidades:

```
R2Lang v2 Architecture:
├── pkg/r2core/     # Núcleo del intérprete (2,590 LOC)
│   ├── lexer.go    # Tokenización (330 LOC)
│   ├── parse.go    # Parsing (678 LOC)
│   ├── environment.go # Scoping (98 LOC)
│   └── [27 archivos AST especializados]
├── pkg/r2libs/     # Bibliotecas incorporadas (3,701 LOC)
│   ├── r2http.go   # Servidor HTTP (410 LOC)
│   ├── r2httpclient.go # Cliente HTTP (324 LOC)
│   ├── r2string.go # Manipulación de strings (194 LOC)
│   └── [15 bibliotecas más]
├── pkg/r2lang/     # Coordinador de alto nivel (45 LOC)
└── pkg/r2repl/     # REPL interactivo (185 LOC)
```

Esta nueva arquitectura elimina el anti-patrón "God Object" y proporciona:
- **Mejor mantenibilidad**: Cada módulo tiene una responsabilidad específica
- **Facilidad de testing**: Componentes independientes
- **Desarrollo paralelo**: Múltiples desarrolladores pueden trabajar simultáneamente
- **Calidad de código**: Pasamos de 6.2/10 a 8.5/10 en métricas de calidad

### Tu Primer Programa

Crea un archivo llamado `hola.r2`:

```r2
func main() {
    print("¡Hola, R2Lang!")
}
```

Ejecútalo:

```bash
# Ejecutar archivo específico
go run main.go hola.r2

# O usar el binario compilado
./r2lang hola.r2
```

Deberías ver: `¡Hola, R2Lang!`

### Comandos Básicos v2

```bash
# Ejecutar programas
go run main.go archivo.r2        # Ejecutar archivo específico
go run main.go                   # Ejecutar main.r2 (si existe)

# Modo interactivo REPL
go run main.go -repl             # REPL con output normal
go run main.go -repl -no-output  # REPL sin output automático

# Testing y desarrollo
go test ./pkg/...                # Ejecutar todos los tests
go test ./pkg/r2core/            # Tests del núcleo
go test ./pkg/r2libs/            # Tests de bibliotecas

# Ejemplos incluidos
go run main.go examples/example01.r2
go run main.go examples/example02.r2
```

## Conceptos Fundamentales

### 1. Variables y Tipos

#### Declaración de Variables

```r2
// Declaración con valor inicial
let nombre = "Juan"
let edad = 25
let activo = true

// Declaración sin valor inicial (undefined)
let resultado
```

#### Tipos de Datos Básicos

```r2
// Números (todos son float64 internamente)
let entero = 42
let decimal = 3.14159
let negativo = -100

// Strings (cadenas de texto)
let saludo = "Hola mundo"
let mensaje = 'También se pueden usar comillas simples'

// Booleanos
let verdadero = true
let falso = false

// Nulo
let vacio = null
let indefinido = nil  // Sinónimo de null
```

#### Verificación de Tipos

```r2
func main() {
    let numero = 42
    let texto = "Hola"
    let bandera = true
    
    print("Tipo de numero:", typeOf(numero))    // float64
    print("Tipo de texto:", typeOf(texto))      // string
    print("Tipo de bandera:", typeOf(bandera))  // bool
}
```

### 2. Operadores

#### Operadores Aritméticos

```r2
func main() {
    let a = 10
    let b = 3
    
    print("Suma:", a + b)        // 13
    print("Resta:", a - b)       // 7
    print("Multiplicación:", a * b)  // 30
    print("División:", a / b)    // 3.3333...
}
```

#### Operadores de Comparación

```r2
func main() {
    let x = 5
    let y = 10
    
    print("x == y:", x == y)     // false
    print("x != y:", x != y)     // true
    print("x < y:", x < y)       // true
    print("x > y:", x > y)       // false
    print("x <= y:", x <= y)     // true
    print("x >= y:", x >= y)     // false
}
```

#### Operadores Lógicos

```r2
func main() {
    let a = true
    let b = false
    
    print("a && b:", a && b)     // false (AND lógico)
    print("a || b:", a || b)     // true (OR lógico)
    print("!a:", !a)             // false (NOT lógico)
}
```

### 3. Strings (Cadenas de Texto)

#### Operaciones Básicas

```r2
func main() {
    let nombre = "Juan"
    let apellido = "Pérez"
    
    // Concatenación
    let nombreCompleto = nombre + " " + apellido
    print("Nombre completo:", nombreCompleto)
    
    // Longitud
    print("Longitud del nombre:", len(nombre))
    
    // Conversión a mayúsculas/minúsculas (usando built-ins)
    print("Mayúsculas:", nombre.upper())
    print("Minúsculas:", apellido.lower())
}
```

#### Métodos de String

```r2
func main() {
    let frase = "Hola mundo desde R2Lang"
    
    // Dividir en palabras
    let palabras = frase.split(" ")
    print("Palabras:", palabras)
    
    // Verificar si contiene texto
    let contiene = frase.contains("mundo")
    print("Contiene 'mundo':", contiene)
    
    // Longitud
    print("Longitud:", frase.length())
}
```

### 4. Entrada y Salida

#### Función print()

```r2
func main() {
    // Imprimir múltiples valores
    print("Nombre:", "Juan", "Edad:", 25)
    
    // Imprimir variables
    let mensaje = "¡Hola!"
    print(mensaje)
    
    // Imprimir resultado de operaciones
    print("2 + 2 =", 2 + 2)
}
```

## Ejercicios Prácticos

### Ejercicio 1: Variables y Operaciones

Crea un programa que:
1. Declare variables para tu nombre, edad y ciudad
2. Calcule tu año de nacimiento
3. Imprima una presentación completa

```r2
func main() {
    // Tu solución aquí
    let nombre = "Tu Nombre"
    let edad = 25
    let ciudad = "Tu Ciudad"
    
    let anoActual = 2024
    let anoNacimiento = anoActual - edad
    
    print("Hola, soy", nombre)
    print("Tengo", edad, "años")
    print("Vivo en", ciudad)
    print("Nací en el año", anoNacimiento)
}
```

### Ejercicio 2: Calculadora Básica

Crea un programa que realice operaciones matemáticas:

```r2
func main() {
    let a = 15
    let b = 4
    
    print("Números:", a, "y", b)
    print("Suma:", a + b)
    print("Resta:", a - b)
    print("Multiplicación:", a * b)
    print("División:", a / b)
    
    // Bonus: Verificar si a es mayor que b
    if (a > b) {
        print(a, "es mayor que", b)
    } else {
        print(a, "no es mayor que", b)
    }
}
```

### Ejercicio 3: Manipulación de Strings

Crea un programa que trabaje con cadenas de texto:

```r2
func main() {
    let frase = "R2Lang es un lenguaje moderno"
    
    print("Frase original:", frase)
    print("Longitud:", len(frase))
    print("En mayúsculas:", frase.upper())
    print("En minúsculas:", frase.lower())
    
    // Separar palabras
    let palabras = frase.split(" ")
    print("Palabras:", palabras)
    print("Primera palabra:", palabras[0])
    print("Última palabra:", palabras[palabras.length() - 1])
}
```

## Buenas Prácticas

### 1. Nombres Descriptivos

```r2
// ❌ Malo
let x = 25
let y = "Juan"

// ✅ Bueno
let edad = 25
let nombreUsuario = "Juan"
```

### 2. Comentarios Útiles

```r2
func main() {
    // Configuración inicial
    let precio = 100
    let descuento = 0.15  // 15% de descuento
    
    // Calcular precio final
    let precioFinal = precio * (1 - descuento)
    print("Precio final:", precioFinal)
}
```

### 3. Organización del Código

```r2
func main() {
    // 1. Declarar variables
    let base = 10
    let altura = 5
    
    // 2. Realizar cálculos
    let area = base * altura
    
    // 3. Mostrar resultados
    print("Área del rectángulo:", area)
}
```

## Errores Comunes

### 1. Variables No Declaradas

```r2
func main() {
    print(nombre)  // ❌ Error: variable no declarada
}

// Solución:
func main() {
    let nombre = "Juan"  // ✅ Declarar primero
    print(nombre)
}
```

### 2. Tipos Incompatibles

```r2
func main() {
    let numero = 5
    let texto = "10"
    
    // ❌ Puede causar comportamiento inesperado
    print(numero + texto)  // Concatenación: "510"
    
    // ✅ Mejor ser explícito
    print("Número:", numero, "Texto:", texto)
}
```

### 3. División por Cero

```r2
func main() {
    let a = 10
    let b = 0
    
    // ❌ Causará error en runtime
    print(a / b)
    
    // ✅ Verificar antes de dividir
    if (b != 0) {
        print("División:", a / b)
    } else {
        print("Error: no se puede dividir por cero")
    }
}
```

## Proyecto del Módulo

### Calculadora Personal

Crea un programa que funcione como una calculadora personal con las siguientes características:

```r2
func main() {
    // Información personal
    let nombre = "Tu Nombre"
    let salarioMensual = 2500
    let gastosFijos = 1200
    let ahorroMensual = salarioMensual - gastosFijos
    
    print("=== CALCULADORA PERSONAL ===")
    print("Usuario:", nombre)
    print()
    
    // Cálculos financieros
    print("FINANZAS MENSUALES:")
    print("Salario:", salarioMensual)
    print("Gastos fijos:", gastosFijos)
    print("Ahorro mensual:", ahorroMensual)
    print()
    
    // Proyecciones anuales
    let ahorroAnual = ahorroMensual * 12
    print("PROYECCIÓN ANUAL:")
    print("Ahorro anual:", ahorroAnual)
    
    // Análisis porcentual
    let porcentajeAhorro = (ahorroMensual / salarioMensual) * 100
    print("Porcentaje de ahorro:", porcentajeAhorro + "%")
    
    // Consejos basados en ahorro
    if (porcentajeAhorro >= 20) {
        print("¡Excelente! Tienes un buen nivel de ahorro.")
    } else if (porcentajeAhorro >= 10) {
        print("Buen trabajo, pero podrías ahorrar un poco más.")
    } else {
        print("Considera revisar tus gastos para ahorrar más.")
    }
}
```

## Resumen del Módulo

En este primer módulo has aprendido:

### Conceptos Clave
- ✅ Qué es R2Lang y sus características
- ✅ Instalación y configuración
- ✅ Variables y tipos de datos
- ✅ Operadores aritméticos, de comparación y lógicos
- ✅ Manipulación básica de strings
- ✅ Entrada y salida con print()

### Habilidades Desarrolladas
- ✅ Crear y ejecutar programas R2Lang básicos
- ✅ Declarar y usar variables de diferentes tipos
- ✅ Realizar operaciones matemáticas y lógicas
- ✅ Trabajar con cadenas de texto
- ✅ Escribir código limpio y bien comentado

### Próximo Módulo

En el **Módulo 2** aprenderás:
- Control de flujo (if/else, while, for)
- Funciones definidas por el usuario
- Arrays y su manipulación
- Manejo básico de errores

### Recursos Adicionales

- **Ejemplos**: Revisa la carpeta `examples/` en el repositorio
- **Documentación**: Lee los archivos en `docs/es/`
- **Práctica**: Experimenta modificando los ejemplos dados

¡Felicitaciones por completar el Módulo 1! Estás listo para continuar con estructuras de control más avanzadas.