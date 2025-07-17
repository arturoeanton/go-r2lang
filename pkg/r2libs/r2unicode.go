package r2libs

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)

func RegisterUnicode(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		"ulen":          r2core.BuiltinFunction(unicodeLength),
		"usubstr":       r2core.BuiltinFunction(unicodeSubstring),
		"uupper":        r2core.BuiltinFunction(unicodeUpper),
		"ulower":        r2core.BuiltinFunction(unicodeLower),
		"utitle":        r2core.BuiltinFunction(unicodeTitle),
		"ureverse":      r2core.BuiltinFunction(unicodeReverse),
		"unormalize":    r2core.BuiltinFunction(unicodeNormalize),
		"ucompare":      r2core.BuiltinFunction(unicodeCompare),
		"uisvalid":      r2core.BuiltinFunction(isValidUTF8),
		"ucharcode":     r2core.BuiltinFunction(getCharCode),
		"ufromcode":     r2core.BuiltinFunction(fromCharCode),
		"uisLetter":     r2core.BuiltinFunction(isLetter),
		"uisDigit":      r2core.BuiltinFunction(isDigit),
		"uisSpace":      r2core.BuiltinFunction(isSpace),
		"uisPunct":      r2core.BuiltinFunction(isPunct),
		"uisUpper":      r2core.BuiltinFunction(isUpper),
		"uisLower":      r2core.BuiltinFunction(isLower),
		"ugetCategory":  r2core.BuiltinFunction(getCategory),
		"uregex":        r2core.BuiltinFunction(unicodeRegex),
		"uregexMatch":   r2core.BuiltinFunction(unicodeRegexMatch),
		"uregexReplace": r2core.BuiltinFunction(unicodeRegexReplace),
	}

	RegisterModule(env, "unicode", functions)
}

// ============================================================
// FUNCIONES BÁSICAS DE STRING UNICODE
// ============================================================

// unicodeLength retorna la longitud en caracteres Unicode (no bytes)
func unicodeLength(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("ulen() requiere exactamente 1 argumento")
	}

	str, ok := args[0].(string)
	if !ok {
		panic("ulen() requiere un string")
	}

	return float64(utf8.RuneCountInString(str))
}

// unicodeSubstring extrae substring basado en caracteres Unicode
func unicodeSubstring(args ...interface{}) interface{} {
	if len(args) < 2 || len(args) > 3 {
		panic("usubstr() requiere 2 o 3 argumentos")
	}

	str, ok := args[0].(string)
	if !ok {
		panic("usubstr() requiere un string como primer argumento")
	}

	start, ok := args[1].(float64)
	if !ok {
		panic("usubstr() requiere un número como segundo argumento")
	}

	runes := []rune(str)
	startIdx := int(start)

	if startIdx < 0 {
		startIdx = 0
	}
	if startIdx >= len(runes) {
		return ""
	}

	endIdx := len(runes)
	if len(args) == 3 {
		length, ok := args[2].(float64)
		if !ok {
			panic("usubstr() requiere un número como tercer argumento")
		}
		endIdx = startIdx + int(length)
		if endIdx > len(runes) {
			endIdx = len(runes)
		}
		if endIdx < startIdx {
			endIdx = startIdx
		}
	}

	return string(runes[startIdx:endIdx])
}

// unicodeUpper convierte a mayúsculas usando reglas Unicode
func unicodeUpper(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("uupper() requiere exactamente 1 argumento")
	}

	str, ok := args[0].(string)
	if !ok {
		panic("uupper() requiere un string")
	}

	return strings.ToUpper(str)
}

// unicodeLower convierte a minúsculas usando reglas Unicode
func unicodeLower(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("ulower() requiere exactamente 1 argumento")
	}

	str, ok := args[0].(string)
	if !ok {
		panic("ulower() requiere un string")
	}

	return strings.ToLower(str)
}

// unicodeTitle convierte a formato título
func unicodeTitle(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("utitle() requiere exactamente 1 argumento")
	}

	str, ok := args[0].(string)
	if !ok {
		panic("utitle() requiere un string")
	}

	return strings.ToTitle(str)
}

// unicodeReverse invierte un string respetando caracteres Unicode
func unicodeReverse(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("ureverse() requiere exactamente 1 argumento")
	}

	str, ok := args[0].(string)
	if !ok {
		panic("ureverse() requiere un string")
	}

	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

// ============================================================
// NORMALIZACIÓN UNICODE
// ============================================================

// unicodeNormalize normaliza un string usando formas Unicode estándar
func unicodeNormalize(args ...interface{}) interface{} {
	if len(args) < 1 || len(args) > 2 {
		panic("unormalize() requiere 1 o 2 argumentos")
	}

	str, ok := args[0].(string)
	if !ok {
		panic("unormalize() requiere un string")
	}

	form := "NFC" // Por defecto
	if len(args) == 2 {
		if formArg, ok := args[1].(string); ok {
			form = formArg
		}
	}

	var normalizer norm.Form
	switch form {
	case "NFC":
		normalizer = norm.NFC
	case "NFD":
		normalizer = norm.NFD
	case "NFKC":
		normalizer = norm.NFKC
	case "NFKD":
		normalizer = norm.NFKD
	default:
		panic("Forma de normalización inválida: " + form + " (usar NFC, NFD, NFKC, NFKD)")
	}

	return normalizer.String(str)
}

