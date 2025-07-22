package r2core

import (
	"fmt"
	"strings"
)

// DSLDefinition represents a DSL definition block
type DSLDefinition struct {
	Token     Token
	Name      *Identifier
	Body      *BlockStatement
	Grammar   *DSLGrammar
	Functions map[string]*FunctionDeclaration
	IsActive  bool
	GlobalEnv *Environment
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
			env.Set("context", context)
			return dsl.evaluateDSLCode(code, env)
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
					return dsl.callDSLFunction(node, args, env)
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
				dsl.Grammar.AddToken(name, pattern)
			}
		}
	}
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

	// Create parser for this DSL
	parser := NewDSLParser(dsl.Grammar)

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

	// Create new environment for function execution that inherits from global environment
	fnEnv := NewInnerEnv(dsl.GlobalEnv)

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

	// Execute function body
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
