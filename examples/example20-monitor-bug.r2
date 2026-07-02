// Función para simular una tarea que toma tiempo
func task(id, duration) {
    std.print("Task", id, "iniciada, duración:", duration, "segundos")
    std.sleep(duration)
    std.print("Task", id, "completada")
}

// Función principal
func main() {
    std.print("Inicio del programa principal")

    // Crear un semáforo con 2 permisos
    let sem = goroutine.semaphore(2)

    // Crear un monitor
    let mon = goroutine.monitor()

    // Iniciar 5 goroutines que intentan ejecutar tareas
    for (let i=1; i<=5; i=i + 1) {
        let id = i
        r2(func (id) {
            goroutine.acquire(sem) // Adquirir permiso del semáforo
            goroutine.lock(mon)    // Adquirir lock del monitor

            // Sección crítica: imprimir el inicio de la tarea
            std.print("Monitor: Task", id, "está ejecutándose")

            // Liberar el lock del monitor
            goroutine.unlock(mon)

            // Ejecutar la tarea
            task(id, 2)

            // Adquirir lock nuevamente para modificar estado
            goroutine.lock(mon)
            std.print("Monitor: Task", id, "ha terminado")
            goroutine.unlock(mon)

            goroutine.release(sem) // Liberar permiso del semáforo
        },id)
    }
    std.print("Fin del programa principal")
}