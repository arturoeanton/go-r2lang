# R2Lang Course - Module 4: Concurrency and Error Handling

## Introduction

In this module, you will learn two advanced aspects of R2Lang: concurrent programming with goroutines and robust error handling. These concepts are fundamental for creating robust and efficient applications.

## Concurrency in R2Lang

### 1. Basic Concurrency Concepts

Concurrency allows multiple tasks to run "at the same time" (in parallel or interleaved). R2Lang uses goroutines, similar to those in Go, to handle concurrency.

#### Your First Goroutine

```r2
func task() {
    print("Executing task in a goroutine")
    sleep(1)  // Simulate time-consuming work
    print("Task completed")
}

func main() {
    print("Starting program")
    
    // Execute function in a goroutine
    r2(task)
    
    print("Continuing with other operations")
    sleep(2)  // Wait for the goroutine to finish
    print("Program finished")
}
```

#### Multiple Goroutines

```r2
func worker(id) {
    print("Worker", id, "started")
    
    for (let i = 1; i <= 3; i++) {
        print("Worker", id, "- task", i)
        sleep(1)
    }
    
    print("Worker", id, "finished")
}

func main() {
    print("Creating workers...")
    
    // Create multiple goroutines
    for (let i = 1; i <= 3; i++) {
        r2(worker, i)
    }
    
    print("All workers created")
    sleep(4)  // Wait for them to finish
    print("Main program finished")
}
```

### 2. Concurrency Patterns

#### Worker Pool Pattern

```r2
func processData(data, workerId) {
    print("Worker", workerId, "processing:", data)
    sleep(1)  // Simulate processing
    print("Worker", workerId, "completed:", data)
}

func createWorkerPool(numWorkers, tasks) {
    print("Creating a pool of", numWorkers, "workers")
    
    for (let i = 0; i < numWorkers; i++) {
        let workerId = i + 1
        
        r2(func() {
            // Each worker processes a portion of the tasks
            let tasksPerWorker = tasks.length() / numWorkers
            let start = i * tasksPerWorker
            let end = start + tasksPerWorker
            
            for (let j = start; j < end && j < tasks.length(); j++) {
                processData(tasks[j], workerId)
            }
        })
    }
}

func main() {
    let tasks = ["Task-A", "Task-B", "Task-C", "Task-D", "Task-E", "Task-F"]
    
    print("Starting parallel processing")
    createWorkerPool(3, tasks)
    
    sleep(4)  // Wait for all to finish
    print("Processing completed")
}
```

#### Producer-Consumer Pattern

```r2
func producer(producerName, quantity) {
    for (let i = 1; i <= quantity; i++) {
        let product = producerName + "-Item-" + i
        print("📦 Produced:", product)
        sleep(0.5)  // Simulate production time
    }
    print("✅ Producer", producerName, "finished")
}

func consumer(consumerName, totalTime) {
    let startTime = 0  // Time simulation
    let timeLimit = totalTime * 2  // 2 seconds per unit
    
    while (startTime < timeLimit) {
        print("🛒 Consumer", consumerName, "processing items...")
        sleep(1)
        startTime++
    }
    print("✅ Consumer", consumerName, "finished")
}

func main() {
    print("=== PRODUCER-CONSUMER PATTERN ===")
    
    // Start producers
    r2(producer, "P1", 3)
    r2(producer, "P2", 4)
    
    // Start consumers  
    r2(consumer, "C1", 3)
    r2(consumer, "C2", 4)
    
    sleep(6)
    print("Simulation finished")
}
```

### 3. Concurrency with Classes

```r2
class ConcurrentCounter {
    let value
    let name
    
    constructor(name) {
        this.name = name
        this.value = 0
    }
    
    increment(amount) {
        for (let i = 0; i < amount; i++) {
            this.value++
            print(this.name, "incremented to:", this.value)
            sleep(0.1)  // Simulate work
        }
    }
    
    decrement(amount) {
        for (let i = 0; i < amount; i++) {
            this.value--
            print(this.name, "decremented to:", this.value)
            sleep(0.1)  // Simulate work
        }
    }
    
    getValue() {
        return this.value
    }
}

func main() {
    let counter = ConcurrentCounter("Counter-1")
    
    print("Initial value:", counter.getValue())
    
    // Concurrent operations
    r2(func() {
        counter.increment(5)
    })
    
    r2(func() {
        counter.decrement(3)
    })
    
    r2(func() {
        counter.increment(2)
    })
    
    sleep(3)
    print("Final value:", counter.getValue())
}
```

