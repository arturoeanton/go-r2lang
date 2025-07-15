# Curso R2Lang - Módulo 9: Testing Avanzado y Debugging

## Introducción

En este módulo aprenderás técnicas avanzadas de testing, debugging y profiling para crear aplicaciones R2Lang robustas y confiables. La nueva arquitectura v2 incluye un sistema de testing BDD integrado y herramientas de debugging avanzadas.

### Sistema de Testing v2

```
pkg/r2libs/r2test.go (Testing Framework):
├── BDD Syntax Support      # Given/When/Then/And
├── Assertion Functions     # assertEqual, assertTrue, etc.
├── Test Runners           # Ejecución automática de tests
├── Mock Support           # Simulación de dependencias
├── Coverage Reporting     # Reporte de cobertura
└── Integration Testing    # Tests de integración
```

El sistema de testing está completamente integrado con el intérprete, permitiendo testing nativo del código R2Lang.

## Testing BDD (Behavior-Driven Development)

### 1. Sintaxis BDD Básica

```r2
TestCase "Calculadora básica" {
    Given {
        let calc = Calculadora()
        assertEqual(calc.resultado, 0)
    }
    
    When {
        calc.sumar(5)
        calc.sumar(3)
    }
    
    Then {
        assertEqual(calc.resultado, 8)
        assertTrue(calc.resultado > 0)
    }
}

TestCase "División por cero" {
    Given {
        let calc = Calculadora()
    }
    
    When {
        let resultado = calc.dividir(10, 0)
    }
    
    Then {
        assertEqual(resultado, null)
        assertTrue(calc.hayError)
    }
}

// Clase a testear
class Calculadora {
    let resultado
    let hayError
    
    constructor() {
        this.resultado = 0
        this.hayError = false
    }
    
    sumar(numero) {
        this.resultado = this.resultado + numero
        return this.resultado
    }
    
    dividir(a, b) {
        if (b == 0) {
            this.hayError = true
            return null
        }
        this.resultado = a / b
        return this.resultado
    }
}

func main() {
    print("Ejecutando tests de Calculadora...")
    // Los TestCase se ejecutan automáticamente
}
```

### 2. Tests con Configuración Compleja

```r2
TestCase "Sistema de usuarios completo" {
    Given {
        let sistema = SistemaUsuarios()
        let usuario1 = {
            nombre: "Ana García",
            email: "ana@test.com",
            edad: 25
        }
        let usuario2 = {
            nombre: "Carlos López", 
            email: "carlos@test.com",
            edad: 30
        }
    }
    
    When {
        sistema.agregarUsuario(usuario1)
        sistema.agregarUsuario(usuario2)
    }
    
    Then {
        assertEqual(sistema.getTotalUsuarios(), 2)
        
        let usuarioEncontrado = sistema.buscarPorEmail("ana@test.com")
        assertEqual(usuarioEncontrado.nombre, "Ana García")
        
        let usuariosAdultos = sistema.filtrarPorEdad(18)
        assertEqual(usuariosAdultos.length(), 2)
    }
    
    And {
        // Verificar que no se pueden agregar usuarios duplicados
        let duplicado = sistema.agregarUsuario(usuario1)
        assertEqual(duplicado, false)
        assertEqual(sistema.getTotalUsuarios(), 2)
    }
}

TestCase "Validación de email" {
    Given {
        let sistema = SistemaUsuarios()
        let usuarioInvalido = {
            nombre: "Usuario Test",
            email: "email-invalido",
            edad: 20
        }
    }
    
    When {
        let resultado = sistema.agregarUsuario(usuarioInvalido)
    }
    
    Then {
        assertEqual(resultado, false)
        assertEqual(sistema.getTotalUsuarios(), 0)
        assertTrue(sistema.getUltimoError().contains("email"))
    }
}

class SistemaUsuarios {
    let usuarios
    let ultimoError
    
    constructor() {
        this.usuarios = []
        this.ultimoError = ""
    }
    
    agregarUsuario(usuario) {
        // Validar email
        if (!this.validarEmail(usuario.email)) {
            this.ultimoError = "Formato de email inválido"
            return false
        }
        
        // Verificar duplicados
        if (this.buscarPorEmail(usuario.email) != null) {
            this.ultimoError = "Usuario ya existe"
            return false
        }
        
        this.usuarios = this.usuarios.push(usuario)
        return true
    }
    
    validarEmail(email) {
        return email.contains("@") && email.contains(".")
    }
    
    buscarPorEmail(email) {
        for (let i = 0; i < this.usuarios.length(); i++) {
            if (this.usuarios[i].email == email) {
                return this.usuarios[i]
            }
        }
        return null
    }
    
    filtrarPorEdad(edadMinima) {
        let resultado = []
        for (let i = 0; i < this.usuarios.length(); i++) {
            if (this.usuarios[i].edad >= edadMinima) {
                resultado = resultado.push(this.usuarios[i])
            }
        }
        return resultado
    }
    
    getTotalUsuarios() {
        return this.usuarios.length()
    }
    
    getUltimoError() {
        return this.ultimoError
    }
}

func main() {
    print("Ejecutando tests del Sistema de Usuarios...")
}
```

