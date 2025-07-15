# Curso R2Lang - Módulo 4: Concurrencia y Manejo de Errores

## Introducción

En este módulo aprenderás dos aspectos avanzados de R2Lang: la programación concurrente con goroutines y el manejo robusto de errores. Estos conceptos son fundamentales para crear aplicaciones robustas y eficientes.

### Mejoras en Concurrencia y Errores v2

La nueva arquitectura modular proporciona mejor soporte para concurrencia y manejo de errores:

```
Concurrency & Error Handling (pkg/r2libs/):
├── r2goroutine.go (237 LOC)    # Goroutines y concurrencia
├── r2std.go (122 LOC)          # Funciones estándar mejoradas
├── r2test.go                   # Testing de concurrencia
└── Error handling (pkg/r2core/):
    ├── try_statement.go        # Try-catch-finally
    ├── throw_statement.go      # Lanzamiento de errores
    └── error_handling.go       # Manejo robusto de errores
```

**Beneficios v2**:
- **Goroutines optimizadas**: Mejor rendimiento y gestión de memoria
- **Error handling robusto**: Stack traces y contexto detallado
- **Testing integrado**: Pruebas de concurrencia automáticas
- **Debugging mejorado**: Herramientas de diagnóstico avanzadas

## Concurrencia en R2Lang

### 1. Conceptos Básicos de Concurrencia

La concurrencia permite que múltiples tareas se ejecuten "al mismo tiempo" (en paralelo o entrelazadas). R2Lang utiliza goroutines, similares a las de Go, para manejar concurrencia.

#### Tu Primera Goroutine

```r2
func tarea() {
    print("Ejecutando tarea en goroutine")
    sleep(1)  // Simular trabajo que toma tiempo
    print("Tarea completada")
}

func main() {
    print("Iniciando programa")
    
    // Ejecutar función en goroutine
    r2(tarea)
    
    print("Continuando con otras operaciones")
    sleep(2)  // Esperar a que termine la goroutine
    print("Programa terminado")
}
```

#### Múltiples Goroutines

```r2
func trabajador(id) {
    print("Trabajador", id, "iniciado")
    
    for (let i = 1; i <= 3; i++) {
        print("Trabajador", id, "- tarea", i)
        sleep(1)
    }
    
    print("Trabajador", id, "terminado")
}

func main() {
    print("Creando trabajadores...")
    
    // Crear múltiples goroutines
    for (let i = 1; i <= 3; i++) {
        r2(trabajador, i)
    }
    
    print("Todos los trabajadores creados")
    sleep(4)  // Esperar a que terminen
    print("Programa principal terminado")
}
```

### 2. Patterns de Concurrencia

#### Worker Pool Pattern

```r2
func procesarDatos(datos, workerId) {
    print("Worker", workerId, "procesando:", datos)
    sleep(1)  // Simular procesamiento
    print("Worker", workerId, "completó:", datos)
}

func crearWorkerPool(numWorkers, tareas) {
    print("Creando pool de", numWorkers, "workers")
    
    for (let i = 0; i < numWorkers; i++) {
        let workerId = i + 1
        
        r2(func() {
            // Cada worker procesa una porción de las tareas
            let tareasPorWorker = tareas.length() / numWorkers
            let inicio = i * tareasPorWorker
            let fin = inicio + tareasPorWorker
            
            for (let j = inicio; j < fin && j < tareas.length(); j++) {
                procesarDatos(tareas[j], workerId)
            }
        })
    }
}

func main() {
    let tareas = ["Tarea-A", "Tarea-B", "Tarea-C", "Tarea-D", "Tarea-E", "Tarea-F"]
    
    print("Iniciando procesamiento paralelo")
    crearWorkerPool(3, tareas)
    
    sleep(4)  // Esperar a que terminen todos
    print("Procesamiento completado")
}
```

#### Producer-Consumer Pattern

