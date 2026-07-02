package r2core

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
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
	Warnings  []string // Grammar warnings from Validate(), populated at Eval time

	// useMu serializes .use() calls on this DSL: semantic actions are
	// registered once (at Eval time) as closures that read
	// currentExecutionEnv, a field shared by every call, matching go-dsl's
	// own instance is documented as unsafe for concurrent Use/Parse calls
	// (see dslbuilder.DSL.Use). Locking here gives R2Lang scripts that call
	// .use() from multiple goroutines (via `go`/`r2`) correct, race-free
	// results instead of one call observing another's environment. It is
	// reentrant so an action function that itself calls .use() on this same
	// DSL (e.g. sub-parsing) doesn't deadlock a plain sync.Mutex against
	// itself; currentExecutionEnv is still saved/restored per call in
	// evaluateDSLCode so a nested .use() can't leak its environment into
	// the outer call once it returns.
	useMu               reentrantMutex
	currentExecutionEnv *Environment // Current execution environment for actions
}

// reentrantMutex is a sync.Mutex that the same goroutine can Lock() again
// without blocking on itself (a plain sync.Mutex is not reentrant and would
// deadlock). Cross-goroutine calls still block on the underlying mutex, so
// mutual exclusion between goroutines is unaffected.
type reentrantMutex struct {
	mu     sync.Mutex
	metaMu sync.Mutex
	owner  int64
	count  int
}

func (m *reentrantMutex) Lock() {
	gid := currentGoroutineID()

	m.metaMu.Lock()
	if m.count > 0 && m.owner == gid {
		m.count++
		m.metaMu.Unlock()
		return
	}
	m.metaMu.Unlock()

	m.mu.Lock()

	m.metaMu.Lock()
	m.owner = gid
	m.count = 1
	m.metaMu.Unlock()
}

func (m *reentrantMutex) Unlock() {
	gid := currentGoroutineID()

	m.metaMu.Lock()
	defer m.metaMu.Unlock()
	if m.count == 0 || m.owner != gid {
		panic("reentrantMutex: Unlock called by a goroutine that doesn't hold the lock")
	}
	m.count--
	if m.count == 0 {
		m.mu.Unlock()
	}
}

// currentGoroutineID extracts the calling goroutine's ID from its stack
// trace header ("goroutine 123 [running]:"). It's only used to detect
// same-goroutine re-entry into reentrantMutex, never for scheduling or
// correctness beyond that, so the usual fragility concerns with parsing
// runtime.Stack output don't apply here.
func currentGoroutineID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	fields := bytes.Fields(buf[:n])
	if len(fields) < 2 {
		return -1
	}
	id, err := strconv.ParseInt(string(fields[1]), 10, 64)
	if err != nil {
		return -1
	}
	return id
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

	// Statically check the grammar: errors mean it can never parse
	// anything (undefined rules, no start rule, ...), so fail loudly right
	// away instead of surfacing a confusing error on the first .use() call.
	// Warnings (e.g. unreachable rules, unregistered actions) don't block
	// the DSL from working, so they're only reported, not fatal.
	warnings, err := dsl.Grammar.Validate()
	if err != nil {
		panic(fmt.Sprintf("DSL '%s' has an invalid grammar: %v", dslName, err))
	}
	dsl.Warnings = warnings
	if len(warnings) > 0 {
		fmt.Fprintf(os.Stderr, "DSL '%s': %d grammar warning(s):\n", dslName, len(warnings))
		for _, w := range warnings {
			fmt.Fprintf(os.Stderr, "  - %s\n", w)
		}
	}

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
		"tokens": func(args ...interface{}) interface{} {
			code, err := dslStringArg(args, 0, "tokens", "code")
			if err != nil {
				return err
			}
			matches, err := dsl.Grammar.DebugTokens(code)
			if err != nil {
				return fmt.Errorf("DSL tokens: %v", err)
			}
			result := make([]interface{}, len(matches))
			for i, m := range matches {
				result[i] = map[string]interface{}{
					"type":  m.TokenType,
					"value": m.Value,
					"start": float64(m.Start),
					"end":   float64(m.End),
				}
			}
			return result
		},
		"ast": func(args ...interface{}) interface{} {
			code, err := dslStringArg(args, 0, "ast", "code")
			if err != nil {
				return err
			}
			node, err := dsl.Grammar.AST(code)
			if err != nil {
				return fmt.Errorf("DSL ast: %v", err)
			}
			return node
		},
		"check": func(args ...interface{}) interface{} {
			code, err := dslStringArg(args, 0, "check", "code")
			if err != nil {
				return err
			}
			valid, errMsg := dsl.Grammar.Check(code)
			return map[string]interface{}{"valid": valid, "error": errMsg}
		},
		"diagnostics": func(args ...interface{}) interface{} {
			code, err := dslStringArg(args, 0, "diagnostics", "code")
			if err != nil {
				return err
			}
			return dsl.Grammar.Diagnostics(code)
		},
		"completions": func(args ...interface{}) interface{} {
			code, err := dslStringArg(args, 0, "completions", "code")
			if err != nil {
				return err
			}
			if len(args) < 2 {
				return fmt.Errorf("DSL completions: 2 arguments (code, offset) are required")
			}
			offset, ok := args[1].(float64)
			if !ok {
				return fmt.Errorf("DSL completions: second argument (offset) must be a number")
			}
			return dsl.Grammar.Completions(code, int(offset))
		},
		"grammar":   dsl.Grammar,
		"functions": dsl.Functions,
		"warnings":  toInterfaceSlice(dsl.Warnings),
	}

	// Register the DSL object in the environment
	env.Set(dslName, dslObject)

	return dslObject
}

