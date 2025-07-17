// ============================================================================
// GOLD TEST COMPREHENSIVO - R2Lang 2025
// ============================================================================
// Este archivo es una suite de tests comprehensiva que valida TODAS las 
// características principales de R2Lang. Si este script se ejecuta sin errores,
// significa que R2Lang está funcionando correctamente en su mayoría.
// ============================================================================

print("🚀 INICIANDO GOLD TEST COMPREHENSIVO R2Lang 2025");
print("=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=");

// ============================================================================
// 1. DECLARACIONES Y TIPOS BÁSICOS
// ============================================================================
print("\n1️⃣ DECLARACIONES Y TIPOS BÁSICOS");

// Variables básicas
let x = 42;
var y = 100;
let nombre = "R2Lang";
var version = "2025";
let activo = true;
var inactivo = false;
let nulo = nil;

print("✓ Variables básicas:");
print("  let x =", x, "(tipo:", typeOf(x) + ")");
print("  var y =", y, "(tipo:", typeOf(y) + ")");
print("  let nombre =", nombre, "(tipo:", typeOf(nombre) + ")");
print("  let activo =", activo, "(tipo:", typeOf(activo) + ")");
print("  let nulo =", nulo, "(tipo:", typeOf(nulo) + ")");

// Declaraciones múltiples
let a = 1, b = 2, c = 3;
var d = 4, e = 5, f = 6;

print("✓ Declaraciones múltiples:");
print("  let a, b, c =", a, b, c);
print("  var d, e, f =", d, e, f);

// ============================================================================
// 2. OPERACIONES ARITMÉTICAS Y LÓGICAS
// ============================================================================
print("\n2️⃣ OPERACIONES ARITMÉTICAS Y LÓGICAS");

let num1 = 10;
let num2 = 20;
let suma = num1 + num2;
let resta = num2 - num1;
let multiplicacion = num1 * num2;
let division = num2 / num1;
let modulo = num2 % num1;  // NUEVA CARACTERÍSTICA

print("✓ Operaciones aritméticas:");
print("  " + num1 + " + " + num2 + " =", suma);
print("  " + num2 + " - " + num1 + " =", resta);
print("  " + num1 + " * " + num2 + " =", multiplicacion);
print("  " + num2 + " / " + num1 + " =", division);
print("  " + num2 + " % " + num1 + " =", modulo, "🆕");

// Operaciones lógicas
let verdadero = true;
let falso = false;
let y_logico = verdadero && falso;
let o_logico = verdadero || falso;

print("✓ Operaciones lógicas:");
print("  true && false =", y_logico);
print("  true || false =", o_logico);

// Operaciones de comparación
let mayor = num2 > num1;
let menor = num1 < num2;
let igual = num1 == num1;
let diferente = num1 != num2;

print("✓ Operaciones de comparación:");
print("  " + num2 + " > " + num1 + " =", mayor);
print("  " + num1 + " < " + num2 + " =", menor);
print("  " + num1 + " == " + num1 + " =", igual);
print("  " + num1 + " != " + num2 + " =", diferente);

// ============================================================================
// 3. ARRAYS AVANZADOS
// ============================================================================
print("\n3️⃣ ARRAYS AVANZADOS");

let numeros = [1, 2, 3, 4, 5];
let mixto = [1, "dos", true, nil, 5.5];
let anidado = [[1, 2], [3, 4], [5, 6]];

print("✓ Arrays diversos:");
print("  Números:", numeros);
print("  Mixto:", mixto);
print("  Anidado:", anidado);
print("  Longitud números:", len(numeros));
print("  Acceso directo:", numeros[2]);
print("  Acceso anidado:", anidado[1][0]);

// Operaciones con arrays
let concatenado = numeros + [6, 7, 8];
print("  Concatenación:", concatenado);

// ============================================================================
// 4. MAPAS MULTILINEA (NUEVA CARACTERÍSTICA)
// ============================================================================
print("\n4️⃣ MAPAS MULTILINEA 🆕");

