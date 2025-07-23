# Reporte de Performance - R2Lang
Fecha: 2025-07-22 21:32:13
Sistema: darwin arm64
CPUs: 8
Versión Go: go1.24.5

## Benchmarks Ejecutados

Para ejecutar estos benchmarks:
```bash
go test -bench=. -benchmem performance_test.go
```

## Casos de Prueba

1. **Operaciones Aritméticas Básicas**: Loop con 1000 iteraciones de cálculos
2. **Operaciones de String**: Concatenación de strings en loop
3. **Operaciones de Array**: Creación y manipulación de arrays
4. **Operaciones de Map**: Creación y acceso a mapas
5. **Llamadas a Funciones**: Fibonacci recursivo
6. **Operaciones con Objetos**: Creación y métodos de objetos
7. **Rendimiento del Lexer**: Análisis léxico de código complejo
8. **Rendimiento del Parser**: Análisis sintáctico
9. **Uso de Memoria**: Creación de estructuras grandes

## Análisis y Mejoras Recomendadas

Los resultados de estos benchmarks se detallan en docs/es/performance.md
