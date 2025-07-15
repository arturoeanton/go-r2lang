# Propuesta: Soporte de Testing Unitario Avanzado en R2Lang

**Versión:** 1.0  
**Fecha:** 2025-07-15  
**Estado:** Propuesta  

## Resumen Ejecutivo

Esta propuesta presenta un sistema integral de testing unitario para R2Lang que extiende las capacidades BDD existentes con funcionalidades modernas de testing, incluyendo mocking, fixtures, cobertura de código y integración continua.

## Problema Actual

R2Lang actualmente tiene soporte básico para testing BDD con `TestCase`, `Given`, `When`, `Then`, pero carece de:

- **Framework de testing unitario estructurado**
- **Mocking y stubbing**
- **Fixtures y setup/teardown**
- **Cobertura de código**
- **Assertions avanzadas**
- **Parallel testing**
- **Test discovery automático**

### Ejemplo de Test Actual
```r2
TestCase "Usuario puede registrarse" {
    Given userService = UserService();
    When result = userService.register("juan", "pass123");
    Then result.success == true;
}
```

## Solución Propuesta

### 1. Framework de Testing Unificado

#### 1.1 Estructura de Test Modules
```r2
// tests/user_service_test.r2
import "assert" as assert;
import "mock" as mock;
import "../src/user_service.r2" as UserService;

describe("UserService", func() {
    let userService;
    
    beforeEach(func() {
        userService = UserService.new();
    });
    
    afterEach(func() {
        userService.cleanup();
    });
    
    it("debería registrar un usuario válido", func() {
        // Arrange
        let userData = {
            name: "Juan",
            email: "juan@test.com",
            password: "secure123"
        };
        
        // Act
        let result = userService.register(userData);
        
        // Assert
        assert.isTrue(result.success);
        assert.equals(result.user.name, "Juan");
        assert.isNotNull(result.user.id);
    });
    
    it("debería fallar con email inválido", func() {
        // Arrange
        let userData = {
            name: "Juan",
            email: "invalid-email",
            password: "secure123"
        };
        
        // Act & Assert
        assert.throws(func() {
            userService.register(userData);
        }, "InvalidEmailError");
    });
});
```

#### 1.2 Testing API Completa
```r2
// Built-in testing functions
describe(description, testSuite);
context(description, testSuite);  // Alias for describe
it(description, testFunction);
test(description, testFunction);  // Alias for it

// Lifecycle hooks
beforeAll(setupFunction);
afterAll(teardownFunction);
beforeEach(setupFunction);
afterEach(teardownFunction);

// Skipping and focusing
describe.skip(description, testSuite);
it.skip(description, testFunction);
describe.only(description, testSuite);
it.only(description, testFunction);

// Async testing
it("should handle async operations", async func() {
    let result = await userService.fetchUserAsync(123);
    assert.equals(result.id, 123);
});
```

### 2. Sistema de Assertions Avanzado

#### 2.1 Assertions Básicas
```r2
// pkg/r2libs/r2assert.go
func RegisterAssert(env *r2core.Environment) {
    env.Set("assert", map[string]interface{}{
        "equals":         assert_equals,
        "notEquals":      assert_notEquals,
        "isTrue":         assert_isTrue,
        "isFalse":        assert_isFalse,
        "isNull":         assert_isNull,
        "isNotNull":      assert_isNotNull,
        "isUndefined":    assert_isUndefined,
        "isEmpty":        assert_isEmpty,
        "isNotEmpty":     assert_isNotEmpty,
        "contains":       assert_contains,
        "notContains":    assert_notContains,
        "startsWith":     assert_startsWith,
        "endsWith":       assert_endsWith,
        "matches":        assert_matches,
        "throws":         assert_throws,
        "notThrows":      assert_notThrows,
        "closeTo":        assert_closeTo,
        "between":        assert_between,
        "hasProperty":    assert_hasProperty,
        "hasLength":      assert_hasLength,
        "instanceOf":     assert_instanceOf,
        "deepEquals":     assert_deepEquals,
    })
}
```

#### 2.2 Assertions con Mensajes Descriptivos
```r2
assert.equals(actual, expected, "Los números deberían ser iguales");
assert.isTrue(user.isActive(), `Usuario ${user.name} debería estar activo`);

// Assertions con contexto automático
assert.that(result)
    .isNotNull("El resultado no debería ser null")
    .hasProperty("success", "El resultado debería tener propiedad success")
    .property("success").isTrue("La operación debería ser exitosa");
```

