// Función para simular una tarea que toma tiempo
func task(id, duration) {
    print("Task", id, "iniciada, duración:", duration, "segundos")
    sleep(duration)
    print("Task", id, "completada")
}

// Función principal
func main() {
    print("Inicio del programa principal")

    // Crear un semáforo con 1 permisos
    let sem = semaphore(1)


    // Iniciar 5 goroutines que intentan ejecutar tareas
    for (let i=1; i<=5; i=i + 1) {
        let id = i
        r2(func (id) {
            acquire(sem) // Adquirir permiso del semáforo
            // Ejecutar la tarea
            task(id, 2)
            release(sem) // Liberar permiso del semáforo
        }, id)
    }

    print("Fin del programa principal")
}