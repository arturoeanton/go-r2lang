package r2libs

import (
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
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

	// Upload endpoint for testing file uploads
	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			http.Error(w, "Failed to parse content type", http.StatusBadRequest)
			return
		}

		if mediaType != "multipart/form-data" {
			http.Error(w, "Expected multipart/form-data", http.StatusBadRequest)
			return
		}

		err = r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
			return
		}

		// Check for file
		file, handler, err := r.FormFile("file_field")
		if err != nil {
			http.Error(w, "Failed to get file from form", http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to read file content", http.StatusInternalServerError)
			return
		}

		// Check for other form data
		dataField := r.FormValue("data_field")

		response := map[string]interface{}{
			"message":      "Upload successful",
			"filename":     handler.Filename,
			"file_content": string(fileBytes),
			"data_field":   dataField,
			"content_type": handler.Header.Get("Content-Type"),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// HEAD endpoint
	mux.HandleFunc("/head", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "HEAD" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	// OPTIONS endpoint
	mux.HandleFunc("/options", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "OPTIONS" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Allow", "GET, POST, PUT, DELETE, HEAD, OPTIONS")
		w.WriteHeader(http.StatusOK)
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

// Test for HEAD and OPTIONS methods
func TestHeadOptions(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	env := r2core.NewEnvironment()
	RegisterRequests(env)

	// Test HEAD request
	headFuncVal, ok := env.Get("head")
	if !ok {
		t.Fatal("head function not found")
	}
	headFunc := headFuncVal.(r2core.BuiltinFunction)
	result := headFunc(server.URL + "/head")

	response := result.(map[string]interface{})
	if response["status_code"] != 200 {
		t.Errorf("Expected status code 200 for HEAD, got %v", response["status_code"])
	}
	if len(response["text"].(string)) > 0 {
		t.Error("Expected empty body for HEAD request")
	}

	// Test OPTIONS request
	optionsFuncVal, ok := env.Get("options")
	if !ok {
		t.Fatal("options function not found")
	}
	optionsFunc := optionsFuncVal.(r2core.BuiltinFunction)
	result = optionsFunc(server.URL + "/options")

	response = result.(map[string]interface{})
	if response["status_code"] != 200 {
		t.Errorf("Expected status code 200 for OPTIONS, got %v", response["status_code"])
	}
}

func TestFileUpload(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	env := r2core.NewEnvironment()
	RegisterRequests(env)

	// Create a temporary file to upload
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	fileContent := "Hello, this is a test file."
	_, err = tmpFile.WriteString(fileContent)
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Test file upload
	postFuncVal, ok := env.Get("post")
	if !ok {
		t.Fatal("post function not found")
	}
	postFunc := postFuncVal.(r2core.BuiltinFunction)

	params := map[string]interface{}{
		"files": map[string]interface{}{
			"file_field": tmpFile.Name(),
		},
		"data": map[string]interface{}{
			"data_field": "some_value",
		},
	}

	result := postFunc(server.URL+"/upload", params)
	response := result.(map[string]interface{})

	if response["status_code"] != 200 {
		t.Errorf("Expected status code 200, got %v", response["status_code"])
	}

	if response["json"] != nil {
		jsonData := response["json"].(map[string]interface{})
		if jsonData["message"] != "Upload successful" {
			t.Errorf("Expected message 'Upload successful', got %v", jsonData["message"])
		}
		if jsonData["file_content"] != fileContent {
			t.Errorf("Expected file content '%s', got '%s'", fileContent, jsonData["file_content"])
		}
		if jsonData["data_field"] != "some_value" {
			t.Errorf("Expected data_field 'some_value', got '%s'", jsonData["data_field"])
		}
	}
}

func TestProxy(t *testing.T) {
	// Create a proxy server
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simply forward the request
		resp, err := http.DefaultTransport.RoundTrip(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()
		for k, vv := range resp.Header {
			for _, v := range vv {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}))
	defer proxy.Close()

	// Create a target server
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Hello from target server",
		})
	}))
	defer target.Close()

	env := r2core.NewEnvironment()
	RegisterRequests(env)

	// Test request through proxy
	getFuncVal, ok := env.Get("get")
	if !ok {
		t.Fatal("get function not found")
	}
	getFunc := getFuncVal.(r2core.BuiltinFunction)
	params := map[string]interface{}{
		"proxies": map[string]interface{}{
			"http": proxy.URL,
		},
	}

	result := getFunc(target.URL, params)
	response := result.(map[string]interface{})

	if response["status_code"] != 200 {
		t.Errorf("Expected status code 200, got %v", response["status_code"])
	}

	if response["json"] != nil {
		jsonData := response["json"].(map[string]interface{})
		if jsonData["message"] != "Hello from target server" {
			t.Errorf("Expected message from target server, got %v", jsonData["message"])
		}
	}
}

