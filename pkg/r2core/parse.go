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
	prevTok  Token
	curTok   Token
	peekTok  Token
	baseDir  string // Directorio base para importaciones
	filename string // Archivo actual para tracking de posición
}

func NewParser(input string) *Parser {
	return NewParserWithFile(input, "")
}

func NewParserWithFile(input string, filename string) *Parser {
	p := &Parser{
		lexer:    NewLexer(input),
		filename: filename,
	}
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
		// Verificar si es destructuring
		if p.peekTok.Value == "[" {
			return p.parseArrayDestructuring()
		}
		if p.peekTok.Value == "{" {
			return p.parseObjectDestructuring()
		}
		return p.parseLetStatement()
	}
	if p.curTok.Value == CONST {
		return p.parseConstStatement()
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

	if p.curTok.Value == DSL {
		return p.parseDSLDefinition()
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

	// Handle compound assignment operators
	if p.curTok.Value == "+=" || p.curTok.Value == "-=" || p.curTok.Value == "*=" || p.curTok.Value == "/=" {
		op := p.curTok.Value[:len(p.curTok.Value)-1] // Remove the '=' to get the binary operator
		p.nextToken()
		right := p.parseExpression()
		if p.curTok.Value == ";" {
			p.nextToken()
		}
		return &GenericAssignStatement{Left: left, Right: &BinaryExpression{Left: left, Op: op, Right: right}}
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
	if p.curTok.Value == IN || p.curTok.Value == "in" {
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

func (p *Parser) parseConstStatement() Node {
	p.nextToken() // "const"

	var declarations []ConstDeclaration

	// Parsear primera declaración
	if p.curTok.Type != TOKEN_IDENT {
		p.except("Variable name expected after 'const'")
	}

	name := p.curTok.Value
	p.nextToken()

	// const requires initialization
	if p.curTok.Value != "=" {
		p.except("const declarations must be initialized")
	}

	p.nextToken()
	value := p.parseExpression()

	declarations = append(declarations, ConstDeclaration{Name: name, Value: value})

	// Parsear declaraciones adicionales separadas por comas
	for p.curTok.Value == "," {
		p.nextToken() // consumir ","

		if p.curTok.Type != TOKEN_IDENT {
			p.except("Variable name expected after ','")
		}

		name = p.curTok.Value
		p.nextToken()

		// const requires initialization
		if p.curTok.Value != "=" {
			p.except("const declarations must be initialized")
		}

		p.nextToken()
		value = p.parseExpression()

		declarations = append(declarations, ConstDeclaration{Name: name, Value: value})
	}

	// Consumir punto y coma opcional
	if p.curTok.Value == ";" {
		p.nextToken()
	}

	// Si solo hay una declaración, usar ConstStatement simple para mantener compatibilidad
	if len(declarations) == 1 {
		return &ConstStatement{Name: declarations[0].Name, Value: declarations[0].Value}
	}

	// Si hay múltiples declaraciones, usar MultipleConstStatement
	return &MultipleConstStatement{Declarations: declarations}
}

// parseFunctionDeclaration => "func nombre(args) { ... }"
func (p *Parser) parseFunctionDeclaration() Node {
	funcToken := p.curTok
	p.nextToken() // consumir "func"
	return p.parseFunctionDeclaratioWithoutFunc(funcToken)
}

func (p *Parser) parseFunctionDeclaratioWithoutFunc(funcToken Token) Node {
	if p.curTok.Type != TOKEN_IDENT {
		p.except("Function name expected after 'func'/'function'")
	}
	funcName := p.curTok.Value
	p.nextToken()
	if p.curTok.Value != "(" {
		p.except("'(' expected after function name")
	}
	params := p.parseFunctionParameters()
	body := p.parseBlockStatement()

	// Convert parameters to args for backward compatibility
	var args []string
	for _, param := range params {
		args = append(args, param.Name)
	}

	return &FunctionDeclaration{
		BaseNode: BaseNode{
			Position: CreatePositionInfo(funcToken, p.filename),
		},
		Name: funcName, Args: args, Params: params, Body: body}
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
		// Check for "else if"
		if p.curTok.Value == "if" {
			// Parse "else if" as nested if statement
			elseIfNode := p.parseIfStatement()
			// Wrap the else-if in a block statement
			alternative = &BlockStatement{Statements: []Node{elseIfNode}}
		} else {
			alternative = p.parseBlockStatement()
		}
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
	if p.peekTok.Value == IN || p.peekTok.Value == "in" {
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
	p.nextToken() // consume collection name

	// Skip to ')'
	if p.curTok.Value != ")" {
		p.except("')' expected after for-in expression")
	}
	p.nextToken() // consume ')'

	// Parse body
	body := p.parseBlockStatement()
	// Create a dummy init that sets the index variable
	init := &LetStatement{Name: indexName, Value: &NumberLiteral{Value: 0}}
	return &ForStatement{Init: init, Body: body, inFlag: true, inArray: collName, inIndexName: indexName}
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
			methodToken := p.curTok
			members = append(members, p.parseFunctionDeclaratioWithoutFunc(methodToken))
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

// parseFunctionArgs => lee identificadores separados por coma hasta ")" (backward compatibility)
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
		p.except("Expected ')' after function arguments")
	}
	p.nextToken() // consumir ")"
	return args
}

// parseFunctionParameters => parsea parámetros con valores por defecto
func (p *Parser) parseFunctionParameters() []Parameter {
	var params []Parameter
	p.nextToken() // consumir "("

	for p.curTok.Value != ")" && p.curTok.Type != TOKEN_EOF {
		if p.curTok.Type == TOKEN_IDENT {
			paramName := p.curTok.Value
			p.nextToken()

			var defaultValue Node
			if p.curTok.Value == "=" {
				p.nextToken() // consumir "="
				defaultValue = p.parseExpression()
			}

			params = append(params, Parameter{
				Name:         paramName,
				DefaultValue: defaultValue,
			})

			if p.curTok.Value == "," {
				p.nextToken() // consumir ","
			} else if p.curTok.Value != ")" {
				p.except("Expected ',' or ')' after parameter")
			}
		} else if p.curTok.Value == ")" {
			// Empty parameter list - break the loop
			break
		} else {
			p.except("Expected parameter name, got: " + p.curTok.Value)
		}
	}

	if p.curTok.Value != ")" {
		p.except("Expected ')' after function parameters")
	}
	p.nextToken() // consumir ")"
	return params
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

	// Operador ternario tiene la precedencia más baja (pero solo si no es ??)
	if p.curTok.Type == TOKEN_SYMBOL && p.curTok.Value == "?" && p.curTok.Type != TOKEN_NULL_COALESCING {
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
	left := p.parseUnaryExpression()
	for (p.curTok.Type == TOKEN_SYMBOL || p.curTok.Type == TOKEN_NULL_COALESCING || p.curTok.Type == TOKEN_PIPE) && isBinaryOp(p.curTok.Value) && getPrecedence(p.curTok.Value) >= precedence {
		op := p.curTok.Value
		opToken := p.curTok // Capture operator token position
		p.nextToken()
		right := p.parseBinaryExpression(getPrecedence(op) + 1)
		left = &BinaryExpression{
			BaseNode: BaseNode{Position: CreatePositionInfo(opToken, p.filename)},
			Left:     left,
			Op:       op,
			Right:    right,
		}
	}
	return left
}

// parseUnaryExpression => parsea operadores unarios como !, -, +
func (p *Parser) parseUnaryExpression() Node {
	// Operador spread ...
	if p.curTok.Type == TOKEN_ELLIPSIS {
		p.nextToken()
		value := p.parseUnaryExpression()
		return &SpreadExpression{Value: value}
	}

	if p.curTok.Type == TOKEN_SYMBOL {
		switch p.curTok.Value {
		case "!", "-", "+", "~":
			operator := p.curTok.Value
			p.nextToken()
			right := p.parseUnaryExpression()
			return &UnaryExpression{Operator: operator, Right: right}
		}
	}
	return p.parseFactor()
}

func getPrecedence(op string) int {
	switch op {
	case "|>": // Pipeline operator (P4) - very low precedence
		return 1
	case "||":
		return 2
	case "??": // Null coalescing (P3) - between || and &&
		return 3
	case "&&":
		return 4
	case "|":
		return 5
	case "^":
		return 6
	case "&":
		return 7
	case "==", "!=", "<", ">", "<=", ">=":
		return 8
	case "<<", ">>":
		return 9
	case "+", "-":
		return 10
	case "*", "/", "%":
		return 11
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

	// Match expression (P3)
	if p.curTok.Type == TOKEN_MATCH {
		return p.parseMatchExpression()
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

	case TOKEN_TRUE:
		node := &BooleanLiteral{Value: true}
		p.nextToken()
		return node

	case TOKEN_FALSE:
		node := &BooleanLiteral{Value: false}
		p.nextToken()
		return node

	case TOKEN_NIL:
		node := &NilLiteral{}
		p.nextToken()
		return node

	case TOKEN_IDENT:
		// Check for arrow function with single parameter: id =>
		if p.peekTok.Type == TOKEN_ARROW {
			return p.parseArrowFunction()
		}
		// normal ident
		idName := p.curTok.Value
		token := p.curTok
		p.nextToken()
		identNode := &Identifier{
			BaseNode: BaseNode{
				Position: CreatePositionInfo(token, p.filename),
			},
			Name: idName,
		}
		return p.parsePostfix(identNode)

	case TOKEN_SYMBOL:
		if p.curTok.Value == "(" {
			// Check if this could be arrow function parameters
			if p.isArrowFunctionParameters() {
				return p.parseArrowFunction()
			}
			// ( expr )
			p.nextToken()
			expr := p.parseExpression()
			if p.curTok.Value != ")" {
				p.except("Expected ')' after ( expr )")
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
		p.except("Expected '(' after 'func' in the anonymous function")
	}
	params := p.parseFunctionParameters()
	body := p.parseBlockStatement()

	// Convert parameters to args for backward compatibility
	var args []string
	for _, param := range params {
		args = append(args, param.Name)
	}

	return &FunctionLiteral{Args: args, Params: params, Body: body}
}

func (p *Parser) parsePostfix(left Node) Node {
	for {
		switch {
		case p.curTok.Type == TOKEN_OPTIONAL_CHAIN && p.curTok.Value == "?.":
			left = p.parseOptionalAccessExpression(left)
		case p.curTok.Value == "(":
			left = p.parseCallExpression(left)
		case p.curTok.Value == ".":
			left = p.parseAccessExpression(left)
		case p.curTok.Value == "[":
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
	var mem string
	if p.curTok.Type == TOKEN_IDENT {
		mem = p.curTok.Value
	} else if p.curTok.Type == TOKEN_USE {
		mem = "use"
	} else {
		p.except("Expected identifier after '.'")
	}
	p.nextToken()
	node := &AccessExpression{Object: left, Member: mem}
	return p.parsePostfix(node)
}

func (p *Parser) parseOptionalAccessExpression(left Node) Node {
	p.nextToken() // "?."
	var mem string
	if p.curTok.Type == TOKEN_IDENT {
		mem = p.curTok.Value
	} else if p.curTok.Type == TOKEN_USE {
		mem = "use"
	} else {
		p.except("Expected identifier after '?.'")
	}
	p.nextToken()
	node := &OptionalAccessExpression{Object: left, Member: mem}
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

// parseArrayLiteral => [ expr, expr, ...] or array comprehension
func (p *Parser) parseArrayLiteral() Node {
	p.nextToken() // "["

	// Check if this is a comprehension by looking ahead
	if p.isArrayComprehension() {
		return p.parseArrayComprehension()
	}

	// Regular array literal
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

	// Skip newlines after opening brace
	for p.curTok.Value == "\n" {
		p.nextToken()
	}

	// Check if this is an object comprehension
	if p.isObjectComprehension() {
		return p.parseObjectComprehension()
	}

	var pairs []MapPair
	if p.curTok.Value == "}" {
		p.nextToken()
		return &MapLiteral{Pairs: pairs}
	}
	for p.curTok.Value != "}" && p.curTok.Type != TOKEN_EOF {
		// Verificar si es spread operator
		if p.curTok.Type == TOKEN_ELLIPSIS {
			p.nextToken() // consumir "..."
			spreadValue := p.parseExpression()
			// Usar una clave especial para spread que no será usada
			pairs = append(pairs, MapPair{
				Key:   &StringLiteral{Value: "..."},
				Value: &SpreadExpression{Value: spreadValue},
			})
		} else {
			var keyNode Node

			// Soportar expresiones como claves estilo JavaScript
			switch p.curTok.Type {
			case TOKEN_STRING:
				keyNode = &StringLiteral{Value: p.curTok.Value}
				p.nextToken()
			case TOKEN_IDENT:
				// En JavaScript: {foo: "bar"} equivale a {"foo": "bar"}
				keyNode = &StringLiteral{Value: p.curTok.Value}
				p.nextToken()
			case TOKEN_NUMBER:
				// En JavaScript: {123: "bar"} es válido
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
				} else if p.curTok.Value == "[" {
					// Soporte para [expression]: value
					p.nextToken() // consumir "["
					keyNode = p.parseExpression()
					if p.curTok.Value != "]" {
						p.except("Expected ']' after computed key expression")
					}
					p.nextToken() // consumir "]"
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
		}

		if p.curTok.Value == "," {
			p.nextToken()
			// Skip newlines after comma
			for p.curTok.Value == "\n" {
				p.nextToken()
			}
		} else if p.curTok.Value == "\n" {
			// Allow newlines instead of commas
			for p.curTok.Value == "\n" {
				p.nextToken()
			}
		} else if p.curTok.Value == "}" {
			break
		} else {
			p.except("Expected ',', newline, or '}' in map-literal")
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

func (p *Parser) parseDSLDefinition() Node {
	token := p.curTok
	p.nextToken() // consume 'dsl'

	if p.curTok.Type != TOKEN_IDENT {
		p.except("Expected identifier after 'dsl'")
	}

	name := &Identifier{Name: p.curTok.Value}
	p.nextToken()

	if p.curTok.Value != "{" {
		p.except("Expected '{' after DSL name")
	}

	body := p.parseBlockStatement()

	return &DSLDefinition{
		Token: token,
		Name:  name,
		Body:  body,
	}
}

// isArrowFunctionParameters checks if current position looks like arrow function parameters
func (p *Parser) isArrowFunctionParameters() bool {
	if p.curTok.Value != "(" {
		return false
	}

	// Check for simple case: () =>
	if p.peekTok.Value == ")" {
		// Look ahead for => after the )
		// We need to look further ahead than just peekTok
		// Save current state
		savedPos := p.lexer.pos
		savedCol := p.lexer.col
		savedLine := p.lexer.line
		savedCurTok := p.curTok
		savedPeekTok := p.peekTok

		// Advance to )
		p.nextToken()
		// Advance past )
		p.nextToken()

		isArrow := p.curTok.Type == TOKEN_ARROW

		// Restore state
		p.lexer.pos = savedPos
		p.lexer.col = savedCol
		p.lexer.line = savedLine
		p.curTok = savedCurTok
		p.peekTok = savedPeekTok

		return isArrow
	}

	// For more complex cases, use string scanning
	pos := p.lexer.pos
	input := p.lexer.input

	// Find the matching closing parenthesis
	parenCount := 1
	i := pos

	// Find matching closing parenthesis
	for i < len(input) && parenCount > 0 {
		if input[i] == '(' {
			parenCount++
		} else if input[i] == ')' {
			parenCount--
		}
		i++
	}

	if parenCount != 0 {
		return false // Unmatched parentheses
	}

	// Skip whitespace after closing )
	for i < len(input) && (input[i] == ' ' || input[i] == '\t' || input[i] == '\r' || input[i] == '\n') {
		i++
	}

	// Check for =>
	return i+1 < len(input) && input[i] == '=' && input[i+1] == '>'
}

// parseArrowFunction parses arrow function expressions
func (p *Parser) parseArrowFunction() Node {
	var params []Parameter

	if p.curTok.Type == TOKEN_IDENT {
		// Single parameter without parentheses: x =>
		paramName := p.curTok.Value
		p.nextToken() // consume parameter name
		params = append(params, Parameter{Name: paramName, DefaultValue: nil})
	} else if p.curTok.Value == "(" {
		// Parameters in parentheses: (x, y) => or () =>
		params = p.parseFunctionParameters()
	} else {
		p.except("Expected parameter or '(' in arrow function")
	}

	// Expect =>
	if p.curTok.Type != TOKEN_ARROW {
		p.except("Expected '=>' in arrow function")
	}
	p.nextToken() // consume "=>"

	// Parse body
	var body Node
	var isExpression bool

	if p.curTok.Value == "{" {
		// Block body: => { statements }
		body = p.parseBlockStatement()
		isExpression = false
	} else {
		// Expression body: => expression
		body = p.parseExpression()
		isExpression = true
	}

	return &ArrowFunction{
		Params:       params,
		Body:         body,
		IsExpression: isExpression,
	}
}

// parseArrayDestructuring parsea: let [a, b, c] = [1, 2, 3]
func (p *Parser) parseArrayDestructuring() Node {
	p.nextToken() // consumir "let"
	p.nextToken() // consumir "["

	var names []string
	for p.curTok.Value != "]" {
		if p.curTok.Type != TOKEN_IDENT {
			p.except("Expected identifier in array destructuring")
		}
		names = append(names, p.curTok.Value)
		p.nextToken()

		if p.curTok.Value == "," {
			p.nextToken()
		} else if p.curTok.Value != "]" {
			p.except("Expected ',' or ']' in array destructuring")
		}
	}
	p.nextToken() // consumir "]"

	if p.curTok.Value != "=" {
		p.except("Expected '=' after array destructuring")
	}
	p.nextToken() // consumir "="

	value := p.parseExpression()

	if p.curTok.Value == ";" {
		p.nextToken()
	}

	return &ArrayDestructuring{Names: names, Value: value}
}

// parseObjectDestructuring parsea: let {name, age} = user
func (p *Parser) parseObjectDestructuring() Node {
	p.nextToken() // consumir "let"
	p.nextToken() // consumir "{"

	var names []string
	for p.curTok.Value != "}" {
		if p.curTok.Type != TOKEN_IDENT {
			p.except("Expected identifier in object destructuring")
		}
		names = append(names, p.curTok.Value)
		p.nextToken()

		if p.curTok.Value == "," {
			p.nextToken()
		} else if p.curTok.Value != "}" {
			p.except("Expected ',' or '}' in object destructuring")
		}
	}
	p.nextToken() // consumir "}"

	if p.curTok.Value != "=" {
		p.except("Expected '=' after object destructuring")
	}
	p.nextToken() // consumir "="

	value := p.parseExpression()

	if p.curTok.Value == ";" {
		p.nextToken()
	}

	return &ObjectDestructuring{Names: names, Value: value}
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

// parseMatchExpression parses match expressions (P3)
func (p *Parser) parseMatchExpression() Node {
	p.nextToken() // consume "match"

	value := p.parseExpression()

	if p.curTok.Value != "{" {
		p.except("Expected '{' after match expression")
	}
	p.nextToken() // consume "{"

	var cases []MatchCase

	for p.curTok.Value != "}" && p.curTok.Type != TOKEN_EOF {
		// Skip newlines
		if p.curTok.Value == "\n" {
			p.nextToken()
			continue
		}

		if p.curTok.Type != TOKEN_CASE {
			p.except("Expected 'case' in match expression")
		}
		p.nextToken() // consume "case"

		// Parse pattern
		pattern := p.parsePattern()

		// Parse optional guard
		var guard Node
		if p.curTok.Type == TOKEN_IDENT && p.curTok.Value == "if" {
			p.nextToken() // consume "if"
			guard = p.parseExpression()
		}

		// Expect '=>'
		if p.curTok.Type != TOKEN_ARROW {
			p.except("Expected '=>' after case pattern")
		}
		p.nextToken() // consume "=>"

		// Parse body
		body := p.parseExpression()

		cases = append(cases, MatchCase{
			Pattern: pattern,
			Guard:   guard,
			Body:    body,
		})

		// Skip optional comma and newlines
		if p.curTok.Value == "," {
			p.nextToken()
		}
		for p.curTok.Value == "\n" {
			p.nextToken()
		}
	}

	if p.curTok.Value != "}" {
		p.except("Expected '}' at end of match expression")
	}
	p.nextToken() // consume "}"

	return &MatchExpression{
		Value: value,
		Cases: cases,
	}
}

// parsePattern parses different types of patterns
func (p *Parser) parsePattern() Pattern {
	switch p.curTok.Type {
	case TOKEN_IDENT:
		if p.curTok.Value == "_" {
			// Wildcard pattern
			p.nextToken()
			return &WildcardPattern{}
		} else {
			// Variable pattern
			name := p.curTok.Value
			p.nextToken()
			return &VariablePattern{Name: name}
		}
	case TOKEN_NUMBER, TOKEN_STRING, TOKEN_TRUE, TOKEN_FALSE, TOKEN_NIL:
		// Literal pattern
		value := p.parseFactor()
		return &LiteralPattern{Value: value}
	case TOKEN_SYMBOL:
		if p.curTok.Value == "[" {
			// Array pattern
			return p.parseArrayPattern()
		}
		if p.curTok.Value == "{" {
			// Object pattern
			return p.parseObjectPattern()
		}
	}

	p.except("Expected pattern in match case")
	return nil
}

// parseArrayPattern parses array destructuring patterns
func (p *Parser) parseArrayPattern() Pattern {
	p.nextToken() // consume "["

	var elements []Pattern
	for p.curTok.Value != "]" && p.curTok.Type != TOKEN_EOF {
		pattern := p.parsePattern()
		elements = append(elements, pattern)

		if p.curTok.Value == "," {
			p.nextToken()
		} else if p.curTok.Value == "]" {
			break
		} else {
			p.except("Expected ',' or ']' in array pattern")
		}
	}

	if p.curTok.Value != "]" {
		p.except("Expected ']' at end of array pattern")
	}
	p.nextToken() // consume "]"

	return &ArrayPattern{Elements: elements}
}

// parseObjectPattern parses object destructuring patterns
func (p *Parser) parseObjectPattern() Pattern {
	p.nextToken() // consume "{"

	fields := make(map[string]Pattern)
	for p.curTok.Value != "}" && p.curTok.Type != TOKEN_EOF {
		if p.curTok.Type != TOKEN_IDENT {
			p.except("Expected field name in object pattern")
		}

		fieldName := p.curTok.Value
		p.nextToken()

		var pattern Pattern
		if p.curTok.Value == ":" {
			p.nextToken() // consume ":"
			pattern = p.parsePattern()
		} else {
			// Shorthand: {name} means {name: name}
			pattern = &VariablePattern{Name: fieldName}
		}

		fields[fieldName] = pattern

		if p.curTok.Value == "," {
			p.nextToken()
		} else if p.curTok.Value == "}" {
			break
		} else {
			p.except("Expected ',' or '}' in object pattern")
		}
	}

	if p.curTok.Value != "}" {
		p.except("Expected '}' at end of object pattern")
	}
	p.nextToken() // consume "}"

	return &ObjectPattern{Fields: fields}
}

// isArrayComprehension checks if the current position indicates array comprehension
func (p *Parser) isArrayComprehension() bool {
	// Save current position
	savedPos := p.lexer.pos
	savedLine := p.lexer.line
	savedCol := p.lexer.col
	savedCurTok := p.curTok
	savedPeekTok := p.peekTok

	// Try to parse as comprehension pattern: expression FOR
	defer func() {
		// Restore position
		p.lexer.pos = savedPos
		p.lexer.line = savedLine
		p.lexer.col = savedCol
		p.curTok = savedCurTok
		p.peekTok = savedPeekTok
	}()

	// Skip the expression part (could be complex)
	depth := 0
	for p.curTok.Type != TOKEN_EOF {
		if p.curTok.Value == "[" || p.curTok.Value == "(" || p.curTok.Value == "{" {
			depth++
		} else if p.curTok.Value == "]" || p.curTok.Value == ")" || p.curTok.Value == "}" {
			depth--
		}

		// If we find "for" at the same depth level, it's a comprehension
		if depth == 0 && p.curTok.Type == TOKEN_IDENT && strings.ToLower(p.curTok.Value) == FOR {
			return true
		}

		// If we hit closing bracket without finding "for", it's a regular array
		if depth < 0 {
			return false
		}

		p.nextToken()
	}

	return false
}

// parseArrayComprehension parses array comprehension: [expr for var in iterable if condition]
func (p *Parser) parseArrayComprehension() Node {
	// Parse the expression
	expression := p.parseExpression()

	// Parse generators and conditions
	generators := []Generator{}
	conditions := []Node{}

	for p.curTok.Value != "]" && p.curTok.Type != TOKEN_EOF {
		if p.curTok.Type == TOKEN_IDENT && strings.ToLower(p.curTok.Value) == FOR {
			// Parse generator: for var in iterable
			p.nextToken() // consume "for"

			if p.curTok.Type != TOKEN_IDENT {
				p.except("Expected variable name after 'for'")
			}
			variable := p.curTok.Value
			p.nextToken()

			if p.curTok.Type != TOKEN_IDENT || strings.ToLower(p.curTok.Value) != IN {
				p.except("Expected 'in' after variable in comprehension")
			}
			p.nextToken() // consume "in"

			iterator := p.parseExpression()

			generators = append(generators, Generator{
				Variable: variable,
				Iterator: iterator,
			})
		} else if p.curTok.Type == TOKEN_IDENT && strings.ToLower(p.curTok.Value) == IF {
			// Parse condition: if condition
			p.nextToken() // consume "if"
			condition := p.parseExpression()
			conditions = append(conditions, condition)
		} else {
			break
		}
	}

	if p.curTok.Value != "]" {
		p.except("Expected ']' at end of array comprehension")
	}
	p.nextToken() // consume "]"

	return &ArrayComprehension{
		Expression: expression,
		Generators: generators,
		Conditions: conditions,
	}
}

// isObjectComprehension checks if this is object comprehension
func (p *Parser) isObjectComprehension() bool {
	// Save current position
	savedPos := p.lexer.pos
	savedLine := p.lexer.line
	savedCol := p.lexer.col
	savedCurTok := p.curTok
	savedPeekTok := p.peekTok

	defer func() {
		// Restore position
		p.lexer.pos = savedPos
		p.lexer.line = savedLine
		p.lexer.col = savedCol
		p.curTok = savedCurTok
		p.peekTok = savedPeekTok
	}()

	// Look for pattern: key: value for var in iterable
	depth := 0
	foundColon := false

	for p.curTok.Type != TOKEN_EOF {
		if p.curTok.Value == "[" || p.curTok.Value == "(" || p.curTok.Value == "{" {
			depth++
		} else if p.curTok.Value == "]" || p.curTok.Value == ")" || p.curTok.Value == "}" {
			depth--
		}

		if depth == 0 && p.curTok.Value == ":" && !foundColon {
			foundColon = true
		}

		// If we find "for" after a colon at the same depth, it's object comprehension
		if depth == 0 && foundColon && p.curTok.Type == TOKEN_IDENT && strings.ToLower(p.curTok.Value) == FOR {
			return true
		}

		if depth < 0 {
			return false
		}

		p.nextToken()
	}

	return false
}

// parseObjectComprehension parses object comprehension: {key: value for var in iterable if condition}
func (p *Parser) parseObjectComprehension() Node {
	// Parse key expression
	keyExpr := p.parseExpression()

	if p.curTok.Value != ":" {
		p.except("Expected ':' after key in object comprehension")
	}
	p.nextToken() // consume ":"

	// Parse value expression
	valueExpr := p.parseExpression()

	// Parse generators and conditions
	generators := []Generator{}
	conditions := []Node{}

	for p.curTok.Value != "}" && p.curTok.Type != TOKEN_EOF {
		if p.curTok.Type == TOKEN_IDENT && strings.ToLower(p.curTok.Value) == FOR {
			// Parse generator
			p.nextToken() // consume "for"

			if p.curTok.Type != TOKEN_IDENT {
				p.except("Expected variable name after 'for'")
			}
			variable := p.curTok.Value
			p.nextToken()

			if p.curTok.Type != TOKEN_IDENT || strings.ToLower(p.curTok.Value) != IN {
				p.except("Expected 'in' after variable in comprehension")
			}
			p.nextToken() // consume "in"

			iterator := p.parseExpression()

			generators = append(generators, Generator{
				Variable: variable,
				Iterator: iterator,
			})
		} else if p.curTok.Type == TOKEN_IDENT && strings.ToLower(p.curTok.Value) == IF {
			// Parse condition
			p.nextToken() // consume "if"
			condition := p.parseExpression()
			conditions = append(conditions, condition)
		} else {
			break
		}
	}

	if p.curTok.Value != "}" {
		p.except("Expected '}' at end of object comprehension")
	}
	p.nextToken() // consume "}"

	return &ObjectComprehension{
		KeyExpr:    keyExpr,
		ValueExpr:  valueExpr,
		Generators: generators,
		Conditions: conditions,
	}
}
