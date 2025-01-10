package r2lang

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// ============================================================
// 1) TOKENS
// ============================================================

const (
	TOKEN_EOF    = "EOF"
	TOKEN_NUMBER = "NUMBER"
	TOKEN_STRING = "STRING"
	TOKEN_IDENT  = "IDENT"
	TOKEN_SYMBOL = "SYMBOL"
	TOKEN_IMPORT = "IMPORT"
	TOKEN_AS     = "AS"

	RETURN   = "return"
	LET      = "let"
	VAR      = "var"
	FUNC     = "func"
	FUNCTION = "function"
	METHOD   = "method"

	IF      = "if"
	WHILE   = "while"
	FOR     = "for"
	IN      = "in"
	OBJECT  = "obj"
	IMPORT  = "import"
	AS      = "as"
	TRY     = "try"
	CATCH   = "catch"
	FINALLY = "finally"
	THROW   = "throw"

	// Nuevos Tokens para la sintaxis de pruebas
	TOKEN_TESTCASE = "TESTCASE"
	TOKEN_GIVEN    = "GIVEN"
	TOKEN_WHEN     = "WHEN"
	TOKEN_THEN     = "THEN"
	TOKEN_AND      = "AND"
)

var (
	wg sync.WaitGroup
)

type Token struct {
	Type  string
	Value string
	Line  int
	Pos   int
}

// ============================================================
// 2) LEXER (reconoce signo, decimales, comentarios, etc.)
// ============================================================

type Lexer struct {
	input        string
	pos          int
	length       int
	currentToken Token
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		pos:    0,
		length: len(input),
	}
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\n' || ch == '\t' || ch == '\r'
}
func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		ch == '_' ||
		ch == '$'
}
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// parseNumberOrSign maneja -2.3, +10, etc.
func (l *Lexer) parseNumberOrSign() Token {
	start := l.pos
	if l.input[l.pos] == '-' || l.input[l.pos] == '+' {
		l.pos++
	}
	hasDigits := false
	for l.pos < l.length && isDigit(l.input[l.pos]) {
		hasDigits = true
		l.pos++
	}
	if l.pos < l.length && l.input[l.pos] == '.' {
		l.pos++
		for l.pos < l.length && isDigit(l.input[l.pos]) {
			hasDigits = true
			l.pos++
		}
	}
	if !hasDigits {
		panic("Número inválido en " + l.input[start:l.pos])
	}
	val := l.input[start:l.pos]
	l.currentToken = Token{Type: TOKEN_NUMBER, Value: val, Line: 0, Pos: l.pos}
	return l.currentToken
}

func (l *Lexer) NextToken() Token {
skipWhitespace:
	for l.pos < l.length {
		ch := l.input[l.pos]
		if isWhitespace(ch) {
			l.pos++
		} else if ch == '/' {
			// Comentarios
			if l.pos+1 < l.length && l.input[l.pos+1] == '/' {
				// comentario de línea
				l.pos += 2
				for l.pos < l.length && l.input[l.pos] != '\n' {
					l.pos++
				}
			} else if l.pos+1 < l.length && l.input[l.pos+1] == '*' {
				// /* ... */
				l.pos += 2
				for l.pos < l.length {
					if l.input[l.pos] == '*' && (l.pos+1 < l.length && l.input[l.pos+1] == '/') {
						l.pos += 2
						break
					}
					l.pos++
				}
			} else {
				break skipWhitespace
			}
		} else {
			break skipWhitespace
		}
	}

	if l.pos >= l.length {
		l.currentToken = Token{Type: TOKEN_EOF, Value: "", Line: 0, Pos: l.pos}
		return l.currentToken
	}

	ch := l.input[l.pos]

	// Números con signo y operadores
	// busca signos + o - seguidos de dígitos y que no estén precedidos por (, [, , o =
	if ch == '-' || ch == '+' {
		// Signo + digitos => parseNumberOrSign
		pos := l.pos - 1
		for pos > 0 && l.input[pos] == ' ' {
			pos--
		}
		if (l.input[pos] == '(' || l.input[pos] == ',' || l.input[pos] == '[' || l.input[pos] == '=') &&
			(l.pos+1 < l.length && isDigit(l.input[l.pos+1])) {
			return l.parseNumberOrSign()
		}
	}

	if ch == '+' {
		nextch := l.input[l.pos+1]
		if nextch == '+' {
			l.currentToken = Token{Type: TOKEN_SYMBOL, Value: "++", Line: 0, Pos: l.pos}
			l.pos += 2
			return l.currentToken
		}
	}

	if ch == '-' {
		nextch := l.input[l.pos+1]
		if nextch == '-' {
			l.currentToken = Token{Type: TOKEN_SYMBOL, Value: "--", Line: 0, Pos: l.pos}
			l.pos += 2
			return l.currentToken
		}
	}

	// Símbolos de 1 caracter
	singleCharSymbols := []string{
		"(", ")", "{", "}", "[", "]", ";", ",", "+", "-", "*", "/", ".", ":",
	}
	for _, s := range singleCharSymbols {
		if string(ch) == s {
			l.pos++
			l.currentToken = Token{Type: TOKEN_SYMBOL, Value: s, Line: 0, Pos: l.pos}
			return l.currentToken
		}
	}

	if string(ch) == "=" {
		if l.pos+1 < l.length && l.input[l.pos+1] == '=' {
			l.pos += 2
			l.currentToken = Token{Type: TOKEN_SYMBOL, Value: "==", Line: 0, Pos: l.pos}
			return l.currentToken
		}
		l.pos++
		l.currentToken = Token{Type: TOKEN_SYMBOL, Value: "=", Line: 0, Pos: l.pos}
		return l.currentToken
	}

	// Operadores relacionales
	if ch == '<' || ch == '>' || ch == '!' || ch == '=' {
		if l.pos+1 < l.length {
			nextCh := l.input[l.pos+1]
			if nextCh == '=' {
				op := string(ch) + string(nextCh)
				l.pos += 2
				l.currentToken = Token{Type: TOKEN_SYMBOL, Value: op, Line: 0, Pos: l.pos}
				return l.currentToken
			}
		}
		l.pos++
		l.currentToken = Token{Type: TOKEN_SYMBOL, Value: string(ch), Line: 0, Pos: l.pos}
		return l.currentToken
	}

	// Cadena
	if ch == '"' || ch == '\'' {
		quote := ch
		start := l.pos
		l.pos++
		for l.pos < l.length && l.input[l.pos] != quote {
			l.pos++
		}
		if l.pos >= l.length {
			panic("Se esperaba comilla de cierre de cadena")
		}
		val := l.input[start+1 : l.pos]
		l.pos++
		l.currentToken = Token{Type: TOKEN_STRING, Value: val, Line: 0, Pos: l.pos}
		return l.currentToken
	}

	// Números sin signo
	if isDigit(ch) {
		return l.parseNumberOrSign()
	}

	// Identificadores
	if isLetter(ch) {
		start := l.pos
		for l.pos < l.length && (isLetter(l.input[l.pos]) || isDigit(l.input[l.pos])) {
			l.pos++
		}
		literal := l.input[start:l.pos]
		switch strings.ToLower(literal) {
		case strings.ToLower(IMPORT):
			l.currentToken = Token{Type: TOKEN_IMPORT, Value: literal, Line: 0, Pos: l.pos}
			return l.currentToken
		case strings.ToLower(AS):
			l.currentToken = Token{Type: TOKEN_AS, Value: literal, Line: 0, Pos: l.pos}
			return l.currentToken
		case "given":
			l.currentToken = Token{Type: TOKEN_GIVEN, Value: "Given", Line: 0, Pos: l.pos}
			return l.currentToken
		case "when":
			l.currentToken = Token{Type: TOKEN_WHEN, Value: "When", Line: 0, Pos: l.pos}
			return l.currentToken
		case "then":
			l.currentToken = Token{Type: TOKEN_THEN, Value: "Then", Line: 0, Pos: l.pos}
			return l.currentToken
		case "and":
			l.currentToken = Token{Type: TOKEN_AND, Value: "And", Line: 0, Pos: l.pos}
			return l.currentToken
		case "testcase":
			l.currentToken = Token{Type: TOKEN_TESTCASE, Value: "TestCase", Line: 0, Pos: l.pos}
			return l.currentToken
		// ... otras palabras clave
		default:
			l.currentToken = Token{Type: TOKEN_IDENT, Value: literal, Line: 0, Pos: l.pos}
			return l.currentToken
		}
	}

	panic(fmt.Sprintf("Caracter inesperado en lexer: %c", ch))
}

