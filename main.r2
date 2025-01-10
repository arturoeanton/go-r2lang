// 3) Ejemplo con for
function div() {

 }

obj Persona   {
    let nombre;
        let edad;

        method constructor(n, e) {
            self.nombre = n;
            self.edad = e;
        }

        method saludar() {
            println("Hola, soy", self.nombre, "y tengo", self.edad, "años.");
        }
}

function main() {
    print("main");
    let p = Persona("Carlos", 30);
    print("p:", p);
    p.saludar();
    div()
    throw "holclkajncdlk a ";

    try{
        for (let i in range(1,4)) {
                print("(range)for i:", i, " value:", $v);
        }
    }catch(e){
        print("catch Error:", e);
    }finally {
        print("Finally");
    }
}
