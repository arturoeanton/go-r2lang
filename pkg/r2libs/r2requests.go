package r2libs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

var (
	globalCookieJar     *cookiejar.Jar
	globalCookieJarOnce sync.Once
)

func initGlobalCookieJar() {
	globalCookieJarOnce.Do(func() {
		jar, _ := cookiejar.New(nil)
		globalCookieJar = jar
	})
}

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
	Client     *http.Client
	Headers    map[string]string
	Auth       *BasicAuth
	Timeout    time.Duration
	Verify     bool
	BaseURL    string
	Proxies    map[string]string
	MaxRetries int
	RetryDelay time.Duration
}

// BasicAuth represents HTTP Basic Authentication
type BasicAuth struct {
	Username string
	Password string
}

func RegisterRequests(env *r2core.Environment) {
	initGlobalCookieJar()

	functions := map[string]r2core.BuiltinFunction{
		"get":       r2core.BuiltinFunction(globalGet),
		"post":      r2core.BuiltinFunction(globalPost),
		"put":       r2core.BuiltinFunction(globalPut),
		"delete":    r2core.BuiltinFunction(globalDelete),
		"patch":     r2core.BuiltinFunction(globalPatch),
		"head":      r2core.BuiltinFunction(globalHead),
		"options":   r2core.BuiltinFunction(globalOptions),
		"session":   r2core.BuiltinFunction(createSession),
		"urlencode": r2core.BuiltinFunction(urlEncode),
		"urldecode": r2core.BuiltinFunction(urlDecode),
	}

	RegisterModule(env, "request", functions)
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
		Client:     &http.Client{Jar: globalCookieJar, Timeout: 30 * time.Second},
		Headers:    make(map[string]string),
		Timeout:    30 * time.Second,
		Verify:     true,
		Proxies:    make(map[string]string),
		MaxRetries: 0,
		RetryDelay: time.Second,
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
	jar, _ := cookiejar.New(nil)
	session := &Session{
		Client:     &http.Client{Jar: jar, Timeout: 30 * time.Second},
		Headers:    make(map[string]string),
		Timeout:    30 * time.Second,
		Verify:     true,
		Proxies:    make(map[string]string),
		MaxRetries: 0,
		RetryDelay: time.Second,
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

	// Prepare request body data. We need to be able to re-read this for retries.
	var bodyBytes []byte
	var contentType string

	if params != nil {
		if files, exists := params["files"]; exists {
			if filesMap, ok := files.(map[string]interface{}); ok {
				var b bytes.Buffer
				writer := multipart.NewWriter(&b)
				for key, value := range filesMap {
					if filePath, ok := value.(string); ok {
						file, err := os.Open(filePath)
						if err != nil {
							panic(fmt.Sprintf("Failed to open file: %v", err))
						}
						defer file.Close()
						part, err := writer.CreateFormFile(key, filepath.Base(filePath))
						if err != nil {
							panic(fmt.Sprintf("Failed to create form file: %v", err))
						}
						_, err = io.Copy(part, file)
						if err != nil {
							panic(fmt.Sprintf("Failed to copy file content: %v", err))
						}
					}
				}
				if data, exists := params["data"]; exists {
					if dataMap, ok := data.(map[string]interface{}); ok {
						for key, value := range dataMap {
							_ = writer.WriteField(key, fmt.Sprint(value))
						}
					}
				}
				writer.Close()
				bodyBytes = b.Bytes()
				contentType = writer.FormDataContentType()
			}
		} else {
			contentType = "application/json" // Default
			if data, exists := params["data"]; exists {
				switch d := data.(type) {
				case string:
					bodyBytes = []byte(d)
					contentType = "text/plain"
				default:
					jsonData, err := json.Marshal(d)
					if err != nil {
						panic(fmt.Sprintf("Failed to encode data: %v", err))
					}
					bodyBytes = jsonData
				}
			}
			if jsonData, exists := params["json"]; exists {
				data, err := json.Marshal(jsonData)
				if err != nil {
					panic(fmt.Sprintf("Failed to encode JSON: %v", err))
				}
				bodyBytes = data
				contentType = "application/json"
			}
		}

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
		if proxies, exists := params["proxies"]; exists {
			if proxiesMap, ok := proxies.(map[string]interface{}); ok {
				proxyURL, err := getProxy(proxiesMap)
				if err != nil {
					panic(fmt.Sprintf("Invalid proxy configuration: %v", err))
				}
				if proxyURL != nil {
					s.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
				}
			}
		}
		if retries, exists := params["retries"]; exists {
			if retriesMap, ok := retries.(map[string]interface{}); ok {
				if max, ok := retriesMap["max"].(float64); ok {
					s.MaxRetries = int(max)
				}
				if delay, ok := retriesMap["delay"].(float64); ok {
					s.RetryDelay = time.Duration(delay) * time.Second
				}
			}
		}
	}

	var resp *http.Response
	var err error

	for i := 0; i <= s.MaxRetries; i++ {
		ctx := context.Background()
		if s.Timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, s.Timeout)
			defer cancel()
		}
		if params != nil {
			if timeout, exists := params["timeout"]; exists {
				if timeoutVal, ok := timeout.(float64); ok {
					duration := time.Duration(timeoutVal) * time.Second
					var cancel context.CancelFunc
					ctx, cancel = context.WithTimeout(ctx, duration)
					defer cancel()
				}
			}
		}

		var bodyReader io.Reader
		if bodyBytes != nil {
			bodyReader = bytes.NewReader(bodyBytes)
		}

		req, reqErr := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
		if reqErr != nil {
			panic(fmt.Sprintf("Failed to create request: %v", reqErr))
		}

		if contentType != "" {
			req.Header.Set("Content-Type", contentType)
		}
		for k, v := range s.Headers {
			req.Header.Set(k, v)
		}
		if params != nil {
			if headers, exists := params["headers"]; exists {
				if headerMap, ok := headers.(map[string]interface{}); ok {
					for k, v := range headerMap {
						req.Header.Set(k, fmt.Sprint(v))
					}
				}
			}
		}

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

		resp, err = s.Client.Do(req)
		if err == nil && resp.StatusCode < 500 {
			break
		}

		if resp != nil {
			resp.Body.Close()
		}

		if i < s.MaxRetries {
			time.Sleep(s.RetryDelay)
		}
	}

	if err != nil {
		panic(fmt.Sprintf("Request failed after %d retries: %v", s.MaxRetries, err))
	}
	defer resp.Body.Close()

	respBody, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		panic(fmt.Sprintf("Failed to read response body: %v", readErr))
	}

	headers := make(map[string]interface{})
	for name, values := range resp.Header {
		if len(values) == 1 {
			headers[name] = values[0]
		} else {
			headers[name] = values
		}
	}

	var jsonResponse interface{}
	if len(respBody) > 0 {
		contentType := resp.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			json.Unmarshal(respBody, &jsonResponse)
		}
	}

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

func getProxy(proxiesMap map[string]interface{}) (*url.URL, error) {
	if httpProxy, ok := proxiesMap["http"]; ok {
		if httpProxyStr, ok := httpProxy.(string); ok {
			return url.Parse(httpProxyStr)
		}
	}
	if httpsProxy, ok := proxiesMap["https"]; ok {
		if httpsProxyStr, ok := httpsProxy.(string); ok {
			return url.Parse(httpsProxyStr)
		}
	}
	return nil, nil
}
