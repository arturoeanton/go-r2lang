import * as vscode from 'vscode';
import * as path from 'path';
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

    // Register providers
    const documentFormattingProvider = vscode.languages.registerDocumentFormattingEditProvider(
        'r2lang',
        new R2LangFormattingProvider()
    );

    const hoverProvider = vscode.languages.registerHoverProvider('r2lang', new R2LangHoverProvider());

    // Add to context subscriptions
    context.subscriptions.push(
        runFileCommand,
        runSelectionCommand,
        openReplCommand,
        documentFormattingProvider,
        hoverProvider
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
    const executablePath = config.get<string>('executablePath', 'r2lang');
    
    const terminal = vscode.window.createTerminal('R2Lang REPL');
    terminal.show();
    terminal.sendText(`"${executablePath}" -repl`);
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
        'TestCase': 'Defines a BDD test case\n\n```r2lang\nTestCase "My Test" {\n    Given setup()\n    When action()\n    Then assertion()\n}\n```',
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