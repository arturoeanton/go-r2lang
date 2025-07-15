# Curso R2Lang - Módulo 7: Bibliotecas Integradas y Programación de Sistemas

## Introducción

En este módulo explorarás las potentes bibliotecas integradas de R2Lang v2 que te permiten interactuar con el sistema operativo, manejar archivos, realizar operaciones de red, y crear aplicaciones más complejas. La nueva arquitectura modular organiza estas bibliotecas en `pkg/r2libs/` con implementaciones especializadas.

### Arquitectura de Bibliotecas v2

```
pkg/r2libs/ (3,701 LOC, 18 archivos especializados):
├── r2http.go (410 LOC)        # Servidor HTTP completo
├── r2httpclient.go (324 LOC)  # Cliente HTTP
├── r2string.go (194 LOC)      # Manipulación de strings
├── r2io.go (194 LOC)         # Operaciones de archivos
├── r2os.go (245 LOC)         # Interface con SO
├── r2math.go (87 LOC)        # Funciones matemáticas
├── r2goroutine.go (237 LOC)  # Concurrencia
├── r2hack.go (509 LOC)       # Cryptografía y seguridad
├── r2print.go (365 LOC)      # Output avanzado
└── [9 bibliotecas más]
```

Cada biblioteca es independiente, con testing completo y documentación especializada.

## Biblioteca de Entrada/Salida (r2io)

### 1. Manejo de Archivos Básico

```r2
func main() {
    // Escribir archivo
    let contenido = "¡Hola mundo desde R2Lang!"
    let resultado = io.writeFile("mi_archivo.txt", contenido)
    
    if (resultado) {
        print("Archivo creado exitosamente")
    } else {
        print("Error al crear archivo")
    }
    
    // Leer archivo
    let contenidoLeido = io.readFile("mi_archivo.txt")
    if (contenidoLeido != null) {
        print("Contenido del archivo:")
        print(contenidoLeido)
    } else {
        print("Error al leer archivo")
    }
}
```

### 2. Operaciones de Archivos Avanzadas

```r2
func gestionarArchivos() {
    // Crear directorio
    let dirCreado = io.mkdir("mi_directorio")
    print("Directorio creado:", dirCreado)
    
    // Verificar si archivo existe
    let existe = io.exists("mi_archivo.txt")
    print("Archivo existe:", existe)
    
    // Obtener información del archivo
    let info = io.fileInfo("mi_archivo.txt")
    if (info != null) {
        print("Tamaño del archivo:", info.size, "bytes")
        print("Última modificación:", info.lastModified)
    }
    
    // Listar archivos en directorio
    let archivos = io.listFiles(".")
    print("Archivos en directorio actual:")
    for (let archivo in archivos) {
        print("- " + archivo)
    }
    
    // Copiar archivo
    let copiado = io.copyFile("mi_archivo.txt", "mi_directorio/copia.txt")
    print("Archivo copiado:", copiado)
    
    // Mover archivo
    let movido = io.moveFile("mi_archivo.txt", "mi_directorio/movido.txt")
    print("Archivo movido:", movido)
    
    // Eliminar archivo
    let eliminado = io.deleteFile("archivo_temporal.txt")
    print("Archivo eliminado:", eliminado)
}

func main() {
    gestionarArchivos()
}
```

### 3. Sistema de Logging

```r2
class Logger {
    let archivo
    let nivel
    
    constructor(archivo, nivel) {
        this.archivo = archivo
        this.nivel = nivel
        
        // Crear archivo de log si no existe
        if (!io.exists(archivo)) {
            io.writeFile(archivo, "=== LOG INICIADO ===\n")
        }
    }
    
    log(mensaje, tipo) {
        let timestamp = "2024-01-15 10:30:00"  // Simulado
        let entrada = "[" + timestamp + "] " + tipo + ": " + mensaje + "\n"
        
        // Leer contenido actual
        let contenidoActual = io.readFile(this.archivo)
        if (contenidoActual == null) {
            contenidoActual = ""
        }
        
        // Agregar nueva entrada
        let nuevoContenido = contenidoActual + entrada
        io.writeFile(this.archivo, nuevoContenido)
        
        // También imprimir en consola
        if (tipo == "ERROR" || this.nivel == "DEBUG") {
            print(entrada)
        }
    }
    
    info(mensaje) {
        this.log(mensaje, "INFO")
    }
    
    warn(mensaje) {
        this.log(mensaje, "WARN")
    }
    
    error(mensaje) {
        this.log(mensaje, "ERROR")
    }
    
    debug(mensaje) {
        if (this.nivel == "DEBUG") {
            this.log(mensaje, "DEBUG")
        }
    }
}

func main() {
    let logger = Logger("aplicacion.log", "INFO")
    
    logger.info("Aplicación iniciada")
    logger.warn("Esta es una advertencia")
    logger.error("Ocurrió un error")
    logger.debug("Mensaje de debug (no se verá)")
    
    // Cambiar nivel de log
    logger.nivel = "DEBUG"
    logger.debug("Ahora sí se ve este debug")
    
    print("Revisa el archivo aplicacion.log para ver los logs")
}
```

