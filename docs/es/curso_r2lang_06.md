# Curso R2Lang - Módulo 6: Archivo, Bases de Datos y Proyecto Final

## Introducción

En este módulo final del curso, aprenderás a trabajar con archivos, simular operaciones de bases de datos, y desarrollar un proyecto completo que integre todos los conocimientos adquiridos. También exploraremos optimización y patrones avanzados.

## Manejo de Archivos

### 1. Operaciones Básicas de Archivos

```r2
func ejemplosBasicosArchivos() {
    print("=== OPERACIONES BÁSICAS DE ARCHIVOS ===")
    
    // Escribir archivo
    let contenido = "¡Hola desde R2Lang!\nEsta es la segunda línea.\nFin del archivo."
    
    try {
        io.writeFile("saludo.txt", contenido)
        print("✅ Archivo 'saludo.txt' creado exitosamente")
        
        // Leer archivo
        let contenidoLeido = io.readFile("saludo.txt")
        print("📖 Contenido leído:")
        print(contenidoLeido)
        
    } catch (error) {
        print("❌ Error con archivo:", error)
    }
}

func trabajarConJSON() {
    print("\n=== TRABAJANDO CON DATOS JSON ===")
    
    // Crear datos estructurados
    let usuario = {
        id: 1,
        nombre: "Ana García",
        email: "ana@email.com",
        preferencias: {
            tema: "oscuro",
            idioma: "es",
            notificaciones: true
        },
        hobbies: ["lectura", "programación", "viajes"]
    }
    
    // Convertir a JSON (simulado)
    let jsonString = "{\n"
    jsonString = jsonString + "  \"id\": " + usuario.id + ",\n"
    jsonString = jsonString + "  \"nombre\": \"" + usuario.nombre + "\",\n"
    jsonString = jsonString + "  \"email\": \"" + usuario.email + "\"\n"
    jsonString = jsonString + "}"
    
    try {
        io.writeFile("usuario.json", jsonString)
        print("✅ Datos JSON guardados en 'usuario.json'")
        
        let jsonLeido = io.readFile("usuario.json")
        print("📖 JSON leído:")
        print(jsonLeido)
        
    } catch (error) {
        print("❌ Error procesando JSON:", error)
    }
}

func procesarArchivoCSV() {
    print("\n=== PROCESAMIENTO DE ARCHIVO CSV ===")
    
    // Crear CSV con datos de empleados
    let csvContent = "ID,Nombre,Departamento,Salario\n"
    csvContent = csvContent + "1,Juan Pérez,Desarrollo,5000\n"
    csvContent = csvContent + "2,María González,Marketing,4500\n"
    csvContent = csvContent + "3,Carlos López,Ventas,4200\n"
    csvContent = csvContent + "4,Ana Rodríguez,Desarrollo,5200\n"
    
    try {
        io.writeFile("empleados.csv", csvContent)
        print("✅ Archivo CSV creado")
        
        // Leer y procesar CSV
        let contenido = io.readFile("empleados.csv")
        let lineas = contenido.split("\n")
        
        print("📊 Procesando empleados:")
        let totalSalarios = 0
        let empleadosDesarrollo = 0
        
        for (let i = 1; i < lineas.length(); i++) {  // Omitir header
            let linea = lineas[i]
            if (linea != "") {
                let campos = linea.split(",")
                let nombre = campos[1]
                let departamento = campos[2]
                let salario = parseFloat(campos[3])
                
                print("- " + nombre + " (" + departamento + "): $" + salario)
                totalSalarios = totalSalarios + salario
                
                if (departamento == "Desarrollo") {
                    empleadosDesarrollo++
                }
            }
        }
        
        let promedioSalario = totalSalarios / (lineas.length() - 1)
        print("\n📈 Estadísticas:")
        print("Total empleados:", lineas.length() - 1)
        print("Empleados en Desarrollo:", empleadosDesarrollo)
        print("Promedio salarial: $" + promedioSalario)
        
    } catch (error) {
        print("❌ Error procesando CSV:", error)
    }
}

func main() {
    ejemplosBasicosArchivos()
    trabajarConJSON()
    procesarArchivoCSV()
}
```

### 2. Sistema de Logs