### 4. Simulating Synchronization

Although R2Lang does not have native synchronization primitives, we can simulate some patterns:

```r2
// Simulating a Mutex using flags
class SimulatedMutex {
    let locked
    let name
    
    constructor(name) {
        this.locked = false
        this.name = name
    }
    
    lock() {
        while (this.locked) {
            sleep(0.01)  // Wait
        }
        this.locked = true
        print("🔒 Lock acquired by", this.name)
    }
    
    unlock() {
        this.locked = false
        print("🔓 Lock released by", this.name)
    }
}

func workWithMutex(id, mutex, sharedResource) {
    print("Process", id, "trying to access the resource")
    
    mutex.lock()
    
    try {
        print("Process", id, "using shared resource")
        sharedResource.value++
        print("Resource updated to:", sharedResource.value)
        sleep(1)  // Simulate resource usage
    } finally {
        mutex.unlock()
        print("Process", id, "finished using the resource")
    }
}

func main() {
    let mutex = SimulatedMutex("Main-Mutex")
    let resource = { value: 0 }
    
    print("Starting concurrent access to shared resource")
    
    for (let i = 1; i <= 3; i++) {
        r2(workWithMutex, i, mutex, resource)
    }
    
    sleep(5)
    print("Final value of the resource:", resource.value)
}
```

## Error Handling

### 1. Basic Try-Catch-Finally

```r2
func riskyOperation(number) {
    if (number < 0) {
        throw "Number cannot be negative"
    }
    
    if (number == 0) {
        throw "Division by zero not allowed"
    }
    
    return 100 / number
}

func main() {
    let numbers = [10, -5, 0, 20]
    
    for (let number in numbers) {
        print("Processing number:", number)
        
        try {
            let result = riskyOperation(number)
            print("Result:", result)
        } catch (error) {
            print("Error caught:", error)
        } finally {
            print("Operation completed for", number)
        }
        print("---")
    }
}
```

### 2. Error Handling in Functions

```r2
func validateAge(age) {
    if (typeOf(age) != "float64") {
        throw "Age must be a number"
    }
    
    if (age < 0) {
        throw "Age cannot be negative"
    }
    
    if (age > 150) {
        throw "Age cannot be greater than 150"
    }
    
    return true
}

func createPerson(name, age) {
    try {
        // Validate input
        if (name == null || name == "") {
            throw "Name is required"
        }
        
        validateAge(age)
        
        // Create person if everything is fine
        let person = {
            name: name,
            age: age,
            createdAt: "Now"
        }
        
        print("Person created:", person.name)
        return person
        
    } catch (error) {
        print("Error creating person:", error)
        return null
    }
}

func main() {
    let peopleData = [
        ["Juan", 25],
        ["", 30],      // Error: empty name
        ["Ana", -5],   // Error: negative age
        ["Carlos", "thirty"],  // Error: non-numeric age
        ["María", 200], // Error: age too high
        ["Luis", 35]
    ]
    
    let validPeople = []
    
    for (let data in peopleData) {
        let name = data[0]
        let age = data[1]
        
        let person = createPerson(name, age)
        if (person != null) {
            validPeople = validPeople.push(person)
        }
    }
    
    print("\nValid people created:", validPeople.length())
    for (let person in validPeople) {
        print("-", person.name, "(" + person.age + " years)")
    }
}
```

### 3. Errors in File Operations

