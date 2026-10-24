package pantech

import (
	"context"
	"net/url"
)

// OperationsService follows asynchronous changes.
type OperationsService struct{ c *Client }

// Get gets an operation: GET /operations/{operation_id}.
// Scope: read.
func (s *OperationsService) Get(ctx context.Context, operationID string, opts ...RequestOption) (*Operation, error) {
	var out Operation
	if err := s.c.Do(ctx, "GET", "/operations/"+url.PathEscape(operationID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
