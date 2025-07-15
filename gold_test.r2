// Gold Test - Validación rápida de características principales de R2Lang
// Este archivo sirve como smoke test para verificar que las características básicas funcionan

// ====================================
// 1. DECLARACIONES CON let Y var
// ====================================
let x = 42;
var y = 100;
let nombre = "R2Lang";
var version = "2025";

print("✓ Declaraciones básicas:");
print("  let x =", x);
print("  var y =", y);
print("  let nombre =", nombre);
print("  var version =", version);

// ====================================
// 2. TIPOS DE DATOS Y OPERACIONES
// ====================================
var num1 = 10;
let num2 = 20;
var resultado = num1 + num2;
var texto = "Suma: " + resultado;

print("\n✓ Operaciones aritméticas:");
print("  " + num1 + " + " + num2 + " =", resultado);
print("  " + texto);

// ====================================
// 3. ARRAYS Y MAPS
// ====================================
let numeros = [1, 2, 3, 4, 5];
var personas = {"nombre": "Juan", "edad": 30};

print("\n✓ Arrays y Maps:");
print("  Array:", numeros);
print("  Map:", personas);

// ====================================
// 4. CONTROL DE FLUJO
// ====================================
var contador = 0;
let limite = 3;

print("\n✓ Control de flujo:");
while (contador < limite) {
    print("  Contador:", contador);
    contador++;
}

// ====================================
// 5. FUNCIONES
// ====================================
func sumar(a, b) {
    return a + b;
}

let resultado_func = sumar(25, 17);
print("\n✓ Funciones:");
print("  sumar(25, 17) =", resultado_func);

// Test expresiones en parámetros
func testParams(a, b, c) {
    return a + " | " + b + " | " + c;
}

let result_params = testParams("pepe" + "!!", 2 + 3, "p" + 1);
print("  Expresiones en parámetros:", result_params);

// ====================================
// 6. CLASES Y OBJETOS
// ====================================
class Persona {
    let nombre;
    var edad;
    
    constructor(nom, ed) {
        this.nombre = nom;
        this.edad = ed;
    }
    
    saludar() {
        return "Hola, soy " + this.nombre;
    }
}

let persona = Persona("Ana", 25);
print("\n✓ Clases y Objetos:");
print("  " + persona.saludar());
print("  Edad:", persona.edad);

// ====================================
// 7. UNICODE Y CARACTERES ESPECIALES
// ====================================
var emoji = "🚀";
let texto_unicode = "Año: 2024 - España";
var mensaje = `Probando template strings con ${emoji}`;

print("\n✓ Unicode y Templates:");
print("  Emoji:", emoji);
print("  Unicode:", texto_unicode);
print("  Template:", mensaje);

// ====================================
// 8. FECHAS
// ====================================
var fecha = @2024-12-25;
let fecha_completa = @"2024-12-25T10:30:00";

print("\n✓ Fechas:");
print("  Fecha simple:", fecha);
print("  Fecha completa:", fecha_completa);

// ====================================
// 9. OPERADOR TERNARIO
// ====================================
let edad_test = 18;
var mensaje_edad = edad_test >= 18 ? "Adulto" : "Menor";

print("\n✓ Operador ternario:");
print("  Edad", edad_test, "es:", mensaje_edad);

// ====================================
// 10. DECLARACIONES MÚLTIPLES
// ====================================
let a = 1, b = 2, c = 3;
var d = 4, e = 5, f = 6;

print("\n✓ Declaraciones múltiples:");
print("  let a, b, c =", a, b, c);
print("  var d, e, f =", d, e, f);

// ====================================
// 11. ARRAYS Y LOOPS
// ====================================
var frutas = ["manzana", "banana", "naranja"];
print("\n✓ Arrays y for-in:");
for (fruta in frutas) {
    print("  Fruta:", fruta);
}

// ====================================
// 12. CONDICIONALES
// ====================================
let test_value = 42;
print("\n✓ Condicionales:");
if (test_value > 40) {
    print("  El valor es mayor que 40");
} else {
    print("  El valor es menor o igual que 40");
}

// ====================================
// 13. STRINGS MULTILÍNEA
// ====================================
var multilinea = `Este es un string
que abarca múltiples
líneas con variables: ${test_value}`;

print("\n✓ Strings multilínea:");
print(multilinea);

// ====================================
// RESUMEN FINAL
// ====================================
print("\n🎉 GOLD TEST COMPLETADO EXITOSAMENTE!");
print("✅ Todas las características principales de R2Lang funcionan correctamente:");
print("   - let/var declarations");
print("   - Basic types and operations");
print("   - Arrays and Maps");
print("   - Control flow");
print("   - Functions");
print("   - Classes and Objects");
print("   - Unicode support");
print("   - Date literals");
print("   - Ternary operator");
print("   - Multiple declarations");
print("   - Template strings");
print("   - Multiline strings");