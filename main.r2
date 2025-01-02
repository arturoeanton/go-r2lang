obj Tareas {
    let lista;

    func agregar(tarea) {
        self.lista["tarea1"] = tarea;
    }

    func listar() {
        println("Listando tareas...");
        println(self.lista);
        return self.lista;
    }
}

let tareas = Tareas();
tareas.lista = {};
httpGet("/agregar/:tarea", func(pathVars) {
    tareas.agregar(pathVars["tarea"]);
    println(tareas.lista);
    return HttpResponse("Tarea agregada");
});

httpGet("/listar", func() {
    tareas.listar();
    return HttpResponse(JSON(tareas));
});

httpServe(":8080");