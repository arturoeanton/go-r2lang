# Ejemplos DSL R2Lang

Esta carpeta contiene ejemplos de **Domain-Specific Languages (DSL)** implementados en R2Lang.

## Ejemplos Disponibles

### ✅ Básicos (Funcionan correctamente)
- **basico.r2**: Ejemplo básico de DSL que procesa palabras simples
- **ejemplo_simple.r2**: DSL para sumar números con sintaxis personalizada
- **primer_dsl.r2**: Tu primer DSL siguiendo el tutorial del manual

### ✅ Funcionales (Funcionan correctamente)
- **dsl_funcional.r2**: DSL para procesamiento de comandos de usuario
- **comando_simple.r2**: DSL para comandos de preparación de bebidas
- **demo_completo.r2**: Demostración completa con múltiples tipos de comandos

### ✅ Avanzados (Funcionan correctamente)
- **calculadora_dsl.r2**: DSL para operaciones matemáticas básicas
- **reglas_negocio_dsl.r2**: DSL para definir reglas de negocio
- **consultas_dsl.r2**: DSL para consultas tipo SQL

### ⚠️ Experimentales (Pueden tener problemas)
- **configuracion_dsl.r2**: DSL para configuración de sistemas (puede colgarse)
- **testing_dsl.r2**: DSL para pruebas BDD (puede colgarse)

## Cómo Ejecutar

```bash
# Ejecutar cualquier ejemplo
go run main.go examples/dsl/ejemplo_simple.r2

# Ejecutar desde la raíz del proyecto
go run main.go examples/dsl/basico.r2
```

## Interpretación de Resultados

Los DSL devuelven un objeto DSLResult con el formato:
```
&{{resultado_procesado} codigo_original {resultado_procesado}}
```

### Ejemplo
```r2
var resultado = MiDSL.use("comando test")
console.log(resultado)
// Salida: &{{Procesado: comando} comando test {Procesado: comando}}
```

## Recursos Adicionales

- **Manual completo**: Ver `docs/es/manual_dsl_v1.md`
- **Documentación**: Consultar la carpeta `docs/`
- **Más ejemplos**: Explorar los archivos `.r2` en esta carpeta

## Características de los DSL

✅ **Sintaxis personalizada**: Define tu propia gramática  
✅ **Tokens y reglas**: Sistema completo de parsing  
✅ **Acciones semánticas**: Funciones para procesar resultados  
✅ **Integración nativa**: Funciona directamente en R2Lang  
✅ **Flexibilidad**: Desde simples hasta complejos  

## Casos de Uso

- **Configuración**: DSL para archivos config
- **Consultas**: DSL para bases de datos
- **Reglas de negocio**: DSL para lógica empresarial
- **Testing**: DSL para pruebas BDD
- **Comandos**: DSL para interfaces de usuario

¡Explora los ejemplos y crea tus propios DSL!