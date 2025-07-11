package r2core

import (
	"os"
	"path/filepath"
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

	// Verificar si ya fue importado
	if env.imported[filePath] {
		return nil // Ya importado, no hacer nada
	}

	// Marcar como importado
	env.imported[filePath] = true

	// Leer el contenido del archivo
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic("Error reading imported file:" + filePath)
	}

	// Crear un nuevo parser con el directorio base actualizado
	parser := NewParser(string(content))
	parser.SetBaseDir(filepath.Dir(filePath))

	// Parsear el programa importado
	importedProgram := parser.ParseProgram()

	// Crear un nuevo entorno para el módulo importado
	moduleEnv := NewInnerEnv(env)
	moduleEnv.Set("currentFile", filePath)
	moduleEnv.Dir = filepath.Dir(filePath)

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
