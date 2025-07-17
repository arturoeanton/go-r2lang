// example32-requests.r2
// Demostración de la librería r2requests (estilo Python Requests)
// Este ejemplo muestra cómo usar el cliente HTTP moderno de R2Lang

print("=== R2REQUESTS LIBRARY DEMO ===\n");

// ==========================================
// 1. GET REQUEST SIMPLE
// ==========================================
print("1. GET Request Simple:");
try {
    let response = get("https://jsonplaceholder.typicode.com/posts/1");
    print("Status Code:", response["status_code"]);
    print("OK:", response["ok"]);
    print("URL:", response["url"]);
    
    if (response["ok"]) {
        let post = response["json"];
        print("Post Title:", post["title"]);
        print("Post Body:", post["body"]);
        print("User ID:", post["userId"]);
    }
} catch (error) {
    print("Error en GET:", error);
}

print("\n==================================================\n");

// ==========================================
// 2. MÚLTIPLES MÉTODOS HTTP
// ==========================================
print("2. Múltiples Métodos HTTP:");
try {
    let getResponse = get("https://jsonplaceholder.typicode.com/posts/2");
    print("GET Status:", getResponse["status_code"]);
    
    let deleteResponse = delete("https://jsonplaceholder.typicode.com/posts/1");
    print("DELETE Status:", deleteResponse["status_code"]);
    
} catch (error) {
    print("Error en métodos HTTP:", error);
}

print("\n==================================================\n");

// ==========================================
// 3. USANDO SESIONES
// ==========================================
print("3. Usando Sesiones para Reutilización:");
try {
    let s = session();
    
    // Múltiples requests con la misma sesión
    let user1 = s["get"]("https://jsonplaceholder.typicode.com/users/1");
    let user2 = s["get"]("https://jsonplaceholder.typicode.com/users/2");
    
    if (user1["ok"] == true) {
        if (user2["ok"] == true) {
            print("Usuario 1:", user1["json"]["name"], "- Email:", user1["json"]["email"]);
            print("Usuario 2:", user2["json"]["name"], "- Email:", user2["json"]["email"]);
        }
    }
    
    // Cerrar sesión
    s["close"]();
} catch (error) {
    print("Error en sesiones:", error);
}

print("\n==================================================\n");

// ==========================================
// 4. MANEJO DE ERRORES HTTP
// ==========================================
print("4. Manejo de Errores HTTP:");
try {
    let response = get("https://jsonplaceholder.typicode.com/posts/999999");
    print("Status Code:", response["status_code"]);
    print("OK:", response["ok"]);
    
    if (response["ok"] == false) {
        print("Error HTTP detectado!");
        print("Texto de error:", response["text"]);
        
        // Intentar raise_for_status (esto causará un panic para errores)
        try {
            response["raise_for_status"]();
        } catch (statusError) {
            print("raise_for_status detectó error:", statusError);
        }
    }
} catch (error) {
    print("Error general:", error);
}

print("\n==================================================\n");

// ==========================================
// 5. UTILIDADES DE URL
// ==========================================
print("5. Utilidades de URL:");
try {
    let originalText = "Hello World! ¿Cómo estás?";
    let encoded = urlencode(originalText);
    let decoded = urldecode(encoded);
    
    print("Original:", originalText);
    print("Encoded:", encoded);
    print("Decoded:", decoded);
    print("Match:", originalText == decoded);
    
} catch (error) {
    print("Error en utilidades URL:", error);
}

print("\n==================================================\n");

// ==========================================
// 6. HEADERS Y INFORMACIÓN DE RESPUESTA
// ==========================================
print("6. Información Detallada de Respuesta:");
try {
    let response = get("https://jsonplaceholder.typicode.com/posts/1");
    
    print("URL final:", response["url"]);
    print("Status Code:", response["status_code"]);
    print("OK:", response["ok"]);
    print("Tiempo transcurrido:", response["elapsed"], "segundos");
    print("Tamaño del contenido:", len(response["text"]), "caracteres");
    
    // Headers más comunes
    print("\nHeaders importantes:");
    let headers = response["headers"];
    if (headers["Content-Type"]) {
        print("Content-Type:", headers["Content-Type"]);
    }
    if (headers["Content-Length"]) {
        print("Content-Length:", headers["Content-Length"]);
    }
    
} catch (error) {
    print("Error obteniendo headers:", error);
}

print("\n==================================================\n");

// ==========================================
// 7. COMPARACIÓN CON CLIENTE HTTP ANTERIOR
// ==========================================
print("7. Comparación con Cliente Anterior:");
print("Método anterior (clientHttpGet):");
try {
    let oldResponse = clientHttpGet("https://jsonplaceholder.typicode.com/posts/1");
    print("Respuesta anterior (string):", len(oldResponse), "caracteres");
    
    let parsedOld = parseJSON(oldResponse);
    print("Título (método anterior):", parsedOld["title"]);
} catch (error) {
    print("Error con método anterior:", error);
}

print("\nMétodo nuevo (r2requests):");
try {
    let newResponse = get("https://jsonplaceholder.typicode.com/posts/1");
    print("Respuesta nueva (objeto):", newResponse["status_code"], newResponse["ok"]);
    print("Título (método nuevo):", newResponse["json"]["title"]);
    print("Más fácil de usar:", newResponse["ok"] ? "SÍ" : "NO");
} catch (error) {
    print("Error con método nuevo:", error);
}

print("\n==================================================\n");

// ==========================================
// 8. EJEMPLO PRÁCTICO: API GITHUB
// ==========================================
print("8. Ejemplo Práctico - API GitHub:");
try {
    let response = get("https://api.github.com/users/octocat");
    
    if (response["ok"]) {
        let user = response["json"];
        print("Usuario GitHub:", user["name"]);
        print("Bio:", user["bio"]);
        print("Seguidores:", user["followers"]);
        print("Repositorios públicos:", user["public_repos"]);
        print("Creado:", user["created_at"]);
    } else {
        print("No se pudo obtener información del usuario");
    }
} catch (error) {
    print("Error en ejemplo práctico:", error);
}

print("\n==================================================\n");

print("=== DEMO COMPLETADA ===");
print("La librería r2requests está funcionando correctamente!");
print("Ventajas sobre el cliente anterior:");
print("- API más intuitiva y familiar");
print("- Manejo automático de JSON");
print("- Mejor manejo de errores");
print("- Soporte para sesiones");
print("- Información detallada de respuestas");
print("- Soporte para todos los métodos HTTP");
print("- Utilidades integradas");