package r2lang

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// r2httpclient.go: Funciones nativas de HTTP y JSON/XML en R2

func RegisterHTTPClient(env *Environment) {

	//========================================
	// 1) httpGet(url) => string (body)
	//========================================
	env.Set("httpGet", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("httpGet necesita (url)")
		}
		url, ok := args[0].(string)
		if !ok {
			panic("httpGet: url debe ser string")
		}
		resp, err := http.Get(url)
		if err != nil {
			panic(fmt.Sprintf("httpGet: error en GET '%s': %v", url, err))
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(fmt.Sprintf("httpGet: error al leer body: %v", err))
		}
		return string(data)
	}))

	//========================================
	// 2) httpPost(url, bodyString) => string (respuesta)
	//========================================
	env.Set("httpPost", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("httpPost necesita (url, bodyString)")
		}
		url, ok1 := args[0].(string)
		bodyStr, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			panic("httpPost: (url, bodyString) deben ser strings")
		}

		resp, err := http.Post(url, "text/plain", bytes.NewBufferString(bodyStr))
		if err != nil {
			panic(fmt.Sprintf("httpPost: error en POST '%s': %v", url, err))
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(fmt.Sprintf("httpPost: error al leer body: %v", err))
		}
		return string(data)
	}))

	//========================================
	// 3) parseJSON(jsonString) => map/array (R2)
	//    Convierte un string JSON a un map[string]interface{} o []interface{}
	//========================================
	env.Set("parseJSON", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("parseJSON necesita (jsonString)")
		}
		js, ok := args[0].(string)
		if !ok {
			panic("parseJSON: argumento debe ser string JSON")
		}
		var result interface{}
		err := json.Unmarshal([]byte(js), &result)
		if err != nil {
			panic(fmt.Sprintf("parseJSON: error al parsear JSON: %v", err))
		}
		// result puede ser map[string]interface{} o []interface{}
		return result
	}))

	//========================================
	// 4) stringifyJSON(value) => string
	//    Convierte un map/array nativo R2 a string JSON
	//========================================
	env.Set("stringifyJSON", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("stringifyJSON necesita (value)")
		}
		val := args[0]
		// val debería ser un map[string]interface{} o []interface{} (o algo anidable)
		data, err := json.Marshal(val)
		if err != nil {
			panic(fmt.Sprintf("stringifyJSON: error al serializar: %v", err))
		}
		return string(data)
	}))

	//========================================
	// 5) httpGetJSON(url) => map/array R2
	//    Hace GET y parsea JSON automáticamente
	//========================================
	env.Set("httpGetJSON", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("httpGetJSON necesita (url)")
		}
		url, ok := args[0].(string)
		if !ok {
			panic("httpGetJSON: url debe ser string")
		}
		resp, err := http.Get(url)
		if err != nil {
			panic(fmt.Sprintf("httpGetJSON: error en GET '%s': %v", url, err))
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(fmt.Sprintf("httpGetJSON: error al leer body: %v", err))
		}
		var result interface{}
		if err := json.Unmarshal(data, &result); err != nil {
			panic(fmt.Sprintf("httpGetJSON: error al parsear JSON: %v", err))
		}
		return result
	}))

	//========================================
	// 6) httpPostJSON(url, value) => map/array R2
	//    Serializa 'value' a JSON, POST, y parsea la respuesta como JSON
	//========================================
	env.Set("httpPostJSON", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("httpPostJSON necesita (url, value)")
		}
		url, ok1 := args[0].(string)
		val := args[1] // supuestamente map/array
		if !ok1 {
			panic("httpPostJSON: url debe ser string")
		}

		// Serializamos 'val' a JSON
		jsData, err := json.Marshal(val)
		if err != nil {
			panic(fmt.Sprintf("httpPostJSON: error al serializar: %v", err))
		}

		resp, err := http.Post(url, "application/json", bytes.NewReader(jsData))
		if err != nil {
			panic(fmt.Sprintf("httpPostJSON: error en POST '%s': %v", url, err))
		}
		defer resp.Body.Close()
		respData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(fmt.Sprintf("httpPostJSON: error al leer respuesta: %v", err))
		}

		var result interface{}
		if err := json.Unmarshal(respData, &result); err != nil {
			panic(fmt.Sprintf("httpPostJSON: error al parsear respuesta JSON: %v", err))
		}
		return result
	}))

	//========================================
	// 7) parseXML(xmlString) => map/array (muy simplificado)
	//========================================
	env.Set("parseXML", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("parseXML necesita (xmlString)")
		}
		s, ok := args[0].(string)
		if !ok {
			panic("parseXML: argumento debe ser string con XML")
		}
		var root xmlNode
		err := xml.Unmarshal([]byte(s), &root)
		if err != nil {
			panic(fmt.Sprintf("parseXML: error al parsear XML: %v", err))
		}
		// Convertimos xmlNode -> map[string]interface{}
		return xmlNodeToMap(&root)
	}))

	//========================================
	// 8) stringifyXML(value) => string
	//    Convierte un map[string]interface{} en un XML muy básico
	//========================================
	env.Set("stringifyXML", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("stringifyXML necesita (value)")
		}
		valMap, ok := args[0].(map[string]interface{})
		if !ok {
			panic("stringifyXML: se esperaba un map[string]interface{}")
		}
		// Podrías tener un "root" con subvalores, etc.
		// Aquí supongo que el map tiene 1 root key
		if len(valMap) != 1 {
			panic("stringifyXML: se espera un map con un root key")
		}

		// Tomamos la primera key como root
		var rootKey string
		var rootVal interface{}
		for k, v := range valMap {
			rootKey = k
			rootVal = v
			break
		}
		// Construimos un xmlNode
		node := mapToXMLNode(rootKey, rootVal)
		out, err := xml.MarshalIndent(node, "", "  ")
		if err != nil {
			panic(fmt.Sprintf("stringifyXML: error al serializar: %v", err))
		}
		// xml.Marshal produce <xmlNode> algo ... </xmlNode>
		// Si queremos algo más refinado, habría que personalizar.
		return string(out)
	}))
}

