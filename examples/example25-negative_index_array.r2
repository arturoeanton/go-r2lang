
func main() {
    let aa = ["a", "b", "c"];
    let bb = [1, 2, 3];
    std.print("aa:", aa[-3]);
    let cc = aa + bb
    std.print("cc:", "hola" + cc );
    let mm = { "nombre": "Carlos", "edad": 30 };
    mm["pp"] = "hola";
    std.print("mm:", mm.nombre, mm.edad, mm.pp);
}
// output:
// aa: a
// cc: [hola a b c 1 2 3]
// mm: Carlos 30 hola

