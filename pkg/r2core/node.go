package r2core

// ============================================================
// 3) AST - Node interface
// ============================================================

type Node interface {
	Eval(env *Environment) interface{}
}
