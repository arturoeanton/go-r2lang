// Test float calculations in template strings
let price = 100
let tax = 0.16

std.print("price = " + price)
std.print("tax = " + tax)
std.print("price * tax = " + (price * tax))
std.print("1 + tax = " + (1 + tax))
std.print("price * (1 + tax) = " + (price * (1 + tax)))

let total = `Price: $${price}, Tax: $${price * tax}, Total: $${price * (1 + tax)}`
std.print("\nTemplate result:")
std.print(total)
std.print("\nExpected:")
std.print("Price: $100, Tax: $16, Total: $116")