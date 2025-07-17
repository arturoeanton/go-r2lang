# Curso R2Lang - Módulo 8: Concurrencia y Programación Paralela

## Introducción

La concurrencia es uno de los aspectos más potentes de R2Lang, permitiendo ejecutar múltiples tareas simultáneamente. En este módulo explorarás el sistema de concurrencia basado en goroutines, patrones de diseño concurrente, y cómo crear aplicaciones paralelas eficientes.

### Arquitectura de Concurrencia v2

```
pkg/r2libs/r2goroutine.go (237 LOC):
├── r2() function           # Crear goroutines
├── Channel implementation  # Comunicación entre goroutines
├── Mutex/RWMutex support  # Sincronización
├── WaitGroup management   # Espera de goroutines
├── Atomic operations      # Operaciones atómicas
└── Worker pool patterns   # Pools de workers
```

El sistema de concurrencia está completamente integrado con el runtime de Go, proporcionando verdadero paralelismo.

## Fundamentos de Concurrencia

### 1. Goroutines Básicas

```r2
func tareaSimple(nombre, duracion) {
    for (let i = 1; i <= 5; i++) {
        print(nombre + " - paso " + i)
        sleep(duracion)
    }
    print(nombre + " completada")
}

func main() {
    print("Iniciando programa concurrente")
    
    // Ejecutar tareas secuencialmente
    print("=== SECUENCIAL ===")
    tareaSimple("Tarea A", 1000)
    tareaSimple("Tarea B", 1000)
    
    print()
    print("=== CONCURRENTE ===")
    
    // Ejecutar tareas concurrentemente
    r2(tareaSimple, "Tarea C", 1000)
    r2(tareaSimple, "Tarea D", 1000)
    r2(tareaSimple, "Tarea E", 1000)
    
    // Esperar a que terminen todas las goroutines
    print("Esperando que terminen todas las tareas...")
    sleep(6000)  // Simular espera
    
    print("Programa terminado")
}
```

### 2. Goroutines con Valores de Retorno

```r2
func calcularCuadrado(numero, resultado) {
    let cuadrado = numero * numero
    print("Cuadrado de " + numero + " es " + cuadrado)
    
    // Simular resultado compartido
    resultado.valor = cuadrado
    resultado.completado = true
}

func calcularFactorial(numero, resultado) {
    let factorial = 1
    for (let i = 1; i <= numero; i++) {
        factorial = factorial * i
    }
    print("Factorial de " + numero + " es " + factorial)
    
    resultado.valor = factorial
    resultado.completado = true
}

func main() {
    print("Calculando valores concurrentemente...")
    
    // Objetos para almacenar resultados
    let resultadoCuadrado = {valor: 0, completado: false}
    let resultadoFactorial = {valor: 0, completado: false}
    
    // Ejecutar cálculos concurrentemente
    r2(calcularCuadrado, 10, resultadoCuadrado)
    r2(calcularFactorial, 8, resultadoFactorial)
    
    // Esperar resultados
    while (!resultadoCuadrado.completado || !resultadoFactorial.completado) {
        sleep(100)
    }
    
    print("Resultados:")
    print("Cuadrado: " + resultadoCuadrado.valor)
    print("Factorial: " + resultadoFactorial.valor)
}
```

### 3. Pool de Workers

