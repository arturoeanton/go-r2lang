// exampleOS.r2

func main() {
    std.print("=== r2os con exec y runProcess ===");

    let dir = os.currentDir();
    std.print("currentDir =>", dir);

    let envs = os.envList();
    std.print("Variables PATH =>", envs["PATH"]);

    // execCmd
    let out1 = os.execCmd("echo 'Hola desde execCmd'");
    std.print("execCmd =>", out1);

    // runProcess
    // Ej: un ping en background (en Linux/Mac)
    // let p = runProcess("ping google.com");
    // sleepMs(2000);
    // killProcess(p);
    // std.print("Matado el proceso ping");
    // o let w = waitProcess(p);

    // Ej: un proceso que termina rápido
    let p2 = os.runProcess("echo 'Proceso background rápido'");
    let w = os.waitProcess(p2);
    std.print("waitProcess =>", w);

    std.print("=== Fin exampleOS.r2 ===");
}