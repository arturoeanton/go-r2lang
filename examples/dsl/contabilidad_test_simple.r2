// DSL Test Simple - Para probar funcionalidad basica
dsl TestVentaDSL {
    token("VENTA", "venta")
    token("USA", "USA")
    token("IMPORTE", "85000")
    
    rule("venta_usa", ["VENTA", "USA", "IMPORTE"], "procesarVenta")
    
    func procesarVenta(venta, region, importe) {
        console.log("VENTA PROCESADA")
        console.log("Region: " + region)
        console.log("Importe: " + importe)
        return "Venta exitosa"
    }
}

func main() {
    console.log("TEST DSL SIMPLE")
    
    let motor = TestVentaDSL
    let resultado = motor.use("venta USA 85000")
    
    console.log("Resultado: " + resultado)
    console.log("TEST COMPLETO")
}