```r2
class WorkerPool {
    let numWorkers
    let tareas
    let activo
    
    constructor(numWorkers) {
        this.numWorkers = numWorkers
        this.tareas = []
        this.activo = false
    }
    
    agregarTarea(tarea) {
        this.tareas = this.tareas.push(tarea)
        print("Tarea agregada. Total en cola: " + this.tareas.length())
    }
    
    worker(id) {
        print("Worker " + id + " iniciado")
        
        while (this.activo) {
            if (this.tareas.length() > 0) {
                // Tomar primera tarea
                let tarea = this.tareas[0]
                this.tareas = this.tareas.slice(1)
                
                print("Worker " + id + " procesando tarea: " + tarea.nombre)
                
                // Procesar tarea
                let resultado = tarea.funcion(tarea.datos)
                print("Worker " + id + " completó: " + tarea.nombre + " -> " + resultado)
                
                // Simular tiempo de procesamiento
                sleep(tarea.duracion)
            } else {
                // No hay tareas, esperar
                sleep(100)
            }
        }
        
        print("Worker " + id + " terminado")
    }
    
    iniciar() {
        this.activo = true
        print("Iniciando pool con " + this.numWorkers + " workers")
        
        // Crear workers
        for (let i = 1; i <= this.numWorkers; i++) {
            r2(this.worker, i)
        }
    }
    
    detener() {
        this.activo = false
        print("Deteniendo worker pool...")
    }
}

func procesarNumero(numero) {
    return numero * numero
}

func procesarTexto(texto) {
    return texto.upper()
}

func main() {
    let pool = WorkerPool(3)
    
    // Agregar tareas
    pool.agregarTarea({
        nombre: "Cuadrado de 5",
        funcion: procesarNumero,
        datos: 5,
        duracion: 1000
    })
    
    pool.agregarTarea({
        nombre: "Cuadrado de 10",
        funcion: procesarNumero,
        datos: 10,
        duracion: 1500
    })
    
    pool.agregarTarea({
        nombre: "Mayúsculas",
        funcion: procesarTexto,
        datos: "hola mundo",
        duracion: 800
    })
    
    pool.agregarTarea({
        nombre: "Cuadrado de 7",
        funcion: procesarNumero,
        datos: 7,
        duracion: 1200
    })
    
    pool.agregarTarea({
        nombre: "Mayúsculas 2",
        funcion: procesarTexto,
        datos: "r2lang concurrente",
        duracion: 600
    })
    
    // Iniciar pool
    pool.iniciar()
    
    // Ejecutar por un tiempo
    sleep(8000)
    
    // Detener pool
    pool.detener()
    
    print("Procesamiento terminado")
}
```

## Patrones de Concurrencia

### 1. Productor-Consumidor

```r2
class ProducerConsumer {
    let buffer
    let maxSize
    let produciendo
    let consumiendo
    
    constructor(maxSize) {
        this.buffer = []
        this.maxSize = maxSize
        this.produciendo = true
        this.consumiendo = true
    }
    
    producer(id) {
        let contador = 1
        
        while (this.produciendo) {
            // Esperar si buffer está lleno
            while (this.buffer.length() >= this.maxSize) {
                sleep(100)
            }
            
            let item = "Item-" + id + "-" + contador
            this.buffer = this.buffer.push(item)
            print("Productor " + id + " creó: " + item + " (Buffer: " + this.buffer.length() + ")")
            
            contador++
            sleep(rand.int(500, 1500))  // Tiempo variable de producción
        }
        
        print("Productor " + id + " terminado")
    }
    
    consumer(id) {
        while (this.consumiendo) {
            // Esperar si buffer está vacío
            while (this.buffer.length() == 0 && this.consumiendo) {
                sleep(100)
            }
            
            if (this.buffer.length() > 0) {
                let item = this.buffer[0]
                this.buffer = this.buffer.slice(1)
                print("Consumidor " + id + " procesó: " + item + " (Buffer: " + this.buffer.length() + ")")
                
                // Simular procesamiento
                sleep(rand.int(800, 2000))
            }
        }
        
        print("Consumidor " + id + " terminado")
    }
    
    iniciar(numProductores, numConsumidores) {
        print("Iniciando " + numProductores + " productores y " + numConsumidores + " consumidores")
        print("Tamaño de buffer: " + this.maxSize)
        
        // Crear productores
        for (let i = 1; i <= numProductores; i++) {
            r2(this.producer, i)
        }
        
        // Crear consumidores
        for (let i = 1; i <= numConsumidores; i++) {
            r2(this.consumer, i)
        }
    }
    
    detener() {
        this.produciendo = false
        sleep(1000)
        this.consumiendo = false
        print("Sistema producer-consumer detenido")
    }
}

func main() {
    let sistema = ProducerConsumer(5)
    
    // Iniciar con 2 productores y 3 consumidores
    sistema.iniciar(2, 3)
    
    // Ejecutar por 15 segundos
    sleep(15000)
    
    // Detener sistema
    sistema.detener()
}
```

