
class Persona{
    let nombre;
    let edad;

    constructor(nombre, edad){
        std.print("Constructor de Persona");

        this.nombre = nombre;
        this.edad = edad;
    }

    velocidad(){
        v = 200 - (this.edad / 2);
        std.print("Velocidad de Persona ", v);
        return v
    }
}

class Empleado extends Persona{
    let salario;

    constructor(nombre, edad, salario){
        std.print("Constructor de Empleado");
        super.constructor(nombre, edad);
        this.salario = salario;
    }

    velocidad(){
        v = 200 - (this.edad / 2);
        std.print("Velocidad de Empleado ", v);
        return v
    }


}

class Velocista extends Empleado {
    constructor(nombre, edad, salario, speed){
        std.print("Constructor de Velocista");
        super.constructor(nombre, edad, salario);
        this.speed = speed;
    }

    velocidad(){
        super.velocidad();
        std.print("Velocidad de Velocista", this.speed);
        return this.speed;
    }

    superVelocidad(){
        std.print("Super Velocidad de Velocista", this.speed * 2);
        return this.speed * 2;
    }
}


func main(){
    std.print("\n👤 1. Clases y Objetos:");
    e =  Velocista("Juan", 30, 1000, 100);
    std.print(e.nombre);
    std.print(e.edad);
    std.print(e.salario);
    e.velocidad();
    e.superVelocidad();

    let fecha0 = @1980-09-12; // Fecha sin hora
    std.print("Fecha 0: ", fecha0);
    let fecha1 = @2024-12-25;                    // Solo fecha
    let fecha2 = @"2024-12-25 14:30:00";         // Fecha y hora
    let fecha3 = @"2024-12-25T14:30:00Z";        // ISO 8601
    let fecha4 = @"2024-12-25T14:30:00-05:00";   // Con zona horaria

    std.print("Fecha 1: ", fecha1);
    std.print("Fecha 2: ", fecha2);
    std.print("Fecha 3: ", fecha3);
    std.print("Fecha 4: ", fecha4);

    let fecha5 = @2024-12-25T14:30:00; // Fecha y hora con zona horaria
    std.print("Fecha 5: ", fecha5);

    let fecha6 = @2024-12-25T14:30:00+02:00; // Fecha y hora con zona horaria positiva
    std.print("Fecha 6: ", fecha6);

    let fecha7 = @2024-12-25T14:30:00-02:00; // Fecha y hora con zona horaria negativa
    std.print("Fecha 7: ", fecha7);

    //fecha mayor 
    if (fecha2 > fecha1) {
        std.print("Fecha 2 es mayor que Fecha 1");
    } else {
        std.print("Fecha 1 es mayor o igual que Fecha 2");
    }


    //fecha menor
    if (fecha0 < fecha2) {
        std.print("Fecha 0 es menor que Fecha 2");
    } else {
        std.print("Fecha 2 es menor o igual que Fecha 0");
    }


    std.print("\n📝 1. Identificadores Unicode:");

    // Español
    let año = 2024;
    let niño = "Antonio";
    let señorita = "María José";

    // Otros idiomas
    let 身長 = 175;          // Japonés: altura
    let имя = "Иван";        // Ruso: nombre
    let اسم = "أحمد";        // Árabe: nombre
    let όνομα = "Γιάννης";   // Griego: nombre
    let prénoms = "Jean-François"; // Francés
    std.print(`El año actual es: ${año}`);
    std.print(`Niño: ${niño}, Señorita: ${señorita}`);
    std.print(`身長 (altura): ${身長}cm`);
    std.print(`Имя (nombre en ruso): ${имя}`);
    std.print(`اسم (nombre en árabe): ${اسم}`);
    std.print(`Όνομα (nombre en griego): ${όνομα}`);
    std.print(`Prénoms français: ${prénoms}`);

    // ========================================
    // 2. STRINGS UNICODE Y ESCAPE SEQUENCES
    // ========================================
    std.print("\n🔤 2. Strings Unicode y Escape Sequences:");

    let emoji_wave = "\U0001F44B";           // 👋
    let emoji_rocket = "\U0001F680";         // 🚀
    let emoji_earth = "\U0001F30D";          // 🌍
    let spanish_chars = "\u00f1\u00e9\u00fa"; // ñéú

    std.print(`Saludando: ${emoji_wave}`);
    std.print(`Cohete: ${emoji_rocket}`);
    std.print(`Tierra: ${emoji_earth}`);
    std.print(`Caracteres españoles: ${spanish_chars}`);
}