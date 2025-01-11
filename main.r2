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

        a = [1,2,3,4,5];
        println(">>>",a.length); // tiene que ser 1,2,3,4,5
        a = a.map(func(v){v*2}).filter(func(v){v<10}).reduce(func(v,c){v+c;});
        print(a); // tiene que ser 20  -> de map 2,4,6,8,10 -> de filter 2,4,6,8 -> de reduce 20

        throw "holclkajncdlk  ";
    }catch(e){
        print("catch Error:", e);
    }finally {
        print("Finally");
    }
}
