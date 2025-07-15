import * as vscode from 'vscode';
import * as path from 'path';
import * as fs from 'fs';
import { exec } from 'child_process';

export function activate(context: vscode.ExtensionContext) {
    console.log('R2Lang extension is now active!');

    // Register commands
    const runFileCommand = vscode.commands.registerCommand('r2lang.runFile', () => {
        runR2LangFile();
    });

    const runSelectionCommand = vscode.commands.registerCommand('r2lang.runSelection', () => {
        runR2LangSelection();
    });

    const openReplCommand = vscode.commands.registerCommand('r2lang.openRepl', () => {
        openR2LangRepl();
    });

    const runTestsCommand = vscode.commands.registerCommand('r2lang.runTests', () => {
        runR2LangTests();
    });

    const runCurrentTestCommand = vscode.commands.registerCommand('r2lang.runCurrentTest', () => {
        runCurrentTestFile();
    });

    const runTestsWithCoverageCommand = vscode.commands.registerCommand('r2lang.runTestsWithCoverage', () => {
        runR2LangTestsWithCoverage();
    });

    const runIndividualTestCommand = vscode.commands.registerCommand('r2lang.runIndividualTest', (testName: string, filePath: string) => {
        runIndividualTest(testName, filePath);
    });

    // Register providers
    const documentFormattingProvider = vscode.languages.registerDocumentFormattingEditProvider(
        'r2lang',
        new R2LangFormattingProvider()
    );

    const hoverProvider = vscode.languages.registerHoverProvider('r2lang', new R2LangHoverProvider());

    // Register CodeLens provider for test execution
    const codeLensProvider = vscode.languages.registerCodeLensProvider('r2lang', new R2LangTestCodeLensProvider());

    // Add to context subscriptions
    context.subscriptions.push(
        runFileCommand,
        runSelectionCommand,
        openReplCommand,
        runTestsCommand,
        runCurrentTestCommand,
        runTestsWithCoverageCommand,
        runIndividualTestCommand,
        documentFormattingProvider,
        hoverProvider,
        codeLensProvider
    );

    // Show welcome message on first use
    const config = vscode.workspace.getConfiguration('r2lang');
    const hasShownWelcome = context.globalState.get('hasShownWelcome', false);
    
    if (!hasShownWelcome) {
        showWelcomeMessage(context);
    }
}

export function deactivate() {
    console.log('R2Lang extension is now deactivated');
}

function runR2LangFile() {
    const editor = vscode.window.activeTextEditor;
    if (!editor) {
        vscode.window.showErrorMessage('No active R2Lang file found');
        return;
    }

    const document = editor.document;
    if (document.languageId !== 'r2lang') {
        vscode.window.showErrorMessage('Current file is not a R2Lang file');
        return;
    }

    const filePath = document.fileName;
    const config = vscode.workspace.getConfiguration('r2lang');
    const executablePath = config.get<string>('executablePath', 'r2lang');

    // Save the file before running
    document.save().then(() => {
        const terminal = vscode.window.createTerminal('R2Lang');
        terminal.show();
        terminal.sendText(`"${executablePath}" "${filePath}"`);
    });
}

function runR2LangSelection() {
    const editor = vscode.window.activeTextEditor;
    if (!editor) {
        vscode.window.showErrorMessage('No active editor found');
        return;
    }

    const selection = editor.selection;
    const selectedText = editor.document.getText(selection);
    
    if (!selectedText.trim()) {
        vscode.window.showErrorMessage('No text selected');
        return;
    }

    // Create a temporary file with selected content
    const config = vscode.workspace.getConfiguration('r2lang');
    const executablePath = config.get<string>('executablePath', 'r2lang');
    
    const tempFilePath = path.join(__dirname, 'temp_selection.r2');
    const fs = require('fs');
    
    fs.writeFileSync(tempFilePath, selectedText);
    
    const terminal = vscode.window.createTerminal('R2Lang Selection');
    terminal.show();
    terminal.sendText(`"${executablePath}" "${tempFilePath}"`);
    
    // Clean up temp file after a delay
    setTimeout(() => {
        try {
            fs.unlinkSync(tempFilePath);
        } catch (error) {
            // Ignore cleanup errors
        }
    }, 5000);
}

function openR2LangRepl() {
    const config = vscode.workspace.getConfiguration('r2lang');
    const replExecutablePath = config.get<string>('replExecutablePath', 'r2repl');
    
    const terminal = vscode.window.createTerminal('R2Lang REPL');
    terminal.show();
    terminal.sendText(`"${replExecutablePath}"`);
}

