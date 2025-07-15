# R2Test - Framework de Testing Avanzado para R2Lang

R2Test es un potente framework de testing diseñado específicamente para R2Lang, que proporciona capacidades de testing unitario, cobertura de código, mocking, fixtures y reporting avanzado.

## Instalación

```bash
go build -o r2test cmd/r2test/main.go
```

## Uso Básico

### Ejecutar todos los tests
```bash
r2test
```

### Ejecutar tests en un directorio específico
```bash
r2test ./tests
```

### Ejecutar tests con cobertura
```bash
r2test -coverage -verbose ./tests
```

### Ejecutar tests con filtro por patrón
```bash
r2test -grep "Calculator" ./tests
```

### Ejecutar tests en paralelo
```bash
r2test -parallel -workers 4
```

## Opciones de Línea de Comandos

### Opciones Básicas
- `-help` - Mostrar ayuda
- `-version` - Mostrar información de versión
- `-config FILE` - Cargar configuración desde archivo JSON
- `-verbose` - Salida detallada
- `-quiet` - Salida mínima (solo errores)
- `-debug` - Habilitar salida de debug
- `-dry-run` - Mostrar qué se ejecutaría sin ejecutar

### Descubrimiento de Tests
- `-dirs DIRS` - Directorios de tests separados por comas (por defecto: `./tests,./test`)
- `-patterns PATTERNS` - Patrones de archivos de test (por defecto: `*_test.r2,*Test.r2,test_*.r2`)
- `-ignore PATTERNS` - Patrones a ignorar (por defecto: `node_modules,vendor`)

### Ejecución de Tests
- `-timeout DURATION` - Timeout para tests (ej: `30s`, `5m`) (por defecto: `30s`)
- `-parallel` - Ejecutar tests en paralelo
- `-workers N` - Máximo de workers en paralelo (por defecto: `4`)
- `-bail` - Parar en el primer test fallido
- `-retries N` - Número de reintentos para tests fallidos

### Filtrado de Tests
- `-grep PATTERN` - Ejecutar solo tests que coincidan con el patrón
- `-tags TAGS` - Ejecutar solo tests con tags específicos
- `-skip PATTERN` - Saltar tests que coincidan con el patrón
- `-only PATTERN` - Ejecutar solo tests que coincidan con el patrón (exclusivo)

### Opciones de Cobertura
- `-coverage` - Habilitar colección de cobertura
- `-coverage-threshold N` - Umbral de cobertura en porcentaje (por defecto: `80`)
- `-coverage-output DIR` - Directorio de salida para cobertura (por defecto: `./coverage`)
- `-coverage-formats LIST` - Formatos de cobertura: `html,json` (por defecto: `html`)
- `-coverage-exclude LIST` - Patrones a excluir de cobertura

### Opciones de Reporting
- `-reporters LIST` - Reporteros: `console,json,junit` (por defecto: `console`)
- `-output DIR` - Directorio de salida para reportes (por defecto: `./test-results`)

### Opciones Avanzadas
- `-watch` - Modo watch - reejecutar tests cuando cambien archivos
- `-fixtures DIR` - Directorio de fixtures (por defecto: `./fixtures`)
- `-cleanup-mocks` - Limpiar mocks después de tests (por defecto: `true`)
- `-isolation` - Ejecutar tests en contextos aislados

## Archivo de Configuración

Puedes usar un archivo JSON para configurar R2Test:

```json
{
  "testDirs": ["./tests", "./examples/testing"],
  "patterns": ["*_test.r2", "*Test.r2", "test_*.r2"],
  "ignore": ["node_modules", "vendor", "*.tmp.r2"],
  "timeout": "30s",
  "parallel": false,
  "maxWorkers": 4,
  "bail": false,
  "coverage": {
    "enabled": true,
    "threshold": 80,
    "output": "./coverage",
    "formats": ["html", "json"],
    "exclude": ["*_test.r2", "vendor/*"]
  },
  "reporters": ["console", "json"],
  "outputDir": "./test-results",
  "verbose": true,
  "retries": 0,
  "watchMode": false,
  "fixtureDir": "./fixtures"
}
```

## Sintaxis de Tests

R2Test utiliza la sintaxis `describe()` e `it()` para estructurar tests:

```r2
describe("Operaciones Matemáticas", func() {
    it("debe sumar dos números correctamente", func() {
        let resultado = 2 + 3;
        assert.equals(resultado, 5);
    });
    
    it("debe restar números correctamente", func() {
        let resultado = 10 - 4;
        assert.equals(resultado, 6);
    });
});
```

### Assertions Disponibles

#### Comparaciones Básicas
- `assert.equals(actual, expected)` - Comparación de igualdad
- `assert.notEquals(actual, expected)` - Comparación de desigualdad
- `assert.true(value)` - Verificar que el valor es verdadero
- `assert.false(value)` - Verificar que el valor es falso

