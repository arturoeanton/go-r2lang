# Manual del Desarrollador: Cliente gRPC Dinámico para R2Lang

## 🎯 Introducción

El cliente gRPC dinámico de R2Lang permite invocar servicios gRPC sin generar código, utilizando únicamente archivos Protocol Buffers (.proto). Esta implementación sigue la misma filosofía exitosa de r2soap, proporcionando una API consistente y empresarial.

## 📋 Características Principales

### ✅ **100% Dinámico**
- No genera código cliente
- Solo requiere archivos .proto
- Parsing automático de servicios y métodos
- Introspección completa de la API

### ✅ **Tipos de Streaming Completos**
- **Unary RPC**: Solicitud-respuesta simple
- **Server Streaming**: Servidor envía múltiples respuestas
- **Client Streaming**: Cliente envía múltiples solicitudes
- **Bidirectional Streaming**: Comunicación bidireccional

### ✅ **Configuración Empresarial**
- Autenticación (Bearer, Basic, mTLS, Custom)
- Configuración TLS/SSL completa
- Manejo de metadatos personalizados
- Timeouts y compresión

## 🚀 Inicio Rápido

### 1. Creación del Cliente

```javascript
// Crear cliente gRPC dinámico
let client = grpc.grpcClient("path/to/service.proto", "server.com:443");

// Con metadatos personalizados
let client = grpc.grpcClient(
    "user-service.proto", 
    "api.company.com:443", 
    {
        "x-api-version": "1.0",
        "x-client-name": "R2Lang-Client"
    }
);
```

### 2. Introspección de Servicios

```javascript
// Listar servicios disponibles
let services = client.listServices();
std.print("Servicios disponibles:", services);

// Listar métodos de un servicio
let methods = client.listMethods("UserService");
std.print("Métodos de UserService:", methods);

// Obtener información de un método
let methodInfo = client.getMethodInfo("UserService", "GetUser");
std.print("Información del método:", methodInfo);
```

### 3. Invocación de Métodos

```javascript
// Llamada unary simple
let user = client.call("UserService", "GetUser", {
    "user_id": "12345"
});

// Respuesta simplificada (solo el resultado)
let userData = client.callSimple("UserService", "GetUser", {
    "user_id": "12345"
});

// Llamada con manejo de errores
let response = client.call("UserService", "GetUser", {
    "user_id": "invalid"
});

if (response.success) {
    std.print("Usuario encontrado:", response.result);
} else {
    std.print("Error:", response.error);
}
```

## 🔧 Configuración Avanzada

### Autenticación

#### Bearer Token (JWT)
```javascript
client.setAuth({
    "type": "bearer",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
});
```

#### Basic Authentication
```javascript
client.setAuth({
    "type": "basic",
    "username": "api_user",
    "password": "secure_password"
});
```

#### mTLS (Mutual TLS)
```javascript
client.setAuth({
    "type": "mtls",
    "certFile": "/path/to/client.crt",
    "keyFile": "/path/to/client.key",
    "caFile": "/path/to/ca.crt"
});
```

#### Custom Metadata
```javascript
client.setAuth({
    "type": "custom",
    "metadata": {
        "x-api-key": "your-api-key",
        "x-tenant-id": "tenant-123",
        "authorization": "Bearer " + token
    }
});
```

### Configuración TLS

```javascript
// Configuración TLS básica
client.setTLSConfig({
    "enabled": true,
    "minVersion": "1.2",
    "skipVerify": false,
    "serverName": "api.company.com"
});

// Configuración TLS estricta
client.setTLSConfig({
    "enabled": true,
    "minVersion": "1.3",
    "skipVerify": false,
    "serverName": "secure-api.company.com"
});
```

### Manejo de Metadatos

```javascript
// Configurar metadatos individuales
client.setMetadata("x-request-id", "req-12345");
client.setMetadata("x-user-role", "admin");

// Configurar múltiples metadatos
client.setMetadata({
    "x-company-id": "CORP-789",
    "x-api-version": "2.0",
    "x-feature-flags": "beta,premium"
});

// Obtener metadatos actuales
let metadata = client.getMetadata();
std.print("Metadatos configurados:", metadata);
```

### Timeouts y Compresión

```javascript
// Configurar timeout (en segundos)
client.setTimeout(60.0);

// Configurar compresión
client.setCompression("gzip");
```

## 📡 Streaming RPC

### Server Streaming

```javascript
// Crear stream de servidor
let userStream = client.callServerStream("UserService", "WatchUsers", {
    "department": "engineering",
    "active_only": true
});

// Configurar callbacks
userStream.onReceive(func(user) {
    std.print("Usuario actualizado:", user.name, user.status);
});

userStream.onError(func(error) {
    std.print("Error en stream:", error);
});

userStream.onClose(func() {
    std.print("Stream cerrado");
});
```

### Client Streaming

```javascript
// Crear stream de cliente
let uploadStream = client.callClientStream("FileService", "UploadFile");

// Enviar múltiples mensajes
let fileChunks = [
    {"filename": "document.pdf", "chunk": "chunk1", "size": 1024},
    {"filename": "document.pdf", "chunk": "chunk2", "size": 1024},
    {"filename": "document.pdf", "chunk": "chunk3", "size": 512}
];

for (let chunk in fileChunks) {
    uploadStream.send(chunk);
}

// Cerrar stream y obtener resultado
uploadStream.closeSend();
```

