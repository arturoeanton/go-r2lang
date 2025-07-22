package r2core

import (
	"fmt"
	"strings"
)

// DSLDefinition represents a DSL definition block
type DSLDefinition struct {
	Token               Token
	Name                *Identifier
	Body                *BlockStatement
	Grammar             *DSLGrammar
	Functions           map[string]*FunctionDeclaration
	IsActive            bool
	GlobalEnv           *Environment
	currentExecutionEnv *Environment // Current execution environment for actions
}

func (dsl *DSLDefinition) Eval(env *Environment) interface{} {
	// Create a new DSL environment
	dslEnv := NewInnerEnv(env)

	// Get DSL name for later registration
	dslName := dsl.Name.Name

	// Initialize grammar
	dsl.Grammar = NewDSLGrammar()

	// Evaluate the body to collect grammar rules and actions
	dsl.collectGrammarDefinitions(dslEnv)

	// Create a DSL object with a 'use' method
	dslObject := map[string]interface{}{
		"use": func(args ...interface{}) interface{} {
			var code string
			var context map[string]interface{}

			if len(args) == 0 {
				return fmt.Errorf("DSL use: at least one argument (code) is required")
			}

			// First argument is always the code string
			if codeStr, ok := args[0].(string); ok {
				code = codeStr
			} else {
				return fmt.Errorf("DSL use: first argument must be a string")
			}

			// Second argument (optional) is the context map
			if len(args) > 1 {
				if ctx, ok := args[1].(map[string]interface{}); ok {
					context = ctx
				} else {
					return fmt.Errorf("DSL use: second argument must be a map")
				}
			}

			if context == nil {
				context = make(map[string]interface{})
			}

			// Create a new inner environment for this DSL execution
			// This ensures isolation and prevents state pollution
			execEnv := NewInnerEnv(env)
			execEnv.Set("context", context)
			return dsl.evaluateDSLCode(code, execEnv)
		},
		"grammar":   dsl.Grammar,
		"functions": dsl.Functions,
	}

	// Register the DSL object in the environment
	env.Set(dslName, dslObject)

	return dslObject
}

func (dsl *DSLDefinition) collectGrammarDefinitions(env *Environment) {
	if dsl.Body != nil {
		for _, stmt := range dsl.Body.Statements {
			switch node := stmt.(type) {
			case *ExprStatement:
				// Handle grammar definitions
				if call, ok := node.Expr.(*CallExpression); ok {
					if id, ok := call.Callee.(*Identifier); ok {
						switch id.Name {
						case "rule":
							dsl.extractRule(call)
						case "token":
							dsl.extractToken(call)
						case "action":
							dsl.extractAction(call, env)
						}
					}
				}
			case *FunctionDeclaration:
				// Store function declarations as semantic actions
				if dsl.Functions == nil {
					dsl.Functions = make(map[string]*FunctionDeclaration)
				}
				dsl.Functions[node.Name] = node

				// Also add as grammar action
				dsl.Grammar.AddAction(node.Name, func(args []interface{}) interface{} {
					// Use current execution environment, or fall back to global environment
					execEnv := dsl.currentExecutionEnv
					if execEnv == nil {
						execEnv = dsl.GlobalEnv
					}
					return dsl.callDSLFunction(node, args, execEnv)
				})
			}
		}
	}
}

func (dsl *DSLDefinition) extractRule(call *CallExpression) {
	if len(call.Args) >= 2 {
		if nameStr, ok := call.Args[0].(*StringLiteral); ok {
			if alternatives, ok := call.Args[1].(*ArrayLiteral); ok {
				var altStrings []string
				for _, alt := range alternatives.Elements {
					if altStr, ok := alt.(*StringLiteral); ok {
						altStrings = append(altStrings, strings.Trim(altStr.Value, "\"'"))
					}
				}
				action := ""
				if len(call.Args) > 2 {
					if actionStr, ok := call.Args[2].(*StringLiteral); ok {
						action = strings.Trim(actionStr.Value, "\"'")
					}
				}
				ruleName := strings.Trim(nameStr.Value, "\"'")

				// Join the alternatives into a single sequence
				sequence := strings.Join(altStrings, " ")
				dsl.Grammar.AddRule(ruleName, []string{sequence}, action)
			}
		}
	}
}

func (dsl *DSLDefinition) extractToken(call *CallExpression) {
	if len(call.Args) >= 2 {
		if nameStr, ok := call.Args[0].(*StringLiteral); ok {
			if patternStr, ok := call.Args[1].(*StringLiteral); ok {
				name := strings.Trim(nameStr.Value, "\"'")
				pattern := strings.Trim(patternStr.Value, "\"'")

				// Auto-detect if this is a keyword token (literal string without regex patterns)
				if dsl.isKeywordToken(pattern) {
					dsl.Grammar.AddKeywordToken(name, pattern)
				} else {
					// Improve pattern by auto-escaping problematic single characters if needed
					improvedPattern := dsl.improveRegexPattern(pattern)
					err := dsl.Grammar.AddToken(name, improvedPattern)
					if err != nil {
						// If the improved pattern fails, try the original
						dsl.Grammar.AddToken(name, pattern)
					}
				}
			}
		}
	}
}

