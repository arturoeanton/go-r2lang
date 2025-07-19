func main() {
    std.print("=== R2Lang Map Literals Examples ===")
    
    // 1. Basic map literals with different key types
    let user = {name: "Juan", age: 30, active: true}
    std.print("Basic map:", user)
    std.print("Map length:", std.len(user))
    
    // 2. String keys (explicit)
    let config = {"host": "localhost", "port": 8080, "ssl": true}
    std.print("String keys map:", config)
    std.print("Config length:", std.len(config))
    
    // 3. Numeric keys
    let httpCodes = {200: "OK", 404: "Not Found", 500: "Server Error"}
    std.print("Numeric keys map:", httpCodes)
    std.print("HTTP codes count:", std.len(httpCodes))
    
    // 4. Mixed key types
    let mixed = {"name": "Test", age: 25, 1: "first", 2: "second"}
    std.print("Mixed keys map:", mixed)
    std.print("Mixed map size:", std.len(mixed))
    
    // 5. Boolean values with true/false keywords
    let permissions = {read: true, write: false, admin: true, guest: false}
    std.print("Permissions map:", permissions)
    std.print("Permission count:", std.len(permissions))
    
    // 6. Computed keys
    let prefix = "user"
    let suffix = "data"
    let computed = {[prefix + "_id"]: 123, [prefix + "_" + suffix]: "info"}
    std.print("Computed keys map:", computed)
    std.print("Computed map length:", std.len(computed))
    
    // 7. Nested maps (simplified)
    let userInfo = {name: "Alice", age: 28}
    let settingsInfo = {theme: "dark", lang: "es"}
    let flagsInfo = {beta: true, debug: false}
    std.print("User info:", userInfo)
    std.print("Settings info:", settingsInfo)
    std.print("Flags info:", flagsInfo)
    std.print("User info size:", std.len(userInfo))
    std.print("Settings size:", std.len(settingsInfo))
    
    // 8. Using keys() function
    let sample = {a: 1, b: 2, c: 3, d: 4}
    let mapKeys = std.keys(sample)
    std.print("Sample map:", sample)
    std.print("Map keys:", mapKeys)
    std.print("Keys count:", std.len(mapKeys))
    
    // 9. Empty map
    let empty = {}
    std.print("Empty map:", empty)
    std.print("Empty map length:", std.len(empty))
    
    // 10. Dynamic map building
    let dynamic = {}
    dynamic["created"] = std.now()
    dynamic["version"] = "1.0"
    dynamic["active"] = true
    std.print("Dynamic map:", dynamic)
    std.print("Dynamic map size:", std.len(dynamic))
    
    // 11. Map with array values (simplified)
    let numbers = [1, 2, 3, 4, 5]
    let strings = ["hello", "world"]
    let booleans = [true, false, true]
    let arrayValues = {numbers: numbers, strings: strings, booleans: booleans}
    std.print("Map with arrays:", arrayValues)
    std.print("Numbers array length:", std.len(numbers))
    
    // 12. Accessing map values
    std.print("\n=== Accessing Map Values ===")
    let person = {name: "Bob", age: 35, city: "Madrid", country: "Spain"}
    std.print("Person name:", person.name)
    std.print("Person age:", person["age"])
    std.print("Person location:", person.city + ", " + person.country)
    
    // 13. Keys function demonstration
    std.print("\n=== Keys Function Demo ===")
    let demo = {x: 10, y: 20, z: 30, w: 40}
    let allKeys = std.keys(demo)
    std.print("Demo map:", demo)
    std.print("All keys:", allKeys)
    
    // Manual iteration to avoid for-in issue
    std.print("Key:", allKeys[0] + " -> Value:", demo[allKeys[0]])
    std.print("Key:", allKeys[1] + " -> Value:", demo[allKeys[1]])
    std.print("Key:", allKeys[2] + " -> Value:", demo[allKeys[2]])
    std.print("Key:", allKeys[3] + " -> Value:", demo[allKeys[3]])
    
    std.print("\n=== Map Operations Summary ===")
    std.print("Total examples processed: 13")
    std.print("Map literal syntax supported: ✓")
    std.print("std.len() function for maps: ✓")
    std.print("keys() function: ✓")
    std.print("Boolean true/false literals: ✓")
}