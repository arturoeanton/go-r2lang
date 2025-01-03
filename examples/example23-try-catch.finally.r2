// 3) Ejemplo con for
func main() {
    try{
        for (let i in range(1,4)) {
                print("(range)for i:", i, " value:", $v);
        }
        let a = 2/0
    }catch(e){
        print("Error:", e);
    }finally {
        print("Finally");
    }
}
