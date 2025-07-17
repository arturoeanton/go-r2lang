package r2libs

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// Test server setup
func setupTestServer() *httptest.Server {
	mux := http.NewServeMux()

	// GET endpoint
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Hello World",
			"method":  r.Method,
		})
	})

	// POST endpoint
	mux.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "POST received",
			"method":  r.Method,
		})
	})

	// PUT endpoint
	mux.HandleFunc("/put", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "PUT received",
			"method":  r.Method,
		})
	})

	// DELETE endpoint
	mux.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "DELETE received",
			"method":  r.Method,
		})
	})

	// JSON endpoint for testing JSON responses
	mux.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"name":  "Test User",
			"age":   30,
			"email": "test@example.com",
		})
	})

	// Status code endpoint
	mux.HandleFunc("/status/404", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})

	// Echo endpoint that returns request data
	mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Read body
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)

		response := map[string]interface{}{
			"method":  r.Method,
			"headers": r.Header,
			"body":    string(body),
			"url":     r.URL.String(),
		}

		json.NewEncoder(w).Encode(response)
	})

	return httptest.NewServer(mux)
}

func TestGlobalGet(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	env := r2core.NewEnvironment()
	RegisterRequests(env)

	// Test GET request
	getFuncVal, ok := env.Get("get")
	if !ok {
		t.Fatal("get function not found")
	}
	getFunc := getFuncVal.(r2core.BuiltinFunction)
	result := getFunc(server.URL + "/test")

	response := result.(map[string]interface{})
	if response["status_code"] != 200 {
		t.Errorf("Expected status code 200, got %v", response["status_code"])
	}

	if !response["ok"].(bool) {
		t.Error("Expected ok to be true")
	}

	// Check JSON response
	if response["json"] != nil {
		jsonData := response["json"].(map[string]interface{})
		if jsonData["message"] != "Hello World" {
			t.Errorf("Expected message 'Hello World', got %v", jsonData["message"])
		}
	}
}

func TestGlobalPost(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	env := r2core.NewEnvironment()
	RegisterRequests(env)

	// Test POST request
	postFuncVal, ok := env.Get("post")
	if !ok {
		t.Fatal("post function not found")
	}
	postFunc := postFuncVal.(r2core.BuiltinFunction)
	result := postFunc(server.URL + "/post")

	response := result.(map[string]interface{})
	if response["status_code"] != 200 {
		t.Errorf("Expected status code 200, got %v", response["status_code"])
	}

	if response["json"] != nil {
		jsonData := response["json"].(map[string]interface{})
		if jsonData["method"] != "POST" {
			t.Errorf("Expected method 'POST', got %v", jsonData["method"])
		}
	}
}

func TestGlobalPut(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	env := r2core.NewEnvironment()
	RegisterRequests(env)

	// Test PUT request
	putFuncVal, ok := env.Get("put")
	if !ok {
		t.Fatal("put function not found")
	}
	putFunc := putFuncVal.(r2core.BuiltinFunction)
	result := putFunc(server.URL + "/put")

	response := result.(map[string]interface{})
	if response["status_code"] != 200 {
		t.Errorf("Expected status code 200, got %v", response["status_code"])
	}

	if response["json"] != nil {
		jsonData := response["json"].(map[string]interface{})
		if jsonData["method"] != "PUT" {
			t.Errorf("Expected method 'PUT', got %v", jsonData["method"])
		}
	}
}

func TestGlobalDelete(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	env := r2core.NewEnvironment()
	RegisterRequests(env)

	// Test DELETE request
	deleteFuncVal, ok := env.Get("delete")
	if !ok {
		t.Fatal("delete function not found")
	}
	deleteFunc := deleteFuncVal.(r2core.BuiltinFunction)
	result := deleteFunc(server.URL + "/delete")

	response := result.(map[string]interface{})
	if response["status_code"] != 200 {
		t.Errorf("Expected status code 200, got %v", response["status_code"])
	}

	if response["json"] != nil {
		jsonData := response["json"].(map[string]interface{})
		if jsonData["method"] != "DELETE" {
			t.Errorf("Expected method 'DELETE', got %v", jsonData["method"])
		}
	}
}

func TestRequestWithParameters(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	env := r2core.NewEnvironment()
	RegisterRequests(env)

	// Test GET request with parameters
	getFuncVal, ok := env.Get("get")
	if !ok {
		t.Fatal("get function not found")
	}
	getFunc := getFuncVal.(r2core.BuiltinFunction)
	params := map[string]interface{}{
		"params": map[string]interface{}{
			"q":    "test",
			"page": 1,
		},
		"headers": map[string]interface{}{
			"User-Agent": "R2Lang-Test",
		},
	}

	result := getFunc(server.URL+"/echo", params)
	response := result.(map[string]interface{})

	if response["status_code"] != 200 {
		t.Errorf("Expected status code 200, got %v", response["status_code"])
	}

	// Check that parameters were included in URL
	url := response["url"].(string)
	if !strings.Contains(url, "q=test") {
		t.Error("Expected URL parameters to be included")
	}
}

func TestRequestWithJSON(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	env := r2core.NewEnvironment()
	RegisterRequests(env)

	// Test POST request with JSON data
	postFuncVal, ok := env.Get("post")
	if !ok {
		t.Fatal("post function not found")
	}
	postFunc := postFuncVal.(r2core.BuiltinFunction)
	jsonData := map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
	}

	params := map[string]interface{}{
		"json": jsonData,
	}

	result := postFunc(server.URL+"/echo", params)
	response := result.(map[string]interface{})

	if response["status_code"] != 200 {
		t.Errorf("Expected status code 200, got %v", response["status_code"])
	}

	// Check that JSON was properly sent
	if response["json"] != nil {
		echoResponse := response["json"].(map[string]interface{})
		if echoResponse["method"] != "POST" {
			t.Error("Expected POST method in echo response")
		}
	}
}

