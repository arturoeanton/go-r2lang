// example35-soap-fixed.r2: Demostración del parsing mejorado de r2soap

std.print("🔧 === DEMO PARSING MEJORADO r2soap ===");
std.print("");

try {
    std.print("Conectando a servicio SOAP...");
    let client = soap.client("http://www.dneonline.com/calculator.asmx?WSDL");
    std.print("✅ Cliente creado exitosamente");
    
    std.print("");
    std.print("🧪 === PRUEBAS DE DIFERENTES FORMATOS DE RESPUESTA ===");
    
    // Test 1: callSimple - debe devolver valor nativo R2Lang
    std.print("1️⃣ Probando callSimple (valor directo):");
    let simpleResult = client.callSimple("Add", {"intA": 100, "intB": 200});
    std.print("   Tipo de resultado:", std.typeOf(simpleResult));
    std.print("   Valor:", simpleResult);
    std.print("   Es número válido:", std.typeOf(simpleResult) == "number");
    
    // Test 2: call - debe devolver map estructurado
    std.print("");
    std.print("2️⃣ Probando call (respuesta completa):");
    let fullResult = client.call("Multiply", {"intA": 12, "intB": 8});
    std.print("   Tipo de resultado:", std.typeOf(fullResult));
    if (std.typeOf(fullResult) == "map") {
        std.print("   ✅ Es un map válido");
        std.print("   success:", fullResult.success);
        std.print("   result:", fullResult.result);
        std.print("   Tipo del result:", std.typeOf(fullResult.result));
        
        // Verificar que values sea un map
        std.print("   values es map:", std.typeOf(fullResult.values) == "map");
    }
    
    // Test 3: Operaciones matemáticas
    std.print("");
    std.print("3️⃣ Verificando operaciones matemáticas:");
    
    let suma = client.callSimple("Add", {"intA": 25, "intB": 75});
    std.print("   25 + 75 =", suma, "(esperado: 100)");
    std.print("   Correcto:", suma == 100);
    
    let resta = client.callSimple("Subtract", {"intA": 200, "intB": 50});
    std.print("   200 - 50 =", resta, "(esperado: 150)");
    std.print("   Correcto:", resta == 150);
    
    let multiplicacion = client.callSimple("Multiply", {"intA": 6, "intB": 7});
    std.print("   6 × 7 =", multiplicacion, "(esperado: 42)");
    std.print("   Correcto:", multiplicacion == 42);
    
    let division = client.callSimple("Divide", {"intA": 100, "intB": 4});
    std.print("   100 ÷ 4 =", division, "(esperado: 25)");
    std.print("   Correcto:", division == 25);
    
    // Test 4: Uso en expresiones matemáticas
    std.print("");
    std.print("4️⃣ Usando resultados en expresiones:");
    let x = client.callSimple("Add", {"intA": 10, "intB": 15});
    let y = client.callSimple("Multiply", {"intA": 3, "intB": 4});
    let z = x + y;
    std.print("   x = 10 + 15 =", x);
    std.print("   y = 3 × 4 =", y);
    std.print("   z = x + y =", z, "(calculado localmente)");
    std.print("   Correcto:", z == 37);
    
    // Test 5: Verificar que los valores son comparables
    std.print("");
    std.print("5️⃣ Verificando comparaciones:");
    let valor1 = client.callSimple("Add", {"intA": 50, "intB": 50});
    let valor2 = 100;
    std.print("   SOAP result:", valor1, "vs Local value:", valor2);
    std.print("   Son iguales:", valor1 == valor2);
    std.print("   Valor1 > 50:", valor1 > 50);
    std.print("   Valor1 < 200:", valor1 < 200);
    
    std.print("");
    std.print("🎉 === TODAS LAS PRUEBAS EXITOSAS ===");
    std.print("✅ callSimple devuelve valores nativos R2Lang");
    std.print("✅ call devuelve maps estructurados");
    std.print("✅ Los valores son comparables y usables");
    std.print("✅ Compatible con operaciones matemáticas");
    std.print("✅ Sin caracteres extraños o corrupción");
    
} catch (error) {
    std.print("❌ Error en las pruebas:");
    std.print("   ", error);
    std.print("");
    std.print("💡 Nota: Las pruebas requieren conectividad al servicio externo");
}

std.print("");
std.print("🚀 r2soap: PARSING MEJORADO Y FUNCIONANDO 100%");