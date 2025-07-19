package r2libs

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// PathObject represents a file system path with a fluent API.
type NativeFunction struct {
	Fn r2core.BuiltinFunction
}

func (nf *NativeFunction) Eval(env *r2core.Environment) interface{} {
	return nf.Fn
}

// PathObject represents a file system path with a fluent API.
type PathObject struct {
	Path string
}

func (p *PathObject) Eval(env *r2core.Environment) interface{} {
	// A PathObject evaluates to itself, making it available for method calls.
	return p
}

// Getattr implements the attribute access for PathObject, returning functions that operate on the path.

// FileStreamObject represents a stream of lines from a file with fluent operations.
type FileStreamObject struct {
	path    string
	filters []r2core.BuiltinFunction
	mappers []r2core.BuiltinFunction
	limit   int
}

func (fs *FileStreamObject) Eval(env *r2core.Environment) interface{} {
	return fs
}

func (fs *FileStreamObject) Getattr(name string) (r2core.Node, bool) {
	switch name {
	case "lines":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			// This method doesn't execute, it just sets up the stream for iteration.
			// The actual reading happens in toArray or saveTo.
			return fs
		}}, true
	case "filter":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("FileStream.filter needs a callback function")
			}
			callback, ok := args[0].(r2core.BuiltinFunction)
			if !ok {
				panic("FileStream.filter: argument must be a function")
			}
			newFs := *fs // Create a copy to avoid modifying the original object
			newFs.filters = append(newFs.filters, callback)
			return &newFs
		}}, true
	case "map":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("FileStream.map needs a callback function")
			}
			callback, ok := args[0].(r2core.BuiltinFunction)
			if !ok {
				panic("FileStream.map: argument must be a function")
			}
			newFs := *fs // Create a copy
			newFs.mappers = append(newFs.mappers, callback)
			return &newFs
		}}, true
	case "limit":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("FileStream.limit needs a number")
			}
			limit, ok := args[0].(float64)
			if !ok {
				panic("FileStream.limit: argument must be a number")
			}
			newFs := *fs // Create a copy
			newFs.limit = int(limit)
			return &newFs
		}}, true
	case "toArray":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			file, err := os.Open(fs.path)
			if err != nil {
				panic(fmt.Sprintf("FileStream.toArray: error opening '%s': %v", fs.path, err))
			}
			defer file.Close()

			var result []interface{}
			scanner := bufio.NewScanner(file)
			linesRead := 0
			for scanner.Scan() {
				line := scanner.Text()

				// Apply filters
				shouldInclude := true
				for _, filter := range fs.filters {
					filterResult := filter(line).(bool) // Assuming filter returns boolean
					if !filterResult {
						shouldInclude = false
						break
					}
				}

				if shouldInclude {
					// Apply mappers
					mappedLine := interface{}(line)
					for _, mapper := range fs.mappers {
						mappedLine = mapper(mappedLine)
					}
					result = append(result, mappedLine)
					linesRead++
				}

				if fs.limit > 0 && linesRead >= fs.limit {
					break
				}
			}
			if err := scanner.Err(); err != nil {
				panic(fmt.Sprintf("FileStream.toArray: error scanning '%s': %v", fs.path, err))
			}
			return result
		}}, true
	case "saveTo":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("FileStream.saveTo needs a destination path")
			}
			destPath, ok := args[0].(string)
			if !ok {
				panic("FileStream.saveTo: argument must be a string")
			}

			file, err := os.Open(fs.path)
			if err != nil {
				panic(fmt.Sprintf("FileStream.saveTo: error opening source '%s': %v", fs.path, err))
			}
			defer file.Close()

			destFile, err := os.Create(destPath)
			if err != nil {
				panic(fmt.Sprintf("FileStream.saveTo: error creating destination '%s': %v", destPath, err))
			}
			defer destFile.Close()

			writer := bufio.NewWriter(destFile)
			scanner := bufio.NewScanner(file)
			linesWritten := 0
			for scanner.Scan() {
				line := scanner.Text()

				// Apply filters
				shouldInclude := true
				for _, filter := range fs.filters {
					filterResult := filter(line).(bool) // Assuming filter returns boolean
					if !filterResult {
						shouldInclude = false
						break
					}
				}

				if shouldInclude {
					// Apply mappers
					mappedLine := interface{}(line)
					for _, mapper := range fs.mappers {
						mappedLine = mapper(mappedLine)
					}
					_, err := writer.WriteString(fmt.Sprintf("%v\n", mappedLine))
					if err != nil {
						panic(fmt.Sprintf("FileStream.saveTo: error writing to '%s': %v", destPath, err))
					}
					linesWritten++
				}

				if fs.limit > 0 && linesWritten >= fs.limit {
					break
				}
			}
			writer.Flush()
			if err := scanner.Err(); err != nil {
				panic(fmt.Sprintf("FileStream.saveTo: error scanning '%s': %v", fs.path, err))
			}
			return nil
		}}, true
	}
	return nil, false
}

