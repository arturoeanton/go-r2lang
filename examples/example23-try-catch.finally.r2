// 3) Ejemplo con for

func main() {


    try{
        let a = 1/0;
    }catch(e){
        print("Error:", e);
    }finally {
        print("Finally");
    }

    try{
        throw "Example of exception";
    }catch(e){
        print("Error2:", e);
    }finally {
        print("Finally2");
    }



}
