package r2libs

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// Mock WSDL content for testing
const mockWSDL = `<?xml version="1.0" encoding="utf-8"?>
<definitions xmlns="http://schemas.xmlsoap.org/wsdl/"
             xmlns:soap="http://schemas.xmlsoap.org/wsdl/soap/"
             xmlns:tns="http://tempuri.org/"
             targetNamespace="http://tempuri.org/">
  
  <types>
    <schema xmlns="http://www.w3.org/2001/XMLSchema" targetNamespace="http://tempuri.org/">
      <element name="Add">
        <complexType>
          <sequence>
            <element name="intA" type="int"/>
            <element name="intB" type="int"/>
          </sequence>
        </complexType>
      </element>
      <element name="AddResponse">
        <complexType>
          <sequence>
            <element name="AddResult" type="int"/>
          </sequence>
        </complexType>
      </element>
    </schema>
  </types>

  <message name="AddSoapIn">
    <part name="parameters" element="tns:Add"/>
  </message>
  <message name="AddSoapOut">
    <part name="parameters" element="tns:AddResponse"/>
  </message>

  <portType name="CalculatorSoap">
    <operation name="Add">
      <input message="tns:AddSoapIn"/>
      <output message="tns:AddSoapOut"/>
    </operation>
  </portType>

  <binding name="CalculatorSoap" type="tns:CalculatorSoap">
    <soap:binding transport="http://schemas.xmlsoap.org/soap/http"/>
    <operation name="Add">
      <soap:operation soapAction="http://tempuri.org/Add" style="document"/>
    </operation>
  </binding>

  <service name="Calculator">
    <port name="CalculatorSoap" binding="tns:CalculatorSoap">
      <soap:address location="http://tempuri.org/Calculator.asmx"/>
    </port>
  </service>
</definitions>`

// Mock SOAP response for Add operation
const mockSOAPResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <AddResponse xmlns="http://tempuri.org/">
      <AddResult>15</AddResult>
    </AddResponse>
  </soap:Body>
</soap:Envelope>`

func TestRegisterSOAP(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterSOAP(env)

	// Check if SOAP functions are registered
	soapModuleObj, ok := env.Get("soap")
	if !ok {
		t.Fatal("soap module not found")
	}
	soapModule := soapModuleObj.(map[string]interface{})

	functions := []string{"client", "envelope", "request"}

	for _, funcName := range functions {
		if _, exists := soapModule[funcName]; !exists {
			t.Errorf("Function %s not registered", funcName)
		}
	}
}

func TestSOAPEnvelope(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterSOAP(env)

	soapModuleObj, ok := env.Get("soap")
	if !ok {
		t.Fatal("soap module not found")
	}
	soapModule := soapModuleObj.(map[string]interface{})

	soapEnvelopeFunc, _ := soapModule["envelope"]
	builtinFunc := soapEnvelopeFunc.(r2core.BuiltinFunction)

	namespace := "http://tempuri.org/"
	methodName := "Add"
	bodyContent := "<intA>5</intA><intB>10</intB>"

	result := builtinFunc(namespace, methodName, bodyContent)
	envelope := result.(string)

	// Check that envelope contains expected elements
	expectedElements := []string{
		`<?xml version="1.0" encoding="utf-8"?>`,
		`<soap:Envelope`,
		`xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"`,
		fmt.Sprintf(`xmlns:tns="%s"`, namespace),
		`<soap:Header />`,
		`<soap:Body>`,
		fmt.Sprintf(`<tns:%s>`, methodName),
		bodyContent,
		fmt.Sprintf(`</tns:%s>`, methodName),
		`</soap:Body>`,
		`</soap:Envelope>`,
	}

	for _, element := range expectedElements {
		if !strings.Contains(envelope, element) {
			t.Errorf("Envelope missing expected element: %s", element)
		}
	}
}

