// examplePrint.r2

func main() {
print("=== r2print Demo ===");

printRepeat("R2 ", 5);

printBox("Hola, R2!", 15);

let arr = [1,2,3];
debugInspect(arr);

let mapa = {};
mapa["clave"] = "123";
debugInspect(mapa);


    // Ejemplo 1: Colores
    printColor("Hola en rojo", "red");
    printColor("Hola en verde", "green");
    printColor("Hola normal", "reset");

    // Ejemplo 2: Barra de progreso
    printProgress("Cargando", 10, 150);

    // Ejemplo 3: Tabla
    let rows = [
        ["Nombre", "Edad", "Ciudad"],
        ["Ana", 30, "Madrid"],
        ["Bob", 25, "Lima"],
        ["Carlos", 35, "Buenos Aires"]
    ];
    printTable(rows);

    // Ejemplo 4: Alineación
    printAlign("Texto Left", "left", 20);
    printAlign("Texto Right", "right", 20);
    printAlign("Texto Center", "center", 20);

print("=== Fin ===");
}