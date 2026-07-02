// Ejemplo del Pipeline Operator |>
// Característica P4 - Composición Fluida de Funciones

std.print("🚀 R2Lang P4: Pipeline Operator |>")
std.print("=============================================")

// 1. Pipeline básico con funciones
std.print("\n🔄 1. Pipeline Básico:")
func double(x) { return x * 2 }
func addTen(x) { return x + 10 }
func square(x) { return x * x }

let result1 = 5 |> double
let result2 = 5 |> double |> addTen
let result3 = 5 |> double |> addTen |> square

std.print("5 |> double =", result1)                    // 10
std.print("5 |> double |> addTen =", result2)          // 20  
std.print("5 |> double |> addTen |> square =", result3) // 400

// 2. Pipeline con funciones lambda
std.print("\n🎯 2. Pipeline con Arrow Functions:")
let result4 = 10 |> (x => x * 3)
let result5 = 10 |> (x => x * 3) |> (x => x + 5)

std.print("10 |> (x => x * 3) =", result4)                    // 30
std.print("10 |> (x => x * 3) |> (x => x + 5) =", result5)    // 35

// 3. Pipeline para procesamiento de datos
std.print("\n📊 3. Procesamiento de Datos:")
func filterPositive(arr) {
    let result = []
    for (let i = 0; i < arr.length(); i++) {
        if (arr[i] > 0) {
            result = result.push(arr[i])
        }
    }
    return result
}

func sumArray(arr) {
    let total = 0
    for (let i = 0; i < arr.length(); i++) {
        total += arr[i]
    }
    return total
}

let numbers = [-2, 5, -1, 8, 3, -4, 7]
let positiveSum = numbers |> filterPositive |> sumArray

std.print("Numbers:", numbers)
std.print("Suma de positivos:", positiveSum)  // 23

// 4. Pipeline para transformación de strings
std.print("\n📝 4. Transformación de Strings:")
func trim(str) { return string.trim(str) }
func toLowerCase(str) { return string.toLowerCase(str) }
func capitalize(str) { return string.capitalize(str) }

let rawText = "  HELLO WORLD  "
let processed = rawText |> trim |> toLowerCase |> capitalize

std.print("Texto original: '" + rawText + "'")
std.print("Texto procesado:", processed)  // "Hello world"

// 5. Pipeline con validación
std.print("\n✅ 5. Pipeline de Validación:")
func validateNotEmpty(str) {
    if (str.length == 0) {
        panic("String cannot be empty")
    }
    return str
}

func validateMinLength(str) {
    if (str.length < 3) {
        panic("String must be at least 3 characters")
    }
    return str
}

func sanitize(str) {
    return regex.replaceAll("[<>]", str, "")
}

let userInput = "abc<script>"
let safeInput = userInput |> validateNotEmpty |> validateMinLength |> sanitize

std.print("Input del usuario:", userInput)
std.print("Input seguro:", safeInput)  // "abcscript"

// 6. Pipeline para cálculos matemáticos
std.print("\n🔢 6. Cálculos Matemáticos Complejos:")
func toRadians(degrees) { return degrees * 3.14159 / 180 }
func sin(x) { return Math.sin ? Math.sin(x) : x }  // Simplified sin
func round(x) { return Math.round ? Math.round(x) : Math.floor(x + 0.5) }

let angle = 45
// let sinValue = angle |> toRadians |> sin |> round
let sinValue = angle |> toRadians  // Simplificado para el ejemplo

std.print("Ángulo:", angle, "grados")
std.print("En radianes:", sinValue)

// 7. Pipeline para configuración de objetos
std.print("\n⚙️  7. Configuración de Objetos:")
func setTheme(config) {
    config.theme = "dark"
    return config
}

func setLanguage(config) {
    config.language = "es"
    return config
}

func enableDebug(config) {
    config.debug = true
    return config
}

let baseConfig = {name: "MyApp", version: "1.0"}
let fullConfig = baseConfig |> setTheme |> setLanguage |> enableDebug

std.print("Configuración base:", baseConfig.name, baseConfig.version)
std.print("Configuración completa:", fullConfig.theme, fullConfig.language, fullConfig.debug)

// 8. Comparación: estilo tradicional vs pipeline
std.print("\n⚖️  8. Comparación de Estilos:")

// ANTES - Llamadas anidadas difíciles de leer
std.print("❌ ANTES (anidado):")
let traditionalResult = square(addTen(double(3)))
std.print("  square(addTen(double(3))) =", traditionalResult)

// DESPUÉS - Pipeline fluido y legible
std.print("✅ DESPUÉS (pipeline):")
let pipelineResult = 3 |> double |> addTen |> square
std.print("  3 |> double |> addTen |> square =", pipelineResult)

// 9. Pipeline con manejo de errores
std.print("\n🛡️  9. Pipeline con Manejo de Errores:")
func safeDivide(divisor) {
    return (dividend) => {
        if (divisor == 0) {
            panic("División por cero")
        }
        return dividend / divisor
    }
}

func positiveOnly(x) {
    if (x <= 0) {
        panic("Solo números positivos")
    }
    return x
}

// Pipeline seguro
let safeResult = 100 |> positiveOnly |> safeDivide(20)
std.print("100 |> positiveOnly |> safeDivide(20) =", safeResult)  // 5

// 10. Casos de uso del mundo real
std.print("\n🌍 10. Caso de Uso Real - API Response:")
func extractData(response) { return response.data }
func filterActive(items) {
    let result = []
    for (let i = 0; i < items.length(); i++) {
        if (items[i].active) {
            result = result.push(items[i])
        }
    }
    return result
}

func mapToNames(items) {
    let result = []
    for (let i = 0; i < items.length(); i++) {
        result = result.push(items[i].name)
    }
    return result
}

let apiResponse = {
    status: "success",
    data: [
        {id: 1, name: "Alice", active: true},
        {id: 2, name: "Bob", active: false},
        {id: 3, name: "Charlie", active: true}
    ]
}

let activeNames = apiResponse |> extractData |> filterActive |> mapToNames
std.print("API Response -> nombres activos:", activeNames)  // ["Alice", "Charlie"]

std.print("\n✅ Pipeline Operator implementado exitosamente!")
std.print("   - Composición fluida de funciones")
std.print("   - Lectura de izquierda a derecha")
std.print("   - Compatible con funciones y lambdas")
std.print("   - Código más expresivo y mantenible")
std.print("   - Ideal para transformaciones de datos")