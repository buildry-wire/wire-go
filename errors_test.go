package wire

import (
	"errors"
	"testing"
)

func TestParseError(t *testing.T) {
	body := []byte(`{"error":{"type":"invalid_request_error","code":"amount_invalid","message":"amount must be positive","param":"amount","request_id":"req_123"}}`)
	err := parseError(400, body)

	var werr *Error
	if !errors.As(err, &werr) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if werr.StatusCode != 400 {
		t.Errorf("StatusCode = %d, want 400", werr.StatusCode)
	}
	if werr.Code != "amount_invalid" || werr.Param != "amount" || werr.RequestID != "req_123" {
		t.Errorf("fields not parsed: %+v", werr)
	}
	if werr.Error() == "" {
		t.Error("Error() is empty")
	}
}

func TestParseErrorFallback(t *testing.T) {
	err := parseError(500, []byte("not json"))
	var werr *Error
	if !errors.As(err, &werr) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if werr.StatusCode != 500 {
		t.Errorf("StatusCode = %d, want 500", werr.StatusCode)
	}
}
