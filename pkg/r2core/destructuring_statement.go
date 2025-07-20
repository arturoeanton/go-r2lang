package r2core

// ArrayDestructuring representa la desestructuración de arrays
// let [a, b, c] = [1, 2, 3]
type ArrayDestructuring struct {
	Names []string
	Value Node
}

func (ad *ArrayDestructuring) Eval(env *Environment) interface{} {
	// Evaluar la expresión del lado derecho
	value := ad.Value.Eval(env)

	// Convertir a array
	arr, ok := value.([]interface{})
	if !ok {
		panic("ArrayDestructuring: right side must be an array")
	}

	// Asignar cada elemento a su variable correspondiente
	for i, name := range ad.Names {
		if name == "_" {
			// Skip underscore assignments
			continue
		}

		var val interface{}
		if i < len(arr) {
			val = arr[i]
		} else {
			val = nil
		}

		env.Set(name, val)
	}

	return nil
}

// ObjectDestructuring representa la desestructuración de objetos
// let {name, age} = user
type ObjectDestructuring struct {
	Names []string
	Value Node
}

func (od *ObjectDestructuring) Eval(env *Environment) interface{} {
	// Evaluar la expresión del lado derecho
	value := od.Value.Eval(env)

	// Convertir a map
	obj, ok := value.(map[string]interface{})
	if !ok {
		panic("ObjectDestructuring: right side must be an object")
	}

	// Asignar cada propiedad a su variable correspondiente
	for _, name := range od.Names {
		if name == "_" {
			// Skip underscore assignments
			continue
		}

		val, exists := obj[name]
		if !exists {
			val = nil
		}

		env.Set(name, val)
	}

	return nil
}
