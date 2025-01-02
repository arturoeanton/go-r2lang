// 3) Ejemplo con for
func main() {
    for (let i=0; i<3; i=i+1) {
        print("for i:", i);
    }

    rr = range(1,4)
    for (let i in rr) {
            print("(range)for i:", i, " value:", $v);
    }

    let arr = [1, 2, 3];
    for (let i=0; i<len(arr); i=i+1) {
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