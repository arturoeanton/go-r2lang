# Propuesta: Librería JWT y JWC para R2Lang - PROPUESTA TÉCNICA COMPLETA

## Resumen Ejecutivo

Se propone el desarrollo de una librería completa para manejo de **JSON Web Tokens (JWT)** y **JSON Web Certificates (JWC)** en R2Lang, siguiendo el mismo patrón de éxito implementado en r2soap. Esta librería proporcionará capacidades empresariales completas para autenticación, autorización, y seguridad basada en tokens.

**🎯 Objetivo**: Crear r2jwt como una librería empresarial robusta para:
- Generación y validación de JWT tokens
- Manejo de JWC (JSON Web Certificates) 
- Integración con r2soap y r2grpc
- Soporte para múltiples algoritmos de firma
- Configuración empresarial y compliance

## Contexto y Motivación

### Necesidad Empresarial

En el ecosistema moderno de microservicios y APIs, JWT es el estándar de facto para:
- **Autenticación stateless** entre servicios
- **Autorización basada en claims** 
- **Single Sign-On (SSO)** empresarial
- **API security** con scopes y roles
- **Federación de identidades** entre organizaciones

### Integración con r2soap y r2grpc

```javascript
// Integración perfecta con clientes existentes
let jwtToken = jwt.create({"user": "admin", "role": "manager"}, secret);

// Usar con r2soap
let soapClient = soapClient("https://secure.company.com/service.wsdl");
soapClient.setAuth({"type": "bearer", "token": jwtToken});

// Usar con r2grpc (futuro)
let grpcClient = grpcClient("service.proto", "secure.company.com:443");
grpcClient.setAuth({"type": "bearer", "token": jwtToken});
```

## Arquitectura y Diseño

### Estructura Principal

```go
// pkg/r2libs/r2jwt.go
type JWTManager struct {
    DefaultAlgorithm string
    DefaultExpiry    time.Duration
    Secrets         map[string]string  // Para múltiples secrets
    PublicKeys      map[string]string  // Para validación RS256
    PrivateKeys     map[string]string  // Para firma RS256
}

type JWTClaim struct {
    Issuer    string                 `json:"iss"`
    Subject   string                 `json:"sub"`
    Audience  []string               `json:"aud"`
    ExpiresAt time.Time              `json:"exp"`
    NotBefore time.Time              `json:"nbf"`
    IssuedAt  time.Time              `json:"iat"`
    JwtID     string                 `json:"jti"`
    Custom    map[string]interface{} `json:"-"`
}
```

### API Completa para R2Lang

#### 1. Creación de Tokens

```javascript
// Creación básica
let token = jwt.create({"user": "john", "role": "admin"}, "secret");

// Creación con configuración completa
let tokenConfig = {
    "payload": {
        "user": "john.doe",
        "role": "admin", 
        "department": "finance",
        "permissions": ["read", "write", "delete"]
    },
    "secret": "enterprise-secret-key",
    "algorithm": "HS256",
    "expiresIn": "1h",
    "issuer": "company.com",
    "audience": ["api.company.com", "admin.company.com"]
};
let enterpriseToken = jwt.createToken(tokenConfig);

// Creación con RS256 (claves asimétricas)
let rsaToken = jwt.createTokenRS256({
    "payload": {"user": "john", "role": "admin"},
    "privateKey": loadPrivateKey("/path/to/private.pem"),
    "expiresIn": "24h"
});
```

#### 2. Validación y Verificación

```javascript
// Validación básica
let isValid = jwt.verify(token, "secret");

// Validación con decodificación
let decoded = jwt.verifyAndDecode(token, "secret");
if (decoded.valid) {
    print("Usuario:", decoded.payload.user);
    print("Rol:", decoded.payload.role);
    print("Expira:", decoded.payload.exp);
}

// Validación con múltiples secrets (key rotation)
let multiSecretResult = jwt.verifyMultiple(token, {
    "current": "current-secret",
    "previous": "previous-secret",
    "legacy": "legacy-secret"
});

// Validación RS256
let rsaResult = jwt.verifyRS256(token, publicKey);
```

#### 3. Manipulación de Claims

```javascript
// Extraer claims sin validar
let claims = jwt.decode(token);
print("Issuer:", claims.iss);
print("Subject:", claims.sub);
print("Custom data:", claims.custom);

// Verificar expiración
let isExpired = jwt.isExpired(token);
let timeToExpiry = jwt.timeToExpiry(token);

// Verificar audience
let isValidAudience = jwt.checkAudience(token, "api.company.com");

// Refresh token
let refreshed = jwt.refresh(token, "secret", "24h");
```

#### 4. JWT Enterprise Features

