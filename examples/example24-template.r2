// Template String Example - R2Lang String Templates Demo
// Demonstrates the new template string functionality with ${} interpolation

func main() {
    std.print("=== R2Lang Template Strings Demo ===");
    
    // 1. Basic interpolation
    let name = "R2Lang";
    let version = "2.0";
    let greeting = `Hello from ${name} version ${version}!`;
    std.print(greeting);
    
    // 2. Arithmetic expressions
    let a = 15;
    let b = 25;
    let mathResult = `${a} + ${b} = ${a + b}`;
    std.print(mathResult);
    
    // 3. Multiline templates
    let title = "Project Report";
    let author = "R2Lang Developer";
    let report = `
${title}
==============
Author: ${author}
Date: July 2025

This is a multiline template string that preserves
formatting and allows interpolation throughout.
    `;
    std.print(report);
    
    // 4. HTML generation
    func generateUserCard(userName, userAge, userEmail) {
        let status = "Adult";
        if (userAge < 18) {
            status = "Minor";
        }
        return `<div class="user-card">
    <h2>${userName}</h2>
    <p>Age: ${userAge}</p>
    <p>Email: ${userEmail}</p>
    <p>Status: ${status}</p>
</div>`;
    }
    
    let userData = generateUserCard("Ana García", 28, "ana@example.com");
    std.print("Generated HTML:");
    std.print(userData);
    
    // 5. SQL query generation
    func buildQuery(table, condition, limit) {
        return `SELECT * 
FROM ${table} 
WHERE ${condition} 
LIMIT ${limit};`;
    }
    
    let query = buildQuery("users", "active = true", 10);
    std.print("Generated SQL:");
    std.print(query);
    
    // 6. JSON-like template
    let userId = 123;
    let status = "active";
    let timestamp = "2025-07-15T10:30:00Z";
    let isActive = true;
    let jsonTemplate = `{
    "id": ${userId},
    "status": "${status}",
    "timestamp": "${timestamp}",
    "success": ${isActive}
}`;
    std.print("JSON Template:");
    std.print(jsonTemplate);
    
    // 7. Configuration file generation
    let serverHost = "localhost";
    let serverPort = 8080;
    let debugMode = true;
    
    let config = `# R2Lang Server Configuration
server.host=${serverHost}
server.port=${serverPort}
debug.enabled=${debugMode}
app.name=${name}
app.version=${version}`;
    
    std.print("Configuration:");
    std.print(config);
    
    // 8. Performance demonstration
    let iterations = 100;
    std.print(`Starting performance test with ${iterations} iterations...`);
    
    for (i = 0; i < iterations; i = i + 1) {
        let perfTest = `Iteration ${i}: Result ${i * 2}`;
    }
    
    let perfResult = `Performance test: ${iterations} template strings created successfully`;
    std.print(perfResult);
    
    // 9. Nested templates
    let innerTemplate = `inner content`;
    let outerTemplate = `Outer: ${innerTemplate} - complete`;
    std.print("Nested template result:");
    std.print(outerTemplate);
    
    // 10. Complex expressions
    let x = 10;
    let y = 20;
    let complexExpr = `Complex calculation: ${x * y + (x - y) / 2}`;
    std.print(complexExpr);
    
    // 11. Boolean interpolation
    let isReady = true;
    let notReady = false;
    let boolTest = `Ready: ${isReady}, Not Ready: ${notReady}`;
    std.print(boolTest);
    
    std.print("=== Template Strings Demo Complete ===");
}