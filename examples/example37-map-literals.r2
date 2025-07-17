func main() {
    print("=== R2Lang Map Literals Examples ===")
    
    // 1. Basic map literals with different key types
    let user = {name: "Juan", age: 30, active: true}
    print("Basic map:", user)
    print("Map length:", len(user))
    
    // 2. String keys (explicit)
    let config = {"host": "localhost", "port": 8080, "ssl": true}
    print("String keys map:", config)
    print("Config length:", len(config))
    
    // 3. Numeric keys
    let httpCodes = {200: "OK", 404: "Not Found", 500: "Server Error"}
    print("Numeric keys map:", httpCodes)
    print("HTTP codes count:", len(httpCodes))
    
    // 4. Mixed key types
    let mixed = {"name": "Test", age: 25, 1: "first", 2: "second"}
    print("Mixed keys map:", mixed)
    print("Mixed map size:", len(mixed))
    
    // 5. Boolean values with true/false keywords
    let permissions = {read: true, write: false, admin: true, guest: false}
    print("Permissions map:", permissions)
    print("Permission count:", len(permissions))
    
    // 6. Computed keys
    let prefix = "user"
    let suffix = "data"
    let computed = {[prefix + "_id"]: 123, [prefix + "_" + suffix]: "info"}
    print("Computed keys map:", computed)
    print("Computed map length:", len(computed))
    
    // 7. Nested maps (simplified)
    let userInfo = {name: "Alice", age: 28}
    let settingsInfo = {theme: "dark", lang: "es"}
    let flagsInfo = {beta: true, debug: false}
    print("User info:", userInfo)
    print("Settings info:", settingsInfo)
    print("Flags info:", flagsInfo)
    print("User info size:", len(userInfo))
    print("Settings size:", len(settingsInfo))
    
    // 8. Using keys() function
    let sample = {a: 1, b: 2, c: 3, d: 4}
    let mapKeys = keys(sample)
    print("Sample map:", sample)
    print("Map keys:", mapKeys)
    print("Keys count:", len(mapKeys))
    
    // 9. Empty map
    let empty = {}
    print("Empty map:", empty)
    print("Empty map length:", len(empty))
    
    // 10. Dynamic map building
    let dynamic = {}
    dynamic["created"] = now()
    dynamic["version"] = "1.0"
    dynamic["active"] = true
    print("Dynamic map:", dynamic)
    print("Dynamic map size:", len(dynamic))
    
    // 11. Map with array values (simplified)
    let numbers = [1, 2, 3, 4, 5]
    let strings = ["hello", "world"]
    let booleans = [true, false, true]
    let arrayValues = {numbers: numbers, strings: strings, booleans: booleans}
    print("Map with arrays:", arrayValues)
    print("Numbers array length:", len(numbers))
    
    // 12. Accessing map values
    print("\n=== Accessing Map Values ===")
    let person = {name: "Bob", age: 35, city: "Madrid", country: "Spain"}
    print("Person name:", person.name)
    print("Person age:", person["age"])
    print("Person location:", person.city + ", " + person.country)
    
    // 13. Keys function demonstration
    print("\n=== Keys Function Demo ===")
    let demo = {x: 10, y: 20, z: 30, w: 40}
    let allKeys = keys(demo)
    print("Demo map:", demo)
    print("All keys:", allKeys)
    
    // Manual iteration to avoid for-in issue
    print("Key:", allKeys[0] + " -> Value:", demo[allKeys[0]])
    print("Key:", allKeys[1] + " -> Value:", demo[allKeys[1]])
    print("Key:", allKeys[2] + " -> Value:", demo[allKeys[2]])
    print("Key:", allKeys[3] + " -> Value:", demo[allKeys[3]])
    
    print("\n=== Map Operations Summary ===")
    print("Total examples processed: 13")
    print("Map literal syntax supported: ✓")
    print("len() function for maps: ✓")
    print("keys() function: ✓")
    print("Boolean true/false literals: ✓")
}