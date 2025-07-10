# Curso R2Lang - Módulo 2: Control de Flujo y Funciones

## Introducción

En este módulo aprenderás a controlar el flujo de ejecución de tus programas y a crear funciones reutilizables. Estas son las herramientas fundamentales para crear programas más complejos y organizados.

## Control de Flujo

### 1. Condicionales (if/else)

#### Sintaxis Básica

```r2
func main() {
    let edad = 18
    
    if (edad >= 18) {
        print("Eres mayor de edad")
    } else {
        print("Eres menor de edad")
    }
}
```

#### Condicionales Múltiples (else if)

```r2
func main() {
    let nota = 85
    
    if (nota >= 90) {
        print("Excelente!")
    } else if (nota >= 80) {
        print("Muy bien!")
    } else if (nota >= 70) {
        print("Bien")
    } else if (nota >= 60) {
        print("Suficiente")
    } else {
        print("Necesitas estudiar más")
    }
}
```

#### Operadores Lógicos en Condicionales

```r2
func main() {
    let edad = 25
    let tienePermiso = true
    let tieneExperiencia = false
    
    // AND (&&)
    if (edad >= 18 && tienePermiso) {
        print("Puede conducir")
    }
    
    // OR (||)
    if (tienePermiso || tieneExperiencia) {
        print("Tiene alguna calificación")
    }
    
    // NOT (!)
    if (!tieneExperiencia) {
        print("Necesita práctica")
    }
    
    // Combinaciones complejas
    if ((edad >= 21 && tienePermiso) || tieneExperiencia) {
        print("Puede conducir vehículos comerciales")
    }
}
```

### 2. Bucles (Loops)

#### While Loop

```r2
func main() {
    let contador = 1
    
    while (contador <= 5) {
        print("Contador:", contador)
        contador++  // Incrementar en 1
    }
    
    print("Bucle terminado")
}
```

#### Ejemplo Práctico: Suma de Números

```r2
func main() {
    let numero = 1
    let suma = 0
    
    while (numero <= 10) {
        suma = suma + numero
        print("Sumando", numero, "- Total:", suma)
        numero++
    }
    
    print("Suma total del 1 al 10:", suma)
}
```

#### For Loop Tradicional

```r2
func main() {
    // Sintaxis: for (inicialización; condición; incremento)
    for (let i = 1; i <= 5; i++) {
        print("Iteración", i)
    }
    
    // Countdown
    for (let i = 10; i >= 1; i--) {
        print("Cuenta regresiva:", i)
    }
    print("¡Despegue!")
}
```

#### For Loop con Arrays (for-in)

```r2
func main() {
    let frutas = ["manzana", "banana", "naranja", "uva"]
    
    // Iterar sobre elementos
    for (let fruta in frutas) {
        print("Fruta:", fruta)
    }
    
    // Iterar con índices
    for (let i = 0; i < frutas.length(); i++) {
        print("Posición", i, ":", frutas[i])
    }
}
```

### 3. Control de Bucles

#### Break - Salir del Bucle

```r2
func main() {
    let numero = 1
    
    while (true) {  // Bucle infinito
        if (numero > 5) {
            break  // Salir del bucle
        }
        print("Número:", numero)
        numero++
    }
    print("Bucle terminado con break")
}
```

#### Continue - Saltar Iteración

```r2
func main() {
    for (let i = 1; i <= 10; i++) {
        if (i % 2 == 0) {  // Si es par
            continue  // Saltar al siguiente
        }
        print("Número impar:", i)
    }
}
```

## Funciones

### 1. Definición y Llamada de Funciones

#### Función Simple

```r2
func saludar() {
    print("¡Hola desde una función!")
}

func main() {
    saludar()  // Llamar la función
    saludar()  // Llamar nuevamente
}
```

#### Funciones con Parámetros

```r2
func saludarPersona(nombre) {
    print("¡Hola", nombre + "!")
}

func sumar(a, b) {
    let resultado = a + b
    print(a, "+", b, "=", resultado)
}

func main() {
    saludarPersona("Ana")
    saludarPersona("Carlos")
    
    sumar(5, 3)
    sumar(10, 25)
}
```

#### Funciones con Valor de Retorno

