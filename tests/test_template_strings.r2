// Test Suite: Template Literals and Multiline Strings
// Validates the implementation of improvement #2 from mejoras_r2lang_dsl.md

std.print("🧪 TEST SUITE: Template Literals and Multiline Strings")
std.print("=====================================================\n")

let tests_passed = 0
let tests_failed = 0

func assert(condition, message) {
    if (condition) {
        std.print("✅ " + message)
        tests_passed = tests_passed + 1
    } else {
        std.print("❌ " + message)
        tests_failed = tests_failed + 1
    }
}

// Test 1: Basic interpolation
std.print("Test 1: Basic variable interpolation")
let name = "R2Lang"
let greeting = `Hello, ${name}!`
assert(greeting == "Hello, R2Lang!", "Basic interpolation works")

// Test 2: Number interpolation
std.print("\nTest 2: Number interpolation")
let version = 2025
let msg = `Version ${version} is here`
assert(msg == "Version 2025 is here", "Number interpolation works")

// Test 3: Expression interpolation
std.print("\nTest 3: Expression interpolation")
let a = 10
let b = 20
let result = `${a} + ${b} = ${a + b}`
assert(result == "10 + 20 = 30", "Expression interpolation works")

// Test 4: Object property interpolation
std.print("\nTest 4: Object property interpolation")
let user = { name: "Juan", age: 30 }
let info = `User: ${user.name}, Age: ${user.age}`
assert(info == "User: Juan, Age: 30", "Object property interpolation works")

// Test 5: Nested object interpolation
std.print("\nTest 5: Nested object interpolation")
let data = {
    company: {
        name: "TechCorp",
        year: 2020
    }
}
let company_info = `${data.company.name} founded in ${data.company.year}`
assert(company_info == "TechCorp founded in 2020", "Nested object interpolation works")

// Test 6: Array element interpolation
std.print("\nTest 6: Array element interpolation")
let items = ["first", "second", "third"]
let list = `Items: ${items[0]}, ${items[1]}, ${items[2]}`
assert(list == "Items: first, second, third", "Array element interpolation works")

// Test 7: Complex expression interpolation
std.print("\nTest 7: Complex expression interpolation")
let price = 100
let tax = 0.16
let total = `Price: $${price}, Tax: $${price * tax}, Total: $${math.round(price * (1 + tax))}`
assert(total == "Price: $100, Tax: $16, Total: $116", "Complex expression interpolation works")

// Test 8: Multiline template strings
std.print("\nTest 8: Multiline template strings")
let multiline = `Line 1
Line 2
Line 3`
let lines = std.split(multiline, "\n")
assert(std.len(lines) == 3, "Multiline string has 3 lines")
assert(lines[0] == "Line 1", "First line correct")
assert(lines[1] == "Line 2", "Second line correct")
assert(lines[2] == "Line 3", "Third line correct")

// Test 9: Multiline with interpolation
std.print("\nTest 9: Multiline with interpolation")
let title = "My Document"
let author = "John Doe"
let doc = `Title: ${title}
Author: ${author}
Date: 2024-01-15`
assert(std.contains(doc, "Title: My Document"), "Title interpolated in multiline")
assert(std.contains(doc, "Author: John Doe"), "Author interpolated in multiline")

// Test 10: HTML template
std.print("\nTest 10: HTML template")
let page_title = "Welcome"
let content = "Hello World"
let html = `<!DOCTYPE html>
<html>
<head>
    <title>${page_title}</title>
</head>
<body>
    <h1>${content}</h1>
</body>
</html>`
assert(std.contains(html, "<title>Welcome</title>"), "HTML title interpolated")
assert(std.contains(html, "<h1>Hello World</h1>"), "HTML content interpolated")

// Test 11: SQL-like template
std.print("\nTest 11: SQL-like template")
let table = "users"
let column = "name"
let value = "John"
let sql = `SELECT * FROM ${table} WHERE ${column} = '${value}'`
assert(sql == "SELECT * FROM users WHERE name = 'John'", "SQL template works")

// Test 12: JSON-like template
std.print("\nTest 12: JSON-like template")
let id = 123
let status = "active"
let json_template = `{
    "id": ${id},
    "status": "${status}"
}`
assert(std.contains(json_template, '"id": 123'), "JSON number interpolated")
assert(std.contains(json_template, '"status": "active"'), "JSON string interpolated")

// Test 13: Empty interpolation
std.print("\nTest 13: Empty string interpolation")
let empty = ""
let with_empty = `Value: ${empty}!`
assert(with_empty == "Value: !", "Empty string interpolation works")

// Test 14: Special characters
std.print("\nTest 14: Special characters in templates")
let special = "a\tb\nc"
let with_special = `Special: ${special}`
assert(std.contains(with_special, "\t"), "Tab character preserved")
assert(std.contains(with_special, "\n"), "Newline character preserved")

// Summary
std.print("\n" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=")
std.print("SUMMARY:")
std.print("Tests passed: " + tests_passed)
std.print("Tests failed: " + tests_failed)
std.print("Total tests: " + (tests_passed + tests_failed))

if (tests_failed == 0) {
    std.print("\n🎉 ALL TESTS PASSED!")
} else {
    std.print("\n⚠️ SOME TESTS FAILED!")
}