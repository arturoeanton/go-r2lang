// Simple P6 test
std.print("Testing P6 features...")

// Test placeholder
func add(a, b) {
    return a + b
}

let addFive = add(5, _)
std.print("Partial application result:", addFive(10))