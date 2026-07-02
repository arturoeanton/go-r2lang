package r2core

import (
	"fmt"
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
