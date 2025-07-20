package r2core

type MapPair struct {
	Key   Node
	Value Node
}

type MapLiteral struct {
	Pairs []MapPair
}

func (ml *MapLiteral) Eval(env *Environment) interface{} {
	m := make(map[string]interface{})

	for _, pair := range ml.Pairs {
		// Verificar si el valor es un spread
		val := pair.Value.Eval(env)
		if sv, isSpread := IsSpreadValue(val); isSpread {
			// Expandir objeto spread
			switch obj := sv.Value.(type) {
			case map[string]interface{}:
				// Expandir todas las propiedades del objeto
				for k, v := range obj {
					m[k] = v
				}
			default:
				// Si no es un objeto, lo tratamos como una propiedad normal
				keyVal := pair.Key.Eval(env)
				keyStr := toString(keyVal)
				m[keyStr] = obj
			}
		} else {
			// Evaluar clave y valor normalmente
			keyVal := pair.Key.Eval(env)
			keyStr := toString(keyVal)
			m[keyStr] = val
		}
	}
	return m
}
