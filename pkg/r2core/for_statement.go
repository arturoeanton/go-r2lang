package r2core

type ForStatement struct {
	Init      Node
	Condition Node
	Post      Node
	Body      *BlockStatement
	inFlag    bool
	inArray   string
	//inMap       string
	inIndexName string
}

func (fs *ForStatement) Eval(env *Environment) interface{} {
	newEnv := NewInnerEnv(env)

	var result interface{}

	var arr []interface{}
	var mapVal map[string]interface{}
	var ok bool
	flagArr := true
	if fs.inFlag {
		var raw interface{}
		if _, ok = fs.Init.(*CallExpression); ok {
			raw = fs.Init.Eval(newEnv)
			newEnv.Set("$c", raw)
		} else {
			raw, _ = newEnv.Get(fs.inArray)
			newEnv.Set("$c", raw)
		}

		arr, ok = raw.([]interface{})
		if !ok {
			flagArr = false
			mapVal, ok = raw.(map[string]interface{})
			if !ok {
				panic("Not an array or map for ‘for’")
			}
		}
	}
	if fs.inFlag {
		if flagArr {
			for i, v := range arr {
				newEnv.Set(fs.inIndexName, float64(i))
				newEnv.Set("$k", float64(i))
				newEnv.Set("$v", v)
				val := fs.Body.Eval(newEnv)
				if rv, ok := val.(ReturnValue); ok {
					return rv
				}
				result = val
				if fs.Post != nil {
					fs.Post.Eval(newEnv)
				}
			}
		} else {
			for k, v := range mapVal {
				newEnv.Set(fs.inIndexName, k)
				newEnv.Set("$k", k)
				newEnv.Set("$v", v)
				val := fs.Body.Eval(newEnv)
				if rv, ok := val.(ReturnValue); ok {
					return rv
				}
				result = val
				if fs.Post != nil {
					fs.Post.Eval(newEnv)
				}
			}
		}
		return result
	}

	if fs.Init != nil {
		fs.Init.Eval(newEnv)
	}

	for {
		condVal := fs.Condition.Eval(newEnv)
		if !toBool(condVal) {
			break
		}
		val := fs.Body.Eval(newEnv)
		if rv, ok := val.(ReturnValue); ok {
			return rv
		}
		result = val
		if fs.Post != nil {
			fs.Post.Eval(newEnv)
		}
	}
	return result
}
