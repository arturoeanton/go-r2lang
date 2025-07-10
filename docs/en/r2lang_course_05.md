# R2Lang Course - Module 5: Testing and Web Development

## Introduction

In this module, you will learn one of the most distinctive features of R2Lang: its integrated testing system with BDD (Behavior Driven Development) syntax, as well as web and API development. These skills will allow you to create complete and well-tested applications.

## Integrated BDD Testing System

### 1. BDD (Behavior Driven Development) Concepts

BDD is a methodology that describes software behavior in a structured natural language:

- **Given**: Establishes the initial context
- **When**: Describes the action being executed
- **Then**: Verifies the expected result
- **And**: Continues the previous step

### 2. Your First TestCase

```r2
// Support functions for testing
func assertEqual(actual, expected) {
    if (actual == expected) {
        print("✅ PASS: Expected value received")
        return true
    } else {
        print("❌ FAIL: Expected", expected, "but received", actual)
        return false
    }
}

func assertTrue(condition) {
    if (condition) {
        print("✅ PASS: Condition is true")
        return true
    } else {
        print("❌ FAIL: Condition is false")
        return false
    }
}

// Function to be tested
func add(a, b) {
    return a + b
}

TestCase "Verify sum of numbers" {
    Given func() {
        print("Preparing data for the sum")
        return "Data prepared"
    }
    
    When func() {
        let result = add(2, 3)
        return "Sum executed: " + result
    }
    
    Then func() {
        let result = add(2, 3)
        assertEqual(result, 5)
        return "Validation completed"
    }
}

func main() {
    print("Running tests...")
}
```

### 3. Advanced TestCase with Setup and Teardown

```r2
// Database simulation
let database = []

func cleanDatabase() {
    database = []
    print("🗑️ Database cleaned")
}

func addUser(name, email) {
    let user = {
        id: database.length() + 1,
        name: name,
        email: email,
        active: true
    }
    database = database.push(user)
    return user
}

func findUser(email) {
    for (let i = 0; i < database.length(); i++) {
        let user = database[i]
        if (user.email == email) {
            return user
        }
    }
    return null
}

func deactivateUser(email) {
    for (let i = 0; i < database.length(); i++) {
        let user = database[i]
        if (user.email == email) {
            user.active = false
            return true
        }
    }
    return false
}

TestCase "Complete user management" {
    Given func() {
        cleanDatabase()
        print("Database prepared for testing")
        return "Setup completed"
    }
    
    When func() {
        let user = addUser("Ana García", "ana@email.com")
        print("User created with ID:", user.id)
        return "User created successfully"
    }
    
    Then func() {
        let user = findUser("ana@email.com")
        assertTrue(user != null)
        assertEqual(user.name, "Ana García")
        assertTrue(user.active)
        return "User found and validated"
    }
    
    And func() {
        let result = deactivateUser("ana@email.com")
        assertTrue(result)
        
        let user = findUser("ana@email.com")
        assertTrue(!user.active)
        return "User deactivated correctly"
    }
}

TestCase "Search for a non-existent user" {
    Given func() {
        cleanDatabase()
        return "Empty database"
    }
    
    When func() {
        let user = findUser("nonexistent@email.com")
        return "Search executed"
    }
    
    Then func() {
        let user = findUser("nonexistent@email.com")
        assertTrue(user == null)
        return "User not found as expected"
    }
}

func main() {
    print("=== USER TEST SUITE ===")
}
```

### 4. Testing Classes and Objects

