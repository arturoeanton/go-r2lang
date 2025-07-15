# R2Lang Testing Guide

## Fixed Issues

### ✅ Grep Filtering Now Works
The `-grep` option now properly filters tests by name/description:

```bash
# Run only tests containing "Calculator" in name or description
r2test -grep "Calculator"

# Run only tests containing "add" in name or description  
r2test -grep "add"

# Run specific test by exact name match
r2test -grep "should add two numbers"
```

### ✅ Verbose Output Now Works
The `-verbose` option now shows detailed test execution information:

```bash
# Run tests with verbose output
r2test -verbose

# Run specific test with verbose output
r2test -verbose -grep "Calculator"
```

## VS Code Extension Enhanced

### CodeLens Test Execution
When you open a test file (ending with `_test.r2`), you'll see:
- **▶️ Run Test Suite** - above each `describe()` block
- **▶️ Run Test** - above each `it()` block

### Smart Binary Detection
The extension automatically finds r2test in:
1. Project root directory (`./r2test`)
2. Command directory (`./cmd/r2test/r2test`) 
3. System PATH (`r2test`)

### Keyboard Shortcuts
- `Ctrl+Shift+T` / `Cmd+Shift+T` - Run all tests
- `Ctrl+F6` / `Cmd+F6` - Run current test file
- `Ctrl+Shift+F6` / `Cmd+Shift+F6` - Run tests with coverage

## Example Usage

### Build the test runner
```bash
go build -o r2test cmd/r2test/main.go
```

### Run all tests
```bash
./r2test
```

### Filter tests by name
```bash
./r2test -grep "Calculator"
```

### Run with verbose output
```bash
./r2test -verbose
```

### Run with coverage
```bash
./r2test -coverage -coverage-formats html,json
```

### Run specific directory
```bash
./r2test ./examples/testing
```

## Test File Example

```r2
describe("Calculator Tests", func() {
    it("should add two numbers", func() {
        let result = 2 + 3;
        assert.equals(result, 5);
    });
    
    it("should subtract numbers", func() {
        let result = 10 - 4;
        assert.equals(result, 6);
    });
});
```

## Configuration

Create a `r2test.config.json` file:

```json
{
  "testDirs": ["./examples/testing", "./tests"],
  "patterns": ["*_test.r2"],
  "verbose": true,
  "coverage": {
    "enabled": true,
    "threshold": 80,
    "formats": ["html", "json"]
  }
}
```

## VS Code Configuration

```json
{
  "r2lang.testExecutablePath": "./r2test",
  "r2lang.enableTestCodeLens": true
}
```

Now the testing framework works as expected with proper filtering and verbose output!