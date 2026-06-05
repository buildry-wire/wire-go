package wire

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPaymentIntentsCreate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/payment_intents" {
			t.Errorf("got %s %s", r.Method, r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["amount"].(float64) != 50000 {
			t.Errorf("amount = %v", body["amount"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"id": "pi_1", "object": "payment_intent", "amount": 50000, "status": "requires_payment_method"})
	}))
	defer srv.Close()

	c := NewClient("sk_test_123", WithBaseURL(srv.URL))
	pi, err := c.PaymentIntents.Create(context.Background(), &PaymentIntentCreateParams{Amount: 50000})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if pi.ID != "pi_1" || pi.Status != "requires_payment_method" {
		t.Errorf("pi = %+v", pi)
	}
}

func TestPaymentIntentsListAutoPaginates(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		after := r.URL.Query().Get("starting_after")
		w.Header().Set("Content-Type", "application/json")
		if after == "" {
			json.NewEncoder(w).Encode(map[string]any{"object": "list", "has_more": true,
				"data": []map[string]any{{"id": "pi_1", "object": "payment_intent"}}})
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"object": "list", "has_more": false,
			"data": []map[string]any{{"id": "pi_2", "object": "payment_intent"}}})
	}))
	defer srv.Close()

	c := NewClient("sk_test_123", WithBaseURL(srv.URL))
	var ids []string
	it := c.PaymentIntents.List(context.Background(), nil)
	for it.Next() {
		ids = append(ids, it.Current().ID)
	}
	if it.Err() != nil {
		t.Fatalf("iter err: %v", it.Err())
	}
	if len(ids) != 2 || ids[0] != "pi_1" || ids[1] != "pi_2" {
		t.Errorf("ids = %v", ids)
	}
}
