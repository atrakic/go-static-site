package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestHealthHandlerGet(t *testing.T) {
	originalSHA := SHA
	SHA = "test-sha"
	t.Cleanup(func() { SHA = originalSHA })

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	http.HandlerFunc(healthHandler).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if got := rr.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Fatalf("unexpected content type: got %q", got)
	}

	if got := rr.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("unexpected cache control header: got %q", got)
	}

	var response healthResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode health response: %v", err)
	}

	if response.SHA != SHA {
		t.Fatalf("unexpected sha: got %q want %q", response.SHA, SHA)
	}

	if _, err := time.Parse(time.RFC3339, response.Time); err != nil {
		t.Fatalf("unexpected time value %q: %v", response.Time, err)
	}
}

func TestHealthHandlerHead(t *testing.T) {
	req := httptest.NewRequest(http.MethodHead, "/health", nil)
	rr := httptest.NewRecorder()

	http.HandlerFunc(healthHandler).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if rr.Body.Len() != 0 {
		t.Fatalf("expected empty body for HEAD request, got %q", rr.Body.String())
	}
}

func TestHealthHandlerMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	rr := httptest.NewRecorder()

	http.HandlerFunc(healthHandler).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}

	if got := rr.Header().Get("Allow"); got != "GET, HEAD" {
		t.Fatalf("unexpected allow header: got %q", got)
	}
}

func TestServerAddress(t *testing.T) {
	originalPort := os.Getenv("PORT")
	t.Cleanup(func() {
		if originalPort == "" {
			if err := os.Unsetenv("PORT"); err != nil {
				t.Fatalf("failed to unset PORT during cleanup: %v", err)
			}
			return
		}
		if err := os.Setenv("PORT", originalPort); err != nil {
			t.Fatalf("failed to restore PORT during cleanup: %v", err)
		}
	})

	tests := []struct {
		name string
		port string
		want string
	}{
		{name: "default", want: ":8080"},
		{name: "numeric port", port: "9090", want: ":9090"},
		{name: "address already prefixed", port: ":9091", want: ":9091"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.port == "" {
				if err := os.Unsetenv("PORT"); err != nil {
					t.Fatalf("failed to unset PORT: %v", err)
				}
			} else {
				if err := os.Setenv("PORT", tt.port); err != nil {
					t.Fatalf("failed to set PORT: %v", err)
				}
			}

			if got := serverAddress(); got != tt.want {
				t.Fatalf("serverAddress() = %q, want %q", got, tt.want)
			}
		})
	}
}
