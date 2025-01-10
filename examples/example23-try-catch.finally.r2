// 3) Ejemplo con for

func main() {


    try{
        for (let i in range(1,4)) {
                print("(range)for i:", i, " value:", $v);
        }
        throw "Exception"

    }catch(e){
        print("Error:", e);
    }finally {
        print("Finally");
    }
}
