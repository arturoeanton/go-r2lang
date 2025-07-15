package r2core

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

// ============================================================
// 1) TOKENS
// ============================================================

const (
	TOKEN_EOF             = "EOF"
	TOKEN_NUMBER          = "NUMBER"
	TOKEN_STRING          = "STRING"
	TOKEN_TEMPLATE_STRING = "TEMPLATE_STRING"
	TOKEN_DATE            = "DATE"
	TOKEN_IDENT           = "IDENT"
	TOKEN_ARROW           = "ARROW"
	TOKEN_SYMBOL          = "SYMBOL"
	TOKEN_IMPORT          = "IMPORT"
	TOKEN_AS              = "AS"

	RETURN   = "return"
	LET      = "let"
	VAR      = "var"
	FUNC     = "func"
	FUNCTION = "function"
	METHOD   = "method"

	IF       = "if"
	WHILE    = "while"
	FOR      = "for"
	IN       = "in"
	OBJECT   = "obj"
	CLASS    = "class"
	EXTENDS  = "extends"
	IMPORT   = "import"
	AS       = "as"
	TRY      = "try"
	CATCH    = "catch"
	FINALLY  = "finally"
	THROW    = "throw"
	BREAK    = "break"
	CONTINUE = "continue"

	// Testing framework tokens - will be replaced with new system
	TOKEN_BREAK    = "BREAK"
	TOKEN_CONTINUE = "CONTINUE"
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

// isLetter ya no se usa - reemplazado por isValidIdentifierStart/Char
func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		ch == '_' ||
		ch == '$'
}
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// Nuevas funciones Unicode para identificadores
func isValidIdentifierStart(r rune) bool {
	return unicode.IsLetter(r) || r == '_' || r == '$'
}

func isValidIdentifierChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '$'
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

	if token, ok := l.parseSymbolToken(ch); ok {
		return token
	}

	if token, ok := l.parseNumberToken(ch); ok {
		return token
	}

	if token, ok := l.parseIdentifierToken(ch); ok {
		return token
	}

	fmt.Fprintf(os.Stderr, "Line: %d,Col: %d\n", l.line, l.col)
	fmt.Fprintf(os.Stderr, "Unexpected character in lexer: %c\n", ch)
	os.Exit(1)
	return Token{}
}

func (l *Lexer) parseSymbolToken(ch byte) (Token, bool) {
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
			return l.parseNumberOrSign(), true
		}
	}

	if ch == '+' {
		nextch := l.input[l.pos+1]
		if nextch == '+' {
			l.currentToken = Token{Type: TOKEN_SYMBOL, Value: "++", Line: l.line, Pos: l.pos, Col: l.col}
			l.pos += 2
			return l.currentToken, true
		}
	}

	if ch == '=' {
		nextch := l.input[l.pos+1]
		if nextch == '>' {
			l.currentToken = Token{Type: TOKEN_ARROW, Value: "=>", Line: l.line, Pos: l.pos, Col: l.col}
			l.pos += 2

			return l.currentToken, true
		}
	}

	if ch == '-' {
		nextch := l.input[l.pos+1]
		if nextch == '-' {
			l.currentToken = Token{Type: TOKEN_SYMBOL, Value: "--", Line: l.line, Pos: l.pos, Col: l.col}
			l.pos += 2
			return l.currentToken, true
		}
	}

	// Símbolos de 1 caracter
	singleCharSymbols := []string{
		"(", ")", "{", "}", "[", "]", ";", ",", "+", "-", "*", "/", ".", ":", "?", "\n",
	}
	for _, s := range singleCharSymbols {
		if string(ch) == s {
			l.nextch()
			l.currentToken = Token{Type: TOKEN_SYMBOL, Value: s, Line: l.line, Pos: l.pos, Col: l.col}
			if s == "\n" {
				l.line++
				l.col = 0
			}
			return l.currentToken, true
		}
	}

	if string(ch) == "=" {
		if l.pos+1 < l.length && l.input[l.pos+1] == '=' {
			l.pos += 2
			l.currentToken = Token{Type: TOKEN_SYMBOL, Value: "==", Line: l.line, Pos: l.pos, Col: l.col}
			return l.currentToken, true
		}
		l.nextch()
		l.currentToken = Token{Type: TOKEN_SYMBOL, Value: "=", Line: l.line, Pos: l.pos, Col: l.col}
		return l.currentToken, true
	}

	// Operadores relacionales
	if ch == '<' || ch == '>' || ch == '!' || ch == '=' {
		if l.pos+1 < l.length {
			nextCh := l.input[l.pos+1]
			if nextCh == '=' {
				op := string(ch) + string(nextCh)
				l.pos += 2
				l.currentToken = Token{Type: TOKEN_SYMBOL, Value: op, Line: l.line, Pos: l.pos, Col: l.col}
				return l.currentToken, true
			}
		}
		l.nextch()
		l.currentToken = Token{Type: TOKEN_SYMBOL, Value: string(ch), Line: l.line, Pos: l.pos, Col: l.col}
		return l.currentToken, true
	}
	return Token{}, false
}

