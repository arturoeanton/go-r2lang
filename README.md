
<div align="center">
  <br />
  <h1>R2Lang</h1>
  <p>
    <b>Write elegant tests, scripts, and applications with a language that blends simplicity and power.</b>
  </p>
  <br />
</div>

<div align="center">

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GitHub stars](https://img.shields.io/github/stars/arturoeanton/go-r2lang.svg?style=social&label=Star)](https://github.com/arturoeanton/go-r2lang)
[![GitHub forks](https://img.shields.io/github/forks/arturoeanton/go-r2lang.svg?style=social&label=Fork)](https://github.com/arturoeanton/go-r2lang)
[![GitHub issues](https://img.shields.io/github/issues/arturoeanton/go-r2lang.svg)](https://github.com/arturoeanton/go-r2lang/issues)
[![Contributors](https://img.shields.io/github/contributors/arturoeanton/go-r2lang.svg)](https://github.com/arturoeanton/go-r2lang/graphs/contributors)

</div>

---

**R2Lang** is a dynamic, interpreted programming language written in Go. It's designed to be simple, intuitive, and powerful, with a syntax heavily inspired by JavaScript and first-class support for **Behavior-Driven Development (BDD)**.

Whether you're writing automation scripts, building a web API, or creating a robust testing suite, R2Lang provides the tools you need in a clean and readable package.

## ✨ Key Features

| Feature                 | Description                                                                                                 | Example                                                              |
| ----------------------- | ----------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------- |
| **🧪 Built-in BDD Testing** | Write tests as specifications with a native `TestCase` syntax. No external frameworks needed.               | `TestCase "User should log in" { Given ... When ... Then ... }`      |
| **🚀 Simple & Familiar**    | If you know JavaScript, you'll feel right at home. This makes it incredibly easy to pick up and start coding. | `let message = "Hello, World!"; print(message);`                     |
| **⚡ Concurrent**          | Leverage the power of Go's goroutines with a simple `r2()` function to run code in parallel.                | `r2(myFunction, "arg1");`                                            |
| **🧱 Object-Oriented**     | Use classes, inheritance (`extends`), and `this` to structure your code in a clean, object-oriented way.    | `class User extends Person { ... }`                                  |
| **🌐 Web Ready**            | Create web servers and REST APIs with a built-in `http` library that feels like Express.js.                 | `http.get("/users", func(req, res) { res.json(...) });`               |
| **🧩 Easily Extendable**   | Written in Go, R2Lang can be easily extended with new native functions and libraries.                       | `env.Set("myNativeFunc", r2lang.NewBuiltinFunction(...));`            |

---

## 🚀 Getting Started

### Prerequisites

- **Go**: Version 1.23 or higher.

### Installation & "Hello, World!"

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/arturoeanton/go-r2lang.git
    cd go-r2lang
    ```

2.  **Build the interpreter:**
    ```bash
    go build -o r2lang main.go
    ```

3.  **Create your first R2Lang file (`hello.r2`):**
    ```r2
    func main() {
        print("Hello, R2Lang! 🚀");
    }
    ```

4.  **Run it!**
    ```bash
    ./r2lang hello.r2
    # Output: Hello, R2Lang! 🚀
    ```

---

## 📚 Documentation & Full Course

Ready to dive deeper? We have a complete, module-by-module course to take you from beginner to expert.

-   [**Read the Full Course (English)**](./docs/en/README.md)
-   [**Leer el Curso Completo (Español)**](./docs/es/README.md)

The documentation covers everything from basic syntax to advanced topics like concurrency, error handling, and web development.

---

## 💖 Contributing

**We are actively looking for contributors!** Whether you're a seasoned developer, a documentation writer, or just enthusiastic about new programming languages, we'd love your help.

Here’s how you can contribute:

1.  **Find an issue:** Check out our [**Issues**](https://github.com/arturoeanton/go-r2lang/issues) and look for `good first issue` or `help wanted` tags.
2.  **Explore the Roadmap:** See our [**Technical Roadmap**](./docs/en/roadmap.md) for long-term goals and big features we need help with.
3.  **Improve Documentation:** Found a typo or a section that could be clearer? Let us know!
4.  **Submit a Pull Request:**
    -   Fork the repository.
    -   Create a new branch (`git checkout -b feature/my-awesome-feature`).
    -   Commit your changes.
    -   Open a Pull Request!

We believe in a welcoming and supportive community. No contribution is too small!

---

## 🗺️ Project Roadmap

We have big plans for R2Lang! Our goal is to make it a fast, reliable, and feature-rich language for a wide range of applications.

Key areas of focus include:

-   **🚀 Performance Revolution:** Implementing a bytecode VM and eventually a JIT compiler for significant speed boosts.
-   **🧠 Advanced Features:** Adding pattern matching, a more sophisticated type system, and advanced concurrency models.
-   **🛠️ Richer Standard Library:** Expanding the built-in libraries for databases, file systems, and more.
-   **📦 Package Manager:** Creating a dedicated package manager for sharing and reusing R2Lang code.

For a detailed look at our plans, check out the [**Technical Roadmap**](./docs/en/roadmap.md) and our [**TODO List**](./TODO.md).

---

## 🤝 Community

-   **Report a Bug:** Found something wrong? Open an [**Issue**](https://github.com/arturoeanton/go-r2lang/issues/new).
-   **Request a Feature:** Have a great idea? Let's discuss it in the [**Issues**](https://github.com/arturoeanton/go-r2lang/issues).
-   **Ask a Question:** Don't hesitate to open an issue for questions and discussions.

---

## 📜 License

R2Lang is licensed under the **Apache License 2.0**. See the [LICENSE](./LICENSE) file for details.
