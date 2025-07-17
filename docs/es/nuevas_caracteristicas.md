# Nuevas Características y Mejoras en R2Lang (r2libs)

Este documento resume las nuevas funciones y mejoras implementadas en las librerías estándar de R2Lang (`r2libs`), con el objetivo de aumentar su madurez, cobertura de tests y funcionalidad.

## Módulos Mejorados:

### 1. `collections`

El módulo `collections` ha sido significativamente expandido para incluir funciones de manipulación de arrays de alto nivel, inspiradas en las colecciones de JavaScript y Python.

**Nuevas Funciones:**

*   `map(array, funcion)`: Aplica una función a cada elemento de un array y devuelve un nuevo array con los resultados.
*   `filter(array, funcion)`: Filtra un array, devolviendo un nuevo array solo con los elementos para los que la función de callback devuelve `true`.
*   `reduce(array, funcion, acumulador_inicial)`: Reduce un array a un único valor aplicando una función acumuladora.
*   `sort(array, [funcion_comparacion])`: Ordena un array. Si se provee una función de comparación, la usa; de lo contrario, usa el ordenamiento estándar (numérico, luego alfabético).
*   `find(array, funcion)`: Devuelve el primer elemento del array que satisface la función de prueba.
*   `contains(array, valor)`: Devuelve `true` si el array contiene el valor especificado.

### 2. `math`

El módulo `math` ha sido enriquecido con constantes y funciones matemáticas adicionales para un conjunto más completo de operaciones.

**Nuevas Constantes:**

*   `PI`: El valor de Pi (π).
*   `E`: El valor de la constante de Euler (e).

**Nuevas Funciones:**

*   `sin(x)`: Seno de `x` (en radianes).
*   `cos(x)`: Coseno de `x` (en radianes).
*   `tan(x)`: Tangente de `x` (en radianes).
*   `sqrt(x)`: Raíz cuadrada de `x`.
*   `abs(x)`: Valor absoluto de `x`.
*   `log(x)`: Logaritmo natural de `x`.
*   `log10(x)`: Logaritmo en base 10 de `x`.
*   `pow(base, exp)`: `base` elevado a la potencia `exp`.

### 3. `io`

El módulo `io` ha sido refactorizado y expandido para ofrecer un manejo de ficheros y directorios más intuitivo y completo. Se han renombrado algunas funciones para mayor consistencia.

**Funciones Renombradas y Mejoradas:**

*   `rmFile(path)`: Anteriormente `removeFile`. Elimina un fichero.
*   `rmDir(path)`: Nueva función. Elimina un directorio y todo su contenido de forma recursiva.
*   `mkdir(path)`: Anteriormente `makeDir`. Crea un directorio.
*   `mkdirAll(path)`: Anteriormente `makeDirs`. Crea un directorio y cualquier directorio padre necesario.
*   `listdir(path)`: Anteriormente `readDir`. Lista el contenido de un directorio.
*   `exists(path)`: Anteriormente `fileExists`. Verifica si una ruta existe (fichero o directorio).

**Nuevas Funciones:**

*   `isdir(path)`: Verifica si una ruta dada es un directorio.
*   `isfile(path)`: Verifica si una ruta dada es un fichero.

### 4. `std`

El módulo `std` ha incorporado utilidades generales para la manipulación de datos y la verificación de tipos.

**Nuevas Funciones:**

*   `deepCopy(valor)`: Realiza una copia profunda (recursiva) de un objeto o array, asegurando que no se compartan referencias internas.
*   `is(valor, tipo)`: Una función de aserción de tipo más robusta. Devuelve `true` si `valor` es del `tipo` especificado (ej. "number", "string", "array", "map", "function", "nil", "date", "duration").

---

## Cobertura de Tests:

Se han añadido tests exhaustivos para todas las nuevas funciones y se han refactorizado los tests existentes para los módulos `collections`, `math`, `io` y `std`. Esto ha resultado en un aumento significativo de la cobertura de código, garantizando la fiabilidad y el correcto funcionamiento de las nuevas características.

**Cobertura Actual de `r2libs`:** 34.1% de las sentencias.

---

## Próximos Pasos Sugeridos:

Aunque la cobertura ha mejorado, aún hay áreas que requieren atención. Se recomienda continuar el trabajo en los módulos con baja cobertura y considerar la implementación de las librerías sugeridas previamente (JSON, CSV, ORM, etc.) para seguir enriqueciendo el ecosistema de R2Lang.