```r2
func multiplicar(a, b) {
    return a * b
}

func esMayorDeEdad(edad) {
    return edad >= 18
}

func obtenerMensaje(nombre, edad) {
    if (esMayorDeEdad(edad)) {
        return nombre + " es mayor de edad"
    } else {
        return nombre + " es menor de edad"
    }
}

func main() {
    let resultado = multiplicar(6, 7)
    print("6 × 7 =", resultado)
    
    let mensaje = obtenerMensaje("Laura", 22)
    print(mensaje)
    
    // Usar función directamente en condicional
    if (esMayorDeEdad(16)) {
        print("Puede votar")
    } else {
        print("No puede votar")
    }
}
```

### 2. Scope (Alcance) de Variables

#### Variables Locales vs Globales

```r2
let variableGlobal = "Soy global"

func miFuncion() {
    let variableLocal = "Soy local"
    print("Dentro de función:")
    print("- Global:", variableGlobal)
    print("- Local:", variableLocal)
}

func main() {
    print("En main:")
    print("- Global:", variableGlobal)
    // print("- Local:", variableLocal)  // ❌ Error: no existe aquí
    
    miFuncion()
}
```

#### Parámetros son Locales

```r2
func modificarParametro(numero) {
    numero = numero + 10
    print("Dentro de función:", numero)
    return numero
}

func main() {
    let miNumero = 5
    print("Antes:", miNumero)
    
    let nuevoNumero = modificarParametro(miNumero)
    print("Después:", miNumero)      // Sigue siendo 5
    print("Retornado:", nuevoNumero) // Es 15
}
```

### 3. Funciones Avanzadas

#### Funciones como Variables

```r2
func sumar(a, b) {
    return a + b
}

func restar(a, b) {
    return a - b
}

func main() {
    // Asignar función a variable
    let operacion = sumar
    let resultado = operacion(10, 5)
    print("Resultado:", resultado)
    
    // Cambiar la operación
    operacion = restar
    resultado = operacion(10, 5)
    print("Resultado:", resultado)
}
```

#### Funciones Anónimas (Lambda)

```r2
func main() {
    // Función anónima asignada a variable
    let duplicar = func(x) {
        return x * 2
    }
    
    print("Duplicar 7:", duplicar(7))
    
    // Función anónima directa
    let resultado = func(a, b) {
        return a * b + 10
    }(5, 3)
    
    print("Resultado:", resultado)  // (5*3)+10 = 25
}
```

## Arrays y Colecciones

### 1. Declaración y Acceso

```r2
func main() {
    // Crear array vacío
    let numeros = []
    
    // Crear array con elementos
    let frutas = ["manzana", "banana", "naranja"]
    let edades = [25, 30, 18, 45]
    
    // Acceder elementos
    print("Primera fruta:", frutas[0])
    print("Última edad:", edades[edades.length() - 1])
    
    // Modificar elementos
    frutas[1] = "fresa"
    print("Frutas modificadas:", frutas)
}
```

### 2. Métodos de Arrays

```r2
func main() {
    let numeros = [1, 2, 3]
    
    // Agregar elementos
    numeros = numeros.push(4)
    numeros = numeros.push(5)
    print("Después de push:", numeros)
    
    // Longitud
    print("Longitud:", numeros.length())
    
    // Buscar elemento
    let posicion = numeros.find(3)
    print("Posición del 3:", posicion)
    
    // Verificar si contiene
    let contiene = numeros.find(10)
    if (contiene != null) {
        print("Contiene 10")
    } else {
        print("No contiene 10")
    }
}
```

### 3. Iteración sobre Arrays

```r2
func imprimirArray(arr, nombre) {
    print("=== " + nombre + " ===")
    for (let i = 0; i < arr.length(); i++) {
        print("Posición", i, ":", arr[i])
    }
}

func main() {
    let colores = ["rojo", "verde", "azul"]
    let numeros = [10, 20, 30, 40]
    
    imprimirArray(colores, "COLORES")
    imprimirArray(numeros, "NÚMEROS")
}
```

## Ejercicios Prácticos

### Ejercicio 1: Calculadora de Promedio

