func main (){
    let arr = [1, 2, 3];
    std.print(arr); //
    arr = arr.add( 4, 5);
    std.print("arr = arr.add( 5, 4);", arr); // [1, 2, 3, 4, 5]

    arr = arr.del( 0,2,4);
    std.print("arr.delete( 0,2,4);",  arr); // [2, 4]

    arr = arr.add( 1, 3, 4, 5);
    std.print("arr = arr.add( 1, 3, 4, 5);", arr); // [2, 4, 1, 3, 4, 5]

    std.print("arr.index(4)", arr.index(4)); // 1
    std.print("arr.indexes(4)", arr.indexes(4)); // [1,4]

    std.print("range 1-3 ", collections.range(1, 3)); // [1, 2, 3]
    std.print("repeat 1-3", collections.repeat(3, "hola")); // ["hola", "hola", "hola"]

    std.print("for loop");
    for (let i = 0; i < arr.len(); i++){
        std.print(i);
    }

    std.print("for in loop");
    for (j in arr){
        std.print($k, arr[j], $v);
    }

   let mapa = { saludo: "hola", despedida: "adios" };

    std.print("for in loop map");
    for ( i in mapa){
        std.print($k, mapa[i], $v);
    }

    arr = arr.insert_at( 2, 0); // [5, 4, 0, 4, 3, 2, 1]
    std.print("arr = arr.insert_at( 2, 0); >>", arr); // [5, 4, 0, 4, 3, 2, 1]
    std.print("arr.indexes(func (v){v==4}).len() >>", arr.indexes(func (v){v==4}).len()); // 2
}