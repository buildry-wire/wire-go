package wire

import (
	"context"
	"net/http"
	"net/url"
)

// Retrieve fetches a charge by id.
func (s *ChargeService) Retrieve(ctx context.Context, id string) (*Charge, error) {
	out := &Charge{}
	if err := s.client.do(ctx, http.MethodGet, "/v1/charges/"+url.PathEscape(id), nil, nil, "", out); err != nil {
		return nil, err
	}
	return out, nil
}

// List returns an auto-paginating iterator over charges.
func (s *ChargeService) List(ctx context.Context, p *ListParams) *Iter[Charge] {
	return &Iter[Charge]{
		ctx:   ctx,
		limit: limitOf(p),
		after: startingAfterOf(p),
		idOf:  func(ch Charge) string { return ch.ID },
		fetch: func(ctx context.Context, after string, limit int) ([]Charge, bool, error) {
			var page List[Charge]
			if err := s.client.do(ctx, http.MethodGet, "/v1/charges", nil, listQuery(after, limit), "", &page); err != nil {
				return nil, false, err
			}
			return page.Data, page.HasMore, nil
		},
	}
}