// dslStringArg validates that args[i] is present and a string, for the
// small introspection methods (tokens/ast/check/diagnostics/completions)
// that all take a code string as their first argument.
func dslStringArg(args []interface{}, i int, method, argName string) (string, error) {
	if len(args) <= i {
		return "", fmt.Errorf("DSL %s: argument (%s) is required", method, argName)
	}
	s, ok := args[i].(string)
	if !ok {
		return "", fmt.Errorf("DSL %s: argument (%s) must be a string", method, argName)
	}
	return s, nil
}

func toInterfaceSlice(warnings []string) []interface{} {
	result := make([]interface{}, len(warnings))
	for i, w := range warnings {
		result[i] = w
	}
	return result
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
						case "keyword":
							dsl.extractKeyword(call)
						case "literal":
							dsl.extractLiteral(call)
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
				var symbols []string
				for _, alt := range alternatives.Elements {
					if altStr, ok := alt.(*StringLiteral); ok {
						trimmed := strings.Trim(altStr.Value, "\"'")
						symbols = append(symbols, strings.Fields(trimmed)...)
					}
				}
				action := ""
				if len(call.Args) > 2 {
					if actionStr, ok := call.Args[2].(*StringLiteral); ok {
						action = strings.Trim(actionStr.Value, "\"'")
					}
				}
				ruleName := strings.Trim(nameStr.Value, "\"'")
				dsl.Grammar.AddRule(ruleName, symbols, action)
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
					if err := dsl.Grammar.AddKeywordToken(name, pattern); err != nil {
						panic(fmt.Sprintf("DSL '%s': invalid keyword token('%s', %q): %v", dsl.Name.Name, name, pattern, err))
					}
				} else {
					// Improve pattern by auto-escaping problematic single characters if needed
					improvedPattern := dsl.improveRegexPattern(pattern)
					err := dsl.Grammar.AddToken(name, improvedPattern)
					if err != nil {
						// If the improved pattern fails, try the original. Its
						// error, not the improved pattern's, is what's actually
						// informative to report if this also fails: silently
						// dropping the token here would otherwise leave it
						// unregistered with no trace, surfacing later (if at
						// all) as a confusing "unknown symbol" error from
						// Validate() instead of the real cause.
						if origErr := dsl.Grammar.AddToken(name, pattern); origErr != nil {
							panic(fmt.Sprintf("DSL '%s': invalid token('%s', %q): %v", dsl.Name.Name, name, pattern, origErr))
						}
					}
				}
			}
		}
	}
}

