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

		"r2": func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("r2 need at least one function as argument")
			}
			// Verificar que el primer argumento sea una función
			fn, ok := args[0].(*UserFunction)
			if !ok {
				panic("r2 first argument must be a function")
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
				fn.Call(args[1:]...)
			}()
			return nil
		},

		"go": func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("go need at least one function as argument")
			}
			// Verificar que el primer argumento sea una función
			fn, ok := args[0].(*UserFunction)
			if !ok {
				panic("go first argument must be a function")
			}
			// Ejecutar la función en una goroutine
			go func() {
				fn.Call(args[1:]...)
			}()
			return nil
		},
	}
	for name, fn := range builtins {
		env.Set(name, fn)
	}
}
