package r2libs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// r2requests.go: Python requests-like HTTP client library for R2Lang

// Response represents an HTTP response with Python requests-like interface
type Response struct {
	URL        string
	StatusCode int
	Headers    map[string]interface{}
	Text       string
	JSON       interface{}
	Content    []byte
	OK         bool
	Elapsed    time.Duration
}

// Session represents an HTTP session for reusing connections and settings
type Session struct {
	Client  *http.Client
	Headers map[string]string
	Cookies map[string]string
	Auth    *BasicAuth
	Timeout time.Duration
	Verify  bool
	BaseURL string
}

// BasicAuth represents HTTP Basic Authentication
type BasicAuth struct {
	Username string
	Password string
}

func RegisterRequests(env *r2core.Environment) {
	// Global request functions
	env.Set("get", r2core.BuiltinFunction(globalGet))
	env.Set("post", r2core.BuiltinFunction(globalPost))
	env.Set("put", r2core.BuiltinFunction(globalPut))
	env.Set("delete", r2core.BuiltinFunction(globalDelete))
	env.Set("patch", r2core.BuiltinFunction(globalPatch))
	env.Set("head", r2core.BuiltinFunction(globalHead))
	env.Set("options", r2core.BuiltinFunction(globalOptions))

	// Session creation
	env.Set("session", r2core.BuiltinFunction(createSession))

	// Utility functions
	env.Set("urlencode", r2core.BuiltinFunction(urlEncode))
	env.Set("urldecode", r2core.BuiltinFunction(urlDecode))
}

// Global request functions that create a new session for each request
func globalGet(args ...interface{}) interface{} {
	return makeRequest("GET", args...)
}

func globalPost(args ...interface{}) interface{} {
	return makeRequest("POST", args...)
}

func globalPut(args ...interface{}) interface{} {
	return makeRequest("PUT", args...)
}

func globalDelete(args ...interface{}) interface{} {
	return makeRequest("DELETE", args...)
}

func globalPatch(args ...interface{}) interface{} {
	return makeRequest("PATCH", args...)
}

func globalHead(args ...interface{}) interface{} {
	return makeRequest("HEAD", args...)
}

func globalOptions(args ...interface{}) interface{} {
	return makeRequest("OPTIONS", args...)
}

// makeRequest creates a temporary session and makes a request
func makeRequest(method string, args ...interface{}) interface{} {
	if len(args) == 0 {
		panic(fmt.Sprintf("%s requires at least a URL", method))
	}

	url, ok := args[0].(string)
	if !ok {
		panic(fmt.Sprintf("%s: URL must be a string", method))
	}

	// Create temporary session
	session := &Session{
		Client:  &http.Client{Timeout: 30 * time.Second},
		Headers: make(map[string]string),
		Cookies: make(map[string]string),
		Timeout: 30 * time.Second,
		Verify:  true,
	}

	// Parse optional parameters
	var params map[string]interface{}
	if len(args) > 1 {
		if p, ok := args[1].(map[string]interface{}); ok {
			params = p
		}
	}

	return session.request(method, url, params)
}

// createSession creates a new Session object
func createSession(args ...interface{}) interface{} {
	session := &Session{
		Client:  &http.Client{Timeout: 30 * time.Second},
		Headers: make(map[string]string),
		Cookies: make(map[string]string),
		Timeout: 30 * time.Second,
		Verify:  true,
	}

	// Return session as a map with methods
	sessionMap := make(map[string]interface{})
	sessionMap["get"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		return session.sessionGet(args...)
	})
	sessionMap["post"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		return session.sessionPost(args...)
	})
	sessionMap["put"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		return session.sessionPut(args...)
	})
	sessionMap["delete"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		return session.sessionDelete(args...)
	})
	sessionMap["patch"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		return session.sessionPatch(args...)
	})
	sessionMap["head"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		return session.sessionHead(args...)
	})
	sessionMap["options"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		return session.sessionOptions(args...)
	})
	sessionMap["close"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		session.Client.CloseIdleConnections()
		return nil
	})

	return sessionMap
}

// Session methods
func (s *Session) sessionGet(args ...interface{}) interface{} {
	return s.sessionRequest("GET", args...)
}

func (s *Session) sessionPost(args ...interface{}) interface{} {
	return s.sessionRequest("POST", args...)
}

func (s *Session) sessionPut(args ...interface{}) interface{} {
	return s.sessionRequest("PUT", args...)
}

func (s *Session) sessionDelete(args ...interface{}) interface{} {
	return s.sessionRequest("DELETE", args...)
}

func (s *Session) sessionPatch(args ...interface{}) interface{} {
	return s.sessionRequest("PATCH", args...)
}

func (s *Session) sessionHead(args ...interface{}) interface{} {
	return s.sessionRequest("HEAD", args...)
}

func (s *Session) sessionOptions(args ...interface{}) interface{} {
	return s.sessionRequest("OPTIONS", args...)
}

func (s *Session) sessionRequest(method string, args ...interface{}) interface{} {
	if len(args) == 0 {
		panic(fmt.Sprintf("%s requires at least a URL", method))
	}

	url, ok := args[0].(string)
	if !ok {
		panic(fmt.Sprintf("%s: URL must be a string", method))
	}

	// Parse optional parameters
	var params map[string]interface{}
	if len(args) > 1 {
		if p, ok := args[1].(map[string]interface{}); ok {
			params = p
		}
	}

	return s.request(method, url, params)
}

