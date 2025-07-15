# R2Lang Unit Testing Framework Examples

This directory contains examples demonstrating the new R2Lang Unit Testing Framework that replaces the previous BDD system.

## Files

- `basic_test.r2` - Basic testing examples showing core assertions and test structure
- `advanced_test.r2` - Advanced examples with hooks, complex scenarios, and edge cases
- `r2test.config.json` - Example configuration file for the testing framework

## New Testing Syntax

The new testing framework uses a `describe()` and `it()` structure similar to modern testing frameworks:

```r2
describe("Feature Name", func() {
    it("should do something specific", func() {
        let result = someOperation();
        assert.equals(result, expectedValue);
    });
});
```

## Key Features

### 1. Test Organization
- `describe()` - Groups related tests together
- `it()` - Individual test cases
- Nested `describe()` blocks for hierarchical organization

### 2. Lifecycle Hooks
- `beforeEach()` - Runs before each test in the suite
- `afterEach()` - Runs after each test in the suite
- `beforeAll()` - Runs once before all tests in the suite
- `afterAll()` - Runs once after all tests in the suite

### 3. Comprehensive Assertions
- `assert.equals(actual, expected)` - Value equality
- `assert.notEquals(actual, expected)` - Value inequality
- `assert.true(value)` - Truthy assertion
- `assert.false(value)` - Falsy assertion
- `assert.nil(value)` - Null check
- `assert.notNil(value)` - Non-null check
- `assert.contains(haystack, needle)` - String contains
- `assert.notContains(haystack, needle)` - String doesn't contain
- `assert.greater(a, b)` - Numeric comparison
- `assert.less(a, b)` - Numeric comparison
- `assert.hasLength(collection, length)` - Collection length
- `assert.empty(collection)` - Empty collection
- `assert.notEmpty(collection)` - Non-empty collection
- `assert.panics(func)` - Exception testing
- `assert.notPanics(func)` - No exception testing

### 4. Configuration Options
The framework supports extensive configuration through JSON files:

- Test discovery patterns
- Timeout settings
- Parallel execution
- Coverage reporting
- Custom reporters
- Tag-based filtering

## Running Tests

### Using the CLI
```bash
# Run all tests
r2test

# Run with specific configuration
r2test -config r2test.config.json

# Run with verbose output
r2test -verbose

# Run tests in parallel
r2test -parallel -workers 4

# Run only tests matching a pattern
r2test -grep "Calculator"

# Run with coverage
r2test -coverage -coverage-dir ./coverage
```

### Using Environment Variables
```bash
# Set timeout
export R2TEST_TIMEOUT=60s

# Enable parallel execution
export R2TEST_PARALLEL=true

# Enable verbose output
export R2TEST_VERBOSE=true

# Enable coverage
export R2TEST_COVERAGE=true
```

## Migration from BDD

The old BDD syntax has been completely removed:

### Old BDD Syntax (Removed)
```r2
TestCase "My Test" {
    Given func() { /* setup */ }
    When func() { /* action */ }
    Then func() { /* assertion */ }
    And func() { /* additional step */ }
}
```

### New Unit Testing Syntax
```r2
describe("My Test Suite", func() {
    beforeEach(func() {
        // setup
    });
    
    it("should perform the action correctly", func() {
        // action and assertion
        let result = performAction();
        assert.equals(result, expected);
    });
});
```

## Best Practices

1. **Descriptive Names**: Use clear, descriptive names for test suites and test cases
2. **Single Responsibility**: Each test should verify one specific behavior
3. **Setup and Teardown**: Use hooks for common setup and cleanup
4. **Independent Tests**: Tests should not depend on each other
5. **Clear Assertions**: Use specific assertions that provide good error messages

## Future Features

The framework is designed to support future enhancements:

- Snapshot testing
- Retry mechanisms for flaky tests
- Performance testing utilities
- Mocking and stubbing system
- Parallel test execution
- Watch mode for development
- CI/CD integration

## Examples Overview

### Basic Examples (`basic_test.r2`)
- Simple mathematical operations
- String manipulation
- Array operations
- Boolean logic
- Function testing
- Error handling

### Advanced Examples (`advanced_test.r2`)
- Class testing with setup/teardown
- File operations (placeholder)
- Async operations (placeholder)
- Complex data structures
- Edge cases
- Performance testing

These examples demonstrate the power and flexibility of the new R2Lang Unit Testing Framework while maintaining simplicity and ease of use.