## Biblioteca de Sistema Operativo (r2os)

### 1. Información del Sistema

```r2
func infoSistema() {
    // Obtener información del sistema operativo
    let os = os.system()
    print("Sistema operativo:", os.name)
    print("Arquitectura:", os.arch)
    print("Versión:", os.version)
    
    // Variables de entorno
    let home = os.getEnv("HOME")
    print("Directorio home:", home)
    
    let path = os.getEnv("PATH")
    print("PATH:", path)
    
    // Establecer variable de entorno
    os.setEnv("MI_VARIABLE", "mi_valor")
    let miVar = os.getEnv("MI_VARIABLE")
    print("Mi variable:", miVar)
    
    // Directorio actual
    let dirActual = os.pwd()
    print("Directorio actual:", dirActual)
    
    // Cambiar directorio
    let cambioDir = os.chdir("/tmp")
    print("Cambio de directorio:", cambioDir)
    print("Nuevo directorio:", os.pwd())
}

func main() {
    infoSistema()
}
```

### 2. Ejecución de Comandos

```r2
func ejecutarComandos() {
    // Ejecutar comando simple
    let resultado = os.exec("ls -la")
    print("Resultado de ls -la:")
    print(resultado)
    
    // Ejecutar comando con parámetros
    let fecha = os.exec("date")
    print("Fecha actual:", fecha)
    
    // Ejecutar comando y capturar código de salida
    let salida = os.execWithCode("echo 'Hola mundo'")
    print("Salida:", salida.output)
    print("Código de salida:", salida.exitCode)
    
    // Comando que falla
    let error = os.execWithCode("comando_inexistente")
    print("Código de salida de comando inválido:", error.exitCode)
}

func main() {
    ejecutarComandos()
}
```

### 3. Utilidades del Sistema

```r2
func utilidadesSistema() {
    // Obtener PID del proceso
    let pid = os.getPid()
    print("PID del proceso:", pid)
    
    // Dormir/pausar ejecución
    print("Pausando por 2 segundos...")
    os.sleep(2000)  // 2000 milisegundos
    print("Continuando...")
    
    // Obtener tiempo actual
    let tiempo = os.time()
    print("Timestamp actual:", tiempo)
    
    // Crear proceso hijo (simulado)
    let proceso = os.fork()
    if (proceso.isChild) {
        print("Soy el proceso hijo")
    } else {
        print("Soy el proceso padre, hijo PID:", proceso.childPid)
    }
}

func main() {
    utilidadesSistema()
}
```

## Biblioteca de Red (r2http y r2httpclient)

### 1. Servidor HTTP Básico

```r2
func crearServidorBasico() {
    // Crear servidor HTTP
    let servidor = http.server(8080, func(request, response) {
        print("Solicitud recibida:", request.method, request.url)
        
        // Responder con HTML
        response.writeHead(200, {"Content-Type": "text/html"})
        response.write("<h1>¡Hola desde R2Lang!</h1>")
        response.write("<p>Servidor funcionando correctamente</p>")
        response.end()
    })
    
    print("Servidor HTTP iniciado en puerto 8080")
    print("Visita http://localhost:8080 en tu navegador")
}

func main() {
    crearServidorBasico()
}
```

### 2. API REST Completa