```r2
func producer(nombreProductor, cantidad) {
    for (let i = 1; i <= cantidad; i++) {
        let producto = nombreProductor + "-Item-" + i
        print("📦 Producido:", producto)
        sleep(0.5)  // Simular tiempo de producción
    }
    print("✅ Productor", nombreProductor, "terminado")
}

func consumer(nombreConsumidor, tiempoTotal) {
    let startTime = 0  // Simulación de tiempo
    let tiempoLimite = tiempoTotal * 2  // 2 segundos por unidad
    
    while (startTime < tiempoLimite) {
        print("🛒 Consumidor", nombreConsumidor, "procesando items...")
        sleep(1)
        startTime++
    }
    print("✅ Consumidor", nombreConsumidor, "terminado")
}

func main() {
    print("=== PRODUCER-CONSUMER PATTERN ===")
    
    // Iniciar productores
    r2(producer, "P1", 3)
    r2(producer, "P2", 4)
    
    // Iniciar consumidores  
    r2(consumer, "C1", 3)
    r2(consumer, "C2", 4)
    
    sleep(6)
    print("Simulación terminada")
}
```

### 3. Concurrencia con Clases

```r2
class ContadorConcurrente {
    let valor
    let nombre
    
    constructor(nombre) {
        this.nombre = nombre
        this.valor = 0
    }
    
    incrementar(cantidad) {
        for (let i = 0; i < cantidad; i++) {
            this.valor++
            print(this.nombre, "incrementado a:", this.valor)
            sleep(0.1)  // Simular trabajo
        }
    }
    
    decrementar(cantidad) {
        for (let i = 0; i < cantidad; i++) {
            this.valor--
            print(this.nombre, "decrementado a:", this.valor)
            sleep(0.1)  // Simular trabajo
        }
    }
    
    obtenerValor() {
        return this.valor
    }
}

func main() {
    let contador = ContadorConcurrente("Contador-1")
    
    print("Valor inicial:", contador.obtenerValor())
    
    // Operaciones concurrentes
    r2(func() {
        contador.incrementar(5)
    })
    
    r2(func() {
        contador.decrementar(3)
    })
    
    r2(func() {
        contador.incrementar(2)
    })
    
    sleep(3)
    print("Valor final:", contador.obtenerValor())
}
```

### 4. Simulación de Sincronización

Aunque R2Lang no tiene primitivas de sincronización nativas, podemos simular algunos patrones:

```r2
// Simulación de Mutex usando flags
class MutexSimulado {
    let bloqueado
    let nombre
    
    constructor(nombre) {
        this.bloqueado = false
        this.nombre = nombre
    }
    
    lock() {
        while (this.bloqueado) {
            sleep(0.01)  // Esperar
        }
        this.bloqueado = true
        print("🔒 Lock adquirido por", this.nombre)
    }
    
    unlock() {
        this.bloqueado = false
        print("🔓 Lock liberado por", this.nombre)
    }
}

func trabajoConMutex(id, mutex, recursoCompartido) {
    print("Proceso", id, "intentando acceder al recurso")
    
    mutex.lock()
    
    try {
        print("Proceso", id, "usando recurso compartido")
        recursoCompartido.valor++
        print("Recurso actualizado a:", recursoCompartido.valor)
        sleep(1)  // Simular uso del recurso
    } finally {
        mutex.unlock()
        print("Proceso", id, "terminó de usar el recurso")
    }
}

func main() {
    let mutex = MutexSimulado("Mutex-Principal")
    let recurso = { valor: 0 }
    
    print("Iniciando acceso concurrente a recurso compartido")
    
    for (let i = 1; i <= 3; i++) {
        r2(trabajoConMutex, i, mutex, recurso)
    }
    
    sleep(5)
    print("Valor final del recurso:", recurso.valor)
}
```

## Manejo de Errores

### 1. Try-Catch-Finally Básico

```r2
func operacionRiesgosa(numero) {
    if (numero < 0) {
        throw "Número no puede ser negativo"
    }
    
    if (numero == 0) {
        throw "División por cero no permitida"
    }
    
    return 100 / numero
}

func main() {
    let numeros = [10, -5, 0, 20]
    
    for (let numero in numeros) {
        print("Procesando número:", numero)
        
        try {
            let resultado = operacionRiesgosa(numero)
            print("Resultado:", resultado)
        } catch (error) {
            print("Error capturado:", error)
        } finally {
            print("Operación completada para", numero)
        }
        print("---")
    }
}
```