```r2
class BankCalculator {
    let balance
    let history
    
    constructor(initialBalance) {
        this.balance = initialBalance
        this.history = []
    }
    
    deposit(amount) {
        if (amount <= 0) {
            throw "Amount must be positive"
        }
        
        this.balance = this.balance + amount
        this.history = this.history.push({
            type: "Deposit",
            amount: amount,
            resultingBalance: this.balance
        })
        
        return this.balance
    }
    
    withdraw(amount) {
        if (amount <= 0) {
            throw "Amount must be positive"
        }
        
        if (amount > this.balance) {
            throw "Insufficient balance"
        }
        
        this.balance = this.balance - amount
        this.history = this.history.push({
            type: "Withdrawal",
            amount: amount,
            resultingBalance: this.balance
        })
        
        return this.balance
    }
    
    getBalance() {
        return this.balance
    }
    
    getHistory() {
        return this.history
    }
}

// Global variable for tests
let calculator

TestCase "Basic bank account operations" {
    Given func() {
        calculator = BankCalculator(1000)
        print("Calculator initialized with balance:", calculator.getBalance())
        return "Calculator ready"
    }
    
    When func() {
        calculator.deposit(500)
        print("Deposit of 500 made")
        return "Deposit completed"
    }
    
    Then func() {
        assertEqual(calculator.getBalance(), 1500)
        let history = calculator.getHistory()
        assertEqual(history.length(), 1)
        assertEqual(history[0].type, "Deposit")
        return "Deposit validated correctly"
    }
    
    And func() {
        calculator.withdraw(300)
        assertEqual(calculator.getBalance(), 1200)
        let history = calculator.getHistory()
        assertEqual(history.length(), 2)
        return "Withdrawal validated correctly"
    }
}

TestCase "Error handling in banking operations" {
    Given func() {
        calculator = BankCalculator(100)
        return "Calculator with low balance initialized"
    }
    
    When func() {
        print("Attempting to withdraw more money than available")
        return "Excessive withdrawal attempt"
    }
    
    Then func() {
        try {
            calculator.withdraw(200)
            assertTrue(false)  // Should not get here
        } catch (error) {
            assertTrue(error.contains("insufficient"))
            assertEqual(calculator.getBalance(), 100)  // Balance did not change
        }
        return "Error handled correctly"
    }
    
    And func() {
        try {
            calculator.deposit(-50)
            assertTrue(false)  // Should not get here
        } catch (error) {
            assertTrue(error.contains("positive"))
        }
        return "Negative amount validation correct"
    }
}

func main() {
    print("=== BANK CALCULATOR TESTS ===")
}
```

### 5. Testing Concurrent Functions

```r2
let globalCounter = 0
let concurrentResults = []

func incrementCounter(times, id) {
    for (let i = 0; i < times; i++) {
        globalCounter++
        concurrentResults = concurrentResults.push({
            worker: id,
            value: globalCounter
        })
        sleep(0.1)  // Simulate work
    }
}

func resetCounters() {
    globalCounter = 0
    concurrentResults = []
}

TestCase "Verify concurrent behavior" {
    Given func() {
        resetCounters()
        print("Counters reset")
        return "Clean initial state"
    }
    
    When func() {
        r2(incrementCounter, 3, "Worker-1")
        r2(incrementCounter, 3, "Worker-2")
        r2(incrementCounter, 3, "Worker-3")
        
        sleep(2)  // Wait for them to finish
        print("Concurrent operations completed")
        return "Concurrent increments executed"
    }
    
    Then func() {
        assertEqual(globalCounter, 9)
        assertEqual(concurrentResults.length(), 9)
        print("Final values validated")
        return "Concurrency verified"
    }
    
    And func() {
        // Verify that there was interleaving (non-deterministic, but likely)
        let workers = []
        for (let i = 0; i < concurrentResults.length(); i++) {
            let result = concurrentResults[i]
            workers = workers.push(result.worker)
        }
        
        print("Sequence of workers:", workers)
        // In a real run, we should see interleaved workers
        return "Interleaving pattern observed"
    }
}

func main() {
    print("=== CONCURRENCY TESTS ===")
}
```

## Web Development and APIs

### 1. Basic HTTP Server

```r2
func handleRoot(req, res) {
    res.send("Hello from R2Lang!")
}

func handleGreeting(req, res) {
    let name = req.query.name || "Anonymous"
    res.send("Hello " + name + "!")
}

func handleInfo(req, res) {
    let info = {
        server: "R2Lang HTTP Server",
        version: "1.0",
        timestamp: "2024-01-01",
        endpoints: ["/", "/greeting", "/info"]
    }
    res.json(info)
}

func main() {
    print("Starting web server on port 8080...")
    
    // Configure routes
    http.get("/", handleRoot)
    http.get("/greeting", handleGreeting)
    http.get("/info", handleInfo)
    
    // Start server
    http.listen(8080)
}
```

### 2. Complete REST API

