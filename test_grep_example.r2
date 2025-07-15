describe("Calculator Tests", func() {
    it("should add two numbers", func() {
        let result = 2 + 3;
        assert.equals(result, 5);
    });
    
    it("should subtract numbers", func() {
        let result = 10 - 4;
        assert.equals(result, 6);
    });
    
    it("should multiply numbers", func() {
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
});