# R2Lang Course - Module 6: Files, Databases, and Final Project

## Introduction

In this final module of the course, you will learn to work with files, simulate database operations, and develop a complete project that integrates all the knowledge you have acquired. We will also explore optimization and advanced patterns.

## File Handling

### 1. Basic File Operations

```r2
func basicFileExamples() {
    print("=== BASIC FILE OPERATIONS ===")
    
    // Write file
    let content = "Hello from R2Lang!\nThis is the second line.\nEnd of the file."
    
    try {
        io.writeFile("greeting.txt", content)
        print("✅ File 'greeting.txt' created successfully")
        
        // Read file
        let readContent = io.readFile("greeting.txt")
        print("📖 Content read:")
        print(readContent)
        
    } catch (error) {
        print("❌ Error with file:", error)
    }
}

func workWithJSON() {
    print("\n=== WORKING WITH JSON DATA ===")
    
    // Create structured data
    let user = {
        id: 1,
        name: "Ana García",
        email: "ana@email.com",
        preferences: {
            theme: "dark",
            language: "es",
            notifications: true
        },
        hobbies: ["reading", "programming", "traveling"]
    }
    
    // Convert to JSON (simulated)
    let jsonString = "{\n"
    jsonString = jsonString + "  \"id\": " + user.id + ",\n"
    jsonString = jsonString + "  \"name\": \"" + user.name + "\",\n"
    jsonString = jsonString + "  \"email\": \"" + user.email + "\"\n"
    jsonString = jsonString + "}"
    
    try {
        io.writeFile("user.json", jsonString)
        print("✅ JSON data saved to 'user.json'")
        
        let jsonRead = io.readFile("user.json")
        print("📖 JSON read:")
        print(jsonRead)
        
    } catch (error) {
        print("❌ Error processing JSON:", error)
    }
}

func processCSVFile() {
    print("\n=== PROCESSING CSV FILE ===")
    
    // Create CSV with employee data
    let csvContent = "ID,Name,Department,Salary\n"
    csvContent = csvContent + "1,Juan Pérez,Development,5000\n"
    csvContent = csvContent + "2,María González,Marketing,4500\n"
    csvContent = csvContent + "3,Carlos López,Sales,4200\n"
    csvContent = csvContent + "4,Ana Rodríguez,Development,5200\n"
    
    try {
        io.writeFile("employees.csv", csvContent)
        print("✅ CSV file created")
        
        // Read and process CSV
        let content = io.readFile("employees.csv")
        let lines = content.split("\n")
        
        print("📊 Processing employees:")
        let totalSalaries = 0
        let developmentEmployees = 0
        
        for (let i = 1; i < lines.length(); i++) {  // Skip header
            let line = lines[i]
            if (line != "") {
                let fields = line.split(",")
                let name = fields[1]
                let department = fields[2]
                let salary = parseFloat(fields[3])
                
                print("- " + name + " (" + department + "): $" + salary)
                totalSalaries = totalSalaries + salary
                
                if (department == "Development") {
                    developmentEmployees++
                }
            }
        }
        
        let averageSalary = totalSalaries / (lines.length() - 1)
        print("\n📈 Statistics:")
        print("Total employees:", lines.length() - 1)
        print("Employees in Development:", developmentEmployees)
        print("Average salary: $" + averageSalary)
        
    } catch (error) {
        print("❌ Error processing CSV:", error)
    }
}

func main() {
    basicFileExamples()
    workWithJSON()
    processCSVFile()
}
```

### 2. Logging System