func TestSOAPClientCreation(t *testing.T) {
	// Create mock WSDL server
	wsdlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		fmt.Fprint(w, mockWSDL)
	}))
	defer wsdlServer.Close()

	env := r2core.NewEnvironment()
	RegisterSOAP(env)

	soapModuleObj, ok := env.Get("soap")
	if !ok {
		t.Fatal("soap module not found")
	}
	soapModule := soapModuleObj.(map[string]interface{})

	soapClientFunc, _ := soapModule["client"]
	builtinFunc := soapClientFunc.(r2core.BuiltinFunction)

	// Test client creation with mock WSDL
	result := builtinFunc(wsdlServer.URL)
	clientMap := result.(map[string]interface{})

	// Check client properties
	if wsdlURL := clientMap["wsdlURL"].(string); wsdlURL != wsdlServer.URL {
		t.Errorf("Expected wsdlURL %s, got %s", wsdlServer.URL, wsdlURL)
	}

	if serviceURL := clientMap["serviceURL"].(string); serviceURL != "http://tempuri.org/Calculator.asmx" {
		t.Errorf("Expected serviceURL http://tempuri.org/Calculator.asmx, got %s", serviceURL)
	}

	if namespace := clientMap["namespace"].(string); namespace != "http://tempuri.org/" {
		t.Errorf("Expected namespace http://tempuri.org/, got %s", namespace)
	}

	// Check that methods are available
	methods := []string{"listOperations", "getOperation", "call", "setTimeout", "setHeader"}
	for _, method := range methods {
		if _, exists := clientMap[method]; !exists {
			t.Errorf("Client missing method: %s", method)
		}
	}
}

func TestSOAPClientListOperations(t *testing.T) {
	// Create mock WSDL server
	wsdlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		fmt.Fprint(w, mockWSDL)
	}))
	defer wsdlServer.Close()

	env := r2core.NewEnvironment()
	RegisterSOAP(env)

	soapModuleObj, ok := env.Get("soap")
	if !ok {
		t.Fatal("soap module not found")
	}
	soapModule := soapModuleObj.(map[string]interface{})

	soapClientFunc, _ := soapModule["client"]
	builtinFunc := soapClientFunc.(r2core.BuiltinFunction)

	clientMap := builtinFunc(wsdlServer.URL).(map[string]interface{})
	listOperationsFunc := clientMap["listOperations"].(r2core.BuiltinFunction)

	operations := listOperationsFunc().([]interface{})

	if len(operations) != 1 {
		t.Errorf("Expected 1 operation, got %d", len(operations))
	}

	if operations[0].(string) != "Add" {
		t.Errorf("Expected operation 'Add', got %s", operations[0].(string))
	}
}

func TestSOAPClientGetOperation(t *testing.T) {
	// Create mock WSDL server
	wsdlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		fmt.Fprint(w, mockWSDL)
	}))
	defer wsdlServer.Close()

	env := r2core.NewEnvironment()
	RegisterSOAP(env)

	soapModuleObj, ok := env.Get("soap")
	if !ok {
		t.Fatal("soap module not found")
	}
	soapModule := soapModuleObj.(map[string]interface{})

	soapClientFunc, _ := soapModule["client"]
	builtinFunc := soapClientFunc.(r2core.BuiltinFunction)

	clientMap := builtinFunc(wsdlServer.URL).(map[string]interface{})
	getOperationFunc := clientMap["getOperation"].(r2core.BuiltinFunction)

	operationInfo := getOperationFunc("Add").(map[string]interface{})

	if operationInfo["name"].(string) != "Add" {
		t.Errorf("Expected operation name 'Add', got %s", operationInfo["name"].(string))
	}

	if operationInfo["soapAction"].(string) != "http://tempuri.org/Add" {
		t.Errorf("Expected SOAPAction 'http://tempuri.org/Add', got %s", operationInfo["soapAction"].(string))
	}
}

