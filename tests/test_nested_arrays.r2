// Test Suite: Nested Arrays and Map Properties
// Validates the implementation of improvement #1 from mejoras_r2lang_dsl.md

std.print("🧪 TEST SUITE: Nested Arrays and Map Properties")
std.print("==============================================\n")

let tests_passed = 0
let tests_failed = 0

func assert(condition, message) {
    if (condition) {
        std.print("✅ " + message)
        tests_passed = tests_passed + 1
    } else {
        std.print("❌ " + message)
        tests_failed = tests_failed + 1
    }
}

// Test 1: Basic map property assignment
std.print("Test 1: Basic map property assignment")
let myObj = { prop: "initial" }
myObj.prop = "updated"
assert(myObj.prop == "updated", "Map property can be updated")

// Test 2: Nested map property assignment
std.print("\nTest 2: Nested map property assignment")
let nested = { 
    level1: { 
        level2: "initial" 
    } 
}
nested.level1.level2 = "updated"
assert(nested.level1.level2 == "updated", "Nested map property can be updated")

// Test 3: Array property in map
std.print("\nTest 3: Array property in map")
let container = { items: [] }
container.items = container.items.push("item1")
container.items = container.items.push("item2")
assert(std.len(container.items) == 2, "Array in map can be updated with push")
assert(container.items[0] == "item1", "First item is correct")
assert(container.items[1] == "item2", "Second item is correct")

// Test 4: Complex nested structure
std.print("\nTest 4: Complex nested structure")
let company = {
    name: "TechCorp",
    departments: []
}

// Add department with employees
company.departments = company.departments.push({
    name: "Engineering",
    employees: []
})

// Add employees to department
company.departments[0].employees = company.departments[0].employees.push({
    name: "Alice",
    role: "Developer"
})

company.departments[0].employees = company.departments[0].employees.push({
    name: "Bob", 
    role: "Manager"
})

assert(std.len(company.departments) == 1, "Department added correctly")
assert(company.departments[0].name == "Engineering", "Department name is correct")
assert(std.len(company.departments[0].employees) == 2, "Employees added correctly")
assert(company.departments[0].employees[0].name == "Alice", "First employee correct")
assert(company.departments[0].employees[1].role == "Manager", "Second employee role correct")

// Test 5: Array element property update
std.print("\nTest 5: Array element property update")
let tasks = [
    { id: 1, status: "pending" },
    { id: 2, status: "pending" }
]
tasks[0].status = "completed"
tasks[1].status = "in_progress"
assert(tasks[0].status == "completed", "Array element property updated")
assert(tasks[1].status == "in_progress", "Second array element property updated")

// Test 6: Mixed updates
std.print("\nTest 6: Mixed property and array updates")
let project = {
    name: "Project X",
    tasks: [],
    metadata: {
        created: "2024-01-01",
        status: "active"
    }
}

project.name = "Project Y"
project.metadata.status = "completed"
project.tasks = project.tasks.push({ id: 1, title: "Task 1" })
project.tasks = project.tasks.push({ id: 2, title: "Task 2" })

assert(project.name == "Project Y", "Project name updated")
assert(project.metadata.status == "completed", "Nested metadata updated")
assert(std.len(project.tasks) == 2, "Tasks array updated")

// Test 7: Accounting use case (from proyecto contable)
std.print("\nTest 7: Accounting use case")
let asiento = {
    id: "AS-001",
    fecha: "2024-01-15",
    movimientos: []
}

// Add movements using the workaround pattern
asiento.movimientos = asiento.movimientos.push({
    cuenta: "1105",
    descripcion: "Clientes",
    tipo: "DEBE",
    monto: 1000
})

asiento.movimientos = asiento.movimientos.push({
    cuenta: "4135",
    descripcion: "Ventas",
    tipo: "HABER",
    monto: 1000
})

assert(std.len(asiento.movimientos) == 2, "Accounting movements added")
assert(asiento.movimientos[0].tipo == "DEBE", "First movement type correct")
assert(asiento.movimientos[1].tipo == "HABER", "Second movement type correct")

// Calculate balance
let totalDebe = 0
let totalHaber = 0
let i = 0
while (i < std.len(asiento.movimientos)) {
    if (asiento.movimientos[i].tipo == "DEBE") {
        totalDebe = totalDebe + asiento.movimientos[i].monto
    } else {
        totalHaber = totalHaber + asiento.movimientos[i].monto
    }
    i = i + 1
}
assert(totalDebe == totalHaber, "Accounting entry is balanced")

// Summary
std.print("\n" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=")
std.print("SUMMARY:")
std.print("Tests passed: " + tests_passed)
std.print("Tests failed: " + tests_failed)
std.print("Total tests: " + (tests_passed + tests_failed))

if (tests_failed == 0) {
    std.print("\n🎉 ALL TESTS PASSED!")
} else {
    std.print("\n⚠️ SOME TESTS FAILED!")
}