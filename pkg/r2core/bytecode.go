package r2core

import (
	"fmt"
)

// OpCode representa un código de operación bytecode
type OpCode uint8

const (
	// Operaciones básicas
	OpConstant     OpCode = iota // Cargar una constante
	OpAdd                        // Suma
	OpSub                        // Resta
	OpMul                        // Multiplicación
	OpDiv                        // División
	OpEqual                      // Igualdad
	OpNotEqual                   // Desigualdad
	OpGreater                    // Mayor que
	OpLess                       // Menor que
	OpGreaterEqual               // Mayor o igual que
	OpLessEqual                  // Menor o igual que

	// Operaciones de variables
	OpGetLocal     // Obtener variable local
	OpSetLocal     // Establecer variable local
	OpGetGlobal    // Obtener variable global
	OpSetGlobal    // Establecer variable global
	OpDefineGlobal // Definir variable global

	// Operaciones de control de flujo
	OpJump          // Salto incondicional
	OpJumpNotTruthy // Salto si es falso
	OpJumpIfFalse   // Salto condicional optimizado
	OpLoop          // Salto hacia atrás (loop)

	// Operaciones de función
	OpCall           // Llamada a función
	OpReturn         // Retorno de función
	OpClosure        // Crear closure
	OpGetBuiltin     // Acceso a funciones built-in
	OpGetFree        // Variables libres en closures
	OpSetFree        // Establecer variables libres
	OpCurrentClosure // Acceso a closure actual

	// Operaciones de array y objetos
	OpArray       // Crear array
	OpHash        // Crear hash/map
	OpIndex       // Acceso por índice
	OpGetProperty // obj.prop
	OpSetProperty // obj.prop = value

	// Operaciones especiales
	OpPop    // Quitar del stack
	OpTrue   // Valor true
	OpFalse  // Valor false
	OpNull   // Valor null
	OpNegate // Negación unaria
	OpBang   // Negación lógica
)

// Instruction representa una instrucción bytecode
type Instruction []byte

// Instructions es una secuencia de instrucciones
type Instructions []byte

// Bytecode contiene las instrucciones y constantes
type Bytecode struct {
	Instructions Instructions
	Constants    []interface{}
}

// Compiler compila AST a bytecode
type Compiler struct {
	instructions Instructions
	constants    []interface{}
	lastInstPos  int
}

// NewCompiler crea un nuevo compilador
func NewCompiler() *Compiler {
	return &Compiler{
		instructions: Instructions{},
		constants:    []interface{}{},
		lastInstPos:  0,
	}
}

// Compile compila un nodo AST a bytecode
func (c *Compiler) Compile(node Node) error {
	switch n := node.(type) {
	case *NumberLiteral:
		constant := n.Value
		c.emit(OpConstant, c.addConstant(constant))

	case *BooleanLiteral:
		if n.Value {
			c.emit(OpTrue)
		} else {
			c.emit(OpFalse)
		}

	case *StringLiteral:
		constant := n.Value
		c.emit(OpConstant, c.addConstant(constant))

	case *BinaryExpression:
		// Intentar constant folding primero
		if c.optimizeConstantFolding(n) {
			return nil
		}

		// Compilar operandos primero
		err := c.Compile(n.Left)
		if err != nil {
			return err
		}

		err = c.Compile(n.Right)
		if err != nil {
			return err
		}

		// Luego emitir la operación
		switch n.Op {
		case "+":
			c.emit(OpAdd)
		case "-":
			c.emit(OpSub)
		case "*":
			c.emit(OpMul)
		case "/":
			c.emit(OpDiv)
		case "==":
			c.emit(OpEqual)
		case "!=":
			c.emit(OpNotEqual)
		case ">":
			c.emit(OpGreater)
		case "<":
			c.emit(OpLess)
		case ">=":
			c.emit(OpGreaterEqual)
		case "<=":
			c.emit(OpLessEqual)
		default:
			return fmt.Errorf("unknown operator %s", n.Op)
		}

	case *ArrayLiteral:
		for _, elem := range n.Elements {
			err := c.Compile(elem)
			if err != nil {
				return err
			}
		}
		c.emit(OpArray, len(n.Elements))

	default:
		return fmt.Errorf("compilation of %T not implemented", node)
	}

	return nil
}

