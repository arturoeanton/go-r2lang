// Test nested array access
let asiento = {
    movimientos: []
}

// Test 1: Manual assignment to map property
try {
    std.print("Test 1: Manual assignment to map property")
    asiento.movimientos = [{cuenta: "1234", monto: 100}]
    std.print("✅ Direct assignment worked!")
    std.print("asiento.movimientos: ")
    std.print(asiento.movimientos)
} catch (e) {
    std.print("❌ Error with direct assignment: " + e)
}

// Test 2: Push and reassign
try {
    std.print("\nTest 2: Push and reassign")
    let newMovimientos = asiento.movimientos.push({cuenta: "5678", monto: 200})
    std.print("Push returned: ")
    std.print(newMovimientos)
    
    // Now reassign
    asiento.movimientos = newMovimientos
    std.print("After reassignment: ")
    std.print(asiento.movimientos)
} catch (e) {
    std.print("❌ Error with push and reassign: " + e)
}

// Test 3: Multiple nested levels
try {
    std.print("\nTest 3: Multiple nested levels")
    let empresa = {
        sucursales: [
            {
                nombre: "Sucursal 1",
                empleados: []
            }
        ]
    }
    
    // Try to push to nested array
    empresa.sucursales[0].empleados = [{nombre: "Juan", id: 1}]
    std.print("✅ Nested assignment worked!")
    std.print("Empleados: ")
    std.print(empresa.sucursales[0].empleados)
} catch (e) {
    std.print("❌ Error with multiple nested levels: " + e)
}