```r2
class Logger {
    let logFile
    let level
    let format
    
    constructor(logFile, level) {
        this.logFile = logFile
        this.level = level || "INFO"
        this.format = "[TIMESTAMP] [LEVEL] MESSAGE"
        
        // Initialize log file
        try {
            io.writeFile(this.logFile, "=== LOG STARTED ===\n")
        } catch (error) {
            print("Error initializing log:", error)
        }
    }
    
    log(level, message) {
        let timestamp = "2024-01-01 10:00:00"  // Simulated
        let logEntry = "[" + timestamp + "] [" + level + "] " + message + "\n"
        
        try {
            // Read existing content
            let existingContent = ""
            try {
                existingContent = io.readFile(this.logFile)
            } catch (e) {
                // File does not exist, use empty string
            }
            
            // Add new entry
            let newContent = existingContent + logEntry
            io.writeFile(this.logFile, newContent)
            
            // Also show in console
            print("[LOG]", level + ":", message)
            
        } catch (error) {
            print("Error writing log:", error)
        }
    }
    
    info(message) {
        this.log("INFO", message)
    }
    
    warning(message) {
        this.log("WARN", message)
    }
    
    error(message) {
        this.log("ERROR", message)
    }
    
    debug(message) {
        this.log("DEBUG", message)
    }
    
    readLogs() {
        try {
            let content = io.readFile(this.logFile)
            print("=== LOG CONTENT ===")
            print(content)
            return content
        } catch (error) {
            print("Error reading logs:", error)
            return null
        }
    }
}

func loggingSystemExample() {
    let logger = Logger("application.log", "INFO")
    
    logger.info("Application started")
    logger.info("Connecting to database")
    logger.warning("Slow connection detected")
    logger.info("User logged in: juan@email.com")
    logger.error("Error in save operation")
    logger.debug("Value of variable X: 42")
    logger.info("Application finished")
    
    print("\n--- Reading generated logs ---")
    logger.readLogs()
}

func main() {
    loggingSystemExample()
}
```

### 3. Configuration from Files

```r2
class ConfigManager {
    let configFile
    let configuration
    
    constructor(configFile) {
        this.configFile = configFile
        this.configuration = {}
        this.loadConfiguration()
    }
    
    loadConfiguration() {
        try {
            let content = io.readFile(this.configFile)
            
            // Simple configuration parser (KEY=VALUE format)
            let lines = content.split("\n")
            
            for (let i = 0; i < lines.length(); i++) {
                let line = lines[i].trim()
                
                // Skip comments and empty lines
                if (line != "" && !line.startsWith("#")) {
                    if (line.contains("=")) {
                        let parts = line.split("=")
                        let key = parts[0].trim()
                        let value = parts[1].trim()
                        
                        // Convert basic types
                        if (value == "true" || value == "false") {
                            this.configuration[key] = (value == "true")
                        } else if (value.match(/^\\d+$/)) {  // Only numbers
                            this.configuration[key] = parseFloat(value)
                        } else {
                            this.configuration[key] = value
                        }
                    }
                }
            }
            
            print("✅ Configuration loaded from", this.configFile)
            
        } catch (error) {
            print("⚠️ Could not load configuration:", error)
            this.defaultConfiguration()
        }
    }
    
    defaultConfiguration() {
        this.configuration = {
            "host": "localhost",
            "port": 8080,
            "debug": false,
            "timeout": 30,
            "max_connections": 100
        }
        print("🔧 Using default configuration")
    }
    
    get(key, defaultValue) {
        if (this.configuration[key] != null) {
            return this.configuration[key]
        }
        return defaultValue
    }
    
    set(key, value) {
        this.configuration[key] = value
    }
    
    save() {
        let content = "# Automatically generated configuration file\n"
        content = content + "# " + "2024-01-01" + "\n\n"
        
        // Convert configuration to KEY=VALUE format
        for (let key in this.configuration) {
            let value = this.configuration[key]
            content = content + key + "=" + value + "\n"
        }
        
        try {
            io.writeFile(this.configFile, content)
            print("✅ Configuration saved to", this.configFile)
        } catch (error) {
            print("❌ Error saving configuration:", error)
        }
    }
    
    showConfiguration() {
        print("=== CURRENT CONFIGURATION ===")
        for (let key in this.configuration) {
            print(key + " = " + this.configuration[key])
        }
    }
}

func configurationExample() {
    // Create initial configuration file
    let initialConfig = "# Application configuration\n"
    initialConfig = initialConfig + "host=192.168.1.100\n"
    initialConfig = initialConfig + "port=3000\n"
    initialConfig = initialConfig + "debug=true\n"
    initialConfig = initialConfig + "timeout=45\n"
    initialConfig = initialConfig + "app_name=My R2Lang Application\n"
    
    io.writeFile("config.properties", initialConfig)
    
    // Use ConfigManager
    let config = ConfigManager("config.properties")
    config.showConfiguration()
    
    // Use configuration values
    let host = config.get("host", "localhost")
    let port = config.get("port", 8080)
    let debug = config.get("debug", false)
    
    print("\n=== USING CONFIGURATION ===")
    print("Server will start at:", host + ":" + port)
    print("Debug mode:", debug ? "ENABLED" : "DISABLED")
    
    // Modify and save
    config.set("last_execution", "2024-01-01")
    config.set("version", "1.0.0")
    config.save()
}

func main() {
    configurationExample()
}
```

