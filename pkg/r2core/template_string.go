package r2core

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TemplatePart represents a part of a template string (text or expression)
type TemplatePart struct {
	IsExpression bool
	Content      string // For literal text
	Expression   Node   // For ${expression} parts
	Format       string // For ${expression:format} parts
}

// TemplateString represents a template string with interpolation
type TemplateString struct {
	Parts []TemplatePart
}

// Eval evaluates the template string by combining text and expressions
func (ts *TemplateString) Eval(env *Environment) interface{} {
	if len(ts.Parts) == 0 {
		return ""
	}

	// Optimization: single text part (no interpolation)
	if len(ts.Parts) == 1 && !ts.Parts[0].IsExpression {
		return ts.Parts[0].Content
	}

	// Use optimized string concatenation for multiple parts
	var parts []string
	for _, part := range ts.Parts {
		if part.IsExpression {
			value := part.Expression.Eval(env)
			if part.Format != "" {
				parts = append(parts, formatValue(value, part.Format))
			} else {
				parts = append(parts, toStringOptimized(value))
			}
		} else {
			parts = append(parts, part.Content)
		}
	}

	// Use existing optimization system
	return OptimizedStringConcat(parts...)
}

// parseTemplateParts parses the encoded template string value from lexer
func parseTemplateParts(encoded string, parser *Parser) []TemplatePart {
	if encoded == "" {
		return []TemplatePart{{IsExpression: false, Content: ""}}
	}

	parts := strings.Split(encoded, "\x00")
	var templateParts []TemplatePart

	for _, part := range parts {
		if strings.HasPrefix(part, "TEXT:") {
			content := part[5:] // Remove "TEXT:" prefix
			templateParts = append(templateParts, TemplatePart{
				IsExpression: false,
				Content:      content,
			})
		} else if strings.HasPrefix(part, "EXPR:") {
			exprStr := part[5:] // Remove "EXPR:" prefix
			if exprStr != "" {
				// Check for format specifier (expression:format)
				// But be careful not to split ternary expressions
				var expression string
				var format string

				colonIndex := findFormatColonIndex(exprStr)
				if colonIndex != -1 {
					expression = exprStr[:colonIndex]
					format = exprStr[colonIndex+1:]
				} else {
					expression = exprStr
				}

				// Create a new parser for the expression
				exprParser := NewParser(expression)
				expr := exprParser.parseExpression()
				templateParts = append(templateParts, TemplatePart{
					IsExpression: true,
					Expression:   expr,
					Format:       format,
				})
			}
		}
	}

	return templateParts
}

// formatValue formats a value according to the specified format string
func formatValue(value interface{}, format string) string {
	if format == "" {
		return toStringOptimized(value)
	}

	switch {
	// Number formatting
	case strings.HasPrefix(format, "$"):
		return formatCurrency(value, format)
	case strings.HasSuffix(format, "%"):
		return formatPercentage(value, format)
	case strings.Contains(format, ".") && strings.HasSuffix(format, "f"):
		return formatFloat(value, format)
	case strings.Contains(format, ","):
		return formatNumberWithCommas(value, format)

	// Date formatting
	case strings.Contains(format, "yyyy") || strings.Contains(format, "MM") || strings.Contains(format, "dd") || strings.Contains(format, "HH") || strings.Contains(format, "mm") || strings.Contains(format, "ss"):
		return formatDate(value, format)

	// String formatting
	case strings.HasPrefix(format, "upper"):
		return strings.ToUpper(toStringOptimized(value))
	case strings.HasPrefix(format, "lower"):
		return strings.ToLower(toStringOptimized(value))
	case strings.HasPrefix(format, "title"):
		return strings.Title(toStringOptimized(value))
	case strings.HasPrefix(format, "trim"):
		return strings.TrimSpace(toStringOptimized(value))

	// Default: treat as printf-style format
	default:
		return formatPrintf(value, format)
	}
}

// formatCurrency formats a number as currency
func formatCurrency(value interface{}, format string) string {
	num := toFloat(value)

	// Extract precision from format like $,.2f
	precision := 2
	if strings.Contains(format, ".") {
		parts := strings.Split(format, ".")
		if len(parts) > 1 {
			precStr := strings.TrimRight(parts[1], "f")
			if p, err := strconv.Atoi(precStr); err == nil && p >= 0 {
				precision = p
			}
		}
	}

	formatted := fmt.Sprintf("%."+strconv.Itoa(precision)+"f", num)

	// Add commas if requested
	if strings.Contains(format, ",") {
		formatted = addCommas(formatted)
	}

	return "$" + formatted
}