// ============================================================
// 3) AST - Node interface
// ============================================================

type Node interface {
	Eval(env *Environment) interface{}
}

type NodeTest interface {
	Eval(env *Environment) interface{}
	EvalStep(env *Environment) interface{}
}

// ============================================================
// 3.1) Nuevos Nodos para TestCase y TestStep
// ============================================================

// TestCase representa un caso de prueba con un nombre y pasos.
type TestCase struct {
	Name  string
	Steps []TestStep
}

type TestStep struct {
	Type    string // "Given", "When", "Then", "And"
	Command Node
}

// Eval ejecuta el caso de prueba.
func (tc *TestCase) Eval(env *Environment) interface{} {
	fmt.Printf("Ejecutando Caso de Prueba: %s\n", tc.Name)
	var previousStepType string

	for _, step := range tc.Steps {
		stepType := step.Type
		if stepType == "And" {
			stepType = previousStepType
		} else {
			previousStepType = stepType
		}
		fmt.Printf("  %s: ", stepType)

		if ce, ok := step.Command.(*CallExpression); ok {
			calleeVal := ce.Callee.Eval(env)
			var argVals []interface{}
			for _, a := range ce.Args {
				argVals = append(argVals, a.Eval(env))
			}
			if currentStep, ok := calleeVal.(*UserFunction); ok {
				out := currentStep.CallStep(env, argVals...)
				if out != nil {
					fmt.Println(out)
				}
			}
			continue
		}

		if fl, ok := step.Command.(*FunctionLiteral); ok {
			currentStep := fl.Eval(env).(*UserFunction)
			out := currentStep.CallStep(env)
			if out != nil {
				fmt.Println(out)
			}
			continue
		}

	}
	fmt.Println("Caso de Prueba Ejecutado Exitosamente.\n")
	return nil
}

func (ts *TestStep) Eval(env *Environment) interface{} {
	// Ejecutar el comando del paso
	return ts.Command.Eval(env)
}

// ============================================================
// UTILS
// ============================================================

func isBinaryOp(op string) bool {
	ops := []string{"+", "-", "*", "/", "<", ">", "<=", ">=", "==", "!="}
	for _, o := range ops {
		if op == o {
			return true
		}
	}
	return false
}

// Queremos un "FunctionLiteral" para soportar func(...) { ... } anónimas
type FunctionLiteral struct {
	Args []string
	Body *BlockStatement
}

func (fl *FunctionLiteral) Eval(env *Environment) interface{} {
	fn := &UserFunction{
		Args:     fl.Args,
		Body:     fl.Body,
		Env:      env, // closure
		IsMethod: false,
	}
	return fn
}

// ============================================================
// 4) STATEMENTS & EXPRESSIONS
// ============================================================

type Program struct {
	Statements []Node
}

func (p *Program) Eval(env *Environment) interface{} {

	var result interface{}
	for _, stmt := range p.Statements {
		val := stmt.Eval(env)

		if rv, ok := val.(ReturnValue); ok {
			return rv.Value
		}
		result = val
	}
	return result
}

type ReturnValue struct {
	Value interface{}
}

// LetStatement => let x = expr;
type LetStatement struct {
	Name  string
	Value Node
}

func (ls *LetStatement) Eval(env *Environment) interface{} {
	var val interface{}
	if ls.Value != nil {
		val = ls.Value.Eval(env)
	}
	env.Set(ls.Name, val)
	return nil
}

type GenericAssignStatement struct {
	Left  Node
	Right Node
}

func (gas *GenericAssignStatement) Eval(env *Environment) interface{} {

	val := gas.Right.Eval(env)
	switch left := gas.Left.(type) {
	case *Identifier:
		env.Set(left.Name, val)
		return val
	case *AccessExpression:
		objVal := left.Object.Eval(env)
		instance, ok := objVal.(*ObjectInstance)
		if !ok {
			panic("No es un objeto-instance para asignar .prop")
		}
		instance.Env.Set(left.Member, val)
		return val
	case *IndexExpression:
		return assignIndexExpression(left, val, env)
	default:
		panic("No se puede asignar a esta expresión")
	}
}

