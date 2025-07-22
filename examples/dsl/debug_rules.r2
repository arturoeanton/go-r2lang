// Debug de reglas DSL

dsl MotorTest {
    token("TEMPLATE_ID", "TPL[0-9]{3}")
    token("TEMPLATE", "template")
    token("CON", "con")
    token("IMPORTE", "[0-9]+")
    
    rule("test_template", ["TEMPLATE", "TEMPLATE_ID", "CON", "IMPORTE"], "testTemplate")
    
    func testTemplate(template, templateId, con, importe) {
        console.log("✅ Template rule funcionando:");
        console.log("   Template ID: " + templateId);
        console.log("   Importe: " + importe);
        return "OK";
    }
}

func main() {
    console.log("🔍 Debug de reglas DSL");
    console.log("======================");
    
    let motor = MotorTest;
    
    console.log("Probando: 'template TPL002 con 8500'");
    let resultado = motor.use("template TPL002 con 8500");
    console.log("Resultado:", resultado);
}