# Curso R2Lang - Módulo 3: Orientación a Objetos

## Introducción a la Programación Orientada a Objetos

La Programación Orientada a Objetos (OOP) es un paradigma que organiza el código en **objetos** que contienen datos (propiedades) y código (métodos). R2Lang soporta OOP con clases, herencia, y encapsulación.

### Conceptos Fundamentales

- **Clase**: Plantilla o blueprint para crear objetos
- **Objeto**: Instancia de una clase
- **Propiedades**: Variables que pertenecen a un objeto
- **Métodos**: Funciones que pertenecen a un objeto
- **Constructor**: Método especial que inicializa un objeto
- **Herencia**: Capacidad de una clase de heredar de otra

## Clases y Objetos Básicos

### 1. Definición de una Clase Simple

```r2
class Persona {
    // Propiedades (se declaran con let)
    let nombre
    let edad
    let email
    
    // Constructor: método especial que se ejecuta al crear un objeto
    constructor(nombre, edad, email) {
        this.nombre = nombre
        this.edad = edad
        this.email = email
        print("Nueva persona creada:", this.nombre)
    }
    
    // Métodos: funciones que pertenecen a la clase
    saludar() {
        print("Hola, soy", this.nombre, "y tengo", this.edad, "años")
    }
    
    cumpleanos() {
        this.edad++
        print(this.nombre, "ahora tiene", this.edad, "años")
    }
    
    cambiarEmail(nuevoEmail) {
        let emailAnterior = this.email
        this.email = nuevoEmail
        print("Email de", this.nombre, "cambió de", emailAnterior, "a", nuevoEmail)
    }
}

func main() {
    // Crear objetos (instanciar la clase)
    let persona1 = Persona("Ana García", 25, "ana@email.com")
    let persona2 = Persona("Carlos López", 30, "carlos@email.com")
    
    // Usar métodos
    persona1.saludar()
    persona2.saludar()
    
    // Modificar propiedades a través de métodos
    persona1.cumpleanos()
    persona2.cambiarEmail("carlos.lopez@gmail.com")
    
    // Acceder propiedades directamente
    print("Nombre de persona1:", persona1.nombre)
    print("Edad de persona2:", persona2.edad)
}
```

### 2. Clase con Métodos Avanzados

```r2
class CuentaBancaria {
    let titular
    let numeroCuenta
    let saldo
    let movimientos
    
    constructor(titular, numeroCuenta, saldoInicial) {
        this.titular = titular
        this.numeroCuenta = numeroCuenta
        this.saldo = saldoInicial
        this.movimientos = []
        
        // Registrar movimiento inicial
        let movimientoInicial = {
            tipo: "Apertura",
            monto: saldoInicial,
            fecha: "Hoy",
            saldoResultante: saldoInicial
        }
        this.movimientos = this.movimientos.push(movimientoInicial)
    }
    
    depositar(monto) {
        if (monto <= 0) {
            print("Error: El monto debe ser positivo")
            return false
        }
        
        this.saldo = this.saldo + monto
        let movimiento = {
            tipo: "Depósito",
            monto: monto,
            fecha: "Hoy",
            saldoResultante: this.saldo
        }
        this.movimientos = this.movimientos.push(movimiento)
        
        print("Depósito exitoso. Nuevo saldo:", this.saldo)
        return true
    }
    
    retirar(monto) {
        if (monto <= 0) {
            print("Error: El monto debe ser positivo")
            return false
        }
        
        if (monto > this.saldo) {
            print("Error: Saldo insuficiente")
            return false
        }
        
        this.saldo = this.saldo - monto
        let movimiento = {
            tipo: "Retiro",
            monto: monto,
            fecha: "Hoy",
            saldoResultante: this.saldo
        }
        this.movimientos = this.movimientos.push(movimiento)
        
        print("Retiro exitoso. Nuevo saldo:", this.saldo)
        return true
    }
    
    consultarSaldo() {
        print("Saldo actual de", this.titular + ":", this.saldo)
        return this.saldo
    }
    
    mostrarMovimientos() {
        print("=== MOVIMIENTOS DE", this.titular, "===")
        print("Cuenta:", this.numeroCuenta)
        
        for (let i = 0; i < this.movimientos.length(); i++) {
            let mov = this.movimientos[i]
            print("- " + mov.tipo + ": $" + mov.monto + " (Saldo: $" + mov.saldoResultante + ")")
        }
        print("Total de movimientos:", this.movimientos.length())
    }
    
    transferir(cuentaDestino, monto) {
        if (this.retirar(monto)) {
            if (cuentaDestino.depositar(monto)) {
                print("Transferencia exitosa de", this.titular, "a", cuentaDestino.titular)
                return true
            } else {
                // Revertir el retiro si el depósito falla
                this.depositar(monto)
                print("Error en transferencia: depósito falló")
                return false
            }
        }
        return false
    }
}

func main() {
    // Crear cuentas bancarias
    let cuentaAna = CuentaBancaria("Ana García", "12345", 1000)
    let cuentaCarlos = CuentaBancaria("Carlos López", "67890", 500)
    
    print()
    
    // Operaciones básicas
    cuentaAna.consultarSaldo()
    cuentaAna.depositar(200)
    cuentaAna.retirar(150)
    
    print()
    
    // Transferencia entre cuentas
    cuentaAna.transferir(cuentaCarlos, 300)
    
    print()
    
    // Mostrar movimientos
    cuentaAna.mostrarMovimientos()
    print()
    cuentaCarlos.mostrarMovimientos()
}
```