### 2. Map-Reduce

```r2
class MapReduce {
    let mappers
    let reducers
    let intermediate
    let resultado
    
    constructor(numMappers, numReducers) {
        this.mappers = numMappers
        this.reducers = numReducers
        this.intermediate = {}
        this.resultado = {}
    }
    
    // Función map: procesar cada elemento
    mapFunction(datos, mapperId) {
        print("Mapper " + mapperId + " procesando " + datos.length() + " elementos")
        
        for (let i = 0; i < datos.length(); i++) {
            let item = datos[i]
            
            // Contar palabras (ejemplo)
            let palabras = item.split(" ")
            for (let j = 0; j < palabras.length(); j++) {
                let palabra = palabras[j].lower()
                
                // Agregar a resultados intermedios
                if (this.intermediate[palabra] == null) {
                    this.intermediate[palabra] = []
                }
                this.intermediate[palabra] = this.intermediate[palabra].push(1)
            }
        }
        
        print("Mapper " + mapperId + " completado")
    }
    
    // Función reduce: combinar resultados
    reduceFunction(palabra, valores, reducerId) {
        print("Reducer " + reducerId + " procesando palabra: " + palabra)
        
        let suma = 0
        for (let i = 0; i < valores.length(); i++) {
            suma = suma + valores[i]
        }
        
        this.resultado[palabra] = suma
        print("Reducer " + reducerId + " completado: " + palabra + " = " + suma)
    }
    
    procesar(datos) {
        print("Iniciando MapReduce con " + this.mappers + " mappers y " + this.reducers + " reducers")
        
        // Dividir datos entre mappers
        let chunkSize = math.ceil(datos.length() / this.mappers)
        
        for (let i = 0; i < this.mappers; i++) {
            let inicio = i * chunkSize
            let fin = math.min((i + 1) * chunkSize, datos.length())
            let chunk = datos.slice(inicio, fin)
            
            r2(this.mapFunction, chunk, i + 1)
        }
        
        // Esperar que terminen todos los mappers
        sleep(3000)
        
        // Fase reduce
        let palabras = Object.keys(this.intermediate)
        let palabrasPorReducer = math.ceil(palabras.length() / this.reducers)
        
        for (let i = 0; i < this.reducers; i++) {
            let inicio = i * palabrasPorReducer
            let fin = math.min((i + 1) * palabrasPorReducer, palabras.length())
            
            for (let j = inicio; j < fin; j++) {
                let palabra = palabras[j]
                let valores = this.intermediate[palabra]
                r2(this.reduceFunction, palabra, valores, i + 1)
            }
        }
        
        // Esperar que terminen todos los reducers
        sleep(2000)
        
        return this.resultado
    }
}

func main() {
    let textos = [
        "hola mundo desde r2lang",
        "r2lang es un lenguaje moderno",
        "programación concurrente en r2lang",
        "hola mundo de la programación",
        "r2lang soporta concurrencia nativa",
        "mundo de la programación concurrente"
    ]
    
    let mapReduce = MapReduce(2, 2)
    let resultado = mapReduce.procesar(textos)
    
    print("=== RESULTADOS ===")
    print("Conteo de palabras:")
    
    // Simular mostrar resultados
    print("hola: 2")
    print("mundo: 3")
    print("r2lang: 3")
    print("programación: 3")
    print("concurrente: 2")
    print("es: 1")
    print("un: 1")
    print("lenguaje: 1")
    print("moderno: 1")
    print("en: 1")
    print("desde: 1")
    print("soporta: 1")
    print("concurrencia: 1")
    print("nativa: 1")
    print("de: 1")
    print("la: 1")
}
```

## Aplicaciones Concurrentes Avanzadas

### 1. Crawler Web Concurrente

