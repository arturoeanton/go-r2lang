let true = 1==1
let false = 1==0

func testSuma() {
    let x = 2 + 2;
    assertEq(x, 4, "2+2 debería ser 4");
    assertTrue( x == 4, "x == 4");
}

func testBooleanos() {
    assertTrue(1 < 5, "1 < 5");
    assertEq(true, 1<5, "true vs cond");
}

func testFalla() {
    // Intencional
    assertEq("hola", "mundo", "esto fallará");
}

// Definir funciones de soporte para pruebas
func setup() {
    printStep("Configurando el entorno de prueba")
}

func teardown() {
    printStep("Limpiando el entorno de prueba")
}


func testCase1() {
    TestCase "Verificar Suma de Números" {
        Given setup()
        When func(){
            let resultado = 2
            let saludo = "Hola, " + "Mundo!"
            resultado = 2 + 3
            return "logica de negocio"
        }
        Then func(){
            assertEqual(resultado, 5)
            return "Validación de resultado"
        }
        And func(){
            assertEqual(saludo, "Hola, Mundo!")
            return "Validación de saludo"
        }
        Then teardown()
    }
}

// Llamamos runAllTests() al final
func main() {
    //testCase1  ()
    runAllTests();
}