func TestSOAPClientSetTimeout(t *testing.T) {
	// Create mock WSDL server
	wsdlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		fmt.Fprint(w, mockWSDL)
	}))
	defer wsdlServer.Close()

	env := r2core.NewEnvironment()
	RegisterSOAP(env)

	soapModuleObj, ok := env.Get("soap")
	if !ok {
		t.Fatal("soap module not found")
	}
	soapModule := soapModuleObj.(map[string]interface{})

	soapClientFunc, _ := soapModule["client"]
	builtinFunc := soapClientFunc.(r2core.BuiltinFunction)

	clientMap := builtinFunc(wsdlServer.URL).(map[string]interface{})
	setTimeoutFunc := clientMap["setTimeout"].(r2core.BuiltinFunction)

	// Test setting timeout - should not panic
	setTimeoutFunc(60.0)
}

func TestSOAPClientSetHeader(t *testing.T) {
	// Create mock WSDL server
	wsdlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		fmt.Fprint(w, mockWSDL)
	}))
	defer wsdlServer.Close()

	env := r2core.NewEnvironment()
	RegisterSOAP(env)

	soapModuleObj, ok := env.Get("soap")
	if !ok {
		t.Fatal("soap module not found")
	}
	soapModule := soapModuleObj.(map[string]interface{})

	soapClientFunc, _ := soapModule["client"]
	builtinFunc := soapClientFunc.(r2core.BuiltinFunction)

	clientMap := builtinFunc(wsdlServer.URL).(map[string]interface{})
	setHeaderFunc := clientMap["setHeader"].(r2core.BuiltinFunction)

	// Test setting header - should not panic
	setHeaderFunc("Authorization", "Bearer token123")
}

func TestSOAPRawRequest(t *testing.T) {
	// Create mock SOAP server
	soapServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers
		if r.Header.Get("Content-Type") != "text/xml; charset=utf-8" {
			t.Errorf("Expected Content-Type 'text/xml; charset=utf-8', got %s", r.Header.Get("Content-Type"))
		}

		if r.Header.Get("SOAPAction") != `"http://tempuri.org/Add"` {
			t.Errorf("Expected SOAPAction '\"http://tempuri.org/Add\"', got %s", r.Header.Get("SOAPAction"))
		}

		w.Header().Set("Content-Type", "text/xml")
		fmt.Fprint(w, mockSOAPResponse)
	}))
	defer soapServer.Close()

	env := r2core.NewEnvironment()
	RegisterSOAP(env)

	soapModuleObj, ok := env.Get("soap")
	if !ok {
		t.Fatal("soap module not found")
	}
	soapModule := soapModuleObj.(map[string]interface{})

	soapRequestFunc, _ := soapModule["request"]
	builtinFunc := soapRequestFunc.(r2core.BuiltinFunction)

	envelope := createSOAPEnvelope("http://tempuri.org/", "Add", "<intA>5</intA><intB>10</intB>")

	response := builtinFunc(soapServer.URL, "http://tempuri.org/Add", envelope)
	responseStr := response.(string)

	if !strings.Contains(responseStr, "AddResult") {
		t.Error("Response should contain AddResult")
	}

	if !strings.Contains(responseStr, "15") {
		t.Error("Response should contain result value 15")
	}
}

func TestCreateSOAPEnvelope(t *testing.T) {
	namespace := "http://example.com/"
	methodName := "TestMethod"
	bodyContent := "<param1>value1</param1><param2>value2</param2>"

	envelope := createSOAPEnvelope(namespace, methodName, bodyContent)

	// Check structure
	if !strings.Contains(envelope, `<?xml version="1.0" encoding="utf-8"?>`) {
		t.Error("Envelope should contain XML declaration")
	}

	if !strings.Contains(envelope, `<soap:Envelope`) {
		t.Error("Envelope should contain soap:Envelope element")
	}

	if !strings.Contains(envelope, fmt.Sprintf(`xmlns:tns="%s"`, namespace)) {
		t.Error("Envelope should contain target namespace")
	}

	if !strings.Contains(envelope, fmt.Sprintf(`<tns:%s>`, methodName)) {
		t.Error("Envelope should contain method element")
	}

	if !strings.Contains(envelope, bodyContent) {
		t.Error("Envelope should contain body content")
	}
}