### Bidirectional Streaming

```javascript
// Crear stream bidireccional
let chatStream = client.callBidirectionalStream("ChatService", "LiveChat");

// Configurar recepción
chatStream.onReceive(func(message) {
    std.print("Mensaje recibido:", message.user, "->", message.text);
    
    // Responder automáticamente
    if (message.text.contains("ping")) {
        chatStream.send({
            "user": "bot",
            "text": "pong!",
            "timestamp": std.now()
        });
    }
});

// Enviar mensaje inicial
chatStream.send({
    "user": "r2lang-user",
    "text": "¡Hola desde R2Lang!",
    "timestamp": std.now()
});

// Cerrar stream después de un tiempo
setTimeout(func() {
    chatStream.close();
}, 10000);
```

## 🎯 Patrones de Uso Comunes

### 1. Servicio de Usuarios

```javascript
// user-service.proto
// service UserService {
//     rpc GetUser(GetUserRequest) returns (GetUserResponse);
//     rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
//     rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
//     rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
// }

let userService = grpc.grpcClient("user-service.proto", "users.company.com:443");

// Configurar autenticación
userService.setAuth({
    "type": "bearer",
    "token": getJWTToken()
});

// CRUD operations
func getUser(userId) {
    return userService.callSimple("UserService", "GetUser", {
        "user_id": userId
    });
}

func createUser(userData) {
    return userService.callSimple("UserService", "CreateUser", userData);
}

func updateUser(userId, updates) {
    return userService.callSimple("UserService", "UpdateUser", {
        "user_id": userId,
        "updates": updates
    });
}

func deleteUser(userId) {
    return userService.callSimple("UserService", "DeleteUser", {
        "user_id": userId
    });
}
```

### 2. Microservicios con Load Balancing

```javascript
// Configurar múltiples instancias
let orderServices = [
    grpc.grpcClient("order-service.proto", "orders1.company.com:443"),
    grpc.grpcClient("order-service.proto", "orders2.company.com:443"),
    grpc.grpcClient("order-service.proto", "orders3.company.com:443")
];

// Configurar autenticación para todos
for (let service in orderServices) {
    service.setAuth({
        "type": "bearer",
        "token": getServiceToken()
    });
    service.setTimeout(30.0);
}

// Load balancer simple
let currentService = 0;
func getNextOrderService() {
    let service = orderServices[currentService];
    currentService = (currentService + 1) % orderServices.length;
    return service;
}

// Crear orden con retry
func createOrder(orderData) {
    let maxRetries = 3;
    for (let i = 0; i < maxRetries; i++) {
        try {
            let service = getNextOrderService();
            return service.callSimple("OrderService", "CreateOrder", orderData);
        } catch (error) {
            std.print("Error en intento", i + 1, ":", error);
            if (i == maxRetries - 1) {
                throw error;
            }
        }
    }
}
```

### 3. Monitoreo y Métricas

```javascript
// Servicio de métricas con streaming
let metricsService = grpc.grpcClient("metrics-service.proto", "metrics.company.com:443");

// Configurar monitoreo en tiempo real
let metricsStream = metricsService.callServerStream("MetricsService", "WatchMetrics", {
    "resources": ["cpu", "memory", "disk", "network"],
    "interval_seconds": 5,
    "threshold_alerts": true
});

// Procesar métricas en tiempo real
metricsStream.onReceive(func(metric) {
    std.print("[METRIC]", metric.timestamp, metric.resource, "->", metric.value, metric.unit);
    
    // Alertas automáticas
    if (metric.alert) {
        sendAlert(metric);
    }
    
    // Almacenar en base de datos
    storeMetric(metric);
});

metricsStream.onError(func(error) {
    std.print("[ERROR] Métricas:", error);
    // Reconectar automáticamente
    setTimeout(func() {
        restartMetricsStream();
    }, 5000);
});
```

## 🔍 Debugging y Troubleshooting

### Logging Detallado

```javascript
// Habilitar logging detallado
client.setMetadata({
    "x-debug": "true",
    "x-log-level": "debug",
    "x-trace-id": "trace-" + std.uuid()
});

// Llamada con información extendida
let response = client.call("UserService", "GetUser", {"user_id": "123"});

std.print("Respuesta completa:", response);
std.print("Metadatos de respuesta:", response.metadata);
std.print("Tiempo de respuesta:", response.response_time);
```

### Manejo de Errores

```javascript
func robustGRPCCall(service, method, params) {
    try {
        let response = service.call(method, params);
        
        if (response.success) {
            return response.result;
        } else {
            // Manejo específico por código de error
            if (response.error.code == "UNAUTHENTICATED") {
                refreshToken();
                return robustGRPCCall(service, method, params);
            } else if (response.error.code == "UNAVAILABLE") {
                // Esperar y reintentar
                sleep(1000);
                return robustGRPCCall(service, method, params);
            } else {
                throw new Error("gRPC Error: " + response.error.message);
            }
        }
    } catch (error) {
        std.print("Error en llamada gRPC:", error);
        throw error;
    }
}
```