### 3. Tests de Integración

```r2
TestCase "API completa de usuarios" {
    Given {
        let servidor = ServidorAPI(3001)
        let cliente = ClienteHTTP("http://localhost:3001")
        
        // Iniciar servidor en background
        r2(servidor.iniciar)
        sleep(1000)  // Esperar que inicie
    }
    
    When {
        // Crear usuario via POST
        let nuevoUsuario = {
            "nombre": "Test User",
            "email": "test@example.com",
            "edad": 28
        }
        
        let respuestaCrear = cliente.post("/usuarios", nuevoUsuario)
        let usuarioCreado = JSON.parse(respuestaCrear.body)
    }
    
    Then {
        assertEqual(respuestaCrear.status, 201)
        assertEqual(usuarioCreado.nombre, "Test User")
        assertTrue(usuarioCreado.id != null)
    }
    
    And {
        // Verificar que el usuario se puede obtener
        let respuestaGet = cliente.get("/usuarios/" + usuarioCreado.id)
        let usuarioObtenido = JSON.parse(respuestaGet.body)
        
        assertEqual(respuestaGet.status, 200)
        assertEqual(usuarioObtenido.email, "test@example.com")
    }
    
    And {
        // Verificar lista de usuarios
        let respuestaLista = cliente.get("/usuarios")
        let listaUsuarios = JSON.parse(respuestaLista.body)
        
        assertEqual(respuestaLista.status, 200)
        assertTrue(listaUsuarios.length() >= 1)
    }
}

class ServidorAPI {
    let puerto
    let usuarios
    let servidor
    
    constructor(puerto) {
        this.puerto = puerto
        this.usuarios = []
    }
    
    iniciar() {
        let self = this
        this.servidor = http.server(this.puerto, func(request, response) {
            self.manejarRequest(request, response)
        })
        print("Servidor API iniciado en puerto " + this.puerto)
    }
    
    manejarRequest(request, response) {
        let path = request.url
        let method = request.method
        
        if (method == "POST" && path == "/usuarios") {
            this.crearUsuario(request, response)
        } else if (method == "GET" && path.startsWith("/usuarios/")) {
            let id = path.substring(10)  // Extraer ID
            this.obtenerUsuario(id, response)
        } else if (method == "GET" && path == "/usuarios") {
            this.listarUsuarios(response)
        } else {
            response.writeHead(404, {"Content-Type": "application/json"})
            response.write('{"error": "Not found"}')
            response.end()
        }
    }
    
    crearUsuario(request, response) {
        let datos = JSON.parse(request.body)
        let nuevoUsuario = {
            id: this.usuarios.length() + 1,
            nombre: datos.nombre,
            email: datos.email,
            edad: datos.edad
        }
        
        this.usuarios = this.usuarios.push(nuevoUsuario)
        
        response.writeHead(201, {"Content-Type": "application/json"})
        response.write(JSON.stringify(nuevoUsuario))
        response.end()
    }
    
    obtenerUsuario(id, response) {
        let usuario = null
        for (let i = 0; i < this.usuarios.length(); i++) {
            if (this.usuarios[i].id == parseInt(id)) {
                usuario = this.usuarios[i]
                break
            }
        }
        
        if (usuario != null) {
            response.writeHead(200, {"Content-Type": "application/json"})
            response.write(JSON.stringify(usuario))
        } else {
            response.writeHead(404, {"Content-Type": "application/json"})
            response.write('{"error": "Usuario no encontrado"}')
        }
        response.end()
    }
    
    listarUsuarios(response) {
        response.writeHead(200, {"Content-Type": "application/json"})
        response.write(JSON.stringify(this.usuarios))
        response.end()
    }
}

class ClienteHTTP {
    let baseUrl
    
    constructor(baseUrl) {
        this.baseUrl = baseUrl
    }
    
    get(path) {
        return http.get(this.baseUrl + path)
    }
    
    post(path, datos) {
        return http.post(this.baseUrl + path, {
            "Content-Type": "application/json",
            "data": JSON.stringify(datos)
        })
    }
}

func main() {
    print("Ejecutando tests de integración...")
}
```