```r2
// In-memory database simulation
let users = []
let nextId = 1

// Utility functions
func generateId() {
    let id = nextId
    nextId++
    return id
}

func findUser(id) {
    for (let i = 0; i < users.length(); i++) {
        let user = users[i]
        if (user.id == id) {
            return { user: user, index: i }
        }
    }
    return null
}

func validateUser(data) {
    if (!data.name || data.name == "") {
        return "Name is required"
    }
    
    if (!data.email || data.email == "") {
        return "Email is required"
    }
    
    // Verify unique email
    for (let i = 0; i < users.length(); i++) {
        if (users[i].email == data.email) {
            return "Email is already in use"
        }
    }
    
    return null
}

// API Handlers
func getAllUsers(req, res) {
    res.json({
        users: users,
        total: users.length()
    })
}

func getUser(req, res) {
    let id = parseInt(req.params.id)
    let result = findUser(id)
    
    if (result == null) {
        res.status(404).json({
            error: "User not found",
            id: id
        })
        return
    }
    
    res.json(result.user)
}

func createUser(req, res) {
    let data = req.body
    
    // Validate data
    let error = validateUser(data)
    if (error != null) {
        res.status(400).json({
            error: error
        })
        return
    }
    
    // Create user
    let newUser = {
        id: generateId(),
        name: data.name,
        email: data.email,
        active: true,
        creationDate: "2024-01-01"
    }
    
    users = users.push(newUser)
    
    res.status(201).json(newUser)
}

func updateUser(req, res) {
    let id = parseInt(req.params.id)
    let data = req.body
    let result = findUser(id)
    
    if (result == null) {
        res.status(404).json({
            error: "User not found"
        })
        return
    }
    
    // Update fields
    let user = result.user
    if (data.name) {
        user.name = data.name
    }
    if (data.email) {
        user.email = data.email
    }
    if (data.active != null) {
        user.active = data.active
    }
    
    res.json(user)
}

func deleteUser(req, res) {
    let id = parseInt(req.params.id)
    let result = findUser(id)
    
    if (result == null) {
        res.status(404).json({
            error: "User not found"
        })
        return
    }
    
    // Delete user (simulate by removing from the array)
    let deletedUser = result.user
    
    // In the current R2Lang, we don't have a remove method, so we'll simulate
    let newUsers = []
    for (let i = 0; i < users.length(); i++) {
        if (users[i].id != id) {
            newUsers = newUsers.push(users[i])
        }
    }
    users = newUsers
    
    res.json({
        message: "User deleted",
        user: deletedUser
    })
}

// Logging middleware
func middleware(req, res, next) {
    print("📥", req.method, req.url, "- IP:", req.ip)
    next()
}

func main() {
    print("🚀 Starting REST API on port 3000...")
    
    // Apply middleware
    http.use(middleware)
    
    // Configure REST routes
    http.get("/api/users", getAllUsers)
    http.get(" /api/users/:id", getUser)
    http.post("/api/users", createUser)
    http.put("/api/users/:id", updateUser)
    http.delete("/api/users/:id", deleteUser)
    
    // Health route
    http.get("/health", func(req, res) {
        res.json({
            status: "OK",
            timestamp: "2024-01-01",
            uptime: "5 minutes"
        })
    })
    
    // Start server
    http.listen(3000)
    print("✅ REST API available at http://localhost:3000")
    print("📋 Available endpoints:")
    print("  GET    /api/users")
    print("  GET    /api/users/:id")
    print("  POST   /api/users")
    print("  PUT    /api/users/:id")
    print("  DELETE /api/users/:id")
    print("  GET    /health")
}
```

### 3. API Testing

```r2
// HTTP client simulation for testing
let baseURL = "http://localhost:3000"

func makeRequest(method, url, body) {
    // HTTP request simulation
    print("📤", method, baseURL + url)
    if (body) {
        print("📦 Body:", body)
    }
    
    // Simulate response based on the endpoint
    if (method == "GET" && url == "/api/users") {
        return {
            status: 200,
            body: { users: [], total: 0 }
        }
    }
    
    if (method == "POST" && url == "/api/users") {
        return {
            status: 201,
            body: { 
                id: 1, 
                name: body.name, 
                email: body.email,
                active: true 
            }
        }
    }
    
    return {
        status: 200,
        body: { message: "Simulated response" }
    }
}

TestCase "REST API - Create user" {
    Given func() {
        print("Preparing data to create a user")
        return "HTTP client ready"
    }
    
    When func() {
        let newUser = {
            name: "Juan Pérez",
            email: "juan@email.com"
        }
        
        let response = makeRequest("POST", "/api/users", newUser)
        print("User created with response:", response.status)
        return "User created via API"
    }
    
    Then func() {
        let newUser = {
            name: "Juan Pérez",
            email: "juan@email.com"
        }
        
        let response = makeRequest("POST", "/api/users", newUser)
        
        assertEqual(response.status, 201)
        assertTrue(response.body.id != null)
        assertEqual(response.body.name, "Juan Pérez")
        assertEqual(response.body.email, "juan@email.com")
        assertTrue(response.body.active)
        
        return "Creation response validated"
    }
}

TestCase "REST API - Get empty list" {
    Given func() {
        print("Clean API with no users")
        return "Initial state"
    }
    
    When func() {
        let response = makeRequest("GET", "/api/users", null)
        return "User list obtained"
    }
    
    Then func() {
        let response = makeRequest("GET", "/api/users", null)
        
        assertEqual(response.status, 200)
        assertTrue(response.body.users != null)
        assertEqual(response.body.total, 0)
        
        return "Empty list validated"
    }
}

func main() {
    print("=== REST API TESTS ===")
}
```

