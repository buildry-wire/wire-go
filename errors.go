package wire

import (
	"encoding/json"
	"fmt"
)

// Error is a typed Wire API error. Use errors.As to extract it.
type Error struct {
	Type                string `json:"type"`
	Code                string `json:"code"`
	Message             string `json:"message"`
	Param               string `json:"param"`
	RequestID           string `json:"request_id"`
	DocURL              string `json:"doc_url"`
	OperatorDeclineCode string `json:"operator_decline_code"`
	StatusCode          int    `json:"-"`
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("wire: %s (type=%s, code=%s, request_id=%s)", e.Message, e.Type, e.Code, e.RequestID)
	}
	return fmt.Sprintf("wire: %s (type=%s, status=%d, request_id=%s)", e.Message, e.Type, e.StatusCode, e.RequestID)
}

// parseError decodes the Wire error envelope; falls back to a generic error.
func parseError(status int, body []byte) error {
	var env struct {
		Error *Error `json:"error"`
	}
	if err := json.Unmarshal(body, &env); err == nil && env.Error != nil {
		env.Error.StatusCode = status
		return env.Error
	}
	return &Error{
		Type:       "api_error",
		Message:    fmt.Sprintf("unexpected response (status %d)", status),
		StatusCode: status,
	}
}
