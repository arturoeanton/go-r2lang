// exampleRand.r2

func main() {
    // Inicializamos semilla con la hora
    randInit();

    print("=== Prueba de r2rand ===");

    // 1) randFloat
    let f = randFloat();
    print("randFloat() =>", f);

    // 2) randInt(1, 10)
    let i = randInt(1, 10);
    print("randInt(1,10) =>", i);

    // 3) randChoice
    let arr = ["rojo","verde","azul","amarillo"];
    let choice = randChoice(arr);
    print("randChoice(...) =>", choice);

    // 4) shuffle
    print("Array original =>", arr);
    shuffle(arr);
    print("Array tras shuffle =>", arr);

    // 5) sample
    let arr2 = [1,2,3,4,5,6,7,8,9];
    let smp = sample(arr2, 3);
    print("sample(...) =>", smp);
    print("Array original no cambia =>", arr2);

    // 6) Sembrar con semilla fija
    randInit(42);
    print("Re-seeded con 42 => reproducible random");
    print("randInt(0,100) =>", randInt(0,100));
    print("randFloat() =>", randFloat());

    print("=== Fin de exampleRand.r2 ===");
}