## Herencia

### 1. Herencia Básica

```r2
// Clase padre (superclase)
class Animal {
    let nombre
    let especie
    let edad
    
    constructor(nombre, especie, edad) {
        this.nombre = nombre
        this.especie = especie
        this.edad = edad
    }
    
    hacerSonido() {
        print(this.nombre, "hace un sonido")
    }
    
    dormir() {
        print(this.nombre, "está durmiendo")
    }
    
    comer() {
        print(this.nombre, "está comiendo")
    }
    
    mostrarInfo() {
        print("Nombre:", this.nombre)
        print("Especie:", this.especie)
        print("Edad:", this.edad, "años")
    }
}

// Clase hija (subclase)
class Perro extends Animal {
    let raza
    let dueno
    
    constructor(nombre, edad, raza, dueno) {
        // Llamar al constructor de la clase padre
        super.constructor(nombre, "Canino", edad)
        this.raza = raza
        this.dueno = dueno
    }
    
    // Sobrescribir método del padre
    hacerSonido() {
        print(this.nombre, "ladra: ¡Guau guau!")
    }
    
    // Métodos específicos de Perro
    buscarPelota() {
        print(this.nombre, "está buscando la pelota")
    }
    
    moverCola() {
        print(this.nombre, "mueve la cola felizmente")
    }
    
    // Sobrescribir método para mostrar información específica
    mostrarInfo() {
        super.mostrarInfo()  // Llamar método del padre
        print("Raza:", this.raza)
        print("Dueño:", this.dueno)
    }
}

class Gato extends Animal {
    let colorPelaje
    let esIndependiente
    
    constructor(nombre, edad, colorPelaje) {
        super.constructor(nombre, "Felino", edad)
        this.colorPelaje = colorPelaje
        this.esIndependiente = true
    }
    
    hacerSonido() {
        print(this.nombre, "maúlla: ¡Miau!")
    }
    
    ronronear() {
        print(this.nombre, "ronronea satisfecho")
    }
    
    afilarGarras() {
        print(this.nombre, "está afilando sus garras")
    }
    
    mostrarInfo() {
        super.mostrarInfo()
        print("Color del pelaje:", this.colorPelaje)
        print("Es independiente:", this.esIndependiente)
    }
}

func main() {
    // Crear instancias
    let miPerro = Perro("Max", 3, "Labrador", "Juan")
    let miGato = Gato("Luna", 2, "Gris")
    
    print("=== INFORMACIÓN DE MASCOTAS ===")
    print()
    
    print("🐕 PERRO:")
    miPerro.mostrarInfo()
    print()
    
    print("🐱 GATO:")
    miGato.mostrarInfo()
    print()
    
    print("=== ACCIONES ===")
    // Métodos heredados
    miPerro.comer()
    miGato.dormir()
    
    // Métodos sobrescritos
    miPerro.hacerSonido()
    miGato.hacerSonido()
    
    // Métodos específicos
    miPerro.buscarPelota()
    miPerro.moverCola()
    
    miGato.ronronear()
    miGato.afilarGarras()
}
```

### 2. Herencia Múltiple Niveles

