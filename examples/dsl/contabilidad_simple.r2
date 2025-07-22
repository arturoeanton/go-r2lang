// DSL Motor Contable - Versión Simplificada para Testing

dsl MotorContable {
    // Tokens básicos
    token("CUENTA", "[0-9]{1,4}")
    token("CONCEPTO", "[A-Za-z][A-Za-z0-9\\s]*")
    token("IMPORTE", "[0-9]+")
    token("FECHA", "[0-9]{1,2}/[0-9]{1,2}/[0-9]{4}")
    
    // Keywords
    token("DEBE", "debe")
    token("HABER", "haber")
    token("ASIENTO", "asiento")
    token("CONTRAPARTIDA", "contrapartida")
    token("POR", "por")
    
    // Regla simple
    rule("operacion_contable", ["ASIENTO", "CUENTA", "DEBE", "IMPORTE", "CONTRAPARTIDA", "CUENTA", "HABER", "IMPORTE", "POR", "CONCEPTO"], "crearAsiento")
    
    func crearAsiento(asiento, cuentaDebe, debe, importeDebe, contrapartida, cuentaHaber, haber, importeHaber, por, concepto) {
        console.log("Asiento creado:");
        console.log("DEBE - Cuenta " + cuentaDebe + ": $" + importeDebe);
        console.log("HABER - Cuenta " + cuentaHaber + ": $" + importeHaber);
        console.log("Concepto: " + concepto);
        return "OK";
    }
}

func main() {
    console.log("DSL Motor Contable - Test Simple");
    
    let motor = MotorContable;
    
    let resultado = motor.use("asiento 1110 debe 15000 contrapartida 4110 haber 15000 por Venta de productos");
    
    console.log("Resultado:", resultado);
}