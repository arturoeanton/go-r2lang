// Simple test for object property assignment

// Test 1: Variable assignment (should work)
let x = 10
x = 20
std.print("x = " + x)

// Test 2: Object creation and access (should work)
let myObj = { prop: "initial" }
std.print("myObj.prop = " + myObj.prop)

// Test 3: Object property assignment (testing)
try {
    myObj.prop = "updated"
    std.print("✅ Assignment worked! myObj.prop = " + myObj.prop)
} catch (e) {
    std.print("❌ Assignment failed: " + e)
    
    // Workaround: reassign entire object
    myObj = { prop: "updated" }
    std.print("Workaround: myObj.prop = " + myObj.prop)
}