func main() {
    // Creamos y escribimos un archivo
    let path = "test.txt";
    writeFile(path, "Hola mundo\nEsta es una prueba.\n");

    // Verificamos si existe
    let existe = fileExists(path);
    print("fileExists?", existe);

    // Leemos su contenido
    let contenido = readFile(path);
    print("Contenido de test.txt:\n", contenido);

    // Lo abrimos en append
    appendFile(path, "Linea agregada!\n");

    // Leemos de nuevo
    let contenido2 = readFile(path);
    print("Contenido tras append:\n", contenido2);

    // Cambiamos nombre
    renameFile("test.txt", "test_renamed.txt");
    print("Renombrado a test_renamed.txt");

    // Creamos un directorio
    makeDir("demoDir");
    print("demoDir creado.");

    // Leemos contenido de directorio actual
    let files = readDir(".");
    print("Archivos en '.' =>", files);

    // Obtenemos la ruta absoluta de test_renamed.txt
    let abs = absPath("test_renamed.txt");
    print("Ruta absoluta =>", abs);

    // Finalmente, lo borramos
    removeFile("test_renamed.txt");
    print("test_renamed.txt borrado.");

    // Quitamos demoDir
    removeFile("demoDir");
    // O si es un directorio, en Linux no se borra con removeFile.
    // Este ejemplo es solo para mostrar que removeFile falla en un directorio
    // (dependiendo de tu SO). Podrías usar "os.RemoveAll" en otra builtin.

    print("Fin del script prueba.r2");
}