```r2
func calcularPromedio(numeros) {
    let suma = 0
    let cantidad = numeros.length()
    
    for (let i = 0; i < cantidad; i++) {
        suma = suma + numeros[i]
    }
    
    return suma / cantidad
}

func main() {
    let notas = [85, 92, 78, 90, 88]
    let promedio = calcularPromedio(notas)
    
    print("Notas:", notas)
    print("Promedio:", promedio)
    
    if (promedio >= 90) {
        print("Calificación: Excelente")
    } else if (promedio >= 80) {
        print("Calificación: Muy Bien")
    } else if (promedio >= 70) {
        print("Calificación: Bien")
    } else {
        print("Calificación: Necesita Mejorar")
    }
}
```

### Ejercicio 2: Números Pares e Impares

```r2
func clasificarNumeros(numeros) {
    let pares = []
    let impares = []
    
    for (let numero in numeros) {
        if (numero % 2 == 0) {
            pares = pares.push(numero)
        } else {
            impares = impares.push(numero)
        }
    }
    
    print("Números originales:", numeros)
    print("Números pares:", pares)
    print("Números impares:", impares)
}

func main() {
    let numeros = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
    clasificarNumeros(numeros)
}
```

### Ejercicio 3: Buscador de Palabras

```r2
func buscarPalabra(palabras, palabraBuscada) {
    let encontradas = []
    
    for (let i = 0; i < palabras.length(); i++) {
        let palabra = palabras[i]
        if (palabra.contains(palabraBuscada)) {
            encontradas = encontradas.push(palabra)
        }
    }
    
    return encontradas
}

func main() {
    let diccionario = ["programación", "programa", "código", "desarrollo", "programador"]
    let buscar = "programa"
    
    let resultados = buscarPalabra(diccionario, buscar)
    
    print("Buscando palabras que contengan:", buscar)
    print("Resultados encontrados:", resultados)
    print("Total encontradas:", resultados.length())
}
```

## Manejo Básico de Errores

### 1. Validación de Parámetros

```r2
func dividir(a, b) {
    if (b == 0) {
        print("Error: No se puede dividir por cero")
        return null
    }
    return a / b
}

func main() {
    let resultado1 = dividir(10, 2)
    let resultado2 = dividir(10, 0)
    
    if (resultado1 != null) {
        print("10 ÷ 2 =", resultado1)
    }
    
    if (resultado2 != null) {
        print("10 ÷ 0 =", resultado2)
    } else {
        print("División por cero no es válida")
    }
}
```

### 2. Validación de Arrays

```r2
func obtenerElemento(array, indice) {
    if (array.length() == 0) {
        print("Error: Array está vacío")
        return null
    }
    
    if (indice < 0 || indice >= array.length()) {
        print("Error: Índice fuera de rango")
        return null
    }
    
    return array[indice]
}

func main() {
    let numeros = [10, 20, 30]
    let vacio = []
    
    print("Elemento válido:", obtenerElemento(numeros, 1))
    print("Índice inválido:", obtenerElemento(numeros, 5))
    print("Array vacío:", obtenerElemento(vacio, 0))
}
```

## Proyecto del Módulo: Sistema de Gestión de Estudiantes

