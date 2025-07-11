package r2core

type ObjectDeclaration struct {
	Name       string
	ParentName string
	Members    []Node
}

func (od *ObjectDeclaration) Eval(env *Environment) interface{} {
	blueprint := make(map[string]interface{})
	if od.ParentName != "" {
		if parent, ok := env.Get(od.ParentName); ok {
			blueprint["super"] = parent
		}
		blueprint["SuperClassName"] = blueprint["ClassName"]
		raw, _ := env.Get(od.ParentName)
		if props, ok := raw.(map[string]interface{}); ok {

			for k, v := range props {
				if k == "ClassName" || k == "SuperClassName" || k == "super" {
					continue
				}
				blueprint[k] = v
			}
		}

	}

	blueprint["ClassName"] = od.Name
	for _, m := range od.Members {
		switch node := m.(type) {
		case *LetStatement:
			if node.Name == "super" || node.Name == "ClassName" || node.Name == "SuperClassName" {
				panic("Cannot redefine 'super'")
			}
			blueprint[node.Name] = nil
		case *FunctionDeclaration:
			if node.Name == "super" || node.Name == "ClassName" || node.Name == "SuperClassName" {
				panic("Cannot redefine 'super'")
			}
			fn := &UserFunction{
				Args:     node.Args,
				Body:     node.Body,
				Env:      nil,
				IsMethod: true,
			}
			blueprint[node.Name] = fn
		}
	}
	env.Set(od.Name, blueprint)
	return nil
}
