// r2repl.go
package r2repl

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
	"github.com/arturoeanton/go-r2lang/pkg/r2libs"
	"github.com/fatih/color"
	"github.com/peterh/liner"
)

// Definir funciones de color
var (
	errorColor   = color.New(color.FgRed, color.Bold).SprintFunc()
	successColor = color.New(color.FgGreen).SprintFunc()
	infoColor    = color.New(color.FgHiBlue).SprintFunc()
	warningColor = color.New(color.FgYellow).SprintFunc()
	keywordColor = color.New(color.FgMagenta, color.Bold).SprintFunc() // Para palabras reservadas
	stringColor  = color.New(color.FgGreen).SprintFunc()               // Para cadenas de texto
)

// Definir las palabras reservadas de R2Lang
var reservedKeywords = []string{
	"let", "var", "func", "function", "class", "method", "if", "else", "return",
	"for", "while", "break", "continue",
	"switch", "case", "default", "import", "export", "var", "const", "true",
	"false", "nil", "null", "testcase", "give", "when", "then", "and", "or",
	"throw", "try", "catch", "finally", "this", "super",
}

// Compilar una expresión regular para palabras reservadas
var keywordRegex = regexp.MustCompile(`\b(` + strings.Join(reservedKeywords, "|") + `)\b`)

// Compilar una expresión regular para cadenas de texto (simples y dobles)
var stringRegex = regexp.MustCompile(`"([^"\\]|\\.)*"|'([^'\\]|\\.)*'`)

// Repl inicia el Read-Eval-Print Loop para R2Lang
func Repl(outputFlag bool) {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)
	loadHistory(line)
	// Deferred (not inside loadHistory): must run when the REPL session
	// itself ends, not right after startup loading finishes, otherwise
	// nothing typed during the session is ever persisted.
	defer saveHistory(line)

	fmt.Println(infoColor("Welcome to R2Lang REPL!"))

	env := createR2Environment()

	var buffer strings.Builder
	var commandHistory []string

	commands := getCommands(&commandHistory)

	for {
		prompt := "> "
		if buffer.Len() > 0 {
			prompt = "... "
		}

		input, err := line.Prompt(prompt)
		if err == liner.ErrPromptAborted {
			if buffer.Len() > 0 {
				buffer.Reset()
				fmt.Println("^C")
				continue
			}
			fmt.Println(successColor("Exiting REPL. See you later!"))
			break
		} else if err == io.EOF {
			// Stdin closed (Ctrl+D, or a piped script ran out of input).
			// Without this case, err falls into the branch below on every
			// subsequent loop iteration (Prompt keeps returning io.EOF once
			// stdin is closed), spinning forever printing the same error.
			fmt.Println()
			fmt.Println(successColor("Exiting REPL. See you later!"))
			break
		} else if err != nil {
			fmt.Println(errorColor("Error reading input:"), err)
			continue
		}

		input = strings.TrimSpace(input)
		line.AppendHistory(input)

		if command, ok := commands[input]; ok {
			if command() {
				break
			}
			continue
		}

		buffer.WriteString(input)
		buffer.WriteString("\n")

		if !isIncomplete(buffer.String()) {
			command := buffer.String()
			buffer.Reset()
			commandHistory = append(commandHistory, command)
			executeCode(env, command, outputFlag)
		}
	}
}

func loadHistory(line *liner.State) {
	if f, err := os.Open(historyFilePath); err == nil {
		line.ReadHistory(f)
		f.Close()
	}
}

func saveHistory(line *liner.State) {
	if f, err := os.Create(historyFilePath); err == nil {
		line.WriteHistory(f)
		f.Close()
	}
}

const historyFilePath = ".r2lang_history"

