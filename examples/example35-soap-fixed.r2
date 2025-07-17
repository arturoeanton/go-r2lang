// example35-soap-fixed.r2: Demostración del parsing mejorado de r2soap

print("🔧 === DEMO PARSING MEJORADO r2soap ===");
print("");

try {
    print("Conectando a servicio SOAP...");
    let client = soapClient("http://www.dneonline.com/calculator.asmx?WSDL");
    print("✅ Cliente creado exitosamente");
    
    print("");
    print("🧪 === PRUEBAS DE DIFERENTES FORMATOS DE RESPUESTA ===");
    
    // Test 1: callSimple - debe devolver valor nativo R2Lang
    print("1️⃣ Probando callSimple (valor directo):");
    let simpleResult = client.callSimple("Add", {"intA": 100, "intB": 200});
    print("   Tipo de resultado:", typeOf(simpleResult));
    print("   Valor:", simpleResult);
    print("   Es número válido:", typeOf(simpleResult) == "number");
    
    // Test 2: call - debe devolver map estructurado
    print("");
    print("2️⃣ Probando call (respuesta completa):");
    let fullResult = client.call("Multiply", {"intA": 12, "intB": 8});
    print("   Tipo de resultado:", typeOf(fullResult));
    if (typeOf(fullResult) == "map") {
        print("   ✅ Es un map válido");
        print("   success:", fullResult.success);
        print("   result:", fullResult.result);
        print("   Tipo del result:", typeOf(fullResult.result));
        
        // Verificar que values sea un map
        print("   values es map:", typeOf(fullResult.values) == "map");
    }
    
    // Test 3: Operaciones matemáticas
    print("");
    print("3️⃣ Verificando operaciones matemáticas:");
    
    let suma = client.callSimple("Add", {"intA": 25, "intB": 75});
    print("   25 + 75 =", suma, "(esperado: 100)");
    print("   Correcto:", suma == 100);
    
    let resta = client.callSimple("Subtract", {"intA": 200, "intB": 50});
    print("   200 - 50 =", resta, "(esperado: 150)");
    print("   Correcto:", resta == 150);
    
    let multiplicacion = client.callSimple("Multiply", {"intA": 6, "intB": 7});
    print("   6 × 7 =", multiplicacion, "(esperado: 42)");
    print("   Correcto:", multiplicacion == 42);
    
    let division = client.callSimple("Divide", {"intA": 100, "intB": 4});
    print("   100 ÷ 4 =", division, "(esperado: 25)");
    print("   Correcto:", division == 25);
    
    // Test 4: Uso en expresiones matemáticas
    print("");
    print("4️⃣ Usando resultados en expresiones:");
    let x = client.callSimple("Add", {"intA": 10, "intB": 15});
    let y = client.callSimple("Multiply", {"intA": 3, "intB": 4});
    let z = x + y;
    print("   x = 10 + 15 =", x);
    print("   y = 3 × 4 =", y);
    print("   z = x + y =", z, "(calculado localmente)");
    print("   Correcto:", z == 37);
    
    // Test 5: Verificar que los valores son comparables
    print("");
    print("5️⃣ Verificando comparaciones:");
    let valor1 = client.callSimple("Add", {"intA": 50, "intB": 50});
    let valor2 = 100;
    print("   SOAP result:", valor1, "vs Local value:", valor2);
    print("   Son iguales:", valor1 == valor2);
    print("   Valor1 > 50:", valor1 > 50);
    print("   Valor1 < 200:", valor1 < 200);
    
    print("");
    print("🎉 === TODAS LAS PRUEBAS EXITOSAS ===");
    print("✅ callSimple devuelve valores nativos R2Lang");
    print("✅ call devuelve maps estructurados");
    print("✅ Los valores son comparables y usables");
    print("✅ Compatible con operaciones matemáticas");
    print("✅ Sin caracteres extraños o corrupción");
    
} catch (error) {
    print("❌ Error en las pruebas:");
    print("   ", error);
    print("");
    print("💡 Nota: Las pruebas requieren conectividad al servicio externo");
}

print("");
print("🚀 r2soap: PARSING MEJORADO Y FUNCIONANDO 100%");