```javascript
// Crear JWT manager empresarial
let jwtManager = jwt.createManager({
    "defaultAlgorithm": "RS256",
    "defaultExpiry": "8h",
    "issuer": "company.com",
    "keyRotation": true
});

// Agregar múltiples secrets para key rotation
jwtManager.addSecret("v1", "secret-version-1");
jwtManager.addSecret("v2", "secret-version-2");
jwtManager.setCurrentSecret("v2");

// Agregar claves RSA
jwtManager.addRSAKeyPair("prod", {
    "private": loadFile("/secure/keys/prod-private.pem"),
    "public": loadFile("/secure/keys/prod-public.pem")
});

// Crear token con manager
let managedToken = jwtManager.createToken({
    "user": "admin",
    "role": "super-admin",
    "scopes": ["users:read", "users:write", "system:admin"]
});
```

#### 5. JWC (JSON Web Certificates) Support

```javascript
// Crear JWC (certificado en formato JWT)
let jwcConfig = {
    "subject": "api.company.com",
    "issuer": "ca.company.com", 
    "publicKey": publicKeyPEM,
    "validFrom": "2024-01-01T00:00:00Z",
    "validTo": "2025-01-01T00:00:00Z",
    "keyUsage": ["digitalSignature", "keyEncipherment"],
    "san": ["api.company.com", "*.api.company.com"]
};
let jwc = jwt.createJWC(jwcConfig, caPrivateKey);

// Validar JWC
let jwcValid = jwt.verifyJWC(jwc, caPublicKey);

// Extraer certificado de JWC
let certInfo = jwt.extractCertificate(jwc);
print("Subject:", certInfo.subject);
print("Valid until:", certInfo.validTo);
```

### Integración con Sistemas de Autenticación

#### 1. OAuth 2.0 / OpenID Connect

```javascript
// Crear token compatible con OAuth 2.0
let oauthToken = jwt.createOAuth({
    "clientId": "app-client-123",
    "userId": "user-456", 
    "scopes": ["profile", "email", "api:read"],
    "grantType": "authorization_code",
    "expiresIn": "3600"
});

// Validar token OAuth
let oauthValid = jwt.verifyOAuth(oauthToken, {
    "checkScopes": ["api:read"],
    "checkClientId": "app-client-123"
});
```

#### 2. SAML to JWT Bridge

```javascript
// Convertir SAML assertion a JWT
let samlAssertion = loadSAMLAssertion();
let jwtFromSAML = jwt.fromSAML(samlAssertion, {
    "mapClaims": {
        "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/name": "user",
        "http://schemas.microsoft.com/ws/2008/06/identity/claims/role": "role"
    }
});
```

#### 3. LDAP Integration

```javascript
// Crear JWT desde autenticación LDAP
func authenticateAndCreateJWT(username, password) {
    // Validar credenciales con LDAP (usando r2ldap future)
    let ldapUser = ldap.authenticate(username, password);
    
    if (ldapUser.valid) {
        return jwt.create({
            "user": ldapUser.username,
            "email": ldapUser.email,
            "groups": ldapUser.groups,
            "department": ldapUser.department
        }, getJWTSecret());
    }
    return null;
}
```

## Casos de Uso Empresariales

### 1. API Gateway Authentication

```javascript
// Middleware de autenticación para APIs
func validateAPIRequest(request) {
    let authHeader = request.headers["Authorization"];
    if (!authHeader || !startsWith(authHeader, "Bearer ")) {
        return {"valid": false, "error": "Missing or invalid authorization header"};
    }
    
    let token = substring(authHeader, 7); // Remove "Bearer "
    let validation = jwt.verifyAndDecode(token, getAPISecret());
    
    if (!validation.valid) {
        return {"valid": false, "error": "Invalid token"};
    }
    
    // Verificar scopes requeridos
    let requiredScope = getRequiredScope(request.path, request.method);
    if (!hasScope(validation.payload.scopes, requiredScope)) {
        return {"valid": false, "error": "Insufficient permissions"};
    }
    
    return {
        "valid": true,
        "user": validation.payload.user,
        "role": validation.payload.role,
        "scopes": validation.payload.scopes
    };
}
```

### 2. Microservices Authentication

```javascript
// Servicio A llama a Servicio B con JWT
func callMicroservice(serviceURL, endpoint, data) {
    // Crear service-to-service token
    let serviceToken = jwt.create({
        "service": "service-a",
        "target": "service-b", 
        "iat": getCurrentTime(),
        "exp": getCurrentTime() + 300 // 5 minutos
    }, getServiceSecret());
    
    // Usar con r2soap o HTTP client
    let client = httpClient();
    client.setHeader("Authorization", "Bearer " + serviceToken);
    
    return client.post(serviceURL + endpoint, data);
}
```

### 3. Session Management

