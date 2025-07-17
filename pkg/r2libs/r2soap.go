package r2libs

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// r2soap.go: Dynamic SOAP client library for R2Lang with WSDL parsing

// SOAPClient represents a dynamic SOAP client
type SOAPClient struct {
	WSDLURL       string
	ServiceURL    string
	Namespace     string
	Operations    map[string]*SOAPOperation
	HTTPTimeout   time.Duration
	Headers       map[string]string
	TLSConfig     *tls.Config
	SkipTLSVerify bool
	Auth          *SOAPAuth
}

// SOAPAuth represents authentication configuration
type SOAPAuth struct {
	Type     string // "basic", "bearer", "certificate"
	Username string
	Password string
	Token    string
	CertFile string
	KeyFile  string
}

// SOAPOperation represents a SOAP operation with its parameters
type SOAPOperation struct {
	Name         string
	SOAPAction   string
	InputMessage string
	Parameters   []SOAPParameter
}

// SOAPParameter represents a parameter for a SOAP operation
type SOAPParameter struct {
	Name string
	Type string
}

// WSDLDefinitions represents the root element of a WSDL document
type WSDLDefinitions struct {
	XMLName   xml.Name       `xml:"definitions"`
	TargetNS  string         `xml:"targetNamespace,attr"`
	Services  []WSDLService  `xml:"service"`
	PortTypes []WSDLPortType `xml:"portType"`
	Bindings  []WSDLBinding  `xml:"binding"`
	Messages  []WSDLMessage  `xml:"message"`
	Types     WSDLTypes      `xml:"types"`
}

// WSDLService represents a WSDL service
type WSDLService struct {
	Name  string     `xml:"name,attr"`
	Ports []WSDLPort `xml:"port"`
}

// WSDLPort represents a WSDL port
type WSDLPort struct {
	Name    string          `xml:"name,attr"`
	Address WSDLSoapAddress `xml:"address"`
}

// WSDLSoapAddress represents a SOAP address
type WSDLSoapAddress struct {
	Location string `xml:"location,attr"`
}

// WSDLPortType represents a WSDL port type
type WSDLPortType struct {
	Name       string          `xml:"name,attr"`
	Operations []WSDLOperation `xml:"operation"`
}

// WSDLOperation represents a WSDL operation
type WSDLOperation struct {
	Name   string      `xml:"name,attr"`
	Input  WSDLMessage `xml:"input"`
	Output WSDLMessage `xml:"output"`
}

// WSDLBinding represents a WSDL binding
type WSDLBinding struct {
	Name       string                 `xml:"name,attr"`
	Type       string                 `xml:"type,attr"`
	Operations []WSDLBindingOperation `xml:"operation"`
}

// WSDLBindingOperation represents a WSDL binding operation
type WSDLBindingOperation struct {
	Name       string         `xml:"name,attr"`
	SOAPAction WSDLSoapAction `xml:"operation"`
}

// WSDLSoapAction represents a SOAP action
type WSDLSoapAction struct {
	SOAPAction string `xml:"soapAction,attr"`
}

// WSDLMessage represents a WSDL message
type WSDLMessage struct {
	Name  string     `xml:"name,attr"`
	Parts []WSDLPart `xml:"part"`
}

// WSDLPart represents a WSDL message part
type WSDLPart struct {
	Name    string `xml:"name,attr"`
	Element string `xml:"element,attr"`
	Type    string `xml:"type,attr"`
}

// WSDLTypes represents WSDL types
type WSDLTypes struct {
	Schemas []WSDLSchema `xml:"schema"`
}

// WSDLSchema represents a WSDL schema
type WSDLSchema struct {
	Elements []WSDLElement `xml:"element"`
}

// WSDLElement represents a WSDL element
type WSDLElement struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

