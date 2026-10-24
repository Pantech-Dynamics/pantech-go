package pantech

import (
	"context"
	"iter"
	"net/url"
)

// SecurityGroupsService manages security groups, the firewalls of standard instances.
type SecurityGroupsService struct{ c *Client }

// List lists security groups: GET /security-groups.
// Scope: read.
func (s *SecurityGroupsService) List(ctx context.Context, params *ListParams, opts ...RequestOption) (*Page[SecurityGroup], error) {
	var out Page[SecurityGroup]
	if err := s.c.Do(ctx, "GET", "/security-groups", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// All iterates over every item of List, fetching pages as needed.
// It stops at the first error, which it yields.
func (s *SecurityGroupsService) All(ctx context.Context, params *ListParams, opts ...RequestOption) iter.Seq2[SecurityGroup, error] {
	var p ListParams
	if params != nil {
		p = *params
	}
	return paginate(p.Cursor, func(cursor string) (*Page[SecurityGroup], error) {
		p.Cursor = cursor
		return s.List(ctx, &p, opts...)
	})
}

// Create creates a security group: POST /security-groups.
// Scope: write.
func (s *SecurityGroupsService) Create(ctx context.Context, req CreateSecurityGroupRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/security-groups", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get gets a security group: GET /security-groups/{security_group_id}.
// Scope: read.
func (s *SecurityGroupsService) Get(ctx context.Context, securityGroupID string, opts ...RequestOption) (*SecurityGroup, error) {
	var out SecurityGroup
	if err := s.c.Do(ctx, "GET", "/security-groups/"+url.PathEscape(securityGroupID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a security group: DELETE /security-groups/{security_group_id}.
// Scope: write.
func (s *SecurityGroupsService) Delete(ctx context.Context, securityGroupID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "DELETE", "/security-groups/"+url.PathEscape(securityGroupID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ReplaceRules replaces a security group's rules: PUT /security-groups/{security_group_id}/rules.
// Scope: write.
func (s *SecurityGroupsService) ReplaceRules(ctx context.Context, securityGroupID string, req ReplaceSecurityGroupRulesRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "PUT", "/security-groups/"+url.PathEscape(securityGroupID)+"/rules", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