```r2
class WebCrawler {
    let maxWorkers
    let visitados
    let cola
    let activo
    let resultados
    
    constructor(maxWorkers) {
        this.maxWorkers = maxWorkers
        this.visitados = {}
        this.cola = []
        this.activo = true
        this.resultados = []
    }
    
    worker(id) {
        print("Crawler worker " + id + " iniciado")
        
        while (this.activo) {
            if (this.cola.length() > 0) {
                // Tomar URL de la cola
                let url = this.cola[0]
                this.cola = this.cola.slice(1)
                
                if (this.visitados[url] == null) {
                    this.visitados[url] = true
                    this.crawlPage(url, id)
                }
            } else {
                // No hay URLs, esperar
                sleep(200)
            }
        }
        
        print("Crawler worker " + id + " terminado")
    }
    
    crawlPage(url, workerId) {
        print("Worker " + workerId + " crawling: " + url)
        
        // Simular descarga de página
        sleep(rand.int(1000, 3000))
        
        // Simular extracción de contenido
        let contenido = {
            url: url,
            titulo: "Título de " + url,
            links: [],
            contenido: "Contenido simulado de " + url,
            timestamp: os.time()
        }
        
        // Simular encontrar links
        let numLinks = rand.int(2, 5)
        for (let i = 0; i < numLinks; i++) {
            let link = url + "/page" + i
            contenido.links = contenido.links.push(link)
            
            // Agregar a cola si no está visitado
            if (this.visitados[link] == null) {
                this.cola = this.cola.push(link)
            }
        }
        
        this.resultados = this.resultados.push(contenido)
        print("Worker " + workerId + " completó: " + url + " (" + contenido.links.length() + " links encontrados)")
    }
    
    iniciar(urlInicial) {
        print("Iniciando crawler con " + this.maxWorkers + " workers")
        
        // Agregar URL inicial
        this.cola = this.cola.push(urlInicial)
        
        // Crear workers
        for (let i = 1; i <= this.maxWorkers; i++) {
            r2(this.worker, i)
        }
    }
    
    detener() {
        this.activo = false
        print("Deteniendo crawler...")
    }
    
    getEstadisticas() {
        return {
            urlsVisitadas: Object.keys(this.visitados).length(),
            urlsEnCola: this.cola.length(),
            paginasCrawled: this.resultados.length()
        }
    }
}

func main() {
    let crawler = WebCrawler(3)
    
    // Iniciar crawler
    crawler.iniciar("https://ejemplo.com")
    
    // Monitorear progreso
    for (let i = 0; i < 10; i++) {
        sleep(2000)
        let stats = crawler.getEstadisticas()
        print("Progreso: " + stats.paginasCrawled + " páginas, " + stats.urlsEnCola + " en cola")
    }
    
    // Detener crawler
    crawler.detener()
    
    let statsFinales = crawler.getEstadisticas()
    print("=== ESTADÍSTICAS FINALES ===")
    print("URLs visitadas: " + statsFinales.urlsVisitadas)
    print("Páginas crawled: " + statsFinales.paginasCrawled)
}
```

### 2. Sistema de Cache Concurrente

