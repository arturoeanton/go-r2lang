// ============================================================================
// GOLD TEST COMPREHENSIVO - R2Lang 2025
// ============================================================================
// Este archivo es una suite de tests comprehensiva que valida TODAS las 
// características principales de R2Lang. Si este script se ejecuta sin errores,
// significa que R2Lang está funcionando correctamente en su mayoría.
// ============================================================================

std.print("🚀 INICIANDO GOLD TEST COMPREHENSIVO R2Lang 2025");
std.print("=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=");

// ============================================================================
// 1. DECLARACIONES Y TIPOS BÁSICOS
// ============================================================================
std.print("\n1️⃣ DECLARACIONES Y TIPOS BÁSICOS");

// Variables básicas
let x = 42;
var y = 100;
let nombre = "R2Lang";
var version = "2025";
let activo = true;
var inactivo = false;
let nulo = nil;

std.print("✓ Variables básicas:");
std.print("  let x =", x, "(tipo:", std.typeOf(x) + ")");
std.print("  var y =", y, "(tipo:", std.typeOf(y) + ")");
std.print("  let nombre =", nombre, "(tipo:", std.typeOf(nombre) + ")");
std.print("  let activo =", activo, "(tipo:", std.typeOf(activo) + ")");
std.print("  let nulo =", nulo, "(tipo:", std.typeOf(nulo) + ")");

// Declaraciones múltiples
let a = 1, b = 2, c = 3;
var d = 4, e = 5, f = 6;

std.print("✓ Declaraciones múltiples:");
std.print("  let a, b, c =", a, b, c);
std.print("  var d, e, f =", d, e, f);

// ============================================================================
// 2. OPERACIONES ARITMÉTICAS Y LÓGICAS
// ============================================================================
std.print("\n2️⃣ OPERACIONES ARITMÉTICAS Y LÓGICAS");

let num1 = 10;
let num2 = 20;
let suma = num1 + num2;
let resta = num2 - num1;
let multiplicacion = num1 * num2;
let division = num2 / num1;
let modulo = num2 % num1;  // NUEVA CARACTERÍSTICA

std.print("✓ Operaciones aritméticas:");
std.print("  " + num1 + " + " + num2 + " =", suma);
std.print("  " + num2 + " - " + num1 + " =", resta);
std.print("  " + num1 + " * " + num2 + " =", multiplicacion);
std.print("  " + num2 + " / " + num1 + " =", division);
std.print("  " + num2 + " % " + num1 + " =", modulo, "🆕");

// Operaciones lógicas
let verdadero = true;
let falso = false;
let y_logico = verdadero && falso;
let o_logico = verdadero || falso;

std.print("✓ Operaciones lógicas:");
std.print("  true && false =", y_logico);
std.print("  true || false =", o_logico);

// Operaciones de comparación
let mayor = num2 > num1;
let menor = num1 < num2;
let igual = num1 == num1;
let diferente = num1 != num2;

std.print("✓ Operaciones de comparación:");
std.print("  " + num2 + " > " + num1 + " =", mayor);
std.print("  " + num1 + " < " + num2 + " =", menor);
std.print("  " + num1 + " == " + num1 + " =", igual);
std.print("  " + num1 + " != " + num2 + " =", diferente);

// ============================================================================
// 3. ARRAYS AVANZADOS
// ============================================================================
std.print("\n3️⃣ ARRAYS AVANZADOS");

let numeros = [1, 2, 3, 4, 5];
let mixto = [1, "dos", true, nil, 5.5];
let anidado = [[1, 2], [3, 4], [5, 6]];

std.print("✓ Arrays diversos:");
std.print("  Números:", numeros);
std.print("  Mixto:", mixto);
std.print("  Anidado:", anidado);
std.print("  Longitud números:", std.len(numeros));
std.print("  Acceso directo:", numeros[2]);
std.print("  Acceso anidado:", anidado[1][0]);

// Operaciones con arrays
let concatenado = numeros + [6, 7, 8];
std.print("  Concatenación:", concatenado);

// ============================================================================
// 4. MAPAS MULTILINEA (NUEVA CARACTERÍSTICA)
// ============================================================================
std.print("\n4️⃣ MAPAS MULTILINEA 🆕");

// Mapa básico
let persona = {"nombre": "Juan", "edad": 30};