func TestSession(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	env := r2core.NewEnvironment()
	RegisterRequests(env)

	// Create session
	sessionFuncVal, ok := env.Get("session")
	if !ok {
		t.Fatal("session function not found")
	}
	sessionFunc := sessionFuncVal.(r2core.BuiltinFunction)
	sessionObj := sessionFunc()

	session := sessionObj.(map[string]interface{})

	// Test session GET
	getFunc := session["get"].(r2core.BuiltinFunction)
	result := getFunc(server.URL + "/test")

	response := result.(map[string]interface{})
	if response["status_code"] != 200 {
		t.Errorf("Expected status code 200, got %v", response["status_code"])
	}

	// Test session POST
	postFunc := session["post"].(r2core.BuiltinFunction)
	result = postFunc(server.URL + "/post")

	response = result.(map[string]interface{})
	if response["status_code"] != 200 {
		t.Errorf("Expected status code 200, got %v", response["status_code"])
	}
}

func TestErrorHandling(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	env := r2core.NewEnvironment()
	RegisterRequests(env)

	// Test 404 error
	getFuncVal, ok := env.Get("get")
	if !ok {
		t.Fatal("get function not found")
	}
	getFunc := getFuncVal.(r2core.BuiltinFunction)
	result := getFunc(server.URL + "/status/404")

	response := result.(map[string]interface{})
	if response["status_code"] != 404 {
		t.Errorf("Expected status code 404, got %v", response["status_code"])
	}

	if response["ok"].(bool) {
		t.Error("Expected ok to be false for 404")
	}

	// Test raise_for_status
	raiseFunc := response["raise_for_status"].(r2core.BuiltinFunction)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected raise_for_status to panic for 404")
		}
	}()
	raiseFunc()
}

func TestURLEncoding(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterRequests(env)

	// Test URL encoding
	encodeFuncVal, ok := env.Get("urlencode")
	if !ok {
		t.Fatal("urlencode function not found")
	}
	encodeFunc := encodeFuncVal.(r2core.BuiltinFunction)
	result := encodeFunc("hello world")

	if result != "hello+world" {
		t.Errorf("Expected 'hello+world', got %v", result)
	}

	// Test URL decoding
	decodeFuncVal, ok := env.Get("urldecode")
	if !ok {
		t.Fatal("urldecode function not found")
	}
	decodeFunc := decodeFuncVal.(r2core.BuiltinFunction)
	result = decodeFunc("hello+world")

	if result != "hello world" {
		t.Errorf("Expected 'hello world', got %v", result)
	}
}

func TestResponseFields(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	env := r2core.NewEnvironment()
	RegisterRequests(env)

	// Test response fields
	getFuncVal, ok := env.Get("get")
	if !ok {
		t.Fatal("get function not found")
	}
	getFunc := getFuncVal.(r2core.BuiltinFunction)
	result := getFunc(server.URL + "/json")

	response := result.(map[string]interface{})

	// Check all expected fields are present
	requiredFields := []string{"url", "status_code", "headers", "text", "json", "content", "ok", "elapsed"}
	for _, field := range requiredFields {
		if _, exists := response[field]; !exists {
			t.Errorf("Expected field %s to exist in response", field)
		}
	}

	// Check json() method
	if jsonFunc, exists := response["json_func"]; exists {
		jsonMethod := jsonFunc.(r2core.BuiltinFunction)
		jsonResult := jsonMethod()

		if jsonResult != nil {
			jsonData := jsonResult.(map[string]interface{})
			if jsonData["name"] != "Test User" {
				t.Errorf("Expected name 'Test User', got %v", jsonData["name"])
			}
		}
	}

	// Check elapsed time is reasonable
	elapsed := response["elapsed"].(float64)
	if elapsed < 0 || elapsed > 10 {
		t.Errorf("Expected elapsed time between 0 and 10 seconds, got %f", elapsed)
	}
}

func TestTimeout(t *testing.T) {
	// Create a slow server
	slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer slowServer.Close()

	env := r2core.NewEnvironment()
	RegisterRequests(env)

	// Test timeout
	getFuncVal, ok := env.Get("get")
	if !ok {
		t.Fatal("get function not found")
	}
	getFunc := getFuncVal.(r2core.BuiltinFunction)
	params := map[string]interface{}{
		"timeout": 1.0, // 1 second timeout
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected timeout to cause panic")
		}
	}()

	getFunc(slowServer.URL, params)
}

func TestAuthenticatedRequest(t *testing.T) {
	// Create server that requires authentication
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || username != "user" || password != "pass" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Authenticated",
		})
	}))
	defer authServer.Close()

	env := r2core.NewEnvironment()
	RegisterRequests(env)

	// Test authenticated request
	getFuncVal, ok := env.Get("get")
	if !ok {
		t.Fatal("get function not found")
	}
	getFunc := getFuncVal.(r2core.BuiltinFunction)
	params := map[string]interface{}{
		"auth": []interface{}{"user", "pass"},
	}

	result := getFunc(authServer.URL, params)
	response := result.(map[string]interface{})

	if response["status_code"] != 200 {
		t.Errorf("Expected status code 200, got %v", response["status_code"])
	}

	if response["json"] != nil {
		jsonData := response["json"].(map[string]interface{})
		if jsonData["message"] != "Authenticated" {
			t.Errorf("Expected message 'Authenticated', got %v", jsonData["message"])
		}
	}
}
