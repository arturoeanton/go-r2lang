# Curso R2Lang - Módulo 5: Testing y Desarrollo Web

## Introducción

En este módulo aprenderás una de las características más distintivas de R2Lang: su sistema de testing integrado con sintaxis BDD (Behavior Driven Development), además de desarrollo web y APIs. Estas habilidades te permitirán crear aplicaciones completas y bien probadas.

### Sistema de Testing y Web v2

La nueva arquitectura modular proporciona un sistema de testing y desarrollo web completamente renovado:

```
Testing & Web Framework (pkg/r2libs/):
├── r2test.go (testing framework)     # Sistema BDD completo
├── r2http.go (410 LOC)              # Servidor HTTP optimizado
├── r2httpclient.go (324 LOC)        # Cliente HTTP robusto
├── r2print.go (365 LOC)             # Output y debugging
└── Testing AST (pkg/r2core/):
    ├── testcase_statement.go         # Sintaxis TestCase nativa
    ├── given_when_then.go            # Pasos BDD
    └── assertion_functions.go        # Funciones de aserción
```

**Mejoras v2**:
- **Testing nativo**: BDD integrado en el lenguaje, no como biblioteca externa
- **HTTP framework**: Servidor y cliente HTTP optimizados con mejor rendimiento
- **Debugging avanzado**: Herramientas integradas para debugging de tests
- **Cobertura automática**: Reportes de cobertura de tests incluidos
- **Paralelización**: Tests pueden ejecutarse en paralelo automáticamente

## Sistema de Testing BDD Integrado

### 1. Conceptos de BDD (Behavior Driven Development)

BDD es una metodología que describe el comportamiento del software en un lenguaje natural estructurado:

- **Given** (Dado): Establece el contexto inicial
- **When** (Cuando): Describe la acción que se ejecuta  
- **Then** (Entonces): Verifica el resultado esperado
- **And** (Y): Continúa con el paso anterior

### 2. Tu Primer TestCase

```r2
// Funciones de soporte para testing
func assertEqual(actual, expected) {
    if (actual == expected) {
        print("✅ PASS: Valor esperado recibido")
        return true
    } else {
        print("❌ FAIL: Esperado", expected, "pero recibido", actual)
        return false
    }
}

func assertTrue(condition) {
    if (condition) {
        print("✅ PASS: Condición verdadera")
        return true
    } else {
        print("❌ FAIL: Condición falsa")
        return false
    }
}

// Función a probar
func sumar(a, b) {
    return a + b
}

TestCase "Verificar suma de números" {
    Given func() {
        print("Preparando datos para la suma")
        return "Datos preparados"
    }
    
    When func() {
        let resultado = sumar(2, 3)
        return "Suma ejecutada: " + resultado
    }
    
    Then func() {
        let resultado = sumar(2, 3)
        assertEqual(resultado, 5)
        return "Validación completada"
    }
}

func main() {
    print("Ejecutando tests...")
}
```

### 3. TestCase Avanzado con Setup y Teardown

