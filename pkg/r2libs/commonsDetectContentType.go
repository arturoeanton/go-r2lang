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
		// Posible XML o HTML. Una página HTML real trae una señal
		// estructural inequívoca (DOCTYPE, <html>, <head> o <body>); sin esa
		// señal, isHTML() sólo detecta si ALGUNA de sus patrones (incluyendo
		// nombres de tag tan genéricos como "<a" o "<p") aparece en
		// cualquier parte del documento, lo cual da falsos positivos
		// constantes contra XML de datos legítimo (p. ej.
		// "<order><a>123</a></order>" se marcaba como text/html sólo por el
		// tag <a>, un nombre de elemento XML perfectamente válido). Por eso
		// una señal débil sólo gana si el documento ni siquiera es XML bien
		// formado.
		if hasStrongHTMLSignal(trimmed) {
			return HTMLType
		}
		if isXML(trimmed) {
			return XMLType
		}
		if isHTML(trimmed) {
			return HTMLType
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

// hasStrongHTMLSignal reporta si input trae una marca estructural
// inequívoca de documento HTML (DOCTYPE o el tag raíz <html>/<head>/<body>),
// a diferencia de isHTML() que también dispara con nombres de tag genéricos
// que son igual de válidos como elementos XML.
func hasStrongHTMLSignal(input string) bool {
	strongPatterns := []string{
		`(?i)^<!DOCTYPE\s+html`,
		`(?i)^<html\b`,
		`(?i)^<head\b`,
		`(?i)^<body\b`,
	}
	for _, pattern := range strongPatterns {
		matched, err := regexp.MatchString(pattern, input)
		if err == nil && matched {
			return true
		}
	}
	return false
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