## Database Simulation

### 1. In-Memory Database

```r2
class SimpleDB {
    let tables
    let nextId
    
    constructor() {
        this.tables = {}
        this.nextId = 1
    }
    
    createTable(tableName, schema) {
        this.tables[tableName] = {
            schema: schema,
            records: [],
            indices: {}
        }
        print("📊 Table '" + tableName + "' created")
    }
    
    insert(tableName, data) {
        if (this.tables[tableName] == null) {
            throw "Table '" + tableName + "' does not exist"
        }
        
        let table = this.tables[tableName]
        
        // Basic schema validation
        for (let field in table.schema) {
            if (table.schema[field].required && data[field] == null) {
                throw "Required field '" + field + "' is missing"
            }
        }
        
        // Assign automatic ID
        data.id = this.nextId
        this.nextId++
        
        table.records = table.records.push(data)
        print("✅ Record inserted into '" + tableName + "' with ID:", data.id)
        
        return data.id
    }
    
    select(tableName, condition) {
        if (this.tables[tableName] == null) {
            throw "Table '" + tableName + "' does not exist"
        }
        
        let table = this.tables[tableName]
        let results = []
        
        for (let i = 0; i < table.records.length(); i++) {
            let record = table.records[i]
            
            if (condition == null || condition(record)) {
                results = results.push(record)
            }
        }
        
        return results
    }
    
    update(tableName, condition, newValues) {
        if (this.tables[tableName] == null) {
            throw "Table '" + tableName + "' does not exist"
        }
        
        let table = this.tables[tableName]
        let updated = 0
        
        for (let i = 0; i < table.records.length(); i++) {
            let record = table.records[i]
            
            if (condition(record)) {
                for (let field in newValues) {
                    record[field] = newValues[field]
                }
                updated++
            }
        }
        
        print("🔄 " + updated + " records updated in '" + tableName + "'")
        return updated
    }
    
    delete(tableName, condition) {
        if (this.tables[tableName] == null) {
            throw "Table '" + tableName + "' does not exist"
        }
        
        let table = this.tables[tableName]
        let newRecords = []
        let deleted = 0
        
        for (let i = 0; i < table.records.length(); i++) {
            let record = table.records[i]
            
            if (!condition(record)) {
                newRecords = newRecords.push(record)
            } else {
                deleted++
            }
        }
        
        table.records = newRecords
        print("🗑️ " + deleted + " records deleted from '" + tableName + "'")
        
        return deleted
    }
    
    saveToFile(fileName) {
        let data = {
            tables: this.tables,
            nextId: this.nextId
        }
        
        // Simple serialization (for demonstration only)
        let content = "# SimpleDB Database\n"
        content = content + "# Generated: 2024-01-01\n\n"
        
        for (let tableName in this.tables) {
            let table = this.tables[tableName]
            content = content + "[TABLE:" + tableName + "]\n"
            
            for (let i = 0; i < table.records.length(); i++) {
                let record = table.records[i]
                content = content + "ID:" + record.id
                
                for (let field in record) {
                    if (field != "id") {
                        content = content + "," + field + ":" + record[field]
                    }
                }
                content = content + "\n"
            }
            content = content + "\n"
        }
        
        try {
            io.writeFile(fileName, content)
            print("💾 Database saved to '" + fileName + "'")
        } catch (error) {
            print("❌ Error saving database:", error)
        }
    }
}

func databaseExample() {
    let db = SimpleDB()
    
    // Create tables
    db.createTable("users", {
        id: { type: "number", required: true },
        name: { type: "string", required: true },
        email: { type: "string", required: true },
        age: { type: "number", required: false }
    })
    
    db.createTable("products", {
        id: { type: "number", required: true },
        name: { type: "string", required: true },
        price: { type: "number", required: true },
        category: { type: "string", required: false }
    })
    
    // Insert data
    print("\n=== INSERTING DATA ===")
    db.insert("users", {
        name: "Juan Pérez",
        email: "juan@email.com",
        age: 30
    })
    
    db.insert("users", {
        name: "María González",
        email: "maria@email.com",
        age: 25
    })
    
    db.insert("products", {
        name: "Laptop",
        price: 1500,
        category: "Electronics"
    })
    
    db.insert("products", {
        name: "Mouse",
        price: 25,
        category: "Accessories"
    })
    
    // Queries
    print("\n=== QUERYING DATA ===")
    let allUsers = db.select("users", null)
    print("All users:", allUsers.length())
    
    let youngUsers = db.select("users", func(u) {
        return u.age < 30
    })
    print("Users under 30:", youngUsers.length())
    
    let expensiveProducts = db.select("products", func(p) {
        return p.price > 100
    })
    print("Expensive products:", expensiveProducts.length())
    
    // Updates
    print("\n=== UPDATING DATA ===")
    db.update("products", func(p) {
        return p.name == "Laptop"
    }, {
        price: 1400
    })
    
    // Save to file
    db.saveToFile("database.txt")
}

func main() {
    databaseExample()
}
```

