# R2Lang Unit Testing Examples

This directory contains comprehensive examples demonstrating the advanced unit testing capabilities of R2Lang, implementing phases 2 and 3 of the unit testing proposal.

## Features Demonstrated

### Phase 1: Core Framework ✅
- Basic test structure with `describe()` and `it()`
- Test lifecycle hooks (`beforeAll`, `afterAll`, `beforeEach`, `afterEach`)
- Comprehensive assertion library
- Test discovery and execution

### Phase 2: Mocking and Fixtures ✅
- Mock object creation and verification
- Spy functionality with call-through capability
- Stub implementation for method replacement
- Test isolation contexts
- Fixture management (JSON, CSV, text files)
- Automatic cleanup and restoration

### Phase 3: Coverage and Reporting ✅
- Line coverage collection
- Statement and branch coverage tracking
- Function coverage monitoring
- Multiple report formats (HTML, JSON, JUnit XML)
- Coverage thresholds and validation
- Uncovered line identification

## Example Files

### 1. `basic_test_example.r2`
Demonstrates the core testing features:
- Simple test suites with nested describes
- Various assertion types
- Lifecycle hooks
- Array, object, and string testing
- Type validation

**Run with:**
```bash
go run main.go examples/unit_testing/basic_test_example.r2
```

### 2. `mocking_example.r2`
Shows advanced mocking and spying capabilities:
- HTTP request mocking
- Database operation stubbing
- Method spy with call-through
- Mock verification and call counting
- Test isolation contexts
- Fixture loading (JSON, CSV, text)

**Key Features:**
- `mock.createMock()` - Create mock objects
- `mock.spyOn()` - Spy on existing methods
- `mock.runInIsolation()` - Run tests in isolated contexts
- `fixtures.createTemporary()` - Create test fixtures
- Automatic mock cleanup and restoration

### 3. `coverage_example.r2`
Demonstrates coverage collection and reporting:
- Line-by-line coverage tracking
- Branch coverage monitoring
- Function coverage validation
- Multiple report generation
- Coverage threshold checking
- Uncovered code identification

**Generated Reports:**
- HTML report with interactive source view
- JSON report for CI/CD integration
- JUnit XML for test result aggregation

## Running the Examples

### Prerequisites
1. Ensure R2Lang is compiled: `go build`
2. Create test directories:
   ```bash
   mkdir -p coverage/html
   mkdir -p test_fixtures
   ```

### Execute Examples
```bash
# Basic testing features
go run main.go examples/unit_testing/basic_test_example.r2

# Mocking and fixtures
go run main.go examples/unit_testing/mocking_example.r2

# Coverage and reporting
go run main.go examples/unit_testing/coverage_example.r2
```

## Test Structure

### Basic Test Syntax
```r2
import "r2test" as test;

test.describe("Component Name", func() {
    test.it("should do something", func() {
        let result = someOperation();
        test.assert("description").that(result).equals(expected);
    });
});
```

### Mocking Example
```r2
import "r2test/mocking" as mock;

// Create mock
let httpMock = mock.createMock("httpService");
httpMock.when("get", "/api/users").returns({users: []});

// Use mock
let result = httpMock.call("get", "/api/users");

// Verify
test.assert("mock called").that(httpMock.wasCalled("get")).isTrue();
```

### Coverage Tracking
```r2
import "r2test/coverage" as coverage;

// Enable coverage
coverage.enable();
coverage.start();

// Record hits during test execution
coverage.recordHit("myfile.r2", 15);

// Generate reports
let stats = coverage.getStats();
```

## Assertion Methods

### Basic Assertions
- `.equals(value)` - Exact equality
- `.isTrue()` / `.isFalse()` - Boolean checks
- `.isNull()` / `.isNotNull()` - Null checks
- `.isUndefined()` / `.isNotUndefined()` - Undefined checks

### Type Assertions
- `.isNumber()` / `.isString()` / `.isBoolean()`
- `.isArray()` / `.isObject()`
- `.isFunction()`

