func main (){
    let arr = [1, 2, 3];
    arr = append(arr, 4, 5);
    print(arr); //

    print("-",arr[0]); // 1
    print(" ",arr[1]); // 2
    print("-",arr[2]); // 3
    print(" ",arr[3]); // 4
    print("-",arr[4]); // 5

    arr = delete(arr, 0,2,4);
    print(arr); // [2, 4]

    arr = append(arr, 1, 3, 4, 5);
    print(arr); // [2, 4, 1, 3, 4, 5]

    arr = sort(arr);
    print(arr); // [1, 2, 3, 4, 4, 5]

    arr = reverse(arr);
    print(arr); // [5, 4, 4, 3, 2, 1]

    print("len:", len(arr)); // 6
    print("distinct:", distinct(arr)); // 0
    print("count of 4 ", count(arr, 4)); // 2
    print("index of 4 ", index(arr, 4)); // 1
    print("indexes of 4 ", indexes(arr, 4)); // 4

    print("range 1-3 ", range(1, 3)); // [4, 4, 3]
    print("repeat 1-3", repeat(3, "hola")); // [4, 4, 3]


    print("insert 0 in position 2:", insert(arr, 0, 2 )); // [5, 4, 0, 4, 3, 2, 1]

    print("for loop");
    for (let i = 0; i < len(arr); i++){
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





}