function showWelcomeMessage(context: vscode.ExtensionContext) {
    const message = 'Welcome to R2Lang! Get started by creating a .r2 file or opening the REPL.';
    const options = ['Create New File', 'Open REPL', 'View Examples', 'Don\'t show again'];
    
    vscode.window.showInformationMessage(message, ...options).then(selection => {
        switch (selection) {
            case 'Create New File':
                createNewR2File();
                break;
            case 'Open REPL':
                openR2LangRepl();
                break;
            case 'View Examples':
                openExamplesPage();
                break;
            case 'Don\'t show again':
                context.globalState.update('hasShownWelcome', true);
                break;
        }
    });
}

function createNewR2File() {
    const template = `func main() {
    print("Hello, R2Lang!");
}`;

    vscode.workspace.openTextDocument({
        content: template,
        language: 'r2lang'
    }).then(doc => {
        vscode.window.showTextDocument(doc);
    });
}

function openExamplesPage() {
    vscode.env.openExternal(vscode.Uri.parse('https://github.com/arturoeanton/go-r2lang/tree/main/examples'));
}

class R2LangFormattingProvider implements vscode.DocumentFormattingEditProvider {
    provideDocumentFormattingEdits(
        document: vscode.TextDocument,
        options: vscode.FormattingOptions,
        token: vscode.CancellationToken
    ): vscode.TextEdit[] {
        // Basic formatting implementation
        const edits: vscode.TextEdit[] = [];
        
        for (let i = 0; i < document.lineCount; i++) {
            const line = document.lineAt(i);
            const trimmedText = line.text.trim();
            
            if (trimmedText && line.text !== trimmedText) {
                const range = new vscode.Range(
                    line.range.start,
                    line.range.end
                );
                edits.push(vscode.TextEdit.replace(range, trimmedText));
            }
        }
        
        return edits;
    }
}

class R2LangHoverProvider implements vscode.HoverProvider {
    provideHover(
        document: vscode.TextDocument,
        position: vscode.Position,
        token: vscode.CancellationToken
    ): vscode.ProviderResult<vscode.Hover> {
        const range = document.getWordRangeAtPosition(position);
        const word = document.getText(range);
        
        // Provide hover information for keywords
        const keywordInfo = getKeywordInfo(word);
        if (keywordInfo) {
            return new vscode.Hover(
                new vscode.MarkdownString(keywordInfo),
                range
            );
        }
        
        return null;
    }
}

function getKeywordInfo(word: string): string | null {
    const keywords: { [key: string]: string } = {
        'func': 'Declares a function\n\n```r2lang\nfunc myFunction(param1, param2) {\n    return param1 + param2;\n}\n```',
        'class': 'Declares a class\n\n```r2lang\nclass MyClass {\n    constructor() {\n        this.value = 0;\n    }\n}\n```',
        'describe': 'Defines a test suite\n\n```r2lang\ndescribe("Calculator", func() {\n    it("should add numbers", func() {\n        assert.equals(2 + 2, 4);\n    });\n});\n```',
        'it': 'Defines a test case\n\n```r2lang\nit("should handle edge cases", func() {\n    assert.true(condition);\n});\n```',
        'assert': 'Test assertion functions\n\n```r2lang\nassert.equals(actual, expected);\nassert.true(condition);\nassert.false(condition);\n```',
        'if': 'Conditional statement\n\n```r2lang\nif (condition) {\n    // code\n}\n```',
        'while': 'Loop statement\n\n```r2lang\nwhile (condition) {\n    // code\n}\n```',
        'for': 'Loop statement\n\n```r2lang\nfor (let i = 0; i < 10; i++) {\n    // code\n}\n```',
        'import': 'Imports a module\n\n```r2lang\nimport "module.r2" as mod\n```',
        'let': 'Declares a variable\n\n```r2lang\nlet variableName = value;\n```',
        'this': 'References the current object instance',
        'super': 'References the parent class',
        'extends': 'Inherits from a parent class\n\n```r2lang\nclass Child extends Parent {\n    // class body\n}\n```'
    };
    
    return keywords[word] || null;
}

// Helper function to find r2test binary
function findR2TestBinary(): string {
    const config = vscode.workspace.getConfiguration('r2lang');
    const configuredPath = config.get<string>('testExecutablePath', '');
    
    // If configured path is absolute, use it
    if (configuredPath && path.isAbsolute(configuredPath)) {
        return configuredPath;
    }
    
    // Try to find r2test in the project root
    const workspaceFolder = vscode.workspace.workspaceFolders?.[0];
    if (workspaceFolder) {
        const projectR2Test = path.join(workspaceFolder.uri.fsPath, 'r2test');
        const projectR2TestExe = path.join(workspaceFolder.uri.fsPath, 'r2test.exe');
        
        if (fs.existsSync(projectR2Test)) {
            return projectR2Test;
        }
        if (fs.existsSync(projectR2TestExe)) {
            return projectR2TestExe;
        }
        
        // Try in cmd/r2test directory
        const cmdR2Test = path.join(workspaceFolder.uri.fsPath, 'cmd', 'r2test', 'r2test');
        const cmdR2TestExe = path.join(workspaceFolder.uri.fsPath, 'cmd', 'r2test', 'r2test.exe');
        
        if (fs.existsSync(cmdR2Test)) {
            return cmdR2Test;
        }
        if (fs.existsSync(cmdR2TestExe)) {
            return cmdR2TestExe;
        }
    }
    
    // Fall back to configured path or system PATH
    return configuredPath || 'r2test';
}

