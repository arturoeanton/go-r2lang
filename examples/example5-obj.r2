// 5) Ejemplo con obj y self
class Persona {
    let nombre;
    let edad;

    method init(n, e) {
        this.nombre = n;
        this.edad = e;
    }

    method saludar() {
        println("Hola, soy", this.nombre, "y tengo", this.edad, "años.");
    }
}

func main() {
    let p = Persona();
    p.init("Carlos", 30);
    p.saludar();
}

// Output:
// Hola, soy Carlos y tengo 30 años.