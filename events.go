package wire

import (
	"context"
	"net/http"
	"net/url"
)

// Retrieve fetches an event by id.
func (s *EventService) Retrieve(ctx context.Context, id string) (*Event, error) {
	out := &Event{}
	if err := s.client.do(ctx, http.MethodGet, "/v1/events/"+url.PathEscape(id), nil, nil, "", out); err != nil {
		return nil, err
	}
	return out, nil
}

// List returns an auto-paginating iterator over events.
func (s *EventService) List(ctx context.Context, p *ListParams) *Iter[Event] {
	return &Iter[Event]{
		ctx:   ctx,
		limit: limitOf(p),
		after: startingAfterOf(p),
		idOf:  func(e Event) string { return e.ID },
		fetch: func(ctx context.Context, after string, limit int) ([]Event, bool, error) {
			var page List[Event]
			if err := s.client.do(ctx, http.MethodGet, "/v1/events", nil, listQuery(after, limit), "", &page); err != nil {
				return nil, false, err
			}
			return page.Data, page.HasMore, nil
		},
	}
}
