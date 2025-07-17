package r2core

type BinaryExpression struct {
	Left  Node
	Op    string
	Right Node
}

func (be *BinaryExpression) Eval(env *Environment) interface{} {
	// Fast-path para operaciones aritméticas simples
	if fastResult := be.tryFastArithmetic(); fastResult != nil {
		return fastResult
	}

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
	// Manejar operaciones con fechas primero
	if dateResult := be.evalDateOperations(lv, rv); dateResult != nil {
		return dateResult
	}

	switch be.Op {
	case "+":
		return addValues(lv, rv)
	case "-":
		return subValues(lv, rv)
	case "*":
		return mulValues(lv, rv)
	case "/":
		return divValues(lv, rv)
	case "%":
		return modValues(lv, rv)
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

// tryFastArithmetic intenta resolver operaciones aritméticas simples sin evaluación completa
func (be *BinaryExpression) tryFastArithmetic() interface{} {
	// Solo optimizar para literals numéricos directos
	leftNum, leftOk := be.Left.(*NumberLiteral)
	rightNum, rightOk := be.Right.(*NumberLiteral)

	if !leftOk || !rightOk {
		return nil // No es fast-path candidate
	}

	// Fast arithmetic para números literales
	switch be.Op {
	case "+":
		return leftNum.Value + rightNum.Value
	case "-":
		return leftNum.Value - rightNum.Value
	case "*":
		return leftNum.Value * rightNum.Value
	case "/":
		if rightNum.Value == 0 {
			panic("Division by zero")
		}
		return leftNum.Value / rightNum.Value
	case "%":
		if rightNum.Value == 0 {
			panic("Modulo by zero")
		}
		return float64(int(leftNum.Value) % int(rightNum.Value))
	case "<":
		return leftNum.Value < rightNum.Value
	case ">":
		return leftNum.Value > rightNum.Value
	case "<=":
		return leftNum.Value <= rightNum.Value
	case ">=":
		return leftNum.Value >= rightNum.Value
	case "==":
		return leftNum.Value == rightNum.Value
	case "!=":
		return leftNum.Value != rightNum.Value
	default:
		return nil // No es operación aritmética simple
	}
}

// evalDateOperations maneja operaciones aritméticas con fechas y duraciones
func (be *BinaryExpression) evalDateOperations(left, right interface{}) interface{} {
	leftDate, leftIsDate := left.(*DateValue)
	rightDate, rightIsDate := right.(*DateValue)
	leftDuration, leftIsDuration := left.(*DurationValue)
	rightDuration, rightIsDuration := right.(*DurationValue)

	switch be.Op {
	case "+":
		if leftIsDate && rightIsDuration {
			// Fecha + Duración = Fecha
			return leftDate.Add(rightDuration)
		}
		if leftIsDuration && rightIsDate {
			// Duración + Fecha = Fecha
			return rightDate.Add(leftDuration)
		}
		if leftIsDuration && rightIsDuration {
			// Duración + Duración = Duración
			return NewDurationValue(leftDuration.Duration + rightDuration.Duration)
		}
	case "-":
		if leftIsDate && rightIsDate {
			// Fecha - Fecha = Duración
			return leftDate.Sub(rightDate)
		}
		if leftIsDate && rightIsDuration {
			// Fecha - Duración = Fecha
			return leftDate.Sub(rightDuration)
		}
		if leftIsDuration && rightIsDuration {
			// Duración - Duración = Duración
			return NewDurationValue(leftDuration.Duration - rightDuration.Duration)
		}
	case "<":
		if leftIsDate && rightIsDate {
			return leftDate.Compare(rightDate) < 0
		}
		if leftIsDuration && rightIsDuration {
			return leftDuration.Duration < rightDuration.Duration
		}
	case "<=":
		if leftIsDate && rightIsDate {
			return leftDate.Compare(rightDate) <= 0
		}
		if leftIsDuration && rightIsDuration {
			return leftDuration.Duration <= rightDuration.Duration
		}
	case ">":
		if leftIsDate && rightIsDate {
			return leftDate.Compare(rightDate) > 0
		}
		if leftIsDuration && rightIsDuration {
			return leftDuration.Duration > rightDuration.Duration
		}
	case ">=":
		if leftIsDate && rightIsDate {
			return leftDate.Compare(rightDate) >= 0
		}
		if leftIsDuration && rightIsDuration {
			return leftDuration.Duration >= rightDuration.Duration
		}
	case "==":
		if leftIsDate && rightIsDate {
			return leftDate.Equals(rightDate)
		}
		if leftIsDuration && rightIsDuration {
			return leftDuration.Duration == rightDuration.Duration
		}
	case "!=":
		if leftIsDate && rightIsDate {
			return !leftDate.Equals(rightDate)
		}
		if leftIsDuration && rightIsDuration {
			return leftDuration.Duration != rightDuration.Duration
		}
	}

	return nil // No es operación de fecha/duración
}