### 2. Manejo de Errores en Funciones

```r2
func validarEdad(edad) {
    if (typeOf(edad) != "float64") {
        throw "Edad debe ser un número"
    }
    
    if (edad < 0) {
        throw "Edad no puede ser negativa"
    }
    
    if (edad > 150) {
        throw "Edad no puede ser mayor a 150"
    }
    
    return true
}

func crearPersona(nombre, edad) {
    try {
        // Validar entrada
        if (nombre == null || nombre == "") {
            throw "Nombre es requerido"
        }
        
        validarEdad(edad)
        
        // Crear persona si todo está bien
        let persona = {
            nombre: nombre,
            edad: edad,
            createdAt: "Ahora"
        }
        
        print("Persona creada:", persona.nombre)
        return persona
        
    } catch (error) {
        print("Error creando persona:", error)
        return null
    }
}

func main() {
    let datosPersonas = [
        ["Juan", 25],
        ["", 30],      // Error: nombre vacío
        ["Ana", -5],   // Error: edad negativa
        ["Carlos", "treinta"],  // Error: edad no numérica
        ["María", 200], // Error: edad muy alta
        ["Luis", 35]
    ]
    
    let personasValidas = []
    
    for (let datos in datosPersonas) {
        let nombre = datos[0]
        let edad = datos[1]
        
        let persona = crearPersona(nombre, edad)
        if (persona != null) {
            personasValidas = personasValidas.push(persona)
        }
    }
    
    print("\nPersonas válidas creadas:", personasValidas.length())
    for (let persona in personasValidas) {
        print("-", persona.nombre, "(" + persona.edad + " años)")
    }
}
```

### 3. Errores en Operaciones de Archivo

```r2
func leerArchivo(nombreArchivo) {
    try {
        print("Intentando leer archivo:", nombreArchivo)
        
        // Simular lectura de archivo
        if (nombreArchivo == "inexistente.txt") {
            throw "Archivo no encontrado"
        }
        
        if (nombreArchivo == "corrupto.txt") {
            throw "Archivo corrupto o dañado"
        }
        
        if (nombreArchivo == "permisos.txt") {
            throw "Sin permisos para leer el archivo"
        }
        
        // Simular contenido del archivo
        let contenido = "Contenido del archivo " + nombreArchivo
        print("Archivo leído exitosamente")
        return contenido
        
    } catch (error) {
        print("Error leyendo archivo:", error)
        throw "Error de archivo: " + error
    }
}

func procesarArchivos(archivos) {
    let procesados = 0
    let errores = 0
    
    for (let archivo in archivos) {
        try {
            let contenido = leerArchivo(archivo)
            print("Procesando contenido de", archivo)
            procesados++
            
        } catch (error) {
            print("No se pudo procesar", archivo + ":", error)
            errores++
            
        } finally {
            print("Finalizando procesamiento de", archivo)
        }
        print("---")
    }
    
    print("RESUMEN:")
    print("Archivos procesados:", procesados)
    print("Errores encontrados:", errores)
}

func main() {
    let archivos = [
        "documento1.txt",
        "inexistente.txt",
        "datos.txt",
        "corrupto.txt",
        "permisos.txt",
        "final.txt"
    ]
    
    procesarArchivos(archivos)
}
```

### 4. Manejo de Errores en Concurrencia

```r2
func tareaConErrores(id, shouldFail) {
    try {
        print("Tarea", id, "iniciada")
        
        if (shouldFail) {
            throw "Error simulado en tarea " + id
        }
        
        // Simular trabajo
        for (let i = 1; i <= 3; i++) {
            print("Tarea", id, "- paso", i)
            sleep(0.5)
        }
        
        print("Tarea", id, "completada exitosamente")
        
    } catch (error) {
        print("ERROR en tarea", id + ":", error)
        
    } finally {
        print("Tarea", id, "finalizando recursos")
    }
}

func supervisorTareas(numeroTareas) {
    let tareasExitosas = 0
    let tareasConError = 0
    
    print("Supervisor iniciando", numeroTareas, "tareas")
    
    for (let i = 1; i <= numeroTareas; i++) {
        // Algunas tareas fallarán (simulado)
        let shouldFail = (i % 3 == 0)  // Cada tercera tarea falla
        
        r2(func() {
            try {
                tareaConErrores(i, shouldFail)
                // No podemos actualizar contadores directamente debido a concurrencia
                print("✅ Tarea", i, "registrada como exitosa")
            } catch (error) {
                print("❌ Tarea", i, "registrada como fallida")
            }
        })
    }
    
    print("Todas las tareas lanzadas")
}

func main() {
    supervisorTareas(6)
    sleep(4)
    print("Supervisión completada")
}
```