// RegisterSOAP registers SOAP functions in R2Lang environment
func RegisterSOAP(env *r2core.Environment) {
	// Create SOAP client from WSDL with optional custom headers
	env.Set("soapClient", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("soapClient requires (wsdlURL, [customHeaders])")
		}

		wsdlURL, ok := args[0].(string)
		if !ok {
			panic("soapClient: wsdlURL must be a string")
		}

		// Parse optional custom headers
		var customHeaders map[string]interface{}
		if len(args) > 1 {
			if headers, ok := args[1].(map[string]interface{}); ok {
				customHeaders = headers
			} else {
				panic("soapClient: customHeaders must be a map")
			}
		}

		client, err := createSOAPClient(wsdlURL, customHeaders)
		if err != nil {
			// Provide more detailed error information
			errorMsg := fmt.Sprintf("soapClient: failed to create client from '%s'", wsdlURL)
			if strings.Contains(err.Error(), "connection reset") {
				errorMsg += " - Network connectivity issue. The WSDL service may be unavailable or blocking requests."
			} else if strings.Contains(err.Error(), "no such host") {
				errorMsg += " - DNS resolution failed. Check the URL and internet connectivity."
			} else if strings.Contains(err.Error(), "timeout") {
				errorMsg += " - Request timeout. The service may be slow or overloaded."
			}
			errorMsg += fmt.Sprintf(" Error: %v", err)
			panic(errorMsg)
		}

		return soapClientToMap(client)
	}))

	// Create simple SOAP envelope
	env.Set("soapEnvelope", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 3 {
			panic("soapEnvelope requires (namespace, methodName, bodyContent)")
		}

		namespace, ok1 := args[0].(string)
		methodName, ok2 := args[1].(string)
		bodyContent, ok3 := args[2].(string)

		if !ok1 || !ok2 || !ok3 {
			panic("soapEnvelope: all parameters must be strings")
		}

		envelope := createSOAPEnvelope(namespace, methodName, bodyContent)
		return envelope
	}))

	// Send raw SOAP request
	env.Set("soapRequest", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 3 {
			panic("soapRequest requires (url, soapAction, envelope)")
		}

		url, ok1 := args[0].(string)
		soapAction, ok2 := args[1].(string)
		envelope, ok3 := args[2].(string)

		if !ok1 || !ok2 || !ok3 {
			panic("soapRequest: all parameters must be strings")
		}

		response, err := sendSOAPRequest(url, soapAction, envelope)
		if err != nil {
			panic(fmt.Sprintf("soapRequest: %v", err))
		}

		return response
	}))
}