```javascript
// Gestión de sesiones con JWT
class JWTSessionManager {
    func createSession(userId, userData) {
        let sessionData = {
            "userId": userId,
            "userData": userData,
            "sessionId": generateUUID(),
            "createdAt": getCurrentTime()
        };
        
        return jwt.create(sessionData, getSessionSecret(), "24h");
    }
    
    func validateSession(sessionToken) {
        let validation = jwt.verifyAndDecode(sessionToken, getSessionSecret());
        
        if (!validation.valid) {
            return {"valid": false, "reason": "Invalid token"};
        }
        
        // Verificar si la sesión debe ser renovada
        let timeToExpiry = jwt.timeToExpiry(sessionToken);
        if (timeToExpiry < 3600) { // Menos de 1 hora
            let refreshedToken = this.refreshSession(sessionToken);
            return {
                "valid": true,
                "user": validation.payload.userData,
                "refreshedToken": refreshedToken
            };
        }
        
        return {
            "valid": true,
            "user": validation.payload.userData
        };
    }
    
    func refreshSession(oldToken) {
        let decoded = jwt.decode(oldToken);
        return this.createSession(decoded.userId, decoded.userData);
    }
}
```

### 4. Role-Based Access Control (RBAC)

```javascript
// Sistema RBAC con JWT
func createRoleToken(user, roles, permissions) {
    let roleData = {
        "user": user.id,
        "username": user.username,
        "roles": roles,
        "permissions": permissions,
        "department": user.department,
        "level": calculateAccessLevel(roles)
    };
    
    return jwt.create(roleData, getRoleSecret(), "8h");
}

func checkPermission(token, requiredPermission) {
    let validation = jwt.verifyAndDecode(token, getRoleSecret());
    
    if (!validation.valid) {
        return false;
    }
    
    let userPermissions = validation.payload.permissions;
    return contains(userPermissions, requiredPermission) || 
           contains(userPermissions, "admin:all");
}

// Decorador para funciones que requieren permisos
func requiresPermission(permission, func_to_execute) {
    return func(token, ...args) {
        if (!checkPermission(token, permission)) {
            throw "Access denied: requires " + permission;
        }
        return func_to_execute(...args);
    };
}

// Uso
let protectedFunction = requiresPermission("users:write", createUser);
protectedFunction(userToken, userData); // Solo ejecuta si tiene permisos
```

## Características de Seguridad

### 1. Key Rotation Support

```javascript
// Gestor de rotación de claves
let keyManager = jwt.createKeyManager({
    "rotationInterval": "30d",
    "keepOldKeys": 3,
    "notifyBeforeExpiry": "7d"
});

// Agregar nueva clave
keyManager.rotateKey("new-secret-v3");

// Validar con múltiples claves (permite transición suave)
let result = keyManager.verifyWithAnyKey(token);
```

### 2. Blacklist/Revocation

```javascript
// Lista negra de tokens
let tokenBlacklist = jwt.createBlacklist();

// Revocar token
tokenBlacklist.revoke(token, "user-requested");

// Verificar si está revocado
let isRevoked = tokenBlacklist.isRevoked(token);

// Limpiar tokens expirados de la blacklist
tokenBlacklist.cleanup();
```

### 3. Rate Limiting

```javascript
// Rate limiting basado en JWT
let rateLimiter = jwt.createRateLimiter({
    "defaultLimit": 100,
    "window": "1h",
    "byUser": true,
    "byRole": {
        "admin": 1000,
        "user": 100,
        "guest": 10
    }
});

func apiCall(token, endpoint) {
    let validation = jwt.verifyAndDecode(token, getSecret());
    let allowed = rateLimiter.checkLimit(validation.payload.user, validation.payload.role);
    
    if (!allowed) {
        throw "Rate limit exceeded";
    }
    
    // Proceder con la llamada
    return processAPICall(endpoint);
}
```

## Algoritmos y Estándares Soportados

### Algoritmos de Firma

```javascript
// HMAC (Shared Secret)
let hsToken = jwt.create(payload, secret, {algorithm: "HS256"});
let hs384Token = jwt.create(payload, secret, {algorithm: "HS384"});
let hs512Token = jwt.create(payload, secret, {algorithm: "HS512"});

// RSA (Asymmetric)
let rsaToken = jwt.createRS256(payload, privateKey);
let rsa384Token = jwt.createRS384(payload, privateKey);
let rsa512Token = jwt.createRS512(payload, privateKey);

// ECDSA (Elliptic Curve)
let ecToken = jwt.createES256(payload, ecPrivateKey);
let ec384Token = jwt.createES384(payload, ecPrivateKey);
let ec512Token = jwt.createES512(payload, ecPrivateKey);

// PSS (Probabilistic Signature Scheme)
let pssToken = jwt.createPS256(payload, privateKey);
```