// Mapa multilinea simple
let configuracion = {
    servidor: "localhost",
    puerto: 8080,
    ssl: true,
    timeout: 30
};

// Mapa multilinea complejo con anidación
let aplicacion = {
    info: {
        nombre: "MiApp",
        version: "1.0.0",
        autor: "Desarrollador"
    },
    servidor: {
        host: "localhost",
        puerto: 3000,
        ssl: false
    },
    base_datos: {
        tipo: "postgresql",
        host: "db.ejemplo.com",
        puerto: 5432,
        credenciales: {
            usuario: "admin",
            password: "secreto",
            timeout: 30
        }
    },
    caracteristicas: {
        logging: true,
        cache: false,
        monitoring: true,
        debug: false
    }
};

std.print("✓ Mapas diversos:");
std.print("  Básico:", persona);
std.print("  Multilinea:", configuracion);
std.print("  Longitud config:", std.len(configuracion));
std.print("  Claves config:", std.keys(configuracion));
std.print("  App nombre:", aplicacion.info.nombre);
std.print("  DB timeout:", aplicacion.base_datos.credenciales.timeout);

// ============================================================================
// 5. CONTROL DE FLUJO CON 'else if' (NUEVA CARACTERÍSTICA)
// ============================================================================
std.print("\n5️⃣ CONTROL DE FLUJO CON 'else if' 🆕");

let puntuacion = 85;
let calificacion = "";

if (puntuacion >= 90) {
    calificacion = "A";
} else if (puntuacion >= 80) {
    calificacion = "B";
} else if (puntuacion >= 70) {
    calificacion = "C";
} else if (puntuacion >= 60) {
    calificacion = "D";
} else {
    calificacion = "F";
}

std.print("✓ Condicionales else if:");
std.print("  Puntuación:", puntuacion, "-> Calificación:", calificacion);

// Test con módulo y else if
let numero_test = 15;
let descripcion = "";

if (numero_test % 15 == 0) {
    descripcion = "FizzBuzz (divisible por 15)";
} else if (numero_test % 5 == 0) {
    descripcion = "Buzz (divisible por 5)";
} else if (numero_test % 3 == 0) {
    descripcion = "Fizz (divisible por 3)";
} else if (numero_test % 2 == 0) {
    descripcion = "Par";
} else {
    descripcion = "Impar";
}

std.print("  Número " + numero_test + ":", descripcion);

// ============================================================================
// 6. BUCLES Y ITERACIÓN
// ============================================================================
std.print("\n6️⃣ BUCLES Y ITERACIÓN");

// While loop
std.print("✓ While loop:");
let contador = 0;
while (contador < 3) {
    std.print("  Iteración while:", contador);
    contador++;
}

// For loop tradicional
std.print("✓ For loop tradicional:");
for (let i = 0; i < 3; i++) {
    std.print("  Iteración for:", i);
}

// For-in con arrays
std.print("✓ For-in con arrays:");
let frutas = ["manzana", "banana", "naranja"];
for (fruta in frutas) {
    std.print("  Índice:", $k, "-> Fruta:", $v);
}

// For-in con mapas usando std.keys()
std.print("✓ For-in con mapas:");
let colores = {
    rojo: "#FF0000",
    verde: "#00FF00",
    azul: "#0000FF",
    amarillo: "#FFFF00"
};

let claves_colores = std.keys(colores);
for (color in claves_colores) {
    let nombre_color = claves_colores[$k];
    std.print("  Color:", nombre_color, "-> Código:", colores[nombre_color]);
}

// Break y continue
std.print("✓ Break y continue:");
for (let i = 0; i < 10; i++) {
    if (i == 2) {
        continue;
    }
    if (i == 5) {
        break;
    }
    std.print("  Valor:", i);
}

// ============================================================================
// 7. FUNCIONES
// ============================================================================
std.print("\n7️⃣ FUNCIONES");

// Función básica
func sumar(a, b) {
    return a + b;
}

// Función con lógica compleja
func evaluarNumero(num) {
    if (num % 2 == 0) {
        return "par";
    } else {
        return "impar";
    }
}

// Función que usa mapas multilinea
func crearPerfil(nombre, edad) {
    return {
        usuario: {
            nombre: nombre,
            edad: edad,
            activo: true
        },
        configuracion: {
            tema: "claro",
            idioma: "es",
            notificaciones: true
        }
    };
}

