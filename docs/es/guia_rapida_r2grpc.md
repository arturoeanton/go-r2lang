# Guía Rápida: Cliente gRPC Dinámico R2Lang

## 🚀 Introducción

Esta guía te ayudará a empezar con el cliente gRPC dinámico de R2Lang en menos de 5 minutos. No necesitas generar código, solo un archivo .proto y ¡listo!

## ⚡ Instalación

El cliente gRPC ya está incluido en R2Lang. Solo necesitas:

```bash
# Tener un archivo .proto
# Tener acceso a un servidor gRPC
# ¡Eso es todo!
```

## 🎯 Ejemplo Básico

### 1. Archivo Proto de Ejemplo

```protobuf
// user.proto
syntax = "proto3";

service UserService {
    rpc GetUser(GetUserRequest) returns (GetUserResponse);
    rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
}

message GetUserRequest {
    string user_id = 1;
}

message GetUserResponse {
    string user_id = 1;
    string name = 2;
    string email = 3;
}

message CreateUserRequest {
    string name = 1;
    string email = 2;
}

message CreateUserResponse {
    string user_id = 1;
    string message = 2;
}
```

### 2. Código R2Lang

```javascript
// user_client.r2

// Crear cliente gRPC
let client = grpc.grpcClient("user.proto", "localhost:9090");

// Ver qué servicios están disponibles
std.print("Servicios:", client.listServices());

// Ver métodos del servicio
std.print("Métodos:", client.listMethods("UserService"));

// Obtener información de un método
let info = client.getMethodInfo("UserService", "GetUser");
std.print("Info GetUser:", info);

// Hacer una llamada simple
let user = client.callSimple("UserService", "GetUser", {
    "user_id": "123"
});

std.print("Usuario:", user);
```

### 3. Ejecutar

```bash
go run main.go user_client.r2
```

## 📡 Ejemplos por Tipo de RPC

### Unary RPC (Simple)

```javascript
// Una solicitud, una respuesta
let user = client.callSimple("UserService", "GetUser", {
    "user_id": "123"
});

std.print("Usuario encontrado:", user.name);
```

### Server Streaming

```javascript
// Una solicitud, múltiples respuestas
let stream = client.callServerStream("UserService", "ListUsers", {
    "department": "engineering"
});

stream.onReceive(func(user) {
    std.print("Usuario:", user.name);
});

stream.onClose(func() {
    std.print("Stream terminado");
});
```

### Client Streaming

```javascript
// Múltiples solicitudes, una respuesta
let stream = client.callClientStream("UserService", "BatchCreateUsers");

// Enviar múltiples usuarios
stream.send({"name": "Ana", "email": "ana@company.com"});
stream.send({"name": "Luis", "email": "luis@company.com"});
stream.send({"name": "María", "email": "maria@company.com"});

// Cerrar y obtener resultado
stream.closeSend();
```

### Bidirectional Streaming

```javascript
// Múltiples solicitudes, múltiples respuestas
let chat = client.callBidirectionalStream("ChatService", "LiveChat");

// Recibir mensajes
chat.onReceive(func(message) {
    std.print("Mensaje:", message.text);
});

// Enviar mensaje
chat.send({
    "user": "r2lang-user",
    "text": "¡Hola mundo!"
});
```

## 🔐 Configuración Rápida

### Autenticación JWT

```javascript
// Configurar token JWT
client.setAuth({
    "type": "bearer",
    "token": "eyJhbGciOiJIUzI1NiIs..."
});

// Ahora todas las llamadas usarán el token
let user = client.callSimple("UserService", "GetUser", {"user_id": "123"});
```

### Headers Personalizados

```javascript
// Agregar headers personalizados
client.setMetadata({
    "x-api-key": "mi-api-key",
    "x-version": "1.0"
});
```

### TLS/SSL

```javascript
// Habilitar TLS
client.setTLSConfig({
    "enabled": true,
    "minVersion": "1.2"
});
```

## 🛠️ Casos de Uso Comunes

### 1. API REST → gRPC

```javascript
// Antes (REST)
// GET /users/123

// Ahora (gRPC)
let user = client.callSimple("UserService", "GetUser", {"user_id": "123"});
```

### 2. Microservicios

```javascript
// Conectar con múltiples microservicios
let userService = grpc.grpcClient("user.proto", "users:9090");
let orderService = grpc.grpcClient("orders.proto", "orders:9090");
let paymentService = grpc.grpcClient("payments.proto", "payments:9090");

// Workflow completo
let user = userService.callSimple("UserService", "GetUser", {"user_id": "123"});
let order = orderService.callSimple("OrderService", "CreateOrder", {
    "user_id": user.user_id,
    "items": [{"product_id": "abc", "quantity": 2}]
});
let payment = paymentService.callSimple("PaymentService", "ProcessPayment", {
    "order_id": order.order_id,
    "amount": order.total
});
```

### 3. Monitoreo en Tiempo Real

```javascript
// Stream de métricas
let metrics = client.callServerStream("MetricsService", "WatchMetrics", {
    "resources": ["cpu", "memory"]
});

metrics.onReceive(func(metric) {
    std.print("CPU:", metric.cpu_usage + "%");
    std.print("Memory:", metric.memory_usage + "%");
});
```