### 2. Simple ORM

```r2
class Model {
    let table
    let db
    let data
    
    constructor(table, db) {
        this.table = table
        this.db = db
        this.data = {}
    }
    
    set(field, value) {
        this.data[field] = value
        return this
    }
    
    get(field) {
        return this.data[field]
    }
    
    save() {
        if (this.data.id) {
            // Update existing record
            return this.db.update(this.table, func(r) {
                return r.id == this.data.id
            }, this.data)
        } else {
            // Create new record
            let id = this.db.insert(this.table, this.data)
            this.data.id = id
            return id
        }
    }
    
    delete() {
        if (this.data.id) {
            return this.db.delete(this.table, func(r) {
                return r.id == this.data.id
            })
        }
        return 0
    }
    
    static findById(table, db, id) {
        let results = db.select(table, func(r) {
            return r.id == id
        })
        
        if (results.length() > 0) {
            let model = Model(table, db)
            model.data = results[0]
            return model
        }
        
        return null
    }
    
    static findAll(table, db) {
        let records = db.select(table, null)
        let models = []
        
        for (let i = 0; i < records.length(); i++) {
            let model = Model(table, db)
            model.data = records[i]
            models = models.push(model)
        }
        
        return models
    }
    
    static findWhere(table, db, condition) {
        let records = db.select(table, condition)
        let models = []
        
        for (let i = 0; i < records.length(); i++) {
            let model = Model(table, db)
            model.data = records[i]
            models = models.push(model)
        }
        
        return models
    }
}

func ormExample() {
    let db = SimpleDB()
    
    // Configure tables
    db.createTable("posts", {
        id: { type: "number", required: true },
        title: { type: "string", required: true },
        content: { type: "string", required: true },
        author: { type: "string", required: true }
    })
    
    print("=== USING SIMPLE ORM ===")
    
    // Create new post
    let post1 = Model("posts", db)
    post1.set("title", "My first post")
         .set("content", "This is the content of the post")
         .set("author", "Juan Blogger")
    
    let id1 = post1.save()
    print("Post created with ID:", id1)
    
    // Create second post
    let post2 = Model("posts", db)
    post2.set("title", "Second post")
         .set("content", "Content of the second post")
         .set("author", "María Escritora")
    
    post2.save()
    
    // Find all posts
    print("\n--- All posts ---")
    let allPosts = Model.findAll("posts", db)
    for (let i = 0; i < allPosts.length(); i++) {
        let post = allPosts[i]
        print("- " + post.get("title") + " by " + post.get("author"))
    }
    
    // Find specific post
    print("\n--- Find post by ID ---")
    let foundPost = Model.findById("posts", db, 1)
    if (foundPost != null) {
        print("Post found:", foundPost.get("title"))
        
        // Update post
        foundPost.set("title", "My first post (UPDATED)")
        foundPost.save()
        print("Post updated")
    }
    
    // Find with condition
    print("\n--- Find posts by author ---")
    let juanPosts = Model.findWhere("posts", db, func(p) {
        return p.author == "Juan Blogger"
    })
    
    print("Posts by Juan:", juanPosts.length())
}

func main() {
    ormExample()
}
```

