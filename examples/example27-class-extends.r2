class Figura {
    constructor() {
        this.color = "red";
    }
    getColor() {
        return this.color;
    }
}

class Circulo extends Figura {
    constructor() {
        super.constructor();
        this.radio = 10;
    }
    getRadio() {
        return this.radio;
    }
}

class Cuadrado extends Figura {
    constructor() {
        super.constructor();
        this.lado = 10;
    }
    getLado() {
        return this.lado;
    }
}

class Triangulo extends Figura {
    constructor() {
        super.constructor();
        this.base = 10;
        this.altura = 10;
    }
    getBase() {
        return this.base;
    }
    getAltura() {
        return this.altura;
    }
}

class Cuadrado2 extends Cuadrado {
    constructor() {
        super.constructor();
        this.lado = 20;
    }
}

function main() {
    std.print("Clases que heredan de Figura");
    std.print("Circulo");
    c =  Circulo();
    std.print(c.getColor());
    std.print(c.getRadio());


    std.print("Cuadrado");
    cu =  Cuadrado();
    std.print(cu.getColor());
    std.print(cu.getLado());

    std.print("Cuadrado2");
    cu2 =  Cuadrado2();
    std.print(cu2.getColor());
    std.print(cu2.getLado());

    std.print("Triangulo");
    t =  Triangulo();
    std.print(t.getColor());
    std.print(t.getBase());
    std.print(t.getAltura());






}