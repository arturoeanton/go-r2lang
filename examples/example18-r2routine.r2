// 5) Ejemplo de uso de 'go' con retardos usando 'sleep'
func main() {
    print("Inicio del programa principal")

    // Goroutine que imprime números con retardos
    r2(func() {
        let i = 1;

        while (i <= 5) {
            print("Goroutine:", i);
            i = i + 1;
            sleep(1); // Esperar 1 segundo
        }
    });

    // Goroutine que imprime mensajes con retardos
    r2(func() {
        let mensajes = ["Hola", "desde", "otra", "goroutine"];
        let i = 0;
        while (i < 4) {
            print("Goroutine:", mensajes[i]);
            i = i + 1;
            sleep(2); // Esperar 2 segundos
        }
    });

    print("Fin del programa principal")
}