package r2libs

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// WebApp represents a web application with modern routing
type WebApp struct {
	routes     map[string]map[string]*r2core.UserFunction
	middleware []interface{}
	templates  map[string]*template.Template
	static     map[string]string
	mu         sync.RWMutex
}

// WebContext represents the request context
type WebContext struct {
	Request  *http.Request
	Response http.ResponseWriter
	Params   map[string]string
	Query    map[string]string
	Body     string
	JSON     interface{}
	Headers  map[string]string
	Method   string
	Path     string
	app      *WebApp
}

var globalApp *WebApp

// RegisterWeb registers the modern web framework
func RegisterWeb(env *r2core.Environment) {
	globalApp = &WebApp{
		routes:    make(map[string]map[string]*r2core.UserFunction),
		templates: make(map[string]*template.Template),
		static:    make(map[string]string),
	}

	// Create web module
	webModule := map[string]interface{}{
		// App creation
		"createApp": func() map[string]interface{} {
			app := &WebApp{
				routes:    make(map[string]map[string]*r2core.UserFunction),
				templates: make(map[string]*template.Template),
				static:    make(map[string]string),
			}

			return map[string]interface{}{
				"get": func(path string, handler interface{}) {
					registerRouteForApp(app, "GET", path, handler)
				},
				"post": func(path string, handler interface{}) {
					registerRouteForApp(app, "POST", path, handler)
				},
				"put": func(path string, handler interface{}) {
					registerRouteForApp(app, "PUT", path, handler)
				},
				"delete": func(path string, handler interface{}) {
					registerRouteForApp(app, "DELETE", path, handler)
				},
				"static": func(path string, dir string) {
					app.mu.Lock()
					defer app.mu.Unlock()
					app.static[path] = dir
				},
				"listen": func(port string) {
					webListenForApp(app, port)
				},
				"use": func(middleware interface{}) {
					app.mu.Lock()
					defer app.mu.Unlock()
					app.middleware = append(app.middleware, middleware)
				},
			}
		},

		// Standalone functions
		"get":    webGet,
		"post":   webPost,
		"put":    webPut,
		"delete": webDelete,
		"static": webStatic,
		"listen": webListen,

		// Response helpers
		"json": func(data interface{}) string {
			b, _ := json.Marshal(data)
			return string(b)
		},
		"html": func(content string) map[string]interface{} {
			return map[string]interface{}{
				"type":    "html",
				"content": content,
			}
		},
		"redirect": func(url string) map[string]interface{} {
			return map[string]interface{}{
				"type": "redirect",
				"url":  url,
			}
		},
		"status": func(code int) map[string]interface{} {
			return map[string]interface{}{
				"type": "status",
				"code": code,
			}
		},

		// Form parsing
		"parseForm": parseFormData,
		"parseJSON": parseJSONData,
	}

	// Register in environment
	env.Set("web", webModule)
}

func webGet(path string, handler interface{}) {
	registerRoute("GET", path, handler)
}

func webPost(path string, handler interface{}) {
	registerRoute("POST", path, handler)
}

func webPut(path string, handler interface{}) {
	registerRoute("PUT", path, handler)
}

func webDelete(path string, handler interface{}) {
	registerRoute("DELETE", path, handler)
}

func webUse(middleware interface{}) {
	globalApp.mu.Lock()
	defer globalApp.mu.Unlock()
	globalApp.middleware = append(globalApp.middleware, middleware)
}

func webStatic(path string, dir string) {
	globalApp.mu.Lock()
	defer globalApp.mu.Unlock()
	globalApp.static[path] = dir
}

func registerRoute(method, path string, handler interface{}) {
	registerRouteForApp(globalApp, method, path, handler)
}

func registerRouteForApp(app *WebApp, method, path string, handler interface{}) {
	app.mu.Lock()
	defer app.mu.Unlock()

	if app.routes[method] == nil {
		app.routes[method] = make(map[string]*r2core.UserFunction)
	}

	switch h := handler.(type) {
	case *r2core.UserFunction:
		app.routes[method][path] = h
	default:
		panic(fmt.Sprintf("web: invalid handler type for %s %s", method, path))
	}
}

func webListen(port string) {
	webListenForApp(globalApp, port)
}

