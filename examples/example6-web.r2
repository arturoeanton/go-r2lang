func handleUser(pathVars, method, body) {
    let id = pathVars["id"];
    let p1 = Persona("Arturo", 44);
    return "Hola usuario con id = " + id + " y method " + method + " y p1.name " + p1.nombre;
}

class Persona {
    let nombre;
    let edad;
    let pp;
    constructor(n, e) {
        this.nombre = n;
        this.edad = e;
    }
    saludar() {
        return sprint("Hola, soy ", this.nombre, " y tengo ", this.edad, " años.");
    }
}

func main(){
    http.handler("GET","/users/:id", handleUser);

    http.handler("GET","/users", func(){
                            let p = Persona("Carlos", 30);


                            let p1 = Persona("Arturo", 44);
                            p.pp = p1;

                            let data = {saludo :   p.saludar(), persona: p};
                            //body = JSON(data);
                            let body = XML("persona",data);
                            return HttpResponse(body);
                            //return HttpResponse(200 ,body);
                            //return HttpResponse("application/xml"  ,body);
                            //return HttpResponse(200, "application/xml"  ,body);
    })


      // Levantamos servidor
      http.serve(":8080");
}