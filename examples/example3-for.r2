// 3) Ejemplo con for
func main() {
    for (let i=0; i<3; i=i+1) {
        print("for i:", i);
    }

    for (let i in range(1,4)) {
            print("(range)for i:", i, " value:", $v);
    }

    let arr = [1, 2, 3];
    for (let i=0; i<arr.len(); i=i+1) {
           print("for i:", i, " value:",  arr[i]);
    }

    for (let i in arr) {
        print("for i:", i, " value:", $v);
    }



}
// Output:
// for i: 0
// for i: 1
// for i: 2
// (range)for i: 1 value: 1
// (range)for i: 2 value: 2
// (range)for i: 3 value: 3
// for i: 0 value: 1
// for i: 1 value: 2
// for i: 2 value: 3
// forin i: 1 value: 1
// forin i: 2 value: 2
// forin i: 3 value: 3
