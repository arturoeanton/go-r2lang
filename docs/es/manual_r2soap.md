# Manual Completo de r2soap - Cliente SOAP Dinámico para R2Lang

## Tabla de Contenidos

1. [Introducción](#introducción)
2. [Instalación y Configuración](#instalación-y-configuración)
3. [Conceptos Fundamentales](#conceptos-fundamentales)
4. [API Completa](#api-completa)
5. [Ejemplos Básicos](#ejemplos-básicos)
6. [Configuración Empresarial](#configuración-empresarial)
7. [Autenticación](#autenticación)
8. [SSL/TLS y Seguridad](#ssltls-y-seguridad)
9. [Manejo de Errores](#manejo-de-errores)
10. [Mejores Prácticas](#mejores-prácticas)
11. [Casos de Uso Avanzados](#casos-de-uso-avanzados)
12. [Troubleshooting](#troubleshooting)
13. [Referencia Técnica](#referencia-técnica)

---

## Introducción

r2soap es un cliente SOAP completamente dinámico para R2Lang que permite interactuar con servicios web SOAP sin necesidad de generar código. Utiliza parsing dinámico de WSDL para descubrir operaciones y crear requests automáticamente.

### Características Principales

- **100% Dinámico**: No requiere generación de código
- **Parsing WSDL Automático**: Descubre operaciones y parámetros automáticamente
- **Parsing de Respuestas Inteligente**: Convierte XML a objetos R2Lang
- **Headers Customizables**: Soporte completo para headers empresariales
- **Autenticación Empresarial**: Basic Auth, Bearer tokens, y más
- **SSL/TLS Completo**: Soporte para servicios seguros
- **Manejo de Errores Robusto**: Detección automática de fallas SOAP
- **Múltiples Formatos de Respuesta**: Full, simple, y raw

### Ventajas sobre Clientes Tradicionales

```javascript
// Tradicional (código generado)
import { CalculatorClient } from './generated/calculator-client';
let client = new CalculatorClient();

// r2soap (100% dinámico)
let client = soapClient("http://service.com/calculator.wsdl");
```

---

## Instalación y Configuración

### Prerequisitos

- R2Lang instalado y configurado
- Conectividad de red para acceder a servicios WSDL
- (Opcional) Certificados SSL para servicios seguros

### Verificación de Instalación

```javascript
// Verificar que r2soap está disponible
print("Funciones SOAP disponibles:");
print("- soapClient:", typeOf(soapClient));
print("- soapEnvelope:", typeOf(soapEnvelope));
print("- soapRequest:", typeOf(soapRequest));
```

---

## Conceptos Fundamentales

### WSDL (Web Services Description Language)

WSDL es un documento XML que describe:
- **Operaciones disponibles**: Métodos que se pueden llamar
- **Parámetros de entrada**: Tipos y nombres de parámetros
- **Respuestas**: Estructura de las respuestas
- **Endpoints**: URLs donde está el servicio
- **Protocolo**: Configuración SOAP

### Cliente SOAP Dinámico

```javascript
// Creación básica
let client = soapClient("http://service.com/service.wsdl");

// Con headers personalizados
let customHeaders = {
    "User-Agent": "MiAplicacion/1.0",
    "X-Company": "MiEmpresa"
};
let client = soapClient("http://service.com/service.wsdl", customHeaders);
```

### Tipos de Respuesta

r2soap ofrece tres tipos de respuesta:

1. **Full Response** (`call`): Objeto completo con metadatos
2. **Simple Response** (`callSimple`): Solo el valor resultado
3. **Raw Response** (`callRaw`): XML crudo sin procesar

---

## API Completa

### Función Principal: soapClient

```javascript
soapClient(wsdlURL, [customHeaders])
```

**Parámetros:**
- `wsdlURL` (string): URL del documento WSDL
- `customHeaders` (object, opcional): Headers HTTP personalizados

**Retorna:** Objeto cliente SOAP con métodos

### Métodos del Cliente

#### Operaciones SOAP

```javascript
// Llamada completa con metadatos
client.call(operationName, parameters)

// Llamada simple (solo resultado)
client.callSimple(operationName, parameters)

// Respuesta XML cruda
client.callRaw(operationName, parameters)
```

#### Descubrimiento de Servicios

```javascript
// Listar operaciones disponibles
let operations = client.listOperations();

// Obtener información de operación específica
let opInfo = client.getOperation("operationName");
```

#### Configuración de Headers

```javascript
// Establecer header individual
client.setHeader("HeaderName", "HeaderValue");

// Establecer múltiples headers
client.setHeader({
    "Header1": "Value1",
    "Header2": "Value2"
});

// Obtener headers actuales
let headers = client.getHeaders();

// Eliminar header específico
client.removeHeader("HeaderName");

// Resetear a defaults del navegador
client.resetHeaders();
```

#### Configuración de Timeouts

```javascript
// Establecer timeout en segundos
client.setTimeout(60.0); // 60 segundos
```

#### Configuración SSL/TLS

```javascript
// Configuración TLS
client.setTLSConfig({
    "minVersion": "1.2",      // TLS 1.2 mínimo
    "skipVerify": false       // Verificar certificados
});

// Para testing (NO producción)
client.setTLSConfig({
    "skipVerify": true        // ⚠️ Solo para testing
});
```

#### Autenticación

```javascript
// Basic Authentication
client.setAuth({
    "type": "basic",
    "username": "usuario",
    "password": "contraseña"
});

// Bearer Token
client.setAuth({
    "type": "bearer",
    "token": "eyJhbGciOiJIUzI1NiIs..."
});
```

#### Propiedades del Cliente

```javascript
// Información del servicio
print("WSDL URL:", client.wsdlURL);
print("Service URL:", client.serviceURL);
print("Namespace:", client.namespace);
```

---

## Ejemplos Básicos

### Ejemplo 1: Cliente Simple

```javascript
// Conectar a servicio de calculadora
let client = soapClient("http://www.dneonline.com/calculator.asmx?WSDL");

// Verificar conexión
print("✅ Conectado a:", client.serviceURL);

// Listar operaciones
let operations = client.listOperations();
print("Operaciones:", operations);

// Realizar operación simple
let result = client.callSimple("Add", {"intA": 10, "intB": 5});
print("10 + 5 =", result);
```

### Ejemplo 2: Respuesta Completa

```javascript
let client = soapClient("http://service.com/calculator.wsdl");

// Obtener respuesta completa
let response = client.call("Multiply", {"intA": 7, "intB": 8});

if (response.success) {
    print("✅ Operación exitosa");
    print("Resultado:", response.result);
    print("Valores extraídos:", response.values);
} else {
    print("❌ Error SOAP:", response.fault);
}
```

### Ejemplo 3: Manejo de Errores

```javascript
try {
    let client = soapClient("http://service.com/service.wsdl");
    let result = client.callSimple("Operation", {"param": "value"});
    print("Resultado:", result);
} catch (error) {
    print("Error:", error);
    
    let errorStr = "" + error;
    if (indexOf(errorStr, "timeout") != -1) {
        print("💡 Sugerencia: Aumentar timeout con client.setTimeout()");
    } else if (indexOf(errorStr, "connection") != -1) {
        print("💡 Sugerencia: Verificar conectividad de red");
    }
}
```

---

## Configuración Empresarial

### Headers Corporativos

```javascript
let client = soapClient("https://internal.company.com/service.wsdl");

// Headers de compliance empresarial
client.setHeader({
    "X-Company-ID": "CORP-12345",
    "X-Department": "Finance",
    "X-User-Role": "Manager",
    "X-Session-ID": generateSessionId(),
    "X-Compliance": "SOX-Approved",
    "X-Request-ID": generateRequestId()
});
```

### Configuración de Timeouts para Servicios Lentos

```javascript
// Para servicios empresariales que pueden ser lentos
client.setTimeout(120.0); // 2 minutos

// Headers adicionales para servicios internos
client.setHeader({
    "Keep-Alive": "timeout=300",
    "Connection": "keep-alive"
});
```

### Ambiente de Desarrollo vs Producción

```javascript
func createEnterpriseClient(wsdlURL, environment) {
    let client = soapClient(wsdlURL);
    
    if (environment == "development") {
        // Configuración para desarrollo
        client.setTimeout(30.0);
        client.setTLSConfig({"skipVerify": true});
        client.setHeader("X-Environment", "dev");
    } else if (environment == "production") {
        // Configuración para producción
        client.setTimeout(60.0);
        client.setTLSConfig({
            "minVersion": "1.2",
            "skipVerify": false
        });
        client.setHeader("X-Environment", "prod");
    }
    
    return client;
}
```

---

## Autenticación

### Basic Authentication

```javascript
let client = soapClient("https://secure-service.com/service.wsdl");

// Configurar Basic Auth
client.setAuth({
    "type": "basic",
    "username": "enterprise_user",
    "password": "secure_password_123"
});

// El header Authorization se agrega automáticamente
let result = client.callSimple("SecureOperation", {"data": "sensitive"});
```

### Bearer Token (OAuth/JWT)

```javascript
let client = soapClient("https://api.company.com/service.wsdl");

// Configurar Bearer token
client.setAuth({
    "type": "bearer",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
});

// Token se incluye automáticamente en requests
let result = client.call("AuthenticatedOperation", {"param": "value"});
```

### Autenticación Personalizada

```javascript
// Para sistemas de autenticación propietarios
let client = soapClient("https://custom-auth.com/service.wsdl");

// Agregar headers de autenticación personalizados
client.setHeader({
    "X-API-Key": "your-api-key",
    "X-Secret": "your-secret",
    "X-Timestamp": getCurrentTimestamp(),
    "X-Signature": generateSignature()
});
```

---

## SSL/TLS y Seguridad

### Configuración SSL Básica

```javascript
let client = soapClient("https://secure.company.com/service.wsdl");

// Configuración SSL estricta para producción
client.setTLSConfig({
    "minVersion": "1.2",      // TLS 1.2 mínimo
    "skipVerify": false       // Siempre verificar certificados
});
```

### Configuración para Certificados Auto-firmados

```javascript
// ⚠️ Solo para desarrollo/testing
let client = soapClient("https://internal-dev.company.com/service.wsdl");

client.setTLSConfig({
    "skipVerify": true,       // Skip verificación para testing
    "minVersion": "1.1"       // Permitir TLS más bajo
});

// Marcar claramente como inseguro
client.setHeader("X-Insecure-Mode", "true");
```

### Configuración Empresarial SSL

```javascript
func createSecureClient(wsdlURL) {
    let client = soapClient(wsdlURL);
    
    // Configuración corporativa estricta
    client.setTLSConfig({
        "minVersion": "1.3",      // Solo TLS 1.3
        "skipVerify": false
    });
    
    // Headers de seguridad adicionales
    client.setHeader({
        "Strict-Transport-Security": "max-age=31536000",
        "X-Content-Type-Options": "nosniff",
        "X-Frame-Options": "DENY"
    });
    
    return client;
}
```

---

## Manejo de Errores

### Tipos de Errores Comunes

#### 1. Errores de Conectividad

```javascript
try {
    let client = soapClient("http://unavailable-service.com/service.wsdl");
} catch (error) {
    let errorStr = "" + error;
    
    if (indexOf(errorStr, "connection reset") != -1) {
        print("🌐 Error: Servicio no disponible o bloqueando requests");
        print("💡 Solución: Verificar headers User-Agent o firewall");
    } else if (indexOf(errorStr, "timeout") != -1) {
        print("⏱️ Error: Timeout de conexión");
        print("💡 Solución: Aumentar timeout o verificar red");
    } else if (indexOf(errorStr, "no such host") != -1) {
        print("🔍 Error: No se pudo resolver DNS");
        print("💡 Solución: Verificar URL y conectividad");
    }
}
```

#### 2. Errores de Autenticación

```javascript
try {
    let result = client.call("SecureOperation", {"data": "test"});
} catch (error) {
    if (indexOf("" + error, "401") != -1) {
        print("🔐 Error de autenticación");
        print("💡 Verificar credenciales con client.setAuth()");
    } else if (indexOf("" + error, "403") != -1) {
        print("🚫 Sin permisos para esta operación");
        print("💡 Contactar administrador del servicio");
    }
}
```

#### 3. Errores SOAP (Faults)

```javascript
let response = client.call("RiskyOperation", {"data": "invalid"});

if (!response.success && response.fault) {
    print("🚨 SOAP Fault detectado:");
    print("Código:", response.fault.code);
    print("Mensaje:", response.fault.message);
    
    // Manejo específico por tipo de fault
    if (response.fault.code == "Client") {
        print("💡 Error del cliente - verificar parámetros");
    } else if (response.fault.code == "Server") {
        print("💡 Error del servidor - contactar soporte");
    }
}
```

### Estrategias de Reintentos

```javascript
func callWithRetry(client, operation, params, maxRetries) {
    for (let attempt = 1; attempt <= maxRetries; attempt++) {
        try {
            print("🔄 Intento", attempt, "de", maxRetries);
            let result = client.callSimple(operation, params);
            print("✅ Éxito en intento", attempt);
            return result;
        } catch (error) {
            print("❌ Fallo intento", attempt, ":", error);
            
            if (attempt < maxRetries) {
                print("⏳ Esperando antes del siguiente intento...");
                sleep(2.0 * attempt); // Backoff exponencial
            }
        }
    }
    print("💥 Falló después de", maxRetries, "intentos");
    return null;
}

// Uso
let result = callWithRetry(client, "UnstableOperation", {"data": "test"}, 3);
```

---

## Mejores Prácticas

### 1. Creación de Clientes Reutilizables

```javascript
func createConfiguredClient(wsdlURL, config) {
    let client = soapClient(wsdlURL, config.headers || {});
    
    // Aplicar configuraciones
    if (config.timeout) {
        client.setTimeout(config.timeout);
    }
    
    if (config.auth) {
        client.setAuth(config.auth);
    }
    
    if (config.tls) {
        client.setTLSConfig(config.tls);
    }
    
    // Headers empresariales estándar
    client.setHeader({
        "X-Client": "R2Lang-SOAP/1.0",
        "X-Request-Time": getCurrentTimestamp()
    });
    
    return client;
}

// Configuración para producción
let prodConfig = {
    "timeout": 60.0,
    "auth": {
        "type": "basic",
        "username": getEnvVar("SOAP_USER"),
        "password": getEnvVar("SOAP_PASS")
    },
    "tls": {
        "minVersion": "1.2",
        "skipVerify": false
    }
};

let client = createConfiguredClient("https://prod.company.com/service.wsdl", prodConfig);
```

### 2. Logging y Auditoría

```javascript
func auditedSOAPCall(client, operation, params, userId) {
    let requestId = "REQ-" + Math.floor(Math.random() * 1000000);
    
    // Agregar headers de auditoría
    client.setHeader({
        "X-Request-ID": requestId,
        "X-User-ID": userId,
        "X-Timestamp": getCurrentTimestamp()
    });
    
    print("📝 Audit Log - Request ID:", requestId);
    print("   Operation:", operation);
    print("   User:", userId);
    print("   Parameters:", Object.keys(params));
    
    try {
        let startTime = getCurrentTime();
        let result = client.call(operation, params);
        let duration = getCurrentTime() - startTime;
        
        print("✅ Audit Log - Success");
        print("   Duration:", duration, "ms");
        print("   Response Type:", typeOf(result));
        
        return result;
    } catch (error) {
        print("❌ Audit Log - Error:", error);
        throw error;
    }
}
```

### 3. Caché de Clientes

```javascript
let clientCache = {};

func getCachedClient(wsdlURL) {
    if (clientCache[wsdlURL]) {
        print("📋 Usando cliente cacheado para:", wsdlURL);
        return clientCache[wsdlURL];
    }
    
    print("🆕 Creando nuevo cliente para:", wsdlURL);
    let client = soapClient(wsdlURL);
    clientCache[wsdlURL] = client;
    return client;
}
```

### 4. Validación de Parámetros

```javascript
func validateAndCall(client, operation, params, schema) {
    // Validar que todos los parámetros requeridos estén presentes
    for (let required in schema.required) {
        if (!params[required]) {
            throw "Parámetro requerido faltante: " + required;
        }
    }
    
    // Validar tipos
    for (let param in params) {
        if (schema.types[param] && typeof(params[param]) != schema.types[param]) {
            throw "Tipo incorrecto para " + param + ". Esperado: " + schema.types[param];
        }
    }
    
    return client.callSimple(operation, params);
}

// Esquema de ejemplo
let addSchema = {
    "required": ["intA", "intB"],
    "types": {
        "intA": "number",
        "intB": "number"
    }
};

let result = validateAndCall(client, "Add", {"intA": 5, "intB": 3}, addSchema);
```

---

## Casos de Uso Avanzados

### 1. Integración con Sistema de Facturación

```javascript
func processInvoice(invoiceData) {
    // Cliente para servicio de facturación empresarial
    let billingClient = soapClient("https://billing.company.com/service.wsdl");
    
    // Configuración empresarial
    billingClient.setAuth({
        "type": "basic",
        "username": getConfig("billing.username"),
        "password": getConfig("billing.password")
    });
    
    billingClient.setHeader({
        "X-Department": "Finance",
        "X-System": "ERP-Integration",
        "X-Version": "2.1"
    });
    
    try {
        // Crear factura
        let invoice = billingClient.call("CreateInvoice", {
            "customerId": invoiceData.customerId,
            "amount": invoiceData.amount,
            "items": serializeItems(invoiceData.items)
        });
        
        if (invoice.success) {
            print("✅ Factura creada:", invoice.result.invoiceId);
            
            // Enviar por email
            let emailResult = billingClient.call("SendInvoiceEmail", {
                "invoiceId": invoice.result.invoiceId,
                "email": invoiceData.customerEmail
            });
            
            return {
                "success": true,
                "invoiceId": invoice.result.invoiceId,
                "emailSent": emailResult.success
            };
        } else {
            throw "Error al crear factura: " + invoice.fault.message;
        }
    } catch (error) {
        print("❌ Error en facturación:", error);
        return {"success": false, "error": error};
    }
}
```

### 2. Cliente Multi-Servicio

```javascript
class EnterpriseSOAPManager {
    func constructor() {
        this.clients = {};
        this.defaultConfig = {
            "timeout": 60.0,
            "tls": {"minVersion": "1.2", "skipVerify": false}
        };
    }
    
    func getClient(serviceName) {
        if (!this.clients[serviceName]) {
            let config = getServiceConfig(serviceName);
            let client = soapClient(config.wsdlURL);
            
            // Aplicar configuración predeterminada
            client.setTimeout(this.defaultConfig.timeout);
            client.setTLSConfig(this.defaultConfig.tls);
            
            // Configuración específica del servicio
            if (config.auth) {
                client.setAuth(config.auth);
            }
            
            if (config.headers) {
                client.setHeader(config.headers);
            }
            
            this.clients[serviceName] = client;
        }
        
        return this.clients[serviceName];
    }
    
    func callService(serviceName, operation, params) {
        let client = this.getClient(serviceName);
        return client.call(operation, params);
    }
}

// Uso
let soapManager = new EnterpriseSOAPManager();

// Llamadas a diferentes servicios
let userInfo = soapManager.callService("UserService", "GetUser", {"userId": 123});
let order = soapManager.callService("OrderService", "CreateOrder", {"userId": 123, "items": []});
let payment = soapManager.callService("PaymentService", "ProcessPayment", {"orderId": order.result.orderId});
```

### 3. Proxy para Múltiples Ambientes

```javascript
func createEnvironmentProxy(environment) {
    let baseURLs = {
        "dev": "https://dev-services.company.com",
        "staging": "https://staging-services.company.com", 
        "prod": "https://services.company.com"
    };
    
    let proxy = {
        "environment": environment,
        "baseURL": baseURLs[environment]
    };
    
    proxy.getClient = func(serviceName) {
        let wsdlURL = proxy.baseURL + "/" + serviceName + "/service.wsdl";
        let client = soapClient(wsdlURL);
        
        // Configuración por ambiente
        if (environment == "dev") {
            client.setTLSConfig({"skipVerify": true});
            client.setHeader("X-Debug", "true");
        } else {
            client.setTLSConfig({"minVersion": "1.2", "skipVerify": false});
        }
        
        client.setHeader("X-Environment", environment);
        return client;
    };
    
    return proxy;
}

// Uso
let devProxy = createEnvironmentProxy("dev");
let prodProxy = createEnvironmentProxy("prod");

let devClient = devProxy.getClient("calculator");
let prodClient = prodProxy.getClient("calculator");
```

---

## Troubleshooting

### Problemas Comunes y Soluciones

#### 1. "Connection Reset by Peer"

**Síntomas:**
```
panic: soapClient: failed to create client: connection reset by peer
```

**Causas:**
- Servicio bloqueando User-Agent no-browser
- Firewall corporativo
- Servicio temporalmente no disponible

**Soluciones:**
```javascript
// Cambiar User-Agent a uno de navegador
let client = soapClient("http://service.com/service.wsdl", {
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
});

// O usar headers corporativos específicos
client.setHeader({
    "User-Agent": "CompanyApp/1.0",
    "X-Requested-With": "SOAPClient"
});
```

#### 2. "Certificate Verification Failed"

**Síntomas:**
```
panic: x509: certificate signed by unknown authority
```

**Soluciones:**
```javascript
// Para testing (NO producción)
client.setTLSConfig({"skipVerify": true});

// Para producción - agregar CA personalizada
client.addCACert("/path/to/company-ca.crt");
```

#### 3. "Operation Not Found"

**Síntomas:**
```
panic: operation 'NonExistent' not found
```

**Diagnóstico:**
```javascript
// Listar operaciones disponibles
let operations = client.listOperations();
print("Operaciones disponibles:", operations);

// Verificar información del servicio
print("Namespace:", client.namespace);
print("Service URL:", client.serviceURL);
```

#### 4. "Invalid SOAP Response"

**Síntomas:**
- Respuestas vacías o malformadas
- Parsing failures

**Diagnóstico:**
```javascript
// Obtener respuesta raw para debug
let rawResponse = client.callRaw("Operation", {"param": "value"});
print("Respuesta XML cruda:", rawResponse);

// Verificar estructura
if (indexOf(rawResponse, "soap:Fault") != -1) {
    print("❌ SOAP Fault detectado en respuesta");
}
```

### Herramientas de Debugging

#### 1. Logging Detallado

```javascript
func debugSOAPCall(client, operation, params) {
    print("🔍 DEBUG - SOAP Call");
    print("   Service URL:", client.serviceURL);
    print("   Operation:", operation);
    print("   Parameters:", params);
    
    // Headers actuales
    let headers = client.getHeaders();
    print("   Headers:", Object.keys(headers));
    
    try {
        let startTime = getCurrentTime();
        let result = client.call(operation, params);
        let duration = getCurrentTime() - startTime;
        
        print("✅ DEBUG - Success");
        print("   Duration:", duration, "ms");
        print("   Success:", result.success);
        print("   Result type:", typeOf(result.result));
        
        return result;
    } catch (error) {
        print("❌ DEBUG - Error:", error);
        throw error;
    }
}
```

#### 2. Validador de WSDL

```javascript
func validateWSDL(wsdlURL) {
    try {
        print("🔍 Validando WSDL:", wsdlURL);
        
        let client = soapClient(wsdlURL);
        print("✅ WSDL válido y parseado");
        
        let operations = client.listOperations();
        print("📋 Operaciones encontradas:", operations.length);
        
        for (let op in operations) {
            let opInfo = client.getOperation(op);
            print("   -", op, "| Action:", opInfo.soapAction);
        }
        
        print("🌐 Endpoint:", client.serviceURL);
        print("📦 Namespace:", client.namespace);
        
        return true;
    } catch (error) {
        print("❌ WSDL inválido:", error);
        return false;
    }
}
```

---

## Referencia Técnica

### Estructura de Respuestas

#### Respuesta Completa (call)

```javascript
{
    "success": true|false,
    "result": valor_extraido,
    "values": {
        "elemento1": "valor1",
        "elemento2": "valor2"
    },
    "body": "contenido_soap_body",
    "raw": "xml_completo_respuesta",
    "cleaned": "xml_sin_namespaces",
    "fault": {  // Solo si success = false
        "code": "codigo_fault",
        "message": "mensaje_error"
    }
}
```

#### Respuesta Simple (callSimple)

```javascript
// Retorna directamente el valor resultado
42  // Para operaciones matemáticas
"string result"  // Para operaciones de texto
true  // Para operaciones booleanas
```

#### Respuesta Raw (callRaw)

```javascript
// Retorna el XML SOAP completo como string
"<?xml version=\"1.0\"?><soap:Envelope>...</soap:Envelope>"
```

### Configuraciones TLS Soportadas

```javascript
{
    "minVersion": "1.0"|"1.1"|"1.2"|"1.3",
    "skipVerify": true|false
}
```

### Tipos de Autenticación

```javascript
// Basic Authentication
{
    "type": "basic",
    "username": "string",
    "password": "string"
}

// Bearer Token
{
    "type": "bearer", 
    "token": "string"
}
```

### Headers por Defecto

```javascript
{
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
    "Accept": "text/xml,application/xml,*/*",
    "Accept-Language": "en-US,en;q=0.9",
    "Accept-Encoding": "gzip, deflate",
    "Connection": "keep-alive"
}
```

### Códigos de Error HTTP Comunes

- **200**: Éxito
- **400**: Bad Request - Parámetros inválidos
- **401**: Unauthorized - Falta autenticación
- **403**: Forbidden - Sin permisos
- **404**: Not Found - Endpoint no existe
- **500**: Internal Server Error - Error del servidor
- **502**: Bad Gateway - Proxy error
- **503**: Service Unavailable - Servicio no disponible

---

## Ejemplos de Integración Empresarial

### Sistema Bancario

```javascript
func bankingIntegration() {
    let bankClient = soapClient("https://bank-api.company.com/banking.wsdl");
    
    // Autenticación bancaria segura
    bankClient.setAuth({
        "type": "basic",
        "username": getSecureConfig("bank.username"),
        "password": getSecureConfig("bank.password")
    });
    
    // Headers de compliance bancario
    bankClient.setHeader({
        "X-Institution-ID": "BANK-12345",
        "X-Regulatory-Compliance": "PCI-DSS",
        "X-Audit-Trail": "enabled"
    });
    
    // SSL estricto para transacciones
    bankClient.setTLSConfig({
        "minVersion": "1.3",
        "skipVerify": false
    });
    
    // Consultar balance
    let balance = bankClient.callSimple("GetAccountBalance", {
        "accountNumber": "1234567890",
        "currency": "USD"
    });
    
    return balance;
}
```

### ERP Integration

```javascript
func erpIntegration(salesOrder) {
    let erpClient = soapClient("https://erp.company.com/sales.wsdl");
    
    // Headers ERP
    erpClient.setHeader({
        "X-ERP-Module": "Sales",
        "X-Transaction-Type": "Order",
        "X-Priority": "High"
    });
    
    // Procesar orden de venta
    let result = erpClient.call("ProcessSalesOrder", {
        "customerId": salesOrder.customerId,
        "items": salesOrder.items,
        "totalAmount": salesOrder.total,
        "currency": salesOrder.currency
    });
    
    if (result.success) {
        // Actualizar inventario
        let inventoryUpdate = erpClient.call("UpdateInventory", {
            "orderId": result.result.orderId,
            "items": salesOrder.items
        });
        
        return {
            "orderId": result.result.orderId,
            "inventoryUpdated": inventoryUpdate.success
        };
    } else {
        throw "Error en ERP: " + result.fault.message;
    }
}
```

---

## Conclusión

r2soap es una solución completa y robusta para integración SOAP en R2Lang. Su arquitectura dinámica elimina la necesidad de generación de código mientras proporciona todas las características necesarias para uso empresarial.

### Características Destacadas

- ✅ **100% Dinámico**: Sin generación de código
- ✅ **Enterprise-Ready**: Headers, autenticación, SSL/TLS
- ✅ **Parsing Inteligente**: Conversión automática a tipos R2Lang
- ✅ **Manejo de Errores**: Detección automática de SOAP faults
- ✅ **Flexible**: Múltiples formatos de respuesta
- ✅ **Seguro**: Configuraciones SSL/TLS empresariales

### Roadmap Futuro

- Soporte para SOAP 1.2
- WS-Security headers
- Attachments MTOM
- Caché de WSDL local
- Métricas y monitoring integrados

---

**Versión del Manual:** 1.0  
**Última Actualización:** 2024  
**Compatibilidad:** R2Lang v2.0+