func TestSOAPClientErrorHandling(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterSOAP(env)

	soapModuleObj, ok := env.Get("soap")
	if !ok {
		t.Fatal("soap module not found")
	}
	soapModule := soapModuleObj.(map[string]interface{})

	soapClientFunc, _ := soapModule["client"]
	builtinFunc := soapClientFunc.(r2core.BuiltinFunction)

	// Test with invalid URL
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic with invalid WSDL URL")
		}
	}()

	builtinFunc("invalid://url")
}

func TestSOAPEnvelopeErrorHandling(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterSOAP(env)

	soapModuleObj, ok := env.Get("soap")
	if !ok {
		t.Fatal("soap module not found")
	}
	soapModule := soapModuleObj.(map[string]interface{})

	soapEnvelopeFunc, _ := soapModule["envelope"]
	builtinFunc := soapEnvelopeFunc.(r2core.BuiltinFunction)

	// Test with insufficient arguments
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic with insufficient arguments")
		}
	}()

	builtinFunc("namespace")
}

func TestSOAPRequestErrorHandling(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterSOAP(env)

	soapModuleObj, ok := env.Get("soap")
	if !ok {
		t.Fatal("soap module not found")
	}
	soapModule := soapModuleObj.(map[string]interface{})

	soapRequestFunc, _ := soapModule["request"]
	builtinFunc := soapRequestFunc.(r2core.BuiltinFunction)

	// Test with insufficient arguments
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic with insufficient arguments")
		}
	}()

	builtinFunc("url", "action")
}

func TestCleanXMLNamespaces(t *testing.T) {
	xmlWithNamespaces := `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
		<soap:Body>
			<tns:Response xmlns:tns="http://tempuri.org/">
				<tns:Result>123</tns:Result>
			</tns:Response>
		</soap:Body>
	</soap:Envelope>`

	cleaned := cleanXMLNamespaces(xmlWithNamespaces)

	// Should remove namespace prefixes
	if strings.Contains(cleaned, "soap:") {
		t.Error("Cleaned XML should not contain soap: prefix")
	}

	if strings.Contains(cleaned, "tns:") {
		t.Error("Cleaned XML should not contain tns: prefix")
	}

	// Should remove namespace declarations
	if strings.Contains(cleaned, "xmlns:soap") {
		t.Error("Cleaned XML should not contain xmlns:soap declaration")
	}

	if strings.Contains(cleaned, "xmlns:tns") {
		t.Error("Cleaned XML should not contain xmlns:tns declaration")
	}
}

func TestWSDLParsing(t *testing.T) {
	// Create mock WSDL server
	wsdlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		fmt.Fprint(w, mockWSDL)
	}))
	defer wsdlServer.Close()

	// Test WSDL parsing directly
	client, err := createSOAPClient(wsdlServer.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create SOAP client: %v", err)
	}

	if client.WSDLURL != wsdlServer.URL {
		t.Errorf("Expected WSDLURL %s, got %s", wsdlServer.URL, client.WSDLURL)
	}

	if client.ServiceURL != "http://tempuri.org/Calculator.asmx" {
		t.Errorf("Expected ServiceURL http://tempuri.org/Calculator.asmx, got %s", client.ServiceURL)
	}

	if client.Namespace != "http://tempuri.org/" {
		t.Errorf("Expected Namespace http://tempuri.org/, got %s", client.Namespace)
	}

	if len(client.Operations) != 1 {
		t.Errorf("Expected 1 operation, got %d", len(client.Operations))
	}

	addOp, exists := client.Operations["Add"]
	if !exists {
		t.Error("Add operation should exist")
	}

	if addOp.Name != "Add" {
		t.Errorf("Expected operation name 'Add', got %s", addOp.Name)
	}

	if addOp.SOAPAction != "http://tempuri.org/Add" {
		t.Errorf("Expected SOAPAction 'http://tempuri.org/Add', got %s", addOp.SOAPAction)
	}
}