func (l *Lexer) parseNumberToken(ch byte) (Token, bool) {
	// Números sin signo
	if isDigit(ch) {
		return l.parseNumberOrSign(), true
	}
	return Token{}, false
}

func (l *Lexer) parseIdentifierToken(ch byte) (Token, bool) {
	// Template String (backticks)
	if ch == '`' {
		return l.parseTemplateString(), true
	}
	// Literal de fecha (@fecha)
	if ch == '@' {
		return l.parseDateLiteral(), true
	}
	// Cadena con soporte Unicode
	if ch == '"' || ch == '\'' {
		val := l.readUnicodeString(ch)
		l.currentToken = Token{Type: TOKEN_STRING, Value: val, Line: l.line, Pos: l.pos, Col: l.col}
		return l.currentToken, true
	}
	// Identificadores Unicode
	r, size := utf8.DecodeRuneInString(l.input[l.pos:])
	if r != utf8.RuneError && isValidIdentifierStart(r) {
		start := l.pos
		l.pos += size
		l.col += size

		for l.pos < l.length {
			r, size := utf8.DecodeRuneInString(l.input[l.pos:])
			if r == utf8.RuneError || !isValidIdentifierChar(r) {
				break
			}
			l.pos += size
			l.col += size
		}
		literal := l.input[start:l.pos]
		switch strings.ToLower(literal) {
		case strings.ToLower(IMPORT):
			l.currentToken = Token{Type: TOKEN_IMPORT, Value: literal, Line: l.line, Pos: l.pos, Col: l.col}
			return l.currentToken, true
		case strings.ToLower(AS):
			l.currentToken = Token{Type: TOKEN_AS, Value: literal, Line: l.line, Pos: l.pos, Col: l.col}
			return l.currentToken, true
		case strings.ToLower(BREAK):
			l.currentToken = Token{Type: TOKEN_BREAK, Value: literal, Line: l.line, Pos: l.pos, Col: l.col}
			return l.currentToken, true
		case strings.ToLower(CONTINUE):
			l.currentToken = Token{Type: TOKEN_CONTINUE, Value: literal, Line: l.line, Pos: l.pos, Col: l.col}
			return l.currentToken, true
			// ... otras palabras clave
		default:
			l.currentToken = Token{Type: TOKEN_IDENT, Value: literal, Line: l.line, Pos: l.pos, Col: l.col}
			return l.currentToken, true
		}
	}
	return Token{}, false
}

// parseTemplateString parses template strings with backticks and ${} interpolation
func (l *Lexer) parseTemplateString() Token {
	start := l.pos
	l.nextch() // Skip opening backtick

	var parts []string
	var currentPart strings.Builder

	for l.pos < l.length && l.input[l.pos] != '`' {
		if l.input[l.pos] == '\\' && l.pos+1 < l.length {
			// Handle escape sequences
			l.nextch()
			switch l.input[l.pos] {
			case '`':
				currentPart.WriteByte('`')
			case '$':
				currentPart.WriteByte('$')
			case '\\':
				currentPart.WriteByte('\\')
			case 'n':
				currentPart.WriteByte('\n')
			case 't':
				currentPart.WriteByte('\t')
			case 'r':
				currentPart.WriteByte('\r')
			default:
				currentPart.WriteByte('\\')
				currentPart.WriteByte(l.input[l.pos])
			}
			l.nextch()
		} else if l.input[l.pos] == '$' && l.pos+1 < l.length && l.input[l.pos+1] == '{' {
			// Found interpolation ${...}
			// Add current text part
			if currentPart.Len() > 0 {
				parts = append(parts, "TEXT:"+currentPart.String())
				currentPart.Reset()
			}

			// Skip ${
			l.pos += 2

			// Find matching }
			braceCount := 1
			exprStart := l.pos
			for l.pos < l.length && braceCount > 0 {
				if l.input[l.pos] == '{' {
					braceCount++
				} else if l.input[l.pos] == '}' {
					braceCount--
				}
				if braceCount > 0 {
					l.nextch()
				}
			}

			if braceCount > 0 {
				panic("Unclosed interpolation in template string")
			}

			// Add expression part
			expr := l.input[exprStart:l.pos]
			parts = append(parts, "EXPR:"+expr)
			l.nextch() // Skip closing }
		} else {
			currentPart.WriteByte(l.input[l.pos])
			l.nextch()
		}
	}

	if l.pos >= l.length {
		panic("Closing backtick of template string expected")
	}

	// Add final text part if any
	if currentPart.Len() > 0 {
		parts = append(parts, "TEXT:"+currentPart.String())
	}

	l.nextch() // Skip closing backtick

	// Encode parts as a single string value
	var encoded strings.Builder
	for i, part := range parts {
		if i > 0 {
			encoded.WriteString("\x00") // Use null separator
		}
		encoded.WriteString(part)
	}

	l.currentToken = Token{Type: TOKEN_TEMPLATE_STRING, Value: encoded.String(), Line: l.line, Pos: start, Col: l.col}
	return l.currentToken
}

