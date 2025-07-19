// example10.r2

func main() {
    std.print("=== Prueba de r2string ===");

    let text = "   Hola Mundo   ";
    let t1 = string.trim(text);
    std.print("trim('   Hola Mundo   ') =>", t1);

    let upper = string.toUpper(t1);
    std.print("toUpper =>", upper);

    let lower = string.toLower(upper);
    std.print("toLower =>", lower);

    let sub = string.substring(lower, 0, 4);
    std.print("substring(lower, 0, 4) =>", sub);

    let idx = string.indexOf(lower, "mundo");
    std.print("indexOf('hola mundo','mundo') =>", idx);

    let rep = string.replace(lower, "mundo", "R2string");
    std.print("replace =>", rep);

    let splitted = string.split(rep, " ");
    std.print("split(...) =>", splitted);

    let joined = string.join(splitted, "-");
    std.print("join(...) =>", joined);

    let starts = string.startsWith(rep, "hola");
    let ends = string.endsWith(rep, "R2string");
    std.print("startsWith =>", starts, " endsWith =>", ends);

    let length = string.lengthOfString(rep);
    std.print("lengthOfString =>", length);

    std.print("=== Fin de example10.r2 ===");
}