package pantech

import (
	"context"
	"iter"
	"net/url"
)

// Policy statement effects. A matching deny wins over any allow.
const (
	EffectAllow = "allow"
	EffectDeny  = "deny"
)

// Member baselines: what a member with no policy may do.
const (
	MemberBaselineRead = "read"
	MemberBaselineNone = "none"
)

// Attachment principal types.
const (
	PrincipalMember = "member"
	PrincipalTeam   = "team"
)

// Simulation decisions.
const (
	DecisionOwner        = "owner"
	DecisionAdmin        = "admin"
	DecisionAllow        = "allow"
	DecisionBaselineRead = "baseline_read"
	DecisionExplicitDeny = "explicit_deny"
	DecisionImplicitDeny = "implicit_deny"
)

// Allow is an allow statement for actions on resources. No resources means
// every resource ("*").
//
//	pantech.Allow([]string{"instances:Read", "instances:Reboot"}, []string{"name:staging-*"})
func Allow(actions, resources []string) PolicyStatement {
	return statement(EffectAllow, actions, resources)
}

// Deny is a deny statement for actions on resources. No resources means every
// resource ("*"). A matching deny wins over any allow, for admins too.
//
//	pantech.Deny([]string{"instances:Delete"}, []string{"name:prod-*"}).WithSid("NoProdDeletes")
func Deny(actions, resources []string) PolicyStatement {
	return statement(EffectDeny, actions, resources)
}

func statement(effect string, actions, resources []string) PolicyStatement {
	if len(resources) == 0 {
		resources = []string{"*"}
	}
	return PolicyStatement{Effect: effect, Actions: append([]string(nil), actions...), Resources: append([]string(nil), resources...)}
}

// WithSid returns the statement with this sid.
func (s PolicyStatement) WithSid(sid string) PolicyStatement {
	s.Sid = sid
	return s
}

// WithSourceIP returns the statement applying only to requests from these CIDRs.
func (s PolicyStatement) WithSourceIP(cidrs ...string) PolicyStatement {
	c := s.conditions()
	c.SourceIP = append(append([]string(nil), c.SourceIP...), cidrs...)
	s.Conditions = c
	return s
}

// WithMFA returns the statement applying only to sessions signed in with MFA.
// API keys never satisfy it.
func (s PolicyStatement) WithMFA() PolicyStatement {
	c := s.conditions()
	c.MFA = Ptr(true)
	s.Conditions = c
	return s
}

// conditions is a copy of the statement's conditions, so the With methods
// never change the statement they were called on.
func (s PolicyStatement) conditions() *PolicyConditions {
	if s.Conditions == nil {
		return &PolicyConditions{}
	}
	c := *s.Conditions
	return &c
}

// IAMService manages teams and access policies: who in the organization may
// do what, on which resources. Owners are always allowed everything.
type IAMService struct{ c *Client }

// ListActions lists the action catalogue policies are written in: GET /iam/actions.
// Scope: read.
func (s *IAMService) ListActions(ctx context.Context, opts ...RequestOption) ([]IAMAction, error) {
	return listAll[IAMAction](ctx, s.c, "/iam/actions", nil, opts...)
}

