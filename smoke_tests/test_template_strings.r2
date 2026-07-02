// Test template strings
std.print("🧪 Testing Template Strings")
std.print("========================\n")

// Test 1: Basic template string
let name = "R2Lang"
let version = 2025
let msg = `Hello from ${name} version ${version}!`
std.print("1. Basic interpolation:")
std.print("   " + msg)

// Test 2: Expression interpolation
let a = 10
let b = 20
let calc = `The sum of ${a} and ${b} is ${a + b}`
std.print("\n2. Expression interpolation:")
std.print("   " + calc)

// Test 3: Multiline template string
let html = `
<!DOCTYPE html>
<html>
<head>
    <title>${name}</title>
</head>
<body>
    <h1>Welcome to ${name}</h1>
    <p>Version: ${version}</p>
</body>
</html>
`
std.print("\n3. Multiline template:")
std.print(html)

// Test 4: Nested objects in template
let user = {
    name: "Juan",
    age: 30,
    city: "Madrid"
}
let profile = `User Profile:
- Name: ${user.name}
- Age: ${user.age}
- City: ${user.city}`
std.print("\n4. Object properties:")
std.print(profile)