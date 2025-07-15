package r2core

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Parser struct {
	lexer *Lexer
	//savTok  Token
	prevTok Token
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
		p.except("A string was expected after ‘import’")
	}

	path := p.curTok.Value
	p.nextToken()

	var alias string
	if p.curTok.Type == TOKEN_AS {
		p.nextToken() // Consumir 'as'
		if p.curTok.Type != TOKEN_IDENT {
			p.except("An identifier was expected after ‘as’")
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
		p.except("A string was expected for the test case name")
	}
	name := p.curTok.Value
	p.nextToken()

	if p.curTok.Value != "{" {
		p.except("‘{’ was expected to start the test case body")
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
			p.except("‘Given’, ‘When’, ‘Then’, or ‘And’ was expected in the test case steps")
		}
		command := p.parseExpression()
		steps = append(steps, TestStep{Type: stepType, Command: command})
		if p.curTok.Value == ";" {
			p.nextToken()
		}
	}
	if p.curTok.Value != "}" {
		p.except("‘}’ was expected at the end of the test case")
	}
	p.nextToken()
	return &TestCase{Name: name, Steps: steps}
}

func (p *Parser) nextToken() {
	p.prevTok = p.curTok
	p.curTok = p.peekTok
	p.peekTok = p.lexer.NextToken()

}

func (p *Parser) ParseProgram() *Program {
	prog := &Program{}
	for p.curTok.Type != TOKEN_EOF {
		if p.curTok.Type == TOKEN_SYMBOL && p.curTok.Value == "\n" {
			p.nextToken()
			continue
		}
		stmt := p.parseStatement()
		prog.Statements = append(prog.Statements, stmt)
	}
	return prog
}

