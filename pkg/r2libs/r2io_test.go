package r2libs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func TestIOFunctions(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterIO(env)

	ioModuleObj, ok := env.Get("io")
	if !ok {
		t.Fatal("io module not found")
	}
	ioModule := ioModuleObj.(map[string]interface{})

	readFileFunc := ioModule["readFile"].(r2core.BuiltinFunction)
	writeFileFunc := ioModule["writeFile"].(r2core.BuiltinFunction)
	appendFileFunc := ioModule["appendFile"].(r2core.BuiltinFunction)
	rmFileFunc := ioModule["rmFile"].(r2core.BuiltinFunction)
	rmDirFunc := ioModule["rmDir"].(r2core.BuiltinFunction)
	renameFileFunc := ioModule["renameFile"].(r2core.BuiltinFunction)
	listdirFunc := ioModule["listDir"].(r2core.BuiltinFunction)
	mkdirFunc := ioModule["mkdir"].(r2core.BuiltinFunction)
	mkdirAllFunc := ioModule["mkdirAll"].(r2core.BuiltinFunction)
	absPathFunc := ioModule["absPath"].(r2core.BuiltinFunction)
	existsFunc := ioModule["exists"].(r2core.BuiltinFunction)
	isdirFunc := ioModule["isDir"].(r2core.BuiltinFunction)
	isfileFunc := ioModule["isFile"].(r2core.BuiltinFunction)

	// Setup a temporary directory for tests
	testDir := filepath.Join(os.TempDir(), "r2lang_io_test")
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	// Test writeFile and readFile
	filePath := filepath.Join(testDir, "test.txt")
	writeFileFunc(filePath, "Hello, R2Lang!")
	content := readFileFunc(filePath).(string)
	if content != "Hello, R2Lang!" {
		t.Errorf("readFile: expected 'Hello, R2Lang!', got %s", content)
	}

	// Test appendFile
	appendFileFunc(filePath, "\nAppended content.")
	content = readFileFunc(filePath).(string)
	if content != "Hello, R2Lang!\nAppended content." {
		t.Errorf("appendFile: unexpected content: %s", content)
	}

	// Test exists, isfile, isdir
	if !existsFunc(filePath).(bool) {
		t.Errorf("exists: expected true for existing file")
	}
	if !isfileFunc(filePath).(bool) {
		t.Errorf("isfile: expected true for file")
	}
	if isdirFunc(filePath).(bool) {
		t.Errorf("isdir: expected false for file")
	}
	if existsFunc(filepath.Join(testDir, "nonexistent.txt")).(bool) {
		t.Errorf("exists: expected false for nonexistent file")
	}

	// Test mkdir and listdir
	subDir := filepath.Join(testDir, "subdir")
	mkdirFunc(subDir)
	if !existsFunc(subDir).(bool) {
		t.Errorf("mkdir: expected true for created directory")
	}
	if !isdirFunc(subDir).(bool) {
		t.Errorf("isdir: expected true for directory")
	}
	if isfileFunc(subDir).(bool) {
		t.Errorf("isfile: expected false for directory")
	}

	// Create another file in subdir for listdir test
	filePath2 := filepath.Join(subDir, "test2.txt")
	writeFileFunc(filePath2, "Content 2")

	files := listdirFunc(testDir).([]interface{})
	foundTestTxt := false
	foundSubdir := false
	for _, f := range files {
		if f.(string) == "test.txt" {
			foundTestTxt = true
		}
		if f.(string) == "subdir" {
			foundSubdir = true
		}
	}
	if !foundTestTxt || !foundSubdir || len(files) != 2 {
		t.Errorf("listdir: unexpected result: %v", files)
	}

	// Test mkdirAll
	deepDir := filepath.Join(testDir, "deep", "nested", "dir")
	mkdirAllFunc(deepDir)
	if !existsFunc(deepDir).(bool) {
		t.Errorf("mkdirAll: expected true for created deep directory")
	}

	// Test renameFile
	newPath := filepath.Join(testDir, "renamed.txt")
	renameFileFunc(filePath, newPath)
	if existsFunc(filePath).(bool) {
		t.Errorf("renameFile: old file should not exist")
	}
	if !existsFunc(newPath).(bool) {
		t.Errorf("renameFile: new file should exist")
	}
	content = readFileFunc(newPath).(string)
	if content != "Hello, R2Lang!\nAppended content." {
		t.Errorf("renameFile: unexpected content in renamed file: %s", content)
	}

	// Test absPath
	abs := absPathFunc("test.txt").(string)
	if !filepath.IsAbs(abs) {
		t.Errorf("absPath: expected absolute path, got %s", abs)
	}

	// Test rmFile
	rmFileFunc(newPath)
	if existsFunc(newPath).(bool) {
		t.Errorf("rmFile: expected file to be removed")
	}

	// Test rmDir
	rmDirFunc(subDir)
	if existsFunc(subDir).(bool) {
		t.Errorf("rmDir: expected directory to be removed")
	}
	// Test rmDir on non-empty directory (should remove recursively)
	rmDirFunc(filepath.Join(testDir, "deep"))
	if existsFunc(filepath.Join(testDir, "deep")).(bool) {
		t.Errorf("rmDir: expected deep directory to be removed")
	}
}
