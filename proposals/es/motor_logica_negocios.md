# Propuesta: Transformar R2Lang en un Motor de Lógica de Negocios

**Fecha:** 2024-07-10
**Autor:** [Tu Nombre/Alias]
**Estado:** Borrador

## 1. Resumen Ejecutivo

Esta propuesta describe la visión y el plan para evolucionar R2Lang de un lenguaje de scripting de propósito general a un **Motor de Lógica de Negocios (Business Logic Engine)** especializado. El objetivo es posicionar R2Lang como una herramienta ideal para que analistas de negocio, expertos en dominios y desarrolladores puedan definir, ejecutar y auditar reglas de negocio complejas de forma segura y declarativa, separando la lógica de negocio del código de la aplicación principal.

Al enfocarse en este nicho, R2Lang puede ofrecer un valor único en industrias como **FinTech, seguros, logística y e-commerce**, donde las reglas cambian constantemente y la auditoría es crucial.

## 2. Idea y Motivación

La lógica de negocio (ej: cálculos de tarifas, evaluación de riesgos, criterios de elegibilidad) es el corazón de muchas aplicaciones. Sin embargo, a menudo está profundamente acoplada al código de la aplicación (ej: Java, C#, Python). Esto presenta varios problemas:

- **Lentitud para cambiar:** Modificar una regla requiere un ciclo completo de desarrollo, pruebas y despliegue.
- **Falta de claridad:** El código imperativo ofusca la intención de negocio de la regla.
- **Brecha de comunicación:** Los analistas de negocio no pueden leer ni validar el código directamente.
- **Riesgo de auditoría:** Es difícil rastrear por qué se tomó una decisión específica.

**La solución propuesta** es transformar R2Lang en un lenguaje donde la lógica de negocio se escribe de forma declarativa y legible. La aplicación principal (escrita en Go, Java, etc.) simplemente carga y ejecuta estos scripts de R2Lang, pasando datos de entrada y recibiendo un resultado claro y auditable.

## 3. Visión del Producto Final

Imaginemos una aseguradora que necesita calcular la prima de un seguro de coche. En lugar de codificarlo en Java, lo definen en un archivo `calculo_prima.r2`:

```r2
// --- Motor de Reglas para Primas de Seguro ---

// Definición de datos de entrada con validaciones
Input { 
    edad_conductor: Number(min: 18, max: 99),
    tipo_vehiculo: String(enum: ["Sedan", "SUV", "Deportivo"]),
    historial_accidentes: Number(min: 0)
}

// Definición de la estructura de salida
Output { 
    prima_base: Number,
    ajuste_por_edad: Number,
    ajuste_por_vehiculo: Number,
    ajuste_por_riesgo: Number,
    prima_final: Number,
    decision: String,
    motivos: Array
}

// Reglas de cálculo
Rule "Calcular Prima Base" {
    Given {
        // No se necesita contexto adicional
    }
    When {
        // La prima base es un valor fijo
        prima_base = 500
    }
    Then {
        Output.prima_base = prima_base
        Output.motivos.push("Prima base establecida en " + prima_base)
    }
}

Rule "Ajuste por Edad del Conductor" {
    Given {
        edad = Input.edad_conductor
    }
    When {
        ajuste = 0
        if (edad < 25) {
            ajuste = 150
        } else if (edad > 65) {
            ajuste = 100
        }
    }
    Then {
        Output.ajuste_por_edad = ajuste
        Output.motivos.push("Ajuste por edad (" + edad + "): " + ajuste)
    }
}

// ... más reglas ...

// Regla de decisión final
Rule "Calcular Prima Final y Decisión" {
    Given {
        // Se ejecuta después de todas las demás
        depende_de: ["Calcular Prima Base", "Ajuste por Edad del Conductor"]
    }
    When {
        prima_final = Output.prima_base + Output.ajuste_por_edad + ...
    }
    Then {
        Output.prima_final = prima_final
        Output.decision = "Aprobado"
    }
}
```

La aplicación principal simplemente haría:

```go
engine := r2lang.NewEngine()
input := map[string]interface{}{"edad_conductor": 22, ...}
result, auditTrail := engine.Execute("calculo_prima.r2", input)

fmt.Println("Prima final:", result["prima_final"]) // 650
fmt.Println("Auditoría:", auditTrail) // ["Prima base establecida en 500", ...]
```

## 4. Épicas y Tareas

### Épica 1: Sintaxis Declarativa para Reglas

*   **Descripción:** Introducir nuevas palabras clave y estructuras en el lenguaje para definir reglas de negocio de forma explícita y legible.
*   **Prioridad:** Alta
*   **Complejidad:** Alta

| Tarea                                       | Descripción                                                                                             | Estimación | Complejidad |
| ------------------------------------------- | ------------------------------------------------------------------------------------------------------- | ---------- | ----------- |
| **T1.1: Implementar `Rule` block**          | Añadir al parser y al AST la capacidad de reconocer y procesar bloques `Rule "Nombre" { ... }`.         | 3-5 días   | Media       |
| **T1.2: Implementar `Given/When/Then`**     | Soportar los bloques `Given`, `When`, y `Then` dentro de una `Rule`.                                      | 2-3 días   | Media       |
| **T1.3: Implementar `Input` y `Output`**    | Crear bloques para definir los esquemas de datos de entrada y salida, con validaciones básicas.         | 4-6 días   | Alta        |
| **T1.4: Control de dependencias de reglas** | Añadir una cláusula `depende_de` para especificar el orden de ejecución de las reglas.                   | 2-4 días   | Alta        |

### Épica 2: Motor de Ejecución y Auditoría

*   **Descripción:** Modificar el evaluador de R2Lang para que pueda ejecutar un conjunto de reglas, gestionar el estado y generar un rastro de auditoría detallado.
*   **Prioridad:** Alta
*   **Complejidad:** Alta

| Tarea                                         | Descripción                                                                                                   | Estimación | Complejidad |
| --------------------------------------------- | ------------------------------------------------------------------------------------------------------------- | ---------- | ----------- |
| **T2.1: Crear un `RuleEngine`**               | Desarrollar un nuevo modo de ejecución que orqueste la evaluación de las reglas en lugar de un script lineal. | 5-8 días   | Alta        |
| **T2.2: Implementar el `Audit Trail`**        | Recolectar información durante la ejecución (qué regla se activó, qué valores se usaron) en un log estructurado. | 3-5 días   | Media       |
| **T2.3: Validación de `Input` y `Output`**    | Validar los datos de entrada contra el esquema `Input` y asegurar que la salida cumpla con el esquema `Output`. | 2-3 días   | Media       |
| **T2.4: Grafo de dependencias**               | Construir un grafo para resolver el orden de ejecución de las reglas basado en la cláusula `depende_de`.      | 3-5 días   | Alta        |

### Épica 3: Sandboxing y Seguridad

*   **Descripción:** Asegurar que los scripts de reglas se ejecuten en un entorno aislado y seguro, sin acceso a recursos no autorizados.
*   **Prioridad:** Media
*   **Complejidad:** Media

| Tarea                                       | Descripción                                                                                             | Estimación | Complejidad |
| ------------------------------------------- | ------------------------------------------------------------------------------------------------------- | ---------- | ----------- |
| **T3.1: Limitar acceso a bibliotecas**      | Por defecto, deshabilitar bibliotecas como `io` y `http` en el modo de motor de reglas.                 | 1-2 días   | Baja        |
| **T3.2: Implementar límites de recursos**   | Añadir límites en el tiempo de ejecución, uso de memoria y profundidad de la pila para evitar abusos. | 3-5 días   | Media       |
| **T3.3: Inyección de funciones seguras**    | Permitir que la aplicación anfitriona inyecte funciones seguras (ej: `lookupCustomerData`) en el scope. | 2-4 días   | Media       |

## 5. Idea de Implementación Mínima (MVP)

El MVP se centraría en la funcionalidad principal para demostrar el valor del concepto:

1.  **Sintaxis `Rule`:** Implementar los bloques `Rule`, `When`, `Then` (sin `Given` ni dependencias complejas).
2.  **Motor Básico:** Un motor que ejecuta todas las reglas en el orden en que aparecen.
3.  **`Input` y `Output` como variables globales:** En lugar de bloques de esquema, usar objetos globales `Input` y `Output`.
4.  **Auditoría Simple:** Un array de strings que registra qué regla se ejecutó.

**Ejemplo de MVP:**

```r2
let Output = { prima: 0, motivos: [] }

Rule "Calcular Prima Base" {
    When {
        Output.prima = 500
    }
    Then {
        Output.motivos.push("Prima base es 500")
    }
}

Rule "Ajuste por Edad" {
    When {
        if (Input.edad < 25) {
            Output.prima = Output.prima + 150
        }
    }
    Then {
        Output.motivos.push("Ajuste por edad aplicado")
    }
}
```

## 6. Casos de Uso

-   **FinTech:** Calcular la elegibilidad para un préstamo, detectar transacciones fraudulentas, determinar tasas de interés.
-   **Seguros:** Calcular primas, procesar reclamaciones, evaluar riesgos de pólizas.
-   **E-commerce:** Lógica de descuentos y promociones, cálculo de costos de envío, reglas de inventario.
-   **Logística:** Determinar la mejor ruta de envío, calcular tarifas, gestionar la asignación de flotas.
-   **Salud:** Triaje de pacientes basado en síntomas, validación de reclamaciones de seguros médicos.

## 7. Conclusión

Transformar R2Lang en un motor de lógica de negocio es una oportunidad estratégica para diferenciar el proyecto y resolver un problema real y valioso en la industria del software. Esta propuesta sienta las bases para un roadmap claro y accionable. Se invita a la comunidad a discutir, refinar y contribuir a esta visión.