// createR2Environment mirrors the registration list in pkg/r2lang/r2lang.go's
// RunCode. Previously this only called r2libs.RegisterLib (which registers
// just the "r2"/"go" concurrency builtins), so the REPL was missing nearly
// the entire standard library documented as globally available — print,
// Math, JSON, string/file/HTTP helpers, etc. all raised "Undeclared
// variable" errors that a script run via `go run main.go script.r2` would
// not.
func createR2Environment() *r2core.Environment {
	env := r2core.NewEnvironment()
	env.Set("false", false)
	env.Set("nil", nil)
	env.Set("null", nil)
	r2libs.RegisterLib(env)
	r2libs.RegisterStd(env)
	r2libs.RegisterIO(env)
	r2libs.RegisterHTTPClient(env)
	r2libs.RegisterRequests(env)
	r2libs.RegisterString(env)
	r2libs.RegisterRegex(env)
	r2libs.RegisterMath(env)
	r2libs.RegisterRand(env)
	r2libs.RegisterTest(env)
	r2libs.RegisterHTTP(env)
	r2libs.RegisterPrint(env)
	r2libs.RegisterOS(env)
	r2libs.RegisterHack(env)
	r2libs.RegisterEncoding(env)
	r2libs.RegisterConcurrency(env)
	r2libs.RegisterSync(env)
	r2libs.RegisterCollections(env)
	r2libs.RegisterValidate(env)
	r2libs.RegisterUnicode(env)
	r2libs.RegisterDate(env)
	r2libs.RegisterDB(env)
	r2libs.RegisterSOAP(env)
	r2libs.RegisterGRPC(env)
	r2libs.RegisterJSON(env)
	r2libs.RegisterXML(env)
	r2libs.RegisterCSV(env)
	r2libs.RegisterJWT(env)
	r2libs.RegisterConsole(env)
	r2libs.RegisterWeb(env)
	return env
}

// getCommands returns the REPL's special (dot-prefixed) commands. Each
// command reports whether the REPL should exit, so that ".exit" can break
// out of the main loop instead of calling os.Exit directly — os.Exit skips
// every deferred cleanup (restoring the terminal's raw mode, saving
// history), which otherwise left the terminal in a broken state after
// exiting the REPL this way.
func getCommands(commandHistory *[]string) map[string]func() bool {
	return map[string]func() bool{
		".exit": func() bool {
			fmt.Println(successColor("Exiting REPL. See you later!"))
			return true
		},
		".showcode": func() bool {
			if len(*commandHistory) == 0 {
				fmt.Println(warningColor("There is no code loaded."))
			} else {
				fmt.Println(infoColor("Code:"))
				for _, cmd := range *commandHistory {
					fmt.Println(highlightSyntax(cmd))
				}
			}
			return false
		},
	}
}

func executeCode(env *r2core.Environment, command string, outputFlag bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(errorColor("Error:"), r)
		}
	}()

	parser := r2core.NewParser(command)
	out := env.Run(parser)
	if outputFlag && out != nil {
		fmt.Println(successColor(out))
	}
}

// unterminatedLiteralMarkers are lexer panic messages that mean "ran out of
// input while still inside a string/template/interpolation", as opposed to a
// genuine syntax error. The REPL treats these as "need another line" too,
// the same way an unclosed brace is handled.
var unterminatedLiteralMarkers = []string{
	"Closing backtick of template string expected",
	"Unclosed interpolation in template string",
	"String sin cerrar: falta comilla de cierre",
}

// isIncomplete reports whether input still needs more lines before it can be
// parsed as a full statement. It tokenizes with the real r2core.Lexer (the
// same lexer the parser uses) and tracks bracket/paren/brace nesting depth
// from the resulting tokens, so brackets/quotes inside strings or comments
// are naturally ignored instead of being miscounted by a raw substring scan.
func isIncomplete(input string) (incomplete bool) {
	defer func() {
		if r := recover(); r != nil {
			if msg, ok := r.(string); ok {
				for _, marker := range unterminatedLiteralMarkers {
					if strings.Contains(msg, marker) {
						incomplete = true
						return
					}
				}
			}
			// Any other lexer error is a genuine syntax error, not a
			// "need more input" signal: let it through so the real
			// error is parsed and reported immediately.
		}
	}()

	lexer := r2core.NewLexer(input)
	depth := 0
	for {
		tok := lexer.NextToken()
		if tok.Type == r2core.TOKEN_EOF {
			break
		}
		if tok.Type == r2core.TOKEN_SYMBOL {
			switch tok.Value {
			case "{", "(", "[":
				depth++
			case "}", ")", "]":
				depth--
			}
		}
	}
	return depth > 0
}

// highlightSyntax aplica colores a palabras reservadas y cadenas de texto en el código
func highlightSyntax(code string) string {
	// Primero, resaltar cadenas de texto
	code = stringRegex.ReplaceAllStringFunc(code, func(str string) string {
		return stringColor(str)
	})

	// Luego, resaltar palabras reservadas
	code = keywordRegex.ReplaceAllStringFunc(code, func(str string) string {
		return keywordColor(str)
	})

	return code
}
