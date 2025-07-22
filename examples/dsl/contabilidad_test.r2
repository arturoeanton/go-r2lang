// DSL Motor Contable - Test Básico de Contexto
// Prueba simplificada para debugging

dsl MotorContable {
    // Tokens básicos
    token("CUENTA", "[0-9]{1,4}")
    token("CONCEPTO", "[A-Za-z][A-Za-z0-9\\s]*")
    token("IMPORTE", "[0-9]+")
    
    // Keywords
    token("DEBE", "debe")
    token("HABER", "haber")
    token("ASIENTO", "asiento")
    token("CONTRAPARTIDA", "contrapartida")
    token("POR", "por")
    
    // Regla simple
    rule("operacion_contable", ["ASIENTO", "CUENTA", "DEBE", "IMPORTE", "CONTRAPARTIDA", "CUENTA", "HABER", "IMPORTE", "POR", "CONCEPTO"], "crearAsiento")
    
    func crearAsiento(asiento, cuentaDebe, debe, importeDebe, contrapartida, cuentaHaber, haber, importeHaber, por, concepto) {
        let ctx = context;
        
        console.log("📝 Debug Contexto:");
        console.log("   Contexto recibido:", ctx);
        
        if (ctx) {
            console.log("   Empresa disponible:", ctx.empresa);
            if (ctx.empresa) {
                console.log("   Razon social:", ctx.empresa.razonSocial);
            }
        }
        
        console.log("📝 Asiento Creado:");
        console.log("   DEBE - Cuenta " + cuentaDebe + ": $" + importeDebe);
        console.log("   HABER - Cuenta " + cuentaHaber + ": $" + importeHaber);
        console.log("   Concepto: " + concepto);
        
        return "OK";
    }
}

func crearContextoSimple() {
    return {
        empresa: {
            razonSocial: "Test Company"
        },
        fecha: "31/12/2024"
    };
}

func main() {
    console.log("🏢 Test DSL Motor Contable con Contexto");
    console.log("=======================================");
    
    let contexto = crearContextoSimple();
    let motor = MotorContable;
    
    console.log("Probando DSL con contexto...");
    let resultado = motor.use("asiento 1110 debe 15000 contrapartida 4110 haber 15000 por Venta de productos", contexto);
    
    console.log("Resultado:", resultado);
    console.log("✅ Test completado");
}