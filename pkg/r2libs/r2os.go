package r2libs

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// r2os.go: Funciones nativas para interactuar con el Sistema Operativo.
// Incluye manejo de procesos (exec, background process, etc.).

// Estructura para guardar referencia a un proceso lanzado en background
type R2Process struct {
	cmd    *exec.Cmd
	killed bool
}

// R2Process Eval => no hace nada, solo devolvemos algo representativo
func (rp *R2Process) Eval(env *r2core.Environment) interface{} {
	return rp // se podría devolver la misma referencia
}

func RegisterOS(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		"exit": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				os.Exit(0)
			}
			code, ok := args[0].(int)
			if !ok {
				panic("exit: arg should be int")
			}
			os.Exit(code)
			return nil
		}),

		"osName": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			// Ejemplo: en Windows a veces "Windows_NT", en Linux no siempre está OS
			// Podrías usar "runtime.GOOS" si deseas algo confiable en Go:
			// return runtime.GOOS
			val, found := os.LookupEnv("OS")
			if !found || val == "" {
				return "unknown"
			}
			return val
		}),

		"currentDir": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			dir, err := os.Getwd()
			if err != nil {
				panic("currentDir: error " + err.Error())
			}
			return dir
		}),

		"chDir": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
		}),

		"setEnv": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
		}),

		"getEnv": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
		}),

		"envList": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			vars := os.Environ()
			m := make(map[string]interface{})
			for _, kv := range vars {
				parts := strings.SplitN(kv, "=", 2)
				if len(parts) == 2 {
					m[parts[0]] = parts[1]
				}
			}
			return m
		}),

		"listDir": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
				panic(fmt.Sprintf("listDir: error reading '%s': %v", path, err))
			}
			arr := make([]interface{}, len(names))
			for i, nm := range names {
				arr[i] = nm
			}
			return arr
		}),

		"absPath": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
		}),

		"execCmd": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
		}),

		"runProcess": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
		}),

		"waitProcess": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
		}),

		"killProcess": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
		}),

		// System Information
		"getPlatform": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			return runtime.GOOS
		}),

		"getArch": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			return runtime.GOARCH
		}),

		"getNumCPU": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			return float64(runtime.NumCPU())
		}),

		"getUser": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			currentUser, err := user.Current()
			if err != nil {
				return fmt.Sprintf("Error getting current user: %v", err)
			}
			return map[string]interface{}{
				"username": currentUser.Username,
				"name":     currentUser.Name,
				"uid":      currentUser.Uid,
				"gid":      currentUser.Gid,
				"homeDir":  currentUser.HomeDir,
			}
		}),

		"getHostname": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			hostname, err := os.Hostname()
			if err != nil {
				return fmt.Sprintf("Error getting hostname: %v", err)
			}
			return hostname
		}),

		"getTempDir": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			return os.TempDir()
		}),

		"getHomeDir": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Sprintf("Error getting home directory: %v", err)
			}
			return home
		}),

		// Process Management
		"getPid": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			return float64(os.Getpid())
		}),

		"getParentPid": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			return float64(os.Getppid())
		}),

		"killPid": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("killPid needs (pid)")
			}
			pid := int(toFloat(args[0]))

			process, err := os.FindProcess(pid)
			if err != nil {
				return fmt.Sprintf("Error finding process %d: %v", pid, err)
			}

			err = process.Kill()
			if err != nil {
				return fmt.Sprintf("Error killing process %d: %v", pid, err)
			}
			return fmt.Sprintf("Process %d killed successfully", pid)
		}),

		"signalProcess": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("signalProcess needs (pid, signal)")
			}
			pid := int(toFloat(args[0]))
			sigName, ok := args[1].(string)
			if !ok {
				panic("signalProcess: signal must be string")
			}

			process, err := os.FindProcess(pid)
			if err != nil {
				return fmt.Sprintf("Error finding process %d: %v", pid, err)
			}

			var sig os.Signal
			switch sigName {
			case "KILL":
				sig = syscall.SIGKILL
			case "TERM":
				sig = syscall.SIGTERM
			case "INT":
				sig = syscall.SIGINT
			case "HUP":
				sig = syscall.SIGHUP
			case "USR1":
				sig = syscall.SIGUSR1
			case "USR2":
				sig = syscall.SIGUSR2
			default:
				return fmt.Sprintf("Unsupported signal: %s", sigName)
			}

			err = process.Signal(sig)
			if err != nil {
				return fmt.Sprintf("Error sending signal %s to process %d: %v", sigName, pid, err)
			}
			return fmt.Sprintf("Signal %s sent to process %d", sigName, pid)
		}),

		// Enhanced Process Execution
		"execWithTimeout": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("execWithTimeout needs (cmdString, timeoutSeconds)")
			}
			cmdLine, ok1 := args[0].(string)
			timeout := int(toFloat(args[1]))
			if !ok1 {
				panic("execWithTimeout: cmd should be string")
			}

			cmd := exec.Command("sh", "-c", cmdLine)

			done := make(chan error, 1)
			var output []byte
			var err error

			go func() {
				output, err = cmd.CombinedOutput()
				done <- err
			}()

			select {
			case err := <-done:
				if err != nil {
					return map[string]interface{}{
						"success": false,
						"output":  string(output),
						"error":   err.Error(),
					}
				}
				return map[string]interface{}{
					"success": true,
					"output":  string(output),
					"error":   nil,
				}
			case <-time.After(time.Duration(timeout) * time.Second):
				cmd.Process.Kill()
				return map[string]interface{}{
					"success": false,
					"output":  string(output),
					"error":   "Command timed out",
				}
			}
		}),

		"execWithEnv": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("execWithEnv needs (cmdString, envMap)")
			}
			cmdLine, ok1 := args[0].(string)
			envMap, ok2 := args[1].(map[string]interface{})
			if !ok1 || !ok2 {
				panic("execWithEnv: arguments should be (string, map)")
			}

			cmd := exec.Command("sh", "-c", cmdLine)

			// Set environment variables
			cmd.Env = os.Environ()
			for key, val := range envMap {
				if valStr, ok := val.(string); ok {
					cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, valStr))
				}
			}

			output, err := cmd.CombinedOutput()
			if err != nil {
				return map[string]interface{}{
					"success": false,
					"output":  string(output),
					"error":   err.Error(),
				}
			}
			return map[string]interface{}{
				"success": true,
				"output":  string(output),
				"error":   nil,
			}
		}),

		// Network and System State
		"getLoadAvg": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			// This is platform specific, simplified version
			if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
				return "Load average not available on this platform"
			}

			cmd := exec.Command("uptime")
			output, err := cmd.Output()
			if err != nil {
				return fmt.Sprintf("Error getting load average: %v", err)
			}
			return string(output)
		}),

		"getMemoryInfo": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			var memInfo map[string]interface{}

			switch runtime.GOOS {
			case "linux":
				cmd := exec.Command("free", "-b")
				output, err := cmd.Output()
				if err == nil {
					lines := strings.Split(string(output), "\n")
					if len(lines) > 1 {
						fields := strings.Fields(lines[1])
						if len(fields) >= 3 {
							total, _ := strconv.ParseFloat(fields[1], 64)
							used, _ := strconv.ParseFloat(fields[2], 64)
							free, _ := strconv.ParseFloat(fields[3], 64)
							memInfo = map[string]interface{}{
								"total": total,
								"used":  used,
								"free":  free,
							}
						}
					}
				}
			case "darwin":
				cmd := exec.Command("vm_stat")
				output, err := cmd.Output()
				if err == nil {
					memInfo = map[string]interface{}{
						"raw_output": string(output),
					}
				}
			default:
				memInfo = map[string]interface{}{
					"platform": runtime.GOOS,
					"message":  "Memory info not implemented for this platform",
				}
			}

			if memInfo == nil {
				memInfo = map[string]interface{}{
					"error": "Failed to get memory information",
				}
			}

			return memInfo
		}),

		// File System Information
		"getDiskUsage": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("getDiskUsage needs (path)")
			}
			path, ok := args[0].(string)
			if !ok {
				panic("getDiskUsage: path must be string")
			}

			var cmd *exec.Cmd
			switch runtime.GOOS {
			case "linux", "darwin":
				cmd = exec.Command("df", "-B1", path)
			case "windows":
				cmd = exec.Command("dir", path)
			default:
				return fmt.Sprintf("Disk usage not implemented for %s", runtime.GOOS)
			}

			output, err := cmd.Output()
			if err != nil {
				return fmt.Sprintf("Error getting disk usage: %v", err)
			}

			if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
				lines := strings.Split(string(output), "\n")
				if len(lines) > 1 {
					fields := strings.Fields(lines[1])
					if len(fields) >= 4 {
						total, _ := strconv.ParseFloat(fields[1], 64)
						used, _ := strconv.ParseFloat(fields[2], 64)
						available, _ := strconv.ParseFloat(fields[3], 64)
						return map[string]interface{}{
							"total":     total,
							"used":      used,
							"available": available,
							"path":      path,
						}
					}
				}
			}

			return map[string]interface{}{
				"raw_output": string(output),
				"path":       path,
			}
		}),

		// System Time
		"getSystemTime": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			now := time.Now()
			return map[string]interface{}{
				"unix":     float64(now.Unix()),
				"iso":      now.Format(time.RFC3339),
				"local":    now.Format("2006-01-02 15:04:05"),
				"utc":      now.UTC().Format("2006-01-02 15:04:05"),
				"timezone": now.Format("MST"),
			}
		}),

		"getUptime": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			var cmd *exec.Cmd
			switch runtime.GOOS {
			case "linux":
				cmd = exec.Command("cat", "/proc/uptime")
			case "darwin":
				cmd = exec.Command("uptime")
			case "windows":
				cmd = exec.Command("systeminfo")
			default:
				return fmt.Sprintf("Uptime not implemented for %s", runtime.GOOS)
			}

			output, err := cmd.Output()
			if err != nil {
				return fmt.Sprintf("Error getting uptime: %v", err)
			}
			return strings.TrimSpace(string(output))
		}),
	}

	RegisterModule(env, "os", functions)
}
