# Changelog

All notable changes to the "R2Lang" extension will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2024-01-XX

### Added
- Initial release of R2Lang VS Code extension
- Complete syntax highlighting for R2Lang language
- Support for all R2Lang keywords, operators, and constructs
- Special highlighting for BDD testing constructs (TestCase, Given, When, Then, And)
- Comprehensive snippet library with 30+ code snippets
- Code execution commands (Run File, Run Selection, Open REPL)
- Language configuration with auto-closing pairs and indentation
- Custom R2Lang Dark theme optimized for the language
- Basic formatting provider
- Hover information for keywords and built-in functions
- File association for .r2 files
- Terminal integration for code execution
- Welcome message for first-time users
- Configuration options for executable path and features

### Language Features
- Function and class declarations syntax highlighting
- Object-oriented programming constructs (class, extends, this, super)
- Control flow statements (if, else, while, for, try, catch, finally)
- Import/export statements
- String literals with escape sequence support
- Numeric literals (integers and floats)
- Comment support (line and block comments)
- BDD testing syntax (TestCase blocks)
- HTTP server route definitions
- Array and object literal syntax
- Operator highlighting (arithmetic, comparison, logical, assignment)

### Code Snippets
- Function declarations (func, main, anonymous functions)
- Class definitions (class, extends, constructor)
- Control structures (if, while, for, try-catch)
- BDD test cases (testcase, test-full)
- HTTP routes (httpget, httppost, httpserver)
- Data structures (array, object)
- Common statements (import, print, variable declarations)
- Comments (block, TODO, FIXME)

### Developer Experience
- Keyboard shortcuts for common actions
- Context menu integration
- Command palette commands
- Auto-completion for language constructs
- Bracket matching and auto-closing
- Code folding support
- Proper indentation handling

### Technical Implementation
- TextMate grammar for syntax highlighting
- TypeScript-based extension with VS Code API
- Language server foundation for future enhancements
- Comprehensive test coverage
- Cross-platform compatibility (Windows, macOS, Linux)

### Documentation
- Comprehensive README with setup instructions
- Code snippet reference guide
- Troubleshooting section
- Configuration options documentation
- Development setup guide
- Contributing guidelines

## [0.0.1] - 2024-01-XX

### Added
- Project setup and initial structure
- Basic package.json configuration
- Development dependencies setup

---

### Legend
- `Added` for new features
- `Changed` for changes in existing functionality
- `Deprecated` for soon-to-be removed features
- `Removed` for now removed features
- `Fixed` for any bug fixes
- `Security` in case of vulnerabilities