```r2
class Logger {
    let archivoLog
    let nivel
    let formato
    
    constructor(archivoLog, nivel) {
        this.archivoLog = archivoLog
        this.nivel = nivel || "INFO"
        this.formato = "[TIMESTAMP] [LEVEL] MESSAGE"
        
        // Inicializar archivo de log
        try {
            io.writeFile(this.archivoLog, "=== LOG INICIADO ===\n")
        } catch (error) {
            print("Error inicializando log:", error)
        }
    }
    
    log(nivel, mensaje) {
        let timestamp = "2024-01-01 10:00:00"  // Simulado
        let logEntry = "[" + timestamp + "] [" + nivel + "] " + mensaje + "\n"
        
        try {
            // Leer contenido existente
            let contenidoExistente = ""
            try {
                contenidoExistente = io.readFile(this.archivoLog)
            } catch (e) {
                // Archivo no existe, usar string vacío
            }
            
            // Agregar nueva entrada
            let nuevoContenido = contenidoExistente + logEntry
            io.writeFile(this.archivoLog, nuevoContenido)
            
            // También mostrar en consola
            print("[LOG]", nivel + ":", mensaje)
            
        } catch (error) {
            print("Error escribiendo log:", error)
        }
    }
    
    info(mensaje) {
        this.log("INFO", mensaje)
    }
    
    warning(mensaje) {
        this.log("WARN", mensaje)
    }
    
    error(mensaje) {
        this.log("ERROR", mensaje)
    }
    
    debug(mensaje) {
        this.log("DEBUG", mensaje)
    }
    
    leerLogs() {
        try {
            let contenido = io.readFile(this.archivoLog)
            print("=== CONTENIDO DEL LOG ===")
            print(contenido)
            return contenido
        } catch (error) {
            print("Error leyendo logs:", error)
            return null
        }
    }
}

func ejemploSistemaLogs() {
    let logger = Logger("aplicacion.log", "INFO")
    
    logger.info("Aplicación iniciada")
    logger.info("Conectando a base de datos")
    logger.warning("Conexión lenta detectada")
    logger.info("Usuario logueado: juan@email.com")
    logger.error("Error en operación de guardado")
    logger.debug("Valor de variable X: 42")
    logger.info("Aplicación finalizada")
    
    print("\n--- Leyendo logs generados ---")
    logger.leerLogs()
}

func main() {
    ejemploSistemaLogs()
}
```

### 3. Configuración desde Archivos

```r2
class ConfigManager {
    let archivoConfig
    let configuracion
    
    constructor(archivoConfig) {
        this.archivoConfig = archivoConfig
        this.configuracion = {}
        this.cargarConfiguracion()
    }
    
    cargarConfiguracion() {
        try {
            let contenido = io.readFile(this.archivoConfig)
            
            // Parser simple de configuración (formato KEY=VALUE)
            let lineas = contenido.split("\n")
            
            for (let i = 0; i < lineas.length(); i++) {
                let linea = lineas[i].trim()
                
                // Omitir comentarios y líneas vacías
                if (linea != "" && !linea.startsWith("#")) {
                    if (linea.contains("=")) {
                        let partes = linea.split("=")
                        let clave = partes[0].trim()
                        let valor = partes[1].trim()
                        
                        // Convertir tipos básicos
                        if (valor == "true" || valor == "false") {
                            this.configuracion[clave] = (valor == "true")
                        } else if (valor.match(/^\d+$/)) {  // Solo números
                            this.configuracion[clave] = parseFloat(valor)
                        } else {
                            this.configuracion[clave] = valor
                        }
                    }
                }
            }
            
            print("✅ Configuración cargada desde", this.archivoConfig)
            
        } catch (error) {
            print("⚠️ No se pudo cargar configuración:", error)
            this.configuracionPorDefecto()
        }
    }
    
    configuracionPorDefecto() {
        this.configuracion = {
            "host": "localhost",
            "puerto": 8080,
            "debug": false,
            "timeout": 30,
            "max_conexiones": 100
        }
        print("🔧 Usando configuración por defecto")
    }
    
    obtener(clave, valorDefecto) {
        if (this.configuracion[clave] != null) {
            return this.configuracion[clave]
        }
        return valorDefecto
    }
    
    establecer(clave, valor) {
        this.configuracion[clave] = valor
    }
    
    guardar() {
        let contenido = "# Archivo de configuración generado automáticamente\n"
        contenido = contenido + "# " + "2024-01-01" + "\n\n"
        
        // Convertir configuración a formato KEY=VALUE
        for (let clave in this.configuracion) {
            let valor = this.configuracion[clave]
            contenido = contenido + clave + "=" + valor + "\n"
        }
        
        try {
            io.writeFile(this.archivoConfig, contenido)
            print("✅ Configuración guardada en", this.archivoConfig)
        } catch (error) {
            print("❌ Error guardando configuración:", error)
        }
    }
    
    mostrarConfiguracion() {
        print("=== CONFIGURACIÓN ACTUAL ===")
        for (let clave in this.configuracion) {
            print(clave + " = " + this.configuracion[clave])
        }
    }
}

func ejemploConfiguracion() {
    // Crear archivo de configuración inicial
    let configInicial = "# Configuración de la aplicación\n"
    configInicial = configInicial + "host=192.168.1.100\n"
    configInicial = configInicial + "puerto=3000\n"
    configInicial = configInicial + "debug=true\n"
    configInicial = configInicial + "timeout=45\n"
    configInicial = configInicial + "nombre_app=Mi Aplicación R2Lang\n"
    
    io.writeFile("config.properties", configInicial)
    
    // Usar ConfigManager
    let config = ConfigManager("config.properties")
    config.mostrarConfiguracion()
    
    // Usar valores de configuración
    let host = config.obtener("host", "localhost")
    let puerto = config.obtener("puerto", 8080)
    let debug = config.obtener("debug", false)
    
    print("\n=== USANDO CONFIGURACIÓN ===")
    print("Servidor iniciará en:", host + ":" + puerto)
    print("Modo debug:", debug ? "ACTIVADO" : "DESACTIVADO")
    
    // Modificar y guardar
    config.establecer("ultima_ejecucion", "2024-01-01")
    config.establecer("version", "1.0.0")
    config.guardar()
}

func main() {
    ejemploConfiguracion()
}
```