```r2
func readFile(fileName) {
    try {
        print("Attempting to read file:", fileName)
        
        // Simulate file reading
        if (fileName == "nonexistent.txt") {
            throw "File not found"
        }
        
        if (fileName == "corrupt.txt") {
            throw "File is corrupt or damaged"
        }
        
        if (fileName == "permissions.txt") {
            throw "No permissions to read the file"
        }
        
        // Simulate file content
        let content = "Content of the file " + fileName
        print("File read successfully")
        return content
        
    } catch (error) {
        print("Error reading file:", error)
        throw "File error: " + error
    }
}

func processFiles(files) {
    let processed = 0
    let errors = 0
    
    for (let file in files) {
        try {
            let content = readFile(file)
            print("Processing content of", file)
            processed++
            
        } catch (error) {
            print("Could not process", file + ":", error)
            errors++
            
        } finally {
            print("Finalizing processing of", file)
        }
        print("---")
    }
    
    print("SUMMARY:")
    print("Files processed:", processed)
    print("Errors found:", errors)
}

func main() {
    let files = [
        "document1.txt",
        "nonexistent.txt",
        "data.txt",
        "corrupt.txt",
        "permissions.txt",
        "final.txt"
    ]
    
    processFiles(files)
}
```

### 4. Error Handling in Concurrency

```r2
func taskWithErrors(id, shouldFail) {
    try {
        print("Task", id, "started")
        
        if (shouldFail) {
            throw "Simulated error in task " + id
        }
        
        // Simulate work
        for (let i = 1; i <= 3; i++) {
            print("Task", id, "- step", i)
            sleep(0.5)
        }
        
        print("Task", id, "completed successfully")
        
    } catch (error) {
        print("ERROR in task", id + ":", error)
        
    } finally {
        print("Task", id, "finalizing resources")
    }
}

func taskSupervisor(numTasks) {
    let successfulTasks = 0
    let failedTasks = 0
    
    print("Supervisor starting", numTasks, "tasks")
    
    for (let i = 1; i <= numTasks; i++) {
        // Some tasks will fail (simulated)
        let shouldFail = (i % 3 == 0)  // Every third task fails
        
        r2(func() {
            try {
                taskWithErrors(i, shouldFail)
                // We can't update counters directly due to concurrency
                print("✅ Task", i, "registered as successful")
            } catch (error) {
                print("❌ Task", i, "registered as failed")
            }
        })
    }
    
    print("All tasks launched")
}

func main() {
    taskSupervisor(6)
    sleep(4)
    print("Supervision completed")
}
```

## Advanced Patterns

### 1. Circuit Breaker Pattern

```r2
class CircuitBreaker {
    let name
    let errorThreshold
    let consecutiveErrors
    let state  // "CLOSED", "OPEN", "HALF_OPEN"
    let lastErrorTime
    
    constructor(name, errorThreshold) {
        this.name = name
        this.errorThreshold = errorThreshold
        this.consecutiveErrors = 0
        this.state = "CLOSED"
        this.lastErrorTime = 0
    }
    
    execute(operation) {
        if (this.state == "OPEN") {
            print("Circuit Breaker OPEN - operation blocked")
            throw "Circuit breaker is open"
        }
        
        try {
            let result = operation()
            this.onSuccess()
            return result
            
        } catch (error) {
            this.onError()
            throw error
        }
    }
    
    onSuccess() {
        this.consecutiveErrors = 0
        if (this.state == "HALF_OPEN") {
            this.state = "CLOSED"
            print("Circuit Breaker returns to CLOSED")
        }
    }
    
    onError() {
        this.consecutiveErrors++
        this.lastErrorTime = 1  // Timestamp simulation
        
        if (this.consecutiveErrors >= this.errorThreshold) {
            this.state = "OPEN"
            print("Circuit Breaker OPEN after", this.consecutiveErrors, "errors")
        }
    }
    
    attemptRecovery() {
        if (this.state == "OPEN") {
            this.state = "HALF_OPEN"
            print("Circuit Breaker in HALF_OPEN mode")
        }
    }
}

func externalOperation(success) {
    if (success) {
        return "Operation successful"
    } else {
        throw "Operation failed"
    }
}

func main() {
    let cb = CircuitBreaker("API-Circuit-Breaker", 3)
    
    // Simulate calls that fail
    let attempts = [false, false, false, false, true]
    
    for (let i = 0; i < attempts.length(); i++) {
        let success = attempts[i]
        
        try {
            let result = cb.execute(func() {
                return externalOperation(success)
            })
            print("Result:", result)
            
        } catch (error) {
            print("Error:", error)
        }
        
        print("Current state:", cb.state)
        print("---")
    }
    
    // Attempt recovery
    print("Attempting recovery...")
    cb.attemptRecovery()
    
    try {
        let result = cb.execute(func() {
            return externalOperation(true)
        })
        print("Result after recovery:", result)
    } catch (error) {
        print("Error in recovery:", error)
    }
}
```

