// Test arrays inside objects
console.log("Testing arrays in objects...")

// Test 1: Direct push
let obj1 = { items: [] }
obj1.items.push("item1")
console.log("After push: " + std.len(obj1.items)) // Should be 1

// Test 2: Index assignment (will fail)
let obj2 = { items: [] }
try {
    obj2.items[0] = "item1"
    console.log("Index assignment worked!")
} catch (e) {
    console.log("Index assignment failed: " + e)
}
console.log("After index assign: " + std.len(obj2.items)) // Will be 0

// Test 3: Working solution
console.log("\nWorking solution:")
let asiento = {
    id: "AS-001",
    movimientos: []
}

// Add movements with push
asiento.movimientos.push({
    cuenta: "1001",
    tipo: "DEBE",
    monto: 100
})

asiento.movimientos.push({
    cuenta: "4001", 
    tipo: "HABER",
    monto: 100
})

console.log("Movements: " + std.len(asiento.movimientos))
let i = 0
while (i < std.len(asiento.movimientos)) {
    let mov = asiento.movimientos[i]
    console.log("  " + mov.tipo + " - " + mov.cuenta + " - $" + mov.monto)
    i = i + 1
}