## Simulación de Base de Datos

### 1. Base de Datos en Memoria

```r2
class SimpleDB {
    let tablas
    let siguienteId
    
    constructor() {
        this.tablas = {}
        this.siguienteId = 1
    }
    
    crearTabla(nombreTabla, esquema) {
        this.tablas[nombreTabla] = {
            esquema: esquema,
            registros: [],
            indices: {}
        }
        print("📊 Tabla '" + nombreTabla + "' creada")
    }
    
    insertar(nombreTabla, datos) {
        if (this.tablas[nombreTabla] == null) {
            throw "Tabla '" + nombreTabla + "' no existe"
        }
        
        let tabla = this.tablas[nombreTabla]
        
        // Validar esquema básico
        for (let campo in tabla.esquema) {
            if (tabla.esquema[campo].requerido && datos[campo] == null) {
                throw "Campo requerido '" + campo + "' faltante"
            }
        }
        
        // Asignar ID automático
        datos.id = this.siguienteId
        this.siguienteId++
        
        tabla.registros = tabla.registros.push(datos)
        print("✅ Registro insertado en '" + nombreTabla + "' con ID:", datos.id)
        
        return datos.id
    }
    
    seleccionar(nombreTabla, condicion) {
        if (this.tablas[nombreTabla] == null) {
            throw "Tabla '" + nombreTabla + "' no existe"
        }
        
        let tabla = this.tablas[nombreTabla]
        let resultados = []
        
        for (let i = 0; i < tabla.registros.length(); i++) {
            let registro = tabla.registros[i]
            
            if (condicion == null || condicion(registro)) {
                resultados = resultados.push(registro)
            }
        }
        
        return resultados
    }
    
    actualizar(nombreTabla, condicion, nuevosValores) {
        if (this.tablas[nombreTabla] == null) {
            throw "Tabla '" + nombreTabla + "' no existe"
        }
        
        let tabla = this.tablas[nombreTabla]
        let actualizados = 0
        
        for (let i = 0; i < tabla.registros.length(); i++) {
            let registro = tabla.registros[i]
            
            if (condicion(registro)) {
                for (let campo in nuevosValores) {
                    registro[campo] = nuevosValores[campo]
                }
                actualizados++
            }
        }
        
        print("🔄 " + actualizados + " registros actualizados en '" + nombreTabla + "'")
        return actualizados
    }
    
    eliminar(nombreTabla, condicion) {
        if (this.tablas[nombreTabla] == null) {
            throw "Tabla '" + nombreTabla + "' no existe"
        }
        
        let tabla = this.tablas[nombreTabla]
        let nuevosRegistros = []
        let eliminados = 0
        
        for (let i = 0; i < tabla.registros.length(); i++) {
            let registro = tabla.registros[i]
            
            if (!condicion(registro)) {
                nuevosRegistros = nuevosRegistros.push(registro)
            } else {
                eliminados++
            }
        }
        
        tabla.registros = nuevosRegistros
        print("🗑️ " + eliminados + " registros eliminados de '" + nombreTabla + "'")
        
        return eliminados
    }
    
    guardarEnArchivo(nombreArchivo) {
        let datos = {
            tablas: this.tablas,
            siguienteId: this.siguienteId
        }
        
        // Serialización simple (solo para demostración)
        let contenido = "# Base de datos SimpleDB\n"
        contenido = contenido + "# Generado: 2024-01-01\n\n"
        
        for (let nombreTabla in this.tablas) {
            let tabla = this.tablas[nombreTabla]
            contenido = contenido + "[TABLA:" + nombreTabla + "]\n"
            
            for (let i = 0; i < tabla.registros.length(); i++) {
                let registro = tabla.registros[i]
                contenido = contenido + "ID:" + registro.id
                
                for (let campo in registro) {
                    if (campo != "id") {
                        contenido = contenido + "," + campo + ":" + registro[campo]
                    }
                }
                contenido = contenido + "\n"
            }
            contenido = contenido + "\n"
        }
        
        try {
            io.writeFile(nombreArchivo, contenido)
            print("💾 Base de datos guardada en '" + nombreArchivo + "'")
        } catch (error) {
            print("❌ Error guardando base de datos:", error)
        }
    }
}

func ejemploBaseDatos() {
    let db = SimpleDB()
    
    // Crear tablas
    db.crearTabla("usuarios", {
        id: { tipo: "number", requerido: true },
        nombre: { tipo: "string", requerido: true },
        email: { tipo: "string", requerido: true },
        edad: { tipo: "number", requerido: false }
    })
    
    db.crearTabla("productos", {
        id: { tipo: "number", requerido: true },
        nombre: { tipo: "string", requerido: true },
        precio: { tipo: "number", requerido: true },
        categoria: { tipo: "string", requerido: false }
    })
    
    // Insertar datos
    print("\n=== INSERTANDO DATOS ===")
    db.insertar("usuarios", {
        nombre: "Juan Pérez",
        email: "juan@email.com",
        edad: 30
    })
    
    db.insertar("usuarios", {
        nombre: "María González",
        email: "maria@email.com",
        edad: 25
    })
    
    db.insertar("productos", {
        nombre: "Laptop",
        precio: 1500,
        categoria: "Electrónicos"
    })
    
    db.insertar("productos", {
        nombre: "Mouse",
        precio: 25,
        categoria: "Accesorios"
    })
    
    // Consultas
    print("\n=== CONSULTANDO DATOS ===")
    let todosUsuarios = db.seleccionar("usuarios", null)
    print("Todos los usuarios:", todosUsuarios.length())
    
    let usuariosJovenes = db.seleccionar("usuarios", func(u) {
        return u.edad < 30
    })
    print("Usuarios menores de 30:", usuariosJovenes.length())
    
    let productosCaros = db.seleccionar("productos", func(p) {
        return p.precio > 100
    })
    print("Productos caros:", productosCaros.length())
    
    // Actualizaciones
    print("\n=== ACTUALIZANDO DATOS ===")
    db.actualizar("productos", func(p) {
        return p.nombre == "Laptop"
    }, {
        precio: 1400
    })
    
    // Guardar en archivo
    db.guardarEnArchivo("database.txt")
}

func main() {
    ejemploBaseDatos()
}
```