func (p *Parser) parseThrowStatement() Node {
	p.nextToken()
	if p.curTok.Type != TOKEN_STRING {
		p.except("A string was expected after ‘throw’")
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
	if p.curTok.Value == BREAK {
		return p.parseBreakStatement()
	}
	if p.curTok.Value == CONTINUE {
		return p.parseContinueStatement()
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

	if p.curTok.Value == CLASS {
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
				p.except("‘(’ was expected after ‘catch’")
			}
			p.nextToken() // consumir "("
			if p.curTok.Type != TOKEN_IDENT {
				p.except("Variable name expected after ‘catch’")
			}
			exceptionVar = p.curTok.Value
			p.nextToken()
			if p.curTok.Value != ")" {
				p.except("‘)’ was expected after the exception variable")
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

func (p *Parser) parseBreakStatement() Node {
	p.nextToken() // consumir "break"
	if p.curTok.Value == ";" {
		p.nextToken()
	}
	return &BreakStatement{}
}

func (p *Parser) parseContinueStatement() Node {
	p.nextToken() // consumir "continue"
	if p.curTok.Value == ";" {
		p.nextToken()
	}
	return &ContinueStatement{}
}

// let x = expr;
func (p *Parser) parseLetStatement() Node {
	p.nextToken() // "let"

	var declarations []LetDeclaration

	// Parsear primera declaración
	if p.curTok.Type != TOKEN_IDENT {
		p.except("Variable name expected after 'let'/'var'")
	}

	name := p.curTok.Value
	p.nextToken()

	// Manejar caso especial para bucles for-in
	if p.curTok.Value == IN {
		return &LetStatement{Name: name, Value: nil}
	}

	var value Node
	if p.curTok.Value == "=" {
		p.nextToken()
		value = p.parseExpression()
	}

	declarations = append(declarations, LetDeclaration{Name: name, Value: value})

	// Parsear declaraciones adicionales separadas por comas
	for p.curTok.Value == "," {
		p.nextToken() // consumir ","

		if p.curTok.Type != TOKEN_IDENT {
			p.except("Variable name expected after ','")
		}

		name = p.curTok.Value
		p.nextToken()

		var value Node
		if p.curTok.Value == "=" {
			p.nextToken()
			value = p.parseExpression()
		}

		declarations = append(declarations, LetDeclaration{Name: name, Value: value})
	}

	// Consumir punto y coma opcional
	if p.curTok.Value == ";" {
		p.nextToken()
	}

	// Si solo hay una declaración, usar LetStatement simple para mantener compatibilidad
	if len(declarations) == 1 {
		return &LetStatement{Name: declarations[0].Name, Value: declarations[0].Value}
	}

	// Si hay múltiples declaraciones, usar MultipleLetStatement
	return &MultipleLetStatement{Declarations: declarations}
}

// parseFunctionDeclaration => "func nombre(args) { ... }"
func (p *Parser) parseFunctionDeclaration() Node {
	p.nextToken() // consumir "func"
	return p.parseFunctionDeclaratioWithoutFunc()
}

func (p *Parser) parseFunctionDeclaratioWithoutFunc() Node {
	if p.curTok.Type != TOKEN_IDENT {
		p.except("Function name expected after 'func'/'function'")
	}
	funcName := p.curTok.Value
	p.nextToken()
	if p.curTok.Value != "(" {
		p.except("'(' expected after function name")
	}
	args := p.parseFunctionArgs()
	body := p.parseBlockStatement()
	return &FunctionDeclaration{Name: funcName, Args: args, Body: body}
}

func (p *Parser) parseIfStatement() Node {
	p.nextToken() // "if"
	if p.curTok.Value != "(" {
		p.except("Expected '(' after 'if'")
	}
	p.nextToken()
	cond := p.parseExpression()
	if p.curTok.Value != ")" {
		p.except("')' expected after if condition")
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
		p.except("‘(’ was expected after ‘while’")
	}
	p.nextToken()
	cond := p.parseExpression()
	if p.curTok.Value != ")" {
		p.except("‘)’ was expected after the condition in ‘while’")
	}
	p.nextToken()
	body := p.parseBlockStatement()
	return &WhileStatement{Condition: cond, Body: body}
}

func (p *Parser) parseForStatement() Node {
	p.nextToken() // "for"
	if p.curTok.Value != "(" {
		p.except("‘(’ was expected after ‘for’")
	}
	p.nextToken()

	// Check for for-in loop
	if p.peekTok.Value == IN {
		return p.parseForInStatement()
	}

	// Standard for loop
	return p.parseStandardForStatement()
}

func (p *Parser) parseForInStatement() Node {
	indexName := p.curTok.Value
	p.nextToken() // consume index name
	p.nextToken() // consume 'in'

	collName := p.curTok.Value
	exp := p.parseExpression()
	for p.curTok.Value != "{" {
		p.nextToken()
	}
	body := p.parseBlockStatement()
	return &ForStatement{Init: exp, Body: body, inFlag: true, inArray: collName, inIndexName: indexName}
}

func (p *Parser) parseStandardForStatement() Node {
	var init Node
	if p.curTok.Type == TOKEN_IDENT && p.peekTok.Value == "=" {
		init = p.parseAssignmentOrExpressionStatement()
	} else if p.curTok.Value == LET || p.curTok.Value == VAR {
		init = p.parseLetStatement()
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
		p.except("‘)’ was expected after the post statement in ‘for’")
	}
	p.nextToken()
	body := p.parseBlockStatement()
	return &ForStatement{Init: init, Condition: condition, Post: post, Body: body, inFlag: false}
}

func (p *Parser) parseObjectDeclaration() Node {
	p.nextToken() // "obj"
	if p.curTok.Type != TOKEN_IDENT {
		p.except("Object name was expected after '" + OBJECT + "'")
	}
	objName := p.curTok.Value
	p.nextToken()

	parentName := p.parseOptionalExtends()

	if p.curTok.Value != "{" {
		p.except("Expected ‘{’ after object name")
	}
	p.nextToken()

	var members []Node
	for p.curTok.Value != "}" && p.curTok.Type != TOKEN_EOF {
		if p.curTok.Type == TOKEN_SYMBOL && p.curTok.Value == "\n" {
			p.nextToken()
			continue
		}
		if p.curTok.Value == LET || p.curTok.Value == VAR {
			members = append(members, p.parseLetStatement())
		} else if p.curTok.Value == FUNC || p.curTok.Value == FUNCTION || p.curTok.Value == METHOD {
			members = append(members, p.parseFunctionDeclaration())
		} else if p.curTok.Type == TOKEN_IDENT {
			members = append(members, p.parseFunctionDeclaratioWithoutFunc())
		} else {
			p.except("Inside " + OBJECT + " only 'let', 'var', 'func', 'function' or 'method' are allowed")
		}
	}
	if p.curTok.Value != "}" {
		p.except("Expected ‘}’ at the end of " + OBJECT)
	}
	p.nextToken()
	return &ObjectDeclaration{Name: objName, Members: members, ParentName: parentName}
}

func (p *Parser) parseOptionalExtends() string {
	if p.curTok.Value != EXTENDS {
		return ""
	}
	p.nextToken() // consume 'extends'
	if p.curTok.Type != TOKEN_IDENT {
		p.except("Expected object name after ‘extends’")
	}
	parentName := p.curTok.Value
	p.nextToken()
	return parentName
}

// parseFunctionArgs => lee identificadores separados por coma hasta ")"
func (p *Parser) parseFunctionArgs() []string {
	var args []string
	p.nextToken() // consumir "("
	for p.curTok.Value != ")" && p.curTok.Type != TOKEN_EOF {
		if p.curTok.Type == TOKEN_IDENT {
			args = append(args, p.curTok.Value)
		} else if p.curTok.Value != "," && p.curTok.Value != ")" {
			p.except("Error parsing args, token:  " + p.curTok.Value)
		}
		p.nextToken()
	}
	if p.curTok.Value != ")" {
		p.except("Expected ‘)’ after function arguments")
	}
	p.nextToken() // consumir ")"
	return args
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	if p.curTok.Value != "{" {
		p.except("Expected ‘{’ to start block")
	}
	p.nextToken()
	var stmts []Node
	for p.curTok.Value != "}" && p.curTok.Type != TOKEN_EOF {
		if p.curTok.Type == TOKEN_SYMBOL && p.curTok.Value == "\n" {
			p.nextToken()
			continue
		}
		stmts = append(stmts, p.parseStatement())
	}
	if p.curTok.Value != "}" {
		p.except("Expected ‘}’ to end block")
	}
	p.nextToken()
	return &BlockStatement{Statements: stmts}
}

// parseExpression => parsea ternarios y binarios
func (p *Parser) parseExpression() Node {
	left := p.parseBinaryExpression(1)

	// Operador ternario tiene la precedencia más baja
	if p.curTok.Type == TOKEN_SYMBOL && p.curTok.Value == "?" {
		p.nextToken() // consumir "?"
		trueExpr := p.parseExpression()
		if p.curTok.Type != TOKEN_SYMBOL || p.curTok.Value != ":" {
			p.except("Expected ':' in ternary expression")
		}
		p.nextToken() // consumir ":"
		falseExpr := p.parseExpression()
		return &TernaryExpression{Condition: left, TrueExpr: trueExpr, FalseExpr: falseExpr}
	}

	return left
}

// parseBinaryExpression => parsea operadores binarios
func (p *Parser) parseBinaryExpression(precedence int) Node {
	left := p.parseFactor()
	for p.curTok.Type == TOKEN_SYMBOL && isBinaryOp(p.curTok.Value) && getPrecedence(p.curTok.Value) >= precedence {
		op := p.curTok.Value
		p.nextToken()
		right := p.parseBinaryExpression(getPrecedence(op) + 1)
		left = &BinaryExpression{Left: left, Op: op, Right: right}
	}
	return left
}

func getPrecedence(op string) int {
	switch op {
	case "||":
		return 1
	case "&&":
		return 2
	case "==", "!=", "<", ">", "<=", ">=":
		return 3
	case "+", "-":
		return 4
	case "*", "/":
		return 5
	default:
		return 0
	}
}

// parseFactor => parsea literales, ident, arrays, maps, paréntesis, O LA FUNCIÓN ANÓNIMA
func (p *Parser) parseFactor() Node {
	// ¿detectamos la anónima "func(x,y){...}"?
	if p.curTok.Type == TOKEN_IDENT && (strings.ToLower(p.curTok.Value) == FUNC || strings.ToLower(p.curTok.Value) == FUNCTION) {
		return p.parseAnonymousFunction()
	}

	switch p.curTok.Type {

	case TOKEN_NUMBER:
		val, err := strconv.ParseFloat(p.curTok.Value, 64)
		if err != nil {
			p.except("Could not parse number: " + p.curTok.Value)
		}
		node := &NumberLiteral{Value: val}
		p.nextToken()
		return node

	case TOKEN_STRING:
		node := &StringLiteral{Value: p.curTok.Value}
		p.nextToken()
		return node

	case TOKEN_TEMPLATE_STRING:
		return p.parseTemplateString()

	case TOKEN_DATE:
		dateValue, err := ParseDateLiteral(p.curTok.Value)
		if err != nil {
			p.except("Invalid date literal: " + p.curTok.Value + " - " + err.Error())
		}
		node := &DateLiteral{Value: dateValue}
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
				p.except("Expected ‘)’ after ( expr )")
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
		p.except("Unexpected symbol in factor: " + p.curTok.Value)
	}
	p.except("Unexpected token in factor: " + p.curTok.Value)
	return nil
}

// parseAnonymousFunction => "func(...args){...}"
func (p *Parser) parseAnonymousFunction() Node {
	// ya vimos p.curTok == "func" (type=ident)
	p.nextToken() // consumir "func"
	if p.curTok.Value != "(" {
		p.except("Expected ‘(’ after ‘func’ in the anonymous function")
	}
	args := p.parseFunctionArgs()
	body := p.parseBlockStatement()
	return &FunctionLiteral{Args: args, Body: body}
}

func (p *Parser) parsePostfix(left Node) Node {
	for {
		switch p.curTok.Value {
		case "(":
			left = p.parseCallExpression(left)
		case ".":
			left = p.parseAccessExpression(left)
		case "[":
			left = p.parseIndexExpression(left)
		default:
			return left
		}
	}
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
		p.except("Expected ‘)’ at the end of function call")
	}
	p.nextToken() // ")"
	return &CallExpression{Callee: left, Args: args}
}

func (p *Parser) parseAccessExpression(left Node) Node {
	p.nextToken() // "."
	if p.curTok.Type != TOKEN_IDENT {
		p.except("Expected identifier after ‘.’")
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
		p.except("Expected ‘]’ at the end of index expression")
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
			p.except("Expected ',' or ']' in array literal")
		}
	}
	if p.curTok.Value != "]" {
		p.except("Expected ']' at the end of array literal")
	}
	p.nextToken()
	return &ArrayLiteral{Elements: elems}
}

func (p *Parser) parseMapLiteral() Node {
	p.nextToken() // "{"
	var pairs []MapPair
	if p.curTok.Value == "}" {
		p.nextToken()
		return &MapLiteral{Pairs: pairs}
	}
	for p.curTok.Value != "}" && p.curTok.Type != TOKEN_EOF {
		var keyNode Node

		// Soportar expresiones como claves
		switch p.curTok.Type {
		case TOKEN_STRING:
			keyNode = &StringLiteral{Value: p.curTok.Value}
			p.nextToken()
		case TOKEN_IDENT:
			keyNode = &StringLiteral{Value: p.curTok.Value}
			p.nextToken()
		case TOKEN_SYMBOL:
			if p.curTok.Value == "(" {
				// Permitir expresiones entre paréntesis como claves
				p.nextToken() // consumir "("
				keyNode = p.parseExpression()
				if p.curTok.Value != ")" {
					p.except("Expected ')' after key expression")
				}
				p.nextToken() // consumir ")"
			} else {
				// Permitir expresiones simples como claves
				keyNode = p.parseExpression()
			}
		default:
			// Permitir expresiones simples como claves
			keyNode = p.parseExpression()
		}

		if p.curTok.Value != ":" {
			p.except("Expected ':' after key in map-literal")
		}
		p.nextToken()
		valNode := p.parseExpression()

		pairs = append(pairs, MapPair{Key: keyNode, Value: valNode})

		if p.curTok.Value == "," {
			p.nextToken()
		} else if p.curTok.Value == "}" {
			break
		} else {
			p.except("Expected ',' or '}' in map-literal")
		}
	}
	if p.curTok.Value != "}" {
		p.except("Expected '}' at the end of map-literal")
	}
	p.nextToken()
	return &MapLiteral{Pairs: pairs}
}

// parseTemplateString parses a template string token into a TemplateString AST node
func (p *Parser) parseTemplateString() Node {
	encoded := p.curTok.Value
	p.nextToken() // consume template string token

	parts := parseTemplateParts(encoded, p)
	return &TemplateString{Parts: parts}
}

func (p *Parser) except(msgErr string) {

	msg := fmt.Sprintln("Parser Exception: Line:", p.curTok.Line, ":", p.curTok.Col, "Error:", msgErr)
	_, err := fmt.Fprint(os.Stderr, msg)
	if err != nil {
		panic(msg)
	}
	os.Exit(1)
	panic(msg)

}
