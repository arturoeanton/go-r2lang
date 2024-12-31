package main

import (
	"fmt"
	"github.com/arturoeanton/go-r2lang/r2lang"
	"io/ioutil"
	"os"
	"sync"
)

func main() {
	var code string
	if len(os.Args) > 1 {
		filename := os.Args[1]
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("Error al leer el archivo %s: %v\n", filename, err)
			os.Exit(1)
		}
		code = string(data)
	} else {
		// intentar main.r2
		if _, err := os.Stat("main.r2"); os.IsNotExist(err) {
			fmt.Println("Error: Debes pasar un archivo .r2 o tener main.r2 en el directorio actual.")
			os.Exit(1)
		}
		data, err := ioutil.ReadFile("main.r2")
		if err != nil {
			fmt.Println("Error al leer main.r2:", err)
			os.Exit(1)
		}
		code = string(data)
	}

	parser := r2lang.NewParser(code)
	ast := parser.ParseProgram()

	// Creamos entorno
	env := r2lang.NewEnvironment()

	wg := sync.WaitGroup{}
	// Builtins básicos (print, println, printf)
	builtins := map[string]r2lang.BuiltinFunction{
		"print": func(args ...interface{}) interface{} {
			for _, a := range args {
				fmt.Print(a, " ")
			}
			fmt.Println()
			return nil
		},
		"println": func(args ...interface{}) interface{} {
			fmt.Println(args...)
			return nil
		},
		"printf": func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("printf necesita al menos un string de formato")
			}
			formatStr, ok := args[0].(string)
			if !ok {
				panic("El primer argumento de printf debe ser un string")
			}
			fmt.Printf(formatStr, args[1:]...)
			return nil
		},
		"go": func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("go necesita al menos una función como argumento")
			}
			// Verificar que el primer argumento sea una función
			fn, ok := args[0].(*r2lang.UserFunction)
			if !ok {
				panic("El argumento de go debe ser una función")
			}
			wg.Add(1)
			// Ejecutar la función en una goroutine
			go func() {
				defer wg.Done()
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("Error en goroutine:", r)
					}
				}()
				fn.Call()
			}()
			return nil
		},
	}
	for name, fn := range builtins {
		env.Set(name, fn)
	}

	// Registrar otras librerías si las tienes:
	r2lang.RegisterStd(env)
	r2lang.RegisterIO(env)
	r2lang.RegisterHTTPClient(env)
	r2lang.RegisterString(env)
	r2lang.RegisterMath(env)
	r2lang.RegisterRand(env)
	r2lang.RegisterTest(env)
	r2lang.RegisterHTTP(env)
	r2lang.RegisterPrint(env)
	r2lang.RegisterOS(env)
	r2lang.RegisterHack(env)
	r2lang.RegisterConcurrency(env)

	// Ejecutar
	ast.Eval(env)

	// Llamar a main() si está
	mainVal, ok := env.Get("main")
	if !ok {
		fmt.Println("Aviso: No existe función main().")
		os.Exit(0)
	}
	mainFn, isFn := mainVal.(*r2lang.UserFunction)
	if !isFn {
		fmt.Println("Error: 'main' no es una función.")
		os.Exit(1)
	}
	mainFn.Call()
	wg.Wait()
}
