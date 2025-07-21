package r2core

// ArrayComprehension represents array comprehension syntax: [expression for variable in iterable if condition]
type ArrayComprehension struct {
	Expression Node        // The expression to evaluate for each element
	Generators []Generator // for clauses (can be multiple)
	Conditions []Node      // if clauses (optional filters)
}

// ObjectComprehension represents object comprehension syntax: {key: value for variable in iterable if condition}
type ObjectComprehension struct {
	KeyExpr    Node        // Key expression
	ValueExpr  Node        // Value expression
	Generators []Generator // for clauses
	Conditions []Node      // if clauses
}

// Generator represents a 'for variable in iterable' clause
type Generator struct {
	Variable string // Loop variable name
	Iterator Node   // Expression that yields iterable
}

// Eval evaluates array comprehension
func (ac *ArrayComprehension) Eval(env *Environment) interface{} {
	return ac.generateElements(env, 0, make(map[string]interface{}))
}

// generateElements recursively generates elements for array comprehension
func (ac *ArrayComprehension) generateElements(env *Environment, genIndex int, bindings map[string]interface{}) []interface{} {
	if genIndex >= len(ac.Generators) {
		// All generators processed, now check conditions and generate element
		newEnv := NewInnerEnv(env)
		for name, val := range bindings {
			newEnv.Set(name, val)
		}

		// Check all conditions
		for _, condition := range ac.Conditions {
			if !toBool(condition.Eval(newEnv)) {
				return []interface{}{} // Empty result if condition fails
			}
		}

		// Generate element
		result := ac.Expression.Eval(newEnv)
		return []interface{}{result}
	}

	// Process current generator
	generator := ac.Generators[genIndex]
	iterable := generator.Iterator.Eval(env)

	var results []interface{}

	// Iterate over the iterable
	switch iter := iterable.(type) {
	case []interface{}:
		for _, item := range iter {
			newBindings := make(map[string]interface{})
			// Copy existing bindings
			for k, v := range bindings {
				newBindings[k] = v
			}
			// Add current binding
			newBindings[generator.Variable] = item

			// Recursively process remaining generators
			subResults := ac.generateElements(env, genIndex+1, newBindings)
			results = append(results, subResults...)
		}
	case map[string]interface{}:
		for _, value := range iter {
			newBindings := make(map[string]interface{})
			for k, v := range bindings {
				newBindings[k] = v
			}
			newBindings[generator.Variable] = value

			subResults := ac.generateElements(env, genIndex+1, newBindings)
			results = append(results, subResults...)
		}
	default:
		panic("Cannot iterate over non-iterable value in comprehension")
	}

	return results
}

// Eval evaluates object comprehension
func (oc *ObjectComprehension) Eval(env *Environment) interface{} {
	result := make(map[string]interface{})
	pairs := oc.generatePairs(env, 0, make(map[string]interface{}))

	for _, pair := range pairs {
		if keyVal, ok := pair.([]interface{}); ok && len(keyVal) == 2 {
			key := toString(keyVal[0])
			value := keyVal[1]
			result[key] = value
		}
	}

	return result
}

// generatePairs recursively generates key-value pairs for object comprehension
func (oc *ObjectComprehension) generatePairs(env *Environment, genIndex int, bindings map[string]interface{}) []interface{} {
	if genIndex >= len(oc.Generators) {
		// All generators processed, check conditions and generate pair
		newEnv := NewInnerEnv(env)
		for name, val := range bindings {
			newEnv.Set(name, val)
		}

		// Check all conditions
		for _, condition := range oc.Conditions {
			if !toBool(condition.Eval(newEnv)) {
				return []interface{}{} // Empty result if condition fails
			}
		}

		// Generate key-value pair
		key := oc.KeyExpr.Eval(newEnv)
		value := oc.ValueExpr.Eval(newEnv)
		return []interface{}{[]interface{}{key, value}}
	}

	// Process current generator
	generator := oc.Generators[genIndex]
	iterable := generator.Iterator.Eval(env)

	var results []interface{}

	// Iterate over the iterable
	switch iter := iterable.(type) {
	case []interface{}:
		for _, item := range iter {
			newBindings := make(map[string]interface{})
			for k, v := range bindings {
				newBindings[k] = v
			}
			newBindings[generator.Variable] = item

			subResults := oc.generatePairs(env, genIndex+1, newBindings)
			results = append(results, subResults...)
		}
	case map[string]interface{}:
		for _, value := range iter {
			newBindings := make(map[string]interface{})
			for k, v := range bindings {
				newBindings[k] = v
			}
			newBindings[generator.Variable] = value

			subResults := oc.generatePairs(env, genIndex+1, newBindings)
			results = append(results, subResults...)
		}
	default:
		panic("Cannot iterate over non-iterable value in comprehension")
	}

	return results
}
