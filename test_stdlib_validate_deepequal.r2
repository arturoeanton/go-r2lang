console.log("== collections.deepEqual / deepClone ==")

let a = { x: 1, y: [1, 2, { z: true }] }
let b = { x: 1, y: [1, 2, { z: true }] }
let c = { x: 1, y: [1, 2, { z: false }] }

console.log("deepEqual(a, b) == true -> " + collections.deepEqual(a, b))
console.log("deepEqual(a, c) == false -> " + collections.deepEqual(a, c))

let cloned = collections.deepClone(a)
cloned.y[0] = 999
console.log("clone independent from original -> " + (a.y[0] == 1) + " / " + (cloned.y[0] == 999))
console.log("deepEqual(a, cloned) == false after mutation -> " + collections.deepEqual(a, cloned))

let cyclic = [1, 2]
cyclic[2] = cyclic
let cyclic2 = [1, 2]
cyclic2[2] = cyclic2
console.log("deepEqual on self-referential arrays -> " + collections.deepEqual(cyclic, cyclic2))

console.log("\n== validate.isEmail / isURL / isIP ==")
console.log("isEmail good -> " + validate.isEmail("user@example.com"))
console.log("isEmail bad -> " + validate.isEmail("not-an-email"))
console.log("isURL good -> " + validate.isURL("https://example.com/path"))
console.log("isURL bad -> " + validate.isURL("not a url"))
console.log("isIP good v4 -> " + validate.isIP("127.0.0.1"))
console.log("isIP good v6 -> " + validate.isIP("::1"))
console.log("isIP bad -> " + validate.isIP("999.999.999.999"))

console.log("\nSmoke test completed")