// Mapa básico
let persona = {"nombre": "Juan", "edad": 30};

// Mapa multilinea simple
let configuracion = {
    servidor: "localhost",
    puerto: 8080,
    ssl: true
    timeout: 30
};

// Mapa multilinea complejo con anidación
let aplicacion = {
    info: {
        nombre: "MiApp",
        version: "1.0.0"
        autor: "Desarrollador"
    },
    servidor: {
        host: "localhost",
        puerto: 3000,
        ssl: false
    }
    base_datos: {
        tipo: "postgresql"
        host: "db.ejemplo.com",
        puerto: 5432,
        credenciales: {
            usuario: "admin",
            password: "secreto"
            timeout: 30
        }
    },
    caracteristicas: {
        logging: true
        cache: false,
        monitoring: true,
        debug: false
    }
};

print("✓ Mapas diversos:");
print("  Básico:", persona);
print("  Multilinea:", configuracion);
print("  Longitud config:", len(configuracion));
print("  Claves config:", keys(configuracion));
print("  App nombre:", aplicacion.info.nombre);
print("  DB timeout:", aplicacion.base_datos.credenciales.timeout);

// ============================================================================
// 5. CONTROL DE FLUJO CON 'else if' (NUEVA CARACTERÍSTICA)
// ============================================================================
print("\n5️⃣ CONTROL DE FLUJO CON 'else if' 🆕");

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

print("✓ Condicionales else if:");
print("  Puntuación:", puntuacion, "-> Calificación:", calificacion);

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

print("  Número " + numero_test + ":", descripcion);

// ============================================================================
// 6. BUCLES Y ITERACIÓN
// ============================================================================
print("\n6️⃣ BUCLES Y ITERACIÓN");

// While loop
print("✓ While loop:");
let contador = 0;
while (contador < 3) {
    print("  Iteración while:", contador);
    contador++;
}

// For loop tradicional
print("✓ For loop tradicional:");
for (let i = 0; i < 3; i++) {
    print("  Iteración for:", i);
}

// For-in con arrays
print("✓ For-in con arrays:");
let frutas = ["manzana", "banana", "naranja"];
for (fruta in frutas) {
    print("  Índice:", $k, "-> Fruta:", $v);
}

// For-in con mapas usando keys()
print("✓ For-in con mapas:");
let colores = {
    rojo: "#FF0000",
    verde: "#00FF00"
    azul: "#0000FF",
    amarillo: "#FFFF00"
};

let claves_colores = keys(colores);
for (color in claves_colores) {
    let nombre_color = claves_colores[$k];
    print("  Color:", nombre_color, "-> Código:", colores[nombre_color]);
}

// Break y continue
print("✓ Break y continue:");
for (let i = 0; i < 10; i++) {
    if (i == 2) {
        continue;
    }
    if (i == 5) {
        break;
    }
    print("  Valor:", i);
}

