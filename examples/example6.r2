func handleUser(pathVars, method, body) {
  let id = pathVars["id"];
  return "Hola usuario con id = " + id + " y method " + method;
}

func main() {
  // GET /users/123 -> handleUser con pathVars = { "id" : "123" }
  httpAddRoute("GET", "/users/:id", "handleUser");
  // Levantamos servidor
  httpServe(":8080");
}