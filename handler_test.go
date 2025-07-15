package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandleRequest_MethodNotAllowed(t *testing.T) {
	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/", nil)
			w := httptest.NewRecorder()

			HandleRequest(w, req)

			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
			}

			if w.Body.String() != "Method not allowed\n" {
				t.Errorf("Expected 'Method not allowed\\n', got '%s'", w.Body.String())
			}
		})
	}
}

func TestHandleRequest_GetSuccess(t *testing.T) {
	// Reset counter to ensure clean state
	counter = NewCounter()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w := httptest.NewRecorder()

	HandleRequest(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	var info Info
	if err := json.Unmarshal(w.Body.Bytes(), &info); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if info.Bucket != "192.168.1.1" {
		t.Errorf("Expected bucket '192.168.1.1', got '%s'", info.Bucket)
	}
	if !info.Allowed {
		t.Error("Expected allowed=true for first request")
	}
	if info.BucketSize != 2 {
		t.Errorf("Expected BucketSize=2, got %d", info.BucketSize)
	}
	if info.Remaining != 1 {
		t.Errorf("Expected Remaining=1, got %d", info.Remaining)
	}
}

func TestHandleRequest_RemoteAddrWithoutPort(t *testing.T) {
	counter = NewCounter()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.1"
	w := httptest.NewRecorder()

	HandleRequest(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var info Info
	if err := json.Unmarshal(w.Body.Bytes(), &info); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if info.Bucket != "192.168.1.1" {
		t.Errorf("Expected bucket '192.168.1.1', got '%s'", info.Bucket)
	}
}

func TestHandleRequest_RateLimitExceeded(t *testing.T) {
	counter = NewCounter()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.2:12345"

	// Make requests up to the limit
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		HandleRequest(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d should succeed, got status %d", i+1, w.Code)
		}
	}

	// Next request should be rate limited
	w := httptest.NewRecorder()
	HandleRequest(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status %d, got %d", http.StatusTooManyRequests, w.Code)
	}

	var info Info
	if err := json.Unmarshal(w.Body.Bytes(), &info); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if info.Allowed {
		t.Error("Expected allowed=false for rate limited request")
	}
	if info.Remaining != 0 {
		t.Errorf("Expected Remaining=0, got %d", info.Remaining)
	}
}

func TestHandleRequest_MultipleClients(t *testing.T) {
	counter = NewCounter()

	client1 := httptest.NewRequest(http.MethodGet, "/", nil)
	client1.RemoteAddr = "192.168.1.3:12345"

	client2 := httptest.NewRequest(http.MethodGet, "/", nil)
	client2.RemoteAddr = "192.168.1.4:12345"

	// Both clients should be able to make requests independently
	w1 := httptest.NewRecorder()
	HandleRequest(w1, client1)

	w2 := httptest.NewRecorder()
	HandleRequest(w2, client2)

	if w1.Code != http.StatusOK {
		t.Errorf("Client 1 should succeed, got status %d", w1.Code)
	}
	if w2.Code != http.StatusOK {
		t.Errorf("Client 2 should succeed, got status %d", w2.Code)
	}

	var info1, info2 Info
	json.Unmarshal(w1.Body.Bytes(), &info1)
	json.Unmarshal(w2.Body.Bytes(), &info2)

	if info1.Bucket == info2.Bucket {
		t.Error("Different clients should have different buckets")
	}
}

func TestHandleRequest_JSONResponse(t *testing.T) {
	counter = NewCounter()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.5:12345"
	w := httptest.NewRecorder()

	HandleRequest(w, req)

	// Verify JSON is properly formatted
	var info Info
	if err := json.Unmarshal(w.Body.Bytes(), &info); err != nil {
		t.Fatalf("Response should be valid JSON: %v", err)
	}

	// Verify all expected fields are present
	if info.Bucket == "" {
		t.Error("Bucket field should not be empty")
	}
	if info.ResetAt.IsZero() {
		t.Error("ResetAt field should not be zero")
	}
	if info.BucketSize == 0 {
		t.Error("BucketSize field should not be zero")
	}

	// Verify JSON is indented (pretty printed)
	prettyJSON, _ := json.MarshalIndent(info, "", "  ")
	if w.Body.String() != string(prettyJSON) {
		t.Error("Response should be pretty-printed JSON")
	}
}

func TestHandleRequest_IPv6Address(t *testing.T) {
	counter = NewCounter()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "[::1]:12345"
	w := httptest.NewRecorder()

	HandleRequest(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var info Info
	if err := json.Unmarshal(w.Body.Bytes(), &info); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if info.Bucket != "::1" {
		t.Errorf("Expected bucket '::1', got '%s'", info.Bucket)
	}
}

func TestHandleRequest_ResetAtField(t *testing.T) {
	counter = NewCounter()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.6:12345"

	// Fill the bucket first
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		HandleRequest(w, req)
	}

	// Next request should be rate limited with ResetAt in the future
	w := httptest.NewRecorder()
	HandleRequest(w, req)

	var info Info
	json.Unmarshal(w.Body.Bytes(), &info)

	now := time.Now()
	if info.ResetAt.Before(now) {
		t.Error("ResetAt should be in the future for rate limited requests")
	}
}