```r2
// Simulación de base de datos
let baseDatos = []

func limpiarBaseDatos() {
    baseDatos = []
    print("🗑️ Base de datos limpiada")
}

func agregarUsuario(nombre, email) {
    let usuario = {
        id: baseDatos.length() + 1,
        nombre: nombre,
        email: email,
        activo: true
    }
    baseDatos = baseDatos.push(usuario)
    return usuario
}

func buscarUsuario(email) {
    for (let i = 0; i < baseDatos.length(); i++) {
        let usuario = baseDatos[i]
        if (usuario.email == email) {
            return usuario
        }
    }
    return null
}

func desactivarUsuario(email) {
    for (let i = 0; i < baseDatos.length(); i++) {
        let usuario = baseDatos[i]
        if (usuario.email == email) {
            usuario.activo = false
            return true
        }
    }
    return false
}

TestCase "Gestión completa de usuarios" {
    Given func() {
        limpiarBaseDatos()
        print("Base de datos preparada para testing")
        return "Setup completado"
    }
    
    When func() {
        let usuario = agregarUsuario("Ana García", "ana@email.com")
        print("Usuario creado con ID:", usuario.id)
        return "Usuario creado exitosamente"
    }
    
    Then func() {
        let usuario = buscarUsuario("ana@email.com")
        assertTrue(usuario != null)
        assertEqual(usuario.nombre, "Ana García")
        assertTrue(usuario.activo)
        return "Usuario encontrado y validado"
    }
    
    And func() {
        let resultado = desactivarUsuario("ana@email.com")
        assertTrue(resultado)
        
        let usuario = buscarUsuario("ana@email.com")
        assertTrue(!usuario.activo)
        return "Usuario desactivado correctamente"
    }
}

TestCase "Búsqueda de usuario inexistente" {
    Given func() {
        limpiarBaseDatos()
        return "Base de datos vacía"
    }
    
    When func() {
        let usuario = buscarUsuario("inexistente@email.com")
        return "Búsqueda ejecutada"
    }
    
    Then func() {
        let usuario = buscarUsuario("inexistente@email.com")
        assertTrue(usuario == null)
        return "Usuario no encontrado como esperado"
    }
}

func main() {
    print("=== SUITE DE TESTS DE USUARIOS ===")
}
```

### 4. Testing de Clases y Objetos

```r2
class CalculadoraBancaria {
    let saldo
    let historial
    
    constructor(saldoInicial) {
        this.saldo = saldoInicial
        this.historial = []
    }
    
    depositar(monto) {
        if (monto <= 0) {
            throw "Monto debe ser positivo"
        }
        
        this.saldo = this.saldo + monto
        this.historial = this.historial.push({
            tipo: "Depósito",
            monto: monto,
            saldoResultante: this.saldo
        })
        
        return this.saldo
    }
    
    retirar(monto) {
        if (monto <= 0) {
            throw "Monto debe ser positivo"
        }
        
        if (monto > this.saldo) {
            throw "Saldo insuficiente"
        }
        
        this.saldo = this.saldo - monto
        this.historial = this.historial.push({
            tipo: "Retiro",
            monto: monto,
            saldoResultante: this.saldo
        })
        
        return this.saldo
    }
    
    obtenerSaldo() {
        return this.saldo
    }
    
    obtenerHistorial() {
        return this.historial
    }
}

// Variable global para tests
let calculadora

TestCase "Operaciones básicas de cuenta bancaria" {
    Given func() {
        calculadora = CalculadoraBancaria(1000)
        print("Calculadora inicializada con saldo:", calculadora.obtenerSaldo())
        return "Calculadora lista"
    }
    
    When func() {
        calculadora.depositar(500)
        print("Depósito de 500 realizado")
        return "Depósito completado"
    }
    
    Then func() {
        assertEqual(calculadora.obtenerSaldo(), 1500)
        let historial = calculadora.obtenerHistorial()
        assertEqual(historial.length(), 1)
        assertEqual(historial[0].tipo, "Depósito")
        return "Depósito validado correctamente"
    }
    
    And func() {
        calculadora.retirar(300)
        assertEqual(calculadora.obtenerSaldo(), 1200)
        let historial = calculadora.obtenerHistorial()
        assertEqual(historial.length(), 2)
        return "Retiro validado correctamente"
    }
}

TestCase "Manejo de errores en operaciones bancarias" {
    Given func() {
        calculadora = CalculadoraBancaria(100)
        return "Calculadora con saldo bajo inicializada"
    }
    
    When func() {
        print("Intentando retirar más dinero del disponible")
        return "Intento de retiro excesivo"
    }
    
    Then func() {
        try {
            calculadora.retirar(200)
            assertTrue(false)  // No debería llegar aquí
        } catch (error) {
            assertTrue(error.contains("insuficiente"))
            assertEqual(calculadora.obtenerSaldo(), 100)  // Saldo no cambió
        }
        return "Error manejado correctamente"
    }
    
    And func() {
        try {
            calculadora.depositar(-50)
            assertTrue(false)  // No debería llegar aquí
        } catch (error) {
            assertTrue(error.contains("positivo"))
        }
        return "Validación de monto negativo correcta"
    }
}

func main() {
    print("=== TESTS DE CALCULADORA BANCARIA ===")
}
```