// formatFloat formats a float with specified precision
func formatFloat(value interface{}, format string) string {
	num := toFloat(value)

	// Extract precision (e.g., ".2f" -> 2)
	precision := 6
	if strings.Contains(format, ".") {
		parts := strings.Split(format, ".")
		if len(parts) > 1 {
			precStr := strings.TrimRight(parts[1], "f")
			if p, err := strconv.Atoi(precStr); err == nil && p >= 0 {
				precision = p
			}
		}
	}

	return fmt.Sprintf("%."+strconv.Itoa(precision)+"f", num)
}

// formatPercentage formats a number as percentage
func formatPercentage(value interface{}, format string) string {
	num := toFloat(value) * 100

	// Extract precision (e.g., ".1%" -> 1)
	precision := 1
	if strings.Contains(format, ".") {
		parts := strings.Split(format, ".")
		if len(parts) > 1 {
			precStr := strings.TrimRight(parts[1], "%")
			if p, err := strconv.Atoi(precStr); err == nil && p >= 0 {
				precision = p
			}
		}
	}

	return fmt.Sprintf("%."+strconv.Itoa(precision)+"f%%", num)
}

// formatNumberWithCommas adds commas to large numbers
func formatNumberWithCommas(value interface{}, format string) string {
	num := toFloat(value)

	var str string
	if num == float64(int64(num)) {
		str = fmt.Sprintf("%.0f", num)
	} else {
		// Use a fixed-point (non-exponential) representation with the
		// minimal digits needed. toStringOptimized/%g switches to
		// scientific notation (e.g. "1.234567891e+06") once a value has
		// enough significant digits, which addCommas cannot group; %f
		// via strconv.FormatFloat never uses exponential notation.
		str = strconv.FormatFloat(num, 'f', -1, 64)
	}
	return addCommas(str)
}

// addCommas adds thousands-separator commas to the integer part of a
// formatted number string. Values that aren't purely numeric (e.g. the
// "NaN"/"+Inf"/"-Inf" strings produced by special float values) are
// returned unchanged instead of having commas spliced into their letters.
func addCommas(s string) string {
	// Handle sign
	sign := ""
	if strings.HasPrefix(s, "-") {
		sign = "-"
		s = s[1:]
	} else if strings.HasPrefix(s, "+") {
		sign = "+"
		s = s[1:]
	}

	// Split at decimal point
	parts := strings.Split(s, ".")
	intPart := parts[0]

	// Only group into thousands when the integer part is actually numeric
	if len(intPart) > 3 && isAllDigits(intPart) {
		var result strings.Builder
		for i, digit := range intPart {
			if i > 0 && (len(intPart)-i)%3 == 0 {
				result.WriteString(",")
			}
			result.WriteRune(digit)
		}
		intPart = result.String()
	}

	result := intPart
	if len(parts) > 1 {
		result += "." + parts[1]
	}

	return sign + result
}

// isAllDigits reports whether s consists entirely of ASCII digits.
func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// formatDate formats a date value
func formatDate(value interface{}, format string) string {
	var t time.Time

	switch v := value.(type) {
	case *DateValue:
		t = v.Time
	case time.Time:
		t = v
	case string:
		// Try to parse common date formats
		layouts := []string{
			"2006-01-02T15:04:05Z07:00",
			"2006-01-02 15:04:05",
			"2006-01-02",
			"15:04:05",
		}
		var err error
		for _, layout := range layouts {
			if t, err = time.Parse(layout, v); err == nil {
				break
			}
		}
		if err != nil {
			return toStringOptimized(value)
		}
	default:
		return toStringOptimized(value)
	}

	// Convert R2Lang format to Go format
	goFormat := convertDateFormat(format)
	return t.Format(goFormat)
}

// convertDateFormat converts R2Lang date format to Go time format
func convertDateFormat(format string) string {
	// Convert common patterns
	format = strings.ReplaceAll(format, "yyyy", "2006")
	format = strings.ReplaceAll(format, "MM", "01")
	format = strings.ReplaceAll(format, "dd", "02")
	format = strings.ReplaceAll(format, "HH", "15")
	format = strings.ReplaceAll(format, "mm", "04")
	format = strings.ReplaceAll(format, "ss", "05")
	return format
}

