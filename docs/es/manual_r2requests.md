# Manual Completo de r2requests - Cliente HTTP para R2Lang

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

r2requests es una librería de cliente HTTP para R2Lang, diseñada para ser intuitiva y potente, inspirada en la popular librería `requests` de Python. Facilita la realización de todo tipo de solicitudes HTTP (GET, POST, PUT, DELETE, etc.) con un manejo sencillo de JSON, formularios, archivos, sesiones y cookies.

### Características Principales

- **API Intuitiva**: Sintaxis clara y fácil de usar para operaciones HTTP comunes.
- **Manejo Automático de JSON**: Envío y recepción de datos JSON de forma transparente.
- **Sesiones Persistentes**: Reutilización de conexiones y manejo automático de cookies.
- **Subida de Archivos**: Soporte para envío de archivos en solicitudes multipart/form-data.
- **Proxies**: Configuración sencilla para solicitudes a través de proxies.
- **Reintentos Automáticos**: Mecanismos de reintento configurables para solicitudes fallidas.
- **Manejo de Errores Robusto**: Detección y gestión de errores HTTP.
- **Utilidades de URL**: Funciones para codificación y decodificación de URLs.

### Ventajas sobre Clientes HTTP Tradicionales

```javascript
// Tradicional (ejemplo hipotético)
let client = new HttpClient();
let response = client.request("GET", "https://api.example.com/data", {
    headers: {"Accept": "application/json"},
    timeout: 10000
});
if (response.statusCode == 200) {
    let data = parseJSON(response.body);
    print(data);
}

// r2requests
let response = get("https://api.example.com/data");
if (response.ok) {
    print(response.json);
}
```

---

## Instalación y Configuración

r2requests está integrado directamente en el entorno de R2Lang. No requiere instalación adicional.

### Prerequisitos

- R2Lang instalado y configurado.
- Conectividad de red para acceder a los endpoints HTTP.

### Verificación de Disponibilidad

```javascript
// Verificar que r2requests está disponible
print("Funciones r2requests disponibles:");
print("- get:", typeOf(get));
print("- post:", typeOf(post));
print("- session:", typeOf(session));
print("- urlencode:", typeOf(urlencode));
```

---

## Conceptos Fundamentales

### Solicitudes Globales vs. Sesiones

r2requests ofrece dos formas principales de realizar solicitudes:

1.  **Funciones Globales (`get`, `post`, `put`, `delete`, `patch`, `head`, `options`):**
    Son funciones de conveniencia para realizar solicitudes rápidas. Cada llamada a una función global utiliza un cliente HTTP temporal que comparte un "cookie jar" global, lo que permite el manejo automático de cookies entre solicitudes globales. Sin embargo, no persisten otras configuraciones como headers o autenticación.

    ```javascript
    // Solicitud GET simple
    let response = get("https://jsonplaceholder.typicode.com/posts/1");
    print(response.json.title);
    ```

2.  **Sesiones (`session()`):**
    Las sesiones son objetos que permiten persistir configuraciones (headers, autenticación, proxies, etc.) y reutilizar conexiones HTTP a lo largo de múltiples solicitudes. Son ideales para interactuar con APIs que requieren mantener un estado (como autenticación basada en cookies o tokens) o para optimizar el rendimiento al realizar muchas solicitudes al mismo host. Cada sesión tiene su propio "cookie jar" independiente.

    ```javascript
    let s = session();
    // Configurar headers para toda la sesión
    s.headers = {"Authorization": "Bearer your_token"};
    
    let user = s.get("https://api.example.com/users/1");
    let posts = s.get("https://api.example.com/users/1/posts");
    
    s.close(); // Cerrar la sesión para liberar recursos
    ```

### Objeto Response

Todas las funciones de solicitud (globales o de sesión) retornan un objeto `Response` que encapsula la respuesta del servidor. Este objeto tiene las siguientes propiedades y métodos:

-   `url` (string): La URL final de la solicitud.
-   `status_code` (number): El código de estado HTTP (ej. 200, 404, 500).
-   `ok` (boolean): `true` si el `status_code` es 2xx, `false` en caso contrario.
-   `headers` (object): Un objeto con los headers de la respuesta.
-   `text` (string): El cuerpo de la respuesta como una cadena de texto.
-   `json` (object/array): El cuerpo de la respuesta parseado como JSON (si el `Content-Type` es `application/json`).
-   `content` (array de bytes): El cuerpo de la respuesta como un array de bytes (útil para datos binarios).
-   `elapsed` (number): El tiempo transcurrido para la solicitud en segundos.
-   `json_func()`: Un método para obtener el cuerpo de la respuesta como JSON.
-   `raise_for_status()`: Un método que lanza un error si la solicitud no fue exitosa (código de estado 2xx).

```javascript
let response = get("https://jsonplaceholder.typicode.com/posts/1");
print("URL:", response.url);
print("Status:", response.status_code);
print("OK:", response.ok);
print("Title:", response.json.title);
print("Content-Type:", response.headers["Content-Type"]);
print("Elapsed:", response.elapsed, "seconds");

try {
    response.raise_for_status();
    print("Solicitud exitosa!");
} catch (error) {
    print("Error HTTP:", error);
}
```

---

## API Completa

### Funciones Globales

-   `get(url, [params])`
-   `post(url, [params])`
-   `put(url, [params])`
-   `delete(url, [params])`
-   `patch(url, [params])`
-   `head(url, [params])`
-   `options(url, [params])`

**Parámetros Comunes para Funciones Globales y de Sesión:**

El parámetro `params` es un objeto opcional que puede contener las siguientes claves:

-   `data` (string/object): Datos a enviar en el cuerpo de la solicitud. Si es un objeto, se enviará como `application/json` por defecto.
-   `json` (object): Datos a enviar como JSON. Establece automáticamente el `Content-Type` a `application/json`.
-   `files` (object): Un objeto donde las claves son los nombres de los campos de archivo y los valores son las rutas a los archivos locales. Se envía como `multipart/form-data`.
-   `headers` (object): Un objeto con headers HTTP personalizados para la solicitud.
-   `params` (object): Un objeto con parámetros de consulta URL.
-   `auth` (array): Un array `[username, password]` para autenticación Basic.
-   `timeout` (number): Tiempo máximo de espera para la solicitud en segundos.
-   `proxies` (object): Un objeto con URLs de proxy (ej. `{"http": "http://proxy.example.com:8080"}`).
-   `retries` (object): Configuración de reintentos (ej. `{"max": 3, "delay": 1.0}`).

### Sesiones

-   `session()`: Crea y retorna un nuevo objeto de sesión.

**Métodos del Objeto Session:**

Una vez creada una sesión, puedes usar los siguientes métodos, que son análogos a las funciones globales pero persisten la configuración de la sesión:

-   `s.get(url, [params])`
-   `s.post(url, [params])`
-   `s.put(url, [params])`
-   `s.delete(url, [params])`
-   `s.patch(url, [params])`
-   `s.head(url, [params])`
-   `s.options(url, [params])`
-   `s.close()`: Cierra la sesión y libera los recursos de conexión.

**Propiedades del Objeto Session:**

-   `s.headers` (object): Un objeto para establecer headers por defecto para todas las solicitudes de la sesión.
-   `s.auth` (array): Un array `[username, password]` para autenticación Basic por defecto para la sesión.
-   `s.timeout` (number): Timeout por defecto para todas las solicitudes de la sesión en segundos.
-   `s.proxies` (object): Proxies por defecto para la sesión.
-   `s.max_retries` (number): Número máximo de reintentos por defecto para la sesión.
-   `s.retry_delay` (number): Retraso entre reintentos por defecto en segundos.

### Utilidades de URL

-   `urlencode(string)`: Codifica una cadena de texto para ser usada en una URL.
-   `urldecode(string)`: Decodifica una cadena de texto de una URL.

---

## Ejemplos Básicos

### Ejemplo 1: Solicitud GET Simple

