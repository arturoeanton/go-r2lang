# Curso R2Lang - Módulo 10: Despliegue y Distribución

## Introducción

En este módulo final aprenderás a desplegar aplicaciones R2Lang en producción, crear sistemas distribuidos, optimizar rendimiento, y mantener aplicaciones en entornos de producción. La arquitectura v2 proporciona herramientas avanzadas para despliegue y monitoreo.

### Herramientas de Despliegue v2

```
Deployment Tools:
├── Build System          # Compilación optimizada
├── Package Management    # Gestión de dependencias
├── Container Support     # Docker integration
├── Cloud Deployment     # AWS, GCP, Azure
├── Monitoring Tools     # Métricas y alertas
├── Load Balancing       # Distribución de carga
└── Auto-scaling         # Escalado automático
```

## Preparación para Producción

### 1. Optimización de Código

```r2
class OptimizadorRendimiento {
    let metricas
    let configuracion
    
    constructor() {
        this.metricas = {
            memoryUsage: 0,
            cpuUsage: 0,
            responseTime: 0,
            throughput: 0
        }
        this.configuracion = {
            cacheSize: 1000,
            maxConnections: 100,
            timeout: 30000,
            compressionEnabled: true,
            loggingLevel: "INFO"
        }
    }
    
    optimizarMemoria() {
        print("Optimizando uso de memoria...")
        
        // Implementar memory pooling
        let memoryPool = this.crearMemoryPool()
        
        // Cleanup de objetos no utilizados
        this.cleanupUnusedObjects()
        
        // Configurar garbage collection
        this.configurarGC()
        
        print("Optimización de memoria completada")
        return true
    }
    
    crearMemoryPool() {
        return {
            buffers: [],
            objects: [],
            strings: [],
            
            getBuffer: func(size) {
                // Reutilizar buffer existente si es posible
                for (let i = 0; i < this.buffers.length(); i++) {
                    if (this.buffers[i].size >= size && !this.buffers[i].inUse) {
                        this.buffers[i].inUse = true
                        return this.buffers[i]
                    }
                }
                
                // Crear nuevo buffer
                let buffer = {
                    size: size,
                    data: new Array(size),
                    inUse: true
                }
                this.buffers = this.buffers.push(buffer)
                return buffer
            },
            
            releaseBuffer: func(buffer) {
                buffer.inUse = false
                // Limpiar datos sensibles
                for (let i = 0; i < buffer.data.length(); i++) {
                    buffer.data[i] = null
                }
            }
        }
    }
    
    cleanupUnusedObjects() {
        // Simular cleanup de objetos no utilizados
        print("Limpiando objetos no utilizados...")
        
        // Remover references circulares
        // Limpiar caches expirados
        // Cerrar conexiones inactivas
        
        return true
    }
    
    configurarGC() {
        // Configurar garbage collection
        print("Configurando garbage collection...")
        
        // Ajustar frecuencia de GC
        // Configurar umbrales de memoria
        // Optimizar algoritmos de limpieza
        
        return true
    }
    
    optimizarRed() {
        print("Optimizando configuración de red...")
        
        // Connection pooling
        let connectionPool = this.crearConnectionPool()
        
        // Configurar timeouts
        this.configurarTimeouts()
        
        // Habilitar compresión
        this.habilitarCompresion()
        
        print("Optimización de red completada")
        return connectionPool
    }
    
    crearConnectionPool() {
        return {
            connections: [],
            maxConnections: this.configuracion.maxConnections,
            activeConnections: 0,
            
            getConnection: func() {
                if (this.activeConnections < this.maxConnections) {
                    let connection = {
                        id: this.activeConnections,
                        created: os.time(),
                        lastUsed: os.time(),
                        inUse: true
                    }
                    this.connections = this.connections.push(connection)
                    this.activeConnections++
                    return connection
                }
                return null
            },
            
            releaseConnection: func(connection) {
                connection.inUse = false
                connection.lastUsed = os.time()
            },
            
            cleanup: func() {
                let now = os.time()
                let newConnections = []
                
                for (let i = 0; i < this.connections.length(); i++) {
                    let conn = this.connections[i]
                    if (now - conn.lastUsed < 300000) {  // 5 minutos
                        newConnections = newConnections.push(conn)
                    } else {
                        this.activeConnections--
                    }
                }
                
                this.connections = newConnections
            }
        }
    }
    
    configurarTimeouts() {
        // Configurar timeouts para diferentes operaciones
        return {
            connection: 10000,  // 10 segundos
            read: 30000,        // 30 segundos
            write: 30000,       // 30 segundos
            idle: 300000        // 5 minutos
        }
    }
    
    habilitarCompresion() {
        // Habilitar compresión gzip para respuestas HTTP
        print("Habilitando compresión gzip...")
        return true
    }
}

func main() {
    let optimizador = OptimizadorRendimiento()
    
    print("=== OPTIMIZANDO APLICACIÓN PARA PRODUCCIÓN ===")
    
    optimizador.optimizarMemoria()
    let connectionPool = optimizador.optimizarRed()
    
    print("Aplicación optimizada para producción")
    print("Connection pool configurado con " + connectionPool.maxConnections + " conexiones máximas")
}
```

