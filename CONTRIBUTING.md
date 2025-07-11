# Contributing to R2Lang

¡Gracias por tu interés en contribuir a R2Lang! Este documento te guiará a través del proceso de contribución y te ayudará a entender cómo puedes ayudar a mejorar el proyecto.

## 🎯 Formas de Contribuir

### Para Principiantes
- 📝 **Documentación**: Mejorar docs, ejemplos, tutoriales
- 🐛 **Bug Reports**: Reportar problemas encontrados
- 🧪 **Testing**: Escribir tests para funcionalidades existentes
- 📚 **Bibliotecas**: Contribuir a `pkg/r2libs/` con nuevas funciones

### Para Desarrolladores Intermedios
- ⚡ **Performance**: Optimizaciones en `pkg/r2core/`
- 🎨 **Features**: Nuevas características del lenguaje
- 🛠️ **Tooling**: VS Code extension, herramientas de desarrollo
- 🔧 **Infrastructure**: CI/CD, builds, deployment

### Para Desarrolladores Avanzados
- 🏗️ **Architecture**: Mejoras arquitecturales cross-module
- 🚀 **Compiler**: Optimizations en parser y evaluator
- 🔒 **Security**: Auditorías y mejoras de seguridad
- 📊 **Performance**: Profiling y optimizaciones avanzadas

## 🏗️ Arquitectura del Proyecto

### Estructura de Módulos

```
R2Lang/
├── pkg/r2core/         # 🔧 Núcleo del intérprete
│   ├── lexer.go        # Tokenización
│   ├── parse.go        # Análisis sintáctico
│   ├── environment.go  # Gestión de variables
│   ├── *.go           # 30 archivos especializados
│   └── README.md      # Documentación del módulo
├── pkg/r2libs/        # 📚 Bibliotecas extensibles
│   ├── r2http.go      # Servidor HTTP
│   ├── r2math.go      # Operaciones matemáticas
│   ├── r2*.go         # 18 bibliotecas
│   └── README.md      # Guía de bibliotecas
├── pkg/r2lang/        # 🎯 Coordinador principal
│   └── r2lang.go      # Orquestación de alto nivel
├── pkg/r2repl/        # 💻 REPL independiente
│   └── r2repl.go      # Interfaz interactiva
├── examples/          # 📝 Ejemplos de uso
├── docs/              # 📖 Documentación
└── vscode_syntax_highlighting/ # 🎨 Extensión VS Code
```

### Principios de Diseño

1. **Single Responsibility**: Cada archivo tiene una responsabilidad clara
2. **Loose Coupling**: Mínimas dependencias entre módulos
3. **High Cohesion**: Funcionalidad relacionada agrupada
4. **Testability**: Cada módulo debe ser testeable independientemente

## 🚀 Getting Started

### 1. Setup del Entorno

```bash
# Clonar el repositorio
git clone https://github.com/arturoeanton/go-r2lang.git
cd go-r2lang

# Instalar dependencias
go mod download

# Verificar que funciona
go build -o r2lang main.go
./r2lang examples/example1-if.r2

# Ejecutar tests existentes
go test ./pkg/...
```

### 2. Primer Contribución - Bug Fix o Documentation

```bash
# Crear una rama para tu contribución
git checkout -b fix/descripcion-del-fix

# Hacer tus cambios
# ...

# Ejecutar tests
go test ./pkg/...

# Verificar que no rompes nada
go build main.go
./main examples/example*.r2

# Commit con mensaje descriptivo
git commit -m "fix: descripción clara del problema resuelto"

# Push y crear PR
git push origin fix/descripcion-del-fix
```

## 📝 Guías por Tipo de Contribución

### 🔧 Contribuir a pkg/r2core/

**Ideal para**: Developers con experiencia en compiladores/intérpretes

#### Archivos Principales:
- `lexer.go`: Tokenización del código fuente
- `parse.go`: Construcción del AST  
- `environment.go`: Gestión de variables y scope
- `*_expression.go`: Evaluación de expresiones
- `*_statement.go`: Ejecución de statements

#### Ejemplo - Agregar Nuevo Tipo de Nodo AST:

```bash
# 1. Crear el archivo del nodo
touch pkg/r2core/my_new_node.go
```

```go
// pkg/r2core/my_new_node.go
package r2core

type MyNewNode struct {
    // Campos específicos del nodo
    Value string
    Line  int
}

func (n *MyNewNode) Eval(env *Environment) interface{} {
    // Lógica de evaluación
    return n.Value
}

func (n *MyNewNode) String() string {
    return fmt.Sprintf("MyNewNode{%s}", n.Value)
}
```

```bash
# 2. Crear tests
touch pkg/r2core/my_new_node_test.go
```

```go
// pkg/r2core/my_new_node_test.go
package r2core

import "testing"

func TestMyNewNode_Eval(t *testing.T) {
    env := NewEnvironment(nil)
    node := &MyNewNode{Value: "test", Line: 1}
    
    result := node.Eval(env)
    if result != "test" {
        t.Errorf("Expected 'test', got %v", result)
    }
}
```

```bash
# 3. Integrar en el parser (si necesario)
# Editar pkg/r2core/parse.go para reconocer la nueva sintaxis

# 4. Ejecutar tests
go test ./pkg/r2core/ -v
```

### 📚 Contribuir a pkg/r2libs/

**Ideal para**: Developers que quieren agregar funcionalidad específica

#### Estructura de una Biblioteca:

```go
// pkg/r2libs/r2mynewlib.go
package r2libs

import (
    "github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func RegisterMyNewLib(env *r2core.Environment) {
    // Función simple
    env.Set("myFunction", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
        // Validar argumentos
        if len(args) != 1 {
            panic("myFunction expects 1 argument")
        }
        
        // Lógica de la función
        return fmt.Sprintf("Processed: %v", args[0])
    }))
    
    // Objeto con múltiples métodos
    myObject := map[string]interface{}{
        "method1": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
            // Lógica del método 1
            return "result1"
        }),
        "method2": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
            // Lógica del método 2
            return "result2"
        }),
    }
    
    env.Set("MyLib", myObject)
}
```

#### Tests para Bibliotecas:

```go
// pkg/r2libs/r2mynewlib_test.go
package r2libs

import (
    "testing"
    "github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func TestMyNewLib_Functions(t *testing.T) {
    env := r2core.NewEnvironment(nil)
    RegisterMyNewLib(env)
    
    // Test función simple
    myFunc, exists := env.Get("myFunction")
    if !exists {
        t.Fatal("myFunction not registered")
    }
    
    result := myFunc.(r2core.BuiltinFunction)("test")
    expected := "Processed: test"
    if result != expected {
        t.Errorf("Expected %s, got %v", expected, result)
    }
}
```

#### Integración en el Sistema:

```go
// pkg/r2lang/r2lang.go - Agregar registro
func RunCode(filename string) {
    // ... código existente ...
    
    // Registrar tu nueva biblioteca
    r2libs.RegisterMyNewLib(env)
    
    // ... resto del código ...
}
```

### 🧪 Contribuir con Tests

**Testing Guidelines:**

1. **Unit Tests**: Cada función debe tener tests
2. **Integration Tests**: Tests que validen interacción entre módulos
3. **Example Tests**: Tests que validen que los ejemplos funcionan

#### Template para Tests:

```go
package r2core // o r2libs

import (
    "testing"
)

func TestFunctionName(t *testing.T) {
    // Arrange - Preparar el test
    env := NewEnvironment(nil)
    input := "test input"
    expected := "expected output"
    
    // Act - Ejecutar la función
    result := FunctionToTest(input, env)
    
    // Assert - Verificar resultado
    if result != expected {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}

func TestFunctionName_EdgeCases(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {"empty input", "", "", false},
        {"nil input", nil, "", true},
        {"large input", "very long string...", "processed", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := FunctionToTest(tt.input)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("FunctionToTest() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if result != tt.expected {
                t.Errorf("FunctionToTest() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### 📖 Contribuir con Documentación

#### Tipos de Documentación:

1. **Code Comments**: Documentar funciones complejas
2. **Module README**: Explicar propósito y uso de cada pkg/
3. **Examples**: Crear ejemplos prácticos en `examples/`
4. **API Documentation**: Documentar funciones públicas
5. **Tutorials**: Guías paso a paso para usuarios

#### Template para Module README:

```markdown
# pkg/r2mymodule/

