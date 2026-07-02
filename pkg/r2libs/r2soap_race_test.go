package r2libs

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// TestSOAPTLSConfigConcurrentAccess adversarially hammers setTLSConfig and the
// request-sending path (call) concurrently to look for data races on the
// shared *tls.Config pointer / other mutable client fields. Run with -race.
func TestSOAPTLSConfigConcurrentAccess(t *testing.T) {
	// Mock SOAP service (also serves as its own "WSDL" service endpoint target)
	var reqCount int64
	soapServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&reqCount, 1)
		w.Header().Set("Content-Type", "text/xml")
		fmt.Fprint(w, mockSOAPResponse)
	}))
	defer soapServer.Close()

	updatedWSDL := fmt_replace(mockWSDL, "http://tempuri.org/Calculator.asmx", soapServer.URL)

	wsdlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		fmt.Fprint(w, updatedWSDL)
	}))
	defer wsdlServer.Close()

	env := r2core.NewEnvironment()
	RegisterSOAP(env)

	soapModuleObj, ok := env.Get("soap")
	if !ok {
		t.Fatal("soap module not found")
	}
	soapModule := soapModuleObj.(map[string]interface{})
	soapClientFunc := soapModule["client"].(r2core.BuiltinFunction)

	clientMap := soapClientFunc(wsdlServer.URL).(map[string]interface{})
	setTLSConfigFunc := clientMap["setTLSConfig"].(r2core.BuiltinFunction)
	callFunc := clientMap["call"].(r2core.BuiltinFunction)
	setHeaderFunc := clientMap["setHeader"].(r2core.BuiltinFunction)
	setTimeoutFunc := clientMap["setTimeout"].(r2core.BuiltinFunction)

	const duration = 400 * time.Millisecond
	deadline := time.Now().Add(duration)

	var wg sync.WaitGroup

	// Writers: repeatedly mutate TLS config (and other mutable fields) with
	// varying values so field writes interleave across goroutines.
	for w := 0; w < 4; w++ {
		w := w
		wg.Add(1)
		go func() {
			defer wg.Done()
			i := 0
			versions := []string{"1.0", "1.1", "1.2", "1.3"}
			for time.Now().Before(deadline) {
				setTLSConfigFunc(map[string]interface{}{
					"skipVerify": i%2 == 0,
					"minVersion": versions[i%len(versions)],
				})
				setHeaderFunc(fmt.Sprintf("X-Writer-%d", w), fmt.Sprintf("%d", i))
				setTimeoutFunc(float64(5 + i%10))
				i++
			}
		}()
	}

	// Readers: repeatedly issue SOAP requests which clone client.TLSConfig
	// under lock and then perform the HTTP round trip.
	for r := 0; r < 8; r++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			params := map[string]interface{}{"intA": 1, "intB": 2}
			for time.Now().Before(deadline) {
				func() {
					defer func() {
						// callOperation panics on transport errors; that's
						// fine for this race test, we only care about -race
						// reports, not business-logic correctness here.
						recover()
					}()
					callFunc("Add", params)
				}()
			}
		}()
	}

	wg.Wait()

	if atomic.LoadInt64(&reqCount) == 0 {
		t.Error("expected at least one SOAP request to reach the server")
	}
}

// fmt_replace is a tiny helper to avoid importing strings twice awkwardly in
// this file (kept local to this test file).
func fmt_replace(s, old, new string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			out = append(out, new...)
			i += len(old)
		} else {
			out = append(out, s[i])
			i++
		}
	}
	return string(out)
}
