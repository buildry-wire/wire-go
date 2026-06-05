package wire

import (
	"context"
	"net/http"
	"net/url"
)

// WebhookEndpointCreateParams are the inputs to Create.
type WebhookEndpointCreateParams struct {
	URL            string   `json:"url"`
	EnabledEvents  []string `json:"enabled_events,omitempty"`
	IdempotencyKey string   `json:"-"`
}

// WebhookEndpointUpdateParams are the partial-update inputs.
type WebhookEndpointUpdateParams struct {
	URL            *string  `json:"url,omitempty"`
	EnabledEvents  []string `json:"enabled_events,omitempty"`
	Status         *string  `json:"status,omitempty"`
	IdempotencyKey string   `json:"-"`
}

// Create registers a webhook endpoint. The signing secret is returned once.
func (s *WebhookEndpointService) Create(ctx context.Context, p *WebhookEndpointCreateParams) (*WebhookEndpoint, error) {
	out := &WebhookEndpoint{}
	key := ""
	if p != nil {
		key = p.IdempotencyKey
	}
	if err := s.client.do(ctx, http.MethodPost, "/v1/webhook_endpoints", p, nil, key, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Retrieve fetches a webhook endpoint by id.
func (s *WebhookEndpointService) Retrieve(ctx context.Context, id string) (*WebhookEndpoint, error) {
	out := &WebhookEndpoint{}
	if err := s.client.do(ctx, http.MethodGet, "/v1/webhook_endpoints/"+url.PathEscape(id), nil, nil, "", out); err != nil {
		return nil, err
	}
	return out, nil
}

// Update modifies a webhook endpoint (Stripe-style POST update).
func (s *WebhookEndpointService) Update(ctx context.Context, id string, p *WebhookEndpointUpdateParams) (*WebhookEndpoint, error) {
	out := &WebhookEndpoint{}
	key := ""
	if p != nil {
		key = p.IdempotencyKey
	}
	if err := s.client.do(ctx, http.MethodPost, "/v1/webhook_endpoints/"+url.PathEscape(id), p, nil, key, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Delete removes a webhook endpoint.
func (s *WebhookEndpointService) Delete(ctx context.Context, id string) (*Deleted, error) {
	out := &Deleted{}
	if err := s.client.do(ctx, http.MethodDelete, "/v1/webhook_endpoints/"+url.PathEscape(id), nil, nil, "", out); err != nil {
		return nil, err
	}
	return out, nil
}

// List returns an auto-paginating iterator over webhook endpoints.
func (s *WebhookEndpointService) List(ctx context.Context, p *ListParams) *Iter[WebhookEndpoint] {
	return &Iter[WebhookEndpoint]{
		ctx:   ctx,
		limit: limitOf(p),
		after: startingAfterOf(p),
		idOf:  func(w WebhookEndpoint) string { return w.ID },
		fetch: func(ctx context.Context, after string, limit int) ([]WebhookEndpoint, bool, error) {
			var page List[WebhookEndpoint]
			if err := s.client.do(ctx, http.MethodGet, "/v1/webhook_endpoints", nil, listQuery(after, limit), "", &page); err != nil {
				return nil, false, err
			}
			return page.Data, page.HasMore, nil
		},
	}
}