// request performs the actual HTTP request
func (s *Session) request(method, reqURL string, params map[string]interface{}) interface{} {
	start := time.Now()

	// Handle base URL
	if s.BaseURL != "" && !strings.HasPrefix(reqURL, "http") {
		reqURL = strings.TrimSuffix(s.BaseURL, "/") + "/" + strings.TrimPrefix(reqURL, "/")
	}

	// Prepare request body
	var body io.Reader
	contentType := "application/json"

	if params != nil {
		// Handle data parameter
		if data, exists := params["data"]; exists {
			switch d := data.(type) {
			case string:
				body = strings.NewReader(d)
				contentType = "text/plain"
			case map[string]interface{}:
				jsonData, err := json.Marshal(d)
				if err != nil {
					panic(fmt.Sprintf("Failed to encode JSON data: %v", err))
				}
				body = bytes.NewReader(jsonData)
				contentType = "application/json"
			default:
				jsonData, err := json.Marshal(d)
				if err != nil {
					panic(fmt.Sprintf("Failed to encode data: %v", err))
				}
				body = bytes.NewReader(jsonData)
				contentType = "application/json"
			}
		}

		// Handle json parameter
		if jsonData, exists := params["json"]; exists {
			data, err := json.Marshal(jsonData)
			if err != nil {
				panic(fmt.Sprintf("Failed to encode JSON: %v", err))
			}
			body = bytes.NewReader(data)
			contentType = "application/json"
		}

		// Handle URL parameters
		if urlParams, exists := params["params"]; exists {
			if paramMap, ok := urlParams.(map[string]interface{}); ok {
				values := url.Values{}
				for k, v := range paramMap {
					values.Add(k, fmt.Sprint(v))
				}
				if strings.Contains(reqURL, "?") {
					reqURL += "&" + values.Encode()
				} else {
					reqURL += "?" + values.Encode()
				}
			}
		}
	}

	// Create request
	ctx := context.Background()
	if s.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.Timeout)
		defer cancel()
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, body)
	if err != nil {
		panic(fmt.Sprintf("Failed to create request: %v", err))
	}

	// Set headers
	req.Header.Set("Content-Type", contentType)
	for k, v := range s.Headers {
		req.Header.Set(k, v)
	}

	// Handle headers from params
	if params != nil {
		if headers, exists := params["headers"]; exists {
			if headerMap, ok := headers.(map[string]interface{}); ok {
				for k, v := range headerMap {
					req.Header.Set(k, fmt.Sprint(v))
				}
			}
		}

		// Handle timeout
		if timeout, exists := params["timeout"]; exists {
			if timeoutVal, ok := timeout.(float64); ok {
				duration := time.Duration(timeoutVal) * time.Second
				ctx, cancel := context.WithTimeout(ctx, duration)
				defer cancel()
				req = req.WithContext(ctx)
			}
		}
	}

	// Set cookies
	for name, value := range s.Cookies {
		req.AddCookie(&http.Cookie{Name: name, Value: value})
	}

	// Handle authentication
	if s.Auth != nil {
		req.SetBasicAuth(s.Auth.Username, s.Auth.Password)
	}
	if params != nil {
		if auth, exists := params["auth"]; exists {
			if authArray, ok := auth.([]interface{}); ok && len(authArray) == 2 {
				username := fmt.Sprint(authArray[0])
				password := fmt.Sprint(authArray[1])
				req.SetBasicAuth(username, password)
			}
		}
	}

	// Make request
	resp, err := s.Client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("Request failed: %v", err))
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Failed to read response body: %v", err))
	}

	// Parse response headers
	headers := make(map[string]interface{})
	for name, values := range resp.Header {
		if len(values) == 1 {
			headers[name] = values[0]
		} else {
			headers[name] = values
		}
	}

	// Try to parse JSON response
	var jsonResponse interface{}
	if len(respBody) > 0 {
		contentType := resp.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			json.Unmarshal(respBody, &jsonResponse)
		}
	}

	// Create response object
	response := &Response{
		URL:        resp.Request.URL.String(),
		StatusCode: resp.StatusCode,
		Headers:    headers,
		Text:       string(respBody),
		JSON:       jsonResponse,
		Content:    respBody,
		OK:         resp.StatusCode >= 200 && resp.StatusCode < 300,
		Elapsed:    time.Since(start),
	}

	// Return response as map
	return responseToMap(response)
}

// responseToMap converts Response struct to map for R2Lang
func responseToMap(resp *Response) map[string]interface{} {
	return map[string]interface{}{
		"url":         resp.URL,
		"status_code": resp.StatusCode,
		"headers":     resp.Headers,
		"text":        resp.Text,
		"json":        resp.JSON,
		"content":     resp.Content,
		"ok":          resp.OK,
		"elapsed":     resp.Elapsed.Seconds(),
		"json_func": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			return resp.JSON
		}),
		"raise_for_status": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if !resp.OK {
				panic(fmt.Sprintf("HTTP %d Error: %s", resp.StatusCode, resp.Text))
			}
			return nil
		}),
	}
}

// Utility functions
func urlEncode(args ...interface{}) interface{} {
	if len(args) == 0 {
		panic("urlencode requires a string argument")
	}

	str, ok := args[0].(string)
	if !ok {
		panic("urlencode: argument must be a string")
	}

	return url.QueryEscape(str)
}

func urlDecode(args ...interface{}) interface{} {
	if len(args) == 0 {
		panic("urldecode requires a string argument")
	}

	str, ok := args[0].(string)
	if !ok {
		panic("urldecode: argument must be a string")
	}

	decoded, err := url.QueryUnescape(str)
	if err != nil {
		panic(fmt.Sprintf("urldecode failed: %v", err))
	}

	return decoded
}
