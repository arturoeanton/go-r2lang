package r2lang

import (
	"fmt"
)

func RegisterLib(env *Environment) {
	builtins := map[string]BuiltinFunction{
		"print": func(args ...interface{}) interface{} {
			for _, a := range args {
				fmt.Print(a, " ")
			}
			fmt.Println()
			return nil
		},
		
		"go": func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("go necesita al menos una función como argumento")
			}
			// Verificar que el primer argumento sea una función
			fn, ok := args[0].(*UserFunction)
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
}