### 2. Retry Pattern

```r2
func retryOperation(operation, maxAttempts, delaySeconds) {
    let attempts = 0
    
    while (attempts < maxAttempts) {
        attempts++
        
        try {
            print("Attempt", attempts, "of", maxAttempts)
            let result = operation()
            print("Operation successful on attempt", attempts)
            return result
            
        } catch (error) {
            print("Error on attempt", attempts + ":", error)
            
            if (attempts >= maxAttempts) {
                print("Maximum number of attempts reached")
                throw "Operation failed after " + maxAttempts + " attempts"
            }
            
            print("Waiting", delaySeconds, "seconds before the next attempt")
            sleep(delaySeconds)
        }
    }
}

func unstableOperation() {
    // Simulate an operation that fails randomly
    let random = math.random()  // Assuming it exists
    
    if (random > 0.7) {  // 30% chance of success
        return "Operation completed successfully"
    } else {
        throw "Temporary network failure"
    }
}

func main() {
    try {
        let result = retryOperation(unstableOperation, 5, 1)
        print("Final result:", result)
        
    } catch (error) {
        print("Operation finally failed:", error)
    }
}
```

## Module Project: Distributed Processing System

```r2
// Simulation of a distributed processing system
// with error handling and concurrency

class Node {
    let id
    let isActive
    let load
    let consecutiveErrors
    
    constructor(id) {
        this.id = id
        this.isActive = true
        this.load = 0
        this.consecutiveErrors = 0
    }
    
    process(task) {
        if (!this.isActive) {
            throw "Node " + this.id + " is inactive"
        }
        
        if (this.load >= 5) {
            throw "Node " + this.id + " is overloaded"
        }
        
        this.load++
        
        try {
            print("Node", this.id, "processing task:", task.name)
            
            // Simulate processing
            sleep(task.duration)
            
            // Simulate possible error
            if (task.name.contains("error")) {
                throw "Error in task: " + task.name
            }
            
            let result = {
                node: this.id,
                task: task.name,
                result: "Processed successfully",
                time: task.duration
            }
            
            this.consecutiveErrors = 0
            print("✅ Node", this.id, "completed:", task.name)
            return result
            
        } catch (error) {
            this.consecutiveErrors++
            print("❌ Error in node", this.id + ":", error)
            
            if (this.consecutiveErrors >= 3) {
                this.isActive = false
                print("🚫 Node", this.id, "deactivated due to consecutive errors")
            }
            
            throw error
            
        } finally {
            this.load--
        }
    }
    
    restart() {
        this.isActive = true
        this.load = 0
        this.consecutiveErrors = 0
        print("🔄 Node", this.id, "restarted")
    }
}

class Coordinator {
    let nodes
    let pendingTasks
    let completedTasks
    let failedTasks
    
    constructor() {
        this.nodes = []
        this.pendingTasks = []
        this.completedTasks = []
        this.failedTasks = []
    }
    
    addNode(node) {
        this.nodes = this.nodes.push(node)
        print("Node", node.id, "added to the cluster")
    }
    
    addTask(task) {
        this.pendingTasks = this.pendingTasks.push(task)
    }
    
    findAvailableNode() {
        for (let i = 0; i < this.nodes.length(); i++) {
            let node = this.nodes[i]
            if (node.isActive && node.load < 5) {
                return node
            }
        }
        return null
    }
    
    processTasks() {
        print("Starting distributed processing")
        
        for (let i = 0; i < this.pendingTasks.length(); i++) {
            let task = this.pendingTasks[i]
            
            r2(func() {
                let processed = false
                let attempts = 0
                let maxAttempts = 3
                
                while (!processed && attempts < maxAttempts) {
                    attempts++
                    
                    try {
                        let node = this.findAvailableNode()
                        
                        if (node == null) {
                            print("⏳ No available nodes, waiting...")
                            sleep(1)
                            continue
                        }
                        
                        let result = node.process(task)
                        this.completedTasks = this.completedTasks.push(result)
                        processed = true
                        
                    } catch (error) {
                        print("Error processing", task.name + ":", error)
                        
                        if (attempts >= maxAttempts) {
                            this.failedTasks = this.failedTasks.push({
                                task: task,
                                error: error,
                                attempts: attempts
                            })
                            processed = true  // To exit the loop
                        } else {
                            sleep(1)  // Wait before retrying
                        }
                    }
                }
            })
        }
    }
    
    showStatistics() {
        print("\n=== CLUSTER STATISTICS ===")
        print("Total nodes:", this.nodes.length())
        
        let activeNodes = 0
        for (let i = 0; i < this.nodes.length(); i++) {
            if (this.nodes[i].isActive) {
                activeNodes++
            }
        }
        
        print("Active nodes:", activeNodes)
        print("Completed tasks:", this.completedTasks.length())
        print("Failed tasks:", this.failedTasks.length())
        
        print("\nNode status:")
        for (let i = 0; i < this.nodes.length(); i++) {
            let node = this.nodes[i]
            let status = node.isActive ? "ACTIVE" : "INACTIVE"
            print("- Node", node.id + ":", status, "- Load:", node.load)
        }
    }
    
    restartInactiveNodes() {
        print("\nRestarting inactive nodes...")
        for (let i = 0; i < this.nodes.length(); i++) {
            let node = this.nodes[i]
            if (!node.isActive) {
                node.restart()
            }
        }
    }
}

func main() {
    // Create coordinator
    let coordinator = Coordinator()
    
    // Create nodes
    for (let i = 1; i <= 4; i++) {
        let node = Node("N" + i)
        coordinator.addNode(node)
    }
    
    // Create tasks (some with simulated errors)
    let tasks = [
        { name: "Task-1", duration: 1 },
        { name: "Task-2", duration: 2 },
        { name: "Task-error-1", duration: 1 },  // Will cause an error
        { name: "Task-3", duration: 1 },
        { name: "Task-4", duration: 2 },
        { name: "Task-error-2", duration: 1 },  // Will cause an error
        { name: "Task-5", duration: 1 },
        { name: "Task-6", duration: 1 }
    ]
    
    for (let task in tasks) {
        coordinator.addTask(task)
    }
    
    // Process tasks
    coordinator.processTasks()
    
    // Wait for them to finish
    sleep(8)
    
    // Show statistics
    coordinator.showStatistics()
    
    // Restart nodes and process failed tasks
    coordinator.restartInactiveNodes()
    
    sleep(2)
    coordinator.showStatistics()
}
```

