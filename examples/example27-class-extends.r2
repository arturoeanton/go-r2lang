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
    print("Clases que heredan de Figura");
    print("Circulo");
    c =  Circulo();
    print(c.getColor());
    print(c.getRadio());


    print("Cuadrado");
    cu =  Cuadrado();
    print(cu.getColor());
    print(cu.getLado());

    print("Cuadrado2");
    cu2 =  Cuadrado2();
    print(cu2.getColor());
    print(cu2.getLado());

    print("Triangulo");
    t =  Triangulo();
    print(t.getColor());
    print(t.getBase());
    print(t.getAltura());






}