// unicodeCompare compara strings usando collation Unicode
func unicodeCompare(args ...interface{}) interface{} {
	if len(args) < 2 || len(args) > 3 {
		panic("ucompare() requiere 2 o 3 argumentos")
	}

	str1, ok := args[0].(string)
	if !ok {
		panic("ucompare() requiere strings")
	}

	str2, ok := args[1].(string)
	if !ok {
		panic("ucompare() requiere strings")
	}

	locale := "en"
	if len(args) == 3 {
		if localeArg, ok := args[2].(string); ok {
			locale = localeArg
		}
	}

	tag := language.Make(locale)
	col := collate.New(tag)

	result := col.CompareString(str1, str2)
	return float64(result)
}

// ============================================================
// VALIDACIÓN Y UTILIDADES
// ============================================================

// isValidUTF8 verifica si un string es UTF-8 válido
func isValidUTF8(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("uisvalid() requiere exactamente 1 argumento")
	}

	str, ok := args[0].(string)
	if !ok {
		panic("uisvalid() requiere un string")
	}

	return utf8.ValidString(str)
}

// getCharCode obtiene el código Unicode del primer carácter
func getCharCode(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("ucharcode() requiere exactamente 1 argumento")
	}

	str, ok := args[0].(string)
	if !ok {
		panic("ucharcode() requiere un string")
	}

	if str == "" {
		panic("ucharcode() requiere un string no vacío")
	}

	r, _ := utf8.DecodeRuneInString(str)
	return float64(r)
}

// fromCharCode crea un string desde un código Unicode
func fromCharCode(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("ufromcode() requiere exactamente 1 argumento")
	}

	code, ok := args[0].(float64)
	if !ok {
		panic("ufromcode() requiere un número")
	}

	r := rune(code)
	if !utf8.ValidRune(r) {
		panic("Código Unicode inválido")
	}

	return string(r)
}

// ============================================================
// CLASIFICACIÓN DE CARACTERES
// ============================================================

func getFirstRune(args []interface{}, funcName string) rune {
	if len(args) != 1 {
		panic(funcName + "() requiere exactamente 1 argumento")
	}

	str, ok := args[0].(string)
	if !ok {
		panic(funcName + "() requiere un string")
	}

	if str == "" {
		panic(funcName + "() requiere un string no vacío")
	}

	r, _ := utf8.DecodeRuneInString(str)
	return r
}

func isLetter(args ...interface{}) interface{} {
	char := getFirstRune(args, "uisLetter")
	return unicode.IsLetter(char)
}

func isDigit(args ...interface{}) interface{} {
	char := getFirstRune(args, "uisDigit")
	return unicode.IsDigit(char)
}

func isSpace(args ...interface{}) interface{} {
	char := getFirstRune(args, "uisSpace")
	return unicode.IsSpace(char)
}

func isPunct(args ...interface{}) interface{} {
	char := getFirstRune(args, "uisPunct")
	return unicode.IsPunct(char)
}

func isUpper(args ...interface{}) interface{} {
	char := getFirstRune(args, "uisUpper")
	return unicode.IsUpper(char)
}

func isLower(args ...interface{}) interface{} {
	char := getFirstRune(args, "uisLower")
	return unicode.IsLower(char)
}

func getCategory(args ...interface{}) interface{} {
	char := getFirstRune(args, "ugetCategory")
	for name, table := range unicode.Categories {
		if unicode.In(char, table) {
			return name
		}
	}
	return "Unknown"
}

// ============================================================
// EXPRESIONES REGULARES UNICODE
// ============================================================

// unicodeRegex busca patrones usando regex con soporte Unicode
func unicodeRegex(args ...interface{}) interface{} {
	if len(args) < 2 {
		panic("uregex() requiere al menos 2 argumentos")
	}

	pattern, ok := args[0].(string)
	if !ok {
		panic("Patrón debe ser string")
	}

	text, ok := args[1].(string)
	if !ok {
		panic("Texto debe ser string")
	}

	// Compilar regex con soporte Unicode
	re, err := regexp.Compile(pattern)
	if err != nil {
		panic("Patrón regex inválido: " + err.Error())
	}

	matches := re.FindAllString(text, -1)
	result := make([]interface{}, len(matches))
	for i, match := range matches {
		result[i] = match
	}

	return result
}

// unicodeRegexMatch verifica si un patrón coincide
func unicodeRegexMatch(args ...interface{}) interface{} {
	if len(args) != 2 {
		panic("uregexMatch() requiere exactamente 2 argumentos")
	}

	pattern, ok := args[0].(string)
	if !ok {
		panic("Patrón debe ser string")
	}

	text, ok := args[1].(string)
	if !ok {
		panic("Texto debe ser string")
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		panic("Patrón regex inválido: " + err.Error())
	}

	return re.MatchString(text)
}

// unicodeRegexReplace reemplaza patrones en texto
func unicodeRegexReplace(args ...interface{}) interface{} {
	if len(args) != 3 {
		panic("uregexReplace() requiere exactamente 3 argumentos")
	}

	pattern, ok := args[0].(string)
	if !ok {
		panic("Patrón debe ser string")
	}

	replacement, ok := args[1].(string)
	if !ok {
		panic("Reemplazo debe ser string")
	}

	text, ok := args[2].(string)
	if !ok {
		panic("Texto debe ser string")
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		panic("Patrón regex inválido: " + err.Error())
	}

	return re.ReplaceAllString(text, replacement)
}