// createSOAPClient creates a new SOAP client from WSDL URL with optional custom headers
func createSOAPClient(wsdlURL string, customHeaders map[string]interface{}) (*SOAPClient, error) {
	// Create HTTP client with TLS configuration
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false, // Default to secure
			MinVersion:         tls.VersionTLS12,
		},
	}
	httpClient := &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}

	req, err := http.NewRequest("GET", wsdlURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set browser-like headers to avoid blocking
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/xml,application/xml,*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	// Fetch WSDL document
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch WSDL: %v", err)
	}
	defer resp.Body.Close()

	wsdlData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read WSDL: %v", err)
	}

	// Parse WSDL
	var wsdl WSDLDefinitions
	err = xml.Unmarshal(wsdlData, &wsdl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse WSDL: %v", err)
	}

	// Extract service URL
	serviceURL := ""
	if len(wsdl.Services) > 0 && len(wsdl.Services[0].Ports) > 0 {
		serviceURL = wsdl.Services[0].Ports[0].Address.Location
	}

	// Extract operations
	operations := make(map[string]*SOAPOperation)

	// Build operation map from portTypes and bindings
	for _, portType := range wsdl.PortTypes {
		for _, op := range portType.Operations {
			operation := &SOAPOperation{
				Name:         op.Name,
				InputMessage: op.Input.Name,
				Parameters:   []SOAPParameter{},
			}

			// Find corresponding binding for SOAPAction
			for _, binding := range wsdl.Bindings {
				for _, bindingOp := range binding.Operations {
					if bindingOp.Name == op.Name {
						operation.SOAPAction = bindingOp.SOAPAction.SOAPAction
						break
					}
				}
			}

			operations[op.Name] = operation
		}
	}

	// Initialize client with browser-like defaults
	defaultHeaders := map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Accept":          "text/xml,application/xml,*/*",
		"Accept-Language": "en-US,en;q=0.9",
		"Accept-Encoding": "gzip, deflate",
		"Connection":      "keep-alive",
	}

	// Override with custom headers if provided
	if customHeaders != nil {
		for key, value := range customHeaders {
			if strValue, ok := value.(string); ok {
				defaultHeaders[key] = strValue
			}
		}
	}

	client := &SOAPClient{
		WSDLURL:       wsdlURL,
		ServiceURL:    serviceURL,
		Namespace:     wsdl.TargetNS,
		Operations:    operations,
		HTTPTimeout:   30 * time.Second,
		Headers:       defaultHeaders,
		TLSConfig:     &tls.Config{},
		SkipTLSVerify: false,
		Auth:          nil,
	}

	return client, nil
}

// soapClientToMap converts SOAPClient to R2Lang map with methods
func soapClientToMap(client *SOAPClient) map[string]interface{} {
	clientMap := make(map[string]interface{})

	// Client properties
	clientMap["wsdlURL"] = client.WSDLURL
	clientMap["serviceURL"] = client.ServiceURL
	clientMap["namespace"] = client.Namespace

	// List available operations
	clientMap["listOperations"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		operations := make([]interface{}, 0, len(client.Operations))
		for name := range client.Operations {
			operations = append(operations, name)
		}
		return operations
	})

	// Get operation info
	clientMap["getOperation"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("getOperation requires (operationName)")
		}

		operationName, ok := args[0].(string)
		if !ok {
			panic("getOperation: operationName must be a string")
		}

		operation, exists := client.Operations[operationName]
		if !exists {
			panic(fmt.Sprintf("getOperation: operation '%s' not found", operationName))
		}

		return map[string]interface{}{
			"name":       operation.Name,
			"soapAction": operation.SOAPAction,
			"message":    operation.InputMessage,
		}
	})

	// Call SOAP operation dynamically with response parsing
	clientMap["call"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("call requires (operationName, parameters)")
		}

		operationName, ok1 := args[0].(string)
		params, ok2 := args[1].(map[string]interface{})

		if !ok1 {
			panic("call: operationName must be a string")
		}
		if !ok2 {
			panic("call: parameters must be a map")
		}

		return client.callOperation(operationName, params)
	})

	// Call SOAP operation and return only the result value (simplified)
	clientMap["callSimple"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("callSimple requires (operationName, parameters)")
		}

		operationName, ok1 := args[0].(string)
		params, ok2 := args[1].(map[string]interface{})

		if !ok1 {
			panic("callSimple: operationName must be a string")
		}
		if !ok2 {
			panic("callSimple: parameters must be a map")
		}

		response := client.callOperation(operationName, params)
		if responseMap, ok := response.(map[string]interface{}); ok {
			// Check if the call was successful
			if success, exists := responseMap["success"]; exists && success == true {
				// Return just the result value if available
				if result, exists := responseMap["result"]; exists && result != nil {
					return result
				}
				// Fall back to first value in values map
				if values, exists := responseMap["values"]; exists {
					if valuesMap, ok := values.(map[string]interface{}); ok && len(valuesMap) > 0 {
						// Return the first non-empty value
						for _, value := range valuesMap {
							if value != nil && value != "" {
								return value
							}
						}
						return valuesMap
					}
				}
			} else {
				// If call failed, return error info
				if fault, exists := responseMap["fault"]; exists {
					return fault
				}
				if err, exists := responseMap["error"]; exists {
					return err
				}
			}
		}
		return response
	})

	// Get response in raw XML format
	clientMap["callRaw"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("callRaw requires (operationName, parameters)")
		}

		operationName, ok1 := args[0].(string)
		params, ok2 := args[1].(map[string]interface{})

		if !ok1 {
			panic("callRaw: operationName must be a string")
		}
		if !ok2 {
			panic("callRaw: parameters must be a map")
		}

		// Call operation and extract raw response
		operation, exists := client.Operations[operationName]
		if !exists {
			panic(fmt.Sprintf("operation '%s' not found", operationName))
		}

		var bodyParts []string
		for key, value := range params {
			bodyParts = append(bodyParts, fmt.Sprintf("<%s>%v</%s>", key, value, key))
		}
		bodyContent := strings.Join(bodyParts, "")
		envelope := createOperationSOAPEnvelope(client.Namespace, operationName, bodyContent)

		response, err := client.sendRequest(operation.SOAPAction, envelope)
		if err != nil {
			panic(fmt.Sprintf("SOAP call failed: %v", err))
		}

		return response
	})

	// Set HTTP timeout
	clientMap["setTimeout"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("setTimeout requires (seconds)")
		}

		seconds, ok := args[0].(float64)
		if !ok {
			panic("setTimeout: seconds must be a number")
		}

		client.HTTPTimeout = time.Duration(seconds) * time.Second
		return nil
	})

	// Set custom headers (single or multiple)
	clientMap["setHeader"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("setHeader requires (name, value) or (headersMap)")
		}

		// If single argument, treat as headers map
		if len(args) == 1 {
			if headersMap, ok := args[0].(map[string]interface{}); ok {
				for name, value := range headersMap {
					if strValue, ok := value.(string); ok {
						client.Headers[name] = strValue
					}
				}
				return nil
			}
		}

		// If two arguments, treat as name-value pair
		if len(args) >= 2 {
			name, ok1 := args[0].(string)
			value, ok2 := args[1].(string)

			if !ok1 || !ok2 {
				panic("setHeader: name and value must be strings")
			}

			client.Headers[name] = value
			return nil
		}

		panic("setHeader: invalid arguments")
	})

	// Get current headers
	clientMap["getHeaders"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		headersMap := make(map[string]interface{})
		for name, value := range client.Headers {
			headersMap[name] = value
		}
		return headersMap
	})

	// Reset headers to browser-like defaults
	clientMap["resetHeaders"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		client.Headers = map[string]string{
			"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			"Accept":          "text/xml,application/xml,*/*",
			"Accept-Language": "en-US,en;q=0.9",
			"Accept-Encoding": "gzip, deflate",
			"Connection":      "keep-alive",
		}
		return nil
	})

	// Remove specific header
	clientMap["removeHeader"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("removeHeader requires (headerName)")
		}

		headerName, ok := args[0].(string)
		if !ok {
			panic("removeHeader: headerName must be a string")
		}

		delete(client.Headers, headerName)
		return nil
	})

	// Configure SSL/TLS settings
	clientMap["setTLSConfig"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("setTLSConfig requires (configMap)")
		}

		configMap, ok := args[0].(map[string]interface{})
		if !ok {
			panic("setTLSConfig: configMap must be a map")
		}

		// Skip certificate verification
		if skipVerify, exists := configMap["skipVerify"]; exists {
			if skip, ok := skipVerify.(bool); ok {
				client.SkipTLSVerify = skip
				client.TLSConfig.InsecureSkipVerify = skip
			}
		}

		// Set minimum TLS version
		if minVersion, exists := configMap["minVersion"]; exists {
			if version, ok := minVersion.(string); ok {
				switch version {
				case "1.0":
					client.TLSConfig.MinVersion = tls.VersionTLS10
				case "1.1":
					client.TLSConfig.MinVersion = tls.VersionTLS11
				case "1.2":
					client.TLSConfig.MinVersion = tls.VersionTLS12
				case "1.3":
					client.TLSConfig.MinVersion = tls.VersionTLS13
				default:
					client.TLSConfig.MinVersion = tls.VersionTLS12 // Default to 1.2
				}
			}
		}

		return nil
	})

	// Set authentication credentials
	clientMap["setAuth"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("setAuth requires (authConfig)")
		}

		authConfig, ok := args[0].(map[string]interface{})
		if !ok {
			panic("setAuth: authConfig must be a map")
		}

		auth := &SOAPAuth{}

		if authType, exists := authConfig["type"]; exists {
			if typeStr, ok := authType.(string); ok {
				auth.Type = typeStr
			}
		}

		if username, exists := authConfig["username"]; exists {
			if userStr, ok := username.(string); ok {
				auth.Username = userStr
			}
		}

		if password, exists := authConfig["password"]; exists {
			if passStr, ok := password.(string); ok {
				auth.Password = passStr
			}
		}

		if token, exists := authConfig["token"]; exists {
			if tokenStr, ok := token.(string); ok {
				auth.Token = tokenStr
			}
		}

		client.Auth = auth
		return nil
	})

	// Add custom CA certificate
	clientMap["addCACert"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("addCACert requires (certPath)")
		}

		certPath, ok := args[0].(string)
		if !ok {
			panic("addCACert: certPath must be a string")
		}

		// This would require reading cert file - placeholder for now
		// In a full implementation, you'd read the cert file and add to TLSConfig.RootCAs
		client.Headers["X-Custom-CA"] = certPath
		return nil
	})

	return clientMap
}