## Integration of Testing and Web

### 1. End-to-End Testing

```r2
// Complete web application simulation
class WebApplication {
    let users
    let sessions
    
    constructor() {
        this.users = []
        this.sessions = []
    }
    
    registerUser(data) {
        // Validate data
        if (!data.name || !data.email || !data.password) {
            throw "Incomplete data"
        }
        
        // Verify unique email
        for (let i = 0; i < this.users.length(); i++) {
            if (this.users[i].email == data.email) {
                throw "Email already registered"
            }
        }
        
        let user = {
            id: this.users.length() + 1,
            name: data.name,
            email: data.email,
            password: data.password,  // In prod, this should be hashed
            active: true
        }
        
        this.users = this.users.push(user)
        return user
    }
    
    login(email, password) {
        for (let i = 0; i < this.users.length(); i++) {
            let user = this.users[i]
            if (user.email == email && user.password == password) {
                let session = {
                    token: "token_" + user.id + "_" + this.sessions.length(),
                    userId: user.id,
                    active: true
                }
                
                this.sessions = this.sessions.push(session)
                return session
            }
        }
        
        throw "Invalid credentials"
    }
    
    logout(token) {
        for (let i = 0; i < this.sessions.length(); i++) {
            let session = this.sessions[i]
            if (session.token == token) {
                session.active = false
                return true
            }
        }
        return false
    }
    
    getCurrentUser(token) {
        for (let i = 0; i < this.sessions.length(); i++) {
            let session = this.sessions[i]
            if (session.token == token && session.active) {
                for (let j = 0; j < this.users.length(); j++) {
                    let user = this.users[j]
                    if (user.id == session.userId) {
                        return user
                    }
                }
            }
        }
        return null
    }
}

let app
let testUser
let testSession

TestCase "Complete user flow - Registration and Login" {
    Given func() {
        app = WebApplication()
        testUser = {
            name: "María González",
            email: "maria@test.com",
            password: "password123"
        }
        print("Web application initialized")
        return "App ready for testing"
    }
    
    When func() {
        let user = app.registerUser(testUser)
        print("User registered:", user.name)
        return "Registration completed"
    }
    
    Then func() {
        let user = app.registerUser(testUser)
        assertTrue(user.id != null)
        assertEqual(user.name, "María González")
        assertEqual(user.email, "maria@test.com")
        assertTrue(user.active)
        return "Registration validated"
    }
    
    And func() {
        let session = app.login(testUser.email, testUser.password)
        testSession = session
        
        assertTrue(session.token != null)
        assertTrue(session.active)
        assertEqual(session.userId, 1)
        
        return "Successful login validated"
    }
}

TestCase "Session management" {
    Given func() {
        print("Using active session from the previous test")
        return "Session available"
    }
    
    When func() {
        let user = app.getCurrentUser(testSession.token)
        print("Current user obtained:", user.name)
        return "Session user obtained"
    }
    
    Then func() {
        let user = app.getCurrentUser(testSession.token)
        assertTrue(user != null)
        assertEqual(user.email, "maria@test.com")
        return "Session user validated"
    }
    
    And func() {
        let result = app.logout(testSession.token)
        assertTrue(result)
        
        let userAfterLogout = app.getCurrentUser(testSession.token)
        assertTrue(userAfterLogout == null)
        
        return "Logout validated"
    }
}

TestCase "Authentication error handling" {
    Given func() {
        print("Preparing error cases")
        return "Error cases ready"
    }
    
    When func() {
        print("Attempting to register a duplicate user")
        return "Duplicate registration attempt"
    }
    
    Then func() {
        try {
            app.registerUser(testUser)  // Same email
            assertTrue(false)  // Should not get here
        } catch (error) {
            assertTrue(error.contains("already registered"))
        }
        return "Duplicate email error handled"
    }
    
    And func() {
        try {
            app.login("nonexistent@test.com", "wrongpass")
            assertTrue(false)  // Should not get here
        } catch (error) {
            assertTrue(error.contains("invalid"))
        }
        return "Invalid credentials error handled"
    }
}

func main() {
    print("=== END-TO-END WEB APPLICATION TESTS ===")
}
```

