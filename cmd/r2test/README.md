# R2Test - Advanced Testing Framework for R2Lang

R2Test is a powerful testing framework designed specifically for R2Lang, providing unit testing, code coverage, mocking, fixtures, and advanced reporting capabilities.

## Installation

```bash
go build -o r2test cmd/r2test/main.go
```

## Basic Usage

### Run all tests
```bash
r2test
```

### Run tests in a specific directory
```bash
r2test ./tests
```

### Run tests with coverage
```bash
r2test -coverage -verbose ./tests
```

### Run tests with pattern filtering
```bash
r2test -grep "Calculator" ./tests
```

### Run tests in parallel
```bash
r2test -parallel -workers 4
```

## Command Line Options

### Basic Options
- `-help` - Show help information
- `-version` - Show version information
- `-config FILE` - Load configuration from JSON file
- `-verbose` - Enable verbose output
- `-quiet` - Quiet output (errors only)
- `-debug` - Enable debug output
- `-dry-run` - Show what would be executed without running

### Test Discovery
- `-dirs DIRS` - Comma-separated test directories (default: `./tests,./test`)
- `-patterns PATTERNS` - Test file patterns (default: `*_test.r2,*Test.r2,test_*.r2`)
- `-ignore PATTERNS` - Patterns to ignore (default: `node_modules,vendor`)

### Test Execution
- `-timeout DURATION` - Test timeout (e.g., `30s`, `5m`) (default: `30s`)
- `-parallel` - Run tests in parallel
- `-workers N` - Maximum parallel workers (default: `4`)
- `-bail` - Stop on first test failure
- `-retries N` - Number of retries for failed tests

### Test Filtering
- `-grep PATTERN` - Run only tests matching pattern
- `-tags TAGS` - Run only tests with specified tags
- `-skip PATTERN` - Skip tests matching pattern
- `-only PATTERN` - Run only tests matching pattern (exclusive)

### Coverage Options
- `-coverage` - Enable coverage collection
- `-coverage-threshold N` - Coverage threshold percentage (default: `80`)
- `-coverage-output DIR` - Coverage output directory (default: `./coverage`)
- `-coverage-formats LIST` - Coverage formats: `html,json` (default: `html`)
- `-coverage-exclude LIST` - Coverage exclude patterns

### Reporting Options
- `-reporters LIST` - Reporters: `console,json,junit` (default: `console`)
- `-output DIR` - Output directory for reports (default: `./test-results`)

### Advanced Options
- `-watch` - Watch mode - rerun tests on file changes
- `-fixtures DIR` - Fixture directory (default: `./fixtures`)
- `-cleanup-mocks` - Cleanup mocks after tests (default: `true`)
- `-isolation` - Run tests in isolation contexts

## Configuration File

You can use a JSON file to configure R2Test:

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

## Test Syntax

R2Test uses `describe()` and `it()` syntax to structure tests:

```r2
describe("Math Operations", func() {
    it("should add two numbers correctly", func() {
        let result = 2 + 3;
        assert.equals(result, 5);
    });
    
    it("should subtract numbers correctly", func() {
        let result = 10 - 4;
        assert.equals(result, 6);
    });
});
```

### Available Assertions

#### Basic Comparisons
- `assert.equals(actual, expected)` - Equality comparison
- `assert.notEquals(actual, expected)` - Inequality comparison
- `assert.true(value)` - Verify value is truthy
- `assert.false(value)` - Verify value is falsy

#### Numeric Comparisons
- `assert.greater(a, b)` - a > b
- `assert.greaterOrEqual(a, b)` - a >= b
- `assert.less(a, b)` - a < b
- `assert.lessOrEqual(a, b)` - a <= b

#### Content Verification
- `assert.contains(string, substring)` - Verify string contains substring
- `assert.notContains(string, substring)` - Verify string doesn't contain substring
- `assert.hasLength(collection, length)` - Verify collection length
- `assert.empty(collection)` - Verify collection is empty
- `assert.notEmpty(collection)` - Verify collection is not empty

#### Null Checks
- `assert.nil(value)` - Verify value is null
- `assert.notNil(value)` - Verify value is not null

#### Error Handling
- `assert.panics(func)` - Verify function throws an error
- `assert.notPanics(func)` - Verify function doesn't throw an error

## Fixtures

R2Test provides fixture support for test data:

```r2
describe("Tests with Fixtures", func() {
    it("should load fixture data", func() {
        let data = fixture.load("users.json");
        assert.notEmpty(data);
        assert.hasLength(data, 3);
    });
});
```

## Mocking

R2Test includes mocking capabilities to isolate code under test:

```r2
describe("Tests with Mocks", func() {
    it("should mock a function", func() {
        let mockFunc = mock.create();
        mock.returns(mockFunc, "mocked result");
        
        let result = mockFunc();
        assert.equals(result, "mocked result");
    });
});
```

## Code Coverage

To generate coverage reports:

```bash
r2test -coverage -coverage-formats html,json -coverage-output ./coverage
```

This will generate:
- `./coverage/html/index.html` - Interactive HTML report
- `./coverage/coverage.json` - Coverage data in JSON

## Reports

R2Test can generate various types of reports:

### JSON Report
```bash
r2test -reporters json -output ./reports
```

### JUnit XML Report
```bash
r2test -reporters junit -output ./reports
```

### Multiple Reports
```bash
r2test -reporters console,json,junit -output ./reports
```

## Usage Examples

### Local Development
```bash
# Run tests with verbose output
r2test -verbose

# Run only specific tests
r2test -grep "Calculator"

# Run tests with coverage
r2test -coverage -coverage-threshold 90
```

### Continuous Integration
```bash
# Run tests with JUnit report for CI
r2test -reporters junit -output ./test-results -bail

# Run tests with coverage for CI
r2test -coverage -coverage-threshold 80 -reporters json,junit -output ./reports
```

### Development with Watch Mode
```bash
# Re-run tests when files change
r2test -watch -verbose
```

## Test File Structure

Test files should:
- End with `*_test.r2`, `*Test.r2`, or `test_*.r2`
- Be located in configured directories (default: `./tests`, `./test`)
- Use `describe()` and `it()` functions to structure tests

```
project/
├── tests/
│   ├── calculator_test.r2
│   ├── utils_test.r2
│   └── fixtures/
│       └── test_data.json
├── src/
│   ├── calculator.r2
│   └── utils.r2
└── r2test.config.json
```

## Editor Integration

R2Test is compatible with editors that support:
- JUnit XML report formats
- Structured JSON reports
- Standard console output

## Troubleshooting

### Tests Not Found
- Verify files end with `*_test.r2`
- Check that they're in configured directories
- Use `-debug` to see what files are being discovered

### Timeouts
- Increase timeout with `-timeout 60s`
- Verify tests don't have infinite loops
- Consider using `-parallel` for better performance

### Coverage Issues
- Verify exclude patterns aren't too broad
- Use `-coverage-exclude` to exclude unnecessary files
- Check that source files are in the correct directory

## Contributing

To contribute to R2Test:
1. Create tests for new features
2. Maintain backward compatibility
3. Update documentation
4. Follow R2Lang code conventions

## License

R2Test is licensed under Apache License 2.0.