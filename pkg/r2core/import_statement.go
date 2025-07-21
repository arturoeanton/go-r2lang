package r2core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ImportStatement representa una declaración de importación con alias.
type ImportStatement struct {
	Path  string
	Alias string // Alias opcional
}

func (is *ImportStatement) Eval(env *Environment) interface{} {
	filePath := is.Path

	// Resolver rutas relativas
	if !filepath.IsAbs(filePath) {
		dir := env.Dir
		filePath = filepath.Join(dir, filePath)
	}

	// Clean the path to get canonical form
	filePath = filepath.Clean(filePath)

	// Check for cyclic imports BEFORE checking if already imported
	if env.IsImportCycle(filePath) {
		chain := env.GetImportChain()
		chain = append(chain, filePath)
		panic(fmt.Sprintf("Cyclic import detected: %s", strings.Join(chain, " -> ")))
	}

	// Verificar si ya fue importado
	if env.imported[filePath] {
		return nil // Ya importado, no hacer nada
	}

	// Add to import stack for cycle detection
	env.PushImport(filePath)
	defer env.PopImport()

	// Marcar como importado
	env.imported[filePath] = true

	// Leer el contenido del archivo
	content, err := os.ReadFile(filePath)
	if err != nil {
		// Create position-aware error message
		currentFile := env.CurrentFile
		if currentFile != "" {
			panic(fmt.Sprintf("%s: Error reading imported file: %s", currentFile, err.Error()))
		} else {
			panic(fmt.Sprintf("Error reading imported file: %s", err.Error()))
		}
	}

	// Crear un nuevo parser con el directorio base actualizado
	parser := NewParserWithFile(string(content), filePath)
	parser.SetBaseDir(filepath.Dir(filePath))

	// Parsear el programa importado
	importedProgram := parser.ParseProgram()

	// Crear un nuevo entorno para el módulo importado
	moduleEnv := NewInnerEnv(env)
	moduleEnv.Set("currentFile", filePath)
	moduleEnv.Dir = filepath.Dir(filePath)
	moduleEnv.CurrentFile = filePath // Set for position-aware errors

	// Evaluar en el entorno del módulo
	importedProgram.Eval(moduleEnv)

	// Obtener los símbolos del módulo importado
	symbols := moduleEnv.store

	// Si hay un alias, asignar los símbolos bajo ese alias
	if is.Alias != "" {
		env.Set(is.Alias, symbols)
	} else {
		// Si no hay alias, exportar directamente
		for k, v := range symbols {
			env.Set(k, v)
		}
	}

	return nil
}
