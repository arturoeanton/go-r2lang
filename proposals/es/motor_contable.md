# Propuesta: Evolucionar R2Lang a un Lenguaje de Motor Contable

**Fecha:** 2024-07-10
**Autor:** [Tu Nombre/Alias]
**Estado:** Borrador

## 1. Resumen Ejecutivo

Esta propuesta detalla la visión para transformar R2Lang en un **Lenguaje de Motor Contable (Accounting Engine Language)**. El objetivo es crear un lenguaje de dominio específico (DSL) que permita a contadores, auditores y desarrolladores definir y procesar transacciones financieras complejas utilizando los principios de la **contabilidad de doble entrada**. 

El lenguaje garantizará que todas las operaciones estén balanceadas (débitos = créditos), generará automáticamente asientos contables y permitirá la creación de informes financieros (Balances, Estados de Resultados) de forma nativa. Esto reduce drásticamente el riesgo de errores manuales, simplifica las auditorías y hace que la lógica financiera sea transparente y verificable.

## 2. Idea y Motivación

La contabilidad es la base de cualquier negocio, pero su implementación en software es a menudo propensa a errores. Los sistemas tradicionales pueden sufrir de:

- **Desbalanceo de cuentas:** Errores de software que resultan en asientos contables que no cuadran.
- **Lógica opaca:** Los cálculos financieros están ocultos en código imperativo, difícil de auditar.
- **Complejidad en informes:** Generar informes financieros requiere consultas complejas y propensas a errores.
- **Falta de inmutabilidad:** Los registros de transacciones a menudo pueden ser modificados, comprometiendo la integridad de la auditoría.

**La solución propuesta** es un lenguaje donde las transacciones no son simples inserciones en una base de datos, sino construcciones de primera clase que el motor de R2Lang entiende y valida. El lenguaje forzará la disciplina contable a nivel de sintaxis, haciendo imposible registrar una transacción desbalanceada.

## 3. Visión del Producto Final

Imaginemos una empresa que registra una venta a crédito. En lugar de múltiples llamadas a una API, el evento de negocio se describe en un único script `venta.r2`:

```r2
// --- Transacción de Venta a Crédito ---

// Definición de las cuentas involucradas (previamente configuradas)
Accounts { 
    Activos.CuentasPorCobrar,
    Ingresos.Ventas,
    Costos.CostoDeVenta,
    Activos.Inventario
}

// Definición de la transacción
Transaction "Venta de Producto A a Cliente B" {
    Date: "2024-07-10",
    Description: "Venta de 10 unidades del Producto A",
    Metadata: {
        cliente_id: "CLI-001",
        producto_id: "PROD-A",
        factura_id: "FACT-2024-105"
    }

    // Asientos contables de la transacción
    Entries {
        // Registro del ingreso por la venta
        Debit(Activos.CuentasPorCobrar, 1500.00, "Aumenta lo que el cliente nos debe"),
        Credit(Ingresos.Ventas, 1500.00, "Registra el ingreso por la venta"),

        // Registro del costo del producto vendido
        Debit(Costos.CostoDeVenta, 600.00, "Registra el costo de la mercancía vendida"),
        Credit(Activos.Inventario, 600.00, "Disminuye el valor del inventario")
    }
}
```

El motor de R2Lang se encargaría de:

1.  **Validar la transacción:** Verificar que la suma de débitos sea igual a la suma de créditos.
2.  **Atomicidad:** Asegurar que todos los asientos se registren o ninguno lo haga.
3.  **Generar el Libro Diario:** Crear los registros contables correspondientes en un formato inmutable.
4.  **Actualizar el Libro Mayor:** Actualizar los saldos de las cuentas `CuentasPorCobrar`, `Ventas`, etc.

Posteriormente, se podrían generar informes de forma nativa:

```r2
// Generar un Balance General
Report BalanceSheet {
    Title: "Balance General al 31 de Julio de 2024",
    Date: "2024-07-31",
    Sections: ["Activos", "Pasivos", "Patrimonio"]
}
```

## 4. Épicas y Tareas

### Épica 1: Sintaxis Contable de Primera Clase

*   **Descripción:** Introducir conceptos contables como `Transaction`, `Debit`, `Credit` y `Accounts` como elementos nativos del lenguaje.
*   **Prioridad:** Alta
*   **Complejidad:** Alta

| Tarea                                       | Descripción                                                                                             | Estimación | Complejidad |
| ------------------------------------------- | ------------------------------------------------------------------------------------------------------- | ---------- | ----------- |
| **T1.1: Implementar `Transaction` block**   | Añadir al parser la capacidad de procesar bloques `Transaction` con metadatos.                          | 3-5 días   | Media       |
| **T1.2: Implementar `Debit` y `Credit`**    | Crear funciones/palabras clave `Debit(Cuenta, Monto)` y `Credit(Cuenta, Monto)`.                       | 2-4 días   | Media       |
| **T1.3: Sistema de Cuentas (`Accounts`)**   | Diseñar una forma de definir y referenciar un plan de cuentas jerárquico (ej: `Activos.Corrientes.Caja`). | 4-7 días   | Alta        |
| **T1.4: Validador de Doble Entrada**        | El motor debe verificar automáticamente que `sum(Debits) == sum(Credits)` para cada transacción.        | 2-3 días   | Media       |

