# Fix Completo: Parámetros Múltiples y Resultados del DSL

## Fecha
2025-01-18

## Resumen
Se corrigieron dos problemas críticos en el sistema DSL de R2Lang:
1. **Formato incorrecto de parámetros**: Los parámetros aparecían como `{+}` en lugar de `+`
2. **Resultados no accesibles**: El resultado del DSL no se retornaba correctamente

## Problemas Identificados

### Problema 1: Formato de Parámetros `{+}`
Las funciones DSL recibían parámetros con formato `{+}` en lugar del valor real `+`, rompiendo la experiencia del desarrollador.

**Ejemplo del problema:**
```r2
func calcular(num1, op, num2) {
    // op llegaba como "{+}" en lugar de "+"
    if (op == "+") {  // Nunca se cumplía
        return n1 + n2
    }
}
```

### Problema 2: Resultado DSL No Accesible
El resultado del DSL se retornaba como `&{<nil> 5 + 3 <nil>}` en lugar del valor calculado.

**Ejemplo del problema:**
```r2
var resultado = calc("5 + 3")
console.log(resultado)  // Imprimía: &{<nil> 5 + 3 <nil>}
// No había forma de acceder al resultado real (8)
```

## Soluciones Implementadas

### 1. Corrección del Formato de Parámetros

**Archivo**: `pkg/r2core/dsl_definition.go`

```go
// Bind parameters to arguments
for i, param := range fn.Args {
    if i < len(args) {
        // If the argument is a ReturnValue, extract its value
        var argValue interface{}
        if retVal, ok := args[i].(*ReturnValue); ok {
            argValue = retVal.Value
        } else if retVal, ok := args[i].(ReturnValue); ok {
            argValue = retVal.Value
        } else {
            argValue = args[i]
        }
        fnEnv.Set(param, argValue)
    } else {
        fnEnv.Set(param, nil)
    }
}
```

**Archivo**: `pkg/r2core/dsl_grammar.go`

```go
// Apply semantic action if available
if alt.Action != "" {
    if action, exists := p.Grammar.Actions[alt.Action]; exists {
        result := action(results)
        // If the result is a ReturnValue, extract its value
        if retVal, ok := result.(*ReturnValue); ok {
            return retVal.Value, nil
        }
        return result, nil
    }
}
```

```go
// Symbol is a rule
result, err := p.parseRule(symbol)
if err != nil {
    return nil, err
}
// If the result is a ReturnValue, extract its value
if retVal, ok := result.(*ReturnValue); ok {
    results = append(results, retVal.Value)
} else {
    results = append(results, result)
}
```

### 2. Mejora del Resultado DSL

**Archivo**: `pkg/r2core/dsl_definition.go`

```go
func (dsl *DSLDefinition) evaluateDSLCode(code string, env *Environment) interface{} {
    // ... parsing logic ...
    
    // Extract the final result from the AST
    var finalResult interface{}
    if retVal, ok := ast.(*ReturnValue); ok {
        finalResult = retVal.Value
    } else if retVal, ok := ast.(ReturnValue); ok {
        finalResult = retVal.Value
    } else {
        finalResult = ast
    }

    // Return the parsed AST with the final result
    return &DSLResult{
        AST:    ast,
        Code:   code,
        Output: finalResult,
    }
}
```

**Archivo**: `pkg/r2core/dsl_grammar.go`

```go
// DSLResult representa el resultado de evaluar código DSL
type DSLResult struct {
    AST    interface{}
    Code   string
    Output interface{}
}

// GetResult returns the final result of the DSL execution
func (r *DSLResult) GetResult() interface{} {
    return r.Output
}

// String returns a string representation of the result
func (r *DSLResult) String() string {
    return fmt.Sprintf("DSL Result: %v", r.Output)
}
```

### 3. Soporte para Acceso a Propiedades DSLResult

**Archivo**: `pkg/r2core/access_expression.go`

```go
case *DSLResult:
    return evalDSLResultAccess(obj, ae.Member, env)
```

```go
func evalDSLResultAccess(dslResult *DSLResult, member string, env *Environment) interface{} {
    switch member {
    case "AST":
        return dslResult.AST
    case "Code":
        return dslResult.Code
    case "Output":
        return dslResult.Output
    case "GetResult":
        return func([]interface{}) interface{} {
            return dslResult.GetResult()
        }
    default:
        return nil
    }
}
```

### 4. Mejora del método String() de ReturnValue

**Archivo**: `pkg/r2core/return_value.go`

```go
import "fmt"

func (rv *ReturnValue) String() string {
    return fmt.Sprintf("%v", rv.Value)
}
```