## Propósito

Breve descripción del propósito del módulo.

## Archivos Principales

- `main_file.go`: Descripción de funcionalidad principal
- `helper_file.go`: Funciones auxiliares
- `types.go`: Definiciones de tipos

## Uso

```go
// Ejemplo de uso básico
example := NewExample()
result := example.Process("input")
```

## API Reference

### func ProcessData(input string) string

Procesa los datos de entrada y retorna resultado formateado.

**Parámetros:**
- `input`: String de entrada a procesar

**Retorna:**
- String procesado

**Ejemplo:**
```go
result := ProcessData("hello")
// result: "Processed: hello"
```

## Testing

```bash
go test ./pkg/r2mymodule/ -v
```
```

## 🎨 Estándares de Código

### Go Style Guidelines

```go
// ✅ CORRECTO: Nombres descriptivos
func ParseExpressionStatement(parser *Parser) *ExpressionStatement {
    return &ExpressionStatement{
        Expression: parser.parseExpression(),
        Line:       parser.currentToken.Line,
    }
}

// ❌ INCORRECTO: Nombres crípticos
func parseExprStmt(p *Parser) *ExprStmt {
    return &ExprStmt{expr: p.parseExpr(), line: p.curTok.Line}
}

// ✅ CORRECTO: Comentarios útiles
// ParseBinaryExpression construye un nodo AST para expresiones binarias
// como "a + b", "x == y", etc. Maneja precedencia de operadores.
func ParseBinaryExpression(left Node, operator string, right Node) *BinaryExpression {
    return &BinaryExpression{Left: left, Op: operator, Right: right}
}

// ❌ INCORRECTO: Comentarios obvios
// ParseBinaryExpression parses binary expressions
func ParseBinaryExpression(left Node, operator string, right Node) *BinaryExpression {
    // ...
}
```

### Error Handling

```go
// ✅ CORRECTO: Error handling explícito
func (env *Environment) Get(name string) (interface{}, error) {
    if val, ok := env.store[name]; ok {
        return val, nil
    }
    if env.outer != nil {
        return env.outer.Get(name)
    }
    return nil, fmt.Errorf("undefined variable: %s", name)
}

// ❌ INCORRECTO: Panic sin contexto
func (env *Environment) Get(name string) interface{} {
    if val, ok := env.store[name]; ok {
        return val
    }
    panic("undefined variable")
}
```

### Testing Standards

```go
// ✅ CORRECTO: Tests descriptivos y completos
func TestEnvironment_Get_ExistingVariable(t *testing.T) {
    env := NewEnvironment(nil)
    env.Set("x", 42)
    
    value, err := env.Get("x")
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    
    if value != 42 {
        t.Errorf("Expected 42, got %v", value)
    }
}

func TestEnvironment_Get_UndefinedVariable(t *testing.T) {
    env := NewEnvironment(nil)
    
    _, err := env.Get("undefined")
    if err == nil {
        t.Error("Expected error for undefined variable")
    }
    
    expectedMsg := "undefined variable: undefined"
    if err.Error() != expectedMsg {
        t.Errorf("Expected error message %q, got %q", expectedMsg, err.Error())
    }
}
```

## 🔄 Proceso de Pull Request

### 1. Antes de Crear el PR

```bash
# Asegúrate que tu código funciona
go build main.go
go test ./pkg/...

# Verifica estilo de código
gofmt -w .
go vet ./...

# Ejecuta ejemplos para verificar compatibilidad
./main examples/example1-if.r2
./main examples/example24-testcase.r2
```

### 2. Creando el PR

**Template de PR:**

```markdown
## Descripción

Breve descripción de los cambios realizados.

## Tipo de Cambio

- [ ] Bug fix (cambio que corrige un issue)
- [ ] Nueva feature (cambio que agrega funcionalidad)
- [ ] Breaking change (cambio que podría romper funcionalidad existente)
- [ ] Documentación
- [ ] Refactoring/optimization

## Testing