```javascript
let response = get("https://jsonplaceholder.typicode.com/todos/1");

if (response.ok) {
    print("Todo Title:", response.json.title);
} else {
    print("Error:", response.status_code);
}
```

### Ejemplo 2: Solicitud POST con JSON

```javascript
let postData = {
    "title": "foo",
    "body": "bar",
    "userId": 1
};

let response = post("https://jsonplaceholder.typicode.com/posts", {
    "json": postData
});

if (response.ok) {
    print("Nuevo Post ID:", response.json.id);
} else {
    print("Error:", response.status_code);
}
```

### Ejemplo 3: Uso de Sesiones para Múltiples Solicitudes

```javascript
let s = session();

// Login (asumiendo que establece una cookie de sesión)
let loginResponse = s.post("https://api.example.com/login", {
    "json": {"username": "user", "password": "password"}
});

if (loginResponse.ok) {
    print("Login exitoso!");
    
    // Acceder a un recurso protegido (la cookie se envía automáticamente)
    let profileResponse = s.get("https://api.example.com/profile");
    if (profileResponse.ok) {
        print("Perfil de usuario:", profileResponse.json.name);
    } else {
        print("Error al obtener perfil:", profileResponse.status_code);
    }
} else {
    print("Error de login:", loginResponse.status_code);
}

s.close();
```

### Ejemplo 4: Subida de Archivos

```javascript
// Crear un archivo temporal para el ejemplo
let tempFilePath = "temp_upload.txt";
writeFile(tempFilePath, "Contenido de prueba para el archivo.");

let response = post("https://httpbin.org/post", {
    "files": {
        "my_file": tempFilePath
    },
    "data": {
        "description": "Un archivo de texto de ejemplo"
    }
});

if (response.ok) {
    print("Archivo subido exitosamente!");
    print("Nombre del archivo:", response.json.files.my_file.name);
    print("Descripción:", response.json.form.description);
} else {
    print("Error al subir archivo:", response.status_code);
}

deleteFile(tempFilePath); // Limpiar archivo temporal
```

---

## Configuración Empresarial

### Headers Corporativos

Puedes establecer headers personalizados para todas las solicitudes de una sesión o para solicitudes individuales.

```javascript
let s = session();
s.headers = {
    "X-Company-ID": "CORP-XYZ",
    "X-Request-Source": "R2Lang-App",
    "User-Agent": "MiAplicacion/1.0 (R2Lang)"
};

let response = s.get("https://api.internal.company.com/data");
// Los headers X-Company-ID, X-Request-Source y User-Agent se enviarán automáticamente
```

### Configuración de Timeouts para Servicios Críticos

Es crucial establecer timeouts para evitar que las solicitudes se queden colgadas indefinidamente.

```javascript
// Timeout global para una sesión
let s = session();
s.timeout = 10.0; // 10 segundos para todas las solicitudes de esta sesión

// Timeout para una solicitud específica (sobrescribe el de la sesión)
let response = s.get("https://api.slowservice.com/data", {
    "timeout": 30.0 // 30 segundos solo para esta solicitud
});
```

### Uso de Proxies en Entornos Corporativos

```javascript
let s = session();
s.proxies = {
    "http": "http://proxy.corp.com:8080",
    "https": "http://proxy.corp.com:8080" // O un proxy HTTPS diferente
};

let response = s.get("https://api.external.com/data"); // La solicitud pasará por el proxy
```

### Reintentos Automáticos para APIs Inestables

```javascript
let s = session();
s.max_retries = 5;
s.retry_delay = 2.0; // 2 segundos de retraso entre reintentos

// La solicitud se reintentará hasta 5 veces con 2 segundos de retraso si falla
let response = s.get("https://api.flakyservice.com/data");
```

---

## Autenticación

r2requests simplifica el manejo de diferentes tipos de autenticación.

### Basic Authentication

```javascript
// Autenticación Basic para una solicitud global
let response = get("https://api.example.com/secure", {
    "auth": ["username", "password"]
});

// Autenticación Basic para toda una sesión
let s = session();
s.auth = ["api_user", "api_key"];
let protectedData = s.get("https://api.example.com/protected");
```

