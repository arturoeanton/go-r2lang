// Test arrow functions with default parameters
let multiply = (a, b = 2) => a * b;
let greetUser = (name = "World") => "Hello " + name + "!";

std.print("multiply(5):", multiply(5));
std.print("multiply(5, 3):", multiply(5, 3));
std.print("greetUser():", greetUser());
std.print("greetUser('Alice'):", greetUser("Alice"));