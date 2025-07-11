# R2Lang Extension for Visual Studio Code

[![Version](https://img.shields.io/badge/version-0.1.0-blue.svg)](https://marketplace.visualstudio.com/items?itemName=r2lang.r2lang)
[![License](https://img.shields.io/badge/license-Apache%202.0-green.svg)](https://opensource.org/licenses/Apache-2.0)

The official Visual Studio Code extension for R2Lang - a dynamic, interpreted programming language with built-in BDD testing support.

## Features

### 🎨 Syntax Highlighting
- **Complete syntax highlighting** for R2Lang keywords, strings, numbers, comments, and operators
- **Special highlighting** for BDD test constructs (`TestCase`, `Given`, `When`, `Then`, `And`)
- **Class and function highlighting** with inheritance support
- **Custom color theme** optimized for R2Lang code

### ✨ Code Snippets
Rich snippet library including:
- **Function declarations** (`func`, `main`)
- **Class definitions** (`class`, `extends`)
- **Control structures** (`if`, `while`, `for`, `try-catch`)
- **BDD test cases** (`testcase`, `test-full`)
- **HTTP server routes** (`httpget`, `httppost`)
- **Common patterns** (`import`, `print`, `array`, `object`)

### 🚀 Code Execution
- **Run current file** (Ctrl+F5 / Cmd+F5)
- **Run selected code** (Ctrl+Shift+F5 / Cmd+Shift+F5)
- **Open REPL** (Ctrl+Shift+R / Cmd+Shift+R)
- **Terminal integration** with automatic file execution

### 🔧 Language Features
- **Auto-closing pairs** for brackets, quotes, and comments
- **Automatic indentation** with smart bracket handling
- **Code folding** with region markers
- **Comment toggling** (line and block comments)
- **Basic formatting** with consistent indentation

### 📖 IntelliSense Support
- **Hover information** for keywords and built-in functions
- **Context-aware suggestions** (coming soon)
- **Error highlighting** (coming soon)
- **Go to definition** (coming soon)

## Installation

### From VS Code Marketplace
1. Open Visual Studio Code
2. Go to the Extensions view (Ctrl+Shift+X)
3. Search for "R2Lang"
4. Click Install

### From VSIX Package
1. Download the latest `.vsix` file from [releases](https://github.com/arturoeanton/go-r2lang/releases)
2. Open VS Code
3. Run `Extensions: Install from VSIX...` from the Command Palette
4. Select the downloaded `.vsix` file

### Building from Source
1. Clone the repository:
   ```bash
   git clone https://github.com/arturoeanton/go-r2lang.git
   cd go-r2lang/vscode_syntax_highlighting
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Compile the extension:
   ```bash
   npm run compile
   ```

4. Package the extension:
   ```bash
   npm run package
   ```

## Requirements

- **R2Lang Interpreter**: You need to have the R2Lang interpreter installed and accessible from your PATH
- **Visual Studio Code**: Version 1.74.0 or higher

### Installing R2Lang
1. Clone the R2Lang repository:
   ```bash
   git clone https://github.com/arturoeanton/go-r2lang.git
   cd go-r2lang
   ```

2. Build the interpreter:
   ```bash
   go build -o r2lang main.go
   ```

3. Add to PATH or configure the extension settings

## Getting Started

### Quick Start
1. Create a new file with `.r2` extension
2. Start typing R2Lang code - syntax highlighting will activate automatically
3. Use `Ctrl+F5` (or `Cmd+F5` on Mac) to run your code

### Your First R2Lang Program
```r2lang
func main() {
    print("Hello, R2Lang!");
    
    let numbers = [1, 2, 3, 4, 5];
    for (let num in numbers) {
        print("Number:", num);
    }
}
```

### BDD Testing Example
```r2lang
TestCase "User Registration" {
    Given func() {
        setupDatabase();
        return "Database initialized";
    }
    When func() {
        let result = registerUser("test@example.com", "password123");
        return "User registration attempted";
    }
    Then func() {
        assertEqual(result.success, true);
        return "Registration should succeed";
    }
    And func() {
        assertEqual(result.user.email, "test@example.com");
        return "Email should be stored correctly";
    }
}
```

## Configuration

Configure the extension through VS Code settings:

### Available Settings

| Setting | Default | Description |
|---------|---------|-------------|
| `r2lang.executablePath` | `"r2lang"` | Path to the R2Lang executable |
| `r2lang.enableCodeLens` | `true` | Enable/disable CodeLens for functions and classes |
| `r2lang.enableAutoCompletion` | `true` | Enable/disable auto-completion suggestions |
| `r2lang.enableDiagnostics` | `true` | Enable/disable syntax and semantic error checking |
| `r2lang.maxNumberOfProblems` | `100` | Maximum number of problems to report per file |

### Example Configuration
```json
{
    "r2lang.executablePath": "/usr/local/bin/r2lang",
    "r2lang.enableCodeLens": true,
    "r2lang.enableAutoCompletion": true
}
```

## Keyboard Shortcuts

| Command | Windows/Linux | macOS | Description |
|---------|---------------|-------|-------------|
| Run File | `Ctrl+F5` | `Cmd+F5` | Execute the current R2Lang file |
| Run Selection | `Ctrl+Shift+F5` | `Cmd+Shift+F5` | Execute selected R2Lang code |
| Open REPL | `Ctrl+Shift+R` | `Cmd+Shift+R` | Open R2Lang interactive shell |

## Code Snippets Reference

### Functions
- `func` → Function declaration
- `main` → Main function
- `afunc` → Anonymous function

### Classes
- `class` → Class declaration
- `extends` → Class with inheritance
- `constructor` → Constructor method

### Control Flow
- `if` → If statement
- `ifelse` → If-else statement
- `while` → While loop
- `for` → For loop
- `forin` → For-in loop

### Testing
- `testcase` → Basic BDD test case
- `test-full` → Full BDD test case with And clause
- `assert` → Assertion statement

### Web Development
- `httpserver` → HTTP server setup
- `httpget` → GET route handler
- `httppost` → POST route handler

### Data Structures
- `array` → Array declaration
- `object` → Object declaration
- `import` → Import statement

## Troubleshooting

### Common Issues

**1. "R2Lang executable not found"**
- Ensure R2Lang is installed and in your PATH
- Configure `r2lang.executablePath` in settings to point to the correct executable

**2. "Syntax highlighting not working"**
- Make sure your file has the `.r2` extension
- Try reloading the window (`Ctrl+Shift+P` → "Developer: Reload Window")

**3. "Code execution fails"**
- Check that the R2Lang interpreter is working: run `r2lang --version` in terminal
- Verify file path doesn't contain special characters
- Ensure file is saved before execution

**4. "Snippets not appearing"**
- Check that IntelliSense is enabled in VS Code settings
- Try typing snippet prefixes and pressing `Ctrl+Space`

### Debug Mode
Enable extension development debugging:
1. Open extension folder in VS Code
2. Press `F5` to start debugging
3. A new VS Code window opens with the extension loaded
4. Check the Debug Console for extension logs

## Contributing

We welcome contributions! Here's how to get started:

### Development Setup
1. Fork and clone the repository
2. Install dependencies: `npm install`
3. Make your changes
4. Test the extension: `F5` in VS Code
5. Submit a pull request

### Areas for Contribution
- **Language Server Protocol** implementation
- **Advanced IntelliSense** features
- **Debugging support**
- **Code formatting** improvements
- **Error detection** and diagnostics
- **Refactoring tools**

### Reporting Issues
Please report bugs and feature requests on [GitHub Issues](https://github.com/arturoeanton/go-r2lang/issues).

## Roadmap

### Version 0.2.0 (Planned)
- [ ] Language Server Protocol support
- [ ] Advanced auto-completion
- [ ] Error detection and diagnostics
- [ ] Go to definition
- [ ] Symbol search

### Version 0.3.0 (Planned)
- [ ] Integrated debugger
- [ ] Code formatting
- [ ] Refactoring tools
- [ ] Test runner integration

### Version 1.0.0 (Planned)
- [ ] Full IntelliSense support
- [ ] Performance optimizations
- [ ] Advanced testing features
- [ ] Plugin ecosystem

## License

This extension is licensed under the [Apache License 2.0](https://opensource.org/licenses/Apache-2.0).

## Links

- **Repository**: [github.com/arturoeanton/go-r2lang](https://github.com/arturoeanton/go-r2lang)
- **Issues**: [github.com/arturoeanton/go-r2lang/issues](https://github.com/arturoeanton/go-r2lang/issues)
- **Examples**: [github.com/arturoeanton/go-r2lang/tree/main/examples](https://github.com/arturoeanton/go-r2lang/tree/main/examples)
- **Documentation**: [R2Lang Documentation](https://github.com/arturoeanton/go-r2lang/tree/main/docs)

---

**Happy coding with R2Lang! 🚀**