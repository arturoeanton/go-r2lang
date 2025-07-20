package r2core

// UnaryExpression represents unary operations like !expr, -expr, +expr
type UnaryExpression struct {
	Operator string
	Right    Node
}

func (ue *UnaryExpression) Eval(env *Environment) interface{} {
	right := ue.Right.Eval(env)

	switch ue.Operator {
	case "!":
		return !isTruthy(right)
	case "-":
		// Handle unary minus
		switch val := right.(type) {
		case float64:
			return -val
		case int:
			return -float64(val)
		default:
			panic("Invalid operand for unary minus: " + toString(val))
		}
	case "+":
		// Handle unary plus (convert to number)
		switch val := right.(type) {
		case float64:
			return val
		case int:
			return float64(val)
		case string:
			// Try to convert string to number using existing toFloat function
			return toFloat(val)
		default:
			panic("Invalid operand for unary plus: " + toString(val))
		}
	case "~":
		// Handle bitwise NOT
		switch val := right.(type) {
		case float64:
			return float64(^int64(val))
		case int:
			return float64(^int64(val))
		default:
			// Try to convert to number first
			num := toFloat(val)
			return float64(^int64(num))
		}
	default:
		panic("Unknown unary operator: " + ue.Operator)
	}
}

// isTruthy determines the truthiness of a value following JavaScript-like rules
func isTruthy(obj interface{}) bool {
	switch obj := obj.(type) {
	case bool:
		return obj
	case nil:
		return false
	case int:
		return obj != 0
	case float64:
		return obj != 0.0
	case string:
		return obj != ""
	case []interface{}:
		return len(obj) > 0
	case map[string]interface{}:
		return len(obj) > 0
	default:
		return true // Objects are truthy
	}
}
