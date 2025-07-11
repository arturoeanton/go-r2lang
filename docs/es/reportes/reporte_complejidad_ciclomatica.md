# Análisis de Complejidad Ciclomática y Plan de Refactorización

## Introducción

Este documento se genera a partir del análisis de la herramienta `gocyclo`, que calcula la complejidad ciclomática del código Go. Una complejidad alta (generalmente > 15) indica que una función tiene demasiados caminos de ejecución (bucles, condicionales), lo que la hace:

- **Difícil de leer y entender.**
- **Difícil de probar**, ya que se necesitan muchos casos de test para cubrir todos los caminos.
- **Propensa a errores** al ser modificada.

El objetivo de este reporte es identificar las funciones problemáticas y proponer un plan de acción claro para refactorizarlas, evaluando el riesgo, la complejidad y el beneficio de cada acción.

## Resumen de Funciones con Alta Complejidad

Las funciones problemáticas se pueden agrupar en dos categorías principales:

1.  **Funciones de Registro de Librerías (`Register...`)**: La gran mayoría de los problemas están aquí. Estas funciones son enormes porque registran todas las funciones nativas de una librería en un solo lugar, usando grandes bloques de `if-else` o `switch`.
2.  **Lógica del Núcleo del Intérprete (`Eval`, `NextToken`, `parse...`)**: Estas son funciones críticas del sistema. Su alta complejidad es más riesgosa porque afecta directamente al comportamiento del lenguaje.

## Plan de Refactorización Detallado

### Categoría 1: Funciones de Registro de Librerías

Estas funciones son el objetivo perfecto para una refactorización de bajo riesgo y alto impacto. La estrategia es la misma para todas ellas.

**Estrategia General:** Aplicar un **patrón data-driven**. En lugar de tener una función gigante con toda la lógica, definiremos las funciones de la librería en una estructura de datos (un `slice` de `structs`, por ejemplo). Una función simple y genérica recorrerá esta estructura para registrar cada función.

| Archivo y Función | Complejidad | Estrategia de Refactorización (Cómo Mejorar) | Riesgo | Complejidad de Arreglo | Beneficio (Qué Ganamos) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `pkg/r2libs/r2hack.go` - `RegisterHack()` | 74 | Crear un `slice` de `structs` donde cada `struct` contenga el nombre de la función (`"print"`, `"println"`, etc.) y la referencia a la función Go que la implementa. Un bucle `for` registrará todo. | **Bajo** | **Bajo** | **Enorme**. Reducimos la complejidad de 74 a ~2. Añadir/quitar funciones será trivial y a prueba de errores. |
| `pkg/r2libs/r2print.go` - `RegisterPrint()` | 60 | Misma estrategia. | **Bajo** | **Bajo** | **Enorme**. Simplificación masiva y mantenibilidad. |
| `pkg/r2libs/r2string.go` - `RegisterString()` | 43 | Misma estrategia. | **Bajo** | **Bajo** | **Muy Alto**. Código más limpio y fácil de extender. |
| `pkg/r2libs/r2os.go` - `RegisterOS()` | 40 | Misma estrategia. | **Bajo** | **Bajo** | **Muy Alto**. |
| `pkg/r2libs/r2http.go` - `RegisterHTTP()` | 37 | Misma estrategia. | **Bajo** | **Bajo** | **Muy Alto**. |
| `pkg/r2libs/r2io.go` - `RegisterIO()` | 35 | Misma estrategia. | **Bajo** | **Bajo** | **Muy Alto**. |
| `pkg/r2libs/r2httpclient.go` - `RegisterHTTPClient()` | 34 | Misma estrategia. | **Bajo** | **Bajo** | **Muy Alto**. |
| `pkg/r2libs/r2go.go` - `RegisterGoInterOp()` | 33 | Misma estrategia. | **Bajo** | **Bajo** | **Muy Alto**. |
| `pkg/r2libs/r2std.go` - `RegisterStd()` | 29 | Misma estrategia. | **Bajo** | **Bajo** | **Alto**. |
| `pkg/r2libs/r2test.go` - `RegisterTest()` | 24 | Misma estrategia. | **Bajo** | **Bajo** | **Alto**. |
| `pkg/r2libs/r2math.go` - `RegisterMath()` | 19 | Misma estrategia. | **Bajo** | **Bajo** | **Alto**. |

