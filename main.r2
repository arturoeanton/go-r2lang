
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


    let fecha0 = @1980-09-12; // Fecha sin hora
    print("Fecha 0: ", fecha0);
    let fecha1 = @2024-12-25;                    // Solo fecha
    let fecha2 = @"2024-12-25 14:30:00";         // Fecha y hora
    let fecha3 = @"2024-12-25T14:30:00Z";        // ISO 8601
    let fecha4 = @"2024-12-25T14:30:00-05:00";   // Con zona horaria

    print("Fecha 1: ", fecha1);
    print("Fecha 2: ", fecha2);
    print("Fecha 3: ", fecha3);
    print("Fecha 4: ", fecha4);

    let fecha5 = @2024-12-25T14:30:00; // Fecha y hora con zona horaria
    print("Fecha 5: ", fecha5);

    let fecha6 = @2024-12-25T14:30:00+02:00; // Fecha y hora con zona horaria positiva
    print("Fecha 6: ", fecha6);

    let fecha7 = @2024-12-25T14:30:00-02:00; // Fecha y hora con zona horaria negativa
    print("Fecha 7: ", fecha7);

    //fecha mayor 
    if (fecha2 > fecha1) {
        print("Fecha 2 es mayor que Fecha 1");
    } else {
        print("Fecha 1 es mayor o igual que Fecha 2");
    }


    //fecha menor
    if (fecha0 < fecha2) {
        print("Fecha 0 es menor que Fecha 2");
    } else {
        print("Fecha 2 es menor o igual que Fecha 0");
    }
}