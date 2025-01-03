// example.r2

func main() {
    let add = func(x, y) {
        return x + y;
    };

    let r = add(10, 20);
    print("Resultado =>", r);

    // Otro ejemplo: función sin args
    let greet = func() {
        print("Hola desde función anónima!");
    };
    greet();

    // Una con variables capturadas (sencillo)
    let base = 100;
    let addBase = func(x) {
        return x + base; // 'base' está en el env
    };
    let res2 = addBase(50);
    print("res2 =>", res2); // 150


        // 1) Función anónima
        let suma = func(a, b) {
            return a + b;
        };
        let r = suma(10, 20);
        print("suma(10,20) =>", r);

        // 2) Otra anónima sin args
        let saluda = func() {
            print("Hola sin args");
        };
        saluda();

        // 3) Prueba de == con numerico
        let eq1 = (2 == 2);
        let eq2 = (2 == 2.0);
        let eq3 = (2.0 == 3.0);
        print("2 == 2 =>", eq1);
        print("2 == 2.0 =>", eq2);
        print("2.0 == 3.0 =>", eq3);
}