```r2
// Clase base
class Vehiculo {
    let marca
    let modelo
    let ano
    let velocidadMaxima
    
    constructor(marca, modelo, ano, velocidadMaxima) {
        this.marca = marca
        this.modelo = modelo
        this.ano = ano
        this.velocidadMaxima = velocidadMaxima
    }
    
    acelerar() {
        print(this.marca, this.modelo, "está acelerando")
    }
    
    frenar() {
        print(this.marca, this.modelo, "está frenando")
    }
    
    mostrarInfo() {
        print("Vehículo:", this.marca, this.modelo, "(" + this.ano + ")")
        print("Velocidad máxima:", this.velocidadMaxima, "km/h")
    }
}

// Clase intermedia
class Automovil extends Vehiculo {
    let numeroPuertas
    let tipoCombustible
    
    constructor(marca, modelo, ano, velocidadMaxima, numeroPuertas, tipoCombustible) {
        super.constructor(marca, modelo, ano, velocidadMaxima)
        this.numeroPuertas = numeroPuertas
        this.tipoCombustible = tipoCombustible
    }
    
    encenderMotor() {
        print("Motor del", this.marca, this.modelo, "encendido")
    }
    
    mostrarInfo() {
        super.mostrarInfo()
        print("Puertas:", this.numeroPuertas)
        print("Combustible:", this.tipoCombustible)
    }
}

// Clase específica
class AutoElectrico extends Automovil {
    let capacidadBateria
    let autonomia
    let nivelCarga
    
    constructor(marca, modelo, ano, velocidadMaxima, numeroPuertas, capacidadBateria, autonomia) {
        super.constructor(marca, modelo, ano, velocidadMaxima, numeroPuertas, "Eléctrico")
        this.capacidadBateria = capacidadBateria
        this.autonomia = autonomia
        this.nivelCarga = 100  // Inicia con carga completa
    }
    
    cargarBateria(porcentaje) {
        this.nivelCarga = this.nivelCarga + porcentaje
        if (this.nivelCarga > 100) {
            this.nivelCarga = 100
        }
        print("Batería cargada al", this.nivelCarga + "%")
    }
    
    verificarAutonomia() {
        let autonomiaActual = (this.autonomia * this.nivelCarga) / 100
        print("Autonomía actual:", autonomiaActual, "km")
        return autonomiaActual
    }
    
    modoEconorico() {
        print("Activando modo económico para", this.marca, this.modelo)
        print("Velocidad limitada y consumo optimizado")
    }
    
    mostrarInfo() {
        super.mostrarInfo()
        print("Capacidad batería:", this.capacidadBateria, "kWh")
        print("Autonomía máxima:", this.autonomia, "km")
        print("Nivel de carga:", this.nivelCarga + "%")
    }
}

func main() {
    let tesla = AutoElectrico("Tesla", "Model 3", 2023, 250, 4, 75, 500)
    
    print("=== AUTO ELÉCTRICO ===")
    tesla.mostrarInfo()
    print()
    
    print("=== ACCIONES ===")
    tesla.encenderMotor()
    tesla.acelerar()
    tesla.verificarAutonomia()
    tesla.modoEconorico()
    tesla.cargarBateria(20)  // Intentar cargar más
    tesla.frenar()
}
```

## Maps y Objetos Avanzados

### 1. Maps (Diccionarios)

```r2
func main() {
    // Crear map vacío
    let persona = {}
    
    // Crear map con datos iniciales
    let estudiante = {
        nombre: "Ana García",
        edad: 20,
        carrera: "Ingeniería de Software",
        semestre: 5,
        materias: ["Programación", "Matemáticas", "Física"]
    }
    
    print("Estudiante:", estudiante.nombre)
    print("Carrera:", estudiante.carrera)
    print("Materias:", estudiante.materias)
    
    // Agregar nuevas propiedades
    estudiante.promedio = 8.5
    estudiante.activo = true
    
    // Modificar propiedades existentes
    estudiante.semestre = 6
    
    print("Promedio:", estudiante.promedio)
    print("Semestre actual:", estudiante.semestre)
}
```

### 2. Maps Dinámicos

