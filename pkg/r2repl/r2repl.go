// r2repl.go
package r2repl

import (
	"fmt"
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
	setupHistory(line)

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
		} else if err != nil {
			fmt.Println(errorColor("Error reading input:"), err)
			continue
		}

		input = strings.TrimSpace(input)
		line.AppendHistory(input)

		if command, ok := commands[input]; ok {
			command()
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

func setupHistory(line *liner.State) {
	historyFile := ".r2lang_history"
	if f, err := os.Open(historyFile); err == nil {
		line.ReadHistory(f)
		f.Close()
	}
	defer func() {
		if f, err := os.Create(historyFile); err == nil {
			line.WriteHistory(f)
			f.Close()
		}
	}()
}

func createR2Environment() *r2core.Environment {
	env := r2core.NewEnvironment()
	env.Set("false", false)
	env.Set("nil", nil)
	env.Set("null", nil)
	r2libs.RegisterLib(env)
	return env
}

func getCommands(commandHistory *[]string) map[string]func() {
	return map[string]func(){
		".exit": func() {
			fmt.Println(successColor("Exiting REPL. See you later!"))
			os.Exit(0)
		},
		".showcode": func() {
			if len(*commandHistory) == 0 {
				fmt.Println(warningColor("There is no code loaded."))
			} else {
				fmt.Println(infoColor("Code:"))
				for _, cmd := range *commandHistory {
					fmt.Println(highlightSyntax(cmd))
				}
			}
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

// isIncomplete verifica si la entrada es una línea completa o si necesita más líneas (multilínea)
func isIncomplete(input string) bool {
	// Implementa una lógica más robusta según la sintaxis de tu lenguaje.
	// Aquí, simplemente verificamos si hay más llaves abiertas que cerradas.
	openBraces := strings.Count(input, "{")
	closeBraces := strings.Count(input, "}")
	return openBraces > closeBraces
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