### 5. Testing de Funciones Concurrentes

```r2
let contadorGlobal = 0
let resultadosConcurrentes = []

func incrementarContador(veces, id) {
    for (let i = 0; i < veces; i++) {
        contadorGlobal++
        resultadosConcurrentes = resultadosConcurrentes.push({
            worker: id,
            valor: contadorGlobal
        })
        sleep(0.1)  // Simular trabajo
    }
}

func resetearContadores() {
    contadorGlobal = 0
    resultadosConcurrentes = []
}

TestCase "Verificar comportamiento concurrente" {
    Given func() {
        resetearContadores()
        print("Contadores reseteados")
        return "Estado inicial limpio"
    }
    
    When func() {
        r2(incrementarContador, 3, "Worker-1")
        r2(incrementarContador, 3, "Worker-2")
        r2(incrementarContador, 3, "Worker-3")
        
        sleep(2)  // Esperar a que terminen
        print("Operaciones concurrentes completadas")
        return "Incrementos concurrentes ejecutados"
    }
    
    Then func() {
        assertEqual(contadorGlobal, 9)
        assertEqual(resultadosConcurrentes.length(), 9)
        print("Valores finales validados")
        return "Concurrencia verificada"
    }
    
    And func() {
        // Verificar que hubo intercalado (no determinístico, pero probable)
        let workers = []
        for (let i = 0; i < resultadosConcurrentes.length(); i++) {
            let resultado = resultadosConcurrentes[i]
            workers = workers.push(resultado.worker)
        }
        
        print("Secuencia de workers:", workers)
        // En ejecución real, deberíamos ver workers intercalados
        return "Patrón de intercalado observado"
    }
}

func main() {
    print("=== TESTS DE CONCURRENCIA ===")
}
```

## Desarrollo Web y APIs

### 1. Servidor HTTP Básico

```r2
func manejarRaiz(req, res) {
    res.send("¡Hola desde R2Lang!")
}

func manejarSaludo(req, res) {
    let nombre = req.query.nombre || "Anónimo"
    res.send("¡Hola " + nombre + "!")
}

func manejarInfo(req, res) {
    let info = {
        servidor: "R2Lang HTTP Server",
        version: "1.0",
        timestamp: "2024-01-01",
        endpoints: ["/", "/saludo", "/info"]
    }
    res.json(info)
}

func main() {
    print("Iniciando servidor web en puerto 8080...")
    
    // Configurar rutas
    http.get("/", manejarRaiz)
    http.get("/saludo", manejarSaludo)
    http.get("/info", manejarInfo)
    
    // Iniciar servidor
    http.listen(8080)
}
```

### 2. API REST Completa

```r2
// Simulación de base de datos en memoria
let usuarios = []
let siguienteId = 1

// Funciones de utilidad
func generarId() {
    let id = siguienteId
    siguienteId++
    return id
}

func encontrarUsuario(id) {
    for (let i = 0; i < usuarios.length(); i++) {
        let usuario = usuarios[i]
        if (usuario.id == id) {
            return { usuario: usuario, indice: i }
        }
    }
    return null
}

func validarUsuario(datos) {
    if (!datos.nombre || datos.nombre == "") {
        return "Nombre es requerido"
    }
    
    if (!datos.email || datos.email == "") {
        return "Email es requerido"
    }
    
    // Verificar email único
    for (let i = 0; i < usuarios.length(); i++) {
        if (usuarios[i].email == datos.email) {
            return "Email ya está en uso"
        }
    }
    
    return null
}

// Handlers de la API
func obtenerTodosLosUsuarios(req, res) {
    res.json({
        usuarios: usuarios,
        total: usuarios.length()
    })
}

func obtenerUsuario(req, res) {
    let id = parseInt(req.params.id)
    let resultado = encontrarUsuario(id)
    
    if (resultado == null) {
        res.status(404).json({
            error: "Usuario no encontrado",
            id: id
        })
        return
    }
    
    res.json(resultado.usuario)
}

func crearUsuario(req, res) {
    let datos = req.body
    
    // Validar datos
    let error = validarUsuario(datos)
    if (error != null) {
        res.status(400).json({
            error: error
        })
        return
    }
    
    // Crear usuario
    let nuevoUsuario = {
        id: generarId(),
        nombre: datos.nombre,
        email: datos.email,
        activo: true,
        fechaCreacion: "2024-01-01"
    }
    
    usuarios = usuarios.push(nuevoUsuario)
    
    res.status(201).json(nuevoUsuario)
}

func actualizarUsuario(req, res) {
    let id = parseInt(req.params.id)
    let datos = req.body
    let resultado = encontrarUsuario(id)
    
    if (resultado == null) {
        res.status(404).json({
            error: "Usuario no encontrado"
        })
        return
    }
    
    // Actualizar campos
    let usuario = resultado.usuario
    if (datos.nombre) {
        usuario.nombre = datos.nombre
    }
    if (datos.email) {
        usuario.email = datos.email
    }
    if (datos.activo != null) {
        usuario.activo = datos.activo
    }
    
    res.json(usuario)
}

func eliminarUsuario(req, res) {
    let id = parseInt(req.params.id)
    let resultado = encontrarUsuario(id)
    
    if (resultado == null) {
        res.status(404).json({
            error: "Usuario no encontrado"
        })
        return
    }
    
    // Eliminar usuario (simular removiendo del array)
    let usuarioEliminado = resultado.usuario
    
    // En R2Lang actual no tenemos método remove, así que simularemos
    let nuevosUsuarios = []
    for (let i = 0; i < usuarios.length(); i++) {
        if (usuarios[i].id != id) {
            nuevosUsuarios = nuevosUsuarios.push(usuarios[i])
        }
    }
    usuarios = nuevosUsuarios
    
    res.json({
        mensaje: "Usuario eliminado",
        usuario: usuarioEliminado
    })
}

// Middleware de logging
func middleware(req, res, next) {
    print("📥", req.method, req.url, "- IP:", req.ip)
    next()
}

func main() {
    print("🚀 Iniciando API REST en puerto 3000...")
    
    // Aplicar middleware
    http.use(middleware)
    
    // Configurar rutas REST
    http.get("/api/usuarios", obtenerTodosLosUsuarios)
    http.get("/api/usuarios/:id", obtenerUsuario)
    http.post("/api/usuarios", crearUsuario)
    http.put("/api/usuarios/:id", actualizarUsuario)
    http.delete("/api/usuarios/:id", eliminarUsuario)
    
    // Ruta de salud
    http.get("/health", func(req, res) {
        res.json({
            status: "OK",
            timestamp: "2024-01-01",
            uptime: "5 minutes"
        })
    })
    
    // Iniciar servidor
    http.listen(3000)
    print("✅ API REST disponible en http://localhost:3000")
    print("📋 Endpoints disponibles:")
    print("  GET    /api/usuarios")
    print("  GET    /api/usuarios/:id")
    print("  POST   /api/usuarios")
    print("  PUT    /api/usuarios/:id")
    print("  DELETE /api/usuarios/:id")
    print("  GET    /health")
}
```

### 3. Testing de APIs

```r2
// Simulación de cliente HTTP para testing
let baseURL = "http://localhost:3000"

func hacerRequest(method, url, body) {
    // Simulación de request HTTP
    print("📤", method, baseURL + url)
    if (body) {
        print("📦 Body:", body)
    }
    
    // Simular respuesta según el endpoint
    if (method == "GET" && url == "/api/usuarios") {
        return {
            status: 200,
            body: { usuarios: [], total: 0 }
        }
    }
    
    if (method == "POST" && url == "/api/usuarios") {
        return {
            status: 201,
            body: { 
                id: 1, 
                nombre: body.nombre, 
                email: body.email,
                activo: true 
            }
        }
    }
    
    return {
        status: 200,
        body: { mensaje: "Respuesta simulada" }
    }
}

TestCase "API REST - Crear usuario" {
    Given func() {
        print("Preparando datos para crear usuario")
        return "Cliente HTTP listo"
    }
    
    When func() {
        let nuevoUsuario = {
            nombre: "Juan Pérez",
            email: "juan@email.com"
        }
        
        let response = hacerRequest("POST", "/api/usuarios", nuevoUsuario)
        print("Usuario creado con respuesta:", response.status)
        return "Usuario creado via API"
    }
    
    Then func() {
        let nuevoUsuario = {
            nombre: "Juan Pérez",
            email: "juan@email.com"
        }
        
        let response = hacerRequest("POST", "/api/usuarios", nuevoUsuario)
        
        assertEqual(response.status, 201)
        assertTrue(response.body.id != null)
        assertEqual(response.body.nombre, "Juan Pérez")
        assertEqual(response.body.email, "juan@email.com")
        assertTrue(response.body.activo)
        
        return "Respuesta de creación validada"
    }
}

TestCase "API REST - Obtener lista vacía" {
    Given func() {
        print("API limpia sin usuarios")
        return "Estado inicial"
    }
    
    When func() {
        let response = hacerRequest("GET", "/api/usuarios", null)
        return "Lista de usuarios obtenida"
    }
    
    Then func() {
        let response = hacerRequest("GET", "/api/usuarios", null)
        
        assertEqual(response.status, 200)
        assertTrue(response.body.usuarios != null)
        assertEqual(response.body.total, 0)
        
        return "Lista vacía validada"
    }
}

func main() {
    print("=== TESTS DE API REST ===")
}
```

## Integración de Testing y Web

### 1. Testing End-to-End

```r2
// Simulación de aplicación web completa
class AplicacionWeb {
    let usuarios
    let sesiones
    
    constructor() {
        this.usuarios = []
        this.sesiones = []
    }
    
    registrarUsuario(datos) {
        // Validar datos
        if (!datos.nombre || !datos.email || !datos.password) {
            throw "Datos incompletos"
        }
        
        // Verificar email único
        for (let i = 0; i < this.usuarios.length(); i++) {
            if (this.usuarios[i].email == datos.email) {
                throw "Email ya registrado"
            }
        }
        
        let usuario = {
            id: this.usuarios.length() + 1,
            nombre: datos.nombre,
            email: datos.email,
            password: datos.password,  // En prod esto debería estar hasheado
            activo: true
        }
        
        this.usuarios = this.usuarios.push(usuario)
        return usuario
    }
    
    iniciarSesion(email, password) {
        for (let i = 0; i < this.usuarios.length(); i++) {
            let usuario = this.usuarios[i]
            if (usuario.email == email && usuario.password == password) {
                let sesion = {
                    token: "token_" + usuario.id + "_" + this.sesiones.length(),
                    usuarioId: usuario.id,
                    activa: true
                }
                
                this.sesiones = this.sesiones.push(sesion)
                return sesion
            }
        }
        
        throw "Credenciales inválidas"
    }
    
    cerrarSesion(token) {
        for (let i = 0; i < this.sesiones.length(); i++) {
            let sesion = this.sesiones[i]
            if (sesion.token == token) {
                sesion.activa = false
                return true
            }
        }
        return false
    }
    
    obtenerUsuarioActual(token) {
        for (let i = 0; i < this.sesiones.length(); i++) {
            let sesion = this.sesiones[i]
            if (sesion.token == token && sesion.activa) {
                for (let j = 0; j < this.usuarios.length(); j++) {
                    let usuario = this.usuarios[j]
                    if (usuario.id == sesion.usuarioId) {
                        return usuario
                    }
                }
            }
        }
        return null
    }
}

let app
let usuarioTest
let sesionTest

TestCase "Flujo completo de usuario - Registro y Login" {
    Given func() {
        app = AplicacionWeb()
        usuarioTest = {
            nombre: "María González",
            email: "maria@test.com",
            password: "password123"
        }
        print("Aplicación web inicializada")
        return "App lista para testing"
    }
    
    When func() {
        let usuario = app.registrarUsuario(usuarioTest)
        print("Usuario registrado:", usuario.nombre)
        return "Registro completado"
    }
    
    Then func() {
        let usuario = app.registrarUsuario(usuarioTest)
        assertTrue(usuario.id != null)
        assertEqual(usuario.nombre, "María González")
        assertEqual(usuario.email, "maria@test.com")
        assertTrue(usuario.activo)
        return "Registro validado"
    }
    
    And func() {
        let sesion = app.iniciarSesion(usuarioTest.email, usuarioTest.password)
        sesionTest = sesion
        
        assertTrue(sesion.token != null)
        assertTrue(sesion.activa)
        assertEqual(sesion.usuarioId, 1)
        
        return "Login exitoso validado"
    }
}

TestCase "Gestión de sesiones" {
    Given func() {
        print("Usando sesión activa del test anterior")
        return "Sesión disponible"
    }
    
    When func() {
        let usuario = app.obtenerUsuarioActual(sesionTest.token)
        print("Usuario actual obtenido:", usuario.nombre)
        return "Usuario de sesión obtenido"
    }
    
    Then func() {
        let usuario = app.obtenerUsuarioActual(sesionTest.token)
        assertTrue(usuario != null)
        assertEqual(usuario.email, "maria@test.com")
        return "Usuario de sesión validado"
    }
    
    And func() {
        let resultado = app.cerrarSesion(sesionTest.token)
        assertTrue(resultado)
        
        let usuarioTrasLogout = app.obtenerUsuarioActual(sesionTest.token)
        assertTrue(usuarioTrasLogout == null)
        
        return "Logout validado"
    }
}

TestCase "Manejo de errores en autenticación" {
    Given func() {
        print("Preparando casos de error")
        return "Casos de error listos"
    }
    
    When func() {
        print("Intentando registrar usuario duplicado")
        return "Intento de registro duplicado"
    }
    
    Then func() {
        try {
            app.registrarUsuario(usuarioTest)  // Mismo email
            assertTrue(false)  // No debería llegar aquí
        } catch (error) {
            assertTrue(error.contains("ya registrado"))
        }
        return "Error de email duplicado manejado"
    }
    
    And func() {
        try {
            app.iniciarSesion("noexiste@test.com", "wrongpass")
            assertTrue(false)  // No debería llegar aquí
        } catch (error) {
            assertTrue(error.contains("inválidas"))
        }
        return "Error de credenciales inválidas manejado"
    }
}

func main() {
    print("=== TESTS END-TO-END DE APLICACIÓN WEB ===")
}
```

## Proyecto del Módulo: Sistema de Blog con Testing

```r2
// Sistema completo de blog con testing BDD

class Post {
    let id
    let titulo
    let contenido
    let autor
    let fechaCreacion
    let fechaModificacion
    let activo
    let comentarios
    
    constructor(id, titulo, contenido, autor) {
        this.id = id
        this.titulo = titulo
        this.contenido = contenido
        this.autor = autor
        this.fechaCreacion = "2024-01-01"
        this.fechaModificacion = "2024-01-01"
        this.activo = true
        this.comentarios = []
    }
    
    agregarComentario(autor, contenido) {
        let comentario = {
            id: this.comentarios.length() + 1,
            autor: autor,
            contenido: contenido,
            fecha: "2024-01-01"
        }
        
        this.comentarios = this.comentarios.push(comentario)
        return comentario
    }
    
    actualizar(nuevoTitulo, nuevoContenido) {
        this.titulo = nuevoTitulo
        this.contenido = nuevoContenido
        this.fechaModificacion = "2024-01-01"
    }
    
    desactivar() {
        this.activo = false
    }
}

class BlogService {
    let posts
    let siguienteId
    
    constructor() {
        this.posts = []
        this.siguienteId = 1
    }
    
    crearPost(titulo, contenido, autor) {
        if (!titulo || !contenido || !autor) {
            throw "Título, contenido y autor son requeridos"
        }
        
        let post = Post(this.siguienteId, titulo, contenido, autor)
        this.siguienteId++
        this.posts = this.posts.push(post)
        
        return post
    }
    
    obtenerPost(id) {
        for (let i = 0; i < this.posts.length(); i++) {
            let post = this.posts[i]
            if (post.id == id && post.activo) {
                return post
            }
        }
        return null
    }
    
    obtenerPostsPorAutor(autor) {
        let posts = []
        for (let i = 0; i < this.posts.length(); i++) {
            let post = this.posts[i]
            if (post.autor == autor && post.activo) {
                posts = posts.push(post)
            }
        }
        return posts
    }
    
    buscarPosts(termino) {
        let posts = []
        for (let i = 0; i < this.posts.length(); i++) {
            let post = this.posts[i]
            if (post.activo && 
                (post.titulo.contains(termino) || post.contenido.contains(termino))) {
                posts = posts.push(post)
            }
        }
        return posts
    }
    
    eliminarPost(id) {
        let post = this.obtenerPost(id)
        if (post != null) {
            post.desactivar()
            return true
        }
        return false
    }
    
    obtenerEstadisticas() {
        let totalPosts = 0
        let totalComentarios = 0
        let autores = []
        
        for (let i = 0; i < this.posts.length(); i++) {
            let post = this.posts[i]
            if (post.activo) {
                totalPosts++
                totalComentarios = totalComentarios + post.comentarios.length()
                
                // Contar autores únicos (simplificado)
                let autorExiste = false
                for (let j = 0; j < autores.length(); j++) {
                    if (autores[j] == post.autor) {
                        autorExiste = true
                        break
                    }
                }
                if (!autorExiste) {
                    autores = autores.push(post.autor)
                }
            }
        }
        
        return {
            totalPosts: totalPosts,
            totalComentarios: totalComentarios,
            totalAutores: autores.length()
        }
    }
}

// Variables globales para testing
let blogService
let postTest

TestCase "Creación y gestión de posts" {
    Given func() {
        blogService = BlogService()
        print("Blog service inicializado")
        return "Sistema de blog listo"
    }
    
    When func() {
        postTest = blogService.crearPost(
            "Mi primer post",
            "Este es el contenido de mi primer post en R2Lang",
            "Juan Blogger"
        )
        print("Post creado con ID:", postTest.id)
        return "Post creado exitosamente"
    }
    
    Then func() {
        assertTrue(postTest != null)
        assertEqual(postTest.id, 1)
        assertEqual(postTest.titulo, "Mi primer post")
        assertEqual(postTest.autor, "Juan Blogger")
        assertTrue(postTest.activo)
        assertEqual(postTest.comentarios.length(), 0)
        return "Post validado correctamente"
    }
    
    And func() {
        let postRecuperado = blogService.obtenerPost(1)
        assertTrue(postRecuperado != null)
        assertEqual(postRecuperado.titulo, postTest.titulo)
        return "Post recuperado correctamente"
    }
}

TestCase "Sistema de comentarios" {
    Given func() {
        print("Usando post existente para comentarios")
        return "Post disponible para comentarios"
    }
    
    When func() {
        let comentario = postTest.agregarComentario(
            "Ana Lectora",
            "¡Excelente post! Me gustó mucho."
        )
        print("Comentario agregado por:", comentario.autor)
        return "Comentario agregado"
    }
    
    Then func() {
        assertEqual(postTest.comentarios.length(), 1)
        let comentario = postTest.comentarios[0]
        assertEqual(comentario.autor, "Ana Lectora")
        assertTrue(comentario.contenido.contains("Excelente"))
        return "Comentario validado"
    }
    
    And func() {
        postTest.agregarComentario("Carlos Lector", "Muy informativo")
        assertEqual(postTest.comentarios.length(), 2)
        return "Múltiples comentarios funcionando"
    }
}

TestCase "Búsqueda y filtros" {
    Given func() {
        // Crear posts adicionales para búsqueda
        blogService.crearPost(
            "Tutorial de R2Lang",
            "Aprende R2Lang desde cero",
            "María Tutora"
        )
        blogService.crearPost(
            "Programación Avanzada",
            "Técnicas avanzadas de programación",
            "Juan Blogger"
        )
        print("Posts adicionales creados para búsqueda")
        return "Dataset de posts preparado"
    }
    
    When func() {
        let postsR2Lang = blogService.buscarPosts("R2Lang")
        print("Búsqueda de 'R2Lang' encontró:", postsR2Lang.length(), "posts")
        return "Búsqueda ejecutada"
    }
    
    Then func() {
        let postsR2Lang = blogService.buscarPosts("R2Lang")
        assertEqual(postsR2Lang.length(), 2)  // "Mi primer post" y "Tutorial de R2Lang"
        return "Búsqueda por término validada"
    }
    
    And func() {
        let postsJuan = blogService.obtenerPostsPorAutor("Juan Blogger")
        assertEqual(postsJuan.length(), 2)  // "Mi primer post" y "Programación Avanzada"
        
        let postsMaria = blogService.obtenerPostsPorAutor("María Tutora")
        assertEqual(postsMaria.length(), 1)  // "Tutorial de R2Lang"
        
        return "Filtro por autor validado"
    }
}

TestCase "Estadísticas del blog" {
    Given func() {
        print("Calculando estadísticas del blog")
        return "Blog con múltiples posts y comentarios"
    }
    
    When func() {
        let stats = blogService.obtenerEstadisticas()
        print("Estadísticas calculadas:", stats)
        return "Estadísticas obtenidas"
    }
    
    Then func() {
        let stats = blogService.obtenerEstadisticas()
        assertEqual(stats.totalPosts, 3)
        assertEqual(stats.totalComentarios, 2)  // Solo el primer post tiene comentarios
        assertEqual(stats.totalAutores, 2)     // Juan Blogger y María Tutora
        return "Estadísticas validadas"
    }
}

TestCase "Eliminación de posts" {
    Given func() {
        print("Preparando eliminación de post")
        return "Posts disponibles para eliminar"
    }
    
    When func() {
        let resultado = blogService.eliminarPost(1)
        print("Post 1 eliminado:", resultado)
        return "Eliminación ejecutada"
    }
    
    Then func() {
        let resultado = blogService.eliminarPost(1)
        assertTrue(resultado)
        
        let postEliminado = blogService.obtenerPost(1)
        assertTrue(postEliminado == null)  // No debería encontrarlo
        
        return "Eliminación validada"
    }
    
    And func() {
        let statsActualizadas = blogService.obtenerEstadisticas()
        assertEqual(statsActualizadas.totalPosts, 2)  // Un post menos
        return "Estadísticas actualizadas tras eliminación"
    }
}

func main() {
    print("=== SUITE COMPLETA DE TESTS DEL BLOG ===")
    print("Ejecutando tests BDD...")
}
```

## Mejores Prácticas

### 1. Testing BDD
- ✅ Usar nombres descriptivos en TestCase
- ✅ Mantener Given-When-Then focalizados
- ✅ Un concepto por paso (Given/When/Then)
- ✅ Usar datos de prueba realistas

### 2. APIs REST
- ✅ Seguir convenciones RESTful
- ✅ Manejar errores con códigos HTTP apropiados
- ✅ Validar entrada de datos
- ✅ Documentar endpoints claramente

### 3. Testing de APIs
- ✅ Probar casos felices y casos de error
- ✅ Verificar códigos de estado HTTP
- ✅ Validar estructura de respuestas
- ✅ Testear autenticación y autorización

## Resumen del Módulo

### Conceptos Aprendidos
- ✅ Sistema de testing BDD integrado
- ✅ TestCase con Given-When-Then-And
- ✅ Testing de clases y objetos
- ✅ Testing de operaciones concurrentes
- ✅ Desarrollo de APIs REST
- ✅ Testing end-to-end
- ✅ Integración testing-desarrollo

### Habilidades Desarrolladas
- ✅ Escribir tests expressivos con BDD
- ✅ Crear APIs REST completas
- ✅ Diseñar suites de tests efectivas
- ✅ Validar comportamiento de aplicaciones
- ✅ Integrar testing en desarrollo
- ✅ Documentar behavior con tests

### Próximo Módulo

En el **Módulo 6** aprenderás:
- Trabajar con archivos y bases de datos
- Deployment y distribución de aplicaciones
- Optimización y performance
- Patrones avanzados de arquitectura

¡Felicitaciones! Ahora dominas el testing BDD y desarrollo web en R2Lang, dos características distintivas del lenguaje.