package pantech

import (
	"context"
	"iter"
	"net/url"
)

// DatabasesService manages managed databases and their orders.
type DatabasesService struct{ c *Client }

// List lists databases: GET /databases.
// Scope: read.
func (s *DatabasesService) List(ctx context.Context, params *ListParams, opts ...RequestOption) (*Page[Database], error) {
	var out Page[Database]
	if err := s.c.Do(ctx, "GET", "/databases", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// All iterates over every item of List, fetching pages as needed.
// It stops at the first error, which it yields.
func (s *DatabasesService) All(ctx context.Context, params *ListParams, opts ...RequestOption) iter.Seq2[Database, error] {
	var p ListParams
	if params != nil {
		p = *params
	}
	return paginate(p.Cursor, func(cursor string) (*Page[Database], error) {
		p.Cursor = cursor
		return s.List(ctx, &p, opts...)
	})
}

// Create creates a database: POST /databases.
// Scope: write.
func (s *DatabasesService) Create(ctx context.Context, req CreateDatabaseRequest, opts ...RequestOption) (*DatabaseOrderAccepted, error) {
	var out DatabaseOrderAccepted
	if err := s.c.Do(ctx, "POST", "/databases", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get gets a database: GET /databases/{database_id}.
// Scope: read.
func (s *DatabasesService) Get(ctx context.Context, databaseID string, opts ...RequestOption) (*Database, error) {
	var out Database
	if err := s.c.Do(ctx, "GET", "/databases/"+url.PathEscape(databaseID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a database: DELETE /databases/{database_id}.
// Scope: write.
func (s *DatabasesService) Delete(ctx context.Context, databaseID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "DELETE", "/databases/"+url.PathEscape(databaseID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ReplaceAccessRules replaces access rules: PUT /databases/{database_id}/access-rules.
// Scope: write.
func (s *DatabasesService) ReplaceAccessRules(ctx context.Context, databaseID string, req ReplaceDatabaseAccessRulesRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "PUT", "/databases/"+url.PathEscape(databaseID)+"/access-rules", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetSecurityGroups attaches security groups: PUT /databases/{database_id}/security-groups.
// Scope: write.
func (s *DatabasesService) SetSecurityGroups(ctx context.Context, databaseID string, req SetDatabaseSecurityGroupsRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "PUT", "/databases/"+url.PathEscape(databaseID)+"/security-groups", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ChangePassword sets the password: PUT /databases/{database_id}/password.
// Scope: write.
func (s *DatabasesService) ChangePassword(ctx context.Context, databaseID string, req ChangeDatabasePasswordRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "PUT", "/databases/"+url.PathEscape(databaseID)+"/password", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ResetPassword resets the password: POST /databases/{database_id}/reset-password.
// Scope: write.
func (s *DatabasesService) ResetPassword(ctx context.Context, databaseID string, opts ...RequestOption) (*DatabasePasswordAccepted, error) {
	var out DatabasePasswordAccepted
	if err := s.c.Do(ctx, "POST", "/databases/"+url.PathEscape(databaseID)+"/reset-password", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Start starts a database: POST /databases/{database_id}/start.
// Scope: write.
func (s *DatabasesService) Start(ctx context.Context, databaseID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/databases/"+url.PathEscape(databaseID)+"/start", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Stop stops a database: POST /databases/{database_id}/stop.
// Scope: write.
func (s *DatabasesService) Stop(ctx context.Context, databaseID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/databases/"+url.PathEscape(databaseID)+"/stop", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ResizeStorage grows a database's storage: POST /databases/{database_id}/resize-storage.
// Scope: write.
func (s *DatabasesService) ResizeStorage(ctx context.Context, databaseID string, req ResizeDatabaseStorageRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/databases/"+url.PathEscape(databaseID)+"/resize-storage", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListSnapshots lists a database's snapshots: GET /databases/{database_id}/snapshots.
// Scope: read.
func (s *DatabasesService) ListSnapshots(ctx context.Context, databaseID string, params *ListParams, opts ...RequestOption) (*Page[Snapshot], error) {
	var out Page[Snapshot]
	if err := s.c.Do(ctx, "GET", "/databases/"+url.PathEscape(databaseID)+"/snapshots", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// AllSnapshots iterates over every item of ListSnapshots, fetching pages as needed.
// It stops at the first error, which it yields.
func (s *DatabasesService) AllSnapshots(ctx context.Context, databaseID string, params *ListParams, opts ...RequestOption) iter.Seq2[Snapshot, error] {
	var p ListParams
	if params != nil {
		p = *params
	}
	return paginate(p.Cursor, func(cursor string) (*Page[Snapshot], error) {
		p.Cursor = cursor
		return s.ListSnapshots(ctx, databaseID, &p, opts...)
	})
}

// CreateSnapshot snapshots a database: POST /databases/{database_id}/snapshots.
// Scope: write.
func (s *DatabasesService) CreateSnapshot(ctx context.Context, databaseID string, req CreateSnapshotRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/databases/"+url.PathEscape(databaseID)+"/snapshots", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetSnapshot gets a database's snapshot: GET /databases/{database_id}/snapshots/{snapshot_id}.
// Scope: read.
func (s *DatabasesService) GetSnapshot(ctx context.Context, databaseID string, snapshotID string, opts ...RequestOption) (*Snapshot, error) {
	var out Snapshot
	if err := s.c.Do(ctx, "GET", "/databases/"+url.PathEscape(databaseID)+"/snapshots/"+url.PathEscape(snapshotID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteSnapshot deletes a database's snapshot: DELETE /databases/{database_id}/snapshots/{snapshot_id}.
// Scope: write.
func (s *DatabasesService) DeleteSnapshot(ctx context.Context, databaseID string, snapshotID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "DELETE", "/databases/"+url.PathEscape(databaseID)+"/snapshots/"+url.PathEscape(snapshotID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListOrders lists database orders: GET /database-orders.
// Scope: read.
func (s *DatabasesService) ListOrders(ctx context.Context, params *PageParams, opts ...RequestOption) (*Page[DatabaseOrder], error) {
	var out Page[DatabaseOrder]
	if err := s.c.Do(ctx, "GET", "/database-orders", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// AllOrders iterates over every item of ListOrders, fetching pages as needed.
// It stops at the first error, which it yields.
func (s *DatabasesService) AllOrders(ctx context.Context, params *PageParams, opts ...RequestOption) iter.Seq2[DatabaseOrder, error] {
	var p PageParams
	if params != nil {
		p = *params
	}
	return paginate(p.Cursor, func(cursor string) (*Page[DatabaseOrder], error) {
		p.Cursor = cursor
		return s.ListOrders(ctx, &p, opts...)
	})
}

// GetOrder gets a database order: GET /database-orders/{order_id}.
// Scope: read.
func (s *DatabasesService) GetOrder(ctx context.Context, orderID string, opts ...RequestOption) (*DatabaseOrder, error) {
	var out DatabaseOrder
	if err := s.c.Do(ctx, "GET", "/database-orders/"+url.PathEscape(orderID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
