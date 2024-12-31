func handleUser(pathVars, method, body) {
    if (false) {
        return "H!ola usuario con id = " + pathVars["id"] + " y method " + method;
    }
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
        println("Hola, soy", self.nombre, "y tengo", self.edad, "años.");
    }
}



func main() {
  // GET /users/123 -> handleUser con pathVars = { "id" : "123" }
  httpGet("/users/:id", handleUser);
  httpGet("/users", func(){
    let p = Persona();
    p.init("Carlos", 30);
    p.saludar();

    let p1 = Persona();
    p1.init("Arturo", 44);
    p.pp = p1;

    // body = JSON(p);

    body = XML("pp",p);

    return HttpResponse(body);
    //return HttpResponse(200 ,body);
    //return HttpResponse("application/xml"  ,body);
    //return HttpResponse(200, "application/xml"  ,body);

  });


  // Levantamos servidor
  httpServe(":8080");
}