## Module Project: Blog System with Testing

```r2
// Complete blog system with BDD testing

class Post {
    let id
    let title
    let content
    let author
    let creationDate
    let modificationDate
    let active
    let comments
    
    constructor(id, title, content, author) {
        this.id = id
        this.title = title
        this.content = content
        this.author = author
        this.creationDate = "2024-01-01"
        this.modificationDate = "2024-01-01"
        this.active = true
        this.comments = []
    }
    
    addComment(author, content) {
        let comment = {
            id: this.comments.length() + 1,
            author: author,
            content: content,
            date: "2024-01-01"
        }
        
        this.comments = this.comments.push(comment)
        return comment
    }
    
    update(newTitle, newContent) {
        this.title = newTitle
        this.content = newContent
        this.modificationDate = "2024-01-01"
    }
    
    deactivate() {
        this.active = false
    }
}

class BlogService {
    let posts
    let nextId
    
    constructor() {
        this.posts = []
        this.nextId = 1
    }
    
    createPost(title, content, author) {
        if (!title || !content || !author) {
            throw "Title, content, and author are required"
        }
        
        let post = Post(this.nextId, title, content, author)
        this.nextId++
        this.posts = this.posts.push(post)
        
        return post
    }
    
    getPost(id) {
        for (let i = 0; i < this.posts.length(); i++) {
            let post = this.posts[i]
            if (post.id == id && post.active) {
                return post
            }
        }
        return null
    }
    
    getPostsByAuthor(author) {
        let posts = []
        for (let i = 0; i < this.posts.length(); i++) {
            let post = this.posts[i]
            if (post.author == author && post.active) {
                posts = posts.push(post)
            }
        }
        return posts
    }
    
    searchPosts(term) {
        let posts = []
        for (let i = 0; i < this.posts.length(); i++) {
            let post = this.posts[i]
            if (post.active && 
                (post.title.contains(term) || post.content.contains(term))) {
                posts = posts.push(post)
            }
        }
        return posts
    }
    
    deletePost(id) {
        let post = this.getPost(id)
        if (post != null) {
            post.deactivate()
            return true
        }
        return false
    }
    
    getStatistics() {
        let totalPosts = 0
        let totalComments = 0
        let authors = []
        
        for (let i = 0; i < this.posts.length(); i++) {
            let post = this.posts[i]
            if (post.active) {
                totalPosts++
                totalComments = totalComments + post.comments.length()
                
                // Count unique authors (simplified)
                let authorExists = false
                for (let j = 0; j < authors.length(); j++) {
                    if (authors[j] == post.author) {
                        authorExists = true
                        break
                    }
                }
                if (!authorExists) {
                    authors = authors.push(post.author)
                }
            }
        }
        
        return {
            totalPosts: totalPosts,
            totalComments: totalComments,
            totalAuthors: authors.length()
        }
    }
}

// Global variables for testing
let blogService
let testPost

TestCase "Post creation and management" {
    Given func() {
        blogService = BlogService()
        print("Blog service initialized")
        return "Blog system ready"
    }
    
    When func() {
        testPost = blogService.createPost(
            "My first post",
            "This is the content of my first post in R2Lang",
            "Juan Blogger"
        )
        print("Post created with ID:", testPost.id)
        return "Post created successfully"
    }
    
    Then func() {
        assertTrue(testPost != null)
        assertEqual(testPost.id, 1)
        assertEqual(testPost.title, "My first post")
        assertEqual(testPost.author, "Juan Blogger")
        assertTrue(testPost.active)
        assertEqual(testPost.comments.length(), 0)
        return "Post validated correctly"
    }
    
    And func() {
        let retrievedPost = blogService.getPost(1)
        assertTrue(retrievedPost != null)
        assertEqual(retrievedPost.title, testPost.title)
        return "Post retrieved correctly"
    }
}

TestCase "Comment system" {
    Given func() {
        print("Using existing post for comments")
        return "Post available for comments"
    }
    
    When func() {
        let comment = testPost.addComment(
            "Ana Lectora",
            "Excellent post! I liked it a lot."
        )
        print("Comment added by:", comment.author)
        return "Comment added"
    }
    
    Then func() {
        assertEqual(testPost.comments.length(), 1)
        let comment = testPost.comments[0]
        assertEqual(comment.author, "Ana Lectora")
        assertTrue(comment.content.contains("Excellent"))
        return "Comment validated"
    }
    
    And func() {
        testPost.addComment("Carlos Lector", "Very informative")
        assertEqual(testPost.comments.length(), 2)
        return "Multiple comments working"
    }
}

TestCase "Search and filters" {
    Given func() {
        // Create additional posts for searching
        blogService.createPost(
            "R2Lang Tutorial",
            "Learn R2Lang from scratch",
            "María Tutora"
        )
        blogService.createPost(
            "Advanced Programming",
            "Advanced programming techniques",
            "Juan Blogger"
        )
        print("Additional posts created for searching")
        return "Post dataset prepared"
    }
    
    When func() {
        let r2langPosts = blogService.searchPosts("R2Lang")
        print("Search for 'R2Lang' found:", r2langPosts.length(), "posts")
        return "Search executed"
    }
    
    Then func() {
        let r2langPosts = blogService.searchPosts("R2Lang")
        assertEqual(r2langPosts.length(), 2)  // "My first post" and "R2Lang Tutorial"
        return "Search by term validated"
    }
    
    And func() {
        let juanPosts = blogService.getPostsByAuthor("Juan Blogger")
        assertEqual(juanPosts.length(), 2)  // "My first post" and "Advanced Programming"
        
        let mariaPosts = blogService.getPostsByAuthor("María Tutora")
        assertEqual(mariaPosts.length(), 1)  // "R2Lang Tutorial"
        
        return "Filter by author validated"
    }
}

TestCase "Blog statistics" {
    Given func() {
        print("Calculating blog statistics")
        return "Blog with multiple posts and comments"
    }
    
    When func() {
        let stats = blogService.getStatistics()
        print("Calculated statistics:", stats)
        return "Statistics obtained"
    }
    
    Then func() {
        let stats = blogService.getStatistics()
        assertEqual(stats.totalPosts, 3)
        assertEqual(stats.totalComments, 2)  // Only the first post has comments
        assertEqual(stats.totalAuthors, 2)     // Juan Blogger and María Tutora
        return "Statistics validated"
    }
}

TestCase "Post deletion" {
    Given func() {
        print("Preparing post deletion")
        return "Posts available for deletion"
    }
    
    When func() {
        let result = blogService.deletePost(1)
        print("Post 1 deleted:", result)
        return "Deletion executed"
    }
    
    Then func() {
        let result = blogService.deletePost(1)
        assertTrue(result)
        
        let deletedPost = blogService.getPost(1)
        assertTrue(deletedPost == null)  // Should not find it
        
        return "Deletion validated"
    }
    
    And func() {
        let updatedStats = blogService.getStatistics()
        assertEqual(updatedStats.totalPosts, 2)  // One less post
        return "Statistics updated after deletion"
    }
}

func main() {
    print("=== COMPLETE BLOG TEST SUITE ===")
    print("Running BDD tests...")
}
```