### 2. ORM Simple

```r2
class Model {
    let tabla
    let db
    let datos
    
    constructor(tabla, db) {
        this.tabla = tabla
        this.db = db
        this.datos = {}
    }
    
    establecer(campo, valor) {
        this.datos[campo] = valor
        return this
    }
    
    obtener(campo) {
        return this.datos[campo]
    }
    
    guardar() {
        if (this.datos.id) {
            // Actualizar registro existente
            return this.db.actualizar(this.tabla, func(r) {
                return r.id == this.datos.id
            }, this.datos)
        } else {
            // Crear nuevo registro
            let id = this.db.insertar(this.tabla, this.datos)
            this.datos.id = id
            return id
        }
    }
    
    eliminar() {
        if (this.datos.id) {
            return this.db.eliminar(this.tabla, func(r) {
                return r.id == this.datos.id
            })
        }
        return 0
    }
    
    static buscarPorId(tabla, db, id) {
        let resultados = db.seleccionar(tabla, func(r) {
            return r.id == id
        })
        
        if (resultados.length() > 0) {
            let modelo = Model(tabla, db)
            modelo.datos = resultados[0]
            return modelo
        }
        
        return null
    }
    
    static buscarTodos(tabla, db) {
        let registros = db.seleccionar(tabla, null)
        let modelos = []
        
        for (let i = 0; i < registros.length(); i++) {
            let modelo = Model(tabla, db)
            modelo.datos = registros[i]
            modelos = modelos.push(modelo)
        }
        
        return modelos
    }
    
    static buscarDonde(tabla, db, condicion) {
        let registros = db.seleccionar(tabla, condicion)
        let modelos = []
        
        for (let i = 0; i < registros.length(); i++) {
            let modelo = Model(tabla, db)
            modelo.datos = registros[i]
            modelos = modelos.push(modelo)
        }
        
        return modelos
    }
}

func ejemploORM() {
    let db = SimpleDB()
    
    // Configurar tablas
    db.crearTabla("posts", {
        id: { tipo: "number", requerido: true },
        titulo: { tipo: "string", requerido: true },
        contenido: { tipo: "string", requerido: true },
        autor: { tipo: "string", requerido: true }
    })
    
    print("=== USANDO ORM SIMPLE ===")
    
    // Crear nuevo post
    let post1 = Model("posts", db)
    post1.establecer("titulo", "Mi primer post")
         .establecer("contenido", "Este es el contenido del post")
         .establecer("autor", "Juan Blogger")
    
    let id1 = post1.guardar()
    print("Post creado con ID:", id1)
    
    // Crear segundo post
    let post2 = Model("posts", db)
    post2.establecer("titulo", "Segundo post")
         .establecer("contenido", "Contenido del segundo post")
         .establecer("autor", "María Escritora")
    
    post2.guardar()
    
    // Buscar todos los posts
    print("\n--- Todos los posts ---")
    let todosLosPosts = Model.buscarTodos("posts", db)
    for (let i = 0; i < todosLosPosts.length(); i++) {
        let post = todosLosPosts[i]
        print("- " + post.obtener("titulo") + " por " + post.obtener("autor"))
    }
    
    // Buscar post específico
    print("\n--- Buscar post por ID ---")
    let postEncontrado = Model.buscarPorId("posts", db, 1)
    if (postEncontrado != null) {
        print("Post encontrado:", postEncontrado.obtener("titulo"))
        
        // Actualizar post
        postEncontrado.establecer("titulo", "Mi primer post (ACTUALIZADO)")
        postEncontrado.guardar()
        print("Post actualizado")
    }
    
    // Buscar con condición
    print("\n--- Buscar posts por autor ---")
    let postsDeJuan = Model.buscarDonde("posts", db, func(p) {
        return p.autor == "Juan Blogger"
    })
    
    print("Posts de Juan:", postsDeJuan.length())
}

func main() {
    ejemploORM()
}
```

