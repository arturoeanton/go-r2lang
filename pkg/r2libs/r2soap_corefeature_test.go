package r2libs

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Regression test: parameter values passed to call/callSimple/callRaw were
// interpolated into the generated SOAP envelope without XML-escaping, so
// values containing "&", "<", ">" produced a malformed envelope that real
// SOAP servers (and any strict XML parser) would reject.
func TestSOAPCallOperationEscapesParameterValues(t *testing.T) {
	client := &SOAPClient{
		Namespace: "http://example.com/ns",
		Operations: map[string]*SOAPOperation{
			"Echo": {Name: "Echo", SOAPAction: "urn:Echo"},
		},
		Headers: map[string]string{},
	}

	var captured string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		captured = string(b)
		w.Header().Set("Content-Type", "text/xml")
		fmt.Fprint(w, `<?xml version="1.0"?><soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body><EchoResponse><Result>ok</Result></EchoResponse></soap:Body></soap:Envelope>`)
	}))
	defer server.Close()
	client.ServiceURL = server.URL

	params := map[string]interface{}{
		"name": "Johnson & Johnson <b>bold</b>",
	}
	client.callOperation("Echo", params)

	if captured == "" {
		t.Fatal("expected the client to send a request body")
	}
	var v interface{}
	if err := xml.Unmarshal([]byte(captured), &v); err != nil {
		t.Fatalf("generated envelope is not valid XML: %v\nenvelope:\n%s", err, captured)
	}
	if want := "Johnson &amp; Johnson &lt;b&gt;bold&lt;/b&gt;"; !strings.Contains(captured, want) {
		t.Fatalf("expected escaped value %q in envelope, got:\n%s", want, captured)
	}
}

// Regression test: SOAP 1.1 servers conventionally report faults using HTTP
// 500 with a well-formed <soap:Envelope>...<Fault> body. sendRequest used to
// treat any non-200 status as a bare transport error, discarding the fault
// body before parseSOAPResponseToR2Lang got a chance to extract structured
// fault.code/fault.message details.
func TestSOAPCallOperationSurfacesFaultOnHTTP500(t *testing.T) {
	client := &SOAPClient{
		Namespace: "http://example.com/ns",
		Operations: map[string]*SOAPOperation{
			"Echo": {Name: "Echo", SOAPAction: "urn:Echo"},
		},
		Headers: map[string]string{},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `<?xml version="1.0"?><soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body><soap:Fault><faultcode>soap:Server</faultcode><faultstring>Something broke</faultstring></soap:Fault></soap:Body></soap:Envelope>`)
	}))
	defer server.Close()
	client.ServiceURL = server.URL

	result := client.callOperation("Echo", map[string]interface{}{})

	rm, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map result, got %T: %v", result, result)
	}
	if rm["success"] != false {
		t.Fatalf("expected success=false for fault, got %v", rm["success"])
	}
	fault, ok := rm["fault"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected structured fault map, got %T: %v", rm["fault"], rm["fault"])
	}
	if fault["message"] != "Something broke" {
		t.Fatalf("expected fault message 'Something broke', got %v", fault["message"])
	}
	if fault["code"] != "soap:Server" {
		t.Fatalf("expected fault code 'soap:Server', got %v", fault["code"])
	}
}

// A genuine non-SOAP HTTP error (e.g. a proxy/gateway error page) must still
// surface as a transport error rather than being misparsed as a SOAP
// response.
func TestSOAPCallOperationNonSOAPHTTPErrorStillErrors(t *testing.T) {
	client := &SOAPClient{
		Namespace: "http://example.com/ns",
		Operations: map[string]*SOAPOperation{
			"Echo": {Name: "Echo", SOAPAction: "urn:Echo"},
		},
		Headers: map[string]string{},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprint(w, `<html><body>502 Bad Gateway</body></html>`)
	}))
	defer server.Close()
	client.ServiceURL = server.URL

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected callOperation to panic on a genuine non-SOAP HTTP error")
		}
	}()
	client.callOperation("Echo", map[string]interface{}{})
}
