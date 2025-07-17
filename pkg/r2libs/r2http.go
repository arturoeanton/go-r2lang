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

// Definiremos una estructura para guardar las rutas
type r2Route struct {
	method  string
	pattern string
	//handlerName string
	handlerfx *r2core.UserFunction
	// Podrías también almacenar un compilado de regexp, si quisieras
}

// Guardamos las rutas en un slice (o array global) para simplificar
var r2Routes []r2Route

func httpHandler(args []interface{}) interface{} {
	if len(args) < 3 {
		panic("handler necesita 3 argumentos: (method, pattern, fx)")
	}
	method, ok1 := args[0].(string)
	pattern, ok2 := args[1].(string)
	handler, ok3 := args[2].(*r2core.UserFunction)
	if !ok1 || !ok2 {
		panic("handler: should be (method, pattern, fx)")
	}

	if !ok3 {
		panic("handler: fx should be a function")
	}

	// Agregamos la ruta a la tabla
	r2Routes = append(r2Routes, r2Route{
		method:    strings.ToUpper(method),
		pattern:   pattern,
		handlerfx: handler,
	})
	return nil
}

func RegisterHTTP(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		"handler": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			return httpHandler(args)
		}),

		"serve": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("serve needs 1 argument: (addr)")
			}
			addr, ok := args[0].(string)
			if !ok {
				panic("serve: argument should be a string")
			}

			// Definimos un único handler en Go para todas las rutas
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				// Leemos body (simple, para requests tipo POST/PUT)
				var bodyStr string
				if r.Method == "POST" || r.Method == "PUT" {
					data, err := io.ReadAll(r.Body)
					if err == nil {
						bodyStr = string(data)
					}
				}
				// Buscamos una ruta que coincida con r.Method y r.URL.Path
				route, pathVars := matchRoute(r2Routes, r.Method, r.URL.Path)
				if route == nil {
					// No match
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprintf(w, "404 Not Found\n")
					return
				}

				r2Handler := route.handlerfx

				// Llamamos la función con [pathVars, method, bodyStr]
				// pathVars lo podemos pasar como map[string]interface{}
				// En R2, el usuario lo recibe como un "obj", o un diccionario nativo:
				pathVarsMap := make(map[string]interface{})
				for k, v := range pathVars {
					pathVarsMap[k] = v
				}
				argsR2 := []interface{}{pathVarsMap, r.Method, bodyStr}
				respVal := r2Handler.Call(argsR2...)

				// Si la respuesta es string, la imprimimos, si no, la convertimos
				respStr, okResp := respVal.(string)
				if !okResp {
					respCustom, okResp := respVal.(*r2core.ObjectInstance)
					var data map[string]interface{}
					if okResp {
						data = respCustom.Env.GetStore()
					} else {
						data, ok = respVal.(map[string]interface{})
						if !ok {
							return
						}
					}
					header, ok := data["header"].(map[string]interface{})
					if ok {
						for k, v := range header {
							w.Header().Set(k, v.(string))
						}
					}
					status, ok := data["status"]
					if ok {
						w.WriteHeader(status.(int))
					}
					body, ok := data["body"]
					if ok {
						fmt.Fprint(w, body.(string))
						return
					}
					fmt.Fprintf(w, "%v", respVal)
					return
				}
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, respStr)
			})

			// Arrancamos el servidor (bloqueante)
			fmt.Println("Listening on ", addr)
			err := http.ListenAndServe(addr, nil)
			if err != nil {
				panic(fmt.Sprintf("serve: error in ListenAndServe: %v", err))
			}
			return nil
		}),

		"vars": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("vars necesita 2 argumentos: (map, key)")
			}
			// Verificamos tipos
			theMap, ok1 := args[0].(map[string]interface{})
			theKey, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("vars: primer argumento debe ser un map, segundo un string")
			}
			// Intentamos extraer
			value, found := theMap[theKey]
			if !found {
				// Retornamos nil si no está la clave
				return nil
			}
			return value
		}),

		"XML": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) != 2 {
				panic("XML needs at least 2 arguments")
			}
			root, ok := args[0].(string)
			if !ok {
				panic("XML: argument must be a string")
			}
			objectInstance, ok := args[1].(*r2core.ObjectInstance)
			var instance map[string]interface{}
			if !ok {
				instance, ok = args[1].(map[string]interface{})
				if !ok {
					panic("XML: argument must be an object or a map")
				}
			} else {
				instance = objectInstance.Env.GetStore()
			}

			instanceClean := removeBehavior(instance)
			// Convertimos a JSON
			data, err := mapToXML(root, instanceClean)
			if err != nil {
				panic(fmt.Sprintf("XML: error in Marshal: %v", err))
			}
			return string(data)
		}),

		"HttpResponse": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("HttpResponse needs at least 1 argument")
			}
			status, ok := args[0].(float64)
			posArgs := 1
			if !ok {
				status = 200
				posArgs = 0
			}
			statusInt := int(status)
			header := make(map[string]interface{})
			body := ""

			if len(args) > posArgs {
				header, ok = args[posArgs].(map[string]interface{})
				if !ok {
					value, ok := args[posArgs].(string)
					if ok {
						if len(args) == (posArgs + 1) {
							body = value
							header = make(map[string]interface{})
							header["Content-Type"] = DetectContentType(body).String()
						} else {
							header = make(map[string]interface{})
							header["Content-Type"] = value
						}
					}

				}
			}
			if len(args) == (posArgs + 2) {
				body, ok = args[posArgs+1].(string)
				if !ok {
					panic("HttpResponse: should be (status, header, body)")
				}
			}
			return map[string]interface{}{
				"status": statusInt,
				"header": header,
				"body":   body,
			}
		}),

		"JSON": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("JSON needs at least 1 argument")
			}
			var instance map[string]interface{}
			objectInstance, ok := args[0].(*r2core.ObjectInstance)
			if !ok {
				instance, ok = args[0].(map[string]interface{})
				if !ok {
					panic("JSON: argument must be an object or a map")
				}
			} else {
				instance = objectInstance.Env.GetStore()
			}
			instanceClean := removeBehavior(instance)
			// Convertimos a JSON
			data, err := json.Marshal(instanceClean)
			if err != nil {
				panic(fmt.Sprintf("JSON: error in Marshal: %v", err))
			}
			return string(data)
		}),
	}

	RegisterModule(env, "http", functions)
}

