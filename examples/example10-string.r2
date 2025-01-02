// example10.r2

func main() {
    print("=== Prueba de r2string ===");

    let text = "   Hola Mundo   ";
    let t1 = trim(text);
    print("trim('   Hola Mundo   ') =>", t1);

    let upper = toUpper(t1);
    print("toUpper =>", upper);

    let lower = toLower(upper);
    print("toLower =>", lower);

    let sub = substring(lower, 0, 4);
    print("substring(lower, 0, 4) =>", sub);

    let idx = indexOf(lower, "mundo");
    print("indexOf('hola mundo','mundo') =>", idx);

    let rep = replace(lower, "mundo", "R2string");
    print("replace =>", rep);

    let splitted = split(rep, " ");
    print("split(...) =>", splitted);

    let joined = join(splitted, "-");
    print("join(...) =>", joined);

    let starts = startsWith(rep, "hola");
    let ends = endsWith(rep, "R2string");
    print("startsWith =>", starts, " endsWith =>", ends);

    let length = lengthOfString(rep);
    print("lengthOfString =>", length);

    print("=== Fin de example10.r2 ===");
}