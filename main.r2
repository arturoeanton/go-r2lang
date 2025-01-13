// 3) Ejemplo con for
function div3(a) {
    a/3
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
        fx : function (j,c)   { return data.aa+j+c; },
        aa : 1,
        bb : 2
    }

    let fx = function (a,j,c)   { return a+j+c; };


    let a = "aaa"

    println(data.fx(2,3))

    p.saludar();
    p2.saludar();
    oo = div3(6)
    println("oo:",oo)
   //


    try{

        for (let i in range(1,4)) {
                print("(range)for i:", i, " value:", $v);
        }


        throw "holclkajncdlk  ";
    }catch(e){
        print("catch Error:", e);
    }finally {
        print("Finally");
    }
}
