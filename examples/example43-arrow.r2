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

// 9. Nested arrow functions
let makeAdder = x => y => x + y;
let add5 = makeAdder(5);
test.assertEq(add5(3), 8, "Nested arrow functions should work");

// 10. Arrow functions returning objects
let makePoint = (x, y) => {
    return {"x": x, "y": y};
};
let point = makePoint(3, 4);
test.assertEq(point["x"], 3, "Arrow function returning object should work");
test.assertEq(point["y"], 4, "Arrow function returning object should work");

// 11. Complex expressions with conditionals
let complex = (a, b = 1) => {
    if (a > 0) {
        return a + b;
    } else {
        return a - b;
    }
};
test.assertEq(complex(5), 6, "Complex arrow function should work");
test.assertEq(complex(-5), -6, "Complex arrow function should work with negative");

// 12. Arrow functions with string concatenation
let makeGreeting = name => "Hello " + name + "!";
test.assertEq(makeGreeting("World"), "Hello World!", "String concat arrow function should work");

// 13. Arrow functions with boolean logic
let isEven = x => x % 2 == 0;
test.assertTrue(isEven(4), "Boolean arrow function should work");
test.assertTrue(!isEven(5), "Boolean arrow function should work for odd numbers");

std.print("All arrow function tests completed successfully!");