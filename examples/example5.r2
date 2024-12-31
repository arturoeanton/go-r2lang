// 5) Ejemplo con obj y self
obj Persona {
    let nombre;
    let edad;

    func init(n, e) {
        self.nombre = n;
        self.edad = e;
    }

    func saludar() {
        println("Hola, soy", self.nombre, "y tengo", self.edad, "años.");
    }
}

func main() {
    let p = Persona();
    p.init("Carlos", 30);
    p.saludar();
}

// Output:
// Hola, soy Carlos y tengo 30 años.