### 2. Sistema de Configuración

```r2
class ConfiguracionApp {
    let config
    let ambiente
    let secretos
    
    constructor(ambiente) {
        this.ambiente = ambiente
        this.config = {}
        this.secretos = {}
        this.cargarConfiguracion()
    }
    
    cargarConfiguracion() {
        // Configuración base
        this.config = {
            app: {
                name: "R2Lang App",
                version: "1.0.0",
                port: 3000,
                host: "localhost"
            },
            database: {
                host: "localhost",
                port: 5432,
                name: "r2lang_db",
                maxConnections: 20
            },
            redis: {
                host: "localhost",
                port: 6379,
                maxConnections: 10
            },
            logging: {
                level: "INFO",
                file: "app.log",
                rotateDaily: true
            }
        }
        
        // Configuración específica por ambiente
        if (this.ambiente == "development") {
            this.config.app.port = 3000
            this.config.logging.level = "DEBUG"
        } else if (this.ambiente == "staging") {
            this.config.app.port = 4000
            this.config.app.host = "staging.example.com"
        } else if (this.ambiente == "production") {
            this.config.app.port = 8080
            this.config.app.host = "app.example.com"
            this.config.logging.level = "WARN"
        }
        
        // Cargar desde variables de entorno
        this.cargarDesdeEnv()
        
        // Cargar secretos
        this.cargarSecretos()
    }
    
    cargarDesdeEnv() {
        // Cargar configuración desde variables de entorno
        let envPort = os.getEnv("PORT")
        if (envPort != null) {
            this.config.app.port = parseInt(envPort)
        }
        
        let envHost = os.getEnv("HOST")
        if (envHost != null) {
            this.config.app.host = envHost
        }
        
        let envDbHost = os.getEnv("DB_HOST")
        if (envDbHost != null) {
            this.config.database.host = envDbHost
        }
        
        let envLogLevel = os.getEnv("LOG_LEVEL")
        if (envLogLevel != null) {
            this.config.logging.level = envLogLevel
        }
    }
    
    cargarSecretos() {
        // En producción, cargar secretos desde un sistema seguro
        if (this.ambiente == "production") {
            this.secretos = {
                database: {
                    username: os.getEnv("DB_USERNAME"),
                    password: os.getEnv("DB_PASSWORD")
                },
                redis: {
                    password: os.getEnv("REDIS_PASSWORD")
                },
                jwt: {
                    secret: os.getEnv("JWT_SECRET")
                },
                encryption: {
                    key: os.getEnv("ENCRYPTION_KEY")
                }
            }
        } else {
            // Valores por defecto para desarrollo
            this.secretos = {
                database: {
                    username: "dev_user",
                    password: "dev_password"
                },
                redis: {
                    password: ""
                },
                jwt: {
                    secret: "dev_jwt_secret"
                },
                encryption: {
                    key: "dev_encryption_key"
                }
            }
        }
    }
    
    get(path) {
        // Obtener valor de configuración usando path notation
        let partes = path.split(".")
        let valor = this.config
        
        for (let i = 0; i < partes.length(); i++) {
            if (valor[partes[i]] != null) {
                valor = valor[partes[i]]
            } else {
                return null
            }
        }
        
        return valor
    }
    
    getSecreto(path) {
        // Obtener secreto usando path notation
        let partes = path.split(".")
        let valor = this.secretos
        
        for (let i = 0; i < partes.length(); i++) {
            if (valor[partes[i]] != null) {
                valor = valor[partes[i]]
            } else {
                return null
            }
        }
        
        return valor
    }
    
    validarConfiguracion() {
        let errores = []
        
        // Validar configuración requerida
        if (this.config.app.port == null) {
            errores = errores.push("Puerto de aplicación requerido")
        }
        
        if (this.config.app.host == null) {
            errores = errores.push("Host de aplicación requerido")
        }
        
        if (this.ambiente == "production") {
            if (this.secretos.database.password == null) {
                errores = errores.push("Password de base de datos requerido en producción")
            }
            
            if (this.secretos.jwt.secret == null) {
                errores = errores.push("JWT secret requerido en producción")
            }
        }
        
        return errores
    }
    
    mostrarConfiguracion() {
        print("=== CONFIGURACIÓN DE APLICACIÓN ===")
        print("Ambiente: " + this.ambiente)
        print("App: " + this.config.app.name + " v" + this.config.app.version)
        print("Host: " + this.config.app.host + ":" + this.config.app.port)
        print("Database: " + this.config.database.host + ":" + this.config.database.port)
        print("Log Level: " + this.config.logging.level)
        print("Secretos cargados: " + Object.keys(this.secretos).length())
    }
}

func main() {
    let ambientes = ["development", "staging", "production"]
    
    for (let ambiente in ambientes) {
        print("=== CONFIGURACIÓN PARA " + ambiente.upper() + " ===")
        
        let config = ConfiguracionApp(ambiente)
        config.mostrarConfiguracion()
        
        let errores = config.validarConfiguracion()
        if (errores.length() > 0) {
            print("ERRORES DE CONFIGURACIÓN:")
            for (let error in errores) {
                print("- " + error)
            }
        } else {
            print("✓ Configuración válida")
        }
        
        print()
    }
}
```