### 3. Sistema de Mocking y Stubbing

#### 3.1 Mock Objects
```r2
// tests/payment_test.r2
import "mock" as mock;
import "../src/payment_service.r2" as PaymentService;

describe("PaymentService", func() {
    it("debería procesar pago con gateway mock", func() {
        // Arrange
        let mockGateway = mock.create("PaymentGateway");
        mock.when(mockGateway.charge)
            .calledWith(100.0, "USD")
            .returns({success: true, transactionId: "tx123"});
        
        let paymentService = PaymentService.new(mockGateway);
        
        // Act
        let result = paymentService.processPayment(100.0, "USD");
        
        // Assert
        assert.isTrue(result.success);
        assert.equals(result.transactionId, "tx123");
        mock.verify(mockGateway.charge).wasCalledOnce();
        mock.verify(mockGateway.charge).wasCalledWith(100.0, "USD");
    });
});
```

#### 3.2 Spies y Stubs
```r2
describe("Analytics Service", func() {
    it("debería enviar métricas", func() {
        // Spy on existing function
        let spy = mock.spy(analytics, "track");
        
        // Act
        userService.login("user123");
        
        // Assert
        assert.that(spy)
            .wasCalledOnce()
            .wasCalledWith("user_login", {userId: "user123"});
    });
    
    it("debería manejar falla de red", func() {
        // Stub network call to fail
        mock.stub(httpClient, "post").throws("NetworkError");
        
        // Act & Assert
        assert.throws(func() {
            userService.syncData();
        }, "NetworkError");
    });
});
```

### 4. Fixtures y Data Management

#### 4.1 Test Fixtures
```r2
// tests/fixtures/users.r2
export let validUser = {
    name: "Juan Pérez",
    email: "juan@example.com",
    age: 30,
    active: true
};

export let invalidUsers = [
    {name: "", email: "invalid", age: -5},
    {name: null, email: "test@example.com", age: 150},
    {name: "Valid", email: "", age: 25}
];

export let createUserWithDefaults = func(overrides = {}) {
    return Object.assign({}, validUser, overrides);
};
```

#### 4.2 Database Fixtures
```r2
// tests/fixtures/database.r2
import "database" as db;

export let setupTestDb = func() {
    db.execute("CREATE TABLE users (id INT, name VARCHAR(100), email VARCHAR(100))");
    db.execute("INSERT INTO users VALUES (1, 'Test User', 'test@example.com')");
};

export let cleanupTestDb = func() {
    db.execute("DROP TABLE IF EXISTS users");
};

export let withTestDb = func(testFunction) {
    return func() {
        setupTestDb();
        try {
            testFunction();
        } finally {
            cleanupTestDb();
        }
    };
};
```

### 5. Cobertura de Código

#### 5.1 Instrumentación de Código
```go
// pkg/r2core/coverage.go
type CoverageCollector struct {
    Files     map[string]*FileCoverage
    Enabled   bool
    OutputDir string
}

type FileCoverage struct {
    Filename      string
    TotalLines    int
    CoveredLines  map[int]bool
    FunctionCalls map[string]int
    Branches      map[string]bool
}

func (cc *CoverageCollector) RecordLineExecution(file string, line int) {
    if !cc.Enabled {
        return
    }
    
    if cc.Files[file] == nil {
        cc.Files[file] = &FileCoverage{
            Filename:     file,
            CoveredLines: make(map[int]bool),
        }
    }
    
    cc.Files[file].CoveredLines[line] = true
}

func (cc *CoverageCollector) GenerateReport() *CoverageReport {
    // Generate HTML/JSON coverage reports
}
```

#### 5.2 Integración con Runtime
```r2
// Ejecutar tests con cobertura
r2 test --coverage ./tests/
r2 test --coverage --output coverage.html ./tests/
r2 test --coverage --threshold 80 ./tests/

// En código R2Lang
coverage.start();
runAllTests();
let report = coverage.stop();

print(`Cobertura total: ${report.totalPercentage}%`);
print(`Líneas cubiertas: ${report.coveredLines}/${report.totalLines}`);
```