```r2
func crearInventario() {
    let inventario = {}
    
    // Función para agregar producto
    let agregarProducto = func(codigo, nombre, precio, cantidad) {
        inventario[codigo] = {
            nombre: nombre,
            precio: precio,
            cantidad: cantidad,
            valorTotal: precio * cantidad
        }
        print("Producto agregado:", nombre)
    }
    
    // Función para mostrar inventario
    let mostrarInventario = func() {
        print("=== INVENTARIO ===")
        let total = 0
        
        // No podemos iterar maps directamente en R2Lang aún
        // Simularemos con códigos conocidos
        if (inventario["001"] != null) {
            let prod = inventario["001"]
            print("001 -", prod.nombre, "- Precio:", prod.precio, "- Cantidad:", prod.cantidad)
            total = total + prod.valorTotal
        }
        
        if (inventario["002"] != null) {
            let prod = inventario["002"]
            print("002 -", prod.nombre, "- Precio:", prod.precio, "- Cantidad:", prod.cantidad)
            total = total + prod.valorTotal
        }
        
        if (inventario["003"] != null) {
            let prod = inventario["003"]
            print("003 -", prod.nombre, "- Precio:", prod.precio, "- Cantidad:", prod.cantidad)
            total = total + prod.valorTotal
        }
        
        print("Valor total del inventario:", total)
    }
    
    return {
        agregar: agregarProducto,
        mostrar: mostrarInventario,
        obtener: func(codigo) {
            return inventario[codigo]
        }
    }
}

func main() {
    let inv = crearInventario()
    
    // Agregar productos
    inv.agregar("001", "Laptop", 1500, 10)
    inv.agregar("002", "Mouse", 25, 50)
    inv.agregar("003", "Teclado", 75, 30)
    
    print()
    inv.mostrar()
    
    print()
    let laptop = inv.obtener("001")
    print("Detalles de laptop:", laptop.nombre, "- Stock:", laptop.cantidad)
}
```

### 3. Composición de Objetos

```r2
class Motor {
    let cilindros
    let potencia
    let combustible
    let encendido
    
    constructor(cilindros, potencia, combustible) {
        this.cilindros = cilindros
        this.potencia = potencia
        this.combustible = combustible
        this.encendido = false
    }
    
    encender() {
        if (!this.encendido) {
            this.encendido = true
            print("Motor encendido -", this.potencia, "HP")
        } else {
            print("Motor ya está encendido")
        }
    }
    
    apagar() {
        if (this.encendido) {
            this.encendido = false
            print("Motor apagado")
        } else {
            print("Motor ya está apagado")
        }
    }
    
    obtenerInfo() {
        return {
            cilindros: this.cilindros,
            potencia: this.potencia,
            combustible: this.combustible,
            estado: this.encendido ? "Encendido" : "Apagado"
        }
    }
}

class Sistema {
    let nombre
    let activo
    
    constructor(nombre) {
        this.nombre = nombre
        this.activo = false
    }
    
    activar() {
        this.activo = true
        print("Sistema", this.nombre, "activado")
    }
    
    desactivar() {
        this.activo = false
        print("Sistema", this.nombre, "desactivado")
    }
}

class Auto {
    let marca
    let modelo
    let motor
    let sistemas
    
    constructor(marca, modelo, motor) {
        this.marca = marca
        this.modelo = modelo
        this.motor = motor
        this.sistemas = {
            aire: Sistema("Aire Acondicionado"),
            navegacion: Sistema("Navegación GPS"),
            sonido: Sistema("Sistema de Sonido")
        }
    }
    
    encender() {
        print("Encendiendo", this.marca, this.modelo)
        this.motor.encender()
        
        // Activar sistemas básicos
        this.sistemas.navegacion.activar()
        print("Auto listo para conducir")
    }
    
    apagar() {
        print("Apagando", this.marca, this.modelo)
        
        // Desactivar sistemas
        this.sistemas.aire.desactivar()
        this.sistemas.navegacion.desactivar()
        this.sistemas.sonido.desactivar()
        
        this.motor.apagar()
        print("Auto apagado completamente")
    }
    
    activarAire() {
        this.sistemas.aire.activar()
    }
    
    activarSonido() {
        this.sistemas.sonido.activar()
    }
    
    mostrarEstado() {
        print("=== ESTADO DEL", this.marca, this.modelo, "===")
        
        let infoMotor = this.motor.obtenerInfo()
        print("Motor:", infoMotor.potencia, "HP -", infoMotor.estado)
        
        print("Aire acondicionado:", this.sistemas.aire.activo ? "ON" : "OFF")
        print("Navegación:", this.sistemas.navegacion.activo ? "ON" : "OFF")
        print("Sonido:", this.sistemas.sonido.activo ? "ON" : "OFF")
    }
}

func main() {
    // Crear motor
    let motorV6 = Motor(6, 300, "Gasolina")
    
    // Crear auto con composición
    let miAuto = Auto("Toyota", "Camry", motorV6)
    
    print("=== PRUEBA DEL AUTO ===")
    miAuto.mostrarEstado()
    print()
    
    miAuto.encender()
    print()
    
    miAuto.activarAire()
    miAuto.activarSonido()
    print()
    
    miAuto.mostrarEstado()
    print()
    
    miAuto.apagar()
}
```