### Autenticación Basada en Headers (Bearer Tokens, API Keys)

Para tokens de portador (Bearer tokens) o claves de API, simplemente configura el header `Authorization` o el header personalizado.

```javascript
// Usando un Bearer Token
let token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...";
let response = get("https://api.example.com/graphql", {
    "headers": {"Authorization": "Bearer " + token}
});

// Usando una API Key
let apiKey = "your_super_secret_api_key";
let s = session();
s.headers = {"X-API-Key": apiKey};
let userData = s.get("https://api.example.com/users");
```

### Autenticación Basada en Cookies (Manejo Automático)

Como se demostró en el ejemplo de sesiones, r2requests maneja automáticamente las cookies. Si un servidor establece una cookie de sesión después de un login, esta se enviará automáticamente en las solicitudes subsiguientes de la misma sesión (o en solicitudes globales si se usa el cookie jar global).

```javascript
let s = session();
// Suponiendo que /login establece una cookie de sesión
s.post("https://api.example.com/login", {"json": {"user": "test", "pass": "123"}});

// La cookie de sesión se enviará automáticamente aquí
let dashboard = s.get("https://api.example.com/dashboard");
print(dashboard.text);
```

---

## SSL/TLS y Seguridad

r2requests utiliza la configuración SSL/TLS subyacente de Go, que es segura por defecto.

### Verificación de Certificados (Por Defecto)

Por defecto, r2requests verifica los certificados SSL/TLS del servidor para asegurar una conexión segura. Si la verificación falla (ej. certificado autofirmado, caducado, o no confiable), la solicitud fallará.

### Deshabilitar Verificación (Solo para Desarrollo/Testing)

**Advertencia:** Deshabilitar la verificación de certificados (`verify: false`) es una práctica **insegura** y solo debe hacerse en entornos de desarrollo o pruebas controlados. **Nunca** en producción.

Actualmente, r2requests no expone una opción directa para deshabilitar la verificación de certificados a nivel de la API de R2Lang. Esto se hace para fomentar prácticas seguras. Si necesitas interactuar con un servidor con un certificado no válido en desarrollo, considera usar un proxy que maneje la validación o añadir el certificado a tu almacén de confianza del sistema.

---

## Manejo de Errores

r2requests proporciona varias formas de manejar errores HTTP y de red.

### Códigos de Estado HTTP

Siempre puedes verificar el `status_code` y la propiedad `ok` del objeto `Response`.

```javascript
let response = get("https://api.example.com/nonexistent");

if (!response.ok) {
    print("Error HTTP:", response.status_code);
    if (response.status_code == 404) {
        print("Recurso no encontrado.");
    } else if (response.status_code == 500) {
        print("Error interno del servidor.");
    }
}
```

### `raise_for_status()`

Este método es útil para lanzar un error si la solicitud no fue exitosa (código de estado 2xx). Esto permite un manejo de errores más conciso.

```javascript
try {
    let response = get("https://api.example.com/data");
    response.raise_for_status(); // Lanza un error si status_code no es 2xx
    print("Datos:", response.json);
} catch (error) {
    print("La solicitud falló:", error);
}
```

### Errores de Conexión y Timeout

Los errores de red (como problemas de DNS, conexión rechazada, o timeouts) resultarán en un `panic` en R2Lang. Debes envolver tus llamadas a `r2requests` en bloques `try-catch` para manejarlos.

```javascript
try {
    // Intentar conectar a un host inexistente o con timeout muy bajo
    let response = get("http://nonexistent-domain.invalid", {"timeout": 0.1}); 
    print(response.text);
} catch (error) {
    print("Error de conexión o timeout:", error);
    let errorStr = "" + error;
    if (indexOf(errorStr, "timeout") != -1) {
        print("La solicitud excedió el tiempo límite.");
    } else if (indexOf(errorStr, "no such host") != -1) {
        print("El dominio no pudo ser resuelto.");
    }
}
```

---

## Mejores Prácticas

### 1. Reutilización de Sesiones

