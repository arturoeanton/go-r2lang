// pruebaHTTP.r2

func main() {
    // 1) GET simple (request.get returns a response map with .text/.json/.status_code/etc.)
    let resp = request.get("https://httpbin.org/get");
    std.print("Texto de GET =>", resp.text);

    // 2) parse JSON (json.parse needs the raw body text, not the response map)
    let jresp = json.parse(resp.text);
    std.print("parseJSON =>", jresp);

    // 3) GET JSON directo (response map already exposes the parsed body as .json)
    let jauto = request.get("https://httpbin.org/json").json;
    std.print("httpGetJSON =>", jauto);

    // 4) POST con body "Hola"
    let postResp = request.post("https://httpbin.org/post", {"data": "Hola desde R2"});
    std.print("POST resp =>", postResp.text);

    // 5) POST JSON
    let data = {};
    data["nombre"] = "Alice";
    data["edad"] = 30;
    let postJsonResp = request.post("https://httpbin.org/post", {"json": data});
    std.print("httpPostJSON =>", postJsonResp.json);

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