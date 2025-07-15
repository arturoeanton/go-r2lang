<div align="center">
  <br />
  <h1>R2Lang</h1>
  <p>
    <b>Escribe pruebas elegantes, scripts y aplicaciones con un lenguaje que combina simplicidad y poder.</b>
  </p>
  <br />
</div>

<div align="center">

[![Go Report Card](https://goreportcard.com/badge/github.com/arturoeanton/go-r2lang)](https://goreportcard.com/report/github.com/arturoeanton/go-r2lang)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GitHub stars](https://img.shields.io/github/stars/arturoeanton/go-r2lang.svg?style=social&label=Star)](https://github.com/arturoeanton/go-r2lang)
[![GitHub forks](https://img.shields.io/github/forks/arturoeanton/go-r2lang.svg?style=social&label=Fork)](https://github.com/arturoeanton/go-r2lang)
[![GitHub issues](https://img.shields.io/github/issues/arturoeanton/go-r2lang.svg)](https://github.com/arturoeanton/go-r2lang/issues)
[![Contributors](https://img.shields.io/github/contributors/arturoeanton/go-r2lang.svg)](https://github.com/arturoeanton/go-r2lang/graphs/contributors)

</div>

---

**R2Lang** es un lenguaje de programación dinámico e interpretado escrito en Go. Está diseñado para ser simple, intuitivo y poderoso, con una sintaxis fuertemente inspirada en JavaScript y soporte de primera clase para **pruebas unitarias integrales**.

Ya sea que estés escribiendo scripts de automatización, construyendo una API web, o creando una suite de pruebas robusta, R2Lang proporciona las herramientas que necesitas en un paquete limpio y legible.

## ✨ Características Principales

| Característica                 | Descripción                                                                                                 | Ejemplo                                                              |
| ----------------------- | ----------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------- |
| **🧪 Framework de Pruebas Integrado** | Framework completo de testing con sintaxis `describe()` e `it()`. No se necesitan frameworks externos.        | `describe("Login de Usuario", func() { it("debería autenticar", func() { ... }) })`      |
| **🚀 Simple y Familiar**    | Si conoces JavaScript, te sentirás como en casa. Esto hace que sea increíblemente fácil de aprender y comenzar a programar. | `let mensaje = "¡Hola, Mundo!"; print(mensaje);`                     |
| **⚡ Concurrente**          | Aprovecha el poder de las goroutinas de Go con una simple función `r2()` para ejecutar código en paralelo.                | `r2(miFuncion, "arg1");`                                            |
| **🧱 Orientado a Objetos**     | Usa clases, herencia (`extends`), y `this` para estructurar tu código de manera limpia y orientada a objetos.    | `class Usuario extends Persona { ... }`                                  |
| **🌐 Listo para Web**            | Crea servidores web y APIs REST con una librería `http` integrada que se siente como Express.js.                 | `http.get("/usuarios", func(req, res) { res.json(...) });`               |
| **🧩 Fácilmente Extensible**   | Escrito en Go, R2Lang puede ser fácilmente extendido con nuevas funciones y librerías nativas.                       | `env.Set("miFuncionNativa", r2lang.NewBuiltinFunction(...));`            |

---

## 🚀 Comenzando

### Prerequisitos

- **Go**: Versión 1.23 o superior.

### Instalación y "¡Hola, Mundo!"

1.  **Clona el repositorio:**
    ```bash
    git clone https://github.com/arturoeanton/go-r2lang.git
    cd go-r2lang
    ```

2.  **Construye el intérprete:**
    ```bash
    go build -o r2lang main.go
    ```

3.  **Crea tu primer archivo R2Lang (`hola.r2`):**
    ```r2
    func main() {
        print("¡Hola, R2Lang! 🚀");
    }
    ```

4.  **¡Ejecútalo!**
    ```bash
    ./r2lang hola.r2
    # Salida: ¡Hola, R2Lang! 🚀
    ```

---

## 🧪 Framework Avanzado de Pruebas Unitarias

R2Lang incluye un framework de testing de nivel profesional con capacidades de nivel empresarial incluyendo **mocking**, **fixtures**, **reportes de cobertura** y **aislamiento de pruebas**:

### Estructura Básica de Pruebas
```r2
import "r2test" as test;

test.describe("Calculadora", func() {
    test.it("debería sumar números correctamente", func() {
        let resultado = 2 + 3;
        test.assert("suma").that(resultado).equals(5);
    });
    
    test.it("debería manejar resta", func() {
        let resultado = 10 - 4;
        test.assert("resta").that(resultado).equals(6);
    });
});
```

### Hooks de Ciclo de Vida de Pruebas
```r2
test.describe("Pruebas de Base de Datos", func() {
    test.beforeAll(func() {
        // Configurar conexión a base de datos
    });
    
    test.beforeEach(func() {
        // Resetear datos de prueba antes de cada test
    });
    
    test.afterEach(func() {
        // Limpiar después de cada test
    });
    
    test.afterAll(func() {
        // Cerrar conexión a base de datos
    });
    
    test.it("debería guardar datos de usuario", func() {
        // Implementación de la prueba
    });
});
```

### Librería Integral de Aserciones
- **Básicas**: `.equals()`, `.isTrue()`, `.isFalse()`, `.isNull()`, `.isNotNull()`
- **Tipos**: `.isNumber()`, `.isString()`, `.isBoolean()`, `.isArray()`, `.isObject()`
- **Comparaciones**: `.isGreaterThan()`, `.isLessThan()`, `.isGreaterThanOrEqual()`
- **Strings**: `.contains()`, `.startsWith()`, `.endsWith()`, `.matches()`
- **Colecciones**: `.hasProperty()`, `.hasLength()`, `.isEmpty()`, `.isNotEmpty()`
- **Excepciones**: `.throws()`, `.doesNotThrow()`, `.withMessage()`

### Mocking y Espionaje
```r2
import "r2test/mocking" as mock;

test.describe("Pruebas de Servicio HTTP", func() {
    test.it("debería mockear llamadas API", func() {
        let httpMock = mock.createMock("servicioHttp");
        httpMock.when("get", "/api/usuarios").returns({usuarios: []});
        
        let resultado = httpMock.call("get", "/api/usuarios");
        test.assert("Resultado API").that(resultado.usuarios).isArray();
        test.assert("Mock llamado").that(httpMock.wasCalled("get")).isTrue();
    });
    
    test.it("debería espiar métodos", func() {
        let spy = mock.spyOn("servicioUsuario.validar", servicioUsuario.validar);
        spy.callThrough();
        
        servicioUsuario.validar("test@email.com");
        
        test.assert("Spy llamado").that(spy.wasCalledWith("test@email.com")).isTrue();
    });
});
```

### Gestión de Fixtures
```r2
import "r2test/fixtures" as fixtures;

test.describe("Pruebas de Datos", func() {
    test.it("debería cargar fixtures JSON", func() {
        let datosUsuario = fixtures.load("datos_usuario.json");
        test.assert("Nombre de usuario").that(datosUsuario.nombre).equals("Juan Pérez");
    });
    
    test.it("debería crear fixtures temporales", func() {
        fixtures.createTemporary("datos_test", {id: 1, nombre: "Test"});
        let datos = fixtures.getData("datos_test");
        test.assert("Datos temporales").that(datos.id).equals(1);
    });
});
```

### Reportes de Cobertura
```r2
import "r2test/coverage" as coverage;
import "r2test/reporters" as reporters;

// Habilitar recolección de cobertura
coverage.enable();
coverage.start();

// Ejecutar pruebas y generar reportes
test.runTests();

// Generar múltiples formatos de reporte
reporters.generateHTMLReport("./cobertura/html");    // Reporte HTML interactivo
reporters.generateJSONReport("./cobertura/datos.json"); // JSON legible por máquina
reporters.generateJUnitReport("./cobertura/junit.xml"); // Integración CI/CD

// Verificar umbrales de cobertura
let stats = coverage.getStats();
if (stats.linePercentage < 80) {
    throw "Cobertura por debajo del umbral del 80%";
}
```

### Aislamiento de Pruebas
```r2
import "r2test/mocking" as mock;

test.describe("Pruebas Aisladas", func() {
    test.it("debería ejecutarse en aislamiento", func() {
        mock.runInIsolation("prueba aislada", func(context) {
            let mockDb = context.createMock("baseDatos");
            mockDb.when("guardar").returns(123);
            
            // La prueba ejecuta en completo aislamiento
            let resultado = mockDb.call("guardar", {datos: "test"});
            test.assert("Resultado aislado").that(resultado).equals(123);
        });
    });
});
```

### Ejecutando Pruebas
```bash
# Ejecutar archivos de prueba específicos
go run main.go examples/unit_testing/basic_test_example.r2
go run main.go examples/unit_testing/mocking_example.r2
go run main.go examples/unit_testing/coverage_example.r2

# Ejecutar todas las pruebas en directorio
go run main.go -test ./tests/

# Generar reportes de cobertura
go run main.go -test -coverage ./tests/
```

### Características Avanzadas
- **Descubrimiento de Pruebas**: Detección automática de archivos de prueba con patrones
- **Ejecución Paralela**: Ejecutar pruebas en paralelo para mayor velocidad
- **Modo Watch**: Re-ejecutar pruebas cuando los archivos cambien
- **Reportes Personalizados**: Sistema de reportes conectables
- **Integración CI/CD**: Salida JUnit XML y JSON para sistemas de build

Para ejemplos completos y características avanzadas, ver [examples/unit_testing/](./examples/unit_testing/).

---

## 📚 Documentación y Curso Completo

¿Listo para profundizar? Tenemos un curso completo, módulo por módulo para llevarte de principiante a experto.

-   [**Leer el Curso Completo (Español)**](./README.md)
-   [**Read the Full Course (English)**](../en/README.md)

La documentación cubre todo desde sintaxis básica hasta temas avanzados como concurrencia, manejo de errores y desarrollo web.

---

## 💖 Contribuyendo

**¡Estamos buscando activamente contribuidores!** Ya seas un desarrollador experimentado, un escritor de documentación, o simplemente alguien entusiasmado por nuevos lenguajes de programación, nos encantaría tu ayuda.

Así es como puedes contribuir:

1.  **Encuentra un issue:** Revisa nuestros [**Issues**](https://github.com/arturoeanton/go-r2lang/issues) y busca las etiquetas `good first issue` o `help wanted`.
2.  **Explora el Roadmap:** Ve nuestro [**Roadmap Técnico**](./roadmap-actualizado.md) para objetivos a largo plazo y características grandes en las que necesitamos ayuda.
3.  **Mejora la Documentación:** ¿Encontraste un error tipográfico o una sección que podría ser más clara? ¡Déjanos saber!
4.  **Envía un Pull Request:**
    -   Haz fork del repositorio.
    -   Crea una nueva rama (`git checkout -b feature/mi-caracteristica-increible`).
    -   Confirma tus cambios.
    -   ¡Abre un Pull Request!

Creemos en una comunidad acogedora y de apoyo. ¡Ninguna contribución es demasiado pequeña!

---

## 🗺️ Roadmap del Proyecto

¡Tenemos grandes planes para R2Lang! Nuestro objetivo es hacer que sea un lenguaje rápido, confiable y rico en características para una amplia gama de aplicaciones.

Las áreas clave de enfoque incluyen:

-   **🚀 Revolución de Rendimiento:** Implementar una VM de bytecode y eventualmente un compilador JIT para aumentos significativos de velocidad.
-   **🧠 Características Avanzadas:** Agregar coincidencia de patrones, un sistema de tipos más sofisticado y modelos avanzados de concurrencia.
-   **🛠️ Librería Estándar Más Rica:** Expandir las librerías integradas para bases de datos, sistemas de archivos y más.
-   **📦 Gestor de Paquetes:** Crear un gestor de paquetes dedicado para compartir y reutilizar código R2Lang.

Para una vista detallada de nuestros planes, revisa el [**Roadmap Técnico**](./roadmap-actualizado.md) y nuestra [**Lista TODO**](../../TODO.md).

---

## 🧪 Características del Framework de Testing

### Implementación Completada

#### Fase 1: Framework Base ✅
- ✅ Estructura básica de pruebas con `describe()` e `it()`
- ✅ Hooks de ciclo de vida (`beforeAll`, `afterAll`, `beforeEach`, `afterEach`)
- ✅ Librería integral de aserciones
- ✅ Descubrimiento y ejecución de pruebas

#### Fase 2: Mocking y Fixtures ✅
- ✅ Sistema de creación y verificación de mocks
- ✅ Funcionalidad de spy con capacidad de call-through
- ✅ Implementación de stubs para reemplazo de métodos
- ✅ Contextos de aislamiento de pruebas
- ✅ Gestión de fixtures (archivos JSON, CSV, texto)
- ✅ Limpieza automática y restauración

#### Fase 3: Cobertura y Reportes ✅
- ✅ Recolección de cobertura de líneas
- ✅ Seguimiento de cobertura de sentencias y ramas
- ✅ Monitoreo de cobertura de funciones
- ✅ Múltiples formatos de reporte (HTML, JSON, JUnit XML)
- ✅ Umbrales de cobertura y validación
- ✅ Identificación de líneas no cubiertas

### Próximas Fases

#### Fase 4: Paralelización y Optimización (Planificada)
- 🔄 Ejecución de pruebas en paralelo
- 🔄 Modo watch para re-ejecución automática
- 🔄 Optimizaciones de rendimiento
- 🔄 Filtrado avanzado por tags y patrones

#### Fase 5: Características Avanzadas (Planificada)
- 📋 Testing de snapshots
- 📋 Mecanismos de retry y timeout
- 📋 Sistema de etiquetado de pruebas
- 📋 Testing de rendimiento

---

## 🤝 Comunidad

-   **Reportar un Bug:** ¿Encontraste algo mal? Abre un [**Issue**](https://github.com/arturoeanton/go-r2lang/issues/new).
-   **Solicitar una Característica:** ¿Tienes una gran idea? Discutámosla en los [**Issues**](https://github.com/arturoeanton/go-r2lang/issues).
-   **Hacer una Pregunta:** No dudes en abrir un issue para preguntas y discusiones.

---

## 📜 Licencia

R2Lang está licenciado bajo la **Licencia Apache 2.0**. Ver el archivo [LICENSE](../../LICENSE) para detalles.