func TestRetries(t *testing.T) {
	var attempts int32
	// Create a server that fails the first 2 times
	flakyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&attempts, 1) <= 2 {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Success on third attempt",
		})
	}))
	defer flakyServer.Close()

	env := r2core.NewEnvironment()
	RegisterRequests(env)

	// Test request with retries
	getFuncVal, ok := env.Get("get")
	if !ok {
		t.Fatal("get function not found")
	}
	getFunc := getFuncVal.(r2core.BuiltinFunction)
	params := map[string]interface{}{
		"retries": map[string]interface{}{
			"max":   3.0,
			"delay": 0.1,
		},
	}

	result := getFunc(flakyServer.URL, params)
	response := result.(map[string]interface{})

	if response["status_code"] != 200 {
		t.Errorf("Expected status code 200, got %v", response["status_code"])
	}

	if response["json"] != nil {
		jsonData := response["json"].(map[string]interface{})
		if jsonData["message"] != "Success on third attempt" {
			t.Errorf("Expected success message, got %v", jsonData["message"])
		}
	}

	if atomic.LoadInt32(&attempts) != 3 {
		t.Errorf("Expected 3 attempts, got %d", atomic.LoadInt32(&attempts))
	}
}

func TestCookieJarAutomatic(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:  "session_id",
			Value: "test_session_123",
			Path:  "/",
		})
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Logged in"))
	})

	mux.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil || cookie.Value != "test_session_123" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Access Denied"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Access Granted"))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	env := r2core.NewEnvironment()
	RegisterRequests(env)

	// Test with global functions (should now handle cookies automatically)
	t.Run("GlobalFunctionsCookieHandling", func(t *testing.T) {
		// First request to login endpoint
		loginResult := globalGet(server.URL + "/login")
		loginResponse := loginResult.(map[string]interface{})
		if loginResponse["status_code"] != 200 {
			t.Fatalf("Expected login status 200, got %v", loginResponse["status_code"])
		}

		// Second request to protected endpoint
		protectedResult := globalGet(server.URL + "/protected")
		protectedResponse := protectedResult.(map[string]interface{})
		if protectedResponse["status_code"] != 200 {
			t.Errorf("Expected protected status 200, got %v", protectedResponse["status_code"])
		}
		if protectedResponse["text"] != "Access Granted" {
			t.Errorf("Expected 'Access Granted', got '%s'", protectedResponse["text"])
		}
	})

	// Test with session functions (should already handle cookies)
	t.Run("SessionFunctionsCookieHandling", func(t *testing.T) {
		sessionFuncVal, ok := env.Get("session")
		if !ok {
			t.Fatal("session function not found")
		}
		sessionFunc := sessionFuncVal.(r2core.BuiltinFunction)
		sessionObj := sessionFunc()
		session := sessionObj.(map[string]interface{})

		sessionGetFunc := session["get"].(r2core.BuiltinFunction)

		// First request to login endpoint
		loginResult := sessionGetFunc(server.URL + "/login")
		loginResponse := loginResult.(map[string]interface{})
		if loginResponse["status_code"] != 200 {
			t.Fatalf("Expected session login status 200, got %v", loginResponse["status_code"])
		}

		// Second request to protected endpoint
		protectedResult := sessionGetFunc(server.URL + "/protected")
		protectedResponse := protectedResult.(map[string]interface{})
		if protectedResponse["status_code"] != 200 {
			t.Errorf("Expected session protected status 200, got %v", protectedResponse["status_code"])
		}
		if protectedResponse["text"] != "Access Granted" {
			t.Errorf("Expected 'Access Granted', got '%s'", protectedResponse["text"])
		}
	})
}
