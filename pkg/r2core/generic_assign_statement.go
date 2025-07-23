package r2core

type GenericAssignStatement struct {
	Left  Node
	Right Node
}

func (gas *GenericAssignStatement) Eval(env *Environment) interface{} {

	val := gas.Right.Eval(env)
	switch left := gas.Left.(type) {
	case *Identifier:
		env.Update(left.Name, val)
		return val
	case *AccessExpression:
		objVal := left.Object.Eval(env)
		switch obj := objVal.(type) {
		case *ObjectInstance:
			obj.Env.Set(left.Member, val)
			return val
		case map[string]interface{}:
			obj[left.Member] = val
			return val
		default:
			panic("Cannot assign to property of non-object type")
		}
		return val
	case *IndexExpression:
		return assignIndexExpression(left, val, env)
	default:
		panic("Cannot assign to this expression")
	}
}
