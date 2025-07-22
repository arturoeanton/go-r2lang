print("=== P6 Features: Partial Application and Currying ===")

print("\n1. Partial Application with Placeholders")

// Basic partial application using underscore placeholder
func add(a, b) {
    return a + b
}

let addFive = add(5, _)
let result1 = addFive(10)
print("add(5, _)(10) =", result1)  // Expected: 15

// Multiple placeholders
func multiply3(a, b, c) {
    return a * b * c
}

let doubleAndMultiply = multiply3(2, _, _)
let result2 = doubleAndMultiply(3, 4)
print("multiply3(2, _, _)(3, 4) =", result2)  // Expected: 24

print("\n2. Explicit Partial Application with partial() function")

// Using the partial() function explicitly
func divide(a, b) {
    return a / b
}

let divideByTwo = partial(divide, _, 2)
let result4 = divideByTwo(20)
print("partial(divide, _, 2)(20) =", result4)  // Expected: 10

print("\n3. Currying Functions")

// Basic currying
func add3(a, b, c) {
    return a + b + c
}

let curriedAdd = curry(add3)
let result5 = curriedAdd(1)(2)(3)
print("curry(add3)(1)(2)(3) =", result5)  // Expected: 6

print("\n=== P6 Features Demo Complete ===")