```r2
// Sistema simple de gestión de estudiantes

func crearEstudiante(nombre, edad, notas) {
    return {
        nombre: nombre,
        edad: edad,
        notas: notas
    }
}

func calcularPromedio(notas) {
    if (notas.length() == 0) {
        return 0
    }
    
    let suma = 0
    for (let nota in notas) {
        suma = suma + nota
    }
    
    return suma / notas.length()
}

func obtenerCalificacion(promedio) {
    if (promedio >= 90) {
        return "A"
    } else if (promedio >= 80) {
        return "B"
    } else if (promedio >= 70) {
        return "C"
    } else if (promedio >= 60) {
        return "D"
    } else {
        return "F"
    }
}

func mostrarEstudiante(estudiante) {
    let promedio = calcularPromedio(estudiante.notas)
    let calificacion = obtenerCalificacion(promedio)
    
    print("=== ESTUDIANTE ===")
    print("Nombre:", estudiante.nombre)
    print("Edad:", estudiante.edad)
    print("Notas:", estudiante.notas)
    print("Promedio:", promedio)
    print("Calificación:", calificacion)
    print()
}

func encontrarMejorEstudiante(estudiantes) {
    if (estudiantes.length() == 0) {
        return null
    }
    
    let mejor = estudiantes[0]
    let mejorPromedio = calcularPromedio(mejor.notas)
    
    for (let i = 1; i < estudiantes.length(); i++) {
        let actual = estudiantes[i]
        let promedioActual = calcularPromedio(actual.notas)
        
        if (promedioActual > mejorPromedio) {
            mejor = actual
            mejorPromedio = promedioActual
        }
    }
    
    return mejor
}

func main() {
    // Crear estudiantes
    let estudiante1 = crearEstudiante("Ana García", 20, [85, 92, 88, 90])
    let estudiante2 = crearEstudiante("Carlos López", 19, [78, 85, 82, 89])
    let estudiante3 = crearEstudiante("María Rodríguez", 21, [95, 98, 92, 96])
    
    let estudiantes = [estudiante1, estudiante2, estudiante3]
    
    // Mostrar todos los estudiantes
    print("REPORTE DE ESTUDIANTES")
    print("======================")
    
    for (let estudiante in estudiantes) {
        mostrarEstudiante(estudiante)
    }
    
    // Encontrar el mejor estudiante
    let mejor = encontrarMejorEstudiante(estudiantes)
    if (mejor != null) {
        print("🏆 MEJOR ESTUDIANTE:")
        mostrarEstudiante(mejor)
    }
    
    // Estadísticas generales
    let totalEstudiantes = estudiantes.length()
    let sumaPromedios = 0
    
    for (let estudiante in estudiantes) {
        sumaPromedios = sumaPromedios + calcularPromedio(estudiante.notas)
    }
    
    let promedioGeneral = sumaPromedios / totalEstudiantes
    
    print("ESTADÍSTICAS GENERALES:")
    print("Total de estudiantes:", totalEstudiantes)
    print("Promedio general:", promedioGeneral)
    print("Calificación general:", obtenerCalificacion(promedioGeneral))
}
```

## Patrones y Buenas Prácticas

### 1. Funciones Pequeñas y Específicas

```r2
// ❌ Función que hace demasiado
func procesarDatos(datos) {
    // Validar, procesar, calcular, formatear, imprimir...
    // 50 líneas de código
}

// ✅ Funciones específicas
func validarDatos(datos) {
    return datos != null && datos.length() > 0
}

func calcularEstadisticas(datos) {
    // Solo calcular
}

func formatearResultados(estadisticas) {
    // Solo formatear
}
```

### 2. Nombres Descriptivos

```r2
// ❌ Nombres poco claros
func calc(x, y) {
    return x * y * 0.15
}

// ✅ Nombres descriptivos
func calcularImpuesto(precio, cantidad) {
    let impuesto = 0.15
    return precio * cantidad * impuesto
}
```

### 3. Validación de Entrada

```r2
func crearPersona(nombre, edad) {
    // Validar parámetros
    if (nombre == null || nombre == "") {
        print("Error: Nombre es requerido")
        return null
    }
    
    if (edad < 0 || edad > 150) {
        print("Error: Edad inválida")
        return null
    }
    
    return {
        nombre: nombre,
        edad: edad
    }
}
```

## Resumen del Módulo

### Conceptos Aprendidos
- ✅ Condicionales (if/else/else if)
- ✅ Bucles (while, for, for-in)
- ✅ Control de bucles (break, continue)
- ✅ Definición y uso de funciones
- ✅ Parámetros y valores de retorno
- ✅ Scope de variables
- ✅ Arrays y sus métodos básicos
- ✅ Funciones anónimas
- ✅ Manejo básico de errores

### Habilidades Desarrolladas
- ✅ Crear programas con lógica condicional
- ✅ Implementar bucles eficientes
- ✅ Escribir funciones reutilizables
- ✅ Trabajar con colecciones de datos
- ✅ Validar entrada y manejar errores
- ✅ Organizar código en funciones pequeñas

### Próximo Módulo

En el **Módulo 3** aprenderás:
- Orientación a objetos (clases y objetos)
- Herencia y polimorfismo
- Maps y objetos avanzados
- Manejo de archivos básico

¡Excelente trabajo completando el Módulo 2! Ya tienes las herramientas fundamentales para crear programas estructurados y funcionales.