// callOperation calls a SOAP operation with given parameters
func (client *SOAPClient) callOperation(operationName string, params map[string]interface{}) interface{} {
	operation, exists := client.Operations[operationName]
	if !exists {
		panic(fmt.Sprintf("operation '%s' not found", operationName))
	}

	// Build SOAP body from parameters with proper namespace
	var bodyParts []string
	for key, value := range params {
		bodyParts = append(bodyParts, fmt.Sprintf("<%s>%v</%s>", key, value, key))
	}
	bodyContent := strings.Join(bodyParts, "")

	// Create SOAP envelope with operation-specific structure
	envelope := createOperationSOAPEnvelope(client.Namespace, operationName, bodyContent)

	// Send SOAP request
	response, err := client.sendRequest(operation.SOAPAction, envelope)
	if err != nil {
		panic(fmt.Sprintf("SOAP call failed: %v", err))
	}

	// Parse response to R2Lang native types
	parsedResponse := parseSOAPResponseToR2Lang(response)
	return parsedResponse
}

// createOperationSOAPEnvelope creates a SOAP envelope specifically for operation calls
func createOperationSOAPEnvelope(namespace, operationName, bodyContent string) string {
	envelope := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
	<soap:Body>
		<%s xmlns="%s">
			%s
		</%s>
	</soap:Body>
</soap:Envelope>`, operationName, namespace, bodyContent, operationName)

	return envelope
}

// sendRequest sends a SOAP request with the client's configuration
func (client *SOAPClient) sendRequest(soapAction, envelope string) (string, error) {
	// Create HTTP client with TLS configuration
	transport := &http.Transport{
		TLSClientConfig: client.TLSConfig,
	}
	httpClient := &http.Client{
		Timeout:   client.HTTPTimeout,
		Transport: transport,
	}

	req, err := http.NewRequest("POST", client.ServiceURL, strings.NewReader(envelope))
	if err != nil {
		return "", err
	}

	// Set SOAP headers
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", fmt.Sprintf(`"%s"`, soapAction))

	// Apply authentication if configured
	if client.Auth != nil {
		switch client.Auth.Type {
		case "basic":
			if client.Auth.Username != "" && client.Auth.Password != "" {
				credentials := base64.StdEncoding.EncodeToString([]byte(client.Auth.Username + ":" + client.Auth.Password))
				req.Header.Set("Authorization", "Basic "+credentials)
			}
		case "bearer":
			if client.Auth.Token != "" {
				req.Header.Set("Authorization", "Bearer "+client.Auth.Token)
			}
		}
	}

	// Set custom headers (includes defaults)
	for name, value := range client.Headers {
		req.Header.Set(name, value)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(responseData))
	}

	return string(responseData), nil
}

// createSOAPEnvelope creates a SOAP 1.1 envelope
func createSOAPEnvelope(namespace, methodName, bodyContent string) string {
	envelope := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tns="%s">
	<soap:Header />
	<soap:Body>
		<tns:%s>
			%s
		</tns:%s>
	</soap:Body>
</soap:Envelope>`, namespace, methodName, bodyContent, methodName)

	return envelope
}

