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
		keyVal := pair.Key.Eval(env)
		
		// Convertir la clave a string
		var keyStr string
		switch k := keyVal.(type) {
		case string:
			keyStr = k
		case float64:
			keyStr = toString(k)
		default:
			keyStr = toString(k)
		}
		
		m[keyStr] = pair.Value.Eval(env)
	}
	return m
}
