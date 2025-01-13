// example_std.r2

// Muestra un uso de algunas funciones
func main() {
    print("Probando funciones de r2std...");

    // 1) typeOf
    let val = 123;
    print("typeOf(val) =>", typeOf(val)); // "float64"

    // 2) len con string
    let texto = "Hola Mundo";
    let l = len(texto);
    print("len(texto) =>", l); // 10

    // 3) sleep
    print("Durmiendo 2 segundos...");
    sleep(2);
    print("Desperté!");

    // 4) parseInt
    let strNum = "42";
    let num = parseInt(strNum);
    print("parseInt('42') =>", num);

    // 5) toString
    let s2 = toString(num);
    print("toString(42) =>", s2);

    // 6) vars / varsSet
    // Creamos un map en R2 (nativo se crea con obj... o con otras técnicas).
    // Para simplificar, definimos un map en Go (si tu lenguaje lo soporta).
    let datos = {};
    // En este intérprete ficticio, un "obj" se implementa con map<string,interface{}>.
    // O si no, pasa un map prehecho.

    datos["nombre"] = "Alice";
    datos["edad"] = 30;

    let e = datos["edad"];
    let n =  datos["nombre"];
    print("datos => nombre:", n, "edad:", e);

    // 7) range
    let r = range(1, 5);
    print("range(1,5) =>", r); // [1, 2, 3, 4]

    // 8) now
    let fecha = now();
    print("Hora actual =>", fecha);



    let splitted = split("uno|dos|tres", "|");
    print("split('uno|dos|tres','|') =>", splitted); // ["uno","dos","tres"]


    print("Fin del script example_std.r2");
}