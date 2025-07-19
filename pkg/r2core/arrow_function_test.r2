// Comprehensive arrow function tests

// 1. Single parameter without parentheses - expression body
let identity = x => x;
test.assertEq(identity(42), 42, "Single param arrow function should work");

// 2. Single parameter with parentheses - expression body
let double = (x) => x * 2;
test.assertEq(double(5), 10, "Single param with parens should work");

// 3. Multiple parameters - expression body
let add = (a, b) => a + b;
test.assertEq(add(3, 4), 7, "Multi param arrow function should work");

// 4. No parameters - expression body
let getFortyTwo = () => 42;
test.assertEq(getFortyTwo(), 42, "No param arrow function should work");

// 5. No parameters - block body
let greetWorld = () => {
    return "Hello World";
};
test.assertEq(greetWorld(), "Hello World", "No param block arrow function should work");

// 6. Single parameter - block body
let square = x => {
    return x * x;
};
test.assertEq(square(4), 16, "Single param block arrow function should work");

// 7. Multiple parameters - block body  
let multiply = (a, b) => {
    let result = a * b;
    return result;
};
test.assertEq(multiply(3, 5), 15, "Multi param block arrow function should work");

// 8. Arrow functions with default parameters
let power = (base, exponent = 2) => {
    let result = 1;
    for (let i = 0; i < exponent; i++) {
        result = result * base;
    }
    return result;
};
test.assertEq(power(3), 9, "Arrow function with default param should work");
test.assertEq(power(2, 3), 8, "Arrow function with provided param should work");

// 9. Arrow functions as variables and in expressions
let numbers = [1, 2, 3, 4, 5];
let doubled = [];
for (let i = 0; i < std.len(numbers); i++) {
    let doubler = x => x * 2;
    doubled[i] = doubler(numbers[i]);
}
test.assertEq(doubled[0], 2, "Arrow function in loop should work");
test.assertEq(doubled[4], 10, "Arrow function in loop should work for all elements");

// 10. Nested arrow functions
let makeAdder = x => y => x + y;
let add5 = makeAdder(5);
test.assertEq(add5(3), 8, "Nested arrow functions should work");

// 11. Arrow functions returning objects
let makePoint = (x, y) => {
    return {"x": x, "y": y};
};
let point = makePoint(3, 4);
test.assertEq(point["x"], 3, "Arrow function returning object should work");
test.assertEq(point["y"], 4, "Arrow function returning object should work");

// 12. Complex expressions
let complex = (a, b = 1) => {
    if (a > 0) {
        return a + b;
    } else {
        return a - b;
    }
};
test.assertEq(complex(5), 6, "Complex arrow function should work");
test.assertEq(complex(-5), -6, "Complex arrow function should work with negative");

std.print("All arrow function tests completed successfully!");