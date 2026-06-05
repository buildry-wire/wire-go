package wire

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChargesRetrieve(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/charges/ch_1" {
			t.Errorf("path = %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"id": "ch_1", "object": "charge", "status": "succeeded", "amount": 50000})
	}))
	defer srv.Close()

	c := NewClient("sk_test_123", WithBaseURL(srv.URL))
	ch, err := c.Charges.Retrieve(context.Background(), "ch_1")
	if err != nil {
		t.Fatalf("Retrieve: %v", err)
	}
	if ch.ID != "ch_1" || ch.Status != "succeeded" {
		t.Errorf("ch = %+v", ch)
	}
}
