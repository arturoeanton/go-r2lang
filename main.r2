// 3) Ejemplo con for
function div() {

 }

class Persona   {
    let nombre;
    let edad;
    constructor(n, e) {
        this.nombre = n;
        this.edad = e;
    }
    saludar() {
        println("Hola, soy", this.nombre, "y tengo", this.edad, "años.");
    }
}

function main() {
    print("main");
    let p = Persona("Carlos", 30);
    let   p2 = Persona("Martin", 32);

    let data = {
        fx : (j,c) =>  { return data.aa+j+c; },
        aa : 1,
        bb : 2
    }

    let fx = (a,j,c) =>  { return a+j+c; };



    println(data.fx(2,3))

    p.saludar();
    p2.saludar();
    div()
   //

    try{

        for (let i in range(1,4)) {
                print("(range)for i:", i, " value:", $v);
        }
        throw "holclkajncdlk a ";
    }catch(e){
        print("catch Error:", e);
    }finally {
        print("Finally");
    }
}
