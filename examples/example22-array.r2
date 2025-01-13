func main (){
    let arr = [1, 2, 3];
    print(arr); //
    arr = arr.add( 4, 5);
    print("arr = arr.add( 5, 4);", arr); // [1, 2, 3, 4, 5]

    arr = arr.del( 0,2,4);
    print("arr.delete( 0,2,4);",  arr); // [2, 4]

    arr = arr.add( 1, 3, 4, 5);
    print("arr = arr.add( 1, 3, 4, 5);", arr); // [2, 4, 1, 3, 4, 5]

    print("arr.index(4)", arr.index(4)); // 1
    print("arr.indexes(4)", arr.indexes(4)); // [1,4]

    print("range 1-3 ", range(1, 3)); // [1, 2, 3]
    print("repeat 1-3", repeat(3, "hola")); // ["hola", "hola", "hola"]

    print("for loop");
    for (let i = 0; i < arr.len(); i++){
        print(i);
    }

    print("for in loop");
    for (let i in arr){
        print($k, arr[i], $v);
    }

   let mapa = { saludo: "hola", despedida: "adios" };

    print("for in loop map");
    for (let i in mapa){
        print($k, mapa[i], $v);
    }

    arr = arr.insert_at( 2, 0); // [5, 4, 0, 4, 3, 2, 1]
    print("arr = arr.insert_at( 2, 0); >>", arr); // [5, 4, 0, 4, 3, 2, 1]
    print("arr.indexes(func (v){v==4}).len() >>", arr.indexes(func (v){v==4}).len()); // 2
}