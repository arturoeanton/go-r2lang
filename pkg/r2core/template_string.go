package r2core

import (
	"strings"
)

// TemplatePart represents a part of a template string (text or expression)
type TemplatePart struct {
	IsExpression bool
	Content      string // For literal text
	Expression   Node   // For ${expression} parts
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
			parts = append(parts, toStringOptimized(value))
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
				// Create a new parser for the expression
				exprParser := NewParser(exprStr)
				expr := exprParser.parseExpression()
				templateParts = append(templateParts, TemplatePart{
					IsExpression: true,
					Expression:   expr,
				})
			}
		}
	}
	
	return templateParts
}