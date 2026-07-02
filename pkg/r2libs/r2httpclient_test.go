package r2libs

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func newHTTPClientModule(t *testing.T) map[string]interface{} {
	t.Helper()
	env := r2core.NewEnvironment()
	RegisterHTTPClient(env)
	mod, ok := env.Get("httpclient")
	if !ok {
		t.Fatal("httpclient module not registered")
	}
	module, ok := mod.(map[string]interface{})
	if !ok {
		t.Fatal("httpclient module has unexpected type")
	}
	return module
}

func httpclientFn(t *testing.T, module map[string]interface{}, name string) r2core.BuiltinFunction {
	t.Helper()
	fn, ok := module[name].(r2core.BuiltinFunction)
	if !ok {
		t.Fatalf("httpclient.%s not found or has wrong type", name)
	}
	return fn
}

// TestClientHttpGetEnforcesBodySizeLimit is a regression test: previously
// clientHttpGet/clientHttpPost/clientHttpGetJSON/clientHttpPostJSON read
// resp.Body via io.ReadAll with no cap at all, so a slow/malicious (or
// just unexpectedly huge) server response would be buffered into memory
// in full, an unbounded-memory DoS vector. Now responses over
// maxHTTPClientResponseBytes must be rejected with a controlled panic
// instead of being fully buffered.
func TestClientHttpGetEnforcesBodySizeLimit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		chunk := make([]byte, 1<<20) // 1MB
		for i := range chunk {
			chunk[i] = 'A'
		}
		flusher, _ := w.(http.Flusher)
		// stream one byte over the limit
		total := int64(0)
		for total < maxHTTPClientResponseBytes+(1<<20) {
			n, err := w.Write(chunk)
			if err != nil {
				return
			}
			total += int64(n)
			if flusher != nil {
				flusher.Flush()
			}
		}
	}))
	defer srv.Close()

	module := newHTTPClientModule(t)
	get := httpclientFn(t, module, "clientHttpGet")

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("BUG: expected clientHttpGet to panic on an oversized response, but it returned normally")
		}
		msg, ok := r.(string)
		if !ok || !strings.Contains(msg, "exceeds maximum allowed size") {
			t.Fatalf("expected a size-limit panic message, got: %v", r)
		}
	}()

	get(srv.URL)
}

// TestClientHttpGetNormalResponse is a sanity check that ordinary
// small responses still work after adding the body-size limit.
func TestClientHttpGetNormalResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello world")
	}))
	defer srv.Close()

	module := newHTTPClientModule(t)
	get := httpclientFn(t, module, "clientHttpGet")

	result := get(srv.URL)
	body, ok := result.(string)
	if !ok || body != "hello world" {
		t.Fatalf("expected 'hello world', got %v", result)
	}
}

// TestClientHttpGetJSONRoundTrip exercises clientHttpGetJSON /
// clientHttpPostJSON against a real httptest server.
func TestClientHttpGetJSONRoundTrip(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"hello":"world","n":42}`)
	}))
	defer srv.Close()

	module := newHTTPClientModule(t)
	getJSON := httpclientFn(t, module, "clientHttpGetJSON")

	result := getJSON(srv.URL)
	obj, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected object result, got %v", result)
	}
	if obj["hello"] != "world" || obj["n"] != 42.0 {
		t.Fatalf("unexpected JSON result: %v", obj)
	}
}

// TestClientHttpGetPropagatesConnectionError makes sure a connection
// failure (e.g. nothing listening on the port) surfaces as a controlled
// panic rather than a nil-pointer dereference or silent empty result.
func TestClientHttpGetPropagatesConnectionError(t *testing.T) {
	module := newHTTPClientModule(t)
	get := httpclientFn(t, module, "clientHttpGet")

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected clientHttpGet to panic on connection failure")
		}
	}()

	// Port 1 is reserved/unlikely to have a listener; expect a dial error.
	get("http://127.0.0.1:1/unreachable")
}

// TestStringifyXMLCircularReferenceDoesNotCrash is a regression test for
// an unrecoverable Go stack overflow: mapToXMLNode recursed over
// map[string]interface{} values without any cycle detection, so a
// self-referential R2Lang map (constructible via a["self"] = a) passed to
// stringifyXML crashed the whole process with a fatal error that
// panic/recover cannot catch. It must now panic with a normal, catchable
// error instead.
func TestStringifyXMLCircularReferenceDoesNotCrash(t *testing.T) {
	module := newHTTPClientModule(t)
	stringifyXML := httpclientFn(t, module, "stringifyXML")

	cyclic := map[string]interface{}{}
	cyclic["self"] = cyclic

	root := map[string]interface{}{"root": cyclic}

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("BUG: expected stringifyXML to panic on a circular map, but it returned normally")
		}
		msg, ok := r.(string)
		if !ok || !strings.Contains(msg, "circular reference detected") {
			t.Fatalf("expected a circular reference panic message, got: %v", r)
		}
	}()

	stringifyXML(root)
}

// TestStringifyXMLSharedNonCyclicReferenceIsNotFlagged makes sure the
// cycle-detection fix doesn't over-trigger on a diamond-shaped (shared but
// non-cyclic) reference: the same sub-map appearing twice as sibling
// values should still serialize fine.
func TestStringifyXMLSharedNonCyclicReferenceIsNotFlagged(t *testing.T) {
	module := newHTTPClientModule(t)
	stringifyXML := httpclientFn(t, module, "stringifyXML")

	shared := map[string]interface{}{"_content": "shared-value"}
	root := map[string]interface{}{
		"root": map[string]interface{}{
			"a": shared,
			"b": shared,
		},
	}

	result := stringifyXML(root)
	out, ok := result.(string)
	if !ok || !strings.Contains(out, "shared-value") {
		t.Fatalf("expected shared non-cyclic reference to serialize normally, got: %v", result)
	}
}

// TestStringifyXMLHandlesInterfaceSlice is a regression test for the bug
// class (recurring throughout this codebase) where a function only
// type-asserted []interface{} for arrays and silently mishandled the
// r2core.InterfaceSlice produced by .map()/.filter()/.sort()/etc.
// mapToXMLNode used to fall into its default branch for InterfaceSlice
// values, collapsing the whole array into a single stringified text node
// (e.g. "[2 4 6]") instead of emitting one child element per item.
func TestStringifyXMLHandlesInterfaceSlice(t *testing.T) {
	module := newHTTPClientModule(t)
	stringifyXML := httpclientFn(t, module, "stringifyXML")

	items := r2core.InterfaceSlice{2.0, 4.0, 6.0}
	root := map[string]interface{}{
		"root": map[string]interface{}{
			"item": items,
		},
	}

	result := stringifyXML(root).(string)
	if strings.Contains(result, "[2 4 6]") {
		t.Fatalf("BUG: InterfaceSlice was collapsed into a single text node: %s", result)
	}
	if strings.Count(result, "<item>") != 3 {
		t.Fatalf("expected 3 separate <item> child elements, got: %s", result)
	}
}
