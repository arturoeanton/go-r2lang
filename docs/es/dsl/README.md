# Documentación DSL de R2Lang

## Última Actualización
2026-07-02

## Nota de migración (2026-07-02)

El motor interno que parseaba los bloques `dsl { }` (tokenizer y parser
hand-rolled en `pkg/r2core/dsl_grammar.go`) fue reemplazado por
[`github.com/arturoeanton/go-dsl`](https://github.com/arturoeanton/go-dsl)
v1.4.0 (`pkg/dslbuilder`). La sintaxis `.r2` (`token()`, `rule()`, `action()`,
`func`, `.use(code, context)`, `.AST`/`.Code`/`.Output`) no cambió — es la
misma API descrita en este manual — pero el backend ahora es una gramática
PEG con memoización (Packrat), soporte de recursión izquierda y tokenización
determinística (prioridad > longitud > orden de declaración), lo que
resuelve las inestabilidades de parsing documentadas en
[Limitaciones y Mejoras Priorizadas](limitaciones_y_mejoras_priorizadas.md).
Una diferencia real de comportamiento: `.AST` ahora expone el árbol de
sintaxis real (`*dslbuilder.Node`) en vez de reusar el valor de `.Output`.

## Descripción

Esta carpeta contiene toda la documentación oficial del sistema DSL (Domain-Specific Language) de R2Lang en español. El sistema DSL permite crear lenguajes específicos de dominio de manera simple y elegante.

## Estructura de Documentación

### 📚 Manuales y Guías

#### [Manual Completo de DSL](manual_dsl_completo.md)
**Última versión - Recomendado**
- Guía completa y actualizada del sistema DSL
- Incluye todos los fixes y mejoras recientes
- Ejemplos prácticos y casos de uso
- Mejores prácticas y solución de problemas
- Referencia API completa

#### [Manual DSL v1](manual_dsl_v1.md)
**Versión histórica**
- Primera versión del manual DSL
- Funcionalidades básicas
- Mantenido por razones de compatibilidad

### 🔧 Fixes y Mejoras

#### [Fix DSL Parámetros y Resultados](fix_dsl_parametros_y_resultados.md)
**Importante - Últimos fixes**
- Solución completa para problema de parámetros `{+}`
- Implementación de resultados DSL accesibles
- Mejoras en la experiencia de desarrollador
- Desempacado de ReturnValue

#### [Fix DSL Parámetros Múltiples](fix_dsl_parametros_multiples.md)
**Histórico**
- Primera versión del fix para parámetros múltiples
- Documento de referencia para el problema inicial

### 📊 Análisis y Comparaciones

#### [Madurez DSL R2Lang](madurez_dsl_r2lang.md)
**Evaluación técnica**
- Análisis de madurez del feature DSL
- Puntuación: 9.0/10 (Beta Estable)
- Criterios de evaluación técnica
- Casos de uso recomendados

#### [Comparación DSL vs Competidores](comparacion_dsl_competidores.md)
**Análisis competitivo**
- Comparación con ANTLR, Lex/Yacc, PEG.js y otros
- Matriz de comparación detallada
- Ventajas y desventajas
- Recomendaciones por caso de uso

## Guía de Inicio Rápido

### 1. Empezar Aquí
Si eres nuevo en DSL de R2Lang, empieza con el [Manual Completo de DSL](manual_dsl_completo.md).

### 2. Ejemplo Básico
```r2
dsl MiPrimerDSL {
    token("NUMERO", "[0-9]+")
    token("SUMA", "\\+")
    
    rule("suma", ["NUMERO", "SUMA", "NUMERO"], "sumar")
    
    func sumar(a, op, b) {
        return parseFloat(a) + parseFloat(b)
    }
}

func main() {
    var calc = MiPrimerDSL.use
    var resultado = calc("5 + 3")
    console.log(resultado.Output)  // 8
}
```

### 3. Próximos Pasos
1. Lee el [Manual Completo](manual_dsl_completo.md) - Sección "Conceptos Básicos"
2. Prueba los ejemplos en la sección "Ejemplos Prácticos"
3. Consulta [Mejores Prácticas](manual_dsl_completo.md#mejores-prácticas)

## Estado del Sistema DSL

### ✅ Funcionalidades Implementadas
- ✅ Definición de tokens con regex
- ✅ Definición de reglas gramaticales
- ✅ Acciones semánticas con funciones R2Lang
- ✅ Parsing robusto con manejo de errores
- ✅ Resultados estructurados (DSLResult)
- ✅ Acceso a propiedades (.Output, .Code, .AST)
- ✅ Representación string mejorada
- ✅ Desempacado correcto de parámetros
- ✅ Tests unitarios completos

### 🔧 Fixes Recientes
- ✅ Problema `{+}` en parámetros → Resuelto
- ✅ Resultados DSL no accesibles → Resuelto
- ✅ Desempacado de ReturnValue → Resuelto
- ✅ Método String() mejorado → Implementado

### 📊 Métricas de Calidad
- **Puntuación General**: 9.0/10
- **Estado**: Beta Estable
- **Cobertura Tests**: >90%
- **Usabilidad**: 10/10
- **Performance**: 8/10

## Casos de Uso Típicos

### 🎯 Ideales para DSL R2Lang
- **Configuraciones personalizadas**
- **Validadores de datos**
- **Sistemas de comandos**
- **Procesadores de texto**
- **Calculadoras especializadas**
- **Lenguajes de scripting simples**

### ⚠️ Considera Alternativas para
- **Lenguajes complejos completos** (usa ANTLR)
- **Performance crítica** (usa Lex/Yacc)
- **Múltiples lenguajes target** (usa ANTLR)

## Soporte y Contribución

### 🐛 Reportar Problemas
1. Revisa la sección [Solución de Problemas](manual_dsl_completo.md#solución-de-problemas)
2. Consulta los [fixes conocidos](fix_dsl_parametros_y_resultados.md)
3. Crea un issue en GitHub

### 📝 Contribuir
1. Lee el [Manual Completo](manual_dsl_completo.md)
2. Revisa los [tests unitarios](../../../pkg/r2core/dsl_test.go)
3. Sigue las [mejores prácticas](manual_dsl_completo.md#mejores-prácticas)

## Historial de Versiones

### v3.0 (2026-07-02)
- ✅ Motor de parsing migrado a `github.com/arturoeanton/go-dsl` v1.4.0
- ✅ Tokenización determinística y parser con memoización (Packrat)
- ✅ `.AST` ahora expone el árbol de sintaxis real, separado de `.Output`
- ✅ Misma superficie pública (`token`, `rule`, `action`, `func`, `.use`)

### v2.0 (2025-01-18)
- ✅ Fix completo para parámetros múltiples
- ✅ Resultados DSL accesibles
- ✅ Método String() mejorado
- ✅ Tests unitarios completos
- ✅ Documentación actualizada

### v1.0 (2024)
- ✅ Implementación inicial del DSL
- ✅ Funcionalidades básicas
- ✅ Documentación inicial

## Enlaces Útiles

### Documentación Técnica
- [Arquitectura R2Lang](../arquitectura-profunda.md)
- [Tests Unitarios](../../../pkg/r2core/dsl_test.go)
- [Ejemplos DSL](../../../examples/dsl/)

### Repositorio
- [GitHub R2Lang](https://github.com/arturoeanton/go-r2lang)
- [Issues](https://github.com/arturoeanton/go-r2lang/issues)
- [Releases](https://github.com/arturoeanton/go-r2lang/releases)

---

**Mantenido por**: Equipo R2Lang  
**Última actualización**: 2025-01-18  
**Versión DSL**: v2.0 Beta Estable