Siempre que realices múltiples solicitudes al mismo host o necesites mantener el estado (cookies, headers), utiliza una sesión. Esto mejora el rendimiento y simplifica el código.

```javascript
// MAL: Múltiples conexiones y sin estado
get("https://api.example.com/users");
get("https://api.example.com/products");

// BIEN: Reutiliza conexión y mantiene estado
let s = session();
s.get("https://api.example.com/users");
s.get("https://api.example.com/products");
s.close();
```

### 2. Manejo de Errores Explícito

Aunque `raise_for_status()` es conveniente, para lógica de negocio compleja, es mejor verificar explícitamente el `status_code` y manejar diferentes escenarios de error.

### 3. Logging y Observabilidad

Integra r2requests con tu sistema de logging para registrar solicitudes, respuestas, errores y tiempos de respuesta.

```javascript
func makeAuditedRequest(method, url, params) {
    print("LOG: Iniciando solicitud", method, url);
    let startTime = now();
    try {
        let response = callRequestFunction(method, url, params); // Función auxiliar para llamar get/post/etc.
        let duration = now() - startTime;
        print("LOG: Solicitud", url, "completada en", duration, "ms. Status:", response.status_code);
        if (!response.ok) {
            print("LOG: ERROR - Solicitud fallida:", response.status_code, response.text);
        }
        return response;
    } catch (error) {
        let duration = now() - startTime;
        print("LOG: ERROR - Solicitud", url, "falló en", duration, "ms:", error);
        throw error;
    }
}

// Uso
let response = makeAuditedRequest("GET", "https://api.example.com/data", null);
```

### 4. Configuración Centralizada

Para aplicaciones empresariales, centraliza la configuración de r2requests (timeouts, proxies, autenticación) para facilitar la gestión y el cambio entre entornos.

```javascript
func createAPIClient(baseURL, config) {
    let s = session();
    s.base_url = baseURL; // Propiedad hipotética para base URL
    if (config.timeout) { s.timeout = config.timeout; }
    if (config.auth) { s.auth = config.auth; }
    if (config.proxies) { s.proxies = config.proxies; }
    if (config.headers) { s.headers = config.headers; }
    return s;
}

let prodClient = createAPIClient("https://api.prod.com", {
    "timeout": 30.0,
    "auth": [getEnv("API_USER"), getEnv("API_PASS")],
    "headers": {"X-Prod-Key": "prod-key"}
});

let devClient = createAPIClient("https://api.dev.com", {
    "timeout": 10.0,
    "headers": {"X-Dev-Key": "dev-key"}
});

let prodData = prodClient.get("/users");
let devData = devClient.get("/users");
```

---

## Casos de Uso Avanzados

### 1. Integración con APIs RESTful Externas

```javascript
func fetchExchangeRates(currency) {
    let s = session();
    s.headers = {"Accept": "application/json"};
    s.timeout = 5.0;

    let response = s.get("https://api.exchangerates.io/latest", {
        "params": {"base": currency}
    });

    if (response.ok) {
        return response.json.rates;
    } else {
        print("Error al obtener tasas de cambio:", response.status_code, response.text);
        return null;
    }
}

let rates = fetchExchangeRates("USD");
if (rates) {
    print("USD a EUR:", rates.EUR);
}
```

### 2. Web Scraping Básico (con manejo de cookies y user-agent)

```javascript
let s = session();
s.headers = {
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
    "Accept-Language": "en-US,en;q=0.9"
};

// Primera solicitud para obtener cookies de sesión
let loginPage = s.get("https://example.com/login");
print("Obtenida página de login. Status:", loginPage.status_code);

// Simular login (asumiendo que el servidor establece más cookies)
let loginResponse = s.post("https://example.com/do_login", {
    "data": {
        "username": "myuser",
        "password": "mypassword",
        "csrf_token": extractCSRFToken(loginPage.text) // Función auxiliar para extraer token
    }
});

if (loginResponse.ok) {
    print("Login exitoso. Accediendo a datos protegidos...");
    let protectedPage = s.get("https://example.com/dashboard");
    print("Contenido del dashboard (primeras 200 chars):", substring(protectedPage.text, 0, 200));
} else {
    print("Fallo el login. Status:", loginResponse.status_code);
}
s.close();

func extractCSRFToken(htmlContent) {
    // Implementar lógica para parsear HTML y extraer el token
    // Esto es un placeholder
    return "dummy_csrf_token"; 
}
```