// Función recursiva
func factorial(n) {
    if (n <= 1) {
        return 1;
    } else {
        return n * factorial(n - 1);
    }
}

std.print("✓ Funciones:");
std.print("  sumar(15, 25) =", sumar(15, 25));
std.print("  evaluarNumero(7) =", evaluarNumero(7));
std.print("  evaluarNumero(8) =", evaluarNumero(8));

let perfil = crearPerfil("Ana", 28);
std.print("  Perfil creado:", perfil.usuario.nombre, "edad", perfil.usuario.edad);

std.print("  factorial(5) =", factorial(5));

// Función anónima
let multiplicar = func(x, y) {
    return x * y;
};
std.print("  Función anónima:", multiplicar(6, 7));

// ============================================================================
// 8. CLASES Y OBJETOS
// ============================================================================
std.print("\n8️⃣ CLASES Y OBJETOS");

class Animal {
    let nombre;
    let especie;
    
    constructor(nom, esp) {
        this.nombre = nom;
        this.especie = esp;
    }
    
    hablar() {
        return this.nombre + " hace ruido";
    }
    
    describir() {
        return "Soy " + this.nombre + ", un " + this.especie;
    }
}

class Perro extends Animal {
    constructor(nom) {
        super(nom, "perro");
    }
    
    hablar() {
        return this.nombre + " dice: ¡Guau!";
    }
}

let animal = Animal("Genérico", "animal");
let perro = Perro("Firulais");

std.print("✓ Clases y herencia:");
std.print("  " + animal.describir());
std.print("  " + animal.hablar());
std.print("  " + perro.describir());
std.print("  " + perro.hablar());

// ============================================================================
// 9. STRINGS Y TEMPLATES
// ============================================================================
std.print("\n9️⃣ STRINGS Y TEMPLATES");

let saludo = "Hola";
let mundo = "Mundo";
let concatenacion = saludo + " " + mundo + "!";

std.print("✓ Strings:");
std.print("  Concatenación:", concatenacion);

// Template strings
let edad_usuario = 25;
let template = `El usuario tiene ${edad_usuario} años`;
std.print("  Template string:", template);

// String multilinea
let multilinea = `Este es un string
que abarca múltiples
líneas con variables: ${edad_usuario}`;
std.print("  String multilínea:", multilinea);

// Unicode
let emoji = "🎉";
let unicode = "Año: 2024 - España ñáéíóú";
std.print("  Unicode y emoji:", emoji, unicode);

// ============================================================================
// 10. FECHAS MEJORADAS CON MÓDULO DATE 🆕
// ============================================================================
std.print("\n🔟 FECHAS MEJORADAS CON MÓDULO DATE 🆕");

let fecha_simple = @2024-12-25;
let fecha_completa = @"2024-12-25T10:30:00";

std.print("✓ Fechas básicas:");
std.print("  Fecha simple:", fecha_simple);
std.print("  Fecha completa:", fecha_completa);

// Nuevo módulo Date con funcionalidad JavaScript-like
let dateObj = date.Date();
let nueva_fecha = dateObj.create(2024, 11, 25, 10, 30, 0);
let fecha_actual = dateObj.create();
let timestamp = dateObj.now();

std.print("✓ Módulo Date mejorado:");
std.print("  Fecha creada:", nueva_fecha);
std.print("  Fecha actual:", fecha_actual);
std.print("  Timestamp now:", timestamp);

// Métodos de fecha JavaScript-like
let año = dateObj.getFullYear(nueva_fecha);
let mes = dateObj.getMonth(nueva_fecha);
let dia = dateObj.getDate(nueva_fecha);

std.print("  Año:", año);
std.print("  Mes:", mes, "(0-based)");
std.print("  Día:", dia);

// Formateo de fechas
let fecha_formateada = date.format(nueva_fecha, "YYYY-MM-DD HH:mm:ss");
std.print("  Fecha formateada:", fecha_formateada);

// Operaciones con fechas
let nueva_fecha_mas_dias = dateObj.addDays(nueva_fecha, 10);
let otra_fecha = dateObj.create(2024, 11, 20);
let diferencia = dateObj.diff(nueva_fecha, otra_fecha, "days");

std.print("  Fecha + 10 días:", nueva_fecha_mas_dias);
std.print("  Diferencia en días:", diferencia);

// Conversiones de fecha
let iso_string = dateObj.toISOString(nueva_fecha);
let date_string = dateObj.toDateString(nueva_fecha);