```r2
class CacheManager {
    let cache
    let maxSize
    let expiracion
    let lectores
    let escritores
    
    constructor(maxSize, expiracionMs) {
        this.cache = {}
        this.maxSize = maxSize
        this.expiracion = expiracionMs
        this.lectores = 0
        this.escritores = 0
    }
    
    // Simular read lock
    readLock() {
        while (this.escritores > 0) {
            sleep(10)
        }
        this.lectores++
    }
    
    readUnlock() {
        this.lectores--
    }
    
    // Simular write lock
    writeLock() {
        while (this.lectores > 0 || this.escritores > 0) {
            sleep(10)
        }
        this.escritores++
    }
    
    writeUnlock() {
        this.escritores--
    }
    
    get(key) {
        this.readLock()
        
        let item = this.cache[key]
        if (item != null) {
            let ahora = os.time()
            if (ahora - item.timestamp < this.expiracion) {
                this.readUnlock()
                return item.valor
            } else {
                this.readUnlock()
                // Eliminar item expirado
                this.delete(key)
                return null
            }
        }
        
        this.readUnlock()
        return null
    }
    
    set(key, valor) {
        this.writeLock()
        
        // Verificar si cache está lleno
        if (Object.keys(this.cache).length() >= this.maxSize) {
            this.evictOldest()
        }
        
        this.cache[key] = {
            valor: valor,
            timestamp: os.time()
        }
        
        this.writeUnlock()
    }
    
    delete(key) {
        this.writeLock()
        delete this.cache[key]
        this.writeUnlock()
    }
    
    evictOldest() {
        let keys = Object.keys(this.cache)
        if (keys.length() > 0) {
            let oldest = keys[0]
            let oldestTime = this.cache[oldest].timestamp
            
            for (let i = 1; i < keys.length(); i++) {
                let key = keys[i]
                if (this.cache[key].timestamp < oldestTime) {
                    oldest = key
                    oldestTime = this.cache[key].timestamp
                }
            }
            
            delete this.cache[oldest]
            print("Evicted oldest item: " + oldest)
        }
    }
    
    stats() {
        this.readLock()
        let stats = {
            size: Object.keys(this.cache).length(),
            maxSize: this.maxSize,
            lectores: this.lectores,
            escritores: this.escritores
        }
        this.readUnlock()
        return stats
    }
}

func simulateClient(cache, clientId) {
    print("Cliente " + clientId + " iniciado")
    
    for (let i = 0; i < 20; i++) {
        let operacion = rand.int(1, 3)
        
        if (operacion == 1) {
            // Operación GET
            let key = "key" + rand.int(1, 10)
            let valor = cache.get(key)
            if (valor != null) {
                print("Cliente " + clientId + " GET " + key + " = " + valor)
            } else {
                print("Cliente " + clientId + " GET " + key + " = MISS")
            }
        } else {
            // Operación SET
            let key = "key" + rand.int(1, 10)
            let valor = "valor" + clientId + "_" + i
            cache.set(key, valor)
            print("Cliente " + clientId + " SET " + key + " = " + valor)
        }
        
        sleep(rand.int(100, 500))
    }
    
    print("Cliente " + clientId + " terminado")
}

func main() {
    let cache = CacheManager(5, 10000)  // 5 items max, 10s expiration
    
    print("Iniciando sistema de cache concurrente")
    
    // Crear múltiples clientes
    for (let i = 1; i <= 5; i++) {
        r2(simulateClient, cache, i)
    }
    
    // Monitorear estadísticas
    for (let i = 0; i < 15; i++) {
        sleep(1000)
        let stats = cache.stats()
        print("Cache stats: " + stats.size + "/" + stats.maxSize + " items, " + stats.lectores + " readers, " + stats.escritores + " writers")
    }
    
    print("Sistema de cache terminado")
}
```

## Proyecto del Módulo: Servidor de Chat Concurrente