// readUnicodeString lee un string con soporte para escape sequences Unicode
func (l *Lexer) readUnicodeString(delimiter byte) string {
	var result strings.Builder
	l.nextch() // saltar comilla inicial

	for l.pos < l.length && l.input[l.pos] != delimiter {
		if l.input[l.pos] == '\\' {
			// Manejar secuencias de escape
			l.nextch()
			if l.pos >= l.length {
				panic("String termina con escape incompleto")
			}

			escaped := l.handleEscape()
			result.WriteString(escaped)
		} else {
			// Verificar que es UTF-8 válido
			r, size := utf8.DecodeRuneInString(l.input[l.pos:])
			if r == utf8.RuneError {
				panic("String contiene UTF-8 inválido")
			}
			result.WriteRune(r)
			l.pos += size
			l.col += size
		}
	}

	if l.pos >= l.length {
		panic("String sin cerrar: falta comilla de cierre")
	}

	l.nextch() // saltar comilla final
	return result.String()
}

// handleEscape maneja secuencias de escape incluyendo Unicode
func (l *Lexer) handleEscape() string {
	if l.pos >= l.length {
		panic("Escape incompleto al final del string")
	}

	ch := l.input[l.pos]
	l.nextch()

	switch ch {
	case 'n':
		return "\n"
	case 't':
		return "\t"
	case 'r':
		return "\r"
	case '\\':
		return "\\"
	case '"':
		return "\""
	case '\'':
		return "'"
	case 'u':
		// \uXXXX - Unicode básico (4 dígitos hex)
		return l.readUnicodeHex(4)
	case 'U':
		// \UXXXXXXXX - Unicode extendido (8 dígitos hex)
		return l.readUnicodeHex(8)
	case 'x':
		// \xXX - Hex básico (2 dígitos hex)
		return l.readUnicodeHex(2)
	default:
		// Escape desconocido, retornar el carácter literal
		return string(ch)
	}
}

// readUnicodeHex lee dígitos hexadecimales y los convierte a Unicode
func (l *Lexer) readUnicodeHex(digits int) string {
	if l.pos+digits > l.length {
		panic("Escape Unicode incompleto")
	}

	hexStr := l.input[l.pos : l.pos+digits]
	l.pos += digits
	l.col += digits

	// Validar que todos son dígitos hex válidos
	for _, ch := range hexStr {
		if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')) {
			panic("Código Unicode inválido: " + hexStr)
		}
	}

	codePoint, err := strconv.ParseInt(hexStr, 16, 32)
	if err != nil {
		panic("Código Unicode inválido: " + hexStr)
	}

	if !utf8.ValidRune(rune(codePoint)) {
		panic("Punto de código Unicode inválido: " + hexStr)
	}

	return string(rune(codePoint))
}

// parseDateLiteral parsea literales de fecha como @2024-12-25 o @"2024-12-25T10:30:00"
func (l *Lexer) parseDateLiteral() Token {
	start := l.pos
	l.nextch() // saltar @

	var dateStr string

	if l.pos < l.length && (l.input[l.pos] == '"' || l.input[l.pos] == '\'') {
		// Fecha con formato específico @"2024-12-25T10:30:00"
		quote := l.input[l.pos]
		l.nextch() // saltar comilla inicial

		strStart := l.pos
		for l.pos < l.length && l.input[l.pos] != quote {
			l.nextch()
		}

		if l.pos >= l.length {
			panic("Closing quote of date literal expected")
		}

		dateStr = l.input[strStart:l.pos]
		l.nextch() // saltar comilla final
	} else {
		// Fecha simple @2024-12-25
		strStart := l.pos
		for l.pos < l.length && (isDigit(l.input[l.pos]) || l.input[l.pos] == '-' || l.input[l.pos] == ':' || l.input[l.pos] == 'T' || l.input[l.pos] == 'Z' || l.input[l.pos] == '+' || l.input[l.pos] == '.') {
			l.nextch()
		}
		dateStr = l.input[strStart:l.pos]
	}

	if dateStr == "" {
		panic("Empty date literal")
	}

	l.currentToken = Token{Type: TOKEN_DATE, Value: dateStr, Line: l.line, Pos: start, Col: l.col}
	return l.currentToken
}
