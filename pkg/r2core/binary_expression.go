package r2core

type BinaryExpression struct {
	Left  Node
	Op    string
	Right Node
}

func (be *BinaryExpression) Eval(env *Environment) interface{} {
	// Temporalmente deshabilitado para evitar recursión infinita
	// TODO: Implementar bytecode compilation de manera segura
	// if isBytecodeCandidate(be) {
	//     return OptimizedEval(be, env)
	// }

	lv := be.Left.Eval(env)

	// Evaluación lazy para operadores lógicos
	switch be.Op {
	case "&&":
		if !toBool(lv) {
			return false // No evaluar right si left es false
		}
		rv := be.Right.Eval(env)
		return toBool(rv)
	case "||":
		if toBool(lv) {
			return true // No evaluar right si left es true
		}
		rv := be.Right.Eval(env)
		return toBool(rv)
	default:
		// Para operadores aritméticos y de comparación, evaluar ambos
		rv := be.Right.Eval(env)
		return be.evaluateArithmeticOp(lv, rv)
	}
}

func (be *BinaryExpression) evaluateArithmeticOp(lv, rv interface{}) interface{} {
	switch be.Op {
	case "+":
		return addValues(lv, rv)
	case "-":
		return subValues(lv, rv)
	case "*":
		return mulValues(lv, rv)
	case "/":
		return divValues(lv, rv)
	case "<":
		return toFloat(lv) < toFloat(rv)
	case ">":
		return toFloat(lv) > toFloat(rv)
	case "<=":
		return toFloat(lv) <= toFloat(rv)
	case ">=":
		return toFloat(lv) >= toFloat(rv)
	case "==":
		return equals(lv, rv)
	case "!=":
		return !equals(lv, rv)
	default:
		panic("Unsupported binary operator: " + be.Op)
	}
}
