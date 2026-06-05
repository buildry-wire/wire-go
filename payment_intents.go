package wire

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// PaymentIntentCreateParams are the inputs to Create.
type PaymentIntentCreateParams struct {
	Amount            int64             `json:"amount"`
	Currency          string            `json:"currency,omitempty"`
	AutomaticOperator *bool             `json:"automatic_operator,omitempty"`
	AllowedOperators  []string          `json:"allowed_operators,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty"`
	IdempotencyKey    string            `json:"-"`
}

// PaymentIntentConfirmParams are the optional inputs to Confirm.
type PaymentIntentConfirmParams struct {
	ReturnURL      string `json:"return_url,omitempty"`
	IdempotencyKey string `json:"-"`
}

// Create starts a new payment intent.
func (s *PaymentIntentService) Create(ctx context.Context, p *PaymentIntentCreateParams) (*PaymentIntent, error) {
	out := &PaymentIntent{}
	key := ""
	if p != nil {
		key = p.IdempotencyKey
	}
	if err := s.client.do(ctx, http.MethodPost, "/v1/payment_intents", p, nil, key, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Retrieve fetches a payment intent by id.
func (s *PaymentIntentService) Retrieve(ctx context.Context, id string) (*PaymentIntent, error) {
	out := &PaymentIntent{}
	if err := s.client.do(ctx, http.MethodGet, "/v1/payment_intents/"+url.PathEscape(id), nil, nil, "", out); err != nil {
		return nil, err
	}
	return out, nil
}

// Confirm submits a payment intent for processing.
func (s *PaymentIntentService) Confirm(ctx context.Context, id string, p *PaymentIntentConfirmParams) (*PaymentIntent, error) {
	out := &PaymentIntent{}
	key := ""
	if p != nil {
		key = p.IdempotencyKey
	}
	if err := s.client.do(ctx, http.MethodPost, "/v1/payment_intents/"+url.PathEscape(id)+"/confirm", p, nil, key, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Cancel cancels a payment intent.
func (s *PaymentIntentService) Cancel(ctx context.Context, id string) (*PaymentIntent, error) {
	out := &PaymentIntent{}
	if err := s.client.do(ctx, http.MethodPost, "/v1/payment_intents/"+url.PathEscape(id)+"/cancel", nil, nil, "", out); err != nil {
		return nil, err
	}
	return out, nil
}

// List returns an auto-paginating iterator over payment intents.
func (s *PaymentIntentService) List(ctx context.Context, p *ListParams) *Iter[PaymentIntent] {
	return &Iter[PaymentIntent]{
		ctx:   ctx,
		limit: limitOf(p),
		after: startingAfterOf(p),
		idOf:  func(pi PaymentIntent) string { return pi.ID },
		fetch: func(ctx context.Context, after string, limit int) ([]PaymentIntent, bool, error) {
			q := listQuery(after, limit)
			var page List[PaymentIntent]
			if err := s.client.do(ctx, http.MethodGet, "/v1/payment_intents", nil, q, "", &page); err != nil {
				return nil, false, err
			}
			return page.Data, page.HasMore, nil
		},
	}
}

// listQuery builds shared pagination query params.
func listQuery(after string, limit int) url.Values {
	q := url.Values{}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if after != "" {
		q.Set("starting_after", after)
	}
	return q
}

func limitOf(p *ListParams) int {
	if p == nil {
		return 0
	}
	return p.Limit
}

func startingAfterOf(p *ListParams) string {
	if p == nil {
		return ""
	}
	return p.StartingAfter
}
