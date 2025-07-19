// pruebaHTTP.r2

func main() {
    // 1) GET simple
    let texto = request.get("https://httpbin.org/get");
    std.print("Texto de GET =>", texto);

    // 2) parse JSON
    let jresp = json.parse(texto);
    std.print("parseJSON =>", jresp);

    // 3) GET JSON directo
    let jauto = request.getJSON("https://httpbin.org/json");
    std.print("httpGetJSON =>", jauto);

    // 4) POST con body "Hola"
    let postResp = request.post("https://httpbin.org/post", "Hola desde R2");
    std.print("POST resp =>", postResp);

    // 5) POST JSON
    let data = {};
    data["nombre"] = "Alice";
    data["edad"] = 30;
    let postJsonResp = request.postJSON("https://httpbin.org/post", data);
    std.print("httpPostJSON =>", postJsonResp);

    // 6) Manejo de XML (didáctico)
    let xmlString = "<root><person name='Bob'><age>25</age></person></root>";
    let parsedXml = xml.parse(xmlString);
    std.print("parsedXml =>", parsedXml);

    // stringifyXML con un map estilo { "root": { ... } }
    let newXmlMap = {};
    newXmlMap["root"] = {};
    newXmlMap["root"]["person"] = {};
    newXmlMap["root"]["person"]["_attrs"] = {};
    newXmlMap["root"]["person"]["_attrs"]["name"] = "Carlos";
    newXmlMap["root"]["person"]["age"] = {};
    newXmlMap["root"]["person"]["age"]["_content"] = 40;


    let xmlOut = xml.stringify(newXmlMap);
    std.print("stringifyXML =>", xmlOut);

    std.print("Fin de pruebaHTTP.r2");
}