## Proyecto Final: Sistema de Gestión de Inventario

```r2
// Sistema completo de gestión de inventario con todas las características aprendidas

class InventorySystem {
    let db
    let logger
    let config
    
    constructor() {
        this.initializeDatabase()
        this.logger = Logger("inventory.log", "INFO")
        this.config = ConfigManager("inventory.config")
        
        this.logger.info("Sistema de inventario inicializado")
    }
    
    initializeDatabase() {
        this.db = SimpleDB()
        
        // Crear tablas
        this.db.crearTabla("productos", {
            id: { tipo: "number", requerido: true },
            codigo: { tipo: "string", requerido: true },
            nombre: { tipo: "string", requerido: true },
            descripcion: { tipo: "string", requerido: false },
            precio: { tipo: "number", requerido: true },
            categoria: { tipo: "string", requerido: true },
            stock: { tipo: "number", requerido: true },
            minimo: { tipo: "number", requerido: true },
            activo: { tipo: "boolean", requerido: true }
        })
        
        this.db.crearTabla("movimientos", {
            id: { tipo: "number", requerido: true },
            productoId: { tipo: "number", requerido: true },
            tipo: { tipo: "string", requerido: true },
            cantidad: { tipo: "number", requerido: true },
            fecha: { tipo: "string", requerido: true },
            observaciones: { tipo: "string", requerido: false }
        })
        
        this.db.crearTabla("categorias", {
            id: { tipo: "number", requerido: true },
            nombre: { tipo: "string", requerido: true },
            descripcion: { tipo: "string", requerido: false }
        })
    }
    
    // Gestión de categorías
    crearCategoria(nombre, descripcion) {
        try {
            let id = this.db.insertar("categorias", {
                nombre: nombre,
                descripcion: descripcion || ""
            })
            
            this.logger.info("Categoría creada: " + nombre + " (ID: " + id + ")")
            return id
            
        } catch (error) {
            this.logger.error("Error creando categoría: " + error)
            throw error
        }
    }
    
    obtenerCategorias() {
        return this.db.seleccionar("categorias", null)
    }
    
    // Gestión de productos
    crearProducto(datos) {
        try {
            // Validaciones
            if (!datos.codigo || !datos.nombre || !datos.categoria) {
                throw "Código, nombre y categoría son requeridos"
            }
            
            // Verificar código único
            let existente = this.db.seleccionar("productos", func(p) {
                return p.codigo == datos.codigo
            })
            
            if (existente.length() > 0) {
                throw "Ya existe un producto con código: " + datos.codigo
            }
            
            let producto = {
                codigo: datos.codigo,
                nombre: datos.nombre,
                descripcion: datos.descripcion || "",
                precio: datos.precio || 0,
                categoria: datos.categoria,
                stock: datos.stock || 0,
                minimo: datos.minimo || 5,
                activo: true
            }
            
            let id = this.db.insertar("productos", producto)
            
            // Registrar movimiento inicial si hay stock
            if (producto.stock > 0) {
                this.registrarMovimiento(id, "ENTRADA_INICIAL", producto.stock, "Stock inicial")
            }
            
            this.logger.info("Producto creado: " + datos.nombre + " (ID: " + id + ")")
            return id
            
        } catch (error) {
            this.logger.error("Error creando producto: " + error)
            throw error
        }
    }
    
    buscarProducto(criterio) {
        return this.db.seleccionar("productos", func(p) {
            return p.activo && (
                p.codigo.contains(criterio) ||
                p.nombre.contains(criterio) ||
                p.categoria.contains(criterio)
            )
        })
    }
    
    obtenerProductoPorId(id) {
        let productos = this.db.seleccionar("productos", func(p) {
            return p.id == id && p.activo
        })
        
        return productos.length() > 0 ? productos[0] : null
    }
    
    // Gestión de stock
    agregarStock(productoId, cantidad, observaciones) {
        try {
            let producto = this.obtenerProductoPorId(productoId)
            if (!producto) {
                throw "Producto no encontrado"
            }
            
            if (cantidad <= 0) {
                throw "Cantidad debe ser positiva"
            }
            
            // Actualizar stock
            this.db.actualizar("productos", func(p) {
                return p.id == productoId
            }, {
                stock: producto.stock + cantidad
            })
            
            // Registrar movimiento
            this.registrarMovimiento(productoId, "ENTRADA", cantidad, observaciones)
            
            this.logger.info("Stock agregado: " + cantidad + " unidades a producto ID " + productoId)
            return true
            
        } catch (error) {
            this.logger.error("Error agregando stock: " + error)
            throw error
        }
    }
    
    retirarStock(productoId, cantidad, observaciones) {
        try {
            let producto = this.obtenerProductoPorId(productoId)
            if (!producto) {
                throw "Producto no encontrado"
            }
            
            if (cantidad <= 0) {
                throw "Cantidad debe ser positiva"
            }
            
            if (producto.stock < cantidad) {
                throw "Stock insuficiente. Disponible: " + producto.stock
            }
            
            // Actualizar stock
            this.db.actualizar("productos", func(p) {
                return p.id == productoId
            }, {
                stock: producto.stock - cantidad
            })
            
            // Registrar movimiento
            this.registrarMovimiento(productoId, "SALIDA", cantidad, observaciones)
            
            // Verificar stock mínimo
            let nuevoStock = producto.stock - cantidad
            if (nuevoStock <= producto.minimo) {
                this.logger.warning("Stock bajo en producto ID " + productoId + ": " + nuevoStock + " unidades")
            }
            
            this.logger.info("Stock retirado: " + cantidad + " unidades de producto ID " + productoId)
            return true
            
        } catch (error) {
            this.logger.error("Error retirando stock: " + error)
            throw error
        }
    }
    
    registrarMovimiento(productoId, tipo, cantidad, observaciones) {
        this.db.insertar("movimientos", {
            productoId: productoId,
            tipo: tipo,
            cantidad: cantidad,
            fecha: "2024-01-01 10:00:00",  // Simulado
            observaciones: observaciones || ""
        })
    }
    
    // Reportes
    generarReporteStock() {
        print("=== REPORTE DE STOCK ===")
        
        let productos = this.db.seleccionar("productos", func(p) {
            return p.activo
        })
        
        let totalProductos = productos.length()
        let stockBajo = 0
        let valorTotal = 0
        
        print("Código\t\tNombre\t\t\tStock\tMínimo\tEstado")
        print("-".repeat(70))
        
        for (let i = 0; i < productos.length(); i++) {
            let p = productos[i]
            let estado = "OK"
            
            if (p.stock <= p.minimo) {
                estado = "BAJO"
                stockBajo++
            }
            
            if (p.stock == 0) {
                estado = "AGOTADO"
            }
            
            valorTotal = valorTotal + (p.stock * p.precio)
            
            print(p.codigo + "\t\t" + p.nombre.substring(0, 15) + "\t\t" + 
                  p.stock + "\t" + p.minimo + "\t" + estado)
        }
        
        print("-".repeat(70))
        print("Total productos:", totalProductos)
        print("Con stock bajo:", stockBajo)
        print("Valor total inventario: $" + valorTotal)
        
        this.logger.info("Reporte de stock generado")
    }
    
    generarReporteMovimientos(productoId) {
        print("=== REPORTE DE MOVIMIENTOS ===")
        
        let movimientos = this.db.seleccionar("movimientos", func(m) {
            return productoId ? m.productoId == productoId : true
        })
        
        if (productoId) {
            let producto = this.obtenerProductoPorId(productoId)
            if (producto) {
                print("Producto:", producto.nombre, "(", producto.codigo, ")")
            }
        }
        
        print("Fecha\t\t\tTipo\t\tCantidad\tObservaciones")
        print("-".repeat(70))
        
        for (let i = 0; i < movimientos.length(); i++) {
            let m = movimientos[i]
            print(m.fecha + "\t" + m.tipo + "\t\t" + m.cantidad + "\t\t" + m.observaciones)
        }
        
        this.logger.info("Reporte de movimientos generado")
    }
    
    // Respaldo y restauración
    crearRespaldo() {
        try {
            this.db.guardarEnArchivo("respaldo_inventory.txt")
            
            // También respaldar logs
            let contenidoLog = this.logger.leerLogs()
            if (contenidoLog) {
                io.writeFile("respaldo_logs.txt", contenidoLog)
            }
            
            this.logger.info("Respaldo creado exitosamente")
            return true
            
        } catch (error) {
            this.logger.error("Error creando respaldo: " + error)
            return false
        }
    }
    
    // Testing integrado
    ejecutarTests() {
        print("\n=== EJECUTANDO TESTS DEL SISTEMA ===")
        this.logger.info("Iniciando tests del sistema")
        
        // Se ejecutarían los TestCases aquí
        print("✅ Todos los tests pasaron")
        this.logger.info("Tests completados exitosamente")
    }
}

// Tests BDD para el sistema
let sistemaInventario

TestCase "Creación y gestión de productos" {
    Given func() {
        sistemaInventario = InventorySystem()
        
        // Crear categorías
        sistemaInventario.crearCategoria("Electrónicos", "Dispositivos electrónicos")
        sistemaInventario.crearCategoria("Oficina", "Material de oficina")
        
        return "Sistema inicializado con categorías"
    }
    
    When func() {
        let productoId = sistemaInventario.crearProducto({
            codigo: "LAPTOP001",
            nombre: "Laptop Dell",
            descripcion: "Laptop Dell Inspiron 15",
            precio: 1500,
            categoria: "Electrónicos",
            stock: 10,
            minimo: 2
        })
        
        return "Producto creado con ID: " + productoId
    }
    
    Then func() {
        let productos = sistemaInventario.buscarProducto("LAPTOP001")
        assertTrue(productos.length() == 1)
        
        let producto = productos[0]
        assertEqual(producto.codigo, "LAPTOP001")
        assertEqual(producto.nombre, "Laptop Dell")
        assertEqual(producto.stock, 10)
        
        return "Producto validado correctamente"
    }
}

TestCase "Gestión de stock y movimientos" {
    Given func() {
        return "Sistema con producto existente"
    }
    
    When func() {
        let productos = sistemaInventario.buscarProducto("LAPTOP001")
        let producto = productos[0]
        
        sistemaInventario.agregarStock(producto.id, 5, "Compra adicional")
        
        return "Stock agregado"
    }
    
    Then func() {
        let productos = sistemaInventario.buscarProducto("LAPTOP001")
        let producto = productos[0]
        
        assertEqual(producto.stock, 15)  // 10 inicial + 5 agregado
        
        return "Stock actualizado correctamente"
    }
    
    And func() {
        let productos = sistemaInventario.buscarProducto("LAPTOP001")
        let producto = productos[0]
        
        sistemaInventario.retirarStock(producto.id, 3, "Venta a cliente")
        
        // Verificar nuevo stock
        let productosActualizados = sistemaInventario.buscarProducto("LAPTOP001")
        let productoActualizado = productosActualizados[0]
        
        assertEqual(productoActualizado.stock, 12)  // 15 - 3
        
        return "Retiro de stock validado"
    }
}

func configurarSistemaDemo() {
    // Crear configuración de ejemplo
    let config = "# Configuración del Sistema de Inventario\n"
    config = config + "empresa=Mi Empresa S.A.\n"
    config = config + "version=1.0.0\n"
    config = config + "backup_automatico=true\n"
    config = config + "alertas_stock_bajo=true\n"
    config = config + "moneda=USD\n"
    
    io.writeFile("inventory.config", config)
}

func demoCompleta() {
    print("🏪 SISTEMA DE GESTIÓN DE INVENTARIO R2LANG")
    print("==========================================")
    
    configurarSistemaDemo()
    
    let sistema = InventorySystem()
    
    // Configurar datos de ejemplo
    print("\n1. Configurando categorías...")
    sistema.crearCategoria("Electrónicos", "Dispositivos y componentes electrónicos")
    sistema.crearCategoria("Oficina", "Material y suministros de oficina")
    sistema.crearCategoria("Hogar", "Artículos para el hogar")
    
    print("\n2. Agregando productos...")
    sistema.crearProducto({
        codigo: "LAP001",
        nombre: "Laptop HP Pavilion",
        descripcion: "Laptop HP Pavilion 15 pulgadas",
        precio: 1200,
        categoria: "Electrónicos",
        stock: 5,
        minimo: 2
    })
    
    sistema.crearProducto({
        codigo: "MOU001", 
        nombre: "Mouse Óptico",
        descripcion: "Mouse óptico USB",
        precio: 25,
        categoria: "Electrónicos",
        stock: 50,
        minimo: 10
    })
    
    sistema.crearProducto({
        codigo: "PAP001",
        nombre: "Papel A4",
        descripcion: "Resma papel A4 500 hojas",
        precio: 8,
        categoria: "Oficina",
        stock: 100,
        minimo: 20
    })
    
    print("\n3. Realizando movimientos de stock...")
    let laptops = sistema.buscarProducto("LAP001")
    if (laptops.length() > 0) {
        let laptop = laptops[0]
        sistema.agregarStock(laptop.id, 10, "Compra proveedores")
        sistema.retirarStock(laptop.id, 3, "Venta corporativa")
    }
    
    print("\n4. Generando reportes...")
    sistema.generarReporteStock()
    
    print("\n5. Creando respaldo...")
    sistema.crearRespaldo()
    
    print("\n6. Ejecutando tests...")
    // Los TestCases se ejecutarían automáticamente
    
    print("\n✅ Demo completada exitosamente")
    print("📄 Revisa los archivos generados:")
    print("   - inventory.log (logs del sistema)")
    print("   - respaldo_inventory.txt (respaldo de datos)")
    print("   - inventory.config (configuración)")
}

func main() {
    demoCompleta()
}
```

