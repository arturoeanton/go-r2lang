package r2libs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// r2io.go: Manejo de archivos e I/O en tu lenguaje R2

func RegisterIO(env *r2core.Environment) {
	// 1) readFile(path) => string (contenido del archivo)
	env.Set("readFile", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("readFile necesita (path)")
		}
		path, ok := args[0].(string)
		if !ok {
			panic("readFile: primer argumento debe ser string (path)")
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			panic(fmt.Sprintf("readFile: error al leer '%s': %v", path, err))
		}
		return string(data)
	}))

	// 2) writeFile(path, contents) => nil
	//    Escribe el contenido en un archivo (sobrescribiendo)
	env.Set("writeFile", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("writeFile necesita 2 argumentos: (path, contents)")
		}
		path, ok1 := args[0].(string)
		contents, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			panic("writeFile: (path, contents) deben ser strings")
		}
		err := ioutil.WriteFile(path, []byte(contents), 0644)
		if err != nil {
			panic(fmt.Sprintf("writeFile: error al escribir '%s': %v", path, err))
		}
		return nil
	}))

	// 3) appendFile(path, contents) => nil
	//    Abre en modo append y agrega el contenido.
	env.Set("appendFile", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
	}))

	// 4) removeFile(path) => nil
	//    Borra un archivo (no directorio)
	env.Set("removeFile", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("removeFile necesita (path)")
		}
		path, ok := args[0].(string)
		if !ok {
			panic("removeFile: path debe ser string")
		}
		err := os.Remove(path)
		if err != nil {
			panic(fmt.Sprintf("removeFile: error al borrar '%s': %v", path, err))
		}
		return nil
	}))

	// 5) renameFile(oldPath, newPath) => nil
	env.Set("renameFile", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
	}))

	// 6) readDir(path) => array de strings con los nombres de entradas (archivos, directorios)
	env.Set("readDir", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("readDir necesita (path)")
		}
		dir, ok := args[0].(string)
		if !ok {
			panic("readDir: path debe ser string")
		}
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			panic(fmt.Sprintf("readDir: error al leer directorio '%s': %v", dir, err))
		}
		var result []interface{}
		for _, f := range files {
			// Ej: "archivo.txt" o "subcarpeta"
			result = append(result, f.Name())
		}
		return result
	}))

	// 7) makeDir(path) => nil
	//    Crea un directorio (equivalente a mkdir)
	env.Set("makeDir", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("makeDir necesita (path)")
		}
		dir, ok := args[0].(string)
		if !ok {
			panic("makeDir: path debe ser string")
		}
		err := os.Mkdir(dir, 0755)
		if err != nil {
			panic(fmt.Sprintf("makeDir: error al crear directorio '%s': %v", dir, err))
		}
		return nil
	}))

	// 8) makeDirs(path) => nil
	//    Crea recursivamente los directorios necesarios (mkdir -p)
	env.Set("makeDirs", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("makeDirs necesita (path)")
		}
		dir, ok := args[0].(string)
		if !ok {
			panic("makeDirs: path debe ser string")
		}
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			panic(fmt.Sprintf("makeDirs: error al crear directorios '%s': %v", dir, err))
		}
		return nil
	}))

	// 9) absPath(path) => string (obtiene ruta absoluta)
	env.Set("absPath", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
	}))

	// 10) fileExists(path) => bool
	env.Set("fileExists", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("fileExists necesita (path)")
		}
		p, ok := args[0].(string)
		if !ok {
			panic("fileExists: path debe ser string")
		}
		_, err := os.Stat(p)
		if os.IsNotExist(err) {
			return false
		}
		return true
	}))
}
