package r2lang

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// r2os.go: Funciones nativas para interactuar con el Sistema Operativo.
// Incluye manejo de procesos (exec, background process, etc.).

// Estructura para guardar referencia a un proceso lanzado en background
type R2Process struct {
	cmd    *exec.Cmd
	killed bool
}

// R2Process Eval => no hace nada, solo devolvemos algo representativo
func (rp *R2Process) Eval(env *Environment) interface{} {
	return rp // se podría devolver la misma referencia
}

func RegisterOS(env *Environment) {

	env.Set("exit", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			os.Exit(0)
		}
		code, ok := args[0].(int)
		if !ok {
			panic("exit: arg should be int")
		}
		os.Exit(code)
		return nil
	}))

	// osName() => Devuelve una string con GOOS o variable OS
	env.Set("osName", BuiltinFunction(func(args ...interface{}) interface{} {
		// Ejemplo: en Windows a veces "Windows_NT", en Linux no siempre está OS
		// Podrías usar "runtime.GOOS" si deseas algo confiable en Go:
		// return runtime.GOOS
		val, found := os.LookupEnv("OS")
		if !found || val == "" {
			return "unknown"
		}
		return val
	}))

	// currentDir() => string con el directorio actual
	env.Set("currentDir", BuiltinFunction(func(args ...interface{}) interface{} {
		dir, err := os.Getwd()
		if err != nil {
			panic("currentDir: error " + err.Error())
		}
		return dir
	}))

	// chDir(path) => cambia el directorio actual
	env.Set("chDir", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("chDir necesita (path)")
		}
		path, ok := args[0].(string)
		if !ok {
			panic("chDir: arg debe ser string")
		}
		err := os.Chdir(path)
		if err != nil {
			panic(fmt.Sprintf("chDir: error cambiando a '%s': %v", path, err))
		}
		return nil
	}))

	// setEnv(key, value)
	env.Set("setEnv", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("setEnv(key, value)")
		}
		k, ok1 := args[0].(string)
		v, ok2 := args[1].(string)
		if !(ok1 && ok2) {
			panic("setEnv: (string, string)")
		}
		err := os.Setenv(k, v)
		if err != nil {
			panic(fmt.Sprintf("setEnv: error => %v", err))
		}
		return nil
	}))

	// getEnv(key) => string o nil
	env.Set("getEnv", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("getEnv(key)")
		}
		k, ok := args[0].(string)
		if !ok {
			panic("getEnv: arg debe ser string")
		}
		val, found := os.LookupEnv(k)
		if !found {
			return nil
		}
		return val
	}))

	// envList() => map con todas las vars
	env.Set("envList", BuiltinFunction(func(args ...interface{}) interface{} {
		vars := os.Environ()
		m := make(map[string]interface{})
		for _, kv := range vars {
			parts := strings.SplitN(kv, "=", 2)
			if len(parts) == 2 {
				m[parts[0]] = parts[1]
			}
		}
		return m
	}))

	// listDir(path) => array con nombres
	env.Set("listDir", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("listDir(path)")
		}
		path, ok := args[0].(string)
		if !ok {
			panic("listDir: arg should be string")
		}
		f, err := os.Open(path)
		if err != nil {
			panic(fmt.Sprintf("listDir: error '%s': %v", path, err))
		}
		defer f.Close()
		names, err := f.Readdirnames(-1)
		if err != nil {
			panic(fmt.Sprintf("listDir: error reading ‘%s’: %v", path, err))
		}
		arr := make([]interface{}, len(names))
		for i, nm := range names {
			arr[i] = nm
		}
		return arr
	}))

	// absPath(path) => string abs
	env.Set("absPath", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("absPath needs (path)")
		}
		path, ok := args[0].(string)
		if !ok {
			panic("absPath: arg should be string")
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			panic(fmt.Sprintf("absPath: error => %v", err))
		}
		return abs
	}))

	// execCmd(cmdString) => string output (sin background, bloqueante)
	// Lanza un comando con "sh -c" (en sistemas tipo Unix)
	env.Set("execCmd", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("execCmd(cmdString)")
		}
		cmdLine, ok := args[0].(string)
		if !ok {
			panic("execCmd: arg should be string")
		}
		out, err := exec.Command("sh", "-c", cmdLine).CombinedOutput()
		if err != nil {
			return fmt.Sprintf("Error:%v\nOutput:\n%s", err, out)
		}
		return string(out)
	}))

	// runProcess(cmdString) => R2Process (no bloquea)
	// Lanza en background. p.ej. let p = runProcess("ping google.com")
	env.Set("runProcess", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("runProcess(cmdString)")
		}
		cmdLine, ok := args[0].(string)
		if !ok {
			panic("runProcess: arg debe ser string")
		}
		// parse: naive approach => "sh -c <cmd>"
		cmd := exec.Command("sh", "-c", cmdLine)
		// Iniciar
		err := cmd.Start()
		if err != nil {
			return fmt.Sprintf("Error al iniciar '%s': %v", cmdLine, err)
		}
		// Creamos R2Process
		rp := &R2Process{cmd: cmd}
		return rp
	}))

	// waitProcess(r2proc) => string con output (o error)
	// si el proceso usara pipes custom, habría que implementarlos.
	// Con exec.Command("sh", "-c", ...) sin capturar std, no tendremos "output".
	env.Set("waitProcess", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("waitProcess needs (r2proc)")
		}
		rp, ok := args[0].(*R2Process)
		if !ok {
			panic("waitProcess: arg is not an R2Process")
		}
		if rp.killed {
			return "error:The process was already kill()ed.."
		}
		err := rp.cmd.Wait()
		if err != nil {
			return fmt.Sprintf("error:Process ended with an error: %v", err)
		}
		return "success"
	}))

	// killProcess(r2proc) => nil
	// Manda una señal de Kill al proceso
	env.Set("killProcess", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("killProcess(r2proc)")
		}
		rp, ok := args[0].(*R2Process)
		if !ok {
			panic("killProcess: arg is not an R2Process")
		}
		if rp.killed {
			return nil // ya kill
		}
		err := rp.cmd.Process.Kill()
		if err != nil {
			return fmt.Sprintf("Error kill => %v", err)
		}
		rp.killed = true
		return nil
	}))
}
