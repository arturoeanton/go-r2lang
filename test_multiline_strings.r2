// Test multiline strings with triple quotes
std.print("🧪 Testing Multiline Strings")
std.print("===========================\n")

// Test triple double quotes
try {
    let sql = """
    SELECT t.*, a.descripcion
    FROM transacciones t
    JOIN asientos a ON t.id = a.transaccion_id
    WHERE t.fecha BETWEEN ? AND ?
    ORDER BY t.fecha DESC
    """
    std.print("✅ Triple quotes work!")
    std.print("SQL query:")
    std.print(sql)
} catch (e) {
    std.print("❌ Triple quotes not supported yet")
    std.print("Error: " + e)
}

// Alternative: Regular multiline with concatenation
std.print("\n📝 Alternative: String concatenation")
let sql2 = "SELECT t.*, a.descripcion\n" +
           "FROM transacciones t\n" +
           "JOIN asientos a ON t.id = a.transaccion_id\n" +
           "WHERE t.fecha BETWEEN ? AND ?\n" +
           "ORDER BY t.fecha DESC"
std.print(sql2)

// Best option: Template strings for multiline
std.print("\n✅ Best option: Template strings")
let sql3 = `
SELECT t.*, a.descripcion
FROM transacciones t
JOIN asientos a ON t.id = a.transaccion_id
WHERE t.fecha BETWEEN ? AND ?
ORDER BY t.fecha DESC
`
std.print(sql3)