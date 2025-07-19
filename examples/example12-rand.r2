// exampleRand.r2

func main() {
    // Inicializamos semilla con la hora
    rand.randInit();

    std.print("=== Prueba de r2rand ===");

    // 1) randFloat
    let f = rand.randFloat();
    std.print("randFloat() =>", f);

    // 2) randInt(1, 10)
    let i = rand.randInt(1, 10);
    std.print("randInt(1,10) =>", i);

    // 3) randChoice
    let arr = ["rojo","verde","azul","amarillo"];
    let choice = rand.randChoice(arr);
    std.print("randChoice(...) =>", choice);

    // 4) shuffle
    std.print("Array original =>", arr);
    rand.shuffle(arr);
    std.print("Array tras shuffle =>", arr);

    // 5) sample
    let arr2 = [1,2,3,4,5,6,7,8,9];
    let smp = rand.sample(arr2, 3);
    std.print("sample(...) =>", smp);
    std.print("Array original no cambia =>", arr2);

    // 6) Sembrar con semilla fija
    rand.randInit(42);
    std.print("Re-seeded con 42 => reproducible random");
    std.print("randInt(0,100) =>", rand.randInt(0,100));
    std.print("randFloat() =>", rand.randFloat());

    std.print("=== Fin de exampleRand.r2 ===");
}