## Patterns Avanzados

### 1. Circuit Breaker Pattern

```r2
class CircuitBreaker {
    let nombre
    let limiteErrores
    let erroresConsecutivos
    let estado  // "CERRADO", "ABIERTO", "SEMI_ABIERTO"
    let tiempoUltimoError
    
    constructor(nombre, limiteErrores) {
        this.nombre = nombre
        this.limiteErrores = limiteErrores
        this.erroresConsecutivos = 0
        this.estado = "CERRADO"
        this.tiempoUltimoError = 0
    }
    
    ejecutar(operacion) {
        if (this.estado == "ABIERTO") {
            print("Circuit Breaker ABIERTO - operación bloqueada")
            throw "Circuit breaker está abierto"
        }
        
        try {
            let resultado = operacion()
            this.onExito()
            return resultado
            
        } catch (error) {
            this.onError()
            throw error
        }
    }
    
    onExito() {
        this.erroresConsecutivos = 0
        if (this.estado == "SEMI_ABIERTO") {
            this.estado = "CERRADO"
            print("Circuit Breaker vuelve a CERRADO")
        }
    }
    
    onError() {
        this.erroresConsecutivos++
        this.tiempoUltimoError = 1  // Simulación de timestamp
        
        if (this.erroresConsecutivos >= this.limiteErrores) {
            this.estado = "ABIERTO"
            print("Circuit Breaker ABIERTO después de", this.erroresConsecutivos, "errores")
        }
    }
    
    intentarRecuperacion() {
        if (this.estado == "ABIERTO") {
            this.estado = "SEMI_ABIERTO"
            print("Circuit Breaker en modo SEMI_ABIERTO")
        }
    }
}

func operacionExterna(exito) {
    if (exito) {
        return "Operación exitosa"
    } else {
        throw "Operación falló"
    }
}

func main() {
    let cb = CircuitBreaker("API-Circuit-Breaker", 3)
    
    // Simular llamadas que fallan
    let intentos = [false, false, false, false, true]
    
    for (let i = 0; i < intentos.length(); i++) {
        let exito = intentos[i]
        
        try {
            let resultado = cb.ejecutar(func() {
                return operacionExterna(exito)
            })
            print("Resultado:", resultado)
            
        } catch (error) {
            print("Error:", error)
        }
        
        print("Estado actual:", cb.estado)
        print("---")
    }
    
    // Intentar recuperación
    print("Intentando recuperación...")
    cb.intentarRecuperacion()
    
    try {
        let resultado = cb.ejecutar(func() {
            return operacionExterna(true)
        })
        print("Resultado después de recuperación:", resultado)
    } catch (error) {
        print("Error en recuperación:", error)
    }
}
```

### 2. Retry Pattern

```r2
func retryOperacion(operacion, maxIntentos, delaySegundos) {
    let intentos = 0
    
    while (intentos < maxIntentos) {
        intentos++
        
        try {
            print("Intento", intentos, "de", maxIntentos)
            let resultado = operacion()
            print("Operación exitosa en intento", intentos)
            return resultado
            
        } catch (error) {
            print("Error en intento", intentos + ":", error)
            
            if (intentos >= maxIntentos) {
                print("Máximo número de intentos alcanzado")
                throw "Operación falló después de " + maxIntentos + " intentos"
            }
            
            print("Esperando", delaySegundos, "segundos antes del siguiente intento")
            sleep(delaySegundos)
        }
    }
}

func operacionInestable() {
    // Simular operación que falla aleatoriamente
    let random = math.random()  // Asumiendo que existe
    
    if (random > 0.7) {  // 30% de probabilidad de éxito
        return "Operación completada exitosamente"
    } else {
        throw "Fallo temporal de red"
    }
}

func main() {
    try {
        let resultado = retryOperacion(operacionInestable, 5, 1)
        print("Resultado final:", resultado)
        
    } catch (error) {
        print("Operación finalmente falló:", error)
    }
}
```

