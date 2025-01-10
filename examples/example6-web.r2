func handleUser(pathVars, method, body) {
  let id = pathVars["id"];
  return "Hola usuario con id = " + id + " y method " + method;
}

obj Persona {
    let nombre;
    let edad;
    let pp;
    func init(n, e) {
        self.nombre = n;
        self.edad = e;
    }
    func saludar() {
        return sprintf("Hola, soy", self.nombre, "y tengo", self.edad, "años.");
    }
}

func main(){
    httpGet("/users/:id", handleUser);

    httpGet("/users", func(){
                            let p = Persona();
                            p.init("Carlos", 30);


                            let p1 = Persona();
                            p1.init("Arturo", 44);
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
      httpServe(":8080");
}