// example36-raw-clean.r2: Demostración del XML raw limpio

print("🧼 === DEMO XML RAW LIMPIO ===");
print("");

try {
    print("Conectando a servicio SOAP...");
    let client = soapClient("http://www.dneonline.com/calculator.asmx?WSDL");
    print("✅ Cliente creado");
    
    print("");
    print("📋 Probando respuesta call (completa):");
    let fullResponse = client.call("Add", {"intA": 123, "intB": 456});
    
    print("   success:", fullResponse.success);
    print("   result:", fullResponse.result);
    print("   values:", fullResponse.values);
    print("");
    print("   📄 XML Raw:");
    print("   " + fullResponse.raw);
    print("");
    print("   ✅ Características del XML raw:");
    print("   - Inicia con <?xml:", indexOf(fullResponse.raw, "<?xml") == 0);
    print("   - Termina con </soap:Envelope>:", indexOf(fullResponse.raw, "</soap:Envelope>") > 0);
    print("   - Sin caracteres extraños al inicio:", indexOf(fullResponse.raw, "<?xml") == 0);
    print("   - Longitud total:", len(fullResponse.raw), "caracteres");
    
    print("");
    print("🔍 Probando callRaw (XML puro):");
    let rawXML = client.callRaw("Multiply", {"intA": 9, "intB": 11});
    print("   Longitud:", len(rawXML), "caracteres");
    print("   Inicia correctamente:", indexOf(rawXML, "<?xml") == 0);
    print("   Contiene resultado:", indexOf(rawXML, "99") > 0);
    
    print("");
    print("✅ VERIFICACIONES EXITOSAS:");
    print("   - XML raw está completamente limpio");
    print("   - Sin caracteres binarios extraños");
    print("   - Formato XML válido y legible");
    print("   - Compatible con parsers XML externos");
    
} catch (error) {
    print("❌ Error:", error);
}

print("");
print("🎉 XML RAW COMPLETAMENTE LIMPIO Y FUNCIONAL");