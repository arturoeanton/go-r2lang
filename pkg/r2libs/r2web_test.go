package r2libs

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func TestMatchWebPattern(t *testing.T) {
	tests := []struct {
		pattern    string
		path       string
		wantMatch  bool
		wantParams map[string]string
	}{
		{"/users/:id", "/users/123", true, map[string]string{"id": "123"}},
		{"/users/:id", "/orders/123", false, nil},
		{"/orders/:id", "/users/123", false, nil},
		{"/users/:id/posts/:postId", "/users/1/posts/2", true, map[string]string{"id": "1", "postId": "2"}},
		{"/users", "/users", true, map[string]string{}},
		{"/users", "/users/1", false, nil},
	}

	for _, tt := range tests {
		params, ok := matchWebPattern(tt.pattern, tt.path)
		if ok != tt.wantMatch {
			t.Errorf("matchWebPattern(%q, %q) match = %v, want %v", tt.pattern, tt.path, ok, tt.wantMatch)
			continue
		}
		if !ok {
			continue
		}
		if len(params) != len(tt.wantParams) {
			t.Errorf("matchWebPattern(%q, %q) params = %v, want %v", tt.pattern, tt.path, params, tt.wantParams)
			continue
		}
		for k, v := range tt.wantParams {
			if params[k] != v {
				t.Errorf("matchWebPattern(%q, %q) params[%q] = %q, want %q", tt.pattern, tt.path, k, params[k], v)
			}
		}
	}
}

func TestStringMapToInterfaceMap(t *testing.T) {
	out := stringMapToInterfaceMap(map[string]string{"a": "1", "b": "2"})
	if len(out) != 2 || out["a"] != "1" || out["b"] != "2" {
		t.Fatalf("unexpected conversion result: %v", out)
	}
	// The whole point of the conversion is that r2core's IndexExpression can
	// index the result; map[string]string cannot.
	env := r2core.NewEnvironment()
	env.Set("m", out)
	parser := r2core.NewParser(`m["a"]`)
	program := parser.ParseProgram()
	if got := program.Eval(env); got != "1" {
		t.Fatalf("indexing converted map failed: got %v", got)
	}
}

func evalUserFunction(t *testing.T, src string) *r2core.UserFunction {
	t.Helper()
	parser := r2core.NewParser(src)
	program := parser.ParseProgram()
	env := r2core.NewEnvironment()
	env.Set("true", true)
	env.Set("false", false)
	env.Set("nil", nil)
	RegisterWeb(env)
	result := program.Eval(env)
	fn, ok := result.(*r2core.UserFunction)
	if !ok {
		t.Fatalf("expected *r2core.UserFunction, got %T (%v)", result, result)
	}
	return fn
}