## Debugging y Profiling

### 1. Técnicas de Debugging

```r2
class Debug {
    let nivel
    let archivo
    
    constructor(nivel) {
        this.nivel = nivel
        this.archivo = "debug.log"
    }
    
    log(mensaje, nivel) {
        if (nivel >= this.nivel) {
            let timestamp = os.time()
            let entry = "[" + timestamp + "] " + mensaje
            
            // Imprimir en consola
            print(entry)
            
            // Escribir a archivo
            let contenido = io.readFile(this.archivo)
            if (contenido == null) {
                contenido = ""
            }
            io.writeFile(this.archivo, contenido + entry + "\n")
        }
    }
    
    trace(mensaje) {
        this.log("TRACE: " + mensaje, 1)
    }
    
    debug(mensaje) {
        this.log("DEBUG: " + mensaje, 2)
    }
    
    info(mensaje) {
        this.log("INFO: " + mensaje, 3)
    }
    
    warn(mensaje) {
        this.log("WARN: " + mensaje, 4)
    }
    
    error(mensaje) {
        this.log("ERROR: " + mensaje, 5)
    }
}

func funcionProblematica(numero) {
    let debug = Debug(2)
    debug.debug("Iniciando función con número: " + numero)
    
    if (numero < 0) {
        debug.warn("Número negativo recibido: " + numero)
        return null
    }
    
    let resultado = 0
    for (let i = 1; i <= numero; i++) {
        resultado = resultado + i
        debug.trace("Iteración " + i + ", resultado parcial: " + resultado)
    }
    
    debug.info("Función completada, resultado: " + resultado)
    return resultado
}

func main() {
    print("=== DEBUGGING DEMO ===")
    
    let resultado1 = funcionProblematica(5)
    print("Resultado 1: " + resultado1)
    
    let resultado2 = funcionProblematica(-3)
    print("Resultado 2: " + resultado2)
    
    let resultado3 = funcionProblematica(10)
    print("Resultado 3: " + resultado3)
}
```

### 2. Profiling de Rendimiento

