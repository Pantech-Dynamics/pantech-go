package pantech

import (
	"context"
	"net/url"
)

// Cloud providers Pantech Migrate can connect; ListProviders says which this
// deployment offers now.
const (
	CloudProviderAWS          = "aws"
	CloudProviderGCP          = "gcp"
	CloudProviderAzure        = "azure"
	CloudProviderHuawei       = "huawei"
	CloudProviderDigitalOcean = "digitalocean"
	CloudProviderHetzner      = "hetzner"
	CloudProviderLinode       = "linode"
	CloudProviderVultr        = "vultr"
	CloudProviderContabo      = "contabo"
	CloudProviderOpenStack    = "openstack"
)

// How Pantech reads a provider: federated (Pantech signs short-lived tokens
// your cloud trusts) or agent (you run the Pantech Migrate agent, which keeps
// your read-only token on your machine).
const (
	CloudAccessFederated = "federated"
	CloudAccessAgent     = "agent"
)

// ListProviders lists every cloud Pantech Migrate can connect, how it reads
// each, and whether a scan reads its spend and its people: GET /cloud-providers.
// Scope: read.
func (s *MigrateService) ListProviders(ctx context.Context, opts ...RequestOption) ([]CloudProvider, error) {
	return listAll[CloudProvider](ctx, s.c, "/cloud-providers", nil, opts...)
}

// CreateAgentEnrollment issues a one-time enrollment code for a connection
// read by the Pantech Migrate agent: POST /cloud-connections/{connection_id}/agent-enrollments.
// The code is shown only in this answer and works once, for an hour; run
// EnrollCommand where the agent will run. Issuing a code replaces any unused
// one. Each call issues a new code, so it is never retried after a lost
// answer. 409 CLOUD_CONNECTION_NOT_AGENT for a federated connection.
// Scope: write.
func (s *MigrateService) CreateAgentEnrollment(ctx context.Context, connectionID string, opts ...RequestOption) (*CloudAgentEnrollment, error) {
	var out CloudAgentEnrollment
	if err := s.c.Do(ctx, "POST", "/cloud-connections/"+url.PathEscape(connectionID)+"/agent-enrollments", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// RevokeAgent ends a connection's Migrate agent at once and returns the
// connection, which then waits for a new agent: DELETE /cloud-connections/{connection_id}/agent.
// Revoking again changes nothing.
// Scope: write.
func (s *MigrateService) RevokeAgent(ctx context.Context, connectionID string, opts ...RequestOption) (*CloudConnection, error) {
	var out CloudConnection
	if err := s.c.Do(ctx, "DELETE", "/cloud-connections/"+url.PathEscape(connectionID)+"/agent", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// SendInvites invites the people a plan lists (items of kind user), by key,
// with each item's role: POST /migrations/{migration_id}/invites. Safe to
// repeat: nobody is invited twice; each person's Outcome says what happened.
// Scope: write.
func (s *MigrateService) SendInvites(ctx context.Context, migrationID string, req SendMigrationInvitesRequest, opts ...RequestOption) (*MigrationInvites, error) {
	var out MigrationInvites
	if err := s.c.Do(ctx, "POST", "/migrations/"+url.PathEscape(migrationID)+"/invites", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
