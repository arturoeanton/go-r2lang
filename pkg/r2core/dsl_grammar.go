package r2core

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/arturoeanton/go-dsl/pkg/dslbuilder"
)

// DSLGrammar adapts R2Lang's dsl{} block builder API (token/rule/action) onto
// a github.com/arturoeanton/go-dsl engine instance.
type DSLGrammar struct {
	engine *dslbuilder.DSL
}

// NewDSLGrammar creates a grammar backed by a fresh go-dsl engine.
func NewDSLGrammar() *DSLGrammar {
	return &DSLGrammar{engine: dslbuilder.New("r2dsl")}
}

// AddRule adds a rule alternative. Multiple calls with the same name create
// alternatives, matching go-dsl's BNF-style semantics. Each element of
// sequence is whitespace-split so callers may pass either one symbol per
// element or a single space-joined sequence (both conventions are used by
// existing callers).
func (g *DSLGrammar) AddRule(name string, sequence []string, action string) {
	var symbols []string
	for _, s := range sequence {
		symbols = append(symbols, strings.Fields(s)...)
	}
	g.engine.Rule(name, symbols, action)
}

// AddToken adds a regex token with default (non-keyword) priority.
func (g *DSLGrammar) AddToken(name, pattern string) error {
	return g.engine.Token(name, pattern)
}

// AddKeywordToken adds a case-insensitive, word-bounded keyword token at
// high priority so it is matched before generic identifier-like tokens.
func (g *DSLGrammar) AddKeywordToken(name, keyword string) error {
	return g.engine.KeywordToken(name, keyword)
}

// AddLiteral adds a token that matches an exact piece of text (an operator
// or punctuation symbol, e.g. "+" or "=="), unlike AddKeywordToken it is
// not word-bounded or case-insensitive, since literals are often not made
// of word characters.
func (g *DSLGrammar) AddLiteral(name, text string) error {
	return g.engine.Token(name, regexp.QuoteMeta(text))
}

// AddAction registers a semantic action for a rule. R2Lang action callbacks
// return a single value; if that value is a Go error it is surfaced as a
// go-dsl action error instead of a successful result.
func (g *DSLGrammar) AddAction(name string, action func([]interface{}) interface{}) {
	g.engine.Action(name, func(args []interface{}) (interface{}, error) {
		result := action(args)
		if err, ok := result.(error); ok {
			return nil, err
		}
		return result, nil
	})
}

// DebugTokens tokenizes code without parsing or evaluating it, useful for
// introspecting how the grammar's tokens match a given input.
func (g *DSLGrammar) DebugTokens(code string) ([]dslbuilder.TokenMatch, error) {
	return g.engine.DebugTokens(code)
}

// Use parses and evaluates code against this grammar, scoping context to
// this call as go-dsl does.
func (g *DSLGrammar) Use(code string, context map[string]interface{}) (*DSLResult, error) {
	result, err := g.engine.Use(code, context)
	if err != nil {
		return nil, err
	}
	return &DSLResult{
		AST:    result.AST,
		Code:   result.Code,
		Output: result.Output,
	}, nil
}

// AST parses code and returns its syntax tree as nested R2Lang maps, without
// running any semantic actions.
func (g *DSLGrammar) AST(code string) (map[string]interface{}, error) {
	node, err := g.engine.ParseAST(code)
	if err != nil {
		return nil, err
	}
	return nodeToMap(node), nil
}

// nodeToMap converts a go-dsl parse tree node into the map/slice shape
// R2Lang scripts can walk directly (property access + array indexing).
func nodeToMap(n *dslbuilder.Node) map[string]interface{} {
	if n == nil {
		return nil
	}
	m := map[string]interface{}{
		"rule":  n.Rule,
		"text":  n.Text(),
		"start": float64(n.Span.Start),
		"end":   float64(n.Span.End),
	}
	if n.IsToken() {
		m["token"] = n.Token.TokenType
	} else {
		m["token"] = ""
	}
	children := make([]interface{}, len(n.Children))
	for i, c := range n.Children {
		children[i] = nodeToMap(c)
	}
	m["children"] = children
	return m
}

// Diagnostics tokenizes and parses code tolerantly, collecting every syntax
// error found instead of stopping at the first one.
func (g *DSLGrammar) Diagnostics(code string) []interface{} {
	diags := g.engine.Diagnostics(code)
	result := make([]interface{}, len(diags))
	for i, d := range diags {
		result[i] = map[string]interface{}{
			"message": d.Message,
			"line":    float64(d.Line),
			"column":  float64(d.Column),
			"token":   d.Token,
		}
	}
	return result
}

// Check reports whether code parses cleanly against this grammar, without
// evaluating it.
func (g *DSLGrammar) Check(code string) (bool, string) {
	if _, err := g.engine.ParseAST(code); err != nil {
		return false, err.Error()
	}
	return true, ""
}

// Completions returns the valid next tokens at a byte offset into code,
// e.g. for editor autocomplete.
func (g *DSLGrammar) Completions(code string, offset int) []interface{} {
	completions := g.engine.Completions(code, offset)
	result := make([]interface{}, len(completions))
	for i, c := range completions {
		result[i] = map[string]interface{}{
			"label":     c.Label,
			"token":     c.Token,
			"isKeyword": c.IsKeyword,
			"detail":    c.Detail,
		}
	}
	return result
}

// Validate statically checks the grammar, returning warnings (informative,
// the grammar still works) and an error joining any structural problems
// that make the grammar unusable.
func (g *DSLGrammar) Validate() ([]string, error) {
	return g.engine.Validate()
}

// DSLResult represents the result of evaluating DSL code: the parse tree
// (AST), the original source (Code), and the evaluated value (Output).
type DSLResult struct {
	AST    interface{}
	Code   string
	Output interface{}
}

// GetResult returns the final evaluated result of the DSL execution.
func (r *DSLResult) GetResult() interface{} {
	return r.Output
}

// String returns a human-readable representation of the result.
func (r *DSLResult) String() string {
	if r.Output == nil {
		return fmt.Sprintf("DSL[%s] -> <no result>", r.Code)
	}

	switch v := r.Output.(type) {
	case string:
		return fmt.Sprintf("DSL[%s] -> \"%s\"", r.Code, v)
	case int, int64, float64:
		return fmt.Sprintf("DSL[%s] -> %v", r.Code, v)
	case bool:
		return fmt.Sprintf("DSL[%s] -> %t", r.Code, v)
	default:
		return fmt.Sprintf("DSL[%s] -> %v", r.Code, v)
	}
}