// sendSOAPRequest sends a raw SOAP request
func sendSOAPRequest(url, soapAction, envelope string) (string, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest("POST", url, strings.NewReader(envelope))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", fmt.Sprintf(`"%s"`, soapAction))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(responseData))
	}

	return string(responseData), nil
}

// parseSOAPResponseToR2Lang parses SOAP response and ensures R2Lang native types
func parseSOAPResponseToR2Lang(response string) map[string]interface{} {
	// Clean the response from any binary/control characters
	cleanResponse := cleanResponseString(response)

	result := map[string]interface{}{
		"success": true,
		"raw":     cleanResponse,
	}

	// Extract SOAP body content with multiple regex patterns
	var bodyContent string
	bodyPatterns := []string{
		`(?s)<soap:Body[^>]*>(.*?)</soap:Body>`,
		`(?s)<Body[^>]*>(.*?)</Body>`,
		`(?s)<.*:Body[^>]*>(.*?)</.*:Body>`,
	}

	for _, pattern := range bodyPatterns {
		bodyRegex := regexp.MustCompile(pattern)
		matches := bodyRegex.FindStringSubmatch(cleanResponse)
		if len(matches) >= 2 {
			bodyContent = matches[1]
			break
		}
	}

	if bodyContent == "" {
		// If no body found, return raw response
		result["success"] = false
		result["error"] = "Could not extract SOAP body"
		return result
	}

	result["body"] = bodyContent

	// Check for SOAP faults
	if strings.Contains(bodyContent, "Fault") || strings.Contains(bodyContent, "fault") {
		result["success"] = false
		result["fault"] = extractSOAPFaultSimple(bodyContent)
		return result
	}

	// Extract values as native R2Lang types
	values := extractValuesAsR2LangTypes(bodyContent)
	result["values"] = values

	// Extract the primary result value (for callSimple)
	primaryResult := extractPrimaryResult(bodyContent)
	if primaryResult != nil {
		result["result"] = primaryResult
	}

	return result
}

