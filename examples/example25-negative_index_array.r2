
func main() {
    let aa = ["a", "b", "c"];
    let bb = [1, 2, 3];
    println("aa:", aa[-3]);
    let cc = aa + bb
    print("cc:", "hola" + cc );
    let mm = { "nombre": "Carlos", "edad": 30 };
    mm["pp"] = "hola";
    print("mm:", mm.nombre, mm.edad, mm.pp);
}
// output:
// aa: a
// cc: [hola a b c 1 2 3]
// mm: Carlos 30 hola

