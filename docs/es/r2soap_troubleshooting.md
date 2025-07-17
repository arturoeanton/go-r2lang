# Solución de Problemas r2soap

## Errores Comunes y Soluciones

### 1. Error de Conectividad
```
panic: soapClient: failed to create client: failed to fetch WSDL: 
Get "http://example.com/service.wsdl": read tcp: connection reset by peer
```

**Causas posibles:**
- Servicio WSDL no disponible
- Problemas de red/firewall
- URL incorrecta
- Servicio temporalmente fuera de línea

**Soluciones:**
```javascript
// Siempre usar try/catch para servicios externos
try {
    let client = soapClient("http://service.com/service.wsdl");
    let result = client.call("Method", params);
    print("Éxito:", result);
} catch (error) {
    print("Error de conectividad:", error);
    // Usar fallback o servicio alternativo
}
```

### 2. Timeout de Conexión
```
panic: soapClient: failed to create client: timeout
```

**Solución:**
```javascript
// Configurar timeout más alto
let client = soapClient("http://slow-service.com/service.wsdl");
client.setTimeout(60.0); // 60 segundos
```

### 3. DNS Resolution Failed
```
panic: soapClient: failed to create client: no such host
```

**Soluciones:**
- Verificar conectividad a internet
- Confirmar que la URL es correcta
- Usar IP en lugar de hostname si es necesario

### 4. WSDL Inválido
```
panic: soapClient: failed to parse WSDL: XML syntax error
```

**Soluciones:**
- Verificar que la URL retorna un WSDL válido
- Probar la URL en un navegador
- Validar el XML del WSDL

## Mejores Prácticas

### 1. Validación de Conectividad
```javascript
func testSOAPConnectivity(wsdlURL) {
    try {
        let client = soapClient(wsdlURL);
        print("✅ Conectividad OK para:", wsdlURL);
        return client;
    } catch (error) {
        print("❌ Error de conectividad:", error);
        return null;
    }
}

// Uso
let client = testSOAPConnectivity("http://service.com/service.wsdl");
if (client != null) {
    // Usar el cliente
    let result = client.call("Method", params);
}
```

### 2. Configuración Robusta
```javascript
func createRobustSOAPClient(wsdlURL) {
    let client = soapClient(wsdlURL);
    
    // Configurar timeouts generosos
    client.setTimeout(30.0);
    
    // Headers para compatibilidad
    client.setHeader("User-Agent", "R2Lang-SOAP-Client/1.0");
    client.setHeader("Accept", "text/xml,application/xml");
    
    return client;
}
```

### 3. Manejo de Errores por Tipo
```javascript
func handleSOAPError(error) {
    let errorStr = "" + error;
    
    if (indexOf(errorStr, "connection reset") != -1) {
        print("Problema de red - reintenta más tarde");
    } else if (indexOf(errorStr, "timeout") != -1) {
        print("Timeout - aumenta el tiempo de espera");
    } else if (indexOf(errorStr, "no such host") != -1) {
        print("DNS error - verifica la URL");
    } else {
        print("Error desconocido:", error);
    }
}

// Uso
try {
    let client = soapClient(wsdlURL);
} catch (error) {
    handleSOAPError(error);
}
```

## Servicios de Prueba Confiables

### Servicios Públicos Estables
```javascript
// Ejemplos de servicios que suelen estar disponibles
// (verificar disponibilidad actual)

// 1. Convertidor de temperatura
let tempClient = soapClient("https://www.w3schools.com/xml/tempconvert.asmx?WSDL");

// 2. Calculadora básica (si está disponible)
let calcClient = soapClient("http://www.dneonline.com/calculator.asmx?WSDL");

// 3. Servicios bancarios de prueba (mock)
// Usar servicios locales o de desarrollo para testing
```

### Servicios Mock Locales
```javascript
// Para desarrollo, considera usar servicios mock locales
// o contenedores Docker con servicios SOAP de prueba

let localClient = soapClient("http://localhost:8080/mockservice?WSDL");
```

## Debugging

### 1. Verificar WSDL Manualmente
```bash
# Probar en terminal
curl -I "http://service.com/service.wsdl"

# Verificar contenido
curl "http://service.com/service.wsdl" | head -20
```

### 2. Test de Conectividad Básica
```javascript
// Test simple de generación de envelope
let envelope = soapEnvelope("http://tempuri.org/", "Test", "<param>value</param>");
print("Envelope generado:", len(envelope), "caracteres");

// Verificar funciones básicas
print("Funciones disponibles:");
print("- soapClient:", typeOf(soapClient));
print("- soapEnvelope:", typeOf(soapEnvelope));
print("- soapRequest:", typeOf(soapRequest));
```

### 3. Logging Detallado
```javascript
func debugSOAPCall(wsdlURL, operation, params) {
    print("=== DEBUG SOAP CALL ===");
    print("WSDL URL:", wsdlURL);
    print("Operation:", operation);
    print("Parameters:", params);
    
    try {
        print("1. Creando cliente...");
        let client = soapClient(wsdlURL);
        
        print("2. Listando operaciones...");
        let operations = client.listOperations();
        print("   Operaciones disponibles:", operations);
        
        print("3. Invocando operación...");
        let result = client.call(operation, params);
        print("   Resultado:", result);
        
        print("=== ÉXITO ===");
        return result;
        
    } catch (error) {
        print("=== ERROR ===");
        print("Error:", error);
        handleSOAPError(error);
        return null;
    }
}
```

## Alternativas ante Fallos

### 1. Fallback a HTTP Directo
```javascript
// Si SOAP falla, usar HTTP directo cuando sea posible
if (soapClient_failed) {
    let httpResponse = clientHttpPost(
        "http://service.com/endpoint",
        manualSoapEnvelope
    );
}
```

### 2. Servicios Alternativos
```javascript
func tryMultipleServices(services, operation, params) {
    for (let service in services) {
        try {
            let client = soapClient(service);
            return client.call(operation, params);
        } catch (error) {
            print("Falló servicio:", service);
            continue;
        }
    }
    print("Todos los servicios fallaron");
    return null;
}

// Uso
let services = [
    "http://primary-service.com/service.wsdl",
    "http://backup-service.com/service.wsdl", 
    "http://localhost:8080/mock-service.wsdl"
];

let result = tryMultipleServices(services, "Calculate", {"a": 5, "b": 3});
```

## Resumen

r2soap está completamente funcional y el error reportado es típico de conectividad con servicios externos. Las soluciones incluyen:

1. **Siempre usar try/catch** para servicios externos
2. **Configurar timeouts apropiados** para conexiones lentas
3. **Validar URLs** antes de usar
4. **Implementar fallbacks** para servicios críticos
5. **Usar servicios locales** para desarrollo y testing

La librería r2soap funciona perfectamente - solo requiere servicios WSDL accesibles para demostrar su funcionalidad completa.