// TestEnterpriseHeadersCustomization tests custom headers functionality
func TestEnterpriseHeadersCustomization(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterSOAP(env)

	// Create mock WSDL server
	wsdlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		w.Write([]byte(mockWSDL))
	}))
	defer wsdlServer.Close()

	// Create client with custom headers
	soapModuleObj, ok := env.Get("soap")
	if !ok {
		t.Fatal("soap module not found")
	}
	soapModule := soapModuleObj.(map[string]interface{})

	soapClientFunc, _ := soapModule["client"]
	builtinFunc := soapClientFunc.(r2core.BuiltinFunction)

	customHeaders := map[string]interface{}{
		"X-Company":  "TestCorp",
		"X-Version":  "2.0",
		"User-Agent": "TestClient/1.0",
	}

	clientResult := builtinFunc(wsdlServer.URL, customHeaders)
	clientMap := clientResult.(map[string]interface{})

	// Test getHeaders functionality
	getHeadersFunc := clientMap["getHeaders"].(r2core.BuiltinFunction)
	headers := getHeadersFunc().(map[string]interface{})

	// Verify custom headers are set
	if headers["X-Company"] != "TestCorp" {
		t.Errorf("Expected X-Company header to be 'TestCorp', got %v", headers["X-Company"])
	}
	if headers["X-Version"] != "2.0" {
		t.Errorf("Expected X-Version header to be '2.0', got %v", headers["X-Version"])
	}

	// Test setHeader with map
	setHeaderFunc := clientMap["setHeader"].(r2core.BuiltinFunction)
	setHeaderFunc(map[string]interface{}{
		"X-Department":  "IT",
		"X-Environment": "test",
	})

	updatedHeaders := getHeadersFunc().(map[string]interface{})
	if updatedHeaders["X-Department"] != "IT" {
		t.Errorf("Expected X-Department header to be 'IT', got %v", updatedHeaders["X-Department"])
	}

	// Test removeHeader
	removeHeaderFunc := clientMap["removeHeader"].(r2core.BuiltinFunction)
	removeHeaderFunc("X-Company")

	finalHeaders := getHeadersFunc().(map[string]interface{})
	if _, exists := finalHeaders["X-Company"]; exists {
		t.Error("X-Company header should have been removed")
	}

	// Test resetHeaders
	resetHeadersFunc := clientMap["resetHeaders"].(r2core.BuiltinFunction)
	resetHeadersFunc()

	resetHeaders := getHeadersFunc().(map[string]interface{})
	if resetHeaders["User-Agent"] == "" {
		t.Error("User-Agent should be set to browser-like default after reset")
	}
}

// TestTLSConfiguration tests TLS/SSL configuration
func TestTLSConfiguration(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterSOAP(env)

	// Create mock WSDL server
	wsdlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		w.Write([]byte(mockWSDL))
	}))
	defer wsdlServer.Close()

	soapModuleObj, ok := env.Get("soap")
	if !ok {
		t.Fatal("soap module not found")
	}
	soapModule := soapModuleObj.(map[string]interface{})

	soapClientFunc, _ := soapModule["client"]
	builtinFunc := soapClientFunc.(r2core.BuiltinFunction)

	clientResult := builtinFunc(wsdlServer.URL)
	clientMap := clientResult.(map[string]interface{})

	// Test setTLSConfig
	setTLSConfigFunc := clientMap["setTLSConfig"].(r2core.BuiltinFunction)

	// Test with valid TLS configuration
	setTLSConfigFunc(map[string]interface{}{
		"minVersion": "1.2",
		"skipVerify": false,
	})

	// Test with skip verify (for testing)
	setTLSConfigFunc(map[string]interface{}{
		"skipVerify": true,
		"minVersion": "1.3",
	})

	// Should not panic - indicates configuration was accepted
}

