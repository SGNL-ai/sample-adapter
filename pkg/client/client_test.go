// Copyright 2025 SGNL.ai, Inc.
package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCustomUserAgent(t *testing.T) {
	// Expected User-Agent value
	tests := map[string]struct {
		inputUserAgent string
		wantUserAgent  string
	}{
		"empty input": {
			inputUserAgent: "",
			wantUserAgent:  "sgnl-adapter",
		},
		"non-empty input": {
			inputUserAgent: "sgnl-myAdapter/1.0.0",
			wantUserAgent:  "sgnl-myAdapter/1.0.0",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Create a test server to capture the request
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check the User-Agent header
				if gotUserAgent := r.Header.Get("User-Agent"); gotUserAgent != tc.wantUserAgent {
					w.WriteHeader(http.StatusBadRequest)
					t.Fatalf("unexpected User-Agent header: got %q, want %q", gotUserAgent, tc.wantUserAgent)

					return
				}

				w.WriteHeader(http.StatusOK)
			}))
			defer testServer.Close()

			// Create an HTTP client with the custom User-Agent
			client := NewSGNLHttpClient(time.Second, tc.inputUserAgent)

			// Fire a request
			resp, err := client.Get(testServer.URL)
			if err != nil {
				t.Fatalf("failed to make request: %v", err)
			}
			defer resp.Body.Close()

			// Ensure the response status is OK
			if resp.StatusCode != http.StatusOK {
				t.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK)
			}
		})
	}
}
