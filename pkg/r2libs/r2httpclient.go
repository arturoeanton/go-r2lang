package r2libs

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// r2httpclient.go: Funciones nativas de HTTP y JSON/XML en R2

func RegisterHTTPClient(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		"clientHttpGet": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("clientHttpGet necesita (url)")
			}
			url, ok := args[0].(string)
			if !ok {
				panic("clientHttpGet: url debe ser string")
			}
			resp, err := http.Get(url)
			if err != nil {
				panic(fmt.Sprintf("clientHttpGet: error en GET '%s': %v", url, err))
			}
			defer resp.Body.Close()
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				panic(fmt.Sprintf("clientHttpGet: error al leer body: %v", err))
			}
			return string(data)
		}),

		"clientHttpPost": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("clientHttpPost necesita (url, bodyString)")
			}
			url, ok1 := args[0].(string)
			bodyStr, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("clientHttpPost: (url, bodyString) deben ser strings")
			}

			resp, err := http.Post(url, "text/plain", bytes.NewBufferString(bodyStr))
			if err != nil {
				panic(fmt.Sprintf("clientHttpPost: error en POST '%s': %v", url, err))
			}
			defer resp.Body.Close()
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				panic(fmt.Sprintf("clientHttpPost: error al leer body: %v", err))
			}
			return string(data)
		}),

		"parseJSON": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
		}),

		"stringifyJSON": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
		}),

		"clientHttpGetJSON": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				panic(fmt.Sprintf("httpGetJSON: error al leer body: %v", err))
			}
			var result interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				panic(fmt.Sprintf("httpGetJSON: error al parsear JSON: %v", err))
			}
			return result
		}),

		"clientHttpPostJSON": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
			respData, err := io.ReadAll(resp.Body)
			if err != nil {
				panic(fmt.Sprintf("httpPostJSON: error al leer respuesta: %v", err))
			}

			var result interface{}
			if err := json.Unmarshal(respData, &result); err != nil {
				panic(fmt.Sprintf("httpPostJSON: error al parsear respuesta JSON: %v", err))
			}
			return result
		}),

		"parseXML": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
		}),

		"stringifyXML": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
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
		}),
	}

	RegisterModule(env, "httpclient", functions)
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
