// Función para simular una tarea que toma tiempo
func task(id, duration) {
    print("Task", id, "iniciada, duración:", duration, "segundos")
    sleep(duration)
    print("Task", id, "completada")
}

// Función principal
func main() {
    print("Inicio del programa principal")

    // Crear un semáforo con 2 permisos
    let sem = semaphore(2)

    // Crear un monitor
    let mon = monitor()

    // Iniciar 5 goroutines que intentan ejecutar tareas
    for (let i=1; i<=5; i=i + 1) {
        let id = i
        r2(func (id) {
            acquire(sem) // Adquirir permiso del semáforo
            lock(mon)    // Adquirir lock del monitor

            // Sección crítica: imprimir el inicio de la tarea
            print("Monitor: Task", id, "está ejecutándose")

            // Liberar el lock del monitor
            unlock(mon)

            // Ejecutar la tarea
            task(id, 2)

            // Adquirir lock nuevamente para modificar estado
            lock(mon)
            print("Monitor: Task", id, "ha terminado")
            unlock(mon)

            release(sem) // Liberar permiso del semáforo
        },id)
    }
    print("Fin del programa principal")
}