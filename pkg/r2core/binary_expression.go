package r2core

type BinaryExpression struct {
	Left  Node
	Op    string
	Right Node
}

func (be *BinaryExpression) Eval(env *Environment) interface{} {
	lv := be.Left.Eval(env)
	rv := be.Right.Eval(env)

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