// GetSettings gets the organization's access setting: GET /iam/settings.
// Scope: read.
func (s *IAMService) GetSettings(ctx context.Context, opts ...RequestOption) (*IAMSettings, error) {
	var out IAMSettings
	if err := s.c.Do(ctx, "GET", "/iam/settings", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateSettings changes the member baseline: PUT /iam/settings. `none`
// takes away members' read access to everything no policy grants. Owners only.
// Scope: write.
func (s *IAMService) UpdateSettings(ctx context.Context, req UpdateIAMSettingsRequest, opts ...RequestOption) (*IAMSettings, error) {
	var out IAMSettings
	if err := s.c.Do(ctx, "PUT", "/iam/settings", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListTeams lists teams: GET /iam/teams.
// Scope: read.
func (s *IAMService) ListTeams(ctx context.Context, params *PageParams, opts ...RequestOption) (*Page[IAMTeam], error) {
	var out Page[IAMTeam]
	if err := s.c.Do(ctx, "GET", "/iam/teams", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// AllTeams iterates over every item of ListTeams, fetching pages as needed.
// It stops at the first error, which it yields.
func (s *IAMService) AllTeams(ctx context.Context, params *PageParams, opts ...RequestOption) iter.Seq2[IAMTeam, error] {
	var p PageParams
	if params != nil {
		p = *params
	}
	return paginate(p.Cursor, func(cursor string) (*Page[IAMTeam], error) {
		p.Cursor = cursor
		return s.ListTeams(ctx, &p, opts...)
	})
}

// CreateTeam creates a team: POST /iam/teams.
// Scope: write.
func (s *IAMService) CreateTeam(ctx context.Context, req CreateIAMTeamRequest, opts ...RequestOption) (*IAMTeam, error) {
	var out IAMTeam
	if err := s.c.Do(ctx, "POST", "/iam/teams", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTeam gets a team: GET /iam/teams/{team_id}.
// Scope: read.
func (s *IAMService) GetTeam(ctx context.Context, teamID string, opts ...RequestOption) (*IAMTeam, error) {
	var out IAMTeam
	if err := s.c.Do(ctx, "GET", teamPath(teamID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateTeam renames a team or changes its description: PATCH /iam/teams/{team_id}.
// Scope: write.
func (s *IAMService) UpdateTeam(ctx context.Context, teamID string, req UpdateIAMTeamRequest, opts ...RequestOption) (*IAMTeam, error) {
	var out IAMTeam
	if err := s.c.Do(ctx, "PATCH", teamPath(teamID), nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteTeam deletes a team: DELETE /iam/teams/{team_id}. Its members lose
// the policies attached to it.
// Scope: delete.
func (s *IAMService) DeleteTeam(ctx context.Context, teamID string, opts ...RequestOption) error {
	return s.c.Do(ctx, "DELETE", teamPath(teamID), nil, nil, nil, opts...)
}

// ListTeamMembers lists a team's members: GET /iam/teams/{team_id}/members.
// Scope: read.
func (s *IAMService) ListTeamMembers(ctx context.Context, teamID string, params *PageParams, opts ...RequestOption) (*Page[IAMTeamMember], error) {
	var out Page[IAMTeamMember]
	if err := s.c.Do(ctx, "GET", teamPath(teamID)+"/members", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// AllTeamMembers iterates over every item of ListTeamMembers, fetching pages as needed.
// It stops at the first error, which it yields.
func (s *IAMService) AllTeamMembers(ctx context.Context, teamID string, params *PageParams, opts ...RequestOption) iter.Seq2[IAMTeamMember, error] {
	var p PageParams
	if params != nil {
		p = *params
	}
	return paginate(p.Cursor, func(cursor string) (*Page[IAMTeamMember], error) {
		p.Cursor = cursor
		return s.ListTeamMembers(ctx, teamID, &p, opts...)
	})
}

// AddTeamMember adds a member to a team: PUT /iam/teams/{team_id}/members/{user_id}.
// Adding one already in the team changes nothing. Admins may add only
// members, not other admins.
// Scope: write.
func (s *IAMService) AddTeamMember(ctx context.Context, teamID, userID string, opts ...RequestOption) (*IAMTeamMember, error) {
	var out IAMTeamMember
	if err := s.c.Do(ctx, "PUT", teamMemberPath(teamID, userID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// RemoveTeamMember removes a member from a team: DELETE /iam/teams/{team_id}/members/{user_id}.
// Scope: write.
func (s *IAMService) RemoveTeamMember(ctx context.Context, teamID, userID string, opts ...RequestOption) error {
	return s.c.Do(ctx, "DELETE", teamMemberPath(teamID, userID), nil, nil, nil, opts...)
}

// ListPolicies lists policies, yours and Pantech's: GET /iam/policies.
// Scope: read.
func (s *IAMService) ListPolicies(ctx context.Context, params *ListIAMPoliciesParams, opts ...RequestOption) (*Page[IAMPolicy], error) {
	var out Page[IAMPolicy]
	if err := s.c.Do(ctx, "GET", "/iam/policies", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// AllPolicies iterates over every item of ListPolicies, fetching pages as needed.
// It stops at the first error, which it yields.
func (s *IAMService) AllPolicies(ctx context.Context, params *ListIAMPoliciesParams, opts ...RequestOption) iter.Seq2[IAMPolicy, error] {
	var p ListIAMPoliciesParams
	if params != nil {
		p = *params
	}
	return paginate(p.Cursor, func(cursor string) (*Page[IAMPolicy], error) {
		p.Cursor = cursor
		return s.ListPolicies(ctx, &p, opts...)
	})
}

// CreatePolicy creates a policy: POST /iam/policies. It applies to nobody
// until attached. An invalid statement answers INVALID_POLICY naming the field.
// Scope: write.
func (s *IAMService) CreatePolicy(ctx context.Context, req CreateIAMPolicyRequest, opts ...RequestOption) (*IAMPolicy, error) {
	var out IAMPolicy
	if err := s.c.Do(ctx, "POST", "/iam/policies", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPolicy gets a policy: GET /iam/policies/{policy_id}.
// Scope: read.
func (s *IAMService) GetPolicy(ctx context.Context, policyID string, opts ...RequestOption) (*IAMPolicy, error) {
	var out IAMPolicy
	if err := s.c.Do(ctx, "GET", policyPath(policyID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdatePolicy replaces a policy's name, description and statements:
// PUT /iam/policies/{policy_id}. Pass WithIfMatch(policy.Version) to replace
// it only if nobody else has changed it since you read it (412
// VERSION_CONFLICT otherwise). Pantech's policies answer IAM_MANAGED_POLICY.
// Scope: write.
func (s *IAMService) UpdatePolicy(ctx context.Context, policyID string, req UpdateIAMPolicyRequest, opts ...RequestOption) (*IAMPolicy, error) {
	var out IAMPolicy
	if err := s.c.Do(ctx, "PUT", policyPath(policyID), nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeletePolicy deletes a policy and detaches it everywhere: DELETE /iam/policies/{policy_id}.
// Scope: delete.
func (s *IAMService) DeletePolicy(ctx context.Context, policyID string, opts ...RequestOption) error {
	return s.c.Do(ctx, "DELETE", policyPath(policyID), nil, nil, nil, opts...)
}

// ListAttachments lists who a policy is attached to: GET /iam/policies/{policy_id}/attachments.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *IAMService) ListAttachments(ctx context.Context, policyID string, opts ...RequestOption) ([]IAMAttachment, error) {
	return listAll[IAMAttachment](ctx, s.c, policyPath(policyID)+"/attachments", nil, opts...)
}

// Attach attaches a policy to a member or a team: POST /iam/policies/{policy_id}/attachments.
// Only owners attach to an admin or to a team containing one (IAM_ADMIN_TARGET).
// Scope: write.
func (s *IAMService) Attach(ctx context.Context, policyID string, req AttachIAMPolicyRequest, opts ...RequestOption) (*IAMAttachment, error) {
	var out IAMAttachment
	if err := s.c.Do(ctx, "POST", policyPath(policyID)+"/attachments", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Detach detaches a policy from a member or a team:
// DELETE /iam/policies/{policy_id}/attachments/{principal_type}/{principal_id}.
// principalType is PrincipalMember or PrincipalTeam.
// Scope: write.
func (s *IAMService) Detach(ctx context.Context, policyID, principalType, principalID string, opts ...RequestOption) error {
	path := policyPath(policyID) + "/attachments/" + url.PathEscape(principalType) + "/" + url.PathEscape(principalID)
	return s.c.Do(ctx, "DELETE", path, nil, nil, nil, opts...)
}

// ListMembers lists the organization's members, with their teams: GET /iam/members.
// Q matches email or name.
// Scope: read.
func (s *IAMService) ListMembers(ctx context.Context, params *ListIAMMembersParams, opts ...RequestOption) (*Page[IAMMember], error) {
	var out Page[IAMMember]
	if err := s.c.Do(ctx, "GET", "/iam/members", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// AllMembers iterates over every item of ListMembers, fetching pages as needed.
// It stops at the first error, which it yields.
func (s *IAMService) AllMembers(ctx context.Context, params *ListIAMMembersParams, opts ...RequestOption) iter.Seq2[IAMMember, error] {
	var p ListIAMMembersParams
	if params != nil {
		p = *params
	}
	return paginate(p.Cursor, func(cursor string) (*Page[IAMMember], error) {
		p.Cursor = cursor
		return s.ListMembers(ctx, &p, opts...)
	})
}

// GetEffectiveAccess says what applies to a member: their role, teams,
// policies and the baseline: GET /iam/members/{user_id}/effective-access.
// Scope: read.
func (s *IAMService) GetEffectiveAccess(ctx context.Context, userID string, opts ...RequestOption) (*IAMEffectiveAccess, error) {
	var out IAMEffectiveAccess
	if err := s.c.Do(ctx, "GET", "/iam/members/"+url.PathEscape(userID)+"/effective-access", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Simulate answers whether a member would be allowed an action, and which
// policy decided, without doing it: POST /iam/simulate. It changes nothing.
// Scope: read.
func (s *IAMService) Simulate(ctx context.Context, req SimulateIAMAccessRequest, opts ...RequestOption) (*IAMSimulation, error) {
	var out IAMSimulation
	if err := s.c.Do(ctx, "POST", "/iam/simulate", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

func teamPath(teamID string) string { return "/iam/teams/" + url.PathEscape(teamID) }

func teamMemberPath(teamID, userID string) string {
	return teamPath(teamID) + "/members/" + url.PathEscape(userID)
}

func policyPath(policyID string) string { return "/iam/policies/" + url.PathEscape(policyID) }
