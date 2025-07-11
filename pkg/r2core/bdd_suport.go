package r2core

import "fmt"

// ============================================================
// 3) AST - Node interface
// ============================================================
// TestCase representa un caso de prueba con un nombre y pasos.

type TestCase struct {
	Name  string
	Steps []TestStep
}

type NodeTest interface {
	Eval(env *Environment) interface{}
	EvalStep(env *Environment) interface{}
}

type TestStep struct {
	Type    string // "Given", "When", "Then", "And"
	Command Node
}

// Eval ejecuta el caso de prueba.
func (tc *TestCase) Eval(env *Environment) interface{} {
	fmt.Printf("Executing Test Case: %s\n", tc.Name)
	var previousStepType string

	for _, step := range tc.Steps {
		stepType := step.Type
		if stepType == "And" {
			stepType = previousStepType
		} else {
			previousStepType = stepType
		}
		fmt.Printf("  %s: ", stepType)

		if ce, ok := step.Command.(*CallExpression); ok {
			calleeVal := ce.Callee.Eval(env)
			var argVals []interface{}
			for _, a := range ce.Args {
				argVals = append(argVals, a.Eval(env))
			}
			if currentStep, ok := calleeVal.(*UserFunction); ok {
				out := currentStep.CallStep(env, argVals...)
				if out != nil {
					fmt.Println(out)
				}
			}
			continue
		}

		if fl, ok := step.Command.(*FunctionLiteral); ok {
			currentStep := fl.Eval(env).(*UserFunction)
			out := currentStep.CallStep(env)
			if out != nil {
				fmt.Println(out)
			}
			continue
		}

	}
	fmt.Println("Test Case Executed Successfully.")
	return nil
}

func (ts *TestStep) Eval(env *Environment) interface{} {
	// Ejecutar el comando del paso
	return ts.Command.Eval(env)
}