// TestRouteDispatcher_PathParamsAndMethods exercises the request pipeline
// end to end (route matching, GET/POST sharing a literal path, ctx.params /
// ctx.query / ctx.headers indexing) via a real HTTP round trip.
func TestRouteDispatcher_PathParamsAndMethods(t *testing.T) {
	app := newWebApp()

	getUsersByID := evalUserFunction(t, `let h = func(ctx) {
		return web.json({ id: ctx.params["id"], q: ctx.query["name"], h: ctx.headers["X-Test"] });
	}; h`)
	getUsers := evalUserFunction(t, `let h = func(ctx) { return "GET users"; }; h`)
	postUsers := evalUserFunction(t, `let h = func(ctx) { return "POST users"; }; h`)

	registerRouteForApp(app, "GET", "/users/:id", getUsersByID)
	registerRouteForApp(app, "GET", "/users", getUsers)
	registerRouteForApp(app, "POST", "/users", postUsers)

	app.mu.RLock()
	routesCopy := make(map[string]map[string]*r2core.UserFunction, len(app.routes))
	for method, routes := range app.routes {
		inner := make(map[string]*r2core.UserFunction, len(routes))
		for path, handler := range routes {
			inner[path] = handler
		}
		routesCopy[method] = inner
	}
	app.mu.RUnlock()

	srv := httptest.NewServer(createRouteDispatcher(app, routesCopy))
	defer srv.Close()

	// Path param + query + header indexing (previously panicked: ctx.params
	// etc. were raw Go map[string]string, unindexable from r2 scripts, and
	// http.ServeMux treated ":id" as a literal segment so this route never
	// matched "/users/123" at all).
	req, _ := http.NewRequest("GET", srv.URL+"/users/123?name=arturo", nil)
	req.Header.Set("X-Test", "hello")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("GET /users/123 status = %d body=%s", resp.StatusCode, body)
	}
	got := string(body)
	if !strings.Contains(got, `"id":"123"`) || !strings.Contains(got, `"q":"arturo"`) || !strings.Contains(got, `"h":"hello"`) {
		t.Fatalf("unexpected body: %s", got)
	}

	// GET and POST sharing the same literal path "/users" (previously
	// panicked inside webListenForApp because http.ServeMux does not allow
	// registering the same literal pattern twice).
	resp, err = http.Get(srv.URL + "/users")
	if err != nil {
		t.Fatal(err)
	}
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != 200 || string(body) != "GET users" {
		t.Fatalf("GET /users = %d %s", resp.StatusCode, body)
	}

	resp, err = http.Post(srv.URL+"/users", "text/plain", nil)
	if err != nil {
		t.Fatal(err)
	}
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != 200 || string(body) != "POST users" {
		t.Fatalf("POST /users = %d %s", resp.StatusCode, body)
	}

	// A method with no handler for a path that does exist under other
	// methods should be 405, not 404.
	req, _ = http.NewRequest("DELETE", srv.URL+"/users", nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("DELETE /users status = %d, want 405", resp.StatusCode)
	}

	// A literal segment that doesn't match a param pattern's sibling
	// literals must not be treated as a match.
	resp, err = http.Get(srv.URL + "/orders/123")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("GET /orders/123 status = %d, want 404", resp.StatusCode)
	}
}

// TestWebAppConcurrentRegistration hammers route/middleware/static
// registration concurrently with each other (run with -race).
func TestWebAppConcurrentRegistration(t *testing.T) {
	app := newWebApp()
	handler := evalUserFunction(t, `let h = func(ctx) { return "ok"; }; h`)

	var wg sync.WaitGroup
	deadline := time.Now().Add(200 * time.Millisecond)

	for i := 0; i < 4; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			n := 0
			for time.Now().Before(deadline) {
				registerRouteForApp(app, "GET", "/r", handler)
				app.mu.Lock()
				app.static["/s"] = "dir"
				app.middleware = append(app.middleware, i)
				app.mu.Unlock()
				n++
			}
		}()
	}

	wg.Wait()
}

// TestWebJSON_ContentType guards against a bug where web.json(data) returned
// a pre-marshaled JSON string instead of the typed {type:"json", data:...}
// response descriptor handleResponse expects — handleResponse's "case
// string:" branch always serves with Content-Type: text/html, so any
// handler using web.json(...) had its response mislabeled instead of being
// served as application/json.
func TestWebJSON_ContentType(t *testing.T) {
	app := newWebApp()

	handler := evalUserFunction(t, `let h = func(ctx) {
		return web.json({ ok: true });
	}; h`)
	registerRouteForApp(app, "GET", "/data", handler)

	app.mu.RLock()
	routesCopy := make(map[string]map[string]*r2core.UserFunction, len(app.routes))
	for method, routes := range app.routes {
		inner := make(map[string]*r2core.UserFunction, len(routes))
		for path, h := range routes {
			inner[path] = h
		}
		routesCopy[method] = inner
	}
	app.mu.RUnlock()

	srv := httptest.NewServer(createRouteDispatcher(app, routesCopy))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/data")
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("expected Content-Type application/json, got %q (body=%s)", ct, body)
	}
	if !strings.Contains(string(body), `"ok":true`) {
		t.Fatalf("unexpected body: %s", body)
	}
}
