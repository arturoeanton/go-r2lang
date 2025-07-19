func main() {
    // Creamos y escribimos un archivo
    let path = "test.txt";
    io.writeFile(path, "Hola mundo\nEsta es una prueba.\n");

    // Verificamos si existe
    let existe = io.exists(path);
    std.print("fileExists?", existe);

    // Leemos su contenido
    let contenido = io.readFile(path);
    std.print("Contenido de test.txt:\n", contenido);

    // Lo abrimos en append
    io.appendFile(path, "Linea agregada!\n");

    // Leemos de nuevo
    let contenido2 = io.readFile(path);
    std.print("Contenido tras append:\n", contenido2);

    // Cambiamos nombre
    io.renameFile("test.txt", "test_renamed.txt");
    std.print("Renombrado a test_renamed.txt");

    // Creamos un directorio
    io.mkdir("demoDir");
    std.print("demoDir creado.");

    // Leemos contenido de directorio actual
    let files = io.listDir(".");
    std.print("Archivos en '.' =>", files);

    // Obtenemos la ruta absoluta de test_renamed.txt
    let abs = io.absPath("test_renamed.txt");
    std.print("Ruta absoluta =>", abs);

    // Finalmente, lo borramos
    io.rmFile("test_renamed.txt");
    std.print("test_renamed.txt borrado.");

    // Quitamos demoDir
    io.rmDir("demoDir");
    std.print("demoDir borrado.");

    std.print("Fin del script prueba.r2");
}