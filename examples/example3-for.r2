// 3) Ejemplo con for
func main() {
    for (let i=0; i<3; i=i+1) {
        std.print("for i:", i);
    }

    let rng = std.range(1,4);
    for (i in rng) {
            std.print("(range)for i:", $k, " value:", $v);
    }

    let arr = [1, 2, 3];
    for (let i=0; i<std.len(arr); i=i+1) {
           std.print("for i:", i, " value:",  arr[i]);
    }

    for (i in arr) {
        std.print("for i:", $k, " value:", $v);
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