```r2
class APIServer {
    let port
    let routes
    
    constructor(port) {
        this.port = port
        this.routes = {
            "GET": {},
            "POST": {},
            "PUT": {},
            "DELETE": {}
        }
    }
    
    get(path, handler) {
        this.routes["GET"][path] = handler
    }
    
    post(path, handler) {
        this.routes["POST"][path] = handler
    }
    
    put(path, handler) {
        this.routes["PUT"][path] = handler
    }
    
    delete(path, handler) {
        this.routes["DELETE"][path] = handler
    }
    
    handleRequest(request, response) {
        let method = request.method
        let path = request.url
        
        print("API Request:", method, path)
        
        // Buscar handler para la ruta
        if (this.routes[method][path] != null) {
            this.routes[method][path](request, response)
        } else {
            // Ruta no encontrada
            response.writeHead(404, {"Content-Type": "application/json"})
            response.write('{"error": "Ruta no encontrada"}')
            response.end()
        }
    }
    
    start() {
        let self = this
        let servidor = http.server(this.port, func(request, response) {
            self.handleRequest(request, response)
        })
        
        print("API Server iniciado en puerto", this.port)
        return servidor
    }
}

// Base de datos simulada
let usuarios = [
    {"id": 1, "nombre": "Ana García", "email": "ana@test.com"},
    {"id": 2, "nombre": "Carlos López", "email": "carlos@test.com"}
]

func main() {
    let api = APIServer(3000)
    
    // Ruta GET para obtener todos los usuarios
    api.get("/usuarios", func(request, response) {
        response.writeHead(200, {"Content-Type": "application/json"})
        response.write(JSON.stringify(usuarios))
        response.end()
    })
    
    // Ruta GET para obtener usuario por ID
    api.get("/usuarios/1", func(request, response) {
        response.writeHead(200, {"Content-Type": "application/json"})
        response.write(JSON.stringify(usuarios[0]))
        response.end()
    })
    
    // Ruta POST para crear usuario
    api.post("/usuarios", func(request, response) {
        // Simulamos agregar usuario
        let nuevoUsuario = {
            "id": usuarios.length() + 1,
            "nombre": "Usuario Nuevo",
            "email": "nuevo@test.com"
        }
        usuarios = usuarios.push(nuevoUsuario)
        
        response.writeHead(201, {"Content-Type": "application/json"})
        response.write(JSON.stringify(nuevoUsuario))
        response.end()
    })
    
    // Ruta DELETE para eliminar usuario
    api.delete("/usuarios/1", func(request, response) {
        response.writeHead(204, {})
        response.end()
    })
    
    api.start()
}
```

### 3. Cliente HTTP

```r2
func clienteHTTP() {
    // GET request simple
    let respuesta = http.get("https://api.github.com/users/octocat")
    if (respuesta.status == 200) {
        print("Usuario de GitHub:")
        print(respuesta.body)
    } else {
        print("Error:", respuesta.status)
    }
    
    // POST request con datos
    let datos = {
        "titulo": "Mi post",
        "contenido": "Contenido del post",
        "autor": "R2Lang"
    }
    
    let respuestaPost = http.post("https://jsonplaceholder.typicode.com/posts", {
        "Content-Type": "application/json",
        "datos": JSON.stringify(datos)
    })
    
    print("POST Response:")
    print("Status:", respuestaPost.status)
    print("Body:", respuestaPost.body)
    
    // Request con headers personalizados
    let headers = {
        "Authorization": "Bearer mi_token",
        "User-Agent": "R2Lang-Client/1.0"
    }
    
    let respuestaConHeaders = http.get("https://api.ejemplo.com/datos", headers)
    print("Response con headers:", respuestaConHeaders.status)
}

func main() {
    clienteHTTP()
}
```

## Biblioteca de Strings (r2string)

### 1. Manipulación Avanzada de Strings

```r2
func manipularStrings() {
    let texto = "  Hola Mundo desde R2Lang  "
    
    // Operaciones básicas
    print("Original:", '"' + texto + '"')
    print("Longitud:", string.length(texto))
    print("Mayúsculas:", string.upper(texto))
    print("Minúsculas:", string.lower(texto))
    print("Sin espacios:", string.trim(texto))
    
    // Búsqueda y reemplazo
    let frase = "R2Lang es genial. R2Lang es poderoso."
    print("Contiene 'genial':", string.contains(frase, "genial"))
    print("Posición de 'es':", string.indexOf(frase, "es"))
    print("Última posición de 'R2Lang':", string.lastIndexOf(frase, "R2Lang"))
    print("Reemplazar:", string.replace(frase, "R2Lang", "Este lenguaje"))
    
    // Dividir y unir
    let palabras = string.split(frase, " ")
    print("Palabras:", palabras)
    let reunidas = string.join(palabras, "-")
    print("Reunidas:", reunidas)
    
    // Subcadenas
    let subcadena = string.substring(frase, 0, 6)
    print("Subcadena (0-6):", subcadena)
    
    // Padding
    let numero = "42"
    print("Pad izquierda:", string.padStart(numero, 5, "0"))
    print("Pad derecha:", string.padEnd(numero, 5, "0"))
}

func main() {
    manipularStrings()
}
```