### 6. Test Runner y Discovery

#### 6.1 Test Discovery Automático
```bash
# Descubrimiento automático de tests
r2 test                          # Ejecuta todos los tests en ./tests/
r2 test ./specific/test/file.r2  # Ejecuta test específico
r2 test --pattern "*_test.r2"    # Patrón personalizado
r2 test --grep "UserService"     # Filtra por descripción
```

#### 6.2 Test Runner Configuración
```r2
// r2test.config.r2
export let config = {
    testDir: "./tests",
    pattern: ["*_test.r2", "*Test.r2"],
    ignore: ["node_modules", "vendor"],
    timeout: 5000,      // 5 segundos por test
    parallel: true,     // Ejecutar tests en paralelo
    maxWorkers: 4,      // Número de workers paralelos
    bail: false,        // Parar en primer fallo
    verbose: true,      // Output detallado
    coverage: {
        enabled: false,
        threshold: 80,
        output: "coverage.html"
    },
    reporters: ["console", "junit", "json"]
};
```

### 7. Parallel Testing

#### 7.1 Test Workers
```go
// pkg/r2core/test_runner.go
type TestRunner struct {
    MaxWorkers int
    Timeout    time.Duration
    Parallel   bool
}

func (tr *TestRunner) RunTestsParallel(testFiles []string) *TestResults {
    if !tr.Parallel {
        return tr.RunTestsSerial(testFiles)
    }
    
    workers := make(chan string, tr.MaxWorkers)
    results := make(chan *TestResult, len(testFiles))
    
    // Start workers
    for i := 0; i < tr.MaxWorkers; i++ {
        go tr.testWorker(workers, results)
    }
    
    // Send work
    for _, file := range testFiles {
        workers <- file
    }
    close(workers)
    
    // Collect results
    var allResults []*TestResult
    for i := 0; i < len(testFiles); i++ {
        allResults = append(allResults, <-results)
    }
    
    return tr.AggregateResults(allResults)
}
```

#### 7.2 Test Isolation
```r2
describe("Parallel Safe Tests", func() {
    // Tests que pueden ejecutarse en paralelo
    it.parallel("test independiente 1", func() {
        // Test que no comparte estado
    });
    
    it.parallel("test independiente 2", func() {
        // Test que no comparte estado
    });
    
    // Tests que requieren ejecución serial
    it.serial("test con estado compartido", func() {
        // Test que modifica estado global
    });
});
```

### 8. Integración con CLI

#### 8.1 Comandos de Testing
```bash
# Comandos básicos
r2 test                              # Ejecutar todos los tests
r2 test --watch                      # Modo watch (re-ejecutar en cambios)
r2 test --parallel                   # Ejecutar en paralelo
r2 test --serial                     # Forzar ejecución serial
r2 test --timeout 10s                # Timeout personalizado

# Filtering
r2 test --grep "UserService"         # Filtrar por descripción
r2 test --pattern "*integration*"    # Patrón de archivos
r2 test --tag slow                   # Ejecutar tests marcados

# Coverage
r2 test --coverage                   # Generar reporte de cobertura
r2 test --coverage --threshold 85    # Fallar si cobertura < 85%

# Output
r2 test --reporter json              # Output JSON
r2 test --reporter junit             # Output JUnit XML
r2 test --verbose                    # Output detallado
r2 test --quiet                      # Output mínimo

# CI/CD Integration
r2 test --ci                         # Modo CI (no colores, output optimizado)
r2 test --fail-fast                  # Parar en primer fallo
r2 test --retry 3                    # Reintentar tests fallidos
```

### 9. Implementación Técnica

#### 9.1 Test Framework Core
```go
// pkg/r2libs/r2test.go
type TestSuite struct {
    Description string
    Tests       []*Test
    Hooks       *TestHooks
    Parent      *TestSuite
    Children    []*TestSuite
}

type Test struct {
    Description string
    Function    r2core.Function
    Skip        bool
    Only        bool
    Parallel    bool
    Timeout     time.Duration
    Tags        []string
}

type TestHooks struct {
    BeforeAll  []r2core.Function
    AfterAll   []r2core.Function
    BeforeEach []r2core.Function
    AfterEach  []r2core.Function
}

func RegisterTest(env *r2core.Environment) {
    env.Set("describe", r2core.BuiltinFunction(describe))
    env.Set("it", r2core.BuiltinFunction(it))
    env.Set("beforeAll", r2core.BuiltinFunction(beforeAll))
    env.Set("afterAll", r2core.BuiltinFunction(afterAll))
    env.Set("beforeEach", r2core.BuiltinFunction(beforeEach))
    env.Set("afterEach", r2core.BuiltinFunction(afterEach))
}
```

