package pantech

import (
	"context"
	"iter"
	"net/http"
	"net/url"
)

// Page is one page of a list: the items, and the cursor of the next page.
type Page[T any] struct {
	Data []T `json:"data"`
	// NextCursor is passed as the cursor parameter to get the next page. It
	// is nil on the last page.
	NextCursor *string `json:"next_cursor"`
}

// HasMore reports whether there is another page.
func (p *Page[T]) HasMore() bool { return p.NextCursor != nil && *p.NextCursor != "" }

// paginate yields every item of every page, starting at cursor first. A
// cursor that repeats ends the loop, so a server bug cannot spin it forever.
func paginate[T any](first string, fetch func(cursor string) (*Page[T], error)) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		cursor := first
		for {
			page, err := fetch(cursor)
			if err != nil {
				var zero T
				yield(zero, err)
				return
			}
			for _, item := range page.Data {
				if !yield(item, nil) {
					return
				}
			}
			if !page.HasMore() || *page.NextCursor == cursor {
				return
			}
			cursor = *page.NextCursor
		}
	}
}

// listAll collects a list the API documents as one page. It follows
// next_cursor anyway, so a future page size cannot silently cut it short.
func listAll[T any](ctx context.Context, c *Client, path string, query url.Values, opts ...RequestOption) ([]T, error) {
	items := []T{}
	for item, err := range paginate("", func(cursor string) (*Page[T], error) {
		q := url.Values{}
		for k, v := range query {
			q[k] = v
		}
		if cursor != "" {
			q.Set("cursor", cursor)
		}
		var page Page[T]
		if err := c.Do(ctx, http.MethodGet, path, q, nil, &page, opts...); err != nil {
			return nil, err
		}
		return &page, nil
	}) {
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