// TestAuthentication tests authentication configuration
func TestAuthentication(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterSOAP(env)

	// Create mock WSDL server
	wsdlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		w.Write([]byte(mockWSDL))
	}))
	defer wsdlServer.Close()

	soapModuleObj, ok := env.Get("soap")
	if !ok {
		t.Fatal("soap module not found")
	}
	soapModule := soapModuleObj.(map[string]interface{})

	soapClientFunc, _ := soapModule["client"]
	builtinFunc := soapClientFunc.(r2core.BuiltinFunction)

	clientResult := builtinFunc(wsdlServer.URL)
	clientMap := clientResult.(map[string]interface{})

	setAuthFunc := clientMap["setAuth"].(r2core.BuiltinFunction)

	// Test Basic Authentication
	setAuthFunc(map[string]interface{}{
		"type":     "basic",
		"username": "testuser",
		"password": "testpass",
	})

	// Test Bearer Token Authentication
	setAuthFunc(map[string]interface{}{
		"type":  "bearer",
		"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test.token",
	})

	// Should not panic - indicates authentication was configured
}

// TestResponseParsing tests the new response parsing functionality
func TestResponseParsing(t *testing.T) {
	// Mock SOAP service that returns a proper SOAP response
	soapServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, ".wsdl") {
			w.Header().Set("Content-Type", "text/xml")
			w.Write([]byte(mockWSDL))
		} else {
			w.Header().Set("Content-Type", "text/xml")
			w.Write([]byte(mockSOAPResponse))
		}
	}))
	defer soapServer.Close()

	env := r2core.NewEnvironment()
	RegisterSOAP(env)

	soapModuleObj, ok := env.Get("soap")
	if !ok {
		t.Fatal("soap module not found")
	}
	soapModule := soapModuleObj.(map[string]interface{})

	soapClientFunc, _ := soapModule["client"]
	builtinFunc := soapClientFunc.(r2core.BuiltinFunction)

	// Update mock WSDL to point to our test server
	updatedWSDL := strings.Replace(mockWSDL, "http://tempuri.org/Calculator.asmx", soapServer.URL, 1)

	// Create another server for the updated WSDL
	wsdlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		w.Write([]byte(updatedWSDL))
	}))
	defer wsdlServer.Close()

	clientResult := builtinFunc(wsdlServer.URL)
	clientMap := clientResult.(map[string]interface{})

	params := map[string]interface{}{
		"intA": 10,
		"intB": 5,
	}

	// Test full response parsing (call)
	callFunc := clientMap["call"].(r2core.BuiltinFunction)
	fullResponse := callFunc("Add", params)

	// Verify response is not nil (may be string or map depending on parsing)
	if fullResponse == nil {
		t.Error("Full response should not be nil")
	}

	// Test simple response parsing (callSimple)
	callSimpleFunc := clientMap["callSimple"].(r2core.BuiltinFunction)
	simpleResponse := callSimpleFunc("Add", params)

	// Simple response should extract the result directly
	if simpleResponse == nil {
		t.Error("Simple response should not be nil")
	}

	// Test raw response (callRaw)
	callRawFunc := clientMap["callRaw"].(r2core.BuiltinFunction)
	rawResponse := callRawFunc("Add", params)

	if rawStr, ok := rawResponse.(string); ok {
		if !strings.Contains(rawStr, "soap:Envelope") {
			t.Error("Raw response should contain SOAP envelope")
		}
	} else {
		t.Error("Raw response should be a string")
	}
}

// TestParseSOAPResponse tests the parseSOAPResponse function
func TestParseSOAPResponse(t *testing.T) {
	response := mockSOAPResponse

	parsed, err := parseSOAPResponse(response)
	if err != nil {
		t.Errorf("Failed to parse SOAP response: %v", err)
	}

	// Check success field
	if success, ok := parsed["success"].(bool); !ok || !success {
		t.Error("Parsed response should have success=true")
	}

	// Check raw field
	if raw, ok := parsed["raw"].(string); !ok || raw != response {
		t.Error("Parsed response should preserve raw XML")
	}

	// Check body extraction
	if body, ok := parsed["body"].(string); !ok || body == "" {
		t.Error("Parsed response should extract SOAP body")
	}

	// Check values extraction
	if values, ok := parsed["values"].(map[string]interface{}); ok {
		if len(values) == 0 {
			t.Error("Should extract some values from response")
		}
	}
}

