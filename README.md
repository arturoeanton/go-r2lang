
# Curso Completo de r2lang generado por ia :P

¡Bienvenido al curso completo de **r2lang**! Aprende desde los conceptos básicos hasta proyectos avanzados en este poderoso y simple lenguaje de programación.

## Índice

1. [Introducción](#introducción)
2. [Conceptos Básicos](#conceptos-básicos)
   - Declaración de Variables
   - Condicionales
   - Ciclos
3. [Conceptos Intermedios](#conceptos-intermedios)
   - Objetos
   - Mapas y Arrays
4. [Conceptos Avanzados](#conceptos-avanzados)
   - Semáforos
   - Pruebas Unitarias
   - Goroutines
   - Semáforos y Monitores
5. [Proyecto Final: Gestor de Tareas Web](#proyecto-final-gestor-de-tareas-web)

---

## Introducción

r2lang es un lenguaje diseñado para ser simple, eficiente y flexible. Este curso cubre desde conceptos básicos como variables y ciclos, hasta temas avanzados como concurrencia y pruebas unitarias.

---

## Conceptos Básicos

### Declaración de Variables

Las variables en r2lang se declaran con la palabra clave `let`:

```r2
let x = 10;
let y = "Hola mundo";
print(x, y);
```

### Condicionales

Usa `if` y `else` para manejar condiciones:

```r2
let x = 15;
if (x > 10) {
    print("x es mayor que 10");
} else {
    print("x es menor o igual a 10");
}
```

### Ciclos

Usa `for` y `while` para iterar:

```r2
for (let i = 0; i < 5; i = i + 1) {
    print("Iteración:", i);
}

let j = 0;
while (j < 3) {
    print("Valor de j:", j);
    j = j + 1;
}
```

---

## Conceptos Intermedios

### Objetos

Los objetos encapsulan datos y métodos relacionados:

```r2
obj Persona {
    let nombre;
    let edad;

    func init(n, e) {
        self.nombre = n;
        self.edad = e;
    }

    func saludar() {
        print("Hola, soy", self.nombre, "y tengo", self.edad, "años.");
    }
}

let p = Persona();
p.init("Carlos", 30);
p.saludar();
```

### Mapas y Arrays

Usa mapas y arrays para almacenar colecciones de datos:

```r2
let mapa = {clave1: "valor1", clave2: "valor2"};
print(mapa["clave1"]);

let array = [1, 2, 3, 4];
print(array[2]);
```

---

## Conceptos Avanzados

### Semáforos

Controla el acceso concurrente con semáforos:

```r2
let sem = semaphore(1);

go(func() {
    acquire(sem);
    print("Acceso exclusivo");
    release(sem);
});
```

### Pruebas Unitarias

Prueba tu código para garantizar calidad:

```r2
func testSuma() {
    let x = 2 + 2;
    assertEq(x, 4, "2+2 debería ser 4");
}

testSuma();
```

### Goroutines

Ejecuta funciones concurrentemente:

```r2
go(func() {
    print("Esta función se ejecuta concurrentemente");
});
```

---

## Proyecto Final: Gestor de Tareas Web

Implementa un servidor web para gestionar tareas. El servidor permite:

- Agregar tareas.
- Listar tareas.
- Eliminar tareas.

### Código

```r2
obj Tareas {
    let lista = [];

    func agregar(tarea) {
        self.lista = append(self.lista, tarea);
    }

    func listar() {
        return self.lista;
    }
}

let tareas = Tareas();

httpGet("/agregar/:tarea", func(pathVars) {
    tareas.agregar(pathVars["tarea"]);
    return HttpResponse("Tarea agregada");
});

httpGet("/listar", func() {
    return HttpResponse(JSON(tareas.listar()));
});

httpServe(":8080");
```

### Instrucciones

1. Ejecuta el código en tu entorno r2lang.
2. Accede al servidor en `http://localhost:8080`:
   - **`/agregar/:tarea`**: Agrega una tarea.
   - **`/listar`**: Lista todas las tareas.

---

## Contribuciones

¡Si tienes ideas o mejoras para este tutorial, no dudes en hacer un fork o enviar un pull request! 🎉


