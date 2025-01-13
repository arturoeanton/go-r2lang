
class Persona{
    let nombre;
    let edad;

    constructor(nombre, edad){
        print("Constructor de Persona");

        this.nombre = nombre;
        this.edad = edad;
    }

    velocidad(){
        v = 200 - (this.edad / 2);
        print("Velocidad de Persona ", v);
        return v
    }
}

class Empleado extends Persona{
    let salario;

    constructor(nombre, edad, salario){
        print("Constructor de Empleado");
        super.constructor(nombre, edad);
        this.salario = salario;
    }

    velocidad(){
        v = 200 - (this.edad / 2);
        print("Velocidad de Empleado ", v);
        return v
    }


}

class Velocista extends Empleado {
    constructor(nombre, edad, salario, speed){
        print("Constructor de Velocista");
        super.constructor(nombre, edad, salario);
        this.speed = speed;
    }

    velocidad(){
        super.velocidad();
        print("Velocidad de Velocista", this.speed);
        return this.speed;
    }

    superVelocidad(){
        print("Super Velocidad de Velocista", this.speed * 2);
        return this.speed * 2;
    }
}


func main(){
    e =  Velocista("Juan", 30, 1000, 100);
    print(e.nombre);
    print(e.edad);
    print(e.salario);
    e.velocidad();
    e.superVelocidad();
}