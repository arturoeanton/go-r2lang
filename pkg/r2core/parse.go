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

// let x = expr;
func (p *Parser) parseLetStatement() Node {
	p.nextToken() // "let"
	if p.curTok.Type != TOKEN_IDENT {
		p.except("Variable name expected after 'let'/'var'")
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
		p.except("Variable assignment expected after variable name")
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
		p.except("‘)’ was expected after the post statement in ‘for’")
	}
	p.nextToken()
	body := p.parseBlockStatement()
	return &ForStatement{Init: init, Condition: condition, Post: post, Body: body, inFlag: false}
}

// obj MiObj { ... }
func (p *Parser) parseObjectDeclaration() Node {
	p.nextToken() // "obj"
	if p.curTok.Type != TOKEN_IDENT {
		p.except("Object name was expected after '" + OBJECT + "'")
	}
	objName := p.curTok.Value
	p.nextToken()
	if p.curTok.Value != "{" && p.curTok.Value != EXTENDS {
		p.except("Expected ‘{’ or ‘extends’ after object name")
	}
	parentName := ""
	if p.curTok.Value == EXTENDS {
		p.nextToken()
		if p.curTok.Type != TOKEN_IDENT {
			p.except("Expected object name after ‘extends’")
		}
		parentName = p.curTok.Value
		p.nextToken()
		if p.curTok.Value != "{" {
			p.except("Expected ‘{’ after object name")
		}
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
			p.except("Expected string or identifier as key in map-literal")
		}

		if p.curTok.Value != ":" {
			p.except("Expected ':' after key in map-literal")
		}
		p.nextToken()
		valNode := p.parseExpression()
		pairs[key] = valNode

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

func (p *Parser) except(msgErr string) {

	msg := fmt.Sprintln("Parser Exception: Line:", p.curTok.Line, ":", p.curTok.Col, "Error:", msgErr)
	_, err := fmt.Fprintf(os.Stderr, msg)
	if err != nil {
		panic(msg)
	}
	os.Exit(1)
	panic(msg)

}
