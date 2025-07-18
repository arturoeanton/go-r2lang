# 👋 Ejemplo 2: Hello World gRPC

**¡100% FUNCIONAL CON SERVIDOR IMPLEMENTADO!** El ejemplo más simple para empezar con gRPC en R2Lang.

## ✨ ¿Qué hace este ejemplo?

- ✅ **Servidor completamente implementado** con método SayHello real
- ✅ **Conecta** al servidor gRPC implementado
- ✅ **Descubre** el servicio Greeter automáticamente
- ✅ **Lista** métodos disponibles  
- ✅ **Realiza** llamadas gRPC y recibe respuestas REALES
- ✅ **Demuestra** configuración básica y avanzada

## 🚀 Instrucciones (2 comandos)

### ⚡ Opción A: Desde la raíz del proyecto

**Terminal 1 - Servidor implementado:**
```bash
cd examples/grpc/example2/server
go run greeter_server.go
```

**Terminal 2 - Cliente R2Lang:**
```bash
go run main.go examples/grpc/example2/hello_world_client.r2
```

### ⚡ Opción B: Desde example2

**Terminal 1 - Servidor implementado:**
```bash
cd examples/grpc/example2/server
go run greeter_server.go
```

**Terminal 2 - Cliente R2Lang:**
```bash
cd examples/grpc/example2
go run ../../../main.go hello_world_client.r2
```

**🎉 ¡AMBAS OPCIONES FUNCIONAN AL 100%!**

## 🎯 Resultado Esperado

```
🚀 Cliente gRPC R2Lang - Hello World
====================================

📋 Paso 1: Conectando al servidor gRPC...
✅ Cliente gRPC creado exitosamente!

📋 Paso 2: Configurando cliente...
✅ Cliente configurado!

📋 Paso 3: Descubriendo servicios...
Servicios encontrados: 1
  📦 Greeter

Métodos del servicio Greeter:
  ⚡ SayHello

📋 Paso 4: Realizando llamada gRPC...
✅ ¡Llamada exitosa!
📨 Respuesta del servidor:
   hello Desarrollador R2Lang

📋 Paso 5: Probando llamadas simplificadas...
✅ Respuesta simple: hello Usuario Simple

📋 Paso 6: Información del método...
📝 Método: SayHello
📥 Tipo entrada: HelloRequest
📤 Tipo salida: HelloReply
🔄 Client Streaming: false
🔄 Server Streaming: false

✅ Cliente cerrado correctamente

🎯 Demo completado!
```

## 📁 Archivos

### 🖥️ Servidor (server/)
- `greeter_server.go` - **Servidor gRPC completamente implementado**
- `greeter.proto` - Definición del protocolo
- `helloworld/` - Código Go generado por protoc
  - `greeter.pb.go` - Mensajes Protocol Buffers
  - `greeter_grpc.pb.go` - Servicio gRPC
- `go.mod` / `go.sum` - Dependencias Go

### 📱 Cliente R2Lang
- `hello_world_client.r2` - **Cliente R2Lang que funciona al 100%**
- `greeter.proto` - Copia del protocolo para referencia

### 📚 Documentación
- `README.md` - Esta guía

## 💡 Lo que demuestra

**¡Este ejemplo es perfecto para empezar con gRPC y ver respuestas REALES!**

### 🎯 Para Desarrolladores:
✅ **Servidor gRPC real** - Implementado completamente con Go  
✅ **Cliente R2Lang dinámico** - Sin generación de código  
✅ **Respuestas reales** - El método SayHello funciona al 100%  
✅ **Configuración mínima** - Solo 2 comandos para ejecutar  
✅ **Conexión end-to-end** - Cliente ↔ Servidor funcionando  
✅ **Fácil de entender** - Código simple y bien documentado  

### 🚀 Flujo completo:
1. **Servidor**: Implementa `SayHello` que devuelve `"hello {name}"`
2. **Cliente**: Se conecta dinámicamente y llama al método
3. **Respuesta**: Recibe la respuesta real del servidor
4. **Introspección**: Descubre servicios y métodos automáticamente

## 🔧 Diferencias con otros ejemplos

- **Ejemplo 1**: Muestra introspección avanzada y manejo de errores
- **Ejemplo 2**: Muestra comunicación end-to-end real con servidor implementado

**¡Ambos ejemplos funcionan al 100% y demuestran diferentes aspectos del cliente gRPC R2Lang!**

## 🎓 Para Desarrolladores

Este ejemplo demuestra que el **cliente gRPC dinámico de R2Lang está listo para producción**:

- ✅ No necesita generación de código
- ✅ Se conecta a cualquier servidor gRPC estándar  
- ✅ Funciona con reflection habilitado
- ✅ Maneja tipos de mensaje dinámicamente
- ✅ Perfecto para microservicios y APIs