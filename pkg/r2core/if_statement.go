package r2core

type IfStatement struct {
	Condition   Node
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ifs *IfStatement) Eval(env *Environment) interface{} {
	condVal := ifs.Condition.Eval(env)
	if toBool(condVal) {
		return ifs.Consequence.Eval(env)
	} else if ifs.Alternative != nil {
		return ifs.Alternative.Eval(env)
	}
	return nil
}
