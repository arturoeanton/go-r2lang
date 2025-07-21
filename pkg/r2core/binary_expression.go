package r2core

import "fmt"

type BinaryExpression struct {
	BaseNode
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

	// Evaluación lazy para operadores lógicos y especiales
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
	case "??": // Null coalescing operator (P3)
		if lv != nil {
			return lv // Return left value if not nil
		}
		rv := be.Right.Eval(env)
		return rv // Return right value if left is nil
	case "|>": // Pipeline operator (P4)
		rv := be.Right.Eval(env)
		return be.evaluatePipeline(lv, rv, env)
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
		return be.divValuesWithPosition(lv, rv)
	case "%":
		return be.modValuesWithPosition(lv, rv)
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
	case "&":
		// Bitwise AND
		return float64(int64(toFloat(lv)) & int64(toFloat(rv)))
	case "|":
		// Bitwise OR
		return float64(int64(toFloat(lv)) | int64(toFloat(rv)))
	case "^":
		// Bitwise XOR
		return float64(int64(toFloat(lv)) ^ int64(toFloat(rv)))
	case "<<":
		// Left shift
		return float64(int64(toFloat(lv)) << uint(int64(toFloat(rv))))
	case ">>":
		// Right shift
		return float64(int64(toFloat(lv)) >> uint(int64(toFloat(rv))))
	default:
		panic("Unsupported binary operator: " + be.Op)
	}
	return 0 // This line should never be reached due to panic above
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
			PanicWithPosition(be.GetPosition(), "Division by zero")
		}
		return leftNum.Value / rightNum.Value
	case "%":
		if rightNum.Value == 0 {
			PanicWithPosition(be.GetPosition(), "Modulo by zero")
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
	case "&":
		// Bitwise AND
		return float64(int64(leftNum.Value) & int64(rightNum.Value))
	case "|":
		// Bitwise OR
		return float64(int64(leftNum.Value) | int64(rightNum.Value))
	case "^":
		// Bitwise XOR
		return float64(int64(leftNum.Value) ^ int64(rightNum.Value))
	case "<<":
		// Left shift
		return float64(int64(leftNum.Value) << uint(int64(rightNum.Value)))
	case ">>":
		// Right shift
		return float64(int64(leftNum.Value) >> uint(int64(rightNum.Value)))
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

// evaluatePipeline handles the pipeline operator |> (P4)
func (be *BinaryExpression) evaluatePipeline(leftValue, rightFunction interface{}, env *Environment) interface{} {
	switch rightFunc := rightFunction.(type) {
	case *UserFunction:
		// Call user function with left value as first argument
		return rightFunc.Call(leftValue)
	case func(...interface{}) interface{}:
		// Built-in function - call with left value
		return rightFunc(leftValue)
	case *FunctionLiteral:
		// Anonymous function - create user function and call
		userFunc := &UserFunction{
			Params: rightFunc.Params,
			Body:   rightFunc.Body,
			Env:    env,
		}
		return userFunc.Call(leftValue)
	default:
		// Try to treat as identifier and look it up in environment
		if identifier, ok := be.Right.(*Identifier); ok {
			if fnValue, exists := env.Get(identifier.Name); exists {
				switch fn := fnValue.(type) {
				case *UserFunction:
					return fn.Call(leftValue)
				case func(...interface{}) interface{}:
					return fn(leftValue)
				default:
					panic(fmt.Sprintf("Pipeline operator |>: expected function, got %s (value: %v)", typeof(fnValue), fnValue))
				}
			} else {
				panic(fmt.Sprintf("Pipeline operator |>: function '%s' not found", identifier.Name))
			}
		}
		panic(fmt.Sprintf("Pipeline operator |>: expected function on right side, got %s", typeof(be.Right)))
	}
	return nil // This line should never be reached due to panics above
}

// divValuesWithPosition performs division with position-aware error reporting
func (be *BinaryExpression) divValuesWithPosition(a, b interface{}) interface{} {
	// Fast path: avoid conversions if already float64
	if af, ok := a.(float64); ok {
		if bf, ok := b.(float64); ok {
			if bf == 0 {
				PanicWithPosition(be.GetPosition(), "Division by zero")
			}
			return af / bf
		}
	}

	den := toFloat(b)
	if den == 0 {
		PanicWithPosition(be.GetPosition(), "Division by zero")
	}
	return toFloat(a) / den
}

// modValuesWithPosition performs modulo with position-aware error reporting
func (be *BinaryExpression) modValuesWithPosition(a, b interface{}) interface{} {
	// Fast path: avoid conversions if already float64
	if af, ok := a.(float64); ok {
		if bf, ok := b.(float64); ok {
			if bf == 0 {
				PanicWithPosition(be.GetPosition(), "Modulo by zero")
			}
			return float64(int64(af) % int64(bf))
		}
	}

	den := toFloat(b)
	if den == 0 {
		PanicWithPosition(be.GetPosition(), "Modulo by zero")
	}
	return float64(int(toFloat(a)) % int(den))
}
