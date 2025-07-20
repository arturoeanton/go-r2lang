// =================================================================
// R2Lang Example 45: Destructuring y Spread Operator (P2 Features)
// =================================================================
// Este ejemplo demuestra las nuevas características P2 implementadas:
// - Array Destructuring
// - Object Destructuring  
// - Spread Operator (...) en arrays, objetos y llamadas a funciones

std.print("=== R2Lang Example 45: Destructuring & Spread Operator ===");
std.print("");

// ----------------------------------------------------------------
// 1. ARRAY DESTRUCTURING
// ----------------------------------------------------------------
std.print("1. ARRAY DESTRUCTURING:");
std.print("  Permite extraer elementos de arrays en variables individuales");
std.print("");

// Ejemplo básico
let numbers = [1, 2, 3, 4, 5];
let [first, second, third] = numbers;
std.print("Array original:", numbers);
std.print("first =", first);   // 1
std.print("second =", second); // 2  
std.print("third =", third);   // 3
std.print("");

// Más variables que elementos
let [a, b, c, d, e, f] = [10, 20];
std.print("Destructuring con más variables que elementos:");
std.print("a =", a);  // 10
std.print("b =", b);  // 20
std.print("c =", c);  // nil
std.print("d =", d);  // nil
std.print("");

// Tipos mixtos
let [name, age, active] = ["Alice", 25, true];
std.print("Destructuring con tipos mixtos:");
std.print("name =", name);     // "Alice"
std.print("age =", age);       // 25
std.print("active =", active); // true
std.print("");

// ----------------------------------------------------------------
// 2. OBJECT DESTRUCTURING  
// ----------------------------------------------------------------
std.print("2. OBJECT DESTRUCTURING:");
std.print("  Permite extraer propiedades de objetos en variables");
std.print("");

// Ejemplo básico
let user = {
    name: "John Doe",
    email: "john@example.com",
    age: 30,
    isActive: true
};

let {name, email, age} = user;
std.print("Object original:", user);
std.print("name =", name);   // "John Doe"
std.print("email =", email); // "john@example.com"
std.print("age =", age);     // 30
std.print("");

// Propiedades faltantes
let {username, password, role} = {username: "admin"};
std.print("Destructuring con propiedades faltantes:");
std.print("username =", username); // "admin"
std.print("password =", password); // nil
std.print("role =", role);         // nil
std.print("");

// ----------------------------------------------------------------
// 3. SPREAD OPERATOR EN ARRAYS
// ----------------------------------------------------------------
std.print("3. SPREAD OPERATOR EN ARRAYS:");
std.print("  Permite expandir arrays en nuevos arrays");
std.print("");

// Combinar arrays
let arr1 = [1, 2, 3];
let arr2 = [4, 5, 6];
let combined = [...arr1, ...arr2];
std.print("arr1 =", arr1);
std.print("arr2 =", arr2);
std.print("combined = [...arr1, ...arr2] =", combined);
std.print("");

// Spread con elementos adicionales
let base = [10, 20];
let extended = [...base, 30, 40, 50];
std.print("base =", base);
std.print("extended = [...base, 30, 40, 50] =", extended);
std.print("");

// Spread al inicio
let prefix = [0, ...arr1];
std.print("prefix = [0, ...arr1] =", prefix);
std.print("");

// Spread en el medio
let middle = [1, ...arr2, 7, 8];
std.print("middle = [1, ...arr2, 7, 8] =", middle);
std.print("");

// ----------------------------------------------------------------
// 4. SPREAD OPERATOR EN OBJETOS
// ----------------------------------------------------------------
std.print("4. SPREAD OPERATOR EN OBJETOS:");
std.print("  Permite expandir objetos en nuevos objetos");
std.print("");

// Combinar objetos
let defaults = {
    theme: "light",
    fontSize: 14,
    autoSave: true
};

let userPrefs = {
    theme: "dark",
    language: "es"
};

let finalConfig = {...defaults, ...userPrefs};
std.print("defaults =", defaults);
std.print("userPrefs =", userPrefs);
std.print("finalConfig = {...defaults, ...userPrefs} =", finalConfig);
std.print("Nota: userPrefs.theme sobrescribe defaults.theme");
std.print("");

// Spread con propiedades adicionales
let person = {name: "Ana", age: 28};
let employee = {...person, id: 12345, department: "Engineering"};
std.print("person =", person);
std.print("employee = {...person, id: 12345, department: 'Engineering'} =", employee);
std.print("");

// ----------------------------------------------------------------
// 5. SPREAD OPERATOR EN LLAMADAS A FUNCIONES
// ----------------------------------------------------------------
std.print("5. SPREAD OPERATOR EN LLAMADAS A FUNCIONES:");
std.print("  Permite expandir arrays como argumentos individuales");
std.print("");

// Función que suma varios números
func sum(a, b, c, d) {
    return a + b + c + d;
}

let values = [5, 10, 15, 20];
let result = sum(...values);
std.print("sum(a, b, c, d) { return a + b + c + d; }");
std.print("values =", values);
std.print("sum(...values) =", result); // 50
std.print("");

// Spread con argumentos adicionales
let partialArgs = [2, 3];
let result2 = sum(1, ...partialArgs, 4);
std.print("partialArgs =", partialArgs);
std.print("sum(1, ...partialArgs, 4) =", result2); // 10
std.print("");

// Función con diferentes aridades
func greet(greeting, name, punctuation) {
    if (punctuation) {
        return greeting + " " + name + punctuation;
    } else if (name) {
        return greeting + " " + name;
    } else {
        return greeting;
    }
}

let greetingArgs1 = ["Hello"];
let greetingArgs2 = ["Hi", "Alice"];
let greetingArgs3 = ["Good morning", "Bob", "!"];

std.print("greet(...args) examples:");
std.print("greet(...['Hello']) =", greet(...greetingArgs1));
std.print("greet(...['Hi', 'Alice']) =", greet(...greetingArgs2));
std.print("greet(...['Good morning', 'Bob', '!']) =", greet(...greetingArgs3));
std.print("");

// ----------------------------------------------------------------
// 6. CASOS DE USO PRÁCTICOS
// ----------------------------------------------------------------
std.print("6. CASOS DE USO PRÁCTICOS:");
std.print("");

// Intercambio de variables usando destructuring
std.print("a) Intercambio de variables con destructuring:");
let x = 10;
let y = 20;
std.print("Antes: x =", x, ", y =", y);
// Intercambio usando destructuring en una nueva declaración
let [newX, newY] = [y, x];
std.print("Después del intercambio: newX =", newX, ", newY =", newY);
std.print("");

// Clonado de arrays con spread
std.print("b) Clonado de arrays:");
let original = [1, 2, 3];
let clone = [...original];
std.print("original =", original);
std.print("clone = [...original] =", clone);
std.print("Son arrays independientes");
std.print("");

// Configuración por defecto con spread
std.print("c) Configuración con valores por defecto:");
func createConfig(options) {
    let defaultConfig = {
        timeout: 5000,
        retries: 3,
        cache: true,
        debug: false
    };
    
    return {...defaultConfig, ...options};
}

let customConfig = createConfig({timeout: 10000, debug: true});
std.print("customConfig =", customConfig);
std.print("");

std.print("=== Fin del Example 45: Destructuring & Spread Operator ===");
std.print("¡Las características P2 han sido implementadas exitosamente!");