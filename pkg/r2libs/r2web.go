package r2libs

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

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

func newWebApp() *WebApp {
	return &WebApp{
		routes:    make(map[string]map[string]*r2core.UserFunction),
		templates: make(map[string]*template.Template),
		static:    make(map[string]string),
	}
}

// RegisterWeb registers the modern web framework
func RegisterWeb(env *r2core.Environment) {
	globalApp = newWebApp()

	// Create web module
	webModule := map[string]interface{}{
		// App creation
		"createApp": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			app := newWebApp()

			return map[string]interface{}{
				"get":    webRouteRegistrar(app, "GET"),
				"post":   webRouteRegistrar(app, "POST"),
				"put":    webRouteRegistrar(app, "PUT"),
				"delete": webRouteRegistrar(app, "DELETE"),
				"static": webStaticRegistrar(app),
				"listen": webListenRegistrar(app),
				"use":    webUseRegistrar(app),
			}
		}),

		// Standalone functions
		"get":    webRouteRegistrar(globalApp, "GET"),
		"post":   webRouteRegistrar(globalApp, "POST"),
		"put":    webRouteRegistrar(globalApp, "PUT"),
		"delete": webRouteRegistrar(globalApp, "DELETE"),
		"static": webStaticRegistrar(globalApp),
		"listen": webListenRegistrar(globalApp),

		// Response helpers
		"json": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("web: json() requires (data)")
			}
			b, err := json.Marshal(args[0])
			if err != nil {
				panic(fmt.Sprintf("web: json() failed to marshal value: %v", err))
			}
			return string(b)
		}),
		"html": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("web: html() requires (content)")
			}
			content, ok := args[0].(string)
			if !ok {
				panic(fmt.Sprintf("web: html() expected string for argument 1, got %T", args[0]))
			}
			return map[string]interface{}{
				"type":    "html",
				"content": content,
			}
		}),
		"redirect": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("web: redirect() requires (url)")
			}
			urlStr, ok := args[0].(string)
			if !ok {
				panic(fmt.Sprintf("web: redirect() expected string for argument 1, got %T", args[0]))
			}
			return map[string]interface{}{
				"type": "redirect",
				"url":  urlStr,
			}
		}),
		"status": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("web: status() requires (code)")
			}
			code, ok := args[0].(float64)
			if !ok {
				panic(fmt.Sprintf("web: status() expected number for argument 1, got %T", args[0]))
			}
			return map[string]interface{}{
				"type": "status",
				"code": int(code),
			}
		}),

		// Form parsing
		"parseForm": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("web: parseForm() requires (body)")
			}
			body, ok := args[0].(string)
			if !ok {
				panic(fmt.Sprintf("web: parseForm() expected string for argument 1, got %T", args[0]))
			}
			return parseFormData(body)
		}),
		"parseJSON": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("web: parseJSON() requires (body)")
			}
			body, ok := args[0].(string)
			if !ok {
				panic(fmt.Sprintf("web: parseJSON() expected string for argument 1, got %T", args[0]))
			}
			return parseJSONData(body)
		}),
	}

	// Register in environment
	env.Set("web", webModule)
}

func webRouteRegistrar(app *WebApp, method string) r2core.BuiltinFunction {
	return func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic(fmt.Sprintf("web: %s() requires (path, handler)", strings.ToLower(method)))
		}
		path, ok := args[0].(string)
		if !ok {
			panic(fmt.Sprintf("web: %s() expected string for argument 1 (path), got %T", strings.ToLower(method), args[0]))
		}
		registerRouteForApp(app, method, path, args[1])
		return nil
	}
}

func webStaticRegistrar(app *WebApp) r2core.BuiltinFunction {
	return func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("web: static() requires (path, dir)")
		}
		path, ok := args[0].(string)
		if !ok {
			panic(fmt.Sprintf("web: static() expected string for argument 1 (path), got %T", args[0]))
		}
		dir, ok := args[1].(string)
		if !ok {
			panic(fmt.Sprintf("web: static() expected string for argument 2 (dir), got %T", args[1]))
		}
		app.mu.Lock()
		defer app.mu.Unlock()
		app.static[path] = dir
		return nil
	}
}

func webListenRegistrar(app *WebApp) r2core.BuiltinFunction {
	return func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("web: listen() requires (port)")
		}
		port, ok := args[0].(string)
		if !ok {
			panic(fmt.Sprintf("web: listen() expected string for argument 1 (port), got %T", args[0]))
		}
		webListenForApp(app, port)
		return nil
	}
}

func webUseRegistrar(app *WebApp) r2core.BuiltinFunction {
	return func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("web: use() requires (middleware)")
		}
		app.mu.Lock()
		defer app.mu.Unlock()
		app.middleware = append(app.middleware, args[0])
		return nil
	}
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

func webListenForApp(app *WebApp, port string) {
	mux := http.NewServeMux()

	// Snapshot routes/static under the read lock so concurrently-registered
	// routes (e.g. from a script using "go") can't race with this read.
	app.mu.RLock()
	routesCopy := make(map[string]map[string]*r2core.UserFunction, len(app.routes))
	for method, routes := range app.routes {
		inner := make(map[string]*r2core.UserFunction, len(routes))
		for path, handler := range routes {
			inner[path] = handler
		}
		routesCopy[method] = inner
	}
	staticCopy := make(map[string]string, len(app.static))
	for prefix, dir := range app.static {
		staticCopy[prefix] = dir
	}
	app.mu.RUnlock()

	// Setup static files first so their specific prefixes take priority over
	// the catch-all route dispatcher registered below.
	for prefix, dir := range staticCopy {
		mux.Handle(prefix, http.StripPrefix(prefix, http.FileServer(http.Dir(dir))))
	}

	// Routes (including ":param" patterns) are matched by a single dispatcher
	// instead of being registered individually on the mux: http.ServeMux
	// treats ":id" as a literal path segment (no wildcard support), and it
	// panics if two methods register the same literal pattern (e.g. GET
	// "/users" and POST "/users").
	mux.HandleFunc("/", createRouteDispatcher(app, routesCopy))

	fmt.Printf("🚀 Web server listening on %s\n", port)
	srv := &http.Server{
		Addr:              port,
		Handler:           mux,
		ReadHeaderTimeout: 15 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		panic(fmt.Sprintf("web: failed to start server: %v", err))
	}
}

