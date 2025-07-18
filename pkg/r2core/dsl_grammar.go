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
	Name    string
	Pattern string
	Regex   *regexp.Regexp
}

// DSLParser parsea código DSL usando una gramática específica
type DSLParser struct {
	Grammar *DSLGrammar
	Tokens  []DSLTokenMatch
	Pos     int
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
	rule := &DSLRule{
		Name:         name,
		Alternatives: []*DSLAlternative{},
	}

	for _, alt := range alternatives {
		sequence := strings.Fields(alt)
		rule.Alternatives = append(rule.Alternatives, &DSLAlternative{
			Sequence: sequence,
			Action:   action,
		})
	}

	g.Rules[name] = rule
	if g.StartRule == "" {
		g.StartRule = name
	}
}

// AddToken agrega un token a la gramática
func (g *DSLGrammar) AddToken(name, pattern string) error {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	g.Tokens[name] = &DSLToken{
		Name:    name,
		Pattern: pattern,
		Regex:   regex,
	}

	return nil
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
		for tokenName, token := range p.Grammar.Tokens {
			if matches := token.Regex.FindStringIndex(code[pos:]); matches != nil && matches[0] == 0 {
				value := code[pos : pos+matches[1]]
				p.Tokens = append(p.Tokens, DSLTokenMatch{
					Type:  tokenName,
					Value: value,
					Start: pos,
					End:   pos + matches[1],
				})
				pos += matches[1]
				matched = true
				break
			}
		}

		if !matched {
			return fmt.Errorf("unexpected character at position %d: %c", pos, code[pos])
		}
	}

	return nil
}

// Parse parsea los tokens usando la gramática
func (p *DSLParser) Parse(code string) (interface{}, error) {
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
			results = append(results, result)
		}
	}

	// Apply semantic action if available
	if alt.Action != "" {
		if action, exists := p.Grammar.Actions[alt.Action]; exists {
			return action(results), nil
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