### Épica 2: Motor de Procesamiento Contable

*   **Descripción:** Construir la lógica del evaluador para que procese transacciones, mantenga el estado del libro mayor y garantice la integridad de los datos.
*   **Prioridad:** Alta
*   **Complejidad:** Alta

| Tarea                                         | Descripción                                                                                                   | Estimación | Complejidad |
| --------------------------------------------- | ------------------------------------------------------------------------------------------------------------- | ---------- | ----------- |
| **T2.1: Implementar el Libro Diario**         | Crear una estructura de datos inmutable (append-only) para almacenar el historial de todas las transacciones. | 4-6 días   | Alta        |
| **T2.2: Implementar el Libro Mayor**          | Mantener los saldos actualizados de todas las cuentas en una estructura de datos eficiente.                   | 4-6 días   | Alta        |
| **T2.3: Atomicidad de Transacciones**         | Asegurar que si una parte de la transacción falla (ej: validación), toda la transacción se revierte.        | 3-5 días   | Media       |
| **T2.4: Soporte para Múltiples Monedas**      | Añadir la capacidad de manejar y convertir entre diferentes monedas, gestionando tasas de cambio.           | 5-8 días   | Alta        |

### Épica 3: Generación de Informes Nativos

*   **Descripción:** Crear una sintaxis declarativa para generar informes financieros estándar directamente desde el libro mayor.
*   **Prioridad:** Media
*   **Complejidad:** Alta

| Tarea                                       | Descripción                                                                                             | Estimación | Complejidad |
| ------------------------------------------- | ------------------------------------------------------------------------------------------------------- | ---------- | ----------- |
| **T3.1: Implementar `Report` block**        | Crear una sintaxis para definir informes como `Report BalanceSheet { ... }`.                            | 3-5 días   | Media       |
| **T3.2: Lógica de Balance General**         | Implementar la lógica para agregar los saldos de las cuentas de Activos, Pasivos y Patrimonio.          | 2-4 días   | Media       |
| **T3.3: Lógica de Estado de Resultados**    | Implementar la lógica para calcular ingresos, costos, gastos y la utilidad neta.                        | 2-4 días   | Media       |
| **T3.4: Exportación de Informes**           | Permitir la exportación de los informes generados a formatos como JSON, CSV o texto plano.              | 1-3 días   | Baja        |

## 5. Idea de Implementación Mínima (MVP)

El MVP se enfocaría en validar el concepto central de transacciones balanceadas:

1.  **Sintaxis `Transaction`:** Implementar `Transaction`, `Debit`, `Credit`.
2.  **Plan de Cuentas Simple:** Usar strings simples para los nombres de las cuentas (ej: `"Activos:Caja"`).
3.  **Validador de Doble Entrada:** El motor debe rechazar cualquier transacción donde los débitos no igualen a los créditos.
4.  **Libro Mayor en Memoria:** Un simple mapa en memoria para mantener los saldos de las cuentas.
5.  **Sin persistencia:** El estado del libro mayor se pierde al reiniciar la aplicación.

**Ejemplo de MVP:**

```r2
Transaction "Venta simple" {
    Debit("CuentasPorCobrar", 100),
    Credit("Ventas", 100)
}

// El motor validaría que 100 == 100 y actualizaría los saldos.
// Una transacción como la siguiente fallaría:
// Transaction "Error" { Debit("Gastos", 99), Credit("Caja", 100) }
```

## 6. Casos de Uso

-   **Software Contable (ERP):** Servir como el núcleo contable para sistemas de planificación de recursos empresariales.
-   **Startups FinTech:** Proporcionar una base sólida y auditable para aplicaciones de pagos, préstamos o gestión de activos.
-   **Contabilidad Personal:** Crear aplicaciones simples para gestionar las finanzas personales con rigor contable.
-   **Auditoría y Cumplimiento:** Utilizar los scripts de R2Lang como una "verdad única" para auditorías, ya que son legibles y autovalidados.
-   **Sistemas de E-commerce:** Registrar de forma fiable los ingresos, costos, inventario y pasivos de cada venta.

## 7. Conclusión

Evolucionar R2Lang hacia un motor contable lo posicionaría en un nicho de alto valor donde la precisión, la auditabilidad y la claridad son primordiales. Este enfoque transformaría a R2Lang de una herramienta de propósito general a una solución especializada y crítica para el negocio. Invitamos a la comunidad a debatir y contribuir a esta emocionante dirección para el proyecto.
