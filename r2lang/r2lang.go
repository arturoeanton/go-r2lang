package r2lang

import (
	"fmt"
	"os"
	"strconv"
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
)

var (
	wg sync.WaitGroup
)

type Token struct {
	Type  string
	Value string
}

// ============================================================
// 2) LEXER (reconoce signo, decimales, comentarios, etc.)
// ============================================================

type Lexer struct {
	input  string
	pos    int
	length int
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
		ch == '_'
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
	return Token{Type: TOKEN_NUMBER, Value: val}
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
		return Token{Type: TOKEN_EOF, Value: ""}
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

	// Símbolos de 1 caracter
	singleCharSymbols := []string{
		"(", ")", "{", "}", "[", "]", ";", ",", "+", "-", "*", "/", ".", ":",
	}
	for _, s := range singleCharSymbols {
		if string(ch) == s {
			l.pos++
			return Token{Type: TOKEN_SYMBOL, Value: s}
		}
	}

	if string(ch) == "=" {
		if l.pos+1 < l.length && l.input[l.pos+1] == '=' {
			l.pos += 2
			return Token{Type: TOKEN_SYMBOL, Value: "=="}
		}
		l.pos++
		return Token{Type: TOKEN_SYMBOL, Value: "="}
	}

	// Operadores relacionales
	if ch == '<' || ch == '>' || ch == '!' || ch == '=' {
		if l.pos+1 < l.length {
			nextCh := l.input[l.pos+1]
			if nextCh == '=' {
				op := string(ch) + string(nextCh)
				l.pos += 2
				return Token{Type: TOKEN_SYMBOL, Value: op}
			}
		}
		l.pos++
		return Token{Type: TOKEN_SYMBOL, Value: string(ch)}
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
		return Token{Type: TOKEN_STRING, Value: val}
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
		return Token{Type: TOKEN_IDENT, Value: l.input[start:l.pos]}
	}

	panic(fmt.Sprintf("Caracter inesperado en lexer: %c", ch))
}

// ============================================================
// 3) AST - Node interface
// ============================================================

type Node interface {
	Eval(env *Environment) interface{}
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

type ForStatement struct {
	Init      Node
	Condition Node
	Post      Node
	Body      *BlockStatement
}

func (fs *ForStatement) Eval(env *Environment) interface{} {
	newEnv := NewInnerEnv(env)
	var result interface{}
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
		return instantiateObject(env, cv)
	default:
		panic("Intento de llamar algo que no es función ni blueprint")
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
}

func (uf *UserFunction) Call(args ...interface{}) interface{} {
	newEnv := NewInnerEnv(uf.Env)
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

type BuiltinFunction func(args ...interface{}) interface{}

type ObjectInstance struct {
	Env *Environment
}

func instantiateObject(env *Environment, blueprint map[string]interface{}) *ObjectInstance {
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
	store map[string]interface{}
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]interface{}),
		outer: nil,
	}
}

func NewInnerEnv(outer *Environment) *Environment {
	return &Environment{
		store: make(map[string]interface{}),
		outer: outer,
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
func (e *Environment) Run(parser *Parser) {
	defer wg.Wait()
	wg = sync.WaitGroup{}

	ast := parser.ParseProgram()
	// Ejecutar
	ast.Eval(e)

	// Llamar a main() si está
	mainVal, ok := e.Get("main")
	if !ok {
		fmt.Println("Aviso: No existe función main().")
		os.Exit(0)
	}
	mainFn, isFn := mainVal.(*UserFunction)
	if !isFn {
		fmt.Println("Error: 'main' no es una función.")
		os.Exit(1)
	}
	mainFn.Call()
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
}

func NewParser(input string) *Parser {
	p := &Parser{lexer: NewLexer(input)}
	p.nextToken()
	p.nextToken()
	return p
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

func (p *Parser) parseStatement() Node {
	if p.curTok.Value == "return" {
		return p.parseReturnStatement()
	}
	if p.curTok.Value == "let" {
		return p.parseLetStatement()
	}
	if p.curTok.Value == "func" {
		// esto parsea "func nombre(...) { ... }" => FunctionDeclaration con nombre
		return p.parseFunctionDeclaration()
	}
	if p.curTok.Value == "if" {
		return p.parseIfStatement()
	}
	if p.curTok.Value == "while" {
		return p.parseWhileStatement()
	}
	if p.curTok.Value == "for" {
		return p.parseForStatement()
	}
	if p.curTok.Value == "obj" {
		return p.parseObjectDeclaration()
	}
	// sino parseAsignmentOrExpressionStatement
	return p.parseAssignmentOrExpressionStatement()
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
		panic("Se esperaba nombre de variable tras 'let'")
	}
	name := p.curTok.Value
	p.nextToken()
	if p.curTok.Value == ";" {
		p.nextToken()
		return &LetStatement{Name: name, Value: nil}
	}
	if p.curTok.Value != "=" {
		panic("Se esperaba '=' tras 'let var'")
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
		panic("Se esperaba nombre de función tras 'func'")
	}
	funcName := p.curTok.Value
	p.nextToken()
	if p.curTok.Value != "(" {
		panic("Se esperaba '(' tras el nombre de la función")
	}
	args := p.parseFunctionArgs()
	body := p.parseBlockStatement()
	return &FunctionDeclaration{Name: funcName, Args: args, Body: body}
}

func (p *Parser) parseIfStatement() Node {
	p.nextToken() // "if"
	if p.curTok.Value != "(" {
		panic("Se esperaba '(' tras 'if'")
	}
	p.nextToken()
	cond := p.parseExpression()
	if p.curTok.Value != ")" {
		panic("Se esperaba ')' tras condición if")
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
	if p.curTok.Value == "let" {
		init = p.parseLetStatement()
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
		if p.curTok.Value == "let" {
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
	return &ForStatement{Init: init, Condition: condition, Post: post, Body: body}
}

// obj MiObj { ... }
func (p *Parser) parseObjectDeclaration() Node {
	p.nextToken() // "obj"
	if p.curTok.Type != TOKEN_IDENT {
		panic("Se esperaba nombre de obj tras 'obj'")
	}
	objName := p.curTok.Value
	p.nextToken()
	if p.curTok.Value != "{" {
		panic("Se esperaba '{' tras nombre del obj")
	}
	p.nextToken()

	var members []Node
	for p.curTok.Value != "}" && p.curTok.Type != TOKEN_EOF {
		if p.curTok.Value == "let" {
			members = append(members, p.parseLetStatement())
		} else if p.curTok.Value == "func" {
			members = append(members, p.parseFunctionDeclaration())
		} else {
			panic("Dentro de obj => 'let' o 'func'")
		}
	}
	if p.curTok.Value != "}" {
		panic("Se esperaba '}' al final de obj")
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
	if p.curTok.Type == TOKEN_IDENT && p.curTok.Value == "func" {
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

func RunCode(input string) {
	env := NewEnvironment()
	env.Set("true", true)
	env.Set("false", false)
	env.Set("nil", nil)
	env.Set("null", nil)
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
	parser := NewParser(input)
	env.Run(parser)
}
