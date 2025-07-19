// example_std.r2

// Muestra un uso de algunas funciones
func main() {
    std.print("Probando funciones de r2std...");

    // 1) typeOf
    let val = 123;
    std.print("typeOf(val) =>", std.typeOf(val)); // "float64"

    // 2) len con string
    let texto = "Hola Mundo";
    let l = std.len(texto);
    std.print("len(texto) =>", l); // 10

    // 3) sleep
    std.print("Durmiendo 2 segundos...");
    std.sleep(2);
    std.print("Desperté!");

    // 4) parseInt
    let strNum = "42";
    let num = std.parseInt(strNum);
    std.print("parseInt('42') =>", num);

    // 5) toString
    let s2 = std.toString(num);
    std.print("toString(42) =>", s2);

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
    std.print("datos => nombre:", n, "edad:", e);

    // 7) range
    let r = collections.range(1, 5);
    std.print("range(1,5) =>", r); // [1, 2, 3, 4]

    // 8) now
    let fecha = std.now();
    std.print("Hora actual =>", fecha);



    let splitted = std.split("uno|dos|tres", "|");
    std.print("split('uno|dos|tres','|') =>", splitted); // ["uno","dos","tres"]


    std.print("Fin del script example_std.r2");
}