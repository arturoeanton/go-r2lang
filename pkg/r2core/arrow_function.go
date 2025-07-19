package r2core

// ArrowFunction represents arrow function expressions like x => x * 2 or (a, b) => a + b
type ArrowFunction struct {
	Params       []Parameter // Function parameters
	Body         Node        // Either an expression or a block statement
	IsExpression bool        // true if body is expression, false if block
}

func (af *ArrowFunction) Eval(env *Environment) interface{} {
	// Convert parameters to args for backward compatibility
	var args []string
	for _, param := range af.Params {
		args = append(args, param.Name)
	}

	// Create UserFunction with lexical scoping
	fn := &UserFunction{
		Args:     args,
		Params:   af.Params,
		Body:     af.createBody(),
		Env:      env, // Lexical scoping - capture current environment
		IsMethod: false,
		code:     "arrow_function",
	}
	return fn
}

// createBody converts the arrow function body to a BlockStatement
func (af *ArrowFunction) createBody() *BlockStatement {
	if af.IsExpression {
		// For expression bodies, wrap in a return statement
		returnStmt := &ReturnStatement{Value: af.Body}
		return &BlockStatement{Statements: []Node{returnStmt}}
	} else {
		// For block bodies, assume it's already a BlockStatement
		if blockStmt, ok := af.Body.(*BlockStatement); ok {
			return blockStmt
		}
		// Fallback: wrap in block
		return &BlockStatement{Statements: []Node{af.Body}}
	}
}
