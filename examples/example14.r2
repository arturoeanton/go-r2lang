// exampleOS.r2

func main() {
    println("=== r2os con exec y runProcess ===");

    let dir = currentDir();
    println("currentDir =>", dir);

    let envs = envList();
    println("Variables PATH =>", envs["PATH"]);

    // execCmd
    let out1 = execCmd("echo 'Hola desde execCmd'");
    println("execCmd =>", out1);

    // runProcess
    // Ej: un ping en background (en Linux/Mac)
    // let p = runProcess("ping google.com");
    // sleepMs(2000);
    // killProcess(p);
    // println("Matado el proceso ping");
    // o let w = waitProcess(p);

    // Ej: un proceso que termina rápido
    let p2 = runProcess("echo 'Proceso background rápido'");
    let w = waitProcess(p2);
    println("waitProcess =>", w);

    println("=== Fin exampleOS.r2 ===");
}