### 3. Consumo de APIs con Paginación

```javascript
func fetchAllUsers(apiBaseUrl, apiKey) {
    let s = session();
    s.headers = {"Authorization": "Bearer " + apiKey};
    let allUsers = [];
    let page = 1;
    let hasMore = true;

    while (hasMore) {
        let response = s.get(apiBaseUrl + "/users", {
            "params": {"page": page, "limit": 100}
        });

        if (response.ok) {
            let users = response.json.data;
            allUsers = allUsers + users; // Concatenar arrays
            
            // Asumiendo que la API indica si hay más páginas
            hasMore = response.json.has_more_pages; 
            page = page + 1;
        } else {
            print("Error fetching users:", response.status_code, response.text);
            hasMore = false;
        }
    }
    s.close();
    return allUsers;
}

let users = fetchAllUsers("https://api.example.com/v1", "your_api_key");
print("Total de usuarios obtenidos:", len(users));
```

---

## Troubleshooting

### Problemas Comunes y Soluciones

#### 1. `panic: Request failed after X retries: context deadline exceeded` o `panic: Failed to create request: Get "URL": dial tcp IP:port: connect: connection refused`

**Síntomas:** La solicitud se cuelga y luego falla con un error de timeout o conexión rechazada.

**Causas:**
- El servidor no responde o está caído.
- Problemas de red (firewall, DNS, conectividad).
- Timeout de solicitud demasiado bajo.
- El servidor bloquea las solicitudes (ej. por User-Agent).

**Soluciones:**
- **Aumentar el `timeout`:**
    ```javascript
    let response = get("http://slow-api.com/data", {"timeout": 30.0}); // 30 segundos
    ```
- **Verificar conectividad:** Asegúrate de que puedes acceder a la URL desde tu entorno.
- **Cambiar `User-Agent`:** Algunos servidores bloquean User-Agents por defecto.
    ```javascript
    let response = get("http://example.com", {"headers": {"User-Agent": "Mozilla/5.0"}});
    ```
- **Configurar `proxies`:** Si estás en una red corporativa, podrías necesitar un proxy.

#### 2. `panic: Failed to encode JSON: json: unsupported type`

**Síntomas:** Error al intentar enviar datos JSON.

**Causas:** Estás pasando un tipo de dato que no puede ser serializado a JSON por R2Lang (ej. una función, un objeto circular).

**Soluciones:** Asegúrate de que los datos que pasas a `json` o `data` sean objetos o arrays simples, strings, números o booleanos.

#### 3. `status_code` 401 (Unauthorized) o 403 (Forbidden)

**Síntomas:** La solicitud retorna un código de estado de autenticación/autorización.

**Causas:**
- Credenciales incorrectas o faltantes.
- Token de autenticación caducado o inválido.
- No tienes permisos para acceder al recurso.

**Soluciones:**
- **Verificar credenciales:** Asegúrate de que `username`, `password`, `token` o `API Key` sean correctos.
- **Refrescar token:** Si usas OAuth, implementa la lógica para refrescar el token.
- **Contactar al administrador de la API:** Si las credenciales son correctas, podrías no tener los permisos necesarios.

### Herramientas de Debugging

-   **Imprimir `response.status_code` y `response.text`:** Siempre verifica el código de estado y el cuerpo de la respuesta para entender el problema.
-   **Usar `httpbin.org`:** Un servicio útil para probar solicitudes HTTP y ver cómo el servidor las recibe.
    ```javascript
    let response = post("https://httpbin.org/post", {"json": {"test": "data"}});
    print(response.json); // Verás los datos que httpbin recibió
    ```
