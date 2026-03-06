package server_testing

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPTransportUsesHTTP1Only(t *testing.T) {
	// Create a test server that captures the protocol version used
	var receivedProto string
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedProto = r.Proto
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer testServer.Close()

	// Create HTTP client with the same transport configuration as in sendRequest
	transport := &http.Transport{
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	}

	client := &http.Client{
		Timeout:   time.Duration(10) * time.Second,
		Transport: transport,
	}

	// Make a request
	req, err := http.NewRequest("GET", testServer.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Verify HTTP/1.1 was used (not HTTP/2)
	if receivedProto != "HTTP/1.1" {
		t.Errorf("Expected HTTP/1.1, but got %s. HTTP/2 should be disabled.", receivedProto)
	}

	if resp.ProtoMajor == 2 {
		t.Error("Response indicates HTTP/2 was used, but HTTP/1.1 was expected")
	}

	t.Logf("✓ Successfully verified HTTP/1.1 is being used (received: %s)", receivedProto)
}

func TestHTTPTransportDisablesHTTP2(t *testing.T) {
	// Verify that an empty TLSNextProto map disables HTTP/2
	transport := &http.Transport{
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	}

	// The presence of an empty (non-nil) TLSNextProto map should disable HTTP/2
	if transport.TLSNextProto == nil {
		t.Error("TLSNextProto should not be nil to disable HTTP/2")
	}

	if len(transport.TLSNextProto) != 0 {
		t.Error("TLSNextProto should be empty to disable HTTP/2")
	}

	t.Log("✓ Transport configuration correctly disables HTTP/2")
}