func webListenForApp(app *WebApp, port string) {
	mux := http.NewServeMux()

	// Setup routes
	for method, routes := range app.routes {
		for path, handler := range routes {
			mux.HandleFunc(path, createHTTPHandler(app, method, path, handler))
		}
	}

	// Setup static files
	for prefix, dir := range app.static {
		mux.Handle(prefix, http.StripPrefix(prefix, http.FileServer(http.Dir(dir))))
	}

	fmt.Printf("🚀 Web server listening on %s\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		panic(fmt.Sprintf("web: failed to start server: %v", err))
	}
}

func createHTTPHandler(app *WebApp, expectedMethod, pattern string, handler *r2core.UserFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != expectedMethod {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Create context
		ctx := createContext(w, r)

		// Extract path parameters
		pathParams := extractPathParams(pattern, r.URL.Path)
		for k, v := range pathParams {
			ctx.Params[k] = v
		}

		// Convert context to R2Lang object
		contextObj := map[string]interface{}{
			"params":  ctx.Params,
			"query":   ctx.Query,
			"body":    ctx.Body,
			"headers": ctx.Headers,
			"method":  ctx.Method,
			"path":    ctx.Path,
			"json": func() interface{} {
				if ctx.JSON != nil {
					return ctx.JSON
				}
				var data interface{}
				json.Unmarshal([]byte(ctx.Body), &data)
				return data
			},
			"send": func(content interface{}) interface{} {
				return handleResponse(w, content)
			},
			"status": func(code int) map[string]interface{} {
				return map[string]interface{}{
					"send": func(content interface{}) interface{} {
						w.WriteHeader(code)
						return handleResponse(w, content)
					},
				}
			},
			"redirect": func(url string) {
				http.Redirect(w, r, url, http.StatusFound)
			},
		}

		// Call handler with context
		result := handler.Call(contextObj)

		// Handle result if not already sent
		if result != nil {
			handleResponse(w, result)
		}
	}
}

func createContext(w http.ResponseWriter, r *http.Request) *WebContext {
	// Parse body
	var body string
	if r.Body != nil {
		bodyBytes := make([]byte, 0)
		buf := make([]byte, 1024)
		for {
			n, err := r.Body.Read(buf)
			if n > 0 {
				bodyBytes = append(bodyBytes, buf[:n]...)
			}
			if err != nil {
				break
			}
		}
		body = string(bodyBytes)
	}

	// Parse query parameters
	query := make(map[string]string)
	for k, v := range r.URL.Query() {
		if len(v) > 0 {
			query[k] = v[0]
		}
	}

	// Parse headers
	headers := make(map[string]string)
	for k, v := range r.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	return &WebContext{
		Request:  r,
		Response: w,
		Params:   make(map[string]string),
		Query:    query,
		Body:     body,
		Headers:  headers,
		Method:   r.Method,
		Path:     r.URL.Path,
		app:      globalApp,
	}
}

func handleResponse(w http.ResponseWriter, content interface{}) interface{} {
	switch v := content.(type) {
	case string:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(v))
	case map[string]interface{}:
		if responseType, ok := v["type"].(string); ok {
			switch responseType {
			case "html":
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.Write([]byte(fmt.Sprintf("%v", v["content"])))
			case "json":
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(v["data"])
			case "redirect":
				http.Redirect(w, nil, fmt.Sprintf("%v", v["url"]), http.StatusFound)
			case "status":
				if code, ok := v["code"].(int); ok {
					w.WriteHeader(code)
				}
			}
		} else {
			// Default to JSON
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(v)
		}
	default:
		// Try to marshal as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(v)
	}
	return nil
}

func extractPathParams(pattern, path string) map[string]string {
	params := make(map[string]string)
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patternParts) != len(pathParts) {
		return params
	}

	for i, part := range patternParts {
		if strings.HasPrefix(part, ":") {
			paramName := strings.TrimPrefix(part, ":")
			params[paramName] = pathParts[i]
		}
	}

	return params
}

func parseFormData(body string) map[string]string {
	result := make(map[string]string)
	if body == "" {
		return result
	}

	values, _ := url.ParseQuery(body)
	for k, v := range values {
		if len(v) > 0 {
			result[k] = v[0]
		}
	}
	return result
}

func parseJSONData(body string) interface{} {
	var result interface{}
	json.Unmarshal([]byte(body), &result)
	return result
}
