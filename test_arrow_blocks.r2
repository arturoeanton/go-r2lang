// Test arrow functions with block bodies
let addAndLog = (a, b) => {
    let result = a + b;
    std.print("Adding", a, "and", b, "=", result);
    return result;
};

let greet = () => {
    std.print("Hello from block arrow function!");
    return "greeting";
};

std.print("Result:", addAndLog(5, 3));
std.print("Greeting result:", greet());