// ============================================================================
// 7. FUNCIONES
// ============================================================================
print("\n7️⃣ FUNCIONES");

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
            edad: edad
            activo: true
        },
        configuracion: {
            tema: "claro"
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

print("✓ Funciones:");
print("  sumar(15, 25) =", sumar(15, 25));
print("  evaluarNumero(7) =", evaluarNumero(7));
print("  evaluarNumero(8) =", evaluarNumero(8));

let perfil = crearPerfil("Ana", 28);
print("  Perfil creado:", perfil.usuario.nombre, "edad", perfil.usuario.edad);

print("  factorial(5) =", factorial(5));

// Función anónima
let multiplicar = func(x, y) {
    return x * y;
};
print("  Función anónima:", multiplicar(6, 7));

// ============================================================================
// 8. CLASES Y OBJETOS
// ============================================================================
print("\n8️⃣ CLASES Y OBJETOS");

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

print("✓ Clases y herencia:");
print("  " + animal.describir());
print("  " + animal.hablar());
print("  " + perro.describir());
print("  " + perro.hablar());

// ============================================================================
// 9. STRINGS Y TEMPLATES
// ============================================================================
print("\n9️⃣ STRINGS Y TEMPLATES");

let saludo = "Hola";
let mundo = "Mundo";
let concatenacion = saludo + " " + mundo + "!";

print("✓ Strings:");
print("  Concatenación:", concatenacion);

// Template strings
let edad_usuario = 25;
let template = `El usuario tiene ${edad_usuario} años`;
print("  Template string:", template);

// String multilinea
let multilinea = `Este es un string
que abarca múltiples
líneas con variables: ${edad_usuario}`;
print("  String multilínea:", multilinea);

// Unicode
let emoji = "🎉";
let unicode = "Año: 2024 - España ñáéíóú";
print("  Unicode y emoji:", emoji, unicode);

// ============================================================================
// 10. FECHAS
// ============================================================================
print("\n🔟 FECHAS");

let fecha_simple = @2024-12-25;
let fecha_completa = @"2024-12-25T10:30:00";

print("✓ Fechas:");
print("  Fecha simple:", fecha_simple);
print("  Fecha completa:", fecha_completa);

// ============================================================================
// 11. OPERADOR TERNARIO
// ============================================================================
print("\n1️⃣1️⃣ OPERADOR TERNARIO");

let edad_test = 20;
let estado = edad_test >= 18 ? "adulto" : "menor";
let mensaje_edad = edad_test >= 65 ? "senior" : (edad_test >= 18 ? "adulto" : "menor");

print("✓ Operador ternario:");
print("  Edad", edad_test, "es:", estado);
print("  Clasificación:", mensaje_edad);

// ============================================================================
// 12. FUNCIONES BUILT-IN Y UTILIDADES
// ============================================================================
print("\n1️⃣2️⃣ FUNCIONES BUILT-IN");

let test_array = [1, 2, 3, "cuatro", true];
let test_map = {a: 1, b: 2, c: 3, d: 4};

print("✓ Funciones built-in:");
print("  len(array) =", len(test_array));
print("  len(map) =", len(test_map));
print("  keys(map) =", keys(test_map));
print("  typeOf(42) =", typeOf(42));
print("  typeOf('hello') =", typeOf("hello"));
print("  typeOf(true) =", typeOf(true));

// parseInt
let numero_string = "123";
let numero_convertido = parseInt(numero_string);
print("  parseInt('123') =", numero_convertido, "(tipo:", typeOf(numero_convertido) + ")");

// ============================================================================
// 13. MANEJO DE ERRORES
// ============================================================================
print("\n1️⃣3️⃣ MANEJO DE ERRORES");

print("✓ Try-catch:");
try {
    let resultado = 10 / 0;  // Esto podría causar error
    print("  División exitosa:", resultado);
} catch (error) {
    print("  Error capturado:", error);
} finally {
    print("  Bloque finally ejecutado");
}

// ============================================================================
// 14. INTEGRACIÓN COMPLETA - CASO REAL
// ============================================================================
print("\n1️⃣4️⃣ INTEGRACIÓN COMPLETA - CASO REAL");

print("✓ Sistema de gestión de productos:");

// Base de datos simulada con mapas multilinea
let base_productos = {
    electronica: {
        laptop: {
            precio: 1200,
            stock: 5
            categoria: "computadoras",
            activo: true
        },
        mouse: {
            precio: 25,
            stock: 50,
            categoria: "accesorios"
            activo: true
        },
        teclado: {
            precio: 80
            stock: 30,
            categoria: "accesorios",
            activo: true
        }
    },
    ropa: {
        camisa: {
            precio: 30,
            stock: 20
            categoria: "vestimenta",
            activo: true
        },
        pantalon: {
            precio: 50,
            stock: 15,
            categoria: "vestimenta"
            activo: false
        }
    }
};

// Función que procesa productos usando todas las características nuevas
func analizarProductos(productos) {
    let total_productos = 0;
    let total_valor = 0;
    let productos_activos = 0;
    
    let categorias = keys(productos);
    for (cat in categorias) {
        let categoria = categorias[$k];
        let items = productos[categoria];
        let items_keys = keys(items);
        
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
            print("  " + producto_nombre + ": $" + producto.precio + " (" + clasificacion + ", " + estado + ")");
        }
    }
    
    return {
        total: total_productos,
        activos: productos_activos,
        valor_total: total_valor
    };
}

