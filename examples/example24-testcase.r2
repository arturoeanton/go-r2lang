import "example21.lib.r2" as e

// Definir funciones de soporte para pruebas
func setup() {
    printStep("Configurando el entorno de prueba")
}

func teardown() {
    printStep("Limpiando el entorno de prueba")
}


func main() {


    TestCase "Verificar Suma de Números" {
        Given setup()
        When func(){
                  e.func2()
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

