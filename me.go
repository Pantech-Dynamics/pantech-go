package pantech

import (
	"context"
)

// MeService describes and revokes the API key making the call.
type MeService struct{ c *Client }

// Get describes this API key: GET /me.
// Scope: read.
func (s *MeService) Get(ctx context.Context, opts ...RequestOption) (*Me, error) {
	var out Me
	if err := s.c.Do(ctx, "GET", "/me", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Revoke revokes the API key making the call: DELETE /me. The key stops
// working at once; a second call fails with 401 because the key is gone.
// Scope: read.
func (s *MeService) Revoke(ctx context.Context, opts ...RequestOption) error {
	return s.c.Do(ctx, "DELETE", "/me", nil, nil, nil, opts...)
}
