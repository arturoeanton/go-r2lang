// example.r2

func main() {

 let a = [1,2,3,4,5];
        println(">>>",a.length); // tiene que ser 1,2,3,4,5
        a = a.map(func (v){ v*2}).filter(func (v){ v<10}).reduce(func (v,c) {v+c;});
        print(a);

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


           a = [1,2,3,4,5];
           println(">>>",a.length); // tiene que ser 1,2,3,4,5
           a = a.map(func (v){ v*2}).filter(func (v){ v<10}).reduce(func (v,c) {v+c;});
           print(a); // tiene que ser 20  -> de map 2,4,6,8,10 -> de filter 2,4,6,8 -> de reduce 20




}