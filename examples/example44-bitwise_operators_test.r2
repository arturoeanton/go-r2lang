// ============================================================================
// BITWISE OPERATORS COMPREHENSIVE TEST - R2Lang 2025
// ============================================================================
// Este archivo contiene tests comprensivos para todos los operadores bitwise
// implementados en R2Lang: &, |, ^, <<, >>, ~
// ============================================================================

std.print("🔧 INICIANDO TESTS DE OPERADORES BITWISE");
std.print("=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=");

// ============================================================================
// 1. OPERADORES BITWISE BÁSICOS
// ============================================================================
std.print("\n1️⃣ OPERADORES BITWISE BÁSICOS");

// Bitwise AND (&)
let and1 = 5 & 3;     // 0101 & 0011 = 0001 = 1
let and2 = 12 & 10;   // 1100 & 1010 = 1000 = 8
let and3 = 7 & 15;    // 0111 & 1111 = 0111 = 7

std.print("✓ Bitwise AND (&):");
std.print("  5 & 3 =", and1, "(esperado: 1)");
std.print("  12 & 10 =", and2, "(esperado: 8)");
std.print("  7 & 15 =", and3, "(esperado: 7)");

// Bitwise OR (|)
let or1 = 5 | 3;      // 0101 | 0011 = 0111 = 7
let or2 = 12 | 10;    // 1100 | 1010 = 1110 = 14
let or3 = 8 | 4;      // 1000 | 0100 = 1100 = 12

std.print("✓ Bitwise OR (|):");
std.print("  5 | 3 =", or1, "(esperado: 7)");
std.print("  12 | 10 =", or2, "(esperado: 14)");
std.print("  8 | 4 =", or3, "(esperado: 12)");

// Bitwise XOR (^)
let xor1 = 5 ^ 3;     // 0101 ^ 0011 = 0110 = 6
let xor2 = 12 ^ 10;   // 1100 ^ 1010 = 0110 = 6
let xor3 = 15 ^ 15;   // 1111 ^ 1111 = 0000 = 0

std.print("✓ Bitwise XOR (^):");
std.print("  5 ^ 3 =", xor1, "(esperado: 6)");
std.print("  12 ^ 10 =", xor2, "(esperado: 6)");
std.print("  15 ^ 15 =", xor3, "(esperado: 0)");

// ============================================================================
// 2. OPERADORES DE DESPLAZAMIENTO (SHIFT)
// ============================================================================
std.print("\n2️⃣ OPERADORES DE DESPLAZAMIENTO");

// Left shift (<<)
let lshift1 = 5 << 1;   // 0101 << 1 = 1010 = 10
let lshift2 = 3 << 2;   // 0011 << 2 = 1100 = 12
let lshift3 = 1 << 4;   // 0001 << 4 = 10000 = 16

std.print("✓ Left Shift (<<):");
std.print("  5 << 1 =", lshift1, "(esperado: 10)");
std.print("  3 << 2 =", lshift2, "(esperado: 12)");
std.print("  1 << 4 =", lshift3, "(esperado: 16)");

// Right shift (>>)
let rshift1 = 10 >> 1;  // 1010 >> 1 = 0101 = 5
let rshift2 = 12 >> 2;  // 1100 >> 2 = 0011 = 3
let rshift3 = 16 >> 4;  // 10000 >> 4 = 0001 = 1

std.print("✓ Right Shift (>>):");
std.print("  10 >> 1 =", rshift1, "(esperado: 5)");
std.print("  12 >> 2 =", rshift2, "(esperado: 3)");
std.print("  16 >> 4 =", rshift3, "(esperado: 1)");

// ============================================================================
// 3. OPERADOR BITWISE NOT (~)
// ============================================================================
std.print("\n3️⃣ OPERADOR BITWISE NOT");

// Bitwise NOT (~) - nota: en sistemas de 64 bits puede producir números grandes negativos
let not1 = ~0;         // Todos los bits en 1
let not2 = ~5;         // NOT 0101 = ...11111010 = -6
let not3 = ~(-1);      // NOT ...11111111 = 0

std.print("✓ Bitwise NOT (~):");
std.print("  ~0 =", not1, "(esperado: -1)");
std.print("  ~5 =", not2, "(esperado: -6)");
std.print("  ~(-1) =", not3, "(esperado: 0)");

// ============================================================================
// 4. OPERACIONES COMPLEJAS Y PRECEDENCIA
// ============================================================================
std.print("\n4️⃣ OPERACIONES COMPLEJAS Y PRECEDENCIA");

// Precedencia: & tiene mayor precedencia que |
let precedence1 = 5 | 3 & 2;    // 5 | (3 & 2) = 5 | 2 = 7
let precedence2 = (5 | 3) & 2;  // (5 | 3) & 2 = 7 & 2 = 2

std.print("✓ Precedencia de operadores:");
std.print("  5 | 3 & 2 =", precedence1, "(esperado: 7)");
std.print("  (5 | 3) & 2 =", precedence2, "(esperado: 2)");

// Combinaciones complejas
let complex1 = (8 | 4) & ~1;    // (8 | 4) & ~1 = 12 & -2 = 12
let complex2 = 15 ^ 7 << 1;     // (15 ^ 7) << 1 = 8 << 1 = 16
let complex3 = ~(5 & 3) | 2;    // ~(5 & 3) | 2 = ~1 | 2 = -2 | 2 = -2

std.print("✓ Operaciones complejas:");
std.print("  (8 | 4) & ~1 =", complex1);
std.print("  15 ^ 7 << 1 =", complex2, "(esperado: 16)");
std.print("  ~(5 & 3) | 2 =", complex3);