// parseSOAPResponse parses a SOAP response and extracts meaningful data
func parseSOAPResponse(response string) (map[string]interface{}, error) {
	// Initialize result with success indicator
	result := map[string]interface{}{
		"success": true,
		"raw":     response,
	}

	// Extract SOAP body content
	bodyRegex := regexp.MustCompile(`<soap:Body[^>]*>(.*?)</soap:Body>`)
	if !bodyRegex.MatchString(response) {
		// Try alternative format without namespace
		bodyRegex = regexp.MustCompile(`<Body[^>]*>(.*?)</Body>`)
	}

	// Try with multiline flag for complex XML
	if !bodyRegex.MatchString(response) {
		bodyRegex = regexp.MustCompile(`(?s)<soap:Body[^>]*>(.*?)</soap:Body>`)
		if !bodyRegex.MatchString(response) {
			bodyRegex = regexp.MustCompile(`(?s)<Body[^>]*>(.*?)</Body>`)
		}
	}

	matches := bodyRegex.FindStringSubmatch(response)
	if len(matches) < 2 {
		return result, fmt.Errorf("could not extract SOAP body")
	}

	bodyContent := matches[1]
	result["body"] = bodyContent

	// Check for SOAP faults
	if strings.Contains(bodyContent, "Fault") || strings.Contains(bodyContent, "fault") {
		result["success"] = false
		faultDetail := extractSOAPFault(bodyContent)
		result["fault"] = faultDetail
		return result, nil
	}

	// Extract operation result values using improved parsing
	values := extractResponseValues(bodyContent)
	result["values"] = values

	// Extract numeric results (common for calculator services)
	numericResult := extractNumericResult(bodyContent)
	if numericResult != nil {
		result["result"] = numericResult
	}

	// Provide cleaned XML for debugging
	result["cleaned"] = cleanXMLNamespaces(bodyContent)

	return result, nil
}