```r2
class Profiler {
    let mediciones
    let activo
    
    constructor() {
        this.mediciones = []
        this.activo = true
    }
    
    marcarInicio(nombre) {
        if (this.activo) {
            let medicion = {
                nombre: nombre,
                inicio: os.time(),
                fin: null,
                duracion: null
            }
            this.mediciones = this.mediciones.push(medicion)
            return this.mediciones.length() - 1  // Retornar índice
        }
        return -1
    }
    
    marcarFin(indice) {
        if (this.activo && indice >= 0 && indice < this.mediciones.length()) {
            let medicion = this.mediciones[indice]
            medicion.fin = os.time()
            medicion.duracion = medicion.fin - medicion.inicio
        }
    }
    
    generarReporte() {
        print("=== REPORTE DE PROFILING ===")
        print("Total de mediciones: " + this.mediciones.length())
        print()
        
        let totalTiempo = 0
        for (let i = 0; i < this.mediciones.length(); i++) {
            let medicion = this.mediciones[i]
            if (medicion.duracion != null) {
                print(medicion.nombre + ": " + medicion.duracion + "ms")
                totalTiempo = totalTiempo + medicion.duracion
            }
        }
        
        print()
        print("Tiempo total: " + totalTiempo + "ms")
        
        // Encontrar la operación más lenta
        let masLenta = null
        let tiempoMasLento = 0
        for (let i = 0; i < this.mediciones.length(); i++) {
            let medicion = this.mediciones[i]
            if (medicion.duracion != null && medicion.duracion > tiempoMasLento) {
                masLenta = medicion
                tiempoMasLento = medicion.duracion
            }
        }
        
        if (masLenta != null) {
            print("Operación más lenta: " + masLenta.nombre + " (" + masLenta.duracion + "ms)")
        }
    }
}

func operacionRapida() {
    let suma = 0
    for (let i = 0; i < 1000; i++) {
        suma = suma + i
    }
    return suma
}

func operacionLenta() {
    let suma = 0
    for (let i = 0; i < 100000; i++) {
        suma = suma + i
        if (i % 10000 == 0) {
            sleep(10)  // Simular operación lenta
        }
    }
    return suma
}

func operacionConArchivos() {
    let contenido = ""
    for (let i = 0; i < 100; i++) {
        contenido = contenido + "Línea " + i + "\n"
    }
    
    io.writeFile("temp_profile.txt", contenido)
    let leido = io.readFile("temp_profile.txt")
    io.deleteFile("temp_profile.txt")
    
    return leido.length()
}

func main() {
    let profiler = Profiler()
    
    print("Iniciando profiling...")
    
    // Medir operación rápida
    let idx1 = profiler.marcarInicio("Operación Rápida")
    let resultado1 = operacionRapida()
    profiler.marcarFin(idx1)
    
    // Medir operación lenta
    let idx2 = profiler.marcarInicio("Operación Lenta")
    let resultado2 = operacionLenta()
    profiler.marcarFin(idx2)
    
    // Medir operación con archivos
    let idx3 = profiler.marcarInicio("Operación con Archivos")
    let resultado3 = operacionConArchivos()
    profiler.marcarFin(idx3)
    
    // Medir operación repetitiva
    let idx4 = profiler.marcarInicio("Operación Repetitiva")
    for (let i = 0; i < 10; i++) {
        operacionRapida()
    }
    profiler.marcarFin(idx4)
    
    // Generar reporte
    profiler.generarReporte()
    
    print()
    print("Resultados:")
    print("Operación rápida: " + resultado1)
    print("Operación lenta: " + resultado2)
    print("Operación archivos: " + resultado3)
}
```

### 3. Tests de Carga y Estrés

