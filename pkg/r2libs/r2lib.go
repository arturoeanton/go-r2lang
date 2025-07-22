package r2libs

import (
	"fmt"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func RegisterLib(env *r2core.Environment) {
	builtins := map[string]r2core.BuiltinFunction{
		"r2": func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("r2 need at least one function as argument")
			}
			// Verificar que el primer argumento sea una función
			fn, ok := args[0].(*r2core.UserFunction)
			if !ok {
				panic("r2 first argument must be a function")
			}
			r2core.Add()
			// Ejecutar la función en una goroutine
			go func() {
				defer r2core.Done()
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
			fn, ok := args[0].(*r2core.UserFunction)
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
