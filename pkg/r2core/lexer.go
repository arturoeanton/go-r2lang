package r2core

import (
	"fmt"
	"os"
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
	TOKEN_ARROW  = "ARROW"
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
	CLASS   = "class"
	EXTENDS = "extends"
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

func Done() {
	wg.Done()
}

func Add() {
	wg.Add(1)
}

func Wait() {
	wg.Wait()
}

type Token struct {
	Type  string
	Value string
	Line  int
	Pos   int
	Col   int
}

// ============================================================
// 2) LEXER (reconoce signo, decimales, comentarios, etc.)
// ============================================================

type Lexer struct {
	input        string
	pos          int
	col          int
	line         int
	length       int
	currentToken Token
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		pos:    0,
		line:   1,
		length: len(input),
	}
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '	' || ch == ''
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

func (l *Lexer) nextch() {
	l.pos++
	l.col++
}

// parseNumberOrSign maneja -2.3, +10, etc.
func (l *Lexer) parseNumberOrSign() Token {
	start := l.pos
	if l.input[l.pos] == '-' || l.input[l.pos] == '+' {
		l.nextch()
	}
	hasDigits := false
	for l.pos < l.length && isDigit(l.input[l.pos]) {
		hasDigits = true
		l.nextch()
	}
	if l.pos < l.length && l.input[l.pos] == '.' {
		l.nextch()
		for l.pos < l.length && isDigit(l.input[l.pos]) {
			hasDigits = true
			l.nextch()
		}
	}
	if !hasDigits {
		panic("Invalid number in " + l.input[start:l.pos])
	}
	val := l.input[start:l.pos]
	l.currentToken = Token{Type: TOKEN_NUMBER, Value: val, Line: l.line, Pos: l.pos, Col: l.col}
	return l.currentToken
}

func (l *Lexer) NextToken() Token {
skipWhitespace:
	for l.pos < l.length {
		ch := l.input[l.pos]
		if isWhitespace(ch) {
			l.nextch()
		} else if ch == '/' {
			// Comentarios
			if l.pos+1 < l.length && l.input[l.pos+1] == '/' {
				// comentario de línea
				l.pos += 2
				for l.pos < l.length && l.input[l.pos] != '\n' {
					l.nextch()
				}
			} else if l.pos+1 < l.length && l.input[l.pos+1] == '*' {
				// /* ... */
				l.pos += 2
				for l.pos < l.length {
					if l.input[l.pos] == '*' && (l.pos+1 < l.length && l.input[l.pos+1] == '/') {
						l.pos += 2
						break
					}
					l.nextch()
				}
			} else {
				break skipWhitespace
			}
		} else {
			break skipWhitespace
		}
	}

	if l.pos >= l.length {
		l.currentToken = Token{Type: TOKEN_EOF, Value: "", Line: l.line, Pos: l.pos, Col: l.col}
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
			l.currentToken = Token{Type: TOKEN_SYMBOL, Value: "++", Line: l.line, Pos: l.pos, Col: l.col}
			l.pos += 2
			return l.currentToken
		}
	}

	if ch == '=' {
		nextch := l.input[l.pos+1]
		if nextch == '>' {
			l.currentToken = Token{Type: TOKEN_ARROW, Value: "=>", Line: l.line, Pos: l.pos, Col: l.col}
			l.pos += 2

			return l.currentToken
		}
	}

	if ch == '-' {
		nextch := l.input[l.pos+1]
		if nextch == '-' {
			l.currentToken = Token{Type: TOKEN_SYMBOL, Value: "--", Line: l.line, Pos: l.pos, Col: l.col}
			l.pos += 2
			return l.currentToken
		}
	}

	// Símbolos de 1 caracter
	singleCharSymbols := []string{
		"(", ")", "{", "}", "[", "]", ";", ",", "+", "-", "*", "/", ".", ":", "\n",
	}
	for _, s := range singleCharSymbols {
		if string(ch) == s {
			l.nextch()
			l.currentToken = Token{Type: TOKEN_SYMBOL, Value: s, Line: l.line, Pos: l.pos, Col: l.col}
			if s == "\n" {
				l.line++
				l.col = 0
			}
			return l.currentToken
		}
	}

	if string(ch) == "=" {
		if l.pos+1 < l.length && l.input[l.pos+1] == '=' {
			l.pos += 2
			l.currentToken = Token{Type: TOKEN_SYMBOL, Value: "==", Line: l.line, Pos: l.pos, Col: l.col}
			return l.currentToken
		}
		l.nextch()
		l.currentToken = Token{Type: TOKEN_SYMBOL, Value: "=", Line: l.line, Pos: l.pos, Col: l.col}
		return l.currentToken
	}

	// Operadores relacionales
	if ch == '<' || ch == '>' || ch == '!' || ch == '=' {
		if l.pos+1 < l.length {
			nextCh := l.input[l.pos+1]
			if nextCh == '=' {
				op := string(ch) + string(nextCh)
				l.pos += 2
				l.currentToken = Token{Type: TOKEN_SYMBOL, Value: op, Line: l.line, Pos: l.pos, Col: l.col}
				return l.currentToken
			}
		}
		l.nextch()
		l.currentToken = Token{Type: TOKEN_SYMBOL, Value: string(ch), Line: l.line, Pos: l.pos, Col: l.col}
		return l.currentToken
	}

	// Cadena
	if ch == '"' || ch == '\'' {
		quote := ch
		start := l.pos
		l.nextch()
		for l.pos < l.length && l.input[l.pos] != quote {
			l.nextch()
		}
		if l.pos >= l.length {
			panic("Closing quote of string expected")
		}
		val := l.input[start+1 : l.pos]
		l.nextch()
		l.currentToken = Token{Type: TOKEN_STRING, Value: val, Line: l.line, Pos: l.pos, Col: l.col}
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
			l.nextch()
		}
		literal := l.input[start:l.pos]
		switch strings.ToLower(literal) {
		case strings.ToLower(IMPORT):
			l.currentToken = Token{Type: TOKEN_IMPORT, Value: literal, Line: l.line, Pos: l.pos, Col: l.col}
			return l.currentToken
		case strings.ToLower(AS):
			l.currentToken = Token{Type: TOKEN_AS, Value: literal, Line: l.line, Pos: l.pos, Col: l.col}
			return l.currentToken
		case "given":
			l.currentToken = Token{Type: TOKEN_GIVEN, Value: "Given", Line: l.line, Pos: l.pos, Col: l.col}
			return l.currentToken
		case "when":
			l.currentToken = Token{Type: TOKEN_WHEN, Value: "When", Line: l.line, Pos: l.pos, Col: l.col}
			return l.currentToken
		case "then":
			l.currentToken = Token{Type: TOKEN_THEN, Value: "Then", Line: l.line, Pos: l.pos, Col: l.col}
			return l.currentToken
		case "and":
			l.currentToken = Token{Type: TOKEN_AND, Value: "And", Line: l.line, Pos: l.pos, Col: l.col}
			return l.currentToken
		case "testcase":
			l.currentToken = Token{Type: TOKEN_TESTCASE, Value: "TestCase", Line: l.line, Pos: l.pos, Col: l.col}
			return l.currentToken
		// ... otras palabras clave
		default:
			l.currentToken = Token{Type: TOKEN_IDENT, Value: literal, Line: l.line, Pos: l.pos, Col: l.col}
			return l.currentToken
		}
	}

	fmt.Fprintf(os.Stderr, "Line: %d,Col: %d\n", l.line, l.col)
	fmt.Fprintf(os.Stderr, "Unexpected character in lexer: %c\n", ch)
	os.Exit(1)
	return Token{}
}