// mapToXML convierte un mapa a XML con un elemento raíz dinámico.
func mapToXML(root string, data map[string]interface{}) ([]byte, error) {
	// Crear un buffer para almacenar el XML
	var xmlData []byte
	// hacer el io writer
	xmlWriter := bytes.NewBuffer(xmlData)
	// Crear un encoder XML que escriba en el buffer
	encoder := xml.NewEncoder(xmlWriter)
	encoder.Indent("", "    ") // Formatear con indentación

	// Escribir el elemento raíz
	startElement := xml.StartElement{Name: xml.Name{Local: root}}
	if err := encoder.EncodeToken(startElement); err != nil {
		return nil, err
	}

	// Función recursiva para procesar el mapa
	err := encodeMap(encoder, data)
	if err != nil {
		return nil, err
	}

	// Cerrar el elemento raíz
	if err := encoder.EncodeToken(startElement.End()); err != nil {
		return nil, err
	}

	// Finalizar el encoder
	if err := encoder.Flush(); err != nil {
		return nil, err
	}
	xmlData = xmlWriter.Bytes()
	return xmlData, nil
}

// encodeMap procesa un mapa y escribe sus elementos como sub-elementos XML.
func encodeMap(encoder *xml.Encoder, data map[string]interface{}) error {
	for key, value := range data {
		// Crear el elemento XML para la clave actual
		startElement := xml.StartElement{Name: xml.Name{Local: key}}

		switch v := value.(type) {
		case map[string]interface{}:
			// Si el valor es otro mapa, anidar elementos
			if err := encoder.EncodeToken(startElement); err != nil {
				return err
			}
			if err := encodeMap(encoder, v); err != nil {
				return err
			}
			if err := encoder.EncodeToken(startElement.End()); err != nil {
				return err
			}
		case []interface{}:
			// Si el valor es una lista, iterar sobre ella
			for _, item := range v {
				if err := encoder.EncodeToken(startElement); err != nil {
					return err
				}
				// Asumimos que los elementos de la lista son básicos o mapas
				switch itemVal := item.(type) {
				case map[string]interface{}:
					if err := encodeMap(encoder, itemVal); err != nil {
						return err
					}
				default:
					// Convertir el valor a string
					if err := encoder.EncodeToken(xml.CharData([]byte(fmt.Sprintf("%v", itemVal)))); err != nil {
						return err
					}
				}
				if err := encoder.EncodeToken(startElement.End()); err != nil {
					return err
				}
			}
		default:
			// Para tipos básicos, escribir el contenido
			if err := encoder.EncodeElement(v, startElement); err != nil {
				return err
			}
		}
	}
	return nil
}

func removeBehavior(instance map[string]interface{}) map[string]interface{} {
	instanceOut := make(map[string]interface{})
	for k, v := range instance {
		if k == "self" || k == "this" {
			continue
		}
		if _, ok := v.(*r2core.UserFunction); ok {
			continue
		}
		if _, ok := v.(*r2core.BuiltinFunction); ok {
			continue
		}

		if subInstance, ok := v.(*r2core.ObjectInstance); ok {
			instanceOut[k] = removeBehavior(subInstance.Env.GetStore())
			continue
		}
		instanceOut[k] = v
	}
	return instanceOut

}

// matchRoute busca en r2Routes la primera que coincida con method y path
// y retorna (rutaEncontrada, mapDeVariables). Si no hay match, retorna (nil, nil).
func matchRoute(routes []r2Route, method, path string) (*r2Route, map[string]string) {
	for _, rt := range routes {
		if rt.method != strings.ToUpper(method) {
			continue
		}
		pathVars, ok := matchPattern(rt.pattern, path)
		if ok {
			return &rt, pathVars
		}
	}
	return nil, nil
}

// matchPattern recibe un pattern (ej "/users/:id") y un path ("/users/123").
// Retorna (mapDeVariables, true) si matchea, o (nil, false) si no.
func matchPattern(pattern, path string) (map[string]string, bool) {
	// Dividimos por "/", comparamos fragmentos. Los ":" se interpretan como variable.
	patParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patParts) != len(pathParts) {
		return nil, false
	}

	vars := make(map[string]string)
	for i := 0; i < len(patParts); i++ {
		p := patParts[i]
		real := pathParts[i]
		if strings.HasPrefix(p, ":") {
			// variable
			varName := p[1:] // sin ':'
			vars[varName] = real
		} else {
			// literal
			if p != real {
				return nil, false
			}
		}
	}
	return vars, true
}
