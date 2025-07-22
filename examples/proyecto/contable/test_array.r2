// Test array assignment
console.log("Testing arrays...")

let arr = []
console.log("Initial array length: " + std.len(arr))

arr[0] = "first"
console.log("After adding first: " + std.len(arr))

arr[1] = "second"
console.log("After adding second: " + std.len(arr))

arr[std.len(arr)] = "third"
console.log("After adding third: " + std.len(arr))

console.log("Element 0: " + arr[0])
console.log("Element 1: " + arr[1])
console.log("Element 2: " + arr[2])

// Test with object
let data = {
    items: [],
    count: 0
}

data.items[0] = "item1"
data.count = 1
console.log("Object array length: " + std.len(data.items))
console.log("Object count: " + data.count)