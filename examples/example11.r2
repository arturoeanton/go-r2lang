// example11.r2

func main() {
    print("=== Prueba de r2math ===");

    // Usamos pi, e
    print("pi =>", pi);
    print("e =>", e);

    // Trig
    let x = 1;
    let s = sin(x);
    let c = cos(x);
    print("sin(1) =>", s, " cos(1) =>", c);

    // log y exp
    let lx = log(x + 10);
    print("log(10+1) =>", lx);
    let ex = exp(2);
    print("exp(2) =>", ex);

    // sqrt, pow
    let sq = sqrt(9);
    print("sqrt(9) =>", sq);
    let pw = pow(2, 8);
    print("pow(2, 8) =>", pw);

    // abs, floor, ceil, round
    let neg = -3.7;
    print("abs(-3.7) =>", abs(neg));
    print("floor(-3.7) =>", floor(neg));
    print("ceil(-3.7) =>", ceil(neg));
    print("round(-3.7) =>", round(neg));

    // max, min
    print("max(10, 20) =>", max(10, 20));
    print("min(10, 20) =>", min(10, 20));

    // hypot
    print("hypot(3,4) =>", hypot(3, 4)); // => 5

    print("=== Fin de example11.r2 ===");
}