// Demostración de arrays anidados con workaround
std.print("🔧 Demostración de Arrays Anidados en R2Lang")
std.print("=========================================\n")

// Crear estructura con array anidado
let asiento = {
    id: "AS-001",
    fecha: "2024-01-15",
    movimientos: []
}

std.print("1️⃣ Estructura inicial creada:")
std.print("   ID: " + asiento.id)
std.print("   Movimientos: " + std.len(asiento.movimientos))

// Workaround: usar push y reasignar
std.print("\n2️⃣ Agregando movimientos...")

// Primer movimiento
asiento.movimientos = asiento.movimientos.push({
    cuenta: "1105",
    descripcion: "Clientes",
    tipo: "DEBE",
    monto: 1160
})
std.print("   ✅ Movimiento 1 agregado")

// Segundo movimiento
asiento.movimientos = asiento.movimientos.push({
    cuenta: "4135",
    descripcion: "Ventas",
    tipo: "HABER",
    monto: 1000
})
std.print("   ✅ Movimiento 2 agregado")

// Tercer movimiento
asiento.movimientos = asiento.movimientos.push({
    cuenta: "2408",
    descripcion: "IVA Débito",
    tipo: "HABER",
    monto: 160
})
std.print("   ✅ Movimiento 3 agregado")

// Mostrar resultado
std.print("\n3️⃣ Asiento completo:")
std.print("   Total movimientos: " + std.len(asiento.movimientos))

let i = 0
let totalDebe = 0
let totalHaber = 0

while (i < std.len(asiento.movimientos)) {
    let mov = asiento.movimientos[i]
    std.print("   - " + mov.descripcion + " (" + mov.tipo + "): $" + mov.monto)
    
    if (mov.tipo == "DEBE") {
        totalDebe = totalDebe + mov.monto
    } else {
        totalHaber = totalHaber + mov.monto
    }
    
    i = i + 1
}

std.print("\n4️⃣ Balance:")
std.print("   Total DEBE:  $" + totalDebe)
std.print("   Total HABER: $" + totalHaber)
std.print("   Cuadrado: " + (totalDebe == totalHaber ? "✅ SI" : "❌ NO"))

// Arrays multinivel
std.print("\n5️⃣ Arrays multinivel:")
let empresa = {
    nombre: "Mi Empresa",
    sucursales: []
}

// Agregar sucursal con empleados
empresa.sucursales = empresa.sucursales.push({
    id: "SUC-001",
    nombre: "Sucursal Norte",
    empleados: []
})

// Agregar empleados a la sucursal
empresa.sucursales[0].empleados = empresa.sucursales[0].empleados.push({
    nombre: "Juan Pérez",
    cargo: "Contador"
})

empresa.sucursales[0].empleados = empresa.sucursales[0].empleados.push({
    nombre: "María García",
    cargo: "Auxiliar"
})

std.print("   Empresa: " + empresa.nombre)
std.print("   Sucursales: " + std.len(empresa.sucursales))
std.print("   - " + empresa.sucursales[0].nombre)
std.print("     Empleados: " + std.len(empresa.sucursales[0].empleados))

let j = 0
while (j < std.len(empresa.sucursales[0].empleados)) {
    let emp = empresa.sucursales[0].empleados[j]
    std.print("     • " + emp.nombre + " - " + emp.cargo)
    j = j + 1
}

std.print("\n✅ Arrays anidados funcionando correctamente con workaround!")