type ExprStatement struct {
	Expr Node
}

func (es *ExprStatement) Eval(env *Environment) interface{} {
	return es.Expr.Eval(env)
}

type IfStatement struct {
	Condition   Node
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ifs *IfStatement) Eval(env *Environment) interface{} {
	condVal := ifs.Condition.Eval(env)
	if toBool(condVal) {
		return ifs.Consequence.Eval(env)
	} else if ifs.Alternative != nil {
		return ifs.Alternative.Eval(env)
	}
	return nil
}

// ImportStatement representa una declaración de importación con alias.
type ImportStatement struct {
	Path  string
	Alias string // Alias opcional
}

func (is *ImportStatement) Eval(env *Environment) interface{} {
	filePath := is.Path

	// Resolver rutas relativas
	if !filepath.IsAbs(filePath) {
		dir := env.dir
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
		panic("Error al leer archivo importado: " + filePath)
	}

	// Crear un nuevo parser con el directorio base actualizado
	parser := NewParser(string(content))
	parser.SetBaseDir(filepath.Dir(filePath))

	// Parsear el programa importado
	importedProgram := parser.ParseProgram()

	// Crear un nuevo entorno para el módulo importado
	moduleEnv := NewInnerEnv(env)
	moduleEnv.Set("currentFile", filePath)
	moduleEnv.dir = filepath.Dir(filePath)

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

type WhileStatement struct {
	Condition Node
	Body      *BlockStatement
}

func (ws *WhileStatement) Eval(env *Environment) interface{} {
	var result interface{}
	for {
		condVal := ws.Condition.Eval(env)
		if !toBool(condVal) {
			break
		}
		val := ws.Body.Eval(env)
		if rv, ok := val.(ReturnValue); ok {
			return rv
		}
		result = val
	}
	return result
}

type TryStatement struct {
	Body         *BlockStatement
	CatchBlock   *BlockStatement
	FinallyBlock *BlockStatement
	ExceptionVar string
}

func (ts *TryStatement) Eval(env *Environment) interface{} {
	defer func() {
		if r := recover(); r != nil {
			if ts.CatchBlock != nil {
				newEnv := NewInnerEnv(env)
				newEnv.Set(ts.ExceptionVar, r)
				ts.CatchBlock.Eval(newEnv)
			}
		}
		if ts.FinallyBlock != nil {
			ts.FinallyBlock.Eval(env)
		}
	}()
	return ts.Body.Eval(env)
}

type ThrowStatement struct {
	Message string
}

func (ts *ThrowStatement) Eval(env *Environment) interface{} {
	panic(ts.Message)
	return nil
}

type ForStatement struct {
	Init        Node
	Condition   Node
	Post        Node
	Body        *BlockStatement
	inFlag      bool
	inArray     string
	inMap       string
	inIndexName string
}

func (fs *ForStatement) Eval(env *Environment) interface{} {
	newEnv := NewInnerEnv(env)

	var result interface{}

	var arr []interface{}
	var mapVal map[string]interface{}
	var ok bool
	flagArr := true
	if fs.inFlag {
		var raw interface{}
		if _, ok = fs.Init.(*CallExpression); ok {
			raw = fs.Init.Eval(newEnv)
			newEnv.Set("$c", raw)
		} else {
			raw, _ = newEnv.Get(fs.inArray)
			newEnv.Set("$c", raw)
		}

		arr, ok = raw.([]interface{})
		if !ok {
			flagArr = false
			mapVal, ok = raw.(map[string]interface{})
			if !ok {
				panic("No es un array ni mapa para for")
			}
		}
	}
	if fs.inFlag {
		if flagArr {
			for i, v := range arr {
				newEnv.Set(fs.inIndexName, float64(i))
				newEnv.Set("$k", float64(i))
				newEnv.Set("$v", v)
				val := fs.Body.Eval(newEnv)
				if rv, ok := val.(ReturnValue); ok {
					return rv
				}
				result = val
				if fs.Post != nil {
					fs.Post.Eval(newEnv)
				}
			}
		} else {
			for k, v := range mapVal {
				newEnv.Set(fs.inIndexName, k)
				newEnv.Set("$k", k)
				newEnv.Set("$v", v)
				val := fs.Body.Eval(newEnv)
				if rv, ok := val.(ReturnValue); ok {
					return rv
				}
				result = val
				if fs.Post != nil {
					fs.Post.Eval(newEnv)
				}
			}
		}
		return result
	}

	if fs.Init != nil {
		fs.Init.Eval(newEnv)
	}

	for {
		condVal := fs.Condition.Eval(newEnv)
		if !toBool(condVal) {
			break
		}
		val := fs.Body.Eval(newEnv)
		if rv, ok := val.(ReturnValue); ok {
			return rv
		}
		result = val
		if fs.Post != nil {
			fs.Post.Eval(newEnv)
		}
	}
	return result
}

type BlockStatement struct {
	Statements []Node
}

func (bs *BlockStatement) Eval(env *Environment) interface{} {
	var result interface{}
	for _, stmt := range bs.Statements {
		val := stmt.Eval(env)
		if rv, ok := val.(ReturnValue); ok {
			return rv
		}
		result = val
	}
	return result
}

type ReturnStatement struct {
	Value Node
}

func (rs *ReturnStatement) Eval(env *Environment) interface{} {
	if rs.Value == nil {
		return ReturnValue{Value: nil}
	}
	val := rs.Value.Eval(env)
	return ReturnValue{Value: val}
}

// Función con nombre
type FunctionDeclaration struct {
	Name string
	Args []string
	Body *BlockStatement
}

func (fd *FunctionDeclaration) Eval(env *Environment) interface{} {
	fn := &UserFunction{
		Args:     fd.Args,
		Body:     fd.Body,
		Env:      env,
		IsMethod: false,
		code:     fd.Name,
	}
	env.Set(fd.Name, fn)
	return nil
}

// obj MiObj { let..., func... }
type ObjectDeclaration struct {
	Name    string
	Members []Node
}

func (od *ObjectDeclaration) Eval(env *Environment) interface{} {
	blueprint := make(map[string]interface{})
	for _, m := range od.Members {
		switch node := m.(type) {
		case *LetStatement:
			blueprint[node.Name] = nil
		case *FunctionDeclaration:
			fn := &UserFunction{
				Args:     node.Args,
				Body:     node.Body,
				Env:      nil,
				IsMethod: true,
			}
			blueprint[node.Name] = fn
		}
	}

	env.Set(od.Name, blueprint)
	return nil
}

// -------------- EXPRESSIONS --------------

type Identifier struct {
	Name string
}

func (id *Identifier) Eval(env *Environment) interface{} {
	val, ok := env.Get(id.Name)
	if !ok {
		panic("Variable no declarada: " + id.Name)
	}
	return val
}

type NumberLiteral struct {
	Value float64
}

func (nl *NumberLiteral) Eval(env *Environment) interface{} {
	return nl.Value
}

type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) Eval(env *Environment) interface{} {
	return sl.Value
}