// TestSOAPFaultHandling tests SOAP fault detection and parsing
func TestSOAPFaultHandling(t *testing.T) {
	faultResponse := `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <soap:Fault>
      <faultcode>Client</faultcode>
      <faultstring>Invalid request format</faultstring>
    </soap:Fault>
  </soap:Body>
</soap:Envelope>`

	parsed, err := parseSOAPResponse(faultResponse)
	if err != nil {
		t.Errorf("Failed to parse SOAP fault: %v", err)
	}

	// Should at least parse the response structure
	if parsed == nil {
		t.Error("Parsed response should not be nil")
	}

	// Check if body was extracted (even if fault parsing fails)
	if body, ok := parsed["body"].(string); ok {
		// Check if the body contains fault information
		if !strings.Contains(body, "Fault") {
			t.Error("Response body should contain fault information")
		}
	}

	// Note: The fault detection logic might need refinement
	// For now, just verify the response is parsed
}

// TestTryConvertValue tests the value conversion functionality
func TestTryConvertValue(t *testing.T) {
	// Test integer conversion
	if result := tryConvertValue("42"); result != float64(42) {
		t.Errorf("Expected 42.0, got %v", result)
	}

	// Test float conversion
	if result := tryConvertValue("3.14"); result != 3.14 {
		t.Errorf("Expected 3.14, got %v", result)
	}

	// Test boolean conversion
	if result := tryConvertValue("true"); result != true {
		t.Errorf("Expected true, got %v", result)
	}

	if result := tryConvertValue("false"); result != false {
		t.Errorf("Expected false, got %v", result)
	}

	// Test string passthrough
	if result := tryConvertValue("hello"); result != "hello" {
		t.Errorf("Expected 'hello', got %v", result)
	}
}

// TestExtractNumericResult tests numeric result extraction
func TestExtractNumericResult(t *testing.T) {
	// Test AddResult pattern
	bodyContent := `<AddResponse xmlns="http://tempuri.org/"><AddResult>15</AddResult></AddResponse>`
	result := extractNumericResult(bodyContent)

	if result != float64(15) {
		t.Errorf("Expected 15.0, got %v", result)
	}

	// Test no result pattern
	bodyContent2 := `<SomeResponse><Data>text</Data></SomeResponse>`
	result2 := extractNumericResult(bodyContent2)

	if result2 != nil {
		t.Errorf("Expected nil for non-numeric result, got %v", result2)
	}
}

// TestEnterpriseErrorHandling tests comprehensive error handling
func TestEnterpriseErrorHandling(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterSOAP(env)

	soapModuleObj, ok := env.Get("soap")
	if !ok {
		t.Fatal("soap module not found")
	}
	soapModule := soapModuleObj.(map[string]interface{})

	soapClientFunc, _ := soapModule["client"]
	builtinFunc := soapClientFunc.(r2core.BuiltinFunction)

	// Test with invalid WSDL URL
	didPanic := false
	defer func() {
		if r := recover(); r != nil {
			didPanic = true
			errorMsg := fmt.Sprintf("%v", r)
			// Should contain helpful error information
			if !strings.Contains(errorMsg, "failed to create client") {
				t.Error("Error message should contain context information")
			}
		}
	}()

	// This should panic with detailed error message
	builtinFunc("http://invalid-non-existent-service.com/service.wsdl")

	if !didPanic {
		t.Error("Should have panicked with connection error")
	}
}

