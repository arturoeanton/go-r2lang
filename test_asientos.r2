// Test minimal para verificar asientos contables
console.log("Testing asientos contables...")

// Array usando push (funciona)
let asientos = []

let asiento1 = {
    id: "AS-001",
    descripcion: "Test Asiento",
    movimientos: []
}

// Agregar movimientos con push
asiento1.movimientos.push({
    cuenta: "1001",
    descripcion: "Caja",
    tipo: "DEBE",
    monto: 100
})

asiento1.movimientos.push({
    cuenta: "4001", 
    descripcion: "Ventas",
    tipo: "HABER",
    monto: 100
})

// Agregar asiento al array
asientos.push(asiento1)

console.log("Asientos creados: " + std.len(asientos))
console.log("Movimientos en asiento 1: " + std.len(asientos[0].movimientos))

// Verificar contenido
let i = 0
while (i < std.len(asientos[0].movimientos)) {
    let mov = asientos[0].movimientos[i]
    console.log("  " + mov.tipo + " - " + mov.cuenta + " - $" + mov.monto)
    i = i + 1
}

// Ahora probar con asignación por índice en array dentro de objeto
console.log("\nTest con asignación por índice en objeto:")
let asiento2 = {
    id: "AS-002",
    descripcion: "Test Asiento 2",
    movimientos: []
}

// Intentar asignar por índice (esto es lo que falla)
asiento2.movimientos[0] = {
    cuenta: "1002",
    descripcion: "Banco",
    tipo: "DEBE", 
    monto: 200
}

console.log("Movimientos en asiento2 después de [0]: " + std.len(asiento2.movimientos))

// Si no funciona, el array seguirá vacío
if (std.len(asiento2.movimientos) == 0) {
    console.log("❌ La asignación por índice NO funciona en arrays dentro de objetos")
    console.log("✅ Solución: usar push() en lugar de asignación por índice")
}