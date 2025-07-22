// Test for boolean access issue

func main() {
    console.log("Testing boolean access...");
    
    // Simple boolean test
    let test1 = {
        activa: 1
    };
    
    console.log("Test 1 - activa:", test1.activa);
    
    // Test with true
    let test2 = {
        activa: true
    };
    
    console.log("Test 2 - activa:", test2.activa);
    
    console.log("Done");
}