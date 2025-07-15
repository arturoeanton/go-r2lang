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
		instance, ok := objVal.(*ObjectInstance)
		if !ok {
			panic("Closing quote of string expected")
		}
		instance.Env.Set(left.Member, val)
		return val
	case *IndexExpression:
		return assignIndexExpression(left, val, env)
	default:
		panic("Cannot assign to this expression")
	}
}