// extractKeyword handles keyword("NAME", "text") — an explicit alternative
// to token() for keywords, so grammars don't have to rely on
// isKeywordToken's auto-detection.
func (dsl *DSLDefinition) extractKeyword(call *CallExpression) {
	if len(call.Args) >= 2 {
		if nameStr, ok := call.Args[0].(*StringLiteral); ok {
			if wordStr, ok := call.Args[1].(*StringLiteral); ok {
				name := strings.Trim(nameStr.Value, "\"'")
				word := strings.Trim(wordStr.Value, "\"'")
				if err := dsl.Grammar.AddKeywordToken(name, word); err != nil {
					panic(fmt.Sprintf("DSL '%s': invalid keyword('%s', %q): %v", dsl.Name.Name, name, word, err))
				}
			}
		}
	}
}

// extractLiteral handles literal("NAME", "text") — an explicit alternative
// to token() for exact-text tokens like operators/punctuation, which
// isKeywordToken always routes through regex handling.
func (dsl *DSLDefinition) extractLiteral(call *CallExpression) {
	if len(call.Args) >= 2 {
		if nameStr, ok := call.Args[0].(*StringLiteral); ok {
			if textStr, ok := call.Args[1].(*StringLiteral); ok {
				name := strings.Trim(nameStr.Value, "\"'")
				text := strings.Trim(textStr.Value, "\"'")
				if err := dsl.Grammar.AddLiteral(name, text); err != nil {
					panic(fmt.Sprintf("DSL '%s': invalid literal('%s', %q): %v", dsl.Name.Name, name, text, err))
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
	dsl.useMu.Lock()
	defer dsl.useMu.Unlock()

	// Make sure we have the global environment
	if dsl.GlobalEnv == nil {
		dsl.GlobalEnv = env
	}

	// Store the execution environment for actions to use. Saved/restored
	// (rather than just set) so a re-entrant .use() call from within an
	// action (allowed by useMu's reentrancy) doesn't leave the outer call's
	// remaining actions running against the inner call's environment.
	savedExecutionEnv := dsl.currentExecutionEnv
	dsl.currentExecutionEnv = env
	defer func() { dsl.currentExecutionEnv = savedExecutionEnv }()

	// Actions run as R2Lang function bodies against currentExecutionEnv, so
	// the context map is threaded through the environment; pass it to go-dsl
	// as well so it's also reachable via GetContext() from action code.
	var context map[string]interface{}
	if ctxVal, exists := env.Get("context"); exists {
		if ctxMap, ok := ctxVal.(map[string]interface{}); ok {
			context = ctxMap
		}
	}

	result, err := dsl.Grammar.Use(code, context)
	if err != nil {
		return fmt.Errorf("DSL parsing error: %v", err)
	}

	// Unwrap a ReturnValue-wrapped output, e.g. from an action whose R2Lang
	// function body ends in an explicit `return`.
	if retVal, ok := result.Output.(*ReturnValue); ok {
		result.Output = retVal.Value
	} else if retVal, ok := result.Output.(ReturnValue); ok {
		result.Output = retVal.Value
	}

	return result
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

	// DSL actions run arbitrary R2Lang function bodies (including, via
	// .use(), themselves) but were invoked outside UserFunction.Call, so
	// they bypassed the same recursion-depth/timeout guards normal calls
	// get — infinite recursion through a DSL action hung the interpreter
	// forever instead of raising the usual clean "Loop infinito" error.
	limiter := fnEnv.GetLimiter()
	if limiter.Enabled {
		if limiter.CheckRecursionDepth() {
			panic(NewRecursionError("max_depth", limiter.CallDepth()))
		}
		if limiter.CheckTimeLimit() {
			panic(NewTimeoutError("function_timeout", fnEnv.GetContext()))
		}
		if limiter.CheckContext() {
			panic(NewTimeoutError("function_context_canceled", fnEnv.GetContext()))
		}
		limiter.EnterFunction(fn.Name)
		defer limiter.ExitFunction()
	}

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
