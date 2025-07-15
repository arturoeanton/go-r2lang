package r2core

import (
	"testing"
)

func TestBytecode_NumberLiteral(t *testing.T) {
	tests := []struct {
		name     string
		input    *NumberLiteral
		expected float64
	}{
		{"positive integer", &NumberLiteral{Value: 42}, 42},
		{"negative integer", &NumberLiteral{Value: -10}, -10},
		{"positive float", &NumberLiteral{Value: 3.14}, 3.14},
		{"zero", &NumberLiteral{Value: 0}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler := NewCompiler()
			err := compiler.Compile(tt.input)
			if err != nil {
				t.Fatalf("compilation failed: %v", err)
			}

			vm := NewVM(compiler.Bytecode())
			err = vm.Run()
			if err != nil {
				t.Fatalf("vm execution failed: %v", err)
			}

			result := vm.LastPoppedStackElem()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestBytecode_BooleanLiteral(t *testing.T) {
	tests := []struct {
		name     string
		input    *BooleanLiteral
		expected bool
	}{
		{"true", &BooleanLiteral{Value: true}, true},
		{"false", &BooleanLiteral{Value: false}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler := NewCompiler()
			err := compiler.Compile(tt.input)
			if err != nil {
				t.Fatalf("compilation failed: %v", err)
			}

			vm := NewVM(compiler.Bytecode())
			err = vm.Run()
			if err != nil {
				t.Fatalf("vm execution failed: %v", err)
			}

			result := vm.LastPoppedStackElem()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestBytecode_StringLiteral(t *testing.T) {
	tests := []struct {
		name     string
		input    *StringLiteral
		expected string
	}{
		{"simple string", &StringLiteral{Value: "hello"}, "hello"},
		{"empty string", &StringLiteral{Value: ""}, ""},
		{"string with spaces", &StringLiteral{Value: "hello world"}, "hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler := NewCompiler()
			err := compiler.Compile(tt.input)
			if err != nil {
				t.Fatalf("compilation failed: %v", err)
			}

			vm := NewVM(compiler.Bytecode())
			err = vm.Run()
			if err != nil {
				t.Fatalf("vm execution failed: %v", err)
			}

			result := vm.LastPoppedStackElem()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestBytecode_BinaryExpression_Arithmetic(t *testing.T) {
	tests := []struct {
		name     string
		left     *NumberLiteral
		op       string
		right    *NumberLiteral
		expected float64
	}{
		{"addition", &NumberLiteral{Value: 5}, "+", &NumberLiteral{Value: 3}, 8},
		{"subtraction", &NumberLiteral{Value: 10}, "-", &NumberLiteral{Value: 4}, 6},
		{"multiplication", &NumberLiteral{Value: 6}, "*", &NumberLiteral{Value: 7}, 42},
		{"division", &NumberLiteral{Value: 15}, "/", &NumberLiteral{Value: 3}, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := &BinaryExpression{
				Left:  tt.left,
				Op:    tt.op,
				Right: tt.right,
			}

			compiler := NewCompiler()
			err := compiler.Compile(expr)
			if err != nil {
				t.Fatalf("compilation failed: %v", err)
			}

			vm := NewVM(compiler.Bytecode())
			err = vm.Run()
			if err != nil {
				t.Fatalf("vm execution failed: %v", err)
			}

			result := vm.LastPoppedStackElem()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestBytecode_BinaryExpression_Comparison(t *testing.T) {
	tests := []struct {
		name     string
		left     *NumberLiteral
		op       string
		right    *NumberLiteral
		expected bool
	}{
		{"equal true", &NumberLiteral{Value: 5}, "==", &NumberLiteral{Value: 5}, true},
		{"equal false", &NumberLiteral{Value: 5}, "==", &NumberLiteral{Value: 3}, false},
		{"not equal true", &NumberLiteral{Value: 5}, "!=", &NumberLiteral{Value: 3}, true},
		{"not equal false", &NumberLiteral{Value: 5}, "!=", &NumberLiteral{Value: 5}, false},
		{"greater true", &NumberLiteral{Value: 10}, ">", &NumberLiteral{Value: 5}, true},
		{"greater false", &NumberLiteral{Value: 3}, ">", &NumberLiteral{Value: 5}, false},
		{"less true", &NumberLiteral{Value: 3}, "<", &NumberLiteral{Value: 5}, true},
		{"less false", &NumberLiteral{Value: 10}, "<", &NumberLiteral{Value: 5}, false},
		{"greater equal true", &NumberLiteral{Value: 5}, ">=", &NumberLiteral{Value: 5}, true},
		{"greater equal false", &NumberLiteral{Value: 3}, ">=", &NumberLiteral{Value: 5}, false},
		{"less equal true", &NumberLiteral{Value: 5}, "<=", &NumberLiteral{Value: 5}, true},
		{"less equal false", &NumberLiteral{Value: 10}, "<=", &NumberLiteral{Value: 5}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := &BinaryExpression{
				Left:  tt.left,
				Op:    tt.op,
				Right: tt.right,
			}

			compiler := NewCompiler()
			err := compiler.Compile(expr)
			if err != nil {
				t.Fatalf("compilation failed: %v", err)
			}

			vm := NewVM(compiler.Bytecode())
			err = vm.Run()
			if err != nil {
				t.Fatalf("vm execution failed: %v", err)
			}

			result := vm.LastPoppedStackElem()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestBytecode_ConstantFolding(t *testing.T) {
	tests := []struct {
		name     string
		left     *NumberLiteral
		op       string
		right    *NumberLiteral
		expected float64
	}{
		{"folding addition", &NumberLiteral{Value: 2}, "+", &NumberLiteral{Value: 3}, 5},
		{"folding multiplication", &NumberLiteral{Value: 4}, "*", &NumberLiteral{Value: 5}, 20},
		{"folding subtraction", &NumberLiteral{Value: 10}, "-", &NumberLiteral{Value: 3}, 7},
		{"folding division", &NumberLiteral{Value: 12}, "/", &NumberLiteral{Value: 4}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := &BinaryExpression{
				Left:  tt.left,
				Op:    tt.op,
				Right: tt.right,
			}

			compiler := NewCompiler()
			err := compiler.Compile(expr)
			if err != nil {
				t.Fatalf("compilation failed: %v", err)
			}

			// Verificar que el constant folding funcionó
			// Debería haber solo una instrucción OpConstant
			bytecode := compiler.Bytecode()
			if len(bytecode.Instructions) != 2 {
				t.Errorf("expected 2 instructions (OpConstant + operand), got %d", len(bytecode.Instructions))
			}

			vm := NewVM(bytecode)
			err = vm.Run()
			if err != nil {
				t.Fatalf("vm execution failed: %v", err)
			}

			result := vm.LastPoppedStackElem()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestBytecode_OptimizedEval(t *testing.T) {
	env := NewEnvironment()

	tests := []struct {
		name     string
		node     Node
		expected interface{}
	}{
		{
			"simple number",
			&NumberLiteral{Value: 42},
			float64(42),
		},
		{
			"simple boolean",
			&BooleanLiteral{Value: true},
			true,
		},
		{
			"simple string",
			&StringLiteral{Value: "hello"},
			"hello",
		},
		{
			"simple arithmetic",
			&BinaryExpression{
				Left:  &NumberLiteral{Value: 5},
				Op:    "+",
				Right: &NumberLiteral{Value: 3},
			},
			float64(8),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := OptimizedEval(tt.node, env)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestBytecode_ArrayLiteral(t *testing.T) {
	array := &ArrayLiteral{
		Elements: []Node{
			&NumberLiteral{Value: 1},
			&NumberLiteral{Value: 2},
			&NumberLiteral{Value: 3},
		},
	}

	compiler := NewCompiler()
	err := compiler.Compile(array)
	if err != nil {
		t.Fatalf("compilation failed: %v", err)
	}

	vm := NewVM(compiler.Bytecode())
	err = vm.Run()
	if err != nil {
		t.Fatalf("vm execution failed: %v", err)
	}

	result := vm.LastPoppedStackElem()
	if result == nil {
		t.Fatal("expected array result, got nil")
	}

	// Verificar que el resultado es un slice
	arr, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected []interface{}, got %T", result)
	}

	if len(arr) != 3 {
		t.Errorf("expected array length 3, got %d", len(arr))
	}
}