// emit agrega una instrucción bytecode
func (c *Compiler) emit(op OpCode, operands ...int) int {
	ins := makeInstruction(op, operands...)
	pos := c.addInstruction(ins)
	c.lastInstPos = pos
	return pos
}

// addInstruction agrega una instrucción a la lista
func (c *Compiler) addInstruction(ins Instruction) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}

// addConstant agrega una constante y retorna su índice
func (c *Compiler) addConstant(obj interface{}) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

// Bytecode retorna el bytecode compilado
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

// makeInstruction crea una instrucción bytecode con operandos
func makeInstruction(op OpCode, operands ...int) []byte {
	instruction := []byte{byte(op)}

	for _, operand := range operands {
		// Por ahora, solo operandos de 8-bit para simplificar
		if operand < 256 {
			instruction = append(instruction, byte(operand))
		} else {
			// Para operandos grandes, usar módulo para ajustar a 8-bit
			instruction = append(instruction, byte(operand%256))
		}
	}

	return instruction
}

// optimizeConstantFolding optimiza expresiones binarias con constantes
func (c *Compiler) optimizeConstantFolding(n *BinaryExpression) bool {
	left, isLeftNumber := n.Left.(*NumberLiteral)
	right, isRightNumber := n.Right.(*NumberLiteral)

	// Solo optimizar si ambos operandos son números
	if !isLeftNumber || !isRightNumber {
		return false
	}

	var result float64
	switch n.Op {
	case "+":
		result = left.Value + right.Value
	case "-":
		result = left.Value - right.Value
	case "*":
		result = left.Value * right.Value
	case "/":
		if right.Value == 0 {
			return false // Evitar división por cero
		}
		result = left.Value / right.Value
	default:
		return false // No optimizar otros operadores por ahora
	}

	// Emitir la constante optimizada
	c.emit(OpConstant, c.addConstant(result))
	return true
}

// peepholeOptimize optimiza patrones de instrucciones
func (c *Compiler) peepholeOptimize() {
	// Implementación básica de peephole optimization
	// Por ahora, solo eliminar push/pop consecutivos innecesarios
	// TODO: Implementar más optimizaciones
}

// VM ejecuta bytecode
type VM struct {
	constants    []interface{}
	instructions Instructions
	stack        []interface{}
	sp           int // stack pointer
}

// NewVM crea una nueva máquina virtual
func NewVM(bytecode *Bytecode) *VM {
	return &VM{
		constants:    bytecode.Constants,
		instructions: bytecode.Instructions,
		stack:        make([]interface{}, 2048), // Stack de 2048 elementos
		sp:           0,
	}
}