### 2. Validación y Formato

```r2
func validarFormatos() {
    let email = "usuario@dominio.com"
    let telefono = "123-456-7890"
    let numero = "12345"
    let url = "https://www.ejemplo.com"
    
    // Validaciones básicas
    print("Email válido:", string.isEmail(email))
    print("Es número:", string.isNumeric(numero))
    print("Es URL:", string.isURL(url))
    
    // Formato de strings
    let nombre = "juan pérez"
    print("Capitalizado:", string.capitalize(nombre))
    print("Título:", string.title(nombre))
    
    // Limpieza de texto
    let textoSucio = "  texto    con    espacios   "
    print("Espacios normalizados:", string.normalizeSpaces(textoSucio))
    
    // Generar strings
    print("Repetir 'A':", string.repeat("A", 5))
    print("String aleatorio:", string.random(10))
}

func main() {
    validarFormatos()
}
```

## Biblioteca de Matemáticas (r2math)

### 1. Operaciones Matemáticas Avanzadas

```r2
func operacionesMatematicas() {
    let numeros = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
    
    // Operaciones básicas
    print("Suma:", math.sum(numeros))
    print("Promedio:", math.average(numeros))
    print("Máximo:", math.max(numeros))
    print("Mínimo:", math.min(numeros))
    
    // Funciones matemáticas
    let angulo = 45
    print("Seno de", angulo, "°:", math.sin(angulo))
    print("Coseno de", angulo, "°:", math.cos(angulo))
    print("Tangente de", angulo, "°:", math.tan(angulo))
    
    // Logaritmos y exponenciales
    let base = 2.718281828  // e
    print("Logaritmo natural de 10:", math.log(10))
    print("Logaritmo base 10:", math.log10(100))
    print("e^2:", math.exp(2))
    print("2^8:", math.pow(2, 8))
    
    // Raíces
    print("Raíz cuadrada de 16:", math.sqrt(16))
    print("Raíz cúbica de 27:", math.cbrt(27))
    
    // Redondeo
    let decimal = 3.14159
    print("Redondeo:", math.round(decimal))
    print("Hacia arriba:", math.ceil(decimal))
    print("Hacia abajo:", math.floor(decimal))
    print("Truncar:", math.trunc(decimal))
}

func main() {
    operacionesMatematicas()
}
```

### 2. Estadísticas y Análisis

```r2
func estadisticas() {
    let datos = [2, 4, 6, 8, 10, 12, 14, 16, 18, 20]
    
    // Estadísticas descriptivas
    print("Datos:", datos)
    print("Media:", math.mean(datos))
    print("Mediana:", math.median(datos))
    print("Moda:", math.mode(datos))
    print("Desviación estándar:", math.stddev(datos))
    print("Varianza:", math.variance(datos))
    
    // Percentiles
    print("Percentil 25:", math.percentile(datos, 25))
    print("Percentil 75:", math.percentile(datos, 75))
    
    // Análisis de distribución
    print("Rango:", math.range(datos))
    print("Cuartil 1:", math.quartile(datos, 1))
    print("Cuartil 3:", math.quartile(datos, 3))
}

func main() {
    estadisticas()
}
```

## Biblioteca de Cryptografía (r2hack)

### 1. Funciones de Hash

```r2
func funcionesHash() {
    let texto = "Hola mundo desde R2Lang"
    
    // Diferentes tipos de hash
    print("Texto original:", texto)
    print("MD5:", hash.md5(texto))
    print("SHA1:", hash.sha1(texto))
    print("SHA256:", hash.sha256(texto))
    print("SHA512:", hash.sha512(texto))
    
    // Verificación de integridad
    let archivoContent = "Contenido importante del archivo"
    let hashOriginal = hash.sha256(archivoContent)
    print("Hash original:", hashOriginal)
    
    // Simular verificación
    let archivoModificado = "Contenido importante del archivo modificado"
    let hashModificado = hash.sha256(archivoModificado)
    print("Hash modificado:", hashModificado)
    print("Integridad mantenida:", hashOriginal == hashModificado)
}

func main() {
    funcionesHash()
}
```

### 2. Codificación y Encriptación