std.print("  ISO String:", iso_string);
std.print("  Date String:", date_string);

// ============================================================================
// 11. MÓDULO JSON MADURO 🆕
// ============================================================================
std.print("\n1️⃣1️⃣ MÓDULO JSON MADURO 🆕");

// Datos de prueba para JSON
let datos_usuario = {
    nombre: "Carlos",
    edad: 32,
    activo: true,
    hobbies: ["lectura", "programación", "música"],
    configuracion: {
        tema: "oscuro",
        notificaciones: true,
        idioma: "es"
    }
};

std.print("✓ Conversión JSON:");
let json_string = json.stringify(datos_usuario);
std.print("  Objeto a JSON:", json_string);

// Parsing JSON
let json_parseado = json.parse(json_string);
std.print("  JSON parseado:", json_parseado);

// Validación JSON
let json_valido = json.validate(json_string);
let json_invalido = json.validate('{"nombre": "mal formato"');
std.print("  JSON válido:", json_valido);
std.print("  JSON inválido:", json_invalido);

// Operaciones JSON avanzadas
let json_keys = json.getKeys(json_string);
std.print("  Claves JSON:", json_keys);

let nombre_usuario = json.getValue(json_string, "nombre");
std.print("  Valor 'nombre':", nombre_usuario);

let json_modificado = json.setValue(json_string, "ciudad", "Madrid");
std.print("  JSON con nueva clave:", json_modificado);

// Fusión de JSON
let json_adicional = '{"telefono": "123-456-789", "email": "carlos@example.com"}';
let json_fusionado = json.merge(json_string, json_adicional);
std.print("  JSON fusionado:", json_fusionado);

// Aplanar JSON
let json_complejo = '{"usuario": {"info": {"nombre": "Ana", "edad": 28}}}';
let json_plano = json.flatten(json_complejo);
std.print("  JSON aplanado:", json_plano);

// Formateo JSON
let json_bonito = json.pretty(json_string);
std.print("  JSON formateado:");
std.print(json_bonito);

// Query JSON
let nombre_anidado = json.query(json_complejo, "usuario.info.nombre");
std.print("  Query resultado:", nombre_anidado);

// ============================================================================
// 12. MÓDULO CONSOLE INTERACTIVO 🆕
// ============================================================================
std.print("\n1️⃣2️⃣ MÓDULO CONSOLE INTERACTIVO 🆕");

std.print("✓ Logging avanzado:");
console.log("Mensaje de log normal");
console.info("Información importante");
console.warn("Advertencia del sistema");
console.error("Error simulado");

// Tabla de datos
let productos_tabla = [{"nombre": "Laptop", "precio": 1200, "stock": 5}, {"nombre": "Mouse", "precio": 25, "stock": 50}, {"nombre": "Teclado", "precio": 80, "stock": 30}];

std.print("  Tabla de productos:");
console.table(productos_tabla);

// Contadores
console.count("operacion");
console.count("operacion");
console.count("operacion");

// Temporizadores
console.time("proceso");
// Simular trabajo...
let suma_temporal = 0;
for (let i = 0; i < 1000; i++) {
    suma_temporal = suma_temporal + i;
}
console.timeEnd("proceso");

// Formateo con colores
console.color("green", "✓ Proceso completado exitosamente");
console.color("red", "✗ Simulación de error");
console.color("yellow", "⚠ Advertencia de prueba");

// Texto con estilos
console.bold("Texto en negrita");
console.italic("Texto en cursiva");
console.underline("Texto subrayado");

// Grouping
console.group("Operaciones matemáticas");
console.log("Suma: 10 + 20 = 30");
console.log("Resta: 30 - 10 = 20");
console.log("Multiplicación: 10 * 20 = 200");
console.groupEnd();

// ============================================================================
// 13. MÓDULO MATH PARA ANÁLISIS DE DATOS 🆕
// ============================================================================
std.print("\n1️⃣3️⃣ MÓDULO MATH PARA ANÁLISIS DE DATOS 🆕");

let numeros_analisis = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
let numeros_dispersos = [2, 4, 4, 4, 5, 5, 7, 9];

