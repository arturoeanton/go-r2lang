// Test arrow functions comprehensive functionality
std.print("=== Testing Arrow Functions ===");

// Basic arrow function with single parameter
let double = x => x * 2;
std.print("double(5):", double(5));

// Arrow function with multiple parameters
let add = (a, b) => a + b;
std.print("add(3, 4):", add(3, 4));

// Arrow function with no parameters
let getHello = () => "Hello World";
std.print("getHello():", getHello());

// Arrow function with block body
let complexCalc = (x, y) => {
    let temp = x * 2;
    let result = temp + y;
    return result;
};
std.print("complexCalc(3, 5):", complexCalc(3, 5));

// Arrow function with default parameters
let greetWithDefault = (name = "World") => "Hello " + name;
std.print("greetWithDefault():", greetWithDefault());
std.print("greetWithDefault('Alice'):", greetWithDefault("Alice"));

// Arrow function used in array operations (if map exists)
let numbers = [1, 2, 3, 4, 5];

// Manual iteration to test arrow functions
let doubled = [];
let i = 0;
while (i < numbers.length) {
    let doubleFunc = x => x * 2;
    doubled[i] = doubleFunc(numbers[i]);
    i++;
}
std.print("Doubled array:", doubled);

// Nested arrow functions
let createMultiplier = factor => x => x * factor;
let triple = createMultiplier(3);
std.print("triple(4):", triple(4));

// Arrow function with string operations
let makeUppercase = str => {
    // Since we don't have built-in toUpperCase, just return with prefix
    return "UPPER: " + str;
};
std.print("makeUppercase('hello'):", makeUppercase("hello"));

// Arrow function with conditionals
let isEven = n => n % 2 == 0;
std.print("isEven(4):", isEven(4));
std.print("isEven(5):", isEven(5));

// Arrow function with object return
let createPoint = (x, y) => {
    return {x: x, y: y};
};
let point = createPoint(10, 20);
std.print("point.x:", point.x);
std.print("point.y:", point.y);

// Arrow function as object method
let calculator = {
    multiply: (a, b) => a * b,
    divide: (a, b) => a / b
};
std.print("calculator.multiply(6, 7):", calculator.multiply(6, 7));
std.print("calculator.divide(15, 3):", calculator.divide(15, 3));

std.print("=== Arrow Functions Test Completed ===");