-   **Logging detallado:** Implementa un logging más granular como se sugiere en las "Mejores Prácticas".

---

## Referencia Técnica

### Objeto `Response`

```javascript
{
    "url": "string",
    "status_code": "number",
    "ok": "boolean",
    "headers": {
        "Header-Name": "string" // o array de strings si hay múltiples
    },
    "text": "string",
    "json": "object | array | null", // Depende del Content-Type y si es JSON válido
    "content": "array de bytes",
    "elapsed": "number", // Tiempo en segundos
    "json_func": "function", // Método para obtener JSON
    "raise_for_status": "function" // Método para lanzar error en caso de fallo HTTP
}
```

### Parámetros de Solicitud

```javascript
{
    "data": "string | object", // Cuerpo de la solicitud (form-urlencoded o JSON si es objeto)
    "json": "object",          // Cuerpo de la solicitud como JSON
    "files": {                 // Para multipart/form-data
        "fieldName": "filePath"
    },
    "headers": {               // Headers HTTP personalizados
        "Header-Name": "string"
    },
    "params": {                // Parámetros de consulta URL
        "paramName": "value"
    },
    "auth": ["username", "password"], // Basic Authentication
    "timeout": "number",       // Timeout en segundos
    "proxies": {               // Configuración de proxy
        "http": "http://proxy.example.com:8080",
        "https": "http://secure-proxy.example.com:8080"
    },
    "retries": {               // Configuración de reintentos
        "max": "number",       // Número máximo de reintentos
        "delay": "number"      // Retraso entre reintentos en segundos
    }
}
```

### Propiedades del Objeto `Session`

```javascript
{
    "headers": "object",       // Headers por defecto para la sesión
    "auth": ["username", "password"], // Autenticación Basic por defecto
    "timeout": "number",       // Timeout por defecto en segundos
    "proxies": "object",       // Proxies por defecto
    "max_retries": "number",   // Número máximo de reintentos por defecto
    "retry_delay": "number",   // Retraso entre reintentos por defecto en segundos
    // Métodos: get, post, put, delete, patch, head, options, close
}
```

### Códigos de Error HTTP Comunes

-   **200 OK**: La solicitud fue exitosa.
-   **201 Created**: La solicitud ha tenido éxito y se ha creado un nuevo recurso.
-   **204 No Content**: La solicitud fue exitosa, pero no hay contenido para enviar en la respuesta.
-   **400 Bad Request**: La solicitud no pudo ser entendida por el servidor debido a una sintaxis inválida.
-   **401 Unauthorized**: La solicitud requiere autenticación.
-   **403 Forbidden**: El servidor entendió la solicitud, pero se niega a autorizarla.
-   **404 Not Found**: El servidor no pudo encontrar el recurso solicitado.
-   **405 Method Not Allowed**: El método de solicitud no es compatible con el recurso.
-   **408 Request Timeout**: El servidor no recibió una respuesta completa a tiempo.
-   **409 Conflict**: La solicitud no pudo ser completada debido a un conflicto con el estado actual del recurso.
-   **429 Too Many Requests**: El usuario ha enviado demasiadas solicitudes en un período de tiempo determinado.
-   **500 Internal Server Error**: El servidor encontró una condición inesperada que le impidió cumplir con la solicitud.
-   **502 Bad Gateway**: El servidor, mientras actuaba como puerta de enlace o proxy, recibió una respuesta inválida de un servidor ascendente.
-   **503 Service Unavailable**: El servidor no está listo para manejar la solicitud.
-   **504 Gateway Timeout**: El servidor, mientras actuaba como puerta de enlace o proxy, no recibió una respuesta a tiempo de un servidor ascendente.

---

## Ejemplos de Integración Empresarial

### 1. Integración con un Sistema de Gestión de Clientes (CRM)

