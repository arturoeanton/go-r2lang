package r2core

import "fmt"

type ReturnValue struct {
	Value interface{}
}

func (rv *ReturnValue) String() string {
	return fmt.Sprintf("%v", rv.Value)
}
