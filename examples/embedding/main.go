// Ejemplo de embedding: un programa Go que registra sus propias funciones y
// structs con el modulo "native" de R2Lang y despues corre un script .r2
// que las usa. Correr con: go run ./examples/embedding
package main

import (
	"fmt"
	"os"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
	"github.com/arturoeanton/go-r2lang/pkg/r2libs"
)

// Greeter es un struct Go comun y corriente. Cualquier campo/metodo
// EXPORTADO queda accesible desde el script via native.setField/getField/
// callMethod.
type Greeter struct {
	Name string
}

func (g *Greeter) Hello() string {
	return "Hello, " + g.Name + "!"
}

// Add es una funcion Go comun. Los numeros que llegan desde R2Lang siempre
// son float64 (no hay distincion int/float en el lenguaje); native.callFunc
// convierte automaticamente a los tipos numericos que la funcion Go espera.
func Add(a, b int) int {
	return a + b
}

func main() {
	// 1) Registrar del lado Go, ANTES de correr cualquier script. Esto es
	// lo unico que un programa anfitrion necesita hacer para exponer su
	// propio codigo a R2Lang.
	r2libs.RegisterNativeFunc("add", Add)
	r2libs.RegisterNativeStruct("Greeter", func() interface{} { return &Greeter{} })

	// 2) Armar el entorno de ejecucion. RegisterGoInterOp registra el
	// modulo "native"; los demas Register* son el resto de la stdlib que el
	// script use (aca solo necesitamos std.print).
	env := r2core.NewEnvironment()
	env.Set("true", true)
	env.Set("false", false)
	env.Set("nil", nil)
	r2libs.RegisterStd(env)
	r2libs.RegisterGoInterOp(env)

	// 3) Correr el script. Podria venir de un archivo (r2core.NewParserWithFile)
	// en vez de un string embebido.
	code := `
std.print("=== Embedding R2Lang: native.* ===")

let sum = native.callFunc("add", 3, 4)
std.print("native.callFunc(\"add\", 3, 4) =", sum)

let g = native.new("Greeter")
native.setField(g, "Name", "R2Lang")
std.print(native.callMethod(g, "Hello"))
`
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, "script error:", r)
			os.Exit(1)
		}
	}()

	parser := r2core.NewParser(code)
	env.Run(parser)
}
