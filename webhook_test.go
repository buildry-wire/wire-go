package wire

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestWebhookVerifyVectors(t *testing.T) {
	raw, err := os.ReadFile("testdata/webhook-signatures.json")
	if err != nil {
		t.Fatalf("read vectors: %v", err)
	}
	var v struct {
		Secret    string `json:"secret"`
		Now       int64  `json:"now"`
		Tolerance int64  `json:"tolerance_seconds"`
		Cases     []struct {
			Name   string `json:"name"`
			Body   string `json:"body"`
			Header string `json:"header"`
			Valid  bool   `json:"valid"`
		} `json:"cases"`
	}
	if err := json.Unmarshal(raw, &v); err != nil {
		t.Fatalf("parse vectors: %v", err)
	}

	w := &WebhooksService{}
	now := time.Unix(v.Now, 0)
	tol := time.Duration(v.Tolerance) * time.Second
	for _, c := range v.Cases {
		ev, err := w.verifyAt([]byte(c.Body), c.Header, v.Secret, tol, now)
		got := err == nil
		if got != c.Valid {
			t.Errorf("case %q: valid=%v, want %v (err=%v)", c.Name, got, c.Valid, err)
		}
		if c.Valid && (ev == nil || ev.Type == "") {
			t.Errorf("case %q: expected parsed event, got %+v", c.Name, ev)
		}
	}
}
