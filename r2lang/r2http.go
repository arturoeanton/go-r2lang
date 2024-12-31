package r2lang

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Definiremos una estructura para guardar las rutas
type r2Route struct {
	method      string
	pattern     string
	handlerName string
	// Podrías también almacenar un compilado de regexp, si quisieras
}

// Guardamos las rutas en un slice (o array global) para simplificar
var r2Routes []r2Route

func RegisterHTTP(env *Environment) {
	// httpAddRoute(method, pattern, handlerName)
	env.Set("httpAddRoute", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 3 {
			panic("httpAddRoute necesita (method, pattern, handlerName)")
		}
		method, ok1 := args[0].(string)
		pattern, ok2 := args[1].(string)
		handler, ok3 := args[2].(string)
		if !ok1 || !ok2 || !ok3 {
			panic("httpAddRoute: todos los argumentos deben ser strings")
		}
		// Agregamos la ruta a la tabla
		r2Routes = append(r2Routes, r2Route{
			method:      strings.ToUpper(method),
			pattern:     pattern,
			handlerName: handler,
		})
		return nil
	}))

	// httpServe(addr)
	env.Set("httpServe", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("httpServe necesita 1 argumento: (addr)")
		}
		addr, ok := args[0].(string)
		if !ok {
			panic("httpServe: addr debe ser string")
		}

		// Definimos un único handler en Go para todas las rutas
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// Leemos body (simple, para requests tipo POST/PUT)
			var bodyStr string
			if r.Method == "POST" || r.Method == "PUT" {
				data, err := ioutil.ReadAll(r.Body)
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
			// Buscamos la función R2 en el env
			handlerVal, found := env.Get(route.handlerName)
			if !found {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Handler %s no existe en R2\n", route.handlerName)
				return
			}
			r2Handler, isFunc := handlerVal.(*UserFunction)
			if !isFunc {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Handler %s no es una función\n", route.handlerName)
				return
			}

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
				respStr = fmt.Sprintf("%v", respVal)
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, respStr)
		})

		// Arrancamos el servidor (bloqueante)
		fmt.Println("Escuchando en", addr)
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			panic(fmt.Sprintf("httpServe: error en ListenAndServe: %v", err))
		}
		return nil
	}))

	// vars(map, key) -> retorna map[key] o nil si no existe
	env.Set("vars", BuiltinFunction(func(args ...interface{}) interface{} {
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
	}))
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
