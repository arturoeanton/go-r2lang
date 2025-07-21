// Test file for position tracking in R2Lang error reporting
// This file tests various error scenarios to ensure proper
// source code position information is included in error messages

func testDivisionByZero() {
    print("Testing division by zero error...")
    result = 10 / 0  // Line 7 - should report this position
    return result
}

func testNestedError() {
    print("Testing nested function error...")
    
    func innerFunction() {
        return 1 / 0  // Line 15 - should report this position in nested context
    }
    
    return innerFunction()
}

func testErrorInLoop() {
    print("Testing error inside loop...")
    for (i = 0; i < 3; i = i + 1) {
        if (i == 2) {
            return 1 / 0  // Line 25 - should report this position inside loop
        }
    }
}

func testErrorInConditional() {
    print("Testing error in conditional...")
    if (true) {
        return 1 / 0  // Line 32 - should report this position in if block
    } else {
        return 2 / 0  // Line 34 - alternative error position
    }
}

// Test function that should work without errors
func testValidOperation() {
    print("Testing valid operation...")
    a = 10
    b = 5
    result = a / b
    return result
}

// Individual test functions to call one at a time
func main() {
    print("=== R2Lang Position Tracking Tests ===")
    print("Comment out tests to run them individually")
    print("")
    
    // Test valid operation first
    print("1. Valid operation test:")
    testValidOperation()
    print("   Success!")
    print("")
    
    // Uncomment ONE test at a time to see position tracking:
    
    print("2. Division by zero test (should show line 7):")
    // testDivisionByZero()
    
    print("3. Nested error test (should show line 15):")
    // testNestedError()
    
    print("4. Error in loop test (should show line 25):")
    // testErrorInLoop()
    
    print("5. Error in conditional test (should show line 32):")
    testErrorInConditional()
    
    print("=== Position Tracking Tests Complete ===")
    print("All error messages should include file name, line number, and column position")
}