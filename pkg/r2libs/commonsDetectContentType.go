package r2libs

import (
	"encoding/json"
	"encoding/xml"
	"regexp"
	"strings"
)

// ContentType representa el tipo de contenido detectado.
type ContentType int

const (
	Unknown ContentType = iota
	JSONType
	XMLType
	HTMLType
	TextType
)

// String convierte el ContentType a una representación en string.
func (c ContentType) String() string {
	switch c {
	case JSONType:
		return "application/json"
	case XMLType:
		return "application/xml"
	case HTMLType:
		return "text/html"
	case TextType:
		return "text/plain"
	default:
		return "text/plain"
	}
}

// DetectContentType determina si una cadena es JSON, XML, HTML o Texto Plano.
func DetectContentType(input string) ContentType {
	trimmed := strings.TrimSpace(input)
	if len(trimmed) == 0 {
		return Unknown
	}

	firstChar := trimmed[0]

	// Heurísticas Iniciales
	switch firstChar {
	case '{', '[':
		// Posible JSON
		if isJSON(trimmed) {
			return JSONType
		}
	case '<':
		// Posible XML o HTML
		if strings.HasPrefix(trimmed, "<!DOCTYPE html") || strings.HasPrefix(strings.ToLower(trimmed), "<html") {
			if isHTML(trimmed) {
				return HTMLType
			}
		}
		if isHTML(trimmed) {
			return HTMLType
		}
		if isXML(trimmed) {
			return XMLType
		}

	}

	// Intentar parsear como JSON
	if isJSON(trimmed) {
		return JSONType
	}

	// Intentar parsear como XML
	if isXML(trimmed) {
		return XMLType
	}

	// Intentar parsear como HTML
	if isHTML(trimmed) {
		return HTMLType
	}

	// Si ninguna coincide, clasificar como Texto Plano
	return TextType
}

// isJSON intenta deserializar la cadena como JSON.
func isJSON(input string) bool {
	var js interface{}
	return json.Unmarshal([]byte(input), &js) == nil
}

// isXML intenta deserializar la cadena como XML.
func isXML(input string) bool {
	var x interface{}
	return xml.Unmarshal([]byte(input), &x) == nil
}

// isHTML intenta determinar si una cadena es HTML utilizando expresiones regulares.
func isHTML(input string) bool {
	// Expresiones regulares básicas para detectar HTML
	htmlPatterns := []string{
		`(?i)<!DOCTYPE\s+html>`, // <!DOCTYPE html>
		`(?i)<html\b`,           // <html
		`(?i)<head\b`,           // <head
		`(?i)<body\b`,           // <body
		`(?i)<div\b`,            // <div
		`(?i)<span\b`,           // <span
		`(?i)<p\b`,              // <p
		`(?i)<a\b`,              // <a
		`(?i)<script\b`,         // <script
		`(?i)<style\b`,          // <style
	}

	for _, pattern := range htmlPatterns {
		matched, err := regexp.MatchString(pattern, input)
		if err != nil {
			continue // Ignorar errores en la expresión regular
		}
		if matched {
			return true
		}
	}

	return false
}