### Comparison Assertions
- `.isGreaterThan(value)` / `.isLessThan(value)`
- `.isGreaterThanOrEqual(value)` / `.isLessThanOrEqual(value)`

### String Assertions
- `.contains(substring)` - String contains check
- `.startsWith(prefix)` / `.endsWith(suffix)`
- `.matches(pattern)` - Regex matching

### Collection Assertions
- `.hasProperty(key)` - Object property check
- `.hasLength(length)` - Array/string length
- `.isEmpty()` / `.isNotEmpty()`

### Exception Assertions
- `.throws()` - Expects function to throw
- `.doesNotThrow()` - Expects function not to throw
- `.withMessage(message)` - Specific error message

## Configuration

### Coverage Configuration
```r2
// Set coverage options
coverage.setBasePath("./src");
coverage.addExcludeGlob("*_test.r2");
coverage.addExcludeGlob("node_modules/*");

// Set thresholds
let threshold = 80.0; // 80% minimum coverage
if (!coverage.meetsThreshold(threshold)) {
    throw "Coverage below threshold";
}
```

### Test Configuration
```r2
// Set fixture base path
fixtures.setBasePath("./test_data");

// Configure test timeouts
test.setDefaultTimeout(5000); // 5 seconds

// Set parallel execution
test.setParallel(true);
test.setMaxWorkers(4);
```

## Best Practices

### 1. Test Organization
- Use descriptive test names
- Group related tests with nested `describe()` blocks
- Use lifecycle hooks for setup and cleanup
- Keep tests independent and isolated

### 2. Mocking Strategy
- Mock external dependencies (HTTP, database, file system)
- Use spies for partial mocking when you need original behavior
- Verify mock interactions to ensure correct usage
- Clean up mocks after each test

### 3. Coverage Goals
- Aim for 80%+ line coverage
- Focus on critical business logic
- Test both success and error paths
- Use coverage reports to identify untested code

### 4. Assertions
- Use descriptive assertion messages
- Test one concept per test case
- Use specific assertions (e.g., `.isTrue()` vs `.equals(true)`)
- Test edge cases and boundary conditions

## Integration with CI/CD

### Coverage Reports
The testing framework generates multiple report formats suitable for CI/CD integration:

- **HTML Reports**: Human-readable coverage reports with source highlighting
- **JSON Reports**: Machine-readable data for custom processing
- **JUnit XML**: Compatible with most CI systems (Jenkins, Azure DevOps, etc.)

### Example CI Configuration (GitHub Actions)
```yaml
- name: Run Tests with Coverage
  run: |
    go run main.go examples/unit_testing/coverage_example.r2
    
- name: Upload Coverage Reports
  uses: actions/upload-artifact@v2
  with:
    name: coverage-reports
    path: coverage/
```

## Troubleshooting

### Common Issues

1. **Mock not working**: Ensure mocks are created before the code under test runs
2. **Coverage not recording**: Check that coverage is enabled and files aren't excluded
3. **Tests failing in isolation**: Verify proper cleanup in `afterEach` hooks
4. **Fixtures not loading**: Check file paths and permissions

### Debug Mode
Enable verbose output for debugging:
```r2
test.setVerbose(true);
coverage.setVerbose(true);
```

## Future Enhancements

The testing framework is designed to be extensible. Future phases may include:

- **Phase 4**: Parallel test execution and performance optimizations
- **Phase 5**: Snapshot testing, retry mechanisms, and performance testing
- **Phase 6**: IDE integration and plugin system

## Contributing

To extend the testing framework:

1. Add new assertion methods in `pkg/r2test/assertions/`
2. Create custom reporters in `pkg/r2test/reporters/`
3. Add fixture loaders in `pkg/r2test/fixtures/`
4. Implement new coverage collectors in `pkg/r2test/coverage/`

See the main project documentation for contribution guidelines.