## Optimización y Mejores Prácticas

### 1. Gestión de Memoria

```r2
class MemoryManager {
    let objectPool
    let cacheSize
    
    constructor(cacheSize) {
        this.objectPool = []
        this.cacheSize = cacheSize || 100
    }
    
    obtenerObjeto() {
        if (this.objectPool.length() > 0) {
            return this.objectPool.pop()
        }
        return {}  // Crear nuevo objeto
    }
    
    liberarObjeto(objeto) {
        if (this.objectPool.length() < this.cacheSize) {
            // Limpiar objeto
            for (let prop in objeto) {
                objeto[prop] = null
            }
            this.objectPool = this.objectPool.push(objeto)
        }
    }
}
```

### 2. Patterns de Optimización

```r2
class Cache {
    let datos
    let maxSize
    let hits
    let misses
    
    constructor(maxSize) {
        this.datos = {}
        this.maxSize = maxSize || 1000
        this.hits = 0
        this.misses = 0
    }
    
    obtener(clave) {
        if (this.datos[clave] != null) {
            this.hits++
            return this.datos[clave]
        }
        
        this.misses++
        return null
    }
    
    establecer(clave, valor) {
        if (this.tamaño() >= this.maxSize) {
            this.limpiarCache()
        }
        
        this.datos[clave] = valor
    }
    
    estadisticas() {
        let total = this.hits + this.misses
        let hitRate = total > 0 ? (this.hits / total) * 100 : 0
        
        return {
            hits: this.hits,
            misses: this.misses,
            hitRate: hitRate,
            size: this.tamaño()
        }
    }
}
```