### Claims Estándar (RFC 7519)

```javascript
let standardClaims = {
    "iss": "issuer.company.com",        // Issuer
    "sub": "user@company.com",          // Subject  
    "aud": ["api.company.com"],         // Audience
    "exp": 1234567890,                  // Expiration Time
    "nbf": 1234567890,                  // Not Before
    "iat": 1234567890,                  // Issued At
    "jti": "unique-jwt-id-123"          // JWT ID
};

// Custom claims
let customClaims = {
    "user_id": 12345,
    "role": "admin",
    "permissions": ["read", "write"],
    "department": "IT",
    "security_level": "high"
};

let completeToken = jwt.create({
    ...standardClaims,
    ...customClaims
}, secret);
```

## Testing y Validación

### Unit Tests Requeridos

```javascript
// tests/jwt_test.r2
func testJWTBasicCreation() {
    let token = jwt.create({"user": "test"}, "secret");
    assert(token != null, "Token should be created");
    assert(len(token) > 0, "Token should not be empty");
}

func testJWTValidation() {
    let payload = {"user": "test", "role": "admin"};
    let token = jwt.create(payload, "secret");
    
    let validation = jwt.verifyAndDecode(token, "secret");
    assert(validation.valid == true, "Token should be valid");
    assert(validation.payload.user == "test", "User should match");
    assert(validation.payload.role == "admin", "Role should match");
}

func testJWTExpiration() {
    let token = jwt.create({"user": "test"}, "secret", "1s");
    sleep(2.0);
    
    let validation = jwt.verifyAndDecode(token, "secret");
    assert(validation.valid == false, "Expired token should be invalid");
}

func testJWTRS256() {
    let privateKey = loadTestPrivateKey();
    let publicKey = loadTestPublicKey();
    
    let token = jwt.createRS256({"user": "test"}, privateKey);
    let validation = jwt.verifyRS256(token, publicKey);
    
    assert(validation.valid == true, "RS256 token should be valid");
}
```

## Roadmap de Implementación

### Fase 1: Funcionalidad Básica (Semanas 1-2)
- Creación y validación básica de JWT
- Algoritmos HS256, HS384, HS512
- Claims estándar
- Integración con r2soap

### Fase 2: Características Empresariales (Semanas 3-4)
- Algoritmos RSA (RS256, RS384, RS512)
- Key rotation y múltiples secrets
- Blacklist y revocación
- Rate limiting

### Fase 3: Características Avanzadas (Semanas 5-6)
- JWC (JSON Web Certificates)
- ECDSA algorithms (ES256, ES384, ES512)
- OAuth 2.0 compliance
- SAML integration

### Fase 4: Enterprise Features (Semanas 7-8)
- LDAP integration
- Advanced RBAC
- Audit logging
- Performance optimization

## Beneficios Empresariales

### Seguridad
- ✅ **Stateless authentication** - No session storage needed
- ✅ **Strong cryptography** - Multiple signing algorithms
- ✅ **Token expiration** - Automatic security
- ✅ **Key rotation** - Enterprise-grade security

### Integración
- ✅ **Seamless r2soap integration** - Immediate value
- ✅ **Future r2grpc integration** - Unified security
- ✅ **Standard compliance** - RFC 7519, OAuth 2.0
- ✅ **Multi-language support** - JWT is universal

### Escalabilidad
- ✅ **Microservices ready** - Service-to-service auth
- ✅ **API Gateway compatible** - Standard bearer tokens
- ✅ **Cloud native** - Stateless design
- ✅ **Performance optimized** - Fast validation

### Flexibilidad
- ✅ **Multiple algorithms** - HS256 to ES512
- ✅ **Custom claims** - Business-specific data
- ✅ **Multiple formats** - JWT, JWC, OAuth
- ✅ **Integration patterns** - Various auth systems

## Conclusión

La implementación de r2jwt proporcionará a R2Lang capacidades de autenticación y autorización de nivel empresarial, completando el ecosistema de integración junto con r2soap y el futuro r2grpc. Esta librería permitirá:

1. **Autenticación moderna** con estándares de la industria
2. **Integración perfecta** con servicios existentes
3. **Seguridad empresarial** con múltiples algoritmos
4. **Flexibilidad de implementación** para diversos casos de uso

**🎯 Resultado**: R2Lang se convertirá en una plataforma completa para desarrollo de aplicaciones empresariales modernas con capacidades de integración SOAP, gRPC y autenticación JWT de primer nivel.

---

**Versión:** 1.0  
**Fecha:** 2024  
**Compatibilidad:** R2Lang v2.0+  
**Dependencias:** pkg/r2core, pkg/r2libs/r2soap