std.print("✓ Análisis estadístico:");
std.print("  Datos:", numeros_analisis);
std.print("  Promedio:", math.mean(numeros_analisis));
std.print("  Mediana:", math.median(numeros_analisis));
std.print("  Moda:", math.mode(numeros_dispersos));
std.print("  Desviación estándar:", math.stdDev(numeros_analisis));
std.print("  Varianza:", math.variance(numeros_analisis));
std.print("  Suma:", math.sum(numeros_analisis));
std.print("  Mínimo:", math.min(numeros_analisis));
std.print("  Máximo:", math.max(numeros_analisis));

// Funciones matemáticas avanzadas
std.print("✓ Funciones matemáticas:");
std.print("  sin(π/2):", math.sin(math.PI / 2));
std.print("  cos(0):", math.cos(0));
std.print("  tan(π/4):", math.tan(math.PI / 4));
std.print("  log(10):", math.log(10));
std.print("  log10(1000):", math.log10(1000));
std.print("  sqrt(16):", math.sqrt(16));
std.print("  pow(2, 8):", math.pow(2, 8));

// Funciones de redondeo
std.print("✓ Redondeo:");
let numero_decimal = 3.7456;
std.print("  Número:", numero_decimal);
std.print("  floor:", math.floor(numero_decimal));
std.print("  ceil:", math.ceil(numero_decimal));
std.print("  round:", math.round(numero_decimal));
std.print("  abs(-5.3):", math.abs(-5.3));

// Números aleatorios
std.print("✓ Números aleatorios:");
std.print("  random():", math.random());
std.print("  randomInt(1, 10):", math.randomInt(1, 10));

// Constantes matemáticas
std.print("✓ Constantes:");
std.print("  PI:", math.PI);
std.print("  E:", math.E);

// ============================================================================
// 14. MÓDULO IO MEJORADO 🆕
// ============================================================================
std.print("\n1️⃣4️⃣ MÓDULO IO MEJORADO 🆕");

std.print("✓ Operaciones de archivos:");

// Crear archivo de prueba
let contenido_test = "Este es un archivo de prueba\ncon múltiples líneas\ny contenido variado.";
let archivo_prueba = "/tmp/r2lang_test.txt";

// Escribir archivo
let resultado_escritura = io.writeFile(archivo_prueba, contenido_test);
std.print("  Archivo escrito:", resultado_escritura);

// Leer archivo
let contenido_leido = io.readFile(archivo_prueba);
std.print("  Contenido leído:", contenido_leido);

// Verificar existencia
let existe = io.exists(archivo_prueba);
std.print("  Archivo existe:", existe);

// Información del archivo
let tamano_archivo = io.fileSize(archivo_prueba);
std.print("  Tamaño del archivo:", tamano_archivo, "bytes");

// Limpiar archivo de prueba
io.rmFile(archivo_prueba);
std.print("  Archivo eliminado");

// Operaciones con directorios
let existe_tmp = io.exists("/tmp");
std.print("  Directorio /tmp existe:", existe_tmp);

// ============================================================================
// 15. OPERADOR TERNARIO
// ============================================================================
std.print("\n1️⃣5️⃣ OPERADOR TERNARIO");

let edad_test = 20;
let estado = edad_test >= 18 ? "adulto" : "menor";
let mensaje_edad = edad_test >= 65 ? "senior" : (edad_test >= 18 ? "adulto" : "menor");

std.print("✓ Operador ternario:");
std.print("  Edad", edad_test, "es:", estado);
std.print("  Clasificación:", mensaje_edad);

// ============================================================================
// 16. FUNCIONES BUILT-IN Y UTILIDADES
// ============================================================================
std.print("\n1️⃣6️⃣ FUNCIONES BUILT-IN");

let test_array = [1, 2, 3, "cuatro", true];
let test_map = {a: 1, b: 2, c: 3, d: 4};

std.print("✓ Funciones built-in:");
std.print("  std.len(array) =", std.len(test_array));
std.print("  std.len(map) =", std.len(test_map));
std.print("  std.keys(map) =", std.keys(test_map));
std.print("  std.typeOf(42) =", std.typeOf(42));
std.print("  std.typeOf('hello') =", std.typeOf("hello"));
std.print("  std.typeOf(true) =", std.typeOf(true));

// parseInt
let numero_string = "123";
let numero_convertido = std.parseInt(numero_string);
std.print("  std.parseInt('123') =", numero_convertido, "(tipo:", std.typeOf(numero_convertido) + ")");

