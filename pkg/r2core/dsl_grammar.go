package r2core

import (
	"fmt"
	"regexp"
	"strings"
)

// DSLGrammar representa la gramática de un DSL
type DSLGrammar struct {
	Rules     map[string]*DSLRule
	Tokens    map[string]*DSLToken
	StartRule string
	Actions   map[string]func([]interface{}) interface{}
}

// DSLRule representa una regla de producción en la gramática
type DSLRule struct {
	Name         string
	Alternatives []*DSLAlternative
}

// DSLAlternative representa una alternativa en una regla
type DSLAlternative struct {
	Sequence []string
	Action   string
}

// DSLToken representa un token en la gramática
type DSLToken struct {
	Name     string
	Pattern  string
	Regex    *regexp.Regexp
	Priority int // Higher priority = matched first (0 = lowest, 100 = highest)
}

// DSLParser parsea código DSL usando una gramática específica
type DSLParser struct {
	Grammar *DSLGrammar
	Tokens  []DSLTokenMatch
	Pos     int
	DSL     *DSLDefinition // Reference to the DSL definition for environment access
}

// DSLTokenMatch representa un token encontrado en el código
type DSLTokenMatch struct {
	Type  string
	Value string
	Start int
	End   int
}

// NewDSLGrammar crea una nueva gramática DSL
func NewDSLGrammar() *DSLGrammar {
	return &DSLGrammar{
		Rules:   make(map[string]*DSLRule),
		Tokens:  make(map[string]*DSLToken),
		Actions: make(map[string]func([]interface{}) interface{}),
	}
}

// AddRule agrega una regla a la gramática
func (g *DSLGrammar) AddRule(name string, alternatives []string, action string) {
	rule, exists := g.Rules[name]
	if !exists {
		rule = &DSLRule{
			Name:         name,
			Alternatives: []*DSLAlternative{},
		}
		g.Rules[name] = rule
		if g.StartRule == "" {
			g.StartRule = name
		}
	}

	// If we have only one alternative, treat it as a sequence
	if len(alternatives) == 1 {
		sequence := strings.Fields(alternatives[0])
		rule.Alternatives = append(rule.Alternatives, &DSLAlternative{
			Sequence: sequence,
			Action:   action,
		})
	} else {
		// Multiple alternatives - each one is a separate sequence
		for _, alt := range alternatives {
			sequence := strings.Fields(alt)
			rule.Alternatives = append(rule.Alternatives, &DSLAlternative{
				Sequence: sequence,
				Action:   action,
			})
		}
	}
}

// AddToken agrega un token a la gramática con prioridad predeterminada
func (g *DSLGrammar) AddToken(name, pattern string) error {
	return g.AddTokenWithPriority(name, pattern, 0) // Default priority
}

// AddTokenWithPriority agrega un token con prioridad específica
func (g *DSLGrammar) AddTokenWithPriority(name, pattern string, priority int) error {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	g.Tokens[name] = &DSLToken{
		Name:     name,
		Pattern:  pattern,
		Regex:    regex,
		Priority: priority,
	}

	return nil
}

// AddKeywordToken agrega un token de keyword con alta prioridad
func (g *DSLGrammar) AddKeywordToken(name, keyword string) error {
	// Keywords have high priority (90) and must match exactly
	pattern := "(?i)\\b" + regexp.QuoteMeta(keyword) + "\\b"
	return g.AddTokenWithPriority(name, pattern, 90)
}

// AddAction agrega una acción semántica
func (g *DSLGrammar) AddAction(name string, action func([]interface{}) interface{}) {
	g.Actions[name] = action
}

// NewDSLParser crea un nuevo parser DSL
func NewDSLParser(grammar *DSLGrammar) *DSLParser {
	return &DSLParser{
		Grammar: grammar,
		Tokens:  []DSLTokenMatch{},
		Pos:     0,
		DSL:     nil, // Will be set when used within a DSL context
	}
}

// NewDSLParserWithContext crea un nuevo parser DSL con contexto de DSL
func NewDSLParserWithContext(grammar *DSLGrammar, dsl *DSLDefinition) *DSLParser {
	return &DSLParser{
		Grammar: grammar,
		Tokens:  []DSLTokenMatch{},
		Pos:     0,
		DSL:     dsl,
	}
}