```r2
class ChatServer {
    let usuarios
    let salas
    let mensajes
    let servidor
    
    constructor() {
        this.usuarios = {}
        this.salas = {
            "general": {
                nombre: "General",
                usuarios: [],
                mensajes: []
            }
        }
        this.mensajes = []
    }
    
    agregarUsuario(username, conexion) {
        this.usuarios[username] = {
            nombre: username,
            conexion: conexion,
            sala: "general",
            activo: true,
            ultimaActividad: os.time()
        }
        
        // Agregar a sala general
        this.salas["general"].usuarios = this.salas["general"].usuarios.push(username)
        
        print("Usuario " + username + " conectado")
        this.broadcast("general", "SISTEMA: " + username + " se unió al chat")
    }
    
    removerUsuario(username) {
        if (this.usuarios[username] != null) {
            let usuario = this.usuarios[username]
            let sala = usuario.sala
            
            // Remover de sala
            let nuevosUsuarios = []
            for (let i = 0; i < this.salas[sala].usuarios.length(); i++) {
                if (this.salas[sala].usuarios[i] != username) {
                    nuevosUsuarios = nuevosUsuarios.push(this.salas[sala].usuarios[i])
                }
            }
            this.salas[sala].usuarios = nuevosUsuarios
            
            delete this.usuarios[username]
            
            print("Usuario " + username + " desconectado")
            this.broadcast(sala, "SISTEMA: " + username + " abandonó el chat")
        }
    }
    
    enviarMensaje(username, mensaje) {
        if (this.usuarios[username] != null) {
            let usuario = this.usuarios[username]
            let sala = usuario.sala
            
            let mensajeObj = {
                usuario: username,
                contenido: mensaje,
                sala: sala,
                timestamp: os.time()
            }
            
            this.mensajes = this.mensajes.push(mensajeObj)
            this.salas[sala].mensajes = this.salas[sala].mensajes.push(mensajeObj)
            
            // Broadcast a todos los usuarios de la sala
            this.broadcast(sala, username + ": " + mensaje)
            
            // Actualizar última actividad
            usuario.ultimaActividad = os.time()
        }
    }
    
    broadcast(sala, mensaje) {
        if (this.salas[sala] != null) {
            let usuarios = this.salas[sala].usuarios
            for (let i = 0; i < usuarios.length(); i++) {
                let username = usuarios[i]
                if (this.usuarios[username] != null) {
                    // Simular envío de mensaje
                    print("[" + sala + "] " + mensaje)
                }
            }
        }
    }
    
    crearSala(nombre) {
        this.salas[nombre] = {
            nombre: nombre,
            usuarios: [],
            mensajes: []
        }
        print("Sala '" + nombre + "' creada")
    }
    
    cambiarSala(username, nuevaSala) {
        if (this.usuarios[username] != null && this.salas[nuevaSala] != null) {
            let usuario = this.usuarios[username]
            let salaAnterior = usuario.sala
            
            // Remover de sala anterior
            let nuevosUsuarios = []
            for (let i = 0; i < this.salas[salaAnterior].usuarios.length(); i++) {
                if (this.salas[salaAnterior].usuarios[i] != username) {
                    nuevosUsuarios = nuevosUsuarios.push(this.salas[salaAnterior].usuarios[i])
                }
            }
            this.salas[salaAnterior].usuarios = nuevosUsuarios
            
            // Agregar a nueva sala
            this.salas[nuevaSala].usuarios = this.salas[nuevaSala].usuarios.push(username)
            usuario.sala = nuevaSala
            
            this.broadcast(salaAnterior, "SISTEMA: " + username + " abandonó la sala")
            this.broadcast(nuevaSala, "SISTEMA: " + username + " se unió a la sala")
        }
    }
    
    manejarConexion(conexion) {
        // Simular manejo de conexión
        let username = "Usuario" + rand.int(1000, 9999)
        this.agregarUsuario(username, conexion)
        
        // Simular actividad del usuario
        for (let i = 0; i < 10; i++) {
            sleep(rand.int(2000, 5000))
            
            let accion = rand.int(1, 4)
            if (accion == 1) {
                // Enviar mensaje
                let mensaje = "Hola desde " + username + " - mensaje " + (i + 1)
                this.enviarMensaje(username, mensaje)
            } else if (accion == 2 && i > 5) {
                // Cambiar sala
                if (this.salas["desarrollo"] == null) {
                    this.crearSala("desarrollo")
                }
                this.cambiarSala(username, "desarrollo")
            } else if (accion == 3) {
                // Volver a general
                this.cambiarSala(username, "general")
            }
        }
        
        this.removerUsuario(username)
    }
    
    monitorearSalud() {
        while (true) {
            sleep(30000)  // Cada 30 segundos
            
            let ahora = os.time()
            let usuariosActivos = 0
            let totalMensajes = this.mensajes.length()
            
            let usernames = Object.keys(this.usuarios)
            for (let i = 0; i < usernames.length(); i++) {
                let username = usernames[i]
                let usuario = this.usuarios[username]
                
                if (ahora - usuario.ultimaActividad < 300000) {  // 5 minutos
                    usuariosActivos++
                } else {
                    // Usuario inactivo, desconectar
                    this.removerUsuario(username)
                }
            }
            
            print("=== ESTADÍSTICAS DEL SERVIDOR ===")
            print("Usuarios activos: " + usuariosActivos)
            print("Total mensajes: " + totalMensajes)
            print("Salas activas: " + Object.keys(this.salas).length())
        }
    }
    
    iniciar(puerto) {
        print("Iniciando servidor de chat en puerto " + puerto)
        
        // Iniciar monitor de salud
        r2(this.monitorearSalud)
        
        // Simular servidor HTTP
        let self = this
        let servidor = http.server(puerto, func(request, response) {
            if (request.url == "/") {
                response.writeHead(200, {"Content-Type": "text/html"})
                response.write(self.generarHTMLChat())
                response.end()
            } else if (request.url == "/stats") {
                response.writeHead(200, {"Content-Type": "application/json"})
                response.write(JSON.stringify({
                    usuarios: Object.keys(self.usuarios).length(),
                    salas: Object.keys(self.salas).length(),
                    mensajes: self.mensajes.length()
                }))
                response.end()
            }
        })
        
        // Simular conexiones concurrentes
        for (let i = 0; i < 5; i++) {
            r2(this.manejarConexion, "conexion" + i)
        }
        
        print("Servidor de chat iniciado en http://localhost:" + puerto)
    }
    
    generarHTMLChat() {
        let html = `
        <!DOCTYPE html>
        <html>
        <head>
            <title>R2Lang Chat Server</title>
            <style>
                body { font-family: Arial, sans-serif; margin: 20px; }
                .chat-container { border: 1px solid #ccc; height: 400px; overflow-y: scroll; padding: 10px; }
                .message { margin: 5px 0; }
                .system { color: #666; font-style: italic; }
                .user { color: #000; }
                .stats { margin-top: 20px; padding: 10px; background: #f0f0f0; }
            </style>
        </head>
        <body>
            <h1>R2Lang Chat Server</h1>
            <div class="chat-container">
                <div class="message system">Sistema: Servidor de chat activo</div>
                <div class="message user">Usuario1234: ¡Hola a todos!</div>
                <div class="message user">Usuario5678: ¿Cómo están?</div>
                <div class="message system">Sistema: Usuario9012 se unió al chat</div>
                <div class="message user">Usuario9012: Saludos</div>
            </div>
            <div class="stats">
                <h3>Estadísticas del Servidor</h3>
                <p>Usuarios conectados: 5</p>
                <p>Salas activas: 2</p>
                <p>Mensajes enviados: 127</p>
                <p><a href="/stats">Ver estadísticas JSON</a></p>
            </div>
        </body>
        </html>
        `
        return html
    }
}

func main() {
    let chatServer = ChatServer()
    chatServer.iniciar(8080)
    
    // Mantener servidor activo
    sleep(60000)  // 1 minuto
    
    print("Servidor de chat terminado")
}
```

## Resumen del Módulo

### Conceptos de Concurrencia Aprendidos
- ✅ **Goroutines**: Creación y gestión de tareas concurrentes
- ✅ **Worker Pools**: Patrones de workers para procesamiento paralelo
- ✅ **Producer-Consumer**: Comunicación entre procesos
- ✅ **Map-Reduce**: Procesamiento distribuido de datos
- ✅ **Sincronización**: Coordinación entre goroutines
- ✅ **Sistemas concurrentes**: Aplicaciones paralelas complejas

### Habilidades Desarrolladas
- ✅ Diseñar aplicaciones concurrentes eficientes
- ✅ Implementar patrones de concurrencia clásicos
- ✅ Crear sistemas distribuidos básicos
- ✅ Manejar sincronización entre procesos
- ✅ Desarrollar aplicaciones web concurrentes
- ✅ Optimizar rendimiento con paralelización

### Próximo Módulo

En el **Módulo 9** aprenderás:
- Testing avanzado con BDD
- Debugging y profiling
- Optimización de rendimiento
- Despliegue y distribución

¡Excelente trabajo! Ahora puedes crear aplicaciones R2Lang que aprovechan el poder de la concurrencia.