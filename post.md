### I created a programming language in Go with built-in BDD testing. Looking for feedback and contributors!

Hey, r/golang!

For the past few months, I've been pouring my free time into a passion project: **R2Lang**, a new, dynamic programming language written entirely in Go.

My main goal was to create a language where testing isn't just a library or an afterthought, but a core, first-class citizen of the syntax. The result is a simple, JavaScript-like language with a native BDD testing framework.

**TL;DR:** I built a JavaScript-like scripting language in Go. Its main feature is a native BDD testing system (`TestCase { Given/When/Then }`). I'm looking for feedback, ideas, and collaborators to help it grow.

**GitHub Repo:** [https://github.com/arturoeanton/go-r2lang](https://github.com/arturoeanton/go-r2lang)

---

### ✨ What is R2Lang?

It's a dynamic, interpreted language designed for scripting, testing, and building simple web APIs. Think of it as a blend of JavaScript's simplicity and Go's concurrency model.

**Key Features:**

*   **🧪 Built-in BDD Testing:** This is the core feature. You can write tests using a clean `Given/When/Then` structure directly in your code, without any external frameworks.
*   **🚀 Simple & Familiar Syntax:** If you know JavaScript, you'll be writing R2Lang in minutes.
*   **⚡ Easy Concurrency:** It leverages Go's goroutines through a simple `r2()` function.
*   **🧱 Object-Oriented:** Supports classes, inheritance, and `this`.
*   **🌐 Web Ready:** Includes a built-in `http` library for creating web servers and REST APIs, inspired by Express.js.

Here’s what the BDD syntax looks like in action:

```r2
// Function to be tested
func add(a, b) {
    return a + b
}

// The test case itself
TestCase "Verify that the add function works correctly" {
    Given func() {
        print("Preparing the numbers for the test.")
        // You can set up context here
        return { a: 5, b: 10 }
    }
    
    When func(context) {
        let result = add(context.a, context.b)
        print("Executing the sum...")
        return result
    }
    
    Then func(result) {
        // assertEqual is a helper, not yet a built-in keyword
        if (result != 15) {
            throw "Assertion failed: Expected 15, got " + result
        }
        print("Validation successful!")
        return "Test passed"
    }
}
```

---

### 💖 How You Can Help

The language is functional, and I've written a full 6-module course to document it. However, it's still a young project with tons of room for improvement. I'd love to get some collaboration to take it to the next level.

I'm looking for all kinds of help:

*   **Go Developers:** To help improve the core interpreter. There are huge opportunities in performance (bytecode VM, JIT), memory management, and implementing new features from the [Roadmap](https://github.com/arturoeanton/go-r2lang/blob/main/docs/en/roadmap.md).
*   **Language Enthusiasts:** To give feedback on the syntax, features, and overall direction of the project. What do you love? What do you hate?
*   **Testers:** I need people to break it! Write some complex scripts, find edge cases, and report bugs in the [Issues](https://github.com/arturoeanton/go-r2lang/issues).
*   **Documentation Writers:** The docs are there, but they can always be improved with more examples and clearer explanations.

---

This has been a solo journey so far, and I'm really excited about the possibility of turning it into a community-driven project.

Check out the [**GitHub repository**](https://github.com/arturoeanton/go-r2lang) to see the code, the [**full documentation**](https://github.com/arturoeanton/go-r2lang/blob/main/docs/en/README.md), and the issue tracker.

Thanks for taking a look! Any feedback, questions, or stars on GitHub would be amazing. Let me know what you think!
