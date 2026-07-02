// Smoke test for the regex module
console.log("regex.test: " + regex.test("\\d+", "order-42"))
console.log("regex.match: " + regex.match("\\d+", "order-42-item-7"))

let all = regex.matchAll("\\d+", "order-42-item-7")
console.log("regex.matchAll length: " + std.len(all))
console.log("regex.matchAll[0]: " + all[0])
console.log("regex.matchAll[1]: " + all[1])

let groups = regex.groups("(\\w+)@(\\w+)\\.com", "user@example.com")
console.log("regex.groups full: " + groups[0])
console.log("regex.groups[1]: " + groups[1])
console.log("regex.groups[2]: " + groups[2])

console.log("regex.replace: " + regex.replace("\\d+", "a1b2c3", "X"))
console.log("regex.replaceAll: " + regex.replaceAll("\\d+", "a1b2c3", "X"))

let parts = regex.split(",\\s*", "a, b,c")
console.log("regex.split length: " + std.len(parts))
console.log("regex.split[0]: " + parts[0])
console.log("regex.split[2]: " + parts[2])

console.log("regex.escape: " + regex.escape("a.b*c?"))

let noMatch = regex.match("\\d+", "no numbers here")
if (noMatch == nil) {
    console.log("regex.match no-match: nil as expected")
} else {
    console.log("regex.match no-match: UNEXPECTED " + noMatch)
}