## Containerización y Despliegue

### 1. Configuración Docker

```r2
class DockerBuilder {
    let appName
    let version
    let ambiente
    
    constructor(appName, version, ambiente) {
        this.appName = appName
        this.version = version
        this.ambiente = ambiente
    }
    
    generarDockerfile() {
        let dockerfile = `# Dockerfile para ${this.appName}
FROM golang:1.21-alpine AS builder

# Instalar dependencias
RUN apk add --no-cache git

# Configurar directorio de trabajo
WORKDIR /app

# Copiar archivos de la aplicación
COPY . .

# Compilar aplicación
RUN go build -o r2lang main.go

# Etapa final
FROM alpine:latest

# Instalar certificados SSL
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copiar binario compilado
COPY --from=builder /app/r2lang .

# Copiar archivos de configuración
COPY --from=builder /app/config ./config
COPY --from=builder /app/examples ./examples

# Exponer puerto
EXPOSE 8080

# Comando por defecto
CMD ["./r2lang", "main.r2"]
`
        
        io.writeFile("Dockerfile", dockerfile)
        print("Dockerfile generado")
        return dockerfile
    }
    
    generarDockerCompose() {
        let compose = `version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - NODE_ENV=${this.ambiente}
      - PORT=8080
      - DB_HOST=database
      - REDIS_HOST=redis
    depends_on:
      - database
      - redis
    volumes:
      - ./logs:/app/logs
    restart: unless-stopped
    
  database:
    image: postgres:15
    environment:
      - POSTGRES_DB=r2lang_db
      - POSTGRES_USER=r2lang
      - POSTGRES_PASSWORD=r2lang_password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    restart: unless-stopped
    
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    restart: unless-stopped
    
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - app
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
`
        
        io.writeFile("docker-compose.yml", compose)
        print("Docker Compose generado")
        return compose
    }
    
    generarNginxConfig() {
        let nginx = `events {
    worker_connections 1024;
}

http {
    upstream app {
        server app:8080;
    }
    
    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    
    # Gzip compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;
    
    server {
        listen 80;
        server_name ${this.appName}.com;
        
        # Redirect HTTP to HTTPS
        return 301 https://$server_name$request_uri;
    }
    
    server {
        listen 443 ssl http2;
        server_name ${this.appName}.com;
        
        # SSL configuration
        ssl_certificate /etc/nginx/ssl/certificate.crt;
        ssl_certificate_key /etc/nginx/ssl/private.key;
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
        ssl_prefer_server_ciphers off;
        
        # Security headers
        add_header X-Frame-Options "SAMEORIGIN" always;
        add_header X-XSS-Protection "1; mode=block" always;
        add_header X-Content-Type-Options "nosniff" always;
        add_header Referrer-Policy "no-referrer-when-downgrade" always;
        add_header Content-Security-Policy "default-src 'self' http: https: data: blob: 'unsafe-inline'" always;
        
        # Rate limiting
        limit_req zone=api burst=20 nodelay;
        
        location / {
            proxy_pass http://app;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection 'upgrade';
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_cache_bypass $http_upgrade;
            
            # Timeouts
            proxy_connect_timeout 30s;
            proxy_send_timeout 30s;
            proxy_read_timeout 30s;
        }
        
        # Static files
        location /static {
            alias /app/static;
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
        
        # Health check
        location /health {
            access_log off;
            proxy_pass http://app;
        }
    }
}
`
        
        io.writeFile("nginx.conf", nginx)
        print("Nginx configuración generada")
        return nginx
    }
    
    generarScriptsDespliegue() {
        // Script de construcción
        let buildScript = `#!/bin/bash
set -e

echo "Construyendo imagen Docker..."
docker build -t ${this.appName}:${this.version} .
docker tag ${this.appName}:${this.version} ${this.appName}:latest

echo "Ejecutando tests..."
docker run --rm ${this.appName}:${this.version} ./r2lang test

echo "Construcción completada"
`
        
        io.writeFile("build.sh", buildScript)
        
        // Script de despliegue
        let deployScript = `#!/bin/bash
set -e

echo "Desplegando ${this.appName} v${this.version}..."

# Backup de la versión actual
docker-compose down
docker tag ${this.appName}:latest ${this.appName}:backup

# Desplegar nueva versión
docker-compose up -d

# Verificar que la aplicación esté funcionando
sleep 10
curl -f http://localhost/health || {
    echo "Deploy failed, rolling back..."
    docker-compose down
    docker tag ${this.appName}:backup ${this.appName}:latest
    docker-compose up -d
    exit 1
}