#### Comparaciones Numéricas
- `assert.greater(a, b)` - a > b
- `assert.greaterOrEqual(a, b)` - a >= b
- `assert.less(a, b)` - a < b
- `assert.lessOrEqual(a, b)` - a <= b

#### Verificaciones de Contenido
- `assert.contains(string, substring)` - Verificar que la cadena contiene la subcadena
- `assert.notContains(string, substring)` - Verificar que la cadena no contiene la subcadena
- `assert.hasLength(collection, length)` - Verificar la longitud de una colección
- `assert.empty(collection)` - Verificar que la colección está vacía
- `assert.notEmpty(collection)` - Verificar que la colección no está vacía

#### Verificaciones de Nulidad
- `assert.nil(value)` - Verificar que el valor es nulo
- `assert.notNil(value)` - Verificar que el valor no es nulo

#### Verificaciones de Errores
- `assert.panics(func)` - Verificar que la función lanza un error
- `assert.notPanics(func)` - Verificar que la función no lanza un error

## Fixtures

R2Test proporciona soporte para fixtures que pueden ser utilizadas en tests:

```r2
describe("Tests con Fixtures", func() {
    it("debe cargar datos de fixture", func() {
        let datos = fixture.load("usuarios.json");
        assert.notEmpty(datos);
        assert.hasLength(datos, 3);
    });
});
```

## Mocking

R2Test incluye capacidades de mocking para aislar código bajo test:

```r2
describe("Tests con Mocks", func() {
    it("debe mockear una función", func() {
        let mockFunc = mock.create();
        mock.returns(mockFunc, "resultado mockeado");
        
        let resultado = mockFunc();
        assert.equals(resultado, "resultado mockeado");
    });
});
```

## Cobertura de Código

Para generar reportes de cobertura:

```bash
r2test -coverage -coverage-formats html,json -coverage-output ./coverage
```

Esto generará:
- `./coverage/html/index.html` - Reporte HTML interactivo
- `./coverage/coverage.json` - Datos de cobertura en JSON

## Reportes

R2Test puede generar varios tipos de reportes:

### Reporte JSON
```bash
r2test -reporters json -output ./reports
```

### Reporte JUnit XML
```bash
r2test -reporters junit -output ./reports
```

### Múltiples Reportes
```bash
r2test -reporters console,json,junit -output ./reports
```

## Ejemplos de Uso

### Desarrollo Local
```bash
# Ejecutar tests con salida detallada
r2test -verbose

# Ejecutar solo tests específicos
r2test -grep "Calculator"

# Ejecutar tests con cobertura
r2test -coverage -coverage-threshold 90
```

### Integración Continua
```bash
# Ejecutar tests con reporte JUnit para CI
r2test -reporters junit -output ./test-results -bail

# Ejecutar tests con cobertura para CI
r2test -coverage -coverage-threshold 80 -reporters json,junit -output ./reports
```

### Desarrollo con Watch Mode
```bash
# Reejecutar tests cuando cambien archivos
r2test -watch -verbose
```

## Estructura de Archivos de Test

Los archivos de test deben:
- Terminar en `*_test.r2`, `*Test.r2`, o `test_*.r2`
- Estar ubicados en directorios configurados (por defecto: `./tests`, `./test`)
- Usar las funciones `describe()` e `it()` para estructurar tests

```
proyecto/
├── tests/
│   ├── calculadora_test.r2
│   ├── utils_test.r2
│   └── fixtures/
│       └── datos_test.json
├── src/
│   ├── calculadora.r2
│   └── utils.r2
└── r2test.config.json
```

## Integración con Editores

R2Test es compatible con editores que soportan:
- Formatos de reporte JUnit XML
- Reportes JSON estructurados
- Salida de consola estándar

## Troubleshooting

### Tests No Encontrados
- Verificar que los archivos terminen en `*_test.r2`
- Comprobar que estén en los directorios configurados
- Usar `-debug` para ver qué archivos se están descubriendo

### Timeouts
- Aumentar el timeout con `-timeout 60s`
- Verificar que los tests no tengan bucles infinitos
- Considerar usar `-parallel` para mejorar el rendimiento

### Problemas de Cobertura
- Verificar que los patrones de exclusión no sean demasiado amplios
- Usar `-coverage-exclude` para excluir archivos innecesarios
- Comprobar que los archivos fuente estén en el directorio correcto

## Contribuir

Para contribuir a R2Test:
1. Crear tests para nuevas funcionalidades
2. Mantener compatibilidad con versiones anteriores
3. Actualizar documentación
4. Seguir las convenciones de código de R2Lang

## Licencia

R2Test está licenciado bajo Apache License 2.0.