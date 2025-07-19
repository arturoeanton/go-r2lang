package r2core

import "fmt"

type IndexExpression struct {
	Left  Node
	Index Node
}

func (ie *IndexExpression) Eval(env *Environment) interface{} {
	leftVal := ie.Left.Eval(env)
	indexVal := ie.Index.Eval(env)

	switch container := leftVal.(type) {
	case map[string]interface{}:
		strKey, ok := indexVal.(string)
		if !ok {
			panic("index must be a string for map")
		}
		vv, found := container[strKey]
		if !found {
			return nil
		}
		return vv
	case []interface{}:
		fIndex, ok := indexVal.(float64)
		if !ok {
			panic("index must be numeric for array")
		}
		idx := int(fIndex)
		if idx < 0 {

			idx = (len(container) + idx)
		}
		if idx < 0 || idx >= len(container) {
			panic(fmt.Sprintf("index out of range: %d len of array %d", idx, len(container)))
		}
		return container[idx]
	case interfaceSlice:
		fIndex, ok := indexVal.(float64)
		if !ok {
			panic("index must be numeric for array")
		}
		idx := int(fIndex)
		if idx < 0 {

			idx = (len(container) + idx)
		}
		if idx < 0 || idx >= len(container) {
			panic(fmt.Sprintf("index out of range: %d len of array %d", idx, len(container)))
		}
		return container[idx]
	default:
		panic("index on something that is neither map nor array")
	}
}