### Testing de Servicios

```javascript
// Test unitario de servicio gRPC
describe("UserService gRPC", func() {
    let client = grpc.grpcClient("user-service.proto", "localhost:9090");
    
    it("should create user successfully", func() {
        let userData = {
            "name": "Test User",
            "email": "test@example.com",
            "department": "engineering"
        };
        
        let response = client.call("UserService", "CreateUser", userData);
        
        assert.equals(response.success, true);
        assert.notNull(response.result.user_id);
        assert.equals(response.result.name, userData.name);
    });
    
    it("should handle invalid user data", func() {
        let invalidData = {
            "name": "",
            "email": "invalid-email"
        };
        
        let response = client.call("UserService", "CreateUser", invalidData);
        
        assert.equals(response.success, false);
        assert.equals(response.error.code, "INVALID_ARGUMENT");
    });
});
```

## 🎛️ Configuración de Producción

### Configuración Empresarial

```javascript
// Configuración para producción
let productionClient = grpc.grpcClient(
    "production-service.proto", 
    "api.company.com:443",
    {
        "x-environment": "production",
        "x-client-version": "1.0.0",
        "x-service-name": "r2lang-client"
    }
);

// Configuración de seguridad estricta
productionClient.setTLSConfig({
    "enabled": true,
    "minVersion": "1.3",
    "skipVerify": false,
    "serverName": "api.company.com"
});

// Autenticación con rotación automática
let tokenRefreshInterval = 3600; // 1 hora
setInterval(func() {
    let newToken = refreshServiceToken();
    productionClient.setAuth({
        "type": "bearer",
        "token": newToken
    });
}, tokenRefreshInterval * 1000);

// Timeout conservador para producción
productionClient.setTimeout(30.0);

// Compresión para optimizar ancho de banda
productionClient.setCompression("gzip");
```

### Monitoreo de Salud

```javascript
// Health check periódico
func healthCheck() {
    try {
        let response = client.call("HealthService", "Check", {});
        if (response.success && response.result.status == "SERVING") {
            return true;
        }
    } catch (error) {
        std.print("Health check failed:", error);
    }
    return false;
}

// Monitoreo continuo
setInterval(func() {
    if (!healthCheck()) {
        std.print("Service unhealthy, triggering alerts");
        sendHealthAlert();
    }
}, 30000); // cada 30 segundos
```

## 📚 Referencia de API

### Métodos del Cliente

| Método | Descripción | Parámetros |
|--------|-------------|------------|
| `listServices()` | Lista servicios disponibles | Ninguno |
| `listMethods(service)` | Lista métodos de un servicio | `service`: nombre del servicio |
| `getMethodInfo(service, method)` | Información de un método | `service`, `method`: nombres |
| `call(service, method, params)` | Llamada unary completa | `service`, `method`, `params`: datos |
| `callSimple(service, method, params)` | Llamada unary simple | `service`, `method`, `params`: datos |
| `callServerStream(service, method, params)` | Stream de servidor | `service`, `method`, `params`: datos |
| `callClientStream(service, method)` | Stream de cliente | `service`, `method`: nombres |
| `callBidirectionalStream(service, method)` | Stream bidireccional | `service`, `method`: nombres |

### Métodos de Configuración

| Método | Descripción | Parámetros |
|--------|-------------|------------|
| `setTimeout(seconds)` | Configura timeout | `seconds`: número |
| `setMetadata(key, value)` | Configura metadata | `key`, `value`: strings |
| `setMetadata(map)` | Configura múltiples metadatos | `map`: objeto |
| `getMetadata()` | Obtiene metadatos actuales | Ninguno |
| `setAuth(config)` | Configura autenticación | `config`: objeto |
| `setTLSConfig(config)` | Configura TLS | `config`: objeto |
| `setCompression(alg)` | Configura compresión | `alg`: string |
| `close()` | Cierra conexión | Ninguno |

### Métodos de Stream

| Método | Descripción | Parámetros |
|--------|-------------|------------|
| `send(message)` | Envía mensaje | `message`: objeto |
| `closeSend()` | Cierra envío | Ninguno |
| `onReceive(callback)` | Callback de recepción | `callback`: función |
| `onError(callback)` | Callback de error | `callback`: función |
| `onClose(callback)` | Callback de cierre | `callback`: función |
| `close()` | Cierra stream | Ninguno |

## 🚀 Conclusión

El cliente gRPC dinámico de R2Lang proporciona una solución empresarial completa para la integración con servicios gRPC modernos. Su API consistente con r2soap facilita la migración y el aprendizaje, mientras que sus características avanzadas de streaming y configuración lo hacen ideal para aplicaciones de producción.

Para más ejemplos y casos de uso, consulte la [Guía Rápida de r2grpc](./guia_rapida_r2grpc.md).

---

**Versión**: 1.0.0  
**Fecha**: 2025-01-18  
**Autor**: R2Lang Team  
**Licencia**: MIT