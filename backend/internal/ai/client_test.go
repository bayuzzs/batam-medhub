package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_IsConfigured(t *testing.T) {
	c1 := NewClient("", "", "", "", 0)
	if c1.IsConfigured() {
		t.Errorf("expected c1 not configured")
	}

	c2 := NewClient("https://api.cloudflare.com", "acc-123", "tok-456", "@cf/meta/llama-3.1-8b-instruct", 5*time.Second)
	if !c2.IsConfigured() {
		t.Errorf("expected c2 configured")
	}
}

func TestClient_InferSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing or invalid authorization header")
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"errors":  []any{},
			"result": map[string]any{
				"response": `{"resolution":"MATCHED","service_code":"MCU_BASIC"}`,
			},
		})
	}))
	defer ts.Close()

	client := NewClient(ts.URL, "test-account", "test-token", "test-model", 5*time.Second)
	res, err := client.Infer(context.Background(), []ChatMessage{
		{Role: "user", Content: "hello"},
	})
	if err != nil {
		t.Fatalf("unexpected infer error: %v", err)
	}

	if res != `{"resolution":"MATCHED","service_code":"MCU_BASIC"}` {
		t.Errorf("unexpected infer response: %s", res)
	}
}

func TestClient_InferErrorResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"errors": []map[string]any{
				{"code": 10000, "message": "Authentication error"},
			},
		})
	}))
	defer ts.Close()

	client := NewClient(ts.URL, "test-account", "invalid-token", "test-model", 5*time.Second)
	_, err := client.Infer(context.Background(), []ChatMessage{
		{Role: "user", Content: "hello"},
	})
	if err == nil {
		t.Fatalf("expected error on 401 response")
	}
}