## Proyecto del Módulo: Sistema de Procesamiento Distribuido

```r2
// Simulación de un sistema de procesamiento distribuido
// con manejo de errores y concurrencia

class Nodo {
    let id
    let activo
    let carga
    let erroresConsecutivos
    
    constructor(id) {
        this.id = id
        this.activo = true
        this.carga = 0
        this.erroresConsecutivos = 0
    }
    
    procesar(tarea) {
        if (!this.activo) {
            throw "Nodo " + this.id + " está inactivo"
        }
        
        if (this.carga >= 5) {
            throw "Nodo " + this.id + " está sobrecargado"
        }
        
        this.carga++
        
        try {
            print("Nodo", this.id, "procesando tarea:", tarea.nombre)
            
            // Simular procesamiento
            sleep(tarea.duracion)
            
            // Simular posible error
            if (tarea.nombre.contains("error")) {
                throw "Error en tarea: " + tarea.nombre
            }
            
            let resultado = {
                nodo: this.id,
                tarea: tarea.nombre,
                resultado: "Procesado exitosamente",
                tiempo: tarea.duracion
            }
            
            this.erroresConsecutivos = 0
            print("✅ Nodo", this.id, "completó:", tarea.nombre)
            return resultado
            
        } catch (error) {
            this.erroresConsecutivos++
            print("❌ Error en nodo", this.id + ":", error)
            
            if (this.erroresConsecutivos >= 3) {
                this.activo = false
                print("🚫 Nodo", this.id, "desactivado por errores consecutivos")
            }
            
            throw error
            
        } finally {
            this.carga--
        }
    }
    
    reiniciar() {
        this.activo = true
        this.carga = 0
        this.erroresConsecutivos = 0
        print("🔄 Nodo", this.id, "reiniciado")
    }
}

class Coordinador {
    let nodos
    let tareasPendientes
    let tareasCompletadas
    let tareasConError
    
    constructor() {
        this.nodos = []
        this.tareasPendientes = []
        this.tareasCompletadas = []
        this.tareasConError = []
    }
    
    agregarNodo(nodo) {
        this.nodos = this.nodos.push(nodo)
        print("Nodo", nodo.id, "agregado al cluster")
    }
    
    agregarTarea(tarea) {
        this.tareasPendientes = this.tareasPendientes.push(tarea)
    }
    
    encontrarNodoDisponible() {
        for (let i = 0; i < this.nodos.length(); i++) {
            let nodo = this.nodos[i]
            if (nodo.activo && nodo.carga < 5) {
                return nodo
            }
        }
        return null
    }
    
    procesarTareas() {
        print("Iniciando procesamiento distribuido")
        
        for (let i = 0; i < this.tareasPendientes.length(); i++) {
            let tarea = this.tareasPendientes[i]
            
            r2(func() {
                let procesada = false
                let intentos = 0
                let maxIntentos = 3
                
                while (!procesada && intentos < maxIntentos) {
                    intentos++
                    
                    try {
                        let nodo = this.encontrarNodoDisponible()
                        
                        if (nodo == null) {
                            print("⏳ No hay nodos disponibles, esperando...")
                            sleep(1)
                            continue
                        }
                        
                        let resultado = nodo.procesar(tarea)
                        this.tareasCompletadas = this.tareasCompletadas.push(resultado)
                        procesada = true
                        
                    } catch (error) {
                        print("Error procesando", tarea.nombre + ":", error)
                        
                        if (intentos >= maxIntentos) {
                            this.tareasConError = this.tareasConError.push({
                                tarea: tarea,
                                error: error,
                                intentos: intentos
                            })
                            procesada = true  // Para salir del loop
                        } else {
                            sleep(1)  // Esperar antes de reintentar
                        }
                    }
                }
            })
        }
    }
    
    mostrarEstadisticas() {
        print("\n=== ESTADÍSTICAS DEL CLUSTER ===")
        print("Nodos totales:", this.nodos.length())
        
        let nodosActivos = 0
        for (let i = 0; i < this.nodos.length(); i++) {
            if (this.nodos[i].activo) {
                nodosActivos++
            }
        }
        
        print("Nodos activos:", nodosActivos)
        print("Tareas completadas:", this.tareasCompletadas.length())
        print("Tareas con error:", this.tareasConError.length())
        
        print("\nEstado de nodos:")
        for (let i = 0; i < this.nodos.length(); i++) {
            let nodo = this.nodos[i]
            let estado = nodo.activo ? "ACTIVO" : "INACTIVO"
            print("- Nodo", nodo.id + ":", estado, "- Carga:", nodo.carga)
        }
    }
    
    reiniciarNodosInactivos() {
        print("\nReiniciando nodos inactivos...")
        for (let i = 0; i < this.nodos.length(); i++) {
            let nodo = this.nodos[i]
            if (!nodo.activo) {
                nodo.reiniciar()
            }
        }
    }
}

func main() {
    // Crear coordinador
    let coordinador = Coordinador()
    
    // Crear nodos
    for (let i = 1; i <= 4; i++) {
        let nodo = Nodo("N" + i)
        coordinador.agregarNodo(nodo)
    }
    
    // Crear tareas (algunas con errores simulados)
    let tareas = [
        { nombre: "Tarea-1", duracion: 1 },
        { nombre: "Tarea-2", duracion: 2 },
        { nombre: "Tarea-error-1", duracion: 1 },  // Causará error
        { nombre: "Tarea-3", duracion: 1 },
        { nombre: "Tarea-4", duracion: 2 },
        { nombre: "Tarea-error-2", duracion: 1 },  // Causará error
        { nombre: "Tarea-5", duracion: 1 },
        { nombre: "Tarea-6", duracion: 1 }
    ]
    
    for (let tarea in tareas) {
        coordinador.agregarTarea(tarea)
    }
    
    // Procesar tareas
    coordinador.procesarTareas()
    
    // Esperar a que terminen
    sleep(8)
    
    // Mostrar estadísticas
    coordinador.mostrarEstadisticas()
    
    // Reiniciar nodos y procesar tareas fallidas
    coordinador.reiniciarNodosInactivos()
    
    sleep(2)
    coordinador.mostrarEstadisticas()
}
```

