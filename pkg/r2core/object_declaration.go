package r2core

type ObjectDeclaration struct {
	Name       string
	ParentName string
	Members    []Node
}

func (od *ObjectDeclaration) Eval(env *Environment) interface{} {
	blueprint := make(map[string]interface{})
	od.setupInheritance(blueprint, env)
	od.addMembers(blueprint, env)
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

func (od *ObjectDeclaration) addMembers(blueprint map[string]interface{}, env *Environment) {
	blueprint["ClassName"] = od.Name
	for _, m := range od.Members {
		switch node := m.(type) {
		case *LetStatement:
			if node.Name == "super" || node.Name == "ClassName" || node.Name == "SuperClassName" {
				panic("Cannot redefine 'super'")
			}
			// A class field's declared default (e.g. "let value = 0;") used
			// to be discarded entirely — every field always started as nil
			// regardless of its initializer, silently breaking any class
			// that relied on a default value without an explicit
			// constructor assignment. Evaluated once here, at class
			// declaration time (not per-instance): correct for the
			// overwhelmingly common case of literal defaults
			// (numbers/strings/booleans/empty arrays); note instances do
			// share the same evaluated array/map reference for a
			// literal-array/map default, but R2Lang's array/map methods are
			// immutable/functional (reassign, not mutate in place — see
			// CLAUDE.md), so this doesn't produce the classic
			// shared-mutable-default bug seen in languages with in-place
			// mutation.
			var defaultValue interface{}
			if node.Value != nil {
				defaultValue = node.Value.Eval(env)
			}
			blueprint[node.Name] = defaultValue
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