## Ejercicios Prácticos

### Ejercicio 1: Sistema de Biblioteca

```r2
class Libro {
    let titulo
    let autor
    let isbn
    let disponible
    let fechaPrestamo
    
    constructor(titulo, autor, isbn) {
        this.titulo = titulo
        this.autor = autor
        this.isbn = isbn
        this.disponible = true
        this.fechaPrestamo = null
    }
    
    prestar() {
        if (this.disponible) {
            this.disponible = false
            this.fechaPrestamo = "Hoy"
            print("Libro", this.titulo, "prestado exitosamente")
            return true
        } else {
            print("Libro", this.titulo, "no está disponible")
            return false
        }
    }
    
    devolver() {
        if (!this.disponible) {
            this.disponible = true
            this.fechaPrestamo = null
            print("Libro", this.titulo, "devuelto exitosamente")
            return true
        } else {
            print("Libro", this.titulo, "ya está disponible")
            return false
        }
    }
    
    mostrarInfo() {
        print("📖", this.titulo, "por", this.autor)
        print("   ISBN:", this.isbn)
        print("   Estado:", this.disponible ? "Disponible" : "Prestado")
        if (!this.disponible) {
            print("   Prestado desde:", this.fechaPrestamo)
        }
    }
}

class Biblioteca {
    let nombre
    let libros
    
    constructor(nombre) {
        this.nombre = nombre
        this.libros = []
    }
    
    agregarLibro(libro) {
        this.libros = this.libros.push(libro)
        print("Libro agregado a", this.nombre)
    }
    
    buscarPorTitulo(titulo) {
        for (let i = 0; i < this.libros.length(); i++) {
            let libro = this.libros[i]
            if (libro.titulo.contains(titulo)) {
                return libro
            }
        }
        return null
    }
    
    mostrarCatalogo() {
        print("=== CATÁLOGO DE", this.nombre, "===")
        print("Total de libros:", this.libros.length())
        print()
        
        for (let i = 0; i < this.libros.length(); i++) {
            this.libros[i].mostrarInfo()
            print()
        }
    }
    
    librosDisponibles() {
        let disponibles = []
        for (let i = 0; i < this.libros.length(); i++) {
            let libro = this.libros[i]
            if (libro.disponible) {
                disponibles = disponibles.push(libro)
            }
        }
        return disponibles
    }
}

func main() {
    let biblioteca = Biblioteca("Biblioteca Central")
    
    // Crear libros
    let libro1 = Libro("El Quijote", "Cervantes", "123456")
    let libro2 = Libro("Cien años de soledad", "García Márquez", "789012")
    let libro3 = Libro("1984", "George Orwell", "345678")
    
    // Agregar a biblioteca
    biblioteca.agregarLibro(libro1)
    biblioteca.agregarLibro(libro2)
    biblioteca.agregarLibro(libro3)
    
    print()
    biblioteca.mostrarCatalogo()
    
    // Realizar préstamos
    print("=== PRÉSTAMOS ===")
    libro1.prestar()
    libro2.prestar()
    
    print()
    print("Libros disponibles:", biblioteca.librosDisponibles().length())
    
    // Devolución
    print()
    print("=== DEVOLUCIONES ===")
    libro1.devolver()
    
    print()
    let libroEncontrado = biblioteca.buscarPorTitulo("1984")
    if (libroEncontrado != null) {
        print("Libro encontrado:")
        libroEncontrado.mostrarInfo()
    }
}
```

## Proyecto del Módulo: Sistema de Gestión Escolar

