// Test array dynamic growth fix
console.log("Testing array dynamic growth...")

let arr = []
console.log("Initial length: " + std.len(arr))

// Test direct index assignment
arr[0] = "first"
console.log("After arr[0]: length = " + std.len(arr) + ", value = " + arr[0])

arr[2] = "third"
console.log("After arr[2]: length = " + std.len(arr) + ", arr[1] = " + arr[1] + ", arr[2] = " + arr[2])

arr[5] = "sixth"
console.log("After arr[5]: length = " + std.len(arr))

// Test with objects
let data = {
    items: []
}

data.items[0] = "item1"
console.log("Object array after [0]: length = " + std.len(data.items))

data.items[3] = "item4"
console.log("Object array after [3]: length = " + std.len(data.items))

console.log("\nArray contents:")
let i = 0
while (i < std.len(arr)) {
    console.log("  arr[" + i + "] = " + arr[i])
    i = i + 1
}

console.log("\nObject array contents:")
i = 0
while (i < std.len(data.items)) {
    console.log("  data.items[" + i + "] = " + data.items[i])
    i = i + 1
}