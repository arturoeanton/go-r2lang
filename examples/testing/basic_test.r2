// Example test file for R2Lang Unit Testing Framework
// This demonstrates the new testing syntax that replaces BDD

// Basic describe/it structure
describe("Basic Math Operations", func() {
    it("should add two numbers correctly", func() {
        let result = 2 + 3;
        assert.equals(result, 5);
    });
    
    it("should subtract numbers correctly", func() {
        let result = 10 - 4;
        assert.equals(result, 6);
    });
    
    it("should multiply numbers correctly", func() {
        let result = 3 * 4;
        assert.equals(result, 12);
    });
});

describe("String Operations", func() {
    it("should concatenate strings", func() {
        let greeting = "Hello" + " " + "World";
        assert.equals(greeting, "Hello World");
    });
    
    it("should check string length", func() {
        let text = "R2Lang";
        assert.hasLength(text, 6);
    });
    
    it("should find substrings", func() {
        let text = "R2Lang Testing Framework";
        assert.contains(text, "Testing");
        assert.notContains(text, "Python");
    });
});

describe("Array Operations", func() {
    it("should create and access arrays", func() {
        let numbers = [1, 2, 3, 4, 5];
        assert.hasLength(numbers, 5);
        assert.equals(numbers[0], 1);
        assert.equals(numbers[4], 5);
    });
    
    it("should detect empty arrays", func() {
        let empty = [];
        assert.empty(empty);
        assert.hasLength(empty, 0);
    });
    
    it("should detect non-empty arrays", func() {
        let data = [1, 2, 3];
        assert.notEmpty(data);
        assert.greater(data.length, 0);
    });
});

describe("Boolean Logic", func() {
    it("should handle true values", func() {
        assert.true(true);
        assert.true(1);
        assert.true("hello");
        assert.true([1, 2, 3]);
    });
    
    it("should handle false values", func() {
        assert.false(false);
        assert.false(0);
        assert.false("");
        assert.false([]);
    });
    
    it("should handle null checks", func() {
        let value = null;
        assert.nil(value);
        
        let notNull = "something";
        assert.notNil(notNull);
    });
});

describe("Functions and Variables", func() {
    it("should define and call functions", func() {
        func multiply(a, b) {
            return a * b;
        }
        
        let result = multiply(6, 7);
        assert.equals(result, 42);
    });
    
    it("should handle variable scoping", func() {
        let outer = "outer";
        
        func innerFunction() {
            let inner = "inner";
            return outer + " " + inner;
        }
        
        let combined = innerFunction();
        assert.equals(combined, "outer inner");
    });
});

describe("Error Handling", func() {
    it("should catch expected errors", func() {
        assert.panics(func() {
            throw "Expected error";
        });
    });
    
    it("should not panic for valid operations", func() {
        assert.notPanics(func() {
            let x = 1 + 1;
            return x;
        });
    });
});

describe("Advanced Assertions", func() {
    it("should compare numbers", func() {
        assert.greater(10, 5);
        assert.greaterOrEqual(10, 10);
        assert.less(5, 10);
        assert.lessOrEqual(5, 5);
    });
    
    it("should handle inequality", func() {
        assert.notEquals(5, 10);
        assert.notEquals("hello", "world");
        assert.notEquals(true, false);
    });
});