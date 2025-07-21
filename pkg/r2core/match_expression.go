package r2core

import (
	"reflect"
)

// MatchExpression represents a match statement (P3)
type MatchExpression struct {
	Value Node
	Cases []MatchCase
}

// MatchCase represents a case in a match expression
type MatchCase struct {
	Pattern   Pattern
	Guard     Node // Optional guard condition (if clause)
	Body      Node
	IsDefault bool // true for case _ (wildcard)
}

// Pattern interface for different types of patterns
type Pattern interface {
	MatchValue(value interface{}, env *Environment) (bool, map[string]interface{})
}

// LiteralPattern matches exact values
type LiteralPattern struct {
	Value Node
}

func (lp *LiteralPattern) MatchValue(value interface{}, env *Environment) (bool, map[string]interface{}) {
	expected := lp.Value.Eval(env)
	matches := isEqual(value, expected)
	return matches, make(map[string]interface{})
}

// VariablePattern binds value to a variable
type VariablePattern struct {
	Name string
}

func (vp *VariablePattern) MatchValue(value interface{}, env *Environment) (bool, map[string]interface{}) {
	bindings := make(map[string]interface{})
	bindings[vp.Name] = value
	return true, bindings
}

// WildcardPattern matches anything (underscore _)
type WildcardPattern struct{}

func (wp *WildcardPattern) MatchValue(value interface{}, env *Environment) (bool, map[string]interface{}) {
	return true, make(map[string]interface{})
}

// ArrayPattern matches array structures
type ArrayPattern struct {
	Elements []Pattern
}

func (ap *ArrayPattern) MatchValue(value interface{}, env *Environment) (bool, map[string]interface{}) {
	arr, ok := value.([]interface{})
	if !ok {
		return false, nil
	}

	if len(arr) != len(ap.Elements) {
		return false, nil
	}

	bindings := make(map[string]interface{})
	for i, pattern := range ap.Elements {
		matches, elementBindings := pattern.MatchValue(arr[i], env)
		if !matches {
			return false, nil
		}
		// Merge bindings
		for name, val := range elementBindings {
			bindings[name] = val
		}
	}

	return true, bindings
}

// ObjectPattern matches object structures
type ObjectPattern struct {
	Fields map[string]Pattern
}

func (op *ObjectPattern) MatchValue(value interface{}, env *Environment) (bool, map[string]interface{}) {
	obj, ok := value.(map[string]interface{})
	if !ok {
		return false, nil
	}

	bindings := make(map[string]interface{})
	for fieldName, pattern := range op.Fields {
		fieldValue, exists := obj[fieldName]
		if !exists {
			return false, nil
		}

		matches, fieldBindings := pattern.MatchValue(fieldValue, env)
		if !matches {
			return false, nil
		}

		// Merge bindings
		for name, val := range fieldBindings {
			bindings[name] = val
		}
	}

	return true, bindings
}

// OrPattern matches any of several patterns (pattern1 | pattern2)
type OrPattern struct {
	Patterns []Pattern
}

func (orp *OrPattern) MatchValue(value interface{}, env *Environment) (bool, map[string]interface{}) {
	for _, pattern := range orp.Patterns {
		matches, bindings := pattern.MatchValue(value, env)
		if matches {
			return true, bindings
		}
	}
	return false, nil
}

// GuardedPattern adds a guard condition to a pattern
type GuardedPattern struct {
	Pattern Pattern
	Guard   Node
}

func (gp *GuardedPattern) MatchValue(value interface{}, env *Environment) (bool, map[string]interface{}) {
	matches, bindings := gp.Pattern.MatchValue(value, env)
	if !matches {
		return false, nil
	}

	// Create new environment with bindings for guard evaluation
	guardEnv := NewInnerEnv(env)
	for name, val := range bindings {
		guardEnv.Set(name, val)
	}

	guardResult := gp.Guard.Eval(guardEnv)
	if toBool(guardResult) {
		return true, bindings
	}

	return false, nil
}

// Eval evaluates the match expression
func (me *MatchExpression) Eval(env *Environment) interface{} {
	value := me.Value.Eval(env)

	for _, matchCase := range me.Cases {
		matches, bindings := matchCase.Pattern.MatchValue(value, env)
		if matches {
			// Create new environment with pattern bindings
			caseEnv := NewInnerEnv(env)
			for name, val := range bindings {
				caseEnv.Set(name, val)
			}

			// Evaluate guard if present
			if matchCase.Guard != nil && !toBool(matchCase.Guard.Eval(caseEnv)) {
				continue
			}

			// Execute case body
			return matchCase.Body.Eval(caseEnv)
		}
	}

	panic("No matching case found in match expression")
}

// isEqual compares two values for equality
func isEqual(a, b interface{}) bool {
	// Handle nil cases
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Use reflect for deep comparison
	return reflect.DeepEqual(a, b)
}