// isKeywordToken detects if a pattern is a simple keyword (no regex metacharacters)
func (dsl *DSLDefinition) isKeywordToken(pattern string) bool {
	// Check if pattern contains regex metacharacters
	regexMetachars := []string{"[", "]", "(", ")", "*", "+", "?", ".", "^", "$", "|", "\\", "{", "}"}
	for _, meta := range regexMetachars {
		if strings.Contains(pattern, meta) {
			return false
		}
	}

	// Single character operators should be treated as regex patterns, not keywords
	if len(pattern) == 1 {
		singleCharOperators := []string{"-", "/", "*", "+", "^", "$", ".", "|", "?", "(", ")", "[", "]", "{", "}"}
		for _, op := range singleCharOperators {
			if pattern == op {
				return false // Force regex pattern handling
			}
		}
	}

	// Simple string with only letters, numbers, and basic punctuation = keyword
	return len(pattern) > 0 && pattern != ""
}

// improveRegexPattern improves regex patterns by escaping problematic single characters
func (dsl *DSLDefinition) improveRegexPattern(pattern string) string {
	// Handle single-character operators that can be problematic in regex
	if len(pattern) == 1 {
		switch pattern {
		case "-":
			return "\\-"
		case "/":
			return "\\/"
		case "^":
			return "\\^"
		case "$":
			return "\\$"
		case ".":
			return "\\."
		case "|":
			return "\\|"
		case "?":
			return "\\?"
		case "*":
			return "\\*"
		case "+":
			return "\\+"
		case "(":
			return "\\("
		case ")":
			return "\\)"
		case "[":
			return "\\["
		case "]":
			return "\\]"
		case "{":
			return "\\{"
		case "}":
			return "\\}"
		}
	}

	// For already escaped patterns, don't double-escape
	if strings.HasPrefix(pattern, "\\") && len(pattern) == 2 {
		return pattern
	}

	return pattern
}

func (dsl *DSLDefinition) extractAction(call *CallExpression, env *Environment) {
	if len(call.Args) >= 2 {
		if nameStr, ok := call.Args[0].(*StringLiteral); ok {
			if fn, ok := call.Args[1].(*FunctionDeclaration); ok {
				name := strings.Trim(nameStr.Value, "\"'")
				dsl.Grammar.AddAction(name, func(args []interface{}) interface{} {
					return dsl.callDSLFunction(fn, args, env)
				})
			}
		}
	}
}

func (dsl *DSLDefinition) evaluateDSLCode(code string, env *Environment) interface{} {
	// Make sure we have the global environment
	if dsl.GlobalEnv == nil {
		dsl.GlobalEnv = env
	}

	// Store the execution environment for actions to use
	dsl.currentExecutionEnv = env

	// Create parser for this DSL - always create a new one to ensure clean state
	parser := NewDSLParserWithContext(dsl.Grammar, dsl)

	// Parse the DSL code
	ast, err := parser.Parse(code)
	if err != nil {
		return fmt.Errorf("DSL parsing error: %v", err)
	}

	// Extract the final result from the AST
	var finalResult interface{}
	if retVal, ok := ast.(*ReturnValue); ok {
		finalResult = retVal.Value
	} else if retVal, ok := ast.(ReturnValue); ok {
		finalResult = retVal.Value
	} else {
		finalResult = ast
	}

	// Return the parsed AST with the final result
	return &DSLResult{
		AST:    ast,
		Code:   code,
		Output: finalResult,
	}
}

func (dsl *DSLDefinition) callDSLFunction(fn *FunctionDeclaration, args []interface{}, env *Environment) interface{} {
	// Use the current execution environment if available, otherwise use the provided environment
	baseEnv := env
	if dsl.currentExecutionEnv != nil {
		baseEnv = dsl.currentExecutionEnv
	} else if dsl.GlobalEnv != nil {
		baseEnv = dsl.GlobalEnv
	}

	// Create new environment for function execution
	fnEnv := NewInnerEnv(baseEnv)

	// Bind parameters to arguments
	for i, param := range fn.Args {
		if i < len(args) {
			// If the argument is a ReturnValue, extract its value
			var argValue interface{}
			if retVal, ok := args[i].(*ReturnValue); ok {
				argValue = retVal.Value
			} else if retVal, ok := args[i].(ReturnValue); ok {
				argValue = retVal.Value
			} else {
				argValue = args[i]
			}
			fnEnv.Set(param, argValue)
		} else {
			fnEnv.Set(param, nil)
		}
	}

	// Execute function body with proper error handling
	result := fn.Body.Eval(fnEnv)

	// Handle return values
	if retVal, ok := result.(*ReturnValue); ok {
		return retVal.Value
	}

	return result
}

func (dsl *DSLDefinition) String() string {
	return fmt.Sprintf("DSL(%s)", dsl.Name.Name)
}