## 🔍 Debugging Rápido

### Ver Respuesta Completa

```javascript
// En lugar de callSimple(), usa call() para ver todo
let response = client.call("UserService", "GetUser", {"user_id": "123"});
std.print("Respuesta completa:", response);

// response.success - true/false
// response.result - datos del usuario
// response.error - detalles del error (si falla)
```

### Manejo de Errores

```javascript
let response = client.call("UserService", "GetUser", {"user_id": "invalid"});

if (response.success) {
    std.print("Usuario:", response.result);
} else {
    std.print("Error:", response.error.code);
    std.print("Mensaje:", response.error.message);
}
```

## 📋 Checklist de Troubleshooting

### ✅ Problemas Comunes

1. **"Service not found"**
   ```javascript
   // Verificar servicios disponibles
   std.print(client.listServices());
   ```

2. **"Method not found"**
   ```javascript
   // Verificar métodos del servicio
   std.print(client.listMethods("UserService"));
   ```

3. **"Connection refused"**
   ```javascript
   // Verificar que el servidor esté corriendo
   // Verificar la dirección y puerto
   ```

4. **"Authentication failed"**
   ```javascript
   // Verificar token o credenciales
   client.setAuth({
       "type": "bearer",
       "token": "token-correcto"
   });
   ```

5. **"TLS handshake failed"**
   ```javascript
   // Para desarrollo local, puedes deshabilitar TLS
   client.setTLSConfig({
       "enabled": false
   });
   ```

## 🎯 Ejemplos Prácticos

### Ejemplo 1: CRUD Básico

```javascript
// crud_example.r2
let client = grpc.grpcClient("user.proto", "localhost:9090");

// CREATE
let newUser = client.callSimple("UserService", "CreateUser", {
    "name": "Juan Pérez",
    "email": "juan@example.com"
});
std.print("Usuario creado:", newUser.user_id);

// READ
let user = client.callSimple("UserService", "GetUser", {
    "user_id": newUser.user_id
});
std.print("Usuario leído:", user.name);

// UPDATE
let updatedUser = client.callSimple("UserService", "UpdateUser", {
    "user_id": user.user_id,
    "name": "Juan Carlos Pérez"
});
std.print("Usuario actualizado:", updatedUser.name);

// DELETE
let result = client.callSimple("UserService", "DeleteUser", {
    "user_id": user.user_id
});
std.print("Usuario eliminado:", result.success);
```

### Ejemplo 2: Chat en Tiempo Real

```javascript
// chat_example.r2
let client = grpc.grpcClient("chat.proto", "localhost:9090");

// Configurar autenticación
client.setAuth({
    "type": "bearer",
    "token": "mi-jwt-token"
});

// Crear stream de chat
let chat = client.callBidirectionalStream("ChatService", "LiveChat");

// Configurar recepción de mensajes
chat.onReceive(func(message) {
    std.print("[" + message.user + "]: " + message.text);
});

// Enviar mensaje de bienvenida
chat.send({
    "user": "r2lang-bot",
    "text": "¡Hola! Soy un bot de R2Lang"
});

// Simular conversación
setTimeout(func() {
    chat.send({
        "user": "r2lang-bot",
        "text": "¿Cómo están todos?"
    });
}, 2000);

// Cerrar chat después de 10 segundos
setTimeout(func() {
    chat.close();
    std.print("Chat cerrado");
}, 10000);
```

### Ejemplo 3: Monitoreo de Sistema

```javascript
// monitoring_example.r2
let client = grpc.grpcClient("monitoring.proto", "localhost:9090");

// Stream de métricas del sistema
let metrics = client.callServerStream("MonitoringService", "WatchSystemMetrics", {
    "interval_seconds": 5,
    "metrics": ["cpu", "memory", "disk"]
});

metrics.onReceive(func(metric) {
    let timestamp = new Date(metric.timestamp * 1000);
    std.print(timestamp + " - " + metric.type + ": " + metric.value + metric.unit);
    
    // Alertas automáticas
    if (metric.type == "cpu" && metric.value > 80) {
        std.print("🚨 ALERTA: CPU alta (" + metric.value + "%)");
    }
    
    if (metric.type == "memory" && metric.value > 90) {
        std.print("🚨 ALERTA: Memoria alta (" + metric.value + "%)");
    }
});

metrics.onError(func(error) {
    std.print("Error en métricas:", error);
});

std.print("Monitoreo iniciado... (presiona Ctrl+C para salir)");
```

## 🚀 Próximos Pasos

1. **Lee el [Manual Completo](./manual_r2grpc.md)** para funcionalidades avanzadas
2. **Practica con tus propios archivos .proto**
3. **Experimenta con streaming RPC**
4. **Configura autenticación y TLS para producción**
5. **Integra con tus microservicios existentes**

## 📚 Recursos Adicionales

- [Manual del Desarrollador](./manual_r2grpc.md)
- [Ejemplos de R2Lang](../../examples/)
- [Documentación de gRPC](https://grpc.io/docs/)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)

---

**¡Feliz codificación con R2Lang y gRPC!** 🎉

---

**Versión**: 1.0.0  
**Fecha**: 2025-01-18  
**Autor**: R2Lang Team