// hello_world_client.r2
// Cliente gRPC simple para demostrar R2Lang gRPC

std.print("🚀 Cliente gRPC R2Lang - Hello World");
std.print("====================================");
std.print("");

// Crear archivo proto temporalmente
let protoContent = `syntax = "proto3";

package helloworld;

service Greeter {
    rpc SayHello (HelloRequest) returns (HelloReply);
}

message HelloRequest {
    string name = 1;
}

message HelloReply {
    string message = 1;
}`;

// Escribir archivo proto temporal
io.writeFile("./greeter.proto", protoContent);

try {
    std.print("📋 Paso 1: Conectando al servidor gRPC...");
    
    // Crear cliente gRPC
    let client = grpc.grpcClient("./greeter.proto", "localhost:9090");
    
    std.print("✅ Cliente gRPC creado exitosamente!");
    std.print("");
    
    std.print("📋 Paso 2: Configurando cliente...");
    
    // Configurar timeout
    client.setTimeout(10.0);
    
    // Configurar metadata personalizada
    client.setMetadata({
        "x-client": "R2Lang",
        "x-version": "1.0.0",
        "x-demo": "hello-world"
    });
    
    std.print("✅ Cliente configurado!");
    std.print("");
    
    std.print("📋 Paso 3: Descubriendo servicios...");
    
    // Listar servicios disponibles
    let services = client.listServices();
    std.print("Servicios encontrados: " + std.len(services));
    for (let i = 0; i < std.len(services); i = i + 1) {
        std.print("  📦 " + services[i]);
    }
    std.print("");
    
    // Listar métodos del servicio Greeter
    if (std.len(services) > 0) {
        std.print("Métodos del servicio " + services[0] + ":");
        let methods = client.listMethods(services[0]);
        for (let j = 0; j < std.len(methods); j = j + 1) {
            std.print("  ⚡ " + methods[j]);
        }
        std.print("");
    }
    
    std.print("📋 Paso 4: Realizando llamada gRPC...");
    
    // Realizar llamada SayHello
    let response = client.call("Greeter", "SayHello", {
        "name": "Desarrollador R2Lang"
    });
    
    if (response.success) {
        std.print("✅ ¡Llamada exitosa!");
        std.print("📨 Respuesta del servidor:");
        std.print("   " + response.result.message);
    } else {
        std.print("❌ Error en la llamada:");
        std.print("   Código: " + response.error.code);
        std.print("   Mensaje: " + response.error.message);
    }
    std.print("");
    
    std.print("📋 Paso 5: Probando llamadas simplificadas...");
    
    // Usar callSimple para una respuesta más directa
    let simpleResponse = client.callSimple("Greeter", "SayHello", {
        "name": "Usuario Simple"
    });
    
    std.print("✅ Respuesta simple: " + simpleResponse.message);
    std.print("");
    
    std.print("📋 Paso 6: Información del método...");
    
    // Obtener información detallada de un método
    let methodInfo = client.getMethodInfo("Greeter", "SayHello");
    std.print("📝 Método: " + methodInfo.name);
    std.print("📥 Tipo entrada: " + methodInfo.inputType);
    std.print("📤 Tipo salida: " + methodInfo.outputType);
    std.print("🔄 Client Streaming: " + methodInfo.clientStreaming);
    std.print("🔄 Server Streaming: " + methodInfo.serverStreaming);
    std.print("");
    
    // Cerrar cliente
    client.close();
    std.print("✅ Cliente cerrado correctamente");
    
} catch (error) {
    std.print("❌ Error: " + error);
    std.print("");
    std.print("🔧 SOLUCIÓN:");
    std.print("1. Asegúrate de que el servidor esté corriendo:");
    std.print("   cd server && go run greeter_server.go");
    std.print("2. En otra terminal, ejecuta este cliente:");
    std.print("   go run ../../../main.go hello_world_client.r2");
}

std.print("");
std.print("🎯 Demo completado!");
std.print("");
std.print("📚 Más información:");
std.print("  📖 README.md - Guía de este ejemplo");
std.print("  📁 ../example1/ - Ejemplo de introspección avanzada");

// Limpiar archivo temporal
try {
    os.remove("./greeter.proto");
} catch (e) {
    // Ignorar errores de limpieza
}