package wire

import "context"

// Iter is a lazy auto-paginator over a cursor-paginated collection.
//
//	it := client.Charges.List(ctx, nil)
//	for it.Next() {
//	    ch := it.Current()
//	}
//	if err := it.Err(); err != nil { ... }
type Iter[T any] struct {
	ctx     context.Context
	limit   int
	after   string
	page    []T
	pos     int
	cur     T
	hasMore bool
	started bool
	err     error
	idOf    func(T) string
	fetch   func(ctx context.Context, after string, limit int) ([]T, bool, error)
}

// Next advances to the next item, fetching the next page when needed. It returns
// false when the collection is exhausted or an error occurred (check Err).
func (it *Iter[T]) Next() bool {
	if it.err != nil {
		return false
	}
	if it.pos >= len(it.page) {
		if it.started && !it.hasMore {
			return false
		}
		page, hasMore, err := it.fetch(it.ctx, it.after, it.limit)
		it.started = true
		if err != nil {
			it.err = err
			return false
		}
		it.page, it.hasMore, it.pos = page, hasMore, 0
		if len(page) == 0 {
			return false
		}
	}
	it.cur = it.page[it.pos]
	it.after = it.idOf(it.cur)
	it.pos++
	return true
}

// Current returns the item the iterator is positioned on.
func (it *Iter[T]) Current() T { return it.cur }

// Err returns the first error encountered while paginating.
func (it *Iter[T]) Err() error { return it.err }
