// Función para simular una tarea que toma tiempo
func task(id, duration) {
    std.print("Task", id, "iniciada, duración:", duration, "segundos")
    std.sleep(duration)
    std.print("Task", id, "completada")
}

// Función principal
func main() {
    std.print("Inicio del programa principal")

    // Crear un semáforo con 1 permisos
    let sem = goroutine.semaphore(1)


    // Iniciar 5 goroutines que intentan ejecutar tareas
    for (let i=1; i<=5; i=i + 1) {
        let id = i
        go.r2(func (id) {
            goroutine.acquire(sem) // Adquirir permiso del semáforo
            // Ejecutar la tarea
            task(id, 2)
            goroutine.release(sem) // Liberar permiso del semáforo
        }, id)
    }

    std.print("Fin del programa principal")
}