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
