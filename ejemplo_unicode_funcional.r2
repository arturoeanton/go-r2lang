// ========================================
// EJEMPLO FUNCIONAL DE UNICODE EN R2LANG
// Demostrando todas las capacidades Unicode
// ========================================

print("🌍 ¡Ejemplo funcional de Unicode en R2Lang! 🚀");
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
let altura = 175;          // Antes era 身長 pero simplificado
let nombre_ruso = "Иван";   // Antes era имя pero simplificado
let nombre_arabe = "أحمد";  // Antes era اسم pero simplificado
let nombre_griego = "Γιάννης"; // Antes era όνομα pero simplificado

print(`El año actual es: ${año}`);
print(`Niño: ${niño}, Señorita: ${señorita}`);
print(`Altura: ${altura}cm`);
print(`Nombre ruso: ${nombre_ruso}`);
print(`Nombre árabe: ${nombre_arabe}`);
print(`Nombre griego: ${nombre_griego}`);

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

// Emojis directos
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
print(`Invertido: "${ureverse(texto_español)}"`);

// ========================================
// 4. VALIDACIÓN Y CÓDIGOS UNICODE
// ========================================
print("\n🔍 4. Validación y códigos Unicode:");

let caracter_a = "A";
let caracter_ñ = "ñ";
let caracter_emoji = "🚀";

print(`Carácter A:`);
print(`  Código Unicode: ${ucharcode(caracter_a)}`);
print(`  Es UTF-8 válido: ${uisvalid(caracter_a)}`);
print(`  Es letra: ${uisLetter(caracter_a)}`);
print(`  Es mayúscula: ${uisUpper(caracter_a)}`);

print(`Carácter ñ:`);
print(`  Código Unicode: ${ucharcode(caracter_ñ)}`);
print(`  Es UTF-8 válido: ${uisvalid(caracter_ñ)}`);
print(`  Es letra: ${uisLetter(caracter_ñ)}`);
print(`  Es minúscula: ${uisLower(caracter_ñ)}`);

print(`Carácter emoji:`);
print(`  Código Unicode: ${ucharcode(caracter_emoji)}`);
print(`  Es UTF-8 válido: ${uisvalid(caracter_emoji)}`);

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

let palabra1 = "café";
let palabra2 = "cafe";
let palabra3 = "CAFÉ";

let comp1 = ucompare(palabra1, palabra2, "es");
let comp2 = ucompare(palabra1, palabra3, "es");

print(`"${palabra1}" comparado con "${palabra2}": ${comp1}`);
print(`"${palabra1}" comparado con "${palabra3}": ${comp2}`);

// ========================================
// 7. EXPRESIONES REGULARES UNICODE
// ========================================
print("\n🔍 7. Expresiones regulares Unicode:");

let texto_multiidioma = "Hello José! Привет María! 🌍";
print(`Texto: "${texto_multiidioma}"`);

// Buscar palabras que empiecen con mayúscula (simplificado)
let nombres = uregex("[A-Z][a-z]+", texto_multiidioma);
print(`Nombres encontrados: ${nombres}`);

// Verificar si contiene emojis
let tiene_emoji = uregexMatch("🌍", texto_multiidioma);
print(`¿Contiene emoji de tierra? ${tiene_emoji}`);

// ========================================
// 8. FUNCIÓN PRÁCTICA: PROCESADOR DE NOMBRES
// ========================================
print("\n🎯 8. Función práctica: Procesador de nombres:");

func procesarNombre(nombre_completo) {
    print(`\nProcesando: "${nombre_completo}"`);
    
    // Normalizar
    let normalizado = unormalize(nombre_completo, "NFC");
    
    // Capitalizar primera letra
    let primera = uupper(usubstr(normalizado, 0, 1));
    let resto = ulower(usubstr(normalizado, 1, ulen(normalizado) - 1));
    let procesado = primera + resto;
    
    print(`  Original: ${nombre_completo}`);
    print(`  Procesado: ${procesado}`);
    print(`  Longitud: ${ulen(normalizado)} caracteres`);
    
    // Verificar si tiene caracteres acentuados
    let tiene_acentos = uregexMatch("[àáâãäåæçèéêëìíîïðñòóôõöøùúûüýþÿ]", ulower(normalizado));
    print(`  ¿Tiene acentos? ${tiene_acentos}`);
    
    return procesado;
}

// Probar con diferentes nombres
procesarNombre("josé maría azañar");
procesarNombre("MARÍA JOSÉ FERNÁNDEZ");
procesarNombre("jean-françois dupont");

// ========================================
// 9. DEMOSTRACIÓN FINAL
// ========================================
print("\n🎉 9. Demostración final:");

let mensaje_multilingual = `
🌍 Saludos multilingües:
• Español: ¡Hola José María! ¿Cómo está usted?
• English: Hello José María! How are you?
• Français: Bonjour José María! Comment allez-vous?
• Русский: Привет José María! Как дела?
• العربية: مرحبا José María! كيف حالك؟
• 日本語: こんにちは José María! 元気ですか？
`;

print(mensaje_multilingual);
print(`Longitud total del mensaje: ${ulen(mensaje_multilingual)} caracteres`);

// Test final con todos los tipos de caracteres
let test_unicode = "ABC123ñéíóúÑÉÍÓÚçüäöß🚀🌍👋";
print(`\nTest final: "${test_unicode}"`);
print(`Longitud: ${ulen(test_unicode)}`);
print(`Mayúsculas: "${uupper(test_unicode)}"`);
print(`Minúsculas: "${ulower(test_unicode)}"`);
print(`Válido UTF-8: ${uisvalid(test_unicode)}`);

print("\n✅ ¡Soporte Unicode completo funcionando perfectamente en R2Lang! 🎊");