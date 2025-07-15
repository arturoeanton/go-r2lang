package r2core

import (
	"testing"
)

func TestBreakStatement_Simple(t *testing.T) {
	// Test break in while loop
	env := NewEnvironment()
	env.Set("true", true)
	env.Set("false", false)

	code := `
		let count = 0;
		let i = 0;
		while (i < 10) {
			if (i == 3) {
				break;
			}
			count = count + 1;
			i = i + 1;
		}
		count;
	`

	parser := NewParser(code)
	ast := parser.ParseProgram()

	result := ast.Eval(env)

	if result.(float64) != 3 {
		t.Errorf("Expected 3, got %f", result.(float64))
	}
}

func TestContinueStatement_Simple(t *testing.T) {
	// Test continue in while loop
	env := NewEnvironment()
	env.Set("true", true)
	env.Set("false", false)

	code := `
		let count = 0;
		let i = 0;
		while (i < 5) {
			i = i + 1;
			if (i == 3) {
				continue;
			}
			count = count + 1;
		}
		count;
	`

	parser := NewParser(code)
	ast := parser.ParseProgram()

	result := ast.Eval(env)

	if result.(float64) != 4 {
		t.Errorf("Expected 4, got %f", result.(float64))
	}
}

func TestBreakStatement_ForLoop_Simple(t *testing.T) {
	// Test break in for loop
	env := NewEnvironment()
	env.Set("true", true)
	env.Set("false", false)

	code := `
		let count = 0;
		for (let i = 0; i < 10; i++) {
			if (i == 3) {
				break;
			}
			count = count + 1;
		}
		count;
	`

	parser := NewParser(code)
	ast := parser.ParseProgram()

	result := ast.Eval(env)

	if result.(float64) != 3 {
		t.Errorf("Expected 3, got %f", result.(float64))
	}
}

func TestContinueStatement_ForLoop_Simple(t *testing.T) {
	// Test continue in for loop
	env := NewEnvironment()
	env.Set("true", true)
	env.Set("false", false)

	code := `
		let count = 0;
		for (let i = 0; i < 5; i++) {
			if (i == 2) {
				continue;
			}
			count = count + 1;
		}
		count;
	`

	parser := NewParser(code)
	ast := parser.ParseProgram()

	result := ast.Eval(env)

	if result.(float64) != 4 {
		t.Errorf("Expected 4, got %f", result.(float64))
	}
}

func TestBreakStatement_ForInLoop_Simple(t *testing.T) {
	// Test break in for-in loop
	env := NewEnvironment()
	env.Set("true", true)
	env.Set("false", false)

	code := `
		let arr = [1, 2, 3, 4, 5];
		let count = 0;
		for (i in arr) {
			if (arr[i] == 3) {
				break;
			}
			count = count + 1;
		}
		count;
	`

	parser := NewParser(code)
	ast := parser.ParseProgram()

	result := ast.Eval(env)

	if result.(float64) != 2 {
		t.Errorf("Expected 2, got %f", result.(float64))
	}
}

func TestContinueStatement_ForInLoop_Simple(t *testing.T) {
	// Test continue in for-in loop
	env := NewEnvironment()
	env.Set("true", true)
	env.Set("false", false)

	code := `
		let arr = [1, 2, 3, 4, 5];
		let count = 0;
		for (i in arr) {
			if (arr[i] == 3) {
				continue;
			}
			count = count + 1;
		}
		count;
	`

	parser := NewParser(code)
	ast := parser.ParseProgram()

	result := ast.Eval(env)

	if result.(float64) != 4 {
		t.Errorf("Expected 4, got %f", result.(float64))
	}
}