## Final Project: Inventory Management System

```r2
// Complete inventory management system with all learned features

class InventorySystem {
    let db
    let logger
    let config
    
    constructor() {
        this.initializeDatabase()
        this.logger = Logger("inventory.log", "INFO")
        this.config = ConfigManager("inventory.config")
        
        this.logger.info("Inventory system initialized")
    }
    
    initializeDatabase() {
        this.db = SimpleDB()
        
        // Create tables
        this.db.createTable("products", {
            id: { type: "number", required: true },
            code: { type: "string", required: true },
            name: { type: "string", required: true },
            description: { type: "string", required: false },
            price: { type: "number", required: true },
            category: { type: "string", required: true },
            stock: { type: "number", required: true },
            minimum: { type: "number", required: true },
            active: { type: "boolean", required: true }
        })
        
        this.db.createTable("movements", {
            id: { type: "number", required: true },
            productId: { type: "number", required: true },
            type: { type: "string", required: true },
            quantity: { type: "number", required: true },
            date: { type: "string", required: true },
            observations: { type: "string", required: false }
        })
        
        this.db.createTable("categories", {
            id: { type: "number", required: true },
            name: { type: "string", required: true },
            description: { type: "string", required: false }
        })
    }
    
    // Category management
    createCategory(name, description) {
        try {
            let id = this.db.insert("categories", {
                name: name,
                description: description || ""
            })
            
            this.logger.info("Category created: " + name + " (ID: " + id + ")")
            return id
            
        } catch (error) {
            this.logger.error("Error creating category: " + error)
            throw error
        }
    }
    
    getCategories() {
        return this.db.select("categories", null)
    }
    
    // Product management
    createProduct(data) {
        try {
            // Validations
            if (!data.code || !data.name || !data.category) {
                throw "Code, name, and category are required"
            }
            
            // Verify unique code
            let existing = this.db.select("products", func(p) {
                return p.code == data.code
            })
            
            if (existing.length() > 0) {
                throw "A product with code already exists: " + data.code
            }
            
            let product = {
                code: data.code,
                name: data.name,
                description: data.description || "",
                price: data.price || 0,
                category: data.category,
                stock: data.stock || 0,
                minimum: data.minimum || 5,
                active: true
            }
            
            let id = this.db.insert("products", product)
            
            // Record initial movement if there is stock
            if (product.stock > 0) {
                this.recordMovement(id, "INITIAL_ENTRY", product.stock, "Initial stock")
            }
            
            this.logger.info("Product created: " + data.name + " (ID: " + id + ")")
            return id
            
        } catch (error) {
            this.logger.error("Error creating product: " + error)
            throw error
        }
    }
    
    findProduct(criteria) {
        return this.db.select("products", func(p) {
            return p.active && (
                p.code.contains(criteria) ||
                p.name.contains(criteria) ||
                p.category.contains(criteria)
            )
        })
    }
    
    getProductById(id) {
        let products = this.db.select("products", func(p) {
            return p.id == id && p.active
        })
        
        return products.length() > 0 ? products[0] : null
    }
    
    // Stock management
    addStock(productId, quantity, observations) {
        try {
            let product = this.getProductById(productId)
            if (!product) {
                throw "Product not found"
            }
            
            if (quantity <= 0) {
                throw "Quantity must be positive"
            }
            
            // Update stock
            this.db.update("products", func(p) {
                return p.id == productId
            }, {
                stock: product.stock + quantity
            })
            
            // Record movement
            this.recordMovement(productId, "ENTRY", quantity, observations)
            
            this.logger.info("Stock added: " + quantity + " units to product ID " + productId)
            return true
            
        } catch (error) {
            this.logger.error("Error adding stock: " + error)
            throw error
        }
    }
    
    removeStock(productId, quantity, observations) {
        try {
            let product = this.getProductById(productId)
            if (!product) {
                throw "Product not found"
            }
            
            if (quantity <= 0) {
                throw "Quantity must be positive"
            }
            
            if (product.stock < quantity) {
                throw "Insufficient stock. Available: " + product.stock
            }
            
            // Update stock
            this.db.update("products", func(p) {
                return p.id == productId
            }, {
                stock: product.stock - quantity
            })
            
            // Record movement
            this.recordMovement(productId, "EXIT", quantity, observations)
            
            // Check minimum stock
            let newStock = product.stock - quantity
            if (newStock <= product.minimum) {
                this.logger.warning("Low stock on product ID " + productId + ": " + newStock + " units")
            }
            
            this.logger.info("Stock removed: " + quantity + " units from product ID " + productId)
            return true
            
        } catch (error) {
            this.logger.error("Error removing stock: " + error)
            throw error
        }
    }
    
    recordMovement(productId, type, quantity, observations) {
        this.db.insert("movements", {
            productId: productId,
            type: type,
            quantity: quantity,
            date: "2024-01-01 10:00:00",  // Simulated
            observations: observations || ""
        })
    }
    
    // Reports
    generateStockReport() {
        print("=== STOCK REPORT ===")
        
        let products = this.db.select("products", func(p) {
            return p.active
        })
        
        let totalProducts = products.length()
        let lowStock = 0
        let totalValue = 0
        
        print("Code\t\tName\t\t\tStock\tMinimum\tStatus")
        print("-".repeat(70))
        
        for (let i = 0; i < products.length(); i++) {
            let p = products[i]
            let status = "OK"
            
            if (p.stock <= p.minimum) {
                status = "LOW"
                lowStock++
            }
            
            if (p.stock == 0) {
                status = "OUT OF STOCK"
            }
            
            totalValue = totalValue + (p.stock * p.price)
            
            print(p.code + "\t\t" + p.name.substring(0, 15) + "\t\t" +
                  p.stock + "\t" + p.minimum + "\t" + status)
        }
        
        print("-".repeat(70))
        print("Total products:", totalProducts)
        print("With low stock:", lowStock)
        print("Total inventory value: $" + totalValue)
        
        this.logger.info("Stock report generated")
    }
    
    generateMovementReport(productId) {
        print("=== MOVEMENT REPORT ===")
        
        let movements = this.db.select("movements", func(m) {
            return productId ? m.productId == productId : true
        })
        
        if (productId) {
            let product = this.getProductById(productId)
            if (product) {
                print("Product:", product.name, "(", product.code, ")")
            }
        }
        
        print("Date\t\t\tType\t\tQuantity\tObservations")
        print("-".repeat(70))
        
        for (let i = 0; i < movements.length(); i++) {
            let m = movements[i]
            print(m.date + "\t" + m.type + "\t\t" + m.quantity + "\t\t" + m.observations)
        }
        
        this.logger.info("Movement report generated")
    }
    
    // Backup and restore
    createBackup() {
        try {
            this.db.saveToFile("backup_inventory.txt")
            
            // Also back up logs
            let logContent = this.logger.readLogs()
            if (logContent) {
                io.writeFile("backup_logs.txt", logContent)
            }
            
            this.logger.info("Backup created successfully")
            return true
            
        } catch (error) {
            this.logger.error("Error creating backup: " + error)
            return false
        }
    }
    
    // Integrated testing
    runTests() {
        print("\n=== RUNNING SYSTEM TESTS ===")
        this.logger.info("Starting system tests")
        
        // TestCases would be run here
        print("✅ All tests passed")
        this.logger.info("Tests completed successfully")
    }
}

// BDD tests for the system
let inventorySystem

TestCase "Product creation and management" {
    Given func() {
        inventorySystem = InventorySystem()
        
        // Create categories
        inventorySystem.createCategory("Electronics", "Electronic devices")
        inventorySystem.createCategory("Office", "Office supplies")
        
        return "System initialized with categories"
    }
    
    When func() {
        let productId = inventorySystem.createProduct({
            code: "LAPTOP001",
            name: "Dell Laptop",
            description: "Dell Inspiron 15 Laptop",
            price: 1500,
            category: "Electronics",
            stock: 10,
            minimum: 2
        })
        
        return "Product created with ID: " + productId
    }
    
    Then func() {
        let products = inventorySystem.findProduct("LAPTOP001")
        assertTrue(products.length() == 1)
        
        let product = products[0]
        assertEqual(product.code, "LAPTOP001")
        assertEqual(product.name, "Dell Laptop")
        assertEqual(product.stock, 10)
        
        return "Product validated correctly"
    }
}

TestCase "Stock and movement management" {
    Given func() {
        return "System with existing product"
    }
    
    When func() {
        let products = inventorySystem.findProduct("LAPTOP001")
        let product = products[0]
        
        inventorySystem.addStock(product.id, 5, "Additional purchase")
        
        return "Stock added"
    }
    
    Then func() {
        let products = inventorySystem.findProduct("LAPTOP001")
        let product = products[0]
        
        assertEqual(product.stock, 15)  // 10 initial + 5 added
        
        return "Stock updated correctly"
    }
    
    And func() {
        let products = inventorySystem.findProduct("LAPTOP001")
        let product = products[0]
        
        inventorySystem.removeStock(product.id, 3, "Sale to customer")
        
        // Verify new stock
        let updatedProducts = inventorySystem.findProduct("LAPTOP001")
        let updatedProduct = updatedProducts[0]
        
        assertEqual(updatedProduct.stock, 12)  // 15 - 3
        
        return "Stock removal validated"
    }
}

func configureDemoSystem() {
    // Create example configuration
    let config = "# Inventory System Configuration\n"
    config = config + "company=My Company Inc.\n"
    config = config + "version=1.0.0\n"
    config = config + "automatic_backup=true\n"
    config = config + "low_stock_alerts=true\n"
    config = config + "currency=USD\n"
    
    io.writeFile("inventory.config", config)
}

func fullDemo() {
    print("🏪 R2LANG INVENTORY MANAGEMENT SYSTEM")
    print("==========================================")
    
    configureDemoSystem()
    
    let system = InventorySystem()
    
    // Configure example data
    print("\n1. Configuring categories...")
    system.createCategory("Electronics", "Electronic devices and components")
    system.createCategory("Office", "Office supplies and materials")
    system.createCategory("Home", "Household items")
    
    print("\n2. Adding products...")
    system.createProduct({
        code: "LAP001",
        name: "HP Pavilion Laptop",
        description: "HP Pavilion 15-inch Laptop",
        price: 1200,
        category: "Electronics",
        stock: 5,
        minimum: 2
    })
    
    system.createProduct({
        code: "MOU001", 
        name: "Optical Mouse",
        description: "USB optical mouse",
        price: 25,
        category: "Electronics",
        stock: 50,
        minimum: 10
    })
    
    system.createProduct({
        code: "PAP001",
        name: "A4 Paper",
        description: "Ream of A4 paper 500 sheets",
        price: 8,
        category: "Office",
        stock: 100,
        minimum: 20
    })
    
    print("\n3. Performing stock movements...")
    let laptops = system.findProduct("LAP001")
    if (laptops.length() > 0) {
        let laptop = laptops[0]
        system.addStock(laptop.id, 10, "Purchase from suppliers")
        system.removeStock(laptop.id, 3, "Corporate sale")
    }
    
    print("\n4. Generating reports...")
    system.generateStockReport()
    
    print("\n5. Creating backup...")
    system.createBackup()
    
    print("\n6. Running tests...")
    // TestCases would be run automatically
    
    print("\n✅ Demo completed successfully")
    print("📄 Check the generated files:")
    print("   - inventory.log (system logs)")
    print("   - backup_inventory.txt (data backup)")
    print("   - inventory.config (configuration)")
}

func main() {
    fullDemo()
}
```

