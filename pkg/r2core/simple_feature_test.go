package r2core

import (
	"testing"
)

func TestBasicArithmetic(t *testing.T) {
	// Test basic arithmetic operations
	testCases := []struct {
		code     string
		expected interface{}
	}{
		{"5 + 3", float64(8)},
		{"10 - 4", float64(6)},
		{"6 * 7", float64(42)},
		{"15 / 3", float64(5)},
		{"10 % 3", float64(1)}, // Modulo operator
		{"7 % 2", float64(1)},
		{"15 % 5", float64(0)},
	}

	for _, tc := range testCases {
		env := NewEnvironment()
		parser := NewParser(tc.code)
		program := parser.ParseProgram()

		result := program.Eval(env)
		if result != tc.expected {
			t.Errorf("Code '%s': expected %v, got %v", tc.code, tc.expected, result)
		}
	}
}

func TestLogicalOperations(t *testing.T) {
	// Test logical operations
	testCases := []struct {
		code     string
		expected interface{}
	}{
		{"true && true", true},
		{"true && false", false},
		{"false && true", false},
		{"false && false", false},
		{"true || true", true},
		{"true || false", true},
		{"false || true", true},
		{"false || false", false},
		{"!true", false},
		{"!false", true},
	}

	for _, tc := range testCases {
		env := NewEnvironment()
		parser := NewParser(tc.code)
		program := parser.ParseProgram()

		result := program.Eval(env)
		if result != tc.expected {
			t.Errorf("Code '%s': expected %v, got %v", tc.code, tc.expected, result)
		}
	}
}

func TestComparisons(t *testing.T) {
	// Test comparison operations
	testCases := []struct {
		code     string
		expected interface{}
	}{
		{"5 > 3", true},
		{"3 > 5", false},
		{"5 < 3", false},
		{"3 < 5", true},
		{"5 >= 5", true},
		{"5 >= 3", true},
		{"3 >= 5", false},
		{"5 <= 5", true},
		{"3 <= 5", true},
		{"5 <= 3", false},
		{"5 == 5", true},
		{"5 == 3", false},
		{"5 != 3", true},
		{"5 != 5", false},
	}

	for _, tc := range testCases {
		env := NewEnvironment()
		parser := NewParser(tc.code)
		program := parser.ParseProgram()

		result := program.Eval(env)
		if result != tc.expected {
			t.Errorf("Code '%s': expected %v, got %v", tc.code, tc.expected, result)
		}
	}
}

func TestSimpleVariables(t *testing.T) {
	// Test simple variable assignments and access
	code := `
	let x = 10;
	let y = 20;
	x + y;
	`

	env := NewEnvironment()
	parser := NewParser(code)
	program := parser.ParseProgram()

	result := program.Eval(env)
	if result != float64(30) {
		t.Errorf("Expected 30, got %v", result)
	}
}

func TestSimpleArrays(t *testing.T) {
	// Test simple array creation and access
	code := `
	let arr = [1, 2, 3, 4, 5];
	arr[2];
	`

	env := NewEnvironment()
	parser := NewParser(code)
	program := parser.ParseProgram()

	result := program.Eval(env)
	if result != float64(3) {
		t.Errorf("Expected 3, got %v", result)
	}
}

func TestSimpleMaps(t *testing.T) {
	// Test simple map creation and access
	code := `
	let person = {name: "Alice", age: 30};
	person.name;
	`

	env := NewEnvironment()
	parser := NewParser(code)
	program := parser.ParseProgram()

	result := program.Eval(env)
	if result != "Alice" {
		t.Errorf("Expected 'Alice', got %v", result)
	}
}

/*TODO: Fix this test case to work with the new function syntax
func TestSimpleFunction(t *testing.T) {
	// Test simple function definition and call
	code := `
	func add(a, b) {
		return a + b;
	}

	add(5, 3);
	`

	env := NewEnvironment()
	parser := NewParser(code)
	program := parser.ParseProgram()

	result := program.Eval(env)
	if retVal, ok := result.(*ReturnValue); ok {
		if retVal.Value != float64(8) {
			t.Errorf("Expected 8, got %v", retVal.Value)
		}
	} else {
		t.Errorf("Expected ReturnValue, got %T", result)
	}
}*/