```r2
class TestCarga {
    let url
    let concurrencia
    let totalRequests
    let resultados
    
    constructor(url, concurrencia, totalRequests) {
        this.url = url
        this.concurrencia = concurrencia
        this.totalRequests = totalRequests
        this.resultados = []
    }
    
    worker(workerId, requestsPorWorker) {
        print("Worker " + workerId + " iniciado - " + requestsPorWorker + " requests")
        
        let resultadosWorker = []
        for (let i = 0; i < requestsPorWorker; i++) {
            let inicio = os.time()
            
            try {
                let respuesta = http.get(this.url)
                let fin = os.time()
                
                resultadosWorker = resultadosWorker.push({
                    worker: workerId,
                    request: i + 1,
                    status: respuesta.status,
                    tiempo: fin - inicio,
                    exito: respuesta.status == 200
                })
            } catch (error) {
                let fin = os.time()
                resultadosWorker = resultadosWorker.push({
                    worker: workerId,
                    request: i + 1,
                    status: 0,
                    tiempo: fin - inicio,
                    exito: false,
                    error: error
                })
            }
        }
        
        // Agregar resultados al array principal
        for (let i = 0; i < resultadosWorker.length(); i++) {
            this.resultados = this.resultados.push(resultadosWorker[i])
        }
        
        print("Worker " + workerId + " completado")
    }
    
    ejecutar() {
        print("=== INICIANDO TEST DE CARGA ===")
        print("URL: " + this.url)
        print("Concurrencia: " + this.concurrencia + " workers")
        print("Total requests: " + this.totalRequests)
        print()
        
        let requestsPorWorker = math.ceil(this.totalRequests / this.concurrencia)
        let inicioTotal = os.time()
        
        // Crear workers
        for (let i = 1; i <= this.concurrencia; i++) {
            r2(this.worker, i, requestsPorWorker)
        }
        
        // Esperar que terminen todos los workers
        while (this.resultados.length() < this.totalRequests) {
            sleep(100)
        }
        
        let finTotal = os.time()
        let tiempoTotal = finTotal - inicioTotal
        
        this.generarReporte(tiempoTotal)
    }
    
    generarReporte(tiempoTotal) {
        print("=== REPORTE DE TEST DE CARGA ===")
        print("Tiempo total: " + tiempoTotal + "ms")
        print("Requests completados: " + this.resultados.length())
        
        // Calcular estadísticas
        let exitosos = 0
        let fallidos = 0
        let tiempoTotal = 0
        let tiempoMin = 999999
        let tiempoMax = 0
        
        for (let i = 0; i < this.resultados.length(); i++) {
            let resultado = this.resultados[i]
            
            if (resultado.exito) {
                exitosos++
            } else {
                fallidos++
            }
            
            tiempoTotal = tiempoTotal + resultado.tiempo
            
            if (resultado.tiempo < tiempoMin) {
                tiempoMin = resultado.tiempo
            }
            if (resultado.tiempo > tiempoMax) {
                tiempoMax = resultado.tiempo
            }
        }
        
        let tiempoPromedio = tiempoTotal / this.resultados.length()
        let tasaExito = (exitosos / this.resultados.length()) * 100
        let requestsPorSegundo = this.resultados.length() / (tiempoTotal / 1000)
        
        print("Requests exitosos: " + exitosos + " (" + tasaExito + "%)")
        print("Requests fallidos: " + fallidos)
        print("Tiempo promedio: " + tiempoPromedio + "ms")
        print("Tiempo mínimo: " + tiempoMin + "ms")
        print("Tiempo máximo: " + tiempoMax + "ms")
        print("Requests por segundo: " + requestsPorSegundo)
        
        // Distribución de códigos de respuesta
        let codigosRespuesta = {}
        for (let i = 0; i < this.resultados.length(); i++) {
            let status = this.resultados[i].status
            if (codigosRespuesta[status] == null) {
                codigosRespuesta[status] = 0
            }
            codigosRespuesta[status]++
        }
        
        print()
        print("Códigos de respuesta:")
        print("200: " + (codigosRespuesta[200] || 0))
        print("404: " + (codigosRespuesta[404] || 0))
        print("500: " + (codigosRespuesta[500] || 0))
        print("0 (error): " + (codigosRespuesta[0] || 0))
    }
}

func main() {
    // Crear servidor de prueba
    let servidor = http.server(8080, func(request, response) {
        // Simular diferentes tipos de respuesta
        let random = rand.int(1, 10)
        
        if (random <= 8) {
            // 80% respuestas exitosas
            response.writeHead(200, {"Content-Type": "application/json"})
            response.write('{"status": "ok", "timestamp": ' + os.time() + '}')
        } else if (random <= 9) {
            // 10% respuestas 404
            response.writeHead(404, {"Content-Type": "application/json"})
            response.write('{"error": "Not found"}')
        } else {
            // 10% respuestas 500
            response.writeHead(500, {"Content-Type": "application/json"})
            response.write('{"error": "Internal server error"}')
        }
        
        response.end()
        
        // Simular tiempo de procesamiento variable
        sleep(rand.int(10, 100))
    })
    
    print("Servidor de prueba iniciado en puerto 8080")
    sleep(2000)  // Esperar que el servidor esté listo
    
    // Ejecutar test de carga
    let testCarga = TestCarga("http://localhost:8080", 5, 100)
    testCarga.ejecutar()
}
```

## Proyecto del Módulo: Suite de Testing Completa