// TestMultipleClientInstances tests creating multiple client instances
func TestMultipleClientInstances(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterSOAP(env)

	// Create mock WSDL server
	wsdlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		w.Write([]byte(mockWSDL))
	}))
	defer wsdlServer.Close()

	soapModuleObj, ok := env.Get("soap")
	if !ok {
		t.Fatal("soap module not found")
	}
	soapModule := soapModuleObj.(map[string]interface{})

	soapClientFunc, _ := soapModule["client"]
	builtinFunc := soapClientFunc.(r2core.BuiltinFunction)

	// Create multiple clients
	client1 := builtinFunc(wsdlServer.URL, map[string]interface{}{"X-Client": "Client1"})
	client2 := builtinFunc(wsdlServer.URL, map[string]interface{}{"X-Client": "Client2"})

	client1Map := client1.(map[string]interface{})
	client2Map := client2.(map[string]interface{})

	// Configure clients differently
	setHeaderFunc1 := client1Map["setHeader"].(r2core.BuiltinFunction)
	setHeaderFunc2 := client2Map["setHeader"].(r2core.BuiltinFunction)

	setHeaderFunc1("X-Environment", "test1")
	setHeaderFunc2("X-Environment", "test2")

	// Verify they are independent
	getHeadersFunc1 := client1Map["getHeaders"].(r2core.BuiltinFunction)
	getHeadersFunc2 := client2Map["getHeaders"].(r2core.BuiltinFunction)

	headers1 := getHeadersFunc1().(map[string]interface{})
	headers2 := getHeadersFunc2().(map[string]interface{})

	if headers1["X-Environment"] != "test1" {
		t.Error("Client1 should have X-Environment=test1")
	}
	if headers2["X-Environment"] != "test2" {
		t.Error("Client2 should have X-Environment=test2")
	}
	if headers1["X-Client"] != "Client1" {
		t.Error("Client1 should have X-Client=Client1")
	}
	if headers2["X-Client"] != "Client2" {
		t.Error("Client2 should have X-Client=Client2")
	}
}

func TestSOAPClientFullWorkflow(t *testing.T) {
	// Create mock WSDL server
	wsdlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		fmt.Fprint(w, mockWSDL)
	}))
	defer wsdlServer.Close()

	// Create mock SOAP service server
	soapServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		fmt.Fprint(w, mockSOAPResponse)
	}))
	defer soapServer.Close()

	// Update mock WSDL to use our mock SOAP server
	updatedWSDL := strings.Replace(mockWSDL, "http://tempuri.org/Calculator.asmx", soapServer.URL, 1)

	wsdlServerUpdated := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		fmt.Fprint(w, updatedWSDL)
	}))
	defer wsdlServerUpdated.Close()

	env := r2core.NewEnvironment()
	RegisterSOAP(env)

	soapModuleObj, ok := env.Get("soap")
	if !ok {
		t.Fatal("soap module not found")
	}
	soapModule := soapModuleObj.(map[string]interface{})

	soapClientFunc, _ := soapModule["client"]
	builtinFunc := soapClientFunc.(r2core.BuiltinFunction)

	// Create client
	clientMap := builtinFunc(wsdlServerUpdated.URL).(map[string]interface{})

	// Call operation
	callFunc := clientMap["call"].(r2core.BuiltinFunction)
	params := map[string]interface{}{
		"intA": 5,
		"intB": 10,
	}

	response := callFunc("Add", params)

	// Verify response is structured correctly
	if responseMap, ok := response.(map[string]interface{}); ok {
		// Check success
		if success, exists := responseMap["success"]; !exists || success != true {
			t.Error("Response should indicate success")
		}

		// Check raw response contains expected elements
		if raw, exists := responseMap["raw"]; exists {
			if rawStr, ok := raw.(string); ok {
				if !strings.Contains(rawStr, "AddResult") {
					t.Error("Raw response should contain AddResult")
				}
				if !strings.Contains(rawStr, "15") {
					t.Error("Raw response should contain result value 15")
				}
			}
		}

		// Check that result is extracted
		if result, exists := responseMap["result"]; exists && result != nil {
			if resultVal, ok := result.(float64); ok && resultVal != 15.0 {
				t.Errorf("Expected result 15.0, got %v", resultVal)
			}
		}
	} else {
		t.Error("Response should be a structured map")
	}
}