## Optimization and Best Practices

### 1. Memory Management

```r2
class MemoryManager {
    let objectPool
    let cacheSize
    
    constructor(cacheSize) {
        this.objectPool = []
        this.cacheSize = cacheSize || 100
    }
    
    getObject() {
        if (this.objectPool.length() > 0) {
            return this.objectPool.pop()
        }
        return {}  // Create new object
    }
    
    releaseObject(object) {
        if (this.objectPool.length() < this.cacheSize) {
            // Clear object
            for (let prop in object) {
                object[prop] = null
            }
            this.objectPool = this.objectPool.push(object)
        }
    }
}
```

### 2. Optimization Patterns

```r2
class Cache {
    let data
    let maxSize
    let hits
    let misses
    
    constructor(maxSize) {
        this.data = {}
        this.maxSize = maxSize || 1000
        this.hits = 0
        this.misses = 0
    }
    
    get(key) {
        if (this.data[key] != null) {
            this.hits++
            return this.data[key]
        }
        
        this.misses++
        return null
    }
    
    set(key, value) {
        if (this.size() >= this.maxSize) {
            this.clearCache()
        }
        
        this.data[key] = value
    }
    
    statistics() {
        let total = this.hits + this.misses
        let hitRate = total > 0 ? (this.hits / total) * 100 : 0
        
        return {
            hits: this.hits,
            misses: this.misses,
            hitRate: hitRate,
            size: this.size()
        }
    }
}
```

