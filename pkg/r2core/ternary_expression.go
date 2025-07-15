package r2core

// TernaryExpression representa un operador ternario condition ? trueExpr : falseExpr
type TernaryExpression struct {
	Condition Node
	TrueExpr  Node
	FalseExpr Node
}

// Eval evalúa la expresión ternaria
func (te *TernaryExpression) Eval(env *Environment) interface{} {
	conditionValue := te.Condition.Eval(env)
	
	// Convertir a boolean
	if toBool(conditionValue) {
		return te.TrueExpr.Eval(env)
	} else {
		return te.FalseExpr.Eval(env)
	}
}