## Best Practices

### 1. BDD Testing
- ✅ Use descriptive names in TestCase
- ✅ Keep Given-When-Then focused
- ✅ One concept per step (Given/When/Then)
- ✅ Use realistic test data

### 2. REST APIs
- ✅ Follow RESTful conventions
- ✅ Handle errors with appropriate HTTP codes
- ✅ Validate data input
- ✅ Document endpoints clearly

### 3. API Testing
- ✅ Test happy paths and error cases
- ✅ Verify HTTP status codes
- ✅ Validate response structures
- ✅ Test authentication and authorization

## Module Summary

### Concepts Learned
- ✅ Integrated BDD testing system
- ✅ TestCase with Given-When-Then-And
- ✅ Testing classes and objects
- ✅ Testing concurrent operations
- ✅ REST API development
- ✅ End-to-end testing
- ✅ Testing-development integration

### Skills Developed
- ✅ Writing expressive tests with BDD
- ✅ Creating complete REST APIs
- ✅ Designing effective test suites
- ✅ Validating application behavior
- ✅ Integrating testing into development
- ✅ Documenting behavior with tests

### Next Module

In **Module 6**, you will learn:
- Working with files and databases
- Application deployment and distribution
- Optimization and performance
- Advanced architectural patterns

Congratulations! You now master BDD testing and web development in R2Lang, two distinctive features of the language.
