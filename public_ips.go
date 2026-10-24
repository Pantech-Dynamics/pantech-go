package pantech

import (
	"context"
	"net/url"
)

// PublicIPsService manages public IPs for VPC instances.
type PublicIPsService struct{ c *Client }

// List lists public IPs: GET /public-ips.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *PublicIPsService) List(ctx context.Context, params *ListPublicIPsParams, opts ...RequestOption) ([]PublicIP, error) {
	return listAll[PublicIP](ctx, s.c, "/public-ips", params.values(), opts...)
}

// Create allocates a public IP: POST /public-ips.
// Scope: write.
func (s *PublicIPsService) Create(ctx context.Context, req CreatePublicIPRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/public-ips", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Quote prices a public IP: GET /public-ips/quote.
// Scope: read.
func (s *PublicIPsService) Quote(ctx context.Context, opts ...RequestOption) (*MonthlyQuote, error) {
	var out MonthlyQuote
	if err := s.c.Do(ctx, "GET", "/public-ips/quote", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get gets a public IP: GET /public-ips/{public_ip_id}.
// Scope: read.
func (s *PublicIPsService) Get(ctx context.Context, publicIPID string, opts ...RequestOption) (*PublicIP, error) {
	var out PublicIP
	if err := s.c.Do(ctx, "GET", "/public-ips/"+url.PathEscape(publicIPID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete releases a public IP: DELETE /public-ips/{public_ip_id}.
// Scope: write.
func (s *PublicIPsService) Delete(ctx context.Context, publicIPID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "DELETE", "/public-ips/"+url.PathEscape(publicIPID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Attach attaches a public IP to an instance: POST /public-ips/{public_ip_id}/attach.
// Scope: write.
func (s *PublicIPsService) Attach(ctx context.Context, publicIPID string, req AttachPublicIPRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/public-ips/"+url.PathEscape(publicIPID)+"/attach", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Detach detaches a public IP from its instance: POST /public-ips/{public_ip_id}/detach.
// Scope: write.
func (s *PublicIPsService) Detach(ctx context.Context, publicIPID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/public-ips/"+url.PathEscape(publicIPID)+"/detach", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// PortForwardingRulesService manages port forwarding on public IPs.
type PortForwardingRulesService struct{ c *Client }

// List lists port forwarding rules: GET /public-ips/{public_ip_id}/port-forwarding-rules.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *PortForwardingRulesService) List(ctx context.Context, publicIPID string, opts ...RequestOption) ([]PortForwardingRule, error) {
	return listAll[PortForwardingRule](ctx, s.c, "/public-ips/"+url.PathEscape(publicIPID)+"/port-forwarding-rules", nil, opts...)
}

// Create adds a port forwarding rule: POST /public-ips/{public_ip_id}/port-forwarding-rules.
// Scope: write.
func (s *PortForwardingRulesService) Create(ctx context.Context, publicIPID string, req CreatePortForwardingRuleRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/public-ips/"+url.PathEscape(publicIPID)+"/port-forwarding-rules", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get gets a port forwarding rule: GET /public-ips/{public_ip_id}/port-forwarding-rules/{rule_id}.
// Scope: read.
func (s *PortForwardingRulesService) Get(ctx context.Context, publicIPID string, ruleID string, opts ...RequestOption) (*PortForwardingRule, error) {
	var out PortForwardingRule
	if err := s.c.Do(ctx, "GET", "/public-ips/"+url.PathEscape(publicIPID)+"/port-forwarding-rules/"+url.PathEscape(ruleID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a port forwarding rule: DELETE /port-forwarding-rules/{rule_id}.
// Scope: write.
func (s *PortForwardingRulesService) Delete(ctx context.Context, ruleID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "DELETE", "/port-forwarding-rules/"+url.PathEscape(ruleID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
