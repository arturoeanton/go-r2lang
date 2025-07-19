// example11.r2

func main() {
    std.print("=== Prueba de r2math ===");

    // Usamos pi, e
    std.print("pi =>", PI);
    std.print("e =>", E);

    // Trig
    let x = 1;
    let s = math.sin(x);
    let c = math.cos(x);
    std.print("sin(1) =>", s, " cos(1) =>", c);

    // log y exp
    let lx = math.log(x + 10);
    std.print("log(10+1) =>", lx);
    let ex = math.exp(2);
    std.print("exp(2) =>", ex);

    // sqrt, pow
    let sq = math.sqrt(9);
    std.print("sqrt(9) =>", sq);
    let pw = math.pow(2, 8);
    std.print("pow(2, 8) =>", pw);

    // abs, floor, ceil, round
    let neg = -3.7;
    std.print("abs(-3.7) =>", math.abs(neg));
    std.print("floor(-3.7) =>", math.floor(neg));
    std.print("ceil(-3.7) =>", math.ceil(neg));
    std.print("round(-3.7) =>", math.round(neg));

    // max, min
    std.print("max(10, 20) =>", math.max(10, 20));
    std.print("min(10, 20) =>", math.min(10, 20));

    // hypot
    std.print("hypot(3,4) =>", math.hypot(3, 4)); // => 5

    std.print("=== Fin de example11.r2 ===");
}