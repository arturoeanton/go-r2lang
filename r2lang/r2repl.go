// r2repl.go
package r2lang

import (
	"fmt"
	"os"
	"regexp"
	"strings"

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

	// Cargar historial si existe
	historyFile := ".r2lang_history"
	if f, err := os.Open(historyFile); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	// Registrar comandos de historial al salir
	defer func() {
		if f, err := os.Create(historyFile); err == nil {
			line.WriteHistory(f)
			f.Close()
		}
	}()

	fmt.Println(infoColor("Welcome to R2Lang REPL!"))

	env := NewEnvironment()
	env.Set("false", false)
	env.Set("nil", nil)
	env.Set("null", nil)

	// Registrar otras librerías si las tienes:
	RegisterLib(env)
	RegisterStd(env)
	RegisterIO(env)
	RegisterHTTPClient(env)
	RegisterString(env)
	RegisterMath(env)
	RegisterRand(env)
	RegisterTest(env)
	RegisterHTTP(env)
	RegisterPrint(env)
	RegisterOS(env)
	RegisterHack(env)
	RegisterConcurrency(env)
	RegisterCollections(env)

	var buffer strings.Builder
	var commandHistory []string
	for {
		// Determinar el prompt basado en si estamos en una entrada multilínea
		prompt := "> "
		if buffer.Len() > 0 {
			prompt = "... "
		}

		// Obtener la entrada del usuario sin colores en el prompt
		input, err := line.Prompt(prompt)
		if err == liner.ErrPromptAborted {
			// Ctrl+C
			if buffer.Len() > 0 {
				buffer.Reset()
				fmt.Println("^C")
				continue
			} else {
				fmt.Println(successColor("Exiting REPL. See you later!"))
				break
			}
		} else if err != nil {
			fmt.Println(errorColor("Error reading input:"), err)
			continue
		}

		input = strings.TrimSpace(input)
		line.AppendHistory(input)

		if input == ".exit" {
			fmt.Println(successColor("Exiting REPL. See you later!"))
			break
		}

		if input == ".showcode" {
			if len(commandHistory) == 0 {
				fmt.Println(warningColor("There is no code loaded."))
			} else {
				fmt.Println(infoColor("Code:"))
				for _, cmd := range commandHistory {
					fmt.Println(highlightSyntax(cmd))
				}
			}
			continue
		}

		// Manejar entrada multilínea
		buffer.WriteString(input)
		buffer.WriteString("\n")

		if !isIncomplete(buffer.String()) {
			command := buffer.String()
			buffer.Reset()

			// Añadir al historial
			commandHistory = append(commandHistory, command)

			// Ejecutar el comando
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println(errorColor("Error:"), r)
					}
				}()

				parser := NewParser(command)
				out := env.Run(parser)
				if outputFlag {
					if out != nil {
						fmt.Println(successColor(out))
					}
				}
			}()
		}
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
