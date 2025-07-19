// example32-requests.r2
// Demostración de la librería r2requests (estilo Python Requests)
// Este ejemplo muestra cómo usar el cliente HTTP moderno de R2Lang

std.print("=== R2REQUESTS LIBRARY DEMO ===\n");

// ==========================================
// 1. GET REQUEST SIMPLE
// ==========================================
std.print("1. GET Request Simple:");
try {
    let response = request.get("https://jsonplaceholder.typicode.com/posts/1");
    std.print("Status Code:", response["status_code"]);
    std.print("OK:", response["ok"]);
    std.print("URL:", response["url"]);
    
    if (response["ok"]) {
        let post = response["json"];
        std.print("Post Title:", post["title"]);
        std.print("Post Body:", post["body"]);
        std.print("User ID:", post["userId"]);
    }
} catch (error) {
    std.print("Error en GET:", error);
}

std.print("\n==================================================\n");

// ==========================================
// 2. MÚLTIPLES MÉTODOS HTTP
// ==========================================
std.print("2. Múltiples Métodos HTTP:");
try {
    let getResponse = request.get("https://jsonplaceholder.typicode.com/posts/2");
    std.print("GET Status:", getResponse["status_code"]);
    
    let deleteResponse = request.delete("https://jsonplaceholder.typicode.com/posts/1");
    std.print("DELETE Status:", deleteResponse["status_code"]);
    
} catch (error) {
    std.print("Error en métodos HTTP:", error);
}

std.print("\n==================================================\n");

// ==========================================
// 3. USANDO SESIONES
// ==========================================
std.print("3. Usando Sesiones para Reutilización:");
try {
    let s = request.session();
    
    // Múltiples requests con la misma sesión
    let user1 = s["get"]("https://jsonplaceholder.typicode.com/users/1");
    let user2 = s["get"]("https://jsonplaceholder.typicode.com/users/2");
    
    if (user1["ok"] == true) {
        if (user2["ok"] == true) {
            std.print("Usuario 1:", user1["json"]["name"], "- Email:", user1["json"]["email"]);
            std.print("Usuario 2:", user2["json"]["name"], "- Email:", user2["json"]["email"]);
        }
    }
    
    // Cerrar sesión
    s["close"]();
} catch (error) {
    std.print("Error en sesiones:", error);
}

std.print("\n==================================================\n");

// ==========================================
// 4. MANEJO DE ERRORES HTTP
// ==========================================
std.print("4. Manejo de Errores HTTP:");
try {
    let response = request.get("https://jsonplaceholder.typicode.com/posts/999999");
    std.print("Status Code:", response["status_code"]);
    std.print("OK:", response["ok"]);
    
    if (response["ok"] == false) {
        std.print("Error HTTP detectado!");
        std.print("Texto de error:", response["text"]);
        
        // Intentar raise_for_status (esto causará un panic para errores)
        try {
            response["raise_for_status"]();
        } catch (statusError) {
            std.print("raise_for_status detectó error:", statusError);
        }
    }
} catch (error) {
    std.print("Error general:", error);
}

std.print("\n==================================================\n");

// ==========================================
// 5. UTILIDADES DE URL
// ==========================================
std.print("5. Utilidades de URL:");
try {
    let originalText = "Hello World! ¿Cómo estás?";
    let encoded = request.urlencode(originalText);
    let decoded = request.urldecode(encoded);
    
    std.print("Original:", originalText);
    std.print("Encoded:", encoded);
    std.print("Decoded:", decoded);
    std.print("Match:", originalText == decoded);
    
} catch (error) {
    std.print("Error en utilidades URL:", error);
}

std.print("\n==================================================\n");

// ==========================================
// 6. HEADERS Y INFORMACIÓN DE RESPUESTA
// ==========================================
std.print("6. Información Detallada de Respuesta:");
try {
    let response = request.get("https://jsonplaceholder.typicode.com/posts/1");
    
    std.print("URL final:", response["url"]);
    std.print("Status Code:", response["status_code"]);
    std.print("OK:", response["ok"]);
    std.print("Tiempo transcurrido:", response["elapsed"], "segundos");
    std.print("Tamaño del contenido:", std.len(response["text"]), "caracteres");
    
    // Headers más comunes
    std.print("\nHeaders importantes:");
    let headers = response["headers"];
    if (headers["Content-Type"]) {
        std.print("Content-Type:", headers["Content-Type"]);
    }
    if (headers["Content-Length"]) {
        std.print("Content-Length:", headers["Content-Length"]);
    }
    
} catch (error) {
    std.print("Error obteniendo headers:", error);
}

std.print("\n==================================================\n");

// ==========================================
// 7. COMPARACIÓN CON CLIENTE HTTP ANTERIOR
// ==========================================
std.print("7. Comparación con Cliente Anterior:");
std.print("Método anterior (clientHttpGet):");
try {
    let oldResponse = request.clientHttpGet("https://jsonplaceholder.typicode.com/posts/1");
    std.print("Respuesta anterior (string):", std.len(oldResponse), "caracteres");
    
    let parsedOld = json.parse(oldResponse);
    std.print("Título (método anterior):", parsedOld["title"]);
} catch (error) {
    std.print("Error con método anterior:", error);
}

std.print("\nMétodo nuevo (r2requests):");
try {
    let newResponse = request.get("https://jsonplaceholder.typicode.com/posts/1");
    std.print("Respuesta nueva (objeto):", newResponse["status_code"], newResponse["ok"]);
    std.print("Título (método nuevo):", newResponse["json"]["title"]);
    std.print("Más fácil de usar:", newResponse["ok"] ? "SÍ" : "NO");
} catch (error) {
    std.print("Error con método nuevo:", error);
}

std.print("\n==================================================\n");

// ==========================================
// 8. EJEMPLO PRÁCTICO: API GITHUB
// ==========================================
std.print("8. Ejemplo Práctico - API GitHub:");
try {
    let response = request.get("https://api.github.com/users/octocat");
    
    if (response["ok"]) {
        let user = response["json"];
        std.print("Usuario GitHub:", user["name"]);
        std.print("Bio:", user["bio"]);
        std.print("Seguidores:", user["followers"]);
        std.print("Repositorios públicos:", user["public_repos"]);
        std.print("Creado:", user["created_at"]);
    } else {
        std.print("No se pudo obtener información del usuario");
    }
} catch (error) {
    std.print("Error en ejemplo práctico:", error);
}

std.print("\n==================================================\n");

std.print("=== DEMO COMPLETADA ===");
std.print("La librería r2requests está funcionando correctamente!");
std.print("Ventajas sobre el cliente anterior:");
std.print("- API más intuitiva y familiar");
std.print("- Manejo automático de JSON");
std.print("- Mejor manejo de errores");
std.print("- Soporte para sesiones");
std.print("- Información detallada de respuestas");
std.print("- Soporte para todos los métodos HTTP");
std.print("- Utilidades integradas");