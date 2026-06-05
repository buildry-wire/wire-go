package wire

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// DefaultBaseURL is the production Wire API endpoint.
const DefaultBaseURL = "https://api.wire.mn"

// Client is a Wire API client. Create it with NewClient.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	maxRetries int
	backoff    time.Duration

	PaymentIntents   *PaymentIntentService
	Charges          *ChargeService
	Events           *EventService
	WebhookEndpoints *WebhookEndpointService
	Webhooks         *WebhooksService
}

// Option configures a Client.
type Option func(*Client)

// WithBaseURL overrides the API base URL (e.g. a local dev server).
func WithBaseURL(u string) Option { return func(c *Client) { c.baseURL = strings.TrimRight(u, "/") } }

// WithHTTPClient injects a custom *http.Client.
func WithHTTPClient(h *http.Client) Option { return func(c *Client) { c.httpClient = h } }

// WithMaxRetries sets the max retry count for 429/5xx/network errors.
func WithMaxRetries(n int) Option { return func(c *Client) { c.maxRetries = n } }

// WithBackoff sets the base backoff duration between retries.
func WithBackoff(d time.Duration) Option { return func(c *Client) { c.backoff = d } }

// WithTimeout sets the per-request timeout on the default HTTP client.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.httpClient.Timeout = d }
}

// NewClient builds a Client with the given API key (sk_live_...).
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:     apiKey,
		baseURL:    DefaultBaseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		maxRetries: 2,
		backoff:    500 * time.Millisecond,
	}
	for _, o := range opts {
		o(c)
	}
	c.PaymentIntents = &PaymentIntentService{client: c}
	c.Charges = &ChargeService{client: c}
	c.Events = &EventService{client: c}
	c.WebhookEndpoints = &WebhookEndpointService{client: c}
	c.Webhooks = &WebhooksService{}
	return c
}

// do executes an API call with auth, optional idempotency key, retries, and JSON
// decoding. idemKey is used only for POST; if empty on a POST it is generated.
func (c *Client) do(ctx context.Context, method, path string, body any, query url.Values, idemKey string, out any) error {
	var raw []byte
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("wire: marshal body: %w", err)
		}
		raw = b
	}

	u := c.baseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	if method == http.MethodPost && idemKey == "" {
		idemKey = newIdempotencyKey()
	}

	var lastErr error
	for attempt := 0; ; attempt++ {
		req, err := http.NewRequestWithContext(ctx, method, u, bytesReader(raw))
		if err != nil {
			return fmt.Errorf("wire: new request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("Accept", "application/json")
		if raw != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		if method == http.MethodPost && idemKey != "" {
			req.Header.Set("Idempotency-Key", idemKey)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("wire: request failed: %w", err)
			if attempt < c.maxRetries && ctx.Err() == nil {
				if sleepErr := c.sleep(ctx, attempt, 0); sleepErr != nil {
					return sleepErr
				}
				continue
			}
			return lastErr
		}

		shouldRetry := resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500
		if shouldRetry && attempt < c.maxRetries {
			retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"))
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			if sleepErr := c.sleep(ctx, attempt, retryAfter); sleepErr != nil {
				return sleepErr
			}
			continue
		}

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if out != nil && len(respBody) > 0 {
				if err := json.Unmarshal(respBody, out); err != nil {
					return fmt.Errorf("wire: decode response: %w", err)
				}
			}
			return nil
		}
		return parseError(resp.StatusCode, respBody)
	}
}

func bytesReader(b []byte) io.Reader {
	if b == nil {
		return nil
	}
	return bytes.NewReader(b)
}

// sleep waits for the backoff interval (exponential) or an explicit Retry-After,
// honoring context cancellation.
func (c *Client) sleep(ctx context.Context, attempt int, retryAfter time.Duration) error {
	d := retryAfter
	if d <= 0 {
		d = c.backoff * (1 << attempt)
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

func parseRetryAfter(h string) time.Duration {
	if h == "" {
		return 0
	}
	if secs, err := strconv.Atoi(h); err == nil {
		return time.Duration(secs) * time.Second
	}
	return 0
}

func newIdempotencyKey() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return "idk_" + strings.ToLower(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b))
}