### Categoría 2: Lógica del Núcleo del Intérprete

Estas funciones son más delicadas. La refactorización aquí requiere más cuidado y es crucial tener una buena cobertura de tests (como los tests de integración que discutimos antes).

**Estrategia General:** **Dividir y conquistar**. Extraer bloques de lógica cohesiva a funciones auxiliares más pequeñas y con un propósito único.

| Archivo y Función | Complejidad | Estrategia de Refactorización (Cómo Mejorar) | Riesgo | Complejidad de Arreglo | Beneficio (Qué Ganamos) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `pkg/r2core/access_expression.go` - `(*AccessExpression).Eval()` | 90 | Esta función es un monstruo. Probablemente contiene un `switch` gigante sobre tipos. Se debe extraer la lógica para cada tipo de acceso a una función auxiliar. Ej: `evalMemberAccess()`, `evalMethodCall()`, `evalArrayAccessOnObject()`. | **Alto** | **Alto** | **Crítico**. Es la función más compleja del proyecto. Simplificarla es la mayor ganancia de calidad posible. Hará que añadir nuevas propiedades a los objetos sea mucho más seguro. |
| `pkg/r2core/lexer.go` - `(*Lexer).NextToken()` | 60 | Aunque ya la hemos arreglado, sigue siendo compleja. Se pueden extraer lógicas específicas a funciones como `parseSymbolToken()`, `parseNumberToken()`, etc. La función principal se convertiría en un despachador que llama a estas funciones auxiliares. | **Medio** | **Medio** | **Muy Alto**. Un lexer más simple es más fácil de extender con nueva sintaxis y menos propenso a bugs sutiles. |
| `pkg/r2core/object_declaration.go` - `(*ObjectDeclaration).Eval()` | 17 | Extraer la lógica de inicialización de miembros y la de vinculación de herencia a funciones separadas. | **Medio** | **Medio** | **Alto**. Clarifica el proceso de instanciación de objetos, facilitando la depuración de la creación de clases. |
| `pkg/r2core/for_statement.go` - `(*ForStatement).Eval()` | 18 | Esta función probablemente maneja tanto el `for` clásico como el `for-in`. Se puede dividir en `evalCStyleFor()` y `evalForIn()`, con la función principal decidiendo a cuál llamar. | **Medio** | **Medio** | **Alto**. Simplifica la lógica de uno de los bucles más importantes del lenguaje. |
| `pkg/r2core/parse.go` - `(*Parser).parseObjectDeclaration()` | 18 | Extraer la lógica de parseo de la herencia (`extends ...`) a una función auxiliar `parseInheritanceClause()`. | **Bajo** | **Bajo** | **Medio**. Mejora la legibilidad del parser. |
| `pkg/r2core/parse.go` - `(*Parser).parseForStatement()` | 17 | Similar a su `Eval`, dividir la lógica de parseo para los diferentes tipos de `for` en funciones auxiliares. | **Bajo** | **Bajo** | **Medio**. |
| `pkg/r2repl/r2repl.go` - `Repl()` | 16 | El bucle principal del REPL probablemente tiene un `switch` o `if` para los comandos. Extraer cada comando a su propia función y usar un `map[string]func()` para despachar los comandos (Command Pattern). | **Bajo** | **Bajo** | **Medio**. Hace que añadir nuevos comandos al REPL sea mucho más limpio y sencillo. |

## Recomendaciones y Plan de Acción

1.  **Empezar por lo fácil y de alto impacto:** Refactorizar **todas las funciones `Register...`** de la Categoría 1. Esto reducirá drásticamente la complejidad reportada con un riesgo mínimo y nos dará una victoria rápida.
2.  **Fortalecer los Tests:** Antes de tocar la Categoría 2, **implementar los tests de integración ("Golden Tests")** que discutimos. Necesitamos una red de seguridad robusta.
3.  **Proceder con Cuidado:** Abordar las funciones de la Categoría 2 una por una, empezando por las menos complejas (ej. `parseObjectDeclaration`) para ir ganando confianza antes de atacar a los gigantes como `AccessExpression.Eval`.

Siguiendo este plan, podemos mejorar sistemáticamente la calidad del código, su mantenibilidad y nuestra capacidad para extender el lenguaje en el futuro de forma segura.