// Tokenize convierte código DSL en tokens
func (p *DSLParser) Tokenize(code string) error {
	code = strings.TrimSpace(code)
	pos := 0

	for pos < len(code) {
		// Skip whitespace
		if code[pos] == ' ' || code[pos] == '\t' || code[pos] == '\n' || code[pos] == '\r' {
			pos++
			continue
		}

		matched := false
		bestMatch := DSLTokenMatch{}
		bestLength := 0

		// Find the best matching token based on priority and length
		// Use a sorted approach to ensure deterministic iteration order
		bestPriority := -1
		tokenNames := make([]string, 0, len(p.Grammar.Tokens))
		for tokenName := range p.Grammar.Tokens {
			tokenNames = append(tokenNames, tokenName)
		}

		// Sort token names for deterministic behavior
		for i := 0; i < len(tokenNames); i++ {
			for j := i + 1; j < len(tokenNames); j++ {
				if tokenNames[i] > tokenNames[j] {
					tokenNames[i], tokenNames[j] = tokenNames[j], tokenNames[i]
				}
			}
		}

		for _, tokenName := range tokenNames {
			token := p.Grammar.Tokens[tokenName]
			if matches := token.Regex.FindStringIndex(code[pos:]); matches != nil && matches[0] == 0 {
				matchLength := matches[1]

				// Priority-based matching: higher priority wins, then longer match wins
				shouldReplace := false
				if token.Priority > bestPriority {
					shouldReplace = true
				} else if token.Priority == bestPriority && matchLength > bestLength {
					shouldReplace = true
				}

				if shouldReplace {
					bestLength = matchLength
					bestPriority = token.Priority
					bestMatch = DSLTokenMatch{
						Type:  tokenName,
						Value: code[pos : pos+matchLength],
						Start: pos,
						End:   pos + matchLength,
					}
					matched = true
				}
			}
		}

		if matched {
			p.Tokens = append(p.Tokens, bestMatch)
			pos += bestLength
		}

		if !matched {
			return fmt.Errorf("unexpected character at position %d: %c", pos, code[pos])
		}
	}

	return nil
}

// Parse parsea los tokens usando la gramática
func (p *DSLParser) Parse(code string) (interface{}, error) {
	// Reset parser state to ensure deterministic behavior
	p.Tokens = []DSLTokenMatch{}
	p.Pos = 0

	err := p.Tokenize(code)
	if err != nil {
		return nil, err
	}

	p.Pos = 0
	return p.parseRule(p.Grammar.StartRule)
}

// parseRule parsea una regla específica
func (p *DSLParser) parseRule(ruleName string) (interface{}, error) {
	rule, exists := p.Grammar.Rules[ruleName]
	if !exists {
		return nil, fmt.Errorf("rule %s not found", ruleName)
	}

	// Intentar cada alternativa
	for _, alt := range rule.Alternatives {
		savedPos := p.Pos
		result, err := p.parseAlternative(alt)
		if err == nil {
			return result, nil
		}
		// Restaurar posición si falla
		p.Pos = savedPos
	}

	return nil, fmt.Errorf("no alternative matched for rule %s", ruleName)
}

// parseAlternative parsea una alternativa específica
func (p *DSLParser) parseAlternative(alt *DSLAlternative) (interface{}, error) {
	var results []interface{}

	for _, symbol := range alt.Sequence {
		if p.Pos >= len(p.Tokens) {
			return nil, fmt.Errorf("unexpected end of input")
		}

		// Check if symbol is a token
		if p.isToken(symbol) {
			if p.Tokens[p.Pos].Type == symbol {
				results = append(results, p.Tokens[p.Pos].Value)
				p.Pos++
			} else {
				return nil, fmt.Errorf("expected token %s, got %s", symbol, p.Tokens[p.Pos].Type)
			}
		} else {
			// Symbol is a rule
			result, err := p.parseRule(symbol)
			if err != nil {
				return nil, err
			}
			// If the result is a ReturnValue, extract its value
			if retVal, ok := result.(*ReturnValue); ok {
				results = append(results, retVal.Value)
			} else {
				results = append(results, result)
			}
		}
	}

	// Apply semantic action if available
	if alt.Action != "" {
		if action, exists := p.Grammar.Actions[alt.Action]; exists {
			result := action(results)
			// If the result is a ReturnValue, extract its value
			if retVal, ok := result.(*ReturnValue); ok {
				return retVal.Value, nil
			}
			return result, nil
		}
	}

	return results, nil
}

// isToken verifica si un símbolo es un token
func (p *DSLParser) isToken(symbol string) bool {
	_, exists := p.Grammar.Tokens[symbol]
	return exists
}

// DSLResult representa el resultado de evaluar código DSL
type DSLResult struct {
	AST    interface{}
	Code   string
	Output interface{}
}

// GetResult returns the final result of the DSL execution
func (r *DSLResult) GetResult() interface{} {
	return r.Output
}

// String returns a string representation of the result
func (r *DSLResult) String() string {
	// If there's no output, show the original code
	if r.Output == nil {
		return fmt.Sprintf("DSL[%s] -> <no result>", r.Code)
	}

	// For simple values, show them directly
	switch v := r.Output.(type) {
	case string:
		return fmt.Sprintf("DSL[%s] -> \"%s\"", r.Code, v)
	case int, int64, float64:
		return fmt.Sprintf("DSL[%s] -> %v", r.Code, v)
	case bool:
		return fmt.Sprintf("DSL[%s] -> %t", r.Code, v)
	default:
		// For complex objects, show type and value
		return fmt.Sprintf("DSL[%s] -> %v", r.Code, v)
	}
}
