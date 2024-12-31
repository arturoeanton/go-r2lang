package main

import (
	"fmt"
	"github.com/arturoeanton/go-r2lang/r2lang"
	"os"
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
		data, err := os.ReadFile("main.r2")
		if err != nil {
			fmt.Println("Error al leer main.r2:", err)
			os.Exit(1)
		}
		code = string(data)
	}

	// Creamos entorno
	env := r2lang.NewEnvironment()
	// Registrar otras librerías si las tienes:
	r2lang.RegisterLib(env)
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

	parser := r2lang.NewParser(code)
	env.Run(parser)

}