type BooleanLiteral struct {
	Value bool
}

func (b *BooleanLiteral) Eval(env *Environment) interface{} {
	return b.Value
}

type ArrayLiteral struct {
	Elements []Node
}

func (al *ArrayLiteral) Eval(env *Environment) interface{} {
	var result []interface{}
	for _, e := range al.Elements {
		result = append(result, e.Eval(env))
	}
	return result
}

type MapLiteral struct {
	Pairs map[string]Node
}

func (ml *MapLiteral) Eval(env *Environment) interface{} {
	m := make(map[string]interface{})
	for k, nd := range ml.Pairs {
		m[k] = nd.Eval(env)
	}
	return m
}

type BinaryExpression struct {
	Left  Node
	Op    string
	Right Node
}

func (be *BinaryExpression) Eval(env *Environment) interface{} {
	lv := be.Left.Eval(env)
	rv := be.Right.Eval(env)

	switch be.Op {
	case "+":
		return addValues(lv, rv)
	case "-":
		return subValues(lv, rv)
	case "*":
		return mulValues(lv, rv)
	case "/":
		return divValues(lv, rv)
	case "<":
		return toFloat(lv) < toFloat(rv)
	case ">":
		return toFloat(lv) > toFloat(rv)
	case "<=":
		return toFloat(lv) <= toFloat(rv)
	case ">=":
		return toFloat(lv) >= toFloat(rv)
	case "==":
		return equals(lv, rv)
	case "!=":
		return !equals(lv, rv)
	default:
		panic("Operador binario no soportado: " + be.Op)
	}
}

// CallExpression => callee(args)
type CallExpression struct {
	Callee Node
	Args   []Node
}

func (ce *CallExpression) Eval(env *Environment) interface{} {
	calleeVal := ce.Callee.Eval(env)
	var argVals []interface{}
	for _, a := range ce.Args {
		argVals = append(argVals, a.Eval(env))
	}
	switch cv := calleeVal.(type) {
	case BuiltinFunction:
		return cv(argVals...)
	case *UserFunction:
		return cv.Call(argVals...)
	case map[string]interface{}:
		// Instanciar un blueprint
		return instantiateObject(env, cv, argVals)
	default:
		panic("Intento de llamar algo que no es función ni blueprint [" + fmt.Sprintf("%T", ce.Callee) + "]")
	}
}

type AccessExpression struct {
	Object Node
	Member string
}

func (ae *AccessExpression) Eval(env *Environment) interface{} {
	objVal := ae.Object.Eval(env)

	// Manejar ObjectInstance
	if instance, ok := objVal.(*ObjectInstance); ok {
		val, exists := instance.Env.Get(ae.Member)
		if !exists {
			panic("El objeto no tiene la propiedad: " + ae.Member)
		}
		return val
	}

	// Manejar map[string]interface{}
	if m, ok := objVal.(map[string]interface{}); ok {
		val, exists := m[ae.Member]
		if !exists {
			panic("El mapa no tiene la clave: " + ae.Member)
		}
		return val
	}

	panic("Acceso a propiedad en tipo no soportado: " + fmt.Sprintf("%T", objVal))
}

type IndexExpression struct {
	Left  Node
	Index Node
}

func (ie *IndexExpression) Eval(env *Environment) interface{} {
	leftVal := ie.Left.Eval(env)
	indexVal := ie.Index.Eval(env)

	switch container := leftVal.(type) {
	case map[string]interface{}:
		strKey, ok := indexVal.(string)
		if !ok {
			panic("índice debe ser string para map")
		}
		vv, found := container[strKey]
		if !found {
			return nil
		}
		return vv
	case []interface{}:
		fIndex, ok := indexVal.(float64)
		if !ok {
			panic("índice debe ser numérico para array")
		}
		idx := int(fIndex)
		if idx < 0 || idx >= len(container) {
			return nil
		}
		return container[idx]
	default:
		panic("índice sobre algo que no es map ni array")
	}
}

// ============================================================
// 5) OBJETO, FUNCIONES, ENV
// ============================================================

type UserFunction struct {
	Args     []string
	Body     *BlockStatement
	Env      *Environment
	IsMethod bool
	code     string
}

func (uf *UserFunction) NativeCall(currentEnv *Environment, args ...interface{}) interface{} {
	newEnv := currentEnv
	if newEnv == nil {
		newEnv = NewInnerEnv(uf.Env)
	}
	if uf.IsMethod {
		if selfVal, ok := uf.Env.Get("self"); ok {
			newEnv.Set("self", selfVal)
		}
	}
	for i, param := range uf.Args {
		if i < len(args) {
			newEnv.Set(param, args[i])
		} else {
			newEnv.Set(param, nil)
		}
	}
	val := uf.Body.Eval(newEnv)
	if rv, ok := val.(ReturnValue); ok {
		return rv.Value
	}
	return val
}

func (uf *UserFunction) Call(args ...interface{}) interface{} {
	tmp := uf.Env.CurrenFx
	uf.Env.CurrenFx = uf.code
	out := uf.NativeCall(nil, args...)
	uf.Env.CurrenFx = tmp
	return out
}

func (uf *UserFunction) CallStep(env *Environment, args ...interface{}) interface{} {
	return uf.NativeCall(env, args...)
}

type BuiltinFunction func(args ...interface{}) interface{}

type ObjectInstance struct {
	Env *Environment
}