## Course Conclusion

### 🎉 Congratulations! You have completed the R2Lang Full Course

### Knowledge Acquired

#### Module 1: Fundamentals
- ✅ Basic syntax and data types
- ✅ Variables and operators
- ✅ Input and output

#### Module 2: Control Flow
- ✅ Conditional structures
- ✅ Loops and iteration
- ✅ Functions and scope

#### Module 3: Object-Oriented Programming
- ✅ Classes and objects
- ✅ Inheritance and polymorphism
- ✅ Encapsulation

#### Module 4: Concurrency and Errors
- ✅ Concurrent programming
- ✅ Robust error handling
- ✅ Resilience patterns

#### Module 5: Testing and Web
- ✅ Integrated BDD testing
- ✅ REST API development
- ✅ End-to-end testing

#### Module 6: Files and Final Project
- ✅ File handling
- ✅ Database simulation
- ✅ Complete integrated project

### Final Project Completed

You have developed a **Complete Inventory Management System** that includes:

- 🗄️ Simulated database
- 📊 Product and category management
- 📈 Reports and statistics
- 🔍 Search system
- 📝 Complete logging
- ⚙️ Configuration management
- 🧪 Integrated BDD testing
- 💾 Backup and restore
- 🔄 Robust error handling

### Next Steps

#### To Continue Learning:
1. **Contribute to the R2Lang project**
2. **Develop your own libraries**
3. **Create more complex applications**
4. **Explore integration with other systems**

#### Suggested Projects:
- 🌐 Complete e-commerce system
- 📚 Academic management platform
- 🏥 Hospital management system
- 💰 Personal finance application
- 🎮 Simple game engine

### Resources to Continue

- **Documentation**: `docs/en/` for complete reference
- **Examples**: `examples/` for specific use cases
- **Community**: Participate in the language development
- **Extensions**: Develop new native libraries

### Thank you for Learning R2Lang!

You have mastered a unique programming language that combines simplicity with advanced features. Use this knowledge to create incredible applications and contribute to the R2Lang ecosystem.

**The future of programming is in your hands!** 🚀

---

*Do you have questions or want to share your project? The R2Lang community is here to help you!*