func TestSimpleIfElse(t *testing.T) {
	// Test simple if-else statements
	code := `
	let x = 10;
	let result = "";
	
	if (x > 5) {
		result = "greater";
	} else {
		result = "smaller";
	}
	
	result;
	`

	env := NewEnvironment()
	parser := NewParser(code)
	program := parser.ParseProgram()

	result := program.Eval(env)
	if result != "greater" {
		t.Errorf("Expected 'greater', got %v", result)
	}
}

func TestSimpleElseIf(t *testing.T) {
	// Test else-if chains
	code := `
	let score = 85;
	let grade = "";
	
	if (score >= 90) {
		grade = "A";
	} else if (score >= 80) {
		grade = "B";
	} else if (score >= 70) {
		grade = "C";
	} else {
		grade = "F";
	}
	
	grade;
	`

	env := NewEnvironment()
	parser := NewParser(code)
	program := parser.ParseProgram()

	result := program.Eval(env)
	if result != "B" {
		t.Errorf("Expected 'B', got %v", result)
	}
}

func TestStringConcatenation(t *testing.T) {
	// Test string concatenation
	code := `
	let first = "Hello";
	let second = "World";
	first + " " + second;
	`

	env := NewEnvironment()
	parser := NewParser(code)
	program := parser.ParseProgram()

	result := program.Eval(env)
	if result != "Hello World" {
		t.Errorf("Expected 'Hello World', got %v", result)
	}
}

func TestMultilineMap(t *testing.T) {
	// Test multiline map (one of our new features)
	code := `
	let config = {
		host: "localhost",
		port: 8080,
		ssl: true,
		timeout: 30
	};
	
	config.host;
	`

	env := NewEnvironment()
	parser := NewParser(code)
	program := parser.ParseProgram()

	result := program.Eval(env)
	if result != "localhost" {
		t.Errorf("Expected 'localhost', got %v", result)
	}
}

func TestMultilineArray(t *testing.T) {
	// Test multiline array (one of our new features)
	code := `
	let users = [
		{name: "Alice", age: 30},
		{name: "Bob", age: 25},
		{name: "Carol", age: 35}
	];
	
	users[1].name;
	`

	env := NewEnvironment()
	parser := NewParser(code)
	program := parser.ParseProgram()

	result := program.Eval(env)
	if result != "Bob" {
		t.Errorf("Expected 'Bob', got %v", result)
	}
}

/* TODO: Fix this test case to work with the new function syntax
func TestBuiltinFunctionsBasic(t *testing.T) {
	// Test basic built-in functions
	testCases := []struct {
		name string
		code string
		test func(interface{}) bool
	}{
		{
			"std.len array",
			"std.len([1, 2, 3])",
			func(result interface{}) bool { return result == float64(3) },
		},
		{
			"std.len string",
			"std.len(\"hello\")",
			func(result interface{}) bool { return result == float64(5) },
		},
		{
			"std.typeOf number",
			"std.typeOf(42)",
			func(result interface{}) bool { return result == "float64" },
		},
		{
			"std.typeOf string",
			"std.typeOf(\"test\")",
			func(result interface{}) bool { return result == "string" },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			env := NewEnvironment()
			parser := NewParser(tc.code)
			program := parser.ParseProgram()

			result := program.Eval(env)
			if !tc.test(result) {
				t.Errorf("Test failed for '%s', got: %v (%T)", tc.code, result, result)
			}
		})
	}
}
*/

func TestNullHandling(t *testing.T) {
	// Test null/nil handling
	testCases := []struct {
		code     string
		expected interface{}
	}{
		{"nil == nil", true},
		{"nil != nil", false},
		{"nil == 5", false},
		{"5 != nil", true},
	}

	for _, tc := range testCases {
		env := NewEnvironment()
		parser := NewParser(tc.code)
		program := parser.ParseProgram()

		result := program.Eval(env)
		if result != tc.expected {
			t.Errorf("Code '%s': expected %v, got %v", tc.code, tc.expected, result)
		}
	}
}