// extractSOAPFault extracts fault information from SOAP response
func extractSOAPFault(bodyContent string) map[string]interface{} {
	fault := make(map[string]interface{})

	// Extract fault code
	codeRegex := regexp.MustCompile(`<faultcode[^>]*>([^<]+)</faultcode>`)
	if matches := codeRegex.FindStringSubmatch(bodyContent); len(matches) > 1 {
		fault["code"] = strings.TrimSpace(matches[1])
	}

	// Extract fault string/message
	stringRegex := regexp.MustCompile(`<faultstring[^>]*>([^<]+)</faultstring>`)
	if matches := stringRegex.FindStringSubmatch(bodyContent); len(matches) > 1 {
		fault["message"] = strings.TrimSpace(matches[1])
	}

	return fault
}

// extractResponseValues extracts key-value pairs from response body
func extractResponseValues(bodyContent string) map[string]interface{} {
	values := make(map[string]interface{})

	// Extract all XML elements and their values
	elementRegex := regexp.MustCompile(`<([^/>\\s]+)[^>]*>([^<]+)</[^>]+>`)
	matches := elementRegex.FindAllStringSubmatch(bodyContent, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			key := strings.TrimSpace(match[1])
			value := strings.TrimSpace(match[2])

			// Try to convert to appropriate type
			if convertedValue := tryConvertValue(value); convertedValue != nil {
				values[key] = convertedValue
			} else {
				values[key] = value
			}
		}
	}

	return values
}

// extractNumericResult extracts numeric results from common patterns
func extractNumericResult(bodyContent string) interface{} {
	// Common result patterns for calculator services
	patterns := []string{
		`<[^>]*Result[^>]*>([^<]+)<`,
		`<result[^>]*>([^<]+)<`,
		`<return[^>]*>([^<]+)<`,
		`<value[^>]*>([^<]+)<`,
	}

	for _, pattern := range patterns {
		regex := regexp.MustCompile(pattern)
		if matches := regex.FindStringSubmatch(bodyContent); len(matches) > 1 {
			value := strings.TrimSpace(matches[1])
			return tryConvertValue(value)
		}
	}

	return nil
}

// extractSOAPFaultSimple extracts basic fault information
func extractSOAPFaultSimple(bodyContent string) map[string]interface{} {
	fault := map[string]interface{}{}

	// Extract fault code
	if codeRegex := regexp.MustCompile(`<(?:\w+:)?faultcode[^>]*>([^<]+)</(?:\w+:)?faultcode>`); codeRegex.MatchString(bodyContent) {
		matches := codeRegex.FindStringSubmatch(bodyContent)
		if len(matches) > 1 {
			fault["code"] = strings.TrimSpace(matches[1])
		}
	}

	// Extract fault string
	if stringRegex := regexp.MustCompile(`<(?:\w+:)?faultstring[^>]*>([^<]+)</(?:\w+:)?faultstring>`); stringRegex.MatchString(bodyContent) {
		matches := stringRegex.FindStringSubmatch(bodyContent)
		if len(matches) > 1 {
			fault["message"] = strings.TrimSpace(matches[1])
		}
	}

	return fault
}

