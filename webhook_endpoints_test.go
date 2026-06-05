package wire

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebhookEndpointsCreateAndDelete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/webhook_endpoints":
			json.NewEncoder(w).Encode(map[string]any{"id": "we_1", "object": "webhook_endpoint",
				"url": "https://m.example/wh", "status": "enabled", "secret": "whsec_abc"})
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/webhook_endpoints/we_1":
			json.NewEncoder(w).Encode(map[string]any{"id": "we_1", "object": "webhook_endpoint", "deleted": true})
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
	}))
	defer srv.Close()

	c := NewClient("sk_test_123", WithBaseURL(srv.URL))
	we, err := c.WebhookEndpoints.Create(context.Background(), &WebhookEndpointCreateParams{
		URL: "https://m.example/wh", EnabledEvents: []string{"*"}})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if we.Secret != "whsec_abc" {
		t.Errorf("secret = %q", we.Secret)
	}
	del, err := c.WebhookEndpoints.Delete(context.Background(), "we_1")
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if !del.Deleted {
		t.Errorf("not deleted: %+v", del)
	}
}