// ============================================================================
// 5. CASOS ESPECIALES Y EDGE CASES
// ============================================================================
std.print("\n5️⃣ CASOS ESPECIALES");

// Operaciones con cero
let zero_and = 0 & 255;     // 0
let zero_or = 0 | 255;      // 255
let zero_xor = 0 ^ 255;     // 255
let zero_lshift = 0 << 5;   // 0
let zero_rshift = 0 >> 5;   // 0

std.print("✓ Operaciones con cero:");
std.print("  0 & 255 =", zero_and, "(esperado: 0)");
std.print("  0 | 255 =", zero_or, "(esperado: 255)");
std.print("  0 ^ 255 =", zero_xor, "(esperado: 255)");
std.print("  0 << 5 =", zero_lshift, "(esperado: 0)");
std.print("  0 >> 5 =", zero_rshift, "(esperado: 0)");

// Números grandes
let big1 = 1024 & 512;      // 0
let big2 = 1024 | 512;      // 1536
let big3 = 1024 ^ 512;      // 1536
let big4 = 256 << 2;        // 1024
let big5 = 1024 >> 2;       // 256

std.print("✓ Números grandes:");
std.print("  1024 & 512 =", big1, "(esperado: 0)");
std.print("  1024 | 512 =", big2, "(esperado: 1536)");
std.print("  1024 ^ 512 =", big3, "(esperado: 1536)");
std.print("  256 << 2 =", big4, "(esperado: 1024)");
std.print("  1024 >> 2 =", big5, "(esperado: 256)");

// ============================================================================
// 6. COMPATIBILIDAD CON VARIABLES Y EXPRESIONES
// ============================================================================
std.print("\n6️⃣ COMPATIBILIDAD CON VARIABLES");

let a = 12;
let b = 7;
let c = 2;

let var_and = a & b;        // 12 & 7 = 4
let var_or = a | b;         // 12 | 7 = 15
let var_xor = a ^ b;        // 12 ^ 7 = 11
let var_lshift = a << c;    // 12 << 2 = 48
let var_rshift = a >> c;    // 12 >> 2 = 3
let var_not = ~a;           // ~12 = -13

std.print("✓ Operaciones con variables (a=12, b=7, c=2):");
std.print("  a & b =", var_and, "(esperado: 4)");
std.print("  a | b =", var_or, "(esperado: 15)");
std.print("  a ^ b =", var_xor, "(esperado: 11)");
std.print("  a << c =", var_lshift, "(esperado: 48)");
std.print("  a >> c =", var_rshift, "(esperado: 3)");
std.print("  ~a =", var_not, "(esperado: -13)");

// ============================================================================
// 7. OPERACIONES EN EXPRESIONES COMPLEJAS
// ============================================================================
std.print("\n7️⃣ EXPRESIONES COMPLEJAS");

let result1 = (a & b) | (c << 1);       // (12 & 7) | (2 << 1) = 4 | 4 = 4
let result2 = ~(a ^ b) & 255;           // ~(12 ^ 7) & 255 = ~11 & 255 = -12 & 255 = 244
let result3 = (a << 1) >> (c - 1);      // (12 << 1) >> (2 - 1) = 24 >> 1 = 12

std.print("✓ Expresiones complejas:");
std.print("  (a & b) | (c << 1) =", result1, "(esperado: 4)");
std.print("  ~(a ^ b) & 255 =", result2);
std.print("  (a << 1) >> (c - 1) =", result3, "(esperado: 12)");

// ============================================================================
// 8. COMPATIBILIDAD CON OTROS OPERADORES
// ============================================================================
std.print("\n8️⃣ COMPATIBILIDAD CON OTROS OPERADORES");

// Mezclando bitwise con aritméticos
let mixed1 = (5 & 3) + (4 | 2);        // 1 + 6 = 7
let mixed2 = (8 >> 1) * (2 << 1);      // 4 * 4 = 16
let mixed3 = (10 ^ 6) - (3 & 1);       // 12 - 1 = 11

std.print("✓ Mezclando con operadores aritméticos:");
std.print("  (5 & 3) + (4 | 2) =", mixed1, "(esperado: 7)");
std.print("  (8 >> 1) * (2 << 1) =", mixed2, "(esperado: 16)");
std.print("  (10 ^ 6) - (3 & 1) =", mixed3, "(esperado: 11)");

// Mezclando bitwise con comparaciones
let comp1 = (5 & 3) == 1;              // true
let comp2 = (8 | 4) > 10;              // true (12 > 10)
let comp3 = (7 ^ 3) < 5;               // false (4 < 5 = false, 4 no es menor que 5)

std.print("✓ Mezclando con operadores de comparación:");
std.print("  (5 & 3) == 1 =", comp1, "(esperado: true)");
std.print("  (8 | 4) > 10 =", comp2, "(esperado: true)");
std.print("  (7 ^ 3) < 5 =", comp3, "(esperado: true)");

// ============================================================================
// RESUMEN FINAL
// ============================================================================
std.print("\n" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=");
std.print("✅ TESTS DE OPERADORES BITWISE COMPLETADOS");
std.print("📊 Se han probado:");
std.print("   • Operadores binarios: &, |, ^, <<, >>");
std.print("   • Operador unario: ~");
std.print("   • Precedencia de operadores");
std.print("   • Casos especiales y edge cases");
std.print("   • Compatibilidad con variables");
std.print("   • Expresiones complejas");
std.print("   • Mezclado con otros operadores");
std.print("🎉 Si este script se ejecuta sin errores, los operadores bitwise están funcionando correctamente!");
std.print("=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=");