echo "Deploy completado exitosamente"
`
        
        io.writeFile("deploy.sh", deployScript)
        
        // Hacer scripts ejecutables
        os.exec("chmod +x build.sh")
        os.exec("chmod +x deploy.sh")
        
        print("Scripts de despliegue generados")
    }
    
    generarKubernetesYAML() {
        let k8s = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: ${this.appName}
spec:
  replicas: 3
  selector:
    matchLabels:
      app: ${this.appName}
  template:
    metadata:
      labels:
        app: ${this.appName}
    spec:
      containers:
      - name: ${this.appName}
        image: ${this.appName}:${this.version}
        ports:
        - containerPort: 8080
        env:
        - name: NODE_ENV
          value: "${this.ambiente}"
        - name: PORT
          value: "8080"
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: ${this.appName}-service
spec:
  selector:
    app: ${this.appName}
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ${this.appName}-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - host: ${this.appName}.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: ${this.appName}-service
            port:
              number: 80
`
        
        io.writeFile("k8s-deployment.yaml", k8s)
        print("Kubernetes deployment generado")
        return k8s
    }
}

func main() {
    let builder = DockerBuilder("r2lang-app", "1.0.0", "production")
    
    print("=== GENERANDO CONFIGURACIÓN DE DESPLIEGUE ===")
    
    builder.generarDockerfile()
    builder.generarDockerCompose()
    builder.generarNginxConfig()
    builder.generarScriptsDespliegue()
    builder.generarKubernetesYAML()
    
    print("✓ Configuración de despliegue generada")
    print("Archivos creados:")
    print("- Dockerfile")
    print("- docker-compose.yml")
    print("- nginx.conf")
    print("- build.sh")
    print("- deploy.sh")
    print("- k8s-deployment.yaml")
}
```

## Monitoreo y Logging

### 1. Sistema de Monitoreo

```r2
class SistemaMonitoreo {
    let metricas
    let alertas
    let configuracion
    
    constructor() {
        this.metricas = {
            http: {
                requests: 0,
                errors: 0,
                responseTime: [],
                statusCodes: {}
            },
            sistema: {
                cpuUsage: 0,
                memoryUsage: 0,
                diskUsage: 0,
                networkIO: 0
            },
            aplicacion: {
                activeUsers: 0,
                dbConnections: 0,
                cacheHits: 0,
                cacheMisses: 0
            }
        }
        
        this.alertas = []
        this.configuracion = {
            intervaloRecoleccion: 30000,  // 30 segundos
            umbralCPU: 80,
            umbralMemoria: 85,
            umbralDisco: 90,
            umbralResponseTime: 1000
        }
    }
    
    recolectarMetricas() {
        // Recolectar métricas del sistema
        this.metricas.sistema.cpuUsage = os.getCpuUsage()
        this.metricas.sistema.memoryUsage = os.getMemoryUsage()
        this.metricas.sistema.diskUsage = os.getDiskUsage()
        this.metricas.sistema.networkIO = os.getNetworkIO()
        
        // Recolectar métricas de aplicación
        this.metricas.aplicacion.activeUsers = this.obtenerUsuariosActivos()
        this.metricas.aplicacion.dbConnections = this.obtenerConexionesDB()
        
        // Calcular métricas derivadas
        this.calcularMetricasDerivadas()
        
        // Verificar alertas
        this.verificarAlertas()
    }
    
    obtenerUsuariosActivos() {
        // Simular obtención de usuarios activos
        return rand.int(50, 200)
    }
    
    obtenerConexionesDB() {
        // Simular obtención de conexiones DB
        return rand.int(5, 20)
    }
    
    calcularMetricasDerivadas() {
        // Calcular throughput
        this.metricas.aplicacion.throughput = this.metricas.http.requests / 60  // req/min
        
        // Calcular error rate
        if (this.metricas.http.requests > 0) {
            this.metricas.aplicacion.errorRate = (this.metricas.http.errors / this.metricas.http.requests) * 100
        }
        
        // Calcular cache hit rate
        let totalCacheRequests = this.metricas.aplicacion.cacheHits + this.metricas.aplicacion.cacheMisses
        if (totalCacheRequests > 0) {
            this.metricas.aplicacion.cacheHitRate = (this.metricas.aplicacion.cacheHits / totalCacheRequests) * 100
        }
    }
    
