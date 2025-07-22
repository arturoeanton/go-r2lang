std.print("Testing P6 features...")

func add(a, b) {
    return a + b
}

// Test placeholder based partial application
let addFive = add(5, _)
std.print("Partial application result:", addFive(10))

// Test curry function  
let curriedAdd = curry(add)
std.print("Curried result:", curriedAdd(5)(10))

// Test explicit partial
let addTen = partial(add, 10)
std.print("Explicit partial result:", addTen(5))