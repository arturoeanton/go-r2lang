// 5) Ejemplo con obj y self
class Persona {
    let nombre;
    let edad;

    constructor(n, e) {
        this.nombre = n;
        this.edad = e;
    }

    saludar() {
        std.print("Hola, soy", this.nombre, "y tengo", this.edad, "años.");
    }
}

func main() {
    let p = Persona("Carlos", 30);
    p.saludar();
}

// Output:
// Hola, soy Carlos y tengo 30 años.