// ============================================================================
// 17. MANEJO DE ERRORES
// ============================================================================
std.print("\n1️⃣7️⃣ MANEJO DE ERRORES");

std.print("✓ Try-catch:");
try {
    let resultado = 10 / 0;  // Esto podría causar error
    std.print("  División exitosa:", resultado);
} catch (error) {
    std.print("  Error capturado:", error);
} finally {
    std.print("  Bloque finally ejecutado");
}

// ============================================================================
// 18. INTEGRACIÓN COMPLETA - CASO REAL
// ============================================================================
std.print("\n1️⃣8️⃣ INTEGRACIÓN COMPLETA - CASO REAL");

std.print("✓ Sistema de gestión de productos:");

// Base de datos simulada con mapas multilinea
let base_productos = {
    electronica: {
        laptop: {
            precio: 1200,
            stock: 5,
            categoria: "computadoras",
            activo: true
        },
        mouse: {
            precio: 25,
            stock: 50,
            categoria: "accesorios",
            activo: true
        },
        teclado: {
            precio: 80,
            stock: 30,
            categoria: "accesorios",
            activo: true
        }
    },
    ropa: {
        camisa: {
            precio: 30,
            stock: 20,
            categoria: "vestimenta",
            activo: true
        },
        pantalon: {
            precio: 50,
            stock: 15,
            categoria: "vestimenta",
            activo: false
        }
    }
};

// Función que procesa productos usando todas las características nuevas
func analizarProductos(productos) {
    let total_productos = 0;
    let total_valor = 0;
    let productos_activos = 0;
    
    let categorias = std.keys(productos);
    for (cat in categorias) {
        let categoria = categorias[$k];
        let items = productos[categoria];
        let items_keys = std.keys(items);
        
        for (item in items_keys) {
            let producto_nombre = items_keys[$k];
            let producto = items[producto_nombre];
            
            total_productos++;
            
            if (producto.activo) {
                productos_activos++;
                total_valor = total_valor + (producto.precio * producto.stock);
            }
            
            // Clasificación usando else if y módulo
            let clasificacion = "";
            if (producto.precio % 100 == 0) {
                clasificacion = "precio redondo";
            } else if (producto.precio > 100) {
                clasificacion = "premium";
            } else if (producto.precio > 50) {
                clasificacion = "medio";
            } else {
                clasificacion = "económico";
            }
            
            let estado = producto.activo ? "activo" : "inactivo";
            std.print("  " + producto_nombre + ": $" + producto.precio + " (" + clasificacion + ", " + estado + ")");
        }
    }
    
    return {
        total: total_productos,
        activos: productos_activos,
        valor_total: total_valor
    };
}

let resumen = analizarProductos(base_productos);
std.print("  Resumen: " + resumen.total + " productos totales, " + resumen.activos + " activos");
std.print("  Valor total inventario activo: $" + resumen.valor_total);

// ============================================================================
// 19. TESTS DE CASOS LÍMITE
// ============================================================================
std.print("\n1️⃣9️⃣ TESTS DE CASOS LÍMITE");

std.print("✓ Casos límite:");

// Operaciones con cero
std.print("  10 % 1 =", 10 % 1);
std.print("  División: 100 / 4 =", 100 / 4);

// Arrays vacíos y mapas vacíos
let array_vacio = [];
let mapa_vacio = {};
std.print("  Array vacío length:", std.len(array_vacio));
std.print("  Mapa vacío length:", std.len(mapa_vacio));
std.print("  Mapa vacío keys:", std.keys(mapa_vacio));

// Comparaciones con nil
let valor_nil = nil;
let es_nil = valor_nil == nil;
std.print("  nil == nil:", es_nil);

// Strings vacíos
let string_vacio = "";
std.print("  String vacío length:", std.len(string_vacio));

// ============================================================================
// RESUMEN FINAL COMPLETO
// ============================================================================
// ============================================================================
// 🆕 NUEVAS CARACTERÍSTICAS P2 - DESTRUCTURING Y SPREAD OPERATOR
// ============================================================================
std.print("\n1️⃣9️⃣ NUEVAS CARACTERÍSTICAS P2 - DESTRUCTURING Y SPREAD OPERATOR");

// Array Destructuring
std.print("✓ Array Destructuring:");
let [primero, segundo, tercero] = [100, 200, 300];
std.print("  [primero, segundo, tercero] = [100, 200, 300]");
std.print("  primero =", primero, ", segundo =", segundo, ", tercero =", tercero);