- [ ] Tests existentes pasan
- [ ] Agregué tests para nueva funcionalidad
- [ ] Tests de integración verificados
- [ ] Ejemplos funcionan correctamente

## Checklist

- [ ] Mi código sigue el style guide del proyecto
- [ ] He agregado comentarios para código complejo
- [ ] He actualizado documentación si es necesario
- [ ] Mis cambios no generan warnings
- [ ] He agregado tests que prueban mi fix/feature
```

### 3. Review Process

1. **Automated Checks**: CI/CD ejecuta tests automáticamente
2. **Code Review**: Maintainers revisan el código
3. **Testing**: Verificación manual si es necesario
4. **Merge**: Una vez aprobado, se hace merge

## 🐛 Reportando Bugs

### Template de Bug Report

```markdown
## Descripción del Bug

Descripción clara y concisa del problema.

## Para Reproducir

Pasos para reproducir el comportamiento:
1. Ejecutar `r2lang example.r2`
2. Ingresar el siguiente código: '...'
3. Presionar Enter
4. Ver error

## Comportamiento Esperado

Descripción clara de lo que esperabas que pasara.

## Código R2Lang

```r2
// Código que causa el problema
func main() {
    // ...
}
```

## Entorno

- OS: [e.g. macOS, Linux, Windows]
- Go version: [e.g. 1.21]
- R2Lang commit: [e.g. abc123]

## Información Adicional

Cualquier otra información que pueda ayudar.
```

## 💡 Requesting Features

### Template de Feature Request

```markdown
## ¿Tu feature request está relacionada a un problema?

Descripción clara del problema. Ej. "Siempre me frustra cuando..."

## Describe la solución que te gustaría

Descripción clara y concisa de lo que quieres que pase.

## Describe alternativas que has considerado

Descripción de cualquier solución o feature alternativa.

## Ejemplo de Uso

```r2
// Cómo se vería usando la nueva feature
func example() {
    newFeature.doSomething();
}
```

## Contexto Adicional

Agrega cualquier otro contexto o screenshots sobre el feature request.
```

## 🤝 Código de Conducta

### Nuestros Estándares

- ✅ **Ser inclusivo**: Usar lenguaje welcoming e inclusivo
- ✅ **Ser respetuoso**: Diferentes viewpoints y experiencias
- ✅ **Aceptar crítica constructiva**: Con gracia
- ✅ **Focalizarse en la comunidad**: Lo que es mejor para la comunidad
- ✅ **Mostrar empatía**: Hacia otros miembros de la comunidad

### Comportamientos Inaceptables

- ❌ Trolling, insultos o comentarios despectivos
- ❌ Ataques personales o políticos
- ❌ Harassment público o privado
- ❌ Publishing private information sin permiso explícito

## 📞 Obtener Ayuda

### Canales de Comunicación

- **GitHub Issues**: Para bugs y feature requests
- **GitHub Discussions**: Para preguntas generales
- **Email**: [maintainer@r2lang.dev] para asuntos privados

### FAQ para Contributors

**Q: ¿Puedo trabajar en múltiples PRs al mismo tiempo?**
A: Sí, pero recomendamos enfocarse en uno a la vez para mejor calidad.

**Q: ¿Cómo puedo agregar una nueva biblioteca?**
A: Crea un archivo `pkg/r2libs/r2mynewlib.go` y sigue el template mostrado arriba.

**Q: ¿Necesito experiencia previa con compiladores?**
A: No para contribuir a bibliotecas o documentación. Para `pkg/r2core/` sí es recomendable.

**Q: ¿Cómo puedo test mi contribución?**
A: Ejecuta `go test ./pkg/...` y también prueba ejemplos manualmente.

**Q: ¿Puedo proponer cambios breaking?**
A: Sí, pero deben estar bien justificados y ser discutidos en un issue primero.

---

## 🙏 Reconocimientos

Gracias a todos los contributors que hacen R2Lang mejor cada día:

- **Core Team**: [Lista de maintainers]
- **Contributors**: [Todos los que han contribuido]
- **Community**: [Usuarios que reportan bugs y dan feedback]

¡Tu contribución hace la diferencia! 🚀