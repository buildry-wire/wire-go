package wire

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

// SignatureHeader is the header carrying the webhook signature.
const SignatureHeader = "WirePayment-Signature"

// DefaultTolerance is the max allowed clock skew between signature timestamp and now.
const DefaultTolerance = 300 * time.Second

// ErrInvalidSignature is returned when a webhook signature does not verify.
var ErrInvalidSignature = errors.New("wire: webhook signature verification failed")

// Verify checks a webhook payload's signature and returns the parsed Event.
// payload must be the raw request body bytes (verify BEFORE JSON parsing).
func (s *WebhooksService) Verify(payload []byte, header, secret string) (*Event, error) {
	return s.verifyAt(payload, header, secret, DefaultTolerance, time.Now())
}

// VerifyWithTolerance is Verify with a custom timestamp tolerance.
func (s *WebhooksService) VerifyWithTolerance(payload []byte, header, secret string, tolerance time.Duration) (*Event, error) {
	return s.verifyAt(payload, header, secret, tolerance, time.Now())
}

// verifyAt is the testable core: it takes an explicit "now".
func (s *WebhooksService) verifyAt(payload []byte, header, secret string, tolerance time.Duration, now time.Time) (*Event, error) {
	ts, v1, ok := parseSignatureHeader(header)
	if !ok {
		return nil, ErrInvalidSignature
	}
	if d := now.Unix() - ts; d > int64(tolerance.Seconds()) || d < -int64(tolerance.Seconds()) {
		return nil, ErrInvalidSignature
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(strconv.FormatInt(ts, 10) + "."))
	mac.Write(payload)
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(v1)) {
		return nil, ErrInvalidSignature
	}
	var ev Event
	if err := json.Unmarshal(payload, &ev); err != nil {
		return nil, err
	}
	return &ev, nil
}

func parseSignatureHeader(h string) (ts int64, v1 string, ok bool) {
	for _, part := range strings.Split(h, ",") {
		k, val, found := strings.Cut(strings.TrimSpace(part), "=")
		if !found {
			continue
		}
		switch k {
		case "t":
			n, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return 0, "", false
			}
			ts = n
		case "v1":
			v1 = val
		}
	}
	if ts == 0 || v1 == "" {
		return 0, "", false
	}
	return ts, v1, true
}
