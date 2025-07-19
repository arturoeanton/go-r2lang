// Comprehensive test suite for all new P0 and P1 syntax features
std.print("=== Testing R2Lang Syntax Improvements ===");

// =================
// P0.1: Logical Negation Operator (!)
// =================
std.print("\n--- Testing Logical Negation Operator (!) ---");

// Basic boolean negation
let isTrue = true;
let isFalse = false;
std.print("!true =", !isTrue);
std.print("!false =", !isFalse);

// Truthiness testing
std.print("!0 =", !0);
std.print("!1 =", !1);
std.print("!'' =", !"");
std.print("!'hello' =", !"hello");

// Double negation
std.print("!!true =", !!true);
std.print("!!false =", !!false);
std.print("!!0 =", !!0);

// In conditionals
if (!isFalse) {
    std.print("Negation works in conditionals");
}

if (!(1 == 2)) {
    std.print("Negation works with complex expressions");
}

// =================
// P0.2: Compound Assignment Operators (+=, -=, *=, /=)
// =================
std.print("\n--- Testing Compound Assignment Operators ---");

// Numeric operations
let num = 10;
std.print("Initial num:", num);

num += 5;
std.print("After += 5:", num);

num -= 3;
std.print("After -= 3:", num);

num *= 2;
std.print("After *= 2:", num);

num /= 4;
std.print("After /= 4:", num);

// String concatenation
let message = "Hello";
std.print("Initial message:", message);

message += " World";
std.print("After += ' World':", message);

message += "!";
std.print("After += '!':", message);

// With variables
let a = 100;
let b = 25;
a += b;
std.print("100 += 25 =", a);

a -= b;
std.print("125 -= 25 =", a);

a *= b;
std.print("100 *= 25 =", a);

a /= b;
std.print("2500 /= 25 =", a);

// =================
// P1.1: Const Declarations
// =================
std.print("\n--- Testing Const Declarations ---");

// Single const declaration
const MY_CONSTANT = 42;
std.print("MY_CONSTANT:", MY_CONSTANT);

// Multiple const declarations
const API_URL = "https://api.example.com", MAX_RETRIES = 3, TIMEOUT_MS = 5000;
std.print("API_URL:", API_URL);
std.print("MAX_RETRIES:", MAX_RETRIES);
std.print("TIMEOUT_MS:", TIMEOUT_MS);

// Const with objects
const CONFIG = {
    debug: true,
    version: "2.0.0",
    features: ["negation", "compound", "const", "defaults"]
};
std.print("CONFIG.debug:", CONFIG.debug);
std.print("CONFIG.version:", CONFIG.version);
std.print("CONFIG.features[0]:", CONFIG.features[0]);

// Const with arrays
const NUMBERS = [1, 2, 3, 4, 5];
std.print("NUMBERS:", NUMBERS);

// =================
// P1.2: Default Parameters
// =================
std.print("\n--- Testing Default Parameters ---");

// Function with one default parameter
func greet(name = "Anonymous") {
    return "Hello, " + name + "!";
}

std.print("greet():", greet());
std.print("greet('Alice'):", greet("Alice"));

// Function with multiple defaults
func createConnection(host = "localhost", port = 8080, secure = false) {
    return {
        host: host,
        port: port,
        secure: secure,
        url: (secure ? "https://" : "http://") + host + ":" + port
    };
}

let conn1 = createConnection();
std.print("Default connection:", conn1);

let conn2 = createConnection("example.com");
std.print("Custom host:", conn2);

let conn3 = createConnection("api.example.com", 443, true);
std.print("Full custom:", conn3);

// Mixed parameters (required + defaults)
func calculateArea(width, height = 10, unit = "px") {
    return width * height + " " + unit + "²";
}

std.print("calculateArea(5):", calculateArea(5));
std.print("calculateArea(5, 8):", calculateArea(5, 8));
std.print("calculateArea(5, 8, 'cm'):", calculateArea(5, 8, "cm"));

// Anonymous function with defaults
let power = func(base, exponent = 2) {
    let result = 1;
    let i = 0;
    while (i < exponent) {
        result *= base;
        i++;
    }
    return result;
};

std.print("power(3):", power(3));
std.print("power(3, 3):", power(3, 3));
std.print("power(2, 4):", power(2, 4));

// Default parameters with expression defaults
func withExpressionDefault(x, y = 10) {
    return x + y;
}

std.print("withExpressionDefault(5):", withExpressionDefault(5));
std.print("withExpressionDefault(5, 3):", withExpressionDefault(5, 3));

// =================
// Combined Feature Testing
// =================
std.print("\n--- Testing Combined Features ---");

// Using all new features together
const MULTIPLIER = 2;
let counter = 0;

func incrementBy(amount = 1, multiply = false) {
    if (!multiply) {
        counter += amount;
    } else {
        counter *= amount;
    }
    return counter;
}

std.print("Initial counter:", counter);
std.print("incrementBy():", incrementBy());
std.print("incrementBy(5):", incrementBy(5));
std.print("incrementBy(MULTIPLIER, true):", incrementBy(MULTIPLIER, true));

// Const with default parameters
const DEFAULT_CONFIG = {
    retries: 3,
    timeout: 1000
};

func makeRequest(url, config = DEFAULT_CONFIG) {
    if (!url) {
        return "Error: URL is required";
    }
    return "Making request to " + url + " with " + config.retries + " retries";
}

std.print("makeRequest('api.example.com'):", makeRequest("api.example.com"));

std.print("\n=== All syntax improvement tests completed successfully! ===");