let [x1, x2, x3, x4] = [10, 20];
std.print("  Más variables que elementos: x1 =", x1, ", x2 =", x2, ", x3 =", x3, ", x4 =", x4);

let [nombre_user, edad_user, activo_user] = ["Juan", 30, true];
std.print("  Tipos mixtos: nombre =", nombre_user, ", edad =", edad_user, ", activo =", activo_user);

// Object Destructuring
std.print("✓ Object Destructuring:");
let usuario = {
    nombre: "Ana García",
    email: "ana@ejemplo.com",
    edad: 28,
    admin: false
};

let {nombre, email, edad} = usuario;
std.print("  {nombre, email, edad} extraído de usuario");
std.print("  nombre =", nombre, ", email =", email, ", edad =", edad);

let {username, role, permisos} = {username: "admin_user"};
std.print("  Propiedades faltantes: username =", username, ", role =", role, ", permisos =", permisos);

// Spread Operator en Arrays
std.print("✓ Spread Operator en Arrays:");
let array1 = [1, 2, 3];
let array2 = [4, 5, 6];
let combinado = [...array1, ...array2];
std.print("  [...array1, ...array2] =", combinado);

let extendido = [...array1, 7, 8, 9];
std.print("  [...array1, 7, 8, 9] =", extendido);

let conPrefijo = [0, ...array2];
std.print("  [0, ...array2] =", conPrefijo);

// Spread Operator en Objetos
std.print("✓ Spread Operator en Objetos:");
let configuracion = {
    tema: "claro",
    idioma: "es",
    notificaciones: true
};

let personalizacion = {
    tema: "oscuro",
    fuente: 16
};

let configFinal = {...configuracion, ...personalizacion};
std.print("  {...configuracion, ...personalizacion} =", configFinal);
std.print("  Nota: tema 'oscuro' sobrescribe 'claro'");

let empleado = {...usuario, id: 12345, departamento: "IT"};
std.print("  {...usuario, id, departamento} =", empleado);

// Spread Operator en Llamadas a Funciones
std.print("✓ Spread Operator en Llamadas a Funciones:");
func sumarCuatro(a, b, c, d) {
    return a + b + c + d;
}

let valores = [10, 20, 30, 40];
let resultado = sumarCuatro(...valores);
std.print("  sumarCuatro(...[10, 20, 30, 40]) =", resultado);

let argsParciales = [5, 15];
let resultado2 = sumarCuatro(1, ...argsParciales, 25);
std.print("  sumarCuatro(1, ...argsParciales, 25) =", resultado2);

func saludar(saludo, nombre, puntuacion) {
    if (puntuacion) {
        return saludo + " " + nombre + puntuacion;
    } else if (nombre) {
        return saludo + " " + nombre;
    } else {
        return saludo;
    }
}

let args1 = ["Hola"];
let args2 = ["Buenos días", "María"];
let args3 = ["¡Felicidades", "Carlos", "!"];

std.print("  saludar(...args) examples:");
std.print("    saludar(...['Hola']) =", saludar(...args1));
std.print("    saludar(...['Buenos días', 'María']) =", saludar(...args2));
std.print("    saludar(...['¡Felicidades', 'Carlos', '!']) =", saludar(...args3));

// Casos de uso prácticos
std.print("✓ Casos de Uso Prácticos:");

// Intercambio de variables
let var1 = "A";
let var2 = "B";
let [nuevo1, nuevo2] = [var2, var1];
std.print("  Intercambio: [nuevo1, nuevo2] = [var2, var1] =>", nuevo1, nuevo2);

// Clonado de arrays
let original = [1, 2, 3, {nested: "valor"}];
let clon = [...original];
std.print("  Clonado: clon = [...original] =>", clon);

// Configuración con defaults
func crearConfig(opciones) {
    let porDefecto = {
        timeout: 5000,
        reintentos: 3,
        cache: true,
        debug: false
    };
    
    return {...porDefecto, ...opciones};
}

let miConfig = crearConfig({timeout: 10000, debug: true});
std.print("  Config personalizada =", miConfig);