```r2
// Clase base para personas
class Persona {
    let nombre
    let edad
    let id
    
    constructor(nombre, edad, id) {
        this.nombre = nombre
        this.edad = edad
        this.id = id
    }
    
    mostrarInfo() {
        print("Nombre:", this.nombre)
        print("Edad:", this.edad)
        print("ID:", this.id)
    }
}

// Estudiante hereda de Persona
class Estudiante extends Persona {
    let grado
    let calificaciones
    let materiaInscrita
    
    constructor(nombre, edad, id, grado) {
        super.constructor(nombre, edad, id)
        this.grado = grado
        this.calificaciones = {}
        this.materiaInscrita = []
    }
    
    inscribirMateria(materia) {
        this.materiaInscrita = this.materiaInscrita.push(materia)
        this.calificaciones[materia] = []
        print(this.nombre, "inscrito en", materia)
    }
    
    agregarCalificacion(materia, calificacion) {
        if (this.calificaciones[materia] != null) {
            let califs = this.calificaciones[materia]
            califs = califs.push(calificacion)
            this.calificaciones[materia] = califs
            print("Calificación", calificacion, "agregada en", materia, "para", this.nombre)
        } else {
            print("Error:", this.nombre, "no está inscrito en", materia)
        }
    }
    
    calcularPromedio(materia) {
        if (this.calificaciones[materia] != null) {
            let califs = this.calificaciones[materia]
            if (califs.length() == 0) {
                return 0
            }
            
            let suma = 0
            for (let i = 0; i < califs.length(); i++) {
                suma = suma + califs[i]
            }
            return suma / califs.length()
        }
        return 0
    }
    
    mostrarInfo() {
        super.mostrarInfo()
        print("Grado:", this.grado)
        print("Materias inscritas:", this.materiaInscrita.length())
        
        for (let i = 0; i < this.materiaInscrita.length(); i++) {
            let materia = this.materiaInscrita[i]
            let promedio = this.calcularPromedio(materia)
            print("-", materia + ":", promedio)
        }
    }
}

// Profesor hereda de Persona
class Profesor extends Persona {
    let materiaEspecialidad
    let estudiantes
    let salario
    
    constructor(nombre, edad, id, materiaEspecialidad, salario) {
        super.constructor(nombre, edad, id)
        this.materiaEspecialidad = materiaEspecialidad
        this.estudiantes = []
        this.salario = salario
    }
    
    asignarEstudiante(estudiante) {
        this.estudiantes = this.estudiantes.push(estudiante)
        print("Estudiante", estudiante.nombre, "asignado al profesor", this.nombre)
    }
    
    calificarEstudiante(estudiante, calificacion) {
        estudiante.agregarCalificacion(this.materiaEspecialidad, calificacion)
    }
    
    mostrarEstudiantes() {
        print("Estudiantes del profesor", this.nombre, "(" + this.materiaEspecialidad + "):")
        for (let i = 0; i < this.estudiantes.length(); i++) {
            let est = this.estudiantes[i]
            let promedio = est.calcularPromedio(this.materiaEspecialidad)
            print("-", est.nombre, "(Promedio:", promedio + ")")
        }
    }
    
    mostrarInfo() {
        super.mostrarInfo()
        print("Especialidad:", this.materiaEspecialidad)
        print("Salario:", this.salario)
        print("Estudiantes a cargo:", this.estudiantes.length())
    }
}

// Escuela como contenedor
class Escuela {
    let nombre
    let estudiantes
    let profesores
    
    constructor(nombre) {
        this.nombre = nombre
        this.estudiantes = []
        this.profesores = []
    }
    
    agregarEstudiante(estudiante) {
        this.estudiantes = this.estudiantes.push(estudiante)
        print("Estudiante", estudiante.nombre, "agregado a", this.nombre)
    }
    
    agregarProfesor(profesor) {
        this.profesores = this.profesores.push(profesor)
        print("Profesor", profesor.nombre, "agregado a", this.nombre)
    }
    
    mostrarResumen() {
        print("=== ESCUELA", this.nombre, "===")
        print("Total estudiantes:", this.estudiantes.length())
        print("Total profesores:", this.profesores.length())
        print()
        
        print("ESTUDIANTES:")
        for (let i = 0; i < this.estudiantes.length(); i++) {
            let est = this.estudiantes[i]
            print("-", est.nombre, "(" + est.grado + "° grado)")
        }
        print()
        
        print("PROFESORES:")
        for (let i = 0; i < this.profesores.length(); i++) {
            let prof = this.profesores[i]
            print("-", prof.nombre, "(" + prof.materiaEspecialidad + ")")
        }
    }
}

func main() {
    // Crear escuela
    let escuela = Escuela("Instituto Tecnológico")
    
    // Crear estudiantes
    let ana = Estudiante("Ana García", 16, "EST001", 10)
    let carlos = Estudiante("Carlos López", 15, "EST002", 9)
    let maria = Estudiante("María Rodríguez", 17, "EST003", 11)
    
    // Crear profesores
    let profMath = Profesor("Dr. González", 45, "PROF001", "Matemáticas", 3500)
    let profFisica = Profesor("Dra. Martínez", 38, "PROF002", "Física", 3200)
    
    // Agregar a la escuela
    escuela.agregarEstudiante(ana)
    escuela.agregarEstudiante(carlos)
    escuela.agregarEstudiante(maria)
    escuela.agregarProfesor(profMath)
    escuela.agregarProfesor(profFisica)
    
    print()
    
    // Inscribir estudiantes en materias
    ana.inscribirMateria("Matemáticas")
    ana.inscribirMateria("Física")
    carlos.inscribirMateria("Matemáticas")
    maria.inscribirMateria("Física")
    
    print()
    
    // Asignar estudiantes a profesores
    profMath.asignarEstudiante(ana)
    profMath.asignarEstudiante(carlos)
    profFisica.asignarEstudiante(ana)
    profFisica.asignarEstudiante(maria)
    
    print()
    
    // Calificar estudiantes
    profMath.calificarEstudiante(ana, 95)
    profMath.calificarEstudiante(ana, 88)
    profMath.calificarEstudiante(carlos, 78)
    profMath.calificarEstudiante(carlos, 85)
    
    profFisica.calificarEstudiante(ana, 92)
    profFisica.calificarEstudiante(maria, 89)
    
    print()
    
    // Mostrar información
    escuela.mostrarResumen()
    print()
    
    print("=== DETALLE DE ESTUDIANTES ===")
    ana.mostrarInfo()
    print()
    carlos.mostrarInfo()
    print()
    
    print("=== ESTUDIANTES POR PROFESOR ===")
    profMath.mostrarEstudiantes()
    print()
    profFisica.mostrarEstudiantes()
}
```

