# 🎯 INSTRUCCIONES DEMO SIIGO - FUNCIONANDO 100%

## ✅ COMANDO QUE FUNCIONA (Usar este)

```bash
go run main.go examples/proyecto/contable/contable_web_simple.r2
```

## 🌐 URLs del Demo

1. **Página Principal**: http://localhost:8080
   - Interfaz visual completa
   - Formularios para procesar ventas y compras
   - Selección de regiones LATAM

2. **Demo Automático**: http://localhost:8080/demo
   - Procesamiento automático multi-región
   - 4 transacciones de ejemplo
   - Comprobantes visuales

3. **API Regiones**: http://localhost:8080/api/regiones
   - JSON con configuración de las 7 regiones
   - Datos técnicos para desarrolladores

4. **API Transacciones**: http://localhost:8080/api/transacciones
   - Historial de transacciones procesadas
   - Total y detalles en JSON

## 🎪 Flow de Demostración para Siigo

### Demo Rápido (5 minutos)
1. Ejecutar: `go run main.go examples/proyecto/contable/contable_web_simple.r2`
2. Abrir: http://localhost:8080
3. Mostrar formulario visual
4. Procesar 1 venta Colombia $100,000
5. Mostrar comprobante generado automáticamente
6. Ir a: http://localhost:8080/demo
7. Mostrar procesamiento automático 4 regiones

### Demo Completo (15 minutos)
1. **Inicio**: Mostrar interfaz principal
2. **Venta Manual**: 
   - Seleccionar región México
   - Importe $50,000 
   - Ver comprobante con IVA 16%
3. **Compra Manual**:
   - Seleccionar región Argentina  
   - Importe $30,000
   - Ver comprobante con IVA 21%
4. **Demo Automático**: 
   - Abrir http://localhost:8080/demo
   - Mostrar 4 transacciones multi-región
5. **APIs**:
   - Mostrar http://localhost:8080/api/regiones
   - Mostrar http://localhost:8080/api/transacciones
6. **Value Proposition**: Resaltar savings en pantalla

## 📊 Funcionalidades Demostradas

### ✅ Procesamiento Visual
- Formularios interactivos HTML
- Selección dropdown de países  
- Comprobantes visuales generados
- Cálculos automáticos de impuestos
- Formato nativo por moneda

### ✅ Multi-Región LATAM
- 🇲🇽 México: 16% IVA, MXN
- 🇨🇴 Colombia: 19% IVA, COP
- 🇦🇷 Argentina: 21% IVA, ARS  
- 🇨🇱 Chile: 19% IVA, CLP
- 🇺🇾 Uruguay: 22% IVA, UYU
- 🇪🇨 Ecuador: 12% IVA, USD
- 🇵🇪 Perú: 18% IVA, PEN

### ✅ DSL Engine
- Procesamiento automático con sintaxis natural
- VentasWeb y ComprasWeb DSL engines
- Validación automática de entrada
- Resultados estructurados en JSON

### ✅ APIs REST
- GET /api/regiones - Configuración regional
- GET /api/transacciones - Historial completo
- Datos listos para integración ERP

## 💡 Value Proposition Mostrado

En la página principal se muestra claramente:

- ✅ **18 meses → 2 meses** por país
- ✅ **$500K → $150K** por localización  
- ✅ **7 codebases → 1 DSL** unificado
- ✅ **ROI: 1,020%** en 3 años

## ❌ NO USAR Estos Archivos (Tienen Errores)

```bash
# ❌ NO FUNCIONA - Problemas con HttpResponse  
go run main.go examples/proyecto/contable/main_web.r2

# ❌ NO FUNCIONA - Import issues
go run main.go examples/proyecto/contable/main.r2
```

## ✅ Archivos que SÍ Funcionan

```bash
# 🎯 PRINCIPAL - POC Web Visual (RECOMENDADO)
go run main.go examples/proyecto/contable/contable_web_simple.r2

# ✅ Demo consola completo  
go run main.go examples/proyecto/contable/demo_completo.r2

# ✅ Demo básico sin web
go run main.go examples/proyecto/contable/main_simple.r2
```

## 🚨 Solución de Problemas

### Si sale error "HttpResponse not found"
- **Causa**: Estás usando archivo incorrecto
- **Solución**: Usar `contable_web_simple.r2` en lugar de `main_web.r2`

### Si no carga la página web
- **Verificar**: Que esté ejecutando `contable_web_simple.r2`
- **URL correcta**: http://localhost:8080 (no https)
- **Puerto libre**: `lsof -i :8080` debe estar libre

### Si aparece "Listening on :8080"
- ✅ **Correcto**: El servidor está funcionando
- **Acción**: Abrir navegador en http://localhost:8080

## 🎯 Comando Final para Demo Siigo

```bash
# Ejecutar este comando y listo:
go run main.go examples/proyecto/contable/contable_web_simple.r2

# Luego abrir navegador en:
http://localhost:8080
```

**🎉 ¡POC Lista para Presentación Siigo!**