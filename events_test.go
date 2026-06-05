package wire

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEventsRetrieve(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/events/evt_1" {
			t.Errorf("path = %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"id": "evt_1", "object": "event", "type": "payment_intent.succeeded"})
	}))
	defer srv.Close()

	c := NewClient("sk_test_123", WithBaseURL(srv.URL))
	ev, err := c.Events.Retrieve(context.Background(), "evt_1")
	if err != nil {
		t.Fatalf("Retrieve: %v", err)
	}
	if ev.Type != "payment_intent.succeeded" {
		t.Errorf("ev = %+v", ev)
	}
}