## Patrones de Diseño en OOP

### 1. Patrón Factory

```r2
class FabricaVehiculos {
    crearVehiculo(tipo, marca, modelo) {
        if (tipo == "auto") {
            return Auto(marca, modelo, 4)
        } else if (tipo == "moto") {
            return Motocicleta(marca, modelo, 2)
        } else if (tipo == "camion") {
            return Camion(marca, modelo, 6)
        } else {
            print("Tipo de vehículo no soportado")
            return null
        }
    }
}
```

### 2. Patrón Observer (Simulado)

```r2
class Notificador {
    let observadores
    
    constructor() {
        this.observadores = []
    }
    
    suscribir(observador) {
        this.observadores = this.observadores.push(observador)
    }
    
    notificar(mensaje) {
        for (let i = 0; i < this.observadores.length(); i++) {
            this.observadores[i].actualizar(mensaje)
        }
    }
}
```

## Resumen del Módulo

### Conceptos Aprendidos
- ✅ Clases y objetos básicos
- ✅ Constructores y métodos
- ✅ Propiedades de clase
- ✅ Herencia con `extends`
- ✅ Uso de `super` para llamar métodos padre
- ✅ Sobrescritura de métodos
- ✅ Maps y objetos dinámicos
- ✅ Composición de objetos
- ✅ Patrones básicos de diseño

### Habilidades Desarrolladas
- ✅ Diseñar clases efectivas
- ✅ Implementar herencia apropiada
- ✅ Crear jerarquías de objetos
- ✅ Usar composición vs herencia
- ✅ Manejar colecciones de objetos
- ✅ Aplicar principios OOP básicos

### Próximo Módulo

En el **Módulo 4** aprenderás:
- Concurrencia y programación paralela
- Manejo avanzado de errores (try/catch/finally)
- Trabajar con archivos
- Bibliotecas integradas avanzadas

¡Excelente trabajo! Ahora dominas los conceptos fundamentales de orientación a objetos en R2Lang.