```r2
class TestSuite {
    let tests
    let resultados
    let configuracion
    
    constructor() {
        this.tests = []
        this.resultados = []
        this.configuracion = {
            timeoutMs: 5000,
            paralelo: false,
            reporteDetallado: true,
            guardarReporte: true
        }
    }
    
    agregarTest(nombre, funcion) {
        this.tests = this.tests.push({
            nombre: nombre,
            funcion: funcion,
            timeout: this.configuracion.timeoutMs
        })
    }
    
    ejecutarTest(test) {
        let inicio = os.time()
        let resultado = {
            nombre: test.nombre,
            exito: false,
            error: null,
            duracion: 0,
            assertions: 0
        }
        
        try {
            // Ejecutar el test
            test.funcion()
            
            resultado.exito = true
            resultado.duracion = os.time() - inicio
            
            if (this.configuracion.reporteDetallado) {
                print("✓ " + test.nombre + " (" + resultado.duracion + "ms)")
            }
        } catch (error) {
            resultado.exito = false
            resultado.error = error
            resultado.duracion = os.time() - inicio
            
            if (this.configuracion.reporteDetallado) {
                print("✗ " + test.nombre + " - ERROR: " + error)
            }
        }
        
        this.resultados = this.resultados.push(resultado)
        return resultado
    }
    
    ejecutarTodos() {
        print("=== EJECUTANDO SUITE DE TESTS ===")
        print("Total de tests: " + this.tests.length())
        print()
        
        let inicioSuite = os.time()
        
        if (this.configuracion.paralelo) {
            // Ejecutar tests en paralelo
            for (let i = 0; i < this.tests.length(); i++) {
                let test = this.tests[i]
                r2(this.ejecutarTest, test)
            }
            
            // Esperar que terminen todos
            while (this.resultados.length() < this.tests.length()) {
                sleep(100)
            }
        } else {
            // Ejecutar tests secuencialmente
            for (let i = 0; i < this.tests.length(); i++) {
                this.ejecutarTest(this.tests[i])
            }
        }
        
        let finSuite = os.time()
        let duracionTotal = finSuite - inicioSuite
        
        this.generarReporte(duracionTotal)
    }
    
    generarReporte(duracionTotal) {
        print()
        print("=== REPORTE DE TESTS ===")
        
        let exitosos = 0
        let fallidos = 0
        let tiempoTotal = 0
        
        for (let i = 0; i < this.resultados.length(); i++) {
            let resultado = this.resultados[i]
            tiempoTotal = tiempoTotal + resultado.duracion
            
            if (resultado.exito) {
                exitosos++
            } else {
                fallidos++
            }
        }
        
        let porcentajeExito = (exitosos / this.resultados.length()) * 100
        
        print("Tests ejecutados: " + this.resultados.length())
        print("Tests exitosos: " + exitosos + " (" + porcentajeExito + "%)")
        print("Tests fallidos: " + fallidos)
        print("Tiempo total: " + duracionTotal + "ms")
        print("Tiempo promedio por test: " + (tiempoTotal / this.resultados.length()) + "ms")
        
        if (fallidos > 0) {
            print()
            print("TESTS FALLIDOS:")
            for (let i = 0; i < this.resultados.length(); i++) {
                let resultado = this.resultados[i]
                if (!resultado.exito) {
                    print("- " + resultado.nombre + ": " + resultado.error)
                }
            }
        }
        
        if (this.configuracion.guardarReporte) {
            this.guardarReporteArchivo(duracionTotal)
        }
    }
    
    guardarReporteArchivo(duracionTotal) {
        let reporte = "=== REPORTE DE TESTS ===\n"
        reporte = reporte + "Fecha: " + os.time() + "\n"
        reporte = reporte + "Duración total: " + duracionTotal + "ms\n\n"
        
        for (let i = 0; i < this.resultados.length(); i++) {
            let resultado = this.resultados[i]
            let status = resultado.exito ? "PASS" : "FAIL"
            reporte = reporte + "[" + status + "] " + resultado.nombre + " (" + resultado.duracion + "ms)\n"
            
            if (!resultado.exito) {
                reporte = reporte + "  Error: " + resultado.error + "\n"
            }
        }
        
        let nombreArchivo = "test_report_" + os.time() + ".txt"
        io.writeFile(nombreArchivo, reporte)
        print("Reporte guardado en: " + nombreArchivo)
    }
}

// Funciones de testing helper
func assertEqual(actual, esperado) {
    if (actual != esperado) {
        throw "Assertion failed: expected " + esperado + " but got " + actual
    }
}

func assertTrue(valor) {
    if (!valor) {
        throw "Assertion failed: expected true but got " + valor
    }
}

func assertFalse(valor) {
    if (valor) {
        throw "Assertion failed: expected false but got " + valor
    }
}

func assertNull(valor) {
    if (valor != null) {
        throw "Assertion failed: expected null but got " + valor
    }
}

func assertNotNull(valor) {
    if (valor == null) {
        throw "Assertion failed: expected non-null value"
    }
}

// Tests de ejemplo
func testCalculadoraBasica() {
    let calc = Calculadora()
    
    assertEqual(calc.sumar(5, 3), 8)
    assertEqual(calc.restar(10, 4), 6)
    assertEqual(calc.multiplicar(6, 7), 42)
    assertEqual(calc.dividir(15, 3), 5)
}

func testCalculadoraErrores() {
    let calc = Calculadora()
    
    let resultado = calc.dividir(10, 0)
    assertNull(resultado)
    assertTrue(calc.hayError)
}

func testStringOperaciones() {
    let texto = "Hola mundo"
    
    assertEqual(texto.length(), 10)
    assertTrue(texto.contains("mundo"))
    assertFalse(texto.contains("xyz"))
    assertEqual(texto.upper(), "HOLA MUNDO")
}

func testArrayOperaciones() {
    let arr = [1, 2, 3, 4, 5]
    
    assertEqual(arr.length(), 5)
    assertEqual(arr[0], 1)
    assertEqual(arr[4], 5)
    
    arr = arr.push(6)
    assertEqual(arr.length(), 6)
    assertEqual(arr[5], 6)
}

func testHTTPCliente() {
    // Test que podría fallar por red
    let respuesta = http.get("http://httpbin.org/get")
    assertEqual(respuesta.status, 200)
    assertTrue(respuesta.body.contains("httpbin"))
}

class Calculadora {
    let hayError
    
    constructor() {
        this.hayError = false
    }
    
    sumar(a, b) {
        this.hayError = false
        return a + b
    }
    
    restar(a, b) {
        this.hayError = false
        return a - b
    }
    
    multiplicar(a, b) {
        this.hayError = false
        return a * b
    }
    
    dividir(a, b) {
        if (b == 0) {
            this.hayError = true
            return null
        }
        this.hayError = false
        return a / b
    }
}

func main() {
    let suite = TestSuite()
    
    // Configurar suite
    suite.configuracion.reporteDetallado = true
    suite.configuracion.paralelo = false
    suite.configuracion.guardarReporte = true
    
    // Agregar tests
    suite.agregarTest("Calculadora Básica", testCalculadoraBasica)
    suite.agregarTest("Calculadora Errores", testCalculadoraErrores)
    suite.agregarTest("String Operaciones", testStringOperaciones)
    suite.agregarTest("Array Operaciones", testArrayOperaciones)
    suite.agregarTest("HTTP Cliente", testHTTPCliente)
    
    // Ejecutar todos los tests
    suite.ejecutarTodos()
}
```

## Resumen del Módulo

### Técnicas de Testing Aprendidas
- ✅ **BDD Testing**: Sintaxis Given/When/Then/And
- ✅ **Unit Testing**: Tests de componentes individuales
- ✅ **Integration Testing**: Tests de sistemas completos
- ✅ **Load Testing**: Tests de carga y rendimiento
- ✅ **Debugging**: Técnicas de depuración avanzadas
- ✅ **Profiling**: Análisis de rendimiento
- ✅ **Test Suites**: Organización y ejecución de tests

### Habilidades Desarrolladas
- ✅ Escribir tests robustos y mantenibles
- ✅ Implementar debugging efectivo
- ✅ Realizar profiling de rendimiento
- ✅ Crear tests de integración complejos
- ✅ Desarrollar herramientas de testing personalizadas
- ✅ Optimizar código basado en métricas
- ✅ Mantener alta calidad de código

### Próximo Módulo

En el **Módulo 10** aprenderás:
- Despliegue y distribución de aplicaciones
- Empaquetado y versionado
- Optimización para producción
- Monitoreo y logging en producción

¡Excelente trabajo! Ahora puedes crear aplicaciones R2Lang robustas con testing comprehensivo y debugging avanzado.