func instantiateObject(env *Environment, blueprint map[string]interface{}, argVals []interface{}) *ObjectInstance {
	objEnv := NewInnerEnv(env)
	instance := &ObjectInstance{Env: objEnv}
	for k, v := range blueprint {
		switch vv := v.(type) {
		case *UserFunction:
			newFn := &UserFunction{
				Args:     vv.Args,
				Body:     vv.Body,
				Env:      objEnv,
				IsMethod: true,
			}
			objEnv.Set(k, newFn)
		default:
			objEnv.Set(k, vv)
		}
	}
	objEnv.Set("self", instance)

	return instance
}

// ============================================================
// 6) ENVIRONMENT
// ============================================================

type Environment struct {
	store    map[string]interface{}
	outer    *Environment
	imported map[string]bool
	dir      string
	CurrenFx string
}

func NewEnvironment() *Environment {
	return &Environment{
		store:    make(map[string]interface{}),
		outer:    nil,
		imported: make(map[string]bool),
	}
}

func NewInnerEnv(outer *Environment) *Environment {
	return &Environment{
		store:    make(map[string]interface{}),
		outer:    outer,
		imported: make(map[string]bool),
		dir:      outer.dir,
	}
}

func (e *Environment) Set(name string, value interface{}) {
	e.store[name] = value
}

func (e *Environment) Get(name string) (interface{}, bool) {
	val, ok := e.store[name]
	if ok {
		return val, true
	}
	if e.outer != nil {
		return e.outer.Get(name)
	}
	return nil, false
}

func (e *Environment) Run(parser *Parser) (result interface{}) {

	defer wg.Wait()
	wg = sync.WaitGroup{}

	ast := parser.ParseProgram()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Exception:", r)
			fmt.Println("Current fx -> ", e.CurrenFx)
		}
	}()

	e.CurrenFx = "."

	// Ejecutar
	result = ast.Eval(e)

	// Llamar a main() si está
	mainVal, ok := e.Get("main")
	if ok {
		mainFn, isFn := mainVal.(*UserFunction)
		if !isFn {
			fmt.Println("Error: 'main' no es una función.")
			os.Exit(1)
		}
		result = mainFn.Call()
	}
	return result
}

// ============================================================
// 7) UTILS
// ============================================================

func toFloat(val interface{}) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case bool:
		if v {
			return 1
		}
		return 0
	case nil:
		return 0
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			panic("No se puede convertir string a número: " + v)
		}
		return f
	}
	panic("No se puede convertir valor a número")
}
func toBool(val interface{}) bool {
	if val == nil {
		return false
	}
	switch v := val.(type) {
	case bool:
		return v
	case float64:
		return v != 0
	case int:
		return v != 0
	case string:
		return v != ""
	}
	return true
}

// Para unificar la lógica numérica en "=="
func isNumeric(v interface{}) bool {
	switch v.(type) {
	case float64, int:
		return true
	}
	return false
}

// Corrige la comparación "=="
func equals(a, b interface{}) bool {
	// Si ambos son numéricos, compare con toFloat
	if isNumeric(a) && isNumeric(b) {
		return toFloat(a) == toFloat(b)
	}
	// sino comparamos string/bool/nil
	switch aa := a.(type) {
	case string:
		if bb, ok := b.(string); ok {
			return aa == bb
		}
	case bool:
		if bb, ok := b.(bool); ok {
			return aa == bb
		}
	case nil:
		return b == nil
	}
	return false
}

func addValues(a, b interface{}) interface{} {
	// Si uno es string => concatenar
	if sa, ok := a.(string); ok {
		return sa + fmt.Sprint(b)
	}
	if sb, ok := b.(string); ok {
		return fmt.Sprint(a) + sb
	}
	return toFloat(a) + toFloat(b)
}
func subValues(a, b interface{}) interface{} {
	return toFloat(a) - toFloat(b)
}
func mulValues(a, b interface{}) interface{} {
	return toFloat(a) * toFloat(b)
}
func divValues(a, b interface{}) interface{} {
	den := toFloat(b)
	if den == 0 {
		panic("División por cero")
	}
	return toFloat(a) / den
}

// Asignación en map/array
func assignIndexExpression(idxExpr *IndexExpression, newVal interface{}, env *Environment) interface{} {
	leftVal := idxExpr.Left.Eval(env)
	indexVal := idxExpr.Index.Eval(env)

	switch container := leftVal.(type) {
	case map[string]interface{}:
		key, ok := indexVal.(string)
		if !ok {
			panic("assignIndexExpression: índice para map debe ser string")
		}
		container[key] = newVal
		return newVal
	case []interface{}:
		idxF, ok := indexVal.(float64)
		if !ok {
			panic("assignIndexExpression: índice array debe ser número")
		}
		idx := int(idxF)
		if idx < 0 {
			panic("Índice negativo en array")
		}
		// auto-extender
		if idx >= len(container) {
			for len(container) <= idx {
				container = append(container, nil)
			}
		}
		container[idx] = newVal
		return newVal
	default:
		panic("No es un map ni array para asignar indice")
	}
}

// ============================================================
// 8) PARSER
// ============================================================

type Parser struct {
	lexer   *Lexer
	curTok  Token
	peekTok Token
	baseDir string // Directorio base para importaciones
}

func NewParser(input string) *Parser {
	p := &Parser{lexer: NewLexer(input)}
	p.nextToken()
	p.nextToken()
	p.baseDir = ""
	return p
}

// SetBaseDir establece el directorio base para importaciones
func (p *Parser) SetBaseDir(dir string) {
	p.baseDir = dir
}

func (p *Parser) parseImportStatement() Node {
	p.nextToken() // Consumir 'import'

	if p.curTok.Type != TOKEN_STRING {
		panic("Se esperaba una cadena de caracteres después de 'import'")
	}

	path := p.curTok.Value
	p.nextToken()

	var alias string
	if p.curTok.Type == TOKEN_AS {
		p.nextToken() // Consumir 'as'
		if p.curTok.Type != TOKEN_IDENT {
			panic("Se esperaba un identificador tras 'as'")
		}
		alias = p.curTok.Value
		p.nextToken()
	}

	if p.curTok.Value == ";" {
		p.nextToken() // Consumir ';'
	}

	return &ImportStatement{Path: path, Alias: alias}
}