// Run ejecuta el bytecode
func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); {
		opcode := OpCode(vm.instructions[ip])

		switch opcode {
		case OpConstant:
			constIndex := int(vm.instructions[ip+1])
			ip += 2

			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}

		case OpAdd:
			ip++
			right := vm.pop()
			left := vm.pop()

			result := addValues(left, right)
			err := vm.push(result)
			if err != nil {
				return err
			}

		case OpSub:
			ip++
			right := vm.pop()
			left := vm.pop()

			result := subValues(left, right)
			err := vm.push(result)
			if err != nil {
				return err
			}

		case OpMul:
			ip++
			right := vm.pop()
			left := vm.pop()

			result := mulValues(left, right)
			err := vm.push(result)
			if err != nil {
				return err
			}

		case OpDiv:
			ip++
			right := vm.pop()
			left := vm.pop()

			result := divValues(left, right)
			err := vm.push(result)
			if err != nil {
				return err
			}

		case OpEqual:
			ip++
			right := vm.pop()
			left := vm.pop()

			result := equals(left, right)
			err := vm.push(result)
			if err != nil {
				return err
			}

		case OpGreater:
			ip++
			right := vm.pop()
			left := vm.pop()

			result := toFloat(left) > toFloat(right)
			err := vm.push(result)
			if err != nil {
				return err
			}

		case OpTrue:
			ip++
			err := vm.push(true)
			if err != nil {
				return err
			}

		case OpFalse:
			ip++
			err := vm.push(false)
			if err != nil {
				return err
			}

		case OpNotEqual:
			ip++
			right := vm.pop()
			left := vm.pop()

			result := !equals(left, right)
			err := vm.push(result)
			if err != nil {
				return err
			}

		case OpLess:
			ip++
			right := vm.pop()
			left := vm.pop()

			result := toFloat(left) < toFloat(right)
			err := vm.push(result)
			if err != nil {
				return err
			}

		case OpGreaterEqual:
			ip++
			right := vm.pop()
			left := vm.pop()

			result := toFloat(left) >= toFloat(right)
			err := vm.push(result)
			if err != nil {
				return err
			}

		case OpLessEqual:
			ip++
			right := vm.pop()
			left := vm.pop()

			result := toFloat(left) <= toFloat(right)
			err := vm.push(result)
			if err != nil {
				return err
			}

		case OpNull:
			ip++
			err := vm.push(nil)
			if err != nil {
				return err
			}

		case OpArray:
			arraySize := int(vm.instructions[ip+1])
			ip += 2

			// Crear array tomando elementos del stack
			array := make([]interface{}, arraySize)
			for i := arraySize - 1; i >= 0; i-- {
				array[i] = vm.pop()
			}

			err := vm.push(array)
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("unknown opcode %d", opcode)
		}
	}

	return nil
}

// push agrega un valor al stack
func (vm *VM) push(obj interface{}) error {
	if vm.sp >= len(vm.stack) {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = obj
	vm.sp++

	return nil
}

// pop quita un valor del stack
func (vm *VM) pop() interface{} {
	if vm.sp == 0 {
		return nil
	}

	obj := vm.stack[vm.sp-1]
	vm.sp--

	return obj
}

// LastPoppedStackElem retorna el último elemento del stack
func (vm *VM) LastPoppedStackElem() interface{} {
	if vm.sp > 0 {
		return vm.stack[vm.sp-1]
	}
	return nil
}

// OptimizedEval evalúa un nodo usando bytecode si es posible
func OptimizedEval(node Node, env *Environment) interface{} {
	// Verificar si el nodo es candidato para bytecode
	if !isBytecodeCandidate(node) {
		return node.Eval(env)
	}

	// Intentar compilación a bytecode para operaciones simples
	compiler := NewCompiler()
	err := compiler.Compile(node)
	if err != nil {
		// Si falla la compilación, usar evaluación normal
		return node.Eval(env)
	}

	// Ejecutar en VM
	vm := NewVM(compiler.Bytecode())
	err = vm.Run()
	if err != nil {
		// Si falla la ejecución, usar evaluación normal
		return node.Eval(env)
	}

	// Retornar el resultado del bytecode
	return vm.LastPoppedStackElem()
}

// isBytecodeCandidate determina si un nodo se beneficia del bytecode
func isBytecodeCandidate(node Node) bool {
	switch n := node.(type) {
	case *BinaryExpression:
		// Solo operaciones numéricas simples se benefician
		return isNumericOp(n.Op) && isSimpleExpression(n.Left) && isSimpleExpression(n.Right)
	case *NumberLiteral, *BooleanLiteral, *StringLiteral:
		return true
	default:
		return false
	}
}

// isNumericOp verifica si es una operación numérica
func isNumericOp(op string) bool {
	return op == "+" || op == "-" || op == "*" || op == "/" ||
		op == ">" || op == "<" || op == ">=" || op == "<=" ||
		op == "==" || op == "!="
}

// isSimpleExpression verifica si es una expresión simple
func isSimpleExpression(node Node) bool {
	switch node.(type) {
	case *NumberLiteral, *BooleanLiteral, *StringLiteral:
		return true
	default:
		return false
	}
}