## Best Practices

### 1. Concurrency
- ✅ Use goroutines for independent tasks
- ✅ Avoid shared state when possible
- ✅ Implement timeouts for operations that might hang
- ✅ Use patterns like worker pools for better control

### 2. Error Handling
- ✅ Always use finally for cleanup
- ✅ Be specific in error messages
- ✅ Implement retry logic for unstable operations
- ✅ Use circuit breakers for external services

### 3. Debugging
- ✅ Add detailed logging in concurrent operations
- ✅ Use unique IDs to track operations
- ✅ Implement health checks for critical components

## Module Summary

### Concepts Learned
- ✅ Concurrent programming with goroutines
- ✅ Concurrency patterns (Worker Pool, Producer-Consumer)
- ✅ Robust error handling with try-catch-finally
- ✅ Resilience patterns (Circuit Breaker, Retry)
- ✅ Debugging concurrent applications

### Skills Developed
- ✅ Designing concurrent systems
- ✅ Implementing effective error handling
- ✅ Creating resilient applications
- ✅ Debugging concurrency issues
- ✅ Applying distributed systems patterns

### Next Module

In **Module 5**, you will learn:
- Integrated testing system (BDD)
- Creating APIs and web services
- Interacting with databases
- Deployment and distribution

Excellent work! You have mastered advanced concepts that will allow you to create robust and scalable applications.