func (p *Parser) parseTestCase() Node {
	p.nextToken() // Consumir 'TestCase'

	if p.curTok.Type != TOKEN_STRING {
		panic("Se esperaba una cadena de caracteres para el nombre del caso de prueba")
	}
	name := p.curTok.Value
	p.nextToken()

	if p.curTok.Value != "{" {
		panic("Se esperaba '{' para iniciar el cuerpo del caso de prueba")
	}
	p.nextToken()

	var steps []TestStep
	for p.curTok.Value != "}" && p.curTok.Type != TOKEN_EOF {
		var stepType string
		switch p.curTok.Type {
		case TOKEN_GIVEN, TOKEN_WHEN, TOKEN_THEN, TOKEN_AND:
			stepType = p.curTok.Value
			p.nextToken()
		default:
			panic("Se esperaba 'Given', 'When', 'Then' o 'And' en los pasos del caso de prueba")
		}
		command := p.parseExpression()
		steps = append(steps, TestStep{Type: stepType, Command: command})
		if p.curTok.Value == ";" {
			p.nextToken()
		}
	}
	if p.curTok.Value != "}" {
		panic("Se esperaba '}' al final del caso de prueba")
	}
	p.nextToken()
	return &TestCase{Name: name, Steps: steps}
}

func (p *Parser) nextToken() {
	p.curTok = p.peekTok
	p.peekTok = p.lexer.NextToken()
}

func (p *Parser) ParseProgram() *Program {
	prog := &Program{}
	for p.curTok.Type != TOKEN_EOF {
		stmt := p.parseStatement()
		prog.Statements = append(prog.Statements, stmt)
	}
	return prog
}

func (p *Parser) parseThrowStatement() Node {
	p.nextToken()
	if p.curTok.Type != TOKEN_STRING {
		panic("Se esperaba una cadena de caracteres para el mensaje de excepción")
	}
	message := fmt.Sprint(p.curTok.Value)
	return &ThrowStatement{Message: message}
}

func (p *Parser) parseStatement() Node {

	if p.curTok.Value == IMPORT {
		return p.parseImportStatement()
	}

	if p.curTok.Type == TOKEN_TESTCASE {
		return p.parseTestCase()
	}

	if p.curTok.Value == TRY {
		return p.parseTryStatement()
	}

	if p.curTok.Value == THROW {
		return p.parseThrowStatement()
	}

	if p.curTok.Value == RETURN {
		return p.parseReturnStatement()
	}
	if p.curTok.Value == LET || p.curTok.Value == VAR {
		return p.parseLetStatement()
	}
	if p.curTok.Value == FUNC || p.curTok.Value == FUNCTION {
		// esto parsea "func nombre(...) { ... }" => FunctionDeclaration con nombre
		return p.parseFunctionDeclaration()
	}
	if p.curTok.Value == IF {
		return p.parseIfStatement()
	}
	if p.curTok.Value == WHILE {
		return p.parseWhileStatement()
	}
	if p.curTok.Value == FOR {
		return p.parseForStatement()
	}

	if p.curTok.Value == OBJECT {
		return p.parseObjectDeclaration()
	}
	// sino parseAsignmentOrExpressionStatement
	return p.parseAssignmentOrExpressionStatement()
}

func (p *Parser) parseTryStatement() Node {
	p.nextToken() // consumir "try"
	body := p.parseBlockStatement()
	exceptionVar := "$e"
	var catchBlock *BlockStatement
	if p.curTok.Value == CATCH {
		p.nextToken() // consumir "catch"

		if p.curTok.Value == "{" {
			catchBlock = p.parseBlockStatement()
		} else {

			if p.curTok.Value != "(" {
				panic("Se esperaba '(' después de 'catch'")
			}
			p.nextToken() // consumir "("
			if p.curTok.Type != TOKEN_IDENT {
				panic("Se esperaba nombre de variable de excepción")
			}
			exceptionVar = p.curTok.Value
			p.nextToken()
			if p.curTok.Value != ")" {
				panic("Se esperaba ')' después del nombre de variable de excepción")
			}
			p.nextToken() // consumir ")"
			catchBlock = p.parseBlockStatement()
		}
	}

	var finallyBlock *BlockStatement
	if p.curTok.Value == FINALLY {
		p.nextToken() // consumir "finally"
		finallyBlock = p.parseBlockStatement()
	}

	return &TryStatement{Body: body, CatchBlock: catchBlock, ExceptionVar: exceptionVar, FinallyBlock: finallyBlock}
}

// parseAssignmentOrExpressionStatement
func (p *Parser) parseAssignmentOrExpressionStatement() Node {
	left := p.parseExpression()
	if p.curTok.Value == "=" {
		p.nextToken()
		right := p.parseExpression()
		if p.curTok.Value == ";" {
			p.nextToken()
		}
		return &GenericAssignStatement{Left: left, Right: right}
	}

	if p.curTok.Value == "++" {
		p.nextToken()
		if p.curTok.Value == ";" {
			p.nextToken()
		}
		return &GenericAssignStatement{Left: left, Right: &BinaryExpression{Left: left, Op: "+", Right: &NumberLiteral{Value: 1}}}
	}

	if p.curTok.Value == "--" {
		p.nextToken()
		if p.curTok.Value == ";" {
			p.nextToken()
		}
		return &GenericAssignStatement{Left: left, Right: &BinaryExpression{Left: left, Op: "-", Right: &NumberLiteral{Value: 1}}}
	}

	if p.curTok.Value == ";" {
		p.nextToken()
	}
	return &ExprStatement{Expr: left}
}

func (p *Parser) parseReturnStatement() Node {
	p.nextToken() // consumir "return"
	if p.curTok.Value == ";" {
		p.nextToken()
		return &ReturnStatement{Value: nil}
	}
	expr := p.parseExpression()
	if p.curTok.Value == ";" {
		p.nextToken()
	}
	return &ReturnStatement{Value: expr}
}

// let x = expr;
func (p *Parser) parseLetStatement() Node {
	p.nextToken() // "let"
	if p.curTok.Type != TOKEN_IDENT {
		panic("Variable name expected after \"let\" or \"var\"")
	}
	name := p.curTok.Value
	p.nextToken()
	if p.curTok.Value == ";" {
		p.nextToken()
		return &LetStatement{Name: name, Value: nil}
	}

	if p.curTok.Value == IN {
		return &LetStatement{Name: name, Value: nil}
	}

	if p.curTok.Value != "=" {
		panic("Variable assignment expected after variable name")
	}
	p.nextToken()
	val := p.parseExpression()
	if p.curTok.Value == ";" {
		p.nextToken()
	}
	return &LetStatement{Name: name, Value: val}
}

