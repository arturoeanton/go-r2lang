# Fix: Problema con Parámetros Múltiples en DSL

## Fecha
2025-01-18

## Problema Identificado

Los DSL (Domain-Specific Languages) en R2Lang tenían un bug crítico donde las funciones semánticas solo recibían el primer parámetro en lugar de todos los parámetros especificados en las reglas.

### Síntomas Observados

1. **En el DSL de calculadora**: La función `calcular(num1, op, num2)` solo recibía `num1`, faltando `op` y `num2`
2. **En el DSL de comando simple**: La función `ejecutar_comando(hacer, bebida)` solo recibía `hacer`, faltando `bebida`
3. **Error generado**: `panic: parseInt: arg should be string` porque se intentaba parsear objetos complejos como strings

### Ejemplo del Bug

```r2
// DSL de calculadora con bug
dsl Calculadora {
    rule("operacion", ["NUMERO", "operador", "NUMERO"], "calcular")
    
    func calcular(num1, op, num2) {
        // Solo recibía num1, op y num2 eran undefined
        var n1 = std.parseInt(num1)  // ✅ Funcionaba
        var n2 = std.parseInt(num2)  // ❌ Fallaba - num2 era undefined
        return n1 + n2
    }
}
```

## Causa Raíz

El problema estaba en el método `extractRule` en `pkg/r2core/dsl_definition.go`. El código interpretaba incorrectamente el array de elementos de una regla como alternativas separadas en lugar de una secuencia única.

### Código Problemático

```go
func (dsl *DSLDefinition) extractRule(call *CallExpression) {
    // ...
    for _, alt := range alternatives.Elements {
        if altStr, ok := alt.(*StringLiteral); ok {
            altStrings = append(altStrings, strings.Trim(altStr.Value, "\"'"))
        }
    }
    // ❌ Problema: Trataba cada elemento como alternativa separada
    dsl.Grammar.AddRule(ruleName, altStrings, action)
}
```

Este código convertía `["NUMERO", "operador", "NUMERO"]` en tres alternativas separadas:
- Alternativa 1: `NUMERO`
- Alternativa 2: `operador`  
- Alternativa 3: `NUMERO`

En lugar de una secuencia: `NUMERO operador NUMERO`

## Solución Implementada

### 1. Corrección del Método `extractRule`

```go
func (dsl *DSLDefinition) extractRule(call *CallExpression) {
    // ...
    // ✅ Solución: Unir elementos en una sola secuencia
    sequence := strings.Join(altStrings, " ")
    dsl.Grammar.AddRule(ruleName, []string{sequence}, action)
}
```

### 2. Mejora del Método `AddRule`

```go
func (g *DSLGrammar) AddRule(name string, alternatives []string, action string) {
    rule, exists := g.Rules[name]
    if !exists {
        rule = &DSLRule{
            Name:         name,
            Alternatives: []*DSLAlternative{},
        }
        g.Rules[name] = rule
        if g.StartRule == "" {
            g.StartRule = name
        }
    }

    // ✅ Mejora: Soporte para múltiples reglas con el mismo nombre
    if len(alternatives) == 1 {
        sequence := strings.Fields(alternatives[0])
        rule.Alternatives = append(rule.Alternatives, &DSLAlternative{
            Sequence: sequence,
            Action:   action,
        })
    } else {
        for _, alt := range alternatives {
            sequence := strings.Fields(alt)
            rule.Alternatives = append(rule.Alternatives, &DSLAlternative{
                Sequence: sequence,
                Action:   action,
            })
        }
    }
}
```

### 3. Corrección del Regex de Tokens

En el ejemplo de calculadora, también se corrigió el regex de números:

```r2
// ❌ Antes: coincidía con cadenas vacías
token("NUMERO", "[0-9]*")

// ✅ Después: requiere al menos un dígito
token("NUMERO", "[0-9]+")
```

## Archivos Modificados

1. **`pkg/r2core/dsl_definition.go`**:
   - Línea 97-100: Corrección del método `extractRule`
   - Línea 72: Eliminación de debug innecesario

2. **`pkg/r2core/dsl_grammar.go`**:
   - Línea 61-91: Mejora del método `AddRule`
   - Eliminación de debug temporal

3. **`examples/dsl/calculadora_dsl.r2`**:
   - Línea 3: Corrección del regex de números

## Tests Unitarios Agregados

Se creó `pkg/r2core/dsl_test.go` con 7 test cases para prevenir regresiones:

1. **`TestDSLBasicFunctionality`**: Verificar funcionalidad básica del DSL
2. **`TestDSLParameterPassing`**: Verificar paso correcto de múltiples parámetros
3. **`TestDSLTokenization`**: Verificar tokenización correcta
4. **`TestDSLRuleSequenceParsing`**: Verificar parsing de secuencias (el bug principal)
5. **`TestDSLMultipleRuleAlternatives`**: Verificar múltiples alternativas
6. **`TestDSLRegexTokens`**: Verificar validación de regex
7. **`TestDSLErrorHandling`**: Verificar manejo de errores

## Verificación del Fix

### Antes del Fix
```bash
go run main.go examples/dsl/calculadora_dsl.r2
# Result: panic: parseInt: arg should be string
```

### Después del Fix
```bash
go run main.go examples/dsl/calculadora_dsl.r2
# Result: 
# === DSL Calculadora ===
# Calculando: 5 + 3
# 5 + 3 = (procesado por DSL)
# Resultado DSL: 8
# ✅ Calculadora DSL funcionando!
```

### Tests Unitarios
```bash
go test ./pkg/r2core/ -v -run TestDSL
# Result: PASS (todos los tests pasan)
```

## Impacto

### ✅ Beneficios
- **Funcionalidad DSL completa**: Ahora los DSL pueden usar todos los parámetros especificados
- **Mejor experiencia de usuario**: Los ejemplos de DSL funcionan correctamente
- **Prevención de regresiones**: Tests unitarios completos
- **Mejor mantenibilidad**: Código más limpio y documentado

### ⚠️ Consideraciones
- **Compatibilidad**: Los DSL existentes que dependían del comportamiento incorrecto pueden necesitar ajustes
- **Performance**: Impacto mínimo en performance, mejora la funcionalidad sin degradar velocidad

## Lecciones Aprendidas

1. **Importancia del Testing**: Este bug existía porque no había tests unitarios para la funcionalidad DSL
2. **Debugging Estratégico**: El uso de logs de debug ayudó a identificar rápidamente el problema
3. **Documentación**: Tener ejemplos funcionales es crucial para detectar bugs en features complejas

## Próximos Pasos

1. **Revisar DSL existentes**: Verificar si otros DSL en el proyecto necesitan ajustes
2. **Documentación**: Actualizar la documentación del DSL con ejemplos funcionales
3. **Performance**: Considerar optimizaciones adicionales para el parser DSL
4. **Features**: Implementar características adicionales como validación de tipos en DSL

## Autor
Claude Code - 2025-01-18

---

**Nota**: Este fix resuelve completamente el problema reportado de parámetros múltiples en DSL y establece una base sólida para futuras mejoras en la funcionalidad DSL de R2Lang.