//===========================================================
// Estructuras de ayuda para parseXML / stringifyXML
//===========================================================

// xmlNode define un nodo genérico para decodificar XML
type xmlNode struct {
	XMLName  xml.Name   `xml:""`
	Content  string     `xml:",chardata"`
	Children []xmlNode  `xml:",any"`
	Attrs    []xml.Attr `xml:",any,attr"`
}

func xmlNodeToMap(n *xmlNode) interface{} {
	// Cada nodo se representa como un map con keys: "_name", "_content", "_attrs", "childName", ...
	m := make(map[string]interface{})
	// Nombre
	m["_name"] = n.XMLName.Local
	// Contenido textual (si hay)
	trimC := strings.TrimSpace(n.Content)
	if trimC != "" {
		m["_content"] = trimC
	}
	// Atributos
	if len(n.Attrs) > 0 {
		attrsMap := make(map[string]interface{})
		for _, at := range n.Attrs {
			attrsMap[at.Name.Local] = at.Value
		}
		m["_attrs"] = attrsMap
	}
	// Hijos
	for _, c := range n.Children {
		childName := c.XMLName.Local
		childVal := xmlNodeToMap(&c)
		// Podríamos agrupar por nombre, si hay múltiples hijos con el mismo
		// para simplicidad, se hace m[childName] = childVal si no existe,
		// o se convierte en array si hay más de 1
		existing, found := m[childName]
		if !found {
			m[childName] = childVal
		} else {
			// si ya hay uno => array
			switch arr := existing.(type) {
			case []interface{}:
				m[childName] = append(arr, childVal)
			default:
				// convertimos en array
				m[childName] = []interface{}{arr, childVal}
			}
		}
	}
	return m
}

// mapToXMLNode crea un xmlNode a partir de un map con keys
func mapToXMLNode(name string, val interface{}) xmlNode {
	node := xmlNode{XMLName: xml.Name{Local: name}}
	// val puede ser map[string]interface{} => interpretamos subnodos
	switch mm := val.(type) {
	case map[string]interface{}:
		// si hay "_content"
		if c, ok := mm["_content"]; ok {
			node.Content = fmt.Sprint(c)
		}
		// si hay "_attrs"
		if a, ok := mm["_attrs"]; ok {
			if attrMap, ok2 := a.(map[string]interface{}); ok2 {
				for k, v := range attrMap {
					node.Attrs = append(node.Attrs, xml.Attr{Name: xml.Name{Local: k}, Value: fmt.Sprint(v)})
				}
			}
		}
		// para cada key, si no es "_content", "_attrs", "_name", se interpretan como hijos
		for k, v := range mm {
			if k == "_content" || k == "_attrs" || k == "_name" {
				continue
			}
			// Si v es array => múltiples subnodos
			switch arr := v.(type) {
			case []interface{}:
				for _, elem := range arr {
					subNode := mapToXMLNode(k, elem)
					node.Children = append(node.Children, subNode)
				}
			default:
				// un solo nodo
				subNode := mapToXMLNode(k, arr)
				node.Children = append(node.Children, subNode)
			}
		}
	default:
		// si val es string / número => lo ponemos en _content
		node.Content = fmt.Sprint(val)
	}
	return node
}
