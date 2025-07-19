// 3) Ejemplo con for

func main() {


    try{
        let a = 1/0;
    }catch(e){
        std.print("Error:", e);
    }finally {
        std.print("Finally");
    }

    try{
        throw "Example of exception";
    }catch(e){
        std.print("Error2:", e);
    }finally {
        std.print("Finally2");
    }



}