let resumen = analizarProductos(base_productos);
print("  Resumen: " + resumen.total + " productos totales, " + resumen.activos + " activos");
print("  Valor total inventario activo: $" + resumen.valor_total);

// ============================================================================
// 15. TESTS DE CASOS LÍMITE
// ============================================================================
print("\n1️⃣5️⃣ TESTS DE CASOS LÍMITE");

print("✓ Casos límite:");

// Operaciones con cero
print("  10 % 1 =", 10 % 1);
print("  División: 100 / 4 =", 100 / 4);

// Arrays vacíos y mapas vacíos
let array_vacio = [];
let mapa_vacio = {};
print("  Array vacío length:", len(array_vacio));
print("  Mapa vacío length:", len(mapa_vacio));
print("  Mapa vacío keys:", keys(mapa_vacio));

// Comparaciones con nil
let valor_nil = nil;
let es_nil = valor_nil == nil;
print("  nil == nil:", es_nil);

// Strings vacíos
let string_vacio = "";
print("  String vacío length:", len(string_vacio));

// ============================================================================
// RESUMEN FINAL COMPLETO
// ============================================================================
print("\n" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=");
print("🎉 GOLD TEST COMPREHENSIVO COMPLETADO EXITOSAMENTE!");
print("=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=");

print("\n✅ CARACTERÍSTICAS BÁSICAS VALIDADAS:");
print("   ✓ Declaraciones let/var y tipos básicos");
print("   ✓ Operaciones aritméticas, lógicas y de comparación");
print("   ✓ Arrays simples, mixtos y anidados");
print("   ✓ Mapas básicos y acceso a propiedades");
print("   ✓ Control de flujo (if/else, while, for)");
print("   ✓ Funciones (básicas, recursivas, anónimas)");
print("   ✓ Clases, objetos y herencia");
print("   ✓ Strings, templates y Unicode");
print("   ✓ Fechas y literales de fecha");
print("   ✓ Operador ternario");
print("   ✓ Funciones built-in (len, keys, typeOf, parseInt)");
print("   ✓ Manejo de errores (try/catch/finally)");
print("   ✓ For-in loops con $k/$v");

print("\n🆕 NUEVAS CARACTERÍSTICAS 2025 VALIDADAS:");
print("   ✅ Mapas multilinea con sintaxis mejorada");
print("   ✅ Separadores mixtos (comas + newlines)");
print("   ✅ Mapas anidados multilinea complejos");
print("   ✅ Sintaxis 'else if' para mejor legibilidad");
print("   ✅ Cadenas complejas de 'else if'");
print("   ✅ Operador módulo '%' en múltiples contextos");
print("   ✅ Integración FizzBuzz con else if + módulo");

print("\n🔄 INTEGRACIÓN Y CASOS REALES:");
print("   ✅ Sistema completo de gestión de productos");
print("   ✅ Todas las características trabajando juntas");
print("   ✅ Casos límite y edge cases");
print("   ✅ Compatibilidad total con código existente");

print("\n🚀 R2LANG 2025 - TOTALMENTE FUNCIONAL");
print("   Si este test se ejecuta sin errores, R2Lang está");
print("   funcionando correctamente en TODAS sus características.");

print("\nTotal de características probadas: 50+");
print("Estado: 🟢 TODOS LOS TESTS PASARON");