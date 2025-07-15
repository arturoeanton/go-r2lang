// Basic Unit Testing Example for R2Lang
// This demonstrates the core testing features with describe/it syntax

import "r2test" as test;

// Example calculator module to test
let Calculator = {
    add: func(a, b) {
        return a + b;
    },
    
    subtract: func(a, b) {
        return a - b;
    },
    
    multiply: func(a, b) {
        return a * b;
    },
    
    divide: func(a, b) {
        if (b == 0) {
            throw "Division by zero";
        }
        return a / b;
    }
};

// Basic test suite
test.describe("Calculator", func() {
    
    test.describe("Addition", func() {
        test.it("should add two positive numbers", func() {
            let result = Calculator.add(2, 3);
            test.assert("addition result").that(result).equals(5);
        });
        
        test.it("should add negative numbers", func() {
            let result = Calculator.add(-5, -3);
            test.assert("negative addition").that(result).equals(-8);
        });
        
        test.it("should handle zero", func() {
            let result = Calculator.add(5, 0);
            test.assert("zero addition").that(result).equals(5);
        });
    });
    
    test.describe("Subtraction", func() {
        test.it("should subtract two numbers", func() {
            let result = Calculator.subtract(10, 4);
            test.assert("subtraction result").that(result).equals(6);
        });
        
        test.it("should handle negative results", func() {
            let result = Calculator.subtract(3, 7);
            test.assert("negative result").that(result).equals(-4);
        });
    });
    
    test.describe("Multiplication", func() {
        test.it("should multiply two numbers", func() {
            let result = Calculator.multiply(6, 7);
            test.assert("multiplication result").that(result).equals(42);
        });
        
        test.it("should handle zero multiplication", func() {
            let result = Calculator.multiply(5, 0);
            test.assert("zero multiplication").that(result).equals(0);
        });
    });
    
    test.describe("Division", func() {
        test.it("should divide two numbers", func() {
            let result = Calculator.divide(15, 3);
            test.assert("division result").that(result).equals(5);
        });
        
        test.it("should handle decimal results", func() {
            let result = Calculator.divide(7, 2);
            test.assert("decimal division").that(result).equals(3.5);
        });
        
        test.it("should throw error for division by zero", func() {
            test.assert("division by zero").throws(func() {
                Calculator.divide(10, 0);
            }).withMessage("Division by zero");
        });
    });
});

// Test lifecycle hooks
test.describe("Lifecycle Hooks Example", func() {
    let counter = 0;
    let setupData = null;
    
    test.beforeAll(func() {
        print("Setting up test suite...");
        setupData = "initialized";
    });
    
    test.afterAll(func() {
        print("Cleaning up test suite...");
        setupData = null;
    });
    
    test.beforeEach(func() {
        counter = 0;
        print("Before each test - counter reset to:", counter);
    });
    
    test.afterEach(func() {
        print("After each test - counter is:", counter);
    });
    
    test.it("should have setup data available", func() {
        test.assert("setup data").that(setupData).equals("initialized");
    });
    
    test.it("should increment counter", func() {
        counter = counter + 1;
        test.assert("counter increment").that(counter).equals(1);
    });
    
    test.it("should start fresh each test", func() {
        test.assert("fresh counter").that(counter).equals(0);
        counter = counter + 5;
        test.assert("modified counter").that(counter).equals(5);
    });
});

// Array and Object testing
test.describe("Array and Object Testing", func() {
    test.it("should test array properties", func() {
        let arr = [1, 2, 3, 4, 5];
        
        test.assert("array length").that(len(arr)).equals(5);
        test.assert("array contains").that(arr).contains(3);
        test.assert("array first element").that(arr[0]).equals(1);
        test.assert("array last element").that(arr[4]).equals(5);
    });
    
    test.it("should test object properties", func() {
        let person = {
            name: "Alice",
            age: 30,
            city: "New York"
        };
        
        test.assert("object name").that(person.name).equals("Alice");
        test.assert("object age").that(person.age).equals(30);
        test.assert("object has property").that(person).hasProperty("city");
        test.assert("object property value").that(person.city).equals("New York");
    });
});

// String testing
test.describe("String Testing", func() {
    test.it("should test string operations", func() {
        let greeting = "Hello, World!";
        
        test.assert("string length").that(len(greeting)).equals(13);
        test.assert("string contains").that(greeting).contains("World");
        test.assert("string starts with").that(greeting).startsWith("Hello");
        test.assert("string ends with").that(greeting).endsWith("World!");
    });
    
    test.it("should test string manipulation", func() {
        let original = "  R2Lang Testing  ";
        let trimmed = trim(original);
        
        test.assert("trimmed string").that(trimmed).equals("R2Lang Testing");
        test.assert("original unchanged").that(original).equals("  R2Lang Testing  ");
    });
});

// Type testing
test.describe("Type Testing", func() {
    test.it("should test different types", func() {
        test.assert("number type").that(42).isNumber();
        test.assert("string type").that("hello").isString();
        test.assert("boolean type").that(true).isBoolean();
        test.assert("array type").that([1, 2, 3]).isArray();
        test.assert("object type").that({key: "value"}).isObject();
        test.assert("null type").that(null).isNull();
    });
    
    test.it("should test null and undefined", func() {
        let undefinedVar;
        
        test.assert("null value").that(null).isNull();
        test.assert("not null").that("value").isNotNull();
        test.assert("undefined value").that(undefinedVar).isUndefined();
        test.assert("not undefined").that("defined").isNotUndefined();
    });
});

// Run the tests
print("Running Basic Unit Tests for R2Lang...");
let results = test.runTests();

print("\n=== TEST RESULTS ===");
print("Total tests:", results.total);
print("Passed:", results.passed);
print("Failed:", results.failed);
print("Skipped:", results.skipped);
print("Success rate:", (results.passed / results.total * 100) + "%");

if (results.failed > 0) {
    print("\nFailed tests:");
    for (let i = 0; i < len(results.failures); i++) {
        let failure = results.failures[i];
        print("  -", failure.test + ":", failure.message);
    }
}