// extractValuesAsR2LangTypes extracts all values ensuring R2Lang native types
func extractValuesAsR2LangTypes(bodyContent string) map[string]interface{} {
	values := map[string]interface{}{}

	// Extract all XML elements and their values
	elementRegex := regexp.MustCompile(`<([^/>\\s]+)[^>]*>([^<]+)</[^>]+>`)
	matches := elementRegex.FindAllStringSubmatch(bodyContent, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			key := strings.TrimSpace(match[1])
			value := strings.TrimSpace(match[2])

			// Convert to R2Lang native type
			values[key] = convertToR2LangType(value)
		}
	}

	return values
}

// extractPrimaryResult extracts the main result value for callSimple
func extractPrimaryResult(bodyContent string) interface{} {
	// Common result patterns for SOAP responses
	patterns := []string{
		`<([^>]*Result[^>]*)>([^<]+)</[^>]*Result[^>]*>`,
		`<([^>]*Response[^>]*)>([^<]+)</[^>]*Response[^>]*>`,
		`<([^>]*Return[^>]*)>([^<]+)</[^>]*Return[^>]*>`,
		`<([^>]*Value[^>]*)>([^<]+)</[^>]*Value[^>]*>`,
	}

	for _, pattern := range patterns {
		regex := regexp.MustCompile(pattern)
		if matches := regex.FindStringSubmatch(bodyContent); len(matches) > 2 {
			return convertToR2LangType(strings.TrimSpace(matches[2]))
		}
	}

	// If no specific result pattern, try to extract any numeric value
	numericRegex := regexp.MustCompile(`>(\d+(?:\.\d+)?)<`)
	if matches := numericRegex.FindStringSubmatch(bodyContent); len(matches) > 1 {
		return convertToR2LangType(matches[1])
	}

	return nil
}

// convertToR2LangType converts string to appropriate R2Lang native type
func convertToR2LangType(value string) interface{} {
	value = strings.TrimSpace(value)

	// Try integer
	if i, err := strconv.Atoi(value); err == nil {
		return float64(i) // R2Lang uses float64 for numbers
	}

	// Try float
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}

	// Try boolean
	switch strings.ToLower(value) {
	case "true":
		return true
	case "false":
		return false
	}

	// Return as string
	return value
}

// tryConvertValue attempts to convert string to appropriate type
func tryConvertValue(value string) interface{} {
	return convertToR2LangType(value)
}

// cleanResponseString removes binary/control characters from response
func cleanResponseString(response string) string {
	// Try to extract clean XML from the response
	xmlStart := strings.Index(response, "<?xml")
	if xmlStart == -1 {
		xmlStart = strings.Index(response, "<soap:Envelope")
	}
	if xmlStart == -1 {
		xmlStart = strings.Index(response, "<Envelope")
	}

	if xmlStart > 0 {
		// Remove everything before the XML declaration
		response = response[xmlStart:]
	}

	// Find the end of the XML
	xmlEnd := strings.LastIndex(response, "</soap:Envelope>")
	if xmlEnd == -1 {
		xmlEnd = strings.LastIndex(response, "</Envelope>")
	}

	if xmlEnd > 0 {
		// Keep only up to the end of the XML
		response = response[:xmlEnd+len("</soap:Envelope>")]
	}

	// Filter out non-printable characters
	var cleaned strings.Builder
	for _, r := range response {
		// Keep XML-safe characters: printable ASCII, newlines, tabs
		if (r >= 32 && r <= 126) || r == '\n' || r == '\r' || r == '\t' {
			cleaned.WriteRune(r)
		}
	}

	return cleaned.String()
}

// cleanXMLNamespaces removes namespace prefixes to simplify parsing
func cleanXMLNamespaces(xmlData string) string {
	// Remove namespace prefixes like "soap:" or "tns:"
	nsRegex := regexp.MustCompile(`\w+:`)
	cleaned := nsRegex.ReplaceAllString(xmlData, "")

	// Remove namespace declarations
	nsDeclarationRegex := regexp.MustCompile(`\s+xmlns[^=]*="[^"]*"`)
	cleaned = nsDeclarationRegex.ReplaceAllString(cleaned, "")

	return cleaned
}
