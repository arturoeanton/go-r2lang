// ========================================
// EJEMPLO COMPLETO DE UNICODE EN R2LANG
// Demostrando todas las capacidades Unicode
// ========================================

print("🌍 ¡Bienvenidos al ejemplo de Unicode en R2Lang! 🚀");
print("==================================================");

// ========================================
// 1. IDENTIFICADORES UNICODE
// ========================================
print("\n📝 1. Identificadores Unicode:");

// Español
let año = 2024;
let niño = "Antonio";
let señorita = "María José";

// Otros idiomas
let 身長 = 175;          // Japonés: altura
let имя = "Иван";        // Ruso: nombre
let اسم = "أحمد";        // Árabe: nombre
let όνομα = "Γιάννης";   // Griego: nombre
let prénoms = "Jean-François"; // Francés

print(`El año actual es: ${año}`);
print(`Niño: ${niño}, Señorita: ${señorita}`);
print(`身長 (altura): ${身長}cm`);
print(`Имя (nombre en ruso): ${имя}`);
print(`اسم (nombre en árabe): ${اسم}`);
print(`Όνομα (nombre en griego): ${όνομα}`);
print(`Prénoms français: ${prénoms}`);

// ========================================
// 2. STRINGS UNICODE Y ESCAPE SEQUENCES
// ========================================
print("\n🔤 2. Strings Unicode y Escape Sequences:");

let emoji_wave = "\U0001F44B";           // 👋
let emoji_rocket = "\U0001F680";         // 🚀
let emoji_earth = "\U0001F30D";          // 🌍
let spanish_chars = "\u00f1\u00e9\u00fa"; // ñéú

print(`Saludando: ${emoji_wave}`);
print(`Cohete: ${emoji_rocket}`);
print(`Tierra: ${emoji_earth}`);
print(`Caracteres españoles: ${spanish_chars}`);

// Emojis complejos (secuencias)
let familia = "👨‍👩‍👧‍👦";
let bandera_españa = "🇪🇸";
print(`Familia: ${familia}`);
print(`Bandera de España: ${bandera_españa}`);

// ========================================
// 3. FUNCIONES BÁSICAS DE UNICODE
// ========================================
print("\n⚙️ 3. Funciones básicas de Unicode:");

let texto_español = "José María Azañar 🇪🇸";
print(`Texto: "${texto_español}"`);
print(`Longitud (ulen): ${ulen(texto_español)} caracteres`);
print(`Longitud normal (len): ${len(texto_español)} bytes`);

// Substring Unicode
let primer_nombre = usubstr(texto_español, 0, 4);
let segundo_nombre = usubstr(texto_español, 5, 5);
print(`Primer nombre: "${primer_nombre}"`);
print(`Segundo nombre: "${segundo_nombre}"`);

// Transformaciones
print(`Mayúsculas: "${uupper(texto_español)}"`);
print(`Minúsculas: "${ulower(texto_español)}"`);
print(`Título: "${utitle("hola mundo")}"`);
print(`Invertido: "${ureverse(texto_español)}"`);

// ========================================
// 4. VALIDACIÓN Y CÓDIGOS UNICODE
// ========================================
print("\n🔍 4. Validación y códigos Unicode:");

let caracteres = ["A", "ñ", "🚀", "José"];
for (let char in caracteres) {
    print(`Carácter: "${char}"`);
    print(`  Código Unicode: ${ucharcode(char)}`);
    print(`  Es UTF-8 válido: ${uisvalid(char)}`);
    print(`  Es letra: ${uisLetter(char)}`);
    print(`  Es mayúscula: ${uisUpper(char)}`);
    print(`  Es minúscula: ${uisLower(char)}`);
    print("");
}

// Crear caracteres desde códigos
let a_from_code = ufromcode(65);    // A
let ñ_from_code = ufromcode(241);   // ñ
let emoji_from_code = ufromcode(128075); // 👋
print(`Desde códigos: ${a_from_code}, ${ñ_from_code}, ${emoji_from_code}`);

// ========================================
// 5. NORMALIZACIÓN UNICODE
// ========================================
print("\n🔄 5. Normalización Unicode:");

let cafe_composed = "café";    // é como un solo carácter
let cafe_decomposed = "cafe\u0301"; // e + combining acute accent

print(`Café compuesto: "${cafe_composed}"`);
print(`Café descompuesto: "${cafe_decomposed}"`);
print(`Longitudes: ${ulen(cafe_composed)} vs ${ulen(cafe_decomposed)}`);