#### 9.2 Test Execution Engine
```go
type TestExecutor struct {
    Suite    *TestSuite
    Reporter TestReporter
    Coverage *CoverageCollector
}

func (te *TestExecutor) ExecuteTest(test *Test) *TestResult {
    start := time.Now()
    
    // Setup context
    ctx := context.WithTimeout(context.Background(), test.Timeout)
    env := r2core.NewEnvironment()
    
    // Execute hooks
    te.executeBeforeHooks(env)
    
    // Execute test
    result := &TestResult{
        Test:      test,
        StartTime: start,
    }
    
    defer func() {
        if r := recover(); r != nil {
            result.Error = fmt.Errorf("Test panic: %v", r)
            result.Status = TestStatusFailed
        }
        result.Duration = time.Since(start)
        te.executeAfterHooks(env)
    }()
    
    // Run test function
    test.Function.Call(env, []interface{}{})
    result.Status = TestStatusPassed
    
    return result
}
```

### 10. Integración con BDD Existente

#### 10.1 Migración Gradual
```r2
// Mantener compatibilidad con sintaxis BDD existente
TestCase "Usuario puede registrarse" {
    Given userService = UserService();
    When result = userService.register("juan", "pass123");
    Then result.success == true;
}

// Equivalente en nuevo framework
describe("Usuario", func() {
    it("puede registrarse", func() {
        // Given
        let userService = UserService();
        
        // When
        let result = userService.register("juan", "pass123");
        
        // Then
        assert.isTrue(result.success);
    });
});
```

#### 10.2 Híbrido BDD + Unit Testing
```r2
// Combinar ambos enfoques según necesidad
describe("UserService Integration", func() {
    // Unit tests tradicionales
    it("debería validar email", func() {
        assert.isFalse(UserService.isValidEmail("invalid"));
    });
    
    // BDD scenarios para flujos complejos
    scenario("User registration flow", func() {
        Given("a new user wants to register", func() {
            userData = {name: "Juan", email: "juan@test.com"};
        });
        
        When("they submit registration form", func() {
            result = userService.register(userData);
        });
        
        Then("they should be registered successfully", func() {
            assert.isTrue(result.success);
            assert.isNotNull(result.user.id);
        });
    });
});
```

## Beneficios

1. **Testing Moderno:** Framework completo con mejores prácticas
2. **Flexibilidad:** Múltiples estilos de testing (unit, BDD, integration)
3. **Productividad:** Mocking, fixtures y assertions avanzadas
4. **Calidad:** Cobertura de código y métricas de testing
5. **CI/CD Ready:** Integración perfecta con pipelines
6. **Parallelización:** Tests más rápidos en proyectos grandes

## Plan de Implementación

### Fase 1: Core Framework
- [ ] Test runner básico con discovery
- [ ] Sistema de assertions
- [ ] Lifecycle hooks (beforeEach, afterEach)
- [ ] CLI básico (`r2 test`)

### Fase 2: Mocking y Fixtures
- [ ] Sistema de mocking y spies
- [ ] Fixture management
- [ ] Async testing support
- [ ] Test isolation

### Fase 3: Cobertura y Reporting
- [ ] Coverage collector
- [ ] Multiple reporters (console, JSON, JUnit)
- [ ] HTML coverage reports
- [ ] CI/CD integration

### Fase 4: Optimización y Paralelización
- [ ] Parallel test execution
- [ ] Watch mode
- [ ] Performance optimizations
- [ ] Advanced filtering

### Fase 5: Ecosistema
- [ ] IDE integrations
- [ ] Plugin system
- [ ] Test generators
- [ ] Documentation completa

## Conclusión

Este sistema de testing unitario proporcionará a R2Lang capacidades de testing de nivel industrial, manteniendo la simplicidad y expresividad que caracterizan al lenguaje, mientras permite escalabilidad para proyectos grandes y equipos de desarrollo.