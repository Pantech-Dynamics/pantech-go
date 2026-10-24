package pantech

import (
	"context"
	"net/url"
)

// LoadBalancersService manages load balancers.
type LoadBalancersService struct{ c *Client }

// List lists load balancers: GET /load-balancers.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *LoadBalancersService) List(ctx context.Context, params *ListLoadBalancersParams, opts ...RequestOption) ([]LoadBalancer, error) {
	return listAll[LoadBalancer](ctx, s.c, "/load-balancers", params.values(), opts...)
}

// Create creates a load balancer: POST /load-balancers.
// Scope: write.
func (s *LoadBalancersService) Create(ctx context.Context, req CreateLoadBalancerRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/load-balancers", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get gets a load balancer: GET /load-balancers/{load_balancer_id}.
// Scope: read.
func (s *LoadBalancersService) Get(ctx context.Context, loadBalancerID string, opts ...RequestOption) (*LoadBalancer, error) {
	var out LoadBalancer
	if err := s.c.Do(ctx, "GET", "/load-balancers/"+url.PathEscape(loadBalancerID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update changes a load balancer: PATCH /load-balancers/{load_balancer_id}.
// Scope: write.
func (s *LoadBalancersService) Update(ctx context.Context, loadBalancerID string, req UpdateLoadBalancerRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "PATCH", "/load-balancers/"+url.PathEscape(loadBalancerID), nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a load balancer: DELETE /load-balancers/{load_balancer_id}.
// Scope: write.
func (s *LoadBalancersService) Delete(ctx context.Context, loadBalancerID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "DELETE", "/load-balancers/"+url.PathEscape(loadBalancerID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
