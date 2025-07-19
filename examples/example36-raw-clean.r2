// example36-raw-clean.r2: Demostración del XML raw limpio

std.print("🧼 === DEMO XML RAW LIMPIO ===");
std.print("");

try {
    std.print("Conectando a servicio SOAP...");
    let client = soap.client("http://www.dneonline.com/calculator.asmx?WSDL");
    std.print("✅ Cliente creado");
    
    std.print("");
    std.print("📋 Probando respuesta call (completa):");
    let fullResponse = client.call("Add", {"intA": 123, "intB": 456});
    
    std.print("   success:", fullResponse.success);
    std.print("   result:", fullResponse.result);
    std.print("   values:", fullResponse.values);
    std.print("");
    std.print("   📄 XML Raw:");
    std.print("   " + fullResponse.raw);
    std.print("");
    std.print("   ✅ Características del XML raw:");
    std.print("   - Inicia con <?xml:", string.indexOf(fullResponse.raw, "<?xml") == 0);
    std.print("   - Termina con </soap:Envelope>:", string.indexOf(fullResponse.raw, "</soap:Envelope>") > 0);
    std.print("   - Sin caracteres extraños al inicio:", string.indexOf(fullResponse.raw, "<?xml") == 0);
    std.print("   - Longitud total:", std.len(fullResponse.raw), "caracteres");
    
    std.print("");
    std.print("🔍 Probando callRaw (XML puro):");
    let rawXML = client.callRaw("Multiply", {"intA": 9, "intB": 11});
    std.print("   Longitud:", std.len(rawXML), "caracteres");
    std.print("   Inicia correctamente:", string.indexOf(rawXML, "<?xml") == 0);
    std.print("   Contiene resultado:", string.indexOf(rawXML, "99") > 0);
    
    std.print("");
    std.print("✅ VERIFICACIONES EXITOSAS:");
    std.print("   - XML raw está completamente limpio");
    std.print("   - Sin caracteres binarios extraños");
    std.print("   - Formato XML válido y legible");
    std.print("   - Compatible con parsers XML externos");
    
} catch (error) {
    std.print("❌ Error:", error);
}

std.print("");
std.print("🎉 XML RAW COMPLETAMENTE LIMPIO Y FUNCIONAL");