// createRouteDispatcher returns a single handler that matches an incoming
// request against every registered route (across all methods) since
// http.ServeMux cannot host ":param" patterns or per-method routes that
// share the same literal path.
func createRouteDispatcher(app *WebApp, routesCopy map[string]map[string]*r2core.UserFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if routes, ok := routesCopy[r.Method]; ok {
			if handler, params, matched := matchRouteSet(routes, r.URL.Path); matched {
				serveRoute(app, handler, params, w, r)
				return
			}
		}

		// The path matches a route registered under a different method:
		// report 405 instead of a plain 404.
		for method, routes := range routesCopy {
			if method == r.Method {
				continue
			}
			if _, _, matched := matchRouteSet(routes, r.URL.Path); matched {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
		}

		http.NotFound(w, r)
	}
}

// matchRouteSet finds the handler registered for path within routes, trying
// an exact match before falling back to ":param" pattern matching.
func matchRouteSet(routes map[string]*r2core.UserFunction, path string) (*r2core.UserFunction, map[string]string, bool) {
	if handler, ok := routes[path]; ok {
		return handler, nil, true
	}
	for pattern, handler := range routes {
		if params, ok := matchWebPattern(pattern, path); ok {
			return handler, params, true
		}
	}
	return nil, nil, false
}

func serveRoute(app *WebApp, handler *r2core.UserFunction, pathParams map[string]string, w http.ResponseWriter, r *http.Request) {
	// Create context
	ctx := createContext(app, w, r)

	for k, v := range pathParams {
		ctx.Params[k] = v
	}

	// Convert context to R2Lang object
	contextObj := map[string]interface{}{
		"params":  stringMapToInterfaceMap(ctx.Params),
		"query":   stringMapToInterfaceMap(ctx.Query),
		"body":    ctx.Body,
		"headers": stringMapToInterfaceMap(ctx.Headers),
		"method":  ctx.Method,
		"path":    ctx.Path,
		"json": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if ctx.JSON != nil {
				return ctx.JSON
			}
			var data interface{}
			json.Unmarshal([]byte(ctx.Body), &data)
			return data
		}),
		"send": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("web: ctx.send() requires (content)")
			}
			return handleResponse(w, r, args[0])
		}),
		"status": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("web: ctx.status() requires (code)")
			}
			code, ok := args[0].(float64)
			if !ok {
				panic(fmt.Sprintf("web: ctx.status() expected number for argument 1, got %T", args[0]))
			}
			return map[string]interface{}{
				"send": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
					if len(args) < 1 {
						panic("web: ctx.status().send() requires (content)")
					}
					w.WriteHeader(int(code))
					return handleResponse(w, r, args[0])
				}),
			}
		}),
		"redirect": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("web: ctx.redirect() requires (url)")
			}
			urlStr, ok := args[0].(string)
			if !ok {
				panic(fmt.Sprintf("web: ctx.redirect() expected string for argument 1, got %T", args[0]))
			}
			http.Redirect(w, r, urlStr, http.StatusFound)
			return nil
		}),
	}

	// Call handler with context
	result := handler.Call(contextObj)

	// Handle result if not already sent
	if result != nil {
		handleResponse(w, r, result)
	}
}

func createContext(app *WebApp, w http.ResponseWriter, r *http.Request) *WebContext {
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
		app:      app,
	}
}

func handleResponse(w http.ResponseWriter, r *http.Request, content interface{}) interface{} {
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
				if err := json.NewEncoder(w).Encode(v["data"]); err != nil {
					panic(fmt.Sprintf("web: failed to encode JSON response: %v", err))
				}
			case "redirect":
				http.Redirect(w, r, fmt.Sprintf("%v", v["url"]), http.StatusFound)
			case "status":
				if code, ok := v["code"].(int); ok {
					w.WriteHeader(code)
				}
			}
		} else {
			// Default to JSON
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(v); err != nil {
				panic(fmt.Sprintf("web: failed to encode JSON response: %v", err))
			}
		}
	default:
		// Try to marshal as JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(v); err != nil {
			panic(fmt.Sprintf("web: failed to encode JSON response: %v", err))
		}
	}
	return nil
}

// matchWebPattern matches a route pattern (e.g. "/users/:id") against a
// request path, requiring every non-":name" segment to match literally.
func matchWebPattern(pattern, path string) (map[string]string, bool) {
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patternParts) != len(pathParts) {
		return nil, false
	}

	params := make(map[string]string)
	for i, part := range patternParts {
		if strings.HasPrefix(part, ":") {
			paramName := strings.TrimPrefix(part, ":")
			params[paramName] = pathParts[i]
			continue
		}
		if part != pathParts[i] {
			return nil, false
		}
	}

	return params, true
}

func stringMapToInterfaceMap(m map[string]string) map[string]interface{} {
	out := make(map[string]interface{}, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
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
