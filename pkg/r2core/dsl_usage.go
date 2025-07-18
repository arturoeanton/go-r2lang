package r2core

import (
	"fmt"
)

// DSLUsage represents a DSL usage block like: var rules = BusinessRules.use()
type DSLUsage struct {
	Token   Token
	DSLName *Identifier
	Body    *BlockStatement
}

func (dsl *DSLUsage) Eval(env *Environment) interface{} {
	// Get the DSL definition from environment
	dslDef, exists := env.Get(dsl.DSLName.Name)
	if !exists {
		return fmt.Errorf("DSL '%s' not found", dsl.DSLName.Name)
	}

	// Get the DSL object and call its use method
	if dslObj, ok := dslDef.(map[string]interface{}); ok {
		if useMethod, ok := dslObj["use"].(func() interface{}); ok {
			return useMethod()
		}
	}

	return fmt.Errorf("DSL '%s' does not have a valid 'use' method", dsl.DSLName.Name)
}

func (dsl *DSLUsage) String() string {
	return fmt.Sprintf("DSLUsage(%s)", dsl.DSLName.Name)
}