```r2
func codificacionEncriptacion() {
    let mensaje = "Mensaje secreto"
    
    // Base64
    let base64Encoded = crypto.base64Encode(mensaje)
    print("Base64 encoded:", base64Encoded)
    let base64Decoded = crypto.base64Decode(base64Encoded)
    print("Base64 decoded:", base64Decoded)
    
    // URL encoding
    let url = "https://ejemplo.com/buscar?q=R2Lang es genial"
    let urlEncoded = crypto.urlEncode(url)
    print("URL encoded:", urlEncoded)
    let urlDecoded = crypto.urlDecode(urlEncoded)
    print("URL decoded:", urlDecoded)
    
    // Hex encoding
    let hexEncoded = crypto.hexEncode(mensaje)
    print("Hex encoded:", hexEncoded)
    let hexDecoded = crypto.hexDecode(hexEncoded)
    print("Hex decoded:", hexDecoded)
}

func main() {
    codificacionEncriptacion()
}
```

## Proyecto del Módulo: Sistema de Monitoreo

```r2
class MonitoreoSistema {
    let logger
    let configuracion
    let metricas
    
    constructor() {
        this.logger = Logger("monitoreo.log", "INFO")
        this.configuracion = {
            "intervalo": 5000,  // 5 segundos
            "alertas": true,
            "maxMemoria": 1024 * 1024 * 100  // 100MB
        }
        this.metricas = {
            "cpu": [],
            "memoria": [],
            "disco": []
        }
    }
    
    obtenerMetricasSistema() {
        // Simular métricas del sistema
        let memoria = os.getMemoryUsage()
        let cpu = os.getCpuUsage()
        let disco = os.getDiskUsage()
        
        return {
            "timestamp": os.time(),
            "memoria": memoria,
            "cpu": cpu,
            "disco": disco
        }
    }
    
    verificarAlertas(metricas) {
        // Verificar memoria
        if (metricas.memoria > this.configuracion.maxMemoria) {
            let mensaje = "ALERTA: Uso de memoria alto: " + metricas.memoria
            this.logger.warn(mensaje)
            
            // Enviar notificación
            this.enviarNotificacion("Memoria Alta", mensaje)
        }
        
        // Verificar CPU
        if (metricas.cpu > 80) {
            let mensaje = "ALERTA: Uso de CPU alto: " + metricas.cpu + "%"
            this.logger.warn(mensaje)
            this.enviarNotificacion("CPU Alta", mensaje)
        }
        
        // Verificar disco
        if (metricas.disco > 90) {
            let mensaje = "ALERTA: Uso de disco alto: " + metricas.disco + "%"
            this.logger.error(mensaje)
            this.enviarNotificacion("Disco Lleno", mensaje)
        }
    }
    
    enviarNotificacion(titulo, mensaje) {
        // Enviar notificación via HTTP
        let datos = {
            "titulo": titulo,
            "mensaje": mensaje,
            "timestamp": os.time()
        }
        
        let respuesta = http.post("https://webhook.site/notificaciones", {
            "Content-Type": "application/json",
            "datos": JSON.stringify(datos)
        })
        
        if (respuesta.status == 200) {
            this.logger.info("Notificación enviada: " + titulo)
        } else {
            this.logger.error("Error enviando notificación: " + respuesta.status)
        }
    }
    
    generarReporte() {
        let reporte = "=== REPORTE DE MONITOREO ===\n"
        reporte = reporte + "Fecha: " + os.time() + "\n\n"
        
        // Estadísticas de CPU
        if (this.metricas.cpu.length() > 0) {
            let promedioCpu = math.average(this.metricas.cpu)
            let maxCpu = math.max(this.metricas.cpu)
            
            reporte = reporte + "CPU:\n"
            reporte = reporte + "  Promedio: " + promedioCpu + "%\n"
            reporte = reporte + "  Máximo: " + maxCpu + "%\n\n"
        }
        
        // Estadísticas de memoria
        if (this.metricas.memoria.length() > 0) {
            let promedioMem = math.average(this.metricas.memoria)
            let maxMem = math.max(this.metricas.memoria)
            
            reporte = reporte + "Memoria:\n"
            reporte = reporte + "  Promedio: " + promedioMem + " MB\n"
            reporte = reporte + "  Máximo: " + maxMem + " MB\n\n"
        }
        
        // Guardar reporte
        let nombreArchivo = "reporte_" + os.time() + ".txt"
        io.writeFile(nombreArchivo, reporte)
        
        this.logger.info("Reporte generado: " + nombreArchivo)
        return nombreArchivo
    }
    
    iniciarMonitoreo() {
        this.logger.info("Sistema de monitoreo iniciado")
        
        // Crear servidor web para dashboard
        let self = this
        let servidor = http.server(8080, func(request, response) {
            if (request.url == "/") {
                response.writeHead(200, {"Content-Type": "text/html"})
                response.write(self.generarDashboard())
                response.end()
            } else if (request.url == "/api/metricas") {
                response.writeHead(200, {"Content-Type": "application/json"})
                response.write(JSON.stringify(self.metricas))
                response.end()
            } else if (request.url == "/api/reporte") {
                let archivo = self.generarReporte()
                response.writeHead(200, {"Content-Type": "text/plain"})
                response.write(io.readFile(archivo))
                response.end()
            }
        })
        
        print("Dashboard disponible en: http://localhost:8080")
        
        // Bucle principal de monitoreo
        while (true) {
            let metricas = this.obtenerMetricasSistema()
            
            // Almacenar métricas
            this.metricas.cpu = this.metricas.cpu.push(metricas.cpu)
            this.metricas.memoria = this.metricas.memoria.push(metricas.memoria)
            this.metricas.disco = this.metricas.disco.push(metricas.disco)
            
            // Mantener solo las últimas 100 métricas
            if (this.metricas.cpu.length() > 100) {
                this.metricas.cpu = this.metricas.cpu.slice(1)
                this.metricas.memoria = this.metricas.memoria.slice(1)
                this.metricas.disco = this.metricas.disco.slice(1)
            }
            
            // Verificar alertas
            if (this.configuracion.alertas) {
                this.verificarAlertas(metricas)
            }
            
            // Log de métricas
            this.logger.debug("CPU: " + metricas.cpu + "%, Memoria: " + metricas.memoria + "MB, Disco: " + metricas.disco + "%")
            
            // Esperar antes de la siguiente medición
            os.sleep(this.configuracion.intervalo)
        }
    }
    
    generarDashboard() {
        let html = `
        <!DOCTYPE html>
        <html>
        <head>
            <title>Monitor R2Lang</title>
            <style>
                body { font-family: Arial, sans-serif; margin: 20px; }
                .metrica { display: inline-block; margin: 20px; padding: 20px; border: 1px solid #ccc; }
                .high { background-color: #ffcccc; }
                .normal { background-color: #ccffcc; }
            </style>
        </head>
        <body>
            <h1>Sistema de Monitoreo R2Lang</h1>
            <div class="metrica normal">
                <h3>CPU</h3>
                <p>Uso actual: 45%</p>
            </div>
            <div class="metrica normal">
                <h3>Memoria</h3>
                <p>Uso actual: 512MB</p>
            </div>
            <div class="metrica normal">
                <h3>Disco</h3>
                <p>Uso actual: 65%</p>
            </div>
            <hr>
            <a href="/api/metricas">Ver métricas JSON</a> | 
            <a href="/api/reporte">Generar reporte</a>
        </body>
        </html>
        `
        return html
    }
}

func main() {
    let monitor = MonitoreoSistema()
    monitor.iniciarMonitoreo()
}
```

## Resumen del Módulo

### Bibliotecas Exploradas
- ✅ **r2io**: Manejo de archivos y directorios
- ✅ **r2os**: Interface con sistema operativo
- ✅ **r2http**: Servidor HTTP y APIs REST
- ✅ **r2httpclient**: Cliente HTTP
- ✅ **r2string**: Manipulación avanzada de strings
- ✅ **r2math**: Operaciones matemáticas y estadísticas
- ✅ **r2hack**: Cryptografía y codificación

### Habilidades Desarrolladas
- ✅ Crear aplicaciones que interactúan con el sistema
- ✅ Desarrollar APIs REST funcionales
- ✅ Implementar clientes HTTP robustos
- ✅ Manejar archivos y directorios eficientemente
- ✅ Aplicar funciones criptográficas
- ✅ Crear sistemas de monitoreo
- ✅ Integrar múltiples bibliotecas en proyectos complejos

### Próximo Módulo

En el **Módulo 8** aprenderás:
- Concurrencia avanzada con goroutines
- Patrones de diseño concurrente
- Comunicación entre procesos
- Optimización de rendimiento

¡Excelente trabajo! Ahora puedes crear aplicaciones R2Lang que interactúan con el sistema operativo y la red.