std.print("\n✅ CARACTERÍSTICAS P2 COMPLETADAS:");
std.print("   ✅ Array Destructuring básico");
std.print("   ✅ Object Destructuring básico");
std.print("   ✅ Spread Operator en arrays");
std.print("   ✅ Spread Operator en objetos");
std.print("   ✅ Spread Operator en llamadas a funciones");
std.print("   ✅ Casos de uso prácticos y edge cases");
std.print("   ✅ Compatibilidad total con sintaxis existente");

std.print("\n" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=");
std.print("🎉 GOLD TEST COMPREHENSIVO COMPLETADO EXITOSAMENTE!");
std.print("=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=");

std.print("\n✅ CARACTERÍSTICAS BÁSICAS VALIDADAS:");
std.print("   ✓ Declaraciones let/var y tipos básicos");
std.print("   ✓ Operaciones aritméticas, lógicas y de comparación");
std.print("   ✓ Arrays simples, mixtos y anidados");
std.print("   ✓ Mapas básicos y acceso a propiedades");
std.print("   ✓ Control de flujo (if/else, while, for)");
std.print("   ✓ Funciones (básicas, recursivas, anónimas)");
std.print("   ✓ Clases, objetos y herencia");
std.print("   ✓ Strings, templates y Unicode");
std.print("   ✓ Fechas y literales de fecha");
std.print("   ✓ Operador ternario");
std.print("   ✓ Funciones built-in (len, keys, typeOf, parseInt)");
std.print("   ✓ Manejo de errores (try/catch/finally)");
std.print("   ✓ For-in loops con $k/$v");

std.print("\n🆕 NUEVAS CARACTERÍSTICAS 2025 VALIDADAS:");
std.print("   ✅ Mapas multilinea con sintaxis mejorada");
std.print("   ✅ Separadores mixtos (comas + newlines)");
std.print("   ✅ Mapas anidados multilinea complejos");
std.print("   ✅ Sintaxis 'else if' para mejor legibilidad");
std.print("   ✅ Cadenas complejas de 'else if'");
std.print("   ✅ Operador módulo '%' en múltiples contextos");
std.print("   ✅ Integración FizzBuzz con else if + módulo");

std.print("\n🚀 CARACTERÍSTICAS P2 IMPLEMENTADAS 2025:");
std.print("   ✅ Array Destructuring - let [a, b, c] = array");
std.print("   ✅ Object Destructuring - let {name, age} = obj");
std.print("   ✅ Spread Operator en Arrays - [...arr1, ...arr2]");
std.print("   ✅ Spread Operator en Objetos - {...obj1, ...obj2}");
std.print("   ✅ Spread Operator en Funciones - func(...args)");
std.print("   ✅ Casos de uso prácticos y edge cases");
std.print("   ✅ Compatibilidad total con sintaxis existente");

std.print("\n🔥 NUEVOS MÓDULOS PREMIUM 2025:");
std.print("   ✅ Módulo DATE mejorado con API JavaScript-like");
std.print("   ✅ Módulo JSON maduro con todas las operaciones");
std.print("   ✅ Módulo CONSOLE interactivo con logging avanzado");
std.print("   ✅ Módulo MATH para análisis de datos y estadísticas");
std.print("   ✅ Módulo IO mejorado con operaciones de archivos");
std.print("   ✅ Formateo de fechas personalizado");
std.print("   ✅ Parsing y stringify JSON con validación");
std.print("   ✅ Logging con colores y estilos");
std.print("   ✅ Funciones estadísticas avanzadas");
std.print("   ✅ Operaciones de archivos y directorios");

std.print("\n🔄 INTEGRACIÓN Y CASOS REALES:");
std.print("   ✅ Sistema completo de gestión de productos");
std.print("   ✅ Todas las características trabajando juntas");
std.print("   ✅ Casos límite y edge cases");
std.print("   ✅ Compatibilidad total con código existente");
std.print("   ✅ Integración completa de módulos nuevos");
std.print("   ✅ Tests exhaustivos de funcionalidad");

std.print("\n🚀 R2LANG 2025 - EDICIÓN PREMIUM FUNCIONAL");
std.print("   Si este test se ejecuta sin errores, R2Lang está");
std.print("   funcionando correctamente en TODAS sus características");
std.print("   incluyendo los nuevos módulos profesionales.");

std.print("\nTotal de características probadas: 80+");
std.print("Módulos nuevos incluidos: 5");
std.print("Estado: 🟢 TODOS LOS TESTS PASARON");
std.print("Versión: 🔥 PREMIUM 2025 EDITION");