    verificarAlertas() {
        let alertasActuales = []
        
        // Verificar CPU
        if (this.metricas.sistema.cpuUsage > this.configuracion.umbralCPU) {
            alertasActuales = alertasActuales.push({
                tipo: "CPU_HIGH",
                mensaje: "Alto uso de CPU: " + this.metricas.sistema.cpuUsage + "%",
                severidad: "WARNING",
                timestamp: os.time()
            })
        }
        
        // Verificar memoria
        if (this.metricas.sistema.memoryUsage > this.configuracion.umbralMemoria) {
            alertasActuales = alertasActuales.push({
                tipo: "MEMORY_HIGH",
                mensaje: "Alto uso de memoria: " + this.metricas.sistema.memoryUsage + "%",
                severidad: "WARNING",
                timestamp: os.time()
            })
        }
        
        // Verificar disco
        if (this.metricas.sistema.diskUsage > this.configuracion.umbralDisco) {
            alertasActuales = alertasActuales.push({
                tipo: "DISK_HIGH",
                mensaje: "Alto uso de disco: " + this.metricas.sistema.diskUsage + "%",
                severidad: "CRITICAL",
                timestamp: os.time()
            })
        }
        
        // Verificar tiempo de respuesta
        if (this.metricas.http.responseTime.length() > 0) {
            let avgResponseTime = math.average(this.metricas.http.responseTime)
            if (avgResponseTime > this.configuracion.umbralResponseTime) {
                alertasActuales = alertasActuales.push({
                    tipo: "RESPONSE_TIME_HIGH",
                    mensaje: "Alto tiempo de respuesta: " + avgResponseTime + "ms",
                    severidad: "WARNING",
                    timestamp: os.time()
                })
            }
        }
        
        // Procesar nuevas alertas
        for (let alerta in alertasActuales) {
            this.procesarAlerta(alerta)
        }
    }
    
    procesarAlerta(alerta) {
        print("🚨 ALERTA: " + alerta.mensaje)
        
        // Agregar a lista de alertas
        this.alertas = this.alertas.push(alerta)
        
        // Enviar notificación
        this.enviarNotificacion(alerta)
        
        // Log de alerta
        this.logAlerta(alerta)
    }
    
    enviarNotificacion(alerta) {
        // Enviar a Slack, email, etc.
        let webhook = "https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
        let payload = {
            text: "🚨 " + alerta.mensaje,
            channel: "#alerts",
            username: "R2Lang Monitor",
            attachments: [{
                color: alerta.severidad == "CRITICAL" ? "danger" : "warning",
                fields: [{
                    title: "Severidad",
                    value: alerta.severidad,
                    short: true
                }, {
                    title: "Tiempo",
                    value: alerta.timestamp,
                    short: true
                }]
            }]
        }
        
        http.post(webhook, {
            "Content-Type": "application/json",
            "data": JSON.stringify(payload)
        })
    }
    
    logAlerta(alerta) {
        let logEntry = "[" + alerta.timestamp + "] " + alerta.severidad + ": " + alerta.mensaje
        
        let logContent = io.readFile("alerts.log")
        if (logContent == null) {
            logContent = ""
        }
        
        io.writeFile("alerts.log", logContent + logEntry + "\n")
    }
    
    generarReporte() {
        print("=== REPORTE DE MONITOREO ===")
        print("Timestamp: " + os.time())
        print()
        
        print("SISTEMA:")
        print("  CPU: " + this.metricas.sistema.cpuUsage + "%")
        print("  Memoria: " + this.metricas.sistema.memoryUsage + "%")
        print("  Disco: " + this.metricas.sistema.diskUsage + "%")
        print("  Red I/O: " + this.metricas.sistema.networkIO + " KB/s")
        print()
        
        print("HTTP:")
        print("  Requests: " + this.metricas.http.requests)
        print("  Errors: " + this.metricas.http.errors)
        print("  Error Rate: " + this.metricas.aplicacion.errorRate + "%")
        if (this.metricas.http.responseTime.length() > 0) {
            print("  Avg Response Time: " + math.average(this.metricas.http.responseTime) + "ms")
        }
        print()
        
        print("APLICACIÓN:")
        print("  Usuarios activos: " + this.metricas.aplicacion.activeUsers)
        print("  Conexiones DB: " + this.metricas.aplicacion.dbConnections)
        print("  Cache Hit Rate: " + this.metricas.aplicacion.cacheHitRate + "%")
        print("  Throughput: " + this.metricas.aplicacion.throughput + " req/min")
        print()
        
        if (this.alertas.length() > 0) {
            print("ALERTAS RECIENTES:")
            for (let i = math.max(0, this.alertas.length() - 5); i < this.alertas.length(); i++) {
                let alerta = this.alertas[i]
                print("  [" + alerta.severidad + "] " + alerta.mensaje)
            }
        }
    }
    