function runR2LangTests() {
    const testExecutablePath = findR2TestBinary();
    
    const terminal = vscode.window.createTerminal('R2Lang Tests');
    terminal.show();
    terminal.sendText(`"${testExecutablePath}"`);
}

function runCurrentTestFile() {
    const editor = vscode.window.activeTextEditor;
    if (!editor) {
        vscode.window.showErrorMessage('No active R2Lang test file found');
        return;
    }

    const document = editor.document;
    if (document.languageId !== 'r2lang') {
        vscode.window.showErrorMessage('Current file is not a R2Lang file');
        return;
    }

    const filePath = document.fileName;
    if (!filePath.endsWith('_test.r2')) {
        vscode.window.showErrorMessage('Current file is not a R2Lang test file (must end with _test.r2)');
        return;
    }

    const testExecutablePath = findR2TestBinary();

    // Save the file before running
    document.save().then(() => {
        const terminal = vscode.window.createTerminal('R2Lang Current Test');
        terminal.show();
        terminal.sendText(`"${testExecutablePath}" "${path.dirname(filePath)}"`);
    });
}

function runR2LangTestsWithCoverage() {
    const testExecutablePath = findR2TestBinary();
    
    const terminal = vscode.window.createTerminal('R2Lang Tests with Coverage');
    terminal.show();
    terminal.sendText(`"${testExecutablePath}" -coverage -coverage-formats html,json`);
}

function runIndividualTest(testName: string, filePath: string) {
    const testExecutablePath = findR2TestBinary();
    
    // Check if the test binary exists
    if (!testExecutablePath.includes('/') && !testExecutablePath.includes('\\')) {
        // It's a system command, let's try to verify it exists
        exec(`"${testExecutablePath}" --version`, (error) => {
            if (error) {
                vscode.window.showErrorMessage(
                    `R2Test executable not found. Please build it first: go build -o r2test cmd/r2test/main.go`
                );
                return;
            }
        });
    }
    
    const terminal = vscode.window.createTerminal(`R2Lang Test: ${testName}`);
    terminal.show();
    terminal.sendText(`"${testExecutablePath}" -grep "${testName}" "${path.dirname(filePath)}"`);
    
    // Show notification
    vscode.window.showInformationMessage(`Running test: ${testName}`);
}

// CodeLens provider for test execution
class R2LangTestCodeLensProvider implements vscode.CodeLensProvider {
    provideCodeLenses(document: vscode.TextDocument, token: vscode.CancellationToken): vscode.CodeLens[] | Thenable<vscode.CodeLens[]> {
        const codeLenses: vscode.CodeLens[] = [];
        
        // Check if test CodeLens is enabled
        const config = vscode.workspace.getConfiguration('r2lang');
        if (!config.get<boolean>('enableTestCodeLens', true)) {
            return codeLenses;
        }
        
        // Only provide CodeLens for test files
        if (!document.fileName.endsWith('_test.r2')) {
            return codeLenses;
        }
        
        const text = document.getText();
        const lines = text.split('\n');
        
        for (let i = 0; i < lines.length; i++) {
            const line = lines[i];
            
            // Look for describe() blocks - handle multi-line definitions
            const describeMatch = line.match(/^\s*describe\s*\(\s*["']([^"']+)["']\s*,/);
            if (describeMatch) {
                const testSuiteName = describeMatch[1];
                const range = new vscode.Range(i, 0, i, line.length);
                
                const command: vscode.Command = {
                    title: '▶️ Run Test Suite',
                    command: 'r2lang.runIndividualTest',
                    arguments: [testSuiteName, document.fileName]
                };
                
                codeLenses.push(new vscode.CodeLens(range, command));
            }
            
            // Look for it() blocks - handle multi-line definitions
            const itMatch = line.match(/^\s*it\s*\(\s*["']([^"']+)["']\s*,/);
            if (itMatch) {
                const testName = itMatch[1];
                const range = new vscode.Range(i, 0, i, line.length);
                
                const command: vscode.Command = {
                    title: '▶️ Run Test',
                    command: 'r2lang.runIndividualTest',
                    arguments: [testName, document.fileName]
                };
                
                codeLenses.push(new vscode.CodeLens(range, command));
            }
        }
        
        return codeLenses;
    }
}