// formatPrintf uses printf-style formatting
func formatPrintf(value interface{}, format string) (result string) {
	defer func() {
		if r := recover(); r != nil {
			// If printf formatting fails, fallback to string conversion
			result = toStringOptimized(value)
		}
	}()

	// R2Lang numbers are always float64, but Go's fmt package requires an
	// integer type for %d/%x/%X/%o/%b/%c and a string for %s/%q. Passing a
	// raw float64 to those verbs doesn't panic; it silently produces a
	// useless placeholder like "%!d(float64=5)" or formats the verb against
	// the float's bit pattern (e.g. %x on 255.0 yields the hex float
	// "0x1.fep+07" instead of "ff"). Convert to the type the verb expects
	// first so numeric values format the way callers actually intend.
	if num, isFloat := value.(float64); isFloat && format != "" {
		switch format[len(format)-1] {
		case 'd', 'x', 'X', 'o', 'b', 'c':
			return fmt.Sprintf("%"+format, int64(num))
		case 's', 'q':
			return fmt.Sprintf("%"+format, toStringOptimized(value))
		}
	}

	return fmt.Sprintf("%"+format, value)
}

// findFormatColonIndex finds the colon that separates expression from format
// Returns -1 if no format colon is found (ignores ternary colons)
func findFormatColonIndex(exprStr string) int {
	// Simple approach: Look for the rightmost colon that appears to be a format specifier
	// Only if it's not inside parentheses (which would indicate nested ternary)
	parenDepth := 0
	inString := false
	stringChar := byte(0)

	// Find potential format colons from right to left
	for i := len(exprStr) - 1; i >= 0; i-- {
		char := exprStr[i]

		// Handle string escaping (backward)
		if i > 0 && exprStr[i-1] == '\\' {
			continue
		}

		// Handle string boundaries (backward)
		if !inString && (char == '"' || char == '\'' || char == '`') {
			inString = true
			stringChar = char
			continue
		}
		if inString && char == stringChar {
			inString = false
			stringChar = 0
			continue
		}

		// Skip characters inside strings
		if inString {
			continue
		}

		// Track parentheses depth
		if char == ')' {
			parenDepth++
		} else if char == '(' {
			parenDepth--
		} else if char == ':' && parenDepth == 0 {
			// This could be a format colon - check if what follows looks like a format
			if i+1 < len(exprStr) && isFormatSpecifier(exprStr[i+1:]) {
				// Also check if there's a question mark before this colon (indicates ternary)
				hasQuestionMark := false
				for j := i - 1; j >= 0; j-- {
					if exprStr[j] == '?' && !isInsideString(exprStr, j) {
						hasQuestionMark = true
						break
					}
				}
				// If there's no question mark before this colon, it's likely a format colon
				if !hasQuestionMark {
					return i
				}
			}
		}
	}

	return -1
}

// isInsideString checks if a position is inside a string literal
func isInsideString(expr string, pos int) bool {
	inString := false
	stringChar := byte(0)

	for i := 0; i < pos && i < len(expr); i++ {
		char := expr[i]

		// Handle escaping
		if i > 0 && expr[i-1] == '\\' {
			continue
		}

		// Handle string boundaries
		if !inString && (char == '"' || char == '\'' || char == '`') {
			inString = true
			stringChar = char
		} else if inString && char == stringChar {
			inString = false
			stringChar = 0
		}
	}

	return inString
}

// isFormatSpecifier checks if a string looks like a format specifier
func isFormatSpecifier(s string) bool {
	if s == "" {
		return false
	}

	// Check specific format patterns
	// Currency formats: $,.2f, $,f
	if strings.HasPrefix(s, "$") {
		return true
	}

	// Percentage formats: .1%, .2%, %
	if strings.HasSuffix(s, "%") {
		return true
	}

	// Float formats: .2f, .3f, .f
	if strings.Contains(s, ".") && strings.HasSuffix(s, "f") {
		return true
	}

	// Comma formatting: ,
	if s == "," {
		return true
	}

	// String formatting
	if s == "upper" || s == "lower" || s == "title" || s == "trim" {
		return true
	}

	// Date formatting patterns
	if strings.Contains(s, "yyyy") || strings.Contains(s, "MM") || strings.Contains(s, "dd") ||
		strings.Contains(s, "HH") || strings.Contains(s, "mm") || strings.Contains(s, "ss") {
		return true
	}

	// Printf-style formatting: d, g, s, x, X
	if len(s) == 1 && (s == "d" || s == "g" || s == "s" || s == "x" || s == "X") {
		return true
	}

	return false
}