```javascript
func createCRMContact(contactData) {
    let s = session();
    s.base_url = "https://api.crm.company.com/v1"; // Propiedad hipotética
    s.headers = {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + getEnv("CRM_API_TOKEN")
    };
    s.timeout = 15.0;
    s.max_retries = 3;
    s.retry_delay = 1.0;

    try {
        let response = s.post("/contacts", {"json": contactData});
        response.raise_for_status();
        print("Contacto CRM creado exitosamente:", response.json.id);
        return response.json;
    } catch (error) {
        print("Error al crear contacto CRM:", error);
        return null;
    } finally {
        s.close();
    }
}

let newContact = {
    "first_name": "Juan",
    "last_name": "Perez",
    "email": "juan.perez@example.com",
    "phone": "123-456-7890"
};
createCRMContact(newContact);
```

### 2. Automatización de Tareas en un Sistema de Tickets (ITSM)

```javascript
func updateTicketStatus(ticketId, newStatus, comments) {
    let s = session();
    s.base_url = "https://itsm.company.com/api"; // Propiedad hipotética
    s.auth = [getEnv("ITSM_USER"), getEnv("ITSM_PASS")];
    s.headers = {"Content-Type": "application/json"};

    try {
        let response = s.put("/tickets/" + ticketId, {
            "json": {
                "status": newStatus,
                "comments": comments
            }
        });
        response.raise_for_status();
        print("Ticket", ticketId, "actualizado a estado:", newStatus);
        return true;
    } catch (error) {
        print("Error al actualizar ticket", ticketId, ":", error);
        return false;
    } finally {
        s.close();
    }
}

updateTicketStatus("TICKET-001", "Closed", "Problema resuelto y verificado.");
```

### 3. Sincronización de Datos con un Almacén de Objetos (S3-like API)

```javascript
func uploadFileToObjectStorage(bucketName, objectKey, filePath) {
    let s = session();
    s.base_url = "https://s3.company.com"; // Propiedad hipotética
    s.headers = {
        "Authorization": "Bearer " + getEnv("S3_ACCESS_TOKEN"),
        "x-amz-acl": "private" // Ejemplo de header específico de S3
    };

    try {
        let response = s.put("/" + bucketName + "/" + objectKey, {
            "files": {"file": filePath} // PUT con cuerpo de archivo
        });
        response.raise_for_status();
        print("Archivo", objectKey, "subido exitosamente al bucket", bucketName);
        return true;
    } catch (error) {
        print("Error al subir archivo a S3:", error);
        return false;
    } finally {
        s.close();
    }
}

// Crear un archivo de ejemplo
let exampleFilePath = "reporte_ventas.pdf";
writeFile(exampleFilePath, "Contenido del reporte PDF...");

uploadFileToObjectStorage("reports-bucket", "2024/q2/ventas.pdf", exampleFilePath);
deleteFile(exampleFilePath);
```

---

## Conclusión

r2requests proporciona una interfaz moderna y potente para interactuar con servicios HTTP en R2Lang. Su diseño intuitivo, combinado con características avanzadas como sesiones, manejo automático de cookies, reintentos y soporte para archivos, lo convierte en una herramienta indispensable para integraciones y automatizaciones empresariales.

### Características Destacadas

-   ✅ **API Sencilla**: Facilita la realización de solicitudes HTTP complejas.
-   ✅ **Sesiones Robustas**: Manejo de estado y reutilización de conexiones.
-   ✅ **Cookies Automáticas**: Gestión transparente de cookies para sesiones y funciones globales.
-   ✅ **Soporte Completo**: GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS.
-   ✅ **Manejo de Datos Flexible**: JSON, formularios, subida de archivos.
-   ✅ **Fiabilidad**: Reintentos automáticos y manejo de timeouts.
-   ✅ **Seguridad**: Verificación de certificados SSL/TLS por defecto.

### Roadmap Futuro

-   Soporte para autenticación OAuth 2.0 y OIDC.
-   Manejo de eventos de progreso para subidas/descargas grandes.
-   Integración con sistemas de métricas y tracing.
-   Soporte para HTTP/2 y HTTP/3.
-   Configuración de certificados de cliente TLS.

---

**Versión del Manual:** 1.0  
**Última Actualización:** 2025-07-17  
**Compatibilidad:** R2Lang v2.0+
