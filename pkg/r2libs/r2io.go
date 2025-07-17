package r2libs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func RegisterIO(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		"readFile": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("readFile needs (path)")
			}
			path, ok := args[0].(string)
			if !ok {
				panic("readFile: first argument must be string (path)")
			}
			data, err := os.ReadFile(path)
			if err != nil {
				panic(fmt.Sprintf("readFile: error reading '%s': %v", path, err))
			}
			return string(data)
		}),

		"readFileBytes": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("readFileBytes needs (path)")
			}
			path, ok := args[0].(string)
			if !ok {
				panic("readFileBytes: first argument must be string (path)")
			}
			data, err := os.ReadFile(path)
			if err != nil {
				panic(fmt.Sprintf("readFileBytes: error reading '%s': %v", path, err))
			}
			result := make([]interface{}, len(data))
			for i, b := range data {
				result[i] = float64(b)
			}
			return result
		}),

		"readLines": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("readLines needs (path)")
			}
			path, ok := args[0].(string)
			if !ok {
				panic("readLines: first argument must be string (path)")
			}
			file, err := os.Open(path)
			if err != nil {
				panic(fmt.Sprintf("readLines: error opening '%s': %v", path, err))
			}
			defer file.Close()

			var lines []interface{}
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				panic(fmt.Sprintf("readLines: error scanning '%s': %v", path, err))
			}
			return lines
		}),

		"writeFile": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("writeFile needs 2 arguments: (path, contents)")
			}
			path, ok1 := args[0].(string)
			contents, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("writeFile: (path, contents) must be strings")
			}
			permissions := os.FileMode(0644)
			if len(args) > 2 {
				if perm, ok := args[2].(float64); ok {
					permissions = os.FileMode(perm)
				}
			}
			err := os.WriteFile(path, []byte(contents), permissions)
			if err != nil {
				panic(fmt.Sprintf("writeFile: error writing '%s': %v", path, err))
			}
			return nil
		}),

		"writeFileBytes": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("writeFileBytes needs 2 arguments: (path, bytes)")
			}
			path, ok1 := args[0].(string)
			if !ok1 {
				panic("writeFileBytes: first argument must be string (path)")
			}
			bytesArray, ok2 := args[1].([]interface{})
			if !ok2 {
				panic("writeFileBytes: second argument must be array")
			}

			data := make([]byte, len(bytesArray))
			for i, b := range bytesArray {
				if num, ok := b.(float64); ok {
					data[i] = byte(num)
				} else {
					panic("writeFileBytes: array must contain numbers")
				}
			}

			permissions := os.FileMode(0644)
			if len(args) > 2 {
				if perm, ok := args[2].(float64); ok {
					permissions = os.FileMode(perm)
				}
			}

			err := os.WriteFile(path, data, permissions)
			if err != nil {
				panic(fmt.Sprintf("writeFileBytes: error writing '%s': %v", path, err))
			}
			return nil
		}),

		"writeLines": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("writeLines needs 2 arguments: (path, lines)")
			}
			path, ok1 := args[0].(string)
			if !ok1 {
				panic("writeLines: first argument must be string (path)")
			}
			linesArray, ok2 := args[1].([]interface{})
			if !ok2 {
				panic("writeLines: second argument must be array")
			}

			var lines []string
			for _, line := range linesArray {
				if str, ok := line.(string); ok {
					lines = append(lines, str)
				} else {
					lines = append(lines, fmt.Sprintf("%v", line))
				}
			}

			content := strings.Join(lines, "\n")
			permissions := os.FileMode(0644)
			if len(args) > 2 {
				if perm, ok := args[2].(float64); ok {
					permissions = os.FileMode(perm)
				}
			}

			err := os.WriteFile(path, []byte(content), permissions)
			if err != nil {
				panic(fmt.Sprintf("writeLines: error writing '%s': %v", path, err))
			}
			return nil
		}),

		"appendFile": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("appendFile needs (path, contents)")
			}
			path, ok1 := args[0].(string)
			contents, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("appendFile: (path, contents) must be strings")
			}

			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				panic(fmt.Sprintf("appendFile: error opening '%s': %v", path, err))
			}
			defer f.Close()

			_, err = f.WriteString(contents)
			if err != nil {
				panic(fmt.Sprintf("appendFile: error writing to '%s': %v", path, err))
			}
			return nil
		}),

		"copyFile": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("copyFile needs (srcPath, destPath)")
			}
			srcPath, ok1 := args[0].(string)
			destPath, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("copyFile: (srcPath, destPath) must be strings")
			}

			srcFile, err := os.Open(srcPath)
			if err != nil {
				panic(fmt.Sprintf("copyFile: error opening source '%s': %v", srcPath, err))
			}
			defer srcFile.Close()

			destFile, err := os.Create(destPath)
			if err != nil {
				panic(fmt.Sprintf("copyFile: error creating destination '%s': %v", destPath, err))
			}
			defer destFile.Close()

			_, err = io.Copy(destFile, srcFile)
			if err != nil {
				panic(fmt.Sprintf("copyFile: error copying from '%s' to '%s': %v", srcPath, destPath, err))
			}
			return nil
		}),

		"moveFile": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("moveFile needs (srcPath, destPath)")
			}
			srcPath, ok1 := args[0].(string)
			destPath, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("moveFile: (srcPath, destPath) must be strings")
			}
			err := os.Rename(srcPath, destPath)
			if err != nil {
				panic(fmt.Sprintf("moveFile: error moving '%s' to '%s': %v", srcPath, destPath, err))
			}
			return nil
		}),

		"rmFile": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("rmFile needs (path)")
			}
			path, ok := args[0].(string)
			if !ok {
				panic("rmFile: path must be string")
			}
			err := os.Remove(path)
			if err != nil {
				panic(fmt.Sprintf("rmFile: error removing '%s': %v", path, err))
			}
			return nil
		}),

		"rmDir": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("rmDir needs (path)")
			}
			path, ok := args[0].(string)
			if !ok {
				panic("rmDir: path must be string")
			}
			err := os.RemoveAll(path)
			if err != nil {
				panic(fmt.Sprintf("rmDir: error removing directory '%s': %v", path, err))
			}
			return nil
		}),

		"renameFile": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("renameFile needs (oldPath, newPath)")
			}
			oldP, ok1 := args[0].(string)
			newP, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("renameFile: (oldPath, newPath) must be strings")
			}
			err := os.Rename(oldP, newP)
			if err != nil {
				panic(fmt.Sprintf("renameFile: error renaming '%s' to '%s': %v", oldP, newP, err))
			}
			return nil
		}),

		"listDir": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("listDir needs (path)")
			}
			dir, ok := args[0].(string)
			if !ok {
				panic("listDir: path must be string")
			}
			files, err := os.ReadDir(dir)
			if err != nil {
				panic(fmt.Sprintf("listDir: error reading directory '%s': %v", dir, err))
			}
			var result []interface{}
			for _, f := range files {
				result = append(result, f.Name())
			}
			return result
		}),

		"listDirDetailed": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("listDirDetailed needs (path)")
			}
			dir, ok := args[0].(string)
			if !ok {
				panic("listDirDetailed: path must be string")
			}
			files, err := os.ReadDir(dir)
			if err != nil {
				panic(fmt.Sprintf("listDirDetailed: error reading directory '%s': %v", dir, err))
			}
			var result []interface{}
			for _, f := range files {
				info, err := f.Info()
				if err != nil {
					continue
				}
				fileInfo := map[string]interface{}{
					"name":    f.Name(),
					"size":    float64(info.Size()),
					"isDir":   info.IsDir(),
					"mode":    float64(info.Mode()),
					"modTime": info.ModTime().Format(time.RFC3339),
				}
				result = append(result, fileInfo)
			}
			return result
		}),

		"mkdir": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("mkdir needs (path)")
			}
			dir, ok := args[0].(string)
			if !ok {
				panic("mkdir: path must be string")
			}
			permissions := os.FileMode(0755)
			if len(args) > 1 {
				if perm, ok := args[1].(float64); ok {
					permissions = os.FileMode(perm)
				}
			}
			err := os.Mkdir(dir, permissions)
			if err != nil {
				panic(fmt.Sprintf("mkdir: error creating directory '%s': %v", dir, err))
			}
			return nil
		}),

		"mkdirAll": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("mkdirAll needs (path)")
			}
			dir, ok := args[0].(string)
			if !ok {
				panic("mkdirAll: path must be string")
			}
			permissions := os.FileMode(0755)
			if len(args) > 1 {
				if perm, ok := args[1].(float64); ok {
					permissions = os.FileMode(perm)
				}
			}
			err := os.MkdirAll(dir, permissions)
			if err != nil {
				panic(fmt.Sprintf("mkdirAll: error creating directories '%s': %v", dir, err))
			}
			return nil
		}),

		"absPath": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("absPath needs (path)")
			}
			p, ok := args[0].(string)
			if !ok {
				panic("absPath: path must be string")
			}
			abs, err := filepath.Abs(p)
			if err != nil {
				panic(fmt.Sprintf("absPath: error with '%s': %v", p, err))
			}
			return abs
		}),

		"baseName": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("baseName needs (path)")
			}
			p, ok := args[0].(string)
			if !ok {
				panic("baseName: path must be string")
			}
			return filepath.Base(p)
		}),

		"dirName": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("dirName needs (path)")
			}
			p, ok := args[0].(string)
			if !ok {
				panic("dirName: path must be string")
			}
			return filepath.Dir(p)
		}),

		"extName": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("extName needs (path)")
			}
			p, ok := args[0].(string)
			if !ok {
				panic("extName: path must be string")
			}
			return filepath.Ext(p)
		}),

		"joinPath": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("joinPath needs at least one path")
			}
			var paths []string
			for _, arg := range args {
				if p, ok := arg.(string); ok {
					paths = append(paths, p)
				} else {
					panic("joinPath: all arguments must be strings")
				}
			}
			return filepath.Join(paths...)
		}),

		"exists": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("exists needs (path)")
			}
			p, ok := args[0].(string)
			if !ok {
				panic("exists: path must be string")
			}
			_, err := os.Stat(p)
			return !os.IsNotExist(err)
		}),

		"isDir": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("isDir needs (path)")
			}
			p, ok := args[0].(string)
			if !ok {
				panic("isDir: path must be string")
			}
			info, err := os.Stat(p)
			if os.IsNotExist(err) {
				return false
			}
			if err != nil {
				panic(fmt.Sprintf("isDir: error getting info for '%s': %v", p, err))
			}
			return info.IsDir()
		}),

		"isFile": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("isFile needs (path)")
			}
			p, ok := args[0].(string)
			if !ok {
				panic("isFile: path must be string")
			}
			info, err := os.Stat(p)
			if os.IsNotExist(err) {
				return false
			}
			if err != nil {
				panic(fmt.Sprintf("isFile: error getting info for '%s': %v", p, err))
			}
			return !info.IsDir()
		}),

		"fileSize": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("fileSize needs (path)")
			}
			p, ok := args[0].(string)
			if !ok {
				panic("fileSize: path must be string")
			}
			info, err := os.Stat(p)
			if err != nil {
				panic(fmt.Sprintf("fileSize: error getting info for '%s': %v", p, err))
			}
			return float64(info.Size())
		}),

		"fileMode": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("fileMode needs (path)")
			}
			p, ok := args[0].(string)
			if !ok {
				panic("fileMode: path must be string")
			}
			info, err := os.Stat(p)
			if err != nil {
				panic(fmt.Sprintf("fileMode: error getting info for '%s': %v", p, err))
			}
			return float64(info.Mode())
		}),

		"fileModTime": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("fileModTime needs (path)")
			}
			p, ok := args[0].(string)
			if !ok {
				panic("fileModTime: path must be string")
			}
			info, err := os.Stat(p)
			if err != nil {
				panic(fmt.Sprintf("fileModTime: error getting info for '%s': %v", p, err))
			}
			return &r2core.DateValue{Time: info.ModTime()}
		}),

		"chmod": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("chmod needs (path, mode)")
			}
			p, ok1 := args[0].(string)
			mode, ok2 := args[1].(float64)
			if !ok1 || !ok2 {
				panic("chmod: (path, mode) must be (string, number)")
			}
			err := os.Chmod(p, os.FileMode(mode))
			if err != nil {
				panic(fmt.Sprintf("chmod: error changing mode for '%s': %v", p, err))
			}
			return nil
		}),

		"walk": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("walk needs (path)")
			}
			root, ok := args[0].(string)
			if !ok {
				panic("walk: path must be string")
			}

			var result []interface{}
			err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				fileInfo := map[string]interface{}{
					"path":    path,
					"name":    info.Name(),
					"size":    float64(info.Size()),
					"isDir":   info.IsDir(),
					"mode":    float64(info.Mode()),
					"modTime": info.ModTime().Format(time.RFC3339),
				}
				result = append(result, fileInfo)
				return nil
			})

			if err != nil {
				panic(fmt.Sprintf("walk: error walking '%s': %v", root, err))
			}
			return result
		}),

		"glob": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("glob needs (pattern)")
			}
			pattern, ok := args[0].(string)
			if !ok {
				panic("glob: pattern must be string")
			}
			matches, err := filepath.Glob(pattern)
			if err != nil {
				panic(fmt.Sprintf("glob: error with pattern '%s': %v", pattern, err))
			}
			result := make([]interface{}, len(matches))
			for i, match := range matches {
				result[i] = match
			}
			return result
		}),

		"findFiles": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("findFiles needs (root, pattern)")
			}
			root, ok1 := args[0].(string)
			pattern, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("findFiles: (root, pattern) must be strings")
			}

			var result []interface{}
			err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					matched, err := filepath.Match(pattern, info.Name())
					if err != nil {
						return err
					}
					if matched {
						result = append(result, path)
					}
				}
				return nil
			})

			if err != nil {
				panic(fmt.Sprintf("findFiles: error searching in '%s': %v", root, err))
			}
			return result
		}),

		"sortFiles": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("sortFiles needs (files)")
			}
			filesArray, ok := args[0].([]interface{})
			if !ok {
				panic("sortFiles: first argument must be array")
			}

			var files []string
			for _, f := range filesArray {
				if str, ok := f.(string); ok {
					files = append(files, str)
				} else {
					panic("sortFiles: array must contain strings")
				}
			}

			sortBy := "name"
			if len(args) > 1 {
				if s, ok := args[1].(string); ok {
					sortBy = s
				}
			}

			switch sortBy {
			case "size":
				sort.Slice(files, func(i, j int) bool {
					info1, err1 := os.Stat(files[i])
					info2, err2 := os.Stat(files[j])
					if err1 != nil || err2 != nil {
						return files[i] < files[j]
					}
					return info1.Size() < info2.Size()
				})
			case "time":
				sort.Slice(files, func(i, j int) bool {
					info1, err1 := os.Stat(files[i])
					info2, err2 := os.Stat(files[j])
					if err1 != nil || err2 != nil {
						return files[i] < files[j]
					}
					return info1.ModTime().Before(info2.ModTime())
				})
			default:
				sort.Strings(files)
			}

			result := make([]interface{}, len(files))
			for i, f := range files {
				result[i] = f
			}
			return result
		}),

		"tempDir": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			return os.TempDir()
		}),

		"tempFile": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			dir := ""
			pattern := "temp*"

			if len(args) > 0 {
				if d, ok := args[0].(string); ok {
					dir = d
				}
			}
			if len(args) > 1 {
				if p, ok := args[1].(string); ok {
					pattern = p
				}
			}

			file, err := os.CreateTemp(dir, pattern)
			if err != nil {
				panic(fmt.Sprintf("tempFile: error creating temp file: %v", err))
			}
			defer file.Close()

			return file.Name()
		}),

		"workingDir": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			wd, err := os.Getwd()
			if err != nil {
				panic(fmt.Sprintf("workingDir: error getting working directory: %v", err))
			}
			return wd
		}),

		"changeDir": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("changeDir needs (path)")
			}
			dir, ok := args[0].(string)
			if !ok {
				panic("changeDir: path must be string")
			}
			err := os.Chdir(dir)
			if err != nil {
				panic(fmt.Sprintf("changeDir: error changing directory to '%s': %v", dir, err))
			}
			return nil
		}),
	}

	RegisterModule(env, "io", functions)
}
