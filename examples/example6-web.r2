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
        self.nombre = n;
        self.edad = e;
    }
    saludar() {
        return sprint("Hola, soy ", self.nombre, " y tengo ", self.edad, " años.");
    }
}

func main(){
    handler("GET","/users/:id", handleUser);

    handler("GET","/users", func(){
                            let p = Persona("Carlos", 30);


                            let p1 = Persona("Arturo", 44);
                            p.pp = p1;

                            data = {saludo :   p.saludar(), persona: p};
                            //body = JSON(data);
                            body = XML("persona",data);
                            return HttpResponse(body);
                            //return HttpResponse(200 ,body);
                            //return HttpResponse("application/xml"  ,body);
                            //return HttpResponse(200, "application/xml"  ,body);
    })


      // Levantamos servidor
      serve(":8080");
}