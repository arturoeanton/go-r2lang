package r2core

// SpreadExpression representa el operador spread ...
// Puede usarse en arrays: [...arr1, 4, 5]
// En objetos: {...obj1, c: 3}
// En función calls: func(...args)
type SpreadExpression struct {
	Value Node
}

func (se *SpreadExpression) Eval(env *Environment) interface{} {
	// El spread expression retorna una estructura especial
	// que será manejada por el contexto donde se use
	return &SpreadValue{
		Value: se.Value.Eval(env),
	}
}

// SpreadValue es un wrapper para indicar que un valor debe ser expandido
type SpreadValue struct {
	Value interface{}
}

// IsSpreadValue verifica si un valor es un SpreadValue
func IsSpreadValue(v interface{}) (*SpreadValue, bool) {
	sv, ok := v.(*SpreadValue)
	return sv, ok
}

// ExpandSpreadInArray expande valores spread en un array
func ExpandSpreadInArray(elements []interface{}) []interface{} {
	var result []interface{}

	for _, elem := range elements {
		if sv, isSpread := IsSpreadValue(elem); isSpread {
			// Expandir el valor spread
			switch val := sv.Value.(type) {
			case []interface{}:
				// Expandir array
				result = append(result, val...)
			case InterfaceSlice:
				// Expandir array
				result = append(result, []interface{}(val)...)
			case map[string]interface{}:
				// Para objetos en arrays, conservamos el objeto completo
				result = append(result, val)
			default:
				// Para tipos primitivos, agregamos el valor directamente
				result = append(result, val)
			}
		} else {
			result = append(result, elem)
		}
	}

	return result
}

// ExpandSpreadInObject expande valores spread en un objeto
//
// Nota: la evaluación real de map/object literals con spread la hace
// MapLiteral.Eval (map_literal.go), que evalúa cada pair.Value con el
// Environment correcto antes de llamar a IsSpreadValue. Esta función se
// conserva como utilidad standalone; evalúa clave/valor con un Environment
// propio (no nil) y usa toString (coerción segura, sin panics de Go por
// aserciones de tipo) para obtener la clave, igual que MapLiteral.Eval.
func ExpandSpreadInObject(pairs []MapPair) map[string]interface{} {
	result := make(map[string]interface{})
	env := NewEnvironment()

	for _, pair := range pairs {
		val := pair.Value.Eval(env)
		if sv, isSpread := IsSpreadValue(val); isSpread {
			// Expandir el valor spread
			switch obj := sv.Value.(type) {
			case map[string]interface{}:
				// Expandir todas las propiedades del objeto
				for k, v := range obj {
					result[k] = v
				}
			case map[string]*Variable:
				for k, v := range obj {
					result[k] = v.Value
				}
			default:
				// Si no es un objeto, lo tratamos como una propiedad normal
				keyStr := toString(pair.Key.Eval(env))
				result[keyStr] = obj
			}
		} else {
			// Evaluar la clave y el valor normalmente
			keyStr := toString(pair.Key.Eval(env))
			result[keyStr] = val
		}
	}

	return result
}

// ExpandSpreadInFunctionCall expande valores spread en argumentos de función
func ExpandSpreadInFunctionCall(args []interface{}) []interface{} {
	var result []interface{}

	for _, arg := range args {
		if sv, isSpread := IsSpreadValue(arg); isSpread {
			// Expandir el valor spread
			switch val := sv.Value.(type) {
			case []interface{}:
				// Expandir array como argumentos individuales
				result = append(result, val...)
			case InterfaceSlice:
				// Expandir array como argumentos individuales
				result = append(result, []interface{}(val)...)
			default:
				// Para tipos no-array, agregar como argumento individual
				result = append(result, val)
			}
		} else {
			result = append(result, arg)
		}
	}

	return result
}