// Normalizar
let cafe_nfc = unormalize(cafe_decomposed, "NFC");
let cafe_nfd = unormalize(cafe_composed, "NFD");
print(`NFC normalizado: "${cafe_nfc}" (${ulen(cafe_nfc)})`);
print(`NFD normalizado: "${cafe_nfd}" (${ulen(cafe_nfd)})`);

// ========================================
// 6. COMPARACIÓN UNICODE
// ========================================
print("\n⚖️ 6. Comparación Unicode:");

let palabras = ["café", "cafe", "CAFÉ", "Café"];
for (let i = 0; i < len(palabras); i++) {
    for (let j = i + 1; j < len(palabras); j++) {
        let palabra1 = palabras[i];
        let palabra2 = palabras[j];
        let comparacion = ucompare(palabra1, palabra2, "es");
        let resultado = "";
        if (comparacion < 0) {
            resultado = "menor que";
        } else {
            if (comparacion > 0) {
                resultado = "mayor que";
            } else {
                resultado = "igual a";
            }
        }
        print(`"${palabra1}" es ${resultado} "${palabra2}"`);
    }
}

// ========================================
// 7. EXPRESIONES REGULARES UNICODE
// ========================================
print("\n🔍 7. Expresiones regulares Unicode:");

let texto_multiidioma = "Hello José! Привет María! こんにちは Ahmed! 🌍";
print(`Texto: "${texto_multiidioma}"`);

// Buscar palabras que empiecen con mayúscula
let nombres = uregex("[A-ZÀ-ÿА-я一-龯][a-zà-ÿа-я]*", texto_multiidioma);
print(`Nombres encontrados: ${nombres}`);

// Verificar si contiene emojis
let tiene_emoji = uregexMatch("🌍", texto_multiidioma);
print(`¿Contiene emojis? ${tiene_emoji}`);

// ========================================
// 8. APLICACIÓN PRÁCTICA: PROCESADOR DE NOMBRES
// ========================================
print("\n🎯 8. Aplicación práctica: Procesador de nombres:");

func procesarNombre(nombre_completo) {
    print(`\nProcesando: "${nombre_completo}"`);
    
    // Normalizar
    let normalizado = unormalize(nombre_completo, "NFC");
    
    // Separar por espacios
    let partes = uregex("[^\\s]+", normalizado);
    
    // Procesar cada parte
    let nombres_procesados = [];
    for (let parte in partes) {
        // Capitalizar primera letra
        let primera = uupper(usubstr(parte, 0, 1));
        let resto = ulower(usubstr(parte, 1, ulen(parte) - 1));
        let procesado = primera + resto;
        nombres_procesados = nombres_procesados + [procesado];
    }
    
    print(`  Partes: ${partes}`);
    print(`  Procesado: ${nombres_procesados}`);
    print(`  Longitud total: ${ulen(normalizado)} caracteres`);
    
    // Verificar caracteres especiales
    let tiene_acentos = uregexMatch("[àáâãäåæçèéêëìíîïðñòóôõöøùúûüýþÿ]", ulower(normalizado));
    print(`  ¿Tiene acentos? ${tiene_acentos}`);
    
    return nombres_procesados;
}

// Probar con diferentes nombres
let nombres_test = ["josé maría azañar", "MARÍA JOSÉ FERNÁNDEZ", "jean-françois dupont", "محمد احمد علي", "山田太郎", "владимир путин"];

for (let nombre in nombres_test) {
    procesarNombre(nombre);
}

// ========================================
// 9. ESTADÍSTICAS FINALES
// ========================================
print("\n📊 9. Estadísticas finales del ejemplo:");

let codigo_ejemplo = `
// Este código demuestra:
// ✅ Identificadores Unicode (español, japonés, ruso, árabe, griego)
// ✅ Strings con emojis y caracteres especiales  
// ✅ Escape sequences Unicode (\u{xxxx})
// ✅ Funciones básicas: ulen, usubstr, uupper, ulower, ureverse
// ✅ Normalización Unicode (NFC, NFD)
// ✅ Comparación sensible a idioma
// ✅ Expresiones regulares Unicode
// ✅ Validación de caracteres
// ✅ Conversión entre códigos y caracteres
// ✅ Aplicación práctica de procesamiento de texto
`;

print(codigo_ejemplo);
print(`\n🎉 ¡Ejemplo completado exitosamente!`);
print(`Total de caracteres en este archivo: ${ulen(codigo_ejemplo)} caracteres`);
print(`¡R2Lang ahora soporta Unicode completamente! 🌍🚀`);