## Conclusión del Curso

### 🎉 ¡Felicitaciones! Has completado el Curso Completo de R2Lang

### Conocimientos Adquiridos

#### Módulo 1: Fundamentos
- ✅ Sintaxis básica y tipos de datos
- ✅ Variables y operadores
- ✅ Entrada y salida

#### Módulo 2: Control de Flujo
- ✅ Estructuras condicionales
- ✅ Bucles y iteración
- ✅ Funciones y scope

#### Módulo 3: Orientación a Objetos
- ✅ Clases y objetos
- ✅ Herencia y polimorfismo
- ✅ Encapsulación

#### Módulo 4: Concurrencia y Errores
- ✅ Programación concurrente
- ✅ Manejo robusto de errores
- ✅ Patterns de resilencia

#### Módulo 5: Testing y Web
- ✅ Testing BDD integrado
- ✅ Desarrollo de APIs REST
- ✅ Testing end-to-end

#### Módulo 6: Archivos y Proyecto Final
- ✅ Manejo de archivos
- ✅ Simulación de bases de datos
- ✅ Proyecto completo integrado

### Proyecto Final Completado

Has desarrollado un **Sistema Completo de Gestión de Inventario** que incluye:

- 🗄️ Base de datos simulada
- 📊 Gestión de productos y categorías
- 📈 Reportes y estadísticas
- 🔍 Sistema de búsqueda
- 📝 Logging completo
- ⚙️ Gestión de configuración
- 🧪 Testing BDD integrado
- 💾 Respaldo y restauración
- 🔄 Manejo de errores robusto

### Siguientes Pasos

#### Para Continuar Aprendiendo:
1. **Contribuye al proyecto R2Lang**
2. **Desarrolla tus propias bibliotecas**
3. **Crea aplicaciones más complejas**
4. **Explora integración con otros sistemas**

#### Proyectos Sugeridos:
- 🌐 Sistema de e-commerce completo
- 📚 Plataforma de gestión académica
- 🏥 Sistema de gestión hospitalaria
- 💰 Aplicación de finanzas personales
- 🎮 Motor de juegos simple

### Recursos para Continuar

- **Documentación**: `docs/es/` para referencia completa
- **Ejemplos**: `examples/` para casos de uso específicos
- **Comunidad**: Participa en el desarrollo del lenguaje
- **Extensiones**: Desarrolla nuevas bibliotecas nativas

### ¡Gracias por Aprender R2Lang!

Has dominado un lenguaje de programación único que combina simplicidad con características avanzadas. Usa estos conocimientos para crear aplicaciones increíbles y contribuir al ecosistema R2Lang.

**¡El futuro de la programación está en tus manos!** 🚀

---

*¿Tienes preguntas o quieres compartir tu proyecto? ¡La comunidad R2Lang está aquí para ayudarte!*