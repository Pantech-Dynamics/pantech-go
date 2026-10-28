package pantech

import (
	"context"
	"fmt"
	"net/url"
)

// Cloud connection statuses. Connected and failed are final until the role
// ARN changes or a retry is asked for.
const (
	CloudConnectionAwaitingRole = "awaiting_role"
	CloudConnectionVerifying    = "verifying"
	CloudConnectionConnected    = "connected"
	CloudConnectionFailed       = "failed"
)

// Cloud scan statuses. Succeeded and failed are final.
const (
	CloudScanQueued    = "queued"
	CloudScanRunning   = "running"
	CloudScanSucceeded = "succeeded"
	CloudScanFailed    = "failed"
)

// MigrateService moves you from another cloud: it connects an AWS account
// through a read-only role (no access keys are shared), scans what runs
// there, and plans the move to Pantech with prices and a Terraform export.
type MigrateService struct{ c *Client }

// ListConnections lists cloud connections: GET /cloud-connections.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *MigrateService) ListConnections(ctx context.Context, opts ...RequestOption) ([]CloudConnection, error) {
	return listAll[CloudConnection](ctx, s.c, "/cloud-connections", nil, opts...)
}

// CreateConnection starts a connection to an AWS account: POST /cloud-connections.
// Its Setup holds the template (and QuickCreateURL) that creates the
// read-only role; send the stack's RoleArn output with UpdateConnection.
// Scope: write.
func (s *MigrateService) CreateConnection(ctx context.Context, req CreateCloudConnectionRequest, opts ...RequestOption) (*CloudConnection, error) {
	var out CloudConnection
	if err := s.c.Do(ctx, "POST", "/cloud-connections", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetConnection gets a cloud connection: GET /cloud-connections/{connection_id}.
// Scope: read.
func (s *MigrateService) GetConnection(ctx context.Context, connectionID string, opts ...RequestOption) (*CloudConnection, error) {
	var out CloudConnection
	if err := s.c.Do(ctx, "GET", "/cloud-connections/"+url.PathEscape(connectionID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateConnection renames a connection, sets its role ARN, or retries a
// failed check: PATCH /cloud-connections/{connection_id}. Setting RoleARN or
// Retry queues a verification; follow it with WaitForConnection.
// Scope: write.
func (s *MigrateService) UpdateConnection(ctx context.Context, connectionID string, req UpdateCloudConnectionRequest, opts ...RequestOption) (*CloudConnection, error) {
	var out CloudConnection
	if err := s.c.Do(ctx, "PATCH", "/cloud-connections/"+url.PathEscape(connectionID), nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteConnection removes a connection and its scans: DELETE /cloud-connections/{connection_id}.
// Migration plans made from it stay. Delete the stack in the other account to
// revoke the role there too.
// Scope: write.
func (s *MigrateService) DeleteConnection(ctx context.Context, connectionID string, opts ...RequestOption) error {
	return s.c.Do(ctx, "DELETE", "/cloud-connections/"+url.PathEscape(connectionID), nil, nil, nil, opts...)
}

// ListScans lists a connection's scans, without their inventories: GET /cloud-connections/{connection_id}/scans.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *MigrateService) ListScans(ctx context.Context, connectionID string, opts ...RequestOption) ([]CloudScan, error) {
	return listAll[CloudScan](ctx, s.c, "/cloud-connections/"+url.PathEscape(connectionID)+"/scans", nil, opts...)
}

// StartScan queues a scan of a connected account: POST /cloud-connections/{connection_id}/scans.
// Follow it with WaitForScan.
// Scope: write.
func (s *MigrateService) StartScan(ctx context.Context, connectionID string, opts ...RequestOption) (*CloudScan, error) {
	var out CloudScan
	if err := s.c.Do(ctx, "POST", "/cloud-connections/"+url.PathEscape(connectionID)+"/scans", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetScan gets a scan: GET /cloud-scans/{scan_id}. Pass
// &GetCloudScanParams{Include: "inventory"} for everything it found.
// Scope: read.
func (s *MigrateService) GetScan(ctx context.Context, scanID string, params *GetCloudScanParams, opts ...RequestOption) (*CloudScan, error) {
	var out CloudScan
	if err := s.c.Do(ctx, "GET", "/cloud-scans/"+url.PathEscape(scanID), params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListMigrations lists migration plans, without their items or estimate: GET /migrations.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *MigrateService) ListMigrations(ctx context.Context, opts ...RequestOption) ([]Migration, error) {
	return listAll[Migration](ctx, s.c, "/migrations", nil, opts...)
}

// CreateMigration builds a plan from a scan that succeeded: POST /migrations.
// Nothing is created on Pantech; the plan carries a price comparison.
// Scope: write.
func (s *MigrateService) CreateMigration(ctx context.Context, req CreateMigrationRequest, opts ...RequestOption) (*Migration, error) {
	var out Migration
	if err := s.c.Do(ctx, "POST", "/migrations", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMigration gets a plan with every item and a price comparison computed
// now: GET /migrations/{migration_id}.
// Scope: read.
func (s *MigrateService) GetMigration(ctx context.Context, migrationID string, opts ...RequestOption) (*Migration, error) {
	var out Migration
	if err := s.c.Do(ctx, "GET", "/migrations/"+url.PathEscape(migrationID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// RenameMigration changes a plan's name: PATCH /migrations/{migration_id}.
// Scope: write.
func (s *MigrateService) RenameMigration(ctx context.Context, migrationID string, req RenameMigrationRequest, opts ...RequestOption) (*Migration, error) {
	var out Migration
	if err := s.c.Do(ctx, "PATCH", "/migrations/"+url.PathEscape(migrationID), nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteMigration deletes a plan: DELETE /migrations/{migration_id}.
// Nothing on Pantech or in the other cloud is touched.
// Scope: write.
func (s *MigrateService) DeleteMigration(ctx context.Context, migrationID string, opts ...RequestOption) error {
	return s.c.Do(ctx, "DELETE", "/migrations/"+url.PathEscape(migrationID), nil, nil, nil, opts...)
}

// UpdateItem changes one item of a plan, named by its Key: PATCH /migrations/{migration_id}/items.
// It returns the whole plan.
// Scope: write.
func (s *MigrateService) UpdateItem(ctx context.Context, migrationID string, req UpdateMigrationItemRequest, opts ...RequestOption) (*Migration, error) {
	var out Migration
	if err := s.c.Do(ctx, "PATCH", "/migrations/"+url.PathEscape(migrationID)+"/items", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ExportTerraform returns the plan's planned items as Terraform (main.tf and
// variables.tf): GET /migrations/{migration_id}/terraform.
// Scope: read.
func (s *MigrateService) ExportTerraform(ctx context.Context, migrationID string, opts ...RequestOption) (*MigrationTerraform, error) {
	var out MigrationTerraform
	if err := s.c.Do(ctx, "GET", "/migrations/"+url.PathEscape(migrationID)+"/terraform", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// CloudFailureError is returned by WaitForConnection and WaitForScan when the
// connection or scan ended as failed. Failure's code is stable and its
// message says what to do.
type CloudFailureError struct {
	// Kind is "connection" or "scan".
	Kind    string
	ID      string
	Failure *CloudFailure
}

func (e *CloudFailureError) Error() string {
	msg := fmt.Sprintf("pantech: cloud %s %s failed", e.Kind, e.ID)
	if e.Failure != nil {
		msg += ": " + e.Failure.Message + " (" + e.Failure.Code + ")"
	}
	return msg
}

// WaitForConnection polls GET /cloud-connections/{id} until the connection is
// connected, and returns it. A failed check returns it with a
// *CloudFailureError. A connection still awaiting its role ARN is waited on
// too, so set the ARN with UpdateConnection first. opts may be nil.
func (s *MigrateService) WaitForConnection(ctx context.Context, connectionID string, opts *WaitOptions) (*CloudConnection, error) {
	var last *CloudConnection
	err := s.c.poll(ctx, opts, func(ctx context.Context) (bool, error) {
		conn, err := s.GetConnection(ctx, connectionID)
		if err != nil {
			return false, err
		}
		last = conn
		switch conn.Status {
		case CloudConnectionConnected:
			return true, nil
		case CloudConnectionFailed:
			return true, &CloudFailureError{Kind: "connection", ID: conn.ID, Failure: conn.Failure}
		}
		return false, nil
	})
	return last, err
}

// WaitForScan polls GET /cloud-scans/{id} until the scan succeeds, and
// returns it (without its inventory: read that with GetScan). A failed scan
// returns it with a *CloudFailureError. opts may be nil.
func (s *MigrateService) WaitForScan(ctx context.Context, scanID string, opts *WaitOptions) (*CloudScan, error) {
	var last *CloudScan
	err := s.c.poll(ctx, opts, func(ctx context.Context) (bool, error) {
		scan, err := s.GetScan(ctx, scanID, nil)
		if err != nil {
			return false, err
		}
		last = scan
		switch scan.Status {
		case CloudScanSucceeded:
			return true, nil
		case CloudScanFailed:
			return true, &CloudFailureError{Kind: "scan", ID: scan.ID, Failure: scan.Failure}
		}
		return false, nil
	})
	return last, err
}