## Mejores Prácticas

### 1. Concurrencia
- ✅ Usar goroutines para tareas independientes
- ✅ Evitar shared state cuando sea posible
- ✅ Implementar timeouts para operaciones que pueden colgarse
- ✅ Usar patterns como worker pools para mejor control

### 2. Manejo de Errores
- ✅ Siempre usar finally para cleanup
- ✅ Ser específico en los mensajes de error
- ✅ Implementar retry logic para operaciones inestables
- ✅ Usar circuit breakers para servicios externos

### 3. Debugging
- ✅ Añadir logging detallado en operaciones concurrentes
- ✅ Usar IDs únicos para rastrear operaciones
- ✅ Implementar health checks para componentes críticos

## Resumen del Módulo

### Conceptos Aprendidos
- ✅ Programación concurrente con goroutines
- ✅ Patterns de concurrencia (Worker Pool, Producer-Consumer)
- ✅ Manejo robusto de errores con try-catch-finally
- ✅ Patterns de resilencia (Circuit Breaker, Retry)
- ✅ Debugging de aplicaciones concurrentes

### Habilidades Desarrolladas
- ✅ Diseñar sistemas concurrentes
- ✅ Implementar manejo de errores efectivo
- ✅ Crear aplicaciones resilientes
- ✅ Debugging de problemas de concurrencia
- ✅ Aplicar patterns de sistemas distribuidos

### Próximo Módulo

En el **Módulo 5** aprenderás:
- Sistema de testing integrado (BDD)
- Creación de APIs y servicios web
- Interacción con bases de datos
- Deployment y distribución

¡Excelente trabajo! Has dominado conceptos avanzados que te permitirán crear aplicaciones robustas y escalables.