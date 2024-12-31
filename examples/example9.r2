// pruebaHTTP.r2

func main() {
    // 1) GET simple
    let texto = clientHttpGet("https://httpbin.org/get");
    print("Texto de GET =>", texto);

    // 2) parse JSON
    let jresp = parseJSON(texto);
    print("parseJSON =>", jresp);

    // 3) GET JSON directo
    let jauto = clientHttpGetJSON("https://httpbin.org/json");
    print("httpGetJSON =>", jauto);

    // 4) POST con body "Hola"
    let postResp = clientHttpPost("https://httpbin.org/post", "Hola desde R2");
    print("POST resp =>", postResp);

    // 5) POST JSON
    let data = {};
    varsSet(data, "nombre", "Alice");
    varsSet(data, "edad", 30);
    let postJsonResp = clientHttpPostJSON("https://httpbin.org/post", data);
    print("httpPostJSON =>", postJsonResp);

    // 6) Manejo de XML (didáctico)
    let xmlString = "<root><person name='Bob'><age>25</age></person></root>";
    let parsedXml = parseXML(xmlString);
    print("parsedXml =>", parsedXml);

    // stringifyXML con un map estilo { "root": { ... } }
    let newXmlMap = {};
    newXmlMap["root"] = {};
    newXmlMap["root"]["person"] = {};
    newXmlMap["root"]["person"]["_attrs"] = {};
    newXmlMap["root"]["person"]["_attrs"]["name"] = "Carlos";
    newXmlMap["root"]["person"]["age"] = {};
    newXmlMap["root"]["person"]["age"]["_content"] = 40;


    let xmlOut = stringifyXML(newXmlMap);
    print("stringifyXML =>", xmlOut);

    print("Fin de pruebaHTTP.r2");
}