func (p *PathObject) Getattr(name string) (r2core.Node, bool) {

	switch name {
	case "readText":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			data, err := os.ReadFile(p.Path)
			if err != nil {
				panic(fmt.Sprintf("Path.readText: error reading '%s': %v", p.Path, err))
			}
			return string(data)
		}}, true
	case "writeText":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("Path.writeText needs (contents)")
			}
			contents, ok := args[0].(string)
			if !ok {
				panic("Path.writeText: contents must be a string")
			}
			err := os.WriteFile(p.Path, []byte(contents), 0644)
			if err != nil {
				panic(fmt.Sprintf("Path.writeText: error writing '%s': %v", p.Path, err))
			}
			return p
		}}, true
	case "exists":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			_, err := os.Stat(p.Path)
			return !os.IsNotExist(err)
		}}, true
	case "isDir":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			info, err := os.Stat(p.Path)
			if err != nil {
				return false
			}
			return info.IsDir()
		}}, true
	case "isFile":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			info, err := os.Stat(p.Path)
			if err != nil {
				return false
			}
			return !info.IsDir()
		}}, true
	case "name":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			return filepath.Base(p.Path)
		}}, true
	case "ext":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			return filepath.Ext(p.Path)
		}}, true
	case "stem":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			base := filepath.Base(p.Path)
			ext := filepath.Ext(base)
			return strings.TrimSuffix(base, ext)
		}}, true
	case "dir":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			return &PathObject{Path: filepath.Dir(p.Path)}
		}}, true
	case "remove":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			err := os.Remove(p.Path)
			if err != nil {
				panic(fmt.Sprintf("Path.remove: error removing '%s': %v", p.Path, err))
			}
			return nil
		}}, true
	case "copyTo":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("Path.copyTo needs (destination)")
			}
			var destPath string
			switch dest := args[0].(type) {
			case string:
				destPath = dest
			case *PathObject:
				destPath = dest.Path
			default:
				panic("Path.copyTo: destination must be a string or a Path object")
			}
			err := copyFileInternal(p.Path, destPath)
			if err != nil {
				panic(fmt.Sprintf("Path.copyTo: error copying to '%s': %v", destPath, err))
			}
			return &PathObject{Path: destPath}
		}}, true
	case "moveTo":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("Path.moveTo needs (destination)")
			}
			var destPath string
			switch dest := args[0].(type) {
			case string:
				destPath = dest
			case *PathObject:
				destPath = dest.Path
			default:
				panic("Path.moveTo: destination must be a string or a Path object")
			}
			err := os.Rename(p.Path, destPath)
			if err != nil {
				panic(fmt.Sprintf("Path.moveTo: error moving to '%s': %v", destPath, err))
			}
			p.Path = destPath // Update the path of the object itself
			return p
		}}, true
	case "sibling":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("Path.sibling needs (name)")
			}
			name, ok := args[0].(string)
			if !ok {
				panic("Path.sibling: name must be a string")
			}
			return &PathObject{Path: filepath.Join(filepath.Dir(p.Path), name)}
		}}, true
	case "withSuffix":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("Path.withSuffix needs (suffix)")
			}
			suffix, ok := args[0].(string)
			if !ok {
				panic("Path.withSuffix: suffix must be a string")
			}
			base := strings.TrimSuffix(p.Path, filepath.Ext(p.Path))
			return &PathObject{Path: base + suffix}
		}}, true
	case "iter":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			info, err := os.Stat(p.Path)
			if err != nil {
				panic(fmt.Sprintf("Path.iter: cannot read path '%s': %v", p.Path, err))
			}
			if !info.IsDir() {
				panic("Path.iter can only be called on a directory")
			}
			files, err := os.ReadDir(p.Path)
			if err != nil {
				panic(fmt.Sprintf("Path.iter: error reading directory '%s': %v", p.Path, err))
			}
			var result []interface{}
			for _, f := range files {
				result = append(result, &PathObject{Path: filepath.Join(p.Path, f.Name())})
			}
			return result
		}}, true
	}
	return nil, false
}

func RegisterIO(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		"Path": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("io.Path needs a path string")
			}
			path, ok := args[0].(string)
			if !ok {
				panic("io.Path: argument must be a string")
			}
			return &PathObject{Path: path}
		}),
		"FileStream": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("io.FileStream needs a path string")
			}
			path, ok := args[0].(string)
			if !ok {
				panic("io.FileStream: argument must be a string")
			}
			return &FileStreamObject{path: path}
		}),

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

		// Stream operations
		"readStream": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("readStream needs (path)")
			}
			path, ok := args[0].(string)
			if !ok {
				panic("readStream: path must be string")
			}

			batchSize := 1024
			if len(args) > 1 {
				if size, ok := args[1].(float64); ok {
					batchSize = int(size)
				}
			}

			file, err := os.Open(path)
			if err != nil {
				panic(fmt.Sprintf("readStream: error opening '%s': %v", path, err))
			}
			defer file.Close()

			var chunks []interface{}
			buffer := make([]byte, batchSize)
			for {
				n, err := file.Read(buffer)
				if n > 0 {
					chunks = append(chunks, string(buffer[:n]))
				}
				if err == io.EOF {
					break
				}
				if err != nil {
					panic(fmt.Sprintf("readStream: error reading '%s': %v", path, err))
				}
			}
			return chunks
		}),

		"writeStream": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("writeStream needs (path, chunks)")
			}
			path, ok1 := args[0].(string)
			chunks, ok2 := args[1].([]interface{})
			if !ok1 || !ok2 {
				panic("writeStream: (path, chunks) must be (string, array)")
			}

			file, err := os.Create(path)
			if err != nil {
				panic(fmt.Sprintf("writeStream: error creating '%s': %v", path, err))
			}
			defer file.Close()

			for _, chunk := range chunks {
				if str, ok := chunk.(string); ok {
					_, err := file.WriteString(str)
					if err != nil {
						panic(fmt.Sprintf("writeStream: error writing to '%s': %v", path, err))
					}
				} else {
					panic("writeStream: all chunks must be strings")
				}
			}
			return nil
		}),

		// File comparison
		"compareFiles": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("compareFiles needs (path1, path2)")
			}
			path1, ok1 := args[0].(string)
			path2, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("compareFiles: paths must be strings")
			}

			data1, err := os.ReadFile(path1)
			if err != nil {
				panic(fmt.Sprintf("compareFiles: error reading '%s': %v", path1, err))
			}
			data2, err := os.ReadFile(path2)
			if err != nil {
				panic(fmt.Sprintf("compareFiles: error reading '%s': %v", path2, err))
			}

			if len(data1) != len(data2) {
				return false
			}

			for i := range data1 {
				if data1[i] != data2[i] {
					return false
				}
			}
			return true
		}),

		// File checksum
		"checksum": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("checksum needs (path)")
			}
			path, ok := args[0].(string)
			if !ok {
				panic("checksum: path must be string")
			}

			algorithm := "sha256"
			if len(args) > 1 {
				if algo, ok := args[1].(string); ok {
					algorithm = algo
				}
			}

			data, err := os.ReadFile(path)
			if err != nil {
				panic(fmt.Sprintf("checksum: error reading '%s': %v", path, err))
			}

			switch algorithm {
			case "md5":
				sum := fmt.Sprintf("%x", md5.Sum(data))
				return sum
			case "sha1":
				sum := fmt.Sprintf("%x", sha1.Sum(data))
				return sum
			case "sha256":
				sum := fmt.Sprintf("%x", sha256.Sum256(data))
				return sum
			default:
				panic(fmt.Sprintf("checksum: unsupported algorithm '%s'", algorithm))
			}
		}),

		// Directory operations
		"createPath": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("createPath needs (path)")
			}
			path, ok := args[0].(string)
			if !ok {
				panic("createPath: path must be string")
			}

			err := os.MkdirAll(filepath.Dir(path), 0755)
			if err != nil {
				panic(fmt.Sprintf("createPath: error creating directories for '%s': %v", path, err))
			}
			return nil
		}),

		// Backup operations
		"backup": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("backup needs (path)")
			}
			path, ok := args[0].(string)
			if !ok {
				panic("backup: path must be string")
			}

			timestamp := time.Now().Format("20060102_150405")
			ext := filepath.Ext(path)
			base := strings.TrimSuffix(path, ext)
			backupPath := fmt.Sprintf("%s_backup_%s%s", base, timestamp, ext)

			err := copyFileInternal(path, backupPath)
			if err != nil {
				panic(fmt.Sprintf("backup: error creating backup: %v", err))
			}
			return backupPath
		}),

		// File watching simulation
		"watchFile": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("watchFile needs (path)")
			}
			path, ok := args[0].(string)
			if !ok {
				panic("watchFile: path must be string")
			}

			info, err := os.Stat(path)
			if err != nil {
				panic(fmt.Sprintf("watchFile: error getting info for '%s': %v", path, err))
			}

			return map[string]interface{}{
				"path":    path,
				"size":    float64(info.Size()),
				"modTime": info.ModTime().Format(time.RFC3339),
				"mode":    float64(info.Mode()),
			}
		}),

		// Batch operations
		"batchCopy": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("batchCopy needs (srcPattern, destDir)")
			}
			pattern, ok1 := args[0].(string)
			destDir, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("batchCopy: arguments must be strings")
			}

			matches, err := filepath.Glob(pattern)
			if err != nil {
				panic(fmt.Sprintf("batchCopy: error with pattern '%s': %v", pattern, err))
			}

			var results []interface{}
			for _, srcPath := range matches {
				filename := filepath.Base(srcPath)
				destPath := filepath.Join(destDir, filename)

				err := copyFileInternal(srcPath, destPath)
				if err == nil {
					results = append(results, map[string]interface{}{
						"src":    srcPath,
						"dest":   destPath,
						"status": "success",
					})
				} else {
					results = append(results, map[string]interface{}{
						"src":    srcPath,
						"dest":   destPath,
						"status": "error",
						"error":  err.Error(),
					})
				}
			}
			return results
		}),

		// File metadata
		"getMetadata": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("getMetadata needs (path)")
			}
			path, ok := args[0].(string)
			if !ok {
				panic("getMetadata: path must be string")
			}

			info, err := os.Stat(path)
			if err != nil {
				panic(fmt.Sprintf("getMetadata: error getting info for '%s': %v", path, err))
			}

			return map[string]interface{}{
				"name":    info.Name(),
				"size":    float64(info.Size()),
				"mode":    float64(info.Mode()),
				"modTime": info.ModTime().Format(time.RFC3339),
				"isDir":   info.IsDir(),
				"abs":     func() string { abs, _ := filepath.Abs(path); return abs }(),
				"ext":     filepath.Ext(path),
				"dir":     filepath.Dir(path),
				"base":    filepath.Base(path),
			}
		}),
	}

	RegisterModule(env, "io", functions)
}

// Helper function for internal file copying
func copyFileInternal(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
