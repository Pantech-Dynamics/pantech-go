package pantech

import (
	"context"
	"iter"
	"net/url"
)

// NetworksService manages VPC networks.
type NetworksService struct{ c *Client }

// List lists VPC networks: GET /networks.
// Scope: read.
func (s *NetworksService) List(ctx context.Context, params *ListParams, opts ...RequestOption) (*Page[Network], error) {
	var out Page[Network]
	if err := s.c.Do(ctx, "GET", "/networks", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// All iterates over every item of List, fetching pages as needed.
// It stops at the first error, which it yields.
func (s *NetworksService) All(ctx context.Context, params *ListParams, opts ...RequestOption) iter.Seq2[Network, error] {
	var p ListParams
	if params != nil {
		p = *params
	}
	return paginate(p.Cursor, func(cursor string) (*Page[Network], error) {
		p.Cursor = cursor
		return s.List(ctx, &p, opts...)
	})
}

// Create creates a VPC network: POST /networks.
// Scope: write.
func (s *NetworksService) Create(ctx context.Context, req CreateNetworkRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/networks", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Quote prices a VPC network: GET /networks/quote.
// Scope: read.
func (s *NetworksService) Quote(ctx context.Context, opts ...RequestOption) (*MonthlyQuote, error) {
	var out MonthlyQuote
	if err := s.c.Do(ctx, "GET", "/networks/quote", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get gets a VPC network: GET /networks/{network_id}.
// Scope: read.
func (s *NetworksService) Get(ctx context.Context, networkID string, opts ...RequestOption) (*Network, error) {
	var out Network
	if err := s.c.Do(ctx, "GET", "/networks/"+url.PathEscape(networkID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a VPC network: DELETE /networks/{network_id}.
// Scope: write.
func (s *NetworksService) Delete(ctx context.Context, networkID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "DELETE", "/networks/"+url.PathEscape(networkID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// SubnetsService manages the subnets of VPC networks.
type SubnetsService struct{ c *Client }

// List lists subnets: GET /networks/{network_id}/subnets.
// Scope: read.
func (s *SubnetsService) List(ctx context.Context, networkID string, params *ListParams, opts ...RequestOption) (*Page[Subnet], error) {
	var out Page[Subnet]
	if err := s.c.Do(ctx, "GET", "/networks/"+url.PathEscape(networkID)+"/subnets", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// All iterates over every item of List, fetching pages as needed.
// It stops at the first error, which it yields.
func (s *SubnetsService) All(ctx context.Context, networkID string, params *ListParams, opts ...RequestOption) iter.Seq2[Subnet, error] {
	var p ListParams
	if params != nil {
		p = *params
	}
	return paginate(p.Cursor, func(cursor string) (*Page[Subnet], error) {
		p.Cursor = cursor
		return s.List(ctx, networkID, &p, opts...)
	})
}

// Create creates a subnet: POST /networks/{network_id}/subnets.
// Scope: write.
func (s *SubnetsService) Create(ctx context.Context, networkID string, req CreateSubnetRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/networks/"+url.PathEscape(networkID)+"/subnets", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get gets a subnet: GET /subnets/{subnet_id}.
// Scope: read.
func (s *SubnetsService) Get(ctx context.Context, subnetID string, opts ...RequestOption) (*Subnet, error) {
	var out Subnet
	if err := s.c.Do(ctx, "GET", "/subnets/"+url.PathEscape(subnetID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a subnet: DELETE /subnets/{subnet_id}.
// Scope: write.
func (s *SubnetsService) Delete(ctx context.Context, subnetID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "DELETE", "/subnets/"+url.PathEscape(subnetID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// FirewallRulesService manages subnet firewall rules.
type FirewallRulesService struct{ c *Client }

// List lists a subnet's firewall rules: GET /subnets/{subnet_id}/firewall-rules.
// Scope: read.
func (s *FirewallRulesService) List(ctx context.Context, subnetID string, opts ...RequestOption) (*FirewallRuleList, error) {
	var out FirewallRuleList
	if err := s.c.Do(ctx, "GET", "/subnets/"+url.PathEscape(subnetID)+"/firewall-rules", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create adds a firewall rule: POST /subnets/{subnet_id}/firewall-rules.
// Scope: write.
func (s *FirewallRulesService) Create(ctx context.Context, subnetID string, req CreateFirewallRuleRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/subnets/"+url.PathEscape(subnetID)+"/firewall-rules", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get gets a firewall rule: GET /subnets/{subnet_id}/firewall-rules/{rule_id}.
// Scope: read.
func (s *FirewallRulesService) Get(ctx context.Context, subnetID string, ruleID string, opts ...RequestOption) (*FirewallRule, error) {
	var out FirewallRule
	if err := s.c.Do(ctx, "GET", "/subnets/"+url.PathEscape(subnetID)+"/firewall-rules/"+url.PathEscape(ruleID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a firewall rule: DELETE /firewall-rules/{rule_id}.
// Scope: write.
func (s *FirewallRulesService) Delete(ctx context.Context, ruleID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "DELETE", "/firewall-rules/"+url.PathEscape(ruleID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
