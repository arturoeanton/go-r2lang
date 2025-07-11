package r2core

type ObjectDeclaration struct {
	Name       string
	ParentName string
	Members    []Node
}

func (od *ObjectDeclaration) Eval(env *Environment) interface{} {
	blueprint := make(map[string]interface{})
	od.setupInheritance(blueprint, env)
	od.addMembers(blueprint)
	env.Set(od.Name, blueprint)
	return nil
}

func (od *ObjectDeclaration) setupInheritance(blueprint map[string]interface{}, env *Environment) {
	if od.ParentName == "" {
		return
	}
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

func (od *ObjectDeclaration) addMembers(blueprint map[string]interface{}) {
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
}
