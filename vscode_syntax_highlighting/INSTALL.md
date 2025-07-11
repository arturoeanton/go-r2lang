# Installation Guide - R2Lang VS Code Extension

## Quick Installation

### Option 1: Install from Source (Recommended for Development)

1. **Clone the repository:**
   ```bash
   git clone https://github.com/arturoeanton/go-r2lang.git
   cd go-r2lang/vscode_syntax_highlighting
   ```

2. **Install Node.js dependencies:**
   ```bash
   npm install
   ```

3. **Compile the extension:**
   ```bash
   npm run compile
   ```

4. **Package the extension:**
   ```bash
   npm run package
   ```
   This creates a `.vsix` file in the current directory.

5. **Install in VS Code:**
   - Open VS Code
   - Press `Ctrl+Shift+P` (or `Cmd+Shift+P` on Mac)
   - Type "Extensions: Install from VSIX..."
   - Select the generated `.vsix` file

### Option 2: Development Mode

For extension development and testing:

1. **Open in VS Code:**
   ```bash
   cd vscode_syntax_highlighting
   code .
   ```

2. **Run in debug mode:**
   - Press `F5` in VS Code
   - A new "Extension Development Host" window opens
   - Create or open a `.r2` file to test the extension

## Prerequisites

### 1. Visual Studio Code
- **Minimum version:** 1.74.0 or higher
- Download from: https://code.visualstudio.com/

### 2. Node.js (for building from source)
- **Minimum version:** 16.x or higher
- Download from: https://nodejs.org/

### 3. R2Lang Interpreter
The extension requires the R2Lang interpreter to be installed:

```bash
# Clone R2Lang repository
git clone https://github.com/arturoeanton/go-r2lang.git
cd go-r2lang

# Build the interpreter (requires Go 1.23+)
go build -o r2lang main.go

# Add to PATH or note the full path for configuration
```

## Configuration

After installation, configure the extension:

1. **Open VS Code Settings:**
   - Press `Ctrl+,` (or `Cmd+,` on Mac)
   - Or go to File → Preferences → Settings

2. **Search for "r2lang"**

3. **Configure the executable path:**
   ```json
   {
     "r2lang.executablePath": "/path/to/your/r2lang"
   }
   ```

### Common Configurations

**Windows:**
```json
{
  "r2lang.executablePath": "C:\\path\\to\\r2lang.exe"
}
```

**macOS:**
```json
{
  "r2lang.executablePath": "/usr/local/bin/r2lang"
}
```

**Linux:**
```json
{
  "r2lang.executablePath": "/usr/local/bin/r2lang"
}
```

## Verification

Test your installation:

1. **Create a test file:**
   - Create a new file: `test.r2`
   - Add this content:
     ```r2lang
     func main() {
         print("Hello, R2Lang!");
     }
     ```

2. **Verify syntax highlighting:**
   - Keywords should be colored (blue for `func`, etc.)
   - Strings should be colored (orange/yellow)
   - Comments should be colored (green)

3. **Test code execution:**
   - Press `Ctrl+F5` (or `Cmd+F5` on Mac)
   - Should see output in the terminal

4. **Test snippets:**
   - Type `func` and press `Tab`
   - Should expand to a function template

## Troubleshooting

### Common Issues

**1. Extension not loading:**
```bash
# Check VS Code version
code --version

# Reinstall dependencies
npm clean-install
npm run compile
```

**2. Syntax highlighting not working:**
- Ensure file has `.r2` extension
- Reload window: `Ctrl+Shift+P` → "Developer: Reload Window"

**3. Code execution fails:**
```bash
# Test R2Lang interpreter directly
r2lang --version
r2lang examples/example1-if.r2

# Check PATH
which r2lang  # Linux/Mac
where r2lang  # Windows
```

**4. Snippets not appearing:**
- Enable IntelliSense in VS Code settings
- Try `Ctrl+Space` to trigger suggestions

### Debug Mode

Enable detailed logging:

1. **Open VS Code Developer Tools:**
   - `Ctrl+Shift+P` → "Developer: Toggle Developer Tools"

2. **Check Console for errors:**
   - Look for R2Lang extension messages
   - Note any error messages

3. **Extension Host Logs:**
   - `Ctrl+Shift+P` → "Developer: Show Logs..."
   - Select "Extension Host"

## Development Setup

For contributing to the extension:

### Requirements
- Node.js 16+
- TypeScript knowledge
- VS Code Extension API familiarity

### Commands
```bash
# Install dependencies
npm install

# Compile TypeScript
npm run compile

# Watch mode (auto-compile)
npm run watch

# Lint code
npm run lint

# Run tests
npm run test

# Package for distribution
npm run package

# Publish to marketplace
npm run publish
```

### Project Structure
```
vscode_syntax_highlighting/
├── src/
│   └── extension.ts          # Main extension code
├── syntaxes/
│   └── r2lang.tmGrammar.json # Syntax highlighting rules
├── snippets/
│   └── r2lang-snippets.json  # Code snippets
├── themes/
│   └── r2lang-dark-color-theme.json
├── language-configuration/
│   └── r2lang-configuration.json
├── package.json              # Extension manifest
├── README.md                 # Documentation
└── CHANGELOG.md              # Version history
```

## Advanced Configuration

### Workspace Settings

For project-specific settings, create `.vscode/settings.json`:

```json
{
  "r2lang.executablePath": "./bin/r2lang",
  "r2lang.enableDiagnostics": true,
  "r2lang.maxNumberOfProblems": 50,
  "files.associations": {
    "*.r2": "r2lang"
  }
}
```

### Task Configuration

Add R2Lang tasks to `.vscode/tasks.json`:

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Run R2Lang File",
      "type": "shell",
      "command": "r2lang",
      "args": ["${file}"],
      "group": {
        "kind": "build",
        "isDefault": true
      },
      "presentation": {
        "echo": true,
        "reveal": "always",
        "focus": false,
        "panel": "shared"
      }
    }
  ]
}
```

### Launch Configuration

Add debugging configuration to `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Run R2Lang",
      "type": "node",
      "request": "launch",
      "program": "r2lang",
      "args": ["${file}"],
      "console": "integratedTerminal",
      "cwd": "${workspaceFolder}"
    }
  ]
}
```

## Support

For help and support:

- **GitHub Issues:** https://github.com/arturoeanton/go-r2lang/issues
- **Documentation:** https://github.com/arturoeanton/go-r2lang/tree/main/docs
- **Examples:** https://github.com/arturoeanton/go-r2lang/tree/main/examples

## Contributing

See [CONTRIBUTING.md](https://github.com/arturoeanton/go-r2lang/blob/main/CONTRIBUTING.md) for development guidelines.