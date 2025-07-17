package r2libs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// r2io.go: Manejo de archivos e I/O en tu lenguaje R2

func RegisterIO(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		"readFile": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("readFile necesita (path)")
			}
			path, ok := args[0].(string)
			if !ok {
				panic("readFile: primer argumento debe ser string (path)")
			}
			data, err := os.ReadFile(path)
			if err != nil {
				panic(fmt.Sprintf("readFile: error al leer '%s': %v", path, err))
			}
			return string(data)
		}),

		"writeFile": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("writeFile necesita 2 argumentos: (path, contents)")
			}
			path, ok1 := args[0].(string)
			contents, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("writeFile: (path, contents) deben ser strings")
			}
			err := os.WriteFile(path, []byte(contents), 0644)
			if err != nil {
				panic(fmt.Sprintf("writeFile: error al escribir '%s': %v", path, err))
			}
			return nil
		}),

		"appendFile": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("appendFile necesita (path, contents)")
			}
			path, ok1 := args[0].(string)
			contents, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("appendFile: (path, contents) deben ser strings")
			}

			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				panic(fmt.Sprintf("appendFile: error al abrir '%s': %v", path, err))
			}
			defer f.Close()

			_, err = f.WriteString(contents)
			if err != nil {
				panic(fmt.Sprintf("appendFile: error al escribir en '%s': %v", path, err))
			}
			return nil
		}),

		"rmFile": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("rmFile necesita (path)")
			}
			path, ok := args[0].(string)
			if !ok {
				panic("rmFile: path debe ser string")
			}
			err := os.Remove(path)
			if err != nil {
				panic(fmt.Sprintf("rmFile: error al borrar '%s': %v", path, err))
			}
			return nil
		}),

		"rmDir": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("rmDir necesita (path)")
			}
			path, ok := args[0].(string)
			if !ok {
				panic("rmDir: path debe ser string")
			}
			err := os.RemoveAll(path)
			if err != nil {
				panic(fmt.Sprintf("rmDir: error al borrar directorio '%s': %v", path, err))
			}
			return nil
		}),

		"renameFile": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("renameFile necesita (oldPath, newPath)")
			}
			oldP, ok1 := args[0].(string)
			newP, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("renameFile: (oldPath, newPath) deben ser strings")
			}
			err := os.Rename(oldP, newP)
			if err != nil {
				panic(fmt.Sprintf("renameFile: error al renombrar '%s' a '%s': %v", oldP, newP, err))
			}
			return nil
		}),

		"listdir": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("listdir necesita (path)")
			}
			dir, ok := args[0].(string)
			if !ok {
				panic("listdir: path debe ser string")
			}
			files, err := os.ReadDir(dir)
			if err != nil {
				panic(fmt.Sprintf("listdir: error al leer directorio '%s': %v", dir, err))
			}
			var result []interface{}
			for _, f := range files {
				result = append(result, f.Name())
			}
			return result
		}),

		"mkdir": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("mkdir necesita (path)")
			}
			dir, ok := args[0].(string)
			if !ok {
				panic("mkdir: path debe ser string")
			}
			err := os.Mkdir(dir, 0755)
			if err != nil {
				panic(fmt.Sprintf("mkdir: error al crear directorio '%s': %v", dir, err))
			}
			return nil
		}),

		"mkdirAll": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("mkdirAll necesita (path)")
			}
			dir, ok := args[0].(string)
			if !ok {
				panic("mkdirAll: path debe ser string")
			}
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				panic(fmt.Sprintf("mkdirAll: error al crear directorios '%s': %v", dir, err))
			}
			return nil
		}),

		"absPath": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("absPath necesita (path)")
			}
			p, ok := args[0].(string)
			if !ok {
				panic("absPath: path debe ser string")
			}
			abs, err := filepath.Abs(p)
			if err != nil {
				panic(fmt.Sprintf("absPath: error con '%s': %v", p, err))
			}
			return abs
		}),

		"exists": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("exists necesita (path)")
			}
			p, ok := args[0].(string)
			if !ok {
				panic("exists: path debe ser string")
			}
			_, err := os.Stat(p)
			return !os.IsNotExist(err)
		}),

		"isdir": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("isdir necesita (path)")
			}
			p, ok := args[0].(string)
			if !ok {
				panic("isdir: path debe ser string")
			}
			info, err := os.Stat(p)
			if os.IsNotExist(err) {
				return false
			}
			if err != nil {
				panic(fmt.Sprintf("isdir: error al obtener información de '%s': %v", p, err))
			}
			return info.IsDir()
		}),

		"isfile": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("isfile necesita (path)")
			}
			p, ok := args[0].(string)
			if !ok {
				panic("isfile: path debe ser string")
			}
			info, err := os.Stat(p)
			if os.IsNotExist(err) {
				return false
			}
			if err != nil {
				panic(fmt.Sprintf("isfile: error al obtener información de '%s': %v", p, err))
			}
			return !info.IsDir()
		}),
	}

	RegisterModule(env, "io", functions)
}