// parseFunctionDeclaration => "func nombre(args) { ... }"
func (p *Parser) parseFunctionDeclaration() Node {
	p.nextToken() // consumir "func"
	if p.curTok.Type != TOKEN_IDENT {
		panic("Function name expected after \"func\"")
	}
	funcName := p.curTok.Value
	p.nextToken()
	if p.curTok.Value != "(" {
		panic("'(' expected after function name")
	}
	args := p.parseFunctionArgs()
	body := p.parseBlockStatement()
	return &FunctionDeclaration{Name: funcName, Args: args, Body: body}
}

func (p *Parser) parseIfStatement() Node {
	p.nextToken() // "if"
	if p.curTok.Value != "(" {
		panic("Expected '(' after 'if'")
	}
	p.nextToken()
	cond := p.parseExpression()
	if p.curTok.Value != ")" {
		panic("')' expected after if condition")
	}
	p.nextToken()
	consequence := p.parseBlockStatement()

	var alternative *BlockStatement
	if p.curTok.Value == "else" {
		p.nextToken()
		alternative = p.parseBlockStatement()
	}
	return &IfStatement{Condition: cond, Consequence: consequence, Alternative: alternative}
}

func (p *Parser) parseWhileStatement() Node {
	p.nextToken() // "while"
	if p.curTok.Value != "(" {
		panic("Se esperaba '(' tras 'while'")
	}
	p.nextToken()
	cond := p.parseExpression()
	if p.curTok.Value != ")" {
		panic("Se esperaba ')'")
	}
	p.nextToken()
	body := p.parseBlockStatement()
	return &WhileStatement{Condition: cond, Body: body}
}

func (p *Parser) parseForStatement() Node {
	p.nextToken() // "for"
	if p.curTok.Value != "(" {
		panic("Se esperaba '(' tras 'for'")
	}
	p.nextToken()

	var init Node
	if p.curTok.Value == LET || p.curTok.Value == VAR {
		init = p.parseLetStatement()
		init.(*LetStatement).Value = Node(&NumberLiteral{Value: 0})
		indexName := init.(*LetStatement).Name
		if p.curTok.Value == IN {
			p.nextToken()

			collName := p.curTok.Value
			exp := p.parseExpression()
			for p.curTok.Value != "{" {
				p.nextToken()
			}
			body := p.parseBlockStatement()
			return &ForStatement{Init: exp, Body: body, inFlag: true, inArray: collName, inIndexName: indexName}
		}

	} else {
		if p.curTok.Type != TOKEN_SYMBOL && !(p.curTok.Type == TOKEN_IDENT && p.peekTok.Value == "=") {
			// no hay init
		} else if p.curTok.Type == TOKEN_IDENT && p.peekTok.Value == "=" {
			init = p.parseAssignmentOrExpressionStatement()
		}
	}

	var condition Node
	if p.curTok.Value != ";" {
		condition = p.parseExpression()
	} else {
		condition = &BooleanLiteral{Value: true}
	}
	if p.curTok.Value == ";" {
		p.nextToken()
	}

	var post Node
	if p.curTok.Value != ")" {
		if p.curTok.Value == LET || p.curTok.Value == VAR {
			post = p.parseLetStatement()
		} else {
			post = p.parseAssignmentOrExpressionStatement()
		}
	}
	if p.curTok.Value != ")" {
		panic("Se esperaba ')' en 'for(...)'")
	}
	p.nextToken()
	body := p.parseBlockStatement()
	return &ForStatement{Init: init, Condition: condition, Post: post, Body: body, inFlag: false}
}

// obj MiObj { ... }
func (p *Parser) parseObjectDeclaration() Node {
	p.nextToken() // "obj"
	if p.curTok.Type != TOKEN_IDENT {
		panic("Se esperaba nombre de obj tras '" + OBJECT + "'")
	}
	objName := p.curTok.Value
	p.nextToken()
	if p.curTok.Value != "{" {
		panic("Se esperaba '{' tras nombre del " + OBJECT)
	}
	p.nextToken()

	var members []Node
	for p.curTok.Value != "}" && p.curTok.Type != TOKEN_EOF {
		if p.curTok.Value == LET || p.curTok.Value == VAR {
			members = append(members, p.parseLetStatement())
		} else if p.curTok.Value == FUNC || p.curTok.Value == FUNCTION || p.curTok.Value == METHOD {
			members = append(members, p.parseFunctionDeclaration())
		} else {
			panic("Dentro de obj => '" + LET + "' o '" + FUNC + "'")
		}
	}
	if p.curTok.Value != "}" {
		panic("Se esperaba '}' al final de " + OBJECT)
	}
	p.nextToken()
	return &ObjectDeclaration{Name: objName, Members: members}
}

// parseFunctionArgs => lee identificadores separados por coma hasta ")"
func (p *Parser) parseFunctionArgs() []string {
	var args []string
	p.nextToken() // consumir "("
	for p.curTok.Value != ")" && p.curTok.Type != TOKEN_EOF {
		if p.curTok.Type == TOKEN_IDENT {
			args = append(args, p.curTok.Value)
		} else if p.curTok.Value != "," && p.curTok.Value != ")" {
			panic("Error parseando args, token: " + p.curTok.Value)
		}
		p.nextToken()
	}
	if p.curTok.Value != ")" {
		panic("Se esperaba ')' tras argumentos de función")
	}
	p.nextToken() // consumir ")"
	return args
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	if p.curTok.Value != "{" {
		panic("Se esperaba '{' para iniciar bloque")
	}
	p.nextToken()
	var stmts []Node
	for p.curTok.Value != "}" && p.curTok.Type != TOKEN_EOF {
		stmts = append(stmts, p.parseStatement())
	}
	if p.curTok.Value != "}" {
		panic("Se esperaba '}' al cerrar bloque")
	}
	p.nextToken()
	return &BlockStatement{Statements: stmts}
}

// parseExpression => parsea binarios
func (p *Parser) parseExpression() Node {
	left := p.parseFactor()
	for p.curTok.Type == TOKEN_SYMBOL && isBinaryOp(p.curTok.Value) {
		op := p.curTok.Value
		p.nextToken()
		right := p.parseFactor()
		left = &BinaryExpression{Left: left, Op: op, Right: right}
	}
	return left
}