## Verificación del Fix

### Antes del Fix
```bash
go run main.go examples/dsl/calculadora_dsl.r2
# Result: 
# Calculando: 5 {+} 3  <- Formato incorrecto
# Resultado DSL: &{<nil> 5 + 3 <nil>}  <- Resultado no accesible
```

### Después del Fix
```bash
go run main.go examples/dsl/calculadora_dsl.r2
# Result:
# Calculando: 5 + 3  <- ✅ Formato correcto
# Resultado DSL: DSL Result: 8  <- ✅ Resultado accesible
```

### Acceso a Propiedades del Resultado
```r2
var resultado = calc("5 + 3")
console.log("Resultado completo:", resultado)  // DSL Result: 8
console.log("Solo el resultado:", resultado.Output)  // 8
```

## Archivos Modificados

1. **`pkg/r2core/dsl_definition.go`**:
   - Línea 168-182: Desempacado de `ReturnValue` en argumentos
   - Línea 147-153: Desempacado de `ReturnValue` en resultado final

2. **`pkg/r2core/dsl_grammar.go`**:
   - Línea 221-227: Desempacado de `ReturnValue` en acciones semánticas
   - Línea 214-220: Desempacado de `ReturnValue` en reglas
   - Línea 251-259: Métodos `GetResult()` y `String()` para `DSLResult`

3. **`pkg/r2core/access_expression.go`**:
   - Línea 43-44: Soporte para acceso a propiedades de `DSLResult`
   - Línea 404-419: Función `evalDSLResultAccess`

4. **`pkg/r2core/return_value.go`**:
   - Línea 9-11: Método `String()` para `ReturnValue`

5. **`examples/dsl/calculadora_mejorada.r2`** (nuevo):
   - Ejemplo de uso del DSL con acceso a resultados

## Características Nuevas

### 1. Acceso a Propiedades del Resultado
```r2
var resultado = calc("5 + 3")
console.log(resultado.Output)    // 8
console.log(resultado.Code)      // "5 + 3"
console.log(resultado.AST)       // AST interno
```

### 2. Método GetResult()
```r2
var resultado = calc("5 + 3")
var valor = resultado.GetResult()  // 8
```

### 3. Representación String Mejorada
```r2
console.log(resultado)  // "DSL Result: 8"
```

## Compatibilidad

### ✅ Totalmente Compatible
- Los DSL existentes siguen funcionando
- Los parámetros ahora se pasan correctamente
- Los resultados son accesibles

### ✅ Mejoras sin Breaking Changes
- Acceso opcional a propiedades del resultado
- Método String() mejorado
- Mejor experiencia de desarrollo

## Casos de Prueba

### Test 1: Operadores Correctos
```r2
dsl Test {
    token("A", "a")
    token("B", "b")
    rule("seq", ["A", "B"], "combine")
    func combine(a, b) {
        return a + b  // Ahora funciona correctamente
    }
}
```

### Test 2: Acceso a Resultados
```r2
var result = Test.use("a b")
console.log(result.Output)  // "ab"
```

### Test 3: Múltiples Parámetros
```r2
dsl Calc {
    rule("op", ["NUM", "OP", "NUM"], "calc")
    func calc(n1, op, n2) {
        // Todos los parámetros llegan correctamente
        return std.parseInt(n1) + std.parseInt(n2)
    }
}
```

## Impacto en Performance

- **Overhead mínimo**: Solo se agrega verificación de tipos
- **Mejora en usabilidad**: Los resultados son directamente accesibles
- **Sin regresiones**: Todos los tests existentes pasan

## Próximos Pasos

1. **Documentación**: Actualizar manual del DSL con ejemplos de acceso a resultados
2. **Tests adicionales**: Agregar más casos de prueba para cobertura completa
3. **Features**: Considerar agregar validación de tipos en tiempo de compilación
4. **Performance**: Optimizar el desempacado de ReturnValue si es necesario

## Conclusión

Los fixes implementados resuelven completamente los problemas reportados:

✅ **Problema 1 resuelto**: Los parámetros ahora se pasan con formato correcto (`+` en lugar de `{+}`)

✅ **Problema 2 resuelto**: Los resultados DSL son accesibles y utilizables

✅ **Experiencia mejorada**: Los desarrolladores pueden trabajar con DSL de forma intuitiva

✅ **Compatibilidad mantenida**: No se rompe código existente

El sistema DSL de R2Lang ahora ofrece una experiencia de desarrollo fluida y resultados accesibles, cumpliendo con las expectativas de los usuarios.

## Autor
Claude Code - 2025-01-18