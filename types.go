package wire

import "encoding/json"

// PaymentIntent is the primary object for accepting a payment.
type PaymentIntent struct {
	ID                string            `json:"id"`
	Object            string            `json:"object"`
	Amount            int64             `json:"amount"`
	Currency          string            `json:"currency"`
	Status            string            `json:"status"`
	ClientSecret      string            `json:"client_secret"`
	AutomaticOperator bool              `json:"automatic_operator"`
	AllowedOperators  []string          `json:"allowed_operators"`
	SelectedOperator  *string           `json:"selected_operator"`
	NextAction        json.RawMessage   `json:"next_action"`
	Metadata          map[string]string `json:"metadata"`
	Livemode          bool              `json:"livemode"`
	Created           int64             `json:"created"`
	ExpiresAt         *int64            `json:"expires_at"`
}

// Charge is a single attempt to move money via an operator.
type Charge struct {
	ID               string  `json:"id"`
	Object           string  `json:"object"`
	PaymentIntent    string  `json:"payment_intent"`
	Operator         string  `json:"operator"`
	OperatorChargeID *string `json:"operator_charge_id"`
	Status           string  `json:"status"`
	Amount           int64   `json:"amount"`
	Fee              int64   `json:"fee"`
	AmountRefunded   int64   `json:"amount_refunded"`
	FailureCode      *string `json:"failure_code"`
	FailureMessage   *string `json:"failure_message"`
	Livemode         bool    `json:"livemode"`
	Created          int64   `json:"created"`
}

// Event is a record of something that happened, delivered via webhooks.
type Event struct {
	ID         string          `json:"id"`
	Object     string          `json:"object"`
	Type       string          `json:"type"`
	APIVersion string          `json:"api_version"`
	Data       json.RawMessage `json:"data"`
	Livemode   bool            `json:"livemode"`
	Created    int64           `json:"created"`
}

// WebhookEndpoint is a merchant-registered URL that receives events.
type WebhookEndpoint struct {
	ID            string   `json:"id"`
	Object        string   `json:"object"`
	URL           string   `json:"url"`
	EnabledEvents []string `json:"enabled_events"`
	Status        string   `json:"status"`
	Secret        string   `json:"secret,omitempty"` // returned only at creation
	Livemode      bool     `json:"livemode"`
	Created       int64    `json:"created"`
}

// Deleted is the response shape for delete operations.
type Deleted struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

// List is one page of a cursor-paginated collection.
type List[T any] struct {
	Object  string `json:"object"`
	Data    []T    `json:"data"`
	HasMore bool   `json:"has_more"`
}

// ListParams are shared pagination parameters.
type ListParams struct {
	Limit         int
	StartingAfter string
	EndingBefore  string
}

// Service handles are defined in their own files; these are their concrete types.
type (
	// PaymentIntentService accesses /v1/payment_intents.
	PaymentIntentService struct{ client *Client }
	// ChargeService accesses /v1/charges.
	ChargeService struct{ client *Client }
	// EventService accesses /v1/events.
	EventService struct{ client *Client }
	// WebhookEndpointService accesses /v1/webhook_endpoints.
	WebhookEndpointService struct{ client *Client }
	// WebhooksService verifies inbound webhook signatures.
	WebhooksService struct{}
)