// parseFactor => parsea literales, ident, arrays, maps, paréntesis, O LA FUNCIÓN ANÓNIMA
func (p *Parser) parseFactor() Node {
	// ¿detectamos la anónima "func(x,y){...}"?
	if p.curTok.Type == TOKEN_IDENT && strings.ToLower(p.curTok.Value) == "func" {
		return p.parseAnonymousFunction()
	}

	switch p.curTok.Type {

	case TOKEN_NUMBER:
		val, err := strconv.ParseFloat(p.curTok.Value, 64)
		if err != nil {
			panic("No se pudo parsear número: " + p.curTok.Value)
		}
		node := &NumberLiteral{Value: val}
		p.nextToken()
		return node

	case TOKEN_STRING:
		node := &StringLiteral{Value: p.curTok.Value}
		p.nextToken()
		return node

	case TOKEN_IDENT:
		// normal ident
		idName := p.curTok.Value
		p.nextToken()
		identNode := &Identifier{Name: idName}
		return p.parsePostfix(identNode)

	case TOKEN_SYMBOL:
		if p.curTok.Value == "(" {
			// ( expr )
			p.nextToken()
			expr := p.parseExpression()
			if p.curTok.Value != ")" {
				panic("Se esperaba ')' tras ( expr )")
			}
			p.nextToken()
			return expr
		}
		if p.curTok.Value == "[" {
			return p.parseArrayLiteral()
		}
		if p.curTok.Value == "{" {
			return p.parseMapLiteral()
		}
		panic("Símbolo inesperado en factor: " + p.curTok.Value)
	}
	panic("Token inesperado en factor: " + p.curTok.Value)
}

// parseAnonymousFunction => "func(...args){...}"
func (p *Parser) parseAnonymousFunction() Node {
	// ya vimos p.curTok == "func" (type=ident)
	p.nextToken() // consumir "func"
	if p.curTok.Value != "(" {
		panic("Se esperaba '(' tras 'func' en la función anónima")
	}
	args := p.parseFunctionArgs()
	body := p.parseBlockStatement()
	return &FunctionLiteral{Args: args, Body: body}
}

func (p *Parser) parsePostfix(left Node) Node {
	for {
		if p.curTok.Value == "(" {
			left = p.parseCallExpression(left)
		} else if p.curTok.Value == "." {
			left = p.parseAccessExpression(left)
		} else if p.curTok.Value == "[" {
			left = p.parseIndexExpression(left)
		} else {
			break
		}
	}
	return left
}

func (p *Parser) parseCallExpression(left Node) Node {
	p.nextToken() // consumir "("
	var args []Node
	for p.curTok.Value != ")" && p.curTok.Type != TOKEN_EOF {
		args = append(args, p.parseExpression())
		if p.curTok.Value == "," {
			p.nextToken()
		}
	}
	if p.curTok.Value != ")" {
		panic("Se esperaba ')' al final de llamada a función")
	}
	p.nextToken() // ")"
	return &CallExpression{Callee: left, Args: args}
}

func (p *Parser) parseAccessExpression(left Node) Node {
	p.nextToken() // "."
	if p.curTok.Type != TOKEN_IDENT {
		panic("Se esperaba identificador tras '.'")
	}
	mem := p.curTok.Value
	p.nextToken()
	node := &AccessExpression{Object: left, Member: mem}
	return p.parsePostfix(node)
}

func (p *Parser) parseIndexExpression(left Node) Node {
	p.nextToken() // "["
	idx := p.parseExpression()
	if p.curTok.Value != "]" {
		panic("Se esperaba ']' en index")
	}
	p.nextToken()
	ie := &IndexExpression{Left: left, Index: idx}
	return p.parsePostfix(ie)
}

// parseArrayLiteral => [ expr, expr, ...]
func (p *Parser) parseArrayLiteral() Node {
	p.nextToken() // "["
	var elems []Node
	for p.curTok.Value != "]" && p.curTok.Type != TOKEN_EOF {
		e := p.parseExpression()
		elems = append(elems, e)
		if p.curTok.Value == "," {
			p.nextToken()
		} else if p.curTok.Value == "]" {
			break
		} else {
			panic("Se esperaba ',' o ']' en array literal")
		}
	}
	if p.curTok.Value != "]" {
		panic("Se esperaba ']' al final de array literal")
	}
	p.nextToken()
	return &ArrayLiteral{Elements: elems}
}

func (p *Parser) parseMapLiteral() Node {
	p.nextToken() // "{"
	pairs := make(map[string]Node)
	if p.curTok.Value == "}" {
		p.nextToken()
		return &MapLiteral{Pairs: pairs}
	}
	for p.curTok.Value != "}" && p.curTok.Type != TOKEN_EOF {
		var key string
		if p.curTok.Type == TOKEN_STRING {
			key = p.curTok.Value
			p.nextToken()
		} else if p.curTok.Type == TOKEN_IDENT {
			key = p.curTok.Value
			p.nextToken()
		} else {
			panic("Se esperaba clave string o identificador en map-literal")
		}

		if p.curTok.Value != ":" {
			panic("Se esperaba ':' tras la clave en map-literal")
		}
		p.nextToken()
		valNode := p.parseExpression()
		pairs[key] = valNode

		if p.curTok.Value == "," {
			p.nextToken()
		} else if p.curTok.Value == "}" {
			break
		} else {
			panic("Se esperaba ',' o '}' en map-literal")
		}
	}
	if p.curTok.Value != "}" {
		panic("Se esperaba '}' al final de map-literal")
	}
	p.nextToken()
	return &MapLiteral{Pairs: pairs}
}

func RunCode(filename string) {

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error al leer el archivo %s: %v\n", filename, err)
		os.Exit(1)
	}
	code := string(data)

	env := NewEnvironment()
	env.Set("true", true)
	env.Set("false", false)
	env.Set("nil", nil)
	env.Set("null", nil)
	env.dir = filepath.Dir(filename)

	// Registrar otras librerías si las tienes:
	RegisterLib(env)
	RegisterStd(env)
	RegisterIO(env)
	RegisterHTTPClient(env)
	RegisterString(env)
	RegisterMath(env)
	RegisterRand(env)
	RegisterTest(env)
	RegisterHTTP(env)
	RegisterPrint(env)
	RegisterOS(env)
	RegisterHack(env)
	RegisterConcurrency(env)
	RegisterCollections(env)
	parser := NewParser(code)
	env.Run(parser)
}
