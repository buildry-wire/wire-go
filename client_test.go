package wire

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestClientAuthAndDecode(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer sk_test_123" {
			t.Errorf("Authorization = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "pi_1", "object": "payment_intent", "amount": 50000})
	}))
	defer srv.Close()

	c := NewClient("sk_test_123", WithBaseURL(srv.URL))
	var out PaymentIntent
	if err := c.do(context.Background(), http.MethodGet, "/v1/payment_intents/pi_1", nil, nil, "", &out); err != nil {
		t.Fatalf("do: %v", err)
	}
	if out.ID != "pi_1" || out.Amount != 50000 {
		t.Errorf("decoded = %+v", out)
	}
}

func TestClientRetriesOn503(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "pi_1", "object": "payment_intent"})
	}))
	defer srv.Close()

	c := NewClient("sk_test_123", WithBaseURL(srv.URL), WithMaxRetries(3), WithBackoff(time.Millisecond))
	var out PaymentIntent
	if err := c.do(context.Background(), http.MethodGet, "/v1/payment_intents/pi_1", nil, nil, "", &out); err != nil {
		t.Fatalf("do: %v", err)
	}
	if calls != 3 {
		t.Errorf("calls = %d, want 3", calls)
	}
}

func TestClientNoRetryOn400(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": map[string]string{"type": "invalid_request_error", "message": "bad"}})
	}))
	defer srv.Close()

	c := NewClient("sk_test_123", WithBaseURL(srv.URL), WithMaxRetries(3), WithBackoff(time.Millisecond))
	err := c.do(context.Background(), http.MethodGet, "/v1/x", nil, nil, "", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if calls != 1 {
		t.Errorf("calls = %d, want 1 (no retry on 4xx)", calls)
	}
}

func TestClientSendsIdempotencyKeyOnPost(t *testing.T) {
	var key string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key = r.Header.Get("Idempotency-Key")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "pi_1", "object": "payment_intent"})
	}))
	defer srv.Close()

	c := NewClient("sk_test_123", WithBaseURL(srv.URL))
	var out PaymentIntent
	if err := c.do(context.Background(), http.MethodPost, "/v1/payment_intents", map[string]any{"amount": 1}, nil, "", &out); err != nil {
		t.Fatalf("do: %v", err)
	}
	if key == "" {
		t.Error("auto Idempotency-Key not sent on POST")
	}
}