    iniciarMonitoreo() {
        print("Iniciando sistema de monitoreo...")
        
        // Crear servidor de métricas
        let self = this
        let servidor = http.server(9090, func(request, response) {
            if (request.url == "/metrics") {
                response.writeHead(200, {"Content-Type": "application/json"})
                response.write(JSON.stringify(self.metricas))
                response.end()
            } else if (request.url == "/health") {
                response.writeHead(200, {"Content-Type": "application/json"})
                response.write('{"status": "healthy"}')
                response.end()
            } else if (request.url == "/alerts") {
                response.writeHead(200, {"Content-Type": "application/json"})
                response.write(JSON.stringify(self.alertas))
                response.end()
            }
        })
        
        print("Servidor de métricas iniciado en puerto 9090")
        
        // Bucle principal de monitoreo
        while (true) {
            this.recolectarMetricas()
            this.generarReporte()
            
            sleep(this.configuracion.intervaloRecoleccion)
        }
    }
}

func main() {
    let monitor = SistemaMonitoreo()
    monitor.iniciarMonitoreo()
}
```

## Proyecto Final: Aplicación Completa de Producción

```r2
class AplicacionProduccion {
    let servidor
    let configuracion
    let monitor
    let logger
    let database
    let cache
    
    constructor() {
        this.configuracion = ConfiguracionApp(os.getEnv("NODE_ENV") || "development")
        this.monitor = SistemaMonitoreo()
        this.logger = Logger("app.log", this.configuracion.get("logging.level"))
        this.database = this.inicializarDatabase()
        this.cache = this.inicializarCache()
        
        this.validarConfiguracion()
    }
    
    validarConfiguracion() {
        let errores = this.configuracion.validarConfiguracion()
        if (errores.length() > 0) {
            for (let error in errores) {
                this.logger.error("Error de configuración: " + error)
            }
            throw "Configuración inválida"
        }
    }
    
    inicializarDatabase() {
        let dbConfig = this.configuracion.get("database")
        
        return {
            host: dbConfig.host,
            port: dbConfig.port,
            name: dbConfig.name,
            maxConnections: dbConfig.maxConnections,
            connectionPool: [],
            
            connect: func() {
                // Simular conexión a base de datos
                return {
                    id: rand.int(1, 1000),
                    connected: true,
                    lastUsed: os.time()
                }
            },
            
            query: func(sql, params) {
                // Simular query
                let startTime = os.time()
                sleep(rand.int(10, 100))  // Simular tiempo de query
                let endTime = os.time()
                
                return {
                    rows: [],
                    time: endTime - startTime,
                    success: true
                }
            }
        }
    }
    
    inicializarCache() {
        let cacheConfig = this.configuracion.get("redis")
        
        return {
            host: cacheConfig.host,
            port: cacheConfig.port,
            storage: {},
            
            get: func(key) {
                return this.storage[key]
            },
            
            set: func(key, value, ttl) {
                this.storage[key] = {
                    value: value,
                    expiry: os.time() + (ttl || 3600000)  // 1 hora por defecto
                }
            },
            
            del: func(key) {
                delete this.storage[key]
            }
        }
    }
    
    crearAPIs() {
        let self = this
        
        return {
            "/api/users": {
                GET: func(request, response) {
                    self.handleGetUsers(request, response)
                },
                POST: func(request, response) {
                    self.handleCreateUser(request, response)
                }
            },
            
            "/api/users/:id": {
                GET: func(request, response) {
                    self.handleGetUser(request, response)
                },
                PUT: func(request, response) {
                    self.handleUpdateUser(request, response)
                },
                DELETE: func(request, response) {
                    self.handleDeleteUser(request, response)
                }
            },
            
            "/api/health": {
                GET: func(request, response) {
                    self.handleHealthCheck(request, response)
                }
            },
            
            "/api/metrics": {
                GET: func(request, response) {
                    self.handleMetrics(request, response)
                }
            }
        }
    }
    
    handleGetUsers(request, response) {
        let startTime = os.time()
        
        try {
            // Verificar cache
            let cachedUsers = this.cache.get("users:all")
            if (cachedUsers != null) {
                this.logger.info("Cache hit for users list")
                response.writeHead(200, {"Content-Type": "application/json"})
                response.write(JSON.stringify(cachedUsers.value))
                response.end()
                return
            }
            
            // Query desde base de datos
            let result = this.database.query("SELECT * FROM users")
            let users = result.rows
            
            // Guardar en cache
            this.cache.set("users:all", users, 300000)  // 5 minutos
            
            response.writeHead(200, {"Content-Type": "application/json"})
            response.write(JSON.stringify(users))
            response.end()
            
            this.logger.info("GET /api/users - 200 - " + (os.time() - startTime) + "ms")
        } catch (error) {
            this.logger.error("Error in GET /api/users: " + error)
            response.writeHead(500, {"Content-Type": "application/json"})
            response.write('{"error": "Internal server error"}')
            response.end()
        }
    }
    
    handleCreateUser(request, response) {
        let startTime = os.time()
        
        try {
            let userData = JSON.parse(request.body)
            
            // Validar datos
            if (!userData.name || !userData.email) {
                response.writeHead(400, {"Content-Type": "application/json"})
                response.write('{"error": "Name and email are required"}')
                response.end()
                return
            }
            
            // Crear usuario en base de datos
            let result = this.database.query(
                "INSERT INTO users (name, email) VALUES (?, ?)",
                [userData.name, userData.email]
            )
            
            if (result.success) {
                let newUser = {
                    id: result.insertId,
                    name: userData.name,
                    email: userData.email,
                    created: os.time()
                }
                
                // Invalidar cache
                this.cache.del("users:all")
                
                response.writeHead(201, {"Content-Type": "application/json"})
                response.write(JSON.stringify(newUser))
                response.end()
                
                this.logger.info("POST /api/users - 201 - " + (os.time() - startTime) + "ms")
            } else {
                throw "Database insert failed"
            }
        } catch (error) {
            this.logger.error("Error in POST /api/users: " + error)
            response.writeHead(500, {"Content-Type": "application/json"})
            response.write('{"error": "Internal server error"}')
            response.end()
        }
    }
    
    handleHealthCheck(request, response) {
        let health = {
            status: "healthy",
            timestamp: os.time(),
            uptime: os.time() - this.startTime,
            version: this.configuracion.get("app.version"),
            environment: this.configuracion.ambiente,
            checks: {
                database: this.checkDatabase(),
                cache: this.checkCache(),
                memory: this.checkMemory()
            }
        }
        
        let allHealthy = true
        for (let check in Object.values(health.checks)) {
            if (!check.healthy) {
                allHealthy = false
                break
            }
        }
        
        let statusCode = allHealthy ? 200 : 503
        response.writeHead(statusCode, {"Content-Type": "application/json"})
        response.write(JSON.stringify(health))
        response.end()
    }
    
    checkDatabase() {
        try {
            let result = this.database.query("SELECT 1")
            return {
                healthy: result.success,
                responseTime: result.time
            }
        } catch (error) {
            return {
                healthy: false,
                error: error
            }
        }
    }
    
    checkCache() {
        try {
            this.cache.set("health:test", "ok", 1000)
            let value = this.cache.get("health:test")
            return {
                healthy: value != null && value.value == "ok"
            }
        } catch (error) {
            return {
                healthy: false,
                error: error
            }
        }
    }
    
    checkMemory() {
        let memoryUsage = os.getMemoryUsage()
        return {
            healthy: memoryUsage < 90,
            usage: memoryUsage + "%"
        }
    }
    
    handleMetrics(request, response) {
        let metrics = this.monitor.metricas
        response.writeHead(200, {"Content-Type": "application/json"})
        response.write(JSON.stringify(metrics))
        response.end()
    }
    
    configurarMiddleware() {
        let self = this
        
        return {
            logging: func(request, response, next) {
                let startTime = os.time()
                
                // Log de request
                self.logger.info("Request: " + request.method + " " + request.url)
                
                // Continuar con siguiente middleware
                next()
                
                // Log de response
                let duration = os.time() - startTime
                self.logger.info("Response: " + response.statusCode + " - " + duration + "ms")
            },
            
            cors: func(request, response, next) {
                response.setHeader("Access-Control-Allow-Origin", "*")
                response.setHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
                response.setHeader("Access-Control-Allow-Headers", "Content-Type, Authorization")
                
                if (request.method == "OPTIONS") {
                    response.writeHead(200)
                    response.end()
                    return
                }
                
                next()
            },
            
            rateLimit: func(request, response, next) {
                // Implementar rate limiting básico
                let clientIP = request.headers["x-forwarded-for"] || request.connection.remoteAddress
                let key = "rate_limit:" + clientIP
                
                let count = self.cache.get(key)
                if (count == null) {
                    self.cache.set(key, 1, 60000)  // 1 minuto
                } else {
                    if (count.value >= 100) {  // 100 requests por minuto
                        response.writeHead(429, {"Content-Type": "application/json"})
                        response.write('{"error": "Too many requests"}')
                        response.end()
                        return
                    }
                    self.cache.set(key, count.value + 1, 60000)
                }
                
                next()
            }
        }
    }
    
    iniciar() {
        this.startTime = os.time()
        let port = this.configuracion.get("app.port")
        let host = this.configuracion.get("app.host")
        
        this.logger.info("Iniciando aplicación en " + host + ":" + port)
        
        // Inicializar APIs
        let apis = this.crearAPIs()
        let middleware = this.configurarMiddleware()
        
        // Crear servidor HTTP
        let self = this
        this.servidor = http.server(port, func(request, response) {
            // Aplicar middleware
            middleware.logging(request, response, func() {
                middleware.cors(request, response, func() {
                    middleware.rateLimit(request, response, func() {
                        self.routeRequest(request, response, apis)
                    })
                })
            })
        })
        
        // Iniciar monitoreo
        r2(this.monitor.iniciarMonitoreo)
        
        // Graceful shutdown
        this.configurarGracefulShutdown()
        
        this.logger.info("Aplicación iniciada exitosamente")
        print("🚀 Aplicación corriendo en http://" + host + ":" + port)
    }
    
    routeRequest(request, response, apis) {
        let path = request.url
        let method = request.method
        
        // Buscar ruta exacta
        if (apis[path] && apis[path][method]) {
            apis[path][method](request, response)
            return
        }
        
        // Buscar rutas con parámetros
        for (let route in Object.keys(apis)) {
            if (route.contains(":")) {
                let pattern = route.replace(/:([^/]+)/g, "([^/]+)")
                let regex = new RegExp("^" + pattern + "$")
                
                if (regex.test(path)) {
                    // Extraer parámetros
                    let matches = path.match(regex)
                    request.params = {}
                    
                    if (apis[route][method]) {
                        apis[route][method](request, response)
                        return
                    }
                }
            }
        }
        
        // Ruta no encontrada
        response.writeHead(404, {"Content-Type": "application/json"})
        response.write('{"error": "Not found"}')
        response.end()
    }
    
    configurarGracefulShutdown() {
        // Configurar manejo de señales para shutdown graceful
        let self = this
        
        process.on("SIGINT", func() {
            self.logger.info("Recibida señal SIGINT, iniciando shutdown...")
            self.shutdown()
        })
        
        process.on("SIGTERM", func() {
            self.logger.info("Recibida señal SIGTERM, iniciando shutdown...")
            self.shutdown()
        })
    }
    
    shutdown() {
        this.logger.info("Cerrando aplicación...")
        
        // Cerrar servidor HTTP
        if (this.servidor) {
            this.servidor.close()
        }
        
        // Cerrar conexiones de base de datos
        this.database.close()
        
        // Cerrar cache
        this.cache.close()
        
        this.logger.info("Aplicación cerrada exitosamente")
        process.exit(0)
    }
}

func main() {
    try {
        let app = AplicacionProduccion()
        app.iniciar()
    } catch (error) {
        print("Error crítico al iniciar aplicación: " + error)
        process.exit(1)
    }
}
```

## Resumen del Curso Completo

### Módulos Completados
- ✅ **Módulo 1**: Fundamentos y variables
- ✅ **Módulo 2**: Control de flujo y funciones
- ✅ **Módulo 3**: Orientación a objetos
- ✅ **Módulo 4**: Manejo de errores y arrays
- ✅ **Módulo 5**: Imports y modularización
- ✅ **Módulo 6**: Testing integrado BDD
- ✅ **Módulo 7**: Bibliotecas integradas
- ✅ **Módulo 8**: Concurrencia y paralelismo
- ✅ **Módulo 9**: Testing avanzado y debugging
- ✅ **Módulo 10**: Despliegue y producción

### Habilidades Desarrolladas
- ✅ **Programación**: Dominio completo del lenguaje R2Lang
- ✅ **Arquitectura**: Diseño de aplicaciones escalables
- ✅ **Testing**: BDD, unit testing, integration testing
- ✅ **Concurrencia**: Programación paralela avanzada
- ✅ **Sistemas**: Interacción con SO, red, archivos
- ✅ **Despliegue**: Containerización, Kubernetes, CI/CD
- ✅ **Monitoreo**: Observabilidad y alertas
- ✅ **Producción**: Aplicaciones enterprise-ready

### Tecnologías Integradas
- ✅ **Docker y Kubernetes**
- ✅ **Nginx y Load Balancing**
- ✅ **Monitoring y Alertas**
- ✅ **CI/CD Pipelines**
- ✅ **Database Integration**
- ✅ **Cache Systems**
- ✅ **Security Best Practices**

### Próximos Pasos

¡Felicitaciones! Has completado el curso completo de R2Lang. Ahora estás preparado para:

1. **Desarrollar aplicaciones profesionales** con R2Lang
2. **Contribuir al proyecto open source** R2Lang
3. **Crear bibliotecas y extensiones** para la comunidad
4. **Enseñar R2Lang** a otros desarrolladores
5. **Construir sistemas distribuidos** complejos

### Recursos Adicionales

- **GitHub**: Contribuye al proyecto en `github.com/arturoeanton/go-r2lang`
- **Documentación**: Explora `docs/` para referencia avanzada
- **Ejemplos**: Revisa `examples/` para casos de uso reales
- **Comunidad**: Únete a las discusiones y ayuda a otros desarrolladores

¡Gracias por completar este curso! R2Lang